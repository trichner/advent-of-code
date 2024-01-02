package in

import (
	"io/fs"
)

func MustOpenExampleTxt(fs fs.FS) fs.File {
	return MustOpen(fs, "example.txt")
}

func MustOpenInputTxt(fs fs.FS) fs.File {
	return MustOpen(fs, "input.txt")
}

func MustOpen(fs fs.FS, f string) fs.File {
	file, err := fs.Open(f)
	if err != nil {
		panic(err)
	}
	return file
}
