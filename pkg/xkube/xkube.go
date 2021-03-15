package xkube

import (
	"fmt"

	"git.basebit.me/enigma/xkube-common/cryptfs"
	"git.basebit.me/enigma/xkube-common/cryptfs/hook"

	"k8s.io/kubernetes/pkg/xkube/internal"
	"k8s.io/kubernetes/pkg/xkube/options"
)

func getConfigFileKey() []byte {
	// TODO(angus): Replace the static plain password with dynamic obscure byte
	// return []byte("PAsWORD_HERE_123")
	return []byte("foo")
}

func Setup(options *options.XOptions, patterns []cryptfs.MatchPattern) error {
	if !internal.XkubeEnabled {
		if options.XConfigFile == "" {
			return fmt.Errorf("X config file not set")
		}
		err := hook.Load(options.XConfigFile, getConfigFileKey(), patterns)
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() error {
	return hook.Unload()
}
