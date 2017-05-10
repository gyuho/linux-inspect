package inspect

import (
	"fmt"
	"strings"

	"github.com/gyuho/linux-inspect/top"
)

// EntryOp defines entry option(filter).
type EntryOp struct {
	ProgramMatchFunc func(string) bool
	program          string

	PID      int64
	TopLimit int

	// for ss
	TCP        bool
	TCP6       bool
	LocalPort  int64
	RemotePort int64

	// for ps
	TopExecPath string
	TopStream   *top.Stream

	// for Proc
	DiskDevice       string
	NetworkInterface string
	ExtraPath        string
}

// OpFunc applies each filter.
type OpFunc func(*EntryOp)

// WithProgramMatch matches command name.
func WithProgramMatch(matchFunc func(string) bool) OpFunc {
	return func(op *EntryOp) { op.ProgramMatchFunc = matchFunc }
}

// WithProgram to filter entries by program name.
func WithProgram(name string) OpFunc {
	return func(op *EntryOp) {
		op.ProgramMatchFunc = func(commandName string) bool {
			return strings.HasSuffix(commandName, name)
		}
		op.program = name
	}
}

// WithPID to filter entries by PIDs.
func WithPID(pid int64) OpFunc {
	return func(op *EntryOp) { op.PID = pid }
}

// WithTopLimit to filter entries with limit.
func WithTopLimit(limit int) OpFunc {
	return func(op *EntryOp) { op.TopLimit = limit }
}

// WithLocalPort to filter entries by local port.
func WithLocalPort(port int64) OpFunc {
	return func(op *EntryOp) { op.LocalPort = port }
}

// WithRemotePort to filter entries by remote port.
func WithRemotePort(port int64) OpFunc {
	return func(op *EntryOp) { op.RemotePort = port }
}

// WithTCP to filter entries by TCP.
// Can be used with 'WithTCP6'.
func WithTCP() OpFunc {
	return func(op *EntryOp) { op.TCP = true }
}

// WithTCP6 to filter entries by TCP6.
// Can be used with 'WithTCP'.
func WithTCP6() OpFunc {
	return func(op *EntryOp) { op.TCP6 = true }
}

// WithTopExecPath configures 'top' command path.
func WithTopExecPath(path string) OpFunc {
	return func(op *EntryOp) { op.TopExecPath = path }
}

// WithTopStream gets the PSEntry from the 'top' stream.
func WithTopStream(str *top.Stream) OpFunc {
	return func(op *EntryOp) { op.TopStream = str }
}

// WithDiskDevice to filter entries by disk device.
func WithDiskDevice(name string) OpFunc {
	return func(op *EntryOp) { op.DiskDevice = name }
}

// WithNetworkInterface to filter entries by disk device.
func WithNetworkInterface(name string) OpFunc {
	return func(op *EntryOp) { op.NetworkInterface = name }
}

// WithExtraPath to filter entries by disk device.
func WithExtraPath(path string) OpFunc {
	return func(op *EntryOp) { op.ExtraPath = path }
}

// applyOpts panics when op.Program != "" && op.PID > 0.
func (op *EntryOp) applyOpts(opts []OpFunc) {
	for _, of := range opts {
		of(op)
	}

	if op.DiskDevice != "" || op.NetworkInterface != "" || op.ExtraPath != "" {
		if (op.program != "" || op.ProgramMatchFunc != nil) || op.TopLimit > 0 || op.LocalPort > 0 || op.RemotePort > 0 || op.TCP || op.TCP6 {
			panic(fmt.Errorf("not-valid Proc fileter; disk device %q or network interface %q or extra path %q", op.DiskDevice, op.NetworkInterface, op.ExtraPath))
		}
	}
	if (op.program != "" || op.ProgramMatchFunc != nil) && op.PID > 0 {
		panic(fmt.Errorf("can't filter both by program(%q or %p) and PID(%d)", op.program, op.ProgramMatchFunc, op.PID))
	}
	if !op.TCP && !op.TCP6 {
		// choose both
		op.TCP, op.TCP6 = true, true
	}
	if op.LocalPort > 0 && op.RemotePort > 0 {
		panic(fmt.Errorf("can't query by both local(%d) and remote(%d) ports", op.LocalPort, op.RemotePort))
	}

	if op.TopExecPath == "" {
		op.TopExecPath = top.DefaultExecPath
	}
}
