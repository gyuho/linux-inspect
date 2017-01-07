package psn

import (
	"fmt"
	"testing"
)

func Test_getRow(t *testing.T) {
	dn, err := GetDevice("/")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	nt, err := GetDefaultInterface()
	if err != nil {
		t.Fatal(err)
	}
	row, err := getRow(1, dn, nt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(psCSVColumns)
	fmt.Println(row)
}
