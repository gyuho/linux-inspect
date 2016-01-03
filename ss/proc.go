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
			fmt.Fprintln(w, "[Kill - panic]", err)
		}
	}()
	for _, v := range ps {
		fmt.Fprintf(w, "[Kill] syscall.Kill -> %s\n", v)
		if err := syscall.Kill(v.PID, syscall.SIGINT); err != nil {
			fmt.Fprintln(w, "[Kill - error]", err)
		}
	}
	fmt.Fprintln(w, "[Kill] Done!")
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

// ListProcess lists all processes.
func ListProcess(opt TransportProtocol) ([]Process, error) {
	return listProcess(opt)
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

	// sort 2D string slice by PROGRAM, PID, LOCAL_ADDR, USER
	by(
		rows,
		makeAscendingFunc(1),
		makeAscendingFunc(2),
		makeAscendingFunc(3),
		makeAscendingFunc(5),
	).Sort(rows)
	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}

// by returns a multiSorter that sorts using the less functions
func by(rows [][]string, lesses ...lessFunc) *multiSorter {
	return &multiSorter{
		data: rows,
		less: lesses,
	}
}

// lessFunc compares between two string slices.
type lessFunc func(p1, p2 *[]string) bool

func makeAscendingFunc(idx int) func(row1, row2 *[]string) bool {
	return func(row1, row2 *[]string) bool {
		return (*row1)[idx] < (*row2)[idx]
	}
}

// multiSorter implements the Sort interface,
// sorting the two dimensional string slices within.
type multiSorter struct {
	data [][]string
	less []lessFunc
}

// Sort sorts the rows according to lessFunc.
func (ms *multiSorter) Sort(rows [][]string) {
	sort.Sort(ms)
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.data)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.data[i], ms.data[j] = ms.data[j], ms.data[i]
}

// Less is part of sort.Interface.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.data[i], &ms.data[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q
			return true
		case less(q, p):
			// p > q
			return false
		}
		// p == q; try next comparison
	}
	return ms.less[k](p, q)
}
