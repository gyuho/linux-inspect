package psn

import (
	"fmt"
	"testing"
)

func TestGetDevice(t *testing.T) {
	s, err := GetDevice("/boot")
	fmt.Println(s, err)
}
