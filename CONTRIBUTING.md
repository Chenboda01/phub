# Contributing to phub

Thanks for your interest in contributing to **phub**.

phub is a fast, keyboard-first terminal project manager whose job is simple:

> Find a project, open it, and start working.

Contributions are welcome, but the project should remain small, fast, understandable, and focused.

---

## 1. Before Contributing

Please read:

```text
README.md
ARCHITECTURE.md
ROADMAP.md
CONTRIBUTING.md
```

Before starting larger work, check the current roadmap milestone.

phub is intentionally built in small steps.

Avoid implementing future features early unless explicitly discussed.

---

## 2. Core Principles

Contributions should preserve these ideas:

* Fast startup
* Keyboard-first interaction
* Local-first behavior
* Read-only project discovery
* Explicit project actions
* Simple architecture
* Minimal dependencies
* Clear error messages
* No unexpected project modification

phub should not become an IDE, Git client, file manager, shell, or AI agent.

It should connect projects to those tools.

---

## 3. Development Setup

phub is written in Go.

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/phub.git
cd phub
```

Verify Go:

```bash
go version
```

Install dependencies:

```bash
go mod download
```

Run phub:

```bash
go run ./cmd/phub
```

Run tests:

```bash
go test ./...
```

Format:

```bash
gofmt -w .
```

Check formatting without modifying files:

```bash
gofmt -l .
```

Static analysis:

```bash
go vet ./...
```

---

## 4. Standard Validation

Before submitting a change, run:

```bash
gofmt -w .
go vet ./...
go test ./...
```

If the project later adopts additional tools such as `golangci-lint`, the documented validation commands should be updated.

Do not claim tests passed unless they were actually run.

---

## 5. Branches

Use focused branch names.

Examples:

```text
feat/project-discovery
feat/fuzzy-search
fix/git-status-parsing
fix/missing-tool-error
docs/configuration
test/discovery-depth
```

Avoid vague names such as:

```text
update
changes
stuff
fixes
```

Keep unrelated changes separate.

---

## 6. Commit Style

Prefer focused commits.

Suggested format:

```text
type(scope): description
```

Examples:

```text
feat(discovery): detect Go projects
feat(ui): add project search
fix(git): handle detached HEAD
fix(launcher): report missing executables
test(storage): cover corrupt project database
docs(readme): clarify tool shortcuts
refactor(project): simplify language detection
```

Common types:

```text
feat
fix
docs
test
refactor
perf
build
ci
chore
```

---

## 7. Pull Requests

A pull request should explain:

### Summary

What changed?

### Motivation

Why is the change useful?

### Scope

What is included and intentionally excluded?

### Validation

List the exact commands run.

Example:

```text
gofmt -w . — completed
go vet ./... — passed
go test ./... — 48 tests passed
```

### Screenshots

For visible TUI changes, include a screenshot or recording when useful.

---

## 8. Pull Request Checklist

Before requesting review:

* [ ] The change fits the current roadmap milestone.
* [ ] The change follows `ARCHITECTURE.md`.
* [ ] The change is focused.
* [ ] No unrelated files were modified.
* [ ] New behavior has tests.
* [ ] Failure behavior is tested where important.
* [ ] No project files are modified unexpectedly.
* [ ] No unnecessary dependency was added.
* [ ] User-facing errors are understandable.
* [ ] Documentation matches behavior.
* [ ] `gofmt` has been run.
* [ ] `go vet ./...` passes.
* [ ] `go test ./...` passes.
* [ ] The final diff has been reviewed.

---

## 9. Coding Style

Write straightforward Go.

Prefer:

* Small functions
* Clear package boundaries
* Explicit error handling
* Standard library functionality
* `context.Context` for cancellable work
* Interfaces only where they improve testing or boundaries
* Simple structs
* Early returns
* Table-driven tests where appropriate

Avoid:

* Deep abstractions
* Generic frameworks
* Hidden global state
* Giant `main.go`
* Clever code that is difficult to understand
* Interfaces for every struct
* Reflection without strong justification
* Large dependency trees
* Shell command construction from strings

---

## 10. Error Handling

Errors should tell the user what went wrong.

Good:

```text
Could not launch Lazygit.

