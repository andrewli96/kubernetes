package hook

import (
	"io"
	"os"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

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

// File.ReadDir
func hookedFileReadDir(file *os.File, n int) (dentries []os.DirEntry, err error) {
	// ATTENTION HERE: the lookup key(fd) should be of the exact same type with key in hashmap
	fd := int(file.Fd())
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedFileReadDirTramp(file, n)
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
	if len(sfile.Dentries) == 0 {
		dents, err := _fs.SFS.Readdir(file.Name())
		if err != nil {
			return nil, err
		}
		sfile.Dentries = dents
	}
	if n <= 0 {
		for _, x := range sfile.Dentries {
			dentries = append(dentries, x)
		}
		sfile.Dentries = make([]os.DirEntry, 0)
		sfile.DirEOF = true
		// When n <= 0, returned error should be nil instead of io.EOF
		return dentries, nil
	}

	if n > len(sfile.Dentries) {
		n = len(sfile.Dentries)
	}
	for _, x := range sfile.Dentries[:n] {
		dentries = append(dentries, x)
	}
	sfile.Dentries = sfile.Dentries[n:]
	sfile.DirEOF = len(sfile.Dentries) == 0
	err = nil
	if sfile.DirEOF {
		err = io.EOF
	}
	return dentries, err
}

//go:noinline
func hookedFileReadDirTramp(file *os.File, n int) (dentries []os.DirEntry, err error) {
	return nil, nil
}

// File.Readdir
func hookedFileReaddir(file *os.File, n int) (fis []os.FileInfo, err error) {
	dentries, err := hookedFileReadDir(file, n)
	if err != nil {
		return nil, err
	}
	for _, x := range dentries {
		fi, err := x.Info()
		if err != nil {
			return nil, err
		}
		fis = append(fis, fi)
	}
	return fis, nil
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
