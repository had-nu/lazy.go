# Contributing to lazy.go

lazy.go is opinionated by design — contributions are held to the same standard.

---

## Before You Open a PR

Open an issue first. Explain what you want to add and why it fits the project. If it makes lazy.go smarter about structure and intent, great. If it adds framework-specific boilerplate or makes the wizard longer without making it more useful, the answer is no.

This is not a democracy. But it is a discussion.

---

## What We Want

- Bug fixes with a clear reproduction case
- New project types that represent meaningfully distinct architectures
- Better default templates (tighter, more idiomatic Go stubs)
- More precise security policy logic
- Documentation that doesn't lie

## What We Don't Want

- Framework adapters (cobra, gin, echo presets — users do that)
- GUI or web wrappers
- Feature flags to disable existing behavior
- PRs with 40 changed files and no discussion

---

## Setup

```bash
git clone https://github.com/had-nu/lazy.go
cd lazy.go
go mod download
make test
```

Requirements: Go 1.22+, `git`, optionally `gh` CLI for GitHub integration tests.

---

## Commit Style

Use [Conventional Commits](https://www.conventionalcommits.org/).

| Prefix | When |
|---|---|
| `feat:` | New capability |
| `fix:` | Bug fix |
| `docs:` | Documentation |
| `refactor:` | No behaviour change |
| `test:` | Tests only |
| `chore:` | Maintenance, deps |
| `ci:` | Pipeline changes |

**One thing per commit. One thing per PR.**

---

## Tests

New code must have tests. Existing tests must not regress. Run:

```bash
make test           # go test -race ./...
make lint           # golangci-lint run
```

Coverage targets by package: `pkg/config` ≥ 80%, `pkg/scaffold` ≥ 70%, `pkg/wizard` ≥ 60%.

---

## Review

All PRs are reviewed before merge. Expect feedback. Address it or argue against it — but don't ignore it.

Merges to `main` go through CI. If CI fails, the PR doesn't merge. No exceptions.
