# Architecture Overview

This document provides a comprehensive overview of the Aspect Workflows Template architecture, explaining how the components work together to generate production-ready Bazel projects.

## System Architecture

```mermaid
graph TB
    subgraph "Template System"
        ScaffoldYAML[scaffold.yaml]
        TemplateFiles[Template Files]
        PostHook[post_scaffold Hook]
    end
    
    subgraph "Generation Process"
        CLI[Scaffold CLI]
        Questions[User Prompts]
        Processing[Template Processing]
        Validation[File Generation]
    end
    
    subgraph "Generated Project"
        Bazel[Bazel Workspace]
        ModuleBazel[MODULE.bazel]
        BuildFiles[BUILD Files]
        Config[Configuration Files]
        Tools[Development Tools]
    end
    
    subgraph "Runtime Components"
        AspectCLI[Aspect CLI]
        Gazelle[Gazelle]
        RulesLint[rules_lint]
        Multitool[Multitool]
        BazelEnv[bazel-env.bzl]
    end
    
    ScaffoldYAML --> CLI
    CLI --> Questions
    Questions --> Processing
    Processing --> TemplateFiles
    Processing --> Validation
    Validation --> PostHook
    PostHook --> Bazel
    
    Bazel --> ModuleBazel
    Bazel --> BuildFiles
    Bazel --> Config
    Bazel --> Tools
    
    Tools --> AspectCLI
    Tools --> Gazelle
    Tools --> RulesLint
    Tools --> Multitool
    Tools --> BazelEnv
```

## Core Components

### 1. Template Configuration (scaffold.yaml)

The central configuration file that defines the entire template system.

**Key Sections:**

- **metadata**: Template versioning requirements
- **messages**: Pre/post generation user guidance
- **questions**: Interactive user prompts
- **features**: Conditional file inclusion rules
- **computed**: Derived values from user selections
- **presets**: Predefined configuration sets

**Flow:**

```mermaid
flowchart LR
    User[User Input] --> Questions
    Questions --> Answers[User Answers]
    Answers --> Scaffold[.Scaffold namespace]
    Scaffold --> Computed[.Computed namespace]
    Computed --> Features[Feature Flags]
    Features --> Files[File Filtering]
```

### 2. Template Files ({{ .ProjectSnake }}/)

The source directory containing templated files that will be processed and copied to the generated project.

**Template Syntax:**

```go
// Conditional inclusion
{{ if .Scaffold.lint }}
  // lint-specific code
{{ end }}

// Computed values
{{ if .Computed.javascript }}
  // JavaScript-specific code
{{ end }}

// Variable substitution
bazel_dep(name = "rules_python", version = "{{ .PythonVersion }}")
```

### 3. Post-Processing (hooks/post_scaffold)

A bash script that runs after file generation to:

1. Format generated Bazel files with buildifier
2. Run custom format tools
3. Update dependency lock files (repin)
4. Initialize package managers (pnpm for JavaScript)

### 4. Bazel Module System (MODULE.bazel)

The generated `MODULE.bazel` uses Bazel's bzlmod system for dependency management.

**Dependency Categories:**

- **Core**: bazel_skylib, rules_multitool
- **Language Rules**: rules_python, rules_go, rules_js, etc.
- **Development**: aspect_rules_lint, buildifier_prebuilt
- **Container**: rules_oci (optional)
- **Utilities**: bazel_env.bzl, bazelrc-preset.bzl

**Extension System:**

```mermaid
graph LR
    ModuleBazel[MODULE.bazel] --> Extensions[Bazel Extensions]
    Extensions --> GoSDK[go_sdk]
    Extensions --> NPM[npm/pnpm]
    Extensions --> Pip[pip]
    Extensions --> Maven[maven]
    Extensions --> OCI[oci]
    Extensions --> LLVM[llvm toolchain]
```

## Language Support Architecture

### Multi-Language Strategy

Each language follows a consistent pattern:

