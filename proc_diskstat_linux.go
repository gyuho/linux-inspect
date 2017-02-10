package psn

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type procDiskstatsColumnIndex int

const (
	proc_diskstats_idx_major_number procDiskstatsColumnIndex = iota
	proc_diskstats_idx_minor_number
	proc_diskstats_idx_device_name

	proc_diskstats_idx_reads_completed
	proc_diskstats_idx_reads_merged
	proc_diskstats_idx_sectors_read
	proc_diskstats_idx_time_spent_on_reading_ms

	proc_diskstats_idx_writes_completed
	proc_diskstats_idx_writes_merged
	proc_diskstats_idx_sectors_written
	proc_diskstats_idx_time_spent_on_writing_ms

	proc_diskstats_idx_ios_in_progress
	proc_diskstats_idx_time_spent_on_ios_ms
	proc_diskstats_idx_weighted_time_spent_on_ios_ms
)

// GetProcDiskstats reads '/proc/diskstats'.
func GetProcDiskstats() ([]DiskStat, error) {
	f, err := openToRead("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dss := []DiskStat{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}
		ds := strings.Fields(strings.TrimSpace(txt))
		if len(ds) < int(proc_diskstats_idx_weighted_time_spent_on_ios_ms+1) {
			return nil, fmt.Errorf("not enough columns at %v", ds)
		}
		d := DiskStat{}

		mn, err := strconv.ParseUint(ds[proc_diskstats_idx_major_number], 10, 64)
		if err != nil {
			return nil, err
		}
		d.MajorNumber = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_minor_number], 10, 64)
		if err != nil {
			return nil, err
		}
		d.MinorNumber = mn

		d.DeviceName = strings.TrimSpace(ds[proc_diskstats_idx_device_name])

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_reads_completed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReadsCompleted = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_reads_merged], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReadsMerged = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_sectors_read], 10, 64)
		if err != nil {
			return nil, err
		}
		d.SectorsRead = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_time_spent_on_reading_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnReadingMs = mn
		d.TimeSpentOnReadingMsParsedTime = humanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_writes_completed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WritesCompleted = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_writes_merged], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WritesMerged = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_sectors_written], 10, 64)
		if err != nil {
			return nil, err
		}
		d.SectorsWritten = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_time_spent_on_writing_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnWritingMs = mn
		d.TimeSpentOnWritingMsParsedTime = humanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_ios_in_progress], 10, 64)
		if err != nil {
			return nil, err
		}
		d.IOsInProgress = mn

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_time_spent_on_ios_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnIOsMs = mn
		d.TimeSpentOnIOsMsParsedTime = humanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[proc_diskstats_idx_weighted_time_spent_on_ios_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WeightedTimeSpentOnIOsMs = mn
		d.WeightedTimeSpentOnIOsMsParsedTime = humanizeDurationMs(mn)

		dss = append(dss, d)
	}

	return dss, nil
}
