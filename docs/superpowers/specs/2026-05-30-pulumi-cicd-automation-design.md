# Pulumi-in-CI/CD Automation — Design

**Date:** 2026-05-30
**Status:** Approved (brainstorming) — pending implementation plan

## Context

The repo-config-as-code initiative (just shipped) manages each repo's GitHub
settings — `delete_branch_on_merge` + parameterized branch protection — via a
standalone Pulumi module (`infrastructure/pulumi/repo_config/`) in every
generated repo, plus the central template-generator stack
(`aspect-workflows-template-infra`). It's all driven through
`bazel run //…/repo_config:{preview,up,destroy,refresh,config,setup}` wrappers,
but **every action is still a manual invocation** — nothing is wired into
CI/CD. A developer must remember to run `:up` after changing config, and a PR
gives no automatic "what would this change?" visibility.

This feature wires Pulumi into GitHub Actions so repo-config changes are
**previewed on PRs** and (optionally) **auto-applied on merge** — a GitOps loop
— while keeping the high-blast-radius central stack safely manual. It adds no
new Pulumi program logic; it only orchestrates the wrappers that already exist.

## Decisions (from brainstorming)

1. **Fully opt-in.** A developer who doesn't use Pulumi gets nothing firing in
   their CI. "Preview on PR" and "apply on merge" are two independent opt-ins.
2. **`:setup` is the single control point.** `bazel run //…/repo_config:setup`
   prompts two questions (enable preview-on-PR? enable apply-on-merge?) and
   persists the answers as **GitHub Actions repo variables**.
3. **Tiered behavior by ownership + risk:**
   - **Generated repos (developers):** both toggles off by default; opt in via `:setup`.
   - **vitruvian-core (owned — dogfood):** both on → full GitOps.
   - **Central template-generator stack:** preview-on-PR only, **never** auto-up
     (a bad `up` rotates all 26 delivery deploy keys — the PRC-F hazard).
4. **Least-privilege GitHub App** for auth (`Administration: write` +
   `Contents: read`); short-lived installation tokens minted per run;
   `PULUMI_ACCESS_TOKEN` for Pulumi Cloud state; App credentials Pulumi-provisioned.
5. **App created via GitHub's App Manifest flow** (browser-side redirect to
   `localhost` — no public endpoint, no UI permission hand-config), driven by a
   `:setup`-integrated / `:create-app` helper. Strongly preferred over manual UI
   App creation for developer experience.

## Architecture

