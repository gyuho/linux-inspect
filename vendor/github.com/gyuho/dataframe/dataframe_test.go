package dataframe

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFrame(t *testing.T) {
	c1 := NewColumn("second1")
	for i := 0; i < 100; i++ {
		d := c1.PushBack(NewStringValue(fmt.Sprintf("%d", i)))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}
	c1.UpdateHeader("aaa")
	if c1.Header() != "aaa" {
		t.Fatalf("expected 'aaa', got %v", c1.Header())
	}
	c1.UpdateHeader("second1")
	if c1.Header() != "second1" {
		t.Fatalf("expected 'second1', got %v", c1.Header())
	}

	c2 := NewColumn("second2")
	for i := 0; i < 100; i++ {
		d := c2.PushBack(NewStringValue(fmt.Sprintf("%d", i)))
		if i+1 != d {
			t.Fatalf("expected %d, got %d", i+1, d)
		}
	}

	fr := New()
	if err := fr.AddColumn(c1); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(c1); err == nil {
		t.Fatal("expected error")
	}
	if err := fr.AddColumn(c2); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(c2); err == nil {
		t.Fatal("expected error")
	}

	if c, err := fr.Column("second1"); c == nil || err != nil {
		t.Fatal(err)
	}
	if c, err := fr.Column("second2"); c == nil || err != nil {
		t.Fatal(err)
	}

	if err := fr.UpdateHeader("second2", "aaa"); err != nil {
		t.Fatal(err)
	}
	if c, err := fr.Column("aaa"); c == nil || err != nil {
		t.Fatal(err)
	}
	if hs := fr.Headers(); !reflect.DeepEqual(hs, []string{"second1", "aaa"}) {
		t.Fatalf("expected equal, got %q != %q", hs, []string{"second1", "aaa"})
	}

	if err := fr.UpdateHeader("aaa", "second2"); err != nil {
		t.Fatal(err)
	}
	if hs := fr.Headers(); !reflect.DeepEqual(hs, []string{"second1", "second2"}) {
		t.Fatalf("expected equal, got %q != %q", hs, []string{"second1", "second2"})
	}
	if ok := fr.DeleteColumn("second1"); !ok {
		t.Fatalf("expected 'true', got %v", ok)
	}
	if cd := fr.Count(); cd != 1 {
		t.Fatalf("expected 1, got %v", cd)
	}
	if ok := fr.DeleteColumn("second1"); ok {
		t.Fatalf("expected 'false', got %v", ok)
	}
	if c, err := fr.Column("second1"); c != nil || err == nil {
		t.Fatalf("expected <nil, 'second1 does not exist'>, but <%v, %v>", c, err)
	}
	if ok := fr.DeleteColumn("second2"); !ok {
		t.Fatalf("expected 'true', got %v", ok)
	}
	if cd := fr.Count(); cd != 0 {
		t.Fatalf("expected 0, got %v", cd)
	}
}

func TestFrameCSVHorizontal(t *testing.T) {
	c1 := NewColumn("A")
	c1.PushBack(NewStringValue(1))
	c2 := NewColumn("B")
	c2.PushBack(NewStringValue(1))
	c2.PushBack(NewStringValue(2))
	c3 := NewColumn("C")
	c3.PushBack(NewStringValue(1))
	c3.PushBack(NewStringValue(2))
	c3.PushBack(NewStringValue(3))
	fr := New()
	if err := fr.AddColumn(c1); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(c2); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(c3); err != nil {
		t.Fatal(err)
	}

	fpath := filepath.Join(os.TempDir(), "test-horizontal-psn-frame.csv")
	defer os.RemoveAll(fpath)

	if err := fr.CSVHorizontal(fpath); err != nil {
		t.Fatal(err)
	}

	bts, err := ioutil.ReadFile(fpath)
	if err != nil {
		t.Fatal(err)
	}
	expected := `A,1
B,1,2
C,1,2,3
`
	if string(bts) != expected {
		t.Fatalf("csv expected %q, got %q", expected, bts)
	}
}

