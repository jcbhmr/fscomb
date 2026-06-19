package xfs

import (
	"io/fs"
	"syscall"
)

var (
	ErrNotDir = errNotDir{}
	ErrIsDir  = errIsDir{}
)

type errNotDir struct{}
type errIsDir struct{}

func (errNotDir) Error() string {
	return "not a directory"
}
func (errIsDir) Error() string {
	return "is a directory"
}

func (errNotDir) Is(target error) bool {
	return target == syscall.ENOTDIR
}
func (errIsDir) Is(target error) bool {
	return target == syscall.EISDIR
}

type SameFileFS interface {
	fs.FS
	SameFile(fi1, fi2 fs.FileInfo) bool
}

func SameFile(fsys fs.FS, fi1, fi2 fs.FileInfo) bool {
	if sameFileFS, ok := fsys.(SameFileFS); ok {
		return sameFileFS.SameFile(fi1, fi2)
	}
	return false
}
