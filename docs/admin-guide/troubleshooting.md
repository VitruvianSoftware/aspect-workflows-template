# Administrative Troubleshooting

This guide covers common administrative issues and their solutions for projects using the Aspect Workflows Template.

## Table of Contents

1. [Build Infrastructure Issues](#build-infrastructure-issues)
2. [Cache Problems](#cache-problems)
3. [CI/CD Failures](#cicd-failures)
4. [Dependency Issues](#dependency-issues)
5. [Performance Problems](#performance-problems)
6. [Access and Permissions](#access-and-permissions)
7. [Deployment Issues](#deployment-issues)

## Build Infrastructure Issues

### Remote Cache Not Working

**Symptoms:**

- Cache hit rate near 0%
- Builds slower in CI than locally
- Cache connection errors

**Diagnosis:**

```bash
# Test cache connectivity
curl -I https://your-cache-server

# Check authentication
gcloud auth list  # For GCS
aws sts get-caller-identity  # For S3

# Enable verbose caching
bazel build --remote_cache=https://... --remote_print_execution_messages //...
```

**Solutions:**

1. **Check authentication:**

```bash
# GCS
gcloud auth application-default login

# AWS
aws configure

# Add to CI
# GitHub Actions
- uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_CREDENTIALS }}
```

2. **Verify cache configuration:**

```bash
# .bazelrc
build --remote_cache=https://storage.googleapis.com/my-cache
build --remote_upload_local_results=true
build --google_default_credentials

# Check it's being used
bazel info | grep remote_cache
```

3. **Check permissions:**

```bash
# GCS - verify bucket permissions
gsutil iam get gs://my-bazel-cache

# Should include: storage.objects.create, storage.objects.get
```

### Bazel Version Mismatch

**Symptoms:**

- "Incompatible Bazel version" errors
- BUILD files not parsing
- Features not available

**Solution:**

```bash
# Check current version
bazel version

# Update .bazelversion
echo "7.0.2" > .bazelversion

# Update CI
# .github/workflows/ci.yml
- uses: bazel-contrib/setup-bazel@0.8.0
  with:
    bazelisk-version: 1.19.0

# Clean and rebuild
bazel clean --expunge
bazel build //...
```

### Disk Space Issues

**Symptoms:**

- "No space left on device"
- Build failures mid-execution
- Cache errors

**Diagnosis:**

```bash
# Check disk usage
df -h

# Check Bazel cache size
du -sh ~/.cache/bazel
du -sh $(bazel info output_base)

# Find large directories
du -sh ~/.cache/bazel/* | sort -rh | head -20
```

**Solutions:**

```bash
# Clean Bazel cache
bazel clean --expunge

# Remove old artifacts
find ~/.cache/bazel -type f -mtime +30 -delete

# Configure cache limits
# .bazelrc
build --disk_cache_gc_idle_time=5m
build --experimental_disk_cache_gc_max_size=50GB

# Move cache to larger disk
export BAZEL_CACHE_DIR=/mnt/large-disk/bazel-cache
bazel build --disk_cache=$BAZEL_CACHE_DIR //...
```

### Out of Memory

**Symptoms:**

- "OutOfMemoryError" during builds
- JVM crashes
- Build hangs

**Diagnosis:**

```bash
# Check memory usage
free -h

# Monitor during build
watch -n 1 free -h

# Check Bazel memory settings
bazel info | grep heap
```

**Solutions:**

```bash
# Limit Bazel memory
# .bazelrc
startup --host_jvm_args=-Xmx4g

# Reduce parallelism
build --jobs=4
build --local_ram_resources=HOST_RAM*.67

# For specific languages
build --worker_max_instances=4

# Disable workers for memory-intensive rules
build --strategy=KotlinCompile=local
```

## Cache Problems

### Low Cache Hit Rate

**Symptoms:**

- Most actions rebuild from scratch
- CI builds take full time
- Cache metrics show <50% hit rate

**Diagnosis:**

```bash
# Profile build with cache stats
bazel build --profile=profile.json //...
bazel analyze-profile profile.json | grep "cache hit"

# Check execution log
bazel build --execution_log_json_file=exec.json //...
jq '[.[] | select(.remoteCacheHit == true)] | length' exec.json
jq '[.[] | select(.remoteCacheHit == false)] | length' exec.json
```

**Solutions:**

1. **Fix non-deterministic builds:**

    ```python
    # Avoid timestamp-dependent rules
    genrule(
        name = "bad",
        outs = ["out.txt"],
        cmd = "date > $@",  # BAD: non-deterministic
    )

    # Use deterministic alternatives
    genrule(
        name = "good",
        outs = ["out.txt"],
        cmd = "echo 'fixed content' > $@",  # GOOD: deterministic
    )
    ```

2. **Check for local environment dependencies:**

    ```bash
    # Find actions depending on local environment
    bazel build --explain=explain.log //...
    grep "local environment" explain.log
    ```

3. **Verify cache key stability:**

    ```bash
    # Build twice, should be identical
    bazel build //... --profile=p1.json
    bazel clean
    bazel build //... --profile=p2.json

    # Compare profiles
    diff <(jq '.[] | .name' p1.json | sort) <(jq '.[] | .name' p2.json | sort)
    ```

### Cache Corruption

**Symptoms:**

- Inconsistent build results
- Random build failures
- "Checksum mismatch" errors

**Solutions:**

```bash
# Clear local cache
bazel clean --expunge

# Clear specific cache entries
bazel clean --expunge_async

# Verify remote cache integrity
# GCS example
gsutil -m ls -r gs://my-bazel-cache | wc -l

# For persistent issues, consider cache rotation
# Increment cache prefix
build --remote_cache=https://storage.googleapis.com/my-cache/v2
```

## CI/CD Failures

### Flaky Tests in CI

**Symptoms:**

- Tests pass locally, fail in CI
- Tests fail randomly
- Different results on retry

**Diagnosis:**

```bash
# Run test multiple times
bazel test --runs_per_test=10 //path/to:flaky_test

# Check for timing issues
bazel test --test_timeout=300 //path/to:flaky_test

# Enable verbose output
bazel test --test_output=all //path/to:flaky_test
```

**Solutions:**

```python
# Mark as flaky temporarily
py_test(
    name = "flaky_test",
    srcs = ["flaky_test.py"],
    flaky = True,  # Allows retries
)

# Add timeout
py_test(
    name = "slow_test",
    srcs = ["slow_test.py"],
    timeout = "long",  # short, moderate, long, eternal
)

# Fix timing issues in test
import time
import unittest

class MyTest(unittest.TestCase):
    def setUp(self):
        # Add retry logic
        for _ in range(3):
            try:
                self.resource = connect_to_resource()
                break
            except ConnectionError:
                time.sleep(1)
```

### CI Running Out of Disk Space

**Symptoms:**

- "No space left" errors in CI
- Builds fail mid-execution
- Cache writes fail

**Solutions:**

```yaml
# GitHub Actions - clean before build
- name: Free disk space
  run: |
    docker system prune -a -f
    sudo rm -rf /usr/share/dotnet
    sudo rm -rf /opt/ghc
    sudo rm -rf /usr/local/share/boost
    df -h

# Use larger runners
runs-on: ubuntu-latest-8-cores  # More disk space

# Or use cleanup action
- name: Cleanup
  uses: jlumbroso/free-disk-space@main
  with:
    tool-cache: true
    android: true
    dotnet: true
```

### Authentication Failures

**Symptoms:**

- "Permission denied" in CI
- Registry push fails
- Remote cache authentication errors

**Solutions:**

```yaml
# GitHub Actions - proper authentication
- name: Authenticate to GCP
  uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_CREDENTIALS }}

- name: Configure Docker
  run: |
    gcloud auth configure-docker gcr.io

# AWS authentication
- name: Configure AWS
  uses: aws-actions/configure-aws-credentials@v4
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    aws-region: us-east-1
```

### Timeout Issues

**Symptoms:**

- Jobs timeout before completion
- Long-running tests killed
- Build interrupted

**Solutions:**

```yaml
# Increase job timeout
jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 60  # Default is 360

# Split into multiple jobs
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build
        run: bazel build //src/...
  
  test:
    needs: build
    runs-on: ubuntu-latest
    strategy:
      matrix:
        shard: [0, 1, 2, 3]
    steps:
      - name: Test
        run: |
          bazel test //... \
            --test_shard_count=4 \
            --test_shard_index=${{ matrix.shard }}
```

## Dependency Issues

### Renovate Not Creating PRs

**Symptoms:**

- No dependency update PRs
- Renovate not running
- Configuration errors

**Diagnosis:**

```bash
# Check Renovate logs
gh api repos/:owner/:repo/actions/workflows | jq '.workflows[] | select(.name == "Renovate")'

# Validate config
npx --package renovate -c renovate-config-validator

# Check rate limits
gh api rate_limit
```

**Solutions:**

```json
// renovate.json - fix common issues
{
  "extends": ["config:recommended"],
  "schedule": ["after 9am and before 5pm every weekday"],
  "prConcurrentLimit": 5,  // Avoid overwhelming team
  "prHourlyLimit": 0,       // No hourly limit
  "recreateClosed": false,  // Don't recreate closed PRs
  "ignorePaths": [
    "**/node_modules/**",
    "**/vendor/**"
  ]
}
```

### Dependency Resolution Conflicts

**Symptoms:**

- "Could not resolve dependencies"
- Version conflicts
- Build fails after dependency update

**Solutions:**

```bash
# Python - check conflicts
pip-compile --dry-run requirements/all.in

# JavaScript - resolve conflicts
pnpm why conflicting-package
pnpm update conflicting-package

# Override if necessary
# package.json
{
  "pnpm": {
    "overrides": {
      "vulnerable-package": "^2.0.0"
    }
  }
}

# Go - check why dependency is needed
go mod why github.com/some/package
go mod graph | grep package-name

# Update to compatible versions
go get -u github.com/some/package@compatible-version
```

### Lockfile Out of Sync

**Symptoms:**

- "Lockfile is out of date"
- Build fails with dependency errors
- Different behavior locally vs CI

**Solutions:**

```bash
# Python
./tools/repin

# JavaScript
pnpm install
git add pnpm-lock.yaml

# Go
go mod tidy
bazel mod tidy
git add go.mod go.sum MODULE.bazel

# Verify in CI
# Add to CI workflow
- name: Check lockfiles
  run: |
    ./tools/repin
    git diff --exit-code
```

## Performance Problems

### Slow Builds

**Symptoms:**

- Builds take >10 minutes
- No improvement with cache
- High CPU usage

**Diagnosis:**

```bash
# Profile build
bazel build --profile=profile.json //...
bazel analyze-profile profile.json

# Generate HTML report
bazel analyze-profile --html profile.json > profile.html
open profile.html

# Find critical path
grep "CRITICAL PATH" profile.json

# Check for rebuilds
bazel build --explain=explain.log //...
grep "action '[^']*' not in cache" explain.log
```

**Solutions:**

1. **Enable remote caching:**

    ```bash
    build --remote_cache=https://cache.example.com
    ```

2. **Increase parallelism:**

    ```bash
    build --jobs=auto
    build --local_cpu_resources=HOST_CPUS*.75
    ```

3. **Split large targets:**

    ```python
    # Instead of one large library
    py_library(
        name = "big_lib",
        srcs = glob(["*.py"]),  # 100+ files
    )

    # Split into smaller libraries
    py_library(
        name = "module_a",
        srcs = glob(["module_a/**/*.py"]),
    )

    py_library(
        name = "module_b",
        srcs = glob(["module_b/**/*.py"]),
    )
    ```

### High Memory Usage

**Symptoms:**

- OOM errors
- Swapping
- System slowdown

**Solutions:**

```bash
# Reduce worker instances
build --worker_max_instances=2

# Limit action memory
build --local_ram_resources=HOST_RAM*.67

# Use remote execution for heavy tasks
build --strategy=JavaCompile=remote
```

### Network Bottlenecks

**Symptoms:**

- Slow dependency downloads
- Remote cache slow
- Timeouts

**Solutions:**

```bash
# Increase timeouts
build --remote_timeout=3600

# Use closer cache region
build --remote_cache=https://cache-us-east1.example.com

# Enable HTTP/2
build --remote_cache=grpcs://cache.example.com

# Check network
ping cache.example.com
curl -w "@curl-format.txt" -o /dev/null https://cache.example.com
```

## Access and Permissions

### Team Members Can't Push

**Symptoms:**

- "Permission denied" on push
- Can't create branches
- Can't merge PRs

**Solutions:**

```bash
# Check user permissions
gh api repos/:owner/:repo/collaborators/:username/permission

# Add user to team
gh api -X PUT repos/:owner/:repo/collaborators/:username \
  -f permission=push

# Check branch protection
gh api repos/:owner/:repo/branches/main/protection

# Ensure user is in correct team
gh api orgs/:org/teams/:team/members
```

### CI Can't Access Secrets

**Symptoms:**

- Environment variables empty in CI
- "Secret not found"
- Authentication failures

**Solutions:**

```bash
# Verify secret exists
gh secret list

# Add secret
gh secret set SECRET_NAME

# Check workflow permissions
# .github/workflows/ci.yml
permissions:
  contents: read
  packages: write
  id-token: write

# Verify secret is accessible in workflow
- name: Debug secrets
  run: |
    echo "Secret length: ${#SECRET_VALUE}"
  env:
    SECRET_VALUE: ${{ secrets.SECRET_VALUE }}
```

### Dependabot Can't Access Private Registry

**Symptoms:**

- Dependabot PRs fail
- "Authentication required"
- Can't fetch private dependencies

**Solutions:**

```yaml
# .github/dependabot.yml
version: 2
registries:
  npm-private:
    type: npm-registry
    url: https://npm.pkg.github.com
    token: ${{ secrets.NPM_TOKEN }}

updates:
  - package-ecosystem: "npm"
    directory: "/"
    registries:
      - npm-private
    schedule:
      interval: "weekly"
```

## Deployment Issues

### Container Push Fails

**Symptoms:**

- "denied: permission denied"
- Registry authentication errors
- Image push timeout

**Solutions:**

```bash
# Verify registry authentication
docker login gcr.io
docker login ghcr.io -u USERNAME -p $GITHUB_TOKEN

# Check image exists
bazel run //app:image
docker images | grep app

# Push manually to debug
bazel run //app:image.push --//tools:registry=gcr.io/project

# Check registry permissions
gcloud projects get-iam-policy PROJECT_ID
```

### Deployment Rollback Needed

**Symptoms:**

- New version has critical bug
- Service degradation
- Need to restore previous version

**Solutions:**

```bash
# Quick rollback script
#!/bin/bash
PREVIOUS_VERSION=$1

echo "Rolling back to $PREVIOUS_VERSION..."

# Kubernetes
kubectl rollout undo deployment/app

# Or specific revision
kubectl rollout undo deployment/app --to-revision=2

# Cloud Run
gcloud run services update-traffic app \
  --to-revisions=$PREVIOUS_VERSION=100

# Verify
kubectl rollout status deployment/app

echo "Rollback complete"
```

### Health Check Failures

**Symptoms:**

- Service marked unhealthy
- Load balancer removes instances
- Deploy fails health check

**Solutions:**

```yaml
# Adjust health check timing
# Kubernetes
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30  # Increase if slow startup
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

# Check logs
kubectl logs deployment/app
kubectl describe pod app-xxx

# Test health endpoint
curl http://localhost:8080/health
```

## Best Practices

### DO

- ✅ Monitor build and cache metrics
- ✅ Set up alerts for failures
- ✅ Document common issues
- ✅ Keep runbooks updated
- ✅ Test fixes in staging first
- ✅ Maintain emergency contacts
- ✅ Regular backup verification

### DON'T

- ❌ Ignore intermittent failures
- ❌ Skip root cause analysis
- ❌ Apply fixes without testing
- ❌ Disable security features to "fix" issues
- ❌ Leave debug logging in production
- ❌ Forget to document solutions

## Emergency Procedures

### Build System Down

```markdown
# Emergency Response: Build System Down

1. **Assess Impact**
   - Which services affected?
   - Can developers build locally?
   - Is production at risk?

2. **Immediate Actions**
   - Post status update
   - Switch to backup cache
   - Enable local-only builds

3. **Investigation**
   - Check system status
   - Review logs
   - Identify root cause

4. **Resolution**
   - Apply fix
   - Verify functionality
   - Post-mortem analysis
```

## Next Steps

- Set up [Monitoring](./monitoring.md) to catch issues early
- Review [Security](./security.md) for security-related issues
- Configure [CI/CD](./ci-cd.md) properly

---

**Back**: [Security](./security.md) | **Next**: [Monitoring](./monitoring.md)
