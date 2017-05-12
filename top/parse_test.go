package top

import "testing"

func TestTop_parseMemoryTxt(t *testing.T) {
	bts, bs, err := parseMemoryTxt("50.883g")
	if err != nil {
		t.Fatal(err)
	}
	if bts != 53687091200 {
		t.Fatalf("bytes expected %d, got %d", 53687091200, bts)
	}
	if bs != "54 GB" {
		t.Fatalf("humanized bytes expected '54 GB', got %q", bs)
	}
}
