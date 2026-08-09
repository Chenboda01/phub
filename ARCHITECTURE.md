# phub Architecture

## 1. Purpose

This document defines the architecture of **phub**, a fast, keyboard-first terminal project manager.

phub has one primary responsibility:

> Find a project, understand its basic development state, and launch the right tool in that project as quickly as possible.

phub is not an IDE, file manager, Git client, AI coding agent, shell, or package manager.

It connects projects to tools that already solve those problems well.

Examples:

```text
phub → project
        ├── nvim
        ├── lazygit
        ├── forge
        ├── yazi
        └── shell
```

The architecture should remain small, modular, fast, and understandable.

---

# 2. Core Principles

## 2.1 Fast Startup

Running:

```bash
phub
```

should feel nearly instantaneous.

Avoid expensive work on the startup path.

Slow metadata such as Git information may be loaded asynchronously.

---

## 2.2 Keyboard First

Every important operation must be accessible without a mouse.

Primary navigation should support familiar Vim-style controls:

```text
j / ↓       down
k / ↑       up
Enter       open
/           search
Esc         back
q           quit
```

---

## 2.3 Projects, Not Files

phub manages **projects**.

Yazi manages files.

Lazygit manages Git.

Neovim edits code.

Forge provides AI coding assistance.

phub should not duplicate those tools.

---

## 2.4 Local First

Project metadata should remain local.

phub should not require:

* An account
* Cloud storage
* Telemetry
* A remote server
* Internet access

The basic application should work completely offline.

---

## 2.5 Read-Only Discovery

Scanning a project should not modify it.

Project discovery may inspect:

```text
.git/
pyproject.toml
package.json
Cargo.toml
go.mod
pom.xml
CMakeLists.txt
Makefile
```

but should not execute project code.

---

## 2.6 Explicit Actions

Launching a tool should happen only because the user requested it.

phub must not automatically:

* Run project scripts
* Install dependencies
* Commit changes
* Push repositories
* Delete files
* Modify configuration
* Start development servers

---

## 2.7 Simple Before Clever

Do not introduce:

* Plugin systems
* Dependency injection frameworks
* Event buses
* Databases
* RPC
* AI
* Semantic indexing
* Complex caching

until actual requirements justify them.

phub should remain understandable.

---

# 3. High-Level Architecture

```text
                         User
                           │
                           ▼
                    ┌─────────────┐
                    │     CLI     │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │     App     │
                    └──────┬──────┘
                           │
              ┌────────────┼─────────────┐
              │            │             │
              ▼            ▼             ▼
         Discovery      Registry       Launcher
              │            │             │
              ▼            ▼             ▼
         Project Info    Storage     External Tools
              │
       ┌──────┼──────┐
       ▼      ▼      ▼
    Git     Language Environment
```

The TUI displays application state.

It should not contain project-discovery or process-launching logic.

---

# 4. Recommended Repository Structure

```text
phub/
├── cmd/
│   └── phub/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── state.go
│   │
│   ├── project/
│   │   ├── project.go
│   │   ├── language.go
│   │   └── status.go
│   │
│   ├── discovery/
│   │   ├── scanner.go
│   │   └── detector.go
│   │
│   ├── git/
│   │   └── git.go
│   │
│   ├── launcher/
│   │   └── launcher.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── storage/
│   │   └── storage.go
│   │
│   └── ui/
│       ├── model.go
│       ├── update.go
│       ├── view.go
│       ├── keys.go
│       └── styles.go
│
├── README.md
├── ARCHITECTURE.md
├── ROADMAP.md
├── CONTRIBUTING.md
├── go.mod
├── go.sum
└── LICENSE
```

Do not create every file immediately.

Create packages as functionality requires them.

---

# 5. Core Project Model

The central data structure is a project.

A possible model:

```go
type Project struct {
    ID       string
    Name     string
    Path     string
    Language Language

    Git      GitInfo
    Env      EnvironmentInfo

    Favorite bool
    LastUsed time.Time
}
```

This model represents information phub knows about a project.

It should not contain UI state.

---

# 6. Project Identity

Projects need stable identities.

A project should initially be identified by its canonical absolute path.

For example:

```text
/home/boda/Projects/Forge
```

Project display names may change.

Paths are therefore more reliable identifiers than project names.

