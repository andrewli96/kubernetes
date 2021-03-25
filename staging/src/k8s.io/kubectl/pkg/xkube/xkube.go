package xkube

import (
	"fmt"

	"git.basebit.me/enigma/xkube-common/cryptfs"
	"git.basebit.me/enigma/xkube-common/cryptfs/hook"
	"k8s.io/klog/v2"
)

func Setup(xConfigFile string, password string, patterns []cryptfs.MatchPattern) error {
	if xConfigFile == "" {
		return fmt.Errorf("X config file not set")
	}
	err := hook.Load(xConfigFile, []byte(password), patterns)
	if err != nil {
		return err
	}
	klog.V(2).Infoln("xkube loaded")
	return nil
}

func Close() error {
	if err := hook.Unload(); err != nil {
		return err
	}
	klog.V(2).Infoln("xkube unloaded")
	return nil
}
