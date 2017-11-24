package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/gyuho/linux-inspect/df"
	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/proc"

	"github.com/coreos/etcd/pkg/netutil"
)

func TestCombine(t *testing.T) {
	dn, err := df.GetDevice("/boot")
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	nm, err := netutil.GetDefaultInterfaces()
	if err != nil {
		fmt.Println(err)
		t.Skip()
	}
	var nt string
	for k := range nm {
		nt = k
		break
	}
	fmt.Println("running with network interface", nt, "and disk device", dn)
	if dn == "overlay" {
		t.Skipf("TODO: overlay is not supported yet")
	}

	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("test-%010d.csv", time.Now().UnixNano()))
	defer os.RemoveAll(fpath)

	epath := filepath.Join(homeDir(), "etcd-client-num")
	defer os.RemoveAll(epath)

	if err = fileutil.ToFile("10", epath); err != nil {
		t.Fatal(err)
	}
	c, err := NewCSV(fpath, 1, dn, nt, epath, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if i > 0 {
			if err := fileutil.ToFile(fmt.Sprintf("%d", 100*i), epath); err != nil {
				t.Fatal(err)
			}
		}

		now := time.Now()
		if err := c.Add(); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("#%d: collected data with %s and %s at %s (c.Add took %v)\n", i, dn, nt, fpath, time.Since(now))

		time.Sleep(time.Second)
	}

	combined := Combine(c.Rows...)
	if combined.UnixSecond != c.Rows[len(c.Rows)-1].UnixSecond {
		t.Fatalf("unix second expected %d, got %d", combined.UnixSecond, c.Rows[len(c.Rows)-1].UnixSecond)
	}

	// test cumulative fields
	if combined.NSEntry.ReceiveBytesNum > c.Rows[len(c.Rows)-1].NSEntry.ReceiveBytesNum {
		t.Fatalf("averaged combined field should have equal or smaller receive combined.NSEntry.ReceiveBytesNum than the last one (last %v, got %v)", c.Rows[len(c.Rows)-1].NSEntry.ReceiveBytesNum, combined.NSEntry.ReceiveBytesNum)
	}
	if combined.NSEntry.TransmitBytesNum > c.Rows[len(c.Rows)-1].NSEntry.TransmitBytesNum {
		t.Fatalf("averaged combined field should have equal or smaller receive combined.NSEntry.TransmitBytesNum than the last one (last %v, got %v)", c.Rows[len(c.Rows)-1].NSEntry.TransmitBytesNum, combined.NSEntry.TransmitBytesNum)
	}

	fmt.Println()
	for i, row := range c.Rows {
		fmt.Printf("row #%d: %+v\n", i, row)
	}
	fmt.Println()
	fmt.Printf("combined: %+v\n", combined)
}

func TestInterpolateMissingOne(t *testing.T) {
	cc1 := &CSV{
		MinUnixNanosecond: 10,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 13,
		MaxUnixSecond:     4,
		Rows: []Proc{
			{
				UnixNanosecond: 10,
				UnixSecond:     1,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  10.0,
					LoadAvg15Minute: 20.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
				TransmitBytesNumDelta: 244,
			},
			{
				UnixNanosecond: 12,
				UnixSecond:     3,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  40.0,
					LoadAvg15Minute: 130.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 50000,
					SectorsWritten:  100000,
				},
				TransmitBytesNumDelta: 74,
			},
			{
				UnixNanosecond: 13,
				UnixSecond:     4,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  50.0,
					LoadAvg15Minute: 150.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 100000,
					SectorsWritten:  200000,
				},
			},
		},
	}

	expectedRowN := int(cc1.MaxUnixSecond - cc1.MinUnixSecond + 1)

	cc2, err := cc1.Interpolate()
	if err != nil {
		t.Fatal(err)
	}

	if len(cc2.Rows) != expectedRowN {
		t.Fatalf("len(cc2.Rows) expected %d, got %d", len(cc2.Rows), expectedRowN)
	}
	if !reflect.DeepEqual(cc1.Rows[0].DSEntry, cc2.Rows[0].DSEntry) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[0].DSEntry, cc2.Rows[0].DSEntry)
	}
	if !reflect.DeepEqual(cc1.Rows[len(cc1.Rows)-1].DSEntry, cc2.Rows[len(cc2.Rows)-1].DSEntry) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[len(cc1.Rows)-1].DSEntry, cc2.Rows[len(cc2.Rows)-1].DSEntry)
	}

	expected2 := &CSV{
		MinUnixNanosecond: 0,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 0,
		MaxUnixSecond:     4,
		Rows: []Proc{
			{
				UnixNanosecond: 0,
				UnixSecond:     1,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  10.0,
					LoadAvg15Minute: 20.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
				TransmitBytesNumDelta: 244,
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     2,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  25.0,
					LoadAvg15Minute: 75.0,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    30000,
					SectorsWritten:     60000,
				},
				NSEntry:               NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:     "0 B",
				TransmitBytesDelta:    "159 B",
				TransmitBytesNumDelta: 159,
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     3,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  40.0,
					LoadAvg15Minute: 130.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 50000,
					SectorsWritten:  100000,
				},
				TransmitBytesNumDelta: 74,
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     4,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  50.0,
					LoadAvg15Minute: 150.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 100000,
					SectorsWritten:  200000,
				},
			},
		},
	}

	if !reflect.DeepEqual(expected2, cc2) {
		for i := range cc2.Rows {
			if !reflect.DeepEqual(expected2.Rows[i], cc2.Rows[i]) {
				t.Fatalf("#%d: expected %+v, got %+v", i, expected2.Rows[i], cc2.Rows[i])
			}
		}
		t.Fatalf("expected %+v, got %+v", expected2, cc2)
	}

	// make sure interpolated CSV was deep-copied
	old := cc1.DiskDevice
	cc2.DiskDevice = "test"
	if cc1.DiskDevice == "test" {
		t.Fatalf("cc1.DiskDevice expected %q, got 'test'", old)
	}
}

