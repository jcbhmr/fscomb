package xfs

import "io/fs"

type copyFS interface {
	fs.FS
	Copy(oldname string, newname string) error
}

func Copy(fsys fs.FS, oldname string, newname string) error {
	if copyFS, ok := fsys.(copyFS); ok {
		return copyFS.Copy(oldname, newname)
	}

}
