# Visual Architecture Diagrams

This document contains comprehensive visual diagrams to help understand the Aspect Workflows Template architecture and workflows.

## Table of Contents

1. [Overall System Architecture](#overall-system-architecture)
2. [Template Generation Flow](#template-generation-flow)
3. [Build System Architecture](#build-system-architecture)
4. [Development Workflow](#development-workflow)
5. [Dependency Management](#dependency-management)
6. [Code Quality Pipeline](#code-quality-pipeline)
7. [Multi-Language Integration](#multi-language-integration)
8. [CI/CD Architecture](#cicd-architecture)

## Overall System Architecture

### High-Level Overview

```mermaid
graph TB
    subgraph "Template Repository"
        A[scaffold.yaml]
        B[Template Files]
        C[post_scaffold Hook]
        D[User Stories]
    end
    
    subgraph "Generation Process"
        E[Scaffold CLI]
        F[User Interaction]
        G[File Processing]
    end
    
    subgraph "Generated Project"
        H[Bazel Workspace]
        I[Development Environment]
        J[CI/CD Configuration]
    end
    
    subgraph "Developer Workflow"
        K[Write Code]
        L[Build/Test]
        M[Format/Lint]
        N[Commit/Deploy]
    end
    
    A --> E
    B --> E
    E --> F
    F --> G
    G --> C
    C --> H
    H --> I
    I --> K
    K --> L
    L --> M
    M --> N
    N --> J
```

### Component Relationships

```mermaid
graph LR
    subgraph "Core Components"
        Bazel[Bazel Build System]
        AspectCLI[Aspect CLI]
        Gazelle[Gazelle Code Gen]
        RulesLint[rules_lint]
    end
    
    subgraph "Language Support"
        Python[Python Rules]
        Go[Go Rules]
        JS[JavaScript Rules]
        Java[Java Rules]
        Kotlin[Kotlin Rules]
        CPP[C++ Rules]
        Rust[Rust Rules]
        Shell[Shell Rules]
    end
    
    subgraph "Development Tools"
        BazelEnv[bazel-env.bzl]
        Direnv[direnv]
        Multitool[Multitool]
    end
    
    subgraph "Quality Tools"
        Formatters[Formatters]
        Linters[Linters]
        Tests[Test Frameworks]
    end
    
    Bazel --> AspectCLI
    Bazel --> Gazelle
    Bazel --> RulesLint
    
    AspectCLI --> Python
    AspectCLI --> Go
    AspectCLI --> JS
    AspectCLI --> Java
    AspectCLI --> Kotlin
    AspectCLI --> CPP
    AspectCLI --> Rust
    AspectCLI --> Shell
    
    BazelEnv --> Direnv
    BazelEnv --> Multitool
    
    RulesLint --> Formatters
    RulesLint --> Linters
    Bazel --> Tests
```

## Template Generation Flow

### Detailed Generation Process

```mermaid
sequenceDiagram
    actor User
    participant Scaffold
    participant Config as scaffold.yaml
    participant Templates
    participant PostHook as post_scaffold
    participant Project
    
    User->>Scaffold: scaffold new [options]
    Scaffold->>Config: Load configuration
    Config-->>Scaffold: Questions, Features, Presets
    
    Scaffold->>User: Display pre-message
    Scaffold->>User: Ask questions
    User-->>Scaffold: Provide answers
    
    Scaffold->>Scaffold: Compute derived values
    Scaffold->>Config: Evaluate feature flags
    Config-->>Scaffold: Active features list
    
    Scaffold->>Templates: Filter by features
    Templates-->>Scaffold: Filtered file list
    
    loop For each template file
        Scaffold->>Templates: Read template
        Templates-->>Scaffold: Template content
        Scaffold->>Scaffold: Process Go template
        Scaffold->>Project: Write processed file
    end
    
    Scaffold->>PostHook: Execute hook
    PostHook->>Project: Format files
    PostHook->>Project: Run buildifier
    PostHook->>Project: Update lockfiles
    PostHook->>Project: Initialize packages
    PostHook-->>Scaffold: Completion status
    
    Scaffold->>User: Display post-message
    Scaffold->>User: Project ready!
```

### Feature Flag Resolution

```mermaid
flowchart TD
    Start[User Selections] --> Questions{Questions Answered}
    Questions --> Scaffold[.Scaffold namespace]
    
    Scaffold --> Lang{Language Selection}
    Lang --> |has check| Computed[.Computed namespace]
    
    Scaffold --> Features{Feature Checks}
    Features --> Lint[Lint Flag]
    Features --> OCI[OCI Flag]
    Features --> Stamp[Stamp Flag]
    Features --> Codegen[Codegen Flag]
    
    Computed --> JS[.Computed.javascript]
    Computed --> Py[.Computed.python]
    Computed --> GoLang[.Computed.go]
    Computed --> JavaLang[.Computed.java]
    
    JS --> Combine[Combine Conditions]
    Py --> Combine
    GoLang --> Combine
    JavaLang --> Combine
    Lint --> Combine
    OCI --> Combine
    Stamp --> Combine
    Codegen --> Combine
    
    Combine --> Filter[File Filtering]
    Filter --> Include{Include File?}
    Include -->|Yes| Copy[Copy to Output]
    Include -->|No| Skip[Skip File]
    
    Copy --> Done[Generation Complete]
    Skip --> Done
```

## Build System Architecture

### Bazel Module System (bzlmod)

```mermaid
graph TB
    subgraph "MODULE.bazel"
        ModFile[Module Declaration]
        BazelDeps[bazel_dep statements]
        Extensions[use_extension calls]
    end
    
    subgraph "Extensions"
        GoSDK[go_sdk]
        NPM[npm/pnpm]
        Pip[pip]
        Maven[maven]
        OCI[oci]
    end
    
    subgraph "External Dependencies"
        GoMod[Go Modules]
        PyPi[Python Packages]
        NPMReg[npm Registry]
        MavenCentral[Maven Central]
        ContainerReg[Container Images]
    end
    
    subgraph "Build Targets"
        Libraries[Libraries]
        Binaries[Binaries]
        Tests[Tests]
        Images[Container Images]
    end
    
    ModFile --> BazelDeps
    BazelDeps --> Extensions
    
    Extensions --> GoSDK
    Extensions --> NPM
    Extensions --> Pip
    Extensions --> Maven
    Extensions --> OCI
    
    GoSDK --> GoMod
    NPM --> NPMReg
    Pip --> PyPi
    Maven --> MavenCentral
    OCI --> ContainerReg
    
    GoMod --> Libraries
    NPMReg --> Libraries
    PyPi --> Libraries
    MavenCentral --> Libraries
    
    Libraries --> Binaries
    Libraries --> Tests
    Binaries --> Images
```

### Build Execution Flow

```mermaid
flowchart LR
    Command[bazel build //...] --> Analysis[Analysis Phase]
    
    Analysis --> Deps[Load Dependencies]
    Analysis --> Rules[Load Rules]
    Analysis --> Graph[Build Action Graph]
    
    Graph --> Execution[Execution Phase]
    
    Execution --> Cache{Check Cache}
    Cache -->|Hit| Reuse[Reuse Cached Output]
    Cache -->|Miss| Execute[Execute Action]
    
    Execute --> Sandbox[Sandboxed Execution]
    Sandbox --> Output[Generate Output]
    Output --> StoreCache[Store in Cache]
    
    Reuse --> Final[Output Artifacts]
    StoreCache --> Final
```

## Development Workflow

### Daily Developer Cycle

```mermaid
stateDiagram-v2
    [*] --> EditCode: Start Work
    
    EditCode --> RunGazelle: Add/Change Files
    RunGazelle --> Build: Generate BUILD
    
    Build --> BuildSuccess: Success
    Build --> BuildFail: Failure
    
    BuildFail --> EditCode: Fix Errors
    
    BuildSuccess --> Test: Run Tests
    Test --> TestSuccess: Pass
    Test --> TestFail: Fail
    
    TestFail --> EditCode: Fix Tests
    
    TestSuccess --> Format: Format Code
    Format --> Lint: Lint Code
    
    Lint --> LintSuccess: No Issues
    Lint --> LintFail: Has Issues
    
    LintFail --> EditCode: Fix Issues
    
    LintSuccess --> Commit: Commit
    Commit --> Push: Push Changes
    Push --> [*]: Complete
```

### Development Environment Setup

```mermaid
sequenceDiagram
    actor Dev as Developer
    participant Git
    participant Direnv
    participant Bazel
    participant Tools
    
    Dev->>Git: Clone repository
    Git-->>Dev: Repository cloned
    
    Dev->>Direnv: cd into project
    Direnv->>Direnv: Check .envrc
    Direnv-->>Dev: Environment not loaded
    
    Dev->>Direnv: direnv allow
    Direnv->>Bazel: Check bazel-out/bazel_env
    Bazel-->>Direnv: Not found
    Direnv-->>Dev: Run bazel run //tools:bazel_env
    
    Dev->>Bazel: bazel run //tools:bazel_env
    Bazel->>Tools: Build tool symlinks
    Tools->>Tools: Create bin/ tree
    Tools-->>Bazel: Complete
    Bazel-->>Dev: Tools ready
    
    Dev->>Direnv: direnv reload
    Direnv->>Direnv: Add bin/ to PATH
    Direnv-->>Dev: Environment loaded
    
    Dev->>Tools: Use tools (format, pnpm, etc)
    Tools-->>Dev: Tools work directly
```

## Dependency Management

### Python Dependency Flow

```mermaid
flowchart TB
    Start[Developer] --> Edit[Edit pyproject.toml]
    Edit --> Repin[Run ./tools/repin]
    
    Repin --> Runtime[Update requirements/runtime.txt]
    Repin --> All[Update requirements/all.txt]
    Repin --> Manifest[Update gazelle_python_manifest]
    
    Runtime --> Parse[pip.parse]
    All --> Parse
    
    Parse --> Wheels[Download Wheels]
    Wheels --> Repo["@pip repository"]
    
    Manifest --> Modules[modules_mapping]
    Modules --> GazelleConfig[Gazelle Configuration]
    
    Repo --> BuildFile[BUILD files]
    GazelleConfig --> BuildFile
    
    BuildFile --> Import[Import in Code]
    Import --> Build[bazel build]
```

### JavaScript Dependency Flow

```mermaid
flowchart TB
    Start[Developer] --> Add[pnpm add package]
    
    Add --> PackageJSON[Update package.json]
    Add --> Lock[Update pnpm-lock.yaml]
    
    PackageJSON --> Translate[npm_translate_lock]
    Lock --> Translate
    
    Translate --> NPMRepo["@npm repository"]
    
    NPMRepo --> Link[npm_link_all_packages]
    Link --> NodeModules[node_modules targets]
    
    NodeModules --> BuildFile[BUILD files]
    BuildFile --> Import[import in code]
    Import --> Build[bazel build]
```

### Go Dependency Flow

```mermaid
flowchart TB
    Start[Developer] --> AddImport[Add import to code]
    
    AddImport --> GoModTidy[go mod tidy]
    GoModTidy --> GoMod[Update go.mod]
    GoModTidy --> GoSum[Update go.sum]
    
    GoMod --> GoDeps[go_deps extension]
    GoSum --> GoDeps
    
    GoDeps --> BazelModTidy[bazel mod tidy]
    BazelModTidy --> UseRepo[Update use_repo]
    
    UseRepo --> Gazelle[bazel run gazelle]
    Gazelle --> BuildFile[Update BUILD files]
    
    BuildFile --> Build[bazel build]
```

## Code Quality Pipeline

### Format and Lint Workflow

```mermaid
flowchart LR
    subgraph "Pre-Commit"
        Stage[Git Staged Files] --> PreHook[Pre-Commit Hook]
        PreHook --> FormatStaged[Format Staged Files]
        FormatStaged --> Success{Success?}
        Success -->|No| Abort[Abort Commit]
        Success -->|Yes| Allow[Allow Commit]
    end
    
    subgraph "Manual Format"
        Manual[Developer] --> FormatCmd[Run format]
        FormatCmd --> Buildifier[Bazel Files]
        FormatCmd --> Prettier[JS/TS Files]
        FormatCmd --> Ruff[Python Files]
        FormatCmd --> More[More Formatters...]
    end
    
    subgraph "Linting"
        LintCmd[aspect lint //...] --> Aspect[Bazel Aspects]
        Aspect --> ESLint[ESLint]
        Aspect --> Ruff[Ruff]
        Aspect --> Nogo[Go nogo]
        Aspect --> PMD[Java PMD]
        Aspect --> Reports[Collect Reports]
        Reports --> Display[Display Results]
    end
    
    Allow --> CI[CI Pipeline]
    Manual --> CI
    CI --> LintCmd
```

### Linting Architecture

```mermaid
graph TB
    subgraph "Bazel Build Graph"
        Target[Target]
        Deps[Dependencies]
    end
    
    subgraph "Linting Aspect"
        Aspect[rules_lint Aspect]
        Config[Linter Config]
    end
    
    subgraph "Linters"
        ESLint[ESLint]
        Ruff[Ruff]
        Nogo[nogo]
        PMD[PMD]
        Ktlint[ktlint]
        ClangTidy[clang-tidy]
    end
    
    subgraph "Outputs"
        Reports[Report Files]
        Cache[Bazel Cache]
    end
    
    Target --> Aspect
    Deps --> Aspect
    Config --> Aspect
    
    Aspect --> ESLint
    Aspect --> Ruff
    Aspect --> Nogo
    Aspect --> PMD
    Aspect --> Ktlint
    Aspect --> ClangTidy
    
    ESLint --> Reports
    Ruff --> Reports
    Nogo --> Reports
    PMD --> Reports
    Ktlint --> Reports
    ClangTidy --> Reports
    
    Reports --> Cache
    Cache --> Display[aspect lint display]
```

## Multi-Language Integration

### Polyglot Project Structure

```mermaid
graph TB
    Root[Project Root]
    
    Root --> Python[python/]
    Root --> Go[go/]
    Root --> Web[web/]
    Root --> Shared[shared/]
    
    Python --> PyAPI[api/]
    Python --> PyLib[lib/]
    
    Go --> GoCmd[cmd/]
    Go --> GoPkg[pkg/]
    
    Web --> WebSrc[src/]
    Web --> WebPkg[packages/]
    
    PyAPI --> Service[Service Targets]
    PyLib --> PyLibTargets[Library Targets]
    GoCmd --> CLITargets[CLI Targets]
    GoPkg --> GoLibTargets[Library Targets]
    WebSrc --> AppTargets[App Targets]
    WebPkg --> JSLibTargets[Library Targets]
    
    Service -.-> GoLibTargets
    AppTargets -.-> PyLibTargets
    CLITargets -.-> JSLibTargets
```

### Cross-Language Dependencies

```mermaid
flowchart LR
    subgraph "Language Boundaries"
        PythonApp[Python Application]
        GoService[Go Service]
        JSFrontend[JS Frontend]
    end
    
    subgraph "Integration Methods"
        API[REST/gRPC API]
        CLI[Command Line]
        Files[Shared Files]
        Proto[Protocol Buffers]
    end
    
    subgraph "Bazel Integration"
        Data[data attribute]
        Deps[deps attribute]
        RunFiles[Runtime Files]
    end
    
    PythonApp -->|calls| API
    API -->|implemented by| GoService
    
    JSFrontend -->|invokes| CLI
    CLI -->|implemented by| PythonApp
    
    PythonApp -->|reads| Files
    GoService -->|writes| Files
    
    PythonApp -->|uses| Proto
    GoService -->|uses| Proto
    JSFrontend -->|uses| Proto
    
    Data --> RunFiles
    Deps --> RunFiles
    RunFiles -->|available at runtime| PythonApp
    RunFiles -->|available at runtime| GoService
    RunFiles -->|available at runtime| JSFrontend
```

## CI/CD Architecture

### Continuous Integration Pipeline

```mermaid
flowchart TB
    Push[Git Push] --> Trigger[Trigger CI]
    
    Trigger --> Checkout[Checkout Code]
    Checkout --> Cache[Restore Cache]
    
    Cache --> Parallel{Parallel Jobs}
    
    Parallel --> Build[Build Job]
    Parallel --> Test[Test Job]
    Parallel --> Lint[Lint Job]
    Parallel --> Format[Format Check Job]
    
    Build --> BuildAll[bazel build //...]
    Test --> TestAll[bazel test //...]
    Lint --> LintAll[aspect lint //...]
    Format --> FormatCheck[format --check]
    
    BuildAll --> BuildStatus{Success?}
    TestAll --> TestStatus{Success?}
    LintAll --> LintStatus{Success?}
    FormatCheck --> FormatStatus{Success?}
    
    BuildStatus -->|Yes| Collect[Collect Results]
    TestStatus -->|Yes| Collect
    LintStatus -->|Yes| Collect
    FormatStatus -->|Yes| Collect
    
    BuildStatus -->|No| Fail[Mark Failed]
    TestStatus -->|No| Fail
    LintStatus -->|No| Fail
    FormatStatus -->|No| Fail
    
    Collect --> Artifacts[Upload Artifacts]
    Artifacts --> SaveCache[Save Cache]
    SaveCache --> Success[Mark Success]
    
    Success --> Deploy{Deploy?}
    Deploy -->|main branch| Production[Deploy to Production]
    Deploy -->|feature branch| Skip[Skip Deploy]
```

### Remote Cache Integration

```mermaid
sequenceDiagram
    participant CI as CI Server
    participant Bazel
    participant Local as Local Cache
    participant Remote as Remote Cache
    participant Build as Build Action
    
    CI->>Bazel: bazel build //...
    Bazel->>Local: Check local cache
    Local-->>Bazel: Not found
    
    Bazel->>Remote: Check remote cache
    Remote-->>Bazel: Cache hit / miss
    
    alt Cache Hit
        Remote->>Bazel: Download artifacts
        Bazel-->>CI: Build complete (cached)
    else Cache Miss
        Bazel->>Build: Execute build action
        Build->>Build: Compile/Process
        Build-->>Bazel: Artifacts ready
        Bazel->>Local: Store locally
        Bazel->>Remote: Upload to remote
        Bazel-->>CI: Build complete (executed)
    end
```

## Summary

These diagrams illustrate:

1. **Overall Architecture**: How components fit together
2. **Generation Flow**: How projects are created from templates
3. **Build System**: How Bazel manages dependencies and builds
4. **Development Workflow**: Daily developer activities
5. **Dependency Management**: How packages are managed per language
6. **Code Quality**: Formatting and linting pipelines
7. **Multi-Language**: How different languages integrate
8. **CI/CD**: Continuous integration and deployment

Use these diagrams as reference when:
- Understanding the system architecture
- Contributing new features
- Debugging issues
- Explaining concepts to team members
- Designing new integrations

---

**Back**: [Documentation Home](./overview.md)
