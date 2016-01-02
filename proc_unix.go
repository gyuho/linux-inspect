package ssn

import (
	"bufio"
	"os"
	"strings"
)

var (
	protocols = map[TransportProtocol]string{
		TCP:  "/proc/net/tcp",
		TCP6: "/proc/net/tcp6",
	}
)

// readProcNet reads '/proc/net/*' in OS.
func readProcNet(opt TransportProtocol) ([][]string, error) {
	f, err := os.OpenFile(protocols[opt], os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fields := [][]string{}

	scanner := bufio.NewScanner(f)

	colN := 0
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fs := strings.Fields(txt)
		if colN == 0 {
			colN = len(fs)
		}

		if len(fs) > colN {
			// join extra columns to the last column with whitespace character
			cf := make([]string, colN)
			copy(cf, fs)
			cf[colN-1] = strings.Join(fs[colN-1:], " ")
			fs = cf
		}

		fields = append(fields, fs)
	}

	if err := scanner.Err(); err != nil {
		return fields, err
	}

	return fields, nil
}
