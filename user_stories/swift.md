# Swift Bazel Starter

    # This is executable Markdown that's tested on CI.
    # How is that possible? See https://gist.github.com/bwoods/1c25cb7723a06a076c2152a2781d4d49
    set -o errexit -o nounset -o xtrace
    alias ~~~=":<<'~~~sh'";:<<'~~~sh'

This repo includes:
- 🧱 Latest version of Bazel and dependencies
- 📦 Curated bazelrc flags via [bazelrc-preset.bzl]
- 🧰 Developer environment setup with [bazel_env.bzl]
- 🎨 `swift-format` (SwiftFormat) using rules_lint
- ✅ Pre-commit hooks for automatic linting and formatting
- 📚 Generic cross-platform Swift via rules_swift

> [!NOTE]
> You can customize languages and features with the interactive wizard in the <code>aspect init</code> command.
> <code>init</code> is an alternative to this starter repo, which was generated using the 'swift' preset.
> See https://docs.aspect.build/cli/overview

## Setup dev environment

First, we recommend you setup a Bazel-based developer environment with direnv.

1. install https://direnv.net/docs/installation.html
1. run <code>direnv allow</code> and follow the prompts to <code>bazel run //tools:bazel_env</code>

This isn't strictly required, but the commands which follow assume that needed tools are on the PATH,
so skipping `direnv` means you're responsible for installing them yourself.

## Try it out

First we create a tiny Swift program:

~~~sh
mkdir -p hello_world
cat >hello_world/main.swift <<EOF
print("Hello from Swift")
EOF
~~~

We don't have any BUILD file generation for Swift yet,
so you're forced to create it manually. We load `swift_binary` from
`//swift:defs.bzl` -- a thin project-local wrapper around rules_swift's
`swift_binary` that bundles the hermetic Swift runtime into the binary's runfiles
on Linux, so `bazel run`/`bazel test` can find libswiftCore at runtime.
~~~sh
touch hello_world/BUILD
buildozer 'new_load //swift:defs.bzl swift_binary' hello_world:__pkg__
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

## Formatting

We can format the code with SwiftFormat. Let's create some intentionally poorly formatted code
(the indentation is wrong and the brace is missing a leading space).

~~~sh
cat >hello_world/main.swift <<EOF
func main(){
print("Hello from Swift")
}
main()
EOF
~~~

Now format it:

~~~sh
format
~~~

Let's verify the code was fixed:

~~~sh
cat hello_world/main.swift
# -> func main() {
# ->     print("Hello from Swift")
# -> }
# -> main()
~~~

<!--
~~~sh
formatted=$(cat hello_world/main.swift)
echo "${formatted}" | grep -q "^func main() {$" && \
echo "${formatted}" | grep -q "^    print(\"Hello from Swift\")$" && \
echo "${formatted}" | grep -q "^}$" || {
    echo >&2 "Code was not properly formatted. Got:"
    echo >&2 "${formatted}"
    exit 1
}
~~~
-->
