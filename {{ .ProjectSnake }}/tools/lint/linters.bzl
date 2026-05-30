{{ if .Scaffold.license }}{{ if eq .Scaffold.license_id `Apache-2.0` }}# Copyright {{ now.Year }} {{ .Scaffold.copyright_holder }}
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
{{ else if eq .Scaffold.license_id `MIT` }}# Copyright (c) {{ now.Year }} {{ .Scaffold.copyright_holder }}
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
{{ else if eq .Scaffold.license_id `BSD-3-Clause` }}# Copyright (c) {{ now.Year }} {{ .Scaffold.copyright_holder }} All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
{{ else }}# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
{{ end }}{{ end }}"Define linter aspects"

{{ if .Computed.cpp -}}
load("@aspect_rules_lint//lint:clang_tidy.bzl", "lint_clang_tidy_aspect")
{{ end -}}
{{ if .Computed.javascript -}}
load("@aspect_rules_lint//lint:eslint.bzl", "lint_eslint_aspect")
{{ end -}}
load("@aspect_rules_lint//lint:lint_test.bzl", "lint_test")
{{ if .Computed.kotlin -}}
load("@aspect_rules_lint//lint:ktlint.bzl", "lint_ktlint_aspect")
{{ end -}}
{{ if .Computed.java -}}
load("@aspect_rules_lint//lint:pmd.bzl", "lint_pmd_aspect")
{{ end -}}
{{ if .Computed.python -}}
load("@aspect_rules_lint//lint:ruff.bzl", "lint_ruff_aspect")
load("@aspect_rules_lint//lint:ty.bzl", "lint_ty_aspect")
{{ end -}}
{{ if .Computed.shell }}
load("@aspect_rules_lint//lint:shellcheck.bzl", "lint_shellcheck_aspect")
{{ end -}}
{{ if .Computed.rust -}}
load("@aspect_rules_lint//lint:clippy.bzl", "lint_clippy_aspect")
{{ end -}}
{{ if .Computed.ruby -}}
load("@aspect_rules_lint//lint:rubocop.bzl", "lint_rubocop_aspect")
{{ end -}}

{{ if .Computed.cpp -}}
clang_tidy = lint_clang_tidy_aspect(
    binary = Label("//tools/lint:clang_tidy"),
    configs = [Label("//:.clang-tidy")],
    lint_target_headers = True,
    angle_includes_are_system = False,
    verbose = False,
)
{{ end -}}
{{ if .Computed.kotlin -}}
ktlint = lint_ktlint_aspect(
    binary = Label("@com_github_pinterest_ktlint//file"),
    editorconfig = Label("//:.editorconfig"),
    baseline_file = Label("//:ktlint-baseline.xml"),
)
{{ end -}}
{{ if .Computed.java -}}
pmd = lint_pmd_aspect(
    binary = Label(":pmd"),
    rulesets = [Label("//:pmd.xml")],
)
{{ end -}}
{{ if .Computed.javascript -}}
eslint = lint_eslint_aspect(
    binary = Label(":eslint"),
    # We trust that eslint will locate the correct configuration file for a given source file.
    # See https://eslint.org/docs/latest/use/configure/configuration-files#cascading-and-hierarchy
    configs = [
        Label("//:eslintrc"),
        # if the repository has nested eslintrc files, they must be added here as well
    ],
)

eslint_test = lint_test(aspect = eslint)
{{ end -}}
{{ if .Computed.python -}}
ruff = lint_ruff_aspect(
    binary = "@multitool//tools/ruff",
    configs = [
        Label("//:pyproject.toml"),
        # if the repository has nested ruff.toml files, they must be added here as well
    ],
)

ruff_test = lint_test(aspect = ruff)

ty = lint_ty_aspect(
    binary = Label("@aspect_rules_lint//lint:ty_bin"),
    config = Label("@//:pyproject.toml"),
)

{{ end -}}
{{ if .Computed.shell }}
shellcheck = lint_shellcheck_aspect(
    binary = "@multitool//tools/shellcheck",
    config = Label("//:.shellcheckrc"),
)

shellcheck_test = lint_test(aspect = shellcheck)
{{ end -}}
{{ if .Computed.rust -}}
clippy = lint_clippy_aspect(
    config = Label("//:.clippy.toml"),
)
{{ end -}}
{{ if .Computed.swift -}}
# Swift linting is a deferred best-effort: aspect_rules_lint does not ship a
# SwiftLint aspect yet, so no lint aspect is wired here. `swift-format`/SwiftFormat
# still provides formatting via //tools/format. Add a SwiftLint aspect here if and
# when rules_lint supports it.
{{ end -}}
{{ if .Computed.ruby -}}
rubocop = lint_rubocop_aspect(
    binary = Label("@bundle//bin:rubocop"),
    configs = [Label("//:.rubocop.yml")],
)
{{ end -}}