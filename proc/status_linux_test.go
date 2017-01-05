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
	fmt.Printf("GetStatus: %+v\n", rs)
}
