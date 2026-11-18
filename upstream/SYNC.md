# Weekly upstream sync

The steps an agent follows to bring everything
[aspect-build/aspect-workflows-template](https://github.com/aspect-build/aspect-workflows-template)
published since the last sync into this fork by merging it, and land it on
`platform-v2.0`. A
scheduled Claude Code routine runs it every Monday morning, Pacific time, in a
fresh cloud session that reads this file from `platform-v2.0`. Anyone can run it
by hand the same way.

The routine lives in James's Claude account, not in this repository: it holds
only the schedule and a pointer to this file. To change how the sync works,
change this file in a pull request; the next run follows the version on
`platform-v2.0`. Read the repository's [`AGENTS.md`](../AGENTS.md) first.

## Goal

`upstream/main` is an ancestor of `platform-v2.0`: everything upstream
published is either merged here or deliberately reverted in the merge, and the
merge is on `platform-v2.0` with its checks green and delivered to the starter
repositories.

## How this fork relates to upstream

- **Merged, not rebased.** This fork tracks upstream by merging `upstream/main`
  into `platform-v2.0` with a real merge commit. Git applies whatever merges
  cleanly; you resolve only the conflicts. Never rebase `platform-v2.0`: it is
  the default branch, the release tags and the 26 starters' deliveries hang off
  its history, and a rebase means force-pushing all of it.
- **History.** From 2026-06-26 (#38, the AXL engine migration copied in rather
  than merged) until the reconnect on 2026-11-17, the two histories were
  disconnected. The reconnect merge (`git merge -s ours`) recorded upstream as
  merged after every upstream commit up to `f595cf0` had been ported or judged
  already present. Its commit message lists each decision. From there on,
  `git merge upstream/main` only offers what is new.
- **`platform-v2.0` is the base branch** (it is the default branch). Ignore this
  fork's `main`.
- **Paths mostly line up.** Upstream's `template/` and ours share a layout, so
  most `template/` changes merge on their own. The root `.github/workflows/`
  does not: ours is rewired around `.github/scripts/normalize-rendered.sh` and
  `deliver.yaml`.

## Where the work happens

- A cloud session's git access is scoped to the branch it was assigned. Make
  your branch by extending that name: `<assigned branch>-<YYYYMMDD>`
  (`git branch --show-current` gives the assigned one). This file is your
  permission to push that branch and open the pull request from it. Do not stop
  to ask. Outside a session use `chore/upstream-sync-<YYYYMMDD>`.
- Plain `git` and `gh` work in this repository. Always pass
  `--repo VitruvianSoftware/aspect-workflows-template` to `gh`, because the
  `upstream` remote makes it guess wrong.
- Fetch **full** history (`git fetch --unshallow origin` if the clone is
  shallow). A shallow clone cannot find the merge base and makes git think
  every upstream commit is new.

## Steps

1. Fetch: `git fetch origin platform-v2.0`, then add and fetch upstream
   (`git remote add upstream https://github.com/aspect-build/aspect-workflows-template.git`
   if missing, then `git fetch upstream main`). Start your branch from
   `origin/platform-v2.0`. If an open pull request already carries an upstream
   sync (title starts `chore(sync): merge upstream` or `fix(sync): merge upstream`),
   drive that one to merged first.
2. List what is new: `git log --reverse --format='%h %ad %s' --date=short
   origin/platform-v2.0..upstream/main`. If it is empty, stop and report
   "caught up at <upstream/main sha>". If it lists commits from June 2026 or
   earlier, the history is not connected (shallow clone, or the reconnect was
   undone): stop and report, do not merge.
3. Merge: `git merge --no-ff --no-commit upstream/main`. Then review every
   change, conflicted or not, with `git diff --cached origin/platform-v2.0`.
   A clean merge can still bring in something we do not want:
   - **Aspect CLI or `setup-aspect` bumps**: take upstream's
     `template/.aspect/version.axl`. For each `aspect-build/setup-aspect@<sha> # vX`
     use (root `.github/workflows/*` and `template/.github/workflows/ci.yaml`),
     move ours to upstream's pin, and confirm the tag resolves to that SHA with
     `gh api repos/aspect-build/setup-aspect/git/ref/tags/<tag>`. Leave
     `aspect-launcher` in `template/tools/tools.lock.json` alone unless
     upstream moved it.
   - **Bazel version** (`template/.bazelversion`): take the newest **stable**
     release in the same major (no `rc`); it may be newer than upstream's pin.
     Read its release notes for removals and flag flips and say what you found
     in the PR. vitruvian-core is a separate repository and not this routine's
     job; if the bump looks relevant there, say so in the report.
   - **`template/` changes**: keep upstream's change and our jinja
     (`{% if ... %}`) and additions (licence tooling, Copybara, Pulumi, RBE,
     backstage presets, the app initializer).
   - **Never take**: upstream's release, tag or publish workflows (we use
     release-please and `deliver.yaml`; delete any such file the merge adds),
     their branding or org names, or anything that reintroduces the
     hay-kot/scaffold engine. Resolve those hunks to ours.
   - When you are not sure a change is wanted, keep ours, and list it for James
     with the upstream SHA. Do not guess.
