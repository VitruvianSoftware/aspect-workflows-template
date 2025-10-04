# Administrator Guide

Welcome to the Aspect Workflows Template Administrator Guide! This guide is for maintainers, DevOps engineers, and administrators who manage and maintain projects based on this template.

## Table of Contents

1. [Maintenance Overview](./maintenance.md) - Regular maintenance tasks
2. [Dependency Management](./dependency-management.md) - Updating dependencies
3. [CI/CD Configuration](./ci-cd.md) - Continuous integration setup
4. [Release Management](./release-management.md) - Version and release processes
5. [Security](./security.md) - Security best practices
6. [Troubleshooting](./troubleshooting.md) - Common administrative issues
7. [Monitoring](./monitoring.md) - Build and performance monitoring

## Overview

As an administrator of projects using the Aspect Workflows Template, you're responsible for:

- Keeping dependencies up to date
- Managing CI/CD pipelines
- Ensuring build reproducibility
- Monitoring build performance
- Handling security updates
- Managing tool versions
- Supporting development teams

## Quick Reference

### Common Administrative Tasks

```bash
# Update all dependencies
./tools/repin

# Update Bazel dependencies
bazel mod tidy

# Update npm dependencies
pnpm update

# Update Go dependencies
go mod tidy

# Test entire project
aspect test //...

# Check for outdated dependencies
bazel query @maven//:outdated  # Java
pnpm outdated                  # JavaScript
```

## Key Responsibilities

### 1. Dependency Management

Maintain secure, up-to-date dependencies across all languages:

```mermaid
graph TB
    Admin[Administrator] --> Check[Check for Updates]
    Check --> Security[Security Patches]
    Check --> Features[New Features]
    Check --> Bugs[Bug Fixes]
    
    Security --> Critical[Critical - Immediate]
    Features --> Scheduled[Scheduled - Sprint]
    Bugs --> Evaluated[Evaluated - Case by Case]
    
    Critical --> Update[Update Dependencies]
    Scheduled --> Update
    Evaluated --> Update
    
    Update --> Test[Test Changes]
    Test --> Pass{Tests Pass?}
    Pass -->|Yes| Deploy[Deploy Update]
    Pass -->|No| Fix[Fix Issues]
    Fix --> Test
```

### 2. Build System Maintenance

Keep the Bazel build system healthy:

- Monitor build times
- Optimize slow builds
- Manage remote cache
- Update Bazel versions
- Review and update .bazelrc configurations

### 3. CI/CD Pipeline Management

Ensure reliable continuous integration:

- Configure build pipelines
- Set up test automation
- Manage secrets and credentials
- Configure deployment processes
- Monitor pipeline health

### 4. Security Management

Maintain security posture:

- Monitor security advisories
- Apply security patches promptly
- Conduct security audits
- Manage access controls
- Review dependency vulnerabilities

## Infrastructure Components

### Build Infrastructure

```mermaid
graph TB
    subgraph "Development"
        Dev[Developer Workstation]
        Local[Local Bazel Cache]
    end
    
    subgraph "CI/CD"
        CI[CI Server]
        Runners[Build Runners]
    end
    
    subgraph "Shared Resources"
        Remote[Remote Cache]
        Registry[Container Registry]
        Artifacts[Artifact Storage]
    end
    
    Dev --> Local
    Dev --> Remote
    CI --> Runners
    Runners --> Remote
    Runners --> Registry
    Runners --> Artifacts
```

### Recommended Infrastructure

1. **Remote Build Cache**
   - Google Cloud Storage
   - AWS S3
   - Azure Blob Storage
   - Self-hosted cache server

2. **Container Registry**
   - Docker Hub
   - Google Artifact Registry
   - AWS ECR
   - GitHub Container Registry

3. **CI/CD Platform**
   - GitHub Actions
   - GitLab CI
   - CircleCI
   - Jenkins

## Maintenance Schedule

### Daily Tasks

- ✅ Monitor CI/CD pipeline status
- ✅ Review failed builds
- ✅ Check security alerts

### Weekly Tasks

- ✅ Review dependency updates (Renovate PRs)
- ✅ Check build performance metrics
- ✅ Review and merge approved PRs
- ✅ Update documentation as needed

### Monthly Tasks

- ✅ Update Bazel version (if new release)
- ✅ Review and update toolchain versions
- ✅ Audit security vulnerabilities
- ✅ Review cache hit rates
- ✅ Performance optimization review

### Quarterly Tasks

- ✅ Major dependency upgrades
- ✅ Review and update CI/CD configuration
- ✅ Security audit
- ✅ Disaster recovery testing
- ✅ Documentation review and updates

## Monitoring and Metrics

### Key Metrics to Track

1. **Build Performance**
   - Build time trends
   - Cache hit rates
   - Test execution time
   - Network bandwidth usage

2. **CI/CD Health**
   - Pipeline success rate
   - Average build duration
   - Queue times
   - Resource utilization

3. **Dependency Health**
   - Outdated dependencies count
   - Security vulnerabilities
   - License compliance
   - Update frequency

4. **Developer Experience**
   - Time to first build
   - Build failure rate
   - Setup documentation effectiveness
   - Support ticket trends

### Monitoring Dashboard Example

```mermaid
graph LR
    subgraph "Metrics Collection"
        Bazel[Bazel Build Events]
        CI[CI Metrics]
        Tests[Test Results]
    end
    
    subgraph "Dashboard"
        BuildTime[Build Time Trends]
        CacheHit[Cache Hit Rates]
        TestPass[Test Pass Rates]
        Security[Security Alerts]
    end
    
    Bazel --> BuildTime
    Bazel --> CacheHit
    CI --> BuildTime
    Tests --> TestPass
    Security --> Alerts[Alert System]
```

