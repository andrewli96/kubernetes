package sqlfs

/*
#cgo LDFLAGS: -lsqlfs-1.0
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64 -DHAVE_LIBSQLCIPHER
#include <stdlib.h>
#include <sqlfs.h>
*/
import "C"

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type File struct {
	fs    *C.struct_sqlfs_t
	fi    *C.struct_fuse_file_info
	cPath *C.char
	flags int

	mu     sync.Mutex
	offset int64
	size   int64
	append bool
	write  bool
}

var _ io.ReadWriteCloser = (*File)(nil)
var _ io.ReaderAt = (*File)(nil)
var _ io.WriterAt = (*File)(nil)
var _ io.Seeker = (*File)(nil)

func openFile(fs *FS, path string, flags int) (*File, error) {
	var err error
	cPath := C.CString(path)
	defer func() {
		if err != nil {
			C.free(unsafe.Pointer(cPath))
		}
	}()

	fi := &C.struct_fuse_file_info{
		flags: C.int(flags),
	}
	rc := C.sqlfs_proc_open(fs.fs, cPath, fi)
	if rc != 0 {
		err = geterrno(rc)
		return nil, err
	}
	var st *syscall.Stat_t
	st, err = fs.Stat(path)
	if err != nil {
		C.sqlfs_proc_release(fs.fs, cPath, fi)
		return nil, err
	}
	return &File{
		fs:     fs.fs,
		fi:     fi,
		cPath:  cPath,
		flags:  flags,
		offset: 0,
		size:   st.Size,
		append: (flags & os.O_APPEND) != 0,
		write:  (flags & (os.O_RDWR | os.O_WRONLY)) != 0,
	}, nil
}

func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	rc := C.sqlfs_proc_release(f.fs, f.cPath, f.fi)
	if rc != 0 {
		return geterrno(rc)
	}
	C.free(unsafe.Pointer(f.cPath))
	return nil
}

func (f *File) Read(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	n, err = f.readAt(p, f.offset)
	if n > 0 {
		f.offset += int64(n)
	}
	return n, err
}

func (f *File) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	n, err = f.writeAt(p, f.offset)
	if n > 0 {
		f.offset += int64(n)
	}
	if f.offset > f.size {
		f.size = f.offset
	}
	return n, err
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.append {
		return 0, fmt.Errorf("non-seekable in APPEND mode")
	}
	switch whence {
	case os.SEEK_SET:
		f.offset = offset
	case os.SEEK_CUR:
		f.offset += offset
	case os.SEEK_END:
		f.offset = f.size + offset
	default:
		panic(fmt.Sprintf("unknown seek whence: %v", whence))
	}
	if f.write {
		if f.offset > f.size {
			f.size = f.offset
		}
	} else {
		// In read mode
		if f.offset > f.size {
			f.offset = f.size
		}
	}

	return f.offset, nil
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.readAt(p, off)
}

func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.writeAt(p, off)
}

func (f *File) readAt(p []byte, off int64) (n int, err error) {
	rc := C.sqlfs_proc_read(f.fs, f.cPath, (*C.char)(unsafe.Pointer(&p[0])), C.size_t(len(p)),
		C.off_t(off), f.fi)
	if rc == 0 {
		return 0, io.EOF
	}
	if rc < 0 {
		return 0, geterrno(rc)
	}
	return int(rc), nil
}

func (f *File) writeAt(p []byte, off int64) (n int, err error) {
	rc := C.sqlfs_proc_write(f.fs, f.cPath, (*C.char)(unsafe.Pointer(&p[0])), C.size_t(len(p)),
		C.off_t(off), f.fi)
	if rc < 0 {
		return 0, geterrno(rc)
	}
	return int(rc), nil
}
