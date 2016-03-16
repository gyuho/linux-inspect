package ps

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gyuho/psn/tablesorter"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

// ListStatus finds the status by specifying the filter.
func ListStatus(filter *Status) ([]Status, error) {
	ds, _ := ioutil.ReadDir("/proc")
	procPaths := []string{}
	for _, f := range ds {
		if f.IsDir() && isInt(f.Name()) {
			procPaths = append(procPaths, filepath.Join("/proc", f.Name(), "status"))
		}
	}
	rs := []Status{}
	if filter.Pid != 0 {
		s, err := StatusByPID(filter.Pid)
		return []Status{s}, err
	} else {
		rc, errc := make(chan Status, len(procPaths)), make(chan error)
		skip := make(chan struct{})
		for _, fpath := range procPaths {
			go func(fpath string, filter *Status) {
				st, err := parseStatus(fpath)
				if err != nil {
					errc <- err
				} else if st.Match(filter) {
					rc <- st
				} else {
					skip <- struct{}{}
				}
			}(fpath, filter)
		}
		cnt := 0
		for cnt != len(procPaths) {
			select {
			case s := <-rc:
				rs = append(rs, s)
			case e := <-errc:
				return nil, e
			case <-skip:
			}
			cnt++
		}
	}
	return rs, nil
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Status is 'proc/pid/status'
// See http://man7.org/linux/man-pages/man5/proc.5.html.
type Status struct {
	// Name is the command run by this process.
	Name string `yaml:"Name"`

	// State is Current state of the process.  One of "R (running)",
	// "S (sleeping)", "D (disk sleep)", "T (stopped)", "T (tracing stop)",
	// "Z (zombie)", or "X (dead)".
	State string `yaml:"State"`

	// Tgid is thread group ID.
	Tgid int `yaml:"Tgid"`
	Ngid int `yaml:"Ngid"`

	// Pid is process ID.
	Pid int `yaml:"Pid"`

	// PPid is parent process ID, which launches the Pid.
	PPid int `yaml:"PPid"`

	// TracerPid is PID of process tracing this process (0 if not
	// being traced).
	TracerPid int `yaml:"TracerPid"`

	Uid string `yaml:"Uid"`
	Gid string `yaml:"Gid"`

	// FDSize is the number of file descriptor slots currently allocated.
	FDSize int `yaml:"FDSize"`

	// Groups is supplementary group list.
	Groups string `yaml:"Groups"`

	// VmPeak is peak virtual memory usage.
	// Vm includes physical memory and swap.
	VmPeak string `yaml:"VmPeak"`
	// VmPeakUint64 is VmPeak in bytes.
	VmPeakUint64 uint64

	// VmSize is current virtual memory usage.
	// VmSize is the total amount of memory required for
	// this process.
	VmSize string `yaml:"VmSize"`
	// VmSizeUint64 is VmSize in bytes.
	VmSizeUint64 uint64

	// VmLck is current mlocked memory.
	VmLck string `yaml:"VmLck"`
	// VmLckUint64 is VmLck in bytes.
	VmLckUint64 uint64

	// VmPin is pinned memory size.
	VmPin string `yaml:"VmPin"`
	// VmPinUint64 is VmPin in bytes.
	VmPinUint64 uint64

	// VmRSS is peak resident set size.
	VmHWM string `yaml:"VmHWM"`
	// VmHWMUint64 is VmHWM in bytes.
	VmHWMUint64 uint64

	// VmRSS is resident set size. VmRSS is the actual
	// amount in memory. Some memory can be swapped out
	// to physical disk. So this is the real memory usage
	// of the process.
	VmRSS string `yaml:"VmRSS"`
	// VmRSSUint64 is VmRSS in bytes.
	VmRSSUint64 uint64

	// VmRSS is size of data segment.
	VmData string `yaml:"VmData"`
	// VmDataUint64 is VmData in bytes.
	VmDataUint64 uint64

	// VmStk is size of stack.
	VmStk string `yaml:"VmStk"`
	// VmStkUint64 is VmStk in bytes.
	VmStkUint64 uint64

	// VmExe is size of text segment.
	VmExe string `yaml:"VmExe"`
	// VmExeUint64 is VmExe in bytes.
	VmExeUint64 uint64

	// VmLib is shared library usage.
	VmLib string `yaml:"VmLib"`
	// VmLibUint64 is VmLib in bytes.
	VmLibUint64 uint64

	// VmPTE is page table entries size.
	VmPTE string `yaml:"VmPTE"`
	// VmPTEUint64 is VmPTE in bytes.
	VmPTEUint64 uint64

	// VmSwap is swap space used.
	VmSwap string `yaml:"VmSwap"`
	// VmSwapUint64 is VmSwap in bytes.
	VmSwapUint64 uint64

	// Threads is the number of threads in process
	// containing this thread (process).
	Threads int `yaml:"Threads"`

	SigQ                     string `yaml:"SigQ"`
	SigPnd                   string `yaml:"SigPnd"`
	ShdPnd                   string `yaml:"ShdPnd"`
	SigBlk                   string `yaml:"SigBlk"`
	SigIgn                   string `yaml:"SigIgn"`
	SigCgt                   string `yaml:"SigCgt"`
	CapInh                   string `yaml:"CapInh"`
	CapPrm                   string `yaml:"CapPrm"`
	CapEff                   string `yaml:"CapEff"`
	CapBnd                   string `yaml:"CapBnd"`
	Seccomp                  int    `yaml:"Seccomp"`
	CpusAllowed              string `yaml:"Cpus_allowed"`
	CpusAllowedList          string `yaml:"Cpus_allowed_list"`
	MemsAllowed              string `yaml:"Mems_allowed"`
	MemsAllowedList          string `yaml:"Mems_allowed_list"`
	VoluntaryCtxtSwitches    int    `yaml:"voluntary_ctxt_switches"`
	NonvoluntaryCtxtSwitches int    `yaml:"nonvoluntary_ctxt_switches"`
}

const statusTmpl = `
-----------------------
[/proc/{{.Pid}}/status]

Name:  {{.Name}}
State: {{.State}}

Pid:  {{.Pid}}
PPid: {{.PPid}}

FDSize:  {{.FDSize}}
Threads: {{.Threads}}

VmRSS:   {{.VmRSS}}
VmSize:  {{.VmSize}}
VmPeak:  {{.VmPeak}}
-----------------------
`

var statusMembers = []string{
	"NAME",
	"STATE",
	"PID",
	"PPID",
	"FD",
	"THREADS",
	"VM_RSS",
	"VM_SIZE",
	"VM_PEAK",

	"VmSizeUint64",
	"VmRSSUint64",
	"VmPeakUint64",
}

// String prints out only parts of the status.
func (s Status) String() string {
	tpl := template.Must(template.New("statusTmpl").Parse(statusTmpl))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, s); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

// WriteToTable writes slice of Status to ASCII table.
func WriteToTable(w io.Writer, top int, sts ...Status) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(statusMembers[:9:9])

	rows := make([][]string, len(sts))
	for i, s := range sts {
		sl := make([]string, len(statusMembers))
		sl[0] = s.Name
		sl[1] = s.State
		sl[2] = strconv.Itoa(s.Pid)
		sl[3] = strconv.Itoa(s.PPid)
		sl[4] = strconv.Itoa(s.FDSize)
		sl[5] = strconv.Itoa(s.Threads)
		sl[6] = s.VmRSS
		sl[7] = s.VmSize
		sl[8] = s.VmPeak

		sl[9] = strconv.Itoa(int(s.VmRSSUint64))
		sl[10] = strconv.Itoa(int(s.VmSizeUint64))
		sl[11] = strconv.Itoa(int(s.VmPeakUint64))

		rows[i] = sl
	}

	tablesorter.By(
		rows,
		tablesorter.MakeDescendingIntFunc(9),  // VM_RSS
		tablesorter.MakeAscendingFunc(0),      // NAME
		tablesorter.MakeDescendingIntFunc(10), // VM_SIZE
		tablesorter.MakeDescendingIntFunc(4),  // FD
	).Sort(rows)

	if top != 0 && len(rows) > top {
		rows = rows[:top:top]
	}
	for _, row := range rows {
		table.Append(row[:9:9])
	}

	table.Render()
}

