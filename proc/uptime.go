package proc

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/pkg/timeutil"
)

// GetUptime reads '/proc/uptime'.
func GetUptime() (Uptime, error) {
	f, err := fileutil.OpenToRead("/proc/uptime")
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
		u.UptimeTotalParsedTime = timeutil.HumanizeDurationSecond(uint64(v))
	}
	if len(fields) > 1 {
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeIdle = v
		u.UptimeIdleParsedTime = timeutil.HumanizeDurationSecond(uint64(v))
	}
	return u, nil
}
