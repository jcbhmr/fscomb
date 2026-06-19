package xos

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"go.jcbhmr.com/xfs"
)

type osDirFS interface {
	fs.FS
	fs.StatFS
	fs.ReadFileFS
	fs.ReadDirFS
	fs.ReadLinkFS
}

type dirFS struct {
	osDirFS
	dir string
}

func DirFS(dir string) xfs.WriteFS {
	return &dirFS{
		osDirFS: os.DirFS(dir).(osDirFS),
		dir:     dir,
	}
}

var errEmptyRoot = errors.New("DirFS with empty root")

func (fsys dirFS) join(name string) (string, error) {
	if fsys.dir == "" {
		return "", errEmptyRoot
	}
	name, err := filepath.Localize(name)
	if err != nil {
		return "", fs.ErrInvalid
	}
	if os.IsPathSeparator(fsys.dir[len(fsys.dir)-1]) {
		return fsys.dir + name, nil
	}
	return fsys.dir + string(os.PathSeparator) + name, nil
}

func osFlag(flag xfs.OpenFlag) int {
	var target int

	read := flag&xfs.FlagRead != 0
	write := flag&xfs.FlagWrite != 0
	if read && write {
		target |= os.O_RDWR
	} else if read {
		target |= os.O_RDONLY
	} else if write {
		target |= os.O_WRONLY
	}

	if flag&xfs.FlagCreate != 0 {
		target |= os.O_CREATE
	}

	return target
}

func (fsys dirFS) OpenFile(name string, flag xfs.OpenFlag, perm fs.FileMode) (fs.File, error) {
	fullname, err := fsys.join(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	f, err := os.OpenFile(fullname, osFlag(flag), perm)
	if err != nil {
		err.(*fs.PathError).Path = name
		return nil, err
	}
	return f, nil
}

func (fsys dirFS) Mkdir(name string, perm fs.FileMode) error {
	fullname, err := fsys.join(name)
	if err != nil {
		return &fs.PathError{Op: "mkdir", Path: name, Err: err}
	}
	err = os.Mkdir(fullname, perm)
	if err != nil {
		err.(*fs.PathError).Path = name
		return err
	}
	return nil
}

func (fsys dirFS) Remove(name string) error {
	fullname, err := fsys.join(name)
	if err != nil {
		return &fs.PathError{Op: "remove", Path: name, Err: err}
	}
	err = os.Remove(fullname)
	if err != nil {
		err.(*fs.PathError).Path = name
		return err
	}
	return nil
}

func (fsys dirFS) Rename(oldname string, newname string) error {
	fulloldname, olderr := fsys.join(oldname)
	fullnewname, newerr := fsys.join(newname)
	err := errors.Join(olderr, newerr)
	if err != nil {
		return &os.LinkError{Op: "rename", Old: oldname, New: newname, Err: err}
	}
	err = os.Rename(fulloldname, fullnewname)
	if err != nil {
		err.(*os.LinkError).Old = oldname
		err.(*os.LinkError).New = newname
		return err
	}
	return nil
}

func (fsys dirFS) Link(oldname string, newname string) error {
	fulloldname, olderr := fsys.join(oldname)
	fullnewname, newerr := fsys.join(newname)
	err := errors.Join(olderr, newerr)
	if err != nil {
		return &os.LinkError{Op: "link", Old: oldname, New: newname, Err: err}
	}
	err = os.Link(fulloldname, fullnewname)
	if err != nil {
		err.(*os.LinkError).Old = oldname
		err.(*os.LinkError).New = newname
		return err
	}
	return nil
}

func (fsys dirFS) Symlink(oldname string, newname string) error {
	fulloldname, olderr := fsys.join(oldname)
	fullnewname, newerr := fsys.join(newname)
	err := errors.Join(olderr, newerr)
	if err != nil {
		return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: err}
	}
	err = os.Symlink(fulloldname, fullnewname)
	if err != nil {
		err.(*os.LinkError).Old = oldname
		err.(*os.LinkError).New = newname
		return err
	}
	return nil
}
