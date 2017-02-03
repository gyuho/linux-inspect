package psn

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	humanize "github.com/dustin/go-humanize"
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

func getWritten(t *testing.T, targetDevice string) (uint64, uint64) {
	dss, err := GetProcDiskstats()
	if err != nil {
		t.Error(err)
	}
	var ds DiskStat
	for _, elem := range dss {
		if elem.DeviceName == targetDevice {
			ds = elem
			break
		}
	}
	if ds.DeviceName == "" {
		t.Skipf("disk stat is not found for device %q", targetDevice)
	}
	return ds.WritesCompleted, ds.SectorsWritten
}

const minSectorSize = 512

// TODO: use tmpfs
// sudo /usr/local/go/bin/go test -v -run TestGetProcDiskstatsSectorWrite
func TestGetProcDiskstatsSectorWrite(t *testing.T) {
	dn, err := GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	fpath := filepath.Join("/boot", "test-temp-file")
	fmt.Println("writing to", fpath, "with device", dn)
	defer os.RemoveAll(fpath)

	oldWritesCompleted, oldSectorWritten := getWritten(t, dn)

	f, err := openToOverwrite(fpath)
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}

	var sum uint64
	for i := 0; i < 1000; i++ {
		n, err := f.Write(bytes.Repeat([]byte{50}, 100*minSectorSize))
		if err != nil {
			t.Fatal(err)
		}
		sum += uint64(n)
	}
	fmt.Println("written", humanize.Bytes(sum))

	if err = f.Sync(); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	newWritesCompleted, newSectorWritten := getWritten(t, dn)

	// sector delta >= write delta
	// because one write size can have 100*minSectorSize
	// e.g. if data 100-byte is written and the sector size is 10-byte
	// then writes completed increases 1,but sector written increases 10
	deltaWrites := newWritesCompleted - oldWritesCompleted
	deltaSector := newSectorWritten - oldSectorWritten
	if deltaSector < deltaWrites {
		t.Fatalf("expected sector delta %d >= writes delta %d", deltaSector, deltaWrites)
	}
	fmt.Printf("writes completed: %d\n", deltaWrites)
	fmt.Printf("sector written: %d\n", deltaSector)
}
