package proc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/gyuho/linux-inspect/pkg/fileutil"

	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
)

// GetStatusByPID reads '/proc/$PID/status' data.
func GetStatusByPID(pid int64) (s Status, err error) {
	d, derr := readStatus(pid)
	if derr != nil {
		return Status{}, derr
	}
	s, err = parseStatus(d)
	if err != nil {
		return s, err
	}

	s.StateParsedStatus = strings.TrimSpace(s.State)

	u, _ := humanize.ParseBytes(s.VmPeak)
	s.VmPeakBytesN = u
	s.VmPeakParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmSize)
	s.VmSizeBytesN = u
	s.VmSizeParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmLck)
	s.VmLckBytesN = u
	s.VmLckParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmPin)
	s.VmPinBytesN = u
	s.VmPinParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmHWM)
	s.VmHWMBytesN = u
	s.VmHWMParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmRSS)
	s.VmRSSBytesN = u
	s.VmRSSParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmData)
	s.VmDataBytesN = u
	s.VmDataParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmStk)
	s.VmStkBytesN = u
	s.VmStkParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmExe)
	s.VmExeBytesN = u
	s.VmExeParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmLib)
	s.VmLibBytesN = u
	s.VmLibParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmPTE)
	s.VmPTEBytesN = u
	s.VmPTEParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmPMD)
	s.VmPMDBytesN = u
	s.VmPMDParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.VmSwap)
	s.VmSwapBytesN = u
	s.VmSwapParsedBytes = humanize.Bytes(u)
	u, _ = humanize.ParseBytes(s.HugetlbPages)
	s.HugetlbPagesBytesN = u
	s.HugetlbPagesParsedBytes = humanize.Bytes(u)

	return s, nil
}

func readStatus(pid int64) ([]byte, error) {
	fpath := fmt.Sprintf("/proc/%d/status", pid)
	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func parseStatus(d []byte) (s Status, err error) {
	err = yaml.Unmarshal(d, &s)
	return s, err
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

// GetProgram returns the program name.
func GetProgram(pid int64) (string, error) {
	// Readlink needs root permission
	// return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	s, err := GetStatusByPID(pid)
	return s.Name, err
}
