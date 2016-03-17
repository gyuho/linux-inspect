package ps

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Flags struct {
	Filter   *Status
	Detailed bool
	Top      int
}

var (
	Command = &cobra.Command{
		Use:   "status",
		Short: "Investigates processes status.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &Status{}}
)

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Name, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Filter.Pid, "pid", "p", 0, "Specify the pid. 0 lists all processes.")
	Command.PersistentFlags().BoolVar(&cmdFlag.Detailed, "detailed", false, "'true' to print out detailed process information.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ps\n\n")
	color.Unset()

	pss, err := ListStatus(cmdFlag.Filter)
	if err != nil {
		return err
	}

	if cmdFlag.Detailed {
		for _, p := range pss {
			fmt.Fprintln(os.Stdout, p.StringDetailed())
		}
	} else {
		WriteToTable(os.Stdout, cmdFlag.Top, pss...)
	}

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
