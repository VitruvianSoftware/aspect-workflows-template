# No Bazel targets here, by design.
#
# This directory is a STANDALONE Pulumi Go module (its own go.mod), gated on
# copybara + copybara_pulumi_auth. Unlike the sibling repo_config module, it is
# intentionally a LIBRARY (package copybara_sync, exporting ManageSyncAuth) with
# no `func main()` and no Pulumi.yaml, and it is deliberately kept out of the
# monorepo's go.work / gazelle setup (see the `# gazelle:exclude infrastructure`
# directive in the root BUILD).
#
# It therefore has NO `bazel run` wrappers: the //tools/pulumi:defs.bzl
# `pulumi_project` macro only fits a runnable `package main` project with a
# Pulumi.yaml (like repo_config), so applying it here would produce :preview /
# :up targets that cannot run. Drive this module manually per README.md --
# either import ManageSyncAuth into a Pulumi stack's main, or add a main.go and
# run `pulumi up` directly.
