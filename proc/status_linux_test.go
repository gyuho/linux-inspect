package proc

import (
	"fmt"
	"testing"
)

func TestGetStatus(t *testing.T) {
	rs, err := GetStatus(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs.VmRSS)
	fmt.Println(rs.VmRSSHumanizedBytes)
	fmt.Printf("GetStatus: %+v\n", rs)
}