func TestNewFromCSV(t *testing.T) {
	if _, err := NewFromCSV([]string{"second"}, "testdata/bench-01-all-aggregated.csv"); err == nil {
		t.Fatal("expected error, got nil")
	}
	fr, err := NewFromCSV(nil, "testdata/bench-01-all-aggregated.csv")
	if err != nil {
		t.Fatal(err)
	}
	cols := []string{"second", "avg_latency_ms_consul", "throughput_consul", "cumulative_throughput_consul", "avg_cpu_consul", "avg_memory_mb_consul", "avg_latency_ms_etcd3", "throughput_etcd3", "cumulative_throughput_etcd3", "avg_cpu_etcd3", "avg_memory_mb_etcd3", "avg_latency_ms_etcd2", "throughput_etcd2", "cumulative_throughput_etcd2", "avg_cpu_etcd2", "avg_memory_mb_etcd2", "avg_latency_ms_zk", "throughput_zk", "cumulative_throughput_zk", "avg_cpu_zk", "avg_memory_mb_zk"}
	if !reflect.DeepEqual(fr.Headers(), cols) {
		t.Fatalf("expected %q, got %q", cols, fr.Headers())
	}
	ac, err := fr.Column("avg_latency_ms_etcd3")
	if err != nil {
		t.Fatal(err)
	}
	if v, err := ac.Value(229); !v.IsNil() || err != nil {
		t.Fatalf("expected <nil, nil>, got <%v, %v>", v.IsNil(), err)
	}
	if v, err := ac.Value(0); v.IsNil() || !v.EqualTo(NewStringValue("4.484004")) || err != nil {
		t.Fatalf("expected <nil, nil>, got <%v, %v>", v, err)
	}
	ac2, err := fr.Column("avg_latency_ms_etcd2")
	if err != nil {
		t.Fatal(err)
	}
	if ac.Count() != ac2.Count() {
		t.Fatalf("expected equal %v != %v", ac.Count(), ac2.Count())
	}

	fpath := "test.csv"
	if err := fr.CSV(fpath); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fpath)

	{
		if _, err := NewFromCSV([]string{"second"}, fpath); err == nil {
			t.Fatal("expected error, got nil")
		}
		fr, err := NewFromCSV(nil, fpath)
		if err != nil {
			t.Fatal(err)
		}
		cols := []string{"second", "avg_latency_ms_consul", "throughput_consul", "cumulative_throughput_consul", "avg_cpu_consul", "avg_memory_mb_consul", "avg_latency_ms_etcd3", "throughput_etcd3", "cumulative_throughput_etcd3", "avg_cpu_etcd3", "avg_memory_mb_etcd3", "avg_latency_ms_etcd2", "throughput_etcd2", "cumulative_throughput_etcd2", "avg_cpu_etcd2", "avg_memory_mb_etcd2", "avg_latency_ms_zk", "throughput_zk", "cumulative_throughput_zk", "avg_cpu_zk", "avg_memory_mb_zk"}
		if !reflect.DeepEqual(fr.Headers(), cols) {
			t.Fatalf("expected %q, got %q", cols, fr.Headers())
		}
		ac, err := fr.Column("avg_latency_ms_etcd3")
		if err != nil {
			t.Fatal(err)
		}
		if v, err := ac.Value(229); !v.IsNil() || err != nil {
			t.Fatalf("expected <nil, nil>, got <%v, %v>", v.IsNil(), err)
		}
		if v, err := ac.Value(0); v.IsNil() || !v.EqualTo(NewStringValue("4.484004")) || err != nil {
			t.Fatalf("expected <nil, nil>, got <%v, %v>", v, err)
		}
		ac2, err := fr.Column("avg_latency_ms_etcd2")
		if err != nil {
			t.Fatal(err)
		}
		if ac.Count() != ac2.Count() {
			t.Fatalf("expected equal %v != %v", ac.Count(), ac2.Count())
		}
	}
}

