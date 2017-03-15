package proc

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
)

type loadAvgColumnIndex int

const (
	load_avg_idx_load_avg_1_minute loadAvgColumnIndex = iota
	load_avg_idx_load_avg_5_minute
	load_avg_idx_load_avg_15_minute
	load_avg_idx_kernel_scheduling_entities_with_slash
	load_avg_idx_pid
)

// GetLoadAvg reads '/proc/loadavg'.
// Expected output is '0.37 0.47 0.39 1/839 31397'.
func GetLoadAvg() (LoadAvg, error) {
	txt, err := readLoadAvg()
	if err != nil {
		return LoadAvg{}, err
	}
	return getLoadAvg(txt)
}

func readLoadAvg() (string, error) {
	f, err := fileutil.OpenToRead("/proc/loadavg")
	if err != nil {
		return "", err
	}
	defer f.Close()

	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bts)), nil
}

func getLoadAvg(txt string) (LoadAvg, error) {
	ds := strings.Fields(txt)
	if len(ds) < 5 {
		return LoadAvg{}, fmt.Errorf("not enough columns at %v", ds)
	}

	lavg := LoadAvg{}

	avg1, err := strconv.ParseFloat(ds[load_avg_idx_load_avg_1_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg1Minute = avg1

	avg5, err := strconv.ParseFloat(ds[load_avg_idx_load_avg_5_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg5Minute = avg5

	avg15, err := strconv.ParseFloat(ds[load_avg_idx_load_avg_15_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg15Minute = avg15

	slashed := strings.Split(ds[load_avg_idx_kernel_scheduling_entities_with_slash], "/")
	if len(slashed) != 2 {
		return LoadAvg{}, fmt.Errorf("expected '/' string in kernel scheduling entities field, got %v", slashed)
	}
	s1, err := strconv.ParseInt(slashed[0], 10, 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.RunnableKernelSchedulingEntities = s1

	s2, err := strconv.ParseInt(slashed[1], 10, 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.CurrentKernelSchedulingEntities = s2

	pid, err := strconv.ParseInt(ds[load_avg_idx_pid], 10, 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.Pid = pid

	return lavg, nil
}
