package dataframe

import (
	"fmt"
	"time"
)

// TimeDefaultLayout is used to parse time.Time.
var TimeDefaultLayout = "2006-01-02 15:04:05 -0700 MST"

// DATA_TYPE defines dataframe data types.
type DATA_TYPE uint8

const (
	// STRING represents Go string or bytes.
	STRING DATA_TYPE = iota

	// TIME represents Go time.Time type.
	TIME
)

func (dt DATA_TYPE) String() string {
	switch dt {
	case STRING:
		return "STRING"
	case TIME:
		return "TIME"
	default:
		panic(fmt.Errorf("DATA_TYPE %d is unknown", dt))
	}
}

// ReflectTypeOf returns the DATA_TYPE.
func ReflectTypeOf(v interface{}) DATA_TYPE {
	switch v.(type) {
	case time.Time:
		return TIME
	default:
		return STRING
	}
}

// ToValue converts to Value.
func ToValue(v interface{}) Value {
	switch ReflectTypeOf(v) {
	case TIME:
		return NewTimeValue(v)
	default:
		return NewStringValue(v)
	}
}
