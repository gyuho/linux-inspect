package kill

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

type Flags struct {
	Force   bool
	CleanUp bool
	Filter  *ss.Process
}

var (
	Command = &cobra.Command{
		Use:   "kill",
		Short: "kill kills programs using syscall. Make sure to specify the flags to find the program.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &ss.Process{}}
)

func init() {
	Command.PersistentFlags().BoolVarP(&cmdFlag.Force, "force", "f", false, "'true' to force kill any programs.")
	Command.PersistentFlags().BoolVar(&cmdFlag.CleanUp, "clean-up", false, "'true' to clean up deleted programs (this overwrites program flag).")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Program, "program", "s", "SPECIFY YOUR PROGRAM HERE", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.Protocol, "protocol", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.LocalIP, "local-ip", "", "Specify the local IP. Empty lists all local IPs.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemoteIP, "remote-ip", "", "Specify the remote IP. Empty lists all remote IPs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemotePort, "remote-port", "", "Specify the remote port. Empty lists all remote ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.State, "state", "", "Specify the state. Empty lists all states.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.User.Username, "username", "", "Specify the user name. Empty lists all user names.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	if cmdFlag.CleanUp {
		cmdFlag.Filter.Program = "deleted)"
	}
	if cmdFlag.Filter.LocalPort != "" && cmdFlag.Filter.Program == "SPECIFY YOUR PROGRAM HERE" {
		cmdFlag.Filter.Program = ""
	}

	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn kill with %q (port %s)\n\n", cmdFlag.Filter.Program, cmdFlag.Filter.LocalPort)
	color.Unset()

	ssr, err := ss.List(cmdFlag.Filter, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}

	ss.WriteToTable(os.Stdout, ssr...)
	fmt.Fprintf(os.Stdout, "\n")

	ss.Kill(os.Stdout, cmdFlag.Force, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
