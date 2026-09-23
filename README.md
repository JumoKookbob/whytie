# Changelog

## [1.1.0] - 2026-09-22

### Added

- `whytie resume` displays the five most recently saved history events

  and related context from saved annotations.

- Related context blocks are displayed once, with an `Inspect` command

  for further lookup.

- Recent-history queries preserve event snapshots, including deleted

  annotations.

### Fixed

- Go raw strings that close and reopen on the same line no longer

  cause string contents to be captured as reasoning comments.

- Updated the version-command test to match the executable version.

### Documentation and licensing

- Updated installation instructions, annotation syntax, commands,

  and known limitations.

- Added the Apache License, Version 2.0.

## [1.0.0]

### Included

- Local SQLite storage for reasoning annotations.

- Questions, decisions, reasons, rejected approaches, failed approaches,

  and important context.

- Multi-language source scanning and adjacent annotation grouping.

- Reasoning lookup by memory ID or source location.

- Best-effort identity preservation across recognized edits and moves.

- Created, changed, moved, and deleted history events.

- Preserved snapshots and historical-location lookup after deletion.

- Git provenance when available, with operation outside Git repositories.



Git records what changed. WhyTie remembers why.

WhyTie is a local-first CLI that preserves the reasoning behind your
code.

Write short comments next to the implementation. Scan your project.
Later, retrieve decisions and their reasons, review their history, or
use whytie resume to revisit recently recorded context.

No AI, account, or server is required.

// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}
Status

This source version is WhyTie v1.2.0.

See GitHub Releases
for published versions and downloadable packages.

Changes in v1.2.0 include:

whytie help, whytie --help, and whytie -h.
A redesigned whytie list grouped by source file with separators,
icons, and terminal colors.
Multiline /_ ... _/ and /\*_ ... _/ WhyTie annotations for
supported slash-comment languages.
Windows install and uninstall scripts so users do not need to create
a manual alias or PATH entry.

WhyTie is under active development. Automated tests and manual checks
cover creation, editing, movement, deletion, persistence, retrieval,
block-comment scanning, and CLI behavior.

What WhyTie does
Captures questions, decisions, reasons, rejected approaches, failed
approaches, and important context from source comments.
Stores reasoning locally in SQLite.
Groups adjacent annotations for readable context.
Retrieves reasoning by memory ID or source location.
Preserves reasoning identity across recognized edits and moves.
Records created, moved, changed, and deleted events.
Keeps historical snapshots after annotations disappear.
Attaches Git provenance when available.
Shows recent recorded events and related context with resume.
Presents current reasoning in a file-oriented list view.

WhyTie retrieves what you recorded. It does not invent explanations or
infer your next task.

Installation
Release packages

The v1.2.0 package targets and archive names are:

Platform Package name

Windows Intel/AMD 64-bit whytie-v1.2.0-windows-amd64.zip
Windows ARM64 whytie-v1.2.0-windows-arm64.zip
Linux x86-64 whytie-v1.2.0-linux-amd64.zip
Linux ARM64 whytie-v1.2.0-linux-arm64.zip
macOS Intel whytie-v1.2.0-darwin-amd64.zip
macOS Apple Silicon whytie-v1.2.0-darwin-arm64.zip

darwin means macOS.

Check GitHub Releases
for available downloads. Prebuilt packages do not require Go. Git is
optional and enables Git provenance when available.

Extract the entire archive, including its documentation and
third_party notices.

Windows

Extract the Windows package, then run:

.\install.cmd

The installer copies whytie.exe to:

%LOCALAPPDATA%\WhyTie\bin

and adds that directory to your user PATH.

After installation, use WhyTie from a terminal:

whytie version
whytie help

Then switch to a project directory:

whytie init
whytie scan .
whytie list
whytie resume

No manual alias is required.

To uninstall WhyTie from the extracted package directory:

.\uninstall.cmd

The uninstaller removes the installed executable and its user PATH
entry. Existing project .whytie directories are not deleted.

The .cmd wrappers invoke the bundled PowerShell scripts for the
installation or removal process without requiring the user to
permanently change their PowerShell execution policy.

