load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "arch_util.go",
        "elf_helper.go",
        "hook.go",
        "utility.go",
        "utility_darwin.go",
        "utility_linux.go",
        "utility_unix.go",
        "utility_windows.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/brahma-adshonor/gohook",
    importpath = "github.com/brahma-adshonor/gohook",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/arch/x86/x86asm:go_default_library"],
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
