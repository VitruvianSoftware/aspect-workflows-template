# Template Swift Support (Spec 2A)

**Status:** Draft for review · 2026-05-24

**Initiative:** Add Swift to the Vitruvian platform — to the `aspect-workflows-template` (this spec, 2A) and to `vitruvian-core`'s `nexus-agent` (Spec 2B). 2A lays the `rules_swift` foundation 2B builds on.

**This spec (2A):** Add **Swift as a first-class, cross-platform language** to `aspect-workflows-template` (branch `platform-v2.0`), following the established add-a-language pattern. Outcome: the template can scaffold Swift projects, and `swift` + `backstage-swift` starter repos exist.

## Goal & success criteria

1. `scaffold new --preset=swift` (and combinations, e.g. `--set langs="Swift Go"`) produces a project whose Swift sample builds + tests under Bazel: `bazel build //swift/... && bazel test //swift/...` green **on Linux** (rules_swift supports Linux; the ubuntu CI must stay green).
2. `swift` and `backstage-swift` are real presets, wired through `scaffold.yaml`, the CI/deliver matrices, the Pulumi `starter_repos` presets, `catalog-info.yaml`, and `user_stories/`.
3. The two new starter repos (`swift`, `backstage-swift`) exist (created by `pulumi up`) and the deliver pipeline pushes to them.

## Out of scope (2A)

- `rules_apple` / macOS `.app` bundling, code signing, macOS CI — that's 2B (nexus-agent in vitruvian-core).
- The `nexus-agent` migration itself — 2B.
- Copybara — Spec 3.
- The Rust LLVM-17 macOS backport — adjacent; tracked separately. It may ride along in the same `MODULE.bazel` edit at the implementer's discretion, but is not required by this spec.

## Decisions (from brainstorming)

- **Generic cross-platform Swift** via `rules_swift` only. No `rules_apple` in the template.
- Sample is a small **CLI/library** that builds on **ubuntu CI** (no macOS runner needed for 2A).
- Swift has **no Gazelle plugin** → the sample's `BUILD` is **hand-written** (no `ENABLE_LANGUAGES` change for swift).
- Follow the contributor guide's checklist, but use the **actual** current `scaffold.yaml` structure (multi-select human labels + `has`), not the guide's stale `split` examples.

## Wiring (files to touch in `aspect-workflows-template`)

1. **`scaffold.yaml`:** add `Swift` to the `langs` question options; `computed.swift: "{{ has \"Swift\" .Scaffold.langs }}"`; a `features` entry (`value: "{{ .Computed.swift }}"`) with globs for the swift sample (`*/swift/**`) + any swift config files; add `swift` and `backstage-swift` presets (mirroring `go`/`backstage-go`).
2. **`{{ .ProjectSnake }}/MODULE.bazel`:** `rules_swift` (current BCR version) gated on `{{- if .Computed.swift }}` — `bazel_dep` + the swift toolchain extension + `register_toolchains`. Verify the exact current `rules_swift` module name/version on BCR during implementation.
3. **`{{ .ProjectSnake }}/swift/hello/`:** hand-written `BUILD` with `swift_library` + `swift_binary` + `swift_test`, plus `hello.swift` and a test file — a "Hello, World" CLI. Feature-gated.
4. **Format/lint:** wire `swift-format` (Apple) into the format tool (`{{ .ProjectSnake }}/tools/format/...`); add SwiftLint via `rules_lint` **if** rules_lint supports Swift — otherwise document it as a best-effort follow-up. Do not block the spec on Swift linting.
5. **`{{ .ProjectSnake }}/README.bazel.md`:** a `{{- if .Computed.swift }}` Swift dev section.
6. **CI/deliver matrices:** add `swift` to `.github/workflows/ci.yaml`'s preset matrix and `swift` + `backstage-swift` to `.github/workflows/deliver.yaml`'s matrix (both on `ubuntu-latest`).
7. **`user_stories/swift.md` + `user_stories/backstage-swift.md`:** required — `deliver.yaml` copies `user_stories/<preset>.md` into each starter repo's README, so a missing file fails delivery.
8. **Pulumi + catalog:** add `swift` + `backstage-swift` to `infrastructure/pulumi/pkg/github_repos/starter_repos.go`'s `presets` list, and to the language branches in `{{ .ProjectSnake }}/catalog-info.yaml` (and `template.yaml` if it enumerates languages).

## Rollout (careful, infra-touching, last)

Adding the two presets to the Pulumi `presets` list means `pulumi up` will create two new starter repos (`swift`, `backstage-swift`) with ED25519 deploy keys + `STARTER_DEPLOY_*` Actions secrets. Do this with the established discipline: `pulumi preview` → confirm it's only the two new repos (+ their keys/secrets), no destructive changes → `pulumi up` → verify the deliver pipeline pushes to both.

## Risks & unknowns

- **`rules_swift` on Linux:** supported, but verify the current BCR version builds the sample on ubuntu cleanly (toolchain download) before wiring CI.
- **`rules_lint` Swift support is uncertain** — SwiftLint integration may be best-effort/deferred; `swift-format` should be straightforward.
- **Dependency interaction:** `rules_swift` is a new `bazel_dep`; confirm it doesn't conflict with existing template deps (different ecosystem, so unlikely).
- **Infra rollout** creates live repos — gated behind `pulumi preview` + explicit confirmation.

## Verification

- Local: `scaffold new --preset=swift --output-dir=/tmp/sw .`, then `cd /tmp/sw && bazel build //swift/... && bazel test //swift/...` green; repeat with a combo (`langs="Swift Go"`).
- Template CI: the `swift` leg of the ci/deliver matrix is green on ubuntu.
- Infra: after `pulumi up`, `swift` + `backstage-swift` repos exist; a deliver run pushes generated content to both.