Linux and macOS

Open a terminal in the extracted package directory:

chmod +x ./whytie
./whytie version
whytie_bin="$(pwd)/whytie"

Then switch to your project directory and run:

"$whytie_bin" init
"$whytie_bin" scan .
"$whytie_bin" list
"$whytie_bin" resume

The executable is a command-line tool; use it from a terminal.

Verification status

All six targets are cross-compiled as part of the release process.

Windows amd64 is runtime-tested locally. Other cross-compiled targets
may not have been runtime-tested on their target operating systems.

Release packages are not publisher-signed. macOS packages are not
notarized.

Build from source on Windows

Prerequisites:

Git
Go compatible with the go.mod requirement
PowerShell

Clone and build the current main branch:

git clone https://github.com/JumoKookbob/whytie.git
Set-Location .\whytie
go build -o .\whytie.exe .\cmd\whytie
.\whytie.exe version

You can run the local executable as .\whytie.exe, or install a release
package to make whytie available on your user PATH.

Building may require network access to download Go dependencies. Normal
scanning and retrieval operate locally.

Quick start

After installing WhyTie, create or open a project and initialize it:

whytie init

Add annotations next to the code:

package demo

// ? How should shutdown handle pending writes?
// + Flush pending writes before shutdown
// < Prevent user data loss
func Shutdown() {}

Scan the project:

whytie scan .

List recorded reasoning:

whytie list

The v1.2 list view groups reasoning by file and visually separates
blocks. When output is written directly to a supported terminal,
reasoning kinds can also be colorized.

Retrieve the question and its associated decision and reason:

whytie why shutdown.go:3

Return to recent recorded context:

whytie resume

Inspect the decision's history:

whytie history shutdown.go:4

Git is optional for this workflow.

Annotation syntax

Each marker must be followed by a space or tab and nonempty text.

Example Meaning

// ? How should this work? Question
// + Use this approach Decision
// < It preserves existing behavior Reason
// - Reject the shared-cache approach Rejected approach
// x The asynchronous attempt lost writes Failed approach
// ! Preserve shutdown ordering Important context

Both forms are accepted:

//? How should this work?
// ? How should this work?

This form is not accepted:

//?How should this work?

Keep related annotations on consecutive lines in the same file:

// ? Which storage should we use?
// + Use SQLite
// < Keep data local
// < Avoid requiring a separate server

The current grouping rules use adjacency. They do not infer semantic
relationships between distant comments.

For why, place a decision immediately after its question, and reasons
immediately after the decision.

Block comments

WhyTie v1.2.0 supports block-style annotations in supported
slash-comment languages.

Single-line block comments are supported:

/_ + Use a linear search for the prototype _/

Multiline block comments and JSDoc-style blocks can contain WhyTie
annotations:

/\*\*

- - Use a linear search for the prototype
- < The current data set is small
  \*/

Ordinary JSDoc without WhyTie markers is ignored:

/\*\*

- Returns a user by ID.
- @param id user ID
  \*/

Comment-like text inside protected string or template-literal contexts
is not treated as WhyTie reasoning by the scanner behavior covered by
the current tests.

Language-specific prefixes

Source files Annotation form

Go, Rust, C/C++, C#, Java, Kotlin, // ? Question and supported
JavaScript/TypeScript, Swift /_ ... _/ blocks

Python, Ruby, shell scripts, # ? Question
PowerShell, YAML

CSS, SCSS, Sass, Less /_ ? Question _/

HTML <!-- ? Question -->

Use standalone comment lines. Trailing annotations after code are not
currently captured.

Language support is lightweight scanning, not complete parsing of every
supported language.

Commands
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

Running whytie without a command also displays help.

Initialize
whytie init

Creates the local .whytie directory. Its internal .gitignore keeps
local database files out of Git.

Scan

From the project root:

whytie scan .

Finds annotations and synchronizes the saved reasoning.

Run a scan after editing, moving, or removing annotations. There is no
automatic background scan.

Generated and dependency directories such as .git, .whytie,
node_modules, vendor, dist, build, and target are skipped.
Filenames containing .min. are also skipped during directory scans.

