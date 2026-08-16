# phub Roadmap

## 1. Purpose

This roadmap defines the planned development sequence for **phub**.

phub should be built in small, usable milestones.

Each milestone should leave the application working.

The goal is not to add many features quickly.

The goal is to make one workflow excellent:

> Find a project, open it, and start working immediately.

---

# 2. Development Principles

Every milestone should:

* Keep phub runnable.
* Keep startup fast.
* Avoid unrelated refactors.
* Include tests for new behavior.
* Preserve local-first behavior.
* Avoid unnecessary dependencies.
* Prefer simple implementations.
* Keep project discovery read-only.
* Keep launching explicit.
* Avoid implementing future milestones early.

Do not build the entire roadmap in one pass.

---

# 3. Version Plan

```text
0.0.1   TUI prototype (implemented)
0.0.2   Project discovery (implemented)
0.0.3   Git and language metadata
0.0.4   Search
0.0.5   Project launching and theme presets (implemented)
0.0.6   Tool integrations
0.0.7   Favorites and recent projects
0.0.8   Configuration and persistence polish
0.0.9   Diagnostics and hardening
0.1.0   First stable usable release
```

---

# Milestone 0.0.1 — TUI Prototype

## Goal

Build the basic phub terminal interface.

Use a small hardcoded project list.

## Required

* Bubble Tea application
* Project list
* Selected-project state
* `j` / `k` navigation
* Arrow-key navigation
* `/` search input placeholder
* `Enter` action placeholder (replaced by real shell launch in 0.0.5)
* `q` quit
* Responsive terminal resizing
* Help text

Example:

```text
┌─ phub ───────────────────────────────────────────┐
│ Search projects...                              │
├─────────────────────────────────────────────────┤
│ > Aster                                         │
│   Dune                                          │
│   Beacon                                     │
│   Cobalt                                       │
├─────────────────────────────────────────────────┤
│ ↑/↓ navigate · enter open · / search · q quit  │
└─────────────────────────────────────────────────┘
```

## Non-Goals

* Real project discovery
* Git
* Configuration
* Persistence
* Tool launching
* Favorites

## Success Criteria

* `phub` starts successfully.
* Navigation feels immediate.
* Terminal resizing works.
* `q` exits cleanly.
* Ctrl+C restores terminal state correctly.

---

# Milestone 0.0.2 — Project Discovery

## Status

Implemented.

## Goal

Replace fake projects with real project scanning.

## Required

Support configured scan roots such as:

```text
~/Projects
~/Code
~/Python
```

When `PHUB_SCAN_ROOTS` is not set, the default scan root is the user's home directory. At startup, the TUI presents a scope chooser with **GitHub only** selected by default. Press Enter to load projects with a local `github.com` remote, or move to **All local projects** and press Enter to keep every discovered project. Press `r` or `R` to scan the default or configured roots again while preserving the selected scope.

The default roots can be overridden with `PHUB_SCAN_ROOTS`. Set `PHUB_SCAN_DEPTH` to a non-negative integer to change the default maximum depth of `4`. On Unix-like systems, separate multiple roots with `:`.

Recognize projects using markers such as:

```text
.git/
pyproject.toml
package.json
Cargo.toml
go.mod
pom.xml
CMakeLists.txt
build.gradle
build.gradle.kts
```

## Required Behavior

* Configurable maximum depth
* Ignore common dependency/build directories
* Prevent duplicate projects
* Handle missing scan roots
* Handle permission errors gracefully
* Use canonical project paths
* Choose GitHub-only or all-local discovery at startup
* Detect GitHub repositories from local Git remotes without network access
* Refresh discovery with `r` or `R` without modifying project files

## Initial Ignore List

```text
.git
node_modules
.venv
venv
target
dist
build
__pycache__
.cache
```

## Success Criteria

Running phub displays real projects discovered from configured directories.

Scanning does not execute project code.

GitHub-only startup filtering uses local `git remote -v` metadata and makes no network requests.

Pressing `r` or `R` discovers projects added since startup.

---

# Milestone 0.0.3 — Project Metadata

## Status

Implemented.

## Goal

Make the project list useful at a glance.

## Required

Detect:

* Project name
* Path
* Primary language
* Git repository status
* Current branch
* Dirty / clean state

Example:

```text
Aster       Python      main      ● dirty
phub        Go          main      ✓ clean
Cobalt     TypeScript  dev       ✓ clean
```

## Language Targets

Initial:

* Python
* JavaScript
* TypeScript
* Go
* Rust
* Java
* C / C++

