# Quick Reference Guide

A quick reference for common tasks and commands when working with the Aspect Workflows Template.

## Template Generation

```bash
# Interactive generation
scaffold new github.com/BlueCentre/aspect-workflows-template

# Using a preset
scaffold new --preset=py --no-prompt github.com/BlueCentre/aspect-workflows-template

# Custom output directory
scaffold new --output-dir=./my-project github.com/BlueCentre/aspect-workflows-template
```

### Available Presets

#### Direct Generation Presets

| Preset | Languages | Features |
|--------|-----------|----------|
| `minimal` | None | Basic structure |
| `py` | Python | Lint, codegen |
| `go` | Go | Lint, codegen, OCI |
| `js` | JavaScript/TypeScript | Lint, codegen |
| `java` | Java | Lint |
| `kotlin` | Kotlin | Lint |
| `cpp` | C/C++ | Lint |
| `rust` | Rust | Lint |
| `shell` | Shell/Bash | Lint |
| `kitchen-sink` | All | All features |

#### Backstage Template Presets

For generating Backstage software templates:

| Preset | Languages | Features |
|--------|-----------|----------|
| `backstage-minimal` | None | Backstage template structure |
| `backstage-py` | Python | Backstage + Lint, codegen |
| `backstage-go` | Go | Backstage + Lint, codegen, OCI |
| `backstage-js` | JavaScript/TypeScript | Backstage + Lint, codegen |
| `backstage-java` | Java | Backstage + Lint |
| `backstage-kotlin` | Kotlin | Backstage + Lint |
| `backstage-cpp` | C/C++ | Backstage + Lint |
| `backstage-rust` | Rust | Backstage + Lint |
| `backstage-shell` | Shell/Bash | Backstage + Lint |
| `backstage-kitchen-sink` | All | Backstage + All features |

## Initial Setup

```bash
# Navigate to project
cd <project-name>

# Enable direnv
direnv allow

# Setup development environment
bazel run //tools:bazel_env

# Enable git hooks
git config core.hooksPath githooks
```

## Building

```bash
# Build everything
bazel build //...

# Build specific target
bazel build //src/app:main

# Build with optimizations
bazel build -c opt //...

# Build for release (with stamping)
bazel build --config=release //...
```

## Testing

```bash
# Run all tests
aspect test //...
# or
bazel test //...

# Run specific test
aspect test //src/app:test

# Run with output
bazel test --test_output=all //src/app:test

# Run only unit tests
bazel test --test_tag_filters=unit //...

# Run tests without cache
bazel test --cache_test_results=no //...
```

## Code Generation (Gazelle)

```bash
# Generate/update BUILD files
bazel run gazelle

# Update Go module dependencies
bazel mod tidy

# Update Python manifest
bazel run //:gazelle_python_manifest.update
```

## Formatting

```bash
# Format all files
format

# Format specific files
format src/app/main.py src/lib/utils.py

# Check formatting without fixing
format --check

# Format only Bazel files
bazel run @buildifier_prebuilt//:buildifier -- -r .
```

## Linting

```bash
# Lint everything
aspect lint //...

# Lint specific package
aspect lint //src/app:all

# Apply automatic fixes
aspect lint --fix //...
```

## Dependency Management

### Python

```bash
# 1. Edit pyproject.toml
vim pyproject.toml

# 2. Update lockfiles
./tools/repin

# 3. Update BUILD files
bazel run gazelle
```

### JavaScript/TypeScript

```bash
# Add dependency
pnpm add <package>

# Add dev dependency
pnpm add -D <package>

# Update all dependencies
pnpm update
```

### Go

```bash
# 1. Add import to code
# 2. Update go.mod
go mod tidy

# 3. Update MODULE.bazel
bazel mod tidy

# 4. Update BUILD files
bazel run gazelle
```

### Java

```bash
# 1. Add to MODULE.bazel maven.install artifacts
vim MODULE.bazel

# 2. Repin dependencies
bazel run @unpinned_maven//:pin
```

## Container Images (if enabled)

```bash
# Build image
bazel build //src/app:image

# Load into Docker
bazel run //src/app:image.load

# Push to registry
bazel run //src/app:image.push

# Build for multiple platforms
bazel build \
  --platforms=//tools/platforms:linux_amd64 \
  --platforms=//tools/platforms:linux_arm64 \
  //src/app:image
```

## Watch Mode

```bash
# Auto-rebuild on changes
ibazel build //src/app:all

# Auto-test on changes
ibazel test //src/app:test

# Auto-run application
ibazel run //src/app:main
```

## Debugging

```bash
# Verbose build output
bazel build -s //src/app:main

# Show command lines
bazel build --subcommands //src/app:main

# Explain rebuild
bazel build --explain=explain.txt //src/app:main

# Profile build
bazel build --profile=profile.json //...
bazel analyze-profile profile.json

# View dependency graph
bazel query --output graph //src/app:main | dot -Tpng > graph.png
```

