package proc

import (
	"fmt"
	"testing"
)

func TestGetUptime(t *testing.T) {
	u, err := GetUptime()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetUptime: %+v\n", u)
}
