package psn

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// DefaultTopPath is the default 'top' command path.
var DefaultTopPath = "/usr/bin/top"

// TopConfig configures 'top' command runs.
type TopConfig struct {
	Exec string

	// MAKE THIS TRUE BY DEFAULT
	// OTHERWISE PARSER HAS TO DEAL WITH HIGHLIGHTED TEXTS
	//
	// BatchMode is true to start 'top' in batch mode, which could be useful
	// for sending output from 'top' to other programs or to a file.
	// In this mode, 'top' will not accept input and runs until the interations
	// limit ('-n' flag) or until killed.
	// It's '-b' flag.
	// BatchMode bool

	// Limit limits the iteration of 'top' commands to run before exit.
	// If 1, 'top' prints out the current processes and exits.
	// It's '-n' flag.
	Limit int

	// IntervalSecond is the delay time between updates.
	// Default is 1 second.
	// It's '-d' flag.
	IntervalSecond float64

	// PID specifies the PID to monitor.
	// It's '-p' flag.
	PID int64

	// Writer stores 'top' command outputs.
	Writer io.Writer

	cmd *exec.Cmd
}

// Flags returns the 'top' command flags.
func (cfg *TopConfig) Flags() (fs []string) {
	// batch mode by default
	fs = append(fs, "-b")

	if cfg.Limit > 0 { // if 1, command just exists after one output
		fs = append(fs, "-n", fmt.Sprintf("%d", cfg.Limit))
	}

	if cfg.IntervalSecond > 0 {
		fs = append(fs, "-d", fmt.Sprintf("%.2f", cfg.IntervalSecond))
	}

	if cfg.PID > 0 {
		fs = append(fs, "-p", fmt.Sprintf("%d", cfg.PID))
	}

	return
}

// process updates with '*exec.Cmd' for the given 'TopConfig'.
func (cfg *TopConfig) createCmd() error {
	if cfg == nil {
		return fmt.Errorf("TopConfig is nil")
	}
	if !exist(cfg.Exec) {
		return fmt.Errorf("%q does not exist", cfg.Exec)
	}
	flags := cfg.Flags()

	c := exec.Command(cfg.Exec, flags...)
	c.Stdout = cfg.Writer
	c.Stderr = cfg.Writer

	cfg.cmd = c
	return nil
}

// GetTop returns all entries in 'top' command.
// If pid<1, it reads all processes in 'top' command.
// This is one-time command.
func GetTop(topPath string, pid int64) ([]TopCommandRow, error) {
	buf := new(bytes.Buffer)
	cfg := &TopConfig{
		Exec:           topPath,
		Limit:          1,
		IntervalSecond: 1,
		PID:            pid,
		Writer:         buf,
		cmd:            nil,
	}
	if err := cfg.createCmd(); err != nil {
		return nil, err
	}

	// run starts the 'top' command and waits for it to complete.
	if err := cfg.cmd.Run(); err != nil {
		return nil, err
	}
	return ParseTopOutput(buf.String())
}
