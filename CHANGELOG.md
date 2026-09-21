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
