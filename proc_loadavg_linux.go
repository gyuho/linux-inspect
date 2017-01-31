package psn

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type procLoadAvgColumnIndex int

const (
	proc_loadavg_idx_load_avg_1_minute procLoadAvgColumnIndex = iota
	proc_loadavg_idx_load_avg_5_minute
	proc_loadavg_idx_load_avg_15_minute
	proc_loadavg_idx_kernel_scheduling_entities_with_slash
	proc_loadavg_idx_pid
)

// GetProcLoadAvg reads '/proc/loadavg'.
// Expected output is '0.37 0.47 0.39 1/839 31397'.
func GetProcLoadAvg() (LoadAvg, error) {
	txt, err := readProcLoadAvg()
	if err != nil {
		return LoadAvg{}, err
	}
	return getProcLoadAvg(txt)
}

func readProcLoadAvg() (string, error) {
	f, err := openToRead("/proc/loadavg")
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

func getProcLoadAvg(txt string) (LoadAvg, error) {
	ds := strings.Fields(txt)
	if len(ds) < 5 {
		return LoadAvg{}, fmt.Errorf("not enough columns at %v", ds)
	}

	lavg := LoadAvg{}

	avg1, err := strconv.ParseFloat(ds[proc_loadavg_idx_load_avg_1_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg1Minute = avg1

	avg5, err := strconv.ParseFloat(ds[proc_loadavg_idx_load_avg_5_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg5Minute = avg5

	avg15, err := strconv.ParseFloat(ds[proc_loadavg_idx_load_avg_15_minute], 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.LoadAvg15Minute = avg15

	slashed := strings.Split(ds[proc_loadavg_idx_kernel_scheduling_entities_with_slash], "/")
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

	pid, err := strconv.ParseInt(ds[proc_loadavg_idx_pid], 10, 64)
	if err != nil {
		return LoadAvg{}, err
	}
	lavg.Pid = pid

	return lavg, nil
}
