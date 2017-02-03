package psn

import (
	"fmt"
	"io/ioutil"
	"time"
)

// Proc represents an entry of various system statistics.
type Proc struct {
	// UnixNanosecond is unix nano second when this Proc row gets created.
	UnixNanosecond int64

	// UnixSecond is the converted Unix seconds from UnixNano.
	UnixSecond int64

	PSEntry PSEntry

	LoadAvg LoadAvg

	DSEntry              DSEntry
	ReadsCompletedDelta  uint64
	SectorsReadDelta     uint64
	WritesCompletedDelta uint64
	SectorsWrittenDelta  uint64

	NSEntry               NSEntry
	ReceiveBytesDelta     string
	ReceivePacketsDelta   uint64
	TransmitBytesDelta    string
	TransmitPacketsDelta  uint64
	ReceiveBytesNumDelta  uint64
	TransmitBytesNumDelta uint64

	// Extra exists to support customized data query.
	Extra []byte
}

type ProcSlice []Proc

func (p ProcSlice) Len() int      { return len(p) }
func (p ProcSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ProcSlice) Less(i, j int) bool {
	if p[i].UnixNanosecond != p[j].UnixNanosecond {
		return p[i].UnixNanosecond < p[j].UnixNanosecond
	}
	return p[i].UnixSecond < p[j].UnixSecond
}

