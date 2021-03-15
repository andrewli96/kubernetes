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

  add_ldflag "${KUBE_GO_PACKAGE}/pkg/xkube.internal.xkubeEnable" "${XKUBE_ENABLE}"

  # The -ldflags parameter takes a single string, so join the output.
  echo "${ldflags[*]-}"
}


xkube::gcflags() {
  local -a gcflags

  if [[ ${XKUBE_ENABLE} == "true" ]]; then
    # Disable inline for for API hooking
    gcflags+="all=-l"
  fi

  echo "${gcflags[*]-}"
}
