# `//:tidy` Build-Hygiene Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an always-on `bazel run //:tidy` (gazelle → py-manifest → format) to every generated repo, a CI gate that fails when BUILD files are stale/unformatted, and a pre-commit gazelle step — so developers fix BUILD hygiene in one command and stale BUILDs never land.

**Architecture:** A `rules_multirun` `multirun` target wraps the existing `//:gazelle` and `//tools/format`; a templated GitHub Actions workflow runs `//:tidy` and `git status --porcelain` as a staleness gate; the existing `githooks/pre-commit` gains a gazelle step. All three live in the generated tree (`{{ .ProjectSnake }}/`) and ship to every preset (no scaffold toggle).

**Tech stack:** Bazel + bzlmod, `rules_multirun` 0.13.0, `aspect_rules_lint` `format_multirun` (already present), gazelle (already present), GitHub Actions, hay-kot scaffold templating.

**Templating notes:**
- `{{ "{{ .Computed.python }}" }}` gates the Python-manifest command.
- Generated-tree workflows escape GHA expressions as `${{ "{{ <expr> }}" }}` → renders to `${{ <expr> }}`.
- Always-shipped workflows wrap the SPDX header in the outer `{{ "{{ if .Scaffold.license }}" }} … {{ "{{ end }}" }}` (see `_repo-config-preview.yaml`), because the repo may have no license.

**Born-tidy invariant (the main risk):** the CI gate fails on a brand-new repo unless the template's shipped BUILD files are already gazelle-fresh and formatted. Task 4 verifies and fixes this.

---

### Task 1: `//:tidy` aggregator target

**Files:**
- Modify: `{{ .ProjectSnake }}/MODULE.bazel` (add `rules_multirun` bazel_dep)
- Modify: `{{ .ProjectSnake }}/BUILD` (add `load` + `multirun` target)

- [ ] **Step 1: Add the `rules_multirun` dependency**

In `{{ .ProjectSnake }}/MODULE.bazel`, next to the other `bazel_dep` lines (e.g. after `aspect_rules_lint`), add:

```starlark
bazel_dep(name = "rules_multirun", version = "0.13.0")
```

- [ ] **Step 2: Define `//:tidy` in the root BUILD**

In `{{ .ProjectSnake }}/BUILD`, add a `load` near the existing `load("@gazelle//:def.bzl", "gazelle")` and a target near the `gazelle(name = "gazelle", …)` definition:

```starlark
load("@rules_multirun//:defs.bzl", "multirun")

# One-command BUILD/source hygiene: regenerate BUILD files (gazelle), refresh the
# Python deps manifest (when Python is enabled), then format everything
# (//tools/format already includes buildifier + every per-language formatter).
# Run sequentially so format sees gazelle's freshly written BUILD files.
multirun(
    name = "tidy",
    commands = [
        ":gazelle",
{{ if .Computed.python }}        ":gazelle_python_manifest.update",
{{ end }}        "//tools/format",
    ],
    jobs = 1,  # sequential
)
```

- [ ] **Step 3: Render and build (non-python preset)**

```bash
OUT=$(mktemp -d)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir="$OUT" --preset=go --no-prompt "$(pwd)"
cd "$OUT"/*/
bazel build //:tidy
```
Expected: target builds (rules_multirun fetches, lock updates), no `:gazelle_python_manifest.update` in the rendered `BUILD` (grep to confirm).

- [ ] **Step 4: Run `//:tidy` — must be a clean no-op (born-tidy)**

```bash
bazel run //:tidy
git -C "$OUT"/*/ status --porcelain
```
Expected: empty `git status` (no changes). If non-empty, it is a born-tidy defect → record the files for Task 4.

- [ ] **Step 5: Render a python preset and confirm the manifest command renders**

```bash
OUT2=$(mktemp -d)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir="$OUT2" --preset=py --no-prompt "$(pwd)"
grep -n "gazelle_python_manifest.update" "$(find "$OUT2" -name BUILD -maxdepth 3 | head -1)"
```
Expected: the line is present (python branch rendered).

- [ ] **Step 6: Commit**

