# WhyTie

**Git records what changed. WhyTie remembers why.**

WhyTie is a local-first CLI for preserving the reasoning behind your code.

Git is excellent at showing **what changed**. But months later, the harder question is often:

> Why did we make this decision?

WhyTie keeps that context close to the code using small, readable annotations:

```go
// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
// x Drop pending writes on exit
// < Fast, but unsafe if work is still buffered

func Shutdown() {}
```

Then scan your project:

```powershell
whytie scan .
```

And come back later with:

```powershell
whytie resume
```

WhyTie stores the reasoning locally, keeps its history, and lets you retrieve it by file location or memory ID.

**No AI, account, hosted service, or server is required.**

---

## Why WhyTie?

Code usually tells you **what it does**.

Git usually tells you **what changed**.

Neither reliably tells you **why a non-obvious decision was made**.

That missing context becomes expensive when you:

- return to a side project after months away,
- revisit a workaround that now looks strange,
- forget why one dependency was chosen over another,
- re-open an architectural decision that was already debated,
- use AI coding tools that can produce code faster than you can remember the reasoning behind it.

WhyTie is designed to preserve that reasoning without requiring a separate knowledge base.

The annotations remain ordinary source comments, so the code stays understandable even without WhyTie.

---

## A 30-second example

Add reasoning next to the implementation:

```go
// ? Which database should we use?
// + Use SQLite
// < The tool should remain local-first with no database server
// x PostgreSQL
// < Adds infrastructure that the current product does not need
```

Scan the project:

```powershell
whytie scan .
```

WhyTie displays the reasoning it found:

```text
📄 storage.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

   12  ❓ QUESTION  Which database should we use?
   13  ✅ DECISION  Use SQLite
   14  └─ 💡 REASON    The tool should remain local-first with no database server
   15  ❌ REJECTED  PostgreSQL
   16  └─ 💡 REASON    Adds infrastructure that the current product does not need
```

Later:

```powershell
whytie why storage.go:13
whytie history storage.go:13
whytie resume
```

---

## WhyTie and AI coding agents

WhyTie does not require AI, but it can work well alongside AI coding agents.

An agent can leave WhyTie annotations when it makes a non-trivial implementation decision, so the reasoning stays close to the code instead of disappearing with the chat session.

For example:

```go
// + Keep SQLite for local storage
// < The application should remain local-first with no database server
// x PostgreSQL
// < Adds infrastructure the current product does not need
```

This gives future developers — and future AI sessions — a compact record of why the implementation looks the way it does.

The repository includes [`AGENTS.md`](AGENTS.md) with the WhyTie annotation rules for AI coding agents, including what should be recorded, what should be ignored, and how existing WhyTie annotations should be handled.

WhyTie does not verify whether an AI-generated explanation is correct. It stores the reasoning that was written into the source.

---

## Status

This source version is **WhyTie v1.2.1**.

