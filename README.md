# WhyTie

**Git records what changed. WhyTie remembers why.**

WhyTie is a local-first CLI that preserves the reasoning behind your code.

Write short comments next to the implementation. Scan your project. Later, retrieve decisions and their reasons, review their history, or use `whytie resume` to revisit recently recorded context.

No AI, account, or server is required.

```go
// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}
```

## Status

This source version is WhyTie v1.2.1.

See [GitHub Releases](https://github.com/JumoKookbob/whytie/releases) for published versions and downloadable packages.

Changes in v1.2.1 include:

- Unified reasoning output across `scan`, `list`, `why`, `history`, and `resume`.
- A consistent file-oriented layout with line numbers, icons, labels, and separators.
- Multiple reasoning blocks from the same source file are grouped under a single file header in `scan`.
- `history` displays the current reasoning entry using the shared WhyTie output style before its recorded history.

WhyTie is under active development. Automated tests and manual checks cover creation, editing, movement, deletion, persistence, retrieval, block-comment scanning, and CLI behavior.

## What WhyTie does

- Captures questions, decisions, reasons, rejected approaches, failed approaches, and important context from source comments.
- Stores reasoning locally in SQLite.
- Groups adjacent annotations for readable context.
- Retrieves reasoning by memory ID or source location.
- Preserves reasoning identity across recognized edits and moves.
- Records created, moved, changed, and deleted events.
- Keeps historical snapshots after annotations disappear.
- Attaches Git provenance when available.
- Shows recent recorded events and related context with `resume`.
- Presents reasoning using a consistent file-oriented CLI layout.

WhyTie retrieves what you recorded. It does not invent explanations or infer your next task.

## Installation

### Release packages

The v1.2.1 package targets and archive names are:

| Platform                 | Package name                      |
| ------------------------ | --------------------------------- |
| Windows Intel/AMD 64-bit | `whytie-v1.2.1-windows-amd64.zip` |
| Windows ARM64            | `whytie-v1.2.1-windows-arm64.zip` |
| Linux x86-64             | `whytie-v1.2.1-linux-amd64.zip`   |
| Linux ARM64              | `whytie-v1.2.1-linux-arm64.zip`   |
| macOS Intel              | `whytie-v1.2.1-darwin-amd64.zip`  |
| macOS Apple Silicon      | `whytie-v1.2.1-darwin-arm64.zip`  |

`darwin` means macOS.

Check [GitHub Releases](https://github.com/JumoKookbob/whytie/releases) for available downloads. Prebuilt packages do not require Go. Git is optional and enables Git provenance when available.

Extract the entire archive, including its documentation and `third_party` notices.

### Windows

Extract the Windows package, then run:

```powershell
.\install.cmd
```

The installer copies `whytie.exe` to:

```text
%LOCALAPPDATA%\WhyTie\bin
```

and adds that directory to your user `PATH`.

After installation, open a new terminal and verify:

```powershell
whytie version
whytie help
```

Then switch to a project directory:

```powershell
whytie init
whytie scan .
whytie list
whytie resume
```

No manual alias is required.

To uninstall WhyTie from the extracted package directory:

```powershell
.\uninstall.cmd
```

The uninstaller removes the installed executable and its user `PATH` entry. Existing project `.whytie` directories are not deleted.

The `.cmd` wrappers invoke the bundled PowerShell scripts for the installation or removal process without requiring the user to permanently change their PowerShell execution policy.

### Linux and macOS

Open a terminal in the extracted package directory:

```bash
chmod +x ./whytie
./whytie version
```

You can run WhyTie directly from that directory or move the executable to a directory on your `PATH`.

For example:

```bash
mkdir -p ~/.local/bin
cp ./whytie ~/.local/bin/whytie
```

Make sure `~/.local/bin` is on your `PATH`.

## Quick start

Move to the root of a project:

```powershell
cd path\to\your-project
```

Initialize WhyTie:

```powershell
whytie init
```

Add reasoning next to code:

```go
// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}
```

Scan the project:

```powershell
whytie scan .
```

WhyTie displays the reasoning it found:

```text
📄 shutdown.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    1  ❓ QUESTION  How should shutdown handle pending writes?

    2  ✅ DECISION  Flush pending writes before shutdown
    3  └─ 💡 REASON    Prevent user data loss
```

List saved reasoning:

```powershell
whytie list
```

Inspect a specific annotation:

```powershell
whytie why shutdown.go:2
```

Review its history:

```powershell
whytie history shutdown.go:2
```

Return to recently recorded context:

```powershell
whytie resume
```

## Reasoning syntax

WhyTie uses small markers inside source comments.

| Marker | Meaning           |
| ------ | ----------------- |
| `?`    | Question          |
| `+`    | Decision          |
| `x`    | Rejected approach |
| `f`    | Failed approach   |
| `<`    | Reason            |
| `!`    | Important context |

Example:

```go
// ? Which database should we use?
// + Use SQLite
// < It fits the local-first design

// x Run PostgreSQL locally
// < Adds infrastructure the tool does not need

// ! Keep the reasoning database out of Git
```

The markers are intentionally small so the source remains readable without WhyTie.

### Block comments

Supported slash-comment languages can also use multiline block comments:

```go
/*
? Which database should we use?
+ Use SQLite
< It fits the local-first design
*/
```

JSDoc-style blocks are also supported when they contain WhyTie markers:

```javascript
/**
 * ? Where should browser interaction live?
 * + Keep UI behavior in a separate JavaScript file
 * < Separating behavior keeps the HTML easier to inspect
 */
```

Ordinary block or JSDoc comments without WhyTie markers are ignored.

Comment-like text inside supported string and template-literal contexts is not treated as WhyTie reasoning by the scanner behavior covered by the current tests.

### Language-specific prefixes

| Source files                                                    | Annotation form                                  |
| --------------------------------------------------------------- | ------------------------------------------------ |
| Go, Rust, C/C++, C#, Java, Kotlin, JavaScript/TypeScript, Swift | `// ? Question` and supported `/* ... */` blocks |
| Python, Ruby, shell scripts, PowerShell, YAML                   | `# ? Question`                                   |
| CSS, SCSS, Sass, Less                                           | `/* ? Question */`                               |
| HTML                                                            | `<!-- ? Question -->`                            |

Use standalone comment lines. Trailing annotations after code are not currently captured.

Language support is lightweight scanning, not complete parsing of every supported language.

## Commands

```text
whytie init
whytie scan <path>
whytie list
whytie why <memory-id|file:line>
whytie history <memory-id|file:line>
whytie resume
whytie version
whytie help
whytie --help
whytie -h
```

Running `whytie` without a command also displays help.

### Initialize

```powershell
whytie init
```

Creates the local `.whytie` directory. Its internal `.gitignore` keeps local database files out of Git.

### Scan

From the project root:

```powershell
whytie scan .
```

Finds annotations and synchronizes the saved reasoning.

Run a scan after editing, moving, or removing annotations. There is no automatic background scan.

Generated and dependency directories such as `.git`, `.whytie`, `node_modules`, `vendor`, `dist`, `build`, and `target` are skipped. Filenames containing `.min.` are also skipped during directory scans.

Reasoning is displayed using the shared WhyTie file-oriented layout. When a source file contains multiple separate reasoning blocks, `scan` groups them under one file header and separates the blocks visually:

```text
📄 calculator.py
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    1  ❓ QUESTION  How should calculations be implemented?

    2  ✅ DECISION  Use small explicit functions
    3  └─ 💡 REASON    Simple functions make this language test easy to understand

────────────────────────────────────────────────────

    9  ❌ FAILED    Avoid hidden global state
   10  └─ 💡 REASON    The calculation does not need persistent state
```

### List

```powershell
whytie list
```

Displays currently stored reasoning.

The output is organized around source files and uses the same line numbers, icons, labels, and visual separators as other reasoning-oriented WhyTie commands.

### Why

```powershell
whytie why shutdown.go:2
whytie why <memory-id>
```

Displays a stored annotation and its associated context using the shared WhyTie reasoning layout.

A question can lead to its adjacent decision and reasons. A decision can lead to its adjacent reasons.

Use the source location shown by WhyTie. Quote paths containing spaces:

```powershell
whytie why "src/my file.go:10"
```

### History

```powershell
whytie history shutdown.go:2
whytie history <memory-id>
```

Displays the current reasoning entry followed by its recorded lifecycle events.

Example:

```text
📄 calculator.py
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    2  ✅ DECISION  Use small explicit functions

History
────────────────────────────────────────────────────
  CREATED  calculator.py:2
```

Recorded lifecycle events include:

- `created`: an annotation was first stored.
- `changed`: its text or kind changed.
- `moved`: its stored file or line changed.
- `deleted`: it was no longer found during synchronization.

Events may include Git commit information when available.

History records changes observed by scans. Changes made and undone between scans are not recorded.

### Resume

```powershell
whytie resume
```

Shows:

- The five most recently saved history events.
- Event snapshots, recording times, and available commit information.
- Related context from currently stored annotations.
- An `Inspect` command for each displayed context block.

Related current context uses the same file-oriented reasoning layout as `scan`, `list`, and `why`.

Related context blocks are displayed once, even when several recent events belong to the same block.

Historical event snapshots and currently stored context are separate. Deleted annotations can appear in the event list, but are not shown as current context.

Resume reads saved data. It does not rescan source files, determine what you last worked on, or mark questions as resolved or unresolved.

"Newest" means most recently inserted into the local history. Within one scan, this does not establish the order in which you actually edited the code.

If saved reasoning exists without history, Resume directs you to `whytie list`.

### Help

```powershell
whytie help
whytie --help
whytie -h
```

Displays command usage and the annotation marker reference.

Running `whytie` with no arguments displays the same help information.

### Version

```powershell
whytie version
```

Displays the version embedded in the executable.

## Reasoning survives recognized changes

After editing or moving annotations, scan again:

```powershell
whytie scan .
```

When reconciliation recognizes the same reasoning, its memory ID is retained and the change is added to its history.

For example, a reason can accumulate lifecycle events:

```text
History
────────────────────────────────────────────────────
  CREATED  shutdown.go:5
  CHANGED  shutdown.go:5
  MOVED    lifecycle.go:5
  DELETED  lifecycle.go:5
```

After deletion, retrieve its recorded history using a historical location:

```powershell
whytie history lifecycle.go:5
```

This displays preserved reasoning; it does not restore source code.

Deleted-history lookup by location is supported. The current ID-based history command requires a memory still present in the current store.

Identity matching is best-effort. Arbitrary refactoring, duplicate annotations, or ambiguous edits may not preserve the intended match.

## Local storage and privacy

Project data is stored in:

```text
.whytie/whytie.db
```

WhyTie uses SQLite and operates locally. Its core workflow requires no account, hosted service, LLM, embeddings, or vector database.

Source code and reasoning are not uploaded by the core workflow.

Source comments remain readable without WhyTie. The local database preserves identities and scan history, including deleted annotations.

The database is not committed to Git by default. A fresh clone does not automatically include its history. Keep a backup of `.whytie` if you need to preserve that local history.

## Git integration

WhyTie attaches Git provenance where available and can operate outside a Git repository.

A displayed commit is provenance attached to the event. It is not proof that the commit caused or validated the recorded decision.

WhyTie does not automatically reconstruct all annotations from the repository's past commits.

## Current limitations

- Scanning is explicit; source changes are not captured automatically.
- Resume shows saved context, not a complete work-session summary.
- Questions have no tracked resolution status.
- Relationships primarily rely on adjacent source annotations.
- Language scanners do not fully model every language's syntax.
- Identity preservation is best-effort, not a guarantee for all refactors.
- Historical locations can become ambiguous when reused by other annotations.
- Local history is not automatically synchronized between machines.
- There is no dedicated export command yet.
- There is no Change Guard, impact analysis, GUI, or editor integration.
- Release packages are not publisher-signed; macOS packages are not notarized.

## Development

Run tests:

```powershell
go test ./...
```

Build on Windows:

```powershell
go build -o .\whytie.exe .\cmd\whytie
```

Run directly from source:

```powershell
go run ./cmd/whytie resume
```

## License

WhyTie is licensed under the Apache License, Version 2.0.

See [LICENSE](LICENSE) for the full terms.

Third-party license and attribution files are provided in [third_party](third_party/README.md) and included in release packages.
