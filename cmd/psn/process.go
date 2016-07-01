package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gyuho/psn/process"
	"github.com/spf13/cobra"
)

type processFlags struct {
	LogPath string

	Filter process.Process
	Top    int

	Kill       bool
	KillParent bool
	CleanUp    bool

	Monitor         bool
	MonitorInterval time.Duration
}

var (
	processCommand = &cobra.Command{
		Use:   "process",
		Short: "Investigates processes status",
		RunE:  processCommandFunc,
	}
	processKillCommand = &cobra.Command{
		Use:   "process-kill",
		Short: "Kills processes",
		RunE:  processKillCommandFunc,
	}
	processMonitorCommand = &cobra.Command{
		Use:   "process-monitor",
		Short: "Monitors processes",
		RunE:  processMonitorCommandFunc,
	}
	processCmdFlag = processFlags{}
)

func init() {
	processCommand.PersistentFlags().StringVarP(&processCmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	processCommand.PersistentFlags().Int64VarP(&processCmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	processCommand.PersistentFlags().IntVarP(&processCmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")

	processKillCommand.PersistentFlags().StringVarP(&processCmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	processKillCommand.PersistentFlags().Int64VarP(&processCmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	processKillCommand.PersistentFlags().BoolVarP(&processCmdFlag.KillParent, "force", "f", false, "'true' to kill processes including its parent processes.")
	processKillCommand.PersistentFlags().BoolVarP(&processCmdFlag.CleanUp, "clean-up", "c", false, "'true' to automatically kill zombie processes. Name must be empty.")

	processMonitorCommand.PersistentFlags().StringVarP(&processCmdFlag.Filter.Stat.Comm, "program", "s", "", "Specify the program. Empty lists all programs.")
	processMonitorCommand.PersistentFlags().Int64VarP(&processCmdFlag.Filter.Stat.Pid, "pid", "i", 0, "Specify the pid. 0 lists all processes.")
	processMonitorCommand.PersistentFlags().IntVarP(&processCmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")
	processMonitorCommand.PersistentFlags().StringVar(&processCmdFlag.LogPath, "log-path", "", "File path to store logs. Empty to print out to stdout. Supports csv file.")
	processMonitorCommand.PersistentFlags().DurationVar(&processCmdFlag.MonitorInterval, "monitor-interval", 10*time.Second, "Monitor interval.")
}

func processCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn process\n\n")
	color.Unset()

	if processCmdFlag.Filter.Stat.Comm != "" {
		processCmdFlag.Filter.Status.Name = processCmdFlag.Filter.Stat.Comm
	}
	if processCmdFlag.Filter.Stat.Pid != 0 {
		processCmdFlag.Filter.Status.Pid = processCmdFlag.Filter.Stat.Pid
	}

	pss, err := process.List(&processCmdFlag.Filter)
	if err != nil {
		return err
	}
	process.WriteToTable(os.Stdout, processCmdFlag.Top, pss...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func processKillCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn process-kill\n\n")
	color.Unset()

	if processCmdFlag.Filter.Stat.Comm != "" {
		processCmdFlag.Filter.Status.Name = processCmdFlag.Filter.Stat.Comm
	}
	if processCmdFlag.Filter.Stat.Pid != 0 {
		processCmdFlag.Filter.Status.Pid = processCmdFlag.Filter.Stat.Pid
	}
	if processCmdFlag.CleanUp && processCmdFlag.Filter.Stat.Comm == "" {
		processCmdFlag.Filter.Stat.State = "Z"
		processCmdFlag.Filter.Status.State = "Z (zombie)"
	}

	pss, err := process.List(&processCmdFlag.Filter)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nerror (%v)\n", err)
		return nil
	}
	process.WriteToTable(os.Stdout, processCmdFlag.Top, pss...)
	process.Kill(os.Stdout, processCmdFlag.KillParent, pss...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func processMonitorCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgBlue)
	fmt.Fprintf(os.Stdout, "\npsn process-monitor\n\n")
	color.Unset()

	if processCmdFlag.Filter.Stat.Comm != "" {
		processCmdFlag.Filter.Status.Name = processCmdFlag.Filter.Stat.Comm
	}
	if processCmdFlag.Filter.Stat.Pid != 0 {
		processCmdFlag.Filter.Status.Pid = processCmdFlag.Filter.Stat.Pid
	}
	if processCmdFlag.CleanUp && processCmdFlag.Filter.Stat.Comm == "" {
		processCmdFlag.Filter.Status.State = "Z (zombie)"
	}

	rFunc := func() error {
		pss, err := process.List(&processCmdFlag.Filter)
		if err != nil {
			return err
		}

		if filepath.Ext(processCmdFlag.LogPath) == ".csv" {
			f, err := openToAppend(processCmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := process.WriteToCSV(f, pss...); err != nil {
				return err
			}
			return err
		}

		var wr io.Writer
		if processCmdFlag.LogPath == "" {
			wr = os.Stdout
		} else {
			f, err := openToAppend(processCmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			wr = f
		}
		process.WriteToTable(wr, processCmdFlag.Top, pss...)
		return nil
	}

	if err := rFunc(); err != nil {
		return err
	}

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

escape:
	for {
		select {
		case <-time.After(processCmdFlag.MonitorInterval):
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
