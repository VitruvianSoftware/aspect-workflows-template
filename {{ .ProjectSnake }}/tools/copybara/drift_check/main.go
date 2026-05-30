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
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// defaultExcluded returns the set of basenames that are never compared.
// These are context-specific files that intentionally never cross the sync
// boundary: BUILD/BUILD.bazel (monorepo-only Bazel files), package-lock.json
// (standalone-only npm lockfile), sync-to-monorepo.yaml (standalone-only
// dispatch workflow), and .git (version control metadata).
// ---------------------------------------------------------------------------

func defaultExcluded() map[string]bool {
	return map[string]bool{
		".git":                  true,
		"BUILD":                 true,
		"BUILD.bazel":           true,
		"package-lock.json":     true,
		"sync-to-monorepo.yaml": true,
	}
}

// ---------------------------------------------------------------------------
// compareDirs recursively walks dirA and dirB, ignoring entries whose
// BASENAME is in excludeBasenames at any depth. It returns the list of
// relative paths that differ (content mismatch, or present in only one side).
// ---------------------------------------------------------------------------

func compareDirs(dirA, dirB string, excludeBasenames map[string]bool) ([]string, error) {
	seen := map[string]bool{}
	var diffs []string

	// Walk dirA: check each file against dirB.
	err := filepath.Walk(dirA, func(pathA string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(pathA)
		if excludeBasenames[base] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dirA, pathA)
		seen[rel] = true
		pathB := filepath.Join(dirB, rel)
		contA, readErrA := os.ReadFile(pathA)
		if readErrA != nil {
			return readErrA
		}
		contB, readErrB := os.ReadFile(pathB)
		if os.IsNotExist(readErrB) {
			diffs = append(diffs, rel)
			return nil
		}
		if readErrB != nil {
			return readErrB
		}
		if !bytes.Equal(contA, contB) {
			diffs = append(diffs, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Walk dirB: find files not in dirA (content diffs already handled above).
	err = filepath.Walk(dirB, func(pathB string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(pathB)
		if excludeBasenames[base] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dirB, pathB)
		if !seen[rel] {
			diffs = append(diffs, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return diffs, nil
}

// ---------------------------------------------------------------------------
// SSH / git helpers (not unit-tested; require real repos + keys)
// ---------------------------------------------------------------------------

// envKey returns the value of ${PREFIX_SYNC_SSH_KEY} for a component name.
// e.g. "mcp-slack" → MCP_SLACK_SYNC_SSH_KEY
func envKey(component string) string {
	prefix := strings.ToUpper(strings.ReplaceAll(component, "-", "_"))
	return os.Getenv(prefix + "_SYNC_SSH_KEY")
}

// writeSSHKey writes the key to ~/.ssh/id_rsa (0600) stripping one trailing newline.
func writeSSHKey(key string) error {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return err
	}
	// Strip a single trailing newline (match bash ${key%$'\n'}).
	key = strings.TrimSuffix(key, "\n")
	return os.WriteFile(filepath.Join(sshDir, "id_rsa"), []byte(key), 0600)
}

// knownHostsPath returns the path to the known_hosts file.
func knownHostsPath() string {
	return filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
}

// scanGitHub runs ssh-keyscan to populate known_hosts with github.com keys.
func scanGitHub() error {
	khPath := knownHostsPath()
	f, err := os.OpenFile(khPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	f.Close()
	cmd := exec.Command("ssh-keyscan", "-t", "rsa,ed25519", "github.com")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	return os.WriteFile(khPath, out, 0600)
}

// gitClone clones org/component into destDir using the written SSH key.
func gitClone(org, component, destDir string) error {
	remote := fmt.Sprintf("git@github.com:%s/%s.git", org, component)
	sshCmd := fmt.Sprintf("ssh -i %s/.ssh/id_rsa -o UserKnownHostsFile=%s",
		os.Getenv("HOME"), knownHostsPath())
	cmd := exec.Command("git", "clone", "--no-tags", remote, destDir)
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND="+sshCmd)
	_, err := cmd.Output()
	return err
}

// isSeeded returns true when the cloned standalone repo has a commit whose
// message contains MONOREPO_REV_ID (indicating it has been synced from the monorepo).
func isSeeded(repoDir string) bool {
	out, err := exec.Command("git", "-C", repoDir,
		"log", "--grep=MONOREPO_REV_ID", "-1", "--format=%H").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

// ---------------------------------------------------------------------------
// run is the main entry point for the CLI logic.
// args: [--org <github_org>] <monorepo_workspace_dir> <component...>
// ---------------------------------------------------------------------------

func run(args []string) int {
	fs := flag.NewFlagSet("drift_check", flag.ContinueOnError)
	org := fs.String("org", "", "GitHub organisation (required)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()
	if *org == "" || len(remaining) < 2 {
		fmt.Fprintln(os.Stderr, "usage: drift_check --org <github_org> <monorepo_workspace_dir> <component...>")
		return 2
	}
	workspaceDir := remaining[0]
	components := remaining[1:]

	excluded := defaultExcluded()

	drift := false
	checked := 0

	for _, comp := range components {
		key := envKey(comp)
		if key == "" {
			fmt.Printf("::warning::[%s] no deploy key configured — skipping\n", comp)
			continue
		}

		if err := writeSSHKey(key); err != nil {
			fmt.Printf("::warning::[%s] failed to write SSH key: %v — skipping\n", comp, err)
			continue
		}
		if err := scanGitHub(); err != nil {
			fmt.Printf("::warning::[%s] ssh-keyscan failed: %v — skipping\n", comp, err)
			continue
		}

		tmpDir, err := os.MkdirTemp("", "drift_check_"+comp+"_*")
		if err != nil {
			fmt.Printf("::warning::[%s] cannot create tmp dir: %v — skipping\n", comp, err)
			continue
		}
		defer os.RemoveAll(tmpDir)

		peerDir := filepath.Join(tmpDir, comp+"-standalone")
		if err := gitClone(*org, comp, peerDir); err != nil {
			fmt.Printf("::warning::[%s] standalone clone failed — skipping\n", comp)
			continue
		}

		if !isSeeded(peerDir) {
			fmt.Printf("[%s] not yet seeded (no MONOREPO_REV_ID in standalone history) — skipping\n", comp)
			continue
		}

		checked++
		monorepoSubtree := filepath.Join(workspaceDir, comp)
		diffs, err := compareDirs(monorepoSubtree, peerDir, excluded)
		if err != nil {
			fmt.Printf("::warning::[%s] compareDirs error: %v — skipping\n", comp, err)
			continue
		}

		if len(diffs) == 0 {
			fmt.Printf("[%s] in sync — %s/%s/ matches the standalone.\n", comp, workspaceDir, comp)
		} else {
			fmt.Printf("::error title=Copybara sync drift (%s)::%s monorepo subtree and standalone have diverged\n", comp, comp)
			fmt.Printf("----- [%s] divergence (< monorepo  > standalone) -----\n", comp)
			for _, d := range diffs {
				monoPath := filepath.Join(monorepoSubtree, d)
				peerPath := filepath.Join(peerDir, d)
				_, monoErr := os.Stat(monoPath)
				_, peerErr := os.Stat(peerPath)
				switch {
				case os.IsNotExist(monoErr):
					fmt.Printf("> %s (only in standalone)\n", d)
				case os.IsNotExist(peerErr):
					fmt.Printf("< %s (only in monorepo)\n", d)
				default:
					fmt.Printf("< %s (content differs)\n", d)
				}
			}
			drift = true
		}
	}

	fmt.Printf("Checked %d seeded component(s).\n", checked)
	if drift {
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
