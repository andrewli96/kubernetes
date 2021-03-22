package sqlfs

/*
#cgo LDFLAGS: -lsqlfs-1.0
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64 -DHAVE_LIBSQLCIPHER
#include <stdlib.h>
#include <string.h>
#include <sqlfs.h>

int readdirnames_filler(void* buf, char* name, struct stat *stbuf, off_t off);
*/
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/mattn/go-pointer"
)

type FS struct {
	fs *C.sqlfs_t
}

func Open(db string) (*FS, error) {
	cDB := C.CString(db)
	defer C.free(unsafe.Pointer(cDB))
	var fs *C.sqlfs_t

	rc := C.sqlfs_open(cDB, &fs)
	if rc == 0 {
		return nil, fmt.Errorf("open sqlfs error")
	}

	return &FS{
		fs: fs,
	}, nil
}

func OpenWithPassword(db string, password []byte) (*FS, error) {
	cDB := C.CString(db)
	defer C.free(unsafe.Pointer(cDB))
	cPassword := C.CBytes(password)
	defer C.free(cPassword)
	defer C.memset(cPassword, 0, C.size_t(len(password)))
	var fs *C.sqlfs_t

	rc := C.sqlfs_open_password(cDB, (*C.char)(cPassword), &fs)
	if rc == 0 {
		return nil, fmt.Errorf("open sqlfs error")
	}

	return &FS{
		fs: fs,
	}, nil
}

// Close closes the file system
func (fs *FS) Close() error {
	rc := C.sqlfs_close(fs.fs)
	if rc == 0 {
		return fmt.Errorf("close sqlfs error")
	}
	return nil
}

// Mkdir creates a directory.
func (fs *FS) Mkdir(path string, mode uint32) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_mkdir(fs.fs, cPath, C.mode_t(mode))
	return geterrno(rc)
}

// Rmdir removes a directory.
func (fs *FS) Rmdir(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_rmdir(fs.fs, cPath)
	return geterrno(rc)
}

//export readdirnames_filler
func readdirnames_filler(buf unsafe.Pointer, name *C.char, stbuf *C.struct_stat, off C.off_t) C.int {
	names := pointer.Restore(buf).(*[]string)
	str := C.GoString(name)
	if str == "." || str == ".." {
		return 0
	}
	*names = append(*names, str)
	return 0
}

// ReadDir reads the named directory.
func (fs *FS) Readdir(path string) ([]os.DirEntry, error) {
	names, err := fs.Readdirnames(path)
	if err != nil {
		return nil, err
	}

	var dentries []os.DirEntry
	for _, name := range names {
		var st syscall.Stat_t
		cPath := C.CString(filepath.Join(path, name))
		rc := C.sqlfs_proc_getattr(fs.fs, cPath, (*C.struct_stat)(unsafe.Pointer(&st)))
		C.free(unsafe.Pointer(cPath))
		if rc != 0 {
			err := geterrno(rc)
			if err == syscall.ENOENT {
				continue
			}
			return nil, err
		}
		dentries = append(dentries, newDirEntry(name, st))
	}
	return dentries, nil
}

// ReadDir reads the named directory.
func (fs *FS) Readdirnames(path string) ([]string, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var names []string
	rc := C.sqlfs_proc_readdir(fs.fs, cPath, pointer.Save(&names),
		C.fuse_fill_dir_t(C.readdirnames_filler), 0, nil)
	if rc != 0 {
		return nil, geterrno(rc)
	}
	return names, nil
}

// Link creates a link.
func (fs *FS) Link(existing, newname string) error {
	cExisting := C.CString(existing)
	defer C.free(unsafe.Pointer(cExisting))
	cNewname := C.CString(newname)
	defer C.free(unsafe.Pointer(cNewname))

	rc := C.sqlfs_proc_link(fs.fs, cExisting, cNewname)
	return geterrno(rc)
}

// Readlink reads a symbolic link.
func (fs *FS) Readlink(path string) (string, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var buf [4096]C.char

	rc := C.sqlfs_proc_readlink(fs.fs, cPath, (*C.char)(&buf[0]), 4096)
	if rc >= 0 {
		return C.GoStringN(&buf[0], rc), nil
	}
	return "", geterrno(rc)
}

