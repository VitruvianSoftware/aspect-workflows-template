"""Targets in the repository root"""

{{ if .Computed.javascript }}
load("@aspect_rules_js//js:defs.bzl", "js_library")
load("@npm//:defs.bzl", "npm_link_all_packages")

# TODO: remove once https://github.com/aspect-build/aspect-cli/issues/560 done
# gazelle:js_npm_package_target_name pkg
npm_link_all_packages(name = "node_modules")
{{ if .Scaffold.lint }}

js_library(
    name = "eslintrc",
    srcs = ["eslint.config.mjs"],
    visibility = ["//:__subpackages__"],
    deps = [
        ":node_modules/@eslint/js",
        ":node_modules/typescript-eslint",
    ],
)

js_library(
    name = "prettier_config",
    srcs = ["prettier.config.cjs"],
    deps = [],
    visibility = ["//tools/format:__pkg__"],
)

js_library(
    name = "prettier_ignore",
    srcs = [".prettierignore"],
    deps = [],
    visibility = ["//tools/format:__pkg__"],
)
{{ end }}
{{ end }}

{{ if .Scaffold.lint}}
exports_files([
{{ if .Computed.cpp }}
    ".clang-tidy",
{{ end }}
{{ if .Computed.ruby }}
    ".rubocop.yml",
{{ end }}
{{ if .Computed.kotlin }}
    "ktlint-baseline.xml",
    ".editorconfig",
{{ end }}
{{ if .Computed.java }}
    "pmd.xml",
{{ end }}
{{ if .Computed.shell }}
    ".shellcheckrc",
{{ end }}
{{ if .Computed.rust }}
    ".clippy.toml",
{{ end }}
], visibility = ["//:__subpackages__"])
{{ end }}

load("@gazelle//:def.bzl", "gazelle")

{{ if .Computed.go }}
# gazelle:prefix github.com/example/project
{{ end }}
{{ if .Computed.javascript }}
# It's faster to avoid type-checking in a devserver when using monorepo packages.
# If you commonly ship your npm packages outside the repo, change this to "npm_package"
# gazelle:js_package_rule_kind js_library
{{ end }}
# We prefer BUILD instead of BUILD.bazel
# gazelle:build_file_name BUILD
# gazelle:exclude githooks/*
# Standalone Pulumi IaC modules under infrastructure/pulumi/* are run via the
# pulumi CLI / //tools/pulumi bazel wrappers — not built by Bazel — so keep
# gazelle out of them (it otherwise collides with the pulumi_project targets).
# gazelle:exclude infrastructure
# tools/pulumi holds a hand-authored macro (defs.bzl) + wrapper scripts; gazelle
# would otherwise generate a bzl_library with an unresolvable rules_shell dep.
# gazelle:exclude tools/pulumi
{{ if .Scaffold.proto -}}
# Workaround https://github.com/bazel-contrib/bazel-gazelle/issues/2001
# gazelle:map_kind proto_library proto_library @protobuf//bazel:proto_library.bzl
{{- end }}
gazelle(
    name = "gazelle",
    gazelle = "@multitool//tools/gazelle",
    env = {
        "ENABLE_LANGUAGES": ",".join([
            "starlark",
            {{- if .Scaffold.proto }}"proto",{{ end }}
            {{- if .Computed.go }}"go",{{ end }}
            {{- if .Computed.kotlin }}"kotlin",{{ end }}
            {{- if .Computed.python }}"python",{{ end }}
            {{- if .Computed.javascript }}"js",{{ end }}
            {{- if .Computed.cpp }}"cc",{{ end }}
        ]),
    }
)

load("@rules_multirun//:defs.bzl", "multirun")

# One-command BUILD/source hygiene. Regenerates BUILD files (gazelle), refreshes
# the Python deps manifest (when Python is enabled), then formats everything
# (//tools/format bundles buildifier + every per-language formatter; it ships
# only when the `lint` feature is on). Sequential so the formatter sees gazelle's
# freshly written BUILD files. Run `bazel run //:tidy`; the Tidy Check CI job
# fails a PR when running this would change anything.
multirun(
    name = "tidy",
    commands = [
        ":gazelle",
{{- if .Computed.python }}
        ":gazelle_python_manifest.update",
{{- end }}
{{- if .Scaffold.lint }}
        "//tools/format",
{{- end }}
    ],
    jobs = 1,  # sequential: gazelle writes BUILD files, then format formats them
)

{{ if .Computed.python }}
load("@pip//:requirements.bzl", "all_whl_requirements")
load("@rules_python_gazelle_plugin//manifest:defs.bzl", "gazelle_python_manifest")
load("@rules_python_gazelle_plugin//modules_mapping:def.bzl", "modules_mapping")

exports_files(
    ["pyproject.toml"],
    visibility = ["//:__subpackages__"],
)

# Produce aspect_rules_py targets rather than rules_python
# gazelle:map_kind py_binary py_binary @aspect_rules_py//py:defs.bzl
# gazelle:map_kind py_library py_library @aspect_rules_py//py:defs.bzl
# gazelle:map_kind py_test py_test //tools/pytest:defs.bzl
#
# Don't walk into virtualenvs when looking for python sources.
# We don't intend to plant BUILD files there.
# gazelle:exclude **/*.venv
#
# Fetches metadata for python packages we depend on.
modules_mapping(
    name = "modules_map",
    wheels = all_whl_requirements,
)

# Provide a mapping from an import to the installed package that provides it.
# Needed to generate BUILD files for .py files.
# This macro produces two targets:
# - //:gazelle_python_manifest.update can be used with `bazel run`
#   to recalculate the manifest
# - //:gazelle_python_manifest.test is a test target ensuring that
#   the manifest doesn't need to be updated
gazelle_python_manifest(
    name = "gazelle_python_manifest",
    modules_mapping = ":modules_map",
    pip_repository_name = "pip",
)

{{ end }}
