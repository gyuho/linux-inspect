// Package timeutil implements time utilities.
package timeutil

import (
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

func NowPST() time.Time {
	tzone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(tzone)
}

func HumanizeDurationMs(ms uint64) string {
	s := humanize.Time(time.Now().Add(-1 * time.Duration(ms) * time.Millisecond))
	if s == "now" {
		s = "0 seconds"
	}
	return strings.TrimSpace(strings.Replace(s, " ago", "", -1))
}

func HumanizeDurationSecond(sec uint64) string {
	s := humanize.Time(time.Now().Add(-1 * time.Duration(sec) * time.Second))
	if s == "now" {
		s = "0 seconds"
	}
	return strings.TrimSpace(strings.Replace(s, " ago", "", -1))
}