A derived ID may later be generated from the canonical path.

---

# 7. Project Discovery

The discovery system finds projects inside configured directories.

Example configuration:

```toml
[scan]

directories = [
    "~/Projects",
    "~/Code",
    "~/Python",
]
```

The discovery layer should expose behavior similar to:

```go
type Scanner interface {
    Scan(ctx context.Context, roots []string) ([]project.Project, error)
}
```

The scanner should not know anything about the TUI.

---

# 8. What Counts as a Project?

A directory may be recognized as a project if it contains one or more project markers.

Initial markers:

```text
.git/
pyproject.toml
package.json
Cargo.toml
go.mod
pom.xml
build.gradle
build.gradle.kts
CMakeLists.txt
```

A `.git` directory or file is a strong project indicator.

Other markers may identify projects that are not Git repositories.

Users must also be able to register projects manually.

---

# 9. Scan Boundaries

Project scanning must avoid expensive recursive traversal.

Do not descend into directories such as:

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

These directories cannot contain projects phub needs to discover during normal scanning.

The ignore list should eventually be configurable.

---

# 10. Scan Depth

Scanning an entire home directory recursively can be extremely expensive.

phub should use a configurable maximum discovery depth.

Example:

```toml
[scan]

max_depth = 4
```

If:

```text
~/Projects
```

is configured, phub might inspect:

```text
~/Projects/Forge
~/Projects/Python/game
~/Projects/Web/website
```

without recursively crawling dependency trees.

---

# 11. Language Detection

Language detection belongs in the project layer, not the UI.

Initial detection rules may use project markers.

Examples:

```text
pyproject.toml       → Python
requirements.txt     → Python

package.json         → JavaScript / TypeScript

go.mod               → Go

Cargo.toml           → Rust

pom.xml              → Java

CMakeLists.txt       → C / C++
```

Where multiple languages are present, phub may choose a primary language.

Early versions should use deterministic rules.

Do not add language-analysis engines until necessary.

---

# 12. Git Integration

Git support should provide lightweight repository information.

Initial information:

```go
type GitInfo struct {
    IsRepository bool
    Branch       string
    Dirty        bool
    Modified     int
    Untracked    int
    Ahead        int
    Behind       int
}
```

phub should use Git primarily for **inspection**.

It is not a Git client.

Detailed Git operations belong in Lazygit or Git itself.

---

# 13. Git Execution

Prefer invoking the installed `git` executable with explicit arguments.

Example:

```go
exec.Command(
    "git",
    "-C",
    projectPath,
    "status",
    "--porcelain",
)
```

Avoid shell execution such as:

```go
exec.Command(
    "sh",
    "-c",
    "cd "+path+" && git status",
)
```

Project paths must never be concatenated into shell commands.

---

# 14. Environment Detection

phub may inspect development environments.

Initial detection can include:

### Python

```text
.venv/
venv/
pyproject.toml
```

### Node

```text
package.json
node_modules/
```

### Rust

```text
Cargo.toml
```

### Go

```text
go.mod
```

Environment detection should remain read-only.

phub should not automatically activate, install, or repair environments.

---

# 15. Project Registry

Discovery answers:

> What projects exist right now?

The registry answers:

> What projects does phub know about?

The registry combines:

```text
Discovered projects
+
Manually added projects
+
Favorites
+
Recent project information
```

Conceptually:

```go
type Registry interface {
    Projects() []project.Project
    Add(project.Project) error
    Remove(string) error
    Favorite(string, bool) error
    MarkUsed(string) error
}
```

The exact interface may evolve.

---

# 16. Persistent Storage

phub should store local application state under:

```text
~/.local/share/phub/
```

Initial storage may be simple JSON.

Example:

```text
~/.local/share/phub/projects.json
```

Possible data:

```json
{
  "projects": [
    {
      "path": "/home/boda/Projects/Forge",
      "favorite": true,
      "last_used": "2026-08-09T12:00:00Z"
    }
  ]
}
```

Do not introduce SQLite unless actual scale or query requirements justify it.

For a few hundred projects, JSON is sufficient.

---

# 17. Configuration

User configuration should live at:

```text
~/.config/phub/config.toml
```

Example:

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

Configuration should have sensible defaults.

phub should work with minimal configuration.

