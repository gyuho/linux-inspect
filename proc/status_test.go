package proc

import (
	"fmt"
	"testing"
)

func TestGetProcStatusByPID(t *testing.T) {
	rs, err := GetProcStatusByPID(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs.VmRSS)
	fmt.Println(rs.VmRSSBytesN)
	fmt.Println(rs.VmRSSParsedBytes)
	fmt.Printf("GetProcStatusByPID: %+v\n", rs)
}

func TestGetProgram(t *testing.T) {
	fds, err := ListFds()
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
