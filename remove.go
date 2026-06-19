package xfs

import (
	"errors"
	"io/fs"
	"path"
)

type RemoveFS interface {
	fs.FS
	Remove(name string) error
}

func Remove(fsys fs.FS, name string) error {
	removeFS, ok := fsys.(RemoveFS)
	if !ok {
		return &fs.PathError{Op: "remove", Path: name, Err: errors.ErrUnsupported}
	}
	return removeFS.Remove(name)
}

type removeAllFS interface {
	fs.FS
	RemoveAll(name string) error
}

func RemoveAll(fsys fs.FS, name string) error {
	if removeAllFS, ok := fsys.(removeAllFS); ok {
		return removeAllFS.RemoveAll(name)
	}

	removeFS, ok := fsys.(RemoveFS)
	if !ok {
		return &fs.PathError{Op: "remove", Path: name, Err: errors.ErrUnsupported}
	}

	err := Remove(removeFS, name)
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	entries, err := ReadDir(fsys, name)
	if err != nil {
		return err
	}
	for _, e := range entries {
		err := RemoveAll(removeFS, path.Join(name, e.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
