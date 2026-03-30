---
name: windows-path-rename
description: Sanitize Git paths that are invalid on Windows (for example :, ", |, ?, * and trailing dot/space) and commit safe renames. Use when clone or checkout on Windows fails with "error: invalid path".
---

# Windows Path Rename

Use this skill when a repository contains file or directory names that cannot be created on Windows and Git reports `error: invalid path`.

## When To Use

- Clone or checkout fails on Windows due to invalid paths.
- Team needs a permanent rename commit so Windows users can work normally.
- You need a deterministic old->new mapping, with collision handling.

## Workflow

1. Ensure a clean branch for rename work.
2. Run the bundled script in preview mode to generate a mapping.
3. Review mapping and conflict suffixes.
4. Apply renames (prefer Linux/WSL working tree mode).
5. Commit with a clear message and push.

## Recommended Execution Modes

- Preferred: run in WSL/Linux where old invalid names are materialized; script uses `git mv`.
- Fallback (Windows-only, no invalid files materialized): generate mapping and apply via index plumbing in a controlled branch.

## Commands

Preview only:

```bash
python .github/skills/windows-path-rename/scripts/sanitize_windows_paths.py --mode preview
```

Apply with `git mv`:

```bash
python .github/skills/windows-path-rename/scripts/sanitize_windows_paths.py --mode apply
```

Write mapping to a file:

```bash
python .github/skills/windows-path-rename/scripts/sanitize_windows_paths.py --mode preview --output .tmp/windows-path-map.tsv
```

## Rename Rules

- `:` -> ` -`
- `"` -> `` (remove)
- `|` -> `_`
- `?` -> `` (remove)
- `*` -> `_`
- `<` -> `(`
- `>` -> `)`
- Remove trailing dots/spaces from each path segment
- Avoid Windows reserved names (`CON`, `PRN`, `AUX`, `NUL`, `COM1..COM9`, `LPT1..LPT9`) by appending `_`
- Preserve file extension when possible
- Resolve collisions by appending ` (n)` before extension

## Guardrails

- Do not rename `.git` internals.
- Do not run apply mode on a dirty working tree.
- Commit only path renames in the rename commit.
- After commit, verify no invalid paths remain with a preview run.
