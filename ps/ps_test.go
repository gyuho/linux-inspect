package ps

import (
	"fmt"
	"testing"
)

func TestListStatus(t *testing.T) {
	rs, err := ListStatus(nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", rs)
	fmt.Println(len(rs))
}

func TestGetStatusByPID(t *testing.T) {
	rs, err := GetStatusByPID(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", rs)
}
