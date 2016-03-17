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

	Filter *Status
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
	cmdFlag = Flags{Filter: &Status{}}
)

func init() {
	Command.PersistentFlags().StringVar(&cmdFlag.LogPath, "log-path", "", "File path to store logs. Empty to print out to stdout. Supports csv file.")

	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Name, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Filter.Pid, "pid", "p", 0, "Specify the pid. 0 lists all processes.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")

	Command.PersistentFlags().BoolVar(&cmdFlag.Kill, "kill", false, "'true' to kill processes that matches the filter.")
	Command.PersistentFlags().BoolVar(&cmdFlag.KillParent, "kill-parent", false, "'true' to kill processes including its parent processes.")
	Command.PersistentFlags().BoolVar(&cmdFlag.CleanUp, "clean-up", false, "'true' to automatically kill zombie processes. Name must be empty.")

	Command.PersistentFlags().BoolVar(&cmdFlag.Monitor, "monitor", false, "'true' to periodically run ps command.")
	Command.PersistentFlags().DurationVar(&cmdFlag.MonitorInterval, "monitor-interval", 10*time.Second, "Monitor interval.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ps\n\n")
	color.Unset()

	if cmdFlag.Kill && cmdFlag.Monitor {
		fmt.Fprintln(os.Stdout, "can't kill and monitor at the same time!")
		os.Exit(1)
	}

	if cmdFlag.Kill {
		if cmdFlag.CleanUp && cmdFlag.Filter.Name == "" && cmdFlag.Filter.State == "" {
			cmdFlag.Filter.State = "Z (zombie)"
		}
	}

	rFunc := func() ([]Status, error) {
		pss, err := List(cmdFlag.Filter)
		if err != nil {
			return nil, err
		}
		if filepath.Ext(cmdFlag.LogPath) == ".csv" {
			fmt.Fprintf(os.Stdout, "File saved at %s\n", cmdFlag.LogPath)
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			if err := WriteToCSV(f, pss...); err != nil {
				return nil, err
			}
			return pss, err
		}
		var wr io.Writer
		if cmdFlag.LogPath == "" {
			wr = os.Stdout
		} else {
			fmt.Fprintf(os.Stdout, "File saved at %s\n", cmdFlag.LogPath)
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			wr = f
		}
		WriteToTable(wr, cmdFlag.Top, pss...)
		return pss, nil
	}

	pss, err := rFunc()
	if err != nil {
		return err
	}

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	if cmdFlag.Kill {
		Kill(os.Stdout, cmdFlag.KillParent, pss...)
	} else if cmdFlag.Monitor {
	escape:
		for {
			select {
			case <-time.After(cmdFlag.MonitorInterval):
				if _, err = rFunc(); err != nil {
					fmt.Fprintf(os.Stdout, "error: %v\n", err)
					break escape
				}
			case sig := <-notifier:
				fmt.Fprintf(os.Stdout, "Received %v\n", sig)
				return nil
			}
		}
	}

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}
