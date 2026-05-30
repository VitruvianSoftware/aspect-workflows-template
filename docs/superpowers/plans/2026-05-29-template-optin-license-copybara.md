# Opt-in Template Features: License + Copybara — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add two opt-in Scaffold features to the monorepo generator — license-header enforcement and Copybara bidirectional sync — each defaulting OFF so unselected output is byte-identical to today.

**Architecture:** Pure Scaffold `questions`/`features`/`presets` gating (the pattern `lint`/`oci`/`proto` already use). License = `addlicense` provisioned via `rules_multitool` + dedicated `//tools/license:{add,check}` Bazel targets + a templated `LICENSE`. Copybara = the proven vitruvian-core engine, re-expressed as Bazel run-targets (`rules_go` `go_binary` for precheck/drift; `bazel run` wrappers around the pinned Copybara Docker image), plus optional Pulumi auth IaC.

**Tech Stack:** Scaffold (hay-kot) Go templates; Bazel + `rules_go` 0.59.0, `aspect_rules_lint` 2.0.0, `rules_multitool` 1.11.1, `bazel_env.bzl` 0.5.0; `addlicense@v1.2.0`; Copybara (`olivr/copybara@sha256:87e2e9089344e64693faebb2ee0ed33b8797358c0420b0fa98325ca611e98679`); Pulumi (Go, `pulumi-github`).

**Spec:** [`docs/superpowers/specs/2026-05-29-template-optin-license-copybara-design.md`](../specs/2026-05-29-template-optin-license-copybara-design.md).

**Conventions:** Repo root `/Users/james/Workspace/gh/application/vitruvian/aspect-workflows-template`, branch `feat/template-optin-license-copybara`. The generated-repo tree is the dir literally named `{{ .ProjectSnake }}/` — in shell, quote as `'{{ .ProjectSnake }}'`. Commit messages: conventional (`feat:`/`docs:`/`test:`), trailer `Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>`.

**Source of truth to port (read these in the OTHER repo `vitruvian-core`):**
- `…/vitruvian-core/tools/copybara/copy.bara.sky`
- `…/vitruvian-core/tools/copybara/conflict-precheck.sh`
- `…/vitruvian-core/.github/workflows/copybara-drift-check.yaml` (drift logic, inline)
- `…/vitruvian-core/.github/workflows/_copybara-{export,import}.yaml` (docker invocation)
- `…/vitruvian-core/infrastructure/pulumi/pkg/copybara_sync/sync.go`

---

## Critical sequencing & guarantees

1. **Byte-identical-when-off is load-bearing.** Every file touched in a *shared* (non-feature-gated) path — `scaffold.yaml`, `{{ .ProjectSnake }}/tools/tools.lock.json`, `{{ .ProjectSnake }}/tools/BUILD`, `{{ .ProjectSnake }}/MODULE.bazel` — must use `{{ if .Scaffold.<feature> }}` so generating with both features OFF reproduces today's output exactly. **Task C.1 is the guard test; run it after every shared-file edit.**
2. **Feature files live behind `features:` globs**, so they're excluded entirely when off.
3. **Copybara implies the Go toolchain** (its precheck/drift are `go_binary`s). Task B.1 makes `rules_go`/`go.mod`/gazelle available when `copybara` is on even if Go isn't a selected language.

---

## File structure

| Path (under `{{ .ProjectSnake }}/` unless noted) | New/Mod | Responsibility | Gate |
|---|---|---|---|
| `scaffold.yaml` (repo root) | Mod | questions, features globs, presets | n/a |
| `tools/tools.lock.json` | Mod | addlicense multitool entry (conditional) | `license` |
| `tools/BUILD` | Mod | `bazel_env` tools map += addlicense (conditional) | `license` |
| `tools/license/BUILD` | New | `:add` (run) + `:check` (test) targets | `license` |
| `tools/license/ignore.txt` | New | addlicense ignore list (data) | `license` |
| `LICENSE` | New | templated full text by `license_id` | `license` |
| `tools/workspace_status.sh`, `tools/scripts/kind-cluster.sh`, `githooks/check-config.sh`, `tools/format/prettier_wrapper.sh`, `eslint.config.mjs` | Mod | conditional license header | `license` (+ their own gates) |
| `MODULE.bazel` | Mod | rules_go gate widened to include copybara | `copybara`/go |
| `tools/copybara/copy.bara.sky` | New | engine config (COMPONENTS seeded) | `copybara` |
| `tools/copybara/conflict_precheck/{main.go,BUILD}` | New | precheck go_binary + go_test | `copybara` |
| `tools/copybara/drift_check/{main.go,BUILD}` | New | drift go_binary + go_test | `copybara` |
| `tools/copybara/sync/{main.go,BUILD}` | New | export/import engine wrapper (docker) | `copybara` |
| `.github/workflows/_copybara-{export,import}.yaml`, `copybara-drift-check.yaml` | New | thin `bazel run` callers | `copybara` |
| `tools/copybara/sync-to-monorepo.yaml` | New | reference dispatch snippet for standalone | `copybara` |
| `docs/copybara-bidi-sync.md` | New | admin runbook | `copybara` |
| `infrastructure/<pulumi>/.../copybara_sync/*` | New | optional auth IaC | `copybara`+`copybara_pulumi_auth` |

---

## Phase 0 — Scaffold foundation

### Task 0.1: Add the questions to `scaffold.yaml`

