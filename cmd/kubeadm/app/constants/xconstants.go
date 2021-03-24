// X Kube Constants
package constants

const (
	// DefaultPath
	XKubeconfigDefaultDir       = "/etc/kubernetes"
	XCertificatesDefaultDir     = "/etc/kubernetes/pki"
	XEtcdCertificatesDefaultDir = "/etc/kubernetes/pki/etcd"
	XKubeletDefaultDir          = "/var/lib/kubelet"
	XManifestDefaultDir         = "/etc/kubernetes/manifests"
	XContainerdDefaultDir       = "/etc/containerd"

	XConfigFileArgumentKey = "x-config"

	// xCommand
	XKubeletCommand               = "xkubelet"
	XKubeApiserverCommand         = "xkube-apiserver"
	XKubeControllerManagerCommand = "xkube-controller-manager"
	XkubeSchedulerCommand         = "xkube-scheduler"
	XEtcdCommand                  = "xetcd"
	XContainerdCommand            = "xcontainerd"

	// xVolume
	XConfigFileMountVolumeName = "xconfig"

	// certs
	XContainerdCertAndKeyBaseName = "containerd"
	XContainerdCertCommonName     = "containerd"
	XContainerdCertName           = "containerd.crt"
	XContainerdKeyName            = "containerd.key"

	XKubeletCRIClientCertAndKeyBaseName = "kubelet-cri-client"
	XKubeletCRIClientCertCommonName     = "kubelet-cri-client"
	XKubeletCRIClientCertName           = "kubelet-cri-client.crt"

	// xContainerd
	XContainerdConfigurationFileName = "config.toml"
)
