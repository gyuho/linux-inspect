package proc

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/schema"

	"github.com/dustin/go-humanize"
)

// GetStatByPID reads '/proc/$PID/stat' data.
func GetStatByPID(pid int64) (s Stat, err error) {
	var d []byte
	d, err = readStat(pid)
	if err != nil {
		return Stat{}, err
	}
	return parseStat(d)
}

func readStat(pid int64) ([]byte, error) {
	fpath := fmt.Sprintf("/proc/%d/stat", pid)
	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func parseStat(d []byte) (s Stat, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(d))
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fds := strings.Fields(txt)
		for i, fv := range fds {
			column := schema.ToField(StatSchema.Columns[i].Name)
			val := reflect.ValueOf(&s).Elem()
			if val.Kind() == reflect.Struct {
				f := val.FieldByName(column)
				if f.IsValid() {
					if f.CanSet() {
						switch StatSchema.Columns[i].Kind {

						case reflect.Uint64:
							uv, uerr := strconv.ParseUint(fv, 10, 64)
							if uerr != nil {
								return Stat{}, fmt.Errorf("%v when parsing %s %v", uerr, column, fv)
							}
							if !f.OverflowUint(uv) {
								f.SetUint(uv)

								fval := val.FieldByName(column + "BytesN")
								if fval.IsValid() {
									if fval.CanSet() {
										fval.SetUint(uv)
									}
								}

								if vv, ok := StatSchema.ColumnsToParse[StatSchema.Columns[i].Name]; ok {
									switch vv {
									case schema.TypeBytes:
										hF := val.FieldByName(column + "ParsedBytes")
										if hF.IsValid() {
											if hF.CanSet() {
												hF.SetString(humanize.Bytes(uv))
											}
										}
									}
								}
							}

						case reflect.Int64:
							iv, ierr := strconv.ParseInt(fv, 10, 64)
							if ierr != nil {
								return Stat{}, fmt.Errorf("%v when parsing %s %v", ierr, column, fv)
							}
							if !f.OverflowInt(iv) {
								f.SetInt(iv)

								fval := val.FieldByName(column + "BytesN")
								if fval.IsValid() {
									if fval.CanSet() {
										fval.SetInt(iv)
									}
								}

								if vv, ok := StatSchema.ColumnsToParse[StatSchema.Columns[i].Name]; ok {
									switch vv {
									case schema.TypeBytes:
										fval := val.FieldByName(column + "ParsedBytes")
										if fval.IsValid() {
											if fval.CanSet() {
												fval.SetString(humanize.Bytes(uint64(iv)))
											}
										}
									}
								}
							}

						case reflect.String:
							f.SetString(fv)

							if vv, ok := StatSchema.ColumnsToParse[StatSchema.Columns[i].Name]; ok {
								switch vv {
								case schema.TypeStatus:
									fval := val.FieldByName(column + "ParsedStatus")
									if fval.IsValid() {
										if fval.CanSet() {
											fval.SetString(convertStatus(fv))
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return s, err
	}
	if strings.HasPrefix(s.Comm, "(") {
		s.Comm = s.Comm[1:]
	}
	if strings.HasSuffix(s.Comm, ")") {
		s.Comm = s.Comm[:len(s.Comm)-1]
	}
	return s, err
}

const statTmpl = `
----------------------------------------
[/proc/{{.Pid}}/stat]

Name:  {{.Comm}}
State: {{.StateParsedStatus}}

Pid:         {{.Pid}}
Ppid:        {{.Ppid}}
NumThreads:  {{.NumThreads}}

Rss:       {{.RssParsedBytes}} ({{.RssBytesN}})
Rsslim:    {{.RsslimParsedBytes}} ({{.RsslimBytesN}})
Vsize:     {{.VsizeParsedBytes}} ({{.VsizeBytesN}})

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