4. Commit the merge with `git commit` (keep it a merge commit). Subject:
   `fix(sync): merge upstream <old sha>..<new sha>` when anything reaches the
   starters (it triggers a release), else `chore(sync): ...`. The body lists
   each upstream commit and what happened to it: merged, adapted (how), or
   resolved to ours (why).
5. Open one pull request to `platform-v2.0` with that title. Body: first line
   `> **Merge with "Create a merge commit", not squash or rebase.**` (a human
   may merge it), then the same table, then the verification. Put a release-please commit override in the
   body (see "What earlier syncs learned").
6. Wait for every check. The 18 Linux preset legs render each preset the way
   delivery does and run build, test, lint, format and gazelle on it; they are
   the authority. macOS legs only run after merge. Read a failing job's log and
   fix the cause with a further commit on the branch; never skip or disable a
   check.
7. Merge it yourself with **`gh pr merge <n> --merge`** once **every** check is
   green. Never `--squash` or `--rebase` a sync: either one drops upstream from
   the history and the next sync is back to hundreds of commits "behind". The
   repository has no required checks and no auto-merge, so you are the gate. If
   a sync cannot be made green after a few honest attempts, leave the PR open
   with one comment naming the blocker, and report that.
8. After merging, verify on the merge commit (pass the **full** SHA to
   `gh run list --commit`; a short one matches nothing):
   - `git merge-base --is-ancestor upstream/main origin/platform-v2.0` succeeds;
   - `CI` is green, including the 17 macOS legs;
   - `Deliver starters` finishes 26 of 26;
   - every starter's `CI` on its new `main` head is green
     (`gh run list --repo VitruvianSoftware/<preset> --commit <sha> --workflow CI`
     for each preset in `deliver.yaml`).
   If release-please opens or updates a release PR, check its version bump and
   that its notes list only our own PRs, then report it. **Do not merge it.**
9. Report: the PR merged, the upstream SHA the fork is caught up with, every
   upstream change resolved to ours and why, and anything James must decide.

## What earlier syncs learned

- **Runners are pinned to `ubuntu-24.04`, never `ubuntu-latest`.** When
  `ubuntu-latest` became Ubuntu 26.04 (from 2026-10-19), it dropped
  `libxml2.so.2`, which the hermetic LLVM 19.1.7 `ld.lld` and the Swift
  toolchain both load, so every preset that links native code (cpp, rust,
  swift, kitchen-sink) failed. If an upstream merge brings `ubuntu-latest`
  back, resolve it to our pin. Move the pin deliberately, after the toolchains
  run on the new image.
- **release-please writes release notes from our merged PRs, not commits**
  (`"changelog-type": "github"` in `release-please-config.json`). With the
  default, commit-based notes, the reconnect (#57) produced a 3.0.3 changelog
  with each of our lines nine times (the override is applied once per commit
  of a merge-committed PR) and upstream's cosmo.zip fix linked to a
  non-existent `#912` here (GitHub ties upstream commits to upstream's PRs).
  GitHub's own notes list only PRs merged into this repository, once each, so
  upstream cannot leak into them. Do not switch it back.

  The version **bump** is still computed from every commit, upstream's
  included, so an upstream `feat:` could make a sync a minor release. Keep a
  commit override at the end of every sync PR body to steer it:

  ```
  BEGIN_COMMIT_OVERRIDE
  fix(sync): merge upstream <old>..<new> (<one-line summary>)
  END_COMMIT_OVERRIDE
  ```

  Use `chore(sync): ...` inside it when nothing reaches the starters. In step 8,
  check the release PR's version is the bump you expect (a patch for `fix`,
  nothing for `chore`) and that its notes list only our PRs. Report anything
  else to James; never merge a release PR.
- **Delivery force-pushes all 26 starters on every push to `platform-v2.0`**,
  even a docs-only change. It is serialised per preset (`concurrency`, #56), so
  two close merges are safe, but each one costs about 40 minutes of runners.
  Land a sync as one PR.
- **Template CI must render what ships.** It once rendered without licence
  headers while delivery added them, so four Ruby starters were red for seven
  weeks while the template stayed green (#52). Both now call
  `.github/scripts/normalize-rendered.sh`; change it, not the two workflows.
- **ttl.sh**, the demo registry the sample image pushes to, sometimes stalls.
  The `Delivery` step and job are capped at 20 minutes and skip with a warning
  when ttl.sh is down (#54). A `Delivery` hang or skip is not your change.
- **A starter can be red while the template is green**, so step 8 checks the
  starters themselves, not just this repository.
- **Shell pitfalls seen here:** zsh does not word-split an unquoted
  `$LIST` in a `for` loop (use an array); `grep` may be a shell function, so
  call `command grep` in pipelines; `git` prints harmless
  `is not a valid attribute name: template/.gitattributes` warnings because
  that file is a jinja template (`export GIT_ATTR_NOSYSTEM=1` does not silence
  them; ignore them).
