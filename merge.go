package xfs

import (
	"errors"
	"fmt"
	"io/fs"
)

type mergeFS struct {
	fsyss []fs.FS
}

func Merge(fsyss ...fs.FS) fs.FS {
	return &mergeFS{
		fsyss: fsyss,
	}
}

func (fsys *mergeFS) Open(name string) (fs.File, error) {
	for _, x := range fsys.fsyss {
		f, err := x.Open(name)
		if err != nil {
			continue
		}
		return &mergeReadDirFile{
			File: f,
			name: name,
			fsys: fsys,
		}, nil
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

func (fsys *mergeFS)

type mergeReadDirFile struct {
	fs.File
	name string
	fsys *mergeFS
}

var _ fs.ReadDirFile = (*mergeReadDirFile)(nil)

func (f *mergeReadDirFile) Stat() (fs.FileInfo, error) {
	return f.File.Stat()
}

func (f *mergeReadDirFile) Read(b []byte) (int, error) {
	return f.File.Read(b)
}

func (f *mergeReadDirFile) Close() error {
	return f.File.Close()
}

var errReadDirNNotSupported = fmt.Errorf("ReadDir(n) where n > 0 not supported by fscomb: %w", fs.ErrInvalid)
// var errNotDir = errors.New("not a directory")

func (f *mergeReadDirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n > 0 {
		return nil, errReadDirNNotSupported
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errNotDir
	}

	return f.fsys.ReadDir(f.name)
}
