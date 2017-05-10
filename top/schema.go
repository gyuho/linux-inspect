package top

import (
	"reflect"

	"github.com/gyuho/linux-inspect/schema"
)

// RowSchema represents a row in 'top' command output.
// Reference http://man7.org/linux/man-pages/man1/top.1.html.
var RowSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "PID", Godoc: "pid of the process", Kind: reflect.Int64},
		{Name: "USER", Godoc: "user name", Kind: reflect.String},
		{Name: "PR", Godoc: "priority", Kind: reflect.String},
		{Name: "NI", Godoc: "nice value of the task", Kind: reflect.String},
		{Name: "VIRT", Godoc: "total amount  of virtual memory used by the task (in KiB)", Kind: reflect.String},
		{Name: "RES", Godoc: "non-swapped physical memory a task is using (in KiB)", Kind: reflect.String},
		{Name: "SHR", Godoc: "amount of shared memory available to a task, not all of which is typically resident (in KiB)", Kind: reflect.String},
		{Name: "S", Godoc: "process status", Kind: reflect.String},
		{Name: "CPUPercent", Godoc: "%CPU", Kind: reflect.Float64},
		{Name: "MEMPercent", Godoc: "%MEM", Kind: reflect.Float64},
		{Name: "TIME", Godoc: "CPU time (TIME+)", Kind: reflect.String},
		{Name: "COMMAND", Godoc: "command", Kind: reflect.String},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"S":    schema.TypeStatus,
		"VIRT": schema.TypeBytes,
		"RES":  schema.TypeBytes,
		"SHR":  schema.TypeBytes,
	},
}
