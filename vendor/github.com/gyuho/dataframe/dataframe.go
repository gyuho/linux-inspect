package dataframe

import (
	"encoding/csv"
	"fmt"
	"sync"
)

// Frame contains data.
type Frame interface {
	// Headers returns the slice of headers in order. Header name is unique among its Frame.
	Headers() []string

	// AddColumn adds a Column to Frame.
	AddColumn(c Column) error

	// Column returns the Column by its header name.
	Column(header string) (Column, error)

	// Columns returns all Columns.
	Columns() []Column

	// Count returns the number of Columns in the Frame.
	Count() int

	// UpdateHeader updates the header name of a Column.
	UpdateHeader(origHeader, newHeader string) error

	// MoveColumn moves the column right before the target index.
	MoveColumn(header string, target int) error

	// DeleteColumn deletes the Column by its header.
	DeleteColumn(header string) bool

	// CSV saves the Frame to a CSV file.
	CSV(fpath string) error

	// CSVHorizontal saves the Frame to a CSV file
	// in a horizontal way. The first column is header.
	// And data are aligned from left to right.
	CSVHorizontal(fpath string) error

	// Rows returns the header and data slices.
	Rows() ([]string, [][]string)

	// Sort sorts the Frame.
	Sort(header string, st SortType, so SortOption) error
}

type frame struct {
	mu       sync.Mutex
	columns  []Column
	headerTo map[string]int
}

// New returns a new Frame.
func New() Frame {
	return &frame{
		columns:  []Column{},
		headerTo: make(map[string]int),
	}
}

// NewFromRows creates Frame from rows.
// Pass 'nil' header if first row is used as header strings.
// Pass 'non-nil' header if the data starts from the first row, without header strings.
func NewFromRows(header []string, rows [][]string) (Frame, error) {
	if len(rows) < 1 {
		return nil, fmt.Errorf("empty row %q", rows)
	}
	fr := New()
	headerN := len(header)
	if headerN > 0 { // use this as header
		// assume no header string at top
		cols := make([]Column, headerN)
		for i := range cols {
			cols[i] = NewColumn(header[i])
		}
		for _, row := range rows {
			rowN := len(row)
			if rowN > headerN {
				return nil, fmt.Errorf("header %q is not specified correctly for %q", header, row)
			}
			for j, v := range row {
				cols[j].PushBack(NewStringValue(v))
			}
			if rowN < headerN { // fill in empty values
				for k := rowN; k < headerN; k++ {
					cols[k].PushBack(NewStringValue(""))
				}
			}
		}
		for _, c := range cols {
			if err := fr.AddColumn(c); err != nil {
				return nil, err
			}
		}
		return fr, nil
	}
	// use first row as header
	// assume header string at top
	header = rows[0]
	headerN = len(header)
	cols := make([]Column, headerN)
	for i := range cols {
		cols[i] = NewColumn(header[i])
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		rowN := len(row)
		if rowN > headerN {
			return nil, fmt.Errorf("header %q is not specified correctly for %q", header, row)
		}
		for j, v := range row {
			cols[j].PushBack(NewStringValue(v))
		}
		if rowN < headerN { // fill in empty values
			for k := rowN; k < headerN; k++ {
				cols[k].PushBack(NewStringValue(""))
			}
		}
	}
	for _, c := range cols {
		if err := fr.AddColumn(c); err != nil {
			return nil, err
		}
	}
	return fr, nil
}

// NewFromCSV creates a new Frame from CSV.
// Pass 'nil' header if first row is used as header strings.
// Pass 'non-nil' header if the data starts from the first row, without header strings.
func NewFromCSV(header []string, fpath string) (Frame, error) {
	f, err := openToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := csv.NewReader(f)

	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to
	// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
	// the number of fields in the first record, so that future records must
	// have the same field count. If FieldsPerRecord is negative, no check is
	// made and records may have a variable number of fields.
	rd.FieldsPerRecord = -1

	rows, err := rd.ReadAll()
	if err != nil {
		return nil, err
	}

	return NewFromRows(header, rows)
}

