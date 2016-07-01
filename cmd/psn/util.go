package main

import "os"

const (
	// privateFileMode grants owner to read/write a file.
	privateFileMode = 0600

	// privateDirMode grants owner to make/remove files inside the directory.
	privateDirMode = 0700
)

func openToRead(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDONLY, privateFileMode)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func openToAppend(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, privateFileMode)
	if err != nil {
		return nil, err
	}
	return f, nil
}
