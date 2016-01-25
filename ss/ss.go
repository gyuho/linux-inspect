package ss

import (
	"fmt"
	"io"
	"net"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"sync"
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
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
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

// List lists all processes. filter is used to return Processes
// that matches the member values in filter struct.
func List(filter *Process, opts ...TransportProtocol) ([]Process, error) {
	rs := []Process{}
	for _, opt := range opts {

		ps, err := list(opt)
		if err != nil {
			return nil, err
		}

		if filter == nil {
			rs = append(rs, ps...)
			continue
		}

		pCh, done := make(chan Process), make(chan struct{})
		fv := *filter

		for _, p := range ps {

			go func(ft, process Process) {

				if !ft.Match(process) {
					done <- struct{}{}
					return
				}
				pCh <- process

			}(fv, p)

		}

		cn := 0
		ns := []Process{}
		for cn != len(ps) {
			select {
			case p := <-pCh:
				ns = append(ns, p)
				cn++
			case <-done:
				cn++
			}
		}

		close(pCh)
		close(done)

		rs = append(rs, ns...)
	}
	return rs, nil
}

// Match returns true if the Process matches the filter.
func (filter Process) Match(p Process) bool {
	if filter.Protocol != "" {
		if p.Protocol != filter.Protocol {
			return false
		}
	}
	// matches the suffix
	if filter.Program != "" {
		if !strings.HasSuffix(p.Program, filter.Program) {
			return false
		}
	}
	if filter.PID != 0 {
		if p.PID != filter.PID {
			return false
		}
	}
	if filter.LocalIP != "" {
		if p.LocalIP != filter.LocalIP {
			return false
		}
	}
	if filter.LocalPort != "" {
		p0 := p.LocalPort
		if !strings.HasPrefix(p0, ":") {
			p0 = ":" + p0
		}
		p1 := filter.LocalPort
		if !strings.HasPrefix(p1, ":") {
			p1 = ":" + p1
		}
		if p0 != p1 {
			return false
		}
	}
	if filter.RemoteIP != "" {
		if p.RemoteIP != filter.RemoteIP {
			return false
		}
	}
	if filter.RemotePort != "" {
		p0 := p.RemotePort
		if !strings.HasPrefix(p0, ":") {
			p0 = ":" + p0
		}
		p1 := filter.RemotePort
		if !strings.HasPrefix(p1, ":") {
			p1 = ":" + p1
		}
		if p0 != p1 {
			return false
		}
	}
	if filter.State != "" {
		if p.State != filter.State {
			return false
		}
	}
	// currently only support user name
	if filter.User.Username != "" {
		if p.User.Username != filter.User.Username {
			return false
		}
	}
	return true
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
		makeAscendingFunc(0), // PROTOCOL
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

// ListPorts lists all ports that are being used.
func ListPorts(filter *Process, opts ...TransportProtocol) map[string]struct{} {
	rm := make(map[string]struct{})
	ps, _ := List(filter, opts...)
	for _, p := range ps {
		rm[p.LocalPort] = struct{}{}
		rm[p.RemotePort] = struct{}{}
	}
	return rm
}

type Ports struct {
	mu        sync.Mutex // guards the following
	beingUsed map[string]struct{}
}

func NewPorts() *Ports {
	p := &Ports{}
	p.beingUsed = make(map[string]struct{})
	return p
}

func (p *Ports) Exist(port string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.beingUsed[port]
	return ok
}

func (p *Ports) Add(ports ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, port := range ports {
		p.beingUsed[port] = struct{}{}
	}
}

func (p *Ports) Free(ports ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, port := range ports {
		delete(p.beingUsed, port)
	}
}

// Refresh refreshes ports that are being used by the OS.
// You would use as follows:
//
//	var (
//		globalPorts     = NewPorts()
//		refreshInterval = 10 * time.Second
//	)
//	func init() {
//		globalPorts.Refresh()
//		go func() {
//			for {
//				select {
//				case <-time.After(refreshInterval):
//					globalPorts.Refresh()
//				}
//			}
//		}()
//	}
//
func (p *Ports) Refresh() {
	p.mu.Lock()
	p.beingUsed = ListPorts(nil, TCP, TCP6)
	p.mu.Unlock()
}

func getFreePort(opts ...TransportProtocol) (string, error) {
	fp := ""
	var errMsg error
	for _, opt := range opts {
		protocolStr := ""
		switch opt {
		case TCP:
			protocolStr = "tcp"
		case TCP6:
			protocolStr = "tcp6"
		}
		addr, err := net.ResolveTCPAddr(protocolStr, "localhost:0")
		if err != nil {
			errMsg = err
			continue
		}
		l, err := net.ListenTCP(protocolStr, addr)
		if err != nil {
			errMsg = err
			continue
		}
		pd := l.Addr().(*net.TCPAddr).Port
		l.Close()
		fp = ":" + strconv.Itoa(pd)
		break
	}
	return fp, errMsg
}

// GetFreePorts returns multiple free ports from OS.
func GetFreePorts(num int, opts ...TransportProtocol) ([]string, error) {
	try := 0
	rm := make(map[string]struct{})
	for len(rm) != num {
		if try > 150 {
			return nil, fmt.Errorf("too many tries to find free ports")
		}
		try++
		p, err := getFreePort(opts...)
		if err != nil {
			return nil, err
		}
		rm[p] = struct{}{}
	}
	rs := []string{}
	for p := range rm {
		rs = append(rs, p)
	}
	sort.Strings(rs)
	return rs, nil
}
