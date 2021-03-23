// +build !go1.16

package hook

import (
	"io"
	"os"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

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

func hookDirOpsEx() error {
	var err error

	defer func() {
		if err != nil {
			unhookDirOpsEx()
		}
	}()

	err = gohook.HookMethod(_dummy, "Readdir", hookedFileReaddir, hookedFileReaddirTramp)
	if err != nil {
		return err
	}

	klog.V(1).Infoln("cryptfs directory ops_ex hooked")

	return nil
}

func unhookDirOpsEx() {
	gohook.UnHookMethod(_dummy, "Readdir")
	klog.V(1).Infoln("cryptfs directory ops_ex unhooked")
}
