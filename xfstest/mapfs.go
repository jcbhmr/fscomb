package xfstest

import (
	"io/fs"
	"testing/fstest"
	"time"
)

type MapFS struct{ fstest.MapFS }

func (fsys MapFS) Open() (fs.File, error) {
}
