package psn

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendSSCSV(t *testing.T) {
	dn, err := GetDevice("/boot")
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

	if err = os.MkdirAll("testdata", 0777); err != nil {
		fmt.Println(err)
		t.Skip()
	}

	fpath := filepath.Join("testdata", fmt.Sprintf("test-%010d.csv", rand.Intn(999999)))
	defer os.RemoveAll(fpath)

	if err := WriteSSCSVHeader(fpath); err != nil {
		fmt.Println(err)
		t.Skip()
	}
	if err := AppendSSCSV(fpath, 1, dn, nt); err != nil {
		fmt.Println(err)
		t.Skip()
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	fmt.Println("CSV fpath:", fpath)
	fmt.Println("CSV contents:", string(b))
}
