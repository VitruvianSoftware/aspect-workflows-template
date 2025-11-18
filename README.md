# Scaffold for basic Aspect Workflows project

## With Aspect CLI
Install the Aspect CLI: https://docs.aspect.build/cli/install

Then add the `init` task with

```
aspect axl add gh:aspect-build/aspect-workflows-template
```

Finally, run `aspect init`

## Manually

Install [scaffold](https://hay-kot.github.io/scaffold/) like so:

```shell
% brew tap hay-kot/scaffold-tap
% brew install scaffold
# OR
% go install github.com/hay-kot/scaffold@latest
```

And then create a new project like so:

```shell
% scaffold new https://github.com/aspect-build/aspect-workflows-template
```

## Two Modes of Operation

This template generator supports **two workflows**:

### 1. Direct Project Generation (Default)

Generate a complete, ready-to-use Bazel monorepo:

```bash
# Interactive mode
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new .

# With preset
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=py --output-dir=my-python-project .
cd my-python-project
bazel test //...
```

**Use this for**: Quickly starting new projects, prototyping, local development

### 2. Backstage Template Generation

Generate a Backstage software template for self-service project creation:

```bash
# Interactive mode - answer "yes" to "Generate as a Backstage template?"
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new .

# With backstage preset
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-py --output-dir=templates/aspect-python .
```

**Use this for**: Creating reusable templates in Backstage, standardizing team workflows

> **Note**: The `SCAFFOLD_SETTINGS_RUN_HOOKS=always` environment variable is required to run the post-generation hooks that create symlinks in the skeleton/ directory for Backstage templates.

## Quick Links

- 📚 [Full Backstage Integration Guide](docs/admin-guide/backstage-integration.md)
- ⚡ [Backstage Quick Start](docs/BACKSTAGE-QUICK-START.md)
- 🔧 [User Guide](docs/user-guide/)
- 🏗️ [Contributor Guide](docs/contributor-guide/)

## Available Presets

### Direct Generation
- `py` - Python
- `js` - JavaScript/TypeScript
- `go` - Go
- `java`, `kotlin`, `cpp`, `rust`, `shell`
- `kitchen-sink` - All languages
- `minimal` - Bare bones

### Backstage Templates
- `backstage-py` - Python template
- `backstage-js` - JavaScript/TypeScript template
- `backstage-go` - Go template
- `backstage-java` - Java template
- `backstage-kotlin` - Kotlin template
- `backstage-cpp` - C++ template
- `backstage-rust` - Rust template
- `backstage-shell` - Shell template
- `backstage-kitchen-sink` - All languages template
- `backstage-minimal` - Bare bones template

## Examples

```bash
# Create a Python microservice (direct)
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=py --output-dir=payment-service .

# Create a Go Backstage template
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-go --output-dir=templates/go-service .

# Create all Backstage templates for your org
for lang in py js go java kotlin; do
  SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-$lang --output-dir=templates/aspect-$lang .
done
```

See [docs/BACKSTAGE-QUICK-START.md](docs/BACKSTAGE-QUICK-START.md) for more examples.
