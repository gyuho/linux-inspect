// generate-proc generates proc struct based on the schema.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/pkg/schemautil"
	"github.com/gyuho/linux-inspect/pkg/timeutil"
	"github.com/gyuho/linux-inspect/proc/schema"
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
	buf.WriteString(`package proc

// updated at ` + timeutil.NowPST().String() + `

`)

	// '/proc/net/dev'
	buf.WriteString(`// NetDev is '/proc/net/dev' in Linux.
// The dev pseudo-file contains network device status information.
type NetDev struct {
`)
	buf.WriteString(schemautil.Generate(schema.NetDev))
	buf.WriteString("}\n\n")

	// '/proc/net/tcp', '/proc/net/tcp6'
	buf.WriteString(`// NetTCP is '/proc/net/tcp', '/proc/net/tcp6' in Linux.
// Holds a dump of the TCP socket table.
type NetTCP struct {
`)
	for _, line := range additionalFieldsNetTCP {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	buf.WriteString(schemautil.Generate(schema.NetTCP))
	buf.WriteString("}\n\n")

	// '/proc/loadavg'
	buf.WriteString(`// LoadAvg is '/proc/loadavg' in Linux.
type LoadAvg struct {
`)
	buf.WriteString(schemautil.Generate(schema.LoadAvg))
	buf.WriteString("}\n\n")

	// '/proc/uptime'
	buf.WriteString(`// Uptime is '/proc/uptime' in Linux.
type Uptime struct {
`)
	buf.WriteString(schemautil.Generate(schema.Uptime))
	buf.WriteString("}\n\n")

	// '/proc/diskstats'
	buf.WriteString(`// DiskStat is '/proc/diskstats' in Linux.
type DiskStat struct {
`)
	buf.WriteString(schemautil.Generate(schema.DiskStat))
	buf.WriteString("}\n\n")

	// '/proc/$PID/io'
	buf.WriteString(`// IO is '/proc/$PID/io' in Linux.
type IO struct {
`)
	buf.WriteString(schemautil.Generate(schema.IO))
	buf.WriteString("}\n\n")

	// '/proc/$PID/stat'
	buf.WriteString(`// Stat is '/proc/$PID/stat' in Linux.
type Stat struct {
`)
	buf.WriteString(schemautil.Generate(schema.Stat))
	buf.WriteString("}\n\n")

	// '/proc/$PID/status'
	buf.WriteString(`// Status is '/proc/$PID/status' in Linux.
type Status struct {
`)
	buf.WriteString(schemautil.Generate(schema.Status))
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := fileutil.ToFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/proc/generated.go")); err != nil {
		panic(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/proc")); err != nil {
		panic(err)
	}
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		panic(err)
	}

	fmt.Println("DONE")
}

var additionalFieldsNetTCP = [...]string{
	"Type string `column:\"type\"`",
}
