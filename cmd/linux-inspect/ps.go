package main

import (
	"fmt"
	"os"

	"github.com/gyuho/linux-inspect/inspect"
	"github.com/gyuho/linux-inspect/top"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type psFlags struct {
	topExecPath string
	limit       int

	program string
	pid     int64
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
	psCommand.PersistentFlags().StringVarP(&psCmdFlag.topExecPath, "top-exec", "t", "", "Specify the top command path.")
	psCommand.PersistentFlags().IntVarP(&psCmdFlag.limit, "limit", "l", 5, "Limit the number results to return.")

	psCommand.PersistentFlags().StringVarP(&psCmdFlag.program, "program", "s", "", "Specify the program name.")
	psCommand.PersistentFlags().Int64VarP(&psCmdFlag.pid, "pid", "p", -1, "Specify the PID.")
}

func psCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ps' to inspect '/proc/$PID/status', 'top' command outpu\n\n")
	color.Unset()

	if psCmdFlag.topExecPath == "" {
		psCmdFlag.topExecPath = top.DefaultExecPath
	}
	pss, err := inspect.GetPS(
		inspect.WithProgram(psCmdFlag.program),
		inspect.WithPID(psCmdFlag.pid),
		inspect.WithTopExecPath(psCmdFlag.topExecPath),
		inspect.WithTopLimit(psCmdFlag.limit),
	)
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
