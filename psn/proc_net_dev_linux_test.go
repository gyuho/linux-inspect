package psn

import (
	"fmt"
	"testing"
)

func TestGetProcNetDev(t *testing.T) {
	ns, err := GetProcNetDev()
	if err != nil {
		t.Error(err)
	}
	for _, nd := range ns {
		fmt.Printf("%+v\n", nd)
	}
}
