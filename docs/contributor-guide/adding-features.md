# Adding Features

This guide explains how to add new features to the Aspect Workflows Template.

## Table of Contents

1. [Overview](#overview)
2. [Feature Types](#feature-types)
3. [Adding a Question](#adding-a-question)
4. [Conditional Feature Implementation](#conditional-feature-implementation)
5. [Examples](#examples)
6. [Testing Features](#testing-features)
7. [Best Practices](#best-practices)

## Overview

Features in the template are capabilities that can be optionally enabled. They're controlled through:

1. **Questions** in `scaffold.yaml` that prompt users
2. **Computed variables** derived from user answers
3. **Feature filters** that include/exclude files
4. **Conditional blocks** in template files

## Feature Types

### Boolean Features

Simple on/off features (e.g., linting, OCI images, stamping):

```yaml
questions:
  - name: enable_monitoring
    description: Enable monitoring and observability tools
    type: bool
    default: false
```

### Choice Features

Select from predefined options:

```yaml
questions:
  - name: database
    description: Database to use
    default: "postgres"
    options:
      - postgres
      - mysql
      - mongodb
      - none
```

### Multi-Select Features

Choose multiple options (like `langs`):

```yaml
questions:
  - name: cloud_providers
    description: Cloud providers to support (space-separated)
    default: "aws"
    options:
      - aws
      - gcp
      - azure
```

## Adding a Question

### Step 1: Define Question in scaffold.yaml

```yaml
# scaffold.yaml

questions:
  # ... existing questions ...
  
  - name: enable_graphql
    description: Enable GraphQL code generation
    type: bool
    default: false
```

### Step 2: Add Computed Variable (if needed)

For complex logic, add a computed variable:

```yaml
computed:
  # ... existing computed vars ...
  
  graphql_enabled: |
    {{- and .Scaffold.EnableGraphql .Computed.TypeScript -}}
```

### Step 3: Add Feature Filter

Control which files are included:

```yaml
features:
  # ... existing features ...
  
  - name: graphql
    if: '{{ .Scaffold.EnableGraphql }}'
    include:
      - "{{ .ProjectSnake }}/graphql/**"
      - "{{ .ProjectSnake }}/codegen.yml"
      - "{{ .ProjectSnake }}/schema.graphql"
```

## Conditional Feature Implementation

### In Template Files

Use Go template conditionals:

```python
# {{ .ProjectSnake }}/BUILD

{{- if .Scaffold.EnableGraphql }}
load("@rules_graphql//graphql:defs.bzl", "graphql_codegen")

graphql_codegen(
    name = "generated",
    schema = "schema.graphql",
    queries = glob(["graphql/**/*.graphql"]),
)
{{- end }}
```

### In MODULE.bazel

Add dependencies conditionally:

```python
{{- if .Scaffold.EnableGraphql }}
# GraphQL code generation
bazel_dep(name = "rules_graphql", version = "0.1.0")

graphql = use_extension("@rules_graphql//:extensions.bzl", "graphql")
graphql.toolchain(version = "16.6.0")
use_repo(graphql, "graphql_toolchains")
{{- end }}
```

### In Configuration Files

```yaml
# {{ .ProjectSnake }}/codegen.yml
{{- if .Scaffold.EnableGraphql }}
schema: ./schema.graphql
documents: ./graphql/**/*.graphql
generates:
  ./generated/graphql.ts:
    plugins:
      - typescript
      - typescript-operations
{{- end }}
```

## Examples

### Example 1: Adding Docker Compose Support

#### Step 1: Add question

```yaml
questions:
  - name: enable_docker_compose
    description: Add Docker Compose configuration
    type: bool
    default: false
```

#### Step 2: Create docker-compose.yml template

```yaml
# {{ .ProjectSnake }}/docker-compose.yml
{{- if .Scaffold.EnableDockerCompose }}
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    {{- if .Computed.Python }}
    command: python -m app.main
    {{- else if .Computed.JavaScript }}
    command: node dist/main.js
    {{- else if .Computed.Go }}
    command: /app/main
    {{- end }}
{{- end }}
```

#### Step 3: Add feature filter

```yaml
features:
  - name: docker_compose
    if: '{{ .Scaffold.EnableDockerCompose }}'
    include:
      - "{{ .ProjectSnake }}/docker-compose.yml"
      - "{{ .ProjectSnake }}/.dockerignore"
```

#### Step 4: Update documentation

```markdown
# {{ .ProjectSnake }}/README.bazel.md

{{- if .Scaffold.EnableDockerCompose }}
## Docker Compose

Start the application with Docker Compose:

\`\`\`bash
docker-compose up
\`\`\`
{{- end }}
```

### Example 2: Adding Database Migrations

#### Step 1: Add question

```yaml
questions:
  - name: enable_migrations
    description: Enable database migration tools
    type: bool
    default: false
  
  - name: migration_tool
    description: Migration tool to use
    if: '{{ .Scaffold.EnableMigrations }}'
    default: "liquibase"
    options:
      - liquibase
      - flyway
      - golang-migrate
```

#### Step 2: Add computed variable

```yaml
computed:
  liquibase: |
    {{- and .Scaffold.EnableMigrations (eq .Scaffold.MigrationTool "liquibase") -}}
  
  flyway: |
    {{- and .Scaffold.EnableMigrations (eq .Scaffold.MigrationTool "flyway") -}}
```

#### Step 3: Add dependencies

```python
# {{ .ProjectSnake }}/MODULE.bazel

{{- if .Computed.Liquibase }}
maven.install(
    artifacts = [
        "org.liquibase:liquibase-core:4.24.0",
    ],
)
{{- else if .Computed.Flyway }}
maven.install(
    artifacts = [
        "org.flywaydb:flyway-core:10.0.0",
    ],
)
{{- end }}
```

#### Step 4: Create migration structure

```python
# {{ .ProjectSnake }}/migrations/
{{- if .Scaffold.EnableMigrations }}
{{- if .Computed.Liquibase }}
db/
  changelog/
    db.changelog-master.xml
    changes/
      001-initial-schema.xml
{{- else if .Computed.Flyway }}
db/
  migration/
    V1__initial_schema.sql
{{- end }}
{{- end }}
```

### Example 3: Adding Observability Stack

#### Step 1: Add question

```yaml
questions:
  - name: observability
    description: Observability components (space-separated)
    default: ""
    options:
      - prometheus
      - grafana
      - jaeger
      - loki
```

#### Step 2: Add computed variables

```yaml
computed:
  prometheus: |
    {{- $obs := split .Scaffold.Observability " " -}}
    {{- has "prometheus" $obs -}}
  
  grafana: |
    {{- $obs := split .Scaffold.Observability " " -}}
    {{- has "grafana" $obs -}}
  
  jaeger: |
    {{- $obs := split .Scaffold.Observability " " -}}
    {{- has "jaeger" $obs -}}
```

#### Step 3: Add configurations

```yaml
# {{ .ProjectSnake }}/observability/prometheus.yml
{{- if .Computed.Prometheus }}
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'app'
    static_configs:
      - targets: ['localhost:8080']
{{- end }}
```

#### Step 4: Add feature filters

```yaml
features:
  - name: prometheus
    if: '{{ .Computed.Prometheus }}'
    include:
      - "{{ .ProjectSnake }}/observability/prometheus.yml"
  
  - name: grafana
    if: '{{ .Computed.Grafana }}'
    include:
      - "{{ .ProjectSnake }}/observability/grafana/**"
```

## Testing Features

### Create User Story

```markdown
# user_stories/monitoring.md

---
name: Monitoring Enabled
description: Project with full observability stack
questions:
  langs: py go
  lint: true
  stamp: false
  oci: true
  enable_monitoring: true
  observability: "prometheus grafana jaeger"
---

# Monitoring Configuration

Tests that monitoring and observability features work correctly.
```

### Test Script

```bash
# Test the new feature
./test.sh monitoring

# Verify generated files
cd test-monitoring
ls observability/
```

### Automated Testing

Add to `test.sh`:

```bash
# test.sh additions

test_monitoring() {
    scaffold new . --preset monitoring
    cd test-monitoring
    
    # Verify monitoring files exist
    assert_file_exists "observability/prometheus.yml"
    assert_file_exists "observability/grafana/dashboards/app.json"
    
    # Verify build works
    bazel build //...
    
    cd ..
}
```

## Feature Dependencies

### Features Requiring Other Features

```yaml
computed:
  # OCI images require at least one language
  oci_valid: |
    {{- and .Scaffold.Oci (or .Computed.Python .Computed.Go .Computed.JavaScript) -}}
```

Use in conditionals:

```python
{{- if .Computed.OciValid }}
load("//tools/oci:py3_image.bzl", "py3_image")
{{- end }}
```

### Feature Conflicts

Document incompatibilities:

```yaml
questions:
  - name: serverless
    description: Deploy as serverless functions
    type: bool
    default: false
  
  - name: kubernetes
    description: Deploy to Kubernetes
    type: bool
    default: false
    # Note: Conflicts with serverless deployment
```

Add validation in `post_scaffold`:

```bash
# hooks/post_scaffold

if [ "$SERVERLESS" = "true" ] && [ "$KUBERNETES" = "true" ]; then
    echo "Warning: Both serverless and kubernetes enabled. Choose one."
fi
```

## Advanced Feature Patterns

### Feature Composition

Combine multiple features:

```yaml
computed:
  full_stack: |
    {{- and .Computed.JavaScript .Computed.Python .Scaffold.EnableGraphql .Scaffold.Oci -}}
```

### Progressive Enhancement

Add enhanced features for power users:

```yaml
questions:
  - name: advanced_mode
    description: Enable advanced features
    type: bool
    default: false
  
  - name: cache_backend
    description: Distributed cache backend
    if: '{{ .Scaffold.AdvancedMode }}'
    default: "redis"
    options:
      - redis
      - memcached
      - none
```

### Feature Presets

Group related features:

```yaml
# scaffold.yaml

presets:
  - name: microservice
    questions:
      langs: go
      lint: true
      stamp: true
      oci: true
      enable_monitoring: true
      observability: "prometheus grafana"
      enable_migrations: true
```

## Documentation

### Update README Template

```markdown
# {{ .ProjectSnake }}/README.bazel.md

{{- if .Scaffold.EnableMyFeature }}
## My Feature

This project includes My Feature support.

### Using My Feature

\`\`\`bash
bazel run //path/to:feature
\`\`\`

### Configuration

Edit `config/feature.yml` to configure the feature.
{{- end }}
```

### Add to Main Docs

Create feature-specific documentation:

```markdown
# docs/features/my-feature.md

# My Feature

## Overview

My Feature provides...

## Configuration

...
```

## Best Practices

### DO

- ✅ Make features optional and independent
- ✅ Provide sensible defaults
- ✅ Document feature in README
- ✅ Test feature in isolation
- ✅ Test feature combinations
- ✅ Use feature filters for file inclusion
- ✅ Add computed variables for complex logic
- ✅ Create user story for the feature
- ✅ Consider feature interactions

### DON'T

- ❌ Make features tightly coupled
- ❌ Enable features by default without consideration
- ❌ Forget to add feature filters
- ❌ Leave undocumented configuration
- ❌ Skip testing feature combinations
- ❌ Hard-code feature-specific values
- ❌ Ignore feature conflicts

## Troubleshooting

### Feature Not Included

```bash
# Check feature filter
grep -A 10 "name: myfeature" scaffold.yaml

# Verify conditional
scaffold new . --preset test --set enable_myfeature=true
ls -la test/
```

### Template Syntax Errors

```bash
# Test template rendering
scaffold new . --dry-run --set enable_myfeature=true

# Check for Go template errors
# Look for unclosed {{ }} blocks
```

### Files Not Generated

```bash
# Verify feature filter includes pattern
# scaffold.yaml
features:
  - name: myfeature
    if: '{{ .Scaffold.EnableMyfeature }}'
    include:
      - "{{ .ProjectSnake }}/path/to/files/**"  # Check this pattern
```

## Checklist

Before submitting a PR for new feature:

- [ ] Added question to `scaffold.yaml`
- [ ] Added computed variable (if needed)
- [ ] Added feature filter
- [ ] Created template files with conditionals
- [ ] Added dependencies to MODULE.bazel
- [ ] Updated README template
- [ ] Created user story
- [ ] Tested feature in isolation
- [ ] Tested with various language combinations
- [ ] Tested feature interactions
- [ ] Added documentation
- [ ] All tests pass

## Next Steps

- Learn about [Testing Changes](./testing.md)
- Review [Template System](./template-system.md) internals
- See [Contribution Workflow](./workflow.md)

---

**Back**: [Adding Languages](./adding-languages.md) | **Next**: [Testing Changes](./testing.md)
