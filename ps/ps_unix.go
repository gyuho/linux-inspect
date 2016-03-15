package ps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/dustin/go-humanize"
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
	rs := []Status{}
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
	return rs, nil
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Status is 'proc/pid/status'
type Status struct {
	Name                     string `yaml:"Name"`
	State                    string `yaml:"State"`
	Tgid                     int    `yaml:"Tgid"`
	Ngid                     int    `yaml:"Ngid"`
	Pid                      int    `yaml:"Pid"`
	PPid                     int    `yaml:"PPid"`
	TracerPid                int    `yaml:"TracerPid"`
	Uid                      string `yaml:"Uid"`
	Gid                      string `yaml:"Gid"`
	FDSize                   int    `yaml:"FDSize"`
	Groups                   string `yaml:"Groups"`
	VmPeak                   string `yaml:"VmPeak"`
	VmPeakUint64             uint64
	VmSize                   string `yaml:"VmSize"`
	VmSizeUint64             uint64
	VmLck                    string `yaml:"VmLck"`
	VmLckUint64              uint64
	VmPin                    string `yaml:"VmPin"`
	VmPinUint64              uint64
	VmHWM                    string `yaml:"VmHWM"`
	VmHWMUint64              uint64
	VmRSS                    string `yaml:"VmRSS"`
	VmRSSUint64              uint64
	VmData                   string `yaml:"VmData"`
	VmDataUint64             uint64
	VmStk                    string `yaml:"VmStk"`
	VmStkUint64              uint64
	VmExe                    string `yaml:"VmExe"`
	VmExeUint64              uint64
	VmLib                    string `yaml:"VmLib"`
	VmLibUint64              uint64
	VmPTE                    string `yaml:"VmPTE"`
	VmPTEUint64              uint64
	VmSwap                   string `yaml:"VmSwap"`
	VmSwapUint64             uint64
	Threads                  int    `yaml:"Threads"`
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

const statusTemplate = `
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

func (s Status) String() string {
	tpl := template.Must(template.New("statusTemplate").Parse(statusTemplate))
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
	return true
}

// GetStatusByPID gets the Status of the pid.
func GetStatusByPID(pid int) (Status, error) {
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
	rs.VmPeakUint64 = u
	u, _ = humanize.ParseBytes(rs.VmSize)
	rs.VmSizeUint64 = u
	u, _ = humanize.ParseBytes(rs.VmLck)
	rs.VmLckUint64 = u
	u, _ = humanize.ParseBytes(rs.VmPin)
	rs.VmPinUint64 = u
	u, _ = humanize.ParseBytes(rs.VmHWM)
	rs.VmHWMUint64 = u
	u, _ = humanize.ParseBytes(rs.VmRSS)
	rs.VmRSSUint64 = u
	u, _ = humanize.ParseBytes(rs.VmData)
	rs.VmDataUint64 = u
	u, _ = humanize.ParseBytes(rs.VmStk)
	rs.VmStkUint64 = u
	u, _ = humanize.ParseBytes(rs.VmExe)
	rs.VmExeUint64 = u
	u, _ = humanize.ParseBytes(rs.VmLib)
	rs.VmLibUint64 = u
	u, _ = humanize.ParseBytes(rs.VmPTE)
	rs.VmPTEUint64 = u
	u, _ = humanize.ParseBytes(rs.VmSwap)
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
