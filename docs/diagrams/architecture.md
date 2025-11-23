# Architecture Diagrams

## Direct Generation Mode (Default)

```
┌─────────────────────────────────────────────────────────────┐
│ User runs: scaffold new --preset=py --output-dir=my-app .   │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  scaffold.yaml                │
        │  backstage: false (default)   │
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌─────────────────────────────────┐
        │  Generate project files         │
        │  - .bazelrc                     │
        │  - BUILD                        │
        │  - MODULE.bazel                 │
        │  - pyproject.toml               │
        │  - tools/                       │
        │  - catalog-info.yaml (Component)│
        └───────────────┬─────────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  Run post_scaffold hook       │
        │  - Format files               │
        │  - Run repin                  │
        │  - Install deps               │
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  ✅ Ready-to-use project      │
        │     cd my-app                 │
        │     bazel test //...          │
        └───────────────────────────────┘
```

## Backstage Template Generation Mode

```
┌──────────────────────────────────────────────────────────────────┐
│ User: scaffold new --preset=backstage-py --output-dir=templates  │
└───────────────────────┬──────────────────────────────────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  scaffold.yaml                │
        │  backstage: true (from preset)│
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌────────────────────────────────┐
        │  Generate template structure   │
        │  - template.yaml               │
        │  - catalog-info.yaml (Location)│
        │  - skeleton/                   │
        │  - All project files at root   │
        └───────────────┬────────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  Run post_scaffold hook       │
        │  - Copy files to skeleton/    │
        │  - Format files               │
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌───────────────────────────────┐
        │  ✅ Backstage template ready  │
        │     Publish to GitHub         │
        │     Register in Backstage     │
        └───────────────────────────────┘
```

## Backstage Template Usage Flow

```
┌─────────────────────────────────────────────────────────────┐
│  Platform Team                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 1. Generate template                                 │   │
│  │    scaffold new --preset=backstage-py                │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │ 2. Publish to GitHub                                 │   │
│  │    git push origin main                              │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │ 3. Register in Backstage catalog                     │   │
│  │    Add catalog-info.yaml URL                         │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                      │
                      │ Template available
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  Developer                                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 4. Browse templates in Backstage UI                  │   │
│  │    Click "Create" → Select "Aspect Workflows - Py"   │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │ 5. Fill form                                         │   │
│  │    - Name: payment-service                           │   │
│  │    - Owner: platform-team                            │   │
│  │    - Features: linting ✓, OCI ✓                      │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │ 6. Backstage creates repo                            │   │
│  │    - Fetch skeleton/                                 │   │
│  │    - Replace ${{ values.* }}                         │   │
│  │    - Create GitHub repo                              │   │
│  │    - Register component                              │   │
│  └──────────────────┬───────────────────────────────────┘   │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐   │
│  │ 7. Clone and develop                                 │   │
│  │    git clone .../payment-service                     │   │
│  │    cd payment-service                                │   │
│  │    bazel test //...                                  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## File Structure Comparison

### Direct Mode
```
my-project/
├── catalog-info.yaml          # Kind: Component
├── .bazelrc
├── BUILD
├── MODULE.bazel
├── pyproject.toml
├── requirements/
│   ├── runtime.txt
│   └── all.txt
├── tools/
│   ├── BUILD
│   ├── format/
│   ├── lint/
│   └── repin
└── githooks/
```

### Backstage Mode
```
templates/aspect-python/
├── catalog-info.yaml          # Kind: Location → template.yaml
├── template.yaml              # Scaffolder definition
├── .bazelrc                   # Source of truth
├── BUILD
├── MODULE.bazel
├── pyproject.toml
├── requirements/
├── tools/
├── githooks/
└── skeleton/                  # Template for scaffolding
    ├── catalog-info.yaml      # Kind: Component (with Backstage variables)
    ├── .bazelrc               # Copied from root
    ├── BUILD                  # Copied from root
    ├── MODULE.bazel           # Copied from root
    ├── pyproject.toml         # Copied from root
    ├── requirements/          # Copied from root
    ├── tools/                 # Copied from root
    └── githooks/              # Copied from root
```

## Decision Flow

```
                    Start
                      │
                      ▼
        ┌─────────────────────────┐
        │ What do you need?       │
        └─────────┬───────┬───────┘
                  │       │
        ┌─────────┘       └─────────┐
        │                           │
        ▼                           ▼
┌───────────────┐         ┌───────────────┐
│ Quick project │         │ Reusable      │
│ for immediate │         │ template for  │
│ development   │         │ Backstage     │
└───────┬───────┘         └───────┬───────┘
        │                         │
        ▼                         ▼
┌───────────────┐         ┌───────────────┐
│ Direct Mode   │         │ Backstage Mode│
│ --preset=py   │         │--preset=      │
│               │         │backstage-py   │
└───────┬───────┘         └───────┬───────┘
        │                         │
        ▼                         ▼
┌───────────────┐         ┌───────────────┐
│ Ready project │         │ Template in   │
│ cd my-project │         │ Backstage     │
│ bazel test    │         │ catalog       │
└───────────────┘         └───────────────┘
```

## Copy Strategy

```
Backstage Template Structure
────────────────────────────

┌─────────────────────────────────────────────┐
│ Root Level (Source of Truth)                │
│                                             │
│  .bazelrc                                   │
│  BUILD                                      │
│  MODULE.bazel                               │
│  pyproject.toml                             │
│  tools/                                     │
│                                             │
└─────────────────┬───────────────────────────┘
                  │
                  │ Copied to (by post_scaffold hook)
                  │
                  ▼
┌─────────────────────────────────────────────┐
│ skeleton/ (Template Files)                  │
│                                             │
│  .bazelrc                                   │
│  BUILD                                      │
│  MODULE.bazel                               │
│  pyproject.toml                             │
│  tools/                                     │
│                                             │
│  catalog-info.yaml  (Backstage variables)   │
│    Uses: ${{ values.name }}                 │
│          ${{ values.owner }}                │
│                                             │
└─────────────────────────────────────────────┘

How it works:
✓ Post-generation hook copies files to skeleton/ (not symlinks - Backstage doesn't follow them)
✓ Some files excluded (template.yaml, root catalog-info.yaml, README.bazel.md)
✓ Backstage-specific files use ${{ values.* }} syntax
```
