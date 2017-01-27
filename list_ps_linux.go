package psn

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/gyuho/dataframe"
	"github.com/olekukonko/tablewriter"
)

// PSEntry is a process entry.
// Simplied from 'Stat' and 'Status'.
type PSEntry struct {
	Program string
	State   string
	PID     int64
	PPID    int64

	CPU    string
	VMRSS  string
	VMSize string

	FD      uint64
	Threads uint64

	VoluntaryCtxtSwitches    uint64
	NonvoluntaryCtxtSwitches uint64

	// extra fields for sorting
	CPUNum    float64
	VMRSSNum  uint64
	VMSizeNum uint64
}

const maxConcurrentProcStat = 32

// GetPS finds all PSEntry by given filter.
func GetPS(opts ...FilterFunc) (pss []PSEntry, err error) {
	ft := &EntryFilter{}
	ft.applyOpts(opts)

	var pids []int64
	switch {
	case ft.ProgramMatchFunc == nil && ft.PID < 1:
		// get all PIDs
		pids, err = ListPIDs()
		if err != nil {
			return
		}

	case ft.PID > 0:
		pids = []int64{ft.PID}

	case ft.ProgramMatchFunc != nil:
		// later to find PIDs by Program
		pids = nil

	default:
		// applyOpts already panic when ft.ProgramMatchFunc != nil && ft.PID > 0
	}

	up, err := GetProcUptime()
	if err != nil {
		return nil, err
	}

	if len(pids) == 0 {
		// find PIDs by Program
		if pids, err = ListPIDs(); err != nil {
			return
		}
	} else {
		ft.ProgramMatchFunc = func(string) bool { return true }
	}

	var pmu sync.RWMutex
	var wg sync.WaitGroup
	wg.Add(len(pids))
	limitc := make(chan struct{}, maxConcurrentProcStat)
	for _, pid := range pids {
		go func(pid int64) {
			defer func() {
				<-limitc
				wg.Done()
			}()

			limitc <- struct{}{}

			stat, err := GetProcStatByPID(pid, up)
			if err != nil {
				log.Printf("GetProcStatByPID error %v for PID %d", err, pid)
				return
			}
			if !ft.ProgramMatchFunc(stat.Comm) {
				return
			}

			pmu.RLock()
			done := ft.TopLimit > 0 && len(pss) >= ft.TopLimit
			pmu.RUnlock()
			if done {
				return
			}

			ent, err := getPSEntry(pid, stat)
			if err != nil {
				log.Printf("getPSEntry error %v for PID %d", err, pid)
				return
			}

			pmu.Lock()
			pss = append(pss, ent)
			pmu.Unlock()
		}(pid)
	}
	wg.Wait()

	if ft.TopLimit > 0 && len(pss) > ft.TopLimit {
		pss = pss[:ft.TopLimit:ft.TopLimit]
	}
	return
}

func getPSEntry(pid int64, stat Stat) (PSEntry, error) {
	status, err := GetProcStatusByPID(pid)
	if err != nil {
		return PSEntry{}, err
	}

	entry := PSEntry{
		Program: stat.Comm,
		State:   stat.StateParsedStatus,

		PID:  stat.Pid,
		PPID: stat.Ppid,

		CPU:    fmt.Sprintf("%3.2f %%", stat.CpuUsage),
		VMRSS:  status.VmRSSParsedBytes,
		VMSize: status.VmSizeParsedBytes,

		FD:      status.FDSize,
		Threads: status.Threads,

		VoluntaryCtxtSwitches:    status.VoluntaryCtxtSwitches,
		NonvoluntaryCtxtSwitches: status.NonvoluntaryCtxtSwitches,

		CPUNum:    stat.CpuUsage,
		VMRSSNum:  status.VmRSSBytesN,
		VMSizeNum: status.VmSizeBytesN,
	}

	if status.Name != "" {
		entry.Program = status.Name
	}
	if status.StateParsedStatus != "" {
		entry.State = status.StateParsedStatus
	}

	return entry, nil
}

const columnsPSToShow = 11

var columnsPSEntry = []string{
	"PROGRAM",

	"STATE",
	"PID",
	"PPID",

	"CPU",
	"VMRSS",
	"VMSIZE",

	"FD",
	"THREADS",

	"VOLUNTARY-CTXT-SWITCHES",
	"NON-VOLUNTARY-CTXT-SWITCHES",

	// extra for sorting
	"CPU-NUM",
	"VMRSS-NUM",
	"VMSIZE-NUM",
}

// ConvertPS converts to rows.
func ConvertPS(nss ...PSEntry) (header []string, rows [][]string) {
	header = columnsPSEntry
	rows = make([][]string, len(nss))
	for i, elem := range nss {
		row := make([]string, len(columnsPSEntry))
		row[0] = elem.Program

		row[1] = elem.State
		row[2] = fmt.Sprintf("%d", elem.PID)
		row[3] = fmt.Sprintf("%d", elem.PPID)

		row[4] = elem.CPU
		row[5] = elem.VMRSS
		row[6] = elem.VMSize

		row[7] = fmt.Sprintf("%d", elem.FD)
		row[8] = fmt.Sprintf("%d", elem.Threads)

		row[9] = fmt.Sprintf("%d", elem.VoluntaryCtxtSwitches)
		row[10] = fmt.Sprintf("%d", elem.NonvoluntaryCtxtSwitches)

		row[11] = fmt.Sprintf("%3.2f", elem.CPUNum)
		row[12] = fmt.Sprintf("%d", elem.VMRSSNum)
		row[13] = fmt.Sprintf("%d", elem.VMSizeNum)

		rows[i] = row
	}
	dataframe.SortBy(
		rows,
		dataframe.Float64DescendingFunc(12), // VMRSSNum
		dataframe.Float64DescendingFunc(11), // CPUNum
		dataframe.Float64DescendingFunc(13), // VMSizeNum
	).Sort(rows)

	return
}

// StringPS converts in print-friendly format.
func StringPS(header []string, rows [][]string, topLimit int) string {
	buf := new(bytes.Buffer)
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(header[:columnsPSToShow:columnsPSToShow])

	if topLimit > 0 && len(rows) > topLimit {
		rows = rows[:topLimit:topLimit]
	}

	for _, row := range rows {
		tw.Append(row[:columnsPSToShow:columnsPSToShow])
	}
	tw.SetAutoFormatHeaders(false)
	tw.SetAlignment(tablewriter.ALIGN_RIGHT)
	tw.Render()

	return buf.String()
}
