// Package schema defines '/etc' schema.
package schema

import (
	"reflect"

	"github.com/gyuho/linux-inspect/pkg/schemautil"
)

// Mtab represents '/etc/mtab'
// (See https://en.wikipedia.org/wiki/Fstab
// and https://en.wikipedia.org/wiki/Mtab).
var Mtab = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"file-system", "file system", reflect.String},
		{"mounted-on", "'mounted on'", reflect.String},
		{"file-system-type", "file system type", reflect.String},
		{"options", "file system type", reflect.String},
		{"dump", "number indicating whether and how often the file system should be backed up by the dump program; a zero indicates the file system will never be automatically backed up", reflect.Int},
		{"pass", "number indicating the order in which the fsck program will check the devices for errors at boot time; this is 1 for the root file system and either 2 (meaning check after root) or 0 (do not check) for all other devices", reflect.Int},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{},
}
