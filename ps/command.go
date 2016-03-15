package ps

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Flags struct {
	Filter *Status
}

var (
	Command = &cobra.Command{
		Use:   "ps",
		Short: "Investigates processes.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &Status{}}
)

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Name, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Filter.Pid, "pid", "p", 0, "Specify the pid. 0 lists all processes.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ps\n\n")
	color.Unset()

	pss, err := ListStatus(cmdFlag.Filter)
	if err != nil {
		return err
	}

	for _, p := range pss {
		fmt.Fprintln(os.Stdout, p)
	}

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
