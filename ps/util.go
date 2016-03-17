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

func open(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0444)
	if err != nil {
		return f, err
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
	return strings.Join(ss, "")
}
