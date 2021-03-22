package sqlfs

import (
	"os"
	"syscall"
	"time"
)

func newFileStat(name string, st syscall.Stat_t) *fileStat {
	fstat := fileStat{
		name: name,
		sys:  st,
	}
	fillFileStatFromSys(&fstat)
	return &fstat
}

type fileStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     syscall.Stat_t
}

func (fs *fileStat) Name() string       { return fs.name }
func (fs *fileStat) Size() int64        { return fs.size }
func (fs *fileStat) Mode() os.FileMode  { return fs.mode }
func (fs *fileStat) ModTime() time.Time { return fs.modTime }
func (fs *fileStat) Sys() interface{}   { return &fs.sys }
func (fs *fileStat) IsDir() bool        { return fs.mode.IsDir() }

func newDirEntry(name string, st syscall.Stat_t) *_DirEntry {
	return &_DirEntry{
		fstat: *newFileStat(name, st),
	}
}

type _DirEntry struct {
	fstat fileStat
}

func (de *_DirEntry) Name() string {
	return de.fstat.name
}

func (de *_DirEntry) IsDir() bool {
	return de.fstat.mode.IsDir()
}

func (de *_DirEntry) Type() os.FileMode {
	return de.fstat.mode
}

func (de *_DirEntry) Info() (os.FileInfo, error) {
	return &de.fstat, nil
}
