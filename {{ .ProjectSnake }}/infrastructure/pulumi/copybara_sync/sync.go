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
{{ end }}{{ end }}// Package copybara_sync manages the GitHub auth resources that back the
// Copybara bidirectional sync between this monorepo and each of its standalone
// component repositories.
//
// For every synced component this creates:
//
//  1. A fresh ED25519 SSH key pair (tls.PrivateKey).
//  2. A WRITE deploy key on the STANDALONE repo, so the monorepo's export
//     workflow can push the component out to its standalone repo.
//  3. An Actions secret in the monorepo holding that key's PRIVATE half
//     (<PREFIX>_SYNC_SSH_KEY), consumed by the export workflow.
//  4. Two Actions secrets in the STANDALONE repo holding GitHub App dispatch
//     credentials (<PREFIX>_DISPATCH_APP_ID / <PREFIX>_DISPATCH_APP_PRIVATE_KEY),
//     consumed by the standalone repo's dispatch workflow to fire a
//     repository_dispatch back into the monorepo (the import trigger).
//
// Additionally, monorepo-wide App credentials are placed as both Actions and
// Dependabot secrets (SYNC_APP_ID / SYNC_APP_PRIVATE_KEY) for the Dependabot
// reconcile / auto-merge automation.
//
// The GitHub App itself is created MANUALLY by the operator; Pulumi only places
// its credentials, which are supplied as Pulumi config secrets (see README.md).
//
// Usage: call ManageSyncAuth from your Pulumi stack's main function, e.g.:
//
//	func main() {
//	    pulumi.Run(func(ctx *pulumi.Context) error {
//	        return copybara_sync.ManageSyncAuth(ctx)
//	    })
//	}
package copybara_sync

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi-tls/sdk/v5/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// githubOrg is the GitHub organisation (or user) that owns BOTH this monorepo
// and all standalone repos listed in syncedProjects below.
//
// EDIT: replace "YOUR_GITHUB_ORG" with your actual org/user before applying.
// See docs/copybara-bidi-sync.md for details.
const githubOrg = "YOUR_GITHUB_ORG" // EDIT: your GitHub org/user

// monorepoRepoName is the GitHub repository name of this monorepo.  The export
// workflow lives here and consumes the <PREFIX>_SYNC_SSH_KEY secret.
const monorepoRepoName = "{{ .ProjectKebab }}"

// syncedProject describes one component that is bidirectionally synced between
// the monorepo and a standalone repository.
//
// To onboard another component later, append an entry here and follow the
// onboarding runbook in docs/copybara-bidi-sync.md (§8f).
type syncedProject struct {
	// Name is the component / monorepo subfolder name (e.g. "my-service"). It is
	// upper-snaked to form the secret-name prefix (my-service → MY_SERVICE) and
	// is also used as the config-key prefix for the App credentials.
	Name string

	// StandaloneRepo is the name of the EXISTING standalone GitHub repository
	// this component is synced with (e.g. "my-service"). This package never
	// creates repositories — the standalone must already exist.
	StandaloneRepo string
}

// syncedProjects is the source of truth for which components have sync auth
// managed by Pulumi.  Seeded from the copybara_components value given at
// generation time; add entries here to onboard new components later.
{{- if .Scaffold.copybara_components }}
var syncedProjects = []syncedProject{
{{- range $c := splitList "," .Scaffold.copybara_components }}{{- $t := trim $c }}{{- if $t }}
	{Name: "{{ $t }}", StandaloneRepo: "{{ $t }}"},
{{- end }}{{- end }}
}
{{- else }}
// No components were seeded at generation time.  Add entries here as you
// onboard components; follow the runbook in docs/copybara-bidi-sync.md (§8f).
var syncedProjects = []syncedProject{}
{{- end }}

// secretPrefix converts a project name into the UPPER_SNAKE prefix used for its
// Actions secret names.  e.g. "my-service" → "MY_SERVICE".
func secretPrefix(projectName string) string {
	return strings.ToUpper(strings.ReplaceAll(projectName, "-", "_"))
}

// configKeyPrefix converts a project name into a safe Pulumi config-key prefix
// for the GitHub App credentials.  Config keys cannot contain '-' segments that
// Pulumi would interpret as namespaces, so hyphens are dropped and the name is
// lower-camel-cased.  e.g. "my-service" → "myService", yielding config keys
// "myServiceDispatchAppId" and "myServiceDispatchAppPrivateKey".
func configKeyPrefix(projectName string) string {
	parts := strings.Split(projectName, "-")
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(p)
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		b.WriteString(p[1:])
	}
	return b.String()
}

