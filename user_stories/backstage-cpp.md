# Backstage C/C++ Bazel Starter

    # This is executable Markdown that's tested on CI.
    set -o errexit -o nounset -o xtrace
    alias ~~~=":<<'~~~sh'";:<<'~~~sh'

This repo includes:
- 🧱 Latest version of Bazel and dependencies
- 📦 Curated bazelrc flags via [bazelrc-preset.bzl]
- 🧰 Developer environment setup with [bazel_env.bzl]
- 🎨 `clang-format` and `clang-tidy`, using rules_lint
- ✅ Pre-commit hooks for automatic linting and formatting
- 🎭 Backstage template skeleton

## Try it out

> Before following these instructions, setup the developer environment by running <code>direnv allow</code> and follow any prompts.
> This ensures that tools we call in the following steps will be on the PATH.

To start with, let’s get the simplest possible program to execute.

Lets create an application that prints “hello world” to the standard out:

~~~sh
mkdir src
>src/hello.c cat <<EOF
#include <stdio.h>

int main(int argc, char const *argv[])
{
    printf("hello world!");
    return 0;
}
EOF
~~~

Next run the BUILD file generator to produce a <code>cc_binary</code> target:

~~~sh
bazel run gazelle
~~~

You can now run the program to see the output.

~~~sh
bazel run src:hello
~~~

## Add a dependency

Many libraries already have <code>BUILD</code> files in the Bazel Central Registry (BCR).
For example, we can search for the <code>libmagic</code>, and find that it’s there:

https://registry.bazel.build/modules/libmagic


To install it, just add to <code>MODULE.bazel</code>

~~~sh
echo 'bazel_dep(name = "libmagic", version = "5.46")' >> MODULE.bazel
~~~

Let's add another program <code>src/magic.c</code> that depends on libmagic:

~~~sh
>src/magic.c cat <<EOF
#include <stdio.h>
#include <stdlib.h>
#include <magic.h>

int main(int argc, char *argv[]) {
    // Check if filename is provided
    if (argc != (int)2) {
        fprintf(stderr, "Usage: %s <filename>\n", argv[0]);
        return 1;
    }

    // Initialize magic handle
    magic_t magic = magic_open(MAGIC_MIME_TYPE);
    if (magic == NULL) {
        fprintf(stderr, "Failed to initialize libmagic\n");
        return 1;
    }

    // Load magic database
    if (magic_load(magic, NULL) != 0) {
        fprintf(stderr, "Cannot load magic database: %s\n", magic_error(magic));
        magic_close(magic);
        return 1;
    }

    // Get MIME type
    const char *mime_type = magic_file(magic, argv[1]);
    if (mime_type == NULL) {
        fprintf(stderr, "Error determining MIME type: %s\n", magic_error(magic));
        magic_close(magic);
        return 1;
    }

    // Print result
    printf("MIME type: %s\n", mime_type);
