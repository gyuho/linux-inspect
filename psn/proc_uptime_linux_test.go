package psn

import (
	"fmt"
	"testing"
)

func TestGetProcUptime(t *testing.T) {
	u, err := GetProcUptime()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetProcUptime: %+v\n", u)
}