Executable "lazygit" was not found in PATH.
```

Bad:

```text
exit status 1
```

Expected failures should not crash the TUI.

Examples include:

* Missing project directory
* Permission denied
* Git not installed
* External tool not installed
* Invalid config
* Corrupted local phub data

---

## 11. Project Discovery

Discovery must remain read-only.

It may inspect project markers such as:

```text
.git/
pyproject.toml
package.json
Cargo.toml
go.mod
pom.xml
CMakeLists.txt
```

It must not:

* Execute project code
* Install dependencies
* Run package scripts
* Modify files
* Follow expensive dependency trees unnecessarily

Discovery changes should include tests for:

* Scan depth
* Ignore rules
* Duplicates
* Missing roots
* Permission errors
* Nested projects

When GitHub-only filtering is involved, tests should use temporary projects and a fake `git` executable or temporary repositories. The filter must inspect local remotes only; it must not contact GitHub or any other network service.

---

## 12. Git Integration

phub uses Git only for lightweight project metadata.

Typical information:

* Repository status
* Branch
* Dirty / clean
* Modified count
* Untracked count
* Ahead / behind later

Use explicit command arguments.

Good:

```go
exec.Command(
    "git",
    "-C",
    projectPath,
    "status",
    "--porcelain",
)
```

Avoid shell composition:

```go
exec.Command(
    "sh",
    "-c",
    "cd "+projectPath+" && git status",
)
```

phub should not implement normal Git workflows that belong in Git or Lazygit.

---

## 13. External Tool Launching

phub may launch tools such as:

```text
nvim
lazygit
forge
yazi
shell
```

Tool launching should:

* Use explicit executable arguments
* Set the project with `cmd.Dir`
* Avoid shell interpolation
* Report missing tools clearly
* Return control to phub cleanly when possible

The Enter action uses phub's embedded PTY terminal. It forwards input to the configured shell and returns to the project list after `exit` or EOF.

Do not automatically install missing tools.

---

## 14. Adding a New Tool Integration

Before adding a built-in integration, ask:

1. Is this something users are likely to launch from a project?
2. Does a dedicated tool already solve the problem?
3. Can it be configured instead of hardcoded?
4. Does it require phub to become responsible for unrelated behavior?

Good examples:

```text
Neovim
Lazygit
Forge
Yazi
shell
```

Poor examples for the core:

```text
Package manager
Docker dashboard
Email client
Web browser
Music player
```

Custom user actions may later support broader workflows without bloating the core.

---

## 15. UI Contributions

The TUI should remain fast and uncluttered.

Important UI principles:

* Project selection is primary.
* Search should be immediate.
* Keyboard actions should be discoverable.
* Small terminals should remain usable.
* Important state must not depend on color alone.
* Theme presets must render an opaque background.
* Nerd Fonts should be optional.
* Avoid permanently empty panels.
* Avoid excessive animation.

Use styling to communicate meaning, not decoration.

---

## 16. Keyboard Bindings

Default shortcuts should remain consistent where practical.

Core navigation:

```text
j / ↓       down
k / ↑       up
Enter       open
r / R       refresh project discovery
/           search
Esc         back / clear
q           quit
?           help
```

At startup, Up/Down selects GitHub-only or all-local discovery and Enter loads that scope. `r` or `R` preserves the selected scope.

Tool actions:

```text
n           Neovim
g           Lazygit
f           Forge
y           Yazi
t           shell
```

Theme selection opens a dropdown with `Ctrl+P`. Up/Down moves the highlight, Enter applies the preset, and Esc cancels. Opening the dropdown alone must not apply a theme.

Changing established shortcuts should be considered a user-facing compatibility change.

---

## 17. Performance

phub should feel instant.

Performance contributions should be based on measurement.

Likely performance-sensitive areas:

* Startup
* Directory scanning
* Git metadata
* Search
* Terminal rendering

Avoid premature complexity such as:

* Databases
* Daemons
* Filesystem watchers
* Persistent metadata caches

until measurements show they are needed.

If metadata is slow, prefer asynchronous loading and bounded concurrency.

---

## 18. Dependencies

Keep dependencies minimal.

A new dependency should solve a meaningful problem.

Before adding one, consider:

* Can the Go standard library handle this?
* Is an existing dependency already enough?
* Is the package actively maintained?
* How large is its transitive dependency tree?
* Is the license compatible?
* Will it make installation or builds more fragile?

Do not add libraries for trivial helpers.

---

## 19. Storage

phub stores only its own local metadata.

Expected location:

```text
~/.local/share/phub/
```

Examples:

* Known projects
* Favorites
* Recent usage
* Storage format version

Project metadata storage must not modify the project itself.

Persistent writes should be atomic where practical.

Corrupted storage should not make phub unusable.

---

## 20. Configuration

User configuration should live under:

```text
~/.config/phub/
```

Current planned configuration:

```text
~/.config/phub/config.toml
```

Config changes should:

* Have safe defaults
* Produce clear validation errors
* Preserve backward compatibility when reasonable
* Avoid silently changing tool behavior

---

## 21. Testing

Tests should not depend on the contributor's actual environment.

Use temporary directories for:

* Project discovery
* Git repositories
* Storage
* Configuration

Do not:

* Modify real user projects
* Modify global Git config
* Launch real Neovim
* Launch real Lazygit
* Launch real Forge
* Launch real shells during normal tests

Abstract launching so command construction can be tested without execution.

---

## 22. Table-Driven Tests

Go table-driven tests are encouraged where they improve clarity.

Example:

```go
func TestDetectLanguage(t *testing.T) {
    tests := []struct {
        name     string
        files    []string
        expected Language
    }{
        {
            name:     "python",
            files:    []string{"pyproject.toml"},
            expected: Python,
        },
        {
            name:     "go",
            files:    []string{"go.mod"},
            expected: Go,
        },
    }

    // ...
}
```

Keep test cases readable.

---

## 23. Documentation

Update documentation in the same change when user-visible behavior changes.

Documentation should clearly distinguish:

* Implemented
* Planned
* Experimental

Do not document planned features as though they already exist.

Examples and commands should be safe to copy.

---

## 24. Generated or AI-Assisted Contributions

AI-assisted code is allowed.

The contributor remains responsible for:

* Understanding the code
* Reviewing the result
* Verifying imports and APIs
* Running tests
* Checking architecture consistency
* Removing unnecessary generated complexity
* Confirming the code actually works

Do not submit large generated changes you cannot explain.

phub should remain small enough that maintainers understand the codebase.

---

## 25. Architecture Changes

Significant changes to the architecture should be discussed before large implementation work.

Examples:

* Introducing a database
* Adding a plugin system
* Adding a background daemon
* Replacing Bubble Tea
* Adding network services
* Adding arbitrary project command execution
* Changing persistent storage format substantially

A complicated solution should demonstrate why the simpler architecture is insufficient.

---

## 26. Scope

If you discover an unrelated problem while working:

* Report it.
* Open an issue if appropriate.
* Do not automatically expand the pull request.

Small blocking fixes are acceptable when required for the requested work, but they should be clearly identified.

---

## 27. Good First Contributions

Useful early contributions include:

* Language detectors
* Project marker tests
* Error-message improvements
* Terminal-size handling
* Documentation fixes
* Tool detection
* Git parsing edge cases
* Search improvements
* Discovery ignore rules
* Unit tests

Avoid beginning with major architecture rewrites.

---

## 28. What phub Should Not Become

Contributions should resist turning phub into:

* An IDE
* A source editor
* A file manager
* A Git client
* An AI coding agent
* A shell
* A terminal emulator
* A package manager
* A system dashboard

Those are different projects.

phub should integrate with them.

---

## 29. Definition of Done

A contribution is complete when:

* The requested behavior works.
* The change is focused.
* New behavior is tested.
* Important failure paths are handled.
* Project discovery remains read-only.
* External commands are constructed safely.
* No unnecessary dependency was added.
* Documentation matches reality.
* Embedded terminal behavior is covered without launching real interactive tools in ordinary unit tests.
* `gofmt` has been run.
* `go vet ./...` passes.
* `go test ./...` passes.
* The final diff has been reviewed.

If a required check cannot be completed, say so clearly.

---

## 30. Final Contribution Principle

When implementation choices conflict, prefer:

```text
Correctness
Simplicity
Responsiveness
Maintainability
Safety
Keyboard usability
Compatibility
Visual polish
Feature count
```

phub should remain a tool that one developer can understand.

> **Find it. Open it. Work.**
