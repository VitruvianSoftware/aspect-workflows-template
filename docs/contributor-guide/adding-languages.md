# Adding Language Support

This guide explains how to add support for a new programming language to the Aspect Workflows Template.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Step-by-Step Process](#step-by-step-process)
4. [Language Components](#language-components)
5. [Testing Language Support](#testing-language-support)
6. [Examples](#examples)
7. [Best Practices](#best-practices)

## Overview

Adding a new language involves several coordinated changes:

1. Add language to selection question in `scaffold.yaml`
2. Add computed variable for the language
3. Add language-specific dependencies to `MODULE.bazel`
4. Add Gazelle configuration
5. Add formatting and linting rules
6. Add package manager configuration
7. Create example user story
8. Test thoroughly

## Prerequisites

Before adding language support, ensure:

- The language has mature Bazel rules (rules_X)
- Bazel rules support bzlmod (MODULE.bazel)
- Formatters and linters are available
- Package manager integration exists (if needed)
- You understand the language's ecosystem

## Step-by-Step Process

### Step 1: Update scaffold.yaml

Add the language to the `langs` question:

```yaml
# scaffold.yaml
questions:
  - name: langs
    description: Programming languages to use (space-separated)
    default: "py"
    options:
      - go
      - py
      - js
      - rust
      - java
      - kotlin
      - cpp
      - shell
      - ruby  # NEW LANGUAGE
```

Add a computed variable:

```yaml
computed:
  # ... existing computed vars ...
  
  ruby: |
    {{- $langs := split .Scaffold.Langs " " -}}
    {{- has "ruby" $langs -}}
```

### Step 2: Add Bazel Dependencies

Update `MODULE.bazel` template to include the language's Bazel rules:

```python
# {{ .ProjectSnake }}/MODULE.bazel

{{- if .Computed.Ruby }}
# Ruby support
bazel_dep(name = "rules_ruby", version = "0.7.0")

ruby = use_extension("@rules_ruby//:extensions.bzl", "ruby")
ruby.toolchain(
    name = "ruby",
    version = "3.2.2",
)
use_repo(ruby, "ruby_toolchains")
register_toolchains("@ruby_toolchains//:all")
{{- end }}
```

### Step 3: Configure Package Manager

If the language has a package manager, add configuration:

```ruby
# {{ .ProjectSnake }}/Gemfile (for Ruby example)
{{- if .Computed.Ruby }}
source "https://rubygems.org"

gem "rake"
gem "minitest"
{{- end }}
```

Add lockfile to gitignore:

```ruby
# {{ .ProjectSnake }}/.gitignore
{{- if .Computed.Ruby }}
Gemfile.lock
{{- end }}
```

### Step 4: Add Gazelle Support

Configure Gazelle to generate BUILD files for the language:

```python
# {{ .ProjectSnake }}/BUILD

{{- if .Computed.Ruby }}
# Gazelle for Ruby
load("@rules_ruby//ruby:defs.bzl", "ruby_gazelle")

ruby_gazelle(
    name = "gazelle_ruby",
)
{{- end }}
```

Update gazelle configuration:

```python
# {{ .ProjectSnake }}/BUILD

gazelle(
    name = "gazelle",
    {{- if .Computed.Ruby }}
    command = "fix",
    gazelle = ":gazelle_ruby",
    {{- end }}
)
```

### Step 5: Add Formatting Support

Add formatter to the format tool configuration:

```python
# {{ .ProjectSnake }}/tools/format/BUILD

{{- if .Computed.Ruby }}
load("@rules_ruby//ruby:defs.bzl", "ruby_binary")

ruby_binary(
    name = "rubocop",
    main = "rubocop",
    deps = ["@gem//rubocop"],
)
{{- end }}
```

Update format multitool:

```python
# {{ .ProjectSnake }}/tools/BUILD

multitool(
    name = "format",
    commands = [
        # ... existing formatters ...
        {{- if .Computed.Ruby }}
        tool(
            name = "rubocop",
            tool = "//tools/format:rubocop",
            args = ["--autocorrect"],
        ),
        {{- end }}
    ],
)
```

### Step 6: Add Linting Support

Configure linter in rules_lint:

```python
# {{ .ProjectSnake }}/tools/lint/BUILD

{{- if .Computed.Ruby }}
load("@rules_ruby//ruby:defs.bzl", "ruby_lint_aspect")

ruby_lint_aspect(
    name = "ruby_lint",
    config = "//:.rubocop.yml",
)
{{- end }}
```

Add to linters list:

```python
# {{ .ProjectSnake }}/tools/lint/linters.bzl

LINTERS = [
    # ... existing linters ...
    {{- if .Computed.Ruby }}
    "//tools/lint:ruby_lint",
    {{- end }}
]
```

### Step 7: Add Example Code

Create sample directory structure:

```ruby
# {{ .ProjectSnake }}/ruby/hello/hello.rb
{{- if .Computed.Ruby }}
# frozen_string_literal: true

module Hello
  def self.greet(name)
    "Hello, #{name}!"
  end
end

puts Hello.greet("World") if __FILE__ == $PROGRAM_NAME
{{- end }}
```

Add BUILD file:

```python
# {{ .ProjectSnake }}/ruby/hello/BUILD
{{- if .Computed.Ruby }}
load("@rules_ruby//ruby:defs.bzl", "ruby_library", "ruby_binary", "ruby_test")

ruby_library(
    name = "hello_lib",
    srcs = ["hello.rb"],
    visibility = ["//visibility:public"],
)

ruby_binary(
    name = "hello",
    main = "hello.rb",
    deps = [":hello_lib"],
)

ruby_test(
    name = "hello_test",
    srcs = ["hello_test.rb"],
    deps = [":hello_lib"],
)
{{- end }}
```

### Step 8: Update Feature Filters

Add feature filters to include/exclude language-specific files:

```yaml
# scaffold.yaml

features:
  # ... existing features ...
  
  - name: ruby
    if: '{{ .Computed.Ruby }}'
    include:
      - "{{ .ProjectSnake }}/ruby/**"
      - "{{ .ProjectSnake }}/Gemfile"
      - "{{ .ProjectSnake }}/.rubocop.yml"
```

### Step 9: Update Documentation

Add language to README:

```markdown
# {{ .ProjectSnake }}/README.bazel.md

{{- if .Computed.Ruby }}
## Ruby Development

### Running Ruby code

\`\`\`bash
bazel run //ruby/hello
\`\`\`

### Testing Ruby code

\`\`\`bash
bazel test //ruby/...
\`\`\`

### Ruby dependencies

Dependencies are managed in `Gemfile`. After updating:

\`\`\`bash
bundle install
bazel sync --only=gems
\`\`\`
{{- end }}
```

### Step 10: Create User Story

Create a test scenario in `user_stories/`:

```markdown
# user_stories/ruby.md

---
name: Ruby Example
description: Ruby project with basic structure
questions:
  langs: ruby
  lint: true
  stamp: false
  oci: false
  codegen: false
---

# Ruby Example

This preset creates a Ruby project with:
- Ruby toolchain via rules_ruby
- Gem dependency management
- RuboCop for linting and formatting
- Minitest for testing
- Example library and binary
```

## Language Components

### Required Components

Every language integration needs:

1. **Bazel Rules**: `bazel_dep(name = "rules_X")`
2. **Toolchain**: Language compiler/interpreter
3. **BUILD Generation**: Gazelle support
4. **Formatter**: Code formatting tool
5. **Linter**: Static analysis tool
6. **Package Manager**: Dependency management (if applicable)
7. **Example Code**: Sample library, binary, and test
8. **Documentation**: Language-specific README section

### Optional Components

Consider adding:

- **IDE Support**: LSP configuration
- **Container Images**: OCI image rules for the language
- **Code Generation**: Protobuf, GraphQL, etc.
- **Benchmarks**: Performance testing
- **Coverage**: Code coverage tools

## Testing Language Support

### Test New Language

```bash
# Test the new language in isolation
./test.sh ruby

# Test with combinations
scaffold new . --preset ruby --set langs="ruby py"
cd test-project
bazel build //...
bazel test //...
```

### Verify Components

Check that all components work:

```bash
# Formatting
format ruby/**/*.rb

# Linting
aspect lint //ruby/...

# Building
bazel build //ruby/...

# Testing
bazel test //ruby/...

# Gazelle
bazel run gazelle
```

### Test User Stories

```bash
# Run all tests including new language
./test.sh

# Or test specific preset
./test.sh ruby
```

## Examples

### Example 1: Adding Swift Support

```yaml
# scaffold.yaml additions
questions:
  - name: langs
    options:
      - swift  # Add to list

computed:
  swift: |
    {{- $langs := split .Scaffold.Langs " " -}}
    {{- has "swift" $langs -}}
```

```python
# MODULE.bazel additions
{{- if .Computed.Swift }}
bazel_dep(name = "rules_swift", version = "1.13.0")
bazel_dep(name = "rules_apple", version = "3.2.0")

swift = use_extension("@rules_swift//swift:extensions.bzl", "swift")
swift.toolchain(version = "5.9")
use_repo(swift, "swift_toolchains")
register_toolchains("@swift_toolchains//:all")
{{- end }}
```

### Example 2: Adding Scala Support

```yaml
# scaffold.yaml
computed:
  scala: |
    {{- $langs := split .Scaffold.Langs " " -}}
    {{- has "scala" $langs -}}
```

```python
# MODULE.bazel
{{- if .Computed.Scala }}
bazel_dep(name = "rules_scala", version = "6.4.0")

scala = use_extension("@rules_scala//scala:extensions.bzl", "scala")
scala.toolchain(
    scala_version = "2.13.12",
)
use_repo(scala, "scala_toolchains")
{{- end }}
```

## Best Practices

### DO

- ✅ Use stable versions of Bazel rules
- ✅ Follow the language's official style guide
- ✅ Provide working example code
- ✅ Include both library and binary examples
- ✅ Add comprehensive tests
- ✅ Document language-specific workflows
- ✅ Test with multiple language combinations
- ✅ Use hermetic toolchains (avoid system dependencies)
- ✅ Add to CI testing matrix

### DON'T

- ❌ Use beta/experimental Bazel rules in main template
- ❌ Add language without formatter/linter
- ❌ Forget to update documentation
- ❌ Skip user story creation
- ❌ Leave TODO comments in generated code
- ❌ Hard-code versions (use variables)
- ❌ Assume language is installed on system

## Common Patterns

### Conditional File Inclusion

```yaml
features:
  - name: language_name
    if: '{{ .Computed.LanguageName }}'
    include:
      - "{{ .ProjectSnake }}/language_dir/**"
      - "{{ .ProjectSnake }}/.language_config"
    exclude:
      - "{{ .ProjectSnake }}/language_dir/**/*.generated.*"
```

### Multi-Version Support

```python
# Allow users to select language version
{{- if .Computed.Language }}
language_toolchain(
    version = "{{ .Scaffold.LanguageVersion | default "1.0.0" }}",
)
{{- end }}
```

### Cross-Language Dependencies

```python
# Language B depends on Language A
{{- if and .Computed.LanguageA .Computed.LanguageB }}
language_b_library(
    name = "cross_lang",
    deps = [
        "//language_a/lib",  # Cross-language dependency
    ],
)
{{- end }}
```

## Troubleshooting

### Rules Not Found

```bash
# Verify rules are in Bazel registry
bazel mod deps | grep rules_language

# Check version availability
# Visit https://registry.bazel.build/modules/rules_language
```

### Gazelle Not Generating

```bash
# Check Gazelle configuration
bazel query //... --output=build | grep gazelle

# Run with verbose output
bazel run gazelle -- -v
```

### Formatter Not Working

```bash
# Test formatter directly
bazel run //tools/format:language_formatter -- --help

# Check multitool configuration
bazel query //tools:format --output=build
```

### Missing Toolchain

```bash
# List registered toolchains
bazel query 'kind("toolchain", @language_toolchains//...)'

# Verify toolchain resolution
bazel cquery --output=build //language/hello
```

## Checklist

Before submitting a PR for new language support:

- [ ] Added language to `scaffold.yaml` questions
- [ ] Added computed variable
- [ ] Added Bazel rules to `MODULE.bazel`
- [ ] Configured package manager (if applicable)
- [ ] Added Gazelle support
- [ ] Added formatter
- [ ] Added linter
- [ ] Created example code (library, binary, test)
- [ ] Added feature filters
- [ ] Updated README template
- [ ] Created user story in `user_stories/`
- [ ] Tested in isolation (`./test.sh language`)
- [ ] Tested with language combinations
- [ ] Tested formatting works
- [ ] Tested linting works
- [ ] Tested BUILD generation works
- [ ] Added documentation
- [ ] All tests pass

## Next Steps

- Read [Adding Features](./adding-features.md) to add capabilities
- See [Testing Guide](./testing.md) for comprehensive testing
- Review [Architecture](./architecture.md) for system overview

---

**Back**: [Contributor Guide](./README.md) | **Next**: [Adding Features](./adding-features.md)
