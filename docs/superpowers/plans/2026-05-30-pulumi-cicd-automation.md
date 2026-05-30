# Pulumi-in-CI/CD Automation — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire the existing `repo_config` Pulumi (and the central stack) into GitHub Actions — opt-in **preview-on-PR** + **apply-on-merge**, gated by `:setup`-managed repo variables, authenticated by a least-privilege GitHub App created via the **Manifest flow** (no public endpoint).

**Architecture:** Templated, runtime-gated GitHub Actions workflows that call the existing `bazel run //infrastructure/pulumi/repo_config:{preview,up}` wrappers. A `:create-app` helper drives GitHub's App Manifest flow (browser-side `localhost` redirect + paste fallback) and `gh secret set`s the credentials; `:setup` prompts the two toggles and `gh variable set`s them. The central stack (template-generator repo, not a Bazel workspace) gets a **preview-only** workflow using the raw `pulumi` CLI.

**Tech Stack:** GitHub Actions, `actions/create-github-app-token`, GitHub App Manifest flow (`POST /app-manifests/{code}/conversions`), existing Bazel `sh_binary` wrappers, `gh` CLI (`variable set` / `secret set`), pulumi-github v6, Scaffold templating.

---

## ⚠️ Templating gotcha (read first)

These workflows are Scaffold (Go-template) files, so literal GitHub Actions `${{ … }}` expressions **must be escaped** or the render breaks (this bit the RBE work). Use the established repo convention (see `_copybara-export.yaml:93`):

- Source: `${{ "{{ secrets.PULUMI_ACCESS_TOKEN }}" }}` → renders to `${{ secrets.PULUMI_ACCESS_TOKEN }}`.
- Pattern: write `${{ "{{ <gha-expr> }}" }}` for every `${{ <gha-expr> }}`.

The central-stack workflow (`infra-preview.yaml`) lives in the **template-generator repo root**, is **not** templated, so its `${{ }}` are written normally (no escaping).

## License headers

`.yml` workflow files are scanned by the license check. Every templated workflow starts with the same per-SPDX `{{ if .Scaffold.license }}…{{ end }}` header block used by `_copybara-export.yaml` (copy it verbatim). The generator's own `infra-preview.yaml` follows whatever `ci.yaml`/`deliver.yaml` do (verify; they currently carry no SPDX header — match that).

---

## File Structure