## Cache Management

```bash
# Clean current project
bazel clean

# Remove all caches
bazel clean --expunge

# Check cache info
bazel info

# Enable remote cache (add to .bazelrc)
echo 'build --remote_cache=https://cache.example.com' >> .bazelrc
```

## Querying

```bash
# List all targets
bazel query //...

# Find all tests
bazel query 'tests(//...)'

# Find targets depending on X
bazel query 'rdeps(//..., //src/lib:utils)'

# Find dependencies of X
bazel query 'deps(//src/app:main)'

# Find specific file targets
bazel query '//src/...:*.py'
```

## Common Issues

### "No such package"

```bash
# Solution: Generate BUILD files
bazel run gazelle
```

### "Module not found" (Python)

```bash
# Solution: Update dependencies
./tools/repin
bazel run gazelle
```

### "Package not found" (Go)

```bash
# Solution: Update modules
go mod tidy
bazel mod tidy
bazel run gazelle
```

### "direnv not working"

```bash
# Solution: Reload environment
direnv allow
bazel run //tools:bazel_env
direnv reload
```

### "Build is slow"

```bash
# Solutions:
# 1. Enable remote cache
echo 'build --remote_cache=https://cache.example.com' >> .bazelrc

# 2. Profile the build
bazel build --profile=profile.json //...
bazel analyze-profile profile.json

# 3. Adjust resource limits
echo 'build --jobs=auto' >> .bazelrc
echo 'build --local_ram_resources=HOST_RAM*.8' >> .bazelrc
```

## Configuration Files

### Key Files by Language

| Language | Config Files |
|----------|--------------|
| Python | `pyproject.toml`, `requirements/*.txt` |
| JavaScript | `package.json`, `pnpm-lock.yaml` |
| Go | `go.mod`, `go.sum` |
| Java/Kotlin | `MODULE.bazel` (maven section) |
| C/C++ | `.clang-tidy` |
| Rust | `Cargo.toml` |

### Bazel Configuration

| File | Purpose |
|------|---------|
| `MODULE.bazel` | Dependency declarations |
| `.bazelrc` | Bazel configuration |
| `.bazelversion` | Bazel version pin |
| `BUILD` | Build target definitions |
| `REPO.bazel` | Repository setup (legacy) |

### Development Environment

| File | Purpose |
|------|---------|
| `.envrc` | direnv configuration |
| `tools/tools.lock.json` | Tool version lockfile |
| `tools/BUILD` | Tool definitions |

## Environment Variables

```bash
# Skip cache during build
export BAZEL_CACHE=false

# Verbose Bazel output
export BAZEL_VERBOSE=1

# Custom Bazel options
export BAZEL_OPTS="--jobs=4"

# Enable debugging
export DEBUG=1
```

## Git Workflow

```bash
# Before committing
format                  # Format code
aspect lint //...       # Check linting
bazel test //...       # Run tests
bazel build //...      # Ensure it builds

# Commit with conventional format
git commit -m "feat: add new feature"
git commit -m "fix: resolve bug in X"
git commit -m "docs: update README"
git commit -m "test: add tests for Y"
```

## CI/CD Commands

```bash
# Build everything (CI)
bazel build --config=ci //...

# Test with coverage
bazel coverage //...

# Build release artifacts
bazel build --config=release //...

# Push containers
bazel run --config=release //src/app:image.push
```

## Performance Tips

1. **Use remote caching** for faster CI builds
2. **Run tests with tags** to run subsets quickly
3. **Use ibazel** for watch mode during development
4. **Profile slow builds** with `--profile`
5. **Limit jobs** if running out of memory: `--jobs=4`
6. **Clean occasionally** if cache gets corrupted: `bazel clean`

## Getting Help

```bash
# Bazel help
bazel help build
bazel help test

# Aspect CLI help
aspect help
aspect help lint

# Tool help
pnpm help
go help mod
```

## Online Resources

- **Bazel Docs**: https://bazel.build/
- **Aspect CLI Docs**: https://docs.aspect.build/
- **Aspect Workflows**: https://aspect.build/workflows
- **rules_lint**: https://github.com/aspect-build/rules_lint
- **Slack**: #aspect-build on [Bazel Slack](https://slack.bazel.build)

## Template Testing

```bash
# Test template generation
./test.sh kitchen-sink

# Test specific preset
./test.sh py
./test.sh go

# Test in custom directory
scaffold new --output-dir=/tmp/test --preset=py --no-prompt .
cd /tmp/test
bazel test //...
```

---

**Quick Navigation:**

- [Getting Started](./user-guide/getting-started.md)
- [Development Workflow](./user-guide/development-workflow.md)
- [Architecture](./contributor-guide/architecture.md)
- [FAQ](./faq.md)
- [Documentation Home](./overview.md)
