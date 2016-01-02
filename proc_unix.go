package ssn

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TransportProtocol int

const (
	TCP TransportProtocol = iota
	TCP6
)

var (
	transportProtocols = map[string]TransportProtocol{
		"tcp":  TCP,
		"tcp6": TCP6,
	}
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

// parseLittleEndianIpv4 parses hexadecimal ipv4 IP addresses.
// For example, it converts '0101007F:0035' into string.
// It assumes that the system has little endian order.
func parseLittleEndianIpv4(s string) (string, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", fmt.Errorf("cannot parse ipv4 %s", s)
	}
	if len(arr[0]) != 8 {
		return "", fmt.Errorf("cannot parse ipv4 host %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", fmt.Errorf("cannot parse ipv4 port %s", arr[1])
	}

	d0, err := strconv.ParseInt(arr[0][6:8], 16, 32)
	if err != nil {
		return "", err
	}
	d1, err := strconv.ParseInt(arr[0][4:6], 16, 32)
	if err != nil {
		return "", err
	}
	d2, err := strconv.ParseInt(arr[0][2:4], 16, 32)
	if err != nil {
		return "", err
	}
	d3, err := strconv.ParseInt(arr[0][0:2], 16, 32)
	if err != nil {
		return "", err
	}
	host := fmt.Sprintf("%d.%d.%d.%d", d0, d1, d2, d3)

	p0, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", err
	}
	port := fmt.Sprintf("%d", p0)

	return host + ":" + port, nil
}

// parseLittleEndianIpv6 parses hexadecimal ipv6 IP addresses.
// For example, it converts '4506012691A700C165EB1DE1F912918C:8BDA'
// into string.
// It assumes that the system has little endian order.
func parseLittleEndianIpv6(s string) (string, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", fmt.Errorf("cannot parse ipv6 %s", s)
	}
	if len(arr[0]) != 32 {
		return "", fmt.Errorf("cannot parse ipv6 host %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", fmt.Errorf("cannot parse ipv6 port %s", arr[1])
	}

	// 32 characters, reverse by 2 characters
	host := ""
	for i := 32; i > 0; i -= 2 {
		host += fmt.Sprintf("%s", arr[0][i-2:i])
		if (len(host)+1)%5 == 0 && len(host) != 39 {
			host += ":"
		}
	}

	p0, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", err
	}
	port := fmt.Sprintf("%d", p0)

	return "[" + host + "]" + ":" + port, nil
}
