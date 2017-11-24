package dataframe

import "time"

// Value represents the value in data frame.
type Value interface {
	// String parses Value to string. It returns false if not possible.
	String() (string, bool)

	// Int64 parses Value to int64. It returns false if not possible.
	Int64() (int64, bool)

	// Uint64 parses Value to uint64. It returns false if not possible.
	Uint64() (uint64, bool)

	// Float64 parses Value to float64. It returns false if not possible.
	Float64() (float64, bool)

	// Time parses Value to time.Time based on the layout. It returns false if not possible.
	Time(layout string) (time.Time, bool)

	// Duration parses Value to time.Duration. It returns false if not possible.
	Duration() (time.Duration, bool)

	// IsNil returns true if the Value is nil.
	IsNil() bool

	// EqualTo returns true if the Value is equal to v.
	EqualTo(v Value) bool

	// Copy copies Value.
	Copy() Value
}
