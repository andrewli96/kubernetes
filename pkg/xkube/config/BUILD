load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["xconfig.go"],
    importpath = "k8s.io/kubernetes/pkg/xkube/config",
    visibility = ["//visibility:public"],
    deps = ["//vendor/git.basebit.me/enigma/sqlfs-go:go_default_library"],
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
