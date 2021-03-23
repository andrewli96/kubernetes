/*
Copyright 2017 The Kubernetes Authors.

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

package preflight

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/version"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	utilsexec "k8s.io/utils/exec"
)

// GetKubeletVersion is helper function that returns version of kubelet available in $PATH
func GetKubeletVersion(execer utilsexec.Interface) (*version.Version, error) {
	kubeletVersionRegex := regexp.MustCompile(`^\s*Kubernetes v((0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)([-0-9a-zA-Z_\.+]*)?)\s*$`)

	command := execer.Command(kubeadmconstants.XKubeletCommand, "--version")
	out, err := command.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute '%s --version'", kubeadmconstants.XKubeletCommand)
	}

	cleanOutput := strings.TrimSpace(string(out))
	subs := kubeletVersionRegex.FindAllStringSubmatch(cleanOutput, -1)
	if len(subs) != 1 || len(subs[0]) < 2 {
		return nil, errors.Errorf("Unable to parse output from Kubelet: %q", cleanOutput)
	}
	return version.ParseSemantic(subs[0][1])
}