---

# 18. Configuration Model

A possible Go structure:

```go
type Config struct {
    General GeneralConfig
    Scan    ScanConfig
}

type GeneralConfig struct {
    Editor      string
    GitUI       string
    AI          string
    FileManager string
}

type ScanConfig struct {
    Directories []string
    MaxDepth    int
}
```

Configuration parsing should happen once during startup.

The rest of the application receives validated configuration.

---

# 19. Launcher

The launcher starts external development tools in the selected project.

Supported initial actions:

```text
shell
nvim
lazygit
forge
yazi
```

Conceptually:

```go
type Launcher interface {
    Launch(ctx context.Context, tool Tool, project project.Project) error
}
```

---

# 20. Tool Configuration

Tool names must not be hardcoded throughout the application.

For example:

```toml
[general]

editor = "nvim"
git_ui = "lazygit"
ai = "forge"
file_manager = "yazi"
```

The UI asks:

```text
Launch editor
```

The launcher determines:

```text
editor → nvim
```

This keeps the UI independent from specific external tools.

---

# 21. Process Execution

Never invoke user-configured tools through a shell unless absolutely necessary.

Preferred:

```go
exec.Command(tool, args...)
```

Avoid:

```go
exec.Command("sh", "-c", command)
```

This prevents project paths and configuration values from becoming shell syntax.

---

# 22. Working Directory

External tools should launch with the selected project as their working directory.

Conceptually:

```go
cmd := exec.Command("nvim", ".")
cmd.Dir = project.Path
```

This is preferable to:

```text
cd /project && nvim
```

---

# 23. TUI Architecture

phub should use Bubble Tea's model:

```text
Model
  ↓
Update(Message)
  ↓
New Model
  ↓
View()
```

The UI should remain a presentation layer.

---

# 24. UI State

UI-specific state may include:

```go
type Model struct {
    projects []project.Project

    selected int
    query    string

    searching bool
    loading   bool

    width  int
    height int
}
```

Do not put persistent project data directly into UI-only structures when it belongs in the project model.

---

# 25. Main Screen

The initial screen should prioritize project selection.

Example:

```text
┌─ phub ────────────────────────────────────────────────────┐
│ Search projects...                                       │
├──────────────────────────────────────────────────────────┤
│                                                         │
│ ★ Forge          Python      main      ● 3 changes      │
│   phub           Go          main      ✓ clean          │
│   ScreenBot      Python      main      ✓ clean          │
│                                                         │
├──────────────────────────────────────────────────────────┤
│ enter open · n nvim · g lazygit · f forge · ? help     │
└──────────────────────────────────────────────────────────┘
```

The interface should not require permanent sidebars.

---

# 26. Search

Project search should operate over the in-memory project list.

Initial searchable fields:

```text
Name
Path
Language
```

Later:

```text
Tags
Branch
Favorite
```

A project collection is small enough that local filtering should be effectively instant.

A sophisticated search index is unnecessary.

---

# 27. Asynchronous Metadata

Some project information may be slow:

* Git status
* Ahead/behind information
* Environment detection
* Project health

The TUI should not freeze while collecting it.

Use Go concurrency where it improves responsiveness.

Example:

```text
Load basic project list
        ↓
Render immediately
        ↓
Fetch Git metadata concurrently
        ↓
Send Bubble Tea messages
        ↓
Update rows
```

---

# 28. Concurrency Limits

Do not launch hundreds of Git processes simultaneously.

Use bounded concurrency.

Conceptually:

```text
Projects
   ↓
Worker pool
   ↓
Git inspection
```

A small worker count such as:

```text
4–8
```

is sufficient initially.

The exact value should be configurable or determined later through measurement.

---

# 29. Caching

Early versions need minimal caching.

Persistent information:

```text
Project path
Favorite
Last used
Manual registration
```

Transient information:

```text
Git status
Environment state
```

Transient metadata can initially be refreshed when phub starts.

Add caching only when measurement shows it is necessary.

---

# 30. Recent Projects

When a project is opened through phub:

```text
Registry.MarkUsed(project)
```

updates its timestamp.

Recent projects can then be sorted by:

```text
LastUsed descending
```

This should not depend on filesystem modification timestamps.

The question is:

> When did the user work on this project through phub?

not:

