package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/linux-inspect/psn"
	"github.com/spf13/cobra"
)

var (
	nsCommand = &cobra.Command{
		Use:   "ns",
		Short: "Inspects '/proc/net/dev'",
		RunE:  nsCommandFunc,
	}
)

func nsCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ns' to inspect '/proc/net/dev'\n\n")
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
