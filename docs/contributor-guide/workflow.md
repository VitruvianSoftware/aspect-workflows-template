# Contribution Workflow

This guide explains the process for contributing changes to the Aspect Workflows Template.

## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Pull Request Process](#pull-request-process)
5. [Code Review](#code-review)
6. [Merging](#merging)
7. [Release Process](#release-process)

## Overview

We follow a standard GitHub workflow with some specific requirements for template projects:

1. Fork and clone
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit Pull Request
6. Address review feedback
7. Merge

## Getting Started

### Fork the Repository

1. Visit https://github.com/aspect-build/aspect-workflows-template
2. Click "Fork" button
3. Clone your fork:

```bash
git clone https://github.com/YOUR-USERNAME/aspect-workflows-template
cd aspect-workflows-template
```

### Set Up Remotes

```bash
# Add upstream remote
git remote add upstream https://github.com/aspect-build/aspect-workflows-template

# Verify remotes
git remote -v
# origin    https://github.com/YOUR-USERNAME/aspect-workflows-template (fetch)
# origin    https://github.com/YOUR-USERNAME/aspect-workflows-template (push)
# upstream  https://github.com/aspect-build/aspect-workflows-template (fetch)
# upstream  https://github.com/aspect-build/aspect-workflows-template (push)
```

### Install Dependencies

```bash
# Install scaffold
brew tap hay-kot/scaffold-tap
brew install scaffold

# Verify installation
scaffold --version
```

## Development Workflow

### 1. Sync with Upstream

Before starting work, sync your fork:

```bash
# Fetch upstream changes
git fetch upstream

# Checkout main branch
git checkout main

# Merge upstream changes
git merge upstream/main

# Push to your fork
git push origin main
```

### 2. Create Feature Branch

```bash
# Create branch from main
git checkout -b feature/my-improvement

# Or for bug fixes
git checkout -b fix/issue-123

# Or for documentation
git checkout -b docs/update-readme
```

Branch naming conventions:

- `feature/` - New features or enhancements
- `fix/` - Bug fixes
- `docs/` - Documentation only
- `refactor/` - Code refactoring
- `test/` - Test additions or fixes

### 3. Make Changes

Edit template files in `{{ .ProjectSnake }}/`:

```bash
# Edit template files
vim '{{ .ProjectSnake }}/MODULE.bazel'

# Edit scaffold configuration
vim scaffold.yaml

# Edit post-processing hook
vim hooks/post_scaffold
```

### 4. Test Changes

**Run full test suite:**

```bash
./test.sh
```

**Test specific presets:**

```bash
./test.sh kitchen-sink
./test.sh minimal
./test.sh py
./test.sh go
./test.sh js
```

**Manual testing:**

```bash
# Generate test project
scaffold new . --preset kitchen-sink

# Enter and test
cd test-kitchen-sink
bazel test //...
aspect lint //...
format --check

# Clean up
cd ..
rm -rf test-kitchen-sink
```

### 5. Commit Changes

Write clear, descriptive commit messages:

```bash
# Stage changes
git add '{{ .ProjectSnake }}/MODULE.bazel'
git add scaffold.yaml

# Commit with descriptive message
git commit -m "feat: add support for Ruby language

- Add Ruby to language options
- Configure rules_ruby in MODULE.bazel
- Add RuboCop formatter and linter
- Create Ruby example code
- Add user story for Ruby

Fixes #123"
```

**Commit message format:**

```text
<type>: <short summary>

<detailed description>

<footer>
```

Types:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Formatting, no code change
- `refactor:` - Code restructuring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

### 6. Push Changes

```bash
# Push branch to your fork
git push origin feature/my-improvement
```

## Pull Request Process

### 1. Create Pull Request

1. Go to your fork on GitHub
2. Click "Pull Request" button
3. Select base: `aspect-build/aspect-workflows-template:main`
4. Select compare: `YOUR-USERNAME/aspect-workflows-template:feature/my-improvement`
5. Click "Create Pull Request"

### 2. PR Title and Description

**Title format:**

```text
feat: add Ruby language support
```

**Description template:**

```markdown
## Description

Brief description of what this PR does.

## Motivation

Why is this change needed? What problem does it solve?

## Changes

- List of changes
- Another change
- One more change

## Testing

How was this tested?

- [ ] Ran `./test.sh`
- [ ] Tested with preset: ruby
- [ ] Tested with combinations: ruby+py, ruby+go
- [ ] Manual testing in generated project

## Checklist

- [ ] Tests pass
- [ ] Documentation updated
- [ ] User story created (if applicable)
- [ ] No breaking changes
- [ ] Follows template conventions

## Related Issues

Closes #123
Related to #456
```

### 3. PR Checklist

Before submitting, ensure:

- [ ] All tests pass (`./test.sh`)
- [ ] Code follows template conventions
- [ ] Documentation is updated
- [ ] User story created for new features
- [ ] Commit messages are clear
- [ ] No unnecessary files included
- [ ] Changes are backwards compatible (or documented)

## Code Review

### What Reviewers Look For

1. **Correctness**: Does it work as intended?
2. **Testing**: Are there adequate tests?
3. **Documentation**: Is it well documented?
4. **Style**: Does it follow conventions?
5. **Compatibility**: Works with all language combinations?
6. **Performance**: No significant performance regressions?

### Addressing Feedback

```bash
# Make requested changes
vim '{{ .ProjectSnake }}/MODULE.bazel'

# Commit changes
git add '{{ .ProjectSnake }}/MODULE.bazel'
git commit -m "fix: address review feedback

- Update Ruby version to 3.2.2
- Add missing test case
- Fix typo in documentation"

# Push changes
git push origin feature/my-improvement
```

### Updating Your Branch

If `main` has changed:

```bash
# Fetch upstream
git fetch upstream

# Rebase on upstream/main
git rebase upstream/main

# Force push (only on your branch!)
git push origin feature/my-improvement --force
```

Or merge instead of rebase:

```bash
git fetch upstream
git merge upstream/main
git push origin feature/my-improvement
```

## Merging

### Merge Requirements

Before merging, ensure:

- [ ] All CI checks pass
- [ ] At least one approval from maintainer
- [ ] All conversations resolved
- [ ] Branch is up to date with main
- [ ] No merge conflicts

### Merge Methods

**Squash and Merge (preferred)**

- Combines all commits into one
- Keeps main history clean
- Use for feature branches

**Rebase and Merge**

- Preserves individual commits
- Use for well-structured commit history

**Merge Commit**

- Creates merge commit
- Use for large features with complex history

## Release Process

### Version Strategy

The template uses semantic versioning:

- **Major** (v2.0.0): Breaking changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Patch** (v1.0.1): Bug fixes

### Creating a Release

Maintainers follow this process:

1. **Update version references** (if any)
2. **Update CHANGELOG.md**:

    ```markdown
    ## [1.1.0] - 2024-01-15

    ### Added
    - Ruby language support (#123)
    - Docker Compose configuration option (#124)

    ### Fixed
    - Python dependency resolution (#125)

    ### Changed
    - Update rules_go to 0.44.0 (#126)
    ```

3. **Create git tag**:

    ```bash
    git tag -a v1.1.0 -m "Release v1.1.0"
    git push upstream v1.1.0
    ```

4. **Create GitHub release**:

   - Go to Releases page
   - Click "Create new release"
   - Select tag v1.1.0
   - Add release notes from CHANGELOG
   - Publish release

## Best Practices

### DO

- ✅ Keep changes focused and atomic
- ✅ Write descriptive commit messages
- ✅ Test thoroughly before submitting
- ✅ Update documentation
- ✅ Respond to review feedback promptly
- ✅ Keep PRs reasonably sized
- ✅ Reference issues in commits/PRs
- ✅ Follow existing code style

### DON'T

- ❌ Submit untested changes
- ❌ Include unrelated changes in PR
- ❌ Force push to main branch
- ❌ Ignore review feedback
- ❌ Make breaking changes without discussion
- ❌ Commit generated files
- ❌ Use generic commit messages

## Common Workflows

### Fixing a Bug

```bash
# Create issue first (or reference existing)
# Create branch
git checkout -b fix/python-deps-issue

# Make fix
vim '{{ .ProjectSnake }}/requirements/all.in'

# Test fix
./test.sh py

# Commit
git commit -m "fix: resolve Python dependency conflict

The transitive dependency issue was caused by...

Fixes #123"

# Push and create PR
git push origin fix/python-deps-issue
```

### Adding Documentation

```bash
# Create branch
git checkout -b docs/add-kotlin-guide

# Add documentation
vim docs/user-guide/languages/kotlin.md

# Commit
git commit -m "docs: add Kotlin language guide

- Add getting started section
- Add dependency management
- Add testing examples"

# Push and create PR
git push origin docs/add-kotlin-guide
```

### Large Feature Development

```bash
# Create feature branch
git checkout -b feature/kubernetes-support

# Work in smaller commits
git commit -m "feat: add kubernetes question to scaffold.yaml"
git commit -m "feat: add kubernetes deployment templates"
git commit -m "feat: add kubernetes documentation"
git commit -m "test: add kubernetes user story"

# Push and create PR
git push origin feature/kubernetes-support
```

## Getting Help

### During Development

- Check existing documentation
- Look at similar features in codebase
- Search closed issues and PRs
- Ask in GitHub Discussions

### During Review

- Respond to feedback
- Ask clarifying questions
- Propose alternative approaches
- Request additional review if needed

### Contact

- **Issues**: Use GitHub Issues for bugs/features
- **Discussions**: Use GitHub Discussions for questions
- **Slack**: Join #aspect-build on Bazel Slack
- **Email**: Contact maintainers for private matters

## Code of Conduct

We follow the [Bazel Community Code of Conduct](https://bazel.build/community/code-of-conduct).

Key points:

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what's best for the community
- Show empathy towards others

## Recognition

Contributors are recognized through:

- Attribution in release notes
- GitHub contributors page
- Special recognition for significant contributions

Thank you for contributing! 🎉

---

**Back**: [Testing Changes](./testing.md) | **Up**: [Contributor Guide](./README.md)
