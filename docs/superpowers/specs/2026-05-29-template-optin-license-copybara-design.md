# Opt-in Template Features: License Enforcement + Copybara Bidi-Sync — Design

**Status:** Approved (brainstorm complete, user waived re-review)
**Date:** 2026-05-29
**Repo:** `aspect-workflows-template` (the Scaffold-based monorepo generator)
**Branch:** `feat/template-optin-license-copybara` (off `platform-v2.0`)

## Goal

Add two **opt-in** features to the monorepo generator, each a Scaffold toggle that
**defaults OFF** so current generated output is byte-for-byte unchanged when unselected:

1. **License enforcement** — generated repos get a license-header check + an `addlicense`
   add-mode capability + a matching top-level `LICENSE`, parameterized by the user's chosen
   SPDX license and copyright holder.
2. **Copybara bidirectional sync** — generated repos get the generalized monorepo↔standalone
   sync engine (proven in `vitruvian-core`), exposed as **Bazel run-targets** rather than loose
   shell scripts, with an optional Pulumi auth-provisioning layer.

## Architecture

Both features ride the template's **existing Scaffold optional-feature mechanism** — the same one
that already gates `lint`, `oci`, `proto`, `stamp`, and `backstage`:

- **`questions:`** entries in `scaffold.yaml` → become `.Scaffold.<name>`. `confirm:` = boolean
  (defaults off), `options:` (no `multi`) = single-select, plain `message:` = free text, `multi: true` =
  multi-select. Conditional sub-questions use `when:`.
- **`features:`** entries `{value: "<go-template bool>", globs: [...]}` — when the value is truthy the
  matching files are *included* in the generated repo; otherwise *excluded*.
- **Inline `{{ if .Scaffold.<name> }}…{{ end }}`** in shared files for small conditional fragments.
- **`presets:`** named answer bundles consumed by the `deliver` CI matrix (and Backstage).

> **Scaffold-syntax verification (planning):** confirm the exact syntax for single-select
> (`options:` without `multi`), free-text prompts, and conditional `when:` against the pinned Scaffold
> version (`metadata.minimum_version: v0.6.1`). Adjust question definitions if the syntax differs.

**Tech stack:** Scaffold (hay-kot) Go templates; Bazel + `rules_go`, `rules_lint`, `rules_multitool`
(`bazel_env`); `addlicense` (`github.com/google/addlicense@v1.2.0`); Copybara (pinned Docker image);
Pulumi (Go, `pulumi-github`).

---

## Feature A — License enforcement

### Questions (added to `scaffold.yaml`)
- `license` — `confirm`, **default off**. Master toggle.
- `license_id` — single-select, asked only `when` `license` is true. Curated set that maps cleanly to
  `addlicense -l`: **Apache-2.0** (default) → `apache`, **MIT** → `mit`, **BSD-3-Clause** → `bsd`,
  **MPL-2.0** → `mpl`. (More can be added later; curated so header + LICENSE stay consistent.)
- `copyright_holder` — free-text, asked `when` `license` is true (e.g. "Acme, Inc.").

### Behaviour when ON (all files gated by `{{ .Scaffold.license }}`)
1. **`addlicense` as a Bazel-managed tool** — provisioned the way the template already provisions
   tooling: prefer `rules_multitool` (`tools/tools.lock.json` + `tools/downloader.cfg`) exposed on PATH
   via `bazel_env`; if no prebuilt release binaries exist for multitool, define a `rules_go` `go_binary`
   instead. **No `go install`, no hand-written `tools/*.sh`.**
2. **Add headers** — wired through Bazel: integrate `addlicense` into the existing `aspect format` /
   `rules_lint` formatter flow (`tools/format/`) so headers are added alongside the other formatters,
   or via a dedicated `bazel run //tools/license:add` target. No bespoke script.
