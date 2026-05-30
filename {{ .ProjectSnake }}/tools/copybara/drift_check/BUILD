load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "drift_check_lib",
    srcs = ["main.go"],
    importpath = "example.com/{{ .ProjectSnake }}/tools/copybara/drift_check",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "drift_check",
    embed = [":drift_check_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "drift_check_test",
    srcs = ["main_test.go"],
    embed = [":drift_check_lib"],
)
