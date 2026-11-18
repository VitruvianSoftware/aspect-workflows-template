# AGENTS.md

Guidance for AI coding agents working on the **Aspect Workflows starter template**
(useful for humans too). This repo is a **scaffolding template** (hay-kot/scaffold),
not a buildable app: it *generates* hermetic Bazel monorepos.

## Layout
- `{{ .ProjectSnake }}/` — the **delivered tree**: everything here is rendered into a
  generated repo. Files use Go-template syntax (`{{ ... }}`).
- `scaffold.yaml` — questions, the `computed` language booleans (`.Computed.python`, …),
  and `features:` globs that gate which files ship for a given selection.
- `user_stories/` — per-preset **executable** smoke tests (run in CI; they also become the
  starter READMEs). See `user_stories/STRUCTURE.md`.
- `docs/` — template *product* docs (admin/user/contributor guides). **Not delivered.**

## Editing the delivered tree — landmines
- **Per-language** content is gated `{{ if .Computed.<lang> }} … {{ end }}`; **per-feature**
  with `{{ if .Scaffold.<feature> }}` (e.g. `copybara`, `lint`, `license`, `oci`).
- `{{` / `}}` are the template delimiters. **Never put a literal `{{` in a delivered file
  unless escaped** — e.g. a Mermaid hexagon `{{…}}` gets eaten by the renderer; use a
  different node shape. Keep delimiters balanced.
- A new per-language file needs a matching `features:` glob in `scaffold.yaml` gated on
  `.Computed.<lang>` (see `docs/contributor-guide/adding-languages.md`).
- **Every render must be a fixed point of `bazel run //:tidy`** (gazelle + format): the
  starters run a Tidy Check on every push, so a render that gazelle/buildifier would
  rewrite fails on delivery. Keep BUILD/MODULE templates buildifier-canonical for every
  preset: loads at the top (sorted), rule attributes sorted (`name` first, `deps`/
  `visibility` last), trailing commas, single blank lines between paragraphs, no trailing
  blank lines. Control whitespace with `{{- if }}`-style trims — the CI "Tidy made
  no changes" step renders every preset and fails on any drift.

## Testing
- CI renders **every preset** (the preset matrix) and runs the `user_stories`.
- **You cannot fully test a starter in isolation** — many bugs only surface once rendered.
  The convention: render via CI, then **port the change into `vitruvian-core`** and build it
  there. That repo is the real test bed for template changes (it was generated from the
  `kitchen-sink` starter and is maintained by hand).

## Conventions
- **Conventional commits** (`feat:`, `fix:`, `chore:`, `docs:` …).
- **How to merge a PR.** Upstream syncs (title starts `fix(sync): merge upstream`
  or `chore(sync): merge upstream`) use **Create a merge commit**
  (`gh pr merge --merge`). Squash or rebase would drop upstream's commits from our
  history, and the next sync would see all of them as new again. Every other PR
  carries only our own commits: **Squash and merge** as usual.
- Generated repos themselves ship an `AGENTS.md` (`{{ .ProjectSnake }}/AGENTS.md`) — keep the
  two in mind: this file is for working *on the template*; that one is for the *generated repo*.
