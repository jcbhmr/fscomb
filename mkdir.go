package xfs

import (
	"errors"
	"io/fs"
	"strings"
)

type MkdirFS interface {
	fs.FS
	Mkdir(name string, perm fs.FileMode) error
}

func Mkdir(fsys fs.FS, name string, perm fs.FileMode) error {
	mkdirFS, ok := fsys.(MkdirFS)
	if !ok {
		return &fs.PathError{Op: "mkdir", Path: name, Err: errors.ErrUnsupported}
	}
	return mkdirFS.Mkdir(name, perm)
}

type mkdirAllFS interface {
	MkdirFS
	MkdirAll(name string, perm fs.FileMode) error
}

func MkdirAll(fsys fs.FS, name string, perm fs.FileMode) error {
	mkdirFS, ok := fsys.(MkdirFS)
	if !ok {
		return &fs.PathError{Op: "mkdir", Path: name, Err: errors.ErrUnsupported}
	}

	if mkdirAllFS, ok := mkdirFS.(mkdirAllFS); ok {
		return mkdirAllFS.MkdirAll(name, perm)
	}

	components := strings.Split(name, "/")
	for i := range components {
		name := strings.Join(components[:i], "/")
		err := Mkdir(fsys, name, perm)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}
	}
	return nil
}
