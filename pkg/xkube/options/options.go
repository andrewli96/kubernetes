package options

import (
	"github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/xkube/internal"
)

type XOptions struct {
	XConfigFile string
}

// NewXOptions creates a new instance of XOptions
func NewXOptions() *XOptions {
	return &XOptions{}
}

// AddFlags adds flags related to XDP specified FlagSet
func (x *XOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&x.XConfigFile, "x-config", x.XConfigFile,
		"The path to the X configuration file. Empty string for no configuration file.")
}

// Validate checks invalid config
func (x *XOptions) Validate() []error {
	if x == nil {
		return nil
	}
	allErrors := []error{}
	if internal.XkubeEnabled {
	}
	return allErrors
}
