package psn

import (
	"fmt"
	"strings"
)

// EntryFilter defines entry filter.
type EntryFilter struct {
	ProgramMatchFunc func(string) bool
	program          string
	PID              int64
	TopLimit         int
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

// applyOpts panics when ft.Program != "" && ft.PID > 0.
func (ft *EntryFilter) applyOpts(opts []FilterFunc) {
	for _, opt := range opts {
		opt(ft)
	}

	if (ft.program != "" || ft.ProgramMatchFunc != nil) && ft.PID > 0 {
		panic(fmt.Errorf("can't filter both by program(%q or %p) and PID(%d)", ft.program, ft.ProgramMatchFunc, ft.PID))
	}
}
