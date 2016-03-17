package ps

import (
	"fmt"
	"testing"
)

func TestGetStat(t *testing.T) {
	rs, err := getStat("/proc/1/stat")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("getStat: %+v\n", rs)
}

func TestGetUptime(t *testing.T) {
	u, err := GetUptime()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetUptime: %+v\n", u)
}

func TestGetCpuUsage(t *testing.T) {
	s, err := getStat("/proc/1/stat")
	if err != nil {
		t.Error(err)
	}
	c, err := s.GetCpuUsage()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetCpuUsage: %+v\n", c)
}
