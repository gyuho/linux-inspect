package ps

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Flags struct {
	LogPath string

	Filter Process
	Top    int

	Kill       bool
	KillParent bool
	CleanUp    bool

	Monitor         bool
	MonitorInterval time.Duration
}

var (
	Command = &cobra.Command{
		Use:   "ps",
		Short: "Investigates processes status.",
		RunE:  CommandFunc,
	}
	KillCommand = &cobra.Command{
		Use:   "ps-kill",
		Short: "Kills processes.",
		RunE:  KillCommandFunc,
	}
	MonitorCommand = &cobra.Command{
		Use:   "ps-monitor",
		Short: "Monitors processes.",
		RunE:  MonitorCommandFunc,
	}
	cmdFlag = Flags{}
)

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().Int64VarP(&cmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")

	KillCommand.PersistentFlags().StringVarP(&cmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	KillCommand.PersistentFlags().Int64VarP(&cmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	KillCommand.PersistentFlags().BoolVarP(&cmdFlag.KillParent, "force", "f", false, "'true' to kill processes including its parent processes.")
	KillCommand.PersistentFlags().BoolVarP(&cmdFlag.CleanUp, "clean-up", "c", false, "'true' to automatically kill zombie processes. Name must be empty.")

	MonitorCommand.PersistentFlags().StringVarP(&cmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	MonitorCommand.PersistentFlags().Int64VarP(&cmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	MonitorCommand.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")
	MonitorCommand.PersistentFlags().StringVar(&cmdFlag.LogPath, "log-path", "", "File path to store logs. Empty to print out to stdout. Supports csv file.")
	MonitorCommand.PersistentFlags().DurationVar(&cmdFlag.MonitorInterval, "monitor-interval", 10*time.Second, "Monitor interval.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ps\n\n")
	color.Unset()

	if cmdFlag.Filter.Stat.Comm != "" {
		cmdFlag.Filter.Status.Name = cmdFlag.Filter.Stat.Comm
	}
	if cmdFlag.Filter.Stat.Pid != 0 {
		cmdFlag.Filter.Status.Pid = cmdFlag.Filter.Stat.Pid
	}

	pss, err := List(&cmdFlag.Filter)
	if err != nil {
		return err
	}
	WriteToTable(os.Stdout, cmdFlag.Top, pss...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func KillCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn ps-kill\n\n")
	color.Unset()

	if cmdFlag.Filter.Stat.Comm != "" {
		cmdFlag.Filter.Status.Name = cmdFlag.Filter.Stat.Comm
	}
	if cmdFlag.Filter.Stat.Pid != 0 {
		cmdFlag.Filter.Status.Pid = cmdFlag.Filter.Stat.Pid
	}
	if cmdFlag.CleanUp && cmdFlag.Filter.Stat.Comm == "" {
		cmdFlag.Filter.Stat.State = "Z"
		cmdFlag.Filter.Status.State = "Z (zombie)"
	}

	pss, err := List(&cmdFlag.Filter)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nerror (%v)\n", err)
		return nil
	}
	WriteToTable(os.Stdout, cmdFlag.Top, pss...)
	Kill(os.Stdout, cmdFlag.KillParent, pss...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func MonitorCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgBlue)
	fmt.Fprintf(os.Stdout, "\npsn ps-monitor\n\n")
	color.Unset()

	if cmdFlag.Filter.Stat.Comm != "" {
		cmdFlag.Filter.Status.Name = cmdFlag.Filter.Stat.Comm
	}
	if cmdFlag.Filter.Stat.Pid != 0 {
		cmdFlag.Filter.Status.Pid = cmdFlag.Filter.Stat.Pid
	}
	if cmdFlag.CleanUp && cmdFlag.Filter.Stat.Comm == "" {
		cmdFlag.Filter.Status.State = "Z (zombie)"
	}

	rFunc := func() error {
		pss, err := List(&cmdFlag.Filter)
		if err != nil {
			return err
		}

		if filepath.Ext(cmdFlag.LogPath) == ".csv" {
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := WriteToCSV(f, pss...); err != nil {
				return err
			}
			return err
		}

		var wr io.Writer
		if cmdFlag.LogPath == "" {
			wr = os.Stdout
		} else {
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			wr = f
		}
		WriteToTable(wr, cmdFlag.Top, pss...)
		return nil
	}

	if err := rFunc(); err != nil {
		return err
	}

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)
escape:
	for {
		select {
		case <-time.After(cmdFlag.MonitorInterval):
			if err := rFunc(); err != nil {
				fmt.Fprintf(os.Stdout, "error: %v\n", err)
				break escape
			}
		case sig := <-notifier:
			fmt.Fprintf(os.Stdout, "Received %v\n", sig)
			return nil
		}
	}

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