See [GitHub Releases](https://github.com/JumoKookbob/whytie/releases) for published versions and downloadable packages.

Changes in v1.2.1 include:

- Unified reasoning output across `scan`, `list`, `why`, `history`, and `resume`.
- A consistent file-oriented layout with line numbers, icons, labels, and separators.
- Multiple reasoning blocks from the same source file are grouped under a single file header in `scan`.
- `history` displays the current reasoning entry using the shared WhyTie output style before its recorded history.

WhyTie is under active development. Automated tests and manual checks cover creation, editing, movement, deletion, persistence, retrieval, block-comment scanning, and CLI behavior.

---

## What WhyTie does

WhyTie can:

- capture questions, decisions, reasons, rejected approaches, failed approaches, and important context from source comments,
- store reasoning locally in SQLite,
- group adjacent annotations into readable reasoning blocks,
- retrieve reasoning by memory ID or source location,
- preserve reasoning identity across recognized edits and moves,
- record `created`, `moved`, `changed`, and `deleted` lifecycle events,
- keep historical snapshots after annotations disappear,
- attach Git provenance when available,
- show recent recorded events and related context with `whytie resume`,
- present reasoning using a consistent file-oriented CLI layout.

WhyTie retrieves what you recorded. It does **not** invent explanations or infer your next task.

---

## Installation

### Release packages

The v1.2.1 package targets and archive names are:

| Platform | Package name |
| --- | --- |
| Windows Intel/AMD 64-bit | `whytie-v1.2.1-windows-amd64.zip` |
| Windows ARM64 | `whytie-v1.2.1-windows-arm64.zip` |
| Linux x86-64 | `whytie-v1.2.1-linux-amd64.zip` |
| Linux ARM64 | `whytie-v1.2.1-linux-arm64.zip` |
| macOS Intel | `whytie-v1.2.1-darwin-amd64.zip` |
| macOS Apple Silicon | `whytie-v1.2.1-darwin-arm64.zip` |

`darwin` means macOS.

Check [GitHub Releases](https://github.com/JumoKookbob/whytie/releases) for available downloads.

Prebuilt packages do not require Go. Git is optional and enables Git provenance when available.

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

Open a new terminal and verify:

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

The `.cmd` wrappers invoke the bundled PowerShell scripts for installation or removal without requiring the user to permanently change their PowerShell execution policy.

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

---

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

---

## Reasoning syntax

WhyTie uses small markers inside source comments.

| Marker | Meaning |
| --- | --- |
| `?` | Question |
| `+` | Decision |
| `x` | Rejected approach |
| `f` | Failed approach |
| `<` | Reason |
| `!` | Important context |

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

### Rejected vs failed

Use `x` when an approach was considered and intentionally rejected:

```go
// x Run PostgreSQL locally
// < Adds infrastructure the tool does not need
```

Use `f` when an approach was actually attempted and did not work:

```go
// f Reconcile entries using line number alone
// < Refactors caused identities to drift when code moved
```

This distinction helps future readers understand whether an option was merely declined or proven problematic in practice.

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

| Source files | Annotation form |
| --- | --- |
| Go, Rust, C/C++, C#, Java, Kotlin, JavaScript/TypeScript, Swift | `// ? Question` and supported `/* ... */` blocks |
| Python, Ruby, shell scripts, PowerShell, YAML | `# ? Question` |
| CSS, SCSS, Sass, Less | `/* ? Question */` |
| HTML | `<!-- ? Question -->` |

Use standalone comment lines. Trailing annotations after code are not currently captured.

Language support is lightweight scanning, not complete parsing of every supported language.

---

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

### `whytie init`

```powershell
whytie init
```

Creates the local `.whytie` directory. Its internal `.gitignore` keeps local database files out of Git.

### `whytie scan`

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

    9  ⚠️ FAILED    Avoid hidden global state
   10  └─ 💡 REASON    The calculation does not need persistent state
```

### `whytie list`

```powershell
whytie list
```

Displays currently stored reasoning.

The output is organized around source files and uses the same line numbers, icons, labels, and visual separators as other reasoning-oriented WhyTie commands.

### `whytie why`

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

### `whytie history`

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

- `created`: an annotation was first stored,
- `changed`: its text or kind changed,
- `moved`: its stored file or line changed,
- `deleted`: it was no longer found during synchronization.

Events may include Git commit information when available.

History records changes observed by scans. Changes made and undone between scans are not recorded.

### `whytie resume`

```powershell
whytie resume
```

Shows:

- the five most recently saved history events,
- event snapshots, recording times, and available commit information,
- related context from currently stored annotations,
- an `Inspect` command for each displayed context block.

Related current context uses the same file-oriented reasoning layout as `scan`, `list`, and `why`.

Related context blocks are displayed once, even when several recent events belong to the same block.

Historical event snapshots and currently stored context are separate. Deleted annotations can appear in the event list, but are not shown as current context.

`resume` reads saved data. It does not rescan source files, determine what you last worked on, or mark questions as resolved or unresolved.

"Newest" means most recently inserted into the local history. Within one scan, this does not establish the order in which you actually edited the code.

If saved reasoning exists without history, `resume` directs you to `whytie list`.

### `whytie help`

```powershell
whytie help
whytie --help
whytie -h
```

Displays command usage and the annotation marker reference.

Running `whytie` with no arguments displays the same help information.

### `whytie version`

```powershell
whytie version
```

Displays the version embedded in the executable.

---

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

---

## Local storage and privacy

Project data is stored in:

```text
.whytie/whytie.db
```

WhyTie uses SQLite and operates locally. Its core workflow requires no account, hosted service, LLM, embeddings, or vector database.

Source code and reasoning are not uploaded by the core workflow.

Source comments remain readable without WhyTie. The local database preserves identities and scan history, including deleted annotations.

The database is not committed to Git by default. A fresh clone does not automatically include its history. Keep a backup of `.whytie` if you need to preserve that local history.

---

## Git integration

WhyTie attaches Git provenance where available and can operate outside a Git repository.

A displayed commit is provenance attached to the event. It is not proof that the commit caused or validated the recorded decision.

WhyTie does not automatically reconstruct all annotations from the repository's past commits.

---

## What WhyTie is not

WhyTie is intentionally narrow.

It is not:

- a replacement for Git,
- a general-purpose documentation platform,
- an issue tracker,
- an AI coding agent,
- an automatic architecture decision generator,
- a complete work-session recorder.

Its job is to preserve small pieces of reasoning close to the code and make them retrievable later.

---

## Current limitations

- Scanning is explicit; source changes are not captured automatically.
- `resume` shows saved context, not a complete work-session summary.
- Questions have no tracked resolution status.
- Relationships primarily rely on adjacent source annotations.
- Language scanners do not fully model every language's syntax.
- Identity preservation is best-effort, not a guarantee for all refactors.
- Historical locations can become ambiguous when reused by other annotations.
- Local history is not automatically synchronized between machines.
- There is no dedicated export command yet.
- There is no Change Guard, impact analysis, GUI, or editor integration.
- Release packages are not publisher-signed; macOS packages are not notarized.

---

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

---

## License

WhyTie is licensed under the Apache License, Version 2.0.

See [LICENSE](LICENSE) for the full terms.

Third-party license and attribution files are provided in [third_party](third_party/README.md) and included in release packages.
