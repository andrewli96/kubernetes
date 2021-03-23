// +build go1.16

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
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// ReadDir reads the named directory.
func (fs *FS) ReadDir(path string) ([]os.DirEntry, error) {
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
