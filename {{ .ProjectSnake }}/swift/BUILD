# Platform config_settings consumed by //swift:defs.bzl to inject the hermetic
# Swift runtime libraries on Linux (see the swift_binary wrapper). macOS gets the
# runtime from Xcode, so it falls through to the select() default (no data dep).
config_setting(
    name = "linux_x86_64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:x86_64",
    ],
    visibility = ["//visibility:public"],
)

config_setting(
    name = "linux_arm64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:arm64",
    ],
    visibility = ["//visibility:public"],
)
