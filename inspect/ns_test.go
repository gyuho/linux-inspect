package inspect

import (
	"fmt"
	"testing"
)

func TestGetNS(t *testing.T) {
	ns, err := GetNS()
	if err != nil {
		t.Skip(err)
	}
	hd, rows := ConvertNS(ns...)
	txt := StringNS(hd, rows, -1)
	fmt.Println(txt)
}
