package ps

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestReadCSVs(t *testing.T) {
	testPaths := []string{}
	for i := 0; i < 3; i++ {
		testCSVPath := fmt.Sprintf("test_%d.csv", i)
		testPaths = append(testPaths, testCSVPath)
		defer os.RemoveAll(testCSVPath)

		for j := 0; j < 3; j++ {
			pss, err := List(&Process{Stat: Stat{Pid: int64(1)}})
			if err != nil {
				t.Fatal(err)
			}
			f, err := openToAppend(testCSVPath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if err := WriteToCSV(f, pss...); err != nil {
				t.Fatal(err)
			}
			fmt.Println("sleeping...")
			time.Sleep(time.Second)
		}
	}

	tb, err := ReadCSVs(ColumnsPS, testPaths...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tb.Columns)
	for _, row := range tb.Rows {
		fmt.Printf("%q\n", row)
	}
	fmt.Println(tb.ToRows())
}

func TestReadCSVsTestdata(t *testing.T) {
	testPaths := []string{"testdata/test-01-etcd-server-1.csv", "testdata/test-01-etcd-server-2.csv", "testdata/test-01-etcd-server-3.csv"}
	tb, err := ReadCSVs(ColumnsPS, testPaths...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tb.Columns)
	for _, row := range tb.Rows {
		fmt.Printf("%q\n", row)
	}
	if err := tb.ToCSV("testdata/test.csv"); err != nil {
		log.Fatal(err)
	}
}

func TestReadCSVFillIn(t *testing.T) {
	tb, err := ReadCSVFillIn("./testdata/missing-monitor.csv")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tb.Columns)
	for _, row := range tb.Rows {
		fmt.Printf("%q\n", row)
	}
	if len(tb.Rows) != 18 {
		t.Fatalf("expected %d rows, got %d rows", len(tb.Rows))
	}
}
