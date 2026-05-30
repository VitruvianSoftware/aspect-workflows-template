"""`bazel run` wrappers for the copybara_sync Pulumi project.

This project (gated on copybara + copybara_pulumi_auth) provisions the
GitHub-App + deploy-key auth IaC used by the Copybara bidi sync (see README.md).
Developers drive it through Bazel instead of the raw `pulumi` CLI:

    bazel run //infrastructure/pulumi/copybara_sync:setup            # one-time bootstrap
    bazel run //infrastructure/pulumi/copybara_sync:preview
    bazel run //infrastructure/pulumi/copybara_sync:up

Pulumi compiles and runs the program itself; the wrappers only launch it.
"""

load("//tools/pulumi:defs.bzl", "pulumi_project")

pulumi_project(
    name = "copybara_sync",
    dir = "infrastructure/pulumi/copybara_sync",
)
