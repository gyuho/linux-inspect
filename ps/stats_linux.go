package ps

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

// Stats represents the proc/$PID/stat file
// specification as documented in http://man7.org/linux/man-pages/man5/proc.5.html.
var Stats = [...]struct {
	Col      string
	Kind     reflect.Kind
	Humanize bool // true if humanized string is needed
}{
	{"pid", reflect.Int64, false},    // the process ID
	{"comm", reflect.String, false},  // filename of the executable (originally in parentheses, automatically removed by this package)
	{"state", reflect.String, false}, // One character that represents the state of the process
	{"ppid", reflect.Int64, false},   // PID of the parent of this process
	{"pgrp", reflect.Int64, false},   // process group ID of the process
	{"session", reflect.Int64, false},
	{"tty_nr", reflect.Int64, false},
	{"tpgid", reflect.Int64, false}, // ID of the foreground process group of the controlling terminal of the process
	{"flags", reflect.Uint64, false},
	{"minflt", reflect.Uint64, false},  // number of minor faults the process has made which have not required loading a memory page from disk.
	{"cminflt", reflect.Uint64, false}, // number of minor faults that the process's waited-for children have made.
	{"majflt", reflect.Uint64, false},  // number of major faults the process has made which have required loading a memory page from disk.
	{"cmajflt", reflect.Uint64, false}, // number of major faults that the process's waited-for children have made.
	{"utime", reflect.Uint64, false},   // Amount of time that this process has been scheduled in user mode, measured in clock ticks.
	{"stime", reflect.Uint64, false},   // Amount of time that this process has been scheduled in kernel mode, measured in clock ticks.
	{"cutime", reflect.Uint64, false},  // Amount of time that this process's waited-for children have been scheduled in user mode.
	{"cstime", reflect.Uint64, false},  // Amount of time that this process's waited-for children have been scheduled in kernel mode.
	{"priority", reflect.Int64, false},
	{"nice", reflect.Int64, false},
	{"num_threads", reflect.Int64, false},
	{"itrealvalue", reflect.Int64, false},
	{"starttime", reflect.Uint64, false}, // time the process started after system boot.
	{"vsize", reflect.Uint64, true},      // Virtual memory size in bytes.
	{"rss", reflect.Int64, true},         // Resident Set Size: number of pages the process has in real memory.
	{"rsslim", reflect.Uint64, true},
	{"startcode", reflect.Uint64, false},
	{"endcode", reflect.Uint64, false},
	{"startstack", reflect.Uint64, false},
	{"kstkesp", reflect.Uint64, false},
	{"kstkeip", reflect.Uint64, false},
	{"signal", reflect.Uint64, false},
	{"blocked", reflect.Uint64, false},
	{"sigignore", reflect.Uint64, false},
	{"sigcatch", reflect.Uint64, false},
	{"wchan", reflect.Uint64, false},
	{"nswap", reflect.Uint64, false},
	{"cnswap", reflect.Uint64, false},
	{"exit_signal", reflect.Int64, false},
	{"processor", reflect.Int64, false}, // CPU number last executed on.
	{"rt_priority", reflect.Uint64, false},
	{"policy", reflect.Uint64, false},
	{"delayacct_blkio_ticks", reflect.Uint64, false},
	{"guest_time", reflect.Uint64, false},
	{"cguest_time", reflect.Int64, false},
	{"start_data", reflect.Uint64, false},
	{"end_data", reflect.Uint64, false},
	{"start_brk", reflect.Uint64, false},
	{"arg_start", reflect.Uint64, false},
	{"arg_end", reflect.Uint64, false},
	{"env_start", reflect.Uint64, false},
	{"env_end", reflect.Uint64, false},
	{"exit_code", reflect.Int64, false},
}

func getStat(fpath string) (Stat, error) {
	_, err := os.Stat(fpath)
	if err != nil {
		return Stat{}, err
	}
	f, err := open(fpath)
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
			column := ToField(Stats[i].Col)
			s := reflect.ValueOf(st).Elem()
			if s.Kind() == reflect.Struct {
				f := s.FieldByName(column)
				if f.IsValid() {
					if f.CanSet() {
						switch Stats[i].Kind {
						case reflect.Uint64:
							value, err := strconv.ParseUint(fv, 10, 64)
							if err != nil {
								return Stat{}, fmt.Errorf("%v when parsing %s %v", err, column, fv)
							}
							if !f.OverflowUint(value) {
								f.SetUint(value)
								if Stats[i].Humanize {
									hF := s.FieldByName(column + "Humanize")
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
								if Stats[i].Humanize {
									hF := s.FieldByName(column + "Humanize")
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
	return process(st), nil
}

func process(s *Stat) Stat {
	if s == nil {
		return Stat{}
	}
	if strings.HasPrefix(s.Comm, "(") {
		s.Comm = s.Comm[1:]
	}
	if strings.HasSuffix(s.Comm, ")") {
		s.Comm = s.Comm[:len(s.Comm)-1]
	}
	cu, _ := s.GetCpuUsage()
	s.CpuUsage = cu
	return *s
}

const statTmpl = `
----------------------
[/proc/{{.Pid}}/stat]

Name:  {{.Comm}}
State: {{.State}}

Pid:         {{.Pid}}
Ppid:        {{.Ppid}}
NumThreads:  {{.NumThreads}}

Rss:       {{.RssHumanize}} ({{.Rss}})
Rsslim:    {{.RsslimHumanize}} ({{.Rsslim}})
Vsize:     {{.VsizeHumanize}} ({{.Vsize}})
CpuUsage:  {{.CpuUsage}} %

Starttime:  {{.Starttime}}
Utime:      {{.Utime}}
Stime:      {{.Stime}}
Cutime:     {{.Cutime}}
Cstime:     {{.Cstime}}
-----------------------
`

const statTmplDetailed = `
----------------------------------------
[/proc/{{.Pid}}/stat]

Name:  {{.Comm}}
State: {{.State}}

Pid:         {{.Pid}}
Ppid:        {{.Ppid}}
NumThreads:  {{.NumThreads}}

Rss:       {{.RssHumanize}} ({{.Rss}})
Rsslim:    {{.RsslimHumanize}} ({{.Rsslim}})
Vsize:     {{.VsizeHumanize}} ({{.Vsize}})
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

func (s Stat) StringDetailed() string {
	tpl := template.Must(template.New("statTmplDetailed").Parse(statTmplDetailed))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, s); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

// GetCpuUsage returns the average CPU usage in percentage.
// http://stackoverflow.com/questions/16726779/how-do-i-get-the-total-cpu-usage-of-an-application-from-proc-pid-stat
func (s Stat) GetCpuUsage() (float64, error) {
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

	u, err := GetUptime()
	if err != nil {
		return 0, err
	}
	tookSec := u.UptimeTotal - (float64(s.Starttime) / float64(hertz))
	return 100 * ((float64(totalSec) / float64(hertz)) / float64(tookSec)), nil
}

// Uptime describes /proc/uptime.
type Uptime struct {
	UptimeTotal float64
	UptimeIdle  float64
}

// GetUptime reads /proc/uptime.
func GetUptime() (Uptime, error) {
	f, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return Uptime{}, err
	}
	fields := strings.Fields(strings.TrimSpace(string(f)))
	u := Uptime{}
	if len(fields) > 0 {
		v, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeTotal = v
	}
	if len(fields) > 1 {
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return Uptime{}, err
		}
		u.UptimeIdle = v
	}
	return u, nil
}
