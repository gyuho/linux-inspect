package proc

import (
	"fmt"
	"io/ioutil"

	"github.com/gyuho/linux-inspect/pkg/fileutil"

	humanize "github.com/dustin/go-humanize"
	yaml "gopkg.in/yaml.v2"
)

// GetIOByPID reads '/proc/$PID/io' data.
func GetIOByPID(pid int64) (s IO, err error) {
	fpath := fmt.Sprintf("/proc/%d/io", pid)
	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		return IO{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return IO{}, err
	}

	rs := IO{}
	if err := yaml.Unmarshal(b, &rs); err != nil {
		return rs, err
	}

	rs.RcharBytesN = rs.Rchar
	rs.RcharParsedBytes = humanize.Bytes(uint64(rs.RcharBytesN))

	rs.WcharBytesN = rs.Wchar
	rs.WcharParsedBytes = humanize.Bytes(uint64(rs.WcharBytesN))

	rs.ReadBytesBytesN = rs.ReadBytes
	rs.ReadBytesParsedBytes = humanize.Bytes(uint64(rs.ReadBytesBytesN))

	rs.WriteBytesBytesN = rs.WriteBytes
	rs.WriteBytesParsedBytes = humanize.Bytes(uint64(rs.WriteBytesBytesN))

	rs.CancelledWriteBytesBytesN = rs.CancelledWriteBytes
	rs.CancelledWriteBytesParsedBytes = humanize.Bytes(uint64(rs.CancelledWriteBytesBytesN))

	return rs, nil
}
