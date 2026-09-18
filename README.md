# WhyTie

**Keep the reason next to the code.**

WhyTie is a local-first CLI for recording and retrieving the reasoning behind code decisions.

Git tells you **what changed**.

WhyTie helps you remember **why it was written that way**.

Instead of maintaining a separate decision log, you write small structured comments directly next to the code:

```go
//? Which database should we use?
//+ Use SQLite
//< The tool should work local-first
//< It should not require a separate server
```

WhyTie scans these comments, stores them locally, and lets you retrieve the reasoning later.

## Status

Current version:

```text
v0.1.1
```

WhyTie is currently an early MVP.

The core workflow is implemented and dogfooded, but the project is still experimental.

## Syntax

WhyTie currently understands three comment types.

```text
//? question
//+ decision
//< reason
```

Example:

```go
//? Which journal mode should SQLite use?
//+ Use WAL mode
//< Reduce unnecessary blocking between reads and writes
//< Keep the design suitable for a future background scanner
```

The comments stay in the source code, close to the implementation they explain.

## Basic workflow

Initialize WhyTie inside a project:

```powershell
whytie init
```

This creates:

```text
.whytie/
```

Scan the project:

```powershell
whytie scan .
```

Example output:

```text
question: Which database should we use?  session.go:48

decision: Use SQLite  session.go:49
└─ reason: The tool should work local-first  session.go:50
└─ reason: It should not require a separate server  session.go:51
```

Then ask why a particular line exists:

```powershell
whytie why session.go:49
```

Output:

```text
decision: Use SQLite  session.go:49
└─ reason: The tool should work local-first  session.go:50
└─ reason: It should not require a separate server  session.go:51
```

You can also query the question:

```powershell
whytie why session.go:48
```

Output:

```text
question: Which database should we use?  session.go:48
└─ decision: Use SQLite  session.go:49
   └─ reason: The tool should work local-first  session.go:50
   └─ reason: It should not require a separate server  session.go:51
```

## Commands

```text
whytie init
whytie scan <path>
whytie list
whytie why <memory-id|file:line>
whytie version
```

### `init`

Initializes a WhyTie repository.

```powershell
whytie init
```

### `scan`

Scans source files for WhyTie comments and synchronizes them with the local repository.

```powershell
whytie scan .
```

### `list`

Lists the reasoning currently stored by WhyTie.

```powershell
whytie list
```

### `why`

Retrieves reasoning by memory ID or source location.

```powershell
whytie why session.go:49
```

### `version`

Prints the current version.

```powershell
whytie version
```

```text
WhyTie v0.1.1
```

## Local-first

WhyTie does not require an account or remote server for its core workflow.

Repository data is stored locally in:

```text
.whytie/whytie.db
```

using SQLite.

The source comments remain the human-readable source of the reasoning, while the local database allows WhyTie to track and query them.

## Editing code

WhyTie is designed for normal source-code editing.

For example, if this:

```go
//+ Use SQLite
```

becomes:

```go
//+ Use SQLite WAL mode
```

and the project is scanned again, WhyTie updates the stored memory instead of requiring you to manually recreate it.

WhyTie also reconciles source-location changes when recognized and removes stale memories during synchronization when their corresponding comments have been removed.

## Example

A larger decision can look like this:

```go
//? Should session closing use Get -> modify -> Save?
//+ Record the end time with a single UPDATE
//< Avoid an unnecessary read-modify-write cycle
func (s *Store) CloseSession(...) error {
    // ...
}
```

Months later:

```powershell
whytie why session.go:393
```

can recover the reasoning without requiring a separate architecture document or issue thread.

## Philosophy

Code preserves implementation.

Git preserves change history.

WhyTie is an experiment in preserving the small decisions and reasons that are otherwise easy to lose while developing software.

The goal is not to replace Git, documentation, comments, or ADRs.

It is to make lightweight reasoning cheap enough to record while you are already coding.

## Development

Run the test suite:

```powershell
go test ./...
```

Build the CLI:

```powershell
go build -o .\whytie.exe .\cmd\whytie
```

Check the version:

```powershell
.\whytie.exe version
```

## Current limitations

WhyTie v0.1.1 is intentionally small.

It currently focuses on:

```text
question -> decision -> reason
```

and source-location-based reasoning retrieval.

More advanced relationships, editor integrations, Git-aware history, automatic reasoning generation, and other higher-level features are outside the current MVP.

## License

No license has been specified yet.
