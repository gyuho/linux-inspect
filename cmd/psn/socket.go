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
	"github.com/gyuho/psn/socket"
	"github.com/spf13/cobra"
)

type socketFlags struct {
	LogPath string

	Top    int
	Filter socket.Process

	Kill    bool
	CleanUp bool

	Monitor         bool
	MonitorInterval time.Duration
}

var (
	socketCommand = &cobra.Command{
		Use:   "socket",
		Short: "Investigates sockets",
		RunE:  socketCommandFunc,
	}
	socketKillCommand = &cobra.Command{
		Use:   "socket-kill",
		Short: "Kills sockets",
		RunE:  socketKillCommandFunc,
	}
	socketMonitorCommand = &cobra.Command{
		Use:   "socket-monitor",
		Short: "Monitors sockets",
		RunE:  socketMonitorCommandFunc,
	}
	socketCmdFlag = socketFlags{}
)

func init() {
	socketCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	socketCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	socketCommand.PersistentFlags().IntVarP(&socketCmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")

	socketKillCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	socketKillCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	socketKillCommand.PersistentFlags().BoolVarP(&socketCmdFlag.CleanUp, "clean-up", "c", false, "'true' to automatically kill deleted processes. Name must be empty.")

	socketMonitorCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.Program, "program", "s", "", "Specify the program. Empty lists all programs.")
	socketMonitorCommand.PersistentFlags().StringVarP(&socketCmdFlag.Filter.LocalPort, "local-port", "l", "", "Specify the local port. Empty lists all local ports.")
	socketMonitorCommand.PersistentFlags().IntVarP(&socketCmdFlag.Top, "top", "t", 0, "Only list the top processes (descending order in memory usage). 0 means all.")
	socketMonitorCommand.PersistentFlags().StringVar(&socketCmdFlag.LogPath, "log-path", "", "File path to store logs. Empty to print out to stdout.")
	socketMonitorCommand.PersistentFlags().DurationVar(&socketCmdFlag.MonitorInterval, "monitor-interval", 10*time.Second, "Monitor interval.")
}

func socketCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\npsn socket\n\n")
	color.Unset()

	ssr, err := socket.List(&socketCmdFlag.Filter, socket.TCP, socket.TCP6)
	if err != nil {
		return err
	}
	socket.WriteToTable(os.Stdout, socketCmdFlag.Top, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func socketKillCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgRed)
	fmt.Fprintf(os.Stdout, "\npsn socket-kill\n\n")
	color.Unset()

	if socketCmdFlag.CleanUp && socketCmdFlag.Filter.Program == "" {
		socketCmdFlag.Filter.Program = "deleted)"
	} else if socketCmdFlag.Filter.LocalPort == "" && socketCmdFlag.Filter.Program == "" { // to prevent killing all
		socketCmdFlag.Filter.Program = "SPECIFY PROGRAM NAME"
	}

	ssr, err := socket.List(&socketCmdFlag.Filter, socket.TCP, socket.TCP6)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nerror (%v)\n", err)
		return nil
	}
	socket.WriteToTable(os.Stdout, socketCmdFlag.Top, ssr...)
	socket.Kill(os.Stdout, ssr...)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDone.\n")
	color.Unset()

	return nil
}

func socketMonitorCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgBlue)
	fmt.Fprintf(os.Stdout, "\npsn socket-monitor\n\n")
	color.Unset()

	rFunc := func() error {
		ssr, err := socket.List(&socketCmdFlag.Filter, socket.TCP, socket.TCP6)
		if err != nil {
			return err
		}

		if filepath.Ext(socketCmdFlag.LogPath) == ".csv" {
			f, err := openToAppend(socketCmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := socket.WriteToCSV(f, ssr...); err != nil {
				return err
			}
			return err
		}

		var wr io.Writer
		if socketCmdFlag.LogPath == "" {
			wr = os.Stdout
		} else {
			f, err := openToAppend(socketCmdFlag.LogPath)
			if err != nil {
				return err
			}
			defer f.Close()
			wr = f
		}
		socket.WriteToTable(wr, socketCmdFlag.Top, ssr...)
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
		case <-time.After(socketCmdFlag.MonitorInterval):
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
