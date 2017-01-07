package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn"
	"github.com/spf13/cobra"
)

type nsFlags struct {
}

var (
	nsCommand = &cobra.Command{
		Use:   "ns",
		Short: "Inspects /proc/net/dev",
		RunE:  nsCommandFunc,
	}
	nsCmdFlag = nsFlags{}
)

func init() {
}

func nsCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ds' to inspect '/proc/net/dev'\n\n")
	color.Unset()

	ns, err := psn.GetNS()
	if err != nil {
		return err
	}
	hd, rows := psn.ConvertNS(ns...)
	txt := psn.StringNS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
