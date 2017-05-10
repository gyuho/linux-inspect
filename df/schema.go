package df

import (
	"reflect"

	"github.com/gyuho/linux-inspect/schema"
)

// RowSchema represents 'df' command output row
// (See https://en.wikipedia.org/wiki/Df_(Unix)
// and https://www.gnu.org/software/coreutils/manual/html_node/df-invocation.html
// and 'df --all --sync --block-size=1024 --output=source,target,fstype,file,itotal,iavail,iused,ipcent,size,avail,used,pcent'
// and the output unit is kilobytes).
var RowSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "file-system", Godoc: "file system ('source')", Kind: reflect.String},
		{Name: "device", Godoc: "device name", Kind: reflect.String},
		{Name: "mounted-on", Godoc: "'mounted on' ('target')", Kind: reflect.String},
		{Name: "file-system-type", Godoc: "file system type ('fstype')", Kind: reflect.String},
		{Name: "file", Godoc: "file name if specified on the command line ('file')", Kind: reflect.String},

		{Name: "inodes", Godoc: "total number of inodes ('itotal')", Kind: reflect.Int64},
		{Name: "ifree", Godoc: "number of available inodes ('iavail')", Kind: reflect.Int64},
		{Name: "iused", Godoc: "number of used inodes ('iused')", Kind: reflect.Int64},
		{Name: "iused-percent", Godoc: "percentage of iused divided by itotal ('ipcent')", Kind: reflect.String},

		{Name: "total-blocks", Godoc: "total number of 1K-blocks ('size')", Kind: reflect.Int64},
		{Name: "available-blocks", Godoc: "number of available 1K-blocks ('avail')", Kind: reflect.Int64},
		{Name: "used-blocks", Godoc: "number of used 1K-blocks ('used')", Kind: reflect.Int64},
		{Name: "used-blocks-percent", Godoc: "percentage of used-blocks divided by total-blocks ('pcent')", Kind: reflect.String},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"total-blocks":     schema.TypeBytes,
		"available-blocks": schema.TypeBytes,
		"used-blocks":      schema.TypeBytes,
	},
}
