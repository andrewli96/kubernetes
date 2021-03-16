package internal

import (
	"k8s.io/klog/v2"
)

var (
	xkubeEnabled = "0" // A compile-time variable which must be of string type
	XkubeEnabled bool
)

func init() {
	XkubeEnabled = xkubeEnabled == "1"
	if XkubeEnabled {
		klog.V(1).Infoln("xkube enabled")
	} else {
		klog.V(1).Infoln("xkube disabled")
	}
}
