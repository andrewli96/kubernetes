package xcontainerd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

var (
	xcontainerdConfig = fmt.Sprintf(`version = 2
root = "/var/lib/containerd"
state = "/run/containerd"
plugin_dir = ""
disabled_plugins = []
required_plugins = []
oom_score = 0

[grpc]
	address = "/run/containerd/containerd.sock"
	tcp_address = ""
	tcp_tls_root_ca = "%s"
	tcp_tls_cert = "%s"
	tcp_tls_key = "%s"
	uid = 0
	gid = 0
	max_recv_message_size = 16777216
	max_send_message_size = 16777216

[ttrpc]
	address = ""
	uid = 0
	gid = 0

[debug]
	address = ""
	uid = 0
	gid = 0
	level = ""

[metrics]
	address = ""
	grpc_histogram = false

[cgroup]
	path = ""

[timeouts]
	"io.containerd.timeout.shim.cleanup" = "5s"
	"io.containerd.timeout.shim.load" = "5s"
	"io.containerd.timeout.shim.shutdown" = "3s"
	"io.containerd.timeout.task.state" = "2s"

[plugins]
	[plugins."io.containerd.gc.v1.scheduler"]
	pause_threshold = 0.02
	deletion_threshold = 0
	mutation_threshold = 100
	schedule_delay = "0s"
	startup_delay = "100ms"
	[plugins."io.containerd.grpc.v1.cri"]
	disable_tcp_service = true
	stream_server_address = "127.0.0.1"
	stream_server_port = "0"
	stream_idle_timeout = "4h0m0s"
	enable_selinux = false
	selinux_category_range = 1024
	sandbox_image = "k8s.gcr.io/pause:3.2"
	stats_collect_period = 10
	systemd_cgroup = false
	enable_tls_streaming = false
	max_container_log_line_size = 16384
	disable_cgroup = false
	disable_apparmor = false
	restrict_oom_score_adj = false
	max_concurrent_downloads = 3
	disable_proc_mount = false
	unset_seccomp_profile = ""
	tolerate_missing_hugetlb_controller = true
	disable_hugetlb_controller = true
	ignore_image_defined_volumes = false
	[plugins."io.containerd.grpc.v1.cri".containerd]
		snapshotter = "overlayfs"
		default_runtime_name = "runc"
		no_pivot = false
		disable_snapshot_annotations = true
		discard_unpacked_layers = false
		[plugins."io.containerd.grpc.v1.cri".containerd.default_runtime]
		runtime_type = ""
		runtime_engine = ""
		runtime_root = ""
		privileged_without_host_devices = false
		base_runtime_spec = ""
		[plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime]
		runtime_type = ""
		runtime_engine = ""
		runtime_root = ""
		privileged_without_host_devices = false
		base_runtime_spec = ""
		[plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
		[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
			runtime_type = "io.containerd.runc.v2"
			runtime_engine = ""
			runtime_root = ""
			privileged_without_host_devices = false
			base_runtime_spec = ""
			[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
	[plugins."io.containerd.grpc.v1.cri".cni]
		bin_dir = "/opt/cni/bin"
		conf_dir = "/etc/cni/net.d"
		max_conf_num = 1
		conf_template = ""
	[plugins."io.containerd.grpc.v1.cri".registry]
		[plugins."io.containerd.grpc.v1.cri".registry.mirrors]
		[plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
			endpoint = ["https://registry-1.docker.io"]
	[plugins."io.containerd.grpc.v1.cri".image_decryption]
		key_model = ""
	[plugins."io.containerd.grpc.v1.cri".x509_key_pair_streaming]
		tls_cert_file = ""
		tls_key_file = ""
	[plugins."io.containerd.internal.v1.opt"]
	path = "/opt/containerd"
	[plugins."io.containerd.internal.v1.restart"]
	interval = "10s"
	[plugins."io.containerd.metadata.v1.bolt"]
	content_sharing_policy = "shared"
	[plugins."io.containerd.monitor.v1.cgroups"]
	no_prometheus = false
	[plugins."io.containerd.runtime.v1.linux"]
	shim = "containerd-shim"
	runtime = "runc"
	runtime_root = ""
	no_shim = false
	shim_debug = false
	[plugins."io.containerd.runtime.v2.task"]
	platforms = ["linux/amd64"]
	[plugins."io.containerd.service.v1.diff-service"]
	default = ["walking"]
	[plugins."io.containerd.snapshotter.v1.devmapper"]
	root_path = ""
	pool_name = ""
	base_image_size = ""
	async_remove = false`,
		filepath.Join(kubeadmconstants.XCertificatesDefaultDir, kubeadmconstants.CACertName),
		filepath.Join(kubeadmconstants.XCertificatesDefaultDir, kubeadmconstants.XContainerdCertName),
		filepath.Join(kubeadmconstants.XCertificatesDefaultDir, kubeadmconstants.XContainerdKeyName))
)

// WriteConfigToDisk writes the kubelet config object down to a file
// Used at "kubeadm init" and "kubeadm upgrade" time
func WriteConfigToDisk(cfg *kubeadmapi.ClusterConfiguration, xcontainerDir string) error {

	xcontainerdConfigBytes, err := getXContainerdConfigBytes()
	if err != nil {
		return err
	}

	return writeConfigBytesToDisk(xcontainerdConfigBytes, xcontainerDir)
}

// writeConfigBytesToDisk writes a byte slice down to disk at the specific location of the kubelet config file
func writeConfigBytesToDisk(b []byte, xcontainerDir string) error {
	configFile := filepath.Join(xcontainerDir, kubeadmconstants.XContainerdConfigurationFileName)
	fmt.Printf("[xcontainerd-start] Writing kubelet configuration to file %q\n", configFile)

	// creates target folder if not already exists
	if err := os.MkdirAll(xcontainerDir, 0700); err != nil {
		return errors.Wrapf(err, "failed to create directory %q", xcontainerDir)
	}

	if err := ioutil.WriteFile(configFile, b, 0644); err != nil {
		return errors.Wrapf(err, "failed to write kubelet configuration to the file %q", configFile)
	}
	return nil
}

func getXContainerdConfigBytes() ([]byte, error) {

	if kubeadmconstants.XCertificatesDefaultDir == "" || kubeadmconstants.CACertName == "" ||
		kubeadmconstants.XContainerdCertName == "" || kubeadmconstants.XContainerdKeyName == "" {
		return nil, errors.New("Failed to load xcontainerd certs path.")
	}

	return []byte(xcontainerdConfig), nil
}
