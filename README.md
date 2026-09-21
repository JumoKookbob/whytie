# WhyTie

**Git records what changed. WhyTie remembers why.**

WhyTie is a local-first CLI that preserves the reasoning behind your code.

Write short comments next to the implementation. Scan your project.
Later, retrieve decisions and their reasons, review their history,
or use `whytie resume` to revisit recently recorded context.

No AI, account, or server is required.

```go
// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}
```

## Status

This source version is WhyTie v1.1.0.

See [GitHub Releases](https://github.com/JumoKookbob/whytie/releases)
for published versions and downloadable packages.

Changes since v1.0.0 include:

- Basic `whytie resume` for recent recorded events and related context.
- A fix for Go raw-string scanning.
- Updated documentation and version-command tests.
- An Apache-2.0 project license.

WhyTie is under active development. Automated tests and manual
checks cover creation, editing, movement, deletion, persistence,
and retrieval of recorded reasoning.

## What WhyTie does

- Captures questions, decisions, reasons, rejected approaches,
  failed approaches, and important context from source comments.
- Stores reasoning locally in SQLite.
- Groups adjacent annotations for readable context.
- Retrieves reasoning by memory ID or source location.
- Preserves reasoning identity across recognized edits and moves.
- Records created, moved, changed, and deleted events.
- Keeps historical snapshots after annotations disappear.
- Attaches Git provenance when available.
- Shows recent recorded events and related context with `resume`.

WhyTie retrieves what you recorded. It does not invent explanations
or infer your next task.

## Installation

### Release packages

The v1.1.0 package targets and archive names are listed below:

| Platform                 | Package name                      |
| ------------------------ | --------------------------------- |
| Windows Intel/AMD 64-bit | `whytie-v1.1.0-windows-amd64.zip` |
| Windows ARM64            | `whytie-v1.1.0-windows-arm64.zip` |
| Linux x86-64             | `whytie-v1.1.0-linux-amd64.zip`   |
| Linux ARM64              | `whytie-v1.1.0-linux-arm64.zip`   |
| macOS Intel              | `whytie-v1.1.0-darwin-amd64.zip`  |
| macOS Apple Silicon      | `whytie-v1.1.0-darwin-arm64.zip`  |

`darwin` means macOS.

Check [GitHub Releases](https://github.com/JumoKookbob/whytie/releases)
for available downloads. Each release lists its published assets.

Prebuilt packages do not require Go. Git is optional and enables
Git provenance when available.

Extract the entire archive, including its documentation and
`third_party` notices.

#### Windows

Open PowerShell in the extracted package directory:

```powershell
.\whytie.exe version
$whytieExe = (Resolve-Path .\whytie.exe).Path
```

Then switch to your project directory and run:

```powershell
& $whytieExe init
& $whytieExe scan .
& $whytieExe resume
```

#### Linux and macOS

Open a terminal in the extracted package directory:

```sh
chmod +x ./whytie
./whytie version
whytie_bin="$(pwd)/whytie"
```

Then switch to your project directory and run:

```sh
"$whytie_bin" init
"$whytie_bin" scan .
"$whytie_bin" resume
```

The executable is a command-line tool; use it from a terminal.

#### Verification status

All six targets have been cross-compiled successfully.
Windows amd64 has been executed locally.
The other five targets have not yet been runtime-tested.

These packages are not publisher-signed. The macOS packages are
not notarized.

### Build from source on Windows

Prerequisites:

- Git
- Go compatible with the `go.mod` requirement: `go 1.27.1`
- PowerShell

Clone and build the current `main` branch:

```powershell
git clone https://github.com/JumoKookbob/whytie.git
Set-Location .\whytie
go build -o .\whytie.exe .\cmd\whytie
.\whytie.exe version
```

Save the executable path for use in other projects:

```powershell
$whytieExe = (Resolve-Path .\whytie.exe).Path
```

In this PowerShell session, invoke it with:

```powershell
& $whytieExe version
```

You can also place `whytie.exe` in a directory on your `PATH`.
The examples below use `whytie` when it is available on `PATH`.

Building may require network access to download Go dependencies.
Normal scanning and retrieval operate locally.

## Quick start

Create a separate demo project using the executable path saved above:

```powershell
$demo = Join-Path $env:TEMP ("whytie-demo-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $demo | Out-Null
Set-Location $demo
```

Create a source file:

```powershell
@'
package demo

// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}
'@ | Set-Content -Encoding UTF8 .\shutdown.go
```

Initialize and scan:

```powershell
& $whytieExe init
& $whytieExe scan .
& $whytieExe list
```

Retrieve the question and its associated decision and reason:

```powershell
& $whytieExe why shutdown.go:3
```

Example output:

```text
question: How should shutdown handle pending writes?  shutdown.go:3
└─ decision: Flush pending writes before shutdown  shutdown.go:4
   └─ reason: Prevent user data loss  shutdown.go:5
```

Return to recorded context:

```powershell
& $whytieExe resume
```

Inspect the decision's history:

```powershell
& $whytieExe history shutdown.go:4
```

Git is optional for this workflow.

## Annotation syntax

Each marker must be followed by a space or tab and nonempty text.

| Go example                                  | Meaning           |
| ------------------------------------------- | ----------------- |
| `// ? How should this work?`                | Question          |
| `// + Use this approach`                    | Decision          |
| `// < It preserves existing behavior`       | Reason            |
| `// - Reject the shared-cache approach`     | Rejected approach |
| `// x The asynchronous attempt lost writes` | Failed approach   |
| `// ! Preserve shutdown ordering`           | Important context |

Both forms are accepted:

```go
//? How should this work?
// ? How should this work?
```

This form is not accepted:

```go
//?How should this work?
```

Keep related annotations on consecutive lines in the same file:

```go
// ? Which storage should we use?
// + Use SQLite
// < Keep data local
// < Avoid requiring a separate server
```

The current grouping rules use adjacency. They do not infer semantic
relationships between distant comments.

For `why`, place a decision immediately after its question,
and reasons immediately after the decision.

### Language-specific prefixes

| Source files                                                    | Annotation form       |
| --------------------------------------------------------------- | --------------------- |
| Go, Rust, C/C++, C#, Java, Kotlin, JavaScript/TypeScript, Swift | `// ? Question`       |
| Python, Ruby, shell scripts, PowerShell, YAML                   | `# ? Question`        |
| CSS, SCSS, Sass, Less                                           | `/* ? Question */`    |
| HTML                                                            | `<!-- ? Question -->` |

Block-style annotations shown above must open and close on the same
line. Use standalone comment lines; trailing annotations after code
are not currently captured.

Language support is lightweight scanning, not complete parsing of
every supported language. See the limitations below.

## Commands

```text
whytie init
whytie scan <path>
whytie list
whytie why <memory-id|file:line>
whytie history <memory-id|file:line>
whytie resume
whytie version
```

### Initialize

```powershell
whytie init
```

Creates the local `.whytie` directory. Its internal `.gitignore`
keeps local database files out of Git.

### Scan

From the project root:

```powershell
whytie scan .
```

Finds annotations and synchronizes the saved reasoning.

Run a scan after editing, moving, or removing annotations.
There is no automatic background scan.

Generated and dependency directories such as `.git`, `.whytie`,
`node_modules`, `vendor`, `dist`, `build`, and `target` are skipped.
Filenames containing `.min.` are also skipped during directory scans.

### List

```powershell
whytie list
```

Displays currently stored reasoning, grouping adjacent annotations.

### Why

```powershell
whytie why shutdown.go:4
whytie why <memory-id>
```

Displays a stored annotation and its associated context.

A question can lead to its adjacent decision and reasons.
A decision can lead to its adjacent reasons.

Use the source location shown by WhyTie. Quote paths containing spaces:

```powershell
whytie why "src/my file.go:10"
```

### History

```powershell
whytie history shutdown.go:4
whytie history <memory-id>
```

Displays recorded lifecycle events:

- `created`: an annotation was first stored.
- `changed`: its text or kind changed.
- `moved`: its stored file or line changed.
- `deleted`: it was no longer found during synchronization.

Events may include Git commit information when available.

History records changes observed by scans. Changes made and undone
between scans are not recorded.

### Resume

```powershell
whytie resume
```

Shows:

- The five most recently saved history events.
- Event snapshots, recording times, and available commit information.
- Related context from currently stored annotations.
- An `Inspect` command for each displayed context block.

Related context blocks are displayed once, even when several recent
events belong to the same block.

Historical event snapshots and currently stored context are separate.
Deleted annotations can appear in the event list, but are not shown
as current context.

Resume reads saved data. It does not rescan source files, determine
what you last worked on, or mark questions as resolved or unresolved.

“Newest” means most recently inserted into the local history.
Within one scan, this does not establish the order in which you
actually edited the code.

If saved reasoning exists without history, Resume directs you to
`whytie list`.

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

When reconciliation recognizes the same reasoning, its memory ID
is retained and the change is added to its history.

For example, a reason can accumulate:

```text
history:
  created  shutdown.go:5
  changed  shutdown.go:5
  moved    lifecycle.go:5
  deleted  lifecycle.go:5
```

After deletion, retrieve its recorded history using a historical
location:

```powershell
whytie history lifecycle.go:5
```

This displays preserved reasoning; it does not restore source code.

Deleted-history lookup by location is supported. The current
ID-based history command requires a memory still present in the
current store.

Identity matching is best-effort. Arbitrary refactoring, duplicate
annotations, or ambiguous edits may not preserve the intended match.

## Local storage and privacy

Project data is stored in:

```text
.whytie/whytie.db
```

WhyTie uses SQLite and operates locally. Its core workflow requires
no account, hosted service, LLM, embeddings, or vector database.

Source code and reasoning are not uploaded by the core workflow.

Source comments remain readable without WhyTie. The local database
preserves identities and scan history, including deleted annotations.

The database is not committed to Git by default. A fresh clone does
not automatically include its history. Keep a backup of `.whytie`
if you need to preserve that local history.

## Git integration

WhyTie attaches Git provenance where available and can operate
outside a Git repository.

A displayed commit is provenance attached to the event. It is not
proof that the commit caused or validated the recorded decision.

WhyTie does not automatically reconstruct all annotations from
the repository's past commits.

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
Third-party license and attribution files are provided in
[third_party](third_party/README.md) and included in release packages.
