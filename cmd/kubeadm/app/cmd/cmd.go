/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"io"
	"path/filepath"
	"syscall"

	"git.basebit.me/enigma/xkube-common/cryptfs"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/alpha"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/upgrade"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/pkg/xkube"
	xkubeoptions "k8s.io/kubernetes/pkg/xkube/options"
	// Register the kubeadm configuration types because CLI flag generation
	// depends on the generated defaults.
)

var (
	hookPaths       []string
	hookParentDirs  []string
	xOptions        *xkubeoptions.XOptions
	xKubeadmOptions *XKubeadmOptions
	sqlfsMkdirMode  uint32
)

type XKubeadmOptions struct {
	X           *xkubeoptions.XOptions
	InitOptions *initOptions
	InitData    *initData
}

// NewKubeadmCommand returns cobra.Command to run kubeadm command
func NewKubeadmCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var rootfsPath string
	xOptions = xkubeoptions.NewXOptions()
	xKubeadmOptions = &XKubeadmOptions{
		X: xOptions,
	}

	cmds := &cobra.Command{
		Use:   "kubeadm",
		Short: "kubeadm: easily bootstrap a secure Kubernetes cluster",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ KUBEADM                                                  │
			    │ Easily bootstrap a secure Kubernetes cluster             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/kubernetes/kubeadm/issues             │
			    └──────────────────────────────────────────────────────────┘

			Example usage:

			    Create a two-machine cluster with one control-plane node
			    (which controls the cluster), and one worker node
			    (where your workloads, like Pods and Deployments run).

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the first machine:                                    │
			    ├──────────────────────────────────────────────────────────┤
			    │ control-plane# kubeadm init                              │
			    └──────────────────────────────────────────────────────────┘

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the second machine:                                   │
			    ├──────────────────────────────────────────────────────────┤
			    │ worker# kubeadm join <arguments-returned-from-init>      │
			    └──────────────────────────────────────────────────────────┘

			    You can then repeat the second step on as many other machines as you like.

		`),
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if rootfsPath != "" {
				if err := kubeadmutil.Chroot(rootfsPath); err != nil {
					return err
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			xkube.Close()
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(newCmdCompletion(out, ""))
	cmds.AddCommand(newCmdConfig(out))
	cmds.AddCommand(newCmdInit(out, nil))
	cmds.AddCommand(newCmdJoin(out, nil))
	cmds.AddCommand(newCmdReset(in, out, nil))
	cmds.AddCommand(newCmdVersion(out))
	cmds.AddCommand(newCmdToken(out, err))
	cmds.AddCommand(upgrade.NewCmdUpgrade(out))
	cmds.AddCommand(alpha.NewCmdAlpha(in, out))
	options.AddKubeadmOtherFlags(cmds.PersistentFlags(), &rootfsPath)
	options.AddXKubeConfigFlag(cmds.PersistentFlags(), &xOptions.XConfigFile)

	// TODO: remove "certs" from "alpha"
	// https://github.com/kubernetes/kubeadm/issues/2291
	cmds.AddCommand(alpha.NewCmdCertsUtility(out))

	return cmds
}

// Find all the files that apiserver needs to operate(read & write)
func getCryptfsHookedFiles(xKubeadmOptions *XKubeadmOptions) ([]cryptfs.MatchPattern, error) {
	var patterns []cryptfs.MatchPattern

	klog.Infof("InitOptions=%+v", xKubeadmOptions.InitOptions)

	// Command Init Hook Procedure
	if xKubeadmOptions.InitData != nil {
		// Hook certs
		for _, cert := range certsphase.GetDefaultCertList() {
			crtFileName := filepath.Join(xKubeadmOptions.InitData.certificatesDir, cert.BaseName+".crt")
			keyFileName := filepath.Join(xKubeadmOptions.InitData.certificatesDir, cert.BaseName+".key")
			hookPaths = append(hookPaths, crtFileName)
			hookPaths = append(hookPaths, keyFileName)
		}

		// Hook Service Account
		saKeyFileName := filepath.Join(xKubeadmOptions.InitData.certificatesDir, kubeadmconstants.ServiceAccountPrivateKeyName)
		saPubFileName := filepath.Join(xKubeadmOptions.InitData.certificatesDir, kubeadmconstants.ServiceAccountPublicKeyName)
		hookPaths = append(hookPaths, saKeyFileName)
		hookPaths = append(hookPaths, saPubFileName)

		// Hook Kubeconfig
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.kubeconfigDir, kubeadmconstants.AdminKubeConfigFileName))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.kubeconfigDir, kubeadmconstants.KubeletBootstrapKubeConfigFileName))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.kubeconfigDir, kubeadmconstants.KubeletKubeConfigFileName))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.kubeconfigDir, kubeadmconstants.ControllerManagerKubeConfigFileName))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.kubeconfigDir, kubeadmconstants.SchedulerKubeConfigFileName))

		// Hook Manifests Pod Spec
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.ManifestDir(), kubeadmconstants.KubeAPIServer+".yaml"))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.ManifestDir(), kubeadmconstants.KubeControllerManager+".yaml"))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.ManifestDir(), kubeadmconstants.KubeScheduler+".yaml"))
		hookPaths = append(hookPaths, filepath.Join(xKubeadmOptions.InitData.ManifestDir(), kubeadmconstants.Etcd+".yaml"))

		// Hook Kubelet
		hookParentDirs = append(hookParentDirs, xKubeadmOptions.InitData.KubeletDir())
	}

	for _, path := range hookPaths {
		if path != "" {
			patterns = append(patterns, cryptfs.MatchPattern{
				Mode:  cryptfs.MATCH_EXACT,
				Value: path,
			})
		}
	}

	for _, path := range hookParentDirs {
		if path != "" {
			patterns = append(patterns, cryptfs.MatchPattern{
				Mode:  cryptfs.MATCH_PARENT,
				Value: path,
			})
		}
	}

	return patterns, nil
}

func xMustMakeDirAll(fs *cryptfs.CryptFs, path string, mode uint32) {
	// Single Thread Disable Lock
	//fs.SFuckingMu.Lock()
	//defer fs.SFuckingMu.Unlock()
	if path == "/" {
		return
	}
	// recursively create from parent's path
	xMustMakeDirAll(fs, filepath.Dir(path), mode)
	// sfs mkdir operation
	err := fs.SFS.Mkdir(path, mode)
	if err == nil {
		return
	} else if err != nil && err == syscall.Errno(syscall.EEXIST) {
		return
	} else {
		klog.Errorf("SQLFS failed to create directory %s. Error: %s.", path, err)
		panic(err)
	}
}
