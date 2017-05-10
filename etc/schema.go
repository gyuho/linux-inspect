package etc

import (
	"reflect"

	"github.com/gyuho/linux-inspect/schema"
)

// MtabSchema represents '/etc/mtab'.
// Reference https://en.wikipedia.org/wiki/Fstab
// and https://en.wikipedia.org/wiki/Mtab).
var MtabSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "file-system", Godoc: "file system", Kind: reflect.String},
		{Name: "mounted-on", Godoc: "'mounted on'", Kind: reflect.String},
		{Name: "file-system-type", Godoc: "file system type", Kind: reflect.String},
		{Name: "options", Godoc: "file system type", Kind: reflect.String},
		{Name: "dump", Godoc: "number indicating whether and how often the file system should be backed up by the dump program; a zero indicates the file system will never be automatically backed up", Kind: reflect.Int},
		{Name: "pass", Godoc: "number indicating the order in which the fsck program will check the devices for errors at boot time; this is 1 for the root file system and either 2 (meaning check after root) or 0 (do not check) for all other devices", Kind: reflect.Int},
	},
	ColumnsToParse: map[string]schema.RawDataType{},
}
