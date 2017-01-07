package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gyuho/psn"
	"github.com/spf13/cobra"
)

type dsFlags struct {
}

var (
	dsCommand = &cobra.Command{
		Use:   "ds",
		Short: "Inspects disk statistics",
		RunE:  dsCommandFunc,
	}
	dsCmdFlag = dsFlags{}
)

func init() {
}

func dsCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ds' to inspect disk statistics\n\n")
	color.Unset()

	ds, err := psn.GetDS()
	if err != nil {
		return err
	}
	hd, rows := psn.ConvertDS(ds...)
	txt := psn.StringDS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
