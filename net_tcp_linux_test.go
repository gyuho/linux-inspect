package psn

import (
	"fmt"
	"testing"
)

func TestGetNetTCP(t *testing.T) {
	fds, err := ListProcFds()
	if err != nil {
		t.Fatal(err)
	}

	fd := fds[len(fds)/5]
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

	ns, err := GetNetTCP(pid, TypeTCP)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetNetTCP TypeTCP: %+v\n", ns)

	nss, err := GetNetTCP(pid, TypeTCP6)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetNetTCP TypeTCP: %+v\n", ns)

	for _, ns := range nss {
		pid2 := SearchInode(fds, ns.Inode)
		pn, err := GetProgram(pid2)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("PID %d for Inode %6s, Program %s\n", pid2, ns.Inode, pn)
		if pn != nm {
			t.Fatalf("program name expected %q, got %q", nm, pn)
		}
	}
}
