package hook

import (
	"fmt"
	"os"
	"syscall"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"

	"git.basebit.me/enigma/sqlfs-go"
)

// Syscall.Unlink
func hookedSyscallUnlink(path string) (err error) {
	klog.V(9).InfoS("xkube:cryptfs: Unlink", "path", path)
	if _fs.Hooked(path) {
		return _fs.SFS.Unlink(path)
	}
	return hookedSyscallUnlinkTramp(path)
}

//go:noinline
func hookedSyscallUnlinkTramp(path string) (err error) {
	return nil
}

// // Syscall.Unlinkat
// func hookedSyscallUnlinkat(dirfd int, path string) error {
// 	klog.V(9).InfoS("xkube:cryptfs: Unlinkat", "dirfd", dirfd, "path", path)
// 	sf, ok := _fs.SFiles.Get(dirfd)
// 	if !ok {
// 		return hookedSyscallUnlinkatTramp(dirfd, path)
// 	}
// 	_ = sf
// 	// TODO(angus): Enhance sqlfs
// 	return cryptfs.ErrUnimplemented
// }

// //go:noinline
// func hookedSyscallUnlinkatTramp(dirfd int, path string) error {
// 	return nil
// }

// Syscall.Read
func hookedSyscallRead(fd int, b []byte) (n int, err error) {
	klog.V(9).InfoS("xkube:cryptfs: Read", "fd", fd)
	sf, ok := _fs.SFiles.Get(fd)
	if !ok {
		return hookedSyscallReadTramp(fd, b)
	}
	return sf.(*sqlfs.File).Read(b)
}

//go:noinline
func hookedSyscallReadTramp(fd int, b []byte) (n int, err error) {
	return 0, nil
}

// Syscall.Pread
func hookedSyscallPread(fd int, b []byte, offset int64) (n int, err error) {
	klog.V(9).InfoS("xkube:cryptfs: Pread", "fd", fd, "offset", offset)
	sf, ok := _fs.SFiles.Get(fd)
	if !ok {
		return hookedSyscallPreadTramp(fd, b, offset)
	}
	return sf.(*sqlfs.File).ReadAt(b, offset)
}

//go:noinline
func hookedSyscallPreadTramp(fd int, b []byte, offset int64) (n int, err error) {
	return 0, nil
}

// Syscall.Write
func hookedSyscallWrite(fd int, b []byte) (n int, err error) {
	// Cannot simply call Infof/Println for infinite recursion
	if klog.V(9).Enabled() {
		hookedSyscallWriteTramp(int(os.Stdout.Fd()), []byte(
			fmt.Sprintf("xkube:cryptfs: Write fd=%d\n", fd)))
	}
	sf, ok := _fs.SFiles.Get(fd)
	if !ok {
		return hookedSyscallWriteTramp(fd, b)
	}
	return sf.(*sqlfs.File).Write(b)
}

//go:noinline
func hookedSyscallWriteTramp(fd int, b []byte) (n int, err error) {
	return 0, nil
}

// Syscall.Pwrite
func hookedSyscallPwrite(fd int, b []byte, offset int64) (n int, err error) {
	// Cannot simply call Infof/Println for potential infinite recursion
	if klog.V(9).Enabled() {
		hookedSyscallWriteTramp(int(os.Stdout.Fd()), []byte(
			fmt.Sprintf("xkube:cryptfs: Pwrite fd=%d offset=%d\n", fd, offset)))
	}
	sf, ok := _fs.SFiles.Get(fd)
	if !ok {
		return hookedSyscallPwriteTramp(fd, b, offset)
	}
	return sf.(*sqlfs.File).WriteAt(b, offset)
}

//go:noinline
func hookedSyscallPwriteTramp(fd int, b []byte, offset int64) (n int, err error) {
	return 0, nil
}

// Syscall.Seek
func hookedSyscallSeek(fd int, offset int64, whence int) (ret int64, err error) {
	klog.V(9).InfoS("xkube:cryptfs: Seek", "fd", fd, "offset", offset, "whence", whence)
	sf, ok := _fs.SFiles.Get(fd)
	if !ok {
		return hookedSyscallSeekTramp(fd, offset, whence)
	}
	return sf.(*sqlfs.File).Seek(offset, whence)
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