> When was some file modified?

---

# 31. Favorites

Favorites are phub metadata.

They must not modify the project.

Example:

```text
★ Forge
★ phub
  Website
```

Favorite state belongs in phub's persistent storage.

---

# 32. Error Handling

Expected failures should not crash the TUI.

Examples:

```text
Project directory disappeared
Git is not installed
Lazygit is not installed
Forge is not installed
Configuration is invalid
Permission denied
```

Display actionable messages.

Example:

```text
Could not launch Lazygit.

Executable "lazygit" was not found in PATH.
```

Avoid:

```text
error: exit status 1
```

when more useful information is available.

---

# 33. Missing External Tools

phub should detect missing configured tools.

Example project actions:

```text
n  Neovim
g  Lazygit
f  Forge
y  Yazi
```

If Forge is unavailable:

```text
f  Forge [not installed]
```

phub should not automatically install it.

---

# 34. CLI Architecture

The CLI and TUI should share the same core services.

For example:

```bash
phub list
```

and the TUI should both use the same project registry.

Do not implement separate discovery systems for CLI and TUI.

Conceptually:

```text
                  Project Registry
                    ↑         ↑
                    │         │
                   CLI       TUI
```

---

# 35. Planned CLI Commands

Initial:

```text
phub
phub PROJECT
phub list
phub scan
phub add PATH
phub remove PROJECT
phub recent
phub --help
phub --version
```

Later:

```text
phub doctor
phub favorites
phub open PROJECT --with TOOL
```

The TUI remains the primary interface.

---

# 36. Dependency Direction

Dependencies should flow inward toward simple core models.

Preferred:

```text
UI ────────────┐
CLI ───────────┤
               ▼
              App
               │
       ┌───────┼────────┐
       ▼       ▼        ▼
   Registry Discovery Launcher
       │       │
       └───┬───┘
           ▼
        Project
```

Avoid:

```text
Project → UI
Discovery → Bubble Tea
Git → UI widgets
Storage → Launcher
```

Core packages should not know how the terminal interface is rendered.

---

# 37. Testing Strategy

phub should be easy to test without launching an interactive terminal.

## 37.1 Project Detection Tests

Test detection for:

* Git repositories
* Python projects
* Go projects
* Rust projects
* JavaScript / TypeScript projects
* Java projects
* C / C++ projects
* Mixed-language repositories
* Directories that are not projects

Use temporary directories.

Do not depend on the contributor's actual filesystem.

---

## 37.2 Discovery Tests

Test:

* Multiple configured roots
* Nested projects
* Maximum scan depth
* Ignored directories
* Missing scan directories
* Permission-denied directories
* Symlinks
* Duplicate projects
* Manually registered projects

Discovery should behave deterministically.

---

## 37.3 Git Tests

Git inspection tests should use temporary repositories.

Test:

* Clean repository
* Modified file
* Untracked file
* Detached HEAD
* Repository with no commits
* Non-Git directory
* Missing Git executable
* Branch names with unusual characters

Do not modify the user's real Git configuration.

---

## 37.4 Storage Tests

Test:

* Saving projects
* Loading projects
* Favorites
* Recent timestamps
* Duplicate prevention
* Missing storage file
* Corrupted storage file
* Atomic writes

Corrupted local phub data should not prevent the application from starting.

phub should recover gracefully where possible.

---

## 37.5 Launcher Tests

Do not open real interactive tools during ordinary tests.

Abstract process launching so tests can inspect commands.

Example:

```go
type CommandRunner interface {
    Run(Command) error
}
```

Tests should verify:

```text
Tool: nvim
Directory: /tmp/project
Args: ["."]
```

rather than launching Neovim.

---

## 37.6 UI Tests

UI behavior worth testing includes:

* Navigation
* Search
* Selection
* Favorites
* Empty project list
* Loading states
* Error messages
* Small terminal sizes
* Missing tools

Core behavior should still be testable independently of Bubble Tea.

---

# 38. State Ownership

Each piece of state should have one clear owner.

Examples:

```text
Config             → config package
Known projects     → registry
Favorites          → registry/storage
Recent usage       → registry/storage
Git metadata       → project metadata
Search query       → UI
Selected row       → UI
Window dimensions  → UI
```

Avoid multiple packages independently maintaining copies of the same state.

