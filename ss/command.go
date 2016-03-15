package ss

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Flags struct {
	Filter *Process
}

var (
	Command = &cobra.Command{
		Use:   "ss",
		Short: "Investigates sockets.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &Process{}}
)

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.Protocol, "protocol", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.LocalIP, "local-ip", "", "Specify the local IP. Empty lists all local IPs.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemoteIP, "remote-ip", "", "Specify the remote IP. Empty lists all remote IPs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemotePort, "remote-port", "", "Specify the remote port. Empty lists all remote ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.State, "state", "", "Specify the state. Empty lists all states.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.User.Username, "username", "", "Specify the user name. Empty lists all user names.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ss\n\n")
	color.Unset()

	ssr, err := List(cmdFlag.Filter, TCP, TCP6)
	if err != nil {
		return err
	}

	WriteToTable(os.Stdout, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