func TestNewFromRows(t *testing.T) {
	fr, err := NewFromCSV(nil, "testdata/bench-01-all-aggregated.csv")
	if err != nil {
		t.Fatal(err)
	}
	header, rows := fr.Rows()
	fr2, err := NewFromRows(header, rows)
	if err != nil {
		t.Fatal(err)
	}
	fpath := "test.csv"
	if err := fr2.CSV(fpath); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fpath)

	fr, err = NewFromCSV(nil, fpath)
	if err != nil {
		t.Fatal(err)
	}
	cols := []string{"second", "avg_latency_ms_consul", "throughput_consul", "cumulative_throughput_consul", "avg_cpu_consul", "avg_memory_mb_consul", "avg_latency_ms_etcd3", "throughput_etcd3", "cumulative_throughput_etcd3", "avg_cpu_etcd3", "avg_memory_mb_etcd3", "avg_latency_ms_etcd2", "throughput_etcd2", "cumulative_throughput_etcd2", "avg_cpu_etcd2", "avg_memory_mb_etcd2", "avg_latency_ms_zk", "throughput_zk", "cumulative_throughput_zk", "avg_cpu_zk", "avg_memory_mb_zk"}
	if !reflect.DeepEqual(fr.Headers(), cols) {
		t.Fatalf("expected %q, got %q", cols, fr.Headers())
	}
	ac, err := fr.Column("avg_latency_ms_etcd3")
	if err != nil {
		t.Fatal(err)
	}
	if v, err := ac.Value(229); !v.IsNil() || err != nil {
		t.Fatalf("expected <nil, nil>, got <%v, %v>", v.IsNil(), err)
	}
	if v, err := ac.Value(0); v.IsNil() || !v.EqualTo(NewStringValue("4.484004")) || err != nil {
		t.Fatalf("expected <nil, nil>, got <%v, %v>", v, err)
	}
	ac2, err := fr.Column("avg_latency_ms_etcd2")
	if err != nil {
		t.Fatal(err)
	}
	if ac.Count() != ac2.Count() {
		t.Fatalf("expected equal %v != %v", ac.Count(), ac2.Count())
	}
}

func TestNewFromColumns(t *testing.T) {
	colA := NewColumn("A")
	colA.PushBack(NewStringValue("1"))
	colB := NewColumn("B")
	colB.PushBack(NewStringValue("1"))
	colB.PushBack(NewStringValue("2"))
	colC := NewColumn("C")
	colC.PushBack(NewStringValue("1"))
	colC.PushBack(NewStringValue("2"))
	colC.PushBack(NewStringValue("3"))

	fr, err := NewFromColumns(NewStringValue(0), colA, colB, colC)
	if err != nil {
		t.Fatal(err)
	}
	header, rows := fr.Rows()
	if !reflect.DeepEqual(header, []string{"A", "B", "C"}) {
		t.Fatalf("expected %v, got %v", []string{"A", "B", "C"}, header)
	}
	expected := [][]string{
		{"1", "1", "1"},
		{"0", "2", "2"},
		{"0", "0", "3"},
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected %v, got %v", expected, rows)
	}
}

func TestNewFromColumnsDifferentRowNumber(t *testing.T) {
	colA := NewColumn("A")
	colA.PushBack(NewStringValue("1"))
	colB := NewColumn("B")
	colB.PushBack(NewStringValue("1"))
	colB.PushBack(NewStringValue("2"))
	colC := NewColumn("C")
	colC.PushBack(NewStringValue("1"))
	colC.PushBack(NewStringValue("2"))
	colC.PushBack(NewStringValue("3"))

	fr, err := NewFromColumns(nil, colA, colB, colC)
	if err != nil {
		t.Fatal(err)
	}
	header, rows := fr.Rows()
	if !reflect.DeepEqual(header, []string{"A", "B", "C"}) {
		t.Fatalf("expected %v, got %v", []string{"A", "B", "C"}, header)
	}
	expected := [][]string{
		{"1", "1", "1"},
		{"", "2", "2"},
		{"", "", "3"},
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected %v, got %v", expected, rows)
	}
}