// Symlink creates a symbolic link.
func (fs *FS) Symlink(existing, newname string) error {
	cExisting := C.CString(existing)
	defer C.free(unsafe.Pointer(cExisting))
	cNewname := C.CString(newname)
	defer C.free(unsafe.Pointer(cNewname))

	rc := C.sqlfs_proc_symlink(fs.fs, cExisting, cNewname)
	return geterrno(rc)
}

// Unlink removes a file, link, or symbolic link.
func (fs *FS) Unlink(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_unlink(fs.fs, cPath)
	return geterrno(rc)
}

// Rename renames a file or directory.
func (fs *FS) Rename(from, to string) error {
	cFrom := C.CString(from)
	defer C.free(unsafe.Pointer(cFrom))
	cTo := C.CString(to)
	defer C.free(unsafe.Pointer(cTo))

	rc := C.sqlfs_proc_rename(fs.fs, cFrom, cTo)
	return geterrno(rc)
}

// Stat gets a file's statistics and attributes.
func (fs *FS) Stat(path string) (*syscall.Stat_t, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var st syscall.Stat_t
	rc := C.sqlfs_proc_getattr(fs.fs, cPath, (*C.struct_stat)(unsafe.Pointer(&st)))
	if rc != 0 {
		return nil, geterrno(rc)
	}
	return &st, nil
}

// Access checks user's permissions for a file.
func (fs *FS) Access(path string, mask int) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_access(fs.fs, cPath, C.int(mask))
	return geterrno(rc)
}

// Chmod changes the mode bits (permissions) of a file/directory.
func (fs *FS) Chmod(path string, mode uint32) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_chmod(fs.fs, cPath, C.mode_t(mode))
	return geterrno(rc)
}

// Chown changes the ownership of a file/directory.
func (fs *FS) Chown(path string, user uint, group uint) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.sqlfs_proc_chown(fs.fs, cPath, C.uid_t(user), C.gid_t(group))
	return geterrno(ret)
}

// Truncate truncates the file to the given size.
func (fs *FS) Truncate(path string, size int64) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.sqlfs_proc_truncate(fs.fs, cPath, C.int64_t(size))
	return geterrno(ret)
}

// Remove removes the named file or (empty) directory.
func (fs *FS) Remove(path string) error {
	// System call interface forces us to know
	// whether name is a file or directory.
	// Try both: it is cheaper on average than
	// doing a Stat plus the right one.
	e := fs.Unlink(path)
	if e == nil {
		return nil
	}
	e1 := fs.Rmdir(path)
	if e1 == nil {
		return nil
	}

	// Both failed: figure out which error to return.
	// OS X and Linux differ on whether unlink(dir)
	// returns EISDIR, so can't use that. However,
	// both agree that rmdir(file) returns ENOTDIR,
	// so we can use that to decide which error is real.
	// Rmdir might also return ENOTDIR if given a bad
	// file path, like /etc/passwd/foo, but in that case,
	// both errors will be ENOTDIR, so it's okay to
	// use the error from unlink.
	if e1 != syscall.ENOTDIR {
		e = e1
	}
	return e
}

// Utime changes file/directory last access and modification times.
func (fs *FS) Utime(path string, actime, modtime int64) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cBuf := C.struct_utimbuf{
		actime:  C.time_t(actime),
		modtime: C.time_t(modtime),
	}
	rc := C.sqlfs_proc_utime(fs.fs, cPath, &cBuf)
	return geterrno(rc)
}

// Mknod makes a block or character special file.
func (fs *FS) Mknod(path string, mode uint32, dev int) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	rc := C.sqlfs_proc_mknod(fs.fs, cPath, C.mode_t(mode), C.dev_t(dev))
	return geterrno(rc)
}

// Statfs gets filesystem statistics.
func (fs *FS) Statfs(path string) (*syscall.Statfs_t, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var stbuf syscall.Statfs_t
	rc := C.sqlfs_proc_statfs(fs.fs, cPath, (*C.struct_statvfs)(unsafe.Pointer(&stbuf)))
	if rc != 0 {
		return nil, geterrno(rc)
	}
	return &stbuf, nil
}

func (fs *FS) Open(path string, flags int) (*File, error) {
	return openFile(fs, path, flags)
}

func geterrno(errno C.int) error {
	if errno == 0 {
		return nil
	}
	if errno < 0 {
		errno = -errno
	}
	return syscall.Errno(errno)
}
