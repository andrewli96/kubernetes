package cryptfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cornelk/hashmap"
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

func New(db string, password []byte, hookedPatterns []MatchPattern) (*CryptFs, error) {
	sfs, err := sqlfs.OpenWithPassword(db, password)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlfs failed")
	}
	var pats []MatchPattern
	for _, pat := range hookedPatterns {
		if pat.Mode != MATCH_EXACT {
			return nil, ErrUnimplemented
		}
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

	SFS    *sqlfs.FS
	SFiles hashmap.HashMap
}

func (fs *CryptFs) Hooked(path string) (matched bool) {
	path = normpath(path)
	for _, x := range fs.HookedPatterns {
		switch x.Mode {
		case MATCH_EXACT:
			return x.Value == path
		case MATCH_PARENT:
			return isParentDir(x.Value, path)
		default:
			panic(fmt.Sprintf("unknown match mode=%d", x.Mode))
		}
	}
	return false
}

func (fs *CryptFs) Close() error {
	if err := fs.SFS.Close(); err != nil {
		return err
	}
	return nil
}

// ReadFile reads contents from encrypted file.
func (fs *CryptFs) ReadFile(filename string) ([]byte, error) {
	if !fs.Hooked(filename) {
		return ioutil.ReadFile(filename)
	}

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
		wd, err := os.Getwd()
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
