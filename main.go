// psn provides utilities to investigate OS processes and sockets.
//
//	Usage:
//	  psn [command]
//
//	Available Commands:
//	  status      Investigates processes status.
//	  ss          Investigates sockets.
//	  kill        Kills programs using syscall. Make sure to specify the flags to find the program.
//	  monitor     Monitors programs.
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

	"github.com/gyuho/psn/kill"
	"github.com/gyuho/psn/monitor"
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
	Command.AddCommand(ss.Command)
	Command.AddCommand(kill.Command)
	Command.AddCommand(monitor.Command)
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
