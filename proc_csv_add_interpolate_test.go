package psn

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/coreos/etcd/pkg/netutil"
)

func TestCombine(t *testing.T) {
	dn, err := GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	nm, err := netutil.GetDefaultInterfaces()
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	var nt string
	for k := range nm {
		nt = k
		break
	}
	fmt.Println("running with network interface", nt)

	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("test-%010d.csv", time.Now().UnixNano()))
	defer os.RemoveAll(fpath)

	epath := filepath.Join(homeDir(), "etcd-client-num")
	defer os.RemoveAll(epath)

	if err := toFile([]byte("10"), epath); err != nil {
		t.Fatal(err)
	}
	c := NewCSV(fpath, 1, dn, nt, epath)
	for i := 0; i < 3; i++ {
		if i > 0 {
			if err := toFile([]byte(fmt.Sprintf("%d", 100*i)), epath); err != nil {
				t.Fatal(err)
			}
		}

		now := time.Now()
		if err := c.Add(); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("#%d: collected data with %s and %s at %s (c.Add took %v)\n", i, dn, nt, fpath, time.Since(now))

		time.Sleep(time.Second)
	}

	combined := Combine(c.Rows...)
	if combined.UnixSecond != c.Rows[len(c.Rows)-1].UnixSecond {
		t.Fatalf("unix second expected %d, got %d", combined.UnixSecond, c.Rows[len(c.Rows)-1].UnixSecond)
	}

	// test cumulative fields
	if combined.NSEntry.ReceiveBytesNum > c.Rows[len(c.Rows)-1].NSEntry.ReceiveBytesNum {
		t.Fatalf("averaged combined field should have equal or smaller receive combined.NSEntry.ReceiveBytesNum than the last one (last %v, got %v)", c.Rows[len(c.Rows)-1].NSEntry.ReceiveBytesNum, combined.NSEntry.ReceiveBytesNum)
	}
	if combined.NSEntry.TransmitBytesNum > c.Rows[len(c.Rows)-1].NSEntry.TransmitBytesNum {
		t.Fatalf("averaged combined field should have equal or smaller receive combined.NSEntry.TransmitBytesNum than the last one (last %v, got %v)", c.Rows[len(c.Rows)-1].NSEntry.TransmitBytesNum, combined.NSEntry.TransmitBytesNum)
	}

	fmt.Println()
	for i, row := range c.Rows {
		fmt.Printf("row #%d: %+v\n", i, row)
	}
	fmt.Println()
	fmt.Printf("combined: %+v\n", combined)
}

func TestInterpolate(t *testing.T) {

}
