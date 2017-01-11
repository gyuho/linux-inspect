package psn

import (
	"encoding/csv"
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
)

// CSV represents CSV data (header, rows, etc.).
type CSV struct {
	FilePath         string
	PID              int64
	DiskDevice       string
	NetworkInterface string

	Header      []string
	HeaderIndex map[string]int

	MinUnixTS int64
	MaxUnixTS int64

	// ExtraPath contains extra information.
	ExtraPath string

	// Rows are sorted by unix seconds.
	Rows []Proc
}

// NewCSV returns a new CSV.
func NewCSV(fpath string, pid int64, diskDevice string, networkInterface string, extraPath string) *CSV {
	return &CSV{
		FilePath:         fpath,
		PID:              pid,
		DiskDevice:       diskDevice,
		NetworkInterface: networkInterface,

		Header:      ProcHeader,
		HeaderIndex: ProcHeaderIndex,

		MinUnixTS: 0,
		MaxUnixTS: 0,

		ExtraPath: extraPath,
		Rows:      []Proc{},
	}
}

// Add is to be called periodically to add a row to CSV.
// It only appends to CSV. And it estimates empty rows by unix seconds.
func (c *CSV) Add() error {
	cur, err := GetProc(
		WithPID(c.PID),
		WithDiskDevice(c.DiskDevice),
		WithNetworkInterface(c.NetworkInterface),
		WithExtraPath(c.ExtraPath),
	)
	if err != nil {
		return err
	}

	// first call; just append and return
	if len(c.Rows) == 0 {
		c.MinUnixTS = cur.UnixTS
		c.MaxUnixTS = cur.UnixTS
		c.Rows = []Proc{cur}
		return nil
	}

	// compare with previous row before append
	prev := c.Rows[len(c.Rows)-1]
	if prev.UnixTS >= cur.UnixTS {
		// ignore data with wrong seconds
		return nil
	}

	// 'Add' only appends, so later unix should be max
	c.MaxUnixTS = cur.UnixTS

	if cur.UnixTS-prev.UnixTS == 1 {
		cur.ReadsCompletedDiff = cur.DSEntry.ReadsCompleted - prev.DSEntry.ReadsCompleted
		cur.SectorsReadDiff = cur.DSEntry.SectorsRead - prev.DSEntry.SectorsRead
		cur.WritesCompletedDiff = cur.DSEntry.WritesCompleted - prev.DSEntry.WritesCompleted
		cur.SectorsWrittenDiff = cur.DSEntry.SectorsWritten - prev.DSEntry.SectorsWritten

		cur.ReceiveBytesNumDiff = cur.NSEntry.ReceiveBytesNum - prev.NSEntry.ReceiveBytesNum
		cur.TransmitBytesNumDiff = cur.NSEntry.TransmitBytesNum - prev.NSEntry.TransmitBytesNum
		cur.ReceivePacketsDiff = cur.NSEntry.ReceivePackets - prev.NSEntry.ReceivePackets
		cur.TransmitPacketsDiff = cur.NSEntry.TransmitPackets - prev.NSEntry.TransmitPackets

		cur.ReceiveBytesDiff = humanize.Bytes(cur.ReceiveBytesNumDiff)
		cur.TransmitBytesDiff = humanize.Bytes(cur.TransmitBytesNumDiff)

		c.Rows = append(c.Rows, cur)
		return nil
	}

	// there are empty rows between; estimate and fill-in
	tsDiff := cur.UnixTS - prev.UnixTS
	nexts := make([]Proc, 0, tsDiff+1)

	// estimate the previous ones based on 'prev' and 'cur'
	mid := prev

	// Extra; just use the previous value
	mid.Extra = prev.Extra

	// PSEntry; just use average since some metrisc might decrease
	mid.PSEntry.FD = prev.PSEntry.FD + (cur.PSEntry.FD-prev.PSEntry.FD)/2
	mid.PSEntry.Threads = prev.PSEntry.Threads + (cur.PSEntry.Threads-prev.PSEntry.Threads)/2
	mid.PSEntry.CPUNum = prev.PSEntry.CPUNum + (cur.PSEntry.CPUNum-prev.PSEntry.CPUNum)/2
	mid.PSEntry.VMRSSNum = prev.PSEntry.VMRSSNum + (cur.PSEntry.VMRSSNum-prev.PSEntry.VMRSSNum)/2
	mid.PSEntry.VMSizeNum = prev.PSEntry.VMSizeNum + (cur.PSEntry.VMSizeNum-prev.PSEntry.VMSizeNum)/2
	mid.PSEntry.CPU = fmt.Sprintf("%3.2f %%", mid.PSEntry.CPUNum)
	mid.PSEntry.VMRSS = humanize.Bytes(mid.PSEntry.VMRSSNum)
	mid.PSEntry.VMSize = humanize.Bytes(mid.PSEntry.VMSizeNum)

	// DSEntry; calculate delta assuming that metrics are cumulative
	mid.ReadsCompletedDiff = (cur.DSEntry.ReadsCompleted - prev.DSEntry.ReadsCompleted) / uint64(tsDiff)
	mid.SectorsReadDiff = (cur.DSEntry.SectorsRead - prev.DSEntry.SectorsRead) / uint64(tsDiff)
	mid.WritesCompletedDiff = (cur.DSEntry.WritesCompleted - prev.DSEntry.WritesCompleted) / uint64(tsDiff)
	mid.SectorsWrittenDiff = (cur.DSEntry.SectorsWritten - prev.DSEntry.SectorsWritten) / uint64(tsDiff)
	timeSpentOnReadingMsDelta := (cur.DSEntry.TimeSpentOnReadingMs - prev.DSEntry.TimeSpentOnReadingMs) / uint64(tsDiff)
	timeSpentOnWritingMsDelta := (cur.DSEntry.TimeSpentOnWritingMs - prev.DSEntry.TimeSpentOnWritingMs) / uint64(tsDiff)

	// NSEntry; calculate delta assuming that metrics are cumulative
	mid.ReceiveBytesNumDiff = (cur.NSEntry.ReceiveBytesNum - prev.NSEntry.ReceiveBytesNum) / uint64(tsDiff)
	mid.ReceiveBytesDiff = humanize.Bytes(mid.ReceiveBytesNumDiff)
	mid.ReceivePacketsDiff = (cur.NSEntry.ReceivePackets - prev.NSEntry.ReceivePackets) / uint64(tsDiff)
	mid.TransmitBytesNumDiff = (cur.NSEntry.TransmitBytesNum - prev.NSEntry.TransmitBytesNum) / uint64(tsDiff)
	mid.TransmitBytesDiff = humanize.Bytes(mid.TransmitBytesNumDiff)
	mid.TransmitPacketsDiff = (cur.NSEntry.TransmitPackets - prev.NSEntry.TransmitPackets) / uint64(tsDiff)

	for i := int64(1); i < tsDiff; i++ {
		ev := mid
		ev.UnixTS = prev.UnixTS + i

		ev.DSEntry.ReadsCompleted += mid.ReadsCompletedDiff * uint64(i)
		ev.DSEntry.SectorsRead += mid.SectorsReadDiff * uint64(i)
		ev.DSEntry.WritesCompleted += mid.WritesCompletedDiff * uint64(i)
		ev.DSEntry.SectorsWritten += mid.SectorsWrittenDiff * uint64(i)
		ev.DSEntry.TimeSpentOnReadingMs += timeSpentOnReadingMsDelta * uint64(i)
		ev.DSEntry.TimeSpentOnWritingMs += timeSpentOnWritingMsDelta * uint64(i)
		ev.DSEntry.TimeSpentOnReading = humanizeDurationMs(ev.DSEntry.TimeSpentOnReadingMs)
		ev.DSEntry.TimeSpentOnWriting = humanizeDurationMs(ev.DSEntry.TimeSpentOnWritingMs)

		ev.NSEntry.ReceiveBytesNum += mid.ReceiveBytesNumDiff * uint64(i)
		ev.NSEntry.ReceiveBytes = humanize.Bytes(ev.NSEntry.ReceiveBytesNum)
		ev.NSEntry.ReceivePackets += mid.ReceivePacketsDiff * uint64(i)
		ev.NSEntry.TransmitBytesNum += mid.TransmitBytesNumDiff * uint64(i)
		ev.NSEntry.TransmitBytes = humanize.Bytes(ev.NSEntry.TransmitBytesNum)
		ev.NSEntry.TransmitPackets += mid.TransmitPacketsDiff * uint64(i)

		nexts = append(nexts, ev)
	}

	// now previous entry is estimated; update 'cur' diff metrics
	realPrev := nexts[len(nexts)-1]

	cur.ReadsCompletedDiff = cur.DSEntry.ReadsCompleted - realPrev.DSEntry.ReadsCompleted
	cur.SectorsReadDiff = cur.DSEntry.SectorsRead - realPrev.DSEntry.SectorsRead
	cur.WritesCompletedDiff = cur.DSEntry.WritesCompleted - realPrev.DSEntry.WritesCompleted
	cur.SectorsWrittenDiff = cur.DSEntry.SectorsWritten - realPrev.DSEntry.SectorsWritten

	cur.ReceiveBytesNumDiff = cur.NSEntry.ReceiveBytesNum - realPrev.NSEntry.ReceiveBytesNum
	cur.TransmitBytesNumDiff = cur.NSEntry.TransmitBytesNum - realPrev.NSEntry.TransmitBytesNum
	cur.ReceivePacketsDiff = cur.NSEntry.ReceivePackets - realPrev.NSEntry.ReceivePackets
	cur.TransmitPacketsDiff = cur.NSEntry.TransmitPackets - realPrev.NSEntry.TransmitPackets

	cur.ReceiveBytesDiff = humanize.Bytes(cur.ReceiveBytesNumDiff)
	cur.TransmitBytesDiff = humanize.Bytes(cur.TransmitBytesNumDiff)

	c.Rows = append(c.Rows, append(nexts, cur)...)
	return nil
}

