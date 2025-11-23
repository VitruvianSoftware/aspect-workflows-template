# Getting Started

This guide will walk you through creating your first project using the Aspect Workflows Template.

## Installation

### 1. Install Scaffold

Scaffold is the tool used to generate projects from this template.

**macOS (using Homebrew):**

```bash
brew tap hay-kot/scaffold-tap
brew install scaffold
```

**Using Go:**

```bash
go install github.com/hay-kot/scaffold@latest
```

### 2. Install Prerequisites

**Required:**

- Git
- direnv: https://direnv.net/ (for development environment management)

**Language-Specific Requirements:**

The template will install hermetic toolchains for most languages, but some initial setup may be needed:

- For JavaScript/TypeScript: Node.js will be managed by Bazel
- For Python: Python interpreter will be managed by Bazel
- For Go: Go SDK will be managed by Bazel
- For Java/Kotlin: JDK will be managed by Bazel
- For C/C++: LLVM toolchain will be managed by Bazel
- For Rust: Rust toolchain will be managed by Bazel

## Creating a New Project

### Interactive Mode (Recommended)

Run scaffold in interactive mode to answer questions about your project:

```bash
scaffold new github.com/BlueCentre/aspect-workflows-template
```

You'll be prompted to answer:

1. **Languages**: Select which programming languages you'll use (required, can select multiple)
2. **Code generation**: Setup copier, yeoman, and scaffold tools? (optional)
3. **Format and linting**: Setup rules_lint for code quality? (recommended)
4. **Version stamping**: Setup automated versioning? (optional)
5. **OCI Containers**: Setup container image building? (optional)

### Using Presets

For common scenarios, use a preset:

```bash
# Minimal project (no languages selected)
scaffold new --preset=minimal --no-prompt github.com/BlueCentre/aspect-workflows-template

# JavaScript/TypeScript project
scaffold new --preset=js --no-prompt github.com/BlueCentre/aspect-workflows-template

# Python project
scaffold new --preset=py --no-prompt github.com/BlueCentre/aspect-workflows-template

# Go project with container support
scaffold new --preset=go --no-prompt github.com/BlueCentre/aspect-workflows-template

# Java project
scaffold new --preset=java --no-prompt github.com/BlueCentre/aspect-workflows-template

# Kitchen sink (all features enabled)
scaffold new --preset=kitchen-sink --no-prompt github.com/BlueCentre/aspect-workflows-template
```

Available presets:

- `minimal` - No languages, basic structure
- `shell` - Bash scripting
- `js` - JavaScript/TypeScript
- `py` - Python
- `go` - Go with OCI support
- `java` - Java
- `kotlin` - Kotlin
- `cpp` - C/C++
- `rust` - Rust
- `kitchen-sink` - All languages and features

### Custom Output Directory

Specify where to create your project:

```bash
scaffold new --output-dir=./my-project github.com/BlueCentre/aspect-workflows-template
```

## Initial Setup

After generating your project, navigate to the generated directory:

```bash
cd <your-project-name>
```

### 1. Initialize Development Environment

Enable direnv to automatically set up your PATH with Bazel-managed tools:

```bash
direnv allow
```

This command will prompt you to run:

```bash
bazel run //tools:bazel_env
```

Follow any instructions printed by this command. Once complete, development tools will be available on your PATH whenever you're in the project directory.

### 2. Verify Installation

Test that everything is working:

```bash
# Run tests
aspect test //...

# Try formatting (if enabled)
format --help

# Check Bazel version
bazel version
```

## Project Structure

Your generated project will have this structure:

```bash
<project-name>/
├── .aspect/              # Aspect CLI configuration
│   └── cli/             # Custom CLI extensions
├── .bazelrc             # Bazel configuration
├── .envrc               # direnv configuration
├── BUILD                # Root BUILD file
├── MODULE.bazel         # Bazel module dependencies
├── REPO.bazel           # Repository setup (legacy)
├── README.bazel.md      # Developer documentation
├── githooks/            # Git hooks for code quality
│   ├── check-config.sh
│   └── pre-commit
├── requirements/        # Python dependencies (if Python selected)
│   ├── all.in
│   ├── all.txt
│   └── BUILD
├── tools/               # Build and development tools
│   ├── BUILD
│   ├── downloader.cfg
│   ├── preset.bazelrc
│   ├── repin            # Dependency update script
│   ├── tools.lock.json  # Locked tool versions
│   ├── format/          # Formatting configuration
│   ├── lint/            # Linting configuration
│   └── oci/             # Container building (if enabled)
├── package.json         # JavaScript dependencies (if JS selected)
├── pnpm-lock.yaml
├── go.mod               # Go dependencies (if Go selected)
└── pyproject.toml       # Python project config (if Python selected)
```

## Next Steps

Now that you have a project set up:

1. **Read the README.bazel.md** in your project for detailed workflows
2. **Add your code** in language-specific directories
3. **Run Gazelle** to generate BUILD files: `bazel run gazelle`
4. **Enable pre-commit hooks**: `git config core.hooksPath githooks`
5. **Start building**: `bazel build //...`

### Adding Your First Code

Depending on your language choice, see the [Language Support Guide](./languages/README.md) which covers all supported languages:

- JavaScript/TypeScript
- Python  
- Go
- Java
- Kotlin
- C/C++
- Rust
- Shell

## Common First Tasks

### Building Everything

```bash
bazel build //...
```

### Running Tests

```bash
aspect test //...
```

### Formatting Code

```bash
# Format all files
format

# Format specific file
format path/to/file
```

### Linting Code

```bash
aspect lint //...
```

## Getting Help

- **In-Project Documentation**: Check `README.bazel.md` in your project
- **Aspect CLI Docs**: https://docs.aspect.build/cli
- **Bazel Documentation**: https://bazel.build/
- **Community Support**: #aspect-build on [Bazel Slack](https://slack.bazel.build)

## Troubleshooting

If you encounter issues:

1. **Check direnv is enabled**: Run `direnv allow` again
2. **Verify tools are installed**: Run `bazel run //tools:bazel_env`
3. **Clear Bazel cache**: `bazel clean --expunge`
4. **See [Troubleshooting Guide](./troubleshooting.md)** for common issues

---

**Next**: [Development Workflow](./development-workflow.md) | **Up**: [User Guide](./README.md)
