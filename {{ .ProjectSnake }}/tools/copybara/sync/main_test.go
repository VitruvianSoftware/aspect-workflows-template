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
	"strings"
	"testing"
)

// fakeRunner is a test double for the runner interface.
// Each call pops the next exitCode from codes; if exhausted it returns the last one.
// callCount records how many times it was called.
type fakeRunner struct {
	codes        []int
	callCount    int
	capturedArgs [][]string
}

func (f *fakeRunner) run(image string, dockerArgs []string) (int, error) {
	f.callCount++
	f.capturedArgs = append(f.capturedArgs, dockerArgs)
	idx := f.callCount - 1
	if idx >= len(f.codes) {
		idx = len(f.codes) - 1
	}
	return f.codes[idx], nil
}

// noopSleeper is injected into tests so backoff sleeps are instant.
func noopSleeper(_ int) {}

// ---------------------------------------------------------------------------
// workflowName tests
// ---------------------------------------------------------------------------

func TestWorkflowName_ImportHyphenated(t *testing.T) {
	got := workflowName("import", "mcp-slack")
	if got != "import_mcp_slack" {
		t.Errorf("workflowName = %q, want %q", got, "import_mcp_slack")
	}
}

func TestWorkflowName_ExportSimple(t *testing.T) {
	got := workflowName("export", "devx")
	if got != "export_devx" {
		t.Errorf("workflowName = %q, want %q", got, "export_devx")
	}
}

func TestWorkflowName_MultiHyphen(t *testing.T) {
	got := workflowName("import", "nexus-agent")
	if got != "import_nexus_agent" {
		t.Errorf("workflowName = %q, want %q", got, "import_nexus_agent")
	}
}

// ---------------------------------------------------------------------------
// exit code classification
// ---------------------------------------------------------------------------

func TestExit0_IsSuccess(t *testing.T) {
	fr := &fakeRunner{codes: []int{0}}
	code := run([]string{"export", "devx"}, fr, noopSleeper)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if fr.callCount != 1 {
		t.Errorf("runner called %d times, want 1", fr.callCount)
	}
}

func TestExit4_IsSuccess(t *testing.T) {
	fr := &fakeRunner{codes: []int{4}}
	code := run([]string{"export", "devx"}, fr, noopSleeper)
	if code != 0 {
		t.Errorf("exit code = %d, want 0 (exit 4 = no-op)", code)
	}
}

// ---------------------------------------------------------------------------
// export: no retry
// ---------------------------------------------------------------------------

func TestExport_Failure_NoRetry(t *testing.T) {
	fr := &fakeRunner{codes: []int{1}}
	code := run([]string{"export", "devx"}, fr, noopSleeper)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if fr.callCount != 1 {
		t.Errorf("runner called %d times, want 1 (no retry for export)", fr.callCount)
	}
}

// ---------------------------------------------------------------------------
// import: retry up to 3 attempts
// ---------------------------------------------------------------------------

func TestImport_Failure_3Attempts(t *testing.T) {
	// All three attempts fail with exit 1.
	fr := &fakeRunner{codes: []int{1, 1, 1}}
	code := run([]string{"import", "devx"}, fr, noopSleeper)
	if code != 1 {
		t.Errorf("exit code = %d, want 1 (all attempts failed)", code)
	}
	if fr.callCount != 3 {
		t.Errorf("runner called %d times, want 3 (max attempts)", fr.callCount)
	}
}

func TestImport_SuccessOnRetry(t *testing.T) {
	// First attempt fails, second succeeds.
	fr := &fakeRunner{codes: []int{1, 0}}
	code := run([]string{"import", "devx"}, fr, noopSleeper)
	if code != 0 {
		t.Errorf("exit code = %d, want 0 (success on retry)", code)
	}
	if fr.callCount != 2 {
		t.Errorf("runner called %d times, want 2", fr.callCount)
	}
}

func TestImport_NoOpOnRetry_IsSuccess(t *testing.T) {
	// First attempt fails, second returns 4 (no-op).
	fr := &fakeRunner{codes: []int{1, 4}}
	code := run([]string{"import", "devx"}, fr, noopSleeper)
	if code != 0 {
		t.Errorf("exit code = %d, want 0 (exit 4 = success)", code)
	}
}

// ---------------------------------------------------------------------------
// unknown direction
// ---------------------------------------------------------------------------

func TestUnknownDirection_Exit2(t *testing.T) {
	fr := &fakeRunner{codes: []int{0}}
	code := run([]string{"sideways", "devx"}, fr, noopSleeper)
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if fr.callCount != 0 {
		t.Errorf("runner should not be called for unknown direction")
	}
}

// ---------------------------------------------------------------------------
// docker args contain the workflow name
// ---------------------------------------------------------------------------

func TestDockerArgs_ContainPinnedImage(t *testing.T) {
	fr := &fakeRunner{codes: []int{0}}
	run([]string{"export", "devx"}, fr, noopSleeper)
	if fr.callCount == 0 {
		t.Fatal("runner was not called")
	}
	// Check that COPYBARA_WORKFLOW appears in the docker args.
	found := false
	for _, arg := range fr.capturedArgs[0] {
		if strings.Contains(arg, "export_devx") {
			found = true
		}
	}
	if !found {
		t.Errorf("docker args %v don't contain workflow name export_devx", fr.capturedArgs[0])
	}
}

// ---------------------------------------------------------------------------
// --options flag is threaded through to COPYBARA_OPTIONS
// ---------------------------------------------------------------------------

func TestOptions_Default(t *testing.T) {
	fr := &fakeRunner{codes: []int{0}}
	run([]string{"export", "devx"}, fr, noopSleeper)
	found := false
	for _, arg := range fr.capturedArgs[0] {
		if arg == "COPYBARA_OPTIONS=--ignore-noop" {
			found = true
		}
	}
	if !found {
		t.Errorf("default --ignore-noop not in docker args: %v", fr.capturedArgs[0])
	}
}

func TestOptions_Custom(t *testing.T) {
	fr := &fakeRunner{codes: []int{0}}
	run([]string{"export", "devx", "--options", "--force --ignore-noop"}, fr, noopSleeper)
	found := false
	for _, arg := range fr.capturedArgs[0] {
		if arg == "COPYBARA_OPTIONS=--force --ignore-noop" {
			found = true
		}
	}
	if !found {
		t.Errorf("custom options not threaded through: %v", fr.capturedArgs[0])
	}
}

// ---------------------------------------------------------------------------
// sleeper receives correct backoff values (attempt*5)
// ---------------------------------------------------------------------------

func TestImport_BackoffValues(t *testing.T) {
	var sleepCalls []int
	sleeper := func(secs int) {
		sleepCalls = append(sleepCalls, secs)
	}
	fr := &fakeRunner{codes: []int{1, 1, 1}}
	run([]string{"import", "devx"}, fr, sleeper)
	// attempt 2 → sleep 2*5=10; attempt 3 → sleep 3*5=15
	// (no sleep before attempt 1, no sleep after final failure)
	want := []int{10, 15}
	if fmt.Sprint(sleepCalls) != fmt.Sprint(want) {
		t.Errorf("sleep calls = %v, want %v", sleepCalls, want)
	}
}
