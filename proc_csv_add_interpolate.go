package psn

import (
	"fmt"
	"sort"

	humanize "github.com/dustin/go-humanize"
)

// Combine combines a list Proc and returns one combined Proc.
// Field values are estimated. UnixNanosecond is reset 0.
// And UnixSecond and other fields that cannot be averaged
// are set with the field value in the last element.
// This is meant to be used to combine Proc rows with duplicate
// unix second timestamps
func Combine(procs ...Proc) Proc {
	if len(procs) < 1 {
		return Proc{}
	}
	if len(procs) == 1 {
		return procs[0]
	}

	lastProc := procs[len(procs)-1]
	combined := lastProc
	combined.UnixNanosecond = 0

	var (
		// for PSEntry
		voluntaryCtxtSwitches    uint64
		nonVoluntaryCtxtSwitches uint64
		cpuNum                   float64
		vmRSSNum                 uint64
		vmSizeNum                uint64

		// for LoadAvg
		loadAvg1Minute                   float64
		loadAvg5Minute                   float64
		loadAvg15Minute                  float64
		runnableKernelSchedulingEntities int64
		currentKernelSchedulingEntities  int64

		// for DSEntry
		readsCompleted       uint64
		sectorsRead          uint64
		writesCompleted      uint64
		sectorsWritten       uint64
		timeSpentOnReadingMs uint64
		timeSpentOnWritingMs uint64

		// for DSEntry delta
		readsCompletedDelta  uint64
		sectorsReadDelta     uint64
		writesCompletedDelta uint64
		sectorsWrittenDelta  uint64

		// for NSEntry
		receivePackets   uint64
		transmitPackets  uint64
		receiveBytesNum  uint64
		transmitBytesNum uint64

		// for NSEntry delta
		receivePacketsDelta   uint64
		transmitPacketsDelta  uint64
		receiveBytesNumDelta  uint64
		transmitBytesNumDelta uint64
	)
	for _, p := range procs {
		// for PSEntry
		voluntaryCtxtSwitches += p.PSEntry.VoluntaryCtxtSwitches
		nonVoluntaryCtxtSwitches += p.PSEntry.NonvoluntaryCtxtSwitches
		cpuNum += p.PSEntry.CPUNum
		vmRSSNum += p.PSEntry.VMRSSNum
		vmSizeNum += p.PSEntry.VMSizeNum

		// for LoadAvg
		loadAvg1Minute += p.LoadAvg.LoadAvg1Minute
		loadAvg5Minute += p.LoadAvg.LoadAvg5Minute
		loadAvg15Minute += p.LoadAvg.LoadAvg15Minute
		runnableKernelSchedulingEntities += p.LoadAvg.RunnableKernelSchedulingEntities
		currentKernelSchedulingEntities += p.LoadAvg.CurrentKernelSchedulingEntities

		// for DSEntry
		readsCompleted += p.DSEntry.ReadsCompleted
		sectorsRead += p.DSEntry.SectorsRead
		writesCompleted += p.DSEntry.WritesCompleted
		sectorsWritten += p.DSEntry.SectorsWritten
		timeSpentOnReadingMs += p.DSEntry.TimeSpentOnReadingMs
		timeSpentOnWritingMs += p.DSEntry.TimeSpentOnWritingMs

		// for DSEntry delta
		readsCompletedDelta += p.ReadsCompletedDelta
		sectorsReadDelta += p.SectorsReadDelta
		writesCompletedDelta += p.WritesCompletedDelta
		sectorsWrittenDelta += p.SectorsWrittenDelta

		// for NSEntry
		receivePackets += p.NSEntry.ReceivePackets
		transmitPackets += p.NSEntry.TransmitPackets
		receiveBytesNum += p.NSEntry.ReceiveBytesNum
		transmitBytesNum += p.NSEntry.TransmitBytesNum

		// for NSEntry delta
		receivePacketsDelta += p.ReceivePacketsDelta
		transmitPacketsDelta += p.TransmitPacketsDelta
		receiveBytesNumDelta += p.ReceiveBytesNumDelta
		transmitBytesNumDelta += p.TransmitBytesNumDelta
	}
	pN := len(procs)

	// for PSEntry
	combined.PSEntry.VoluntaryCtxtSwitches = uint64(voluntaryCtxtSwitches) / uint64(pN)
	combined.PSEntry.NonvoluntaryCtxtSwitches = uint64(nonVoluntaryCtxtSwitches) / uint64(pN)
	combined.PSEntry.CPUNum = float64(cpuNum) / float64(pN)
	combined.PSEntry.CPU = fmt.Sprintf("%3.2f %%", combined.PSEntry.CPUNum)
	combined.PSEntry.VMRSSNum = uint64(vmRSSNum) / uint64(pN)
	combined.PSEntry.VMRSS = humanize.Bytes(combined.PSEntry.VMRSSNum)
	combined.PSEntry.VMSizeNum = uint64(vmSizeNum) / uint64(pN)
	combined.PSEntry.VMSize = humanize.Bytes(combined.PSEntry.VMSizeNum)

	// for LoadAvg
	combined.LoadAvg.LoadAvg1Minute = float64(loadAvg1Minute) / float64(pN)
	combined.LoadAvg.LoadAvg5Minute = float64(loadAvg5Minute) / float64(pN)
	combined.LoadAvg.LoadAvg15Minute = float64(loadAvg15Minute) / float64(pN)
	combined.LoadAvg.RunnableKernelSchedulingEntities = int64(loadAvg15Minute) / int64(pN)
	combined.LoadAvg.CurrentKernelSchedulingEntities = int64(loadAvg15Minute) / int64(pN)

	// for DSEntry
	combined.DSEntry.ReadsCompleted = uint64(readsCompleted) / uint64(pN)
	combined.DSEntry.SectorsRead = uint64(sectorsRead) / uint64(pN)
	combined.DSEntry.WritesCompleted = uint64(writesCompleted) / uint64(pN)
	combined.DSEntry.SectorsRead = uint64(sectorsRead) / uint64(pN)
	combined.DSEntry.TimeSpentOnReadingMs = uint64(timeSpentOnReadingMs) / uint64(pN)
	combined.DSEntry.TimeSpentOnReading = humanizeDurationMs(combined.DSEntry.TimeSpentOnReadingMs)
	combined.DSEntry.TimeSpentOnWritingMs = uint64(timeSpentOnWritingMs) / uint64(pN)
	combined.DSEntry.TimeSpentOnWriting = humanizeDurationMs(combined.DSEntry.TimeSpentOnWritingMs)
	combined.ReadsCompletedDelta = uint64(readsCompletedDelta) / uint64(pN)
	combined.SectorsReadDelta = uint64(sectorsReadDelta) / uint64(pN)
	combined.WritesCompletedDelta = uint64(writesCompletedDelta) / uint64(pN)
	combined.SectorsWrittenDelta = uint64(sectorsWrittenDelta) / uint64(pN)

	// for NSEntry
	combined.NSEntry.ReceiveBytesNum = uint64(receiveBytesNum) / uint64(pN)
	combined.NSEntry.TransmitBytesNum = uint64(transmitBytesNum) / uint64(pN)
	combined.NSEntry.ReceivePackets = uint64(receivePackets) / uint64(pN)
	combined.NSEntry.TransmitPackets = uint64(transmitPackets) / uint64(pN)
	combined.NSEntry.ReceiveBytes = humanize.Bytes(combined.NSEntry.ReceiveBytesNum)
	combined.NSEntry.TransmitBytes = humanize.Bytes(combined.NSEntry.TransmitBytesNum)
	combined.ReceivePacketsDelta = uint64(receivePacketsDelta) / uint64(pN)
	combined.TransmitPacketsDelta = uint64(transmitPacketsDelta) / uint64(pN)
	combined.ReceiveBytesNumDelta = uint64(receiveBytesNumDelta) / uint64(pN)
	combined.ReceiveBytesDelta = humanize.Bytes(combined.ReceiveBytesNumDelta)
	combined.TransmitBytesNumDelta = uint64(transmitBytesNumDelta) / uint64(pN)
	combined.TransmitBytesDelta = humanize.Bytes(combined.TransmitBytesNumDelta)

	return combined
}

