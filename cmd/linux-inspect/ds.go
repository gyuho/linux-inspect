package main

import (
	"fmt"
	"os"

	"github.com/gyuho/linux-inspect/inspect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	dsCommand = &cobra.Command{
		Use:   "ds",
		Short: "Inspects '/proc/diskstats'",
		RunE:  dsCommandFunc,
	}
)

func dsCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ds' to inspect '/proc/diskstats'\n\n")
	color.Unset()

	ds, err := inspect.GetDS()
	if err != nil {
		return err
	}
	hd, rows := inspect.ConvertDS(ds...)
	txt := inspect.StringDS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
