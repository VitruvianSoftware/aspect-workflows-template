# Aspect Workflows Template Infrastructure

This directory contains the Pulumi Infrastructure-as-Code (IaC) configuration to manage the main `aspect-workflows-template` repository and all of its external starter repositories. It automates the creation of the repositories and manages the distribution of SSH deploy keys between the main repository (as Actions Secrets) and the starter repositories (as Deploy Keys).

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- [Go](https://go.dev/doc/install)
- GitHub Personal Access Token (PAT) with appropriate permissions (`repo`, `admin:org` if managing org-level settings).

## Configuration

When running Pulumi commands (like `pulumi up`, `pulumi preview`, or `pulumi import`), the Pulumi GitHub provider needs to know which GitHub organization or user account to target.

### 1. Environment Variables (Temporary)

You can set the required environment variables in your terminal before running Pulumi:

```bash
export GITHUB_TOKEN="your_github_token_here"
export GITHUB_OWNER="VitruvianSoftware"
```

### 2. Pulumi Configuration (Persistent)

For a more permanent setup, you can save these settings directly into your Pulumi stack's configuration (e.g., `Pulumi.dev.yaml`):

```bash
# Set the target organization or user
pulumi config set github:owner VitruvianSoftware

# Securely set your GitHub token (this will be encrypted in the Pulumi state)
pulumi config set --secret github:token "your_github_token_here"
```

This ensures that Pulumi will automatically target the correct organization (`VitruvianSoftware`) for all GitHub resources managed by this stack without needing to export environment variables every time.

## Usage

1. **Initialize the stack** (if not already done):
   ```bash
   pulumi stack init dev
   ```

2. **Preview changes**:
   ```bash
   pulumi preview
   ```

3. **Apply changes**:
   ```bash
   pulumi up
   ```
