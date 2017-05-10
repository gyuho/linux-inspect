package main

import (
	"fmt"
	"os"

	"github.com/gyuho/linux-inspect/inspect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ssFlags struct {
	topExecPath string
	limit       int

	program   string
	protocol  string
	localPort int64
}

var (
	ssCommand = &cobra.Command{
		Use:   "ss",
		Short: "Inspects '/proc/net/tcp,tcp6'",
		RunE:  ssCommandFunc,
	}
	ssCmdFlag = ssFlags{}
)

func init() {
	ssCommand.PersistentFlags().StringVarP(&ssCmdFlag.topExecPath, "top-exec", "t", "", "Specify the top command path.")
	ssCommand.PersistentFlags().IntVarP(&ssCmdFlag.limit, "limit", "l", 5, "Limit the number results to return.")

	ssCommand.PersistentFlags().StringVarP(&ssCmdFlag.protocol, "protocol", "c", "tcp", "Specify the protocol ('tcp' or 'tcp6').")
	ssCommand.PersistentFlags().StringVarP(&ssCmdFlag.program, "program", "s", "", "Specify the program name.")
	ssCommand.PersistentFlags().Int64VarP(&ssCmdFlag.localPort, "local-port", "l", -1, "Specify the local port.")
}

func ssCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ss' to inspect '/proc/net/tcp,tcp6'\n\n")
	color.Unset()

	topt := inspect.WithTCP()
	if ssCmdFlag.protocol == "tcp6" {
		topt = inspect.WithTCP6()
	} else if ssCmdFlag.protocol != "tcp" {
		fmt.Fprintf(os.Stderr, "unknown protocol %q\n", ssCmdFlag.protocol)
		os.Exit(233)
	}
	sss, err := inspect.GetSS(
		topt,
		inspect.WithTopExecPath(ssCmdFlag.topExecPath),
		inspect.WithTopLimit(ssCmdFlag.limit),
		inspect.WithProgram(ssCmdFlag.program),
		inspect.WithLocalPort(ssCmdFlag.localPort),
	)
	if err != nil {
		return err
	}
	hd, rows := inspect.ConvertSS(sss...)
	txt := inspect.StringSS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
