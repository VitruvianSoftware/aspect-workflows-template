# Dependency Management

This guide covers managing external dependencies in projects generated from the Aspect Workflows Template.

## Table of Contents

1. [Overview](#overview)
2. [Python Dependencies](#python-dependencies)
3. [JavaScript/TypeScript Dependencies](#javascripttypescript-dependencies)
4. [Go Dependencies](#go-dependencies)
5. [Java/Kotlin Dependencies](#javakotlin-dependencies)
6. [C/C++ Dependencies](#cc-dependencies)
7. [Rust Dependencies](#rust-dependencies)
8. [Updating Dependencies](#updating-dependencies)
9. [Security and Auditing](#security-and-auditing)

## Overview

The template provides hermetic dependency management for all supported languages, ensuring reproducible builds across different environments.

### Key Concepts

- **Lockfiles**: Pin exact dependency versions
- **Hermetic**: All dependencies fetched through Bazel
- **Reproducible**: Same inputs = same outputs
- **Cached**: Dependencies cached locally and remotely

### Dependency Flow

```shell
Developer adds dependency
    ↓
Update declaration file (pyproject.toml, package.json, etc.)
    ↓
Run update tool (repin, pnpm install, etc.)
    ↓
Lockfile updated
    ↓
Bazel fetches and caches dependencies
    ↓
Use in BUILD files
```

## Python Dependencies

### Adding Dependencies

1. **Add to pyproject.toml:**

    ```bash
    vim pyproject.toml
    ```

    ```toml
    [project]
    dependencies = [
        "requests>=2.31.0",
        "numpy>=1.24.0",
    ]

    [project.optional-dependencies]
    dev = [
        "pytest>=7.4.0",
        "mypy>=1.5.0",
    ]
    ```

2. **Update lockfiles:**

    ```bash
    ./tools/repin
    ```

    This script runs:

    - `bazel run //requirements:runtime.update` - Updates runtime dependencies
    - `bazel run //requirements:requirements.all.update` - Updates all dependencies
    - `bazel run //:gazelle_python_manifest.update` - Updates import mapping

3. **Update BUILD files:**

    ```bash
    bazel run gazelle
    ```

### Using Dependencies in BUILD Files

```python
py_binary(
    name = "my_app",
    srcs = ["main.py"],
    deps = [
        "@pip//requests",
        "@pip//numpy",
    ],
)
```

### Python Dependency Structure

```bash
requirements/
├── all.in          # All dependencies (aggregates others)
├── all.txt         # Locked all dependencies
├── runtime.txt     # Locked runtime dependencies
├── test.in         # Test dependencies
├── tools.in        # Development tools
└── BUILD           # Bazel targets
```

### Direct vs Transitive Dependencies

```toml
# pyproject.toml - Only declare direct dependencies
[project]
dependencies = [
    "requests>=2.31.0",  # Direct dependency
    # Don't list: urllib3, certifi, etc. (transitive)
]
```

### Version Constraints

```toml
# Exact version
"requests==2.31.0"

# Minimum version
"requests>=2.31.0"

# Compatible release
"requests~=2.31.0"  # >= 2.31.0, < 2.32.0

# Version range
"requests>=2.28.0,<3.0.0"
```

## JavaScript/TypeScript Dependencies

### Adding Dependencies

```bash
# Add runtime dependency
pnpm add lodash

# Add dev dependency
pnpm add -D @types/lodash

# Add specific version
pnpm add lodash@4.17.21

# Add from workspace
pnpm add @myorg/shared-lib
```

### Using Dependencies in BUILD Files

```python
load("@npm//:defs.bzl", "npm_link_all_packages")

npm_link_all_packages(name = "node_modules")

js_library(
    name = "my_lib",
    srcs = ["index.ts"],
    deps = [
        ":node_modules/lodash",
        ":node_modules/@types/lodash",
    ],
)
```

### Package Management Files

```
project/
├── package.json         # Dependency declarations
├── pnpm-lock.yaml      # Locked versions
├── pnpm-workspace.yaml # Workspace configuration
└── .npmrc              # npm/pnpm configuration
```

### Workspace Dependencies

For monorepo packages:

```json
// package.json
{
  "dependencies": {
    "@myorg/shared": "workspace:*"
  }
}
```

### Version Ranges

```json
{
  "dependencies": {
    "lodash": "^4.17.21",    // Compatible (^)
    "react": "~18.2.0",      // Patch updates (~)
    "vue": "3.3.4",          // Exact version
    "axios": ">=1.4.0"       // Minimum version
  }
}
```

### Updating JavaScript Dependencies

```bash
# Update all dependencies
pnpm update

# Update specific package
pnpm update lodash

# Update to latest
pnpm update lodash --latest

# Check for outdated
pnpm outdated
```

## Go Dependencies

### Adding Dependencies

1. **Add import to Go code:**

```go
package main

import (
    "github.com/spf13/cobra"
)
```

2. **Update go.mod:**

```bash
go mod tidy
```

3. **Update MODULE.bazel:**

```bash
bazel mod tidy
```

4. **Update BUILD files:**

```bash
bazel run gazelle
```

### Using Dependencies in BUILD Files

Gazelle automatically adds dependencies:

```python
go_library(
    name = "mylib",
    srcs = ["lib.go"],
    importpath = "github.com/myorg/myproject/mylib",
    deps = [
        "@com_github_spf13_cobra//:cobra",
    ],
)
```

### Go Module Files

```bash
project/
├── go.mod              # Module definition and dependencies
├── go.sum              # Checksums
├── MODULE.bazel        # Bazel module with use_repo
└── tools/
    └── tools.go        # Tools dependencies
```

### Version Selection

```go
// go.mod
module github.com/myorg/myproject

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/stretchr/testify v1.8.4
)
```

### Go Module Commands

```bash
# Add dependency (auto-detects version)
go get github.com/spf13/cobra

# Add specific version
go get github.com/spf13/cobra@v1.8.0

# Update dependency
go get -u github.com/spf13/cobra

# Update all dependencies
go get -u ./...

# Remove unused dependencies
go mod tidy

# View dependency graph
go mod graph
```

### Updating MODULE.bazel

After `go mod tidy`, update Bazel:

```bash
bazel mod tidy
```

This updates `use_repo()` declarations in MODULE.bazel:

```python
use_repo(
    go_deps,
    "com_github_spf13_cobra",
    "com_github_stretchr_testify",
)
```

## Java/Kotlin Dependencies

### Adding Dependencies

1. **Edit MODULE.bazel:**

```python
maven = use_extension("@rules_jvm_external//:extensions.bzl", "maven")
maven.install(
    artifacts = [
        "com.google.guava:guava:32.1.3-jre",
        "org.junit.jupiter:junit-jupiter:5.10.0",
    ],
    lock_file = "//:maven_install.json",
)
```

2. **Update lockfile:**

```bash
bazel run @unpinned_maven//:pin
```

### Using Dependencies in BUILD Files

```python
java_library(
    name = "mylib",
    srcs = ["MyLib.java"],
    deps = [
        "@maven//:com_google_guava_guava",
    ],
)

java_test(
    name = "mylib_test",
    srcs = ["MyLibTest.java"],
    deps = [
        ":mylib",
        "@maven//:org_junit_jupiter_junit_jupiter",
    ],
)
```

### Maven Dependency Format

```python
# Format: "group:artifact:version"
"com.google.guava:guava:32.1.3-jre"

# With classifier
"com.example:artifact:1.0:tests"

# With exclusions
maven.artifact(
    group = "com.example",
    artifact = "artifact",
    version = "1.0",
    exclusions = [
        "com.unwanted:dependency"
    ],
)
```

### Maven Repositories

```python
maven.install(
    artifacts = [...],
    repositories = [
        "https://repo1.maven.org/maven2",
        "https://jcenter.bintray.com",
    ],
)
```

### Checking for Updates

```bash
# Query outdated dependencies
bazel query @maven//:outdated

# Or use Maven tools
mvn versions:display-dependency-updates
```

## C/C++ Dependencies

### System Dependencies

C/C++ typically uses:

- System libraries (avoid for hermetic builds)
- Bazel modules
- http_archive for external projects

### Adding External Libraries

```python
# MODULE.bazel
bazel_dep(name = "abseil-cpp", version = "20230802.0")

# Or use http_archive
http_archive = use_repo_rule("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_google_googletest",
    urls = ["https://github.com/google/googletest/archive/v1.14.0.tar.gz"],
    strip_prefix = "googletest-1.14.0",
    sha256 = "...",
)
```

### Using in BUILD Files

```python
cc_library(
    name = "mylib",
    srcs = ["mylib.cc"],
    hdrs = ["mylib.h"],
    deps = [
        "@abseil-cpp//absl/strings",
        "@com_google_googletest//:gtest",
    ],
)
```

### Package Managers

For development, you can use:

- vcpkg
- Conan
- System package manager (apt, brew)

But prefer Bazel-managed dependencies for hermetic builds.

## Rust Dependencies

### Adding Dependencies

1. **Edit Cargo.toml:**

    ```toml
    [dependencies]
    serde = { version = "1.0", features = ["derive"] }
    tokio = { version = "1.35", features = ["full"] }

    [dev-dependencies]
    criterion = "0.5"
    ```

2. **Dependencies are fetched by Cargo/Bazel automatically**

### Using in BUILD Files

```python
load("@rules_rust//rust:defs.bzl", "rust_binary")

rust_binary(
    name = "myapp",
    srcs = ["src/main.rs"],
    deps = [
        # Cargo dependencies available automatically
    ],
)
```

### Cargo Features

```toml
[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.35", features = ["rt-multi-thread"] }
```

### Development Commands

```bash
# Add dependency
cargo add serde

# Update dependencies
cargo update

# Check for updates
cargo outdated
```

## Updating Dependencies

### Automated Updates with Renovate

The template includes `renovate.json`:

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "config:recommended",
    ":automergeMinor"
  ]
}
```

Renovate automatically:

- Creates PRs for dependency updates
- Groups related updates
- Auto-merges minor/patch updates
- Respects version constraints

### Manual Update Strategies

**Check for updates:**

```bash
# Python
pip list --outdated

# JavaScript
pnpm outdated

# Go
go list -u -m all

# Java
# Use Maven version plugin or Renovate
```

**Update all dependencies:**

```bash
# Python
vim pyproject.toml  # Update versions
./tools/repin

# JavaScript
pnpm update --latest

# Go
go get -u ./...
go mod tidy
bazel mod tidy

# Java
vim MODULE.bazel  # Update versions
bazel run @unpinned_maven//:pin
```

### Update Testing

After updating dependencies:

```bash
# Build everything
bazel build //...

# Run all tests
bazel test //...

# Check for issues
aspect lint //...
```

## Security and Auditing

### Checking for Vulnerabilities

**Python:**

```bash
pip-audit requirements/all.txt
```

**JavaScript:**

```bash
pnpm audit
pnpm audit --fix  # Auto-fix if possible
```

**Go:**

```bash
go list -json -m all | nancy sleuth
```

### Dependency Scanning in CI

```yaml
# GitHub Actions example
- name: Security Audit
  run: |
    pnpm audit --audit-level=moderate
    pip-audit requirements/all.txt
```

### License Compliance

Check dependency licenses:

```bash
# Python
pip-licenses

# JavaScript
pnpm licenses list

# Go
go-licenses check ./...
```

### Pinning for Security

Always use lockfiles:

- ✅ `pnpm-lock.yaml` for JavaScript
- ✅ `requirements/*.txt` for Python
- ✅ `go.sum` for Go
- ✅ `maven_install.json` for Java

## Best Practices

**DO**:

- ✅ Use lockfiles for all languages
- ✅ Review dependency updates before merging
- ✅ Run tests after updating dependencies
- ✅ Keep dependencies up to date regularly
- ✅ Use automated tools like Renovate
- ✅ Audit dependencies for security
- ✅ Document why specific versions are pinned
- ✅ Minimize dependency count

**DON'T**:

- ❌ Commit with outdated lockfiles
- ❌ Use wildcard versions in production
- ❌ Skip security updates
- ❌ Add unnecessary dependencies
- ❌ Ignore Renovate PRs indefinitely
- ❌ Mix package managers (use project's choice)
- ❌ Vendor dependencies without good reason

## Troubleshooting

### Python "Module not found"

```bash
# Ensure dependency is in pyproject.toml
grep "package-name" pyproject.toml

# Regenerate lockfiles
./tools/repin

# Update BUILD files
bazel run gazelle

# Verify in lockfile
grep "package-name" requirements/all.txt
```

### JavaScript "Cannot find module"

```bash
# Install the package
pnpm add package-name

# Verify in package.json
grep "package-name" package.json

# Rebuild
bazel clean
bazel build //...
```

### Go "Package not found"

```bash
# Update go.mod
go mod tidy

# Update MODULE.bazel
bazel mod tidy

# Update BUILD files
bazel run gazelle
```

### Version Conflicts

```bash
# Python - check for conflicts
pip-compile --dry-run requirements/all.in

# JavaScript - check for conflicts
pnpm why package-name

# Go - view selected versions
go mod why -m github.com/some/package
```

## Next Steps

- Review [Troubleshooting Guide](./troubleshooting.md)
- Learn about [Building and Testing](./building-testing.md)
- Check [FAQ](../faq.md) for common questions

---

**Back**: [Formatting and Linting](./formatting-linting.md) | **Up**: [User Guide](./README.md)
