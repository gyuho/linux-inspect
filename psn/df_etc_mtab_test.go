package psn

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

func TestGetEtcMtab(t *testing.T) {
	mss, err := GetEtcMtab()
	if err != nil {
		t.Error(err)
	}
	for _, ms := range mss {
		fmt.Printf("%+v\n", ms)
	}
}

func TestGetDevice(t *testing.T) {
	s, err := GetDevice("/boot")
	fmt.Println(s, err)
}
