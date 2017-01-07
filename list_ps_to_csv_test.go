package psn

import (
	"fmt"
	"testing"
)

func Test_getRow(t *testing.T) {
	nt, err := GetDefaultInterface()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nt)
}
