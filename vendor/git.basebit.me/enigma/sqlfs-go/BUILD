load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "file.go",
        "sqlfs.go",
        "stat.go",
        "stat_darwin.go",
        "stat_linux.go",
    ],
    cgo = True,
    clinkopts = ["-lsqlfs-1.0"],
    copts = ["-D_FILE_OFFSET_BITS=64 -DHAVE_LIBSQLCIPHER"],
    importmap = "k8s.io/kubernetes/vendor/git.basebit.me/enigma/sqlfs-go",
    importpath = "git.basebit.me/enigma/sqlfs-go",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/mattn/go-pointer:go_default_library"],
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