```mermaid
graph TB
    User[User Selects Language] --> Computed[Computed Flag Set]
    Computed --> Deps[Dependencies Added]
    Computed --> Rules[Rule Files Included]
    Computed --> Tools[Language Tools Added]
    
    Deps --> ModuleBazel[MODULE.bazel]
    Rules --> BuildFiles[BUILD Files]
    Tools --> ToolsBuild[tools/BUILD]
    
    ModuleBazel --> Gazelle[Gazelle Integration]
    BuildFiles --> Gazelle
    ToolsBuild --> BazelEnv[bazel-env.bzl]
```

### Language Components

Each language integration includes:

1. **MODULE.bazel entries**: Language rules and toolchains
2. **BUILD file templates**: Language-specific targets
3. **Tool exports**: CLI tools added to PATH via bazel-env
4. **Gazelle support**: Automatic BUILD file generation
5. **Linting rules**: Language-specific linters
6. **Package management**: Dependency resolution

## Development Environment Architecture

### bazel-env.bzl Integration

```mermaid
graph TB
    Developer[Developer] --> Direnv[direnv]
    Direnv --> BazelEnv[bazel run //tools:bazel_env]
    BazelEnv --> BinTree[bin/ tree generation]
    BinTree --> Path[Added to PATH]
    
    Path --> Tools[Development Tools]
    Tools --> Buildifier
    Tools --> Format
    Tools --> Gazelle
    Tools --> LanguageTools[Language CLIs]
    
    BazelEnv --> Toolchains[Bazel Toolchains]
    Toolchains --> Python
    Toolchains --> Node
    Toolchains --> Go
    Toolchains --> Java
```

**Benefits:**

- No manual tool installation required
- Version consistency across team
- Hermetic tool execution
- Automatic PATH management via direnv

## Code Quality Architecture

### Formatting and Linting Pipeline

```mermaid
flowchart LR
    Developer[Developer] --> PreCommit[Git Pre-Commit Hook]
    PreCommit --> Format[format tool]
    
    Format --> Buildifier[Bazel files]
    Format --> Prettier[JS/TS files]
    Format --> Ruff[Python files]
    Format --> Gofmt[Go files]
    Format --> More[Language-specific formatters]
    
    Developer --> Manual[Manual Commands]
    Manual --> AspectLint[aspect lint //...]
    AspectLint --> Linters[Language Linters]
    
    Linters --> ESLint[JavaScript]
    Linters --> Ruff[Python]
    Linters --> Nogo[Go]
    Linters --> PMD[Java]
    Linters --> Ktlint[Kotlin]
    Linters --> ClangTidy[C/C++]
```

### rules_lint Integration

