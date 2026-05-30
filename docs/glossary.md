# Glossary — repository tiers

This platform produces and uses **three tiers** of repositories. Using these terms consistently avoids confusion about which repo a change applies to.

| Tier | Term | What it is | Example |
|---|---|---|---|
| 1 | **Template Generator** | `aspect-workflows-template` (this repo): the [Scaffold](https://hay-kot.github.io/scaffold/) meta-template — `scaffold.yaml` plus the **Generated tree** `{{ .ProjectSnake }}/` — **and** the `deliver` pipeline that produces the Starter Templates. | `aspect-workflows-template` |
| 2 | **Starter Template** (*Starter*) | A per-preset repo **we** generate from the Template Generator via the `deliver` pipeline (`scaffold new --preset=…`) and publish for developers to start from. One per preset. | `kitchen-sink`, `go`, `py`, `backstage-swift`, … |
| 3 | **Project Repo** (*Developer Repo*) | The actual project a **developer** creates from a Starter Template, where they do their real work. | `vitruvian-core` (created from the `kitchen-sink` Starter) |

**Helper term — Generated tree:** the `{{ .ProjectSnake }}/` directory inside the Template Generator. It's the Scaffold-templated content that becomes a Starter (and flows on into Project Repos).

## Flow

```
Template Generator  ──(we run deliver / scaffold)──▶  Starter Templates  ──(a developer creates from one)──▶  Project Repos
   (this repo)                                          (kitchen-sink, …)                                       (vitruvian-core, …)
```

## Important: no sync between tiers

- The Generated tree's content reaches Starters automatically via the `deliver` pipeline.
- A **Project Repo is a one-time copy** of a Starter at creation time — **Starters and Project Repos do not sync.** A change added to the Template Generator reaches *new* Project Repos (created afterward) but **not existing** ones, which must be updated by a **manual port**.

> Note: this is distinct from the Copybara bidirectional sync, which keeps a *monorepo's component subtrees* in step with their *standalone component repos* — it does **not** connect Starters to Project Repos.