// GetProc returns current 'Proc' data.
// PID is required.
// Disk device, network interface, extra path are optional.
func GetProc(opts ...FilterFunc) (Proc, error) {
	ft := &EntryFilter{}
	ft.applyOpts(opts)

	if ft.PID == 0 {
		return Proc{}, fmt.Errorf("unknown PID %d", ft.PID)
	}
	ts := time.Now().UnixNano()
	proc := Proc{UnixNanosecond: ts, UnixSecond: ConvertUnixNano(ts)}

	errc := make(chan error)
	go func() {
		// get process stats
		ets, err := GetPS(WithPID(ft.PID))
		if err != nil {
			errc <- err
			return
		}
		if len(ets) != 1 {
			errc <- fmt.Errorf("len(PID=%d entries) != 1 (got %d)", ft.PID, len(ets))
			return
		}
		proc.PSEntry = ets[0]
		errc <- nil
	}()

	go func() {
		lvg, err := GetProcLoadAvg()
		if err != nil {
			errc <- err
			return
		}
		proc.LoadAvg = lvg
		errc <- nil
	}()

	if ft.DiskDevice != "" {
		go func() {
			// get diskstats
			ds, err := GetDS()
			if err != nil {
				errc <- err
				return
			}
			for _, elem := range ds {
				if elem.Device == ft.DiskDevice {
					proc.DSEntry = elem
					break
				}
			}
			errc <- nil
		}()
	}

	if ft.NetworkInterface != "" {
		go func() {
			// get network I/O stats
			ns, err := GetNS()
			if err != nil {
				errc <- err
				return
			}
			for _, elem := range ns {
				if elem.Interface == ft.NetworkInterface {
					proc.NSEntry = elem
					break
				}
			}
			errc <- nil
		}()
	}

	if ft.ExtraPath != "" {
		go func() {
			f, err := openToRead(ft.ExtraPath)
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
	}

	cnt := 0
	for cnt != len(opts)+1 { // include load avg query
		err := <-errc
		if err != nil {
			return Proc{}, err
		}
		cnt++
	}

	if ft.DiskDevice != "" {
		if proc.DSEntry.Device == "" {
			return Proc{}, fmt.Errorf("disk device %q was not found", ft.DiskDevice)
		}
	}
	if ft.NetworkInterface != "" {
		if proc.NSEntry.Interface == "" {
			return Proc{}, fmt.Errorf("network interface %q was not found", ft.NetworkInterface)
		}
	}
	return proc, nil
}

var (
	// ProcHeader lists all Proc CSV columns.
	ProcHeader = append([]string{"UNIX-NANOSECOND", "UNIX-SECOND"}, columnsPSEntry...)

	// ProcHeaderIndex maps each Proc column name to its index in row.
	ProcHeaderIndex = make(map[string]int)
)

func init() {
	// more columns to 'ProcHeader'
	ProcHeader = append(ProcHeader,
		"LOAD-AVERAGE-1-MINUTE",
		"LOAD-AVERAGE-5-MINUTE",
		"LOAD-AVERAGE-15-MINUTE",
	)
	ProcHeader = append(ProcHeader, columnsDSEntry...)
	ProcHeader = append(ProcHeader, columnsNSEntry...)
	ProcHeader = append(ProcHeader,
		"READS-COMPLETED-DELTA",
		"SECTORS-READ-DELTA",
		"WRITES-COMPLETED-DELTA",
		"SECTORS-WRITTEN-DELTA",

		"RECEIVE-BYTES-DELTA",
		"RECEIVE-PACKETS-DELTA",
		"TRANSMIT-BYTES-DELTA",
		"TRANSMIT-PACKETS-DELTA",
		"RECEIVE-BYTES-NUM-DELTA",
		"TRANSMIT-BYTES-NUM-DELTA",

		"EXTRA",
	)

	for i, v := range ProcHeader {
		ProcHeaderIndex[v] = i
	}
}

// ToRow converts 'Proc' to string slice.
// Make sure to change this whenever 'Proc' fields are updated.
func (p *Proc) ToRow() (row []string) {
	row = make([]string, len(ProcHeader))
	row[0] = fmt.Sprintf("%d", p.UnixNanosecond) // UNIX-NANOSECOND
	row[1] = fmt.Sprintf("%d", p.UnixSecond)     // UNIX-SECOND

	row[2] = p.PSEntry.Program                       // PROGRAM
	row[3] = p.PSEntry.State                         // STATE
	row[4] = fmt.Sprintf("%d", p.PSEntry.PID)        // PID
	row[5] = fmt.Sprintf("%d", p.PSEntry.PPID)       // PPID
	row[6] = p.PSEntry.CPU                           // CPU
	row[7] = p.PSEntry.VMRSS                         // VMRSS
	row[8] = p.PSEntry.VMSize                        // VMSIZE
	row[9] = fmt.Sprintf("%d", p.PSEntry.FD)         // FD
	row[10] = fmt.Sprintf("%d", p.PSEntry.Threads)   // THREADS
	row[11] = fmt.Sprintf("%d", p.PSEntry.Threads)   // VOLUNTARY-CTXT-SWITCHES
	row[12] = fmt.Sprintf("%d", p.PSEntry.Threads)   // NON-VOLUNTARY-CTXT-SWITCHES
	row[13] = fmt.Sprintf("%3.2f", p.PSEntry.CPUNum) // CPU-NUM
	row[14] = fmt.Sprintf("%d", p.PSEntry.VMRSSNum)  // VMRSS-NUM
	row[15] = fmt.Sprintf("%d", p.PSEntry.VMSizeNum) // VMSIZE-NUM

	row[16] = fmt.Sprintf("%3.2f", p.LoadAvg.LoadAvg1Minute)  // LOAD-AVERAGE-1-MINUTE
	row[17] = fmt.Sprintf("%3.2f", p.LoadAvg.LoadAvg5Minute)  // LOAD-AVERAGE-5-MINUTE
	row[18] = fmt.Sprintf("%3.2f", p.LoadAvg.LoadAvg15Minute) // LOAD-AVERAGE-15-MINUTE

	row[19] = p.DSEntry.Device                                  // DEVICE
	row[20] = fmt.Sprintf("%d", p.DSEntry.ReadsCompleted)       // READS-COMPLETED
	row[21] = fmt.Sprintf("%d", p.DSEntry.SectorsRead)          // SECTORS-READ
	row[22] = p.DSEntry.TimeSpentOnReading                      // TIME(READS)
	row[23] = fmt.Sprintf("%d", p.DSEntry.WritesCompleted)      // WRITES-COMPLETED
	row[24] = fmt.Sprintf("%d", p.DSEntry.SectorsWritten)       // SECTORS-WRITTEN
	row[25] = p.DSEntry.TimeSpentOnWriting                      // TIME(WRITES)
	row[26] = fmt.Sprintf("%d", p.DSEntry.TimeSpentOnReadingMs) // MILLISECONDS(READS)
	row[27] = fmt.Sprintf("%d", p.DSEntry.TimeSpentOnWritingMs) // MILLISECONDS(WRITES)

	row[28] = p.NSEntry.Interface                           // INTERFACE
	row[29] = p.NSEntry.ReceiveBytes                        // RECEIVE-BYTES
	row[30] = fmt.Sprintf("%d", p.NSEntry.ReceivePackets)   // RECEIVE-PACKETS
	row[31] = p.NSEntry.TransmitBytes                       // TRANSMIT-BYTES
	row[32] = fmt.Sprintf("%d", p.NSEntry.TransmitPackets)  // TRANSMIT-PACKETS
	row[33] = fmt.Sprintf("%d", p.NSEntry.ReceiveBytesNum)  // RECEIVE-BYTES-NUM
	row[34] = fmt.Sprintf("%d", p.NSEntry.TransmitBytesNum) // TRANSMIT-BYTES-NUM

	row[35] = fmt.Sprintf("%d", p.ReadsCompletedDelta)  // READS-COMPLETED-DELTA
	row[36] = fmt.Sprintf("%d", p.SectorsReadDelta)     // SECTORS-READ-DELTA
	row[37] = fmt.Sprintf("%d", p.WritesCompletedDelta) // WRITES-COMPLETED-DELTA
	row[38] = fmt.Sprintf("%d", p.SectorsWrittenDelta)  // SECTORS-WRITTEN-DELTA

	row[39] = p.ReceiveBytesDelta                        // RECEIVE-BYTES-DELTA
	row[40] = fmt.Sprintf("%d", p.ReceivePacketsDelta)   // RECEIVE-PACKETS-DELTA
	row[41] = p.TransmitBytesDelta                       // TRANSMIT-BYTES-DELTA
	row[42] = fmt.Sprintf("%d", p.TransmitPacketsDelta)  // TRANSMIT-PACKETS-DELTA
	row[43] = fmt.Sprintf("%d", p.ReceiveBytesNumDelta)  // RECEIVE-BYTES-NUM-DELTA
	row[44] = fmt.Sprintf("%d", p.TransmitBytesNumDelta) // TRANSMIT-BYTES-NUM-DELTA

	row[45] = string(p.Extra) // EXTRA

	return
}
