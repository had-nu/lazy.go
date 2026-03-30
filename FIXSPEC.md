# FIXSPEC — lazy.go security review

> Source: `/go-security` review · 30 Mar 2026  
> Scope: `main.go`, `pkg/github/auth.go`, `pkg/scaffold/generator.go`, `pkg/wizard/flow.go`, `pkg/security/policies.go`

---

## FIX-001 — Command injection via shell in `lookupEnv`

**Severity:** CRITICAL  
**File:** `pkg/github/auth.go`  
**Lines:** `lookupEnv` var declaration (approx. L22–L30)

### What is wrong

`lookupEnv` spawns a shell to read an environment variable:

```go
var lookupEnv = func(key string) string {
    cmd := exec.Command("sh", "-c", fmt.Sprintf("echo ${%s}", key))
    out, err := cmd.Output()
    if err != nil {
        return ""
    }
    return strings.TrimSpace(string(out))
}
```

`key` is caller-controlled. If a caller passes a value containing shell metacharacters — `$(...)`, `` ` `` , `;` — the shell executes arbitrary commands as the current user. The comment claims this exists for testability, but testability in Go is achieved by making the variable substitutable, not by invoking a subprocess.

`os.LookupEnv` is a direct syscall wrapper. No shell, no process, no attack surface.

### Fix

```go
// lookupEnv reads a single environment variable.
// Extracted as a variable so tests can substitute it without spawning a shell.
var lookupEnv = func(key string) string {
    v, _ := os.LookupEnv(key)
    return v
}
```

Remove the `strings` import from this file if it is no longer used after this change. Verify with `go vet ./pkg/github/...`.

### Test impact

No test changes required. The variable remains substitutable in tests:

```go
// in test
oldLookup := auth.LookupEnv  // export the var if needed, or test via TokenFromEnv
auth.LookupEnv = func(key string) string { return "fake-token" }
defer func() { auth.LookupEnv = oldLookup }()
```

### Verification

```bash
go vet ./pkg/github/...
gosec ./pkg/github/...        # G204 should no longer fire on this path
grep -n "exec.Command" pkg/github/auth.go  # must return zero lines
```

---

## FIX-002 — Duplicated security enforcement logic

**Severity:** HIGH  
**Files:** `pkg/wizard/flow.go` (L48–L56) and `pkg/security/policies.go` (`EnforceSecurity`)

### What is wrong

Security enforcement is written twice with identical logic:

```go
// pkg/wizard/flow.go — BuildConfig()
if cfg.IsSecure() {
    cfg.Features.StaticAnalysis = true
    cfg.Features.SAST = true
    cfg.Features.Tests = true
    if cfg.Features.GitHubActions {
        cfg.Features.Dependabot = true
    }
}
```

```go
// pkg/security/policies.go — EnforceSecurity()
func EnforceSecurity(cfg *config.ProjectConfig) {
    if !ShouldEnableSecurity(cfg) {
        return
    }
    cfg.Features.StaticAnalysis = true
    cfg.Features.SAST = true
    cfg.Features.Tests = true
    if cfg.Features.GitHubActions {
        cfg.Features.Dependabot = true
    }
}
```

Two copies of a security rule will diverge. When a new mandatory control is added to `EnforceSecurity` — say, `Features.Linting = true` for production projects — the copy in `BuildConfig` will silently miss it. The bug surfaces only when a user replays a config via `--from` (which calls `BuildConfig` but not the TUI path that reaches `main.go`'s `security.EnforceSecurity`).

### Fix

Remove the duplicated block from `BuildConfig`. `EnforceSecurity` is already called in `runGeneration` in `main.go`, but `BuildConfig` is the correct place to enforce it so that configs constructed programmatically (tests, `--from` replay) are also covered. Call it once, at the end of `BuildConfig`:

```go
// pkg/wizard/flow.go — BuildConfig()
func BuildConfig(state WizardState) *config.ProjectConfig {
    cfg := &config.ProjectConfig{
        Name:        state.ProjectName,
        ModulePath:  state.ModulePath,
        Description: state.Description,
        Author:      state.Author,
        Type:        config.ProjectType(state.ProjectType),
        Visibility:  config.Visibility(state.Visibility),
        Criticality: config.CriticalityLevel(state.Criticality),
        License:     config.LicenseType(state.License),
        Features: config.Features{
            Docker:         state.Features["docker"],
            GitHubActions:  state.Features["github_actions"],
            Linting:        state.Features["linting"],
            StaticAnalysis: state.Features["static_analysis"],
            Dependabot:     state.Features["dependabot"],
            Tests:          state.Features["tests"],
            SAST:           state.Features["sast"],
        },
        GitHub: config.GitHubConfig{
            Enabled:    state.GitHubEnable,
            PushOnInit: state.GitHubPush,
        },
    }

    if cfg.License == "" || cfg.License == "auto" {
        cfg.License = SuggestLicense(cfg)
    }

    // Single source of truth for security enforcement.
    // EnforceSecurity is a no-op for experimental projects.
    security.EnforceSecurity(cfg)

    return cfg
}
```

This introduces an import of `pkg/security` into `pkg/wizard`. Verify there is no circular import: `security` imports only `config`, `wizard` imports `config` and now `security` — no cycle.

Remove the call to `security.EnforceSecurity(cfg)` from `main.go`'s `runGeneration` — it is now redundant.

### Verification

```bash
go build ./...                # catches import cycles immediately
go test ./pkg/wizard/...
go test ./pkg/security/...
```

Add a regression test in `pkg/wizard/flow_test.go`:

```go
func TestBuildConfig_EnforcesSecurityForProduction(t *testing.T) {
    state := wizard.NewWizardState()
    state.ProjectName = "svc"
    state.ModulePath = "github.com/x/svc"
    state.ProjectType = string(config.ProjectTypeAPI)
    state.Criticality = string(config.CriticalityProduction)
    state.Features = map[string]bool{"github_actions": true}

    cfg := wizard.BuildConfig(state)

    if !cfg.Features.StaticAnalysis {
        t.Error("production project must have StaticAnalysis")
    }
    if !cfg.Features.SAST {
        t.Error("production project must have SAST")
    }
    if !cfg.Features.Tests {
        t.Error("production project must have Tests")
    }
    if !cfg.Features.Dependabot {
        t.Error("production project with GH Actions must have Dependabot")
    }
}
```

---

## FIX-003 — Dead code: deprecated `strings.Title` in funcMap

**Severity:** MEDIUM  
**File:** `pkg/scaffold/generator.go`  
**Lines:** `funcMap` declaration

### What is wrong

```go
var funcMap = template.FuncMap{
    "upper":   strings.ToUpper,
    "lower":   strings.ToLower,
    "title":   strings.Title, //nolint:staticcheck
    "replace": strings.ReplaceAll,
    "join":    strings.Join,
}
```

`strings.Title` is deprecated since Go 1.18. The `//nolint:staticcheck` suppression is a flag that dead code is being carried. A grep across all `.tmpl` files confirms `title` is never called from any template — the entry is unreachable.

