package proc

import (
	"fmt"
	"strconv"
	"strings"
)

func convertStatus(s string) string {
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

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
