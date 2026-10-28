# Termonaut

A Go CLI that turns terminal activity into local usage statistics, XP, levels,
achievements, and colorful ASCII avatars.

**Project direction (2026-10-28):** preserve it as a resumable personal project.
There is no committed expansion roadmap. Core tracking and progression exist;
several dashboard panels and API endpoints still contain placeholders. The latest
synchronized upstream code (`7d20cf1`) also has confirmed build blockers. See
[implementation status](docs/STATUS.md) before relying on displayed metrics.
This is a repository direction, not a claim that GitHub has been archived or
scheduled workflows disabled.

## Pick it up again

- **AI coding:** start with [AGENTS.md](AGENTS.md).
- **Returning maintainer:** follow [RESUME.md](docs/RESUME.md).
- **Try it locally:** use the isolated [quick start](docs/QUICK_START.md).
- **Understand the code:** read [architecture](docs/ARCHITECTURE.md).
- **Build and verify:** read [development](docs/DEVELOPMENT.md).
- **All documentation:** [documentation index](docs/README.md).

## Build

Requires the Go toolchain declared in `go.mod` and a C compiler for
`github.com/mattn/go-sqlite3` (CGO). The module declares Go 1.23.0 and toolchain
Go 1.24.4; newer compatible toolchains may work.

```sh
go build -o /tmp/termonaut-dev ./cmd/termonaut
/tmp/termonaut-dev --help
go test ./...
```

Build from source when resuming; binaries in the checkout may be old artifacts.
Run the isolated quick start before installing any shell hooks.

## What it contains

- Shell hooks that invoke a background `log-command` CLI process.
- Local SQLite command/session storage and XP/achievement progression.
- Basic statistics plus Bubble Tea / Lip Gloss terminal interfaces.
- DiceBear avatars rendered with `ascii-image-converter`, with cache/fallbacks.
- Optional GitHub profile/badge/heatmap generation and repository sync.
- An experimental HTTP API with incomplete endpoints.

Data defaults to `~/.termonaut/`. Avatar fetching uses a remote service;
GitHub sync is disabled by default. Local storage does not mean every feature
is offline. Recorded command volume is an activity measure, not proof of
productivity or successful command execution.

## History and license

[Changelog](docs/CHANGELOG.md) · [Historical documentation](docs/archive/README.md)
· [MIT license](LICENSE)

Old documentation is retained for recovery, not used as current acceptance criteria.
