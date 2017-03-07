// Package schema defines 'df' command output schema.
package schema

import (
	"reflect"

	"github.com/gyuho/linux-inspect/pkg/schemautil"
)

// DfCommandRow represents 'df' command output row
// (See https://en.wikipedia.org/wiki/Df_(Unix)
// and https://www.gnu.org/software/coreutils/manual/html_node/df-invocation.html
// and 'df --all --sync --block-size=1024 --output=source,target,fstype,file,itotal,iavail,iused,ipcent,size,avail,used,pcent'
// and the output unit is kilobytes).
var DfCommandRow = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"file-system", "file system ('source')", reflect.String},
		{"device", "device name", reflect.String},
		{"mounted-on", "'mounted on' ('target')", reflect.String},
		{"file-system-type", "file system type ('fstype')", reflect.String},
		{"file", "file name if specified on the command line ('file')", reflect.String},

		{"inodes", "total number of inodes ('itotal')", reflect.Int64},
		{"ifree", "number of available inodes ('iavail')", reflect.Int64},
		{"iused", "number of used inodes ('iused')", reflect.Int64},
		{"iused-percent", "percentage of iused divided by itotal ('ipcent')", reflect.String},

		{"total-blocks", "total number of 1K-blocks ('size')", reflect.Int64},
		{"available-blocks", "number of available 1K-blocks ('avail')", reflect.Int64},
		{"used-blocks", "number of used 1K-blocks ('used')", reflect.Int64},
		{"used-blocks-percent", "percentage of used-blocks divided by total-blocks ('pcent')", reflect.String},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"total-blocks":     schemautil.TypeBytes,
		"available-blocks": schemautil.TypeBytes,
		"used-blocks":      schemautil.TypeBytes,
	},
}
