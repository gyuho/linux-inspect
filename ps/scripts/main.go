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

// Stat is metrics in linux proc/$PID/stat.
type Stat struct {
`)
	for i := range ps.Stats {
		fieldName := ps.ToField(ps.Stats[i].Col)
		var typeStr string
		switch ps.Stats[i].Kind {
		case reflect.Uint64:
			typeStr = "uint64"
		case reflect.Int64:
			typeStr = "int64"
		case reflect.String:
			typeStr = "string"
		}
		buf.WriteString(fmt.Sprintf("\t%s\t%s `column:\"%s\"`\n", fieldName, typeStr, ps.Stats[i].Col))
		if ps.Stats[i].Humanize {
			buf.WriteString(fmt.Sprintf("\t%sHumanize\tstring `column:\"%s_humanized\"`\n", fieldName, ps.Stats[i].Col))
		}
	}
	for _, line := range additionalFields {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}

	buf.WriteString("}\n\n")
	txt := buf.String()
	if err := toFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/psn/ps/generated_linux.go")); err != nil {
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
