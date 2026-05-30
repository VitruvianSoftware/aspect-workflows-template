"""
License header enforcement via addlicense.

Usage:
  bazel run //tools/license:add    — add/fix missing license headers
  bazel run //tools/license:check  — verify all headers present (exit non-zero if missing)

Both targets are `bazel run` targets (not `bazel test`). A generated launcher
script cd's to $BUILD_WORKSPACE_DIRECTORY before calling addlicense so the tool
scans the real workspace tree (Bazel 7+ runs binaries from the runfiles
directory, not the workspace root; $BUILD_WORKSPACE_DIRECTORY holds the
workspace path).

NOTE: `:check` must be invoked via `bazel run`, not `bazel test`, because
`bazel test` runs in a sandbox with no access to the live source tree. CI should
invoke `bazel run //tools/license:check` (exit non-zero = missing headers).
"""

load("//tools/license:defs.bzl", "addlicense_run_target")

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

# Globs passed as -ignore flags to addlicense (generated / tool-managed files
# that should not receive license headers).  Defined as a Starlark string so
# they're inlined verbatim into the generated launcher shell script.
_IGNORE_FLAGS = " ".join([
    "-ignore '**/*.lock'",
    "-ignore '**/pnpm-lock.yaml'",
    "-ignore '**/package-lock.json'",
    "-ignore '**/Cargo.lock'",
    "-ignore 'MODULE.bazel.lock'",
    "-ignore '**/gazelle_python.yaml'",
    "-ignore '**/*-baseline.xml'",
    "-ignore '**/.release-please-manifest.json'",
    "-ignore 'bazel-*/**'",
    "-ignore '**/node_modules/**'",
    "-ignore '**/*.venv/**'",
    "-ignore '.git/**'",
    "-ignore '**/BUILD'",
    "-ignore '**/BUILD.bazel'",
])

# SPDX ID → addlicense -l value mapping (resolved at scaffold generation time).
_LICENSE_FLAG = "{{ if eq .Scaffold.license_id `Apache-2.0` }}apache{{ else if eq .Scaffold.license_id `MIT` }}mit{{ else if eq .Scaffold.license_id `BSD-3-Clause` }}bsd{{ else }}mpl{{ end }}"

# Copyright holder (provided at scaffold generation time).
_COPYRIGHT = "{{ .Scaffold.copyright_holder }}"

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

# bazel run //tools/license:add — write/fix missing license headers.
addlicense_run_target(
    name = "add",
    check_mode = False,
    copyright = _COPYRIGHT,
    ignore_flags = _IGNORE_FLAGS,
    license_flag = _LICENSE_FLAG,
    visibility = ["//visibility:public"],
)

# bazel run //tools/license:check — exit non-zero if any headers are missing.
addlicense_run_target(
    name = "check",
    check_mode = True,
    copyright = _COPYRIGHT,
    ignore_flags = _IGNORE_FLAGS,
    license_flag = _LICENSE_FLAG,
    visibility = ["//visibility:public"],
)
