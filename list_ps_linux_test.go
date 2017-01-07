package psn

import (
	"fmt"
	"testing"
)

func TestGetPS(t *testing.T) {
	ns, err := GetPS(WithTopLimit(3))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertPS(ns...)
	txt := StringPS(hd, rows, -1)
	fmt.Println(txt)
}

func TestGetPSWithFilter(t *testing.T) {
	ns, err := GetPS(WithPID(1))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertPS(ns...)
	txt := StringPS(hd, rows, -1)
	fmt.Println(txt)
}
