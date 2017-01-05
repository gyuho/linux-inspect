package psn

import (
	"os"
	"strconv"
)

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func openToRead(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func openToAppend(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func openToOverwrite(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return f, nil
}
