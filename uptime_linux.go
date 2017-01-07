package psn

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetUptime reads '/proc/uptime'.
func GetUptime() (Uptime, error) {
	f, err := openToRead("/proc/uptime")
	if err != nil {
		return Uptime{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return Uptime{}, err
	}
	fields := strings.Fields(strings.TrimSpace(string(b)))

	u := Uptime{}
	if len(fields) > 0 {
		v, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeTotal = v
		u.UptimeTotalParsedTime = humanizeDurationSecond(uint64(v))
	}
	if len(fields) > 1 {
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeIdle = v
		u.UptimeIdleParsedTime = humanizeDurationSecond(uint64(v))
	}
	return u, nil
}
