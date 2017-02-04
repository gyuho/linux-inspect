package psn

import (
	"fmt"
	"os"
	"testing"
)

func TestGetPS(t *testing.T) {
	ns, err := GetPS(WithTopLimit(3))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertPS(ns...)
	txt := StringPS(hd, rows, -1)
	fmt.Println(txt)
}

func TestGetPSWithFilter(t *testing.T) {
	pid := int64(os.Getpid())

	ns, err := GetPS(WithPID(pid))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertPS(ns...)
	txt := StringPS(hd, rows, -1)
	fmt.Println(txt)
}

func TestGetPSWithTopStream(t *testing.T) {
	pid := int64(os.Getpid())

	cfg := &TopConfig{
		Exec:           DefaultTopPath,
		IntervalSecond: 1,
		PID:            pid,
	}
	str, err := cfg.StartStream()
	if err != nil {
		t.Fatal(err)
	}

	ns, err := GetPS(WithPID(pid), WithTopStream(str))
	if err != nil {
		t.Fatal(err)
	}
	hd, rows := ConvertPS(ns...)
	txt := StringPS(hd, rows, -1)
	fmt.Println("ps:")
	fmt.Println(txt)

	if err = str.Stop(); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-str.ErrChan():
		t.Fatal(err)
	default:
		fmt.Println("'top' stopped")
	}

	rm := str.Latest()
	for _, row := range rm {
		fmt.Printf("%+v\n", row)
	}
	fmt.Println("total", len(rm), "processes")
}
