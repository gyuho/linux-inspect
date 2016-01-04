// psn is an utility to investigate sockets.
package main

import (
	"fmt"
	"os"

	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

const (
	cliName        = "psn"
	cliDescription = "psn provides utilities to investigate OS processes and sockets."
)

type flags struct {
}

var (
	w            = os.Stdout
	globalFlags  = flags{}
	globalFilter = &ss.Process{}

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"sssn", "sns", "sn"},
		RunE:       rootCommandFunc,
	}
	killCmd = &cobra.Command{
		Use:   "kill",
		Short: "kill kills programs using syscall. Make sure to specify the flags to find the program.",
		RunE:  killCommandFunc,
	}
)

func rootCommandFunc(cmd *cobra.Command, args []string) error {
	ps, err := ss.List(globalFilter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}
	ss.WriteToTable(w, ps...)
	return nil
}

func killCommandFunc(cmd *cobra.Command, args []string) error {
	ps, err := ss.List(globalFilter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "Killing the following processes...")
	ss.WriteToTable(w, ps...)
	ss.Kill(w, ps...)
	return nil
}

func init() {
	rootCmd.AddCommand(killCmd)

	rootCmd.PersistentFlags().StringVarP(&globalFilter.Protocol, "protocol", "t", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.LocalIP, "local-ip", "l", "", "Specify the local IP. Empty lists all local IPs.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.LocalPort, "local-port", "p", "", "Specify the local port. Empty lists all local ports.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.RemoteIP, "remote-ip", "r", "", "Specify the remote IP. Empty lists all remote IPs.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.RemotePort, "remote-port", "m", "", "Specify the remote port. Empty lists all remote ports.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.State, "state", "a", "", "Specify the state. Empty lists all states.")
	rootCmd.PersistentFlags().StringVarP(&globalFilter.User.Username, "username", "u", "", "Specify the user name. Empty lists all user names.")
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(w, err)
		os.Exit(1)
	}
}
