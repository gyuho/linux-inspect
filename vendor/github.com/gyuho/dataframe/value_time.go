package dataframe

import (
	"fmt"
	"time"
)

// GoTime defines time data types.
type GoTime time.Time

// NewTimeValue takes any interface and returns Value.
func NewTimeValue(v interface{}) Value {
	switch t := v.(type) {
	case time.Time:
		return GoTime(t)
	default:
		panic(fmt.Errorf("%v(%T) is not supported yet", v, v))
	}
}

// NewTimeValueNil returns an empty value.
func NewTimeValueNil() Value {
	return GoTime(time.Time{})
}

func (gt GoTime) String() (string, bool) {
	return time.Time(gt).String(), true
}

func (gt GoTime) Int64() (int64, bool) {
	return 0, false
}

func (gt GoTime) Uint64() (uint64, bool) {
	return 0, false
}

func (gt GoTime) Float64() (float64, bool) {
	return 0, false
}

func (gt GoTime) Time(layout string) (time.Time, bool) {
	return time.Time(gt), true
}

func (gt GoTime) Duration() (time.Duration, bool) {
	return time.Duration(0), false
}

func (gt GoTime) IsNil() bool {
	return time.Time(gt).IsZero()
}

func (gt GoTime) EqualTo(v Value) bool {
	tv, ok := v.(GoTime)
	return ok && time.Time(gt).Equal(time.Time(tv))
}

func (gt GoTime) Copy() Value {
	return gt
}
