# `//:tidy` — one-command BUILD/source hygiene + staleness enforcement

**Date:** 2026-05-31
**Status:** Approved (design)
**Scope:** aspect-workflows-template generated tree (`{{ .ProjectSnake }}/`), all presets

## Problem

Generated repos are Bazel monorepos. Keeping them correct requires two routine
chores:

1. **Regenerating BUILD files** after adding/moving/deleting source files —
   `bazel run //:gazelle` (plus `//:gazelle_python_manifest.update` when Python
   is enabled).
2. **Formatting** BUILD files and sources — `bazel run //tools/format`
   (`format_multirun`, which already includes buildifier/Starlark plus every
   per-language formatter).

Today:

- **Formatting is well covered.** `//tools/format` exists and the `githooks/pre-commit`
  hook already formats staged files and fails if the formatter changed them.
- **Gazelle staleness is NOT covered.** The pre-commit hook never runs gazelle,
  and there is no in-repo CI check for it. A developer who adds a source file and
  forgets `bazel run //:gazelle` ships stale BUILD files; the failure surfaces
  later as a confusing build break (often in a teammate's checkout or in a
  downstream Aspect Workflows run), not at the point of the mistake.
- **There is no single "fix it all" command.** Developers must remember and run
  `//:gazelle`, `//:gazelle_python_manifest.update`, and `//tools/format`
  separately.

## Goal

Give every generated repo:

1. A single `bazel run //:tidy` that makes BUILD files and formatting correct in
   one shot.
2. A CI gate that fails a PR (with the exact fix command) when BUILD files are
   stale or unformatted.
3. A pre-commit hook that catches gazelle staleness locally, before the commit
   lands.

This is **always-on** for every preset — every generated repo is a Bazel
workspace, and the pieces `tidy` wraps (`//:gazelle`, `//tools/format`) are
already always-on. No new scaffold toggle.

## Design

### Component 1 — `//:tidy` aggregator

Add a `multirun` target to the generated-tree root `BUILD`
(`{{ .ProjectSnake }}/BUILD`) that runs, sequentially (`jobs = 1`):

1. `//:gazelle`
2. `//:gazelle_python_manifest.update` — **templated, only emitted when
   `{{ "{{ .Computed.python }}" }}`** (the manifest target only exists for Python)
3. `//tools/format`

```starlark
load("@rules_multirun//:defs.bzl", "multirun")

multirun(
    name = "tidy",
    commands = [
        ":gazelle",
        # (python only) ":gazelle_python_manifest.update",
        "//tools/format",
    ],
    jobs = 1,  # sequential: gazelle regenerates BUILD files, then format formats them
)
```

