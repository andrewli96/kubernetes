package cryptfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"

	"git.basebit.me/enigma/sqlfs-go"
)

type MatchMode int

const (
	MATCH_EXACT MatchMode = iota
	MATCH_PARENT
)

type MatchPattern struct {
	Mode  MatchMode
	Value string
}

func NewWithPlainSQLite(db string, hookedPatterns []MatchPattern) (*CryptFs, error) {
	sfs, err := sqlfs.Open(db)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlfs failed")
	}
	defer func() {
		if err != nil {
			sfs.Close()
		}
	}()
	return newCryptFS(sfs, hookedPatterns)
}

func New(db string, password []byte, hookedPatterns []MatchPattern) (*CryptFs, error) {
	sfs, err := sqlfs.OpenWithPassword(db, password)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlfs failed")
	}
	defer func() {
		if err != nil {
			sfs.Close()
		}
	}()
	return newCryptFS(sfs, hookedPatterns)
}

func newCryptFS(sfs *sqlfs.FS, hookedPatterns []MatchPattern) (*CryptFs, error) {
	var pats []MatchPattern
	for _, pat := range hookedPatterns {
		if !filepath.IsAbs(pat.Value) {
			return nil, errors.Wrapf(filepath.ErrBadPattern, "invalid pattern=%v", pat)
		}
		pats = append(pats, MatchPattern{
			Mode:  pat.Mode,
			Value: filepath.Clean(pat.Value),
		})
	}
	fs := &CryptFs{
		HookedPatterns: pats,
		SFS:            sfs,
	}
	return fs, nil
}

type CryptFs struct {
	HookedPatterns []MatchPattern

	SFS        *sqlfs.FS
	SFuckingMu sync.Mutex // REMOVE it if sqlfs race condition resolved
}

func (fs *CryptFs) Hooked(path string) (matched bool) {
	path = normpath(path)
	found := false
	for _, x := range fs.HookedPatterns {
		switch x.Mode {
		case MATCH_EXACT:
			found = x.Value == path
		case MATCH_PARENT:
			found = isParentDir(x.Value, path)
		default:
			panic(fmt.Sprintf("unknown match mode=%d", x.Mode))
		}
		if found {
			break
		}
	}
	return found
}

func (fs *CryptFs) Close() error {
	fs.SFuckingMu.Lock()
	defer fs.SFuckingMu.Unlock()

	if err := fs.SFS.Close(); err != nil {
		return err
	}
	return nil
}

// ReadFile reads contents from encrypted file.
func (fs *CryptFs) ReadDir(path string) ([]os.DirEntry, error) {
	if !fs.Hooked(path) {
		return os.ReadDir(path)
	}
	fs.SFuckingMu.Lock()
	defer fs.SFuckingMu.Unlock()
	return fs.SFS.Readdir(path)
}

// ReadFile reads contents from encrypted file.
func (fs *CryptFs) ReadFile(filename string) ([]byte, error) {
	if !fs.Hooked(filename) {
		return ioutil.ReadFile(filename)
	}

	fs.SFuckingMu.Lock()
	defer fs.SFuckingMu.Unlock()

	sf, err := fs.SFS.Open(filename, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	var size int
	// WARNING(angus): Potential inconsistency for lack of fstat
	if st, err := fs.SFS.Stat(filename); err == nil {
		size64 := st.Size
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	} else {
		return nil, err
	}
	size++ // one byte for final read at EOF

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := sf.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}

func normpath(path string) string {
	if !filepath.IsAbs(path) {
		wd, err := syscall.Getwd()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(wd, path)
	}
	return filepath.Clean(path)
}

func isParentDir(parent, path string) bool {
	p, err := filepath.Rel(parent, path)
	return err == nil && !strings.HasPrefix(p, "../")
}
