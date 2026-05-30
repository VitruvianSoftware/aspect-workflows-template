# Copybara Sync — Optional Pulumi Auth IaC

This directory contains a **self-contained Pulumi Go module** that provisions the
GitHub authentication resources required by the Copybara bidirectional sync.

It is generated only when `copybara_pulumi_auth` is enabled at scaffold time.
If you are provisioning auth manually instead, see
[docs/copybara-bidi-sync.md](../../docs/copybara-bidi-sync.md) (§6).

---

## What this provisions

For each component listed in `syncedProjects` in `sync.go`:

| Resource | Where | Purpose |
|---|---|---|
| `tls.NewPrivateKey` (ED25519) | Pulumi state | Keypair for export push auth |
| `github.NewRepositoryDeployKey` (write) | Standalone repo | Allows the monorepo export workflow to push |
| `github.NewActionsSecret` `<PREFIX>_SYNC_SSH_KEY` | Monorepo | Private half of the deploy key |
| `github.NewActionsSecret` `<PREFIX>_DISPATCH_APP_ID` | Standalone repo | GitHub App ID for the dispatch workflow |
| `github.NewActionsSecret` `<PREFIX>_DISPATCH_APP_PRIVATE_KEY` | Standalone repo | GitHub App private key |

Plus monorepo-wide secrets for the shared dispatch App (both Actions and Dependabot):
`SYNC_APP_ID`, `SYNC_APP_PRIVATE_KEY`.

The GitHub App is created **manually** by the operator (GitHub has no headless
App-creation API). Pulumi only places the credentials.

---

## One-time setup

### 1. Edit `sync.go`

- Set `githubOrg` to your GitHub organisation or user (search for `YOUR_GITHUB_ORG`).
- Verify `syncedProjects` lists all components you want managed. If you specified
  `copybara_components` at generation time they are already seeded.
- Update the `module` path in `go.mod` to match your actual GitHub org.

### 2. Create the GitHub App (manual, one-time)

1. Go to **GitHub → Settings → Developer settings → GitHub Apps → New GitHub App**.
2. Give it a name (e.g. `<your-org>-copybara-sync`).
3. Grant **Contents: Read & write** and **Metadata: Read** (minimum required).
4. Install the App on **this monorepo only** (least-privilege — it fires
   `repository_dispatch` into the monorepo).
5. Note the **App ID** and generate a **private key** (`.pem` file).

### 3. Set Pulumi config secrets

```bash
# Shared App credentials (monorepo-wide):
pulumi config set --secret syncAppId          <App ID>
pulumi config set --secret syncAppPrivateKey  < /path/to/app.pem

# Per-component App credentials (same App, one entry per component):
# Replace "myService" with the lower-camelCase form of each component name
# (e.g. "my-service" → "myService", "api" → "api").
pulumi config set --secret myServiceDispatchAppId          <App ID>
pulumi config set --secret myServiceDispatchAppPrivateKey  < /path/to/app.pem
```

### 4. Set the GitHub provider token

```bash
export GITHUB_TOKEN=<your PAT or token with repo + admin:repo_hook scope>
```

### 5. Integrate into your Pulumi stack

This module is intentionally a **standalone Go module** so the monorepo's Bazel /
gazelle / go.work setup does not pull it in. Integrate it in one of two ways:

**Option A — call from an existing stack:**

```go
import copybara_sync "github.com/YOUR_GITHUB_ORG/{{ .ProjectKebab }}/infrastructure/pulumi/copybara_sync"

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        return copybara_sync.ManageSyncAuth(ctx)
    })
}
```

**Option B — run as a standalone program:**

```bash
cd infrastructure/pulumi/copybara_sync
go mod tidy
# Add a main.go that calls ManageSyncAuth (see example above)
pulumi up --stack dev
```

---

## Ongoing operations

- **Rotate a deploy key:** change the `tls.NewPrivateKey` resource (e.g. rename it
  to force re-creation), then `pulumi up`. Pulumi replaces the keypair and updates
  the deploy key and secret atomically.
- **Rotate the App private key:** generate a new key in the App's GitHub settings,
  then `pulumi config set --secret <key> < new.pem` and `pulumi up`.
- **Add a component:** append to `syncedProjects` in `sync.go`, add its config
  secrets, and `pulumi up`. Then follow §8f of the runbook to seed the baselines.
- **Remove a component:** remove it from `syncedProjects` and run `pulumi up`.
  Pulumi destroys the deploy key and secrets.

See [docs/copybara-bidi-sync.md](../../docs/copybara-bidi-sync.md) for the full
onboarding runbook and troubleshooting guide.
