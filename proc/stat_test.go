package proc

import (
	"fmt"
	"testing"
)

func TestGetStatByPID(t *testing.T) {
	s, err := GetStatByPID(1)
	if err != nil {
		t.Skip(err)
	}
	if s.Rss != s.RssBytesN {
		t.Fatalf("got different Rss; %d != %d", s.Rss, s.RssBytesN)
	}
	fmt.Printf("GetStatByPID: %+v\n", s)
}