---

# 39. Application Layer

The application layer coordinates core services.

Possible structure:

```go
type App struct {
    Registry   Registry
    Discovery  DiscoveryService
    Launcher   Launcher
    Config     Config
}
```

The application layer may expose operations such as:

```go
func (a *App) Scan(ctx context.Context) error
func (a *App) Projects() []project.Project
func (a *App) Open(projectID string, tool Tool) error
func (a *App) ToggleFavorite(projectID string) error
```

The UI calls these operations instead of reaching directly into storage or discovery internals.

---

# 40. Project Sorting

phub should support predictable project ordering.

Initial default:

```text
Favorites first
then
Recently used
then
Alphabetical
```

Possible future sort modes:

```text
Name
Last used
Language
Git status
Path
```

Sorting logic belongs outside rendering code.

---

# 41. Project Filtering

Filtering should operate on the in-memory project collection.

Initial filtering:

```text
Name contains query
Path contains query
Language contains query
```

Matching should be case-insensitive.

Future fuzzy matching may use a small dedicated library if necessary.

Do not add a heavyweight search engine.

---

# 42. Startup Flow

A normal startup should look like:

```text
Load configuration
      ↓
Load project registry
      ↓
Render known projects
      ↓
Refresh discovery
      ↓
Refresh project metadata
      ↓
Update UI incrementally
```

The user should not have to wait for every repository to be inspected before seeing the interface.

---

# 43. First Run

If no configuration or project database exists:

```text
Welcome to phub.

No projects found yet.

Scan:
  ~/Projects
  ~/Code
  ~/Python

[Enter] Scan
[a] Add directory
[q] Quit
```

phub should provide useful defaults without forcing the user to manually create TOML files.

---

# 44. Refresh Behavior

Users should be able to refresh project information.

Potential key:

```text
R
```

Refresh may update:

* Project discovery
* Git status
* Tool availability
* Environment information

It should not modify projects.

---

# 45. Project Removal

Removing a project from phub means:

> Forget this project from phub.

It must **not** mean:

> Delete the directory.

Example confirmation:

```text
Remove Forge from phub?

The project files will NOT be deleted.

[y/N]
```

Actual project deletion should not exist in early versions.

---

# 46. Manual Projects

Users may add projects outside configured scan roots.

Example:

```bash
phub add ~/Desktop/special-project
```

Manual projects should remain registered even if they are outside normal scan directories.

If the directory disappears, phub should mark it unavailable rather than silently forgetting it.

---

# 47. Project Availability

A project may have a status such as:

```go
type Availability int

const (
    Available Availability = iota
    Missing
    Inaccessible
)
```

The UI may display:

```text
Forge        ✓
OldProject   missing
PrivateRepo  permission denied
```

This is better than crashing or silently dropping projects.

---

# 48. External Tool Detection

At startup or lazily, phub may inspect `PATH` for configured tools.

Example:

```go
exec.LookPath("nvim")
```

Possible status:

```text
Neovim     ✓
Lazygit    ✓
Forge      ✓
Yazi       ✗
```

Tool detection should be cached for the session.

---

# 49. Tool Launch Semantics

Different tools behave differently.

### Editor

```text
nvim .
```

### Lazygit

```text
lazygit
```

with project directory as working directory.

### Yazi

```text
yazi .
```

### Forge

```text
forge
```

with project directory as working directory.

### Shell

Start the user's configured shell in the project directory.

Tool-specific arguments belong in the launcher configuration, not scattered throughout the UI.

---

# 50. Terminal Suspension

When opening an interactive external program from the TUI, phub may need to suspend Bubble Tea temporarily.

The intended lifecycle:

```text
phub TUI
   ↓
Suspend
   ↓
Run nvim / lazygit / forge / shell
   ↓
Tool exits
   ↓
Resume phub
   ↓
Refresh project metadata
```

This should feel seamless.

---

# 51. Terminal Compatibility

phub should work in common terminal emulators.

Do not depend on:

* Kitty-only features
* Alacritty-only behavior
* Nerd Fonts
* Truecolor being available
* Mouse support

Nerd Font icons may be optional.

Important information must remain understandable in plain ASCII.

---

# 52. Styling

phub should have a restrained visual style.

Use styling for meaning:

