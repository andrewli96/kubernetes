// +build go1.16

package hook

import (
	"os"

	"git.basebit.me/enigma/sqlfs-go"
)

type _SFile struct {
	Dir      bool
	File     *sqlfs.File
	Dentries []os.DirEntry // For Go 1.16+
	DirEOF   bool
}
