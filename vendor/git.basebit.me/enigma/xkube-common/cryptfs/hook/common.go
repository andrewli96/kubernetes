package hook

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"

	"git.basebit.me/enigma/sqlfs-go"

	"git.basebit.me/enigma/xkube-common/cryptfs"
)

// Syscall.Open
func hookedSyscallOpen(name string, mode int, perm uint32) (fd int, err error) {
	if mode&syscall.O_DIRECTORY != 0 {
		return -1, cryptfs.ErrUnimplemented
	}

	if !filepath.IsAbs(name) {
		wd, err := os.Getwd()
		if err != nil {
			return -1, err
		}
		name = filepath.Join(wd, name)
	}
	name = filepath.Clean(name)
	klog.V(8).InfoS("xkube:cryptfs: Open file", "name", name, "mode", mode, "perm", perm)

	if _fs.Hooked(name) {
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		f, err := _fs.SFS.Open(name, mode)
		if err != nil {
			return -1, err
		}
		fd, err := hookedSyscallOpenTramp("/dev/null", os.O_RDONLY, 0)
		if err != nil {
			return fd, err
		}
		_fs.SFiles.Set(fd, f)
		return fd, nil
	}
	return hookedSyscallOpenTramp(name, mode, perm)
}

//go:noinline
func hookedSyscallOpenTramp(name string, mode int, perm uint32) (fd int, err error) {
	return 0, nil
}

// Syscall.Close
func hookedSyscallClose(fd int) error {
	klog.V(8).InfoS("xkube:cryptfs: Close", "fd", fd)
	err := hookedSyscallCloseTramp(fd)
	if err != nil {
		return err
	}
	sf, ok := _fs.SFiles.Get(fd)
	if ok {
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		sf.(*sqlfs.File).Close()
		_fs.SFiles.Del(fd)
	}
	return nil
}

//go:noinline
func hookedSyscallCloseTramp(fd int) error {
	return nil
}

// Syscall.Stat
func hookedSyscallStat(path string, stat *syscall.Stat_t) (err error) {
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(wd, path)
	}
	path = filepath.Clean(path)
	klog.V(9).InfoS("xkube:cryptfs: Stat", "path", path)

	if _fs.Hooked(path) {
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		f, err := _fs.SFS.Stat(path)
		if err != nil {
			return err
		}
		*stat = *f
		return nil
	}
	return hookedSyscallStatTramp(path, stat)
}

//go:noinline
func hookedSyscallStatTramp(path string, stat *syscall.Stat_t) (err error) {
	return nil
}

func hookCommonOps() error {
	var err error

	defer func() {
		if err != nil {
			unhookCommonOps()
		}
	}()

	err = gohook.Hook(syscall.Stat, hookedSyscallStat, hookedSyscallStatTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Close, hookedSyscallClose, hookedSyscallCloseTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Open, hookedSyscallOpen, hookedSyscallOpenTramp)
	if err != nil {
		return err
	}
	klog.V(1).Infoln("cryptfs common ops hooked")

	return nil
}

func unhookCommonOps() {
	gohook.UnHook(syscall.Open)
	gohook.UnHook(syscall.Close)
	gohook.UnHook(syscall.Stat)
	klog.V(1).Infoln("cryptfs common ops unhooked")
}
