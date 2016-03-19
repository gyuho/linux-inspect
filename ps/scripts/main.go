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

	"github.com/gyuho/psn/ps"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString(`package ps

	` + "// updated at " + nowPST().String() + `

// Process represents '/proc' in linux.
type Process struct {
	Stat   Stat
	Status Status
}

`)

	buf.WriteString(`// Stat is 'proc/$PID/stat' in linux.
type Stat struct {
`)

	for i := range ps.StatList {
		fieldName := ps.ToField(ps.StatList[i].Col)
		var typeStr string
		switch ps.StatList[i].Kind {
		case reflect.Uint64:
			typeStr = "uint64"
		case reflect.Int64:
			typeStr = "int64"
		case reflect.String:
			typeStr = "string"
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`column:\"%s\"`\n", fieldName, typeStr, ps.StatList[i].Col))
		if ps.StatList[i].Humanize {
			buf.WriteString(fmt.Sprintf("\t%sHumanize\tstring\t`column:\"%s_humanized\"`\n", fieldName, ps.StatList[i].Col))
		}
	}
	for _, line := range additionalFields {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	buf.WriteString("}\n\n")

	buf.WriteString(`// Status is 'proc/$PID/status' in linux.
type Status struct {
`)
	for i := range ps.StatusListYAML {
		fieldName := ps.ToField(ps.StatusListYAML[i].Col)
		var typeStr string
		switch ps.StatusListYAML[i].Kind {
		case reflect.Uint64:
			typeStr = "uint64"
		case reflect.Int64:
			typeStr = "int64"
		case reflect.String:
			typeStr = "string"
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`yaml:\"%s\"`\n", fieldName, typeStr, ps.StatusListYAML[i].Col))
		if ps.StatusListYAML[i].Bytes {
			buf.WriteString(fmt.Sprintf("\t%sBytes\tuint64\t`yaml:\"%s_bytes\"`\n", fieldName, ps.StatusListYAML[i].Col))
		}
	}
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := toFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/ps/process_generated_linux.go")); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/ps")); err != nil {
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

var additionalFields = [...]string{
	"CpuUsage float64 `column:\"cpu_usage\"`",
}
