package psn

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// parses KiB strings, returns bytes in int64, and humanized bytes.
//
//  KiB = kibibyte = 1024 bytes
//  MiB = mebibyte = 1024 KiB = 1,048,576 bytes
//  GiB = gibibyte = 1024 MiB = 1,073,741,824 bytes
//  TiB = tebibyte = 1024 GiB = 1,099,511,627,776 bytes
//  PiB = pebibyte = 1024 TiB = 1,125,899,906,842,624 bytes
//  EiB = exbibyte = 1024 PiB = 1,152,921,504,606,846,976 bytes
//
func parseTopCommandKiB(s string) (bts uint64, hs string, err error) {
	s = strings.TrimSpace(s)
	switch {
	// suffix 'm' means megabytes
	case strings.HasSuffix(s, "m"):
		ns := s[:len(s)-1]
		var mib float64
		mib, err = strconv.ParseFloat(ns, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(mib) * 1024 * 1024

	// suffix 'g' means gigabytes
	case strings.HasSuffix(s, "g"):
		ns := s[:len(s)-1]
		var gib float64
		gib, err = strconv.ParseFloat(ns, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(gib) * 1024 * 1024 * 1024

	default:
		var kib float64
		kib, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(kib) * 1024
	}

	hs = humanize.Bytes(bts)
	return
}

// TopRowHeaders is the headers in 'top' output.
var TopRowHeaders = []string{
	"PID",
	"USER",
	"PR",
	"NI",
	"VIRT",
	"RES",
	"SHR",
	"S",
	"%CPU",
	"%MEM",
	"TIME+",
	"COMMAND",
}

type topCommandOutputRowIdx int

const (
	top_command_output_row_idx_pid topCommandOutputRowIdx = iota
	top_command_output_row_idx_user
	top_command_output_row_idx_pr
	top_command_output_row_idx_ni
	top_command_output_row_idx_virt
	top_command_output_row_idx_res
	top_command_output_row_idx_shr
	top_command_output_row_idx_s
	top_command_output_row_idx_cpu
	top_command_output_row_idx_mem
	top_command_output_row_idx_time
	top_command_output_row_idx_command
)

var bytesToSkip = [][]byte{
	{116, 111, 112, 32, 45},                     // 'top -'
	{84, 97, 115, 107, 115, 58, 32},             // 'Tasks: '
	{37, 67, 112, 117, 40, 115, 41, 58, 32},     // '%Cpu(s): '
	{67, 112, 117, 40, 115, 41, 58, 32},         // 'Cpu(s): '
	{75, 105, 66, 32, 77, 101, 109, 32, 58, 32}, // 'KiB Mem : '
	{75, 105, 66, 32, 83, 119, 97, 112, 58, 32}, // 'KiB Swap: '
	{77, 101, 109, 58, 32},                      // 'Mem: '
	{80, 73, 68, 32},                            // 'PID '
}

func topRowToSkip(data []byte) bool {
	for _, prefix := range bytesToSkip {
		if bytes.HasPrefix(data, prefix) {
			return true
		}
	}
	return false
}

// ParseTopOutput parses 'top' command output and returns the rows.
func ParseTopOutput(s string) ([]TopCommandRow, error) {
	lines := strings.Split(s, "\n")
	rows := make([][]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if topRowToSkip([]byte(line)) {
			continue
		}

		row := strings.Fields(strings.TrimSpace(line))
		if len(row) != len(TopRowHeaders) {
			return nil, fmt.Errorf("unexpected row column number %v (expected %v)", row, TopRowHeaders)
		}
		rows = append(rows, row)
	}

	type result struct {
		row TopCommandRow
		err error
	}
	rc := make(chan result, len(rows))
	for _, row := range rows {
		go func(row []string) {
			tr, err := parseTopRow(row)
			rc <- result{row: tr, err: err}
		}(row)
	}

	tcRows := make([]TopCommandRow, 0, len(rows))
	for len(tcRows) != len(rows) {
		select {
		case rs := <-rc:
			if rs.err != nil {
				return nil, rs.err
			}
			tcRows = append(tcRows, rs.row)
		}
	}
	return tcRows, nil
}

func parseTopRow(row []string) (TopCommandRow, error) {
	trow := TopCommandRow{
		USER: strings.TrimSpace(row[top_command_output_row_idx_user]),
	}

	pv, err := strconv.ParseInt(row[top_command_output_row_idx_pid], 10, 64)
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.PID = pv

	trow.PR = strings.TrimSpace(row[top_command_output_row_idx_pr])
	trow.NI = strings.TrimSpace(row[top_command_output_row_idx_ni])

	virt, virtTxt, err := parseTopCommandKiB(row[top_command_output_row_idx_virt])
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.VIRT = row[top_command_output_row_idx_virt]
	trow.VIRTBytesN = virt
	trow.VIRTParsedBytes = virtTxt

	res, resTxt, err := parseTopCommandKiB(row[top_command_output_row_idx_res])
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.RES = row[top_command_output_row_idx_res]
	trow.RESBytesN = res
	trow.RESParsedBytes = resTxt

	shr, shrTxt, err := parseTopCommandKiB(row[top_command_output_row_idx_shr])
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.SHR = row[top_command_output_row_idx_shr]
	trow.SHRBytesN = shr
	trow.SHRParsedBytes = shrTxt

	trow.S = row[top_command_output_row_idx_s]
	trow.SParsedStatus = convertProcStatus(row[top_command_output_row_idx_s])

	cnum, err := strconv.ParseFloat(row[top_command_output_row_idx_cpu], 64)
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.CPUPercent = cnum

	mnum, err := strconv.ParseFloat(row[top_command_output_row_idx_mem], 64)
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.MEMPercent = mnum

	trow.TIME = row[top_command_output_row_idx_time]

	return trow, nil
}
