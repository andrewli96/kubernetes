package internal

import (
	"k8s.io/klog/v2"
)

var (
	xkubeEnable  = "false" // A compile-time variable which must be of string type
	XkubeEnabled bool
)

func init() {
	XkubeEnabled = xkubeEnable == "true"
	if XkubeEnabled {
		klog.V(1).Infoln("xkube enabled")
	} else {
		klog.V(1).Infoln("xkube disabled")
	}
}
