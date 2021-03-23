// +build go1.16

package cryptfs

import (
	"os"
)

// ReadDir reads named directory.
func (fs *CryptFs) ReadDir(path string) ([]os.DirEntry, error) {
	if !fs.Hooked(path) {
		return os.ReadDir(path)
	}
	fs.SFuckingMu.Lock()
	defer fs.SFuckingMu.Unlock()
	return fs.SFS.ReadDir(path)
}
