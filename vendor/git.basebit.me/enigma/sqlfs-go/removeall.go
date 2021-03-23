package sqlfs

import (
	"path/filepath"
	"syscall"
)

const (
	PathSeparator     = '/' // OS-specific path separator
	PathListSeparator = ':' // OS-specific path list separator
)

// IsPathSeparator reports whether c is a directory separator character.
func IsPathSeparator(c uint8) bool {
	return PathSeparator == c
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
func (fs *FS) RemoveAll(path string) error {
	if path == "" {
		// fail silently to retain compatibility with previous behavior
		// of RemoveAll. See issue 28830.
		return nil
	}

	// The rmdir system call does not permit removing ".",
	// so we don't permit it either.
	if endsWithDot(path) {
		return syscall.EINVAL
	}

	// Simple case: if Remove works, we're done.
	err := fs.Remove(path)
	if err == nil || err == syscall.ENOENT {
		return nil
	}

	return fs.removeAllFrom(path)
}

func (fs *FS) removeAllFrom(path string) error {
	// Simple case: if Unlink (aka remove) works, we're done.
	err := fs.Unlink(path)
	if err == nil || err == syscall.ENOENT {
		return nil
	}

	// EISDIR means that we have a directory, and we need to
	// remove its contents.
	// EPERM or EACCES means that we don't have write permission on
	// the parent directory, but this entry might still be a directory
	// whose contents need to be removed.
	// Otherwise just return the error.
	if err != syscall.EISDIR && err != syscall.EPERM && err != syscall.EACCES {
		return err
	}

	names, err := fs.Readdirnames(path)
	if err != nil {
		if err == syscall.ENOENT {
			return nil
		}
	}
	for _, x := range names {
		err := fs.removeAllFrom(filepath.Join(path, x))
		if err != nil {
			return err
		}
	}

	// Remove the directory itself.
	err = fs.Rmdir(path)
	if err == syscall.ENOENT {
		err = nil
	}
	return err
}

// endsWithDot reports whether the final component of path is ".".
func endsWithDot(path string) bool {
	if path == "." {
		return true
	}
	if len(path) >= 2 && path[len(path)-1] == '.' && IsPathSeparator(path[len(path)-2]) {
		return true
	}
	return false
}
