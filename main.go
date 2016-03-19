// psn provides utilities to investigate OS processes and sockets.
//
//	Usage:
//	  psn [command]
//
//	Available Commands:
//	  ps          Investigates processes status.
//	  ps-kill     Kills processes.
//	  ps-monitor  Monitors processes.
//	  ss          Investigates sockets.
//	  ss-kill     Kills sockets.
//	  ss-monitor  Monitors sockets.
//
//	Flags:
//	  -h, --help   help for psn
//
//	Use "psn [command] --help" for more information about a command.
//
package main

import (
	"fmt"
	"os"

	"github.com/gyuho/psn/ps"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:        "psn",
		Short:      "psn provides utilities to investigate OS processes and sockets.",
		SuggestFor: []string{"pssn", "psns", "snp"},
	}
)

func init() {
	Command.AddCommand(ps.Command)
	Command.AddCommand(ps.KillCommand)
	Command.AddCommand(ps.MonitorCommand)
	Command.AddCommand(ss.Command)
	Command.AddCommand(ss.KillCommand)
	Command.AddCommand(ss.MonitorCommand)
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := Command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
