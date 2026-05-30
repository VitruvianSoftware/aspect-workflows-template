#!/usr/bin/env bash
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
{{ end }}{{ end }}set -e

# Locate the prettier binary in the runfiles
# It is expected to be in the same directory as this script (in runfiles)
DIR="$(dirname "$0")"

# The target name is "prettier_bin", so the binary should be named "prettier_bin"
# However, aspect_rules_js might put it in a subdirectory or give it a different name structure in runfiles.
# Based on observation, it is in "prettier_bin_/prettier_bin".

PRETTIER_BIN="$DIR/prettier_bin_/prettier_bin"

if [[ ! -f "$PRETTIER_BIN" ]]; then
    # Fallback: try to find it in the current directory
    FOUND=$(find "$DIR" -name "prettier_bin" -type f | head -n 1)
    if [[ -n "$FOUND" ]]; then
        PRETTIER_BIN="$FOUND"
    else
        echo "ERROR: Could not find prettier_bin in $DIR" >&2
        exit 1
    fi
fi

ARGS=()
for arg in "$@"; do
  if [[ -L "$arg" ]]; then
    # Skip symlinks to avoid Prettier errors
    continue
  fi
  ARGS+=("$arg")
done

exec "$PRETTIER_BIN" "${ARGS[@]}"
