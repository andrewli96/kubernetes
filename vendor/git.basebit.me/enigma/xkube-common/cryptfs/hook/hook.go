package hook

import (
	"git.basebit.me/enigma/xkube-common/cryptfs"
)

var (
	_fs *cryptfs.CryptFs
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

func Unload() error {
	if _fs == nil {
		return nil
	}
	unhookCommonOps()
	// unhookDirOps()
	unhookFileOps()

	if err := _fs.Close(); err != nil {
		return err
	}
	_fs = nil
	return nil
}
