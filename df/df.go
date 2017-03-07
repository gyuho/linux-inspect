package df

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/gyuho/linux-inspect/pkg/fileutil"
)

// GetDf returns entries in 'df' command.
// Pass '' target to list all information.
func GetDf(dfPath string, target string) ([]DfCommandRow, error) {
	o, err := ReadDf(dfPath, target)
	if err != nil {
		return nil, err
	}
	return ParseDfOutput(o)
}

// GetDfDefault returns entries in 'df' command.
// Pass '' target to list all information.
func GetDfDefault(target string) ([]DfCommandRow, error) {
	o, err := ReadDf(DefaultDfPath, target)
	if err != nil {
		return nil, err
	}
	return ParseDfOutput(o)
}

// DefaultDfPath is the default 'df' command path.
var DefaultDfPath = "/bin/df"

// DfFlags is 'df --all --sync --block-size=1024 --output=source,target,fstype,file,itotal,iavail,iused,ipcent,size,avail,used,pcent'.
var DfFlags = []string{"--all", "--sync", "--block-size=1024", "--output=source,target,fstype,file,itotal,iavail,iused,ipcent,size,avail,used,pcent"}

// ReadDfDefault reads Linux 'df' command output.
// Pass '' target to list all information.
func ReadDfDefault(target string) (string, error) {
	return ReadDf(DefaultDfPath, target)
}

// ReadDf reads Linux 'df' command output.
// Pass '' target to list all information.
func ReadDf(dfPath string, target string) (string, error) {
	buf := new(bytes.Buffer)
	err := readDf(dfPath, target, buf)
	o := strings.TrimSpace(buf.String())
	return o, err
}

func readDf(dfPath string, target string, w io.Writer) error {
	if !fileutil.Exist(dfPath) {
		return fmt.Errorf("%q does not exist", dfPath)
	}
	if target != "" {
		DfFlags = append(DfFlags, strings.TrimSpace(target))
	}
	cmd := exec.Command(dfPath, DfFlags...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

// DfRowHeaders is the headers in 'df' output.
var DfRowHeaders = []string{
	"Filesystem",

	// Mounted on
	"Mounted",
	"on",

	"Type",
	"File",
	"Inodes",
	"IFree",
	"IUsed",
	"IUse%",
	"1K-blocks",
	"Avail",
	"Used",
	"Use%",
}

type dfCommandOutpudrowIdx int

const (
	df_command_output_row_idx_file_system dfCommandOutpudrowIdx = iota
	df_command_output_row_idx_mounted_on
	df_command_output_row_idx_file_system_type
	df_command_output_row_idx_file
	df_command_output_row_idx_inodes
	df_command_output_row_idx_ifree
	df_command_output_row_idx_iused
	df_command_output_row_idx_iused_percent
	df_command_output_row_idx_total_blocks
	df_command_output_row_idx_available_blocks
	df_command_output_row_idx_used_blocks
	df_command_output_row_idx_used_blocks_percentage
)

// ParseDfOutput parses 'df' command output and returns the rows.
func ParseDfOutput(s string) ([]DfCommandRow, error) {
	lines := strings.Split(s, "\n")
	rows := make([][]string, 0, len(lines))
	headerFound := false
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		ds := strings.Fields(strings.TrimSpace(line))
		if ds[0] == "Filesystem" { // header line
			if !reflect.DeepEqual(ds, DfRowHeaders) {
				return nil, fmt.Errorf("unexpected 'df' command header order (%v, expected %v, output: %q)", ds, DfRowHeaders, s)
			}
			headerFound = true
			continue
		}

		if !headerFound {
			continue
		}

		row := strings.Fields(strings.TrimSpace(line))
		if len(row) != len(DfRowHeaders)-1 {
			return nil, fmt.Errorf("unexpected row column number %v (expected %v)", row, DfRowHeaders)
		}
		rows = append(rows, row)
	}

	type result struct {
		row DfCommandRow
		err error
	}
	rc := make(chan result, len(rows))
	for _, row := range rows {
		go func(row []string) {
			tr, err := parseDfRow(row)
			rc <- result{row: tr, err: err}
		}(row)
	}

	tcRows := make([]DfCommandRow, 0, len(rows))
	for len(tcRows) != len(rows) {
		select {
		case rs := <-rc:
			if rs.err != nil {
				return nil, rs.err
			}
			tcRows = append(tcRows, rs.row)
		}
	}
	rm := make(map[string]DfCommandRow)
	for _, row := range tcRows {
		rm[row.MountedOn] = row
	}
	rrs := make([]DfCommandRow, 0, len(rm))
	for _, row := range rm {
		rrs = append(rrs, row)
	}
	return rrs, nil
}

func parseDfRow(row []string) (DfCommandRow, error) {
	drow := DfCommandRow{
		FileSystem:        strings.TrimSpace(row[df_command_output_row_idx_file_system]),
		MountedOn:         strings.TrimSpace(row[df_command_output_row_idx_mounted_on]),
		FileSystemType:    strings.TrimSpace(row[df_command_output_row_idx_file_system_type]),
		File:              strings.TrimSpace(row[df_command_output_row_idx_file]),
		IusedPercent:      strings.TrimSpace(strings.Replace(row[df_command_output_row_idx_iused_percent], "%", " %", -1)),
		UsedBlocksPercent: strings.TrimSpace(strings.Replace(row[df_command_output_row_idx_used_blocks_percentage], "%", " %", -1)),
	}
	drow.Device = filepath.Base(drow.FileSystem)

	ptxt := strings.TrimSpace(row[df_command_output_row_idx_inodes])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err := strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.Inodes = iv

	ptxt = strings.TrimSpace(row[df_command_output_row_idx_ifree])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err = strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.Ifree = iv

	ptxt = strings.TrimSpace(row[df_command_output_row_idx_iused])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err = strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.Iused = iv

	ptxt = strings.TrimSpace(row[df_command_output_row_idx_total_blocks])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err = strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.TotalBlocks = iv
	drow.TotalBlocksBytesN = iv * 1024
	drow.TotalBlocksParsedBytes = humanize.Bytes(uint64(drow.TotalBlocksBytesN))

	ptxt = strings.TrimSpace(row[df_command_output_row_idx_available_blocks])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err = strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.AvailableBlocks = iv
	drow.AvailableBlocksBytesN = iv * 1024
	drow.AvailableBlocksParsedBytes = humanize.Bytes(uint64(drow.AvailableBlocksBytesN))

	ptxt = strings.TrimSpace(row[df_command_output_row_idx_used_blocks])
	if ptxt == "-" {
		ptxt = "0"
	}
	iv, err = strconv.ParseInt(ptxt, 10, 64)
	if err != nil {
		return DfCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	drow.UsedBlocks = iv
	drow.UsedBlocksBytesN = iv * 1024
	drow.UsedBlocksParsedBytes = humanize.Bytes(uint64(drow.UsedBlocksBytesN))

	return drow, nil
}

// GetDevice returns the device name where dir is mounted.
func GetDevice(target string) (string, error) {
	drows, err := GetDfDefault(target)
	if err != nil {
		return "", err
	}
	if len(drows) != 1 {
		return "", fmt.Errorf("expected 1 df row at %q (got %+v)", target, drows)
	}
	return drows[0].Device, nil
}
