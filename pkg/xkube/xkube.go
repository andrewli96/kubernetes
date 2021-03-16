package xkube

import (
	"fmt"

	"k8s.io/klog/v2"

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
		return nil
	}
	if options.XConfigFile == "" {
		return fmt.Errorf("X config file not set")
	}
	err := hook.Load(options.XConfigFile, getConfigFileKey(), patterns)
	if err != nil {
		return err
	}
	klog.Infoln("xkube loaded")
	return nil
}

func Close() error {
	if !internal.XkubeEnabled {
		return nil
	}

	if err := hook.Unload(); err != nil {
		return err
	}
	klog.Infoln("xube unloaded")
	return nil
}
