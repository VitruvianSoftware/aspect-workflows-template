# Aspect Workflows Template - AI Coding Agent Guide

This is a **Bazel monorepo template generator** using [Scaffold](https://hay-kot.github.io/scaffold/) for multi-language project scaffolding.

## Architecture Overview

- **Template Engine**: Go templates in `{{ .ProjectSnake }}/` directory process via Scaffold CLI
- **Configuration**: `scaffold.yaml` defines questions, features, computed values, and file inclusion rules
- **Generated Projects**: Use Bazel with bzlmod (MODULE.bazel), not WORKSPACE
- **Post-Processing**: `hooks/post_scaffold` formats files, runs repin, initializes package managers

### Key Namespaces

- `.Scaffold.*` - User answers (e.g., `.Scaffold.lint`, `.Scaffold.langs`)
- `.Computed.*` - Derived flags (e.g., `.Computed.javascript`, `.Computed.python`)
- `.ProjectSnake` - Output directory name

## Critical Developer Workflows

### Testing Template Changes

```bash
# Test specific preset (kitchen-sink, py, go, js, java, kotlin, cpp, rust, shell, minimal)
./test.sh kitchen-sink

# Manual test with custom output
scaffold new --output-dir=/tmp/test --preset=py --no-prompt .
cd /tmp/test && bazel test //...
```

### Modifying Templates

1. Edit files in `{{ .ProjectSnake }}/` with Go template syntax
2. Update `scaffold.yaml` features if adding/removing conditional files
3. Test generation with affected presets using `./test.sh`
4. `hooks/post_scaffold` runs buildifier, format, and repin automatically

### Generated Project Workflow

**Setup**:
```bash
cd <generated-project>
direnv allow                    # Enable environment
bazel run //tools:bazel_env     # Setup PATH with tools
git config core.hooksPath githooks
```

**Build & Test**:
```bash
bazel build //...               # Build everything
aspect test //...               # Run tests (requires Aspect CLI)
bazel run gazelle               # Generate/update BUILD files
```

**Formatting & Linting**:
```bash
format                          # Format all (from direnv PATH)
format path/to/file            # Format single file
aspect lint //...              # Run all linters (aspect CLI)
aspect lint --fix //...        # Apply auto-fixes
```

**Dependency Management**:
- Python: Edit `pyproject.toml` → `./tools/repin` → `bazel run gazelle`
- Go: `go mod tidy` → `bazel mod tidy` → `bazel run gazelle`
- JS: `pnpm add <pkg>` → auto-updates `pnpm-lock.yaml`
- Java: Edit `MODULE.bazel` maven.install → `bazel run @unpinned_maven//:pin`

## Project-Specific Conventions

### Bazel Patterns

- **BUILD over BUILD.bazel**: `# gazelle:build_file_name BUILD` directive enforced
- **bzlmod only**: No WORKSPACE files, all deps in MODULE.bazel
- **Gazelle prefixes**: Go uses `# gazelle:prefix github.com/example/project`
- **Rules mapping**: Python uses aspect_rules_py not rules_python: `# gazelle:map_kind py_binary py_binary @aspect_rules_py//py:defs.bzl`

### File Structure Conventions

- `tools/` - Development tooling, exported via bazel-env.bzl to PATH
- `.bazelrc` - Imports `tools/preset.bazelrc`, language-specific configs
- `requirements/` (Python) - `runtime.txt`, `all.txt` lockfiles, updated by `./tools/repin`
- `.envrc` - Sets up direnv to add Bazel-managed tools to PATH

### Linting Architecture

Uses [rules_lint](https://github.com/aspect-build/rules_lint) with Bazel aspects:
- Linters run as cached Bazel actions (not external processes)
- Configuration in `tools/lint/linters.bzl` defines aspects per language
- Aspect CLI provides `aspect lint` command for better UX than raw Bazel
- Report files are Bazel outputs, cached like any other action

### Code Generation (Gazelle)

- Python: Requires `modules_mapping` from `gazelle_python_manifest` (updated by repin)
- Go: Reads go.mod, requires `bazel mod tidy` to update MODULE.bazel use_repo
- JS: Uses `npm_link_all_packages` to expose npm dependencies
- **Always run after dependency changes**: `bazel run gazelle`

## Integration Points

### bazel-env.bzl System

- `bazel run //tools:bazel_env` generates `bin/` tree with symlinks to Bazel-managed tools
- `.envrc` adds this bin tree to PATH via direnv
- Enables running `format`, `buildifier`, `pnpm`, `go`, etc. without manual installation
- Tool versions locked in `tools/tools.lock.json` (multitool lockfile)

### Multitool Integration

- `@multitool//tools/<name>` - Hermetic CLI tools (ruff, shellcheck, etc.)
- Lockfile at `tools/tools.lock.json`
- Used in both BUILD files and exported to PATH via bazel-env

### Container Images (Optional)

- `tools/oci/py3_image.bzl` and `tools/oci/go_image.bzl` wrap rules_oci
- Base images pulled via oci.pull() in MODULE.bazel
- Build: `bazel build //app:image`, Load: `bazel run //app:image.load`

## Common Pitfalls

1. **Don't use WORKSPACE**: Projects are bzlmod-only (MODULE.bazel)
2. **Don't manually install tools**: Use `bazel run //tools:bazel_env` and direnv
3. **Don't forget repin**: After Python/Java dependency changes, run `./tools/repin`
4. **Don't skip gazelle**: After code changes importing new packages, run `bazel run gazelle`
5. **Template syntax**: Use `{{ if .Computed.javascript }}` not `{{ if .Scaffold.langs contains ... }}`
6. **Feature flags**: Files only included if `features:` globs match and value is true

## Documentation Locations

- User workflows: `docs/user-guide/`
- Template internals: `docs/contributor-guide/architecture.md`, `docs/contributor-guide/template-system.md`
- Quick commands: `docs/quick-reference.md`
- Generated project guide: `{{ .ProjectSnake }}/README.bazel.md`

## Key Files to Understand

- `scaffold.yaml` - Complete template configuration, questions, features, computed values
- `MODULE.bazel` - Shows all Bazel dependencies and extension system usage
- `hooks/post_scaffold` - Post-generation formatting and initialization steps
- `tools/BUILD` - Tool exports and bazel_env configuration
- `.bazelrc` - Bazel configuration including lint aspects, language-specific flags
