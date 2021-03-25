package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
)

func Normpath(path string) string {
	if strings.TrimSpace(path) == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		wd, err := syscall.Getwd()
		if err != nil {
			panic(err)
		}
		if !filepath.IsAbs(wd) {
			panic(fmt.Sprintf("Getwd: Expected absolute path but get %s", wd))
		}
		path = filepath.Join(wd, path)
	}
	return filepath.Clean(path)
}
