package proc

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gyuho/psn/proc/schema"
)

// GetStat reads /proc/$PID/stat data.
func GetStat(pid int64, up Uptime) (Stat, error) {
	var (
		s   Stat
		err error
	)
	for i := 0; i < 3; i++ {
		s, err = getStat(pid, up)
		if err == nil {
			return s, nil
		}
		log.Println("retrying;", err)
		time.Sleep(5 * time.Millisecond)
	}
	return s, err
}

func getStat(pid int64, up Uptime) (Stat, error) {
	fpath := fmt.Sprintf("/proc/%d/stat", pid)
	f, err := openToRead(fpath)
	if err != nil {
		return Stat{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	st := &Stat{}
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fds := strings.Fields(txt)
		for i, fv := range fds {
			column := schema.ToField(schema.Stat[i].Name)
			s := reflect.ValueOf(st).Elem()
			if s.Kind() == reflect.Struct {
				f := s.FieldByName(column)
				if f.IsValid() {
					if f.CanSet() {
						switch schema.Stat[i].Kind {

						case reflect.Uint64:
							value, err := strconv.ParseUint(fv, 10, 64)
							if err != nil {
								return Stat{}, fmt.Errorf("%v when parsing %s %v", err, column, fv)
							}
							if !f.OverflowUint(value) {
								f.SetUint(value)

								bF := s.FieldByName(column + "BytesN")
								if bF.IsValid() {
									if bF.CanSet() {
										bF.SetUint(value)
									}
								}

								if schema.Stat[i].HumanizedBytes {
									hF := s.FieldByName(column + "HumanizedBytes")
									if hF.IsValid() {
										if hF.CanSet() {
											hF.SetString(humanize.Bytes(value))
										}
									}
								}
							}

						case reflect.Int64:
							value, err := strconv.ParseInt(fv, 10, 64)
							if err != nil {
								return Stat{}, fmt.Errorf("%v when parsing %s %v", err, column, fv)
							}
							if !f.OverflowInt(value) {
								f.SetInt(value)

								bF := s.FieldByName(column + "BytesN")
								if bF.IsValid() {
									if bF.CanSet() {
										bF.SetInt(value)
									}
								}

								if schema.Stat[i].HumanizedBytes {
									hF := s.FieldByName(column + "HumanizedBytes")
									if hF.IsValid() {
										if hF.CanSet() {
											if value > 0 {
												hF.SetString(humanize.Bytes(uint64(value)))
											}
										}
									}
								}
							}

						case reflect.String:
							f.SetString(fv)
						}
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Stat{}, err
	}
	return st.update(up)
}

func (s *Stat) update(up Uptime) (Stat, error) {
	if s == nil {
		return Stat{}, nil
	}
	if strings.HasPrefix(s.Comm, "(") {
		s.Comm = s.Comm[1:]
	}
	if strings.HasSuffix(s.Comm, ")") {
		s.Comm = s.Comm[:len(s.Comm)-1]
	}
	cu, err := s.getCPU(up)
	if err != nil {
		return Stat{}, err
	}
	s.CpuUsage = cu
	return *s, nil
}

// getCPU returns the average CPU usage in percentage.
// http://stackoverflow.com/questions/16726779/how-do-i-get-the-total-cpu-usage-of-an-application-from-proc-pid-stat
func (s Stat) getCPU(up Uptime) (float64, error) {
	totalSec := s.Utime + s.Stime
	totalSec += s.Cutime + s.Cstime

	out, err := exec.Command("/usr/bin/getconf", "CLK_TCK").Output()
	if err != nil {
		return 0, err
	}
	ot := strings.TrimSpace(strings.Replace(string(out), "\n", "", -1))
	hertz, err := strconv.ParseUint(ot, 10, 64)
	if err != nil || hertz == 0 {
		return 0, err
	}

	tookSec := up.UptimeTotal - (float64(s.Starttime) / float64(hertz))
	if hertz == 0 || tookSec == 0.0 {
		return 0.0, nil
	}
	return 100 * ((float64(totalSec) / float64(hertz)) / float64(tookSec)), nil
}

const statTmpl = `
----------------------------------------
[/proc/{{.Pid}}/stat]

Name:  {{.Comm}}
State: {{.State}}

Pid:         {{.Pid}}
Ppid:        {{.Ppid}}
NumThreads:  {{.NumThreads}}

Rss:       {{.RssHumanizedBytes}} ({{.RssBytesN}})
Rsslim:    {{.RsslimHumanizedBytes}} ({{.RsslimBytesN}})
Vsize:     {{.VsizeHumanizedBytes}} ({{.VsizeBytesN}})
CpuUsage:  {{.CpuUsage}} %

Starttime:  {{.Starttime}}
Utime:      {{.Utime}}
Stime:      {{.Stime}}
Cutime:     {{.Cutime}}
Cstime:     {{.Cstime}}

Session:   {{.Session}}
TtyNr:     {{.TtyNr}}
Tpgid:     {{.Tpgid}}
Flags:     {{.Flags}}

minflt:    {{.Minflt}}
cminflt:   {{.Cminflt}}
majflt:    {{.Majflt}}
cmajflt:   {{.Cmajflt}}

priority:     {{.Priority}}
nice:         {{.Nice}}
itrealvalue:  {{.Itrealvalue}}

startcode:    {{.Startcode}}
endcode:      {{.Endcode}}
startstack:   {{.Startstack}}
lstkesp:      {{.Kstkesp}}
lstkeip:      {{.Kstkeip}}
signal:       {{.Signal}}
blocked:      {{.Blocked}}
sigignore:    {{.Sigignore}}
sigcatch:     {{.Sigcatch}}
wchan:        {{.Wchan}}
nswap:        {{.Nswap}}
cnswap:       {{.Cnswap}}
exitSignal:   {{.ExitSignal}}
processor:    {{.Processor}}
rt_priority:  {{.RtPriority}}
policy:       {{.Policy}}

delayacct_blkio_ticks:
{{.DelayacctBlkioTicks}}

guest_time:   {{.GuestTime}}
cguest_time:  {{.CguestTime}}
start_data:   {{.StartData}}
end_data:     {{.EndData}}
start_brk:    {{.StartBrk}}
arg_start:    {{.ArgStart}}
arg_end:      {{.ArgEnd}}
env_start:    {{.EnvStart}}
env_end:      {{.EnvEnd}}
exit_code:    {{.ExitCode}}
----------------------------------------
`

func (s Stat) String() string {
	tpl := template.Must(template.New("statTmpl").Parse(statTmpl))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, s); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
