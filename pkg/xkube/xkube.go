package xkube

import (
	"k8s.io/kubernetes/pkg/xkube/internal"
	"k8s.io/kubernetes/pkg/xkube/options"
)

func Init(options *options.XOptions) error {
	if !internal.XkubeEnabled {
	}
	return nil
}
