package psn

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
		}
		log.Println(err)
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
	rs.VmPeakBytesN = u
	rs.VmPeakHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmSize)
	rs.VmSizeBytesN = u
	rs.VmSizeHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmLck)
	rs.VmLckBytesN = u
	rs.VmLckHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmPin)
	rs.VmPinBytesN = u
	rs.VmPinHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmHWM)
	rs.VmHWMBytesN = u
	rs.VmHWMHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmRSS)
	rs.VmRSSBytesN = u
	rs.VmRSSHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmData)
	rs.VmDataBytesN = u
	rs.VmDataHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmStk)
	rs.VmStkBytesN = u
	rs.VmStkHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmExe)
	rs.VmExeBytesN = u
	rs.VmExeHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmLib)
	rs.VmLibBytesN = u
	rs.VmLibHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmPTE)
	rs.VmPTEBytesN = u
	rs.VmPTEHumanizedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmSwap)
	rs.VmSwapBytesN = u
	rs.VmSwapHumanizedBytes = humanize.Bytes(u)

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

VmPeak:  {{.VmPeakHumanizedBytes}}
VmSize:  {{.VmSizeHumanizedBytes}}
VmLck:   {{.VmLckHumanizedBytes}}
VmPin:   {{.VmPinHumanizedBytes}}
VmHWM:   {{.VmHWMHumanizedBytes}}
VmRSS:   {{.VmRSSHumanizedBytes}}
VmData:  {{.VmDataHumanizedBytes}}
VmStk:   {{.VmStkHumanizedBytes}}
VmExe:   {{.VmExeHumanizedBytes}}
VmLib:   {{.VmLibHumanizedBytes}}
VmPTE:   {{.VmPTEHumanizedBytes}}
VmSwap:  {{.VmSwapHumanizedBytes}}

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
