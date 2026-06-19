package xfs

import (
	"io/fs"
)

type filterFS struct {
	fs.FS
	filter func(name string) bool
}

func Filter(fsys fs.FS, filter func(name string) bool) fs.FS {
	return &filterFS{
		FS:     fsys,
		filter: filter,
	}
}

func (fsys *filterFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	if !fsys.filter(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return fsys.FS.Open(name)
}