## Git Targets

Initial:

* Is repository
* Branch
* Clean / dirty
* Modified-file count
* Untracked-file count

## Success Criteria

Metadata failures do not prevent the project from appearing.

A missing Git executable produces a useful warning rather than a crash.

---

# Milestone 0.0.4 — Search

## Goal

Make project navigation extremely fast.

## Required

Typing `/` activates search.

Search matches:

* Project name
* Path
* Language

Search should update results as the user types.

Example:

```text
Search: for

Aster
terraform-lab
```

## Nice to Have

Basic fuzzy matching.

Do not add a heavyweight indexing system.

## Success Criteria

Searching hundreds of projects feels instantaneous.

Esc clears search and returns to the full project list.

---

# Milestone 0.0.5 — Project Opening

## Status

Implemented with an embedded PTY-backed terminal and theme presets. Git and language metadata reload asynchronously after the terminal exits.

## Goal

Allow the user to enter a selected project.

## Required

Pressing `Enter` should open the configured shell in phub's embedded terminal with the project directory as its working directory. Typing `exit` or pressing `Ctrl+D` should return to phub.

No manual `cd` should be necessary.

Example:

```text
Select Aster
Press Enter
```

Result:

```text
shell working directory:
~/Projects/Aster
```

## Required Behavior

* Use project directory as `cmd.Dir`.
* Do not construct shell commands using string concatenation.
* Render shell output through a terminal emulator so ANSI applications remain usable.
* Forward keyboard input to the PTY in order.
* Restore the project list when the child process exits.
* Refresh available project metadata afterward when the metadata milestone is implemented.

## Success Criteria

The user can enter a real project directly from phub.

## Theme Presets

The TUI also provides an 18-item non-transparent color dropdown. Press `Ctrl+P` to open it, use Up/Down to highlight a preset, press Enter to apply it, or press Esc to cancel. Opening the picker alone never changes the theme.

```text
> Red Background
  Red Theme
  Red Combo
  ...
  Purple Combo
```

Background presets change the terminal background, Theme presets change the foreground while keeping an opaque neutral background, and Combo presets change both.

---

# Milestone 0.0.6 — Tool Integrations

## Goal

Launch common development tools from the selected project.

## Initial Actions

```text
n   Neovim
g   Lazygit
f   Forge
y   Yazi
t   Shell
```

## Required

* Detect whether configured tools exist.
* Show missing tools clearly.
* Use the project as the working directory.
* Suspend and resume the TUI cleanly.

Example:

```text
n → nvim .
g → lazygit
f → forge
y → yazi .
```

## Success Criteria

A user can move from phub into their normal tools without manually navigating directories.

---

# Milestone 0.0.7 — Favorites and Recent Projects

## Goal

Make commonly used projects easier to reach.

## Favorites

Allow:

```text
*   toggle favorite
```

Display:

```text
★ Aster
★ phub
  Cobalt
```

Favorites should sort before ordinary projects by default.

## Recent Projects

Opening a project updates its last-used timestamp.

Possible display:

```text
Aster        now
phub         1h
Cobalt      2d
```

## Success Criteria

Favorite and recent state persists across phub restarts.

Removing a project from phub does not delete its files.

---

# Milestone 0.0.8 — Configuration and Persistence

## Goal

Make phub comfortable for daily use.

## Configuration

Use:

```text
~/.config/phub/config.toml
```

Initial settings:

```toml
[general]
editor = "nvim"
git_ui = "lazygit"
ai = "forge"
file_manager = "yazi"

[scan]
directories = [
    "~/Projects",
    "~/Code",
    "~/Python",
]

max_depth = 4
```

## Persistent Data

Use:

```text
~/.local/share/phub/
```

Store:

* Known projects
* Manual projects
* Favorites
* Last-used timestamps
* Storage format version

## Required

* Safe defaults
* Atomic writes
* Corrupt-storage recovery
* Configuration validation
* Helpful errors

## Success Criteria

A new user can install phub, configure scan roots, and use it daily without touching source code.

---

# Milestone 0.0.9 — Diagnostics and Hardening

## Goal

Make phub dependable.

## Add

```bash
phub --help
phub --version
phub doctor
phub list
phub add PATH
phub remove PROJECT
phub scan
phub recent
```

## `phub doctor`

Check:

* Config directory
* Data directory
* Git
* Neovim
* Lazygit
* Forge
* Yazi
* Scan roots
* Storage health

Example:

```text
phub doctor

Git        ✓
Neovim     ✓
Lazygit    ✓
Forge      ✓
Yazi       ✗ not found

Scan roots
~/Projects ✓
~/Python   ✓
```

## Hardening

Test:

* Missing projects
* Deleted projects
* Permission-denied directories
* Corrupt config
* Corrupt storage
* Missing external tools
* Unusual paths
* Unusual Git branch names
* Terminal resize
* Small terminal dimensions

## Success Criteria

Expected failures are understandable and do not crash phub.

---

# Milestone 0.1.0 — First Usable Release

## Goal

Declare the first genuinely usable phub release.

phub 0.1.0 should reliably support:

1. Starting quickly.
2. Discovering configured project directories.
3. Showing projects in a responsive TUI.
4. Searching projects.
5. Showing language and Git status.
6. Opening a shell in a project.
7. Opening Neovim.
8. Opening Lazygit.
9. Opening Aster.
10. Opening Yazi.
11. Remembering favorites.
12. Remembering recently used projects.
13. Loading user configuration.
14. Recovering gracefully from common errors.
15. Running basic CLI commands without the TUI.

---

# 4. Testing Roadmap

Testing should grow with the project.

## 0.0.1

Test:

* TUI state transitions
* Navigation
* Quit behavior

## 0.0.2

Add:

* Project marker detection
* Scan depth
* Ignore rules
* Duplicate prevention

## 0.0.3

Add:

* Language detection
* Git status parsing
* Non-Git directories

## 0.0.4

Add:

* Search filtering
* Case-insensitive matching

## 0.0.5–0.0.6

Add:

* Launcher command construction
* Working-directory handling
* Missing executable behavior
* Enter-to-shell behavior
* Embedded PTY terminal lifecycle and exit handoff
* Ctrl+P dropdown navigation and theme selection
* Opaque background rendering for all theme presets
* Refresh-key command construction and project-list replacement
* Startup scope selection and scope-preserving refresh

Do not launch real Neovim, Lazygit, or Forge during normal tests.

## 0.0.7–0.0.8

Add:

* Favorites
* Recent timestamps
* Storage persistence
* Atomic writes
* Configuration parsing
* Corruption handling

---

# 5. Performance Roadmap

Do not optimize prematurely.

Initial expectations:

```text
Project count:        10–500
Search:               immediate
Startup first render: fast
Git metadata:         asynchronous
```

If Git inspection becomes slow:

1. Measure it.
2. Add bounded concurrency.
3. Consider session caching only if necessary.

Do not introduce a database or daemon merely for performance speculation.

---

# 6. Explicitly Deferred

Do not build these before 0.1.0:

* AI project recommendations
* Remote GitHub API integration
* GitLab integration
* Repository cloning
* SSH project management
* Plugin system
* Docker dashboard
* Package management
* Built-in terminal
* Built-in editor
* Built-in Git client
* Built-in file manager
* Semantic code indexing
* Cloud sync
* Collaboration
* Background daemon
* Automatic dependency installation
* Automatic project repair

These ideas may be reconsidered after phub's core workflow is proven useful.

---

# 7. Possible Post-0.1 Roadmap

After actual daily usage, consider only features that solve observed problems.

Possible directions:

## 0.2

* Project tags
* Custom sorting
* Better fuzzy search
* Project groups
* Custom project actions

## 0.3

* Environment detection
* Test-command detection
* Project health view
* Better Git metadata

## 0.4

* tmux integration
* Zellij integration
* Workspace/session launching

## 0.5

* Remote repository support
* Optional GitHub integration
* Repository cloning

These are not commitments.

---

# 8. Definition of Done

A milestone is complete only when:

* The required behavior works.
* Tests cover the important behavior.
* Existing functionality still works.
* No unnecessary future features were added.
* Error handling is reasonable.
* The TUI remains responsive.
* Documentation matches reality.
* The change is understandable and maintainable.
* No project files are modified unexpectedly.

---

# 9. Priority Rule

When choosing between competing implementation options:

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

Do not sacrifice simplicity for features that phub does not yet need.

---

# 10. Final Goal

phub 0.1.0 should make this workflow feel effortless:

```text
$ phub
```

Type:

```text
for
```

Select:

```text
Aster
```

Press:

```text
n
```

Neovim opens in Aster.

Or press:

```text
g
```

Lazygit opens.

Or press:

```text
f
```

Forge opens.

No hunting through directories.

No repeated `cd`.

No unnecessary friction.

> **Find it. Open it. Work.**
