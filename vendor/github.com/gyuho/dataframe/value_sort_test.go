package dataframe

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestByStringAscending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%d", i)))
	}
	sort.Sort(ByStringAscending(vs))
	if !vs[0].EqualTo(NewStringValue("0")) {
		t.Fatalf("expected '0', got %v", vs[0])
	}
}

func TestByStringDescending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%d", i)))
	}
	sort.Sort(ByStringDescending(vs))
	if !vs[0].EqualTo(NewStringValue("99")) {
		t.Fatalf("expected '99', got %v", vs[0])
	}
}

func TestByStringDescendingTime(t *testing.T) {
	now := time.Now()
	vs := []Value{}
	for i := 0; i < 100; i++ {
		t := now.Add(time.Duration(i) * time.Second).String()
		vs = append(vs, NewStringValue(t))
	}
	sort.Sort(ByStringDescending(vs))
	if !vs[0].EqualTo(NewStringValue(now.Add(time.Duration(99) * time.Second).String())) {
		t.Fatalf("expected '99', got %v", vs[0])
	}
}

func TestByStringDescendingTimeColumn(t *testing.T) {
	now := time.Now()
	vs := []Value{}
	for i := 0; i < 100; i++ {
		t := now.Add(time.Duration(i) * time.Second)
		vs = append(vs, NewTimeValue(t))
	}
	sort.Sort(ByStringDescending(vs))
	if !vs[0].EqualTo(NewTimeValue(now.Add(time.Duration(99) * time.Second))) {
		t.Fatalf("expected '99', got %v", vs[0])
	}
}

func TestByFloat64Ascending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%d", i)))
	}
	sort.Sort(ByFloat64Ascending(vs))
	if !vs[0].EqualTo(NewStringValue("0")) {
		t.Fatalf("expected '0', got %v", vs[0])
	}
}

func TestByFloat64Descending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%d", i)))
	}
	vs = append(vs, NewStringValue("199.9"))
	sort.Sort(ByFloat64Descending(vs))
	if !vs[0].EqualTo(NewStringValue("199.9")) {
		t.Fatalf("expected '199.9', got %v", vs[0])
	}
}

func TestByDurationAscending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%s", time.Duration(i)*time.Second)))
	}
	sort.Sort(ByDurationAscending(vs))
	if !vs[0].EqualTo(NewStringValue("0s")) {
		t.Fatalf("expected '0s', got %v", vs[0])
	}
}

func TestByDurationDescending(t *testing.T) {
	vs := []Value{}
	for i := 0; i < 100; i++ {
		vs = append(vs, NewStringValue(fmt.Sprintf("%s", time.Duration(i)*time.Second)))
	}
	vs = append(vs, NewStringValue("200h"))
	sort.Sort(ByDurationDescending(vs))
	if !vs[0].EqualTo(NewStringValue("200h")) {
		t.Fatalf("expected '200h', got %v", vs[0])
	}
}