// NewFromColumns combines multiple columns into one data frame.
// If zero Value is not nil, it makes all columns have the same row number
// by inserting zero values where the row number is short compared to the
// one with the msot row number. The columns are deep-copied to the Frame.
func NewFromColumns(zero Value, cols ...Column) (Frame, error) {
	maxEndIndex := 0
	columns := make([]Column, len(cols))
	for i, col := range cols {
		columns[i] = col.Copy()

		if i == 0 {
			maxEndIndex = col.Count()
		}
		if maxEndIndex < col.Count() {
			maxEndIndex = col.Count()
		}
	}
	// this is index, so decrement by 1 to make it as valid index
	maxEndIndex--
	maxSize := maxEndIndex + 1

	if zero != nil {
		// make all columns have same row number
		for _, col := range columns {
			rNum := col.Count()
			if rNum < maxSize { // fill-in with zero values
				for i := 0; i < maxSize-rNum; i++ {
					col.PushBack(zero)
				}
			}
			if rNum > maxSize {
				return nil, fmt.Errorf("something wrong with minimum end index %d (%q has %d rows)", maxEndIndex, col.Header(), rNum)
			}
		}
		// double-check
		rNum := columns[0].Count()
		for _, col := range columns {
			if rNum != col.Count() {
				return nil, fmt.Errorf("%q has %d rows (expected %d rows as %q)", col.Header(), col.Count(), rNum, columns[0].Header())
			}
		}
	}

	combined := New()
	for _, col := range columns {
		if err := combined.AddColumn(col); err != nil {
			return nil, err
		}
	}
	return combined, nil
}

func (f *frame) Headers() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	rs := make([]string, len(f.headerTo))
	for k, v := range f.headerTo {
		rs[v] = k
	}
	return rs
}

func (f *frame) AddColumn(c Column) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	header := c.Header()
	if _, ok := f.headerTo[header]; ok {
		return fmt.Errorf("%q already exists", header)
	}
	f.columns = append(f.columns, c)
	f.headerTo[header] = len(f.columns) - 1
	return nil
}

func (f *frame) Column(header string) (Column, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	idx, ok := f.headerTo[header]
	if !ok {
		return nil, fmt.Errorf("%q does not exist", header)
	}
	return f.columns[idx], nil
}

func (f *frame) Columns() []Column {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.columns
}

func (f *frame) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.columns)
}

func (f *frame) UpdateHeader(origHeader, newHeader string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	idx, ok := f.headerTo[origHeader]
	if !ok {
		return fmt.Errorf("%q does not exist", origHeader)
	}
	if _, ok := f.headerTo[newHeader]; ok {
		return fmt.Errorf("%q already exists", newHeader)
	}
	f.columns[idx].UpdateHeader(newHeader)
	f.headerTo[newHeader] = idx
	delete(f.headerTo, origHeader)
	return nil
}

func (f *frame) MoveColumn(header string, target int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if target < 0 || target > len(f.headerTo) {
		return fmt.Errorf("%d is out of range", target)
	}

	oldi, ok := f.headerTo[header]
	if !ok {
		return fmt.Errorf("%q does not exist", header)
	}
	if target == oldi {
		// no need to insert
		return nil
	}

	var copied []Column
	switch {
	case target < oldi: // move somewhere to left
		// e.g. arr1, oldi 7, target 2
		// 0  1 | 2  3  4  5  6  [7]  8  9
		// 1. copy[:2]
		// 2. arr2[2] = arr1[7]
		// 3. copy[3:7]
		// 4. copy[8:]
		copied = make([]Column, target)
		if target == 0 {
			copied = []Column{}
		} else {
			copy(copied, f.columns[:target])
		}
		copied = append(copied, f.columns[oldi])
		// at this point, moved until 'target' index
		for i, c := range f.columns {
			if i < target || i == oldi { // already moved
				continue
			}
			copied = append(copied, c)
		}

	case oldi < target: // move somewhere to right
		// e.g. arr1, oldi 2, target 8
		// 0  1 [2] 3  4  5  6  7 | 8  9
		// 1. copy[:2]
		// 2. copy[3:8]
		// 3. arr2[7] = arr1[2]
		// 4. copy[8:]
		copied = make([]Column, oldi)
		if oldi == 0 {
			copied = []Column{}
		} else {
			copy(copied, f.columns[:oldi])
		}
		copied = append(copied, f.columns[oldi+1:target]...)
		for i, c := range f.columns {
			if i != oldi && i < target { // already moved
				continue
			}
			copied = append(copied, c)
		}
	}
	f.columns = copied

	// update column index
	for i, col := range f.columns {
		f.headerTo[col.Header()] = i
	}
	return nil
}

