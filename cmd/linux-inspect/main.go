// linux-inspect inspects Linux processes, sockets (ps, ss, netstat).
//
//	Usage:
//	linux-inspect [command]
//
//	Available Commands:
//	ds          Inspects '/proc/diskstats'
//	ns          Inspects '/proc/net/dev'
//	ps          Inspects '/proc/$PID/stat,status'
//	ss          Inspects '/proc/net/tcp,tcp6'
//
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	command = &cobra.Command{
		Use:        "linux-inspect",
		Short:      "linux-inspect inspects Linux processes, sockets (ps, ss, netstat).",
		SuggestFor: []string{"linux-inspects", "linuxinspect", "linux-inspec"},
	}
)

func init() {
	command.AddCommand(dsCommand)
	command.AddCommand(nsCommand)
	command.AddCommand(psCommand)
	command.AddCommand(ssCommand)
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
