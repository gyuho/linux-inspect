package proc

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestGetNetTCPByPID(t *testing.T) {
	fds, err := ListFds()
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
		t.Skip(err)
	}
	fmt.Println("GetProgram:", nm)

	ns, err := GetNetTCPByPID(pid, TypeTCP)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetNetTCPByPID TypeTCP: %+v\n", ns)

	nss, err := GetNetTCPByPID(pid, TypeTCP6)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetNetTCPByPID TypeTCP: %+v\n", ns)

	for _, ns := range nss {
		pid2 := searchInode(fds, ns.Inode)
		if pid2 < 0 {
			continue
		}
		pn, err := GetProgram(pid2)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("PID %d for Inode %6s, Program %s\n", pid2, ns.Inode, pn)
		if pn != nm {
			fmt.Printf("program name expected %q, got %q\n", nm, pn)
		}
	}
}

// searchInode finds the matching process to the given inode.
func searchInode(fds []string, inode string) (pid int64) {
	var mu sync.RWMutex

	var wg sync.WaitGroup
	wg.Add(len(fds))
	for _, fd := range fds {
		go func(fdpath string) {
			defer wg.Done()

			mu.RLock()
			done := pid != 0
			mu.RUnlock()
			if done {
				return
			}

			// '/proc/[pid]/fd' contains type:[inode]
			sym, err := os.Readlink(fdpath)
			if err != nil {
				return
			}
			if !strings.Contains(strings.TrimSpace(sym), inode) {
				return
			}

			pd, err := pidFromFd(fdpath)
			if err != nil {
				return
			}
			mu.Lock()
			pid = pd
			mu.Unlock()
		}(fd)
	}
	wg.Wait()

	if pid == 0 {
		pid = -1
	}
	return
}