**Files:** Modify `scaffold.yaml` (repo root) — append to the `questions:` list (after the existing `oci` question, before `features:`).

- [ ] **Step 1: Add the questions.** Verify single-select / text / `when:` syntax against Scaffold v0.6.1 first (`scaffold --version`; check `https://hay-kot.github.io/scaffold/schema.json`). Add:

```yaml
  - name: license
    prompt:
      confirm: Add license-header enforcement?
      description: "Generates a LICENSE, an addlicense header check (bazel test //tools/license:check), and a `bazel run //tools/license:add` helper"
  - name: license_id
    when: "{{ .Scaffold.license }}"
    prompt:
      options:
        - Apache-2.0
        - MIT
        - BSD-3-Clause
        - MPL-2.0
      message: License (SPDX)
  - name: copyright_holder
    when: "{{ .Scaffold.license }}"
    prompt:
      message: Copyright holder (e.g. "Acme, Inc.")
  - name: copybara
    prompt:
      confirm: Add Copybara monorepo<->standalone bidirectional sync?
      description: "Bazel-run sync engine (export/import/drift/conflict-precheck) for publishing subtrees to standalone repos"
  - name: copybara_components
    when: "{{ .Scaffold.copybara }}"
    prompt:
      message: "Initial sync component subdirs, comma-separated (blank = none, add later via the runbook)"
  - name: copybara_pulumi_auth
    when: "{{ .Scaffold.copybara }}"
    prompt:
      confirm: "Also scaffold the GitHub-App + deploy-key Pulumi auth IaC?"
```

- [ ] **Step 2: Validate Scaffold parses it.**

Run: `scaffold new --preset=minimal --no-prompt /tmp/sc-foundation-test && echo OK`
Expected: `OK` (generation still succeeds with both new confirms defaulting off). Clean up: `rm -rf /tmp/sc-foundation-test`.

- [ ] **Step 3: Commit.**

```bash
git add scaffold.yaml
git commit -m "feat(scaffold): add license + copybara opt-in questions"
```

### Task 0.2: Add `features:` globs + `computed` helper

**Files:** Modify `scaffold.yaml` — append to `features:` list and `computed:` map.

- [ ] **Step 1: Add the feature gates.** Append to `features:`:

```yaml
  - value: "{{ .Scaffold.license }}"
    globs:
      - "*/LICENSE"
      - "*/tools/license/*"
  - value: "{{ .Scaffold.copybara }}"
    globs:
      - "*/tools/copybara/**"
      - "**/_copybara-export.yaml"
      - "**/_copybara-import.yaml"
      - "**/copybara-drift-check.yaml"
      - "*/docs/copybara-bidi-sync.md"
  - value: "{{ and .Scaffold.copybara .Scaffold.copybara_pulumi_auth }}"
    globs:
      - "*/infrastructure/**/copybara_sync/**"
```

Append to `computed:`:
```yaml
  go_or_copybara: "{{ or (has \"Go\" .Scaffold.langs) .Scaffold.copybara }}"
```

- [ ] **Step 2: Verify generation still works both ways.**

Run: `scaffold new --preset=minimal --no-prompt /tmp/sc-feat-off && echo OFF_OK && rm -rf /tmp/sc-feat-off`
Expected: `OFF_OK` (globs match nothing yet since files don't exist; harmless).

- [ ] **Step 3: Commit.**

```bash
git add scaffold.yaml
git commit -m "feat(scaffold): gate license + copybara files via features globs"
```

---

## Phase A — License feature

### Task A.1: Provision `addlicense` via multitool (conditional)

**Files:** Modify `{{ .ProjectSnake }}/tools/tools.lock.json`; modify `{{ .ProjectSnake }}/tools/BUILD`.

`addlicense` v1.2.0 ships tarballs: `addlicense_v1.2.0_<OS>_<ARCH>.tar.gz` for `Linux/macOS` × `x86_64/arm64`, with `checksums.txt` (fetch SHAs at impl time with `gh release view v1.2.0 -R google/addlicense`). `rules_multitool` schema supports `kind:"archive"` + `file:"addlicense"`.

- [ ] **Step 1: Add a conditional addlicense entry to `tools.lock.json`.** Scaffold templates this file, so wrap in `{{ if .Scaffold.license }}`. Place the entry so JSON stays valid when off (no trailing comma) and on. Pattern (insert as the LAST top-level key; note the leading comma is inside the conditional):

```json
  "buf": { "binaries": [ ... ] }{{ if .Scaffold.license }},
  "addlicense": {
    "binaries": [
      { "kind": "archive", "file": "addlicense", "os": "linux",  "cpu": "x86_64", "url": "https://github.com/google/addlicense/releases/download/v1.2.0/addlicense_v1.2.0_Linux_x86_64.tar.gz",  "sha256": "<FILL>" },
      { "kind": "archive", "file": "addlicense", "os": "linux",  "cpu": "arm64",  "url": "https://github.com/google/addlicense/releases/download/v1.2.0/addlicense_v1.2.0_Linux_arm64.tar.gz",   "sha256": "<FILL>" },
      { "kind": "archive", "file": "addlicense", "os": "macos",  "cpu": "x86_64", "url": "https://github.com/google/addlicense/releases/download/v1.2.0/addlicense_v1.2.0_macOS_x86_64.tar.gz",  "sha256": "<FILL>" },
      { "kind": "archive", "file": "addlicense", "os": "macos",  "cpu": "arm64",  "url": "https://github.com/google/addlicense/releases/download/v1.2.0/addlicense_v1.2.0_macOS_arm64.tar.gz",   "sha256": "<FILL>" }
    ]
  }{{ end }}
