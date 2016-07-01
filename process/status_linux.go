package process

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
)

// GetStatus reads /proc/$PID/status data.
func GetStatus(pid int64) (Status, error) {
	var (
		s   Status
		err error
	)
	for i := 0; i < 3; i++ {
		s, err = getStatus(pid)
		if err == nil {
			return s, nil
		} else {
			log.Println(err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	return s, err
}

func getStatus(pid int64) (Status, error) {
	fpath := fmt.Sprintf("/proc/%d/status", pid)
	f, err := openToRead(fpath)
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
	rs.VmPeakBytes = u
	u, _ = humanize.ParseBytes(rs.VmSize)
	rs.VmSize = humanize.Bytes(u)
	rs.VmSizeBytes = u
	u, _ = humanize.ParseBytes(rs.VmLck)
	rs.VmLck = humanize.Bytes(u)
	rs.VmLckBytes = u
	u, _ = humanize.ParseBytes(rs.VmPin)
	rs.VmPin = humanize.Bytes(u)
	rs.VmPinBytes = u
	u, _ = humanize.ParseBytes(rs.VmHWM)
	rs.VmHWM = humanize.Bytes(u)
	rs.VmHWMBytes = u
	u, _ = humanize.ParseBytes(rs.VmRSS)
	rs.VmRSS = humanize.Bytes(u)
	rs.VmRSSBytes = u
	u, _ = humanize.ParseBytes(rs.VmData)
	rs.VmData = humanize.Bytes(u)
	rs.VmDataBytes = u
	u, _ = humanize.ParseBytes(rs.VmStk)
	rs.VmStk = humanize.Bytes(u)
	rs.VmStkBytes = u
	u, _ = humanize.ParseBytes(rs.VmExe)
	rs.VmExe = humanize.Bytes(u)
	rs.VmExeBytes = u
	u, _ = humanize.ParseBytes(rs.VmLib)
	rs.VmLib = humanize.Bytes(u)
	rs.VmLibBytes = u
	u, _ = humanize.ParseBytes(rs.VmPTE)
	rs.VmPTE = humanize.Bytes(u)
	rs.VmPTEBytes = u
	u, _ = humanize.ParseBytes(rs.VmSwap)
	rs.VmSwap = humanize.Bytes(u)
	rs.VmSwapBytes = u

	return rs, nil
}

const statusTmpl = `
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

Threads:   {{.Threads}}

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
	tpl := template.Must(template.New("statusTmpl").Parse(statusTmpl))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, s); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