// ManageSyncAuth provisions the sync-auth resources for every synced project.
//
// Resources created per component:
//   - tls.NewPrivateKey         — ED25519 keypair
//   - github.NewRepositoryDeployKey — WRITE deploy key on the standalone repo
//   - github.NewActionsSecret   — <PREFIX>_SYNC_SSH_KEY in the monorepo
//   - github.NewActionsSecret   — <PREFIX>_DISPATCH_APP_ID in the standalone
//   - github.NewActionsSecret   — <PREFIX>_DISPATCH_APP_PRIVATE_KEY in the standalone
//
// Monorepo-wide resources (the single shared GitHub App):
//   - github.NewDependabotSecret — SYNC_APP_ID (monorepo)
//   - github.NewDependabotSecret — SYNC_APP_PRIVATE_KEY (monorepo)
//   - github.NewActionsSecret   — SYNC_APP_ID (monorepo)
//   - github.NewActionsSecret   — SYNC_APP_PRIVATE_KEY (monorepo)
//
// Required Pulumi config secrets (set with `pulumi config set --secret`):
//   - syncAppId              — GitHub App ID for the dispatch App
//   - syncAppPrivateKey      — GitHub App private key for the dispatch App
//   - <camel>DispatchAppId   — per-component App ID (same App, keyed per component)
//   - <camel>DispatchAppPrivateKey — per-component App private key
func ManageSyncAuth(ctx *pulumi.Context) error {
	cfg := config.New(ctx, "")

	for _, project := range syncedProjects {
		prefix := secretPrefix(project.Name)
		cfgPrefix := configKeyPrefix(project.Name)

		// 1. Create a fresh ED25519 key pair for the export push.
		privateKey, err := tls.NewPrivateKey(ctx, fmt.Sprintf("%s-sync-key", project.Name), &tls.PrivateKeyArgs{
			Algorithm: pulumi.String("ED25519"),
		})
		if err != nil {
			return err
		}

		// 2. Install the PUBLIC half as a WRITE deploy key on the STANDALONE repo
		//    so the monorepo's export workflow can push to it.
		_, err = github.NewRepositoryDeployKey(ctx, fmt.Sprintf("%s-standalone-deploy-key", project.Name), &github.RepositoryDeployKeyArgs{
			Title:      pulumi.String("copybara-sync (write)"),
			Repository: pulumi.String(project.StandaloneRepo),
			Key:        privateKey.PublicKeyOpenssh,
			ReadOnly:   pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// 3. Store the PRIVATE half as an Actions secret in the MONOREPO, where
		//    the export workflow reads it (<PREFIX>_SYNC_SSH_KEY).
		_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-sync-ssh-key-secret", project.Name), &github.ActionsSecretArgs{
			Repository:     pulumi.String(monorepoRepoName),
			SecretName:     pulumi.String(fmt.Sprintf("%s_SYNC_SSH_KEY", prefix)),
			PlaintextValue: privateKey.PrivateKeyOpenssh,
		})
		if err != nil {
			return err
		}

		// 4. Place the GitHub App dispatch credentials as Actions secrets in the
		//    STANDALONE repo, where its dispatch workflow reads them to fire a
		//    repository_dispatch back into the monorepo (the import trigger).
		//    Values come from Pulumi config secrets; the App is created manually.
		dispatchAppID := cfg.RequireSecret(fmt.Sprintf("%sDispatchAppId", cfgPrefix))
		dispatchAppPrivateKey := cfg.RequireSecret(fmt.Sprintf("%sDispatchAppPrivateKey", cfgPrefix))

		_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-dispatch-app-id-secret", project.Name), &github.ActionsSecretArgs{
			Repository:     pulumi.String(project.StandaloneRepo),
			SecretName:     pulumi.String(fmt.Sprintf("%s_DISPATCH_APP_ID", prefix)),
			PlaintextValue: dispatchAppID,
		})
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-dispatch-app-private-key-secret", project.Name), &github.ActionsSecretArgs{
			Repository:     pulumi.String(project.StandaloneRepo),
			SecretName:     pulumi.String(fmt.Sprintf("%s_DISPATCH_APP_PRIVATE_KEY", prefix)),
			PlaintextValue: dispatchAppPrivateKey,
		})
		if err != nil {
			return err
		}
	}

	// Monorepo-wide App credentials for the Dependabot reconcile + auto-merge
	// automation.  Dependabot-triggered workflow runs cannot read normal Actions
	// secrets, so the App id/key are placed as Dependabot secrets; the Actions-
	// secret twins cover non-Dependabot-context steps.  Both reuse the single
	// shared sync App.
	appID := cfg.RequireSecret("syncAppId")
	appKey := cfg.RequireSecret("syncAppPrivateKey")

	_, err := github.NewDependabotSecret(ctx, "monorepo-sync-app-id-dependabot", &github.DependabotSecretArgs{
		Repository:     pulumi.String(monorepoRepoName),
		SecretName:     pulumi.String("SYNC_APP_ID"),
		PlaintextValue: appID,
	})
	if err != nil {
		return err
	}
	_, err = github.NewDependabotSecret(ctx, "monorepo-sync-app-key-dependabot", &github.DependabotSecretArgs{
		Repository:     pulumi.String(monorepoRepoName),
		SecretName:     pulumi.String("SYNC_APP_PRIVATE_KEY"),
		PlaintextValue: appKey,
	})
	if err != nil {
		return err
	}
	_, err = github.NewActionsSecret(ctx, "monorepo-sync-app-id-actions", &github.ActionsSecretArgs{
		Repository:     pulumi.String(monorepoRepoName),
		SecretName:     pulumi.String("SYNC_APP_ID"),
		PlaintextValue: appID,
	})
	if err != nil {
		return err
	}
	_, err = github.NewActionsSecret(ctx, "monorepo-sync-app-key-actions", &github.ActionsSecretArgs{
		Repository:     pulumi.String(monorepoRepoName),
		SecretName:     pulumi.String("SYNC_APP_PRIVATE_KEY"),
		PlaintextValue: appKey,
	})
	if err != nil {
		return err
	}

	return nil
}
