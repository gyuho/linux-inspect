package psn

import (
	"bytes"
	"fmt"

	"github.com/gyuho/dataframe"
	"github.com/olekukonko/tablewriter"
)

// DSEntry represents disk statistics.
// Simplied from 'DiskStat'.
type DSEntry struct {
	Device string

	ReadsCompleted     uint64
	SectorsRead        uint64
	TimeSpentOnReading string

	WritesCompleted    uint64
	SectorsWritten     uint64
	TimeSpentOnWriting string

	// extra fields for sorting
	TimeSpentOnReadingMs uint64
	TimeSpentOnWritingMs uint64
}

var columnsDSEntry = []string{
	"DEVICE",
	"READS-COMPLETED", "SECTORS-READ", "TIME(READS)",
	"WRITES-COMPLETED", "SECTORS-WRITE", "TIME(WRITES)",

	// extra for sorting
	"MILLISECONDS(READS)",
	"MILLISECONDS(WRITES)",
}

// GetDS lists all disk statistics.
func GetDS() ([]DSEntry, error) {
	ss, err := GetDiskstats()
	if err != nil {
		return nil, err
	}
	ds := make([]DSEntry, len(ss))
	for i := range ss {
		ds[i] = DSEntry{
			Device: ss[i].DeviceName,

			ReadsCompleted:     ss[i].ReadsCompleted,
			SectorsRead:        ss[i].SectorsRead,
			TimeSpentOnReading: ss[i].TimeSpentOnReadingMsParsedTime,

			WritesCompleted:    ss[i].WritesCompleted,
			SectorsWritten:     ss[i].SectorsWritten,
			TimeSpentOnWriting: ss[i].TimeSpentOnWritingMsParsedTime,

			TimeSpentOnReadingMs: ss[i].TimeSpentOnReadingMs,
			TimeSpentOnWritingMs: ss[i].TimeSpentOnWritingMs,
		}
	}
	return ds, nil
}

// ConvertDS converts to rows.
func ConvertDS(dss ...DSEntry) (header []string, rows [][]string) {
	header = columnsDSEntry
	rows = make([][]string, len(dss))
	for i, elem := range dss {
		row := make([]string, len(columnsDSEntry))
		row[0] = elem.Device

		row[1] = fmt.Sprintf("%d", elem.ReadsCompleted)
		row[2] = fmt.Sprintf("%d", elem.SectorsRead)
		row[3] = elem.TimeSpentOnReading

		row[4] = fmt.Sprintf("%d", elem.WritesCompleted)
		row[5] = fmt.Sprintf("%d", elem.SectorsWritten)
		row[6] = elem.TimeSpentOnWriting

		row[7] = fmt.Sprintf("%d", elem.TimeSpentOnReadingMs)
		row[8] = fmt.Sprintf("%d", elem.TimeSpentOnWritingMs)

		rows[i] = row
	}
	dataframe.SortBy(
		rows,
		dataframe.NumberDescendingFunc(5), // SECTORS-WRITE
		dataframe.NumberDescendingFunc(4), // SECTORS-READ
	).Sort(rows)

	return
}

// StringDS converts in print-friendly format.
func StringDS(header []string, rows [][]string, topLimit int) string {
	buf := new(bytes.Buffer)
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(header[:7:7])

	if topLimit > 0 && len(rows) > topLimit {
		rows = rows[:topLimit:topLimit]
	}

	for _, row := range rows {
		tw.Append(row[:7:7])
	}
	tw.Render()

	return buf.String()
}
