# Aspect Workflows Template Documentation

Welcome to the comprehensive documentation for the Aspect Workflows Template. This template provides a robust foundation for creating multi-language monorepo projects using Bazel and Aspect Workflows.

> 💡 **Tip**: If you're viewing this in Backstage, use the navigation on the left. If you're on GitHub, use the links below.

---

## 🚀 Quick Navigation

| I want to... | Start here |
|-------------|------------|
| **Create my first project** | [Getting Started](./user-guide/getting-started.md) |
| **See common commands** | [Quick Reference](./quick-reference.md) |
| **Understand the architecture** | [Architecture Overview](./contributor-guide/architecture.md) |
| **Fix an issue** | [Troubleshooting](./user-guide/troubleshooting.md) |
| **Find answers** | [FAQ](./faq.md) |
| **See visual diagrams** | [Diagrams](./diagrams.md) |

---

## 📚 Documentation by Role

### 👤 [User Guide](./user-guide/README.md)

For developers using the template to create and work with projects.

| Document | Description |
|----------|-------------|
| [Getting Started](./user-guide/getting-started.md) | Installation and first project |
| [Development Workflow](./user-guide/development-workflow.md) | Daily development tasks |
| [Building & Testing](./user-guide/building-testing.md) | Build and test commands |
| [Formatting & Linting](./user-guide/formatting-linting.md) | Code quality tools |
| [Dependency Management](./user-guide/dependency-management.md) | Managing dependencies |
| [Language Support](./user-guide/languages/README.md) | Multi-language details |
| [Troubleshooting](./user-guide/troubleshooting.md) | Common issues |

### 🔧 [Contributor Guide](./contributor-guide/README.md)

For developers contributing to the template itself.

| Document | Description |
|----------|-------------|
| [Architecture](./contributor-guide/architecture.md) | System design overview |
| [Template System](./contributor-guide/template-system.md) | How scaffold.yaml works |
| [Adding Languages](./contributor-guide/adding-languages.md) | Add new language support |
| [Adding Features](./contributor-guide/adding-features.md) | Extend functionality |
| [Testing](./contributor-guide/testing.md) | Testing strategies |
| [Workflow](./contributor-guide/workflow.md) | Development workflow |

### 🛠️ [Administrator Guide](./admin-guide/README.md)

For maintainers managing projects and infrastructure.

| Document | Description |
|----------|-------------|
| [Backstage Integration](./admin-guide/backstage-integration.md) | Self-service templates |
| [CI/CD](./admin-guide/ci-cd.md) | Pipeline configuration |
| [Dependency Management](./admin-guide/dependency-management.md) | Keeping deps updated |
| [Maintenance](./admin-guide/maintenance.md) | Routine maintenance |
| [Security](./admin-guide/security.md) | Security practices |
| [Troubleshooting](./admin-guide/troubleshooting.md) | Admin troubleshooting |

---

## 📖 Reference Documents

| Document | Description |
|----------|-------------|
| [Quick Reference](./quick-reference.md) | Command cheatsheet for all common tasks |
| [FAQ](./faq.md) | Frequently asked questions and answers |
| [Visual Diagrams](./diagrams.md) | Architecture and workflow diagrams |
| [Backstage Quick Start](./BACKSTAGE-QUICK-START.md) | Fast path for Backstage users |

---

## ✨ What This Template Provides

| Category | Features |
|----------|----------|
| **Languages** | Python, JavaScript/TypeScript, Go, Java, Kotlin, C/C++, Rust, Shell |
| **Build System** | Bazel with bzlmod, Gazelle code generation |
| **Code Quality** | Formatting (Ruff, Prettier, gofumpt), Linting (rules_lint) |
| **Testing** | Integrated test frameworks per language |
| **Containers** | OCI image building with rules_oci |
| **CI/CD** | Version stamping, release automation |
| **Developer Experience** | direnv integration, hermetic tooling via bazel-env |

---

## 🔗 External Resources

- **[Scaffold Documentation](https://hay-kot.github.io/scaffold/)** - Template engine
- **[Bazel Documentation](https://bazel.build/docs)** - Build system
- **[Aspect CLI](https://docs.aspect.build/cli/)** - Enhanced Bazel CLI
- **[rules_lint](https://github.com/aspect-build/rules_lint)** - Linting framework
- **[Backstage](https://backstage.io/docs)** - Developer portal

---

## 📄 License

Apache 2.0 - See [LICENSE](../LICENSE) for details.
