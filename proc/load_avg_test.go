package proc

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetLoadAvg(t *testing.T) {
	txt, err := readLoadAvg()
	if err != nil {
		t.Fatal(err)
	}
	lv, err := getLoadAvg(txt)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v with %q\n", lv, txt)

	if !strings.Contains(txt, fmt.Sprint(lv.LoadAvg15Minute)) {
		t.Fatalf("'/proc/loadavg' expected LoadAvg15Minute %f, got %q", lv.LoadAvg15Minute, txt)
	}
	if !strings.Contains(txt, fmt.Sprint(lv.Pid)) {
		t.Fatalf("'/proc/loadavg' expected pid %d, got %q", lv.Pid, txt)
	}
}
