// +build go1.16

package hook

import (
	"io"
	"os"

	"github.com/brahma-adshonor/gohook"
	"k8s.io/klog/v2"
)

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
		dents, err := _fs.SFS.ReadDir(file.Name())
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

func hookDirOpsEx() error {
	var err error

	defer func() {
		if err != nil {
			unhookDirOpsEx()
		}
	}()

	err = gohook.HookMethod(_dummy, "ReadDir", hookedFileReadDir, hookedFileReadDirTramp)
	if err != nil {
		return err
	}

	klog.V(1).Infoln("cryptfs directory ops_ex hooked")

	return nil
}

func unhookDirOpsEx() {
	gohook.UnHookMethod(_dummy, "ReadDir")
	klog.V(1).Infoln("cryptfs directory ops_ex unhooked")
}
