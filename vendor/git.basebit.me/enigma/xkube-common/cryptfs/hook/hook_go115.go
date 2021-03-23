// +build !go1.16

package hook

type _SFile struct {
	Dir      bool
	File     *sqlfs.File
	Dentries []os.DirEntry
	DirEOF   bool
}