func (f *frame) DeleteColumn(header string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	idx, ok := f.headerTo[header]
	if !ok {
		return false
	}
	if idx == 0 && len(f.headerTo) == 1 {
		f.headerTo = make(map[string]int)
		f.columns = []Column{}
		return true
	}

	copy(f.columns[idx:], f.columns[idx+1:])
	f.columns = f.columns[:len(f.columns)-1 : len(f.columns)-1]

	// update headerTo
	f.headerTo = make(map[string]int)
	for i, c := range f.columns {
		f.headerTo[c.Header()] = i
	}
	return true
}

func (f *frame) Rows() ([]string, [][]string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	headers := make([]string, len(f.headerTo))
	for k, v := range f.headerTo {
		headers[v] = k
	}

	var rowN int
	for _, col := range f.columns {
		n := col.Count()
		if rowN < n {
			rowN = n
		}
	}

	rows := make([][]string, rowN)
	colN := len(f.columns)
	for rowIdx := 0; rowIdx < rowN; rowIdx++ {
		row := make([]string, colN)
		for colIdx, col := range f.columns { // rowIdx * colIdx
			v, err := col.Value(rowIdx)
			var elem string
			if err == nil {
				elem, _ = v.String()
			}
			row[colIdx] = elem
		}
		rows[rowIdx] = row
	}

	return headers, rows
}

func (f *frame) CSV(fpath string) error {
	file, err := openToOverwrite(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(file)

	headers, rows := f.Rows()
	if err := wr.Write(headers); err != nil {
		return err
	}
	if err := wr.WriteAll(rows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}

func (f *frame) CSVHorizontal(fpath string) error {
	var rows [][]string
	for _, col := range f.columns {
		row := []string{col.Header()}
		row = append(row, col.Rows()...)
		rows = append(rows, row)
	}

	file, err := openToOverwrite(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(file)
	if err := wr.WriteAll(rows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}

// Sort sorts the data frame.
// TODO: use tree?
func (f *frame) Sort(header string, st SortType, so SortOption) error {
	f.mu.Lock()
	idx, ok := f.headerTo[header]
	if !ok {
		f.mu.Unlock()
		return fmt.Errorf("%q does not exist", header)
	}
	f.mu.Unlock()

	var lesses []LessFunc
	switch st {
	case SortType_String:
		switch so {
		case SortOption_Ascending:
			lesses = []LessFunc{StringAscendingFunc(idx)}

		case SortOption_Descending:
			lesses = []LessFunc{StringDescendingFunc(idx)}
		}

	case SortType_Float64:
		switch so {
		case SortOption_Ascending:
			lesses = []LessFunc{Float64AscendingFunc(idx)}

		case SortOption_Descending:
			lesses = []LessFunc{Float64DescendingFunc(idx)}
		}

	case SortType_Duration:
		switch so {
		case SortOption_Ascending:
			lesses = []LessFunc{DurationAscendingFunc(idx)}

		case SortOption_Descending:
			lesses = []LessFunc{DurationDescendingFunc(idx)}
		}
	}

	headers, rows := f.Rows()
	SortBy(
		rows,
		lesses...,
	).Sort(rows)

	nf, err := NewFromRows(headers, rows)
	if err != nil {
		return err
	}
	v, ok := nf.(*frame)
	if !ok {
		return fmt.Errorf("cannot type assert on frame")
	}
	// *f = *v
	f.columns = v.columns
	f.headerTo = v.headerTo
	return nil
}
