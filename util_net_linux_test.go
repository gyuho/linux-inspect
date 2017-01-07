package psn

import (
	"fmt"
	"testing"
)

func TestGetDefaultInterface(t *testing.T) {
	s, err := GetDefaultInterface()
	fmt.Println(s, err)
}
