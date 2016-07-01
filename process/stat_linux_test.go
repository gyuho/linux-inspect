package process

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

func TestGetUptime(t *testing.T) {
	u, err := GetUptime()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetUptime: %+v\n", u)
}
