package dataframe

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestColumn(t *testing.T) {
	c := NewColumn("second")
	if c.Header() != "second" {
		t.Fatalf("expected 'second', got %v", c.Header())
	}
	for i := 0; i < 100; i++ {
		d := c.PushBack(NewStringValue(i))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	if c.Count() != 100 {
		t.Fatalf("expected '100', got %v", c.Count())
	}

	if err := c.Set(10, NewStringValue(10000)); err != nil {
		t.Fatal(err)
	}
	if v, err := c.Value(10); err != nil || !v.EqualTo(NewStringValue(10000)) {
		t.Fatalf("expected '10', got %v(%v)", v, err)
	}

	if err := c.Set(10, NewStringValue(10)); err != nil {
		t.Fatal(err)
	}
	if v, err := c.Value(10); err != nil || !v.EqualTo(NewStringValue(10)) {
		t.Fatalf("expected '10', got %v(%v)", v, err)
	}
	idx, ok := c.FindFirst(NewStringValue(10))
	if !ok || idx != 10 {
		t.Fatalf("expected 10, got %d", idx)
	}
	bv, ok := c.Back()
	if !ok || !bv.EqualTo(NewStringValue(99)) {
		t.Fatalf("expected '99', got %v", bv)
	}
	bv, ok = c.PopBack()
	if !ok || !bv.EqualTo(NewStringValue(99)) {
		t.Fatalf("expected '99', got %v", bv)
	}
	fv, ok := c.Front()
	if !ok || !fv.EqualTo(NewStringValue(0)) {
		t.Fatalf("expected '0', got %v", fv)
	}
	fv, ok = c.PopFront()
	if !ok || !fv.EqualTo(NewStringValue(0)) {
		t.Fatalf("expected '0', got %v", fv)
	}
	dv, err := c.Delete(1)
	if err != nil || !dv.EqualTo(NewStringValue(2)) {
		t.Fatalf("expected '2', got %v(%v)", dv, err)
	}
	fidx, ok := c.FindFirst(NewStringValue(2))
	if fidx != -1 || ok {
		t.Fatalf("expected -1, false, got %v %v", fidx, ok)
	}

	if pv := c.PushFront(NewStringValue("A")); pv != 98 {
		t.Fatalf("expected '98', got %v", pv)
	}
	if vv, ok := c.PopFront(); !ok || !vv.EqualTo(NewStringValue("A")) {
		t.Fatalf("expected 'A', got %v", vv)
	}
}

func TestColumnTyped(t *testing.T) {
	col := NewColumnTyped("TIME", TIME)
	size, err := col.PushFrontTyped(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if size != 1 {
		t.Fatalf("size expected 1, got %d", size)
	}
	_, err = col.PushFrontTyped("a")
	fmt.Println(err)
	if err == nil {
		t.Fatalf("error expected, gut got error %v", err)
	}
}

func TestColumnRow(t *testing.T) {
	c := NewColumn("A")
	for i := 0; i < 3; i++ {
		d := c.PushBack(NewStringValue(i))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	if pv := c.PushFront(NewStringValue("Google")); pv != 4 {
		t.Fatalf("expected '4', got %v", pv)
	}
	if pv := c.PushBack(NewStringValue("Amazon")); pv != 5 {
		t.Fatalf("expected '5', got %v", pv)
	}
	if n := c.Count(); n != 5 {
		t.Fatalf("c.Count() expected 5, got %d", n)
	}
	expected := []string{"Google", "0", "1", "2", "Amazon"}
	rows := c.Rows()
	if !reflect.DeepEqual(expected, rows) {
		t.Fatalf("rows expected %+v, got %+v", expected, rows)
	}
}

func TestColumnRowInt64s(t *testing.T) {
	c := NewColumn("A")
	expected := []int64{100, 50, 7, 10, 5}
	for i, v := range expected {
		if pv := c.PushBack(NewStringValue(v)); pv != i+1 {
			t.Fatalf("expected '%d', got %v", i+1, pv)
		}
	}
	if n := c.Count(); n != len(expected) {
		t.Fatalf("c.Count() expected %d, got %d", len(expected), n)
	}
	rows, ok := c.Int64s()
	if !ok {
		t.Fatalf("ok expected true, got %v", ok)
	}
	if !reflect.DeepEqual(expected, rows) {
		t.Fatalf("rows expected %+v, got %+v", expected, rows)
	}
}

func TestColumnRowFloat64s(t *testing.T) {
	c := NewColumn("A")
	expected := []float64{10.01, 5.05, 1.7, 8, 111101, 5.5, 7, 10}
	for i, v := range expected {
		if pv := c.PushBack(NewStringValue(v)); pv != i+1 {
			t.Fatalf("expected '%d', got %v", i+1, pv)
		}
	}
	if n := c.Count(); n != len(expected) {
		t.Fatalf("c.Count() expected %d, got %d", len(expected), n)
	}
	rows, ok := c.Float64s()
	if !ok {
		t.Fatalf("ok expected true, got %v", ok)
	}
	if !reflect.DeepEqual(expected, rows) {
		t.Fatalf("rows expected %+v, got %+v", expected, rows)
	}
}

func TestColumnNonNil(t *testing.T) {
	c := NewColumn("second")
	if c.Header() != "second" {
		t.Fatalf("expected 'second', got %v", c.Header())
	}
	for i := 0; i < 100; i++ {
		d := c.PushBack(NewStringValue(i))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	c.PushFront(NewStringValue(""))
	c.PushBack(NewStringValue(""))

	fv, ok := c.Front()
	if !ok || !fv.EqualTo(NewStringValue("")) {
		t.Fatalf("expected '0', got %v", fv)
	}
	fv, ok = c.FrontNonNil()
	if !ok || !fv.EqualTo(NewStringValue(0)) {
		t.Fatalf("expected '0', got %v", fv)
	}
	bv, ok := c.Back()
	if !ok || !bv.EqualTo(NewStringValue("")) {
		t.Fatalf("expected '99', got %v", bv)
	}
	bv, ok = c.BackNonNil()
	if !ok || !bv.EqualTo(NewStringValue(99)) {
		t.Fatalf("expected '99', got %v", bv)
	}
}

func TestColumnAppends(t *testing.T) {
	c := NewColumn("second")
	c.PushBack(NewStringValue(1))
	c.PushBack(NewStringValue(2))
	if err := c.Appends(NewStringValue(1000), 1000); err != nil {
		t.Fatal(err)
	}
	s := c.Count()
	if s != 1000 {
		t.Fatalf("expected '1000', got %v", s)
	}
	fv, ok := c.Front()
	if !ok || !fv.EqualTo(NewStringValue(1)) {
		t.Fatalf("expected '1', got %v", fv)
	}
	fv, ok = c.FrontNonNil()
	if !ok || !fv.EqualTo(NewStringValue(1)) {
		t.Fatalf("expected '1', got %v", fv)
	}
	bv, ok := c.Back()
	if !ok || !bv.EqualTo(NewStringValue(1000)) {
		t.Fatalf("expected '1000', got %v", bv)
	}
	bv, ok = c.BackNonNil()
	if !ok || !bv.EqualTo(NewStringValue(1000)) {
		t.Fatalf("expected '1000', got %v", bv)
	}
}

func TestColumnAppendsNil(t *testing.T) {
	c := NewColumn("second")
	c.PushBack(NewStringValue(1))
	c.PushBack(NewStringValue(2))
	if err := c.Appends(NewStringValue(""), 1000); err != nil {
		t.Fatal(err)
	}
	s := c.Count()
	if s != 1000 {
		t.Fatalf("expected '1000', got %v", s)
	}
	fv, ok := c.Front()
	if !ok || !fv.EqualTo(NewStringValue(1)) {
		t.Fatalf("expected '1', got %v", fv)
	}
	bv, ok := c.Back()
	if !ok || !bv.EqualTo(NewStringValue("")) {
		t.Fatalf("expected '', got %v", bv)
	}
}

func TestColumnDeletes(t *testing.T) {
	c := NewColumn("second")
	for i := 0; i < 100; i++ {
		d := c.PushBack(NewStringValue(fmt.Sprintf("%d", i)))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	idx, ok := c.FindFirst(NewStringValue(60))
	if idx != 60 || !ok {
		t.Fatalf("expected 60, true, got %d %v", idx, ok)
	}
	if err := c.Deletes(50, 70); err != nil {
		t.Fatal(err)
	}
	idx, ok = c.FindFirst(NewStringValue(70))
	if idx != 50 || !ok {
		t.Fatalf("expected 50, true, got %d %v", idx, ok)
	}
	if c.Count() != 80 {
		t.Fatalf("expected 80, got %d", c.Count())
	}
	idx, ok = c.FindFirst(NewStringValue(60))
	if idx != -1 || ok {
		t.Fatalf("expected -1, false, got %d %v", idx, ok)
	}
}

func TestColumnKeep(t *testing.T) {
	c := NewColumn("second")
	for i := 0; i < 100; i++ {
		d := c.PushBack(NewStringValue(i))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	if err := c.Keep(50, 70); err != nil {
		t.Fatal(err)
	}
	idx, ok := c.FindFirst(NewStringValue(50))
	if idx != 0 || !ok {
		t.Fatalf("expected 0, true, got %d %v", idx, ok)
	}
	idx, ok = c.FindFirst(NewStringValue(69))
	if idx != 19 || !ok {
		t.Fatalf("expected 19, true, got %d %v", idx, ok)
	}
	idx, ok = c.FindFirst(NewStringValue(70))
	if idx != -1 || ok {
		t.Fatalf("expected -1, false, got %d %v", idx, ok)
	}
	if c.Count() != 20 {
		t.Fatalf("expected 20, got %d", c.Count())
	}
	idx, ok = c.FindFirst(NewStringValue(90))
	if idx != -1 || ok {
		t.Fatalf("expected -1, false, got %d %v", idx, ok)
	}
}

func TestColumnSortByStringAscending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(fmt.Sprintf("%d", i)))
	}
	c.SortByStringAscending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue(0)) {
		t.Fatalf("expected '0', got %v", fv)
	}
}

func TestColumnSortByStringDescending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(i))
	}
	c.SortByStringDescending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue(99)) {
		t.Fatalf("expected '99', got %v", fv)
	}
}

func TestColumnSortByFloat64Ascending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(i))
	}
	c.SortByFloat64Ascending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue(0)) {
		t.Fatalf("expected '0', got %v", fv)
	}
}

func TestColumnSortByFloat64Descending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(i))
	}
	c.PushBack(NewStringValue("199.9"))
	c.SortByFloat64Descending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue(199.9)) {
		t.Fatalf("expected '199.9', got %v", fv)
	}
}

func TestColumnSortByDurationAscending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(time.Duration(i) * time.Second))
	}
	c.SortByDurationAscending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue("0s")) {
		t.Fatalf("expected '0s', got %v", fv)
	}
}

func TestColumnSortByDurationDescending(t *testing.T) {
	c := NewColumn("column")
	for i := 0; i < 100; i++ {
		c.PushBack(NewStringValue(time.Duration(i) * time.Second))
	}
	c.PushBack(NewStringValue(200 * time.Hour))
	c.SortByDurationDescending()
	fv, err := c.Value(0)
	if err != nil {
		t.Fatal(err)
	}
	if !fv.EqualTo(NewStringValue(200 * time.Hour)) {
		t.Fatalf("expected '200h', got %v", fv)
	}
}
