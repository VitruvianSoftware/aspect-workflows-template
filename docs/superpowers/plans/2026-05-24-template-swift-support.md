# Template Swift Support Implementation Plan (Spec 2A)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Swift as a first-class, cross-platform language to `aspect-workflows-template` (branch `platform-v2.0`) — `rules_swift`, format/lint, scaffold/preset/CI/Pulumi/catalog wiring, and a `user_stories/swift.md` verification script — so `scaffold new --preset=swift` yields a project that builds on Linux CI, and `swift` + `backstage-swift` starter repos exist.

**Architecture:** Mirror the existing **`rust`** language wiring — it's the closest analog (compiled, no Gazelle, manual `BUILD`, `rules_lint` format/lint). The template ships **no static per-language sample**; per-preset verification is an **executable `user_stories/<preset>.md`** script (run in CI) that creates+builds+runs a tiny program. Swift uses **`rules_swift` only** (no `rules_apple` — that's Spec 2B). All work is on `platform-v2.0`.

**Tech Stack:** scaffold (hay-kot/scaffold) Go-templating, Bazel/bzlmod, `rules_swift`, `swift-format`, `rules_lint`, GitHub Actions, Pulumi (Go).

**Prerequisites:** on `platform-v2.0`; confirm the current `rules_swift` BCR version + its bzlmod toolchain usage at <https://registry.bazel.build/modules/rules_swift> (the contributor guide's `1.13.0` is stale).

> **Correction vs. the spec:** the spec referenced a `{{ .ProjectSnake }}/swift/hello/` sample. The template has **no** static per-language samples — verification is the executable `user_stories/<preset>.md` (see `user_stories/rust.md`). This plan uses that mechanism. End goal/success criteria are unchanged.

---

### Task 1: scaffold.yaml — Swift language option + presets

**Files:** Modify `scaffold.yaml`

- [ ] **Step 1:** In the `langs` question's `options:` list, add `- Swift` (place near `Rust`; match the existing human-label style).
- [ ] **Step 2:** Under `computed:`, add (mirroring the `rust:` line): `swift: "{{ has \"Swift\" .Scaffold.langs }}"`.
- [ ] **Step 3:** Under `presets:`, add `swift:` and `backstage-swift:` blocks copied from `rust:`/`backstage-rust:`, with `langs: ['Swift']` (and `lint: true`; `backstage-swift` also `backstage: true`).
- [ ] **Step 4:** Commit: `feat(scaffold): add Swift language + swift/backstage-swift presets`.

### Task 2: MODULE.bazel — rules_swift toolchain (gated)

**Files:** Modify `{{ .ProjectSnake }}/MODULE.bazel`

- [ ] **Step 1:** Read the existing `{{- if .Computed.rust }} … {{- end }}` rules_rust block (~line 148) as the pattern. After it, add a `{{- if .Computed.swift }} … {{- end }}` block that adds `bazel_dep(name = "rules_swift", version = "<current BCR version>")` and registers the Swift toolchain exactly per `rules_swift`'s bzlmod docs (`use_extension` + `swift.toolchain` / `register_toolchains` / `use_repo` as required). Verify the exact extension label + version against the BCR registry page.
- [ ] **Step 2:** Verify resolution: `scaffold new --preset=swift --no-prompt /tmp/sw` (or `--output-dir`), then in the output `bazel mod deps | grep rules_swift` shows it resolved.
- [ ] **Step 3:** Commit: `feat(module): add rules_swift toolchain gated on .Computed.swift`.

### Task 3: Format + lint (swift-format; SwiftLint best-effort)

**Files:** Modify `{{ .ProjectSnake }}/tools/format/BUILD`, `{{ .ProjectSnake }}/tools/lint/linters.bzl` (+ `tools/lint/BUILD`), and any `.swift-format` config

- [ ] **Step 1:** Find how `rustfmt` is wired into the format tool and mirror it to add `swift-format`, gated on `{{- if .Computed.swift }}`. Add a `.swift-format` config file at the template root if swift-format needs one (gate its inclusion via a `features` entry in `scaffold.yaml`).
- [ ] **Step 2:** Check whether `rules_lint` supports Swift/SwiftLint (its docs/`example`). If yes, mirror the `clippy` lint wiring in `linters.bzl` for SwiftLint, gated on swift. If no, add a brief comment noting Swift linting is a deferred best-effort follow-up — do **not** block the task on it.
- [ ] **Step 3:** Commit: `feat(tools): swift-format (+ SwiftLint if supported) for Swift`.

### Task 4: README section + user_stories verification scripts

**Files:** Modify `{{ .ProjectSnake }}/README.bazel.md`; Create `user_stories/swift.md` and `user_stories/backstage-swift.md`

- [ ] **Step 1:** Add a `{{- if .Computed.swift }} … {{- end }}` Swift dev section to `README.bazel.md`, mirroring the Rust section (build/test/format commands).
- [ ] **Step 2:** Create `user_stories/swift.md` mirroring `user_stories/rust.md`'s executable-markdown structure (same `set -o errexit -o nounset -o xtrace` + `~~~sh` literate-script header). Body:

