package psn

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

// GetDevice returns the device name where dir is mounted.
// It parses '/etc/mtab'.
// TODO: use this just to make it easierto find!
func GetDevice(mounted string) (string, error) {
	f, err := openToRead("/etc/mtab")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fields := strings.Fields(txt)
		if len(fields) < 2 {
			continue
		}

		dev := strings.TrimSpace(fields[0])
		at := strings.TrimSpace(fields[1])
		if mounted == at {
			return filepath.Base(dev), nil
		}
	}

	return "", fmt.Errorf("no device found, mounted at %q", mounted)
}
