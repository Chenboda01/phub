# phub

> Your projects. One terminal.

**phub** is a fast, keyboard-first terminal project manager for developers.

Instead of remembering where every project lives, repeatedly `cd`-ing through directories, checking Git status manually, activating environments, and opening the same development tools over and over, phub gives you one place to find, inspect, and launch your projects.

```text
$ phub
```

```text
┌─ phub ───────────────────────────────────────────────────────────┐
│ Search projects...                                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PROJECT          LANG        BRANCH       STATUS        LAST   │
│                                                                 │
│  Forge            Python      main         ● 3 changes   now    │
│  ScreenBot        Python      main         ✓ clean       2d     │
│  Website          TypeScript  dev          ↑ 2           4d     │
│  67-Challenge     Python      main         ✓ clean       12d    │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ enter open   n nvim   g lazygit   f forge   t terminal   ? help│
└─────────────────────────────────────────────────────────────────┘
```

## Why phub?

Developers often have projects scattered across:

```text
~/Projects/
~/Code/
~/Python/
~/Documents/
~/GitHub/
~/School/
~/Experiments/
```

Finding the project is only the beginning.

Then comes:

```bash
cd ~/Projects/Forge
source .venv/bin/activate
git status
nvim .
```

Or:

```bash
cd ~/Code/website
lazygit
```

Or:

```bash
cd ~/Python/game
forge
```

phub turns that into:

```bash
phub
```

Select the project.

Press a key.

Done.

---

# Philosophy

phub should be:

* **Fast** — opening phub should feel instant.
* **Keyboard-first** — normal workflows should not require a mouse.
* **Useful** — every visible element should help manage projects.
* **Simple** — phub is a project manager, not an IDE.
* **Tool-friendly** — phub launches the tools you already use.
* **Local-first** — project information stays on your computer.
* **Predictable** — phub should never modify a project unexpectedly.

phub does not try to replace:

* Neovim
* Yazi
* Lazygit
* Git
* Forge
* Your terminal
* Your shell

It connects them.

---

# Core Idea

phub treats a **project** as more than a directory.

A project may have:

```text
Path
Language
Git repository
Current branch
Git status
Development environment
Run command
Test command
Editor
Recently used time
Favorite status
Custom actions
```

phub discovers that information and puts it in one place.

---

# Project Discovery

phub can scan configured directories for projects.

Example:

```bash
phub scan ~/Projects
```

It may discover:

```text
~/Projects/Forge
~/Projects/website
~/Projects/game
~/Projects/phub
```

A directory may be recognized as a project through files such as:

```text
.git/
pyproject.toml
package.json
Cargo.toml
go.mod
pom.xml
Makefile
CMakeLists.txt
Dockerfile
```

Users may also register a directory manually:

```bash
phub add ~/some/project
```

---

# Main Interface

The default interface shows projects rather than individual files.

```text
┌─ Projects ──────────────────────────────────────────────────────┐
│                                                                │
│  Forge                                                         │
│  Python · main · 3 modified · used now                         │
│                                                                │
│  phub                                                          │
│  Go · main · clean · used 1h ago                               │
│                                                                │
│  Website                                                       │
│  TypeScript · dev · ahead 2 · used yesterday                   │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

Projects should be searchable immediately by typing.

---

# Project View

Selecting a project opens its details.

```text
┌─ Forge ─────────────────────────────────────────────────────────┐
│                                                               │
│ Path          ~/Projects/Forge                                │
│ Language      Python                                          │
│ Git           main                                            │
│ Status        3 modified                                      │
│ Environment   .venv                                           │
│ Tests         pytest                                          │
│ Last used     2 minutes ago                                   │
│                                                               │
│ Actions                                                       │
│                                                               │
│ [enter] Open                                                  │
│ [n]     Neovim                                                │
│ [g]     Lazygit                                               │
│ [f]     Forge                                                 │
│ [t]     Terminal                                              │
│ [r]     Run                                                   │
│ [x]     Tests                                                 │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

---

# Quick Launch

phub should make common actions extremely fast.

Open the TUI:

```bash
phub
```

Jump directly to a project:

```bash
phub forge
```

Open a project in Neovim:

```bash
phub open forge --with nvim
```

Open Lazygit:

```bash
phub open forge --with lazygit
```

Open Forge:

```bash
phub open forge --with forge
```

Open a shell in the project:

```bash
phub open forge --with shell
```

The exact CLI may evolve as phub develops.

---

# Keyboard Controls

