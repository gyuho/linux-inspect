package proc

import (
	"fmt"
	"testing"
)

func TestGetStat(t *testing.T) {
	rs, err := GetStat(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetStat: %+v\n", rs)
}
