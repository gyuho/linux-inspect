// generate generates ps struct based on the schema.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gyuho/psn/proc/schema"
)

func generate(cols ...schema.Column) string {
	buf := new(bytes.Buffer)
	for i := range cols {
		tagstr := "yaml"
		if !cols[i].YAMLTag {
			tagstr = "column"
		}
		buf.WriteString(fmt.Sprintf(
			"\t%s\t%s\t`%s:\"%s\"`\n",
			schema.ToField(cols[i].Name),
			goType(cols[i].Kind),
			tagstr,
			cols[i].Name,
		))
		if cols[i].HumanizedSeconds {
			buf.WriteString(fmt.Sprintf("\t%sHumanizedTime\tstring\t`%s:\"%s_humanized_time\"`\n",
				schema.ToField(cols[i].Name),
				tagstr,
				cols[i].Name,
			))
		} else if cols[i].HumanizedBytes {
			if cols[i].Kind == reflect.String {
				buf.WriteString(fmt.Sprintf("\t%sBytesN\tuint64\t`%s:\"%s_bytes_n\"`\n",
					schema.ToField(cols[i].Name),
					tagstr,
					cols[i].Name,
				))
			}
			buf.WriteString(fmt.Sprintf("\t%sHumanizedBytes\tstring\t`%s:\"%s_humanized_bytes\"`\n",
				schema.ToField(cols[i].Name),
				tagstr,
				cols[i].Name,
			))
		}
	}
	return buf.String()
}

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString(`package proc

	` + "// updated at " + nowPST().String() + `

// Proc represents '/proc' in linux.
type Proc struct {
	DiskStats DiskStats
	Stat      Stat
	Status    Status
	IO        IO
}

`)

	// 'proc/uptime'
	buf.WriteString(`// Uptime is 'proc/uptime' in linux.
type Uptime struct {
`)
	buf.WriteString(generate(schema.Uptime...))
	buf.WriteString("}\n\n")

	// 'proc/diskstats'
	buf.WriteString(`// DiskStats is 'proc/diskstats' in linux.
type DiskStats struct {
`)
	buf.WriteString(generate(schema.DiskStats...))
	buf.WriteString("}\n\n")

	// 'proc/$PID/stat'
	buf.WriteString(`// Stat is 'proc/$PID/stat' in linux.
type Stat struct {
`)
	buf.WriteString(generate(schema.Stat...))
	for _, line := range additionalFields {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	buf.WriteString("}\n\n")

	// 'proc/$PID/status'
	buf.WriteString(`// Status is 'proc/$PID/status' in linux.
type Status struct {
`)
	buf.WriteString(generate(schema.Status...))
	buf.WriteString("}\n\n")

	// 'proc/$PID/io'
	buf.WriteString(`// IO is 'proc/$PID/io' in linux.
type IO struct {
`)
	buf.WriteString(generate(schema.IO...))
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := toFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/proc/generated_linux.go")); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/proc")); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		log.Fatal(err)
	}
}

func goType(tp reflect.Kind) string {
	switch tp {
	case reflect.Float64:
		return "float64"
	case reflect.Uint64:
		return "uint64"
	case reflect.Int64:
		return "int64"
	case reflect.String:
		return "string"
	default:
		panic(fmt.Errorf("unknown type %q", tp.String()))
	}
}

func nowPST() time.Time {
	tzone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(tzone)
}

func openToRead(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0444)
	if err != nil {
		return f, err
	}
	return f, nil
}

func toFile(txt, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		return err
	}
	return nil
}

var additionalFields = [...]string{
	"CpuUsage float64 `column:\"cpu_usage\"`",
}
