# User Guide

Welcome to the Aspect Workflows Template User Guide! This guide will help you get started with creating and working with projects generated from this template.

## Table of Contents

1. [Getting Started](./getting-started.md) - Install and create your first project
2. [Language-Specific Guides](./languages/README.md) - Working with different languages
3. [Development Workflow](./development-workflow.md) - Day-to-day development tasks
4. [Building and Testing](./building-testing.md) - Compile, test, and verify your code
5. [Formatting and Linting](./formatting-linting.md) - Code quality and style
6. [Dependency Management](./dependency-management.md) - Managing external dependencies
7. [Troubleshooting](./troubleshooting.md) - Common issues and solutions

## Quick Start

```bash
# Install scaffold
brew tap hay-kot/scaffold-tap
brew install scaffold

# Create a new project
scaffold new github.com/BlueCentre/aspect-workflows-template

# Enter your project
cd <your-project-name>

# Set up development environment
direnv allow
bazel run //tools:bazel_env

# Start developing!
aspect test //...
```

## What You'll Learn

### As a User

- How to generate a new project from this template
- How to set up your local development environment
- How to work with different programming languages
- How to run common development tasks (build, test, format, lint)
- How to manage dependencies for your project
- How to create container images for deployment
- How to version and stamp your releases

### Who This Guide Is For

This guide is designed for:

- **New developers** starting a project from the template
- **Team members** joining an existing project
- **Developers** evaluating Bazel and Aspect Workflows
- **Anyone** wanting to understand how to use the generated projects

## Prerequisites

Before using this template, you should have:

- Basic familiarity with command-line tools
- Understanding of the programming language(s) you plan to use
- Git installed on your system
- One of the following:
  - macOS or Linux development environment
  - Windows with WSL2

## What's Included in Generated Projects

Projects generated from this template include:

- ✅ **Bazel Build System**: Fast, reproducible builds
- ✅ **Aspect CLI**: Enhanced Bazel developer experience
- ✅ **Development Environment Setup**: Automated tool installation via direnv
- ✅ **Multi-Language Support**: Choose from 8 programming languages
- ✅ **Code Quality Tools**: Formatting and linting pre-configured
- ✅ **Testing Framework**: Language-specific test runners integrated
- ✅ **Dependency Management**: Hermetic dependency resolution
- ✅ **Git Hooks**: Pre-commit hooks for code quality
- ✅ **Container Support** (optional): OCI image building
- ✅ **Version Stamping** (optional): Automated versioning

## Navigation

- **Next**: [Getting Started](./getting-started.md)
- **Up**: [Documentation Home](../overview.md)
