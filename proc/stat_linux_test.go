package proc

import (
	"fmt"
	"testing"
)

func TestGetStat(t *testing.T) {
	u, err := GetUptime()
	if err != nil {
		t.Fatal(err)
	}
	rs, err := GetStat(1, u)
	if err != nil {
		t.Error(err)
	}
	if rs.Rss != rs.RssBytesN {
		t.Fatalf("got different Rss; %d != %d", rs.Rss, rs.RssBytesN)
	}
	fmt.Printf("GetStat: %+v\n", rs)
}
