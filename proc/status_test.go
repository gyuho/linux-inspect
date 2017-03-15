package proc

import (
	"fmt"
	"testing"
)

func TestGetProcStatusByPID(t *testing.T) {
	rs, err := GetProcStatusByPID(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs.VmRSS)
	fmt.Println(rs.VmRSSBytesN)
	fmt.Println(rs.VmRSSParsedBytes)
	fmt.Printf("GetProcStatusByPID: %+v\n", rs)
}
