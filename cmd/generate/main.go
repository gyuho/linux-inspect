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

	"github.com/gyuho/psn/schema"
)

func generate(raw schema.RawData) string {
	tagstr := "yaml"
	if !raw.IsYAML {
		tagstr = "column"
	}

	buf := new(bytes.Buffer)
	for _, col := range raw.Columns {
		if col.Godoc != "" {
			buf.WriteString(fmt.Sprintf("\t// %s is %s.\n", schema.ToField(col.Name), col.Godoc))
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`%s:\"%s\"`\n",
			schema.ToField(col.Name),
			schema.GoType(col.Kind),
			tagstr,
			schema.ToYAMLField(col.Name),
		))

		// additional parsed column
		if v, ok := raw.ColumnsToParse[col.Name]; ok {
			switch v {
			case schema.TypeBytes:
				ntstr := "uint64"
				if col.Kind == reflect.Int64 {
					ntstr = "int64"
				}
				buf.WriteString(fmt.Sprintf("\t%sBytesN\t%s\t`%s:\"%s_bytes_n\"`\n",
					schema.ToField(col.Name),
					ntstr,
					tagstr,
					schema.ToYAMLField(col.Name),
				))
				buf.WriteString(fmt.Sprintf("\t%sParsedBytes\tstring\t`%s:\"%s_parsed_bytes\"`\n",
					schema.ToField(col.Name),
					tagstr,
					schema.ToYAMLField(col.Name),
				))

			case schema.TypeTimeSeconds:
				buf.WriteString(fmt.Sprintf("\t%sParsedTime\tstring\t`%s:\"%s_parsed_time\"`\n",
					schema.ToField(col.Name),
					tagstr,
					schema.ToYAMLField(col.Name),
				))

			case schema.TypeIPAddress:
				buf.WriteString(fmt.Sprintf("\t%sParsedIPAddress\tstring\t`%s:\"%s_parsed_ip_address\"`\n",
					schema.ToField(col.Name),
					tagstr,
					schema.ToYAMLField(col.Name),
				))

			case schema.TypeStatus:
				buf.WriteString(fmt.Sprintf("\t%sParsedStatus\tstring\t`%s:\"%s_parsed_status\"`\n",
					schema.ToField(col.Name),
					tagstr,
					col.Name,
				))

			default:
				panic(fmt.Errorf("unknown parse type %d", raw.ColumnsToParse[col.Name]))
			}
		}
	}

	return buf.String()
}

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString(`package psn

// updated at ` + nowPST().String() + `

// Proc represents '/proc' in Linux.
type Proc struct {
	NetStats  NetStats
	Uptime    Uptime
	DiskStats DiskStats
	IO        IO
	Stat      Stat
	Status    Status
}

`)

	// '/proc/net/tcp', '/proc/net/tcp6'
	buf.WriteString(`// NetStat is '/proc/net/tcp', '/proc/net/tcp6' in Linux.
type NetStat struct {
`)
	buf.WriteString(generate(schema.NetStat))
	buf.WriteString("}\n\n")

	// '/proc/uptime'
	buf.WriteString(`// Uptime is '/proc/uptime' in Linux.
type Uptime struct {
`)
	buf.WriteString(generate(schema.Uptime))
	buf.WriteString("}\n\n")

	// '/proc/diskstats'
	buf.WriteString(`// DiskStats is '/proc/diskstats' in Linux.
type DiskStats struct {
`)
	buf.WriteString(generate(schema.DiskStat))
	buf.WriteString("}\n\n")

	// '/proc/$PID/io'
	buf.WriteString(`// IO is '/proc/$PID/io' in Linux.
type IO struct {
`)
	buf.WriteString(generate(schema.IO))
	buf.WriteString("}\n\n")

	// '/proc/$PID/stat'
	buf.WriteString(`// Stat is '/proc/$PID/stat' in Linux.
type Stat struct {
`)
	buf.WriteString(generate(schema.Stat))
	for _, line := range additionalFieldsStat {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	buf.WriteString("}\n\n")

	// '/proc/$PID/status'
	buf.WriteString(`// Status is '/proc/$PID/status' in Linux.
type Status struct {
`)
	buf.WriteString(generate(schema.Status))
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := toFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/generated_linux.go")); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn")); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		log.Fatal(err)
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

var additionalFieldsStat = [...]string{
	"CpuUsage float64 `column:\"cpu_usage\"`",
}
