package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gyuho/psn/ss"
	"github.com/spf13/cobra"
)

type Flags struct {
	LogPath  string
	Interval time.Duration
	Filter   *ss.Process
}

var (
	Command = &cobra.Command{
		Use:   "monitor",
		Short: "monitor monitors programs.",
		RunE:  CommandFunc,
	}
	cmdFlag = Flags{Filter: &ss.Process{}}
)

func init() {
	Command.PersistentFlags().StringVar(&cmdFlag.LogPath, "log-path", "", "Specify the log path to write the monitor results.")
	Command.PersistentFlags().DurationVar(&cmdFlag.Interval, "interval", 10*time.Second, "Interval to repeat monitor scans.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.Protocol, "protocol", "", "'tcp' or 'tcp6'. Empty lists all protocols.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.LocalIP, "local-ip", "", "Specify the local IP. Empty lists all local IPs.")
	Command.PersistentFlags().StringVarP(&cmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemoteIP, "remote-ip", "", "Specify the remote IP. Empty lists all remote IPs.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.RemotePort, "remote-port", "", "Specify the remote port. Empty lists all remote ports.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.State, "state", "", "Specify the state. Empty lists all states.")
	Command.PersistentFlags().StringVar(&cmdFlag.Filter.User.Username, "username", "", "Specify the user name. Empty lists all user names.")
}

func CommandFunc(cmd *cobra.Command, args []string) error {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	rFunc := func() error {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "\npsn monitor at %s\n\n", time.Now())
		color.Unset()

		ssr, err := ss.List(cmdFlag.Filter, ss.TCP, ss.TCP6)
		if err != nil {
			return err
		}

		ss.WriteToTable(os.Stdout, ssr...)
		fmt.Fprintf(os.Stdout, "\n")

		return nil
	}
	if cmdFlag.LogPath != "" {
		rFunc = func() error {
			f, err := openToAppend(cmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()

			color.Set(color.FgRed)
			fmt.Fprintf(f, "\npsn monitor at %s\n\n", time.Now())
			color.Unset()

			ssr, err := ss.List(cmdFlag.Filter, ss.TCP, ss.TCP6)
			if err != nil {
				return err
			}

			ss.WriteToTable(f, ssr...)
			fmt.Fprintf(f, "\n")

			color.Set(color.FgGreen)
			fmt.Fprintf(f, "\nDone.\n")
			color.Unset()

			return nil
		}
	}

	var err error
	if err = rFunc(); err != nil {
		return err
	}

escape:
	for {
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
