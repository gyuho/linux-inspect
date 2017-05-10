package top

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	now := time.Now()
	rows, err := Get(DefaultExecPath, 0)
	if err != nil {
		t.Skip(err)
	}
	for _, elem := range rows {
		fmt.Printf("%+v\n", elem)
	}
	fmt.Printf("found %d entrines in %v", len(rows), time.Since(now))
}
