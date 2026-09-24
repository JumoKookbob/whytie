# AGENTS.md

## Purpose

When making non-trivial implementation decisions, preserve the reasoning using WhyTie annotations close to the relevant code.

WhyTie annotations are for information that a future developer or AI agent may need in order to understand **why** the code looks the way it does.

Do not use WhyTie annotations as ordinary code comments.

---

## WhyTie markers

Use these markers inside the source language's normal comment syntax:

| Marker | Meaning |
| --- | --- |
| `?` | Unresolved question |
| `+` | Decision |
| `x` | Rejected approach |
| `f` | Failed approach |
| `<` | Reason |
| `!` | Important context |

Example:

```go
// ? Which database should we use?
// + Use SQLite
// < The application should remain local-first with no database server
// x PostgreSQL
// < Adds infrastructure the current product does not need
```

---

## What to record

Add WhyTie annotations when the reasoning is likely to matter later.

Record:

- architectural decisions,
- meaningful implementation tradeoffs,
- dependency or technology choices,
- intentionally unusual or non-obvious code,
- temporary implementations that are expected to change later,
- compatibility or platform constraints,
- security, privacy, performance, or reliability constraints,
- rejected alternatives,
- failed approaches that should not be repeated,
- unresolved design questions,
- important context that is not obvious from the code itself.

Prefer recording decisions that answer a future question such as:

> Why was this done this way instead of the obvious alternative?

---

## What not to record

Do **not** add WhyTie annotations for:

- obvious implementation details,
- trivial code,
- routine variable or function choices,
- comments that merely restate what the code does,
- every refactor,
- every bug fix,
- formatting changes,
- generated code,
- information already obvious from nearby code or stable project documentation.

Avoid annotation noise.

A WhyTie annotation should earn its place by preserving reasoning that would otherwise be easy to lose.

---

## How to write annotations

### Keep them concise

Prefer:

```go
// + Use SQLite
// < Keep the application local-first with no external database server
```

Avoid long essays inside source files.

If a decision requires extensive documentation, keep the WhyTie annotation concise and point to the relevant document when useful.

### Record why, not what

Bad:

```go
// + Loop over all items
// < This loop iterates over the items
```

Good:

```go
// + Use a linear scan for v1
// < Expected datasets are small and simplicity matters more than indexing
```

### Keep annotations close to the code they explain

Place the reasoning immediately above or near the relevant implementation.

Do not collect unrelated WhyTie annotations in a distant file unless the reasoning itself is project-wide.

---

## Decision and reason pairing

For meaningful decisions, prefer pairing `+` with one or more `<` reasons.

```go
// + Keep manual snapshots available in the basic workflow
// < The product should remain useful without automation
```

When an alternative was considered and intentionally rejected, use `x`:

```go
// x Require PostgreSQL
// < A separate database server conflicts with the local-first design
```

When an approach was actually attempted and did not work, use `f`:

```go
// f Match reasoning identity by line number alone
// < Refactors caused identities to drift when code moved
```

Do not use `f` for an option that was merely considered.

Do not use `x` for an approach that was tried and failed.

---

## Questions

Use `?` only for meaningful unresolved questions that may affect implementation or architecture.

```go
// ? Should session history remain device-local or support export later?
```

Do not create `?` annotations for questions that can be answered immediately from the codebase.

When a question is resolved, update the nearby reasoning so the source reflects the current decision.

---

## Important context

Use `!` for constraints or context that future work must not accidentally violate.

```go
// ! Keep the local reasoning database out of Git
```

Use `!` sparingly.

It should communicate something important enough that losing it could cause a regression, incorrect redesign, or repeated mistake.

---

## Existing WhyTie annotations

Treat existing WhyTie annotations as durable project reasoning.

Do not remove, rewrite, or "clean up" them merely because the surrounding code is being refactored.

You may change an existing WhyTie annotation when:

- the underlying decision has changed,
- the recorded reason is no longer true,
- the annotation is factually incorrect,
- the associated implementation no longer exists,
- the user explicitly asks for it to be changed or removed.

When replacing an earlier decision, preserve enough context to explain why the new decision supersedes it.

Example:

```go
// + Use batched writes instead of immediate writes
// < Profiling showed immediate writes dominate save latency
// x Immediate writes
// < Simpler, but too slow for the current workload
```

Do not silently erase useful historical reasoning just to make comments shorter.

---

## AI-generated reasoning

Do not invent project history.

Only write a WhyTie annotation when the reasoning is supported by:

- the current task,
- the surrounding code,
- repository documentation,
- test behavior,
- user instructions,
- or a decision made during the current implementation work.

If the reason for an existing design is unknown, do not fabricate one.

Use a question instead when appropriate:

```go
// ? Why is this retry delay fixed instead of configurable?
```

WhyTie stores what is written into the source. It does not verify that an AI-generated explanation is correct.

---

## Annotation density

Default to fewer, higher-value annotations.

A normal implementation change may require no WhyTie annotation at all.

A useful rule:

> If a future developer can understand the reason directly from the code, do not add a WhyTie annotation.

Add one when the code shows **what** happened but not **why this choice was made**.

---

## Language syntax

Use the normal comment syntax for the file type.

Examples:

### Go, Rust, C/C++, C#, Java, Kotlin, JavaScript/TypeScript, Swift

```go
// + Use SQLite
// < Keep storage local
```

### Python, Ruby, shell, PowerShell, YAML

```python
# + Use SQLite
# < Keep storage local
```

### HTML

```html
<!-- + Keep the settings panel server-rendered -->
<!-- < Avoid adding a frontend framework for this small surface -->
```

### CSS and similar block-comment languages

```css
/* + Keep layout rules in this file */
/* < The component has no runtime styling requirements */
```

Use standalone comment lines.

Do not place WhyTie markers inside strings or string-like literals.

---

## Block comments

Where supported, WhyTie annotations may be written in block comments:

```go
/*
? Which database should we use?
+ Use SQLite
< Keep the application local-first
*/
```

JSDoc-style blocks may also contain WhyTie markers:

```javascript
/**
 * + Keep browser behavior in a separate JavaScript file
 * < Separating behavior keeps the HTML easier to inspect
 */
```

---

## During implementation

When working on a task:

1. Understand the requested change.
2. Inspect existing WhyTie annotations near the affected code.
3. Preserve existing reasoning unless the underlying decision changes.
4. Make the implementation change.
5. Add a WhyTie annotation only if a meaningful decision, tradeoff, rejection, failure, question, or constraint should survive the current session.
6. Keep the annotation short and close to the relevant code.
7. If WhyTie is available in the project, scan after meaningful annotation changes.

Do not add annotations merely to demonstrate compliance with this file.

---

## Before finishing a task

Check:

- Did I make a non-obvious decision that future work may question?
- Did I reject an alternative worth remembering?
- Did I try something that failed and should not be repeated?
- Did I discover an important constraint?
- Did I leave an unresolved question that matters?
- Did I accidentally remove or contradict an existing WhyTie annotation?
- Are the annotations concise and about **why**, not merely **what**?

If the answer to all of these is no, no new WhyTie annotation is required.

---

## Core principle

**Git records what changed. WhyTie remembers why.**

Preserve only the reasoning worth remembering.
