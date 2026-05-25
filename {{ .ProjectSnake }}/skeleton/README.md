# Skeleton Directory for Backstage Templates

{{- if .Scaffold.backstage }}
This directory contains the actual code template that will be scaffolded when users create a new project from Backstage.

## Structure

When Backstage mode is enabled, this directory mirrors the parent directory structure but uses Backstage template variables instead of Scaffold variables.

## File Synchronization

At generation time the `post_scaffold` hook **copies** the rendered project files from the parent directory into this `skeleton/` directory. Real copies are used rather than symlinks, because Backstage's scaffolder does not follow symlinks — you don't manage these copies by hand:

```bash
# The post_scaffold hook copies parent files into skeleton/ automatically.
# Excluded from the copy: template.yaml, the root catalog-info.yaml,
# skeleton/ itself, README.md, and README.bazel.md
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
