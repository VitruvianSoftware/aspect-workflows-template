# Frequently Asked Questions (FAQ)

## General Questions

### What is the Aspect Workflows Template?

The Aspect Workflows Template is a project generator (scaffold) that creates production-ready Bazel monorepos with support for multiple programming languages, integrated tooling, and best practices built-in.

### Why should I use this template instead of starting from scratch?

Starting a Bazel project from scratch requires:

- Understanding Bazel's module system (bzlmod)
- Configuring language rules and toolchains
- Setting up formatting and linting
- Integrating package managers
- Configuring development environments

This template provides all of this pre-configured and tested, saving days or weeks of setup time.

### What's the difference between this and rules_* repositories?

Bazel `rules_*` repositories (like `rules_python`, `rules_go`) provide the build rules for specific languages. This template integrates multiple rules repositories together with:

- Pre-configured MODULE.bazel
- Development environment setup
- Code quality tools
- Package management
- Best practice configurations

Think of it as a "batteries included" starter kit that uses the rules_* repositories.

## Installation and Setup

### Do I need to install Bazel first?

No! The template uses `.bazelversion` to automatically download the correct Bazel version when you first run `bazel` commands. However, you do need to install the Bazel binary itself - see [bazel.build](https://bazel.build/install).

### What is direnv and why do I need it?

[direnv](https://direnv.net/) automatically modifies your shell environment when you enter a project directory. The template uses it to add Bazel-managed tools to your PATH, eliminating the need to manually install development tools.

**Without direnv:**

```bash
bazel run @pnpm -- install
bazel run //tools:format -- my-file.js
```

**With direnv:**

```bash
pnpm install        # runs Bazel-managed pnpm
format my-file.js   # runs Bazel-managed formatter
```

### Can I use this template without direnv?

Yes, but it's less convenient. You'll need to run tools via `bazel run` instead of directly:

```bash
bazel run @pnpm -- install
bazel run //tools/format:format -- path/to/file
```

### How do I update the template after generating a project?

The template is a starting point - once generated, the project is yours to modify. To get updates:

1. Check the template repository for new features
2. Manually apply relevant changes to your project
3. Consider using `git remote add template <url>` to track template updates

There's no automatic update mechanism by design - you have full control.

## Language Support

### Which languages are supported?

- JavaScript & TypeScript
- Python
- Go
- Java
- Kotlin
- C & C++
- Rust
- Shell (Bash)

### Can I add multiple languages to one project?

Yes! You can select multiple languages during project generation. The template is designed for polyglot monorepos.

### Can I add a language after project generation?

Not easily. The template is designed to be run once. To add a language:

1. Generate a new project with the desired languages
2. Copy the relevant sections from the new project's:
   - MODULE.bazel (language dependencies)
   - BUILD files
   - tools/BUILD (language tools)
   - Configuration files

### Why isn't my favorite language supported?

The template focuses on languages with mature Bazel rules. To add a new language:

1. Check if Bazel rules exist for that language
2. Submit a feature request or PR to the template repository
3. Provide configuration examples and test cases

## Dependency Management

### How do I add a new dependency?

It depends on the language:

**Python:**

```bash
# Add to pyproject.toml
vim pyproject.toml
# Update lockfiles
./tools/repin
# Update BUILD files
bazel run gazelle
```

**JavaScript:**

```bash
pnpm add <package>
# BUILD files are auto-updated
```

**Go:**

```bash
go get <package>
go mod tidy
bazel mod tidy
bazel run gazelle
```

**Java:**

```bash
# Edit MODULE.bazel to add to maven.install artifacts
vim MODULE.bazel
# Repin maven dependencies
bazel run @unpinned_maven//:pin
```

### What is the `repin` script?

The `repin` script updates dependency lockfiles after you modify dependency specifications. It runs:

- `requirements.update` for Python
- `gazelle_python_manifest.update` for Python imports

### Can I use a different package manager?

The template uses specific package managers for each language:

- **Python**: pip (managed by rules_python)
- **JavaScript**: pnpm (managed by rules_js)
- **Go**: Go modules (native)
- **Java**: Maven (via rules_jvm_external)

Switching package managers would require significant modification of the template.

## Development Workflow

### What is Gazelle?

Gazelle is a BUILD file generator that automatically creates and updates BUILD files based on your source code. When you add new files or change imports, run:

```bash
bazel run gazelle
```

Gazelle analyzes your code and generates the appropriate `BUILD` targets.

### How do I format my code?

If you enabled formatting during project generation:

```bash
# Format all files
format

# Format specific files
format path/to/file.py path/to/file.js
```

Or with Bazel directly:

```bash
bazel run //tools/format
```

### How do I lint my code?

```bash
# Lint everything
aspect lint //...

# Lint specific targets
aspect lint //src/app:all
```

The linting is integrated into Bazel, so results are cached and only changed files are re-linted.

### What are the pre-commit hooks?

Pre-commit hooks automatically format your code before each commit. Enable them with:

```bash
git config core.hooksPath githooks
```

Now every commit will automatically format staged files.

## Building and Testing

### How do I build my project?

```bash
# Build everything
bazel build //...

# Build specific target
bazel build //src/app:my_app

# Build with optimizations
bazel build -c opt //...
```

### How do I run tests?

```bash
# Run all tests
bazel test //...

# Run specific test
bazel test //src:my_test

# Use Aspect CLI for better UX
aspect test //...
```

### What's the difference between `bazel` and `aspect` commands?

`aspect` is an enhanced CLI built on top of `bazel` that provides:

- Better error messages
- Interactive test output
- `aspect lint` command
- Performance insights
- Improved developer experience

Both work with the same workspace. Use `aspect` when available for a better experience.

### How do I run my application?

```bash
# Run a binary target
bazel run //src:my_app

# With arguments
bazel run //src:my_app -- --arg1 --arg2

# Watch mode (rebuilds on file changes)
ibazel run //src:my_app
```

## Containers and Deployment

### How do I build a container image?

If you enabled OCI support during generation:

```bash
# Build an image
bazel build //src:my_app_image

# Load into Docker
bazel run //src:my_app_image.load

# Push to registry
bazel run //src:my_app_image.push
```

### What base images are used?

- **Go**: Distroless base (`gcr.io/distroless/base`)
- **Python**: Ubuntu (`ubuntu:latest`)

You can customize base images in MODULE.bazel's `oci.pull` sections.

### How do I configure image registries?

Add configuration to your .bazelrc:

```bash
# .bazelrc
build --@rules_oci//oci:registry=docker.io/myorg
```

Or pass at runtime:

```bash
bazel run //src:my_app_image.push --@rules_oci//oci:registry=gcr.io/myproject
```

## Troubleshooting

### Bazel is using too much disk space

Bazel's cache can grow large. Clean it with:

```bash
# Clean current project
bazel clean

# Remove all caches (drastic)
bazel clean --expunge

# Remove only external dependencies
bazel clean --expunge_async
```

### My build is slow

Common causes:

1. **No remote cache**: Set up remote caching
2. **Recompiling everything**: Check cache hit rates
3. **Too many local jobs**: Adjust `--jobs` flag
4. **Large dependencies**: Profile with `--profile`

```bash
# Profile a build
bazel build --profile=profile.json //...
bazel analyze-profile profile.json
```

### Tests are failing in CI but pass locally

Possible causes:

1. **Non-hermetic tests**: Tests depend on local environment
2. **Race conditions**: Tests have timing issues
3. **Platform differences**: Test assumes specific OS/architecture
4. **Missing dependencies**: Test doesn't declare all dependencies

Make tests hermetic by declaring all dependencies and avoiding system state.

### I'm getting "no such package" errors

This usually means:

1. **Missing BUILD file**: Run `bazel run gazelle`
2. **Visibility issue**: Add target to visibility list
3. **Typo in label**: Check target names
4. **External dependency not loaded**: Check MODULE.bazel

### direnv isn't loading the environment

```bash
# Allow direnv for this directory
direnv allow

# Generate the bin tree
bazel run //tools:bazel_env

# Reload direnv
direnv reload
```

## Advanced Topics

### Can I use remote caching?

Yes! Add to your .bazelrc:

```bash
# .bazelrc
build --remote_cache=https://your-cache-server
build --remote_upload_local_results=true
```

Popular options:

- Google Cloud Storage
- AWS S3
- Bazel Remote Cache
- BuildBuddy

### Can I use remote execution?

Yes, if you have access to a remote execution service:

```bash
# .bazelrc
build --remote_executor=grpcs://your-remote-executor
```

Services:

- BuildBuddy
- EngFlow
- BuildFarm

### How do I configure cross-compilation?

The template includes platform support. To build for different platforms:

```bash
# Build for Linux ARM64
bazel build --platforms=//tools/platforms:linux_arm64 //...

# Build for multiple platforms
bazel build --platforms=//tools/platforms:linux_amd64,//tools/platforms:linux_arm64 //...
```

### Can I use this in a monorepo with existing code?

Yes, but carefully:

1. Generate the template in a new directory
2. Copy the generated structure into your monorepo
3. Merge configuration files (MODULE.bazel, .bazelrc, etc.)
4. Test thoroughly

Consider generating a separate workspace first to understand the structure.

## Best Practices

### What should I commit to version control?

**Commit:**

- ✅ All source code
- ✅ BUILD files
- ✅ MODULE.bazel and other config files
- ✅ Lockfiles (pnpm-lock.yaml, go.sum, requirements.txt)
- ✅ .bazelrc and .bazelversion
- ✅ tools/tools.lock.json

**Don't commit:**

- ❌ bazel-* symlinks
- ❌ node_modules/ (managed by Bazel)
- ❌ .venv/ (virtual environments)
- ❌ Build outputs

### How should I structure my code?

Follow language conventions:

- **Go**: `cmd/`, `pkg/`, `internal/`
- **Python**: Package structure with `__init__.py`
- **JavaScript**: `packages/` for multiple packages
- **Java/Kotlin**: Standard package structure

Keep BUILD files close to source code - typically one per directory.

### Should I use rules_python or aspect_rules_py?

The template uses both:

- **rules_python**: Official Python rules, manages dependencies
- **aspect_rules_py**: Enhanced developer experience, better performance

aspect_rules_py builds on rules_python, providing additional features.

## Getting Help

### Where can I ask questions?

1. **GitHub Issues**: Bug reports and feature requests
2. **GitHub Discussions**: Questions and community support
3. **Bazel Slack**: #aspect-build channel
4. **Documentation**: Check the docs first!

### How do I report a bug?

1. Check if the issue is already reported
2. Create a minimal reproduction
3. Include:
   - Template version
   - Bazel version
   - Operating system
   - Full error messages
   - Steps to reproduce

### How can I contribute?

See the [Contributor Guide](./contributor-guide/README.md) for:

- How to submit PRs
- Coding standards
- Testing requirements
- Development workflow

---

**Back**: [Documentation Home](./overview.md)
