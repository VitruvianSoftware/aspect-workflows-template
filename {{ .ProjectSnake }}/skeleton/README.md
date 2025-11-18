# Skeleton Directory for Backstage Templates

{{- if .Scaffold.backstage }}
This directory contains the actual code template that will be scaffolded when users create a new project from Backstage.

## Structure

When Backstage mode is enabled, this directory mirrors the parent directory structure but uses Backstage template variables instead of Scaffold variables.

## File Synchronization

Most files in this skeleton/ directory should be symlinks to the parent directory to avoid duplication:

```bash
# Create symlinks for files that don't need Backstage variables
ln -s ../.bazelrc .bazelrc
ln -s ../.bazelversion .bazelversion
ln -s ../BUILD BUILD
ln -s ../MODULE.bazel MODULE.bazel
# ... etc
```

Files that need Backstage-specific templating (like README.bazel.md with project name) should be separate copies with `${{ "{{" }} values.VAR {{ "}}" }}` variables.

## Backstage Variables

Use these variables in skeleton files:

- `${{ "{{" }} values.name {{ "}}" }}` - Project name
- `${{ "{{" }} values.description {{ "}}" }}` - Project description  
- `${{ "{{" }} values.owner {{ "}}" }}` - Project owner
- `${{ "{{" }} values.destination.owner {{ "}}" }}` - GitHub organization
- `${{ "{{" }} values.destination.repo {{ "}}" }}` - Repository name
- `${{ "{{" }} values.languages {{ "}}" }}` - Selected languages (array)
- `${{ "{{" }} values.enableLinting {{ "}}" }}` - Boolean for linting
- `${{ "{{" }} values.enableOCI {{ "}}" }}` - Boolean for OCI
- `${{ "{{" }} values.enableStamping {{ "}}" }}` - Boolean for version stamping
- `${{ "{{" }} values.enableCodegen {{ "}}" }}` - Boolean for code generation

## Conditional Files

Use Nunjucks syntax for conditional content:

```
{% if values.enableLinting %}
# Linting configuration
{% endif %}
```

See ../template.yaml for the full parameter definitions.
{{- else }}
# Skeleton Directory

This directory is only created when generating Backstage templates.

To create a Backstage template, regenerate with:
```bash
scaffold new --output-dir=/tmp/my-backstage-template --preset=py .
# and answer "yes" to "Generate as a Backstage template?"
```
{{- end }}
