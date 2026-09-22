"""Project Swift build rules.

Wraps rules_swift's `swift_binary` to bundle the hermetic Swift runtime on Linux.

The hermetic Swift toolchain (configured via the `swift` bzlmod extension in
//MODULE.bazel) downloads a standalone toolchain from swift.org and sets each
binary's RUNPATH to the toolchain's `lib/swift/<os>` dir -- but it does NOT stage
the runtime `.so`s (libswiftCore.so, libswiftSwiftOnoneSupport.so, ...) into the
binary's runfiles. As a result a bare `bazel run` / `bazel test` of a swift_binary
fails to load the Swift runtime. This wrapper adds the toolchain files as a data
dependency on Linux so the runtime is present in runfiles. macOS resolves the
runtime through Xcode, so no data dep is added there.

It also turns off `swift.static_stdlib`. rules_swift 4.x enables that feature by
default on Linux, and the static link then needs host `libcurl`/`libxml2`
development libraries (`-lcurl -lxml2`), which breaks the "no host packages"
promise of the hermetic toolchain. Linking the runtime dynamically (staged via
the data dep above) keeps builds self-contained.

Remove this wrapper (and the config_settings in //swift:BUILD) once rules_swift
stages the runtime libraries automatically.
"""

load("@build_bazel_rules_swift//swift:swift_binary.bzl", _swift_binary = "swift_binary")

_SWIFT_RUNTIME = select({
    "//swift:linux_x86_64": ["@swift_toolchain_ubuntu24.04//:files"],
    "//swift:linux_arm64": ["@swift_toolchain_ubuntu24.04-aarch64//:files"],
    "//conditions:default": [],
})

def swift_binary(name, data = [], features = [], **kwargs):
    """Drop-in `swift_binary` that bundles the hermetic Swift runtime on Linux."""
    _swift_binary(
        name = name,
        data = data + _SWIFT_RUNTIME,
        features = features + ["-swift.static_stdlib"],
        **kwargs
    )
