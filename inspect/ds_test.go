package inspect

import (
	"fmt"
	"testing"
)

func TestGetDS(t *testing.T) {
	ds, err := GetDS()
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertDS(ds...)
	txt := StringDS(hd, rows, -1)
	fmt.Println(txt)
}