// WriteToTable writes slice of Status to a csv file.
func WriteToCSV(firstCSV bool, f *os.File, sts ...Status) error {
	wr := csv.NewWriter(f)
	if firstCSV { // write header
		if err := wr.Write(append([]string{"timestamp"}, statusMembers...)); err != nil {
			return err
		}
	}

	rows := make([][]string, len(sts))
	for i, s := range sts {
		sl := make([]string, len(statusMembers))
		sl[0] = s.Name
		sl[1] = s.State
		sl[2] = strconv.Itoa(s.Pid)
		sl[3] = strconv.Itoa(s.PPid)
		sl[4] = strconv.Itoa(s.FDSize)
		sl[5] = strconv.Itoa(s.Threads)
		sl[6] = s.VmRSS
		sl[7] = s.VmSize
		sl[8] = s.VmPeak

		sl[9] = strconv.Itoa(int(s.VmRSSUint64))
		sl[10] = strconv.Itoa(int(s.VmSizeUint64))
		sl[11] = strconv.Itoa(int(s.VmPeakUint64))

		rows[i] = sl
	}

	tablesorter.By(
		rows,
		tablesorter.MakeDescendingIntFunc(9),  // VM_RSS
		tablesorter.MakeAscendingFunc(0),      // NAME
		tablesorter.MakeDescendingIntFunc(10), // VM_SIZE
		tablesorter.MakeDescendingIntFunc(4),  // FD
	).Sort(rows)

	// adding timestamp
	ts := time.Now().String()[:19]
	nrows := make([][]string, len(rows))
	for i, row := range rows {
		nrows[i] = append([]string{ts}, row...)
	}
	if err := wr.WriteAll(nrows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}

const statusTmplDetailed = `
----------------------------------------
[/proc/{{.Pid}}/status]

Name:  {{.Name}}
State: {{.State}}

Tgid:      {{.Tgid}}
Ngid:      {{.Ngid}}
Pid:       {{.Pid}}
PPid:      {{.PPid}}
TracerPid: {{.TracerPid}}

FDSize:  {{.FDSize}}

VmPeak:  {{.VmPeak}}
VmSize:  {{.VmSize}}
VmLck:   {{.VmLck}}
VmPin:   {{.VmPin}}
VmHWM:   {{.VmHWM}}
VmRSS:   {{.VmRSS}}
VmData:  {{.VmData}}
VmStk:   {{.VmStk}}
VmExe:   {{.VmExe}}
VmLib:   {{.VmLib}}
VmPTE:   {{.VmPTE}}
VmSwap:  {{.VmSwap}}

Threads: {{.Threads}}

Groups: {{.Groups}}
Uid:    {{.Uid}}
Gid:    {{.Gid}}

SigQ:   {{.SigQ}}
SigPnd: {{.SigPnd}}
ShdPnd: {{.ShdPnd}}
SigBlk: {{.SigBlk}}
SigIgn: {{.SigIgn}}
SigCgt: {{.SigCgt}}
CapInh: {{.CapInh}}
CapPrm: {{.CapPrm}}
CapEff: {{.CapEff}}
CapBnd: {{.CapBnd}}

Seccomp: {{.Seccomp}}

Cpus_allowed:      {{.CpusAllowed}}
Cpus_allowed_list: {{.CpusAllowedList}}

Mems_allowed:      {{.MemsAllowed}}
Mems_allowed_list: {{.MemsAllowedList}}

voluntary_ctxt_switches:
{{.VoluntaryCtxtSwitches}}

nonvoluntary_ctxt_switches:
{{.NonvoluntaryCtxtSwitches}}
----------------------------------------
`

func (s Status) StringDetailed() string {
	tpl := template.Must(template.New("statusTmplDetailed").Parse(statusTmplDetailed))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, s); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

// Match matches status to filter.
// Only supports name, pid, ppid.
func (s *Status) Match(filter *Status) bool {
	if s == nil {
		return false
	}
	if filter == nil {
		return true // no need to compare
	}
	if filter.Name != "" {
		if filter.Name != s.Name {
			return false
		}
	}
	if filter.Pid != 0 {
		if filter.Pid != s.Pid {
			return false
		}
	}
	if filter.PPid != 0 {
		if filter.PPid != s.PPid {
			return false
		}
	}
	if filter.State != "" {
		if filter.State != s.State {
			return false
		}
	}
	return true
}

// StatusByPID gets the Status of the pid.
func StatusByPID(pid int) (Status, error) {
	fpath := fmt.Sprintf("/proc/%d/status", pid)
	return parseStatus(fpath)
}

func parseStatus(fpath string) (Status, error) {
	_, err := os.Stat(fpath)
	if err != nil {
		return Status{}, err
	}
	f, err := open(fpath)
	if err != nil {
		return Status{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return Status{}, err
	}
	rs := Status{}
	if err := yaml.Unmarshal(b, &rs); err != nil {
		return Status{}, err
	}
	u, _ := humanize.ParseBytes(rs.VmPeak)
	rs.VmPeak = humanize.Bytes(u)
	rs.VmPeakUint64 = u
	u, _ = humanize.ParseBytes(rs.VmSize)
	rs.VmSize = humanize.Bytes(u)
	rs.VmSizeUint64 = u
	u, _ = humanize.ParseBytes(rs.VmLck)
	rs.VmLck = humanize.Bytes(u)
	rs.VmLckUint64 = u
	u, _ = humanize.ParseBytes(rs.VmPin)
	rs.VmPin = humanize.Bytes(u)
	rs.VmPinUint64 = u
	u, _ = humanize.ParseBytes(rs.VmHWM)
	rs.VmHWM = humanize.Bytes(u)
	rs.VmHWMUint64 = u
	u, _ = humanize.ParseBytes(rs.VmRSS)
	rs.VmRSS = humanize.Bytes(u)
	rs.VmRSSUint64 = u
	u, _ = humanize.ParseBytes(rs.VmData)
	rs.VmData = humanize.Bytes(u)
	rs.VmDataUint64 = u
	u, _ = humanize.ParseBytes(rs.VmStk)
	rs.VmStk = humanize.Bytes(u)
	rs.VmStkUint64 = u
	u, _ = humanize.ParseBytes(rs.VmExe)
	rs.VmExe = humanize.Bytes(u)
	rs.VmExeUint64 = u
	u, _ = humanize.ParseBytes(rs.VmLib)
	rs.VmLib = humanize.Bytes(u)
	rs.VmLibUint64 = u
	u, _ = humanize.ParseBytes(rs.VmPTE)
	rs.VmPTE = humanize.Bytes(u)
	rs.VmPTEUint64 = u
	u, _ = humanize.ParseBytes(rs.VmSwap)
	rs.VmSwap = humanize.Bytes(u)
	rs.VmSwapUint64 = u
	return rs, nil
}

func open(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0444)
	if err != nil {
		return f, err
	}
	return f, nil
}

// Kill kills all processes in arguments.
func Kill(w io.Writer, parent bool, sts ...Status) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(w, "Kill:", err)
		}
	}()

	pidToKill := make(map[int]string)
	for _, s := range sts {
		pidToKill[s.Pid] = s.Name
		if parent {
			pidToKill[s.PPid] = s.Name
		}
	}
	pids := []int{}
	for pid := range pidToKill {
		pids = append(pids, pid)
	}
	sort.Ints(pids)

	for _, pid := range pids {
		fmt.Fprintf(w, "syscall.Kill: %s [PID: %d]\n", pidToKill[pid], pid)
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			fmt.Fprintf(w, "error when sending syscall.SIGTERM (%v):", err)

			shell := os.Getenv("SHELL")
			if len(shell) == 0 {
				shell = "sh"
			}
			args := []string{shell, "-c", fmt.Sprintf("sudo kill -9 %d", pid)}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = w
			cmd.Stderr = w
			fmt.Fprintf(w, "Starting: %q\n", cmd.Args)
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(w, "error when 'sudo kill' (%v)", err)
			}
			if err := cmd.Wait(); err != nil {
				fmt.Fprintf(w, "Start(%s) cmd.Wait returned %v\n", cmd.Path, err)
			}
		}
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			fmt.Fprintf(w, "error when sending syscall.SIGKILL (%v):", err)

			shell := os.Getenv("SHELL")
			if len(shell) == 0 {
				shell = "sh"
			}
			args := []string{shell, "-c", fmt.Sprintf("sudo kill -9 %d", pid)}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = w
			cmd.Stderr = w
			fmt.Fprintf(w, "Starting: %q\n", cmd.Args)
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(w, "error when 'sudo kill' (%v)", err)
			}
			if err := cmd.Wait(); err != nil {
				fmt.Fprintf(w, "Start(%s) cmd.Wait returned %v\n", cmd.Path, err)
			}
		}
	}
}
