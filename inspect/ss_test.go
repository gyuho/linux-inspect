package inspect

import (
	"fmt"
	"testing"
)

func TestGetSS(t *testing.T) {
	ss, err := GetSS(WithTCP(), WithTopLimit(2))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertSS(ss...)
	txt := StringSS(hd, rows, -1)
	fmt.Println(txt)
}

func TestGetSSWithFilter(t *testing.T) {
	ss, err := GetSS(WithPID(1))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertSS(ss...)
	txt := StringSS(hd, rows, -1)
	fmt.Println(txt)
}
