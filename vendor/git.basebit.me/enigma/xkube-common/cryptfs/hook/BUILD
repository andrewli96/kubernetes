load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "dir.go",
        "dir_go115.go",
        "dir_go116.go",
        "file.go",
        "hook.go",
        "hook_go115.go",
        "hook_go116.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/git.basebit.me/enigma/xkube-common/cryptfs/hook",
    importpath = "git.basebit.me/enigma/xkube-common/cryptfs/hook",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/git.basebit.me/enigma/sqlfs-go:go_default_library",
        "//vendor/git.basebit.me/enigma/xkube-common/cryptfs:go_default_library",
        "//vendor/github.com/brahma-adshonor/gohook:go_default_library",
        "//vendor/github.com/cornelk/hashmap:go_default_library",
        "//vendor/k8s.io/klog/v2:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
