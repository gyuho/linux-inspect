package ps

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestReadCSV(t *testing.T) {
	testCSVPath := "test.csv"
	defer os.RemoveAll(testCSVPath)

	for i := 0; i < 5; i++ {
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

	tb, err := ReadCSV(testCSVPath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("columns:", tb.Columns)
	for _, row := range tb.Rows {
		fmt.Println(row)
	}
}

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

	tb, err := ReadCSVs(testPaths...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tb.Columns)
	for _, row := range tb.Rows {
		fmt.Printf("%q\n", row)
	}
	fmt.Println(tb.ToRows())
}
