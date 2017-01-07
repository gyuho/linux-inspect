package psn

import (
	"encoding/csv"
	"fmt"
	"time"
)

var psCSVColumns = append([]string{"unix-ts"}, columnsPSEntry...)

func init() {
	// more columns to 'psCSVColumns'
	psCSVColumns = append(psCSVColumns, columnsDSEntry...)
	psCSVColumns = append(psCSVColumns, columnsNSEntry...)
}

func getRow(pid int64, diskDevice string, networkInterface string) ([]string, error) {
	// get process stats
	ets, err := GetPS(WithPID(pid))
	if err != nil {
		return nil, err
	}
	if len(ets) != 1 {
		return nil, fmt.Errorf("len(PID=%d entries) != 1 (got %d)", pid, len(ets))
	}
	entry := ets[0]

	// get diskstats
	ds, err := GetDS()
	if err != nil {
		return nil, err
	}
	var dentry DSEntry
	for _, elem := range ds {
		if elem.Device == diskDevice {
			dentry = elem
			break
		}
	}
	if dentry.Device == "" {
		return nil, fmt.Errorf("disk device %q was not found", diskDevice)
	}

	// get network I/O stats
	ns, err := GetNS()
	if err != nil {
		return nil, err
	}
	var nentry NSEntry
	for _, elem := range ns {
		if elem.Interface == networkInterface {
			nentry = elem
			break
		}
	}
	if nentry.Interface == "" {
		return nil, fmt.Errorf("network interface %q was not found", networkInterface)
	}

	row := make([]string, len(psCSVColumns))
	row[0] = fmt.Sprintf("%d", time.Now().Unix()) // unix-ts

	row[1] = entry.Program                       // PROGRAM
	row[2] = entry.State                         // STATE
	row[3] = fmt.Sprintf("%d", entry.PID)        // PID
	row[4] = fmt.Sprintf("%d", entry.PPID)       // PPID
	row[5] = entry.CPU                           // CPU
	row[6] = entry.VMRSS                         // VMRSS
	row[7] = entry.VMSize                        // VMSIZE
	row[8] = fmt.Sprintf("%d", entry.FD)         // FD
	row[9] = fmt.Sprintf("%d", entry.Threads)    // THREADS
	row[10] = fmt.Sprintf("%3.2f", entry.CPUNum) // CPU-NUM
	row[11] = fmt.Sprintf("%d", entry.VMRSSNum)  // VMRSS-NUM
	row[12] = fmt.Sprintf("%d", entry.VMSizeNum) // VMSIZE-NUM

	row[13] = dentry.Device                                  // DEVICE
	row[14] = fmt.Sprintf("%d", dentry.ReadsCompleted)       // READS-COMPLETED
	row[15] = fmt.Sprintf("%d", dentry.SectorsRead)          // SECTORS-READ
	row[16] = dentry.TimeSpentOnReading                      // TIME(READS)
	row[17] = fmt.Sprintf("%d", dentry.WritesCompleted)      // WRITES-COMPLETED
	row[18] = fmt.Sprintf("%d", dentry.SectorsWritten)       // SECTORS-WRITE
	row[19] = dentry.TimeSpentOnWriting                      // TIME(WRITES)
	row[20] = fmt.Sprintf("%d", dentry.TimeSpentOnReadingMs) // MILLISECONDS(READS)
	row[21] = fmt.Sprintf("%d", dentry.TimeSpentOnWritingMs) // MILLISECONDS(WRITES)

	row[22] = nentry.Interface                           // INTERFACE
	row[23] = nentry.ReceiveBytes                        // RECEIVE-BYTES
	row[24] = fmt.Sprintf("%d", nentry.ReceivePackets)   // RECEIVE-PACKETS
	row[25] = nentry.TransmitBytes                       // TRANSMIT-BYTES
	row[26] = fmt.Sprintf("%d", nentry.TransmitPackets)  // TRANSMIT-PACKETS
	row[27] = fmt.Sprintf("%d", nentry.ReceiveBytesNum)  // RECEIVE-BYTES-NUM
	row[28] = fmt.Sprintf("%d", nentry.TransmitBytesNum) // TRANSMIT-BYTES-NUM

	return row, nil
}

// WriteSSCSVHeader writes header to a CSV file.
// Meant to be called only once. Optional columns
// (diskstats, network I/O stats) may be added.
func WriteSSCSVHeader(fpath string) error {
	f, err := openToAppend(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	if err := wr.Write(psCSVColumns); err != nil {
		return err
	}
	wr.Flush()
	return wr.Error()
}

// AppendSSCSV is called periodically (usually per-second)
// to collect all data wit psn utilities and keeps appending
// rows to a CSV file. Optional columns (diskstats, network I/O stats)
// may be added.
func AppendSSCSV(fpath string, pid int64, diskDevice string, networkInterface string) error {
	f, err := openToAppend(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	row, err := getRow(pid, diskDevice, networkInterface)
	if err != nil {
		return err
	}
	wr := csv.NewWriter(f)
	if err := wr.Write(row); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}
