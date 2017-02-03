package psn

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestGetProcDiskstats(t *testing.T) {
	dss, err := GetProcDiskstats()
	if err != nil {
		t.Error(err)
	}
	for _, ds := range dss {
		if ds.ReadsCompleted == 0 {
			continue
		}
		fmt.Printf("%s %d\n", ds.DeviceName, ds.ReadsCompleted)
	}
}

/*
TODO: want to see disk writes increases the sector writes!

just run with

sudo /usr/local/go/bin/go test -v -run TestGetProcDiskstatsSectorWrite
*/
func TestGetProcDiskstatsSectorWrite(t *testing.T) {
	dn, err := GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	fpath := filepath.Join("/boot", "test-temp-file")
	fmt.Println("writing to", fpath, "with device", dn)
	// defer os.RemoveAll(fpath)

	ds1, err := GetProcDiskstats()
	if err != nil {
		t.Error(err)
	}
	var diskStat1 DiskStat
	for _, elem := range ds1 {
		if elem.DeviceName == dn {
			diskStat1 = elem
			break
		}
	}
	if diskStat1.DeviceName == "" {
		t.Skipf("disk stat is not found for device %q", dn)
	}

	f, err := openToOverwrite(fpath)
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}

	const minSectorSize = 512
	const walPageBytes = 8 * minSectorSize
	n, err := f.Write(bytes.Repeat([]byte{50}, 10*walPageBytes))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("wrote", n, "bytes")

	// if err = f.Close(); err != nil {
	// 	t.Fatal(err)
	// }
	if err = fdatasync(f); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	ds2, err := GetProcDiskstats()
	if err != nil {
		t.Error(err)
	}
	var diskStat2 DiskStat
	for _, elem := range ds2 {
		if elem.DeviceName == dn {
			diskStat2 = elem
			break
		}
	}
	if diskStat2.DeviceName == "" {
		t.Skipf("disk stat is not found for device %q", dn)
	}

	fmt.Printf("diskStat1 %+v\n", diskStat1)
	fmt.Printf("diskStat2 %+v\n", diskStat2)
}

// fdatasync flushes all data buffers of a file onto the disk.
// Fsync is required to update the metadata, such as access time.
// Fsync always does two write operations: one for writing new data
// to disk. Another for updating the modification time stored in its
// inode. If the modification time is not a part of the transaction,
// syscall.Fdatasync can be used to avoid unnecessary inode disk writes.
//
// (etcd pkg.fileutil.Fdatasync)
func fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}
