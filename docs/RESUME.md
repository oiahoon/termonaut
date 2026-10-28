# Resume this project

## Five-minute restart

1. Read [README](../README.md), [STATUS](STATUS.md), and [AGENTS](../AGENTS.md).
2. Inspect `git status --short`, `git log -5 --oneline`, and `go.mod`.
3. Build from source and run `go test ./...` using [DEVELOPMENT](DEVELOPMENT.md).
4. Run the disposable-HOME smoke test. Do not install hooks to establish a baseline.
5. Choose one concrete outcome and trace its active code path before editing.

The latest upstream baseline is `7d20cf1`; it has confirmed compile blockers.
Earlier local `823ed6a` passed tests before synchronization. See STATUS for both.
The last documented baseline was reviewed on 2026-10-28. Recheck mutable facts;
do not assume the checkout, installed binary, remote releases, or dependencies
still match this snapshot.

## What is worth preserving

The distinctive personal experience is command activity → XP/level progression →
colorful terminal avatar. SQLite storage, command classification, ASCII conversion,
and terminal layout can also be reused independently. Continuing development is
optional; there is no obligation to finish historical roadmaps.

## Candidate next work, if useful

These are suggestions, not approved implementation tasks:

1. Restore compilation of synchronized upstream `7d20cf1` and run its expanded tests.
2. Make displayed information honest: connect or label static dashboard panels.
3. Fix capture fidelity/session identity and test shell compatibility.
4. Repair the basic build target, minimal TUI route, JSON output, and config wiring.
5. Add focused tests around XP/streak/privacy and persistence failures.
6. Reassess API/GitHub/multiple UI variants only after the small core is useful.

Avoid introducing more features before choosing what you personally want to use.
No changes to application code were made by the documentation consolidation.

## Copyable prompts

### Reorient without editing

> Read AGENTS.md, docs/STATUS.md, and docs/ARCHITECTURE.md. Check the current
> checkout and relevant source. Explain the active command capture → SQLite → XP
> → dashboard path, distinguish placeholders from stored data, and report what
> changed since the documented baseline. Do not edit files or install hooks.

### Resume a focused fix

> Read AGENTS.md and docs/RESUME.md. Fix [specific behavior] in its active code
> path, preserve unrelated changes, use disposable application data, add a focused
> regression test where appropriate, and update docs/STATUS.md. Report local test
> evidence separately from unverified shell/network behavior.

### Keep only the personal core

> Read the current architecture and status. Propose a minimal tracking + XP +
> avatar scope. Inspect callers and references before proposing removals. Present
> the concrete reduction plan before deleting features or changing stored data.

## Historical recovery

Older specs, duplicated setup guides, AI guidance, phase summaries, and release
notes are preserved under [archive](archive/README.md), with checksums and original
paths. Consult a specific historical document only when needed. Old plans and
“complete” reports do not override source or current verification.

At cleanup start, `GITHUB_RELEASE_DESCRIPTION.md` was already untracked. It was
left untouched; it is not an authoritative release state or part of the archive.
