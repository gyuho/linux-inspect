package psn

import (
	"fmt"
	"testing"
)

func TestGetNetstat(t *testing.T) {
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

	ns, err := GetNetstat(pid, TCP)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetNetstat TCP: %+v\n", ns)

	ns, err = GetNetstat(pid, TCP6)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetNetstat TCP6: %+v\n", ns)

	pid2 := SearchInode(fds, ns[0].Inode)
	fmt.Println("PID from Inode:", pid2)
}
