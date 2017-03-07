package etc

import (
	"fmt"
	"testing"
)

func TestGetMtab(t *testing.T) {
	mss, err := GetMtab()
	if err != nil {
		t.Error(err)
	}
	for _, ms := range mss {
		fmt.Printf("%+v\n", ms)
	}
}
