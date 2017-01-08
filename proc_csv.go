package psn

import (
	"encoding/csv"
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Proc represents an entry of various system statistics.
type Proc struct {
	UnixTS  int64
	PSEntry PSEntry

	DSEntry             DSEntry
	ReadsCompletedDiff  uint64
	SectorsReadDiff     uint64
	WritesCompletedDiff uint64
	SectorsWrittenDiff  uint64

	NSEntry              NSEntry
	ReceiveBytesDiff     string
	ReceivePacketsDiff   uint64
	TransmitBytesDiff    string
	TransmitPacketsDiff  uint64
	ReceiveBytesNumDiff  uint64
	TransmitBytesNumDiff uint64
}

// GetProc returns current 'Proc' data.
func GetProc(pid int64, diskDevice string, networkInterface string) (Proc, error) {
	proc := Proc{UnixTS: time.Now().Unix()}

	errc := make(chan error)
	go func() {
		// get process stats
		ets, err := GetPS(WithPID(pid))
		if err != nil {
			errc <- err
			return
		}
		if len(ets) != 1 {
			errc <- fmt.Errorf("len(PID=%d entries) != 1 (got %d)", pid, len(ets))
			return
		}
		proc.PSEntry = ets[0]
		errc <- nil
	}()
	go func() {
		// get diskstats
		ds, err := GetDS()
		if err != nil {
			errc <- err
			return
		}
		for _, elem := range ds {
			if elem.Device == diskDevice {
				proc.DSEntry = elem
				break
			}
		}
		errc <- nil
	}()
	go func() {
		// get network I/O stats
		ns, err := GetNS()
		if err != nil {
			errc <- err
			return
		}
		for _, elem := range ns {
			if elem.Interface == networkInterface {
				proc.NSEntry = elem
				break
			}
		}
		errc <- nil
	}()

	cnt := 0
	for cnt != 3 {
		err := <-errc
		if err != nil {
			return Proc{}, err
		}
		cnt++
	}

	if proc.DSEntry.Device == "" {
		return Proc{}, fmt.Errorf("disk device %q was not found", diskDevice)
	}
	if proc.NSEntry.Interface == "" {
		return Proc{}, fmt.Errorf("network interface %q was not found", networkInterface)
	}
	return proc, nil
}

var (
	// ProcCSVHeader lists all Proc CSV columns.
	ProcCSVHeader = append([]string{"UNIX-TS"}, columnsPSEntry...)

	// ProcCSVHeaderIndex maps each column name to its index in row.
	ProcCSVHeaderIndex = make(map[string]int)
)

func init() {
	// more columns to 'ProcCSVHeader'
	ProcCSVHeader = append(ProcCSVHeader, columnsDSEntry...)
	ProcCSVHeader = append(ProcCSVHeader, columnsNSEntry...)
	ProcCSVHeader = append(ProcCSVHeader,
		"READS-COMPLETED-DIFF",
		"SECTORS-READ-DIFF",
		"WRITES-COMPLETED-DIFF",
		"SECTORS-WRITTEN-DIFF",

		"RECEIVE-BYTES-DIFF",
		"RECEIVE-PACKETS-DIFF",
		"TRANSMIT-BYTES-DIFF",
		"TRANSMIT-PACKETS-DIFF",
		"RECEIVE-BYTES-NUM-DIFF",
		"TRANSMIT-BYTES-NUM-DIFF",
	)

	for i, v := range ProcCSVHeader {
		ProcCSVHeaderIndex[v] = i
	}
}

// ToRow converts 'Proc' to string slice.
func (p *Proc) ToRow() (row []string) {
	row = make([]string, len(ProcCSVHeader))
	row[0] = fmt.Sprintf("%d", p.UnixTS) // UNIX-TS

	row[1] = p.PSEntry.Program                       // PROGRAM
	row[2] = p.PSEntry.State                         // STATE
	row[3] = fmt.Sprintf("%d", p.PSEntry.PID)        // PID
	row[4] = fmt.Sprintf("%d", p.PSEntry.PPID)       // PPID
	row[5] = p.PSEntry.CPU                           // CPU
	row[6] = p.PSEntry.VMRSS                         // VMRSS
	row[7] = p.PSEntry.VMSize                        // VMSIZE
	row[8] = fmt.Sprintf("%d", p.PSEntry.FD)         // FD
	row[9] = fmt.Sprintf("%d", p.PSEntry.Threads)    // THREADS
	row[10] = fmt.Sprintf("%3.2f", p.PSEntry.CPUNum) // CPU-NUM
	row[11] = fmt.Sprintf("%d", p.PSEntry.VMRSSNum)  // VMRSS-NUM
	row[12] = fmt.Sprintf("%d", p.PSEntry.VMSizeNum) // VMSIZE-NUM

	row[13] = p.DSEntry.Device                                  // DEVICE
	row[14] = fmt.Sprintf("%d", p.DSEntry.ReadsCompleted)       // READS-COMPLETED
	row[15] = fmt.Sprintf("%d", p.DSEntry.SectorsRead)          // SECTORS-READ
	row[16] = p.DSEntry.TimeSpentOnReading                      // TIME(READS)
	row[17] = fmt.Sprintf("%d", p.DSEntry.WritesCompleted)      // WRITES-COMPLETED
	row[18] = fmt.Sprintf("%d", p.DSEntry.SectorsWritten)       // SECTORS-WRITTEN
	row[19] = p.DSEntry.TimeSpentOnWriting                      // TIME(WRITES)
	row[20] = fmt.Sprintf("%d", p.DSEntry.TimeSpentOnReadingMs) // MILLISECONDS(READS)
	row[21] = fmt.Sprintf("%d", p.DSEntry.TimeSpentOnWritingMs) // MILLISECONDS(WRITES)

	row[22] = p.NSEntry.Interface                           // INTERFACE
	row[23] = p.NSEntry.ReceiveBytes                        // RECEIVE-BYTES
	row[24] = fmt.Sprintf("%d", p.NSEntry.ReceivePackets)   // RECEIVE-PACKETS
	row[25] = p.NSEntry.TransmitBytes                       // TRANSMIT-BYTES
	row[26] = fmt.Sprintf("%d", p.NSEntry.TransmitPackets)  // TRANSMIT-PACKETS
	row[27] = fmt.Sprintf("%d", p.NSEntry.ReceiveBytesNum)  // RECEIVE-BYTES-NUM
	row[28] = fmt.Sprintf("%d", p.NSEntry.TransmitBytesNum) // TRANSMIT-BYTES-NUM

	row[29] = fmt.Sprintf("%d", p.ReadsCompletedDiff)  // READS-COMPLETED-DIFF
	row[30] = fmt.Sprintf("%d", p.SectorsReadDiff)     // SECTORS-READ-DIFF
	row[31] = fmt.Sprintf("%d", p.WritesCompletedDiff) // WRITES-COMPLETED-DIFF
	row[32] = fmt.Sprintf("%d", p.SectorsWrittenDiff)  // SECTORS-WRITTEN-DIFF

	row[33] = p.ReceiveBytesDiff                        // RECEIVE-BYTES-DIFF
	row[34] = fmt.Sprintf("%d", p.ReceivePacketsDiff)   // RECEIVE-PACKETS-DIFF
	row[35] = p.TransmitBytesDiff                       // TRANSMIT-BYTES-DIFF
	row[36] = fmt.Sprintf("%d", p.TransmitPacketsDiff)  // TRANSMIT-PACKETS-DIFF
	row[37] = fmt.Sprintf("%d", p.ReceiveBytesNumDiff)  // RECEIVE-BYTES-NUM-DIFF
	row[38] = fmt.Sprintf("%d", p.TransmitBytesNumDiff) // TRANSMIT-BYTES-NUM-DIFF

	return
}

// CSV represents CSV data (header, rows, etc.).
type CSV struct {
	FilePath         string
	PID              int64
	DiskDevice       string
	NetworkInterface string

	Header      []string
	HeaderIndex map[string]int
	MinUnixTS   int64
	MaxUnixTS   int64

	// Rows are sorted by unix timestamps.
	Rows []Proc
}

// NewCSV returns a new CSV.
func NewCSV(fpath string, pid int64, diskDevice string, networkInterface string) *CSV {
	return &CSV{
		FilePath:         fpath,
		PID:              pid,
		DiskDevice:       diskDevice,
		NetworkInterface: networkInterface,

		Header:      ProcCSVHeader,
		HeaderIndex: ProcCSVHeaderIndex,

		MinUnixTS: 0,
		MaxUnixTS: 0,

		Rows: []Proc{},
	}
}

// Add is to be called periodically to add a row to CSV.
// It only appends to CSV. And it estimates empty rows by unix timestamps(seconds).
func (c *CSV) Add() error {
	cur, err := GetProc(c.PID, c.DiskDevice, c.NetworkInterface)
	if err != nil {
		return err
	}
	if len(c.Rows) == 0 {
		c.MinUnixTS = cur.UnixTS
		c.MaxUnixTS = cur.UnixTS
		c.Rows = []Proc{cur}
		return nil
	}
	prev := c.Rows[len(c.Rows)-1]
	if prev.UnixTS >= cur.UnixTS {
		// ignore data with wrong timestamps
		return nil
	}
	c.MaxUnixTS = cur.UnixTS

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

	var nexts []Proc

	// see if there are empty rows between
	if (cur.UnixTS - prev.UnixTS) > 1 {
		tsDiff := cur.UnixTS - prev.UnixTS
		nexts = make([]Proc, 0, tsDiff+1)

		// estimate the previous ones based on 'prev' and 'cur'
		prev2 := prev

		// PSEntry; just use average since some metrisc might decrease
		prev2.PSEntry.FD = prev.PSEntry.FD + (cur.PSEntry.FD-prev.PSEntry.FD)/2
		prev2.PSEntry.Threads = prev.PSEntry.Threads + (cur.PSEntry.Threads-prev.PSEntry.Threads)/2
		prev2.PSEntry.CPUNum = prev.PSEntry.CPUNum + (cur.PSEntry.CPUNum-prev.PSEntry.CPUNum)/2
		prev2.PSEntry.VMRSSNum = prev.PSEntry.VMRSSNum + (cur.PSEntry.VMRSSNum-prev.PSEntry.VMRSSNum)/2
		prev2.PSEntry.VMSizeNum = prev.PSEntry.VMSizeNum + (cur.PSEntry.VMSizeNum-prev.PSEntry.VMSizeNum)/2
		prev2.PSEntry.CPU = fmt.Sprintf("%3.2f %%", prev2.PSEntry.CPUNum)
		prev2.PSEntry.VMRSS = humanize.Bytes(prev2.PSEntry.VMRSSNum)
		prev2.PSEntry.VMSize = humanize.Bytes(prev2.PSEntry.VMSizeNum)

		// DSEntry; calculate delta assuming that metrics are cumulative
		prev2.ReadsCompletedDiff = (cur.DSEntry.ReadsCompleted - prev.DSEntry.ReadsCompleted) / uint64(tsDiff)
		prev2.SectorsReadDiff = (cur.DSEntry.SectorsRead - prev.DSEntry.SectorsRead) / uint64(tsDiff)
		prev2.WritesCompletedDiff = (cur.DSEntry.WritesCompleted - prev.DSEntry.WritesCompleted) / uint64(tsDiff)
		prev2.SectorsWrittenDiff = (cur.DSEntry.SectorsWritten - prev.DSEntry.SectorsWritten) / uint64(tsDiff)
		timeSpentOnReadingMsDelta := (cur.DSEntry.TimeSpentOnReadingMs - prev.DSEntry.TimeSpentOnReadingMs) / uint64(tsDiff)
		timeSpentOnWritingMsDelta := (cur.DSEntry.TimeSpentOnWritingMs - prev.DSEntry.TimeSpentOnWritingMs) / uint64(tsDiff)

		// NSEntry; calculate delta assuming that metrics are cumulative
		prev2.ReceiveBytesNumDiff = (cur.NSEntry.ReceiveBytesNum - prev.NSEntry.ReceiveBytesNum) / uint64(tsDiff)
		prev2.ReceiveBytesDiff = humanize.Bytes(prev2.ReceiveBytesNumDiff)
		prev2.ReceivePacketsDiff = (cur.NSEntry.ReceivePackets - prev.NSEntry.ReceivePackets) / uint64(tsDiff)
		prev2.TransmitBytesNumDiff = (cur.NSEntry.TransmitBytesNum - prev.NSEntry.TransmitBytesNum) / uint64(tsDiff)
		prev2.TransmitBytesDiff = humanize.Bytes(prev2.TransmitBytesNumDiff)
		prev2.TransmitPacketsDiff = (cur.NSEntry.TransmitPackets - prev.NSEntry.TransmitPackets) / uint64(tsDiff)

		for i := int64(1); i < tsDiff; i++ {
			ev := prev2
			ev.UnixTS = prev.UnixTS + i

			ev.DSEntry.ReadsCompleted += prev2.ReadsCompletedDiff * uint64(i)
			ev.DSEntry.SectorsRead += prev2.SectorsReadDiff * uint64(i)
			ev.DSEntry.WritesCompleted += prev2.WritesCompletedDiff * uint64(i)
			ev.DSEntry.SectorsWritten += prev2.SectorsWrittenDiff * uint64(i)
			ev.DSEntry.TimeSpentOnReadingMs += timeSpentOnReadingMsDelta * uint64(i)
			ev.DSEntry.TimeSpentOnWritingMs += timeSpentOnWritingMsDelta * uint64(i)
			ev.DSEntry.TimeSpentOnReading = humanizeDurationMs(ev.DSEntry.TimeSpentOnReadingMs)
			ev.DSEntry.TimeSpentOnWriting = humanizeDurationMs(ev.DSEntry.TimeSpentOnWritingMs)

			ev.NSEntry.ReceiveBytesNum += prev2.ReceiveBytesNumDiff * uint64(i)
			ev.NSEntry.ReceiveBytes = humanize.Bytes(ev.NSEntry.ReceiveBytesNum)
			ev.NSEntry.ReceivePackets += prev2.ReceivePacketsDiff * uint64(i)
			ev.NSEntry.TransmitBytesNum += prev2.TransmitBytesNumDiff * uint64(i)
			ev.NSEntry.TransmitBytes = humanize.Bytes(ev.NSEntry.TransmitBytesNum)
			ev.NSEntry.TransmitPackets += prev2.TransmitPacketsDiff * uint64(i)

			nexts = append(nexts, ev)
		}
		nexts = append(nexts, cur)
	} else {
		nexts = []Proc{cur}
	}

	c.Rows = append(c.Rows, nexts...)
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
	// TODO
	return nil, nil
}