```
(Replace `<FILL>` with real sha256 from `checksums.txt`. Adjust the preceding key/brace to match the actual last entry in the file.)

- [ ] **Step 2: Add addlicense to the `bazel_env` tools map (conditional).** In `{{ .ProjectSnake }}/tools/BUILD`, inside the `bazel_env(... tools = {...})`, add a conditional line:

```python
    tools = {
        # ... existing entries ...
        {{- if .Scaffold.license }}
        "addlicense": "@multitool//tools/addlicense",
        {{- end }}
    } | MULTITOOL_TOOLS | ...
```
(Match the file's exact existing structure; if `tools` is built via `|` merges, add a conditional `{"addlicense": ...}` dict to the merge chain instead.)

- [ ] **Step 3: Verify byte-identical when off + valid when on.**

```bash
scaffold new --preset=minimal --no-prompt /tmp/sc-lic-off
python3 -c "import json;json.load(open('/tmp/sc-lic-off/tools/tools.lock.json'))" && echo OFF_VALID
# diff tools.lock.json against current generated baseline (see Task C.1 for the baseline harness)
scaffold new --preset=go --no-prompt --snapshot=license=true /tmp/sc-lic-on 2>/dev/null || true
# (exact non-interactive feature-set flag verified in Task C.1; here just assert JSON validity when present)
rm -rf /tmp/sc-lic-off /tmp/sc-lic-on
```
Expected: `OFF_VALID`; when on, `tools.lock.json` is valid JSON containing `addlicense`.

- [ ] **Step 4: Commit.**

```bash
git add "{{ .ProjectSnake }}/tools/tools.lock.json" "{{ .ProjectSnake }}/tools/BUILD"
git commit -m "feat(license): provision addlicense via multitool (gated)"
```

### Task A.2: License `add`/`check` Bazel targets + ignore list

**Files:** Create `{{ .ProjectSnake }}/tools/license/BUILD`, `{{ .ProjectSnake }}/tools/license/ignore.txt`.

- [ ] **Step 1: Create the ignore list** `{{ .ProjectSnake }}/tools/license/ignore.txt` (one glob per line; comprehensive, harmless across languages):

```
**/BUILD
**/BUILD.bazel
**/*.lock
**/*.lock.json
**/pnpm-lock.yaml
**/package-lock.json
**/Cargo.lock
**/MODULE.bazel.lock
**/gazelle_python.yaml
**/*-baseline.xml
**/.release-please-manifest.json
bazel-*/**
**/node_modules/**
**/*.venv/**
.git/**
```

- [ ] **Step 2: Create `tools/license/BUILD`** with a runnable `add` target and a `check` test. Both call the bazel_env-provided `addlicense`. Implementer chooses `native_binary` + `sh_test`-free wiring; the simplest robust form uses `rules_multitool`'s binary alias as the executable and a Bazel `*_test`. Target shape (templated copyright/license from scaffold vars):

```python
# Targets to add (fix) and check MIT/Apache/etc. headers. addlicense comes from
# @multitool//tools/addlicense (provisioned in tools.lock.json, license-gated).
load("@bazel_skylib//rules:native_binary.bzl", "native_binary", "native_test")

_ADDLICENSE_ARGS = [
    "-c", "{{ .Scaffold.copyright_holder }}",
    "-l", "{{ if eq .Scaffold.license_id \"Apache-2.0\" }}apache{{ else if eq .Scaffold.license_id \"MIT\" }}mit{{ else if eq .Scaffold.license_id \"BSD-3-Clause\" }}bsd{{ else }}mpl{{ end }}",
]

native_binary(
    name = "add",
    src = "@multitool//tools/addlicense",
    args = _ADDLICENSE_ARGS + ["-ignore", "@$(location ignore.txt)", "."],  # see note
    data = ["ignore.txt"],
    visibility = ["//visibility:public"],
)

native_test(
    name = "check",
    src = "@multitool//tools/addlicense",
    args = _ADDLICENSE_ARGS + ["-check"] + IGNORES + ["."],
    data = ["ignore.txt"],
)
```
**Implementation note:** addlicense takes repeated `-ignore <glob>` flags, not a file. The implementer should either (a) expand `ignore.txt` into repeated `-ignore` args in the BUILD (read the file at load time via a `.bzl` helper), or (b) generate the `IGNORES` list inline in the BUILD as a Starlark list literal and drop `ignore.txt`. Prefer (b) for simplicity — make `IGNORES = ["-ignore", "**/BUILD", "-ignore", "**/BUILD.bazel", ...]` a Starlark list in the BUILD and delete `ignore.txt` from Step 1. The check must run from the workspace root; if `native_test` sandboxing prevents whole-tree scanning, fall back to a `genrule`-free `sh_test` that `cd $BUILD_WORKSPACE_DIRECTORY` — but try `native_test` first.

- [ ] **Step 3: Smoke-test the targets build (in a generated repo — see Task C.2 harness).** Deferred to C.2 (needs a generated repo). Here just `buildifier`-lint the BUILD locally:

