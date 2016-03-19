package ps

import (
	"fmt"
	"os"
	"testing"
)

func TestList(t *testing.T) {
	pss, err := List(nil)
	if err != nil {
		t.Error(err)
	}
	WriteToTable(os.Stdout, 0, pss...)
	fmt.Println(len(pss))
}
