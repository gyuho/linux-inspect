package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ps"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

type Flags struct {
	LogPath       string
	Interval      time.Duration
	Top           int
	FilterStatus  *ps.Status
	FilterProcess *ss.Process
}

var (
	Command = &cobra.Command{
		Use:   "monitor",
		Short: "Monitors programs.",
	}
	PsCommand = &cobra.Command{
		Use:   "ps",
		Short: "Monitors using 'ps'.",
		RunE:  PsCommandFunc,
	}
	SsCommand = &cobra.Command{
		Use:   "ss",
		Short: "Monitors using 'ss'.",
		RunE:  SsCommandFunc,
	}
	cmdFlag = Flags{FilterStatus: &ps.Status{}, FilterProcess: &ss.Process{}}
)

func init() {
	Command.PersistentFlags().StringVarP(&cmdFlag.LogPath, "log-path", "f", "", "Specify the log path to write the monitor results. Extend with .csv to record in comma separate format.")
	Command.PersistentFlags().DurationVarP(&cmdFlag.Interval, "interval", "i", 10*time.Second, "Interval to repeat monitor scans.")
	Command.PersistentFlags().IntVarP(&cmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")

	Command.AddCommand(PsCommand)
	PsCommand.PersistentFlags().StringVarP(&cmdFlag.FilterStatus.Name, "program", "s", "", "Specify the program. Empty lists all programs.")
	PsCommand.PersistentFlags().IntVarP(&cmdFlag.FilterStatus.Pid, "pid", "p", 0, "Specify the pid. 0 lists all processes.")

	Command.AddCommand(SsCommand)
	SsCommand.PersistentFlags().StringVarP(&cmdFlag.FilterProcess.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.Protocol, "protocol", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.LocalIP, "local-ip", "", "Specify the local IP. Empty lists all local IPs.")
	SsCommand.PersistentFlags().StringVarP(&cmdFlag.FilterProcess.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.RemoteIP, "remote-ip", "", "Specify the remote IP. Empty lists all remote IPs.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.RemotePort, "remote-port", "", "Specify the remote port. Empty lists all remote ports.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.State, "state", "", "Specify the state. Empty lists all states.")
	SsCommand.PersistentFlags().StringVar(&cmdFlag.FilterProcess.User.Username, "username", "", "Specify the user name. Empty lists all user names.")
}

func PsCommandFunc(cmd *cobra.Command, args []string) error {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	rFunc := func() error {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "\npsn ps monitor at %s\n\n", time.Now())
		color.Unset()

		pss, err := ps.ListStatus(cmdFlag.FilterStatus)
		if err != nil {
			return err
		}

		ps.WriteToTable(os.Stdout, cmdFlag.Top, pss...)
		fmt.Fprintf(os.Stdout, "\n")

		return nil
	}
	firstCSV := true
	if cmdFlag.LogPath != "" {
		rFunc = func() error {
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()

			color.Set(color.FgRed)
			fmt.Fprintf(f, "\npsn ps monitor at %s\n\n", time.Now())
			color.Unset()

			pss, err := ps.ListStatus(cmdFlag.FilterStatus)
			if err != nil {
				return err
			}

			ps.WriteToTable(f, cmdFlag.Top, pss...)
			fmt.Fprintf(f, "\n")

			color.Set(color.FgGreen)
			fmt.Fprintf(f, "\nDone.\n")
			color.Unset()

			return nil
		}
		if filepath.Ext(cmdFlag.LogPath) == ".csv" {
			rFunc = func() error {
				pss, err := ps.ListStatus(cmdFlag.FilterStatus)
				if err != nil {
					return err
				}

				f, err := openToAppend(cmdFlag.LogPath)
				if err != nil {
					return err
				}
				defer f.Close()

				if err := ps.WriteToCSV(firstCSV, f, pss...); err != nil {
					return err
				}
				firstCSV = false
				return nil
			}
		}
	}

	fmt.Fprintf(os.Stdout, "Monitoring file saved at %v\n", cmdFlag.LogPath)
	var err error
	if err = rFunc(); err != nil {
		return err
	}

escape:
	for {
		fmt.Fprintf(os.Stdout, "Running 'psn monitor ps' at %v\n", time.Now())
		select {
		case <-time.After(cmdFlag.Interval):
			if err = rFunc(); err != nil {
				fmt.Fprintf(os.Stdout, "error: %v\n", err)
				break escape
			}
		case sig := <-notifier:
			fmt.Fprintf(os.Stdout, "Received %v\n", sig)
			return nil
		}
	}

	return err
}

func SsCommandFunc(cmd *cobra.Command, args []string) error {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	rFunc := func() error {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "\npsn ss monitor at %s\n\n", time.Now())
		color.Unset()

		ssr, err := ss.List(cmdFlag.FilterProcess, ss.TCP, ss.TCP6)
		if err != nil {
			return err
		}

		ss.WriteToTable(os.Stdout, cmdFlag.Top, ssr...)
		fmt.Fprintf(os.Stdout, "\n")

		return nil
	}
	firstCSV := true
	if cmdFlag.LogPath != "" {
		rFunc = func() error {
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()

			color.Set(color.FgRed)
			fmt.Fprintf(f, "\npsn ss monitor at %s\n\n", time.Now())
			color.Unset()

			ssr, err := ss.List(cmdFlag.FilterProcess, ss.TCP, ss.TCP6)
			if err != nil {
				return err
			}

			ss.WriteToTable(f, cmdFlag.Top, ssr...)
			fmt.Fprintf(f, "\n")

			color.Set(color.FgGreen)
			fmt.Fprintf(f, "\nDone.\n")
			color.Unset()

			return nil
		}
		if filepath.Ext(cmdFlag.LogPath) == ".csv" {
			rFunc = func() error {
				ssr, err := ss.List(cmdFlag.FilterProcess, ss.TCP, ss.TCP6)
				if err != nil {
					return err
				}

				f, err := openToAppend(cmdFlag.LogPath)
				if err != nil {
					return err
				}
				defer f.Close()

				if err := ss.WriteToCSV(firstCSV, f, ssr...); err != nil {
					return err
				}
				firstCSV = false
				return nil
			}
		}
	}

	var err error
	if err = rFunc(); err != nil {
		return err
	}

escape:
	for {
		fmt.Fprintf(os.Stdout, "Running 'psn monitor ss' at %v\n", time.Now())
		select {
		case <-time.After(cmdFlag.Interval):
			if err = rFunc(); err != nil {
				break escape
			}
		case sig := <-notifier:
			fmt.Fprintf(os.Stdout, "Received %v\n", sig)
			return nil
		}
	}

	return err
}
