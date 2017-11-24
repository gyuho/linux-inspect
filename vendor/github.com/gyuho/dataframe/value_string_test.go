package dataframe

import (
	"testing"
	"time"
)

func TestStringValue(t *testing.T) {
	v1 := NewStringValue("1")
	if v, ok := v1.Float64(); !ok {
		t.Fatalf("expected number 1, got %v", v)
	}

	v2 := NewStringValue("2.2")
	if v, ok := v2.Float64(); !ok || v != 2.2 {
		t.Fatalf("expected number 2.2, got %v", v)
	}

	v2c := v2.Copy()
	if !v2.EqualTo(v2c) {
		t.Fatalf("expected equal, got %v", v2.EqualTo(v2c))
	}

	v3t := time.Now()
	v3 := NewStringValue(v3t)
	if v, ok := v3.Time(TimeDefaultLayout); !ok || !v.Equal(v3t) {
		t.Fatalf("expected time %q, got %q", v3t, v)
	}

	v4t := time.Now().String()
	v4 := NewStringValue(v4t)
	if v, ok := v4.Time(TimeDefaultLayout); !ok || v4t != v.String() {
		t.Fatalf("expected time %q, got %q", v4t, v)
	}

	v5d := time.Hour
	v5 := NewStringValue(v5d)
	if v, ok := v5.Duration(); !ok || v != v5d {
		t.Fatalf("expected duration %s, got %v", v5d, v)
	}

	if !NewStringValue("hello").EqualTo(NewStringValue("hello")) {
		t.Fatal("EqualTo expected 'true' for 'hello' == 'hello' but got false")
	}
}

func TestStringValueInt(t *testing.T) {
	v := NewStringValue("1")
	fv, ok := v.Float64()
	if !ok || fv != 1.0 {
		t.Fatalf("expected number 1, got %f(%v)", fv, v)
	}

	iv, ok := v.Int64()
	if !ok || iv != 1 {
		t.Fatalf("expected number 1, got %d(%v)", iv, v)
	}

	uv, ok := v.Uint64()
	if !ok || uv != 1 {
		t.Fatalf("expected number 1, got %d(%v)", uv, v)
	}
}

func TestNewStringValueNil(t *testing.T) {
	v := NewStringValueNil()
	if !v.IsNil() {
		t.Fatalf("expected nil, got %v", v)
	}
}
