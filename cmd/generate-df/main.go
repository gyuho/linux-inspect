// generate-psn-schema generates psn struct based on the schema.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gyuho/linux-inspect/df/schema"
	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/pkg/schemautil"
	"github.com/gyuho/linux-inspect/pkg/timeutil"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	exp := filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect")
	if wd != exp {
		panic(fmt.Errorf("must be run in repo root %q, but run at %q", exp, wd))
	}

	buf := new(bytes.Buffer)
	buf.WriteString(`package df

// updated at ` + timeutil.NowPST().String() + `

`)

	// 'df'
	buf.WriteString(`// DfCommandRow is 'df' command output row in Linux.
type DfCommandRow struct {
`)
	buf.WriteString(schemautil.Generate(schema.DfCommandRow))
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := fileutil.ToFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/df/generated.go")); err != nil {
		panic(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/df")); err != nil {
		panic(err)
	}
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		panic(err)
	}

	fmt.Println("DONE")
}
