package ps

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Table represents CSV tables.
type Table struct {
	// Columns maps its column name to its column index.
	Columns map[string]int
	Rows    [][]string

	minTS int64
	maxTS int64
}

var (
	columnsPS = make(map[string]int)
)

func init() {
	for i, v := range append([]string{"unix_ts"}, ProcessTableColumns...) {
		columnsPS[v] = i
	}
}

// ReadCSV reads a csv file for ps csv file.
// It assumes the results from one single program.
func ReadCSV(fpath string) (Table, error) {
	f, err := openToRead(fpath)
	if err != nil {
		return Table{}, err
	}
	defer f.Close()

	rd := csv.NewReader(f)

	// in case that rows have different number of fields
	rd.FieldsPerRecord = -1

	rows, err := rd.ReadAll()
	if err != nil {
		return Table{}, err
	}
	if len(rows) > 0 && len(rows[0]) > 0 {
		if rows[0][0] == "unix_ts" {
			rows = rows[1:]
		}
	}
	min, err := strconv.ParseInt(rows[0][0], 10, 64)
	if err != nil {
		return Table{}, err
	}
	max, err := strconv.ParseInt(rows[len(rows)-1][0], 10, 64)
	if err != nil {
		return Table{}, err
	}
	return Table{
		Columns: columnsPS,
		Rows:    rows,
		minTS:   min,
		maxTS:   max,
	}, nil
}

// ReadCSVs reads multiple csv files, including only unix timestamps,
// CPU usage, and VmRSS in MB. It assumes the results from one single program.
func ReadCSVs(fpaths ...string) (Table, error) {
	tbs := []Table{}
	var minTS, maxTS int64
	for i, fpath := range fpaths {
		tb, err := ReadCSV(fpath)
		if err != nil {
			return Table{}, err
		}
		tbs = append(tbs, tb)

		if i == 0 {
			minTS = tb.minTS
			maxTS = tb.maxTS
		}
		if minTS > tb.minTS {
			minTS = tb.minTS
		}
		if maxTS < tb.maxTS {
			maxTS = tb.maxTS
		}
	}

	ntb := Table{}
	ntb.Columns = make(map[string]int)
	ntb.Columns["unix_tx"] = 0
	for i := range fpaths {
		ntb.Columns[fmt.Sprintf("cpu_%d", i+1)] = 2*i + 1
		ntb.Columns[fmt.Sprintf("memory_mb_%d", i+1)] = 2*i + 2
	}

	cIdx, mIdx := tbs[0].Columns["CpuUsageFloat64"], tbs[0].Columns["VmRSSBytes"]

	nrows := [][]string{}
	for i, tb := range tbs {
		if minTS < tb.minTS { // need to fill in top rows
			// push-front from minTS to tb.minTS-1
			rows := make([][]string, tb.minTS-minTS)
			for i := range rows {
				emptyRow := append([]string{fmt.Sprintf("%d", minTS+int64(i))}, strings.Split(strings.Repeat("0.00_", len(ProcessTableColumns)), "_")...)
				rows[i] = emptyRow
			}
			tb.Rows = append(rows, tb.Rows...)
		}
		if maxTS > tb.maxTS { // need to fill in bottom rows
			// push-back from tb.maxTS+1 to maxTS
			rows := make([][]string, maxTS-tb.maxTS)
			for i := range rows {
				emptyRow := append([]string{fmt.Sprintf("%d", tb.maxTS+int64(i))}, strings.Split(strings.Repeat("0.00_", len(ProcessTableColumns)), "_")...)
				rows[i] = emptyRow
			}
			tb.Rows = append(tb.Rows, rows...)
		}

		if i == 0 {
			for _, row := range tb.Rows {
				if row[0] == "unix_ts" {
					continue
				}
				mf, _ := strconv.ParseFloat(row[mIdx], 64)
				nrows = append(nrows, []string{row[0], row[cIdx], fmt.Sprintf("%.2f", mf/(1000*1000))}) // memory usage in mega-bytes
			}
		} else {
			for rowNum := range nrows {
				trow := tb.Rows[rowNum]
				if trow[0] == "unix_ts" {
					continue
				}
				mf, _ := strconv.ParseFloat(trow[mIdx], 64)
				nrows[rowNum] = append(nrows[rowNum], []string{trow[cIdx], fmt.Sprintf("%.2f", mf/(1000*1000))}...)
			}
		}
	}
	ntb.Rows = nrows

	return ntb, nil
}

// ToRows returns rows from the table.
func (t Table) ToRows() [][]string {
	cols := make([]string, len(t.Columns))
	for k, v := range t.Columns {
		cols[v] = k
	}
	return append([][]string{cols}, t.Rows...)
}

// ToCSV saves the table to csv file.
func (t Table) ToCSV(fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	wr := csv.NewWriter(f)

	cols := make([]string, len(t.Columns))
	for k, v := range t.Columns {
		cols[v] = k
	}
	if err := wr.Write(cols); err != nil {
		return err
	}

	if err := wr.WriteAll(t.Rows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}
