// psn is an utility to investigate sockets.
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

const (
	cliName        = "psn"
	cliDescription = "psn provides utilities to investigate OS processes and sockets."
)

type GlobalFlag struct {
	GlobalFilter *ss.Process
}

var (
	w          = os.Stdout
	globalFlag = GlobalFlag{GlobalFilter: &ss.Process{}}

	rootCommand = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"sssn", "sns", "sn"},
		RunE:       rootCommandFunc,
	}
	ssCommand = &cobra.Command{
		Use:   "ss",
		Short: "ss investigates sockets.",
		RunE:  ssCommandFunc,
	}
	killCommand = &cobra.Command{
		Use:   "kill",
		Short: "kill kills programs using syscall. Make sure to specify the flags to find the program.",
		RunE:  killCommandFunc,
	}
)

func rootCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgBlue)
	fmt.Fprintf(w, "\npsn is listing all ps and ss data:\n\n")
	color.Unset()

	// TODO:
	// psr, err := ps.List(...)

	ssr, err := ss.List(globalFlag.GlobalFilter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}

	ss.WriteToTable(w, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(w, "\nDone.\n")
	color.Unset()

	return nil
}

func ssCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(w, "\npsn ss is listing:\n\n")
	color.Unset()

	ssr, err := ss.List(globalFlag.GlobalFilter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}

	ss.WriteToTable(w, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(w, "\nDone.\n")
	color.Unset()

	return nil
}

func killCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgRed)
	fmt.Fprintf(w, "\npsn is killing:\n\n")
	color.Unset()

	// TODO:
	// psr, err := ps.List(...)
	// ps.WriteToTable(...)
	// ps.Kill(...)

	ssr, err := ss.List(globalFlag.GlobalFilter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}

	ss.WriteToTable(w, ssr...)
	fmt.Fprintf(w, "\n")

	ss.Kill(w, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(w, "\nDone.\n")
	color.Unset()

	return nil
}

func init() {
	rootCommand.AddCommand(ssCommand)
	rootCommand.AddCommand(killCommand)

	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.Protocol, "protocol", "t", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.LocalIP, "local-ip", "l", "", "Specify the local IP. Empty lists all local IPs.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.LocalPort, "local-port", "p", "", "Specify the local port. Empty lists all local ports.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.RemoteIP, "remote-ip", "r", "", "Specify the remote IP. Empty lists all remote IPs.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.RemotePort, "remote-port", "m", "", "Specify the remote port. Empty lists all remote ports.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.State, "state", "a", "", "Specify the state. Empty lists all states.")
	rootCommand.PersistentFlags().StringVarP(&globalFlag.GlobalFilter.User.Username, "username", "u", "", "Specify the user name. Empty lists all user names.")
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(w, err)
		os.Exit(1)
	}
}
