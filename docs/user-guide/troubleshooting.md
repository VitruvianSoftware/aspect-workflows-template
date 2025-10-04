# Troubleshooting Guide

This guide covers common issues you might encounter when using the Aspect Workflows Template and how to resolve them.

## Table of Contents

1. [Template Generation Issues](#template-generation-issues)
2. [Environment Setup Issues](#environment-setup-issues)
3. [Build Issues](#build-issues)
4. [Dependency Issues](#dependency-issues)
5. [Test Issues](#test-issues)
6. [Formatting and Linting Issues](#formatting-and-linting-issues)
7. [Performance Issues](#performance-issues)
8. [IDE Integration Issues](#ide-integration-issues)

## Template Generation Issues

### Scaffold command not found

**Problem**: Running `scaffold new` returns "command not found"

**Solution**:

```bash
# Install scaffold
brew tap hay-kot/scaffold-tap
brew install scaffold

# Or with Go
go install github.com/hay-kot/scaffold@latest

# Verify installation
scaffold --version
```

### Template generation fails with parse errors

**Problem**: Generation fails with Go template syntax errors

**Solution**:

1. Ensure you're using the correct template URL:

   ```bash
   scaffold new github.com/aspect-build/aspect-workflows-template
   ```

2. Check your scaffold version:

   ```bash
   scaffold --version  # Should be >= v0.6.1
   ```

3. Try with a preset to bypass interactive mode:

   ```bash
   scaffold new --preset=minimal --no-prompt github.com/aspect-build/aspect-workflows-template
   ```

### Post-scaffold hook fails

**Problem**: Project generates but post_scaffold hook fails

**Solution**:

1. Check that Bazel is installed:

   ```bash
   bazel --version
   ```

2. Ensure you have write permissions in the output directory
3. Check the error message - it may indicate missing dependencies
4. Try running the hook steps manually:

   ```bash
   cd <generated-project>
   bazel run @buildifier_prebuilt//:buildifier -- -r .
   bazel run //tools/format
   ./tools/repin
   ```

## Environment Setup Issues

### direnv not loading environment

**Problem**: After running `direnv allow`, tools are still not on PATH

**Solution**:

```bash
# Check direnv is installed
direnv --version

# Allow the directory
direnv allow

# Check if .envrc is being processed
direnv status

# Generate the bin tree
bazel run //tools:bazel_env

# Reload direnv
direnv reload

# Verify PATH includes bazel-out directory
echo $PATH | grep bazel-out
```

### bazel_env generation fails

**Problem**: `bazel run //tools:bazel_env` fails

**Solution**:

```bash
# Check Bazel can build the target
bazel build //tools:bazel_env

# Check for errors in tools/BUILD file
cat tools/BUILD

# Try cleaning and rebuilding
bazel clean
bazel run //tools:bazel_env
```

### Tools not available after setup

**Problem**: Even with direnv loaded, commands like `format` or `pnpm` are not found

**Solution**:

```bash
# Check the bin tree was created
ls -la bazel-out/bazel_env-opt/bin/tools/bazel_env/bin/

# Verify direnv is watching the bin tree
cat .envrc

# Manually add to PATH for testing
export PATH="$(pwd)/bazel-out/bazel_env-opt/bin/tools/bazel_env/bin:$PATH"

# Verify tools are available
which format
which pnpm  # if JavaScript enabled
which python  # if Python enabled
```

## Build Issues

### "No such package" error

**Problem**: `bazel build //pkg:target` fails with "no such package: 'pkg'"

**Solution**:

```bash
# Generate BUILD files with Gazelle
bazel run gazelle

# Verify BUILD file was created
ls pkg/BUILD

# If still missing, create BUILD file manually or check .bazelignore
cat .bazelignore
```

### "Target not declared" error

**Problem**: BUILD file exists but target is not found

**Solution**:

```bash
# Run Gazelle to update BUILD file
bazel run gazelle

# Check BUILD file has the target
cat pkg/BUILD

# Query available targets in package
bazel query //pkg:all

# Check for typos in target name
bazel query //pkg:*
```

### MODULE.bazel errors

**Problem**: Errors loading MODULE.bazel

**Solution**:

```bash
# Check MODULE.bazel syntax
bazel info workspace  # Will error if MODULE.bazel is invalid

# For Go dependencies, update MODULE.bazel
bazel mod tidy

# Check for duplicate bazel_dep declarations
grep -n "bazel_dep" MODULE.bazel

# Clear external cache if corrupted
bazel clean --expunge
```

### Compilation errors after updating dependencies

**Problem**: Build fails after updating dependencies

**Solution**:

```bash
# Clean and rebuild
bazel clean
bazel build //...

# For Python, regenerate lockfiles
./tools/repin
bazel run gazelle

# For JavaScript, ensure lockfile is in sync
pnpm install
bazel build //...

# For Go, ensure all modules are in MODULE.bazel
go mod tidy
bazel mod tidy
bazel run gazelle
```

## Dependency Issues

### Python: "No module named 'X'"

**Problem**: Python imports fail at runtime or test time

**Solution**:

```bash
# Add dependency to pyproject.toml
vim pyproject.toml

# Regenerate lockfiles
./tools/repin

# Update BUILD files
bazel run gazelle

# If Gazelle didn't add dependency, add manually to BUILD:
# deps = ["@pip//package_name"],

# Verify package is in lockfile
grep "package_name" requirements/all.txt
```

### JavaScript: "Cannot find module 'X'"

**Problem**: JavaScript imports fail

**Solution**:

```bash
# Install the package
pnpm add package-name

# Verify it's in package.json
grep "package-name" package.json

# Rebuild
bazel build //path/to:target

# If still failing, check BUILD file has deps
cat path/to/BUILD

# Clear node_modules and rebuild
rm -rf node_modules
pnpm install
bazel clean
bazel build //...
```

### Go: "Package not found"

**Problem**: Go imports fail to resolve

**Solution**:

```bash
# Update go.mod with the import
go mod tidy

# Update MODULE.bazel
bazel mod tidy

# Regenerate BUILD files
bazel run gazelle

# Verify use_repo includes the package
grep "use_repo" MODULE.bazel

# If missing, add manually to MODULE.bazel use_repo()
```

### Maven: Java dependency not found

**Problem**: Java dependencies can't be resolved

**Solution**:

```bash
# Check maven.install in MODULE.bazel
vim MODULE.bazel

# Repin Maven dependencies
bazel run @unpinned_maven//:pin

# Verify maven_install.json was updated
git diff maven_install.json

# Clean and rebuild
bazel clean
bazel build //...
```

## Test Issues

### Tests failing in Bazel but pass locally

**Problem**: `bazel test` fails but running test directly works

**Solution**:

1. **Check test is hermetic** - doesn't depend on:
   - System environment variables
   - Files outside the workspace
   - Network access (unless explicitly allowed)
   - System time or timezone

2. **Add missing dependencies**:

   ```python
   py_test(
       name = "my_test",
       srcs = ["my_test.py"],
       deps = [
           ":my_lib",
           "@pip//pytest",  # Don't forget test frameworks
       ],
       data = [
           "testdata/input.json",  # Include test data
       ],
   )
   ```

3. **Check for race conditions**:

   ```bash
   # Run test multiple times
   bazel test --runs_per_test=10 //pkg:test
   ```

4. **View test logs**:

   ```bash
   bazel test //pkg:test
   cat bazel-testlogs/pkg/test/test.log
   ```

### Tests are cached and not running

**Problem**: Test output says "cached" and doesn't re-run

**Solution**:

```bash
# Disable test caching
bazel test --cache_test_results=no //pkg:test

# Or force re-run
bazel test --nocache_test_results //pkg:test

# Check test is properly marked as test (not build)
bazel query --output=build //pkg:test
```

### Test timeouts

**Problem**: Tests timeout before completion

**Solution**:

```python
# Increase timeout in BUILD file
py_test(
    name = "slow_test",
    timeout = "long",  # short, moderate, long, or eternal
    srcs = ["slow_test.py"],
)
```

Or override on command line:

```bash
bazel test --test_timeout=300 //pkg:test
```

## Formatting and Linting Issues

### Format command not working

**Problem**: Running `format` returns "command not found"

**Solution**:

```bash
# Ensure direnv is loaded
direnv allow
direnv reload

# Regenerate bazel_env
bazel run //tools:bazel_env

# Try running directly via Bazel
bazel run //tools/format

# Check format is in tools/BUILD
grep "format" tools/BUILD
```

### Pre-commit hook not running

**Problem**: Git commits without running formatter

**Solution**:

```bash
# Enable git hooks
git config core.hooksPath githooks

# Verify hook is executable
chmod +x githooks/pre-commit

# Test the hook
./githooks/pre-commit

# Check hook script is correct
cat githooks/pre-commit
```

### Linting fails with "aspect: command not found"

**Problem**: `aspect lint` doesn't work

**Solution**:

```bash
# Aspect CLI may not be installed
# Use bazel directly instead
bazel build //... --aspects @aspect_rules_lint//lint:lint.bzl%lint

# Or install Aspect CLI
# See https://docs.aspect.build/cli/install

# Alternatively, use individual linters
bazel run //tools/lint:eslint -- path/to/file.js
```

### Lint errors that can't be fixed

**Problem**: Linter reports errors that seem incorrect

**Solution**:

1. Check linter configuration files:
   - Python: `pyproject.toml`, `.ruff.toml`
   - JavaScript: `eslint.config.mjs`
   - Java: `pmd.xml`
   - Kotlin: `ktlint-baseline.xml`

2. Add exceptions to baseline files:

   ```bash
   # Kotlin example
   aspect lint //... --fix
   # This updates ktlint-baseline.xml
   ```

3. Disable specific rules if needed:

   ```python
   # In code comments
   # pylint: disable=rule-name
   # ruff: noqa: E501
   ```

## Performance Issues

### Slow builds

**Problem**: `bazel build //...` takes too long

**Solution**:

```bash
# 1. Profile the build
bazel build --profile=profile.json //...
bazel analyze-profile profile.json

# 2. Check cache hit rates
bazel info

# 3. Enable remote cache (if available)
echo 'build --remote_cache=https://cache.example.com' >> .bazelrc

# 4. Limit parallelism if system is overloaded
bazel build --jobs=4 //...

# 5. Build only what changed
bazel build //path/to/changed:target
```

### Bazel using too much disk space

**Problem**: Workspace is consuming excessive disk space

**Solution**:

```bash
# Check disk usage
bazel info output_base
du -sh $(bazel info output_base)

# Clean build outputs
bazel clean

# Remove all caches (drastic)
bazel clean --expunge

# Clean only external dependencies
bazel clean --expunge_async

# Configure disk cache limits in .bazelrc
echo 'build --disk_cache=~/.cache/bazel --experimental_disk_cache_gc_max_size=50GB' >> .bazelrc
```

### Bazel using too much memory

**Problem**: Bazel crashes or system runs out of memory

**Solution**:

```bash
# Limit memory usage
echo 'build --local_ram_resources=HOST_RAM*.5' >> .bazelrc

# Reduce concurrent actions
echo 'build --jobs=4' >> .bazelrc
echo 'build --local_cpu_resources=4' >> .bazelrc

# For large builds, use streaming output
bazel build --experimental_stream_log=build.log //...
```

## IDE Integration Issues

### VSCode not recognizing imports

**Problem**: Editor shows import errors but code builds fine

**Solution**:

**Python**:

```bash
# Create/update .vscode/settings.json
cat > .vscode/settings.json <<EOF
{
  "python.analysis.extraPaths": [
    "bazel-bin",
    "bazel-out/k8-fastbuild/bin"
  ]
}
EOF
```

**JavaScript**:

```bash
# Run pnpm install to create node_modules
pnpm install
```

**Go**:

```bash
# Generate go.work for IDE
go work init
go work use .
```

### IntelliJ not finding Bazel targets

**Problem**: IntelliJ shows errors or can't find targets

**Solution**:

```bash
# Install Bazel plugin for IntelliJ
# Then sync project: Tools > Bazel > Sync Project

# Or generate .idea project
bazel run @rules_intellij//intellij:project

# Restart IntelliJ after sync
```

### Language server errors

**Problem**: Language server shows errors

**Solution**:

**C/C++**:

```bash
# Generate compile_commands.json
bazel run @hedron_compile_commands//:refresh_all

# VSCode will automatically use this
```

**Python**:

```bash
# Use bazel-managed Python
which python  # Should point to bazel-out/... if direnv loaded

# Or configure explicitly
echo 'python.pythonPath: "bazel-bin/..."' >> .vscode/settings.json
```

## Getting Additional Help

### Check Logs

```bash
# Bazel output
bazel build //... 2>&1 | tee build.log

# Test logs
cat bazel-testlogs/path/to/test/test.log

# Verbose output
bazel build -s //...  # Shows all commands
```

### Useful Debugging Commands

```bash
# Show what Bazel will build
bazel query //...

# Show dependencies
bazel query 'deps(//pkg:target)'

# Show reverse dependencies
bazel query 'rdeps(//..., //pkg:target)'

# Show why target was rebuilt
bazel build --explain=explain.log //...
cat explain.log

# Check for rule errors
bazel query --output=build //pkg:target
```

### Community Support

If you're still stuck:

1. **Search existing issues**: [GitHub Issues](https://github.com/aspect-build/aspect-workflows-template/issues)
2. **Ask in Slack**: #aspect-build on [Bazel Slack](https://slack.bazel.build)
3. **GitHub Discussions**: [Ask a question](https://github.com/aspect-build/aspect-workflows-template/discussions)
4. **Check documentation**:
   - [Bazel Docs](https://bazel.build/)
   - [Aspect CLI Docs](https://docs.aspect.build/)
   - [FAQ](../faq.md)

### Reporting Bugs

When reporting issues, include:

1. **Template version** (git commit hash)
2. **Bazel version**: `bazel --version`
3. **Operating system**: `uname -a`
4. **Preset used**: minimal, py, go, etc.
5. **Full error message**
6. **Steps to reproduce**
7. **Relevant BUILD files**

---

**Back**: [User Guide](./README.md) | **FAQ**: [FAQ](../faq.md) | **Quick Reference**: [Quick Reference](../quick-reference.md)
