package hook

import (
	"os"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"

	"git.basebit.me/enigma/xkube-common/cryptfs"
)

// Syscall.Mkdir
func hookedSyscallMkdir(path string, mode uint32) (err error) {
	klog.V(9).InfoS("xkube:cryptfs: Mkdir", "path", path, "mode", mode)
	if _fs.Hooked(path) {
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
		return _fs.SFS.Rmdir(path)
	}
	return hookedSyscallRmdirTramp(path)
}

//go:noinline
func hookedSyscallRmdirTramp(path string) (err error) {
	return nil
}

// // File.ReadDir
// func hookedFileReadDir(file *os.File, n int) (dentries []os.DirEntry, err error) {
// 	klog.V(9).InfoS("xkube:cryptfs: File.ReadDir", "fd", file.Fd(), "n", n)
// 	_, ok := _fs.SFiles.Get(file.Fd())
// 	if !ok {
// 		return hookedFileReadDirTramp(file, n)
// 	}
// 	names, err := _fs.SFS.Readdirnames(file.Name())
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, name := range names {
// 		st, err := _fs.SFS.Stat(filepath.Join(file.Name(), name))
// 		if err != nil {
// 			continue
// 		}
// 		dentries = append(dentries, newDirEntry(name, *st))
// 	}
// 	return dentries, nil
// }

// //go:noinline
// func hookedFileReadDirTramp(file *os.File, n int) (dentries []os.DirEntry, err error) {
// 	return nil, nil
// }

// File.Readdir
func hookedFileReaddir(file *os.File, n int) (dentries []os.FileInfo, err error) {
	klog.V(9).InfoS("xkube:cryptfs: File.Readdir", "fd", file.Fd(), "n", n)
	_, ok := _fs.SFiles.Get(file.Fd())
	if !ok {
		return hookedFileReaddirTramp(file, n)
	}
	// TODO(angus): Enhance sqlfs-go to support stateful readdir
	if n <= 0 {
		// TODO(angus): Read all entries
	}
	return nil, cryptfs.ErrUnimplemented
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

	// err = gohook.HookMethod(_dummy, "ReadDir", hookedFileReadDir, hookedFileReadDirTramp)
	// if err != nil {
	// 	return err
	// }

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
	// gohook.UnHookMethod(_dummy, "ReadDir")
	gohook.UnHookMethod(_dummy, "Readdir")
	klog.V(1).Infoln("cryptfs directory ops unhooked")
}
