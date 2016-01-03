package ss

import (
	"fmt"
	"io"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

type TransportProtocol int

const (
	TCP TransportProtocol = iota
	TCP6
)

var (
	stringToProtocol = map[string]TransportProtocol{
		"tcp":  TCP,
		"tcp6": TCP6,
	}
	protocolToString = map[TransportProtocol]string{
		TCP:  "tcp",
		TCP6: "tcp6",
	}
)

// Process describes OS processes.
type Process struct {
	Protocol string

	Program string
	PID     int

	LocalIP   string
	LocalPort string

	RemoteIP   string
	RemotePort string

	State string

	User user.User
}

var processMembers = []string{
	"PROTOCOL",
	"PROGRAM",
	"PID",
	"LOCAL_ADDR",
	"REMOTE_ADDR",
	"USER",
}

// String describes Process type.
func (p Process) String() string {
	return fmt.Sprintf("Protocol: '%s' | Program: '%s' | PID: '%d' | LocalAddr: '%s%s' | RemoteAddr: '%s%s' | State: '%s' | User: '%+v'",
		p.Protocol,
		p.Program,
		p.PID,
		p.LocalIP,
		p.LocalPort,
		p.RemoteIP,
		p.RemotePort,
		p.State,
		p.User,
	)
}

// Kill kills all processes in arguments.
func Kill(w io.Writer, ps ...Process) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(w, "Kill:", err)
		}
	}()

	pidToKill := make(map[int]string)
	pids := []int{}
	for _, p := range ps {
		if _, ok := pidToKill[p.PID]; !ok {
			pidToKill[p.PID] = p.Program
			pids = append(pids, p.PID)
		}
	}
	sort.Ints(pids)

	for _, pid := range pids {
		fmt.Fprintf(w, "syscall.Kill: %s [PID: %d]\n", pidToKill[pid], pid)
		if err := syscall.Kill(pid, syscall.SIGINT); err != nil {
			fmt.Fprintln(w, "Kill:", err)
		}
	}
}

// parseLittleEndianIpv4 parses hexadecimal ipv4 IP addresses.
// For example, it converts '0101007F:0035' into string.
// It assumes that the system has little endian order.
func parseLittleEndianIpv4(s string) (string, string, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("cannot parse ipv4 %s", s)
	}
	if len(arr[0]) != 8 {
		return "", "", fmt.Errorf("cannot parse ipv4 ip %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", "", fmt.Errorf("cannot parse ipv4 port %s", arr[1])
	}

	d0, err := strconv.ParseInt(arr[0][6:8], 16, 32)
	if err != nil {
		return "", "", err
	}
	d1, err := strconv.ParseInt(arr[0][4:6], 16, 32)
	if err != nil {
		return "", "", err
	}
	d2, err := strconv.ParseInt(arr[0][2:4], 16, 32)
	if err != nil {
		return "", "", err
	}
	d3, err := strconv.ParseInt(arr[0][0:2], 16, 32)
	if err != nil {
		return "", "", err
	}
	ip := fmt.Sprintf("%d.%d.%d.%d", d0, d1, d2, d3)

	p0, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", "", err
	}
	port := fmt.Sprintf("%d", p0)

	return ip, ":" + port, nil
}

// parseLittleEndianIpv6 parses hexadecimal ipv6 IP addresses.
// For example, it converts '4506012691A700C165EB1DE1F912918C:8BDA'
// into string.
// It assumes that the system has little endian order.
func parseLittleEndianIpv6(s string) (string, string, error) {
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("cannot parse ipv6 %s", s)
	}
	if len(arr[0]) != 32 {
		return "", "", fmt.Errorf("cannot parse ipv6 ip %s", arr[0])
	}
	if len(arr[1]) != 4 {
		return "", "", fmt.Errorf("cannot parse ipv6 port %s", arr[1])
	}

	// 32 characters, reverse by 2 characters
	ip := ""
	for i := 32; i > 0; i -= 2 {
		ip += fmt.Sprintf("%s", arr[0][i-2:i])
		if (len(ip)+1)%5 == 0 && len(ip) != 39 {
			ip += ":"
		}
	}

	p0, err := strconv.ParseInt(arr[1], 16, 32)
	if err != nil {
		return "", "", err
	}
	port := fmt.Sprintf("%d", p0)

	return ip, ":" + port, nil
}

// List lists all processes.
func List(opt TransportProtocol) ([]Process, error) {
	return list(opt)
}

// ListProgram lists all processes running a specific program.
func ListProgram(opt TransportProtocol, program string) ([]Process, error) {
	ps, err := list(opt)
	if err != nil {
		return nil, err
	}
	ns := []Process{}
	for _, p := range ps {
		if strings.HasSuffix(p.Program, program) {
			ns = append(ns, p)
		}
	}
	return ns, nil
}

// ListTcpPorts lists all TCP ports that are being used.
func ListTcpPorts() map[string]struct{} {
	ps4, _ := List(TCP)
	ps6, _ := List(TCP6)
	rm := make(map[string]struct{})
	for _, p := range ps4 {
		rm[p.LocalPort] = struct{}{}
	}
	for _, p := range ps6 {
		rm[p.RemotePort] = struct{}{}
	}
	return rm
}

// WriteToTable writes slice of Processes to ASCII table.
func WriteToTable(w io.Writer, ps ...Process) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(processMembers)

	rows := make([][]string, len(ps))
	for i, p := range ps {
		sl := make([]string, len(processMembers))
		sl[0] = p.Protocol
		sl[1] = p.Program
		sl[2] = strconv.Itoa(p.PID)
		sl[3] = p.LocalIP + p.LocalPort
		sl[4] = p.RemoteIP + p.RemotePort
		sl[5] = p.User.Name
		rows[i] = sl
	}

	by(
		rows,
		makeAscendingFunc(1), // PROGRAM
		makeAscendingFunc(2), // PID
		makeAscendingFunc(3), // LOCAL_ADDR
		makeAscendingFunc(5), // USER
	).Sort(rows)

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}