Dead code with a linter suppression is worse than dead code alone: it trains the reader to ignore suppression comments, which are the last line of defence when `gosec` or `staticcheck` catch a real issue.

### Fix

Remove the entry entirely:

```go
var funcMap = template.FuncMap{
    "upper":   strings.ToUpper,
    "lower":   strings.ToLower,
    "replace": strings.ReplaceAll,
    "join":    strings.Join,
}
```

If `title` is needed in a future template, the correct replacement is:

```go
import "golang.org/x/text/cases"
import "golang.org/x/text/language"

"title": cases.Title(language.English).String,
```

Do not add this dependency now. Add it when a template actually calls it.

### Verification

```bash
grep -r "title" pkg/scaffold/templates/   # must return zero matches
staticcheck ./pkg/scaffold/...             # no suppressed warnings remain
go test ./pkg/scaffold/...
```

---

## FIX-004 — String concatenation in loop: `indentStr`

**Severity:** MEDIUM  
**File:** `main.go`  
**Function:** `indentStr`

### What is wrong

```go
func indentStr(depth int) string {
    s := ""
    for i := 0; i < depth; i++ {
        s += "│   "
    }
    return s + "├── "
}
```

Each `s += "│   "` allocates a new string. For the depth values in a typical generated project tree (< 10), the performance impact is zero — per Pike's Rule 1. The issue is stylistic correctness: Go has `strings.Builder` for exactly this pattern, and using the wrong tool signals unfamiliarity with the stdlib to anyone reading the code.

Secondary: `for i := 0; i < depth; i++` can be written as `for range depth` since Go 1.22, which is already the project's minimum version.

### Fix

```go
func indentStr(depth int) string {
    var sb strings.Builder
    for range depth {
        sb.WriteString("│   ")
    }
    sb.WriteString("├── ")
    return sb.String()
}
```

`strings` is already imported in `main.go` — no new import needed.

### Verification

```bash
go vet ./...
go test ./...   # printTree is cosmetic, no unit test needed — visual inspection sufficient
```

---

## Application order

Apply in this sequence to minimise conflict risk:

```
FIX-003  →  FIX-004  →  FIX-001  →  FIX-002
```

FIX-003 and FIX-004 are single-file, zero-dependency changes. FIX-001 touches one file but removes an import. FIX-002 last because it reorganises the call graph across two packages and requires the regression test before committing.

Commit each fix separately per the project's Conventional Commits policy:

```
fix(scaffold): remove deprecated strings.Title from funcMap
fix(main): use strings.Builder in indentStr
fix(github): replace shell-based lookupEnv with os.LookupEnv
fix(wizard): remove duplicated security enforcement from BuildConfig
```
