package psn

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func humanizeDurationMs(ms uint64) string {
	s := humanize.Time(time.Now().Add(-1 * time.Duration(ms) * time.Millisecond))
	if s == "now" {
		s = "0 seconds"
	}
	return strings.TrimSpace(strings.Replace(s, " ago", "", -1))
}

func humanizeDurationSecond(sec uint64) string {
	s := humanize.Time(time.Now().Add(-1 * time.Duration(sec) * time.Second))
	if s == "now" {
		s = "0 seconds"
	}
	return strings.TrimSpace(strings.Replace(s, " ago", "", -1))
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

func toFile(data []byte, fpath string) error {
	f, err := openToOverwrite(fpath)
	if err != nil {
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
	}
	_, err = f.Write(data)
	f.Close()
	return err
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
