package psn

import (
	"fmt"
	"strings"
)

// EntryFilter defines entry filter.
type EntryFilter struct {
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
	TopCommandPath string
	TopStream      *TopStream

	// for Proc
	DiskDevice       string
	NetworkInterface string
	ExtraPath        string
}

// FilterFunc applies each filter.
type FilterFunc func(*EntryFilter)

// WithProgramMatch matches command name.
func WithProgramMatch(matchFunc func(string) bool) FilterFunc {
	return func(ft *EntryFilter) { ft.ProgramMatchFunc = matchFunc }
}

// WithProgram to filter entries by program name.
func WithProgram(name string) FilterFunc {
	return func(ft *EntryFilter) {
		ft.ProgramMatchFunc = func(commandName string) bool {
			return strings.HasSuffix(commandName, name)
		}
		ft.program = name
	}
}

// WithPID to filter entries by PIDs.
func WithPID(pid int64) FilterFunc {
	return func(ft *EntryFilter) { ft.PID = pid }
}

// WithTopLimit to filter entries with limit.
func WithTopLimit(limit int) FilterFunc {
	return func(ft *EntryFilter) { ft.TopLimit = limit }
}

// WithLocalPort to filter entries by local port.
func WithLocalPort(port int64) FilterFunc {
	return func(ft *EntryFilter) { ft.LocalPort = port }
}

// WithRemotePort to filter entries by remote port.
func WithRemotePort(port int64) FilterFunc {
	return func(ft *EntryFilter) { ft.RemotePort = port }
}

// WithTCP to filter entries by TCP.
// Can be used with 'WithTCP6'.
func WithTCP() FilterFunc {
	return func(ft *EntryFilter) { ft.TCP = true }
}

// WithTCP6 to filter entries by TCP6.
// Can be used with 'WithTCP'.
func WithTCP6() FilterFunc {
	return func(ft *EntryFilter) { ft.TCP6 = true }
}

// WithTopCommandPath configures 'top' command path.
func WithTopCommandPath(path string) FilterFunc {
	return func(ft *EntryFilter) { ft.TopCommandPath = path }
}

// WithTopStream gets the PSEntry from the 'top' stream.
func WithTopStream(str *TopStream) FilterFunc {
	return func(ft *EntryFilter) { ft.TopStream = str }
}

// WithDiskDevice to filter entries by disk device.
func WithDiskDevice(name string) FilterFunc {
	return func(ft *EntryFilter) { ft.DiskDevice = name }
}

// WithNetworkInterface to filter entries by disk device.
func WithNetworkInterface(name string) FilterFunc {
	return func(ft *EntryFilter) { ft.NetworkInterface = name }
}

// WithExtraPath to filter entries by disk device.
func WithExtraPath(path string) FilterFunc {
	return func(ft *EntryFilter) { ft.ExtraPath = path }
}

// applyOpts panics when ft.Program != "" && ft.PID > 0.
func (ft *EntryFilter) applyOpts(opts []FilterFunc) {
	for _, opt := range opts {
		opt(ft)
	}

	if ft.DiskDevice != "" || ft.NetworkInterface != "" || ft.ExtraPath != "" {
		if (ft.program != "" || ft.ProgramMatchFunc != nil) || ft.TopLimit > 0 || ft.LocalPort > 0 || ft.RemotePort > 0 || ft.TCP || ft.TCP6 {
			panic(fmt.Errorf("not-valid Proc fileter; disk device %q or network interface %q or extra path %q", ft.DiskDevice, ft.NetworkInterface, ft.ExtraPath))
		}
	}
	if (ft.program != "" || ft.ProgramMatchFunc != nil) && ft.PID > 0 {
		panic(fmt.Errorf("can't filter both by program(%q or %p) and PID(%d)", ft.program, ft.ProgramMatchFunc, ft.PID))
	}
	if !ft.TCP && !ft.TCP6 {
		// choose both
		ft.TCP, ft.TCP6 = true, true
	}
	if ft.LocalPort > 0 && ft.RemotePort > 0 {
		panic(fmt.Errorf("can't query by both local(%d) and remote(%d) ports", ft.LocalPort, ft.RemotePort))
	}

	if ft.TopCommandPath == "" {
		ft.TopCommandPath = DefaultTopPath
	}
}
