# Backstage Ruby Bazel Starter

    # This is executable Markdown that's tested on CI.
    set -o errexit -o nounset -o xtrace
    alias ~~~=":<<'~~~sh'";:<<'~~~sh'

This repo includes:
- 🧱 Latest version of Bazel and dependencies
- 📦 Curated bazelrc flags via [bazelrc-preset.bzl]
- 🧰 Developer environment setup with [bazel_env.bzl]
- 🎨 `rubocop` and `standard`, using rules_lint
- ✅ Pre-commit hooks for automatic linting and formatting
- 🎭 Backstage template skeleton

## Try it out

> Before following these instructions, setup the developer environment by running <code>direnv allow</code> and follow any prompts.
> This ensures that tools we call in the following steps will be on the PATH.

Write a simple Ruby application:

~~~sh
mkdir app
>app/hello.rb cat <<'EOF'
# frozen_string_literal: true

puts "Hello from Bazel + Ruby!"
EOF
~~~

Write a BUILD file:

~~~sh
>app/BUILD cat <<EOF
load("@rules_ruby//ruby:defs.bzl", "rb_binary")

rb_binary(
    name = "hello",
    srcs = ["hello.rb"],
    main = "hello.rb",
)
EOF
~~~

Run it to see the result:

~~~sh
output=$(bazel run //app:hello | tail -1)
~~~

Let's verify the application output matches expectation:

~~~sh
echo "${output}" | grep -q "Hello from Bazel + Ruby!" || {
    echo >&2 "Wanted output containing 'Hello from Bazel + Ruby!' but got '${output}'"
    exit 1
}
~~~
