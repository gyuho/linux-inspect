package psn

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/etcd/pkg/netutil"
)

func TestCombine(t *testing.T) {
	dn, err := GetDevice("/boot")
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
	fmt.Println("running with network interface", nt)

	fpath := filepath.Join(os.TempDir(), fmt.Sprintf("test-%010d.csv", time.Now().UnixNano()))
	defer os.RemoveAll(fpath)

	epath := filepath.Join(homeDir(), "etcd-client-num")
	defer os.RemoveAll(epath)

	if err := toFile([]byte("10"), epath); err != nil {
		t.Fatal(err)
	}
	c := NewCSV(fpath, 1, dn, nt, epath)
	for i := 0; i < 3; i++ {
		if i > 0 {
			if err := toFile([]byte(fmt.Sprintf("%d", 100*i)), epath); err != nil {
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

func TestInterpolate(t *testing.T) {
	cc1 := &CSV{
		MinUnixNanosecond: 10,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 100,
		MaxUnixSecond:     10,
		Rows: []Proc{
			Proc{
				UnixNanosecond: 10,
				UnixSecond:     1,
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  5.0,
					LoadAvg15Minute: 15.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
			},
			Proc{
				UnixNanosecond: 100,
				UnixSecond:     10,
				LoadAvg: LoadAvg{
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
	if !reflect.DeepEqual(cc1.Rows[0], cc2.Rows[0]) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[0], cc2.Rows[0])
	}
	if !reflect.DeepEqual(cc1.Rows[len(cc1.Rows)-1], cc2.Rows[len(cc2.Rows)-1]) {
		t.Fatalf("first row expected %+v, got %+v", cc1.Rows[len(cc1.Rows)-1], cc2.Rows[len(cc2.Rows)-1])
	}

	expected2 := &CSV{
		MinUnixNanosecond: 0,
		MinUnixSecond:     1,
		MaxUnixNanosecond: 0,
		MaxUnixSecond:     10,
		Rows: []Proc{
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     1,
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  5.0,
					LoadAvg15Minute: 15.0,
				},
				DSEntry: DSEntry{
					WritesCompleted: 10000,
					SectorsWritten:  20000,
				},
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     2,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  9.5,
					LoadAvg15Minute: 28.5,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    19000,
					SectorsWritten:     38000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     3,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  14,
					LoadAvg15Minute: 42,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    28000,
					SectorsWritten:     56000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     4,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  18.5,
					LoadAvg15Minute: 55.5,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    37000,
					SectorsWritten:     74000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     5,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  23,
					LoadAvg15Minute: 69,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    46000,
					SectorsWritten:     92000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     6,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  27.5,
					LoadAvg15Minute: 82.5,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    55000,
					SectorsWritten:     110000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     7,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  32,
					LoadAvg15Minute: 96,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    64000,
					SectorsWritten:     128000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     8,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  36.5,
					LoadAvg15Minute: 109.5,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    73000,
					SectorsWritten:     146000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     9,
				PSEntry:        PSEntry{CPU: "0.00 %", VMRSS: "0 B", VMSize: "0 B"},
				LoadAvg: LoadAvg{
					LoadAvg1Minute:  41,
					LoadAvg15Minute: 123,
				},
				DSEntry: DSEntry{
					TimeSpentOnReading: "0 seconds",
					TimeSpentOnWriting: "0 seconds",
					WritesCompleted:    82000,
					SectorsWritten:     164000,
				},
				NSEntry:            NSEntry{ReceiveBytes: "0 B", TransmitBytes: "0 B"},
				ReceiveBytesDelta:  "0 B",
				TransmitBytesDelta: "0 B",
			},
			Proc{
				UnixNanosecond: 0,
				UnixSecond:     10,
				LoadAvg: LoadAvg{
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

	if !reflect.DeepEqual(cc2, expected2) {
		for _, row := range cc2.Rows {
			fmt.Printf("%+v\n", row)
		}
		t.Fatal("unexpected CSV")
	}
}
