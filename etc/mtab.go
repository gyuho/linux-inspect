package etc

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
)

const mtabPath = "/etc/mtab"

type columnIndex int

const (
	idx_file_system columnIndex = iota
	idx_mounted_on
	idx_file_system_type
	idx_options
	idx_dump
	idx_pass
)

// GetMtab returns '/etc/mtab' information.
func GetMtab() ([]Mtab, error) {
	if !fileutil.Exist(mtabPath) {
		return nil, fmt.Errorf("%q does not exist", mtabPath)
	}
	f, err := fileutil.OpenToRead(mtabPath)
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
		if len(ms) < int(idx_pass+1) {
			return nil, fmt.Errorf("not enough columns at %v", ms)
		}

		m := Mtab{
			FileSystem:     strings.TrimSpace(ms[idx_file_system]),
			MountedOn:      strings.TrimSpace(ms[idx_mounted_on]),
			FileSystemType: strings.TrimSpace(ms[idx_file_system_type]),
			Options:        strings.TrimSpace(ms[idx_options]),
		}

		mn, err := strconv.ParseInt(ms[idx_dump], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Dump = int(mn)

		mn, err = strconv.ParseInt(ms[idx_dump], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Pass = int(mn)

		mss = append(mss, m)
	}

	return mss, nil
}