3. **Check / enforcement** — rides the generated repo's existing `bazel test //...` + `aspect` flow (a
   Bazel check target or the rules_lint format-check), so it is enforced wherever the repo already runs
   Bazel CI. The check ignores generated/tool-managed files, mirroring `vitruvian-core`'s list:
   `**/BUILD`, `**/BUILD.bazel`, lockfiles (`pnpm-lock.yaml`, `**/package-lock.json`, `**/Cargo.lock`,
   `MODULE.bazel.lock`), `**/gazelle_python.yaml`, `**/*-baseline.xml` (e.g. `ktlint-baseline.xml`),
   `**/.release-please-manifest.json`, `bazel-*/**`, `**/node_modules/**`, `**/*.venv/**`, `.git/**`.
   (A comprehensive static list is harmless for languages a given repo doesn't use.)
4. **`LICENSE` file** — a single templated `LICENSE` whose content switches on `license_id` (full
   standard text for each of the 4 licenses, with `{{ .Scaffold.copyright_holder }}` + year filled in).
   Replaces the template's current static Apache-2.0 `LICENSE` **only when the feature is on**.

### Behaviour when OFF
No license question files included; the template's current Apache-2.0 `LICENSE` is untouched; no check.
Output byte-identical to today.

### Explicitly NOT doing (YAGNI)
- No standalone `.github/workflows/license-check.yml` (the check rides `bazel test`/`aspect`).
- No language-by-language hand-maintained header logic (`addlicense` generates headers).

---

## Feature B — Copybara bidirectional sync

### Questions (added to `scaffold.yaml`)
- `copybara` — `confirm`, **default off**. Master toggle.
- `copybara_components` — free text, comma-separated, **blank allowed**, asked `when` `copybara` is true
  (e.g. "api,worker"). Seeds the engine's component list; blank → machinery-only.
- `copybara_pulumi_auth` — `confirm`, **default off**, nested (asked `when` `copybara` is true). Gates the
  Pulumi auth-provisioning layer.

### Bazel-as-entrypoint principle
Every operation is a **Bazel run-target** under `//tools/copybara`; CI/devs invoke `bazel run`, never a
loose shell script. The external Copybara engine is *executed by* a Bazel target (Bazel is the
entrypoint; the engine stays external).

### Behaviour when ON (files gated by `{{ .Scaffold.copybara }}`)
- **`//tools/copybara:conflict_precheck`** and **`//tools/copybara:drift_check`** — `rules_go`
  `go_binary` targets that port the logic currently in `vitruvian-core`'s `conflict-precheck.sh` /
  drift-check shell (shelling out to `git` via `os/exec`). Unit-tested via `go_test` (`bazel test`).
  Faithful 1:1 port of the validated logic.
- **`//tools/copybara:export`** and **`//tools/copybara:import`** — Bazel run-targets that execute the
  external engine by running the **pinned Copybara Docker image**
  (`olivr/copybara@sha256:87e2e9089344e64693faebb2ee0ed33b8797358c0420b0fa98325ca611e98679`, the version
  proven in `vitruvian-core`). (Maximally-hermetic alternative — building Copybara from source as a
  `bazel_dep` — is noted as a future upgrade, not the default: heavy dep + version-behaviour risk.)
- **`tools/copybara/copy.bara.sky`** — the generalized engine config (Copybara Starlark, i.e. data, not
  a script). Its `COMPONENTS` list is **seeded from `copybara_components`** (empty list if blank).
- **`.github/workflows/_copybara-export.yaml`, `_copybara-import.yaml`, `copybara-drift-check.yaml`** —
  thin callers that `bazel run` the targets above (no committed shell logic).
- **`docs/copybara-bidi-sync.md`** — the admin runbook (generalized from `vitruvian-core`), including the
  "onboard a component" flow used to add per-component caller workflows **post-generation** (Approach 1:
  callers are not generated at gen time).
- **Reference `sync-to-monorepo.yaml` snippet** under `tools/copybara/` — the dispatch file lives in the
  *standalone* repo (which the template cannot write to); the runbook documents installing it there.

### Behaviour when ON + `copybara_pulumi_auth` ON
- **`infrastructure/<pulumi>/.../copybara_sync/*`** — parameterized GitHub-App + per-component deploy-key
  IaC, generalized from `vitruvian-core`'s `sync.go`.

### Behaviour when OFF
No `tools/copybara/`, no copybara workflows, no runbook, no auth IaC. Output byte-identical to today.

### Native-Bazel scope note
The Copybara *engine* is inherently external (a pinned Docker image; `copy.bara.sky` is Copybara's own
Starlark; the sync is a git operation). Bazel is the entrypoint and the orchestration logic
(precheck/drift) is native Go, but the engine itself is not reimplemented in Bazel. This is the correct
boundary for "execute something even if it's external."

---

## Presets (`deliver` CI matrix coverage)
- Extend or add a license-on preset (e.g. set `license: true`, `license_id: Apache-2.0`,
  `copyright_holder: "Example Org"` on an existing preset such as `kitchen-sink`, or a dedicated
  `kitchen-sink-licensed`).
- Add a `copybara-go` preset (`langs: [Go]`, `copybara: true`, blank components, `copybara_pulumi_auth:
  false`) so the matrix generates + builds a Copybara-enabled repo.

---

## Testing & validation

1. **Born-green via `deliver`:** each new/extended preset generates a repo that must pass
   `bazel build //... && bazel test //...` (which now includes the license check and the copybara
   Go-target tests) and `actionlint` on generated workflows.
2. **"Zero impact when off" guard:** generate with both features OFF and diff against current output —
   must be **byte-identical**. This is the load-bearing guarantee.
3. **Local validation (primary gate given the CI-runner blocker, see Risks):** generate sample repos
   locally and run `bazel build //...`, `bazel test //...`, `bazel run //tools/copybara:drift_check`,
   `aspect format --check` (or equivalent), and `actionlint` directly — none of which require the
   org's CI runners.
4. **Unit tests:** `go_test` for `conflict_precheck` and `drift_check`.

---

## Edge cases & risks

1. **Born-green requirement (key):** when `license` is ON, the generated *boilerplate source* must
   already satisfy the check or the repo's first `bazel test` is red. Two candidate mechanisms — **(a)** a
   Scaffold-templated header partial prepended to boilerplate (per-language comment style × SPDX text;
   no generation-time dependency, bounded maintenance), or **(b)** a generation-time `aspect format` pass
   (Bazel-entrypointed, zero maintenance, but needs Bazel/network during generation). **Lean: (a)** to
   keep generation fast/offline. **Resolve concretely in planning.**
2. **license × copybara interaction:** with both on, the generated copybara workflow/`.bzl`/Go files also
   need headers — the same partial/format pass must cover them.
3. **`copybara_components` blank:** machinery-only; drift-check over zero components is a green no-op.
4. **Docker dependency:** `:export`/`:import` need Docker (CI runners have it; document for local use).
5. **`addlicense` provisioning:** prefer multitool; fall back to `rules_go` `go_binary` if no prebuilt
   release binaries (verify in planning).
6. **CI-runner blocker (external):** the VitruvianSoftware org's GitHub Actions runner allocation is
   currently failing (billing/quota; tracked as task #77). The `deliver` matrix validation therefore
   may not run until that's resolved (or the repos are made public for free Actions). **Local validation
   is the primary gate**; CI matrix validation is confirmed once runners are available.

---

## Out of scope
- Backporting `vitruvian-core`'s Dependabot setup (the template uses Renovate — separate decision).
- Auto-generating per-component caller workflows at gen time (Approach 2's post-gen hook) — deferred.
- Building Copybara from source via Bazel (hermetic engine) — deferred future upgrade.

## Open decisions to finalize in planning
1. Born-green mechanism: (a) templated header partial vs (b) generation-time `aspect format`.
2. License add/check wiring: `rules_lint` formatter entry vs dedicated `bazel run`/`bazel test` target.
   **Lean:** a dedicated `bazel run //tools/license:add` + `bazel test //tools/license:check` target
   (most predictable, fewest assumptions about `rules_lint`'s formatter extensibility); only fold into
   the `rules_lint` format flow if it cleanly accepts `addlicense` as a formatter.
3. `addlicense` provisioning: `rules_multitool` vs `rules_go` `go_binary`.
4. Exact Scaffold question syntax for single-select / text / `when:` at the pinned version.
