// X Kube Constants
package constants

const (
	// DefaultPath
	XKubeconfigDefaultDir       = "/etc/kubernetes"
	XCertificatesDefaultDir     = "/etc/kubernetes/pki"
	XEtcdCertificatesDefaultDir = "/etc/kubernetes/pki/etcd"
	XKubeletDefaultDir          = "/var/lib/kubelet"
	XManifestDefaultDir         = "/etc/kubernetes/manifests"

	XConfigFileArgumentKey = "x-config"

	// xCommand
	XKubeletCommand               = "xkubelet"
	XKubeApiserverCommand         = "xkube-apiserver"
	XKubeControllerManagerCommand = "xkube-controller-manager"
	XkubeSchedulerCommand         = "xkube-scheduler"
	XEtcdCommand                  = "xetcd"

	// xVolume
	XConfigFileMountVolumeName = "xconfig"
)
