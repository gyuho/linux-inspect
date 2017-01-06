package psn

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ListPIDs reads all PIDs in '/proc'.
func ListPIDs() ([]int64, error) {
	ds, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	pids := make([]int64, 0, len(ds))
	for _, f := range ds {
		if f.IsDir() && isInt(f.Name()) {
			id, err := strconv.ParseInt(f.Name(), 10, 64)
			if err != nil {
				return nil, err
			}
			pids = append(pids, id)
		}
	}
	return pids, nil
}

// ListProcFds reads '/proc/*/fd/*' to grab process IDs.
func ListProcFds() ([]string, error) {
	// returns the names of all files matching pattern
	// or nil if there is no matching file
	fs, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func pidFromFd(s string) (int64, error) {
	// get 5261 from '/proc/5261/fd/69'
	return strconv.ParseInt(strings.Split(s, "/")[2], 10, 64)
}

// GetProgram returns the program name.
func GetProgram(pid int64) (string, error) {
	return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
}
