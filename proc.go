package psn

import (
	"fmt"
	"io/ioutil"
	"time"
)

// Proc represents an entry of various system statistics.
type Proc struct {
	UnixTS int64

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

	// Extra exists to support customized data query.
	Extra []byte
}

// GetProc returns current 'Proc' data.
func GetProc(pid int64, diskDevice string, networkInterface string, extraPath string) (Proc, error) {
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
	go func() {
		f, err := openToRead(extraPath)
		if err != nil {
			errc <- err
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			errc <- err
			return
		}
		proc.Extra = b
		errc <- nil
	}()

	cnt := 0
	for cnt != 4 {
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
	// ProcHeader lists all Proc CSV columns.
	ProcHeader = append([]string{"UNIX-TS"}, columnsPSEntry...)

	// ProcHeaderIndex maps each Proc column name to its index in row.
	ProcHeaderIndex = make(map[string]int)
)

func init() {
	// more columns to 'ProcHeader'
	ProcHeader = append(ProcHeader, columnsDSEntry...)
	ProcHeader = append(ProcHeader, columnsNSEntry...)
	ProcHeader = append(ProcHeader,
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

		"EXTRA",
	)

	for i, v := range ProcHeader {
		ProcHeaderIndex[v] = i
	}
}

// ToRow converts 'Proc' to string slice.
func (p *Proc) ToRow() (row []string) {
	row = make([]string, len(ProcHeader))
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

	row[39] = string(p.Extra) // EXTRA

	return
}
