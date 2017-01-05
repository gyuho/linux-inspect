package proc

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// GetUptime reads '/proc/uptime'.
func GetUptime() (Uptime, error) {
	f, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return Uptime{}, err
	}
	fields := strings.Fields(strings.TrimSpace(string(f)))

	now := time.Now()

	u := Uptime{}
	if len(fields) > 0 {
		v, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeTotal = v
		u.UptimeTotalHumanizedTime = humanize.Time(now.Add(-1 * time.Duration(u.UptimeTotal) * time.Second))
	}
	if len(fields) > 1 {
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeIdle = v
		u.UptimeIdleHumanizedTime = humanize.Time(now.Add(-1 * time.Duration(u.UptimeIdle) * time.Second))
	}
	return u, nil
}
