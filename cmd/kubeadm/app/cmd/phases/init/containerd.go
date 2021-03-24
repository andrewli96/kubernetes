package phases

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	xcontainerdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/xcontainerd"
)

var (
	xContainerdStartPhaseExample = cmdutil.Examples(`
		# Write xcontainerd configurations, then (re)start xcontainerd on the node.
		kubeadm init phase xcontainerd-start
		`)
)

// NewXContainerdPhase creates a kubeadm workflow phase that start xcontainerd on a node.
func NewXContainerdPhase() workflow.Phase {
	return workflow.Phase{
		Name:    "xcontainerd-start",
		Short:   "Write xcontainerd configurations and (re)start the xcontainerd",
		Long:    "Write xcontainerd configurations, then (re)start xcontainerd.",
		Example: xContainerdStartPhaseExample,
		Run:     runXContainerdStart,
	}
}

// runXContainerdStart executes xcontainerd start logic.
func runXContainerdStart(c workflow.RunData) error {
	data, ok := c.(InitData)
	if !ok {
		return errors.New("xcontainerd-start phase invoked with an invalid data struct")
	}

	// First off, configure the xcontainerd. In this short timeframe, kubeadm is trying to stop/restart the xcontainerd
	// Try to stop the xcontainerd service so no race conditions occur when configuring it
	if !data.DryRun() {
		klog.V(1).Infoln("Stopping the xcontainerd")
		xcontainerdphase.TryStopXContainerd()
	}

	// Write the xcontainerd configuration file to sqlfs.
	if err := xcontainerdphase.WriteConfigToDisk(&data.Cfg().ClusterConfiguration, kubeadmconstants.XContainerdDefaultDir); err != nil {
		return errors.Wrap(err, "error writing xcontainer configuration to sqlfs")
	}

	// Try to start the kubelet service in case it's inactive
	if !data.DryRun() {
		fmt.Println("[xcontainerd-start] Starting the xcontainerd")
		xcontainerdphase.TryStartXContainerd()
	}

	return nil
}
