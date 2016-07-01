// psn provides utilities to investigate OS processes and sockets.
//
//	Usage:
//	  psn [command]
//
//	Available Commands:
//	  process         Investigates processes status
//	  process-kill    Kills processes
//	  process-monitor Monitors processes
//	  socket          Investigates sockets
//	  socket-kill     Kills sockets
//	  socket-monitor  Monitors sockets
//
//	Flags:
//	  -h, --help[=false]: help for psn
//
//	Use "psn [command] --help" for more information about a command.
//
package main

import (
	"fmt"
	"os"

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
	Command.AddCommand(processCommand)
	Command.AddCommand(processKillCommand)
	Command.AddCommand(processMonitorCommand)
	Command.AddCommand(socketCommand)
	Command.AddCommand(socketKillCommand)
	Command.AddCommand(socketMonitorCommand)
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