Run: `bazel run //{{ .ProjectSnake }}-noop 2>/dev/null; buildifier --mode=check "{{ .ProjectSnake }}/tools/license/BUILD" && echo BUILDIFIER_OK` (or `~/go/bin/buildifier`).
Expected: `BUILDIFIER_OK`.

- [ ] **Step 4: Commit.**

```bash
git add "{{ .ProjectSnake }}/tools/license/"
git commit -m "feat(license): add bazel run :add + bazel test :check targets"
```

### Task A.3: Templated `LICENSE`

**Files:** Create `{{ .ProjectSnake }}/LICENSE`.

- [ ] **Step 1: Create `LICENSE`** as a templated file switching on `license_id`. Embed the four standard texts. Structure:

```
{{- if eq .Scaffold.license_id "MIT" -}}
MIT License

Copyright (c) {{ now.Format "2006" }} {{ .Scaffold.copyright_holder }}

Permission is hereby granted, free of charge, ...
{{- else if eq .Scaffold.license_id "Apache-2.0" -}}
                                 Apache License
                           Version 2.0, January 2004
... (full Apache 2.0 text) ...
Copyright {{ now.Format "2006" }} {{ .Scaffold.copyright_holder }}
{{- else if eq .Scaffold.license_id "BSD-3-Clause" -}}
BSD 3-Clause License

Copyright (c) {{ now.Format "2006" }}, {{ .Scaffold.copyright_holder }}
... (full text) ...
{{- else -}}
Mozilla Public License Version 2.0
... (full MPL-2.0 text) ...
{{- end -}}
```
(Use the canonical SPDX full texts. Verify Scaffold's date helper — `now.Format` or `{{ .Computed.year }}`; if no date helper, drop the year or compute it in a `computed:` entry.)

- [ ] **Step 2: Verify each variant renders.** Generate once per license id (non-interactive snapshot mechanism confirmed in C.1) and assert the LICENSE first line matches. Defer the full matrix to Task C.2; here generate one (MIT) and eyeball.

- [ ] **Step 3: Commit.**

```bash
git add "{{ .ProjectSnake }}/LICENSE"
git commit -m "feat(license): templated LICENSE by SPDX id"
```

### Task A.4: Born-green headers on existing source files

**Files:** Modify the committed source files that the check would scan (only these exist): `{{ .ProjectSnake }}/tools/workspace_status.sh`, `{{ .ProjectSnake }}/tools/scripts/kind-cluster.sh`, `{{ .ProjectSnake }}/githooks/check-config.sh`, `{{ .ProjectSnake }}/tools/format/prettier_wrapper.sh`, `{{ .ProjectSnake }}/eslint.config.mjs`.

- [ ] **Step 1: Define a reusable header partial.** Scaffold supports `{{ template }}` partials or inline conditionals. Simplest: an inline conditional header block per comment-style. For shell files (`#` after the shebang):

```bash
#!/usr/bin/env bash
{{- if .Scaffold.license }}
# Copyright (c) {{ now.Format "2006" }} {{ .Scaffold.copyright_holder }}
# SPDX-License-Identifier: {{ .Scaffold.license_id }}
{{- end }}
# ... existing body ...
```
For `.mjs` (`//`):
```javascript
{{- if .Scaffold.license }}
// Copyright (c) {{ now.Format "2006" }} {{ .Scaffold.copyright_holder }}
// SPDX-License-Identifier: {{ .Scaffold.license_id }}
{{- end }}
// ... existing body ...
```
**Note:** addlicense's check accepts an SPDX-ID header OR its generated header; keep this partial's format consistent with what `addlicense -check -l <x>` accepts (verify: run `addlicense -check` over a sample headered file during C.2; if it rejects the SPDX-only form, switch the partial to addlicense's exact template output for that license).

- [ ] **Step 2: Apply the partial to all five files** (after shebang for shell; top for mjs).

- [ ] **Step 3: Verify off = unchanged.** Covered by Task C.1 guard.

- [ ] **Step 4: Commit.**

```bash
git add "{{ .ProjectSnake }}/tools/workspace_status.sh" "{{ .ProjectSnake }}/tools/scripts/kind-cluster.sh" "{{ .ProjectSnake }}/githooks/check-config.sh" "{{ .ProjectSnake }}/tools/format/prettier_wrapper.sh" "{{ .ProjectSnake }}/eslint.config.mjs"
git commit -m "feat(license): conditional headers on generated source (born-green)"
```

---

## Phase B — Copybara feature

### Task B.1: Make the Go toolchain available when `copybara` is on

**Files:** Modify `{{ .ProjectSnake }}/MODULE.bazel`, `{{ .ProjectSnake }}/.aspect/gazelle/*` if gated, and the `features:`/go-glob in `scaffold.yaml` (the `*/go.mod`, `*/go.sum`, `*/tools/tools.go` glob).

The copybara Go binaries need `rules_go` + a `go.mod`. Today both are gated on `{{ .Computed.go }}`. Widen to `{{ .Computed.go_or_copybara }}` (added in Task 0.2).

- [ ] **Step 1: In `scaffold.yaml`**, change the Go `features` entry value from `"{{ .Computed.go }}"` to `"{{ .Computed.go_or_copybara }}"` for the globs that the copybara Go targets need: `*/go.mod`, `*/go.sum`, `*/tools/tools.go`, and `*/.aspect/gazelle/*go*` (the go gazelle config). Leave language-demo globs Go-only if any.

- [ ] **Step 2: In `{{ .ProjectSnake }}/MODULE.bazel`**, widen the `rules_go`/`gazelle`/`go_sdk`/`go_deps` `{{- if .Computed.go }}` guards to `{{- if .Computed.go_or_copybara }}`.

- [ ] **Step 3: Verify generation matrices.**

```bash
# copybara on, Go NOT selected -> go.mod + rules_go present
scaffold new --preset=minimal --no-prompt <copybara=true> /tmp/sc-cb-nogo  # exact flag per C.1
test -f /tmp/sc-cb-nogo/go.mod && grep -q rules_go /tmp/sc-cb-nogo/MODULE.bazel && echo CB_GO_OK
# both off -> NO go.mod (unchanged)
scaffold new --preset=minimal --no-prompt /tmp/sc-min && test ! -f /tmp/sc-min/go.mod && echo MIN_OK
rm -rf /tmp/sc-cb-nogo /tmp/sc-min
```
Expected: `CB_GO_OK` and `MIN_OK`.

- [ ] **Step 4: Commit.**

```bash
git add scaffold.yaml "{{ .ProjectSnake }}/MODULE.bazel"
git commit -m "feat(copybara): pull in Go toolchain when copybara enabled"
```

### Task B.2: Templated `copy.bara.sky`

**Files:** Create `{{ .ProjectSnake }}/tools/copybara/copy.bara.sky`.

Port `…/vitruvian-core/tools/copybara/copy.bara.sky`. Keep the structure verbatim (COMPONENTS list of dicts; `[_define_component(_c) for _c in COMPONENTS]`; `_make_skip_guard`; `_monorepo_only`; export/import `core.workflow`s; `ITERATIVE`; `experimental_custom_rev_id`; label format `[A-Z][A-Z_0-9]{1,30}_REV_ID`). Parameterize the constants from scaffold vars.

- [ ] **Step 1: Create the file**, replacing the hardcoded org/repo with template vars and seeding `COMPONENTS`:

```python
# (header partial — Task A.3 style — if .Scaffold.license)
MONOREPO = "https://github.com/{{ .Scaffold.github_org | default "YOUR_ORG" }}/{{ .ProjectKebab }}.git"
STANDALONE_SSH_PREFIX = "git@github.com:{{ .Scaffold.github_org | default "YOUR_ORG" }}/"
AUTHOR = "{{ .ProjectName }} Sync <sync@example.com>"
MONOREPO_REV_ID = "MONOREPO_REV_ID"

COMPONENTS = [
{{- range $c := splitList "," .Scaffold.copybara_components }}{{ if $c }}
    {"name": "{{ trim $c }}", "standalone_rev_id": "{{ upper (replace "-" "_" (trim $c)) }}_REV_ID", "standalone_only": [".github/workflows/sync-to-monorepo.yaml"]},
{{- end }}{{- end }}
]
# ... rest verbatim from vitruvian-core (the _define_component fn, skip guards, globs, core.move) ...
```
(If Scaffold lacks `splitList`/`trim`/`replace`/`upper` helpers, verify its template funcs and adjust; worst case seed an empty `COMPONENTS = []` and document adding entries in the runbook.) Confirm whether a `github_org` question is needed — if so add it to Task 0.1; otherwise use a documented placeholder the operator edits.

- [ ] **Step 2: Validate Starlark offline.** Copy the proven offline check: `docker run --rm --network none -v "$PWD":/src -w /src olivr/copybara@sha256:87e2e90… copybara info /src/{{ .ProjectSnake }}/tools/copybara/copy.bara.sky` against a *generated* repo (defer to C.2) — here just confirm no `for` at top level and balanced parens by eye + `python3 -c "compile(open('…').read(),'x','exec')" || true` (Starlark≈Python syntactically for a smoke check).

- [ ] **Step 3: Commit.**

```bash
git add "{{ .ProjectSnake }}/tools/copybara/copy.bara.sky"
git commit -m "feat(copybara): templated copy.bara.sky (COMPONENTS seeded)"
```

### Task B.3: `conflict_precheck` Go binary (TDD)

**Files:** Create `{{ .ProjectSnake }}/tools/copybara/conflict_precheck/main.go`, `…/conflict_precheck/main_test.go`, `…/conflict_precheck/BUILD`.

Port `conflict-precheck.sh` (algorithm below) to Go, shelling out to `git`. CLI: `conflict_precheck <export|import> <component> <monorepo_dir> <standalone_dir> [standalone_only...]`.

**Algorithm (faithful port):**
- `EXPORT_LABEL="MONOREPO_REV_ID"`; `IMPORT_LABEL = upper(replace(component,"-","_")) + "_REV_ID"`.
- `latestRev(repo,label)`: `git -C repo log -1 --grep="<label>:" --format=%B`, extract via regex `(?m)^<label>: ([0-9a-f]{7,40})`, first match (empty if none).
- `genuineCommits(repo,range,peerLabel,pathspecs...)`: `git -C repo log <range> --no-merges --invert-grep --grep="<peerLabel>:" --format='  %h %s' -- <pathspecs...>`.
- `monorepoOnlyExcludes(c)` = `[":(exclude,glob)c/**/BUILD", ":(exclude,glob)c/**/BUILD.bazel", ":(exclude)c/BUILD", ":(exclude)c/BUILD.bazel"]`.
- **export:** `base=latestRev(MONO,IMPORT_LABEL)`; empty→print "no import baseline", exit 0. `git -C STD cat-file -e base^{commit}`; fail→`::warning::`, exit 0. `scan=["."]+[":(exclude)"+f for f in standalone_only]`; `pending=genuineCommits(STD, base+"..HEAD", EXPORT_LABEL, scan...)`.
- **import:** `base=latestRev(STD,EXPORT_LABEL)`; empty→exit 0. `git -C MONO cat-file -e base^{commit}`; fail→`::warning::`, exit 0. `pending=genuineCommits(MONO, base+"..HEAD", IMPORT_LABEL, "c/", monorepoOnlyExcludes(c)...)`.
- unknown direction→exit 2. If `pending` non-empty→print `::error title=Copybara conflict pre-check::…` + commits, exit 1; else print OK, exit 0.

- [ ] **Step 1: Write failing tests** `main_test.go` using temp git repos (`t.TempDir()`, `git init`, commit with/without label-bearing messages). Cover: export with no baseline→exit 0; export with a genuine standalone commit→exit 1; import clean→exit 0; rev-id extraction regex; label derivation (`mcp-slack`→`MCP_SLACK_REV_ID`).

```go
func TestImportLabelDerivation(t *testing.T) {
    if got := importLabel("mcp-slack"); got != "MCP_SLACK_REV_ID" {
        t.Fatalf("got %q", got)
    }
}
// + TestLatestRevExtractsSha, TestExportNoBaselineExitsZero, TestExportGenuineChangeFails (temp repos)
```

- [ ] **Step 2: Run, verify fail.** `cd <generated repo>` then `bazel test //tools/copybara/conflict_precheck:conflict_precheck_test` (or `go test ./tools/copybara/conflict_precheck/`). Expected: FAIL (undefined funcs).
- [ ] **Step 3: Implement `main.go`** per the algorithm (use `os/exec` for git, `regexp`, `fmt` to stderr for `::error::`/`::warning::`, `os.Exit`).
- [ ] **Step 4: Run, verify pass.** Expected: PASS.
- [ ] **Step 5: BUILD** (`go_binary` + `go_test`, gazelle-generatable):
```python
load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
# (gazelle will generate; ensure go_binary name = "conflict_precheck", visibility public)
```
- [ ] **Step 6: Commit.** `git add "{{ .ProjectSnake }}/tools/copybara/conflict_precheck/" && git commit -m "feat(copybara): conflict_precheck go_binary (port of bash precheck)"`

### Task B.4: `drift_check` Go binary (TDD)

**Files:** Create `{{ .ProjectSnake }}/tools/copybara/drift_check/{main.go,main_test.go,BUILD}`.

Port the drift logic (inline in `…/vitruvian-core/.github/workflows/copybara-drift-check.yaml`). CLI: `drift_check <monorepo_workspace_dir> <component...>` (SSH key + org from env/flags).

**Algorithm:** for each component: resolve `<PREFIX>_SYNC_SSH_KEY` env (PREFIX = upper/replace-`-`_`); empty→`::warning::` skip. Write key `~/.ssh/id_rsa` (chmod 600, strip trailing newline), `ssh-keyscan github.com`. `git clone --no-tags git@github.com:<org>/<c>.git <tmp>`; fail→`::warning::` skip. Seed gate: `git -C <tmp> log --grep=MONOREPO_REV_ID -1 --format=%H` empty→skip "not seeded". Compare: recursive dir diff of `<workspace>/<c>` vs `<tmp>` excluding basenames `.git`,`BUILD`,`BUILD.bazel`,`package-lock.json`,`sync-to-monorepo.yaml`; differences→`::error title=Copybara sync drift (<c>)::` + the diff, set drift. After loop: exit 1 if any drift.

- [ ] **Step 1: Failing tests** — `TestNoDriftIdenticalTrees`, `TestDriftDetectedOnContentDiff`, `TestExcludedBasenamesIgnored` using two temp dirs (skip the clone/SSH path by factoring the comparison into a pure `compareDirs(a, b, excludes)` func and testing that directly).
- [ ] **Step 2: Run, verify fail.**
- [ ] **Step 3: Implement** — factor `compareDirs` (pure, tested) + the clone/seed-gate orchestration (`os/exec`).
- [ ] **Step 4: Run, verify pass.**
- [ ] **Step 5: BUILD** (`go_binary` `drift_check` + `go_test`).
- [ ] **Step 6: Commit.** `git commit -m "feat(copybara): drift_check go_binary (port of inline drift logic)"`

### Task B.5: `sync` engine wrapper Go binary (export/import)

**Files:** Create `{{ .ProjectSnake }}/tools/copybara/sync/{main.go,main_test.go,BUILD}`.

Wrap the pinned Copybara Docker image. CLI: `sync <export|import> <component> [--options=...] [--force]`. Replicates `_copybara-{export,import}.yaml`: build `WF = "{export|import}_" + replace(component,"-","_")`; run the exact `docker run … olivr/copybara@sha256:87e2e9089344e64693faebb2ee0ed33b8797358c0420b0fa98325ca611e98679 copybara` with the volume/env wiring; treat exit 0 and 4 (NO_OP) as success; **import** retries up to 3× with `attempt*5`s backoff on other non-zero.

- [ ] **Step 1: Failing tests** — `TestWorkflowName` (`import`,`mcp-slack`→`import_mcp_slack`), `TestExitCodeClassification` (0,4→success; 1→fail/retry), `TestRetryStopsAtMax` (inject a fake runner func). Factor the docker call behind an interface so tests don't need Docker.
- [ ] **Step 2: Run, verify fail.**
- [ ] **Step 3: Implement** (constants for image digest + volume/env flags; `runner` interface; retry loop for import).
- [ ] **Step 4: Run, verify pass.**
- [ ] **Step 5: BUILD** (`go_binary` `export` + `import` via two `go_binary`s sharing the lib, OR one `sync` binary; the workflow callers in B.6 must match the chosen target names — use `//tools/copybara/sync:export` and `:import` as `go_binary` aliases over the shared `go_library`).
- [ ] **Step 6: Commit.** `git commit -m "feat(copybara): sync engine wrapper (bazel run -> pinned copybara image)"`

### Task B.6: Thin workflow callers (`bazel run`)

**Files:** Create `{{ .ProjectSnake }}/.github/workflows/_copybara-export.yaml`, `_copybara-import.yaml`, `copybara-drift-check.yaml`.

These replace vitruvian-core's docker-in-workflow with `bazel run` calls. Each sets up SSH/git-credentials (port the auth-setup step from `_copybara-export.yaml`), then `bazel run //tools/copybara/sync:export -- <component> ...` (export) / `:import` (import) / `bazel run //tools/copybara/drift_check -- ...`. Conflict-precheck step becomes `bazel run //tools/copybara/conflict_precheck -- ...`.

- [ ] **Step 1: Create `_copybara-export.yaml`** (reusable `workflow_call`, inputs `component`/`standalone_only`/`copybara_options`, secret `sync_ssh_key`, `permissions: contents: read`): auth-setup step (verbatim port), pre-check step `bazel run //tools/copybara/conflict_precheck -- export "${{ inputs.component }}" "$GITHUB_WORKSPACE" "$PEER" $STANDALONE_ONLY`, sync step `bazel run //tools/copybara/sync:export -- "${{ inputs.component }}" --options="$CB_OPTS"`. (Include the license header partial if `.Scaffold.license`.)
- [ ] **Step 2: Create `_copybara-import.yaml`** (mirror; `contents: write`; the retry lives in the Go binary now, so the workflow is a single `bazel run :import`).
- [ ] **Step 3: Create `copybara-drift-check.yaml`** (`schedule` + `workflow_dispatch` + `workflow_run`; `bazel run //tools/copybara/drift_check -- "$GITHUB_WORKSPACE" <components…>`; per-component SSH key env).
- [ ] **Step 4: actionlint.** `~/go/bin/actionlint <each generated file>` (validate against a *generated* repo in C.2; here actionlint the templated files after stripping `{{ }}` is unreliable — defer real lint to C.2, but eyeball-validate YAML structure).
- [ ] **Step 5: Commit.** `git commit -m "feat(copybara): thin bazel-run workflow callers"`

> **Note — per-component caller workflows are NOT generated** (Approach 1). The runbook (B.7) documents adding `copybara-export-<c>.yaml`/`copybara-import-<c>.yaml` per component post-generation.

### Task B.7: Runbook + standalone dispatch snippet

**Files:** Create `{{ .ProjectSnake }}/docs/copybara-bidi-sync.md`, `{{ .ProjectSnake }}/tools/copybara/sync-to-monorepo.yaml`.

- [ ] **Step 1: Port the runbook** from `…/vitruvian-core/docs/copybara-bidi-sync.md`, generalized (template the org/repo; keep §8b seed recipe, §8f onboard-a-component flow incl. creating per-component caller workflows, §9 troubleshooting, §12 if license on). Replace bash-script references with the `bazel run //tools/copybara:*` targets.
- [ ] **Step 2: Create the standalone dispatch snippet** `sync-to-monorepo.yaml` (the file the operator installs in each *standalone* repo to fire `repository_dispatch <component>-import`), as a documented reference (it is NOT a monorepo workflow).
- [ ] **Step 3: Commit.** `git commit -m "docs(copybara): runbook + standalone dispatch snippet"`

### Task B.8: Optional Pulumi auth IaC (gated)

**Files:** Create `{{ .ProjectSnake }}/infrastructure/<pulumi pkg>/copybara_sync/sync.go` (match the template's existing infrastructure/ layout; if the template has no Pulumi pkg dir, create `infrastructure/pulumi/pkg/copybara_sync/sync.go` and a minimal `BUILD`/go wiring, gated).

Port `…/vitruvian-core/infrastructure/pulumi/pkg/copybara_sync/sync.go`, parameterized: `syncedProjects` seeded from `copybara_components`; org from the template var. Resources per project: `tls.NewPrivateKey` (ED25519), `github.NewRepositoryDeployKey` (write) on standalone, `github.NewActionsSecret` `<PREFIX>_SYNC_SSH_KEY` in monorepo, two dispatch-app secrets in standalone; plus monorepo-wide `SYNC_APP_ID`/`SYNC_APP_PRIVATE_KEY` (Dependabot + Actions).

- [ ] **Step 1: Create the gated Pulumi file** (everything under `infrastructure/**/copybara_sync/**`, gated by `copybara`+`copybara_pulumi_auth` per Task 0.2). Seed `syncedProjects` from `copybara_components`.
- [ ] **Step 2: `gofmt`/`buildifier` lint** the file(s).
- [ ] **Step 3: Commit.** `git commit -m "feat(copybara): optional pulumi auth IaC (gated)"`

---

## Phase C — Presets, validation, final review

### Task C.1: Byte-identical-when-off guard + non-interactive feature flags

**Files:** none (validation); may add a tiny `test.sh`-style check if the repo wants one (it has `test.sh`).

- [ ] **Step 1: Determine the non-interactive way to set feature answers.** Check `scaffold new --help` for `--snapshot`/`--preset`/answer flags. Document the exact invocation to generate with `license`/`copybara` on/off non-interactively (used by all earlier verify steps).
- [ ] **Step 2: Baseline diff.** Generate `minimal` (and `kitchen-sink`) with both features OFF from the feature branch and from `platform-v2.0`; diff the trees — must be **identical**.

```bash
git stash || true
for ref in platform-v2.0 feat/template-optin-license-copybara; do
  git checkout "$ref" 2>/dev/null
  scaffold new --preset=kitchen-sink --no-prompt "/tmp/ks-$ref"
done
git checkout feat/template-optin-license-copybara
diff -r "/tmp/ks-platform-v2.0" "/tmp/ks-feat-template-optin-license-copybara" && echo "BYTE_IDENTICAL_OFF" || echo "DRIFT — investigate"
rm -rf /tmp/ks-*
```
Expected: `BYTE_IDENTICAL_OFF`. (If drift, a shared-file conditional leaked — fix before proceeding.)

- [ ] **Step 3: Commit** (if any guard script added).

### Task C.2: Local generate-and-build validation (primary gate)

Because `deliver.yaml` force-pushes to per-preset repos (missing) and the org's CI runners are currently unavailable (spec Risk 6), validation is **local**.

- [ ] **Step 1: License on — generate + build + check.**
```bash
scaffold new --preset=go --no-prompt <license=true,license_id=MIT,copyright_holder="Example Org"> /tmp/gen-lic
cd /tmp/gen-lic
test "$(head -1 LICENSE)" = "MIT License" && echo LICENSE_OK
bazel run //tools/license:add        # headers the tree
bazel test //tools/license:check     # must pass after add
cd - && rm -rf /tmp/gen-lic
```
Expected: `LICENSE_OK`, `:add` succeeds, `:check` is green.

- [ ] **Step 2: Copybara on (no components) — build + unit tests + drift no-op.**
```bash
scaffold new --preset=go --no-prompt <copybara=true,copybara_components=""> /tmp/gen-cb
cd /tmp/gen-cb
bazel test //tools/copybara/...      # go_tests pass
bazel run //tools/copybara/drift_check -- "$PWD"   # no components -> green no-op
~/go/bin/actionlint .github/workflows/_copybara-export.yaml .github/workflows/_copybara-import.yaml .github/workflows/copybara-drift-check.yaml
cd - && rm -rf /tmp/gen-cb
```
Expected: copybara go_tests pass; drift no-op exits 0; actionlint clean.

- [ ] **Step 3: Copybara on, Go NOT selected — toolchain still present + builds.**
```bash
scaffold new --preset=minimal --no-prompt <copybara=true> /tmp/gen-cb-min
cd /tmp/gen-cb-min && test -f go.mod && bazel test //tools/copybara/... && echo CB_MIN_OK
cd - && rm -rf /tmp/gen-cb-min
```
Expected: `CB_MIN_OK`.

- [ ] **Step 4: copy.bara.sky offline validation** in the copybara-on repo: `docker run --rm --network none -v /tmp/gen-cb:/src olivr/copybara@sha256:87e2e90… copybara info /src/tools/copybara/copy.bara.sky` (or `validate`) — confirm it parses (no top-level `for`, labels valid).

### Task C.3: Presets (local only) + final self-review

**Files:** Modify `scaffold.yaml` (`presets:` block). **Do NOT** add to `deliver.yaml`'s push matrix (those targets don't exist + CI is blocked).

- [ ] **Step 1: Add local presets** so `scaffold new --preset=…` is easy:
```yaml
  license-go:
    langs: ['Go']
    lint: true
    license: true
    license_id: Apache-2.0
    copyright_holder: "Example Org"
  copybara-go:
    langs: ['Go']
    lint: true
    copybara: true
    copybara_components: ""
    copybara_pulumi_auth: false
```
- [ ] **Step 2: Generate from each preset** and rerun the C.2 checks. Expected: green.
- [ ] **Step 3: Commit.** `git commit -m "feat(scaffold): local presets for license + copybara"`
- [ ] **Step 4: Final code review** — dispatch a code-review subagent over the whole branch diff (`git diff platform-v2.0...HEAD`). Address findings.
- [ ] **Step 5: Finish** via superpowers:finishing-a-development-branch (PR to `platform-v2.0`; note the `deliver`-matrix CI validation is deferred until org Actions runners return — task #77).

---

## Notes on iteration
- The biggest risk is **leaking a non-conditional change into a shared file** → run Task C.1's byte-identical guard after every Phase-0/A.1/B.1 edit.
- `addlicense -check` header-format acceptance (Task A.4 note): if the SPDX-only partial isn't accepted, switch the partial to addlicense's exact generated header for the chosen license.
- Scaffold template-func availability (`splitList`, `trim`, `upper`, `replace`, `now`): verify against v0.6.1; adjust seeding/date logic if missing.
- Copybara `github_org` is unknown at generation — either add a `github_org` question (Task 0.1) or ship a documented placeholder the operator edits post-gen.
