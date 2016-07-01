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

	"github.com/gyuho/psn/process"
)

func main() {
	buf := new(bytes.Buffer)
	buf.WriteString(`package process

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

	for i := range process.StatList {
		fieldName := process.ToField(process.StatList[i].Col)
		var typeStr string
		switch process.StatList[i].Kind {
		case reflect.Uint64:
			typeStr = "uint64"
		case reflect.Int64:
			typeStr = "int64"
		case reflect.String:
			typeStr = "string"
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`column:\"%s\"`\n", fieldName, typeStr, process.StatList[i].Col))
		if process.StatList[i].Humanize {
			buf.WriteString(fmt.Sprintf("\t%sHumanize\tstring\t`column:\"%s_humanized\"`\n", fieldName, process.StatList[i].Col))
		}
	}
	for _, line := range additionalFields {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	buf.WriteString("}\n\n")

	buf.WriteString(`// Status is 'proc/$PID/status' in linux.
type Status struct {
`)
	for i := range process.StatusListYAML {
		fieldName := process.ToField(process.StatusListYAML[i].Col)
		var typeStr string
		switch process.StatusListYAML[i].Kind {
		case reflect.Uint64:
			typeStr = "uint64"
		case reflect.Int64:
			typeStr = "int64"
		case reflect.String:
			typeStr = "string"
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s\t`yaml:\"%s\"`\n", fieldName, typeStr, process.StatusListYAML[i].Col))
		if process.StatusListYAML[i].Bytes {
			buf.WriteString(fmt.Sprintf("\t%sBytes\tuint64\t`yaml:\"%s_bytes\"`\n", fieldName, process.StatusListYAML[i].Col))
		}
	}
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := toFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/process/process_generated_linux.go")); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/process")); err != nil {
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