```text
Selected      → emphasized
Favorite      → star
Clean Git     → success indication
Dirty Git     → warning indication
Error         → error indication
Disabled tool → dimmed
```

Avoid excessive gradients, animations, borders, and decorative elements.

Speed and clarity matter more than spectacle.

---

# 53. Responsive Layout

phub must remain usable in smaller terminals.

Large:

```text
PROJECT      LANG       BRANCH     STATUS       LAST USED
Forge        Python     main       3 changes    now
```

Small:

```text
Forge
Python · main · dirty
```

The layout should adapt rather than overflow.

---

# 54. No-Color Mode

phub should respect environments such as:

```text
NO_COLOR=1
```

The interface must remain understandable without color.

Do not use color as the only indication of:

* Selection
* Error
* Dirty repository
* Disabled action

---

# 55. Performance Targets

Initial performance goals:

```text
Startup to first render: effectively instant
Search filtering: immediate
Navigation: no visible lag
Project list: hundreds of projects comfortably
Git refresh: asynchronous
```

Do not optimize based on guesses.

Measure first.

---

# 56. Scaling Expectations

phub is designed primarily for:

```text
10–500 projects
```

It does not initially need to optimize for:

```text
100,000 repositories
```

Architecture should reflect realistic usage.

Simple data structures are sufficient.

---

# 57. Filesystem Watchers

Do not add filesystem watchers in the first versions.

A manual refresh and startup refresh are sufficient.

Watchers add:

* Platform differences
* Resource consumption
* Complexity
* Event storms
* Harder tests

They may be added later if real usage demonstrates the need.

---

# 58. Network Access

Core phub functionality requires no network access.

Initial features should not automatically contact:

* GitHub
* GitLab
* Package registries
* Update servers
* Analytics systems

Future remote integrations must be explicit and optional.

---

# 59. Privacy

phub may know:

* Project names
* Project paths
* Languages
* Git state
* Recently opened projects

This information should stay local by default.

Telemetry should not be required.

If telemetry is ever introduced, it must be opt-in and documented.

---

# 60. Security

phub has a much smaller attack surface than Forge, but several rules still matter.

## Never Construct Shell Commands from Project Paths

Bad:

```go
exec.Command("sh", "-c", "cd "+project.Path+" && nvim")
```

Good:

```go
cmd := exec.Command("nvim", ".")
cmd.Dir = project.Path
```

---

## Treat Project Metadata as Untrusted

A project may have unusual or malicious:

* File names
* Branch names
* Paths
* Configuration values

These values must be rendered as data, not interpreted as terminal commands.

---

## Avoid Arbitrary Automatic Execution

Detecting:

```text
package.json
```

must not cause phub to automatically execute:

```text
npm install
npm run dev
```

Running project-defined commands should always be an explicit user action.

---

# 61. Project-Specific Configuration

A future `.phub.toml` may define project actions.

Example:

```toml
name = "Forge"

[commands]
run = ["python", "-m", "forge"]
test = ["pytest"]
```

If implemented:

* Parse arguments as arrays.
* Do not interpret commands through a shell by default.
* Do not automatically run them.
* Show commands clearly before state-changing actions.

This feature is not necessary for v0.1.

---

# 62. Data Migration

Persistent data formats may change as phub evolves.

Storage should include a format version.

Example:

```json
{
  "version": 1,
  "projects": []
}
```

If the format changes:

```text
Old format
   ↓
Migration
   ↓
New format
```

Do not silently discard user favorites or recent-project history.

---

# 63. Atomic Storage

Persistent metadata should use atomic replacement.

Recommended:

```text
Write temporary file
      ↓
Flush
      ↓
Rename into place
```

This prevents a crash from leaving a half-written project database.

---

# 64. Logging

phub should not produce noisy log files during ordinary usage.

Debug mode may provide structured diagnostics.

Example:

```bash
phub --debug
```

Useful information:

```text
Loaded config
Scanned 3 roots
Found 27 projects
Git metadata refreshed in 241ms
Could not inspect one directory: permission denied
```

Do not log unnecessary project file contents.

---

# 65. Doctor Command

A later command:

```bash
phub doctor
```

may inspect:

```text
✓ config directory
✓ data directory
✓ Git
✓ Neovim
✓ Lazygit
✓ Forge
✓ Yazi
✓ configured scan roots
```