The template uses [rules_lint](https://github.com/aspect-build/rules_lint) which:

1. Runs linters as Bazel aspects
2. Caches linter results
3. Produces report files
4. Integrates with Aspect CLI for UX

## Dependency Management Architecture

### Python Dependencies

```mermaid
graph TB
    PyProject[pyproject.toml] --> Repin[./tools/repin]
    Repin --> Runtime[requirements/runtime.txt]
    Repin --> All[requirements/all.txt]
    Repin --> Manifest[gazelle_python_manifest]
    
    Runtime --> Pip[pip.parse]
    All --> Pip
    Pip --> WheelRepo["@pip repository"]
    
    Manifest --> Gazelle[gazelle for Python]
    Gazelle --> BuildGen[BUILD file generation]
```

### JavaScript Dependencies

```mermaid
graph TB
    PackageJson[package.json] --> PNPM[pnpm]
    PNPM --> Lock[pnpm-lock.yaml]
    Lock --> NPMTranslate[npm.npm_translate_lock]
    NPMTranslate --> NPMRepo["@npm repository"]
    
    NPMRepo --> LinkAll[npm_link_all_packages]
    LinkAll --> NodeModules[node_modules]
```

### Go Dependencies

```mermaid
graph TB
    GoMod[go.mod] --> GoModTidy[go mod tidy]
    GoModTidy --> GoSum[go.sum]
    GoSum --> GoDeps[go_deps extension]
    GoDeps --> BazelModTidy[bazel mod tidy]
    BazelModTidy --> UseRepo[use_repo declarations]
```

## Build File Generation (Gazelle)

### Gazelle Architecture

```mermaid
graph TB
    Source[Source Files] --> Gazelle[bazel run gazelle]
    
    Gazelle --> GoPlugin[Go Plugin]
    Gazelle --> PyPlugin[Python Plugin]
    Gazelle --> JSPlugin[JavaScript Plugin]
    
    GoPlugin --> GoConfig[gazelle:prefix]
    PyPlugin --> PyManifest[modules_mapping]
    JSPlugin --> JSConfig[gazelle:js_npm_package_target_name]
    
    GoPlugin --> GoBuild[Go BUILD targets]
    PyPlugin --> PyBuild[Python BUILD targets]
    JSPlugin --> JSBuild[JS BUILD targets]
```

### Configuration

Gazelle is configured via directives in BUILD files:

```python
# gazelle:prefix github.com/example/project
# gazelle:build_file_name BUILD
# gazelle:map_kind py_binary py_binary @aspect_rules_py//py:defs.bzl
# gazelle:exclude **/*.venv
```

## Container Building Architecture (Optional)

### OCI Image Generation

```mermaid
graph TB
    Binary[Application Binary] --> ImageRule[*_image rule]
    
    ImageRule --> PyImage[py3_image]
    ImageRule --> GoImage[go_image]
    
    BaseImage[Base Image] --> OCI[rules_oci]
    PyImage --> OCI
    GoImage --> OCI
    
    OCI --> Tarball[OCI Tarball]
    Tarball --> Registry[Container Registry]
```

### Image Types

1. **Python Images**: Ubuntu base with Python runtime
2. **Go Images**: Distroless base for minimal images

## Version Stamping Architecture (Optional)

```mermaid
graph TB
    Git[Git Repository] --> WorkspaceStatus[workspace_status.sh]
    WorkspaceStatus --> Vars[Stamp Variables]
    
    Vars --> Commit[STABLE_GIT_COMMIT]
    Vars --> Version[STABLE_MONOREPO_VERSION]
    
    Build[bazel build --config=release] --> Stamp[Stamping]
    Stamp --> Vars
    Stamp --> Template[expand_template]
    Template --> Output[Stamped Artifacts]
```

## Extensibility Points

### 1. Custom Rules

Add custom rules in `tools/` directories:

```bash
tools/
├── custom_rule.bzl
└── BUILD
```

### 2. Aspect CLI Extensions

Add custom commands in `.aspect/cli/`:

```bash
.aspect/cli/
├── go_image.star
├── py3_image.star
└── custom_command.star
```

### 3. Linter Configuration

Extend linting in `tools/lint/linters.bzl`:

```python
my_linter = lint_my_linter_aspect(
    binary = Label(":my_linter"),
    configs = [Label("//:.my-linter-config")],
)
```

### 4. Format Tool Integration

Add formatters in `tools/format/BUILD`:

```python
multirun(
    name = "format",
    commands = [
        ":my_formatter",
        # ...
    ],
)
```

## Performance Considerations

### Caching Strategy

- **Remote Caching**: Can be configured in `.bazelrc`
- **Local Cache**: Bazel's built-in caching
- **Action Cache**: Linter results cached as Bazel actions
- **Tool Lockfiles**: Deterministic tool versions

### Scalability

- **Incremental Builds**: Bazel rebuilds only changed targets
- **Parallel Execution**: Multiple actions run concurrently
- **Lazy Loading**: Bzlmod loads only needed dependencies
- **Optimized Toolchains**: Hermetic toolchains reduce setup time

## Security Considerations

- **Hermetic Builds**: All dependencies fetched through Bazel
- **Checksum Verification**: Integrity checks on downloaded artifacts
- **Lockfiles**: Pin exact dependency versions
- **Sandboxing**: Bazel runs actions in sandboxes

---

**Next**: [Template System](./template-system.md) | **Back**: [Contributor Guide](./README.md)
