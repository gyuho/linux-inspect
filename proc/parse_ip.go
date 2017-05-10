package proc

import (
	"fmt"
	"strconv"
	"strings"
)

// parseLittleEndianIpv4 parses hexadecimal ipv4 IP addresses.
// For example, it converts '0101007F:0035' into string.
// It assumes that the system has little endian order.
func parseLittleEndianIpv4(s string) (string, int64, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", 0, fmt.Errorf("cannot parse ipv4 %s", s)
	}
	if len(arr[0]) != 8 {
		return "", 0, fmt.Errorf("cannot parse ipv4 ip %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", 0, fmt.Errorf("cannot parse ipv4 port %s", arr[1])
	}

	d0, err := strconv.ParseInt(arr[0][6:8], 16, 32)
	if err != nil {
		return "", 0, err
	}
	d1, err := strconv.ParseInt(arr[0][4:6], 16, 32)
	if err != nil {
		return "", 0, err
	}
	d2, err := strconv.ParseInt(arr[0][2:4], 16, 32)
	if err != nil {
		return "", 0, err
	}
	d3, err := strconv.ParseInt(arr[0][0:2], 16, 32)
	if err != nil {
		return "", 0, err
	}
	ip := fmt.Sprintf("%d.%d.%d.%d", d0, d1, d2, d3)

	port, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}

// parseLittleEndianIpv6 parses hexadecimal ipv6 IP addresses.
// For example, it converts '4506012691A700C165EB1DE1F912918C:8BDA'
// into string. It assumes that the system has little endian order.
func parseLittleEndianIpv6(s string) (string, int64, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", 0, fmt.Errorf("cannot parse ipv6 %s", s)
	}
	if len(arr[0]) != 32 {
		return "", 0, fmt.Errorf("cannot parse ipv6 ip %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", 0, fmt.Errorf("cannot parse ipv6 port %s", arr[1])
	}

	// 32 characters, reverse by 2 characters
	ip := ""
	for i := 32; i > 0; i -= 2 {
		ip += fmt.Sprintf("%s", arr[0][i-2:i])
		if (len(ip)+1)%5 == 0 && len(ip) != 39 {
			ip += ":"
		}
	}

	port, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}
