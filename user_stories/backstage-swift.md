# Backstage Swift Bazel Starter

    # This is executable Markdown that's tested on CI.
    set -o errexit -o nounset -o xtrace
    alias ~~~=":<<'~~~sh'";:<<'~~~sh'

This repo includes:
- 🧱 Latest version of Bazel and dependencies
- 📦 Curated bazelrc flags via [bazelrc-preset.bzl]
- 🧰 Developer environment setup with [bazel_env.bzl]
- 🎨 `swift-format` (SwiftFormat), using rules_lint
- ✅ Pre-commit hooks for automatic linting and formatting
- 📚 Generic cross-platform Swift via rules_swift
- 🎭 Backstage template skeleton

## Try it out

> Before following these instructions, setup the developer environment by running <code>direnv allow</code> and follow any prompts.
> This ensures that tools we call in the following steps will be on the PATH.

First we create a tiny Swift program

~~~sh
mkdir -p hello_world
cat >hello_world/main.swift <<EOF
print("Hello from Swift")
EOF
~~~

We don't have any BUILD file generation for Swift yet,
so you're forced to create it manually.

~~~sh
touch hello_world/BUILD
buildozer 'new_load @build_bazel_rules_swift//swift:swift_binary.bzl swift_binary' hello_world:__pkg__
buildozer 'new swift_binary hello_world' hello_world:__pkg__
buildozer 'add srcs main.swift' hello_world:hello_world
~~~

Now you can run the program and assert that it produces the expected output.

~~~sh
output="$(bazel run hello_world | tail -1)"

[ "${output}" = "Hello from Swift" ] || {
    echo >&2 "Wanted output 'Hello from Swift' but got '${output}'"
    exit 1
}
~~~
