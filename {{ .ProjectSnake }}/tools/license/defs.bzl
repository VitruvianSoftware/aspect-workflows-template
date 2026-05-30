"""
Starlark rule and macro for addlicense run targets.

Provides `addlicense_run_target`, a macro that generates a `bazel run`
target wrapping addlicense so that it operates on the workspace root
regardless of where Bazel's runfiles directory lives.
"""

def _addlicense_launcher_impl(ctx):
    """Create a launcher shell script that cd's to the workspace before running addlicense."""
    addlicense = ctx.executable.addlicense
    out = ctx.actions.declare_file(ctx.attr.name + ".sh")

    # The launcher script:
    #  1. cd's to $BUILD_WORKSPACE_DIRECTORY (set by `bazel run` to the workspace root)
    #  2. Resolves the addlicense binary via the BASH_RUNFILES_DIR / RUNFILES_MANIFEST_FILE
    #     fallback chain from the Bazel bash runfiles library; simplified here by using
    #     a path derived from the script's own location.
    #
    # We embed the runfiles-relative path of the addlicense binary (computed at build
    # time via ctx.file.addlicense.short_path) so the launcher can find it at runtime
    # without needing an external runfiles library.
    #
    # short_path for external deps is "../<repo>/<path>" (relative from main repo root).
    # $RUNFILES_DIR is the runfiles tree root, so the path under it is "<repo>/<path>"
    # (strip the leading "../").
    short_path = addlicense.short_path
    runfiles_rel = short_path[3:] if short_path.startswith("../") else short_path

    ctx.actions.write(
        output = out,
        is_executable = True,
        content = """\
#!/usr/bin/env bash
set -euo pipefail

# Resolve addlicense binary from runfiles.
# $RUNFILES_DIR is set by Bazel for the executable's runfiles tree root.
# The manifest fallback handles environments where RUNFILES_DIR is not set.
if [[ -n "${RUNFILES_DIR:-}" ]]; then
  ADDLICENSE="$RUNFILES_DIR/%s"
elif [[ -f "${BASH_SOURCE[0]}.runfiles_manifest" ]]; then
  ADDLICENSE=$(grep -m1 '^%s ' "${BASH_SOURCE[0]}.runfiles_manifest" | awk '{print $2}')
else
  echo "Cannot locate addlicense binary" >&2; exit 1
fi

# cd to the workspace root (bazel run sets BUILD_WORKSPACE_DIRECTORY).
cd "$BUILD_WORKSPACE_DIRECTORY"

exec "$ADDLICENSE" %s .
""" % (runfiles_rel, runfiles_rel, ctx.attr.addlicense_args),
    )

    runfiles = ctx.runfiles(files = [addlicense])
    runfiles = runfiles.merge(ctx.attr.addlicense[DefaultInfo].default_runfiles)

    return [DefaultInfo(
        executable = out,
        runfiles = runfiles,
    )]

_addlicense_launcher = rule(
    implementation = _addlicense_launcher_impl,
    attrs = {
        "addlicense": attr.label(
            executable = True,
            cfg = "exec",
            mandatory = True,
            doc = "The addlicense binary label.",
        ),
        "addlicense_args": attr.string(
            mandatory = True,
            doc = "Flags to pass to addlicense (copyright, license, check, ignores).",
        ),
    },
    executable = True,
    doc = "Generates a workspace-scanning addlicense launcher for `bazel run`.",
)

def addlicense_run_target(name, copyright, license_flag, ignore_flags, check_mode, visibility = None):
    """Generate a `bazel run` target that invokes addlicense over the workspace.

    The generated target cd's to $BUILD_WORKSPACE_DIRECTORY (set by `bazel run`)
    before scanning, so it operates on the real source tree rather than the
    Bazel runfiles directory.

    Args:
        name: target name.
        copyright: copyright holder string (passed as addlicense -c flag).
        license_flag: addlicense -l value (apache / mit / bsd / mpl).
        ignore_flags: space-joined string of -ignore flag pairs.
        check_mode: if True, pass -check (verify only; don't modify files).
        visibility: standard Bazel visibility list.
    """
    check_arg = "-check " if check_mode else ""
    args = "-c '{copyright}' -l {lid} {check}{ignores}".format(
        copyright = copyright,
        lid = license_flag,
        check = check_arg,
        ignores = ignore_flags,
    )

    _addlicense_launcher(
        name = name,
        addlicense = "@multitool//tools/addlicense",
        addlicense_args = args,
        visibility = visibility,
    )
