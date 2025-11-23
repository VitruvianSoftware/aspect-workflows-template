# Aspect Workflows Template

[![Bazel](https://img.shields.io/badge/Bazel-43A047?logo=bazel&logoColor=white)](https://bazel.build)
[![Scaffold](https://img.shields.io/badge/Scaffold-Template-blue)](https://hay-kot.github.io/scaffold/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

> **Multi-language Bazel monorepo template generator** powered by [Scaffold](https://hay-kot.github.io/scaffold/)

Generate production-ready, multi-language Bazel monorepos with best practices built-in. Supports **8 programming languages**, integrated tooling, container builds, and optional Backstage integration for self-service project creation.

---

## 📚 Documentation

| Guide | Description |
|-------|-------------|
| **[📖 Full Documentation](docs/overview.md)** | Complete documentation index |
| [Getting Started](docs/user-guide/getting-started.md) | Create your first project |
| [Quick Reference](docs/quick-reference.md) | Command cheatsheet |
| [FAQ](docs/faq.md) | Frequently asked questions |
| [Visual Diagrams](docs/diagrams.md) | Architecture diagrams |

### By Role

| Role | Guide | Topics |
|------|-------|--------|
| 👤 **Developer** | [User Guide](docs/user-guide/README.md) | Building, testing, formatting, linting |
| 🔧 **Contributor** | [Contributor Guide](docs/contributor-guide/README.md) | Template system, adding features |
| 🛠️ **Admin** | [Admin Guide](docs/admin-guide/README.md) | CI/CD, maintenance, Backstage integration |

---

## 🚀 Quick Start

### Install Scaffold

```shell
brew tap hay-kot/scaffold-tap && brew install scaffold
# OR
go install github.com/hay-kot/scaffold@latest
```

### Generate a Project

```bash
# Interactive mode
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new github.com/BlueCentre/aspect-workflows-template

# With preset (non-interactive)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=py --output-dir=my-project \
  github.com/BlueCentre/aspect-workflows-template

cd my-project
bazel test //...
```

---

## 🎯 Two Modes of Operation

### 1. Direct Project Generation (Default)

Generate a complete, ready-to-use Bazel monorepo:

```bash
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=py --output-dir=my-python-project .
cd my-python-project
bazel test //...
```

**Use this for**: Starting new projects, prototyping, local development

### 2. Backstage Template Generation

Generate a [Backstage](https://backstage.io) software template for self-service project creation:

```bash
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-py --output-dir=templates/aspect-python .
```

**Use this for**: Creating reusable templates in Backstage, standardizing team workflows

> 📖 See [Backstage Quick Start](docs/BACKSTAGE-QUICK-START.md) and [Backstage Integration Guide](docs/admin-guide/backstage-integration.md) for details.

---

## 📦 Available Presets

<table>
<tr>
<th>Direct Generation</th>
<th>Backstage Templates</th>
</tr>
<tr>
<td>

| Preset | Languages |
|--------|-----------|
| `py` | Python |
| `js` | JavaScript/TypeScript |
| `go` | Go |
| `java` | Java |
| `kotlin` | Kotlin |
| `cpp` | C/C++ |
| `rust` | Rust |
| `shell` | Shell (Bash) |
| `kitchen-sink` | All languages |
| `minimal` | No languages |

</td>
<td>

| Preset | Languages |
|--------|-----------|
| `backstage-py` | Python |
| `backstage-js` | JavaScript/TypeScript |
| `backstage-go` | Go |
| `backstage-java` | Java |
| `backstage-kotlin` | Kotlin |
| `backstage-cpp` | C/C++ |
| `backstage-rust` | Rust |
| `backstage-shell` | Shell (Bash) |
| `backstage-kitchen-sink` | All languages |
| `backstage-minimal` | No languages |

</td>
</tr>
</table>

---

## ✨ Features

| Category | Features |
|----------|----------|
| **Languages** | Python, JavaScript/TypeScript, Go, Java, Kotlin, C/C++, Rust, Shell |
| **Build System** | Bazel with bzlmod, Gazelle code generation |
| **Code Quality** | Formatting (Ruff, Prettier, gofumpt, etc.), Linting (rules_lint) |
| **Testing** | Integrated test frameworks per language |
| **Containers** | OCI image building with rules_oci |
| **CI/CD** | Version stamping, release automation |
| **Developer Experience** | direnv integration, hermetic tooling via bazel-env |

---

## 🏗️ Generated Project Structure

```
my-project/
├── .bazelrc              # Bazel configuration
├── BUILD                 # Root build file
├── MODULE.bazel          # Bazel module dependencies
├── pyproject.toml        # Python configuration (if Python enabled)
├── package.json          # JS configuration (if JS enabled)
├── go.mod                # Go configuration (if Go enabled)
├── requirements/         # Python lockfiles
├── tools/
│   ├── format/           # Formatting configuration
│   ├── lint/             # Linting configuration
│   └── repin             # Dependency update script
└── githooks/             # Git hooks for pre-commit
```

---

## 🔧 Common Commands

After generating a project:

```bash
# Setup
direnv allow                    # Enable environment
bazel run //tools:bazel_env     # Install tools to PATH

# Development
bazel build //...               # Build everything
bazel test //...                # Run tests
bazel run gazelle               # Regenerate BUILD files

# Code Quality
format                          # Format all files
aspect lint //...               # Run linters (requires Aspect CLI)
aspect lint --fix //...         # Auto-fix lint issues

# Dependencies
./tools/repin                   # Update Python/Java lockfiles
pnpm install                    # Update JS dependencies
go mod tidy && bazel mod tidy   # Update Go dependencies
```

> 📖 See [Quick Reference](docs/quick-reference.md) for the complete command list.

---

## 📖 Learn More

- **[Full Documentation](docs/overview.md)** - Complete documentation index
- **[Architecture Overview](docs/contributor-guide/architecture.md)** - How the template system works
- **[Troubleshooting](docs/user-guide/troubleshooting.md)** - Common issues and solutions
- **[Visual Diagrams](docs/diagrams.md)** - Architecture and workflow diagrams

---

## 🤝 Contributing

See the [Contributor Guide](docs/contributor-guide/README.md) to learn how to:
- Add new language support
- Extend the template system
- Test your changes

---

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) for details.