List
whytie list

Displays currently stored reasoning.

In v1.2.0, the output is organized around source files, includes line
numbers and visual separators, and uses icons for reasoning kinds. When
writing directly to a supported terminal, WhyTie can add colors without
inserting ANSI color codes into non-terminal output such as test
buffers.

Why
whytie why shutdown.go:4
whytie why <memory-id>

Displays a stored annotation and its associated context.

A question can lead to its adjacent decision and reasons. A decision can
lead to its adjacent reasons.

Use the source location shown by WhyTie. Quote paths containing spaces:

whytie why "src/my file.go:10"
History
whytie history shutdown.go:4
whytie history <memory-id>

Displays recorded lifecycle events:

created: an annotation was first stored.
changed: its text or kind changed.
moved: its stored file or line changed.
deleted: it was no longer found during synchronization.

Events may include Git commit information when available.

History records changes observed by scans. Changes made and undone
between scans are not recorded.

Resume
whytie resume

Shows:

The five most recently saved history events.
Event snapshots, recording times, and available commit information.
Related context from currently stored annotations.
An Inspect command for each displayed context block.

Related context blocks are displayed once, even when several recent
events belong to the same block.

Historical event snapshots and currently stored context are separate.
Deleted annotations can appear in the event list, but are not shown as
current context.

Resume reads saved data. It does not rescan source files, determine what
you last worked on, or mark questions as resolved or unresolved.

"Newest" means most recently inserted into the local history. Within one
scan, this does not establish the order in which you actually edited the
code.

If saved reasoning exists without history, Resume directs you to
whytie list.

Help
whytie help
whytie --help
whytie -h

Displays command usage and the annotation marker reference.

Running whytie with no arguments displays the same help information.

Version
whytie version

Displays the version embedded in the executable.

Reasoning survives recognized changes

After editing or moving annotations, scan again:

whytie scan .

When reconciliation recognizes the same reasoning, its memory ID is
retained and the change is added to its history.

For example, a reason can accumulate:

history:
created shutdown.go:5
changed shutdown.go:5
moved lifecycle.go:5
deleted lifecycle.go:5

After deletion, retrieve its recorded history using a historical
location:

whytie history lifecycle.go:5

This displays preserved reasoning; it does not restore source code.

Deleted-history lookup by location is supported. The current ID-based
history command requires a memory still present in the current store.

Identity matching is best-effort. Arbitrary refactoring, duplicate
annotations, or ambiguous edits may not preserve the intended match.

Local storage and privacy

Project data is stored in:

.whytie/whytie.db

WhyTie uses SQLite and operates locally. Its core workflow requires no
account, hosted service, LLM, embeddings, or vector database.

Source code and reasoning are not uploaded by the core workflow.

Source comments remain readable without WhyTie. The local database
preserves identities and scan history, including deleted annotations.

The database is not committed to Git by default. A fresh clone does not
automatically include its history. Keep a backup of .whytie if you
need to preserve that local history.

Git integration

WhyTie attaches Git provenance where available and can operate outside a
Git repository.

A displayed commit is provenance attached to the event. It is not proof
that the commit caused or validated the recorded decision.

WhyTie does not automatically reconstruct all annotations from the
repository's past commits.

Current limitations
Scanning is explicit; source changes are not captured automatically.
Resume shows saved context, not a complete work-session summary.
Questions have no tracked resolution status.
Relationships primarily rely on adjacent source annotations.
Language scanners do not fully model every language's syntax.
Identity preservation is best-effort, not a guarantee for all
refactors.
Historical locations can become ambiguous when reused by other
annotations.
Local history is not automatically synchronized between machines.
There is no dedicated export command yet.
There is no Change Guard, impact analysis, GUI, or editor
integration.
Release packages are not publisher-signed; macOS packages are not
notarized.
Development

Run tests:

go test ./...

Build on Windows:

go build -o .\whytie.exe .\cmd\whytie

Run directly from source:

go run ./cmd/whytie resume
License

WhyTie is licensed under the Apache License, Version 2.0.

See LICENSE for the full terms.

Third-party license and attribution files are provided in
third_party and included in release packages.
