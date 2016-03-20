package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gyuho/psn/ps"
)

func main() {
	rsPaths := []string{}
	tb, err := ps.ReadCSVs(rsPaths...)
	if err != nil {
		log.Fatal(err)
	}

	if err != toCSV(tb.Columns, tb.Rows, "results.csv"); err != nil {
		log.Fatal(err)
	}
}

func toCSV(header []string, rows [][]string, fpath string) error {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	// func NewWriter(w io.Writer) *Writer
	wr := csv.NewWriter(f)

	if err := wr.Write(header); err != nil {
		return err
	}

	if err := wr.WriteAll(rows); err != nil {
		return err
	}

	wr.Flush()
	return wr.Error()
}
