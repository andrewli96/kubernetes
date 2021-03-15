load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "blocks.go",
        "blocks_amd64.s",
        "blocks_arm.s",
        "hash.go",
        "hash128.go",
        "hash128_amd64.s",
        "hash128_arm.s",
        "hash_amd64.s",
        "hash_arm.s",
        "hash_asm.go",
        "siphash.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/dchest/siphash",
    importpath = "github.com/dchest/siphash",
    visibility = ["//visibility:public"],
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
