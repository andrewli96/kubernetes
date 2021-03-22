package sqlfs

import (
	"os"
	"syscall"
	"time"
)

func fillFileStatFromSys(fs *fileStat) {
	fs.size = fs.sys.Size
	fs.modTime = timespecToTime(fs.sys.Mtimespec)
	fs.mode = os.FileMode(fs.sys.Mode & 0777)
	switch fs.sys.Mode & syscall.S_IFMT {
	case syscall.S_IFBLK, syscall.S_IFWHT:
		fs.mode |= os.ModeDevice
	case syscall.S_IFCHR:
		fs.mode |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		fs.mode |= os.ModeDir
	case syscall.S_IFIFO:
		fs.mode |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		fs.mode |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		fs.mode |= os.ModeSocket
	}
	if fs.sys.Mode&syscall.S_ISGID != 0 {
		fs.mode |= os.ModeSetgid
	}
	if fs.sys.Mode&syscall.S_ISUID != 0 {
		fs.mode |= os.ModeSetuid
	}
	if fs.sys.Mode&syscall.S_ISVTX != 0 {
		fs.mode |= os.ModeSticky
	}
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// For testing.
func atime(fi os.FileInfo) time.Time {
	return timespecToTime(fi.Sys().(*syscall.Stat_t).Atimespec)
}
