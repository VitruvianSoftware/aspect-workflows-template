#!/usr/bin/env bash
# Post-render normalisation of a freshly rendered starter.
#
# Run from INSIDE the rendered directory. This is the single definition of
# "what a delivered starter ships": .github/workflows/deliver.yaml publishes its
# output and .github/workflows/ci.yaml checks it, so the two cannot drift apart.
# They did drift: CI rendered without licence headers, so a defect the headers
# introduce shipped to the Ruby starters while the template's CI stayed green.
#
#   usage: normalize-rendered.sh <preset> [<task-suffix>]
#
#   NORMALIZE_GIT_NAME / NORMALIZE_GIT_EMAIL  identity for the rendered commit.
#   <task-suffix>  when given, tasks run as `<task>-<suffix>` (CI uses this to
#                  keep per-leg task names unique); delivery omits it.
set -euo pipefail

preset="${1:?usage: normalize-rendered.sh <preset> [<task-suffix>]}"
suffix="${2:-}"
name="${NORMALIZE_GIT_NAME:-CI}"
email="${NORMALIZE_GIT_EMAIL:-ci@aspect.build}"
git_id=(-c "user.name=${name}" -c "user.email=${email}")

# `aspect <task> --task:name <task>-<suffix>` only when a suffix was supplied.
task_name() { if [ -n "$suffix" ]; then printf -- '--task:name %s-%s' "$1" "$suffix"; fi; }
quiet=(--github-status-comments:enabled=false --github-status-checks:enabled=false)

# repin: the template ships placeholder lockfiles; regenerate them.
if [ -x ./tools/repin ]; then ./tools/repin; fi
bazel mod deps --lockfile_mode=update || true

# git init FIRST: gazelle/format need a git tree to enumerate and apply to.
git init -q
git add -A
git "${git_id[@]}" commit -q -m template

# Normalise: buildifier, then gazelle.
bazel run --run_in_cwd @buildifier_prebuilt//buildifier -- -r . || true
# shellcheck disable=SC2046  # task_name is intentionally word-split
aspect gazelle $(task_name gazelle) "${quiet[@]}" --check-only=false || true

# Licence headers go on BEFORE the formatter runs. addlicense inserts the header
# plus a trailing blank line, but a file that already began with a blank line
# after its shebang leaves two in a row; the formatter must see (and fix) the
# headered tree, because that is what the starter's own CI formats.
bazel run //tools/license:add || true

case "${preset}" in
  # no wired source formatter (minimal has no languages; scala/ruby have none)
  minimal | scala | ruby | backstage-minimal | backstage-scala | backstage-ruby) ;;
  *)
    # shellcheck disable=SC2046
    aspect format $(task_name format) "${quiet[@]}" --scope=all || true
    ;;
esac

# Ruby has no source formatter, and addlicense puts its header straight under a
# `# frozen_string_literal` magic comment (and leaves a doubled blank line).
# Autocorrect exactly those two layout cops and nothing else.
if [ -f Gemfile ] && grep -q rubocop Gemfile; then
  rb_files="$(git ls-files '*.rb')"
  if [ -n "${rb_files}" ]; then
    printf '%s\n' "${rb_files}" | xargs bazel run --run_in_cwd @bundle//bin:rubocop -- \
      -a --only Layout/EmptyLineAfterMagicComment,Layout/EmptyLines \
      || echo "::warning::rubocop layout autocorrect failed; Ruby lint may flag licence-header layout"
  fi
fi

git add -A
git "${git_id[@]}" commit -q --amend --no-edit || true
