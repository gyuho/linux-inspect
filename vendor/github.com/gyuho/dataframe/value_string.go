package dataframe

import (
	"fmt"
	"strconv"
	"time"
)

// String defines string data types.
type String string

// NewStringValue takes any interface and returns Value.
func NewStringValue(v interface{}) Value {
	switch t := v.(type) {
	case string:
		return String(t)
	case []byte:
		return String(t)
	case bool:
		return String(fmt.Sprintf("%v", t))
	case int:
		return String(strconv.FormatInt(int64(t), 10))
	case int8:
		return String(strconv.FormatInt(int64(t), 10))
	case int16:
		return String(strconv.FormatInt(int64(t), 10))
	case int32:
		return String(strconv.FormatInt(int64(t), 10))
	case int64:
		return String(strconv.FormatInt(t, 10))
	case uint:
		return String(strconv.FormatUint(uint64(t), 10))
	case uint8: // byte is an alias for uint8
		return String(strconv.FormatUint(uint64(t), 10))
	case uint16:
		return String(strconv.FormatUint(uint64(t), 10))
	case uint32:
		return String(strconv.FormatUint(uint64(t), 10))
	case float32:
		return String(strconv.FormatFloat(float64(t), 'f', -1, 64))
	case float64:
		return String(strconv.FormatFloat(t, 'f', -1, 64))
	case time.Time:
		return String(t.String())
	case time.Duration:
		return String(t.String())
	default:
		panic(fmt.Errorf("%v(%T) is not supported yet", v, v))
	}
}

// NewStringValueNil returns an empty value.
func NewStringValueNil() Value {
	return String("")
}

func (s String) String() (string, bool) {
	return string(s), true
}

func (s String) Int64() (int64, bool) {
	iv, err := strconv.ParseInt(string(s), 10, 64)
	return iv, err == nil
}

func (s String) Uint64() (uint64, bool) {
	iv, err := strconv.ParseUint(string(s), 10, 64)
	return iv, err == nil
}

func (s String) Float64() (float64, bool) {
	f, err := strconv.ParseFloat(string(s), 64)
	return f, err == nil
}

func (s String) Time(layout string) (time.Time, bool) {
	t, err := time.Parse(layout, string(s))
	return t, err == nil
}

func (s String) Duration() (time.Duration, bool) {
	d, err := time.ParseDuration(string(s))
	return d, err == nil
}

func (s String) IsNil() bool {
	return len(s) == 0
}

func (s String) EqualTo(v Value) bool {
	tv, ok := v.(String)
	return ok && s == tv
}

func (s String) Copy() Value {
	return s
}