GitHub Actions workflows drive the **existing** `bazel run //…/repo_config:{preview,up}`
wrappers (generated repos) and the raw `pulumi` CLI (central stack — the template
root isn't a Bazel workspace). Each workflow is **gated by a repo Actions
variable**, so it ships everywhere but stays inert until `:setup` enables it.
No re-scaffold is ever required to change a repo's automation state.

Three valid end-states per generated repo:

| State | `REPO_CONFIG_PREVIEW_ENABLED` | `REPO_CONFIG_AUTO_APPLY` |
|---|---|---|
| Don't use Pulumi (default) | off | off |
| Preview only | on | off |
| Full GitOps | on | on |

## Components

### A. Generated-tree workflows (templated, always shipped, runtime-gated)

**`{{ .ProjectSnake }}/.github/workflows/_repo-config-preview.yaml`**
- Trigger: `pull_request`, `paths: ['infrastructure/pulumi/repo_config/**']`.
- Guard: `if: vars.REPO_CONFIG_PREVIEW_ENABLED == 'true'`.
- Steps: checkout → mint App token (`actions/create-github-app-token`) →
  export `GITHUB_TOKEN` + `GITHUB_OWNER` + `PULUMI_ACCESS_TOKEN` →
  `bazel run //infrastructure/pulumi/repo_config:preview -- --stack dev` →
  post the diff as a sticky PR comment. Read-only; no mutation.

**`{{ .ProjectSnake }}/.github/workflows/_repo-config-apply.yaml`**
- Trigger: `push` to the default branch, `paths: ['infrastructure/pulumi/repo_config/**']`.
- Guard: `if: vars.REPO_CONFIG_AUTO_APPLY == 'true'`.
- Steps: checkout → mint App token → `bazel run //infrastructure/pulumi/repo_config:up -- --stack dev --yes`.

Both carry the per-SPDX license header where applicable (workflows are YAML; the
existing license check ignores `**/BUILD` and docs but DOES scan `.yml` — match
the existing copybara workflows' header convention).

### B. `:setup` enhancement (`{{ .ProjectSnake }}/tools/pulumi/pulumi_setup.sh`)

After the existing bootstrap (prereqs → login → stack select/init → adopt hint →
`go mod tidy`), add an **opt-in configuration** section (repo_config only):
- Prompt: *"Enable Pulumi preview on pull requests? [y/N]"* →
  `gh variable set REPO_CONFIG_PREVIEW_ENABLED --body <true|false>`.
- Prompt: *"Enable Pulumi apply on merge (auto-up)? [y/N]"* →
  `gh variable set REPO_CONFIG_AUTO_APPLY --body <true|false>`.
- Print the **secrets** the workflows require and `gh secret set` hints (values
  never echoed): `PULUMI_ACCESS_TOKEN`, plus the App ID + private key (or a note
  to reuse an org-level App).
- Idempotent and non-interactive-safe (defaults to "no" on EOF, like the existing
  stack prompt). Requires the developer's `gh` authed with admin (same bar as the
  App/secret setup).

### C. Central-stack preview (template-generator repo — NOT the generated tree)

**`.github/workflows/infra-preview.yaml`** (alongside `ci.yaml`, `deliver.yaml`)
- Trigger: `pull_request`, `paths: ['infrastructure/pulumi/**']`.
- Steps: checkout → set up Go + pulumi → mint App token → `cd infrastructure/pulumi
  && pulumi preview --stack dev` → post diff as a PR comment.
- **No apply workflow.** `pulumi up` on the central stack stays a gated manual
  step (preserves the deploy keys).

### D. Auth (GitHub App, Pulumi-provisioned)

- A dedicated App (working name **`vitruvian-pulumi`**) with **Administration:
  write** (branch protection + repo settings) + **Contents: read**; installed on
  repos that opt in.
- Workflows mint short-lived installation tokens via `actions/create-github-app-token`
  and export them as `GITHUB_TOKEN` for the Pulumi github provider.
- `PULUMI_ACCESS_TOKEN` stored as an Actions secret for Pulumi Cloud state.
- App **credentials** placed via a small Pulumi auth module mirroring the existing
  `copybara_sync` auth IaC (reproducible, not click-ops).

**App bootstrap — Manifest flow (no public endpoint, no UI hand-config).**
A `:setup`-integrated helper (or a dedicated `//tools/pulumi:create-app` target)
drives GitHub's App Manifest flow:
1. Build a manifest (name, `default_permissions: {administration: write,
   contents: read}`, events, `redirect_url: http://localhost:<port>/callback`).
2. Open the browser to an auto-submitting local HTML form that POSTs the manifest
   to `https://github.com/settings/apps/new` (user) or
   `https://github.com/organizations/<org>/settings/apps/new` (org).
3. The operator clicks **"Create GitHub App"** once — permissions are pre-filled.
4. GitHub redirects the **browser** to `localhost:<port>/callback?code=…`; a
   throwaway local server captures the `code`. **Fallback:** operator copies the
   `code` from the URL bar and pastes it (works with no local server at all).
5. `POST /app-manifests/<code>/conversions` (unauthenticated; **outbound HTTPS
   only**) returns the App `id`, private-key `pem`, and secrets.
6. Pulumi places those credentials as repo secrets; CI then mints installation tokens.

Because the redirect is **browser-side**, this needs **no public endpoint and no
inbound network** — only outbound HTTPS to github.com. App creation + installation
need the operator's browser once; CI is fully headless. **Two-host note:** if the
shell host has no browser, use the manual code-paste fallback or run the bootstrap
on the host that has the browser.

## Reuse (don't reinvent)

- The `tools/pulumi` `bazel run` wrappers (already shipped) — workflows just call them.
- The copybara workflows' App-token-minting + sticky-comment patterns.
- The `copybara_sync` Pulumi auth module as the template for the `vitruvian-pulumi`
  App-credential provisioning.
- `actions/create-github-app-token` (already used by the copybara dispatch flow).

## Verification

1. **Render:** generate a preset — both `_repo-config-*.yaml` render cleanly (no
   stray `{{ }}`), born-green headers; `deliver` matrix stays green (workflows
   present but inert because the vars are unset).
2. **App bootstrap + setup:** run the manifest-flow helper (localhost redirect) →
   confirm it creates the App and returns an `id` + `pem` with no public endpoint
   (and the manual code-paste fallback also works); Pulumi places them as secrets.
   Then in vitruvian-core run `:setup`, answer yes/yes → `gh variable list` shows
   both vars `true`; confirm App token + `PULUMI_ACCESS_TOKEN` secrets exist.
3. **Preview:** open a PR touching `repo_config/` → preview workflow runs, posts
   the diff comment, mutates nothing.
4. **Apply:** merge → apply workflow runs `:up`, applies, CI green; verify the
   intended repo-setting change landed.
5. **Toggle off:** set a var to `false` → the corresponding workflow no-ops.
6. **Central:** open a PR touching `infrastructure/pulumi/` in the template repo →
   `infra-preview.yaml` posts the diff; confirm NO central apply workflow exists.

## Rollout

1. Implement + verify in the template generator (`platform-v2.0`); merge →
   `deliver` ships the (inert) workflows + `:setup` enhancement to the 26 Starters.
2. **Dogfood on vitruvian-core** (standing practice): run `:setup` (yes/yes),
   exercise preview → merge → apply end-to-end; fix any gap back in the template
   and re-deliver before calling it done.
3. Add `infra-preview.yaml` to the template-generator repo for the central stack.

## Non-goals / open items (resolve during writing-plans)

- **No auto-up on the central stack** (deploy-key rotation) — preview only.
- **copybara_sync automation** is out of scope (separate concern; see the flagged
  `pulumi_project`-macro cleanup task).
- **App-credential IaC placement:** extend `repo_config` vs a small dedicated
  `pulumi_auth` module — decide while planning.
- **Preview comment mechanism:** sticky-comment action vs check annotation.
- **BYO App vs shared org App** for generated-repo developers — document both;
  default to "bring your own App, set its ID + private key as secrets."
- **Secret/variable bootstrap in CI runners for deliver:** confirm the inert
  workflows never require secrets at scaffold/deliver time (guards must short-circuit
  before any secret read).
