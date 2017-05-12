package top

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// parses memory bytes in top command,
// returns bytes in int64, and humanized bytes.
//
//  KiB = kibibyte = 1024 bytes
//  MiB = mebibyte = 1024 KiB = 1,048,576 bytes
//  GiB = gibibyte = 1024 MiB = 1,073,741,824 bytes
//  TiB = tebibyte = 1024 GiB = 1,099,511,627,776 bytes
//  PiB = pebibyte = 1024 TiB = 1,125,899,906,842,624 bytes
//  EiB = exbibyte = 1024 PiB = 1,152,921,504,606,846,976 bytes
//
func parseMemoryTxt(s string) (bts uint64, hs string, err error) {
	s = strings.TrimSpace(s)

	switch {
	case strings.HasSuffix(s, "m"): // suffix 'm' means megabytes
		ns := s[:len(s)-1]
		var mib float64
		mib, err = strconv.ParseFloat(ns, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(mib) * 1024 * 1024

	case strings.HasSuffix(s, "g"): // gigabytes
		ns := s[:len(s)-1]
		var gib float64
		gib, err = strconv.ParseFloat(ns, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(gib) * 1024 * 1024 * 1024

	case strings.HasSuffix(s, "t"): // terabytes
		ns := s[:len(s)-1]
		var tib float64
		tib, err = strconv.ParseFloat(ns, 64)
		if err != nil {
			return 0, "", err
		}
		bts = uint64(tib) * 1024 * 1024 * 1024 * 1024

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

// Headers is the headers in 'top' output.
var Headers = []string{
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

type commandOutputRowIdx int

const (
	command_output_row_idx_pid commandOutputRowIdx = iota
	command_output_row_idx_user
	command_output_row_idx_pr
	command_output_row_idx_ni
	command_output_row_idx_virt
	command_output_row_idx_res
	command_output_row_idx_shr
	command_output_row_idx_s
	command_output_row_idx_cpu
	command_output_row_idx_mem
	command_output_row_idx_time
	command_output_row_idx_command
)

var bytesToSkip = [][]byte{
	{116, 111, 112, 32, 45},                     // 'top -'
	{84, 97, 115, 107, 115, 58, 32},             // 'Tasks: '
	{37, 67, 112, 117, 40, 115, 41, 58, 32},     // '%Cpu(s): '
	{67, 112, 117, 40, 115, 41, 58, 32},         // 'Cpu(s): '
	{75, 105, 66, 32, 77, 101, 109, 32, 58, 32}, // 'KiB Mem : '
	{75, 105, 66, 32, 83, 119, 97, 112, 58, 32}, // 'KiB Swap: '
	{77, 101, 109, 58, 32},                      // 'Mem: '
	{83, 119, 97, 112, 58, 32},                  // 'Swap: '
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

// Parse parses 'top' command output and returns the rows.
func Parse(s string) ([]Row, error) {
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
		if len(row) != len(Headers) {
			return nil, fmt.Errorf("unexpected row column number %v (expected %v)", row, Headers)
		}
		rows = append(rows, row)
	}

	type result struct {
		row Row
		err error
	}
	rc := make(chan result, len(rows))
	for _, row := range rows {
		go func(row []string) {
			tr, err := parseRow(row)
			rc <- result{row: tr, err: err}
		}(row)
	}

	tcRows := make([]Row, 0, len(rows))
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

func parseRow(row []string) (Row, error) {
	trow := Row{
		USER: strings.TrimSpace(row[command_output_row_idx_user]),
	}

	pv, err := strconv.ParseInt(row[command_output_row_idx_pid], 10, 64)
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.PID = pv

	trow.PR = strings.TrimSpace(row[command_output_row_idx_pr])
	trow.NI = strings.TrimSpace(row[command_output_row_idx_ni])

	virt, virtTxt, err := parseMemoryTxt(row[command_output_row_idx_virt])
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.VIRT = row[command_output_row_idx_virt]
	trow.VIRTBytesN = virt
	trow.VIRTParsedBytes = virtTxt

	res, resTxt, err := parseMemoryTxt(row[command_output_row_idx_res])
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.RES = row[command_output_row_idx_res]
	trow.RESBytesN = res
	trow.RESParsedBytes = resTxt

	shr, shrTxt, err := parseMemoryTxt(row[command_output_row_idx_shr])
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.SHR = row[command_output_row_idx_shr]
	trow.SHRBytesN = shr
	trow.SHRParsedBytes = shrTxt

	trow.S = row[command_output_row_idx_s]
	trow.SParsedStatus = parseStatus(row[command_output_row_idx_s])

	cnum, err := strconv.ParseFloat(row[command_output_row_idx_cpu], 64)
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.CPUPercent = cnum

	mnum, err := strconv.ParseFloat(row[command_output_row_idx_mem], 64)
	if err != nil {
		return Row{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.MEMPercent = mnum

	trow.TIME = row[command_output_row_idx_time]

	return trow, nil
}

func parseStatus(s string) string {
	ns := strings.TrimSpace(s)
	if len(s) > 1 {
		ns = ns[:1]
	}
	switch ns {
	case "D":
		return "D (uninterruptible sleep)"
	case "R":
		return "R (running)"
	case "S":
		return "S (sleeping)"
	case "T":
		return "T (stopped by job control signal)"
	case "t":
		return "t (stopped by debugger during trace)"
	case "Z":
		return "Z (zombie)"
	default:
		return fmt.Sprintf("unknown process %q", s)
	}
}
