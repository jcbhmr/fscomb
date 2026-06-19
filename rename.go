package xfs

import "io/fs"

type RenameFS interface {
	fs.FS
	Rename(oldname, newname string) error
}

func Rename(fsys fs.FS, oldname, newname string) error {
	if renameFS, ok := fsys.(RenameFS); ok {
		return renameFS.Rename(oldname, newname)
	}

	err := Copy(fsys, oldname, newname)
	if err != nil {
		return err
	}
	return Remove(fsys, oldname)
}
