# MEMORY.md

## Project Context
- **Repo:** `bluecentre/aspect-workflows-template`
- **Purpose:** Scaffold for a new Bazel project using Aspect CLI and Aspect Workflows.
- **Features Setup:** Code generation, format/linting (`rules_lint`), version stamping, OCI containers (`rules_oci`).
- **Languages Supported:** JS/TS, Python, Go, Java, Kotlin, C/C++, Rust, Shell.

## Key Files
- `scaffold.yaml`: Contains the scaffolding configuration and feature globs.
- `.aspect/`: Aspect CLI configuration.
- `tools/`: Workspace status scripts, OCI container definitions, platform definitions.