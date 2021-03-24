package xcontainerd

import (
	"fmt"

	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/initsystem"
)

var (
	//kubeletCommand = "kubelet"
	//XKubelet
	xcontainerdCommand = kubeadmconstants.XContainerdCommand
)

// TryStopKubelet attempts to bring down the kubelet service momentarily
func TryStopXContainerd() {
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[xcontainerd-start] no supported init system detected, won't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists(xcontainerdCommand) {
		fmt.Println("[xcontainerd-start] couldn't detect a xcontainerd service, can't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
	}

	// This runs "systemctl daemon-reload && systemctl stop xcontainerd"
	if err := initSystem.ServiceStop(xcontainerdCommand); err != nil {
		fmt.Printf("[xcontainerd-start] WARNING: unable to stop the xcontainerd service momentarily: [%v]\n", err)
	}
}

// TryRestartXContainerd attempts to restart the xcontainerd service
func TryRestartXContainerd() {
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[xcontainerd-start] no supported init system detected, won't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists(xcontainerdCommand) {
		fmt.Println("[xcontainerd-start] couldn't detect a xcontainerd service, can't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
	}

	// This runs "systemctl daemon-reload && systemctl restart xcontainerd"
	if err := initSystem.ServiceRestart(xcontainerdCommand); err != nil {
		fmt.Printf("[xcontainerd-start] WARNING: unable to restart the xcontainerd service momentarily: [%v]\n", err)
	}
}

// TryRestartXContainerd attempts to restart the xcontainerd service
func TryStartXContainerd() {
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[xcontainerd-start] no supported init system detected, won't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists(xcontainerdCommand) {
		fmt.Println("[xcontainerd-start] couldn't detect a xcontainerd service, can't make sure the xcontainerd not running for a short period of time while setting up configuration for it.")
	}

	// This runs "systemctl daemon-reload && systemctl restart xcontainerd"
	if err := initSystem.ServiceStart(xcontainerdCommand); err != nil {
		fmt.Printf("[xcontainerd-start] WARNING: unable to stop the xcontainerd service momentarily: [%v]\n", err)
	}
}
