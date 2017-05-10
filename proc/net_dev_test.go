package proc

import (
	"fmt"
	"testing"
)

func TestGetNetDev(t *testing.T) {
	nds, err := GetNetDev()
	if err != nil {
		t.Skip(err)
	}
	for _, nd := range nds {
		fmt.Printf("%+v\n", nd)
	}
}
