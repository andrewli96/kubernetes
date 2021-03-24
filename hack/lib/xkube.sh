#!/usr/bin/env bash


xkube::ldflags() {
  local -a ldflags
  function add_ldflag() {
    local key=${1}
    local val=${2}
    ldflags+=(
      "-X '${key}=${val}'"
    )
  }

  add_ldflag "${KUBE_GO_PACKAGE}/pkg/xkube/internal.xkubeEnabled" "${XKUBE_ENABLED}"

  if [[ $(go env GOOS) != "darwin" ]]; then
    ldflags+=' -extldflags "-static -lsqlcipher -ltomcrypt -ldl -lpthread"'
  fi

  # The -ldflags parameter takes a single string, so join the output.
  echo "${ldflags[*]-}"
}


xkube::gcflags() {
  local -a gcflags

  if [[ ${XKUBE_ENABLED} == "1" ]]; then
    # Disable inline for for API hooking
    gcflags+="all=-l"
  fi

  echo "${gcflags[*]-}"
}



xkube::gotags() {
# Turn on build tag "osusergo" to eliminiate warnings like below:
# /usr/bin/ld: /tmp.k8s/go-link-977064025/000019.o: in function `mygetgrouplist':
# /usr/local/go/src/os/user/getgrouplist_unix.go:16: warning: Using 'getgrouplist' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking

# Turn on build tag "netgo" to eliminiate warnings like below:
# /usr/bin/ld: /tmp.k8s/go-link-004397297/000004.o: in function `_cgo_26061493d47f_C2func_getaddrinfo':
# /tmp/go-build/cgo-gcc-prolog:58: warning: Using 'getaddrinfo' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking

  if [[ $(go env GOOS) != "darwin" ]]; then
    echo "osusergo,netgo"
  fi
}
