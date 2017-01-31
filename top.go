package psn

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// GetTop returns all entries in 'top' command.
// If pid<1, it reads all processes in 'top' command.
func GetTop(topPath string, pid int64) ([]TopCommandRow, error) {
	o, err := ReadTop(topPath, pid)
	if err != nil {
		return nil, err
	}
	return ParseTopOutput(o)
}

// GetTopDefault returns all entries in 'top' command.
// If pid<1, it reads all processes in 'top' command.
func GetTopDefault(pid int64) ([]TopCommandRow, error) {
	o, err := ReadTop(DefaultTopPath, pid)
	if err != nil {
		return nil, err
	}
	return ParseTopOutput(o)
}

// DefaultTopPath is the default 'top' command path.
var DefaultTopPath = "/usr/bin/top"

// ReadTopDefault reads Linux 'top' command output.
func ReadTopDefault(pid int64) (string, error) {
	return ReadTop(DefaultTopPath, pid)
}

// ReadTop reads Linux 'top' command output.
func ReadTop(topPath string, pid int64) (string, error) {
	buf := new(bytes.Buffer)
	err := readTop(topPath, pid, buf)
	o := strings.TrimSpace(buf.String())
	return o, err
}

func readTop(topPath string, pid int64, w io.Writer) error {
	if !exist(topPath) {
		return fmt.Errorf("%q does not exist", topPath)
	}
	topFlags := []string{"-b", "-n", "1"}
	if pid > 0 {
		topFlags = append(topFlags, "-p", fmt.Sprint(pid))
	}
	cmd := exec.Command(topPath, topFlags...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func convertProcStatus(s string) string {
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

// parses KiB strings, returns bytes in int64, and humanized bytes.
//
//  KiB = kibibyte = 1024 bytes
//  MiB = mebibyte = 1024 KiB = 1,048,576 bytes
//  GiB = gibibyte = 1024 MiB = 1,073,741,824 bytes
//  TiB = tebibyte = 1024 GiB = 1,099,511,627,776 bytes
//  PiB = pebibyte = 1024 TiB = 1,125,899,906,842,624 bytes
//  EiB = exbibyte = 1024 PiB = 1,152,921,504,606,846,976 bytes
//
func parseKiBInTop(s string) (bts uint64, hs string, err error) {
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

// ParseTopOutput parses 'top' command output and returns the rows.
func ParseTopOutput(s string) ([]TopCommandRow, error) {
	lines := strings.Split(s, "\n")
	rows := make([][]string, 0, len(lines))
	headerFound := false
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		ds := strings.Fields(strings.TrimSpace(line))
		if ds[0] == "PID" { // header line
			if !reflect.DeepEqual(ds, TopRowHeaders) {
				return nil, fmt.Errorf("unexpected 'top' command header order (%v, expected %v)", ds, TopRowHeaders)
			}
			headerFound = true
			continue
		}

		if !headerFound {
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
	rc := make(chan result)
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

	virt, virtTxt, err := parseKiBInTop(row[top_command_output_row_idx_virt])
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.VIRT = row[top_command_output_row_idx_virt]
	trow.VIRTBytesN = virt
	trow.VIRTParsedBytes = virtTxt

	res, resTxt, err := parseKiBInTop(row[top_command_output_row_idx_res])
	if err != nil {
		return TopCommandRow{}, fmt.Errorf("parse error %v (row %v)", err, row)
	}
	trow.RES = row[top_command_output_row_idx_res]
	trow.RESBytesN = res
	trow.RESParsedBytes = resTxt

	shr, shrTxt, err := parseKiBInTop(row[top_command_output_row_idx_shr])
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
