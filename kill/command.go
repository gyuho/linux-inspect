package kill

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ps"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

type Flags struct {
	Force         bool
	CleanUp       bool
	FilterStatus  *ps.Status
	FilterProcess *ss.Process
}

var (
	Command = &cobra.Command{
		Use:   "kill",
		Short: "Kills programs using syscall. Make sure to specify the flags to find the program.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{FilterStatus: &ps.Status{}, FilterProcess: &ss.Process{}}
)

func init() {
	Command.PersistentFlags().BoolVarP(&cmdFlag.Force, "force", "f", false, "'true' to force kill parent processes.")
	Command.PersistentFlags().BoolVar(&cmdFlag.CleanUp, "clean-up", false, "'true' to clean up deleted programs (this overwrites program flag).")
	Command.PersistentFlags().StringVarP(&cmdFlag.FilterProcess.Program, "program", "s", "SPECIFY YOUR PROGRAM HERE", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().StringVarP(&cmdFlag.FilterProcess.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	if cmdFlag.CleanUp {
		cmdFlag.FilterProcess.Program = "deleted)"
		cmdFlag.FilterStatus.State = "Z (zombie)"
	}
	if cmdFlag.FilterProcess.LocalPort != "" && cmdFlag.FilterProcess.Program == "SPECIFY YOUR PROGRAM HERE" {
		cmdFlag.FilterProcess.Program = ""
	}
	if cmdFlag.FilterProcess.Program != "" && !cmdFlag.CleanUp {
		cmdFlag.FilterStatus.Name = cmdFlag.FilterProcess.Program
	}

	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn kill with %q, %q (port %s)\n\n", cmdFlag.FilterProcess.Program, cmdFlag.FilterStatus.Name, cmdFlag.FilterProcess.LocalPort)
	color.Unset()

	if cmdFlag.FilterStatus.State != "" || cmdFlag.FilterStatus.Name != "" {
		pss, err := ps.ListStatus(cmdFlag.FilterStatus)
		if err != nil {
			return err
		}
		ps.WriteToTable(os.Stdout, 0, pss...)
		ps.Kill(os.Stdout, cmdFlag.Force, pss...)
		fmt.Fprintf(os.Stdout, "\n")
	}

	ssr, err := ss.List(cmdFlag.FilterProcess, ss.TCP, ss.TCP6)
	if err != nil {
		return err
	}
	ss.WriteToTable(os.Stdout, 0, ssr...)
	fmt.Fprintf(os.Stdout, "\n")
	ss.Kill(os.Stdout, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
