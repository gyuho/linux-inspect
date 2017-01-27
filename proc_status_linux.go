package psn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
)

// GetProcStatusByPID reads '/proc/$PID/status' data.
func GetProcStatusByPID(pid int64) (s Status, err error) {
	return parseProcStatusByPID(pid)
}

func rawProcStatus(pid int64) (Status, error) {
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
		return rs, err
	}
	return rs, nil
}

func parseProcStatusByPID(pid int64) (Status, error) {
	rs, err := rawProcStatus(pid)
	if err != nil {
		return rs, err
	}

	rs.StateParsedStatus = strings.TrimSpace(rs.State)

	u, _ := humanize.ParseBytes(rs.VmPeak)
	rs.VmPeakBytesN = u
	rs.VmPeakParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmSize)
	rs.VmSizeBytesN = u
	rs.VmSizeParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmLck)
	rs.VmLckBytesN = u
	rs.VmLckParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmPin)
	rs.VmPinBytesN = u
	rs.VmPinParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmHWM)
	rs.VmHWMBytesN = u
	rs.VmHWMParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmRSS)
	rs.VmRSSBytesN = u
	rs.VmRSSParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmData)
	rs.VmDataBytesN = u
	rs.VmDataParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmStk)
	rs.VmStkBytesN = u
	rs.VmStkParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmExe)
	rs.VmExeBytesN = u
	rs.VmExeParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmLib)
	rs.VmLibBytesN = u
	rs.VmLibParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmPTE)
	rs.VmPTEBytesN = u
	rs.VmPTEParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmPMD)
	rs.VmPMDBytesN = u
	rs.VmPMDParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.VmSwap)
	rs.VmSwapBytesN = u
	rs.VmSwapParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(rs.HugetlbPages)
	rs.HugetlbPagesBytesN = u
	rs.HugetlbPagesParsedBytes = humanize.Bytes(u)

	return rs, nil
}

const statusTmpl = `
----------------------------------------
[/proc/{{.Pid}}/status]

Name:   {{.Name}}
Umask:  {{.Umask}}
State:  {{.StateParsedStatus}}

Tgid:      {{.Tgid}}
Ngid:      {{.Ngid}}
Pid:       {{.Pid}}
PPid:      {{.PPid}}
TracerPid: {{.TracerPid}}

FDSize:  {{.FDSize}}

VmPeak:  {{.VmPeakParsedBytes}}
VmSize:  {{.VmSizeParsedBytes}}
VmLck:   {{.VmLckParsedBytes}}
VmPin:   {{.VmPinParsedBytes}}
VmHWM:   {{.VmHWMParsedBytes}}
VmRSS:   {{.VmRSSParsedBytes}}
VmData:  {{.VmDataParsedBytes}}
VmStk:   {{.VmStkParsedBytes}}
VmExe:   {{.VmExeParsedBytes}}
VmLib:   {{.VmLibParsedBytes}}
VmPTE:   {{.VmPTEParsedBytes}}
VmPMD:   {{.VmPMDParsedBytes}}
VmSwap:  {{.VmSwapParsedBytes}}
HugetlbPages:  {{.HugetlbPagesParsedBytes}}

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
