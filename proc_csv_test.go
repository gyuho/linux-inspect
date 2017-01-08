package psn

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcCSV(t *testing.T) {
	dn, err := GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	nt, err := GetDefaultInterface()
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}

	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("test-%010d.csv", rand.Intn(999999)))
	defer os.RemoveAll(fpath)

	c := NewCSV(fpath, 1, dn, nt)
	for i := 0; i < 5; i++ {
		fmt.Printf("#%d: collecting data with %s and %s at %s\n", i, dn, nt, fpath)
		if err := c.Add(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	fmt.Println("CSV contents:", string(b))
}
