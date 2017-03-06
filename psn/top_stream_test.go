package psn

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTopStartTopStream(t *testing.T) {
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

	time.Sleep(2 * time.Second)

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
