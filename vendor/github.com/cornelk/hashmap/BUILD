load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "hashmap.go",
        "hashmap_get.go",
        "list.go",
        "listelement.go",
        "util.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/cornelk/hashmap",
    importpath = "github.com/cornelk/hashmap",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/dchest/siphash:go_default_library"],
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