// Interpolate interpolates missing rows in CSV assuming CSV is to be collected for every second.
// 'Missing' means unix seconds in rows are not continuous.
// It fills in the empty rows by estimating the averages.
// It returns a new copy of CSV. And the new copy sets all unix nanoseconds to 0.,
// since it's now aggregated by the unix "second".
func (c *CSV) Interpolate() (cc *CSV, err error) {
	if c == nil || len(c.Rows) < 2 {
		// no need to interpolate
		return
	}

	// copy the original CSV data
	cc = &(*c)

	// find missing rows, assuming CSV is to be collected every second
	if cc.MinUnixSecond == cc.MaxUnixSecond {
		// no need to interpolate
		return
	}

	// min unix second is 5, max is 7
	// then the expected row number is 7-5+1=3
	expectedRowN := cc.MaxUnixSecond - cc.MinUnixSecond + 1
	secondToAllProcs := make(map[int64][]Proc)
	for _, row := range cc.Rows {
		if _, ok := secondToAllProcs[row.UnixSecond]; ok {
			secondToAllProcs[row.UnixSecond] = append(secondToAllProcs[row.UnixSecond], row)
		} else {
			secondToAllProcs[row.UnixSecond] = []Proc{row}
		}
	}
	if int64(len(cc.Rows)) == expectedRowN && len(cc.Rows) == len(secondToAllProcs) {
		// all rows have distinct unix second
		// and they are all continuous unix seconds
		return
	}

	// interpolate cases
	//
	// case #1. If duplicate rows are found (equal/different unix nanoseconds, equal unix seconds),
	//          combine those into one row with its average.
	//
	// case #2. If some rows are discontinuous in unix seconds, there are missing rows.
	//          Fill in those rows with average estimates.

	// case #1, find duplicate rows!
	// It finds duplicates by unix second! Not by unix nanoseconds!
	secondToProc := make(map[int64]Proc)
	for sec, procs := range secondToAllProcs {
		if len(procs) == 0 {
			return nil, fmt.Errorf("empty row found at unix second %d", sec)
		}

		if len(procs) == 1 {
			secondToProc[sec] = procs[0]
			continue // no need to combine
		}

		// procs conflicted on unix second,
		// we want to combine those into one
		secondToProc[sec] = Combine(procs...)
	}

	// sort and reset the unix second
	rows2 := make([]Proc, 0, len(secondToProc))
	allUnixSeconds := make([]int64, 0, len(secondToProc))
	for _, row := range secondToProc {
		row.UnixNanosecond = 0
		rows2 = append(rows2, row)
		allUnixSeconds = append(allUnixSeconds, row.UnixSecond)
	}
	sort.Sort(ProcSlice(rows2))

	cc.Rows = rows2
	cc.MinUnixNanosecond = rows2[0].UnixNanosecond
	cc.MinUnixSecond = rows2[0].UnixSecond
	cc.MaxUnixNanosecond = rows2[len(rows2)-1].UnixNanosecond
	cc.MaxUnixSecond = rows2[len(rows2)-1].UnixSecond

	// case #2, find missing rows!
	// if unix seconds have discontinued ranges, it's missing some rows!
	missingTS := make(map[int64]struct{})
	for unixSecond := cc.MinUnixSecond; unixSecond <= cc.MaxUnixSecond; unixSecond++ {
		_, ok := secondToProc[unixSecond]
		if !ok {
			missingTS[unixSecond] = struct{}{}
		}
	}
	if len(missingTS) == 0 {
		// now all rows have distinct unix second
		// and there's no missing unix seconds
		return
	}

	// now we need to estimate the Proc for missingTS
	// fmt.Printf("total %d points available, missing %d points\n", len(allUnixSeconds), len(missingTS))
	bds := buildBoundaries(allUnixSeconds)
	_ = bds

	return
}

// ConvertUnixNano unix nanoseconds to unix second.
func ConvertUnixNano(unixNano int64) (unixSec int64) {
	return int64(unixNano / 1e9)
}