func TestDataFrameFindFirst(t *testing.T) {
	fr, err := NewFromCSV(nil, "testdata/bench-01-etcd-1-monitor.csv")
	if err != nil {
		t.Fatal(err)
	}
	col, err := fr.Column("unix_ts")
	if err != nil {
		t.Fatal(err)
	}
	if col.Count() != 362 {
		t.Fatalf("expected 362, got %d", col.Count())
	}
	minTS := "1458758226"
	idx, ok := col.FindFirst(NewStringValue(minTS))
	if idx != 361 || !ok {
		t.Fatalf("expected 361, true, got %d, %v", idx, ok)
	}
	v, err := col.Value(idx)
	if err != nil {
		t.Fatal(err)
	}
	if !v.EqualTo(NewStringValue(minTS)) {
		t.Fatalf("unexpected: %v != %v", v, minTS)
	}

	if err := col.Deletes(0, 2); err != nil {
		t.Fatal(err)
	}
	if col.Count() != 360 {
		t.Fatalf("expected 360, got %d", col.Count())
	}
	{
		idx, ok := col.FindFirst(NewStringValue(minTS))
		if idx != 359 || !ok {
			t.Fatalf("expected 359, true, got %d, %v", idx, ok)
		}
		v, err := col.Value(idx)
		if err != nil {
			t.Fatal(err)
		}
		if !v.EqualTo(NewStringValue(minTS)) {
			t.Fatalf("unexpected: %v != %v", v, minTS)
		}
	}
}

func TestMoveColumn(t *testing.T) {
	fr := New()
	if err := fr.AddColumn(NewColumn("0")); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(NewColumn("1")); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(NewColumn("2")); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(NewColumn("3")); err != nil {
		t.Fatal(err)
	}

	if err := fr.MoveColumn("1", 0); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"1", "0", "2", "3"}, fr.Headers()) {
		t.Fatalf("header expected %+v, got %+v", []string{"1", "0", "2", "3"}, fr.Headers())
	}

	if err := fr.MoveColumn("1", 3); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"0", "2", "1", "3"}, fr.Headers()) {
		t.Fatalf("header expected %+v, got %+v", []string{"0", "2", "1", "3"}, fr.Headers())
	}

	if err := fr.MoveColumn("1", 1); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"0", "1", "2", "3"}, fr.Headers()) {
		t.Fatalf("header expected %+v, got %+v", []string{"0", "1", "2", "3"}, fr.Headers())
	}

	if err := fr.MoveColumn("3", 1); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"0", "3", "1", "2"}, fr.Headers()) {
		t.Fatalf("header expected %+v, got %+v", []string{"0", "3", "1", "2"}, fr.Headers())
	}

	if err := fr.AddColumn(NewColumn("A")); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(NewColumn("B")); err != nil {
		t.Fatal(err)
	}
	if err := fr.AddColumn(NewColumn("C")); err != nil {
		t.Fatal(err)
	}
	if err := fr.MoveColumn("A", 1); err != nil {
		t.Fatal(err)
	}
	if err := fr.MoveColumn("B", 1); err != nil {
		t.Fatal(err)
	}
	if err := fr.MoveColumn("C", 1); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"0", "C", "B", "A", "3", "1", "2"}, fr.Headers()) {
		t.Fatalf("header expected %+v, got %+v", []string{"0", "C", "B", "A", "3", "1", "2"}, fr.Headers())
	}
}

func TestSort(t *testing.T) {
	fr, err := NewFromCSV(nil, "testdata/bench-01-all-aggregated.csv")
	if err != nil {
		t.Fatal(err)
	}
	if err := fr.Sort("second", SortType_Float64, SortOption_Descending); err != nil {
		t.Fatal(err)
	}
	fpath := "test.csv"
	if err := fr.CSV(fpath); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fpath)
}
