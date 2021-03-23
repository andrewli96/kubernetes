package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

// Syscall.Open
func hookedSyscallOpen(name string, mode int, perm uint32) (fd int, err error) {
	if !filepath.IsAbs(name) {
		wd, err := os.Getwd()
		if err != nil {
			return -1, err
		}
		name = filepath.Join(wd, name)
	}
	name = filepath.Clean(name)
	if !_fs.Hooked(name) {
		return hookedSyscallOpenTramp(name, mode, perm)
	}

	klog.V(8).InfoS("xkube:cryptfs: Open file", "name", name, "mode", mode, "perm", perm)
	dir := false
	if st, err := os.Stat(name); err != nil {
		if !os.IsNotExist(err) {
			return -1, err
		}
		if mode&os.O_CREATE == 0 {
			return -1, err
		}
	} else {
		dir = st.IsDir()
	}
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	sfile := &_SFile{
		Dir:    dir,
		DirEOF: false,
	}
	if !dir {
		f, err := _fs.SFS.Open(name, mode)
		if err != nil {
			return -1, err
		}
		sfile.File = f
	}
	fd, err = hookedSyscallOpenTramp("/dev/null", os.O_RDONLY, 0)
	if err != nil {
		return fd, err
	}
	_sfiles.Set(fd, sfile)
	return fd, nil
}

//go:noinline
func hookedSyscallOpenTramp(name string, mode int, perm uint32) (fd int, err error) {
	return 0, nil
}

// Syscall.Close
func hookedSyscallClose(fd int) error {
	err := hookedSyscallCloseTramp(fd)
	if err != nil {
		return err
	}
	v, ok := _sfiles.Get(fd)
	if ok {
		klog.V(8).InfoS("xkube:cryptfs: Close", "fd", fd)
		sfile := v.(*_SFile)
		_fs.SFuckingMu.Lock()
		defer _fs.SFuckingMu.Unlock()
		if !sfile.Dir {
			sfile.File.Close()
		}
		_sfiles.Del(fd)
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
		wd, err := syscall.Getwd()
		if err != nil {
			return err
		}
		if !filepath.IsAbs(wd) {
			panic(fmt.Sprintf("Getwd: Expected absolute path but get %s", wd))
		}
		path = filepath.Join(wd, path)
	}
	path = filepath.Clean(path)
	if !_fs.Hooked(path) {
		return hookedSyscallStatTramp(path, stat)
	}
	klog.V(9).InfoS("xkube:cryptfs: Stat", "path", path)
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	f, err := _fs.SFS.Stat(path)
	if err != nil {
		return err
	}
	*stat = *f
	return nil
}

//go:noinline
func hookedSyscallStatTramp(path string, stat *syscall.Stat_t) (err error) {
	return nil
}

// os.RemoveAll
func hookedOSRemoveAll(path string) (err error) {
	if !filepath.IsAbs(path) {
		wd, err := syscall.Getwd()
		if err != nil {
			return err
		}
		if !filepath.IsAbs(wd) {
			panic(fmt.Sprintf("Getwd: Expected absolute path but get %s", wd))
		}
		path = filepath.Join(wd, path)
	}
	path = filepath.Clean(path)
	if !_fs.Hooked(path) {
		return hookedOSRemoveAllTramp(path)
	}
	klog.V(9).InfoS("xkube:cryptfs: RemoveAll", "path", path)

	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return _fs.SFS.RemoveAll(path)
}

//go:noinline
func hookedOSRemoveAllTramp(path string) (err error) {
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
	err = gohook.Hook(os.RemoveAll, hookedOSRemoveAll, hookedOSRemoveAllTramp)
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
	gohook.UnHook(os.RemoveAll)
	klog.V(1).Infoln("cryptfs common ops unhooked")
}
