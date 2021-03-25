package hook

import (
	"os"
	"syscall"

	"git.basebit.me/enigma/xkube-common/cryptfs/utils"
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
	path = utils.Normpath(path)
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
	path = utils.Normpath(path)
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

var _dummy *os.File

func hookDirOps() error {
	var err error

	defer func() {
		if err != nil {
			unhookDirOps()
		}
	}()

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
	klog.V(1).Infoln("cryptfs directory ops unhooked")
}
