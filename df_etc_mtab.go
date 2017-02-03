package psn

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type etcMtabColumnIndex int

const (
	etc_mtab_idx_file_system etcMtabColumnIndex = iota
	etc_mtab_idx_mounted_on
	etc_mtab_idx_file_system_type
	etc_mtab_idx_options
	etc_mtab_idx_dump
	etc_mtab_idx_pass
)

// GetEtcMtab returns '/etc/mtab' information.
func GetEtcMtab() ([]Mtab, error) {
	f, err := openToRead("/etc/mtab")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mss := []Mtab{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}
		ms := strings.Fields(strings.TrimSpace(txt))
		if len(ms) < int(etc_mtab_idx_pass+1) {
			return nil, fmt.Errorf("not enough columns at %v", ms)
		}

		m := Mtab{
			FileSystem:     strings.TrimSpace(ms[etc_mtab_idx_file_system]),
			MountedOn:      strings.TrimSpace(ms[etc_mtab_idx_mounted_on]),
			FileSystemType: strings.TrimSpace(ms[etc_mtab_idx_file_system_type]),
			Options:        strings.TrimSpace(ms[etc_mtab_idx_options]),
		}

		mn, err := strconv.ParseInt(ms[etc_mtab_idx_dump], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Dump = int(mn)

		mn, err = strconv.ParseInt(ms[etc_mtab_idx_dump], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Pass = int(mn)

		mss = append(mss, m)
	}

	return mss, nil
}
