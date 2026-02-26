# lazy.go

> Because setting up a Go project correctly is mostly boring, repetitive nonsense that eats into your actual thinking time — and frankly, you should stop doing it by hand.

---

## What is this?

It's a TUI wizard that asks you a few pointed questions about what you're actually building, then scaffolds the right Go project structure for it.

Not the generic, copy-paste-from-Stack-Overflow structure. Not the "throw everything in main.go" structure. The *correct* structure — with proper separation of `bin`, `cmd`, `internal`, and `pkg`, a Makefile that works, a CI pipeline that isn't security theatre, and a license you actually chose for a reason.

You don't have to be smart about it. You just have to answer the questions honestly.

---

## Why does this exist?

Because I'm tired of watching people spend the first hour of a new Go project doing one of three things:

1. **Copying a half-broken template from GitHub** that was written in 2019, references deprecated tools, and has a Dockerfile with `FROM ubuntu:latest` in it like an animal.
2. **Flat-out ignoring all conventions** and then wondering why their 6000-line `main.go` is hard to maintain.
3. **Over-engineering it from day one** with a plugin system, DI framework, and 14 directories for a tool that parses CSV files.

lazy.go forces you to think for five seconds about what you're actually building — a library? An API? A security tool? — and generates a structure *appropriate for that thing*. No more, no less.

---

## Install

```bash
go install github.com/had-nu/lazy.go@latest
```

Or build it yourself, it compiles in under 3 seconds because it doesn't pull in half of npm:

```bash
git clone https://github.com/had-nu/lazy.go
cd lazy.go
make build
```

---

## Usage

### Start the wizard

```bash
lazy.go init
```

You'll get a terminal UI that asks you:

- What's this project called?
- What *type* of project is it? (CLI, API, microservice, library, security tool, worker)
- Who's it for? (internal, open source, commercial)
- How bad is it if this breaks in production?
- What do you need? (Docker, CI, linting, SAST, Dependabot...)
- What license?
- Should I create the GitHub repo and push it now?

Based on your answers, it generates a coherent, opinionated project skeleton and hands it back to you. Then you write actual code, which is the interesting part.

### Replay from config

```bash
lazy.go init --from lazygo.yml
```

The wizard exports a `lazygo.yml` when it finishes. Use it to reproduce the same structure, version-control your architectural decisions, or set up CI automation.

### Validate a config

```bash
lazy.go validate lazygo.yml
```

---

## What gets generated

### For a REST API

```
myservice/
├── cmd/server/main.go          ← graceful shutdown, signal handling
├── internal/
│   ├── handler/                ← HTTP handlers
│   ├── service/                ← business logic
│   ├── repository/             ← data access
│   ├── middleware/             ← logging, recovery
│   └── config/                 ← env-based config
├── api/openapi.yaml
├── Makefile                    ← build, test, lint targets
├── Dockerfile                  ← multi-stage, scratch final image
├── .github/workflows/ci.yml    ← actually runs tests
├── .golangci.yml
├── CONTRIBUTING.md
├── SECURITY.md
└── lazygo.yml                  ← your architectural config
```

### For a library

```
mylib/
├── pkg/mylib/
│   ├── mylib.go
│   └── mylib_test.go
└── README.md, LICENSE, go.mod
```

No `cmd/`, no `internal/`, no 11 empty directories "for future use". Just the library.

### For a security tool / production system

Automatically activates:
- `gosec` + `staticcheck` + `govulncheck` in CI
- Race detector (`go test -race`) 
- Dependabot
- `SECURITY.md` with a responsible disclosure policy
- `.golangci.yml` tuned for security-relevant linters

Because if you told the wizard it's a security tool, it believes you.

---

## The `lazygo.yml`

```yaml
project:
  name: sentinel
  module_path: github.com/user/sentinel
  type: api
  license: apache-2.0
  visibility: public
  criticality: production
features:
  docker: true
  github_actions: true
  static_analysis: true
  sast: true
  dependabot: true
  tests: true
github:
  enabled: true
  push_on_init: true
```

This file is the point. It makes your initial architectural decisions explicit and reproducible. You can check it into source control, use it in CI, or hand it to a new teammate so they understand what this project is supposed to be at a glance.

---

## What this is NOT

- It's not a code generator. It generates *structure*, not logic. Your business logic is still your problem.
- It's not magic. It won't fix a bad architecture decision you made before running it.
- It's not opinionated about your framework. Use `net/http`, use `chi`, use whatever. The scaffold is framework-agnostic.
- It's not going to hold your hand after generation. Once the files are there, you're on your own. That's the deal.

---

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md). Use Conventional Commits. Write tests. Don't open a PR with 47 changed files — split it up.

If you found a bug, open an issue. If you have a feature idea that aligns with the philosophy of "less boilerplate, more coherence", open an issue. If you want to add support for a framework-specific template, open an issue first before writing 800 lines of code nobody asked for.

---

## License

Apache-2.0. See [LICENSE](LICENSE).

---

*"The structure is the message."*
