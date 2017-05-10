package main

import (
	"fmt"
	"os"

	"github.com/gyuho/linux-inspect/inspect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type psFlags struct {
	program string
	pid     int64
	top     int
}

var (
	psCommand = &cobra.Command{
		Use:   "ps",
		Short: "Inspects '/proc/$PID/status', 'top' command output",
		RunE:  psCommandFunc,
	}
	psCmdFlag psFlags
)

func init() {
	psCommand.PersistentFlags().StringVarP(&psCmdFlag.program, "program", "s", "", "Specify the program name.")
	psCommand.PersistentFlags().Int64VarP(&psCmdFlag.pid, "pid", "p", -1, "Specify the PID.")
	psCommand.PersistentFlags().IntVarP(&psCmdFlag.top, "top", "t", 5, "Limit the number results to return.")
}

func psCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ps' to inspect '/proc/$PID/status', 'top' command outpu\n\n")
	color.Unset()

	pss, err := inspect.GetPS(inspect.WithProgram(psCmdFlag.program), inspect.WithPID(psCmdFlag.pid), inspect.WithTopLimit(psCmdFlag.top))
	if err != nil {
		return err
	}
	hd, rows := inspect.ConvertPS(pss...)
	txt := inspect.StringPS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
