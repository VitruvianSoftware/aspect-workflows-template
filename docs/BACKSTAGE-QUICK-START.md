# Backstage Mode Quick Reference

## Generate Project Directly (Default)

```bash
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=py --output-dir=my-project .
```

Result: Ready-to-use Bazel project in `my-project/`

## Generate Backstage Template

```bash
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-py --output-dir=templates/aspect-python .
```

Result: Backstage template in `templates/aspect-python/`

> **Important**: Use `SCAFFOLD_SETTINGS_RUN_HOOKS=always` to ensure post-generation hooks run, creating the skeleton/ symlinks for Backstage templates.

## Comparison

| Feature | Direct Mode | Backstage Mode |
|---------|------------|----------------|
| **Command** | `--preset=py` | `--preset=backstage-py` |
| **Output** | Runnable project | Template definition |
| **catalog-info.yaml** | Kind: Component | Kind: Location |
| **Extra Files** | None | template.yaml, skeleton/ |
| **Usage** | Immediate coding | Register in Backstage |

## Backstage Template Structure

```
aspect-python/
├── catalog-info.yaml      # Points to template.yaml
├── template.yaml          # Scaffolder definition
└── skeleton/              # Code template
    ├── catalog-info.yaml  # For generated projects
    └── ... (symlinks)
```

## Presets

### Direct Generation
- `py`, `js`, `go`, `java`, `kotlin`, `cpp`, `rust`, `shell`
- `kitchen-sink` (all languages)
- `minimal` (bare bones)

### Backstage Templates
- `backstage-py`, `backstage-js`, `backstage-go`
- `backstage-java`, `backstage-kotlin`, `backstage-cpp`, `backstage-rust`, `backstage-shell`
- `backstage-kitchen-sink` (all languages)
- `backstage-minimal` (bare bones)

## Publishing Backstage Template

```bash
cd templates/aspect-python
git init
git add .
git commit -m "Add Aspect Python template"
git remote add origin https://github.com/YOUR-ORG/aspect-python-template.git
git push -u origin main
```

## Registering in Backstage

### Option 1: UI
1. Backstage → Create → Register Component
2. Enter: `https://github.com/YOUR-ORG/aspect-python-template/blob/main/catalog-info.yaml`

### Option 2: app-config.yaml
```yaml
catalog:
  locations:
    - type: url
      target: https://github.com/YOUR-ORG/aspect-python-template/blob/main/catalog-info.yaml
```

## Common Commands

```bash
# List available presets
scaffold list

# Generate with custom answers
SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --output-dir=output .
# Answer "yes" to "Generate as a Backstage template?"

# Generate all Backstage templates
for lang in py js go; do
  SCAFFOLD_SETTINGS_RUN_HOOKS=always scaffold new --preset=backstage-$lang --output-dir=templates/aspect-$lang .
done

# Test template structure
cd templates/aspect-python
ls -la skeleton/  # Check symlinks
cat template.yaml # Verify parameters
```

## Variables

### In Direct Mode
```go
// Uses Scaffold Go templates
{{ .ProjectSnake }}
{{ .Scaffold.langs }}
```

### In Backstage Mode (skeleton/ files)
```yaml
# Uses Backstage/Nunjucks templates
${{ values.name }}
${{ values.owner }}
{% if values.enableLinting %}...{% endif %}
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Skeleton symlinks broken | Re-generate with `SCAFFOLD_SETTINGS_RUN_HOOKS=always` |
| Template not in Backstage | Check catalog logs, verify Location kind |
| Variables not replaced | Use `${{ values.* }}` in skeleton files |
| Missing features | Check parameter names in template.yaml |

## Documentation

- Full guide: `docs/admin-guide/backstage-integration.md`
- Scaffold docs: https://hay-kot.github.io/scaffold/
- Backstage templates: https://backstage.io/docs/features/software-templates/
