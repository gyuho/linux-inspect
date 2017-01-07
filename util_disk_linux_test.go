package psn

import (
	"fmt"
	"testing"
)

func TestGetDevice(t *testing.T) {
	s, err := GetDevice("/")
	fmt.Println(s, err)
}
