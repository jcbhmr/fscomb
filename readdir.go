package xfs

import (
	"errors"
	"io/fs"
)

func ReadDir(fsys fs.FS, name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(fsys, name)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, pathError) && pathError.Op == "readdir" && pathError.Path == name && pathError.Err.Error() == "not implemented" {
			pathError.Err = errors.ErrUnsupported
		}
	}
	return entries, err
}