Example:

```text
phub doctor

Git        ✓ /usr/bin/git
Neovim     ✓ /usr/bin/nvim
Lazygit    ✓ /usr/bin/lazygit
Forge      ✓ /usr/local/bin/forge
Yazi       ✗ not found

Projects
~/Projects ✓
~/Python   ✓
~/Code     ✗ directory does not exist
```

This helps diagnose setup problems without putting diagnostic logic into the TUI.

---

# 66. CLI Exit Codes

CLI commands should use meaningful exit status codes.

```text
0   success
1   ordinary failure
2   invalid command or arguments
```

Do not invent many custom exit codes until there is a real need.

---

# 67. Versioning

Before stable release:

```text
0.x.y
```

Suggested progression:

```text
0.0.1   TUI prototype
0.0.2   discovery
0.0.3   Git metadata
0.0.4   search
0.0.5   launching
0.0.6   integrations
0.0.7   favorites/recent
0.1.0   first usable release
```

Do not inflate version numbers merely because the codebase is large.

---

# 68. Dependencies

Keep dependencies small.

Likely initial dependencies:

```text
Bubble Tea
Bubbles
Lip Gloss
TOML parser
```

Potential fuzzy-search dependency may be added if justified.

Do not add dependencies for functionality easily provided by the Go standard library.

---

# 69. Error Model

Errors should be categorized enough to provide useful feedback.

Examples:

```go
ErrProjectNotFound
ErrToolNotFound
ErrInvalidConfig
ErrPermissionDenied
ErrStorageCorrupt
```

Not every error needs a custom type.

Use custom errors where the UI or CLI needs different behavior.

---

# 70. Context Cancellation

Long-running operations should accept `context.Context`.

Examples:

```go
Scan(ctx)
InspectGit(ctx)
Launch(ctx)
```

This allows:

* User cancellation
* Clean shutdown
* Timeouts
* Future asynchronous work

---

# 71. Graceful Shutdown

When phub exits:

* Cancel active scans.
* Stop metadata workers.
* Finish or safely abandon local state writes.
* Restore terminal state.
* Do not leave background child processes accidentally running.

Ctrl+C should exit cleanly.

---

# 72. Signal Handling

The CLI entry point may handle:

```text
SIGINT
SIGTERM
```

and cancel the root context.

Bubble Tea should remain responsible for terminal restoration.

---

# 73. Architecture Evolution

Do not preserve this document blindly if real usage proves part of the design wrong.

When changing a major architecture choice:

1. Identify the current limitation.
2. Propose the simpler alternative.
3. Consider migration.
4. Update tests.
5. Update this document.

Architecture exists to help the project, not restrict necessary improvement.

---

# 74. Explicit Early Non-Goals

Until the core project-management workflow is excellent, do not build:

* AI project recommendations
* Cloud project sync
* Built-in source editor
* Built-in Git client
* Built-in file manager
* Package installer
* Docker dashboard
* SSH manager
* Plugin marketplace
* Semantic code indexing
* Project collaboration
* Background daemon
* Web UI
* Automatic project repair
* Automatic dependency installation

Use dedicated tools for those jobs.

---

# 75. v0.1 Architectural Boundary

For v0.1, phub needs only:

```text
CLI
TUI
Configuration
Project discovery
Project registry
Language detection
Basic Git metadata
Search
Favorites
Recent projects
Tool launcher
Local storage
```

Everything else is optional.

---

# 76. Success Criteria

The architecture is successful when the following workflow is simple and reliable:

```text
$ phub
```

phub opens quickly.

The user sees their projects.

They type:

```text
for
```

Forge becomes selected.

They press:

```text
n
```

Neovim opens in the Forge project.

They exit Neovim.

phub returns.

Git status refreshes.

They press:

```text
g
```

Lazygit opens in the same project.

No manual `cd`.

No searching through folders.

No complicated setup.

That is the core phub experience.

---

# 77. Architectural Priority

When implementation choices conflict, use this priority:

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

phub should remain small enough that one developer can understand the whole application.

---

# 78. Final Principle

phub should not try to become the tool that does everything.

It should become the tool that gets you to the right project and the right tool immediately.

> **Find it. Open it. Work.**

