# Security

This guide covers security best practices for projects using the Aspect Workflows Template.

## Table of Contents

1. [Overview](#overview)
2. [Dependency Security](#dependency-security)
3. [Secrets Management](#secrets-management)
4. [Access Control](#access-control)
5. [Container Security](#container-security)
6. [Supply Chain Security](#supply-chain-security)
7. [Security Auditing](#security-auditing)
8. [Incident Response](#incident-response)

## Overview

Security considerations for Bazel projects:

- Hermetic builds prevent supply chain attacks
- Dependency pinning ensures reproducibility
- Build caching requires secure storage
- Secrets must never be in version control

### Security Layers

```mermaid
graph TB
    subgraph "Development"
        Code[Source Code] --> Review[Code Review]
        Review --> Scan[Security Scan]
    end
    
    subgraph "Build"
        Scan --> Hermetic[Hermetic Build]
        Hermetic --> Deps[Verified Dependencies]
        Deps --> Cache[Secure Cache]
    end
    
    subgraph "Deploy"
        Cache --> Sign[Signed Artifacts]
        Sign --> Deploy[Secure Deployment]
        Deploy --> Runtime[Runtime Protection]
    end
    
    subgraph "Monitor"
        Runtime --> Audit[Security Audit]
        Audit --> Alert[Alerts]
        Alert --> Response[Incident Response]
    end
```

## Dependency Security

### Automated Vulnerability Scanning

**GitHub Dependabot:**

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    reviewers:
      - "security-team"
    labels:
      - "dependencies"
      - "security"
  
  - package-ecosystem: "pip"
    directory: "/"
    schedule:
      interval: "weekly"
  
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

**Security scanning workflow:**

```yaml
# .github/workflows/security-scan.yml
name: Security Scan

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  push:
    branches: [main]
  pull_request:

jobs:
  dependency-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Python vulnerability scan
        run: |
          pip install pip-audit
          pip-audit --format json requirements/all.txt > python-vulns.json
      
      - name: JavaScript vulnerability scan
        run: |
          pnpm audit --json > js-vulns.json
          pnpm audit --audit-level=moderate
      
      - name: Go vulnerability scan
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
      
      - name: Upload results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: security-scan-results
          path: '*-vulns.json'
```

### Manual Security Audits

```bash
# Python - pip-audit
pip-audit requirements/all.txt --format json

# JavaScript - pnpm audit
pnpm audit --audit-level=moderate

# Go - govulncheck
govulncheck ./...

# Ruby - bundler-audit
bundle audit check --update

# Rust - cargo-audit
cargo audit
```

### Vulnerability Response Matrix

| Severity | Response Time | Action Required |
|----------|---------------|-----------------|
| Critical (CVSS 9.0-10.0) | 0-4 hours | Immediate patch, emergency deploy |
| High (CVSS 7.0-8.9) | 24 hours | Priority patch, scheduled deploy |
| Medium (CVSS 4.0-6.9) | 1 week | Regular patch cycle |
| Low (CVSS 0.1-3.9) | Next sprint | Backlog, update when convenient |

### Dependency Pinning

**Always pin versions in production:**

```python
# requirements/runtime.txt
requests==2.31.0
urllib3==2.1.0
certifi==2023.11.17
```

```json
// package.json
{
  "dependencies": {
    "express": "4.18.2",
    "lodash": "4.17.21"
  }
}
```

```go
// go.mod
require (
    github.com/spf13/cobra v1.8.0
    github.com/stretchr/testify v1.8.4
)
```

## Secrets Management

### Never Commit Secrets

**.gitignore configuration:**

```gitignore
# Secrets
*.key
*.pem
*.p12
*.pfx
.env
.env.local
secrets.yaml
credentials.json

# SSH keys
id_rsa
id_ecdsa
id_ed25519

# API keys
**/api_keys.txt
**/tokens.txt
```

### Environment Variables

```bash
# Use environment variables for secrets
export DATABASE_URL="postgresql://user:pass@host/db"
export API_KEY="your-api-key"

# Load from .env file (gitignored)
set -a
source .env
set +a
```

### GitHub Secrets

```yaml
# .github/workflows/deploy.yml
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy with secrets
        env:
          API_KEY: ${{ secrets.API_KEY }}
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: |
          ./deploy.sh
```

**Managing secrets:**

```bash
# Add secret via CLI
gh secret set API_KEY --body "your-api-key"

# Or via web UI
# https://github.com/owner/repo/settings/secrets/actions
```

### Secret Scanning

```yaml
# .github/workflows/secret-scan.yml
name: Secret Scanning

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Run gitleaks
        uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Secrets in Bazel

**Use --action_env for secrets:**

```bash
# Pass secret as environment variable
bazel build --action_env=API_KEY //app

# Or in .bazelrc (for local only, not committed)
build --action_env=API_KEY
```

**Secrets file pattern:**

```python
# BUILD file
genrule(
    name = "deploy",
    srcs = ["deploy.sh"],
    outs = ["deploy.log"],
    cmd = """
        API_KEY=$${API_KEY} \
        $(location deploy.sh) > $@
    """,
)
```

### Vault Integration

**Using HashiCorp Vault:**

```bash
#!/bin/bash
# scripts/get-secrets.sh

# Login to Vault
vault login -method=github

# Read secrets
export DATABASE_URL=$(vault kv get -field=url secret/database)
export API_KEY=$(vault kv get -field=key secret/api)

# Run application
./app
```

## Access Control

### Repository Access

**GitHub teams and roles:**

```markdown
# CODEOWNERS
# Global owners
* @org/admins

# Specific paths
/src/security/ @org/security-team
/.github/ @org/devops-team
/docs/ @org/tech-writers

# Language-specific
*.py @org/python-team
*.go @org/go-team
```

### Branch Protection

**Recommended settings:**

```yaml
# Branch protection for main
- Require pull request reviews (2 reviewers)
- Require status checks to pass
- Require branches to be up to date
- Require conversation resolution
- Require linear history
- Include administrators
- Restrict force pushes
- Restrict deletions
```

**Configure via CLI:**

```bash
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["ci","security-scan"]}' \
  --field required_pull_request_reviews='{"required_approving_review_count":2}' \
  --field enforce_admins=true
```

### Code Review Requirements

**Security-focused reviews:**

```markdown
# Pull Request Security Checklist

- [ ] No hardcoded secrets or credentials
- [ ] Dependencies are up to date
- [ ] No new security vulnerabilities introduced
- [ ] Input validation implemented
- [ ] Error messages don't leak sensitive info
- [ ] Logging doesn't include secrets
- [ ] Authentication/authorization unchanged (or reviewed)
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (proper escaping)
- [ ] CSRF protection in place
```

## Container Security

### Secure Base Images

```python
# Use distroless images
# BUILD file
oci_image(
    name = "app_image",
    base = "@distroless_python",
    entrypoint = ["/app/main"],
    tars = [":app_layer"],
)
```

### Image Scanning

```yaml
# .github/workflows/container-security.yml
name: Container Security

on:
  push:
    branches: [main]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build image
        run: bazel build //app:image
      
      - name: Scan with Trivy
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'bazel/app:image'
          format: 'sarif'
          output: 'trivy-results.sarif'
      
      - name: Upload scan results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
      
      - name: Scan with Grype
        uses: anchore/scan-action@v3
        with:
          image: 'bazel/app:image'
          fail-build: true
          severity-cutoff: high
```

### Runtime Security

```yaml
# Kubernetes security context
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  
  containers:
  - name: app
    image: myapp:latest
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
          - ALL
```

## Supply Chain Security

### SLSA (Supply chain Levels for Software Artifacts)

**Build provenance:**

```yaml
# .github/workflows/slsa.yml
name: SLSA Provenance

on:
  release:
    types: [published]

permissions:
  contents: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      artifacts: ${{ steps.build.outputs.artifacts }}
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build
        id: build
        run: |
          bazel build --config=release //...
          echo "artifacts=bazel-bin/app/app" >> $GITHUB_OUTPUT
      
      - uses: slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v1.9.0
        with:
          base64-subjects: ${{ steps.build.outputs.artifacts }}
```

### Signed Commits

```bash
# Configure GPG signing
git config --global user.signingkey YOUR_GPG_KEY
git config --global commit.gpgsign true

# Sign commits
git commit -S -m "feat: add feature"

# Verify signature
git log --show-signature
```

### Signed Artifacts

```bash
# Sign release artifacts with GPG
gpg --armor --detach-sign bazel-bin/app/app

# Verify signature
gpg --verify app.asc app
```

### Software Bill of Materials (SBOM)

```bash
# Generate SBOM with syft
syft packages bazel/app:image -o spdx-json > sbom.json

# Generate for Bazel project
bazel query 'deps(//...)' --output=xml > dependencies.xml
```

## Security Auditing

### Audit Logging

```bash
# Enable Bazel audit logging
build --experimental_execution_log_file=execution.log
build --experimental_execution_log_json_file=execution.json

# Analyze build actions
jq '.[] | select(.type == "action")' execution.json
```

### Access Logs

```bash
# Monitor repository access
gh api repos/:owner/:repo/traffic/clones
gh api repos/:owner/:repo/traffic/views

# Check recent commits
gh api repos/:owner/:repo/commits --jq '.[] | {author: .commit.author.name, date: .commit.author.date}'
```

### Security Audit Checklist

**Monthly audit:**

- [ ] Review all dependencies for vulnerabilities
- [ ] Check for exposed secrets in git history
- [ ] Review access controls and permissions
- [ ] Audit container images
- [ ] Review CI/CD security
- [ ] Check for outdated dependencies
- [ ] Review security alerts and incidents
- [ ] Update security documentation

### Automated Compliance Checks

```yaml
# .github/workflows/compliance.yml
name: Compliance Check

on:
  schedule:
    - cron: '0 0 1 * *'  # Monthly

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Check for secrets in history
        run: |
          git log --all --pretty=format: --name-only | sort -u | \
            xargs -I {} sh -c 'echo "Checking {}"; grep -i "password\|secret\|key" {} || true'
      
      - name: License compliance
        run: |
          pnpm licenses list
          pip-licenses
      
      - name: Dependency audit
        run: |
          pnpm audit
          pip-audit requirements/all.txt
          govulncheck ./...
```

## Incident Response

### Incident Response Plan

```mermaid
sequenceDiagram
    participant D as Detector
    participant I as Incident Team
    participant S as Security Team
    participant U as Users
    
    D->>I: Report Security Issue
    I->>I: Assess Severity
    
    alt Critical
        I->>S: Immediate Escalation
        S->>S: Emergency Response
        S->>U: Security Advisory
    else High
        I->>S: Urgent Notification
        S->>S: Priority Fix
    else Medium/Low
        I->>S: Standard Process
    end
    
    S->>S: Develop Fix
    S->>S: Test Fix
    S->>U: Deploy Fix
    S->>U: Post-mortem Report
```

### Security Incident Template

```markdown
# Security Incident Report

## Incident Summary
- **ID**: INC-2024-001
- **Severity**: Critical/High/Medium/Low
- **Status**: Open/In Progress/Resolved
- **Reporter**: Name
- **Date Reported**: 2024-01-15
- **Date Resolved**: TBD

## Description
Brief description of the security incident.

## Impact
- Systems affected
- Data exposed
- Users impacted

## Root Cause
Technical explanation of the vulnerability.

## Timeline
- 2024-01-15 09:00: Incident detected
- 2024-01-15 09:15: Security team notified
- 2024-01-15 10:00: Fix developed
- 2024-01-15 11:00: Fix deployed

## Resolution
Steps taken to resolve the incident.

## Prevention
Measures to prevent recurrence.

## Lessons Learned
Key takeaways from the incident.
```

### Emergency Contacts

```markdown
# Security Contacts

## Security Team
- Primary: security@example.com
- On-call: +1-555-0100

## Escalation Path
1. Security Engineer (on-call)
2. Security Lead
3. CISO
4. CTO

## External Contacts
- GitHub Security: security@github.com
- Cloud Provider Security: support@...
```

## Best Practices

### DO

- ✅ Use hermetic builds
- ✅ Pin all dependencies
- ✅ Scan for vulnerabilities daily
- ✅ Never commit secrets
- ✅ Use signed commits
- ✅ Enable branch protection
- ✅ Require code reviews
- ✅ Scan container images
- ✅ Maintain audit logs
- ✅ Have incident response plan

### DON'T

- ❌ Commit credentials
- ❌ Use wildcard dependency versions
- ❌ Skip security scans
- ❌ Run containers as root
- ❌ Disable security features for convenience
- ❌ Ignore security alerts
- ❌ Share credentials across environments
- ❌ Use outdated base images

## Troubleshooting

### False Positives

```bash
# Mark false positive in Dependabot
# Via web UI or .github/dependabot.yml

# Suppress in pip-audit
pip-audit --ignore-vuln VULN-ID

# Suppress in npm audit
pnpm audit --ignore-advisory ADVISORY-ID
```

### Leaked Secrets

```bash
# If secret leaked in git history

# 1. Rotate the secret immediately
./rotate-secret.sh API_KEY

# 2. Remove from git history
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch path/to/secret' \
  --prune-empty --tag-name-filter cat -- --all

# 3. Force push (coordinate with team!)
git push --force --all

# 4. Notify team
# 5. Update secret scanners
```

## Next Steps

- Review [Troubleshooting](./troubleshooting.md) guide
- Set up [Monitoring](./monitoring.md) for security events
- Configure [CI/CD](./ci-cd.md) security scans

---

**Back**: [Release Management](./release-management.md) | **Next**: [Troubleshooting](./troubleshooting.md)
