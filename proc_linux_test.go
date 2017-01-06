package psn

import (
	"fmt"
	"testing"
)

func TestGetProgram(t *testing.T) {
	fds, err := ListProcFds()
	if err != nil {
		t.Fatal(err)
	}

	fd := fds[len(fds)/2]
	fmt.Println("Chosen FD:", fd)

	pid, err := pidFromFd(fd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Chosen PID:", pid)

	nm, err := GetProgram(pid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("GetProgram:", nm)
}