The TUI should support familiar Vim-style navigation.

```text
j / ↓       Move down
k / ↑       Move up
g g         First project
G           Last project
Enter       Open project
/           Search
Esc         Clear / go back
q           Quit
?           Help
```

Project actions:

```text
n           Open in Neovim
g           Open Lazygit
f           Open Forge
t           Open terminal
r           Run project
x           Run tests
e           Edit project configuration
*           Favorite
```

Keyboard bindings should eventually be configurable.

---

# Search

Search should be one of phub's strongest features.

Press:

```text
/
```

and type:

```text
for
```

Results update immediately:

```text
Forge
Forge-Test
terraform-experiments
```

Search may match:

* Project name
* Path
* Language
* Tags
* Git branch

Future versions may support filters:

```text
lang:python
git:dirty
favorite:true
branch:main
```

---

# Git Awareness

phub should understand Git without trying to replace Git tools.

Example:

```text
Forge

Branch        main
Working tree  3 modified
Ahead         2
Behind        0
```

Status indicators might include:

```text
✓ clean
● modified
? untracked
↑ ahead
↓ behind
↕ diverged
```

For detailed Git operations, phub should launch a dedicated tool such as Lazygit.

phub should not automatically:

* Commit
* Push
* Pull
* Reset
* Rebase
* Delete branches

---

# Language Detection

phub should detect common project types.

Initial targets:

### Python

Detect:

```text
pyproject.toml
requirements.txt
setup.py
Pipfile
```

Possible environments:

```text
.venv/
venv/
```

### JavaScript / TypeScript

Detect:

```text
package.json
package-lock.json
pnpm-lock.yaml
yarn.lock
bun.lock
```

### Go

Detect:

```text
go.mod
```

### Rust

Detect:

```text
Cargo.toml
```

### Java

Detect:

```text
pom.xml
build.gradle
build.gradle.kts
```

### C / C++

Detect:

```text
CMakeLists.txt
Makefile
meson.build
```

More project types can be added later.

---

# Environment Detection

phub should display useful development-environment information without automatically changing it.

Example:

```text
Python
  Version       3.12.7
  Environment   .venv
  Package tool  uv

Git
  Branch        main
  Status        clean

Testing
  pytest

Formatting
  Ruff
```

phub may eventually detect:

* Python virtual environments
* pyenv
* Node versions
* npm
* pnpm
* Yarn
* Bun
* Cargo
* Go toolchain
* Docker
* Dev containers

---

# Recent Projects

phub should remember which projects you actually use.

```bash
phub recent
```

Example:

```text
Forge          now
phub           18 minutes ago
ScreenBot      yesterday
Website        3 days ago
```

Recent activity should be stored locally.

Opening a project through phub updates its last-used timestamp.

---

# Favorites

Frequently used projects can be pinned.

```text
★ Forge
★ phub
  ScreenBot
  Website
```

Favorites appear before ordinary projects unless another sort is selected.

---

# Project Health

A future project-health view may show:

```text
Forge

Git             ✓
Environment     ✓
Dependencies    ✓
Tests           ✓
Uncommitted     3 files
Updates         unknown

Health          GOOD
```

phub should not claim a project is healthy merely because the directory exists.

Health indicators must correspond to real checks.

---

# Custom Project Actions

Projects may define custom actions.

Example:

```toml
[actions]
run = ["python", "app.py"]
test = ["pytest"]
format = ["ruff", "format", "."]
check = ["ruff", "check", "."]
```

Then phub can display:

```text
[r] Run
[x] Test
[c] Check
```

Commands should be shown before execution when they may modify project state.

---

# Configuration

User configuration should live at:

```text
~/.config/phub/config.toml
```

Example:

```toml
[general]
editor = "nvim"
terminal = "alacritty"

[scan]
directories = [
    "~/Projects",
    "~/Code",
    "~/Python",
]

[tools]
git = "lazygit"
ai = "forge"
files = "yazi"
```

Project-specific configuration may later live in:

```text
.phub.toml
```

Example:

```toml
name = "Forge"

[commands]
run = ["forge"]
test = ["pytest"]
check = ["ruff", "check", "."]
```

---

# Tool Integration

phub should work well with existing terminal tools.

### Neovim

```text
n → nvim .
```

### Lazygit

```text
g → lazygit
```

### Forge

```text
f → forge
```

### Yazi

```text
y → yazi
```

### Terminal

```text
t → open shell in project
```

phub should allow users to replace these tools through configuration.

---

# Safety

phub is a project launcher and manager.

