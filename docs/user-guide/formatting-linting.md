# Formatting and Linting

This guide covers code formatting and linting in projects generated from the Aspect Workflows Template.

## Table of Contents

1. [Code Formatting](#code-formatting)
2. [Code Linting](#code-linting)
3. [Pre-Commit Hooks](#pre-commit-hooks)
4. [Language-Specific Tools](#language-specific-tools)
5. [Configuration](#configuration)
6. [CI Integration](#ci-integration)

## Code Formatting

### Overview

The template includes integrated code formatting via `rules_lint`. Formatting is available through a unified `format` command.

### Using the Format Command

```bash
# Format all files in the workspace
format

# Format specific files
format src/app/main.py src/lib/utils.py

# Format files matching a pattern
format $(find src -name "*.py")

# Check formatting without modifying files
format --check

# Format only staged git files
git diff --cached --name-only --diff-filter=AM | xargs format
```

### Format Command Setup

The `format` command is available via direnv:

```bash
# Ensure direnv is loaded
direnv allow

# Regenerate if needed
bazel run //tools:bazel_env

# Verify format is available
which format
```

### Language-Specific Formatters

The template includes formatters for each language:

| Language | Formatter | Style |
|----------|-----------|-------|
| Python | Ruff | PEP 8 |
| JavaScript/TypeScript | Prettier | Standard |
| Go | gofmt | Go standard |
| Java | google-java-format | Google style |
| Kotlin | ktlint | Kotlin official |
| C/C++ | clang-format | LLVM |
| Rust | rustfmt | Rust standard |
| Shell | shfmt | Shell standard |
| Bazel | Buildifier | Bazel standard |

## Code Linting

### Overview

Linting runs automatically via Bazel aspects, caching results for performance.

### Using Aspect Lint

```bash
# Lint everything
aspect lint //...

# Lint specific package
aspect lint //src/app:all

# Lint specific target
aspect lint //src/app:main

# Apply automatic fixes
aspect lint --fix //...
```

### Using Bazel Directly

If Aspect CLI is not available:

```bash
# Run linters via Bazel
bazel build //... \
  --aspects=@aspect_rules_lint//lint:lint.bzl%lint \
  --output_groups=+rules_lint_report

# View lint reports
find bazel-bin -name "*_lint_report.txt" -exec cat {} \;
```

### Understanding Lint Output

```shell
Lint results for //src/app:main:

src/app/main.py:10:5: E501 line too long (88 > 79 characters)
src/app/main.py:15:1: W293 blank line contains whitespace
src/app/util.py:5:1: F401 'sys' imported but unused

Summary:
* errors: 2
* warnings: 1
```

### Language-Specific Linters

| Language | Linter | Configuration |
|----------|--------|---------------|
| Python | Ruff | `pyproject.toml` |
| JavaScript/TypeScript | ESLint | `eslint.config.mjs` |
| Go | nogo (vet + staticcheck) | `tools/lint/BUILD` |
| Java | PMD | `pmd.xml` |
| Kotlin | ktlint | `ktlint-baseline.xml` |
| C/C++ | clang-tidy | `.clang-tidy` |
| Rust | clippy | `Cargo.toml` |
| Shell | shellcheck | `.shellcheckrc` |

## Pre-Commit Hooks

### Enabling Pre-Commit Hooks

```bash
# Enable git hooks
git config core.hooksPath githooks

# Verify hook is executable
ls -l githooks/pre-commit
```

### How Pre-Commit Works

The pre-commit hook automatically formats staged files:

```bash
#!/usr/bin/env bash
# githooks/pre-commit
git diff --cached --diff-filter=AM --name-only -z | \
  xargs --null --no-run-if-empty bazel run //:format --
```

### Testing Pre-Commit Hook

```bash
# Test the hook manually
./githooks/pre-commit

# Make a test commit
echo "test" >> README.md
git add README.md
git commit -m "test"  # Hook runs automatically
```

### Bypassing Pre-Commit

In rare cases where you need to skip the hook:

```bash
# Skip pre-commit hook (use sparingly!)
git commit --no-verify -m "emergency fix"
```

## Language-Specific Tools

### Python

#### Formatter: Ruff

```bash
# Format Python files
format src/**/*.py

# Configure in pyproject.toml
[tool.ruff]
line-length = 88
target-version = "py311"
```

#### Linter: Ruff

```bash
# Lint Python code
aspect lint //src/...

# Configure rules
[tool.ruff.lint]
select = ["E", "F", "I", "N"]
ignore = ["E501"]
```

### JavaScript/TypeScript

#### Formatter: Prettier

```bash
# Format JS/TS files
format packages/**/*.ts

# Configure in prettier.config.cjs
module.exports = {
  semi: true,
  trailingComma: 'all',
  singleQuote: true,
  printWidth: 100,
};
```

#### Linter: ESLint

```bash
# Lint JavaScript/TypeScript
aspect lint //packages/...

# Configure in eslint.config.mjs
import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';

export default [
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
];
```

### Go

#### Formatter: gofmt

```bash
# Format Go files
format $(find . -name "*.go")

# Or use go fmt directly
go fmt ./...
```

#### Linter: nogo

```bash
# Lint Go code
aspect lint //...

# nogo runs automatically on build
bazel build //...
```

### Java

#### Formatter: google-java-format

```bash
# Format Java files
format src/**/*.java
```

#### Linter: PMD

```bash
# Lint Java code
aspect lint //src/...

# Configure in pmd.xml
<?xml version="1.0"?>
<ruleset name="Custom Rules">
  <rule ref="category/java/bestpractices.xml" />
</ruleset>
```

### Kotlin

#### Formatter: ktlint

```bash
# Format Kotlin files
format src/**/*.kt

# Configure in .editorconfig
[*.kt]
indent_size = 4
```

#### Linter: ktlint

```bash
# Lint Kotlin code
aspect lint //src/...

# Update baseline for known issues
aspect lint //... --fix
```

### C/C++

#### Formatter: clang-format

```bash
# Format C++ files
format src/**/*.cpp

# Configure style in .clang-format
BasedOnStyle: LLVM
IndentWidth: 4
```

#### Linter: clang-tidy

```bash
# Lint C++ code
aspect lint //src/...

# Configure in .clang-tidy
Checks: '-*,clang-analyzer-*,bugprone-*'
```

### Rust

#### Formatter: rustfmt

```bash
# Format Rust files
format src/**/*.rs

# Or use cargo fmt
cargo fmt
```

#### Linter: clippy

```bash
# Lint Rust code
cargo clippy
```

### Shell

#### Formatter: shfmt

```bash
# Format shell scripts
format scripts/**/*.sh
```

#### Linter: shellcheck

```bash
# Lint shell scripts
aspect lint //scripts/...

# Configure in .shellcheckrc
disable=SC2086,SC2046
```

## Configuration

### Global Configuration

Linting and formatting are configured per-language in the project root:

```bash
project/
├── .clang-tidy          # C++ linting
├── .editorconfig        # Kotlin formatting
├── .shellcheckrc        # Shell linting
├── eslint.config.mjs    # JavaScript linting
├── ktlint-baseline.xml  # Kotlin lint baseline
├── pmd.xml              # Java linting
├── prettier.config.cjs  # JavaScript formatting
└── pyproject.toml       # Python config
```

### Per-File Overrides

Use in-file comments to override rules:

**Python:**

```python
# ruff: noqa: E501
def very_long_function_name():
    pass
```

**JavaScript:**

```javascript
// eslint-disable-next-line no-console
console.log('Debug message');
```

**Java:**

```java
@SuppressWarnings("PMD.AvoidPrintStackTrace")
public void debugMethod() {
    e.printStackTrace();
}
```

### Ignoring Files

Add patterns to ignore files:

**Python:**

```toml
# pyproject.toml
[tool.ruff]
exclude = [
    "generated/*",
    "*.pyi",
]
```

**JavaScript:**

```javascript
// eslint.config.mjs
export default [
  {
    ignores: ['dist/*', '*.min.js']
  }
];
```

## CI Integration

### CI Format Check

```bash
# Check formatting in CI (fails if not formatted)
format --check

# Or use Bazel
bazel test //tools/format:format_test
```

### CI Lint Check

```bash
# Lint in CI
aspect lint //...

# With Bazel
bazel build //... \
  --aspects=@aspect_rules_lint//lint:lint.bzl%lint \
  --output_groups=rules_lint_report
```

### GitHub Actions Example

```yaml
name: Code Quality

on: [push, pull_request]

jobs:
  format:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check formatting
        run: |
          direnv allow
          bazel run //tools:bazel_env
          format --check
  
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run linters
        run: aspect lint //...
```

### Pre-Merge Requirements

Configure branch protection to require:

1. ✅ Format check passes
2. ✅ Lint check passes
3. ✅ All tests pass
4. ✅ Code review approval

## Best Practices

**DO**:

- ✅ Run `format` before committing
- ✅ Enable pre-commit hooks for automatic formatting
- ✅ Address lint warnings promptly
- ✅ Use `--fix` for auto-fixable issues
- ✅ Configure linters to match team standards
- ✅ Keep configuration files in version control
- ✅ Run format and lint in CI

**DON'T**:

- ❌ Disable linters globally without team agreement
- ❌ Commit unformatted code
- ❌ Ignore linter warnings
- ❌ Use `--no-verify` to bypass hooks routinely
- ❌ Have inconsistent formatting across files
- ❌ Mix multiple formatting styles

## Troubleshooting

### Format command not found

```bash
# Ensure direnv is loaded
direnv allow
direnv reload

# Regenerate bazel_env
bazel run //tools:bazel_env

# Check PATH
echo $PATH | grep bazel-out
```

### Pre-commit hook not running

```bash
# Enable hooks
git config core.hooksPath githooks

# Make hook executable
chmod +x githooks/pre-commit

# Test manually
./githooks/pre-commit
```

### Linting errors that seem incorrect

```bash
# Check linter configuration
cat pyproject.toml  # Python
cat eslint.config.mjs  # JavaScript
cat pmd.xml  # Java

# Update configuration to disable specific rules
# See language-specific sections above
```

### Format command hangs

```bash
# Format may be processing many files
# Try formatting specific files instead
format src/app/main.py

# Or format in smaller batches
format src/app/*.py
format src/lib/*.py
```

## Advanced Topics

### Custom Linting Rules

Add custom rules for your team:

**ESLint plugin:**

```javascript
// eslint.config.mjs
import customPlugin from './tools/eslint-plugin-custom';

export default [
  {
    plugins: {
      custom: customPlugin
    },
    rules: {
      'custom/no-console-log': 'error'
    }
  }
];
```

### Formatting in IDE

Most IDEs can be configured to use project formatters:

**VSCode:**

```json
{
  "editor.formatOnSave": true,
  "[python]": {
    "editor.defaultFormatter": "charliermarsh.ruff"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

**IntelliJ:**

- Settings → Editor → Code Style
- Import scheme from project configuration

### Batch Reformatting

When introducing formatting to existing code:

```bash
# Format entire codebase
format

# Commit as separate change
git add -A
git commit -m "style: apply formatting to entire codebase"

# Update baseline for linters
aspect lint //... --fix
git commit -am "lint: update lint baseline"
```

## Next Steps

- Learn about [Dependency Management](./dependency-management.md)
- Check [Troubleshooting](./troubleshooting.md) for issues
- Review [Development Workflow](./development-workflow.md)

---

**Back**: [Building and Testing](./building-testing.md) | **Next**: [Dependency Management](./dependency-management.md)
