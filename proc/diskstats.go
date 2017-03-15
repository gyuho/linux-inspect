package proc

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/pkg/timeutil"
)

type diskstatsColumnIndex int

const (
	diskstats_idx_major_number diskstatsColumnIndex = iota
	diskstats_idx_minor_number
	diskstats_idx_device_name

	diskstats_idx_reads_completed
	diskstats_idx_reads_merged
	diskstats_idx_sectors_read
	diskstats_idx_time_spent_on_reading_ms

	diskstats_idx_writes_completed
	diskstats_idx_writes_merged
	diskstats_idx_sectors_written
	diskstats_idx_time_spent_on_writing_ms

	diskstats_idx_ios_in_progress
	diskstats_idx_time_spent_on_ios_ms
	diskstats_idx_weighted_time_spent_on_ios_ms
)

// GetDiskstats reads '/proc/diskstats'.
func GetDiskstats() ([]DiskStat, error) {
	f, err := fileutil.OpenToRead("/proc/diskstats")
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
		if len(ds) < int(diskstats_idx_weighted_time_spent_on_ios_ms+1) {
			return nil, fmt.Errorf("not enough columns at %v", ds)
		}
		d := DiskStat{}

		mn, err := strconv.ParseUint(ds[diskstats_idx_major_number], 10, 64)
		if err != nil {
			return nil, err
		}
		d.MajorNumber = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_minor_number], 10, 64)
		if err != nil {
			return nil, err
		}
		d.MinorNumber = mn

		d.DeviceName = strings.TrimSpace(ds[diskstats_idx_device_name])

		mn, err = strconv.ParseUint(ds[diskstats_idx_reads_completed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReadsCompleted = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_reads_merged], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReadsMerged = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_sectors_read], 10, 64)
		if err != nil {
			return nil, err
		}
		d.SectorsRead = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_time_spent_on_reading_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnReadingMs = mn
		d.TimeSpentOnReadingMsParsedTime = timeutil.HumanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[diskstats_idx_writes_completed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WritesCompleted = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_writes_merged], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WritesMerged = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_sectors_written], 10, 64)
		if err != nil {
			return nil, err
		}
		d.SectorsWritten = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_time_spent_on_writing_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnWritingMs = mn
		d.TimeSpentOnWritingMsParsedTime = timeutil.HumanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[diskstats_idx_ios_in_progress], 10, 64)
		if err != nil {
			return nil, err
		}
		d.IOsInProgress = mn

		mn, err = strconv.ParseUint(ds[diskstats_idx_time_spent_on_ios_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TimeSpentOnIOsMs = mn
		d.TimeSpentOnIOsMsParsedTime = timeutil.HumanizeDurationMs(mn)

		mn, err = strconv.ParseUint(ds[diskstats_idx_weighted_time_spent_on_ios_ms], 10, 64)
		if err != nil {
			return nil, err
		}
		d.WeightedTimeSpentOnIOsMs = mn
		d.WeightedTimeSpentOnIOsMsParsedTime = timeutil.HumanizeDurationMs(mn)

		dss = append(dss, d)
	}

	return dss, nil
}
