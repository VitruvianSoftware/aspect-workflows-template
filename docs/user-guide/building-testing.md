# Building and Testing

This guide covers building and testing code in projects generated from the Aspect Workflows Template.

## Table of Contents

1. [Building Code](#building-code)
2. [Running Tests](#running-tests)
3. [Test Organization](#test-organization)
4. [Advanced Build Techniques](#advanced-build-techniques)
5. [Continuous Testing](#continuous-testing)
6. [Build Optimization](#build-optimization)

## Building Code

### Basic Build Commands

```bash
# Build everything in the workspace
bazel build //...

# Build a specific target
bazel build //src/app:main

# Build a specific package
bazel build //src/app:all

# Build with verbose output
bazel build -s //src/app:main
```

### Build Configurations

Projects include several pre-configured build modes:

```bash
# Debug build (default) - fast compilation, debug symbols
bazel build //src/app:main

# Optimized build - slower compilation, optimized for performance
bazel build -c opt //src/app:main

# Release build - optimized + stamping
bazel build --config=release //src/app:main
```

### Understanding Build Output

```bash
# Build and see where output goes
bazel build //src/app:main
ls -l bazel-bin/src/app/main

# Show build output location
bazel cquery --output=files //src/app:main
```

### Platform-Specific Builds

```bash
# Build for specific platform
bazel build --platforms=//tools/platforms:linux_amd64 //src/app:main

# Build for multiple platforms
bazel build \
  --platforms=//tools/platforms:linux_amd64 \
  --platforms=//tools/platforms:linux_arm64 \
  //src/app:main

# Cross-compile for different OS
bazel build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //...
```

## Running Tests

### Basic Test Commands

```bash
# Run all tests
bazel test //...

# Use Aspect CLI for better output
aspect test //...

# Run tests in a specific package
bazel test //src/app:all

# Run a specific test
bazel test //src/app:app_test
```

### Test Output Options

```bash
# Show all test output
bazel test --test_output=all //src/app:app_test

# Show only errors
bazel test --test_output=errors //...

# Stream output in real-time
bazel test --test_output=streamed //src/app:app_test

# Show summary only (default)
bazel test //...
```

### Test Filtering

```bash
# Run tests with specific tags
bazel test --test_tag_filters=unit //...

# Exclude tests with certain tags
bazel test --test_tag_filters=-integration //...

# Combine filters
bazel test --test_tag_filters=unit,-slow //...

# Run tests matching pattern
bazel test //src/...  # All tests under src/
```

### Test Results

```bash
# View test logs
bazel test //src/app:app_test
cat bazel-testlogs/src/app/app_test/test.log

# View test XML results
cat bazel-testlogs/src/app/app_test/test.xml

# Check test status
bazel test //... && echo "All tests passed!"
```

## Test Organization

### Test Types and Tags

Organize tests by type using tags:

```python
# Unit tests
py_test(
    name = "unit_test",
    srcs = ["unit_test.py"],
    tags = ["unit", "fast"],
)

# Integration tests
py_test(
    name = "integration_test",
    srcs = ["integration_test.py"],
    tags = ["integration", "requires-docker"],
)

# End-to-end tests
py_test(
    name = "e2e_test",
    srcs = ["e2e_test.py"],
    tags = ["e2e", "slow"],
)
```

### Test Data and Fixtures

Include test data files:

```python
py_test(
    name = "data_test",
    srcs = ["data_test.py"],
    data = [
        "testdata/input.json",
        "testdata/expected.json",
        "//other/package:test_data",
    ],
)
```

Access test data in code:

```python
import os
import unittest

class DataTest(unittest.TestCase):
    def test_with_data(self):
        # Get path to test data
        testdata_dir = os.path.join(
            os.path.dirname(__file__),
            "testdata"
        )
        input_file = os.path.join(testdata_dir, "input.json")
        
        with open(input_file) as f:
            data = f.read()
        # Use data...
```

### Language-Specific Testing

#### Python Tests

```python
# BUILD file
py_test(
    name = "calculator_test",
    srcs = ["calculator_test.py"],
    deps = [
        ":calculator",
        "@pip//pytest",
    ],
)
```

```python
# calculator_test.py
import pytest
from calculator import add, subtract

def test_add():
    assert add(2, 3) == 5

def test_subtract():
    assert subtract(5, 3) == 2
```

#### JavaScript/TypeScript Tests

```javascript
// calculator.test.ts
import { add, subtract } from './calculator';

describe('Calculator', () => {
    test('adds numbers', () => {
        expect(add(2, 3)).toBe(5);
    });
    
    test('subtracts numbers', () => {
        expect(subtract(5, 3)).toBe(2);
    });
});
```

#### Go Tests

```go
// calculator_test.go
package calculator

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

## Advanced Build Techniques

### Selective Building

```bash
# Build only changed targets
bazel build //...

# Build targets that depend on X
bazel build $(bazel query 'rdeps(//..., //src/lib:utils)')

# Build specific file types
bazel build $(bazel query 'kind("py_binary", //...)')
```

### Build with Custom Flags

```bash
# Add custom defines
bazel build --define=version=1.0.0 //src/app:main

# Set custom configuration
bazel build --//config:feature=enabled //src/app:main

# Override toolchain
bazel build --python_version=3.11 //src/app:main
```

### Dependency Analysis

```bash
# View target dependencies
bazel query --output=graph 'deps(//src/app:main)' | dot -Tpng > deps.png

# Find why target depends on X
bazel query 'somepath(//src/app:main, //lib:utils)'

# Check for unused dependencies
bazel cquery 'deps(//src/app:main)' --output=build

# List all dependencies
bazel query 'deps(//src/app:main)'
```

### Build Event Protocol

Track detailed build information:

```bash
# Generate build event log
bazel build --build_event_json_file=bep.json //...

# Generate build event binary
bazel build --build_event_binary_file=bep.bin //...

# Analyze build events
cat bep.json | jq '.targetComplete'
```

## Continuous Testing

### Watch Mode with ibazel

Install ibazel for automatic rebuilds:

```bash
# Install ibazel
npm install -g @bazel/ibazel
# or
go install github.com/bazelbuild/bazel-watcher/cmd/ibazel@latest
```

Use watch mode:

```bash
# Watch and rebuild
ibazel build //src/app:all

# Watch and retest
ibazel test //src/app:app_test

# Watch and run
ibazel run //src/app:main
```

### Test Sharding

Split large test suites:

```python
py_test(
    name = "large_test",
    srcs = ["large_test.py"],
    shard_count = 4,  # Split into 4 shards
)
```

Run with sharding:

```bash
# Run all shards
bazel test //src:large_test

# Run specific shard
bazel test //src:large_test --test_sharding_strategy=explicit
```

### Flaky Test Handling

```bash
# Rerun flaky tests
bazel test --flaky_test_attempts=3 //...

# Run tests multiple times to detect flakiness
bazel test --runs_per_test=10 //src:potentially_flaky_test

# Mark tests as flaky in BUILD
py_test(
    name = "flaky_test",
    srcs = ["flaky_test.py"],
    flaky = True,  # Allow retries
)
```

## Build Optimization

### Caching Strategies

#### Local Cache

```bash
# Check cache info
bazel info output_base

# Size of cache
du -sh $(bazel info output_base)

# Clean cache
bazel clean
```

#### Remote Cache

Configure in `.bazelrc`:

```bash
# Add to .bazelrc
build --remote_cache=https://cache.example.com
build --remote_upload_local_results=true
```

Or use per-command:

```bash
bazel build --remote_cache=https://cache.example.com //...
```

#### Disk Cache

```bash
# Configure disk cache
bazel build --disk_cache=~/.cache/bazel //...

# Add to .bazelrc
build --disk_cache=~/.cache/bazel
build --experimental_disk_cache_gc_max_size=50GB
```

### Profiling Builds

```bash
# Generate profile
bazel build --profile=profile.json //...

# Analyze profile
bazel analyze-profile profile.json

# Generate HTML report
bazel analyze-profile --html profile.json > profile.html
open profile.html
```

### Parallel Execution

```bash
# Auto-detect CPU count
bazel build --jobs=auto //...

# Limit parallel jobs
bazel build --jobs=4 //...

# Configure in .bazelrc
build --jobs=auto
```

### Memory Management

```bash
# Limit memory usage
bazel build --local_ram_resources=HOST_RAM*.67 //...

# Limit per-action memory
bazel build --local_resources=4096,.5,1.0 //...

# Add to .bazelrc
build --local_ram_resources=HOST_RAM*.67
```

### Incremental Builds

Bazel automatically handles incremental builds:

```bash
# First build (full)
bazel build //...

# Edit one file
vim src/app/main.py

# Second build (only rebuilds what changed)
bazel build //...  # Much faster!
```

### Build Without Tests

```bash
# Build only, skip test targets
bazel build //...

# Build and test separately
bazel build //... && bazel test //...
```

## Troubleshooting Builds

### Build Failures

```bash
# Verbose build output
bazel build -s //src/app:main

# Show full command lines
bazel build --subcommands //src/app:main

# Show why target was rebuilt
bazel build --explain=explain.txt //src/app:main
cat explain.txt
```

### Test Failures

```bash
# Run with verbose output
bazel test --test_output=all --test_arg=-v //src/app:app_test

# Get detailed logs
bazel test //src/app:app_test
cat bazel-testlogs/src/app/app_test/test.log

# Run test without cache
bazel test --cache_test_results=no //src/app:app_test
```

### Common Issues

**Build is slow:**

- Enable remote caching
- Profile the build
- Reduce parallelism if memory-bound

**Tests are cached incorrectly:**

```bash
bazel test --cache_test_results=no //...
```

**Clean build needed:**

```bash
bazel clean
bazel build //...
```

**Corrupted cache:**

```bash
bazel clean --expunge
```

## CI/CD Integration

### CI Build Commands

```bash
# Full CI build
bazel build --config=ci //...

# Run tests with coverage
bazel coverage //...

# Generate coverage report
genhtml bazel-out/_coverage/_coverage_report.dat -o coverage_html
```

### Build Flags for CI

Add to `.bazelrc`:

```bash
# CI configuration
build:ci --announce_rc
build:ci --verbose_failures
build:ci --show_timestamps
build:ci --test_output=errors
build:ci --keep_going
```

## Best Practices

**DO**:

- ✅ Use `bazel test //...` to run all tests regularly
- ✅ Tag tests appropriately (unit, integration, e2e)
- ✅ Enable remote caching for teams
- ✅ Profile slow builds
- ✅ Keep test data in `testdata/` directories
- ✅ Use watch mode during active development

**DON'T**:

- ❌ Commit bazel-* symlinks
- ❌ Rely on system state in tests
- ❌ Make tests depend on order of execution
- ❌ Use `bazel clean` routinely (only when needed)
- ❌ Ignore flaky tests
- ❌ Skip CI tests to save time

## Next Steps

- Learn about [Formatting and Linting](./formatting-linting.md)
- Explore [Dependency Management](./dependency-management.md)
- Check [Troubleshooting](./troubleshooting.md) for common issues

---

**Back**: [User Guide](./README.md) | **Next**: [Formatting and Linting](./formatting-linting.md)
