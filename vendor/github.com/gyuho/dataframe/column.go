package dataframe

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Column represents column-based data.
type Column interface {
	// Count returns the number of rows of the Column.
	Count() int

	// Header returns the header of the Column.
	Header() string

	// Rows returns all the data in string slice.
	Rows() []string

	// Uint64s returns all the data in int64 slice.
	Uint64s() ([]uint64, bool)

	// Int64s returns all the data in int64 slice.
	Int64s() ([]int64, bool)

	// Float64s returns all the data in float64 slice.
	Float64s() ([]float64, bool)

	// Times returns all the data in time.Time slice.
	Times(layout string) ([]time.Time, bool)

	// UpdateHeader updates the header of the Column.
	UpdateHeader(header string)

	// Value returns the Value in the row. It returns error if the row
	// is out of index range.
	Value(row int) (Value, error)

	// Set overwrites the value
	Set(row int, v Value) error

	// FindFirst finds the first Value, and returns the row number.
	// It returns -1 and false if the value does not exist.
	FindFirst(v Value) (int, bool)

	// FindLast finds the last Value, and returns the row number.
	// It returns -1 and false if the value does not exist.
	FindLast(v Value) (int, bool)

	// Front returns the first row Value.
	Front() (Value, bool)

	// FrontNonNil returns the first non-nil Value from the first row.
	FrontNonNil() (Value, bool)

	// Back returns the last row Value.
	Back() (Value, bool)

	// BackNonNil returns the first non-nil Value from the last row.
	BackNonNil() (Value, bool)

	// PushFront adds a Value to the front of the Column.
	// This does not prevent inserting wrong data types.
	// Assumes all data are string.
	PushFront(v Value) int

	// PushFrontTyped adds a Value to the front of the Column.
	// It returns error if the value doesn't match the type of the column.
	PushFrontTyped(v interface{}) (int, error)

	// PushBack appends the Value to the Column.
	// This does not prevent inserting wrong data types.
	// Assumes all data are string.
	PushBack(v Value) int

	// PushBackTyped appends the Value to the Column.
	// It returns error if the value doesn't match the type of the column.
	PushBackTyped(v interface{}) (int, error)

	// Delete deletes a row by index.
	Delete(row int) (Value, error)

	// Deletes deletes rows by index [start, end).
	Deletes(start, end int) error

	// Keep keeps the rows by index [start, end).
	Keep(start, end int) error

	// PopFront deletes the value at front.
	PopFront() (Value, bool)

	// PopBack deletes the last value.
	PopBack() (Value, bool)

	// Appends adds the Value to the Column until it reaches the target size.
	Appends(v Value, targetSize int) error

	// Copy deep-copies a column.
	Copy() Column

	// SortByStringAscending sorts Column in string ascending order.
	SortByStringAscending()

	// SortByStringDescending sorts Column in string descending order.
	SortByStringDescending()

	// SortByFloat64Ascending sorts Column in number(float) ascending order.
	SortByFloat64Ascending()

	// SortByFloat64Descending sorts Column in number(float) descending order.
	SortByFloat64Descending()

	// SortByDurationAscending sorts Column in time.Duration ascending order.
	SortByDurationAscending()

	// SortByDurationDescending sorts Column in time.Duration descending order.
	SortByDurationDescending()
}

type column struct {
	mu       sync.Mutex
	dataType DATA_TYPE
	header   string
	size     int
	data     []Value
}

// NewColumn creates a new Column.
func NewColumn(hd string) Column {
	return &column{
		dataType: STRING,
		header:   hd,
		size:     0,
		data:     []Value{},
	}
}

// NewColumnTyped creates a new Column with data type.
func NewColumnTyped(hd string, tp DATA_TYPE) Column {
	return &column{
		dataType: tp,
		header:   hd,
		size:     0,
		data:     []Value{},
	}
}

func (c *column) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.size
}

func (c *column) Header() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.header
}

func (c *column) Rows() (rows []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rows = make([]string, len(c.data))
	for i := range c.data {
		v, _ := c.data[i].String()
		rows[i] = v
	}
	return
}

func (c *column) Uint64s() (rows []uint64, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rows = make([]uint64, len(c.data))
	for i := range c.data {
		var v uint64
		v, ok = c.data[i].Uint64()
		if !ok {
			break
		}
		rows[i] = v
	}
	return
}

func (c *column) Int64s() (rows []int64, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rows = make([]int64, len(c.data))
	for i := range c.data {
		var v int64
		v, ok = c.data[i].Int64()
		if !ok {
			break
		}
		rows[i] = v
	}
	return
}

func (c *column) Float64s() (rows []float64, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rows = make([]float64, len(c.data))
	for i := range c.data {
		var v float64
		v, ok = c.data[i].Float64()
		if !ok {
			break
		}
		rows[i] = v
	}
	return
}

func (c *column) Times(layout string) (rows []time.Time, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rows = make([]time.Time, len(c.data))
	for i := range c.data {
		var v time.Time
		v, ok = c.data[i].Time(layout)
		if !ok {
			break
		}
		rows[i] = v
	}
	return
}

func (c *column) UpdateHeader(header string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.header = header
}

func (c *column) Value(row int) (Value, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if row > c.size-1 {
		return nil, fmt.Errorf("index out of range (got %d for size %d)", row, c.size)
	}
	return c.data[row], nil
}

