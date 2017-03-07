// Package schema defines 'top' command output schema.
package schema

import (
	"reflect"

	"github.com/gyuho/linux-inspect/pkg/schemautil"
)

// TopCommandRow represents a row in 'top' command output.
// (See http://man7.org/linux/man-pages/man1/top.1.html).
var TopCommandRow = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"PID", "pid of the process", reflect.Int64},
		{"USER", "user name", reflect.String},
		{"PR", "priority", reflect.String},
		{"NI", "nice value of the task", reflect.String},
		{"VIRT", "total amount  of virtual memory used by the task (in KiB)", reflect.String},
		{"RES", "non-swapped physical memory a task is using (in KiB)", reflect.String},
		{"SHR", "amount of shared memory available to a task, not all of which is typically resident (in KiB)", reflect.String},
		{"S", "process status", reflect.String},
		{"CPUPercent", "%CPU", reflect.Float64},
		{"MEMPercent", "%MEM", reflect.Float64},
		{"TIME", "CPU time (TIME+)", reflect.String},
		{"COMMAND", "command", reflect.String},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"S":    schemautil.TypeStatus,
		"VIRT": schemautil.TypeBytes,
		"RES":  schemautil.TypeBytes,
		"SHR":  schemautil.TypeBytes,
	},
}
