package hook

import (
	"io"
	"os"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

// // Syscall.Unlinkat
// func hookedSyscallUnlinkat(dirfd int, path string) error {
// 	sf, ok := _sfiles.Get(dirfd)
// 	if !ok {
// 		return hookedSyscallUnlinkatTramp(dirfd, path)
// 	}
// 	klog.V(9).InfoS("xkube:cryptfs: Unlinkat", "dirfd", dirfd, "path", path)
// 	_fs.SFuckingMu.Lock()
// 	defer _fs.SFuckingMu.Unlock()
// 	_ = sf
// 	// TODO(angus): Enhance sqlfs
// 	// return cryptfs.ErrUnimplemented
// 	return nil
// }

// //go:noinline
// func hookedSyscallUnlinkatTramp(dirfd int, path string) error {
// 	return nil
// }

// Syscall.Mkdir
func hookedSyscallMkdir(path string, mode uint32) (err error) {
	klog.V(9).InfoS("xkube:cryptfs: Mkdir", "path", path, "mode", mode)
	if _fs.Hooked(path) {
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		return _fs.SFS.Mkdir(path, mode)
	}
	return hookedSyscallMkdirTramp(path, mode)
}

//go:noinline
func hookedSyscallMkdirTramp(path string, mode uint32) (err error) {
	return nil
}

// Syscall.Rmdir
func hookedSyscallRmdir(path string) (err error) {
	klog.V(9).InfoS("xkube:cryptfs: Rmdir", "path", path)
	if _fs.Hooked(path) {
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		return _fs.SFS.Rmdir(path)
	}
	return hookedSyscallRmdirTramp(path)
}

//go:noinline
func hookedSyscallRmdirTramp(path string) (err error) {
	return nil
}

// File.Readdir
func hookedFileReaddir(file *os.File, n int) (fis []os.FileInfo, err error) {
	// ATTENTION HERE: the lookup key(fd) should be of the exact same type with key in hashmap
	fd := int(file.Fd())
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedFileReaddirTramp(file, n)
	}
	sfile := v.(*_SFile)
	if !sfile.Dir {
		return nil, os.ErrInvalid
	}
	klog.V(9).InfoS("xkube:cryptfs: File.ReadDir", "fd", file.Fd(), "n", n)
	if sfile.DirEOF {
		return nil, io.EOF
	}

	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	if len(sfile.FileInfos) == 0 {
		xs, err := _fs.SFS.Readdir(file.Name())
		if err != nil {
			return nil, err
		}
		sfile.FileInfos = xs
	}
	if n <= 0 {
		for _, x := range sfile.FileInfos {
			fis = append(fis, x)
		}
		sfile.FileInfos = make([]os.FileInfo, 0)
		sfile.DirEOF = true
		// When n <= 0, returned error should be nil instead of io.EOF
		return fis, nil
	}

	if n > len(sfile.FileInfos) {
		n = len(sfile.FileInfos)
	}
	for _, x := range sfile.FileInfos[:n] {
		fis = append(fis, x)
	}
	sfile.FileInfos = sfile.FileInfos[n:]
	sfile.DirEOF = len(sfile.FileInfos) == 0
	err = nil
	if sfile.DirEOF {
		err = io.EOF
	}
	return fis, err
}

//go:noinline
func hookedFileReaddirTramp(file *os.File, n int) (dentries []os.FileInfo, err error) {
	return nil, nil
}

var _dummy *os.File

func hookDirOps() error {
	var err error

	defer func() {
		if err != nil {
			unhookDirOps()
		}
	}()

	err = gohook.HookMethod(_dummy, "ReadDir", hookedFileReadDir, hookedFileReadDirTramp)
	if err != nil {
		return err
	}

	err = gohook.HookMethod(_dummy, "Readdir", hookedFileReaddir, hookedFileReaddirTramp)
	if err != nil {
		return err
	}

	err = gohook.Hook(syscall.Mkdir, hookedSyscallMkdir, hookedSyscallMkdirTramp)
	if err != nil {
		return err
	}

	err = gohook.Hook(syscall.Rmdir, hookedSyscallRmdir, hookedSyscallRmdirTramp)
	if err != nil {
		return err
	}
	klog.V(1).Infoln("cryptfs directory ops hooked")

	return nil
}

func unhookDirOps() {
	gohook.UnHook(syscall.Mkdir)
	gohook.UnHook(syscall.Rmdir)
	gohook.UnHookMethod(_dummy, "ReadDir")
	gohook.UnHookMethod(_dummy, "Readdir")
	klog.V(1).Infoln("cryptfs directory ops unhooked")
}
