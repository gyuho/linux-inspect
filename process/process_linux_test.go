package process

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	pss, err := List(nil)
	if err != nil && !strings.Contains(err.Error(), "too many open files") {
		t.Error(err)
	}
	WriteToTable(os.Stdout, 0, pss...)
	fmt.Println(len(pss))
}
