// X Kube Constants
package constants

const (
	// DefaultPath
	XKubeconfigDefaultDir       = "/etc/kubernetes"
	XCertificatesDefaultDir     = "/etc/kubernetes/pki"
	XEtcdCertificatesDefaultDir = "/etc/kubernetes/pki/etcd"
	XKubeletDefaultDir          = "/var/lib/kubelet"
	XKubeletCertDefaultDir      = "/var/lib/kubelet/pki"
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

	// xImages
	XRegistry = "docker-reg.basebit.me:5000"

	XKubeApiserverImageName         = "kube/xkube-apiserver-amd64"
	XKubeControllerManagerImageName = "kube/xkube-controller-manager-amd64"
	XkubeSchedulerImageName         = "kube/xkube-scheduler-amd64"
	XEtcdImageName                  = "kube/xetcd"
	XEtcdImageTag                   = "v3.4.15-9-g4bfc37a0a"
)
