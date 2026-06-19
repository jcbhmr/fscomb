package xfs

import (
	"errors"
	"io/fs"
	"time"
)

type Flag uint32

const (
	FlagReadOnly  Flag = 0x00
	FlagWriteOnly Flag = 0x01
	FlagReadWrite Flag = 0x02
	FlagAppend    Flag = 0x400
	FlagTruncate  Flag = 0x200
	FlagCreate    Flag = 0x40
	FlagExclusive Flag = 0x80
	FlagSync      Flag = 0x101000
)

type OpenFileFS interface {
	fs.FS
	OpenFile(name string, flag Flag, perm fs.FileMode) (fs.File, error)
}

func OpenFile(fsys fs.FS, name string, flag Flag, perm fs.FileMode) (fs.File, error) {
	if openFileFS, ok := fsys.(OpenFileFS); ok {
		return openFileFS.OpenFile(name, flag, perm)
	}
	if flag != FlagReadOnly {
		return nil, &fs.PathError{Op: "open", Path: name, Err: errors.ErrUnsupported}
	}
	return fsys.Open(name)
}

type createFS interface {
	fs.FS
	Create(name string) (fs.File, error)
}

func Create(fsys fs.FS, name string) (fs.File, error) {
	if createFS, ok := fsys.(createFS); ok {
		return createFS.Create(name)
	}

	return OpenFile(fsys, name, FlagReadWrite|FlagCreate|FlagTruncate, 0o666)
}

type ChtimesFS interface {
	Chtimes(name string, atime, mtime time.Time) error
}

type WriteFile interface {
	fs.File
	Write([]byte) (int, error)
}

type ChmodFile interface {
	fs.File
	Chmod(mode fs.FileMode) error
}

type ChownFile interface {
	fs.File
	Chown(uid, gid int) error
}

type TruncateFile interface {
	fs.File
	Truncate(size int64) error
}