// Save saves CSV to disk.
func (c *CSV) Save() error {
	f, err := openToAppend(c.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	if err := wr.Write(c.Header); err != nil {
		return err
	}

	rows := make([][]string, len(c.Rows))
	for i, row := range c.Rows {
		rows[i] = row.ToRow()
	}
	if err := wr.WriteAll(rows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}

// ReadCSV reads a CSV file and convert to 'CSV'.
func ReadCSV(fpath string) (*CSV, error) {
	f, err := openToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := csv.NewReader(f)

	// in case that rows have different number of fields
	rd.FieldsPerRecord = -1

	rows, err := rd.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) <= 1 {
		return nil, fmt.Errorf("expected len(rows)>1, got %d", len(rows))
	}
	if rows[0][0] != "UNIX-TS" {
		return nil, fmt.Errorf("expected header at top, got %+v", rows[0])
	}

	// remove header
	rows = rows[1:len(rows):len(rows)]
	min, err := strconv.ParseInt(rows[0][0], 10, 64)
	if err != nil {
		return nil, err
	}
	max, err := strconv.ParseInt(rows[len(rows)-1][0], 10, 64)
	if err != nil {
		return nil, err
	}
	c := &CSV{
		FilePath:         fpath,
		PID:              0,
		DiskDevice:       "",
		NetworkInterface: "",

		Header:      ProcHeader,
		HeaderIndex: ProcHeaderIndex,
		MinUnixTS:   min,
		MaxUnixTS:   max,

		Rows: make([]Proc, 0, len(rows)),
	}
	for _, row := range rows {
		ts, err := strconv.ParseInt(row[ProcHeaderIndex["UNIX-TS"]], 10, 64)
		if err != nil {
			return nil, err
		}
		pid, err := strconv.ParseInt(row[ProcHeaderIndex["PID"]], 10, 64)
		if err != nil {
			return nil, err
		}
		ppid, err := strconv.ParseInt(row[ProcHeaderIndex["PPID"]], 10, 64)
		if err != nil {
			return nil, err
		}
		fd, err := strconv.ParseUint(row[ProcHeaderIndex["FD"]], 10, 64)
		if err != nil {
			return nil, err
		}
		threads, err := strconv.ParseUint(row[ProcHeaderIndex["THREADS"]], 10, 64)
		if err != nil {
			return nil, err
		}
		cpuNum, err := strconv.ParseFloat(row[ProcHeaderIndex["CPU-NUM"]], 64)
		if err != nil {
			return nil, err
		}
		vmRssNum, err := strconv.ParseUint(row[ProcHeaderIndex["VMRSS-NUM"]], 10, 64)
		if err != nil {
			return nil, err
		}
		vmSizeNum, err := strconv.ParseUint(row[ProcHeaderIndex["VMSIZE-NUM"]], 10, 64)
		if err != nil {
			return nil, err
		}

		readsCompleted, err := strconv.ParseUint(row[ProcHeaderIndex["READS-COMPLETED"]], 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsRead, err := strconv.ParseUint(row[ProcHeaderIndex["SECTORS-READ"]], 10, 64)
		if err != nil {
			return nil, err
		}
		writesCompleted, err := strconv.ParseUint(row[ProcHeaderIndex["WRITES-COMPLETED"]], 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsWritten, err := strconv.ParseUint(row[ProcHeaderIndex["SECTORS-WRITTEN"]], 10, 64)
		if err != nil {
			return nil, err
		}
		timeSpentOnReadingMs, err := strconv.ParseUint(row[ProcHeaderIndex["MILLISECONDS(READS)"]], 10, 64)
		if err != nil {
			return nil, err
		}
		timeSpentOnWritingMs, err := strconv.ParseUint(row[ProcHeaderIndex["MILLISECONDS(WRITES)"]], 10, 64)
		if err != nil {
			return nil, err
		}

		readsCompletedDiff, err := strconv.ParseUint(row[ProcHeaderIndex["READS-COMPLETED-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsReadDiff, err := strconv.ParseUint(row[ProcHeaderIndex["SECTORS-READ-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		writesCompletedDiff, err := strconv.ParseUint(row[ProcHeaderIndex["WRITES-COMPLETED-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsWrittenDiff, err := strconv.ParseUint(row[ProcHeaderIndex["SECTORS-WRITTEN-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}

		receivePackets, err := strconv.ParseUint(row[ProcHeaderIndex["RECEIVE-PACKETS"]], 10, 64)
		if err != nil {
			return nil, err
		}
		transmitPackets, err := strconv.ParseUint(row[ProcHeaderIndex["TRANSMIT-PACKETS"]], 10, 64)
		if err != nil {
			return nil, err
		}
		receiveBytesNum, err := strconv.ParseUint(row[ProcHeaderIndex["RECEIVE-BYTES-NUM"]], 10, 64)
		if err != nil {
			return nil, err
		}
		transmitBytesNum, err := strconv.ParseUint(row[ProcHeaderIndex["TRANSMIT-BYTES-NUM"]], 10, 64)
		if err != nil {
			return nil, err
		}

		receivePacketsDiff, err := strconv.ParseUint(row[ProcHeaderIndex["RECEIVE-PACKETS-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		transmitPacketsDiff, err := strconv.ParseUint(row[ProcHeaderIndex["TRANSMIT-PACKETS-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		receiveBytesNumDiff, err := strconv.ParseUint(row[ProcHeaderIndex["RECEIVE-BYTES-NUM-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}
		transmitBytesNumDiff, err := strconv.ParseUint(row[ProcHeaderIndex["TRANSMIT-BYTES-NUM-DIFF"]], 10, 64)
		if err != nil {
			return nil, err
		}

		proc := Proc{
			UnixTS: ts,
			PSEntry: PSEntry{
				Program:   row[ProcHeaderIndex["PROGRAM"]],
				State:     row[ProcHeaderIndex["STATE"]],
				PID:       pid,
				PPID:      ppid,
				CPU:       row[ProcHeaderIndex["CPU"]],
				VMRSS:     row[ProcHeaderIndex["VMRSS"]],
				VMSize:    row[ProcHeaderIndex["VMSIZE"]],
				FD:        fd,
				Threads:   threads,
				CPUNum:    cpuNum,
				VMRSSNum:  vmRssNum,
				VMSizeNum: vmSizeNum,
			},

			DSEntry: DSEntry{
				Device:               row[ProcHeaderIndex["DEVICE"]],
				ReadsCompleted:       readsCompleted,
				SectorsRead:          sectorsRead,
				TimeSpentOnReading:   row[ProcHeaderIndex["TIME(READS)"]],
				WritesCompleted:      writesCompleted,
				SectorsWritten:       sectorsWritten,
				TimeSpentOnWriting:   row[ProcHeaderIndex["TIME(WRITES)"]],
				TimeSpentOnReadingMs: timeSpentOnReadingMs,
				TimeSpentOnWritingMs: timeSpentOnWritingMs,
			},
			ReadsCompletedDiff:  readsCompletedDiff,
			SectorsReadDiff:     sectorsReadDiff,
			WritesCompletedDiff: writesCompletedDiff,
			SectorsWrittenDiff:  sectorsWrittenDiff,

			NSEntry: NSEntry{
				Interface:        row[ProcHeaderIndex["INTERFACE"]],
				ReceiveBytes:     row[ProcHeaderIndex["RECEIVE-BYTES"]],
				ReceivePackets:   receivePackets,
				TransmitBytes:    row[ProcHeaderIndex["TRANSMIT-BYTES"]],
				TransmitPackets:  transmitPackets,
				ReceiveBytesNum:  receiveBytesNum,
				TransmitBytesNum: transmitBytesNum,
			},
			ReceiveBytesDiff:     row[ProcHeaderIndex["RECEIVE-BYTES-DIFF"]],
			ReceivePacketsDiff:   receivePacketsDiff,
			TransmitBytesDiff:    row[ProcHeaderIndex["TRANSMIT-BYTES-DIFF"]],
			TransmitPacketsDiff:  transmitPacketsDiff,
			ReceiveBytesNumDiff:  receiveBytesNumDiff,
			TransmitBytesNumDiff: transmitBytesNumDiff,

			Extra: []byte(row[ProcHeaderIndex["EXTRA"]]),
		}
		c.PID = proc.PSEntry.PID
		c.DiskDevice = proc.DSEntry.Device
		c.NetworkInterface = proc.NSEntry.Interface

		c.Rows = append(c.Rows, proc)
	}

	return c, nil
}
