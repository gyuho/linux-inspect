package df

import (
	"fmt"
	"testing"
)

func TestGetDfDefault(t *testing.T) {
	dfs, err := GetDfDefault("")
	if err != nil {
		t.Skip(err)
	}
	for _, df := range dfs {
		fmt.Printf("%+v\n", df)
	}

	dfs, err = GetDfDefault(".")
	if err != nil {
		t.Skip(err)
	}
	for _, df := range dfs {
		fmt.Printf("%+v\n", df)
	}
}

func TestGetDevice(t *testing.T) {
	s, err := GetDevice("/boot")
	fmt.Println(s, err)
}
