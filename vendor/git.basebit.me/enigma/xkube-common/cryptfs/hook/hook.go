package hook

import (
	"os"

	"github.com/cornelk/hashmap"

	"git.basebit.me/enigma/sqlfs-go"
	"git.basebit.me/enigma/xkube-common/cryptfs"
)

type _SFile struct {
	Dir      bool
	File     *sqlfs.File
	Dentries []os.DirEntry
	DirEOF   bool
}

var (
	_fs     *cryptfs.CryptFs
	_sfiles hashmap.HashMap // A lockfree hash map

)

func LoadWithPlainSQLite(db string, hookedPatterns []cryptfs.MatchPattern) error {
	if _fs != nil {
		return cryptfs.ErrLoaded
	}
	fs, err := cryptfs.NewWithPlainSQLite(db, hookedPatterns)
	if err != nil {
		return err
	}
	_fs = fs

	if err := hookFileOps(); err != nil {
		Unload()
		return err
	}
	// if err := hookDirOps(); err != nil {
	// 	Unload()
	// 	return err
	// }
	if err := hookCommonOps(); err != nil {
		Unload()
		return err
	}
	return nil
}

func Load(db string, password []byte, hookedPatterns []cryptfs.MatchPattern) error {
	if _fs != nil {
		return cryptfs.ErrLoaded
	}
	fs, err := cryptfs.New(db, password, hookedPatterns)
	if err != nil {
		return err
	}
	_fs = fs

	if err := hookFileOps(); err != nil {
		Unload()
		return err
	}
	if err := hookDirOps(); err != nil {
		Unload()
		return err
	}
	if err := hookCommonOps(); err != nil {
		Unload()
		return err
	}
	return nil
}

func Unload() error {
	if _fs == nil {
		return nil
	}
	unhookCommonOps()
	unhookDirOps()
	unhookFileOps()

	if err := _fs.Close(); err != nil {
		return err
	}
	_fs = nil
	return nil
}
