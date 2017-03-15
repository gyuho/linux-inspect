package proc

import (
	"fmt"
	"testing"
)

func TestListPIDs(t *testing.T) {
	pids, err := ListPIDs()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ListPIDs:", pids)
}
