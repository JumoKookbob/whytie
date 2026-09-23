# Changelog

## [1.2.0] - 2026-09-23

### Added

- Added `whytie help`, `whytie --help`, and `whytie -h`.
- Running `whytie` without a command now displays the help screen.
- Added Windows install and uninstall scripts.
- Windows users can install WhyTie onto their user `PATH` without manually creating an alias.
- Added multiline `/* ... */` and `/** ... */` annotation support for supported slash-comment languages.

### Changed

- Redesigned `whytie list` to make recorded reasoning easier to scan.
- List output is now grouped by source file.
- Added visual separators between reasoning blocks.
- Added icons for questions, decisions, reasons, rejected approaches, failed approaches, and important context.
- Added terminal colors when output is written directly to a supported terminal.
- File paths and line numbers are now presented before their associated reasoning instead of being repeated after every annotation.

### Fixed

- Block-style reasoning annotations can now be captured from multiline block comments instead of requiring the opening and closing delimiters on the same line.
- Ordinary JSDoc comments without WhyTie markers are ignored.
- Comment-like text inside supported string and template-literal contexts is not treated as WhyTie reasoning.

## [1.1.0] - 2026-09-22

### Added

- `whytie resume` displays the five most recently saved history events and related context from saved annotations.
- Related context blocks are displayed once, with an `Inspect` command for further lookup.
- Recent-history queries preserve event snapshots, including deleted annotations.

### Fixed

- Go raw strings that close and reopen on the same line no longer cause string contents to be captured as reasoning comments.
- Updated the version-command test to match the executable version.

### Documentation and licensing

- Updated installation instructions, annotation syntax, commands, and known limitations.
- Added the Apache License, Version 2.0.

## [1.0.0]

### Included

- Local SQLite storage for reasoning annotations.
- Questions, decisions, reasons, rejected approaches, failed approaches, and important context.
- Multi-language source scanning and adjacent annotation grouping.
- Reasoning lookup by memory ID or source location.
- Best-effort identity preservation across recognized edits and moves.
- Created, changed, moved, and deleted history events.
- Preserved snapshots and historical-location lookup after deletion.
- Git provenance when available, with operation outside Git repositories.