- **Create** `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/create_app.sh` — App Manifest-flow bootstrap (manifest → browser form → paste/localhost code → conversions exchange → `gh secret set`).
- **Modify** `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/defs.bzl` — add a standalone `:create-app` `sh_binary` (not per-project; it's repo-level).
- **Modify** `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/BUILD` — `exports_files(["create_app.sh"])`.
- **Modify** `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/pulumi_setup.sh` — add the two opt-in prompts + `gh variable set` + secret hints + optional `create-app` call.
- **Create** `{{ "{{ .ProjectSnake }}" }}/.github/workflows/_repo-config-preview.yaml`.
- **Create** `{{ "{{ .ProjectSnake }}" }}/.github/workflows/_repo-config-apply.yaml`.
- **Create** `.github/workflows/infra-preview.yaml` (generator repo root).
- **Modify** `{{ "{{ .ProjectSnake }}" }}/infrastructure/pulumi/repo_config/README.md` — document the automation + manifest flow + the two variables.

**Resolved open items:** (a) credentials are placed via `gh` by the operator (using their authenticated `gh`), NOT a Pulumi module (avoids the chicken-and-egg of needing a token to provision the token); (b) preview comments use `marocchino/sticky-pull-request-comment` (single updating comment); (c) **single shared org App** for VitruvianSoftware's repos — `:create-app` runs **once** against the org manifest endpoint and sets **org-level** `PULUMI_APP_ID` (variable) + `APP_PRIVATE_KEY` (secret) + `PULUMI_ACCESS_TOKEN` (secret), which all org repos inherit; per-repo `:setup` then only sets the two **toggle** variables and ensures the App is installed on the repo. (App ID is non-sensitive → an org *variable*; the private key + Pulumi token are org *secrets*.)

---

## Task 1: App Manifest-flow bootstrap helper

> **Shared org App:** this target is run **once per GitHub org** (e.g. VitruvianSoftware). It creates ONE App and sets **org-level** credentials that all org repos inherit — it is NOT run per-repo.

**Files:**
- Create: `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/create_app.sh`
- Modify: `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/defs.bzl` (add `create_app` sh_binary), `tools/pulumi/BUILD` (export the script)

- [ ] **Step 1: Write `create_app.sh`** (per-SPDX `#` license header like `pulumi_setup.sh`, then):
  - `set -euo pipefail`; require `gh` + `python3`/`jq` + a browser-capable operator.
  - Read target `OWNER` (org or user) — arg or prompt; default to `gh api user --jq .login`.
  - Build the manifest JSON (heredoc): `name`, `url`, `redirect_url: http://localhost:8723/cb`, `public: false`, `default_permissions: {administration: write, contents: read}`, `default_events: []`.
  - Write a temp HTML file with an auto-submitting `<form method="post" action="https://github.com/settings/apps/new?state=…">` (org → `…/organizations/$OWNER/settings/apps/new`) carrying `<input name="manifest" value='…'>`; `open`/`xdg-open` it (print the path if no opener).
  - Capture the `code`: try a one-shot `nc -l 8723` localhost capture; **fallback** — prompt `read -r CODE` ("paste the `code=` value from the browser URL").
  - Exchange: `gh api --method POST /app-manifests/$CODE/conversions` → parse `.id`, `.pem`, `.slug`.
  - Set **org-level** credentials (shared across all org repos): `gh variable set PULUMI_APP_ID --org "$OWNER" --visibility all --body "$id"`; `printf '%s' "$pem" | gh secret set APP_PRIVATE_KEY --org "$OWNER" --visibility all` (never echo the pem); prompt for + `gh secret set PULUMI_ACCESS_TOKEN --org "$OWNER" --visibility all` (read silently).
  - Print the App's install URL (`https://github.com/apps/<slug>/installations/new`) and instruct the operator to install it on the org / target repos.
- [ ] **Step 2: shellcheck** — `shellcheck create_app.sh` → clean (mirror the RBE `setup.sh` lint bar).
- [ ] **Step 3: Add the Bazel target** in `defs.bzl` (repo-level, outside `pulumi_project`):
  ```python
  def pulumi_create_app(name = "create-app", visibility = ["//visibility:public"]):
      sh_binary(name = name, srcs = ["//tools/pulumi:create_app.sh"], visibility = visibility)
  ```
  and `exports_files(["create_app.sh"])` in `tools/pulumi/BUILD`; call `pulumi_create_app()` from a suitable always-shipped package (e.g. `tools/pulumi/BUILD` itself via a small `BUILD` target, or `//infrastructure/pulumi/repo_config:create-app`). Keep it license-header-clean.
- [ ] **Step 4: Render + build check** — generate a preset, `bazel build //tools/pulumi:create-app` (or chosen label) succeeds; no stray `{{ }}`.
- [ ] **Step 5: Commit** — `feat(pulumi): :create-app manifest-flow App bootstrap helper`.

## Task 2: `:setup` opt-in prompts + variable/secret wiring

**Files:** Modify `{{ "{{ .ProjectSnake }}" }}/tools/pulumi/pulumi_setup.sh`

- [ ] **Step 1:** After the existing `go mod tidy` block (repo_config case), add an **automation opt-in** section:
  - `read -r -p "Enable Pulumi preview on pull requests? [y/N] " ANS` → `gh variable set REPO_CONFIG_PREVIEW_ENABLED --body $([ "$ANS" = y ] && echo true || echo false)`.
  - Same for `"Enable Pulumi apply on merge (auto-up)? [y/N]"` → `REPO_CONFIG_AUTO_APPLY`.
  - Credentials are **org-level** (set once by `:create-app`), so `:setup` does NOT set per-repo creds. It instead: (a) checks the org provides `PULUMI_APP_ID` (variable) + `APP_PRIVATE_KEY` + `PULUMI_ACCESS_TOKEN` (secrets) via `gh variable list --org`/`gh secret list --org` — if missing, points the operator at `bazel run //…:create-app`; (b) verifies/prints the App-installation URL so the org App is installed on this repo.
  - EOF-safe defaults (`|| true`, default `N`) so non-interactive `bazel test` never hangs.
- [ ] **Step 2: shellcheck** clean.
- [ ] **Step 3: Render check** — generate a preset; `pulumi_setup.sh` renders cleanly; born-green header intact.
- [ ] **Step 4: Commit** — `feat(pulumi): :setup prompts for PR-preview / merge-apply opt-in`.

## Task 3: `_repo-config-preview.yaml` (templated, gated)

**Files:** Create `{{ "{{ .ProjectSnake }}" }}/.github/workflows/_repo-config-preview.yaml`

- [ ] **Step 1: Write the workflow** — license header block (copy from `_copybara-export.yaml`), then (GHA expressions escaped `${{ "{{ … }}" }}`):
  ```yaml
  name: Repo Config Preview
  on:
      pull_request:
          paths: ["infrastructure/pulumi/repo_config/**"]
  jobs:
      preview:
          if: ${{ "{{ vars.REPO_CONFIG_PREVIEW_ENABLED == 'true' }}" }}
          runs-on: ubuntu-latest
          steps:
              - uses: actions/checkout@v6
              - id: app-token
                uses: actions/create-github-app-token@v2
                with:
                    app-id: ${{ "{{ vars.PULUMI_APP_ID }}" }}
                    private-key: ${{ "{{ secrets.APP_PRIVATE_KEY }}" }}
              - name: bazel run preview
                env:
                    GITHUB_TOKEN: ${{ "{{ steps.app-token.outputs.token }}" }}
                    GITHUB_OWNER: ${{ "{{ github.repository_owner }}" }}
                    PULUMI_ACCESS_TOKEN: ${{ "{{ secrets.PULUMI_ACCESS_TOKEN }}" }}
                run: bazel run //infrastructure/pulumi/repo_config:preview -- --stack dev --diff 2>&1 | tee /tmp/preview.txt
              - uses: marocchino/sticky-pull-request-comment@v2
                with:
                    header: repo-config-preview
                    path: /tmp/preview.txt
  ```
  (Decide `APP_ID` as a var vs secret during impl — App ID is non-sensitive, a `var` is fine; the private key is always a secret. Keep consistent with `:setup`/`:create-app`.)
- [ ] **Step 2: Render check** — generate a preset; `cat` the rendered file → exact GHA expressions present, no `{{ }}` leakage, header rendered.
- [ ] **Step 3: License check** — `addlicense -check` (or `bazel test //tools/license:check`) clean on the rendered tree.
- [ ] **Step 4: Commit** — `feat(pulumi): PR-preview workflow for repo_config (gated)`.

## Task 4: `_repo-config-apply.yaml` (templated, gated)

**Files:** Create `{{ "{{ .ProjectSnake }}" }}/.github/workflows/_repo-config-apply.yaml`

- [ ] **Step 1: Write the workflow** — same header + token-mint pattern; trigger and guard differ:
  ```yaml
  name: Repo Config Apply
  on:
      push:
          branches: ["${{ "{{ github.event.repository.default_branch }}" }}"]
          paths: ["infrastructure/pulumi/repo_config/**"]
  jobs:
      apply:
          if: ${{ "{{ vars.REPO_CONFIG_AUTO_APPLY == 'true' }}" }}
          runs-on: ubuntu-latest
          steps:
              - uses: actions/checkout@v6
              - id: app-token
                uses: actions/create-github-app-token@v2
                with: { app-id: "${{ "{{ vars.PULUMI_APP_ID }}" }}", private-key: "${{ "{{ secrets.APP_PRIVATE_KEY }}" }}" }
              - env:
                    GITHUB_TOKEN: ${{ "{{ steps.app-token.outputs.token }}" }}
                    GITHUB_OWNER: ${{ "{{ github.repository_owner }}" }}
                    PULUMI_ACCESS_TOKEN: ${{ "{{ secrets.PULUMI_ACCESS_TOKEN }}" }}
                run: bazel run //infrastructure/pulumi/repo_config:up -- --stack dev --yes
  ```
  (Note: `branches:` on the default branch — confirm the literal-vs-expression form during impl; a static `["main"]` plus a comment is acceptable if the expression form is awkward.)
- [ ] **Step 2: Render + license check** — clean, no brace leakage.
- [ ] **Step 3: Commit** — `feat(pulumi): apply-on-merge workflow for repo_config (gated)`.

## Task 5: Central-stack preview workflow (generator repo)

**Files:** Create `.github/workflows/infra-preview.yaml` (repo root — NOT templated)

- [ ] **Step 1: Write the workflow** (plain `${{ }}`, no escaping):
  - `on: pull_request: paths: ["infrastructure/pulumi/**"]`.
  - Steps: checkout → `actions/setup-go` → install pulumi → mint App token (or a `PULUMI_GITHUB_TOKEN` secret) → `cd infrastructure/pulumi && GITHUB_TOKEN=… GITHUB_OWNER=VitruvianSoftware pulumi preview --stack dev` → sticky comment.
  - **No apply job.**
- [ ] **Step 2:** Confirm it does NOT run on `push` (preview only) and the diff renders in a PR comment.
- [ ] **Step 3: Commit** — `feat(pulumi): preview-on-PR for the central infra stack`.

## Task 6: Verify across presets + dogfood vitruvian-core

- [ ] **Step 1:** Generate `kitchen-sink` + a copybara preset → both `_repo-config-*.yaml` render cleanly (no stray `{{ }}`), born-green; `:setup`/`:create-app` render; `bazel build //...` of the new targets passes.
- [ ] **Step 2:** Commit + PR to `platform-v2.0`; **watch CI to green** (retry transient CDN 502s); gated merge → **watch deliver** (26 Starters, workflows inert because vars unset).
- [ ] **Step 3: Dogfood (standing practice).** Port to vitruvian-core; `bazel run //…:create-app` (manifest flow → App + secrets), install the App, `bazel run //…:repo_config:setup` answering **yes/yes** → `gh variable list` shows both `true`.
- [ ] **Step 4:** Open a PR in vitruvian-core touching `repo_config/` → preview workflow posts the diff comment; merge → apply workflow runs `:up`, applies, CI green. Toggle a var off → workflow no-ops.
- [ ] **Step 5:** Add `infra-preview.yaml` to the generator repo; open a PR touching `infrastructure/pulumi/` → confirm the preview comment; confirm no apply workflow exists.
- [ ] **Step 6:** Fix any gap back in the template + re-deliver before calling done.

---

## Self-Review

- **Spec coverage:** preview-on-PR (T3), apply-on-merge (T4), `:setup` opt-in via repo vars (T2), manifest-flow App bootstrap with no public endpoint (T1), central preview-only (T5), dogfood vitruvian-core (T6) — all covered.
- **Placeholder scan:** none; the two impl-time decisions (APP_ID var-vs-secret; default-branch literal-vs-expression) are explicitly flagged inline, not left vague.
- **Type/name consistency:** variables `REPO_CONFIG_PREVIEW_ENABLED` / `REPO_CONFIG_AUTO_APPLY`; secrets `APP_PRIVATE_KEY` / `PULUMI_ACCESS_TOKEN`; `PULUMI_APP_ID` (var) used consistently across T1–T4.
- **Escaping:** every templated GHA `${{ }}` uses the `${{ "{{ … }}" }}` form; T5 (non-templated) uses plain form. Called out up front.

## Execution Handoff

Plan saved to `docs/superpowers/plans/2026-05-30-pulumi-cicd-automation.md`. Two execution options when ready:
1. **Subagent-Driven (recommended)** — fresh subagent per task, two-stage review between tasks.
2. **Inline Execution** — executing-plans, batched with checkpoints.

(Gated steps remain gated: the deliver fan-out merge and any live `pulumi up` pause for confirmation, per standing practice.)