`rules_multirun` is already in the module graph transitively (via
`aspect_rules_lint`'s `format_multirun`). If the `load(...)` requires it as a
direct dependency, add `bazel_dep(name = "rules_multirun", ...)` to
`MODULE.bazel` pinned to the version already resolved in `MODULE.bazel.lock`.

Order rationale: gazelle first (it writes BUILD files), format last (it formats
the freshly generated BUILD files and all sources).

### Component 2 — CI staleness gate

New templated workflow `{{ .ProjectSnake }}/.github/workflows/tidy-check.yaml`,
modeled on the existing `license-check.yml`:

- Same born-green SPDX header block (`{{ "{{ if .Scaffold.license }}" }}` … per
  `license_id`) as the other templated workflows.
- Same SHA-pinned `bazel-contrib/setup-bazel` action and the same triggers as
  `license-check.yml` (`pull_request`, and `push` to the default branch).
- Steps: check out, set up Bazel, `bazel run //:tidy`, then `git diff --exit-code`.
- On failure: print `git diff --stat` and the exact remediation —
  `bazel run //:tidy && git add -A` — then exit non-zero so the check is red.

This catches stale/unformatted BUILD files at PR time with an actionable message
rather than as a downstream break. It complements (does not replace) Aspect
Workflows' build/test CI.

### Component 3 — pre-commit gazelle

Edit `{{ .ProjectSnake }}/githooks/pre-commit`. Keep the existing staged-file
format step unchanged. Add a gazelle step (after the format step):

- Run gazelle (prefer the `bazel_env` binary if present, else `bazel run //:gazelle`),
  matching the existing hook's "use the prebuilt binary if available" pattern.
- If gazelle modified any BUILD files (`git diff` over `**/BUILD` / `**/BUILD.bazel`
  is non-empty), print "gazelle updated BUILD files — review and stage them",
  show `git diff --stat`, and `exit 1` — the same review-and-restage UX the hook
  already uses for formatting.

Accepted trade-off: this adds a gazelle run to each commit (a few seconds).
Developers can bypass with `git commit --no-verify` when needed. The hook only
runs when `core.hooksPath` is set to `githooks` (the repo's documented setup;
`check-config.sh` already nudges this).

## Born-tidy invariant (primary correctness risk)

The CI gate means **every freshly generated repo must already be clean under
`bazel run //:tidy`** (empty `git diff`). If the template ships any BUILD file
that gazelle would rewrite or that is unformatted, `tidy-check` fails on a brand
new repo on day one.

Therefore a required part of this work is making the generated tree **born-tidy**:
render each preset, run `bazel run //:tidy`, and confirm the diff is empty. Any
file gazelle/format would change is a pre-existing latent inconsistency in the
template's shipped BUILD files and must be fixed at the source (the same
born-green principle used for the license headers). The existing
`# gazelle:exclude` directives (`infrastructure`, `tools/pulumi`, `githooks/*`,
`**/*.venv`) are respected by `//:gazelle` and therefore by `//:tidy`.

## Files

- **Modify:** `{{ .ProjectSnake }}/BUILD` — add the `multirun` `tidy` target (+ its `load`).
- **Modify:** `{{ .ProjectSnake }}/MODULE.bazel` — add `rules_multirun` `bazel_dep` only if the `load` requires a direct dep.
- **Create:** `{{ .ProjectSnake }}/.github/workflows/tidy-check.yaml` — the CI gate (born-green).
- **Modify:** `{{ .ProjectSnake }}/githooks/pre-commit` — add the gazelle staleness step.
- **Modify:** `docs/user-guide/development-workflow.md` — document `bazel run //:tidy` and the CI gate (cross-reference it from `docs/user-guide/formatting-linting.md`).
- **Possibly fix:** any shipped BUILD files that are not born-tidy (discovered during verification).

## Verification

1. Render `kitchen-sink`, one or two single-language presets, `minimal`, and a
   `python`-bearing preset (exercises the manifest branch). For each:
   - `//:tidy` builds and `bazel run //:tidy` runs cleanly with an **empty**
     `git diff` (born-tidy).
   - `tidy-check.yaml` renders to valid YAML, no stray `{{ "{{ }}" }}`, with the
     correct SPDX header for the chosen `license_id`.
   - `githooks/pre-commit` renders to valid bash.
2. Negative check: introduce a deliberately stale BUILD (e.g., delete a `srcs`
   entry or add an untracked Go file without running gazelle), confirm
   `bazel run //:tidy` fixes it and that `git diff --exit-code` would fail before
   the fix.
3. `deliver`-matrix (13 presets) stays green on `platform-v2.0`.
4. Dogfood on vitruvian-core: port the target/workflow/hook, run `//:tidy`,
   confirm born-tidy and a green `tidy-check`.

## Rollout

1. Implement + verify on `platform-v2.0`.
2. PR → (gated) merge → `deliver` fans out `//:tidy` + `tidy-check.yaml` + the
   updated hook to the 26 Starters.
3. Port to vitruvian-core and verify end-to-end; fix any gap back in the template
   and re-deliver (standing practice).

## Out of scope (YAGNI)

- No new scaffold question/toggle (always-on).
- No changes to Aspect Workflows' own build/test CI.
- No RBE for `tidy` (it mutates the local tree; RBE adds nothing).
- Not folding a full `//:tidy` into the pre-commit hook (only gazelle is added;
  the staged-file format step stays scoped and fast).