## Backup and Disaster Recovery

### What to Back Up

1. **Source Code**: Already in Git
2. **Build Cache**: Can be rebuilt (optional backup)
3. **Credentials and Secrets**: Secure backup required
4. **CI/CD Configuration**: Version controlled
5. **Documentation**: Version controlled

### Recovery Procedures

```mermaid
flowchart TB
    Disaster[Disaster Event] --> Assess[Assess Impact]
    Assess --> Type{Type?}
    
    Type -->|Lost Cache| RebuildCache[Rebuild from Source]
    Type -->|Lost Credentials| RestoreCreds[Restore from Secure Backup]
    Type -->|Lost CI Config| RestoreCI[Restore from VCS]
    Type -->|Lost Source| RestoreGit[Restore from Git Remote]
    
    RebuildCache --> Verify[Verify Builds]
    RestoreCreds --> Verify
    RestoreCI --> Verify
    RestoreGit --> Verify
    
    Verify --> Test[Run Full Test Suite]
    Test --> Pass{Pass?}
    Pass -->|Yes| Normal[Resume Normal Operations]
    Pass -->|No| Debug[Debug Issues]
    Debug --> Test
```

## Automated Dependency Updates

### Renovate Configuration

The template includes `renovate.json` for automated dependency updates:

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "config:recommended",
    ":automergeMinor"
  ]
}
```

**Configuration Options:**

- Auto-merge minor updates
- Group related updates
- Schedule updates
- Custom PR descriptions

### Dependency Update Flow

```mermaid
sequenceDiagram
    participant R as Renovate Bot
    participant G as GitHub
    participant CI as CI/CD
    participant M as Maintainer
    
    R->>G: Check for updates
    G->>R: Return current versions
    R->>R: Compare with latest
    R->>G: Create PR for updates
    G->>CI: Trigger CI build
    CI->>CI: Run tests
    CI->>G: Report results
    G->>M: Notify of PR
    M->>G: Review and approve
    G->>G: Merge PR (auto or manual)
```

## Tool Version Management

### Multitool Lockfile

Tools are managed via `tools/tools.lock.json`:

```json
{
  "gazelle": {
    "version": "0.45.0",
    "checksums": {...}
  }
}
```

**Updating Tools:**

1. Edit `tools/BUILD` to reference new version
2. Run `bazel run @multitool//:update_lock`
3. Test with the new version
4. Commit the updated lockfile

## Performance Optimization

### Build Performance Tuning

1. **Enable Remote Caching**

   ```bash
   # Add to .bazelrc
   build --remote_cache=https://your-cache-server
   ```

2. **Optimize Memory Usage**

   ```bash
   # Add to .bazelrc
   build --jobs=auto
   build --local_ram_resources=HOST_RAM*.8
   ```

3. **Profile Builds**

   ```bash
   bazel build --profile=profile.json //...
   bazel analyze-profile profile.json
   ```

4. **Use Build Event Protocol**

   ```bash
   bazel build --build_event_json_file=bep.json //...
   ```

### Cache Optimization

```mermaid
graph TB
    Build[Build Action] --> Check{In Cache?}
    Check -->|Yes| Hit[Cache Hit - Fast]
    Check -->|No| Miss[Cache Miss - Build]
    Miss --> Execute[Execute Action]
    Execute --> Store[Store in Cache]
    Store --> Output[Output Artifacts]
    Hit --> Output
    
    Remote[Remote Cache] --> Check
    Local[Local Cache] --> Check
```

## Security Best Practices

### 1. Dependency Security

- Enable Renovate security updates
- Use `bazel mod graph` to audit dependencies
- Scan container images for vulnerabilities
- Keep toolchains updated

### 2. Credential Management

- Never commit secrets to version control
- Use CI/CD secret management
- Rotate credentials regularly
- Use service accounts for automation

### 3. Access Control

- Implement least privilege access
- Review access permissions quarterly
- Use branch protection rules
- Require code reviews for sensitive changes

## Support and Escalation

### Support Tiers

```mermaid
graph TB
    Issue[Issue Reported] --> Triage[Triage]
    Triage --> Level{Severity?}
    
    Level -->|P0 - Critical| Immediate[Immediate Response]
    Level -->|P1 - High| SameDay[Same Day Response]
    Level -->|P2 - Medium| NextDay[Next Day Response]
    Level -->|P3 - Low| Scheduled[Scheduled Response]
    
    Immediate --> Resolve[Resolve]
    SameDay --> Resolve
    NextDay --> Resolve
    Scheduled --> Resolve
    
    Resolve --> Document[Document Solution]
    Document --> Close[Close Issue]
```

### Escalation Path

1. **Level 1**: Development team
2. **Level 2**: Build system administrators
3. **Level 3**: Template maintainers
4. **Level 4**: Aspect Build support (enterprise customers)

## Resources

### Documentation

- [Bazel Documentation](https://bazel.build/)
- [Aspect CLI Docs](https://docs.aspect.build/)
- [rules_lint](https://github.com/aspect-build/rules_lint)
- [Renovate Docs](https://docs.renovatebot.com/)

### Community

- **Bazel Slack**: #aspect-build channel
- **GitHub Issues**: Report bugs and feature requests
- **GitHub Discussions**: Ask questions and share knowledge

### Training

- Aspect Build workshops
- Bazel training courses
- Internal knowledge sharing sessions

---

**Next**: [Maintenance Overview](./maintenance.md) | **Up**: [Documentation Home](../README.md)