```bash
git add "{{ .ProjectSnake }}/MODULE.bazel" "{{ .ProjectSnake }}/BUILD"
git commit -m "feat(tidy): add //:tidy multirun (gazelle + py-manifest + format)"
```

---

### Task 2: CI staleness gate (`tidy-check.yaml`)

**Files:**
- Create: `{{ .ProjectSnake }}/.github/workflows/tidy-check.yaml`

- [ ] **Step 1: Create the workflow (born-green, always-shipped header)**

Header block must match the always-shipped pattern (outer `{{ "{{ if .Scaffold.license }}" }}`); copy the exact header block from `{{ .ProjectSnake }}/.github/workflows/_repo-config-preview.yaml` (lines 1–39, ending in `{{ "{{ end }}{{ end }}" }}`), then the body:

```yaml
name: Tidy Check

# Fails when BUILD files are stale (gazelle) or unformatted. The single fix is
# `bazel run //:tidy`. Complements Aspect Workflows' build/test CI.
on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  tidy-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Install Bazel
        uses: bazel-contrib/setup-bazel@c5acdfb288317d0b5c0bbd7a396a3dc868bb0f86 # 0.19.0
        with:
          bazelisk-cache: true
          repository-cache: true
          bazelrc: |
            common --announce_rc --color=yes --curses=yes
            common --show_progress_rate_limit=60 --show_timestamps

      - name: Run //:tidy and verify no changes
        run: |
          bazel run //:tidy
          if [ -n "$(git status --porcelain)" ]; then
            echo "::error::BUILD files are stale or unformatted. Run 'bazel run //:tidy' and commit the result."
            git status --porcelain
            git diff
            exit 1
          fi
```

- [ ] **Step 2: Render and validate YAML + header + no stray delimiters**

```bash
OUT=$(mktemp -d)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir="$OUT" --preset=license-go --no-prompt "$(pwd)"
F="$(find "$OUT" -name tidy-check.yaml)"
python3 -c "import yaml;yaml.safe_load(open('$F'));print('YAML OK')"
head -1 "$F"                 # expect a license comment line (Apache for license-go)
grep -nE '\{\{|\}\}' "$F" || echo "no stray scaffold delimiters"
```
Expected: YAML OK; header present; the only `${{ }}` are none here (this workflow has no GHA expressions) → grep prints "no stray scaffold delimiters". Render a no-license preset (`go`) too and confirm the header block renders empty (file starts at `name: Tidy Check`).

- [ ] **Step 3: Commit**

```bash
git add "{{ .ProjectSnake }}/.github/workflows/tidy-check.yaml"
git commit -m "feat(tidy): CI staleness gate (tidy-check.yaml)"
```

---

### Task 3: pre-commit gazelle step

**Files:**
- Modify: `{{ .ProjectSnake }}/githooks/pre-commit`

- [ ] **Step 1: Append a gazelle staleness step after the existing format block**

Add to the end of `{{ .ProjectSnake }}/githooks/pre-commit` (after the existing staged-file format `xargs … _` block):

```bash

# Regenerate BUILD files and fail if any changed, so stale BUILD files never get
# committed. Prefer the bazel_env-built gazelle if present (avoids a bazel build),
# matching the formatter step above; otherwise fall back to `bazel run`.
if [ -e "bazel-out/bazel_env-opt/bin/tools/bazel_env/bin/gazelle" ]; then
  bazel-out/bazel_env-opt/bin/tools/bazel_env/bin/gazelle
else
  bazel run //:gazelle
fi
if ! git diff --quiet -- '*BUILD' '*BUILD.bazel'; then
  echo "❌ gazelle updated BUILD files."
  echo "Please review and stage the changes before committing again."
  git diff --stat -- '*BUILD' '*BUILD.bazel'
  exit 1
