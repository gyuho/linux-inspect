package ps

import (
	"os"
	"strconv"
	"strings"
)

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func openToRead(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0444)
	if err != nil {
		return f, err
	}
	return f, nil
}

func openToAppend(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}

func openToOverwrite(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}

func ToField(s string) string {
	cs := strings.Split(s, "_")
	var ss []string
	for _, v := range cs {
		if len(v) > 0 {
			ss = append(ss, strings.Title(v))
		}
	}
	return strings.TrimSpace(strings.Join(ss, ""))
}
