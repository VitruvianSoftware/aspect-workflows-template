# TOOLS.md - Tool Notes

## Aspect CLI
- Use `aspect` instead of `bazel` where possible for enhanced features.
- Common commands:
  - `aspect build //...`
  - `aspect test //...`
  - `aspect lint //...`
  - `aspect run //path/to:target`

## Scaffold
- This repo was generated using `hay-kot/scaffold`.

## Rules
- Familiarize with `rules_oci` for containers, `rules_lint` for formatting/linting, and language-specific rules (e.g., `gazelle` for Go/Python/JS).