It should not unexpectedly modify projects.

By default, phub may safely:

* Scan directories
* Read project metadata
* Read Git status
* Detect languages
* Detect environments
* Store phub-local metadata
* Launch user-configured tools

Actions that can modify project state should be explicit.

phub should never automatically:

* Delete projects
* Delete files
* Commit changes
* Push Git repositories
* Install packages
* Modify project configuration
* Execute unknown project scripts
* Run commands as root

Project configuration and repository contents should be treated as untrusted input.

---

# Performance

phub should feel instant.

Project scanning should avoid repeatedly walking enormous directories.

Directories such as these should normally be ignored:

```text
.git/
node_modules/
.venv/
venv/
target/
dist/
build/
__pycache__/
```

Git and environment information may be refreshed asynchronously.

The interface should remain responsive while metadata is being collected.

Future versions may maintain a lightweight local project index.

---

# Local Data

phub should keep its own state under:

```text
~/.local/share/phub/
```

Possible contents:

```text
projects.json
recent.json
cache/
```

Temporary cache should be safe to delete.

User configuration remains under:

```text
~/.config/phub/
```

phub should not require an account or cloud service.

---

# Planned CLI

```text
phub
phub PROJECT
phub list
phub add PATH
phub remove PROJECT
phub scan PATH
phub recent
phub favorites
phub open PROJECT
phub doctor
phub --version
phub --help
```

The CLI should remain useful even when the interactive TUI is unavailable.

---

# Technology

phub is planned as a terminal application written in **Go**.

The TUI may use **Bubble Tea** and related Charm ecosystem libraries.

Why Go?

* Fast startup
* Excellent filesystem APIs
* Easy concurrency
* Simple deployment
* Single executable
* Strong terminal-tool ecosystem
* Good cross-platform support

The project should avoid unnecessary dependencies.

---

# Initial Architecture

A possible starting structure:

```text
phub/
├── cmd/
│   └── phub/
│       └── main.go
│
├── internal/
│   ├── app/
│   ├── project/
│   ├── discovery/
│   ├── git/
│   ├── launcher/
│   ├── config/
│   ├── storage/
│   └── ui/
│
├── tests/
├── README.md
├── go.mod
└── LICENSE
```

The architecture should remain simple until real complexity requires additional layers.

---

# v0.1 Scope

The first version should be intentionally small.

## Required

* [ ] Interactive TUI
* [ ] Scan configured project directories
* [ ] Detect Git repositories
* [ ] Detect basic project language
* [ ] Fuzzy project search
* [ ] Recent projects
* [ ] Favorite projects
* [ ] Show Git branch and dirty/clean status
* [ ] Open project in terminal
* [ ] Open project in Neovim
* [ ] Open Lazygit
* [ ] Open Forge
* [ ] TOML configuration
* [ ] `phub --help`
* [ ] `phub --version`

## Not Required for v0.1

* Package management
* File management
* Git operations
* AI features
* Cloud sync
* Plugin system
* Project deletion
* Dependency installation
* Remote repositories
* Containers
* Complex environment management

The first goal is simple:

> Find any project and enter its development environment in seconds.

---

# Future Ideas

Possible later features include:

* Project tags
* Workspaces
* Multiple project groups
* Test status
* Project health checks
* Environment activation
* Custom actions
* GitHub integration
* Repository cloning
* SSH projects
* Docker projects
* Project templates
* Session restoration
* Terminal multiplexers
* `tmux` integration
* Zellij integration
* Project notes
* Command history
* Plugin API

These features should be added only when the core project-navigation experience is excellent.

---

# Non-Goals

phub is not intended to become:

* An IDE
* A file manager
* A Git client
* An AI coding agent
* A package manager
* A shell
* A terminal emulator

Dedicated tools already solve those problems well.

phub's job is to connect projects to those tools.

---

# Example Workflow

Start:

```bash
phub
```

Search:

```text
for
```

Select:

```text
Forge
```

Press:

```text
n
```

Result:

```text
Neovim opens in ~/Projects/Forge
```

Or:

```text
g
```

Result:

```text
Lazygit opens in ~/Projects/Forge
```

Or:

```text
f
```

Result:

```text
Forge opens with ~/Projects/Forge as its workspace
```

That's phub.

---

# Project Status

**Early development.**

The current priority is building a small, fast, reliable project launcher before expanding into advanced project management.

---

# Guiding Principle

phub should answer one question extremely well:

> **What do I want to work on right now?**

And then get out of the way.
