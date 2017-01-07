package psn

import (
	"fmt"
	"testing"
)

func TestGetDiskstats(t *testing.T) {
	ds, err := GetDiskstats()
	if err != nil {
		t.Error(err)
	}
	for _, ds := range ds {
		if ds.ReadsCompleted == 0 {
			continue
		}
		fmt.Printf("%s %d\n", ds.DeviceName, ds.ReadsCompleted)
	}
}
