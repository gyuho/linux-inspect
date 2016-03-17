package ps

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	rs, err := List(nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", rs)
	fmt.Println(len(rs))
}

func TestStatusByPID(t *testing.T) {
	rs, err := StatusByPID(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", rs)
}
