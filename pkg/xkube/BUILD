load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["xkube.go"],
    importpath = "k8s.io/kubernetes/pkg/xkube",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/xkube/internal:go_default_library",
        "//pkg/xkube/options:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//pkg/xkube/config:all-srcs",
        "//pkg/xkube/internal:all-srcs",
        "//pkg/xkube/options:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