func TestInterpolateMissingALot(t *testing.T) {
	cc1 := &CSV{
		MinUnixNanosecond: 10,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 100,
		MaxUnixSecond:     10,
		Rows: []Proc{
			{
				UnixNanosecond: 10,
				UnixSecond:     1,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  5.0,
					LoadAvg15Minute: 15.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
				WriteBytesDelta: 50,
			},
			{
				UnixNanosecond: 100,
				UnixSecond:     10,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  50.0,
					LoadAvg15Minute: 150.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 100000,
					SectorsWritten:  200000,
				},
				WriteBytesDelta: 5000,
			},
		},
	}

	expectedRowN := int(cc1.MaxUnixSecond - cc1.MinUnixSecond + 1)

	cc2, err := cc1.Interpolate()
	if err != nil {
		t.Fatal(err)
	}

	if len(cc2.Rows) != expectedRowN {
		t.Fatalf("len(cc2.Rows) expected %d, got %d", len(cc2.Rows), expectedRowN)
	}
	if !reflect.DeepEqual(cc1.Rows[0].DSEntry, cc2.Rows[0].DSEntry) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[0].DSEntry, cc2.Rows[0].DSEntry)
	}
	if !reflect.DeepEqual(cc1.Rows[len(cc1.Rows)-1].DSEntry, cc2.Rows[len(cc2.Rows)-1].DSEntry) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[len(cc1.Rows)-1].DSEntry, cc2.Rows[len(cc2.Rows)-1].DSEntry)
	}

	expected2 := &CSV{
		MinUnixNanosecond: 0,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 0,
		MaxUnixSecond:     10,
		Rows: []Proc{
			{
				UnixNanosecond: 0,
				UnixSecond:     1,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  5.0,
					LoadAvg15Minute: 15.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
				WriteBytesDelta: 50,
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     2,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  10,
					LoadAvg15Minute: 30,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    20000,
					SectorsWritten:     40000,
				},
				WriteBytesDelta:    600,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     3,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  15,
					LoadAvg15Minute: 45,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    30000,
					SectorsWritten:     60000,
				},
				WriteBytesDelta:    1150,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     4,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  20.0,
					LoadAvg15Minute: 60.0,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    40000,
					SectorsWritten:     80000,
				},
				WriteBytesDelta:    1700,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     5,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  25,
					LoadAvg15Minute: 75,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    50000,
					SectorsWritten:     100000,
				},
				WriteBytesDelta:    2250,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     6,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  30.0,
					LoadAvg15Minute: 90.0,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    60000,
					SectorsWritten:     120000,
				},
				WriteBytesDelta:    2800,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     7,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  35,
					LoadAvg15Minute: 105,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    70000,
					SectorsWritten:     140000,
				},
				WriteBytesDelta:    3350,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     8,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  40.0,
					LoadAvg15Minute: 120.0,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    80000,
					SectorsWritten:     160000,
				},
				WriteBytesDelta:    3900,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     9,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  45,
					LoadAvg15Minute: 135,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    90000,
					SectorsWritten:     180000,
				},
				WriteBytesDelta:    4450,
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			{
				UnixNanosecond: 0,
				UnixSecond:     10,
				LoadAvg: proc.LoadAvg{
					LoadAvg1Minute:  50.0,
					LoadAvg15Minute: 150.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 100000,
					SectorsWritten:  200000,
				},
				WriteBytesDelta: 5000,
			},
		},
	}

	if !reflect.DeepEqual(expected2, cc2) {
		for i := range cc2.Rows {
			if !reflect.DeepEqual(expected2.Rows[i], cc2.Rows[i]) {
				t.Fatalf("#%d: expected %+v, got %+v", i, expected2.Rows[i], cc2.Rows[i])
			}
		}
		t.Fatalf("expected %+v, got %+v", expected2, cc2)
	}

	// make sure interpolated CSV was deep-copied
	old := cc1.DiskDevice
	cc2.DiskDevice = "test"
	if cc1.DiskDevice == "test" {
		t.Fatalf("cc1.DiskDevice expected %q, got 'test'", old)
	}
}
