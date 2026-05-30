# SOP: Rotating secrets & API keys

Standard operating procedure for rotating the credentials this repo depends on —
primarily the **BuildBuddy RBE API key** used by the remote-build helper
(`bazel run //tools/remote:setup`; see [remote-build.md](remote-build.md)){{ if .Scaffold.copybara }}, plus the
Copybara sync secrets{{ end }}. The same mechanics apply to any credential stored
as a GitHub Actions secret.

## When to rotate

- **Suspected exposure** — the secret showed up in a log, screenshot, chat
  message, committed file, or was handled insecurely. Rotate immediately.
- **Offboarding** — someone with access to the secret (or its source account)
  leaves.
- **Schedule** — on a regular cadence (e.g. quarterly) as hygiene.

## This repo's secrets

List them with `gh secret list` (run from the repo):

| Secret | Purpose | Source of truth |
|---|---|---|
| `BUILDBUDDY_API_KEY` | RBE: authenticates Bazel to BuildBuddy (remote execution + cache + BES) | BuildBuddy → Org API keys |
{{- if .Scaffold.copybara }}
| `SYNC_APP_ID`, `SYNC_APP_PRIVATE_KEY` | Copybara sync dispatch (GitHub App) | the GitHub App's settings |
| `<COMPONENT>_SYNC_SSH_KEY` (one per synced component) | Copybara per-component deploy keys | each component's standalone-repo deploy keys |
{{- end }}

> A freshly generated repo has none of these set yet — they appear once you run
> the remote-build helper and enable RBE in CI{{ if .Scaffold.copybara }}, or wire up Copybara sync{{ end }}.

The detailed steps below are written for `BUILDBUDDY_API_KEY`.{{ if .Scaffold.copybara }} For the
Copybara secrets, only **how you create the new credential** differs — see the
brief section at the end.{{ end }}

---

## Rotating `BUILDBUDDY_API_KEY`

### Where the key lives (everything you must update)

| Location | Holds the key? | In git? |
|---|---|---|
| BuildBuddy org → Org API keys | source of truth | — |
| `BUILDBUDDY_API_KEY` GitHub Actions repo secret | yes (encrypted) | no |
| each developer's `user.bazelrc` | yes (`--remote_header` / `--bes_header`) | **no — git-ignored** |
| `tools/remote.bazelrc` | **no** — endpoints/platform only | yes (committed; safe — no secret) |

### Steps

**1. Create the new key** (BuildBuddy — requires **org admin**)
- In the BuildBuddy app: **Settings → Org API keys → Create new API key**.
- Name it and grant the capabilities RBE needs (**Cache** + **Remote execution**).
- Copy the value — BuildBuddy shows it **once**. Keep it out of chat, logs, and
  screenshots.

**2. Update the CI secret** — pipe via stdin so the key never lands in shell
history or a process listing (run from the repo so `gh` targets it):
```sh
printf '%s' '<NEW_KEY>' | gh secret set BUILDBUDDY_API_KEY
```
(Or via the UI: repo → **Settings → Secrets and variables → Actions →
`BUILDBUDDY_API_KEY` → Update**.)

**3. Update local dev config** — each developer who builds with `--config=remote`
locally rewrites the key in their git-ignored `user.bazelrc`. Easiest is to
re-run the helper (it replaces the prior key line idempotently):
```sh
bazel run //tools/remote:setup -- --provider buildbuddy --no-ci
# paste the new key at the hidden prompt (or `export BUILDBUDDY_API_KEY=...` first)
```
Never commit `user.bazelrc`.

> **Shortcut:** running the full helper *without* `--no-ci` does steps **2 and 3
> together** — it rewrites your local `user.bazelrc` **and** re-sets the
> `BUILDBUDDY_API_KEY` CI secret via `gh` (needs `gh` authenticated + repo admin).

**4. Verify the new key works — _before_ revoking the old one:**
- **CI (authoritative):** re-run the latest Bazel job (or push a trivial commit).
  It should authenticate and stream to BuildBuddy — the job stays green and the
  build appears in the org's build history.
- **Local (optional, macOS-safe):**
  ```sh
  bazel build --config=remote --remote_executor= //...
  ```
  This builds locally but streams to BuildBuddy with the new key (cache + BES, no
  remote execution), so it works from a Mac and confirms auth + that a build
  shows up. Full remote *execution* is verified by CI (Linux).

**5. Revoke the old key** — BuildBuddy → **Org API keys → (old key) → Delete**.
Do this only **after** step 4 passes, so in-flight builds aren't cut off.

### Security rules

- The key lives **only** in the git-ignored `user.bazelrc` and the GitHub Actions
  secret. `tools/remote.bazelrc` is committed but contains **no** secret
  (endpoints + exec platform only).
- Never commit the key; never paste it into chat, PRs, logs, or screenshots. If
  it is exposed, treat it as compromised and rotate immediately.
- The key is **org-scoped** — rotating affects everyone using it, so coordinate
  the swap.
{{- if .Scaffold.copybara }}

---

## Rotating the Copybara secrets (brief)

Same **create → update secret → verify → revoke** flow; only creation differs:

- **`<COMPONENT>_SYNC_SSH_KEY`** — generate a new deploy key pair (`ssh-keygen`),
  add the **public** key to the component's standalone repo under **Settings →
  Deploy keys** (allow write), then
  `printf '%s' "$(cat newkey)" | gh secret set <COMPONENT>_SYNC_SSH_KEY`. Delete
  the old deploy key after a sync run succeeds.
- **`SYNC_APP_ID` / `SYNC_APP_PRIVATE_KEY`** — generate a new private key in the
  GitHub App's settings, set `SYNC_APP_PRIVATE_KEY` (stdin), verify a dispatch,
  then delete the old private key. If these are managed by the Pulumi
  `copybara_sync` IaC, prefer rotating them there.
{{- end }}
