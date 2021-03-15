package cryptfs

import (
	"github.com/pkg/errors"
)

var (
	ErrLoaded        = errors.New("cryptfs loaded")
	ErrUnimplemented = errors.New("unimplemented")
)
