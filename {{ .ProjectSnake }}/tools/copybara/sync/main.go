{{ if .Scaffold.license }}{{ if eq .Scaffold.license_id `Apache-2.0` }}// Copyright {{ now.Year }} {{ .Scaffold.copyright_holder }}
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
{{ else if eq .Scaffold.license_id `MIT` }}// Copyright (c) {{ now.Year }} {{ .Scaffold.copyright_holder }}
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
{{ else if eq .Scaffold.license_id `BSD-3-Clause` }}// Copyright (c) {{ now.Year }} {{ .Scaffold.copyright_holder }} All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
{{ else }}// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
{{ end }}{{ end }}package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	pinnedImage    = "olivr/copybara@sha256:87e2e9089344e64693faebb2ee0ed33b8797358c0420b0fa98325ca611e98679"
	defaultOptions = "--ignore-noop"
	maxAttempts    = 3
)

// runner is the interface the sync logic uses to invoke Docker.
// The real implementation shells out; tests inject a fake.
type runner interface {
	run(image string, dockerArgs []string) (exitCode int, err error)
}

// dockerRunner is the real implementation: calls `docker run`.
type dockerRunner struct{}

func (dockerRunner) run(image string, dockerArgs []string) (int, error) {
	args := append([]string{"run"}, dockerArgs...)
	args = append(args, image, "copybara")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

// workflowName derives the Copybara workflow identifier:
// ("import", "mcp-slack") → "import_mcp_slack".
func workflowName(direction, component string) string {
	return direction + "_" + strings.ReplaceAll(component, "-", "_")
}

// isSuccess treats exit codes 0 and 4 (NO_OP) as success.
func isSuccess(code int) bool {
	return code == 0 || code == 4
}

// resolveWorkspace returns the workspace directory: prefers BUILD_WORKSPACE_DIRECTORY
// (set by `bazel run`), then GITHUB_WORKSPACE, then current working directory.
func resolveWorkspace() string {
	if v := os.Getenv("BUILD_WORKSPACE_DIRECTORY"); v != "" {
		return v
	}
	if v := os.Getenv("GITHUB_WORKSPACE"); v != "" {
		return v
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// run is the testable core. args = [direction, component, (--options "<opts>")?]
// r is the runner used to invoke Docker; sleeper is called with seconds to sleep (for retry backoff).
// Returns the process exit code.
func run(args []string, r runner, sleeper func(int)) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sync <export|import> <component> [--options \"<copybara cli opts>\"]")
		return 2
	}

	direction := args[0]
	component := args[1]
	options := defaultOptions

	// Parse optional --options flag.
	for i := 2; i < len(args)-1; i++ {
		if args[i] == "--options" {
			options = args[i+1]
		}
	}

	switch direction {
	case "export", "import":
		// valid
	default:
		fmt.Fprintf(os.Stderr, "unknown direction: %s (use export|import)\n", direction)
		return 2
	}

	wf := workflowName(direction, component)
	workspace := resolveWorkspace()
	home := os.Getenv("HOME")

	dockerArgs := []string{
		"--rm",
		"-v", workspace + ":/usr/src/app",
		"-v", home + "/.ssh/id_rsa:/root/.ssh/id_rsa",
		"-v", home + "/.ssh/known_hosts:/root/.ssh/known_hosts",
		"-v", workspace + "/tools/copybara/copy.bara.sky:/root/copy.bara.sky",
		"-v", home + "/.git-credentials:/root/.git-credentials",
		"-v", home + "/.gitconfig:/root/.gitconfig",
		"-e", "COPYBARA_CONFIG=/root/copy.bara.sky",
		"-e", "COPYBARA_WORKFLOW=" + wf,
		"-e", "COPYBARA_OPTIONS=" + options,
	}

	attempt := 1
	for {
		exitCode, err := r.run(pinnedImage, dockerArgs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "docker error: %v\n", err)
			return 1
		}

		if isSuccess(exitCode) {
			fmt.Printf("Copybara OK (exit %d; 4 = no-op/up-to-date)\n", exitCode)
			return 0
		}

		// Failure path.
		if direction != "import" || attempt >= maxAttempts {
			fmt.Fprintf(os.Stderr, "Copybara failed (exit %d", exitCode)
			if direction == "import" && attempt >= maxAttempts {
				fmt.Fprintf(os.Stderr, " after %d attempts", maxAttempts)
			}
			fmt.Fprintln(os.Stderr, ")")
			return exitCode
		}

		// Import retry: sleep attempt*5 seconds then retry.
		attempt++
		fmt.Printf("Copybara exit %d (attempt %d/%d) — retrying after re-fetch (likely a concurrent import race)\n",
			exitCode, attempt-1, maxAttempts)
		sleeper(attempt * 5)
	}
}

func main() {
	sleeper := func(secs int) { time.Sleep(time.Duration(secs) * time.Second) }
	os.Exit(run(os.Args[1:], dockerRunner{}, sleeper))
}
