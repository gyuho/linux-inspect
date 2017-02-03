package psn

import (
	"fmt"
	"testing"
)

func TestGetEtcMtab(t *testing.T) {
	mss, err := GetEtcMtab()
	if err != nil {
		t.Error(err)
	}
	for _, ms := range mss {
		fmt.Printf("%+v\n", ms)
	}
}