```sh
mkdir -p hello_world
cat >hello_world/main.swift <<'EOF'
print("Hello from Swift")
EOF
```
```sh
# Swift has no Gazelle; create the BUILD manually with buildozer.
touch hello_world/BUILD
buildozer 'new_load @rules_swift//swift:swift_binary.bzl swift_binary' hello_world:__pkg__
buildozer 'new swift_binary hello_world' hello_world:__pkg__
buildozer 'add srcs main.swift' hello_world:hello_world
```
```sh
output="$(bazel run hello_world | tail -1)"
[ "${output}" = "Hello from Swift" ] || { echo >&2 "Wanted 'Hello from Swift' but got '${output}'"; exit 1; }
```
Then a `format` section like rust.md (write poorly-formatted Swift, run `format`, assert it was fixed). Verify the exact `rules_swift` load label (`@rules_swift//swift:swift_binary.bzl` vs `@build_bazel_rules_swift//...`) against the version pinned in Task 2.
- [ ] **Step 3:** Create `user_stories/backstage-swift.md` mirroring `user_stories/backstage-rust.md`.
- [ ] **Step 4:** Commit: `docs(user-stories): add swift + backstage-swift verification scripts`.

### Task 5: CI matrix — verify Swift builds (no infra needed)

**Files:** Modify `.github/workflows/ci.yaml`

- [ ] **Step 1:** Add `swift` to the `test` job's `preset` matrix list. `ci.yaml` scaffolds + builds + runs the user-story script; it does **not** push to starter repos, so this verifies Swift end-to-end without the starter repo existing yet.
- [ ] **Step 2:** Commit and push `platform-v2.0`. Watch the run (`gh run watch … --exit-status`); confirm the `swift` matrix leg is green on ubuntu. If the Swift toolchain fails to build on Linux, debug (`superpowers:systematic-debugging`) before continuing.

### Task 6: Infra rollout — Pulumi presets + create the two starters (careful)

**Files:** Modify `infrastructure/pulumi/pkg/github_repos/starter_repos.go`

- [ ] **Step 1:** Add `"swift"` and `"backstage-swift"` to the `presets` slice.
- [ ] **Step 2:** `export GITHUB_TOKEN="$(gh auth token)" GITHUB_OWNER=VitruvianSoftware` then `pulumi preview --stack dev --json`. Confirm the diff is **only**: 2 new `Repository` (create), 2 new `tls PrivateKey` (create), 2 new `RepositoryDeployKey` (create), 2 new `ActionsSecret` (create) — and **no** deletes/replaces on existing repos. STOP and escalate if anything else appears.
- [ ] **Step 3:** `pulumi up --yes --stack dev`; verify `swift` + `backstage-swift` repos exist (`gh api repos/VitruvianSoftware/swift`), with `is_template=true`, deploy keys, and `STARTER_DEPLOY_SWIFT` / `STARTER_DEPLOY_BACKSTAGE_SWIFT` secrets. Commit: `feat(iac): add swift + backstage-swift starter repos`.

### Task 7: Deliver matrix + catalog — push to the new starters

**Files:** Modify `.github/workflows/deliver.yaml`, `{{ .ProjectSnake }}/catalog-info.yaml` (and `{{ .ProjectSnake }}/template.yaml` if it enumerates languages)

- [ ] **Step 1:** Add `swift` and `backstage-swift` to `deliver.yaml`'s `preset` matrix.
- [ ] **Step 2:** Add `swift`/`backstage-swift` branches to `catalog-info.yaml` (mirror the existing rust branches: name, tags, repo link) and to `template.yaml`'s language conditionals if present.
- [ ] **Step 3:** Commit and push `platform-v2.0`. Watch the deliver run; confirm the `swift` and `backstage-swift` legs scaffold and force-push to the new starter repos (per the standing always-watch preference).

### Task 8: Commit the planning docs

**Files:** `docs/superpowers/specs/2026-05-24-template-swift-support-design.md`, `docs/superpowers/plans/2026-05-24-template-swift-support.md`

- [ ] **Step 1:** These were left uncommitted during brainstorming. Commit them on `platform-v2.0` alongside the implementation (e.g., as part of Task 1's commit or a dedicated `docs:` commit) so they're not a standalone deliver-triggering commit.

---

## Notes on iteration

The genuinely uncertain pieces are flagged inline: the exact `rules_swift` version + bzlmod toolchain incantation (Task 2), the `rules_swift` load label in the user-story (Task 4), and whether `rules_lint` supports SwiftLint (Task 3, best-effort). Everything else is a faithful mirror of the existing `rust` wiring — read rust's actual lines in each file and substitute Swift equivalents. Don't invent rules that don't exist; verify against `rules_swift`/`rules_lint` docs.

## Deferred (not this plan)

- **Spec 2B:** `nexus-agent` into vitruvian-core — the Swift macOS `.app` via `rules_swift` + `rules_apple` (bundle/sign/macOS CI) + its Node/shell parts.
- **Rust LLVM-17 macOS backport** to this template's `MODULE.bazel` — may ride along in Task 2's edit, or stay separate.
