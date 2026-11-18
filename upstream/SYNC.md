# Weekly upstream sync

The steps an agent follows to bring everything
[aspect-build/aspect-workflows-template](https://github.com/aspect-build/aspect-workflows-template)
published since the last sync into this fork and land it on `platform-v2.0`. A
scheduled Claude Code routine runs it every Monday morning, Pacific time, in a
fresh cloud session that reads this file from `platform-v2.0`. Anyone can run it
by hand the same way.

The routine lives in James's Claude account, not in this repository: it holds
only the schedule and a pointer to this file. To change how the sync works,
change this file in a pull request; the next run follows the version on
`platform-v2.0`. Read the repository's [`AGENTS.md`](../AGENTS.md) first.

## Goal

Every upstream commit has a row in [`ledger.tsv`](ledger.tsv), none is
`pending`, and every port is merged to `platform-v2.0` with its checks green and
delivered to the starter repositories.

## How this fork relates to upstream

- **No shared history.** On 2026-06-26 (#38) we copied upstream's new template
  engine (`render.axl`, `template/`, `template-config.json`) in rather than
  merging it, so `git merge-base` with upstream is a 2026-05-23 commit and git
  reports dozens of commits "behind" that are already here. Never rebase or
  merge `upstream/main`. The ledger is the record of what has been taken.
- **`platform-v2.0` is the base branch** (it is the default branch). Ignore this
  fork's `main`.
- **Paths mostly line up.** Upstream's `template/` and ours have the same layout,
  so `git show <sha> -- template/ | git apply` often works. The root
  `.github/workflows/` does not line up: ours is rewired around
  `.github/scripts/normalize-rendered.sh`, so upstream CI changes are ported by
  hand or skipped.

## Ledger

`ledger.tsv` has one row per upstream commit after the 2026-05-23 merge-base,
oldest first. Columns: `upstream_sha`, `date`, `status`, `ref`, `note`,
`subject`. Status is one of:

| status | meaning | ref |
|---|---|---|
| `pending` | not yet decided | `-` |
| `ported` | change made in this fork | the pull request, `#NN` |
| `present` | already here (verified, say how in the note) | the PR that brought it |
| `superseded` | a later change makes it moot, e.g. a newer version pin | that PR |
| `skipped` | deliberately not taken; the note says why | `-` or a related PR |

`present` needs evidence: the patch reverse-applies
(`git show <sha> | git apply --check -R`), or you read both files and say so in
the note. Never mark a row `present` on a hunch.

## Where the work happens

- A cloud session's git access is scoped to the branch it was assigned. Make
  your branch by extending that name: `<assigned branch>-<YYYYMMDD>`
  (`git branch --show-current` gives the assigned one). This file is your
  permission to push that branch and open the pull request from it. Do not stop
  to ask. Outside a session use `chore/upstream-sync-<YYYYMMDD>`.
- Plain `git` and `gh` work in this repository. Always pass
  `--repo VitruvianSoftware/aspect-workflows-template` to `gh`, because the
  `upstream` remote makes it guess wrong.

## Steps

1. Fetch: `git fetch origin platform-v2.0`, and add and fetch upstream
   (`git remote add upstream https://github.com/aspect-build/aspect-workflows-template.git`
   if missing, then `git fetch upstream main`). Start your branch from
   `origin/platform-v2.0`. If an open pull request already carries an upstream
   sync (title starts `chore(sync): port upstream`), drive that one to merged
   first.
2. Add a `pending` row for every commit in `git log --reverse upstream/main`
   that is newer than the last ledger row and not in it. If nothing is pending,
   stop and report "caught up at <upstream sha>".
3. Decide each pending row, oldest first:
   - **Aspect CLI or `setup-aspect` bumps** collapse into one change: set
     `template/.aspect/version.axl` to upstream's latest CLI, and every
     `aspect-build/setup-aspect@<sha> # vX` use (root `.github/workflows/*` and
     `template/.github/workflows/ci.yaml`) to upstream's latest. Confirm the
     tag resolves to that SHA with
     `gh api repos/aspect-build/setup-aspect/git/ref/tags/<tag>`. Earlier bump
     rows become `superseded` by this PR. Leave `aspect-launcher` in
     `template/tools/tools.lock.json` alone unless upstream moved it.
   - **Bazel version** (`template/.bazelversion`): take it only if it is a
     stable release (no `rc`), and read its release notes for removals and
     flag flips. Say in the PR what you found. vitruvian-core is a separate
     repository and not this routine's job; if the bump looks relevant there,
     say so in the report.
   - **`template/` changes**: apply the upstream diff, then resolve it against
     our jinja (`{% if ... %}` blocks). Keep our additions (licence tooling,
     Copybara, Pulumi, RBE, backstage presets, the app initializer).
   - **Never take**: upstream's release, tag or publish workflows (we use
     release-please and `deliver.yaml`), their branding or org names, or
     anything that reintroduces the hay-kot/scaffold engine.
   - When you are not sure a change is wanted, mark it `skipped` with the
     reason and list it for James. Do not guess.
4. Commit in conventional format. Use `fix:` or `feat:` only when the change
   reaches the starters (it triggers a release); use `ci:` or `chore:` for
   changes to this repository's own workflows or docs. Each commit body names
   the upstream SHAs it covers.
5. Open one pull request to `platform-v2.0` titled
   `chore(sync): port upstream <first sha>..<last sha>` (use `fix(sync):` if
   anything reaches the starters). Body: a table of each upstream commit, its
   status and what changed here, then the verification. Update the ledger rows
   with the PR number in the same PR.
6. Wait for every check. The 18 Linux preset legs render each preset the way
   delivery does and run build, test, lint, format and gazelle on it; they are
   the authority. macOS legs only run after merge. Read a failing job's log and
   fix the cause; never skip or disable a check.
7. Merge it yourself with `gh pr merge <n> --squash` once **every** check is
   green. The repository has no required checks and no auto-merge, so nothing
   else stops a red merge: you are the gate. If a port cannot be made green
   after a few honest attempts, leave the PR open with one comment naming the
   blocker, and report that.
8. After merging, verify on the merge commit (pass the **full** SHA to
   `gh run list --commit`; a short one matches nothing):
   - `CI` is green, including the 17 macOS legs;
   - `Deliver starters` finishes 26 of 26;
   - every starter's `CI` on its new `main` head is green
     (`gh run list --repo VitruvianSoftware/<preset> --commit <sha> --workflow CI`
     for each preset in `deliver.yaml`).
   If release-please opens a release PR, **do not merge it**. Report it.
9. Report: the PR merged, the upstream SHA the fork is caught up with, every row
   still `pending` or newly `skipped` and why, and anything James must decide.

## What earlier syncs learned

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