func (c *column) Set(row int, v Value) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if row > c.size-1 {
		return fmt.Errorf("index out of range (got %d for size %d)", row, c.size)
	}
	c.data[row] = v
	return nil
}

func (c *column) FindFirst(v Value) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.data {
		if c.data[i].EqualTo(v) {
			return i, true
		}
	}
	return -1, false
}

func (c *column) FindLast(v Value) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var idx int
	for i := range c.data {
		if c.data[i].EqualTo(v) {
			idx = i
		}
	}
	if idx != 0 {
		return idx, true
	}
	return -1, false
}

func (c *column) Front() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	v := c.data[0]
	return v, true
}

func (c *column) FrontNonNil() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	for _, v := range c.data {
		if !v.IsNil() {
			return v, true
		}
	}
	return nil, false
}

func (c *column) Back() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	v := c.data[c.size-1]
	return v, true
}

func (c *column) BackNonNil() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	for i := c.size - 1; i > 0; i-- {
		v := c.data[i]
		if !v.IsNil() {
			return v, true
		}
	}
	return nil, false
}

func (c *column) PushFront(v Value) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	temp := make([]Value, c.size+1)
	temp[0] = v
	copy(temp[1:], c.data)
	c.data = temp
	c.size++
	return c.size
}

func (c *column) PushFrontTyped(v interface{}) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var value Value
	switch expected := c.dataType; expected {
	case STRING:
		value = NewStringValue(v)
	default:
		t := ReflectTypeOf(v)
		if expected != t { // column is typed
			return -1, fmt.Errorf("column %q expected data type %q, got %q", c.header, expected, t)
		}
		value = ToValue(v)
	}

	temp := make([]Value, c.size+1)
	temp[0] = value
	copy(temp[1:], c.data)
	c.data = temp
	c.size++
	return c.size, nil
}

func (c *column) PushBack(v Value) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = append(c.data, v)
	c.size++
	return c.size
}

func (c *column) PushBackTyped(v interface{}) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var value Value
	switch expected := c.dataType; expected {
	case STRING:
		value = NewStringValue(v)
	default:
		t := ReflectTypeOf(v)
		if expected != t { // column is typed
			return -1, fmt.Errorf("column %q expected data type %q, got %q", c.header, expected, t)
		}
		value = ToValue(v)
	}

	c.data = append(c.data, value)
	c.size++
	return c.size, nil
}

func (c *column) Delete(row int) (Value, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if row > c.size-1 {
		return nil, fmt.Errorf("index out of range (got %d for size %d)", row, c.size)
	}
	v := c.data[row]
	copy(c.data[row:], c.data[row+1:])
	c.data = c.data[:len(c.data)-1 : len(c.data)-1]
	c.size--
	return v, nil
}

func (c *column) Deletes(start, end int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if start < 0 || end < 0 || start > end {
		return fmt.Errorf("wrong range %d %d", start, end)
	}
	if start > c.size {
		return fmt.Errorf("index out of range (start %d, size %d)", start, c.size)
	}
	if end > c.size {
		return fmt.Errorf("index out of range (end %d, size %d)", end, c.size)
	}
	if start == end {
		return nil
	}

	delta := end - start
	c.size = c.size - delta
	var nds []Value
	for i := range c.data {
		if i >= start && i < end {
			continue
		}
		nds = append(nds, c.data[i])
	}
	c.data = nds
	return nil
}

func (c *column) Keep(start, end int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if start < 0 || end < 0 || start > end {
		return fmt.Errorf("wrong range %d %d", start, end)
	}
	if start > c.size {
		return fmt.Errorf("index out of range (start %d, size %d)", start, c.size)
	}
	if end > c.size {
		return fmt.Errorf("index out of range (end %d, size %d)", end, c.size)
	}
	if start == end {
		return nil
	}

	delta := end - start
	c.size = delta
	var nds []Value
	for _, v := range c.data[start:end] {
		nds = append(nds, v)
	}
	c.data = nds
	return nil
}

func (c *column) PopFront() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	v := c.data[0]
	c.data = c.data[1:len(c.data):len(c.data)]
	c.size--
	return v, true
}

func (c *column) PopBack() (Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size == 0 {
		return nil, false
	}
	v := c.data[c.size-1]
	c.data = c.data[:len(c.data)-1 : len(c.data)-1]
	c.size--
	return v, true
}

func (c *column) Appends(v Value, targetSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size > 0 && c.size > targetSize {
		return fmt.Errorf("cannot append with target size %d, which is less than the column size %d (can't overwrite)", targetSize, c.size)
	}

	for i := c.size; i < targetSize; i++ {
		c.data = append(c.data, v)
		c.size++
	}
	return nil
}

func (c *column) Copy() Column {
	c2 := &column{
		header: c.header,
		size:   c.size,
		data:   make([]Value, len(c.data)),
	}
	for i := range c.data {
		c2.data[i] = c.data[i].Copy()
	}
	return c2
}

func (c *column) SortByStringAscending()    { sort.Sort(ByStringAscending(c.data)) }
func (c *column) SortByStringDescending()   { sort.Sort(ByStringDescending(c.data)) }
func (c *column) SortByFloat64Ascending()   { sort.Sort(ByFloat64Ascending(c.data)) }
func (c *column) SortByFloat64Descending()  { sort.Sort(ByFloat64Descending(c.data)) }
func (c *column) SortByDurationAscending()  { sort.Sort(ByDurationAscending(c.data)) }
func (c *column) SortByDurationDescending() { sort.Sort(ByDurationDescending(c.data)) }
