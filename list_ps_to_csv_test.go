package psn

import (
	"fmt"
	"testing"

	"github.com/coreos/etcd/pkg/netutil"
)

func Test_getRow(t *testing.T) {
	nt, err := netutil.GetDefaultInterface()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nt)
}
