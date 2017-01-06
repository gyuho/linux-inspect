package psn

import (
	"fmt"
	"testing"
)

func TestGetNetDev(t *testing.T) {
	ns, err := GetNetDev()
	if err != nil {
		t.Error(err)
	}
	for _, nd := range ns {
		fmt.Printf("%+v\n", nd)
	}
}
