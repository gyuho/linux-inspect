package psn

import (
	"fmt"
	"testing"
)

func TestListPIDs(t *testing.T) {
	pids, err := ListPIDs()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ListPIDs:", pids)
}

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

func TestGetDevice(t *testing.T) {
	s, err := GetDevice("/")
	fmt.Println(s, err)
}
