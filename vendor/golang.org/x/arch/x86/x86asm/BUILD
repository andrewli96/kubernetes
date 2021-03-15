load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "decode.go",
        "gnu.go",
        "inst.go",
        "intel.go",
        "plan9x.go",
        "tables.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/golang.org/x/arch/x86/x86asm",
    importpath = "golang.org/x/arch/x86/x86asm",
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
