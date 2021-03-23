package hook

import (
	"fmt"
	"os"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

// Syscall.Unlink
func hookedSyscallUnlink(path string) (err error) {
	if !_fs.Hooked(path) {
		return hookedSyscallUnlinkTramp(path)
	}
	klog.V(9).InfoS("xkube:cryptfs: Unlink", "path", path)
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return _fs.SFS.Unlink(path)
}

//go:noinline
func hookedSyscallUnlinkTramp(path string) (err error) {
	return nil
}

// Syscall.Read
func hookedSyscallRead(fd int, b []byte) (n int, err error) {
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedSyscallReadTramp(fd, b)
	}
	sfile := v.(*_SFile)
	if sfile.Dir {
		return -1, os.ErrInvalid
	}
	klog.V(9).InfoS("xkube:cryptfs: Read(hooked)", "fd", fd)
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return sfile.File.Read(b)
}

//go:noinline
func hookedSyscallReadTramp(fd int, b []byte) (n int, err error) {
	return 0, nil
}

// Syscall.Pread
func hookedSyscallPread(fd int, b []byte, offset int64) (n int, err error) {
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedSyscallPreadTramp(fd, b, offset)
	}
	sfile := v.(*_SFile)
	if sfile.Dir {
		return -1, os.ErrInvalid
	}
	klog.V(9).InfoS("xkube:cryptfs: Pread(hooked)", "fd", fd, "offset", offset)
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return sfile.File.ReadAt(b, offset)
}

//go:noinline
func hookedSyscallPreadTramp(fd int, b []byte, offset int64) (n int, err error) {
	return 0, nil
}

// Syscall.Write
func hookedSyscallWrite(fd int, b []byte) (n int, err error) {
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedSyscallWriteTramp(fd, b)
	}
	sfile := v.(*_SFile)
	if sfile.Dir {
		return -1, os.ErrInvalid
	}
	// Cannot simply call Infof/Println for infinite recursion
	if klog.V(9).Enabled() {
		hookedSyscallWriteTramp(int(os.Stdout.Fd()), []byte(
			fmt.Sprintf("xkube:cryptfs: Write(hooked) fd=%d\n", fd)))
	}
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return sfile.File.Write(b)
}

//go:noinline
func hookedSyscallWriteTramp(fd int, b []byte) (n int, err error) {
	return 0, nil
}

// Syscall.Pwrite
func hookedSyscallPwrite(fd int, b []byte, offset int64) (n int, err error) {
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedSyscallPwriteTramp(fd, b, offset)
	}
	sfile := v.(*_SFile)
	if sfile.Dir {
		return -1, os.ErrInvalid
	}
	// Cannot simply call Infof/Println for potential infinite recursion
	if klog.V(9).Enabled() {
		hookedSyscallWriteTramp(int(os.Stdout.Fd()), []byte(
			fmt.Sprintf("xkube:cryptfs: Pwrite(hooked) fd=%d offset=%d\n", fd, offset)))
	}
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return sfile.File.WriteAt(b, offset)
}

//go:noinline
func hookedSyscallPwriteTramp(fd int, b []byte, offset int64) (n int, err error) {
	return 0, nil
}

// Syscall.Seek
func hookedSyscallSeek(fd int, offset int64, whence int) (ret int64, err error) {
	v, ok := _sfiles.Get(fd)
	if !ok {
		return hookedSyscallSeekTramp(fd, offset, whence)
	}
	sfile := v.(*_SFile)
	if sfile.Dir {
		return -1, os.ErrInvalid
	}
	klog.V(9).InfoS("xkube:cryptfs: Seek(hooked)", "fd", fd, "offset", offset, "whence", whence)
	_fs.SFuckingMu.Lock()
	defer _fs.SFuckingMu.Unlock()
	return sfile.File.Seek(offset, whence)
}

//go:noinline
func hookedSyscallSeekTramp(fd int, offset int64, whence int) (ret int64, err error) {
	return 0, nil
}

func hookFileOps() error {
	var err error

	defer func() {
		if err != nil {
			unhookFileOps()
		}
	}()

	err = gohook.Hook(syscall.Read, hookedSyscallRead, hookedSyscallReadTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Pread, hookedSyscallPread, hookedSyscallPreadTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Write, hookedSyscallWrite, hookedSyscallWriteTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Pwrite, hookedSyscallPwrite, hookedSyscallPwriteTramp)
	if err != nil {
		return err
	}
	err = gohook.Hook(syscall.Seek, hookedSyscallSeek, hookedSyscallSeekTramp)
	if err != nil {
		return err
	}

	err = gohook.Hook(syscall.Unlink, hookedSyscallUnlink, hookedSyscallUnlinkTramp)
	if err != nil {
		return err
	}

	klog.V(1).Infoln("cryptfs file ops hooked")
	return nil
}

func unhookFileOps() {
	gohook.UnHook(syscall.Read)
	gohook.UnHook(syscall.Pread)
	gohook.UnHook(syscall.Write)
	gohook.UnHook(syscall.Pwrite)
	gohook.UnHook(syscall.Seek)
	gohook.UnHook(syscall.Unlink)
	klog.V(1).Infoln("cryptfs file ops unhooked")
}