fi
```

- [ ] **Step 2: Render and shellcheck/bash-validate**

```bash
OUT=$(mktemp -d)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir="$OUT" --preset=kitchen-sink --no-prompt "$(pwd)"
H="$(find "$OUT" -path '*githooks/pre-commit')"
bash -n "$H" && echo "bash syntax OK"
grep -nE '\{\{|\}\}' "$H" || echo "no stray scaffold delimiters"
```
Expected: bash syntax OK; no stray delimiters. (The hook only carries a license header when `license` is on — confirm header renders correctly on `license-go` and is absent on `kitchen-sink`.)

- [ ] **Step 3: Commit**

```bash
git add "{{ .ProjectSnake }}/githooks/pre-commit"
git commit -m "feat(tidy): run gazelle in pre-commit to catch stale BUILD files"
```

---

### Task 4: Born-tidy verification + fixes (the correctness gate)

**Files:**
- Possibly modify: any shipped `{{ .ProjectSnake }}/**/BUILD` files that gazelle/format would rewrite.

- [ ] **Step 1: For each representative preset, render → `//:tidy` → assert clean**

```bash
for p in minimal go py js kitchen-sink swift; do
  OUT=$(mktemp -d)
  SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir="$OUT" --preset="$p" --no-prompt "$(pwd)" >/dev/null 2>&1
  R="$(ls -d "$OUT"/*/)"
  ( cd "$R" && bazel run //:tidy >/dev/null 2>&1; echo "[$p] $(git status --porcelain | wc -l | tr -d ' ') changed files"; git status --porcelain )
done
```
Expected: `0 changed files` for every preset.

- [ ] **Step 2: If any preset reports changes, fix at the source**

For each changed file, inspect the diff: if gazelle rewrote a shipped `BUILD`, update the template's source `BUILD` so the rendered output is already gazelle-fresh; if it was a formatting change, apply buildifier/format to the template source. Re-run Step 1 until every preset is `0 changed files`. (No code can be written generically here — fix exactly what the diff shows.)

- [ ] **Step 3: No commit unless fixes were needed** (born-tidy verification is a check; commit only the fixes from Step 2, if any)

```bash
git status --short  # if Task-4 fixes exist:
git add -A && git commit -m "fix(tidy): make shipped BUILD files born-tidy"
```

---

### Task 5: Document `//:tidy`

**Files:**
- Modify: `{{ .ProjectSnake }}/docs/user-guide/development-workflow.md` (or the generator's equivalent if docs live at generator root — confirm path at exec time)
- Modify: `{{ .ProjectSnake }}/docs/user-guide/formatting-linting.md` (cross-reference)

- [ ] **Step 1: Add a "Keeping BUILD files tidy" subsection**

Add prose: `bazel run //:tidy` regenerates BUILD files (gazelle), refreshes the Python deps manifest, and formats everything in one step; the `Tidy Check` CI job fails a PR if anything is stale, with the same command as the fix; the pre-commit hook runs gazelle automatically. Cross-link from `formatting-linting.md`.

- [ ] **Step 2: Commit**

```bash
git add "{{ .ProjectSnake }}/docs/user-guide/development-workflow.md" "{{ .ProjectSnake }}/docs/user-guide/formatting-linting.md"
git commit -m "docs(tidy): document bazel run //:tidy and the Tidy Check gate"
```

---

### Task 6: Ship + dogfood

- [ ] **Step 1: Push branch + open PR to `platform-v2.0`**, body summarizing the three components + born-tidy result.
- [ ] **Step 2: Watch CI** (13-preset deliver/validation matrix) to green.
- [ ] **Step 3: GATED — merge** only on explicit user go-ahead (fan-out to 26 Starters). Watch the `deliver` run.
- [ ] **Step 4: Port to vitruvian-core** (add `rules_multirun` dep + `//:tidy` + `tidy-check.yaml` + pre-commit gazelle), run `bazel run //:tidy` to confirm born-tidy, open a PR, confirm `Tidy Check` is green. Fix any gap back in the template and re-deliver.

---

## Self-review

- **Spec coverage:** Component 1 → Task 1; Component 2 → Task 2; Component 3 → Task 3; born-tidy invariant → Task 4; docs → Task 5; rollout/dogfood → Task 6. All spec sections covered.
- **Placeholder scan:** Task 4 Step 2 deliberately cannot show generic code (the fix depends on the actual diff) — it is an exact procedure, not a placeholder. Task 5 prose is descriptive (docs). No TBD/TODO.
- **Consistency:** target label `//:tidy`, command `bazel run //:tidy`, and the `git status --porcelain` staleness check are used identically in Tasks 1, 2, 4, 6.
