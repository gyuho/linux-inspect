package ss

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Flags struct {
	LogPath string

	Top    int
	Filter *Process

	Kill    bool
	CleanUp bool

	Monitor         bool
	MonitorInterval time.Duration
}

var (
	Command = &cobra.Command{
		Use:   "ss",
		Short: "Investigates sockets.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &Process{}}
)

func init() {
	Command.PersistentFlags().StringVar(&cmdFlag.LogPath, "log-path", "", "File path to store logs. Empty to print out to stdout.")

	Command.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")

	Command.PersistentFlags().BoolVar(&cmdFlag.Kill, "kill", false, "'true' to kill processes that matches the filter.")
	Command.PersistentFlags().BoolVar(&cmdFlag.CleanUp, "clean-up", false, "'true' to automatically kill zombie processes. Name must be empty.")

	Command.PersistentFlags().BoolVar(&cmdFlag.Monitor, "monitor", false, "'true' to periodically run ps command.")
	Command.PersistentFlags().DurationVar(&cmdFlag.MonitorInterval, "monitor-interval", 10*time.Second, "Monitor interval.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn ss\n\n")
	color.Unset()

	if cmdFlag.Kill && cmdFlag.Monitor {
		fmt.Fprintln(os.Stdout, "can't kill and monitor at the same time!")
		os.Exit(1)
	}

	if cmdFlag.Kill {
		if cmdFlag.CleanUp && cmdFlag.Filter.Program == "" {
			cmdFlag.Filter.Program = "deleted)"
		} else if cmdFlag.Filter.LocalPort == "" && cmdFlag.Filter.Program == "" { // to prevent killing all
			cmdFlag.Filter.Program = "SPECIFY PROGRAM NAME"
		}
	}

	rFunc := func() ([]Process, error) {
		ssr, err := List(cmdFlag.Filter, TCP, TCP6)
		if err != nil {
			return nil, err
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
		WriteToTable(wr, cmdFlag.Top, ssr...)
		return ssr, nil
	}

	ssr, err := rFunc()
	if err != nil {
		return err
	}

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	if cmdFlag.Kill {
		Kill(os.Stdout, ssr...)
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
