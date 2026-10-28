# Architecture

## Execution model

Termonaut is a Go CLI, not a required background daemon. Cobra dispatches
commands from `cmd/termonaut`; Viper reads TOML configuration. Shell hooks spawn
`termonaut log-command` in the background for each captured command.

```text
Shell hook → log-command → sanitizer → SQLite command record
                                      → category/streak → XP/achievements
SQLite → stats/analytics → text output / Bubble Tea dashboard
Progress + avatar request → DiceBear PNG → ascii-image-converter → colored text
Progress + stats → optional GitHub generation / git commit and push
```

## Read the code in this order

| Responsibility | Source |
| --- | --- |
| CLI registration and command capture | [main.go](../cmd/termonaut/main.go) |
| Shell hook generation and installation | [shell.go](../internal/shell/shell.go) |
| Privacy filtering | [sanitizer.go](../internal/privacy/sanitizer.go) |
| Schema, connections, basic queries | [database.go](../internal/database/database.go) |
| Command queries | [commands.go](../internal/database/commands.go) |
| StoreCommandWithXP orchestration | [database/gamification.go](../internal/database/gamification.go) |
| Classification and XP | [classifier.go](../internal/categories/classifier.go), [xp.go](../internal/gamification/xp.go) |
| Basic and advanced statistics | [stats](../internal/stats), [analytics](../internal/analytics) |
| TUI selection | [tui.go](../cmd/termonaut/tui.go) |
| Enhanced TUI update/view/data loading | [dashboard.go](../internal/tui/enhanced/dashboard.go) |
| Avatar orchestration, caching, conversion | [avatar](../internal/avatar) |
| GitHub generation and scheduled sync | [sync.go](../internal/github/sync.go) |
| Experimental HTTP API | [server.go](../internal/api/server.go) |
| Data structures | [command.go](../pkg/models/command.go) |

## Persistence and configuration

Default config: `~/.termonaut/config.toml`. Default database:
`~/.termonaut/termonaut.db`, configurable through `data_dir`.
See [config.go](../internal/config/config.go) for actual defaults and loading.
Configuration loading creates defaults if no file exists.

The synchronized upstream code configures WAL, a ten-second busy timeout,
five open/two idle connections, and an LRU cache dependency that is missing from
the checkout. These are source settings, not a verified working configuration. Tables are created through `CREATE TABLE IF NOT EXISTS` on startup:
`commands`, `sessions`, `user_progress`, `achievements`, `daily_stats`.
There is no standalone versioned migration system in this initialization path.
Commands carry timestamp, session, command text, exit code, cwd, and duration;
field existence does not imply the shell hook supplies real values.

Command insertion and the subsequent progress updates are separate operations;
`StoreCommandWithXP` suppresses some progress-update failures after storing the
command. Do not assume the entire pipeline is one atomic transaction.

## Avatars

The manager derives a deterministic seed from the request and applies level-based
parameters. It fetches SVG for storage and a PNG endpoint for conversion.
`ConvertSVGToASCII` rewrites `/svg` to `/png`; it is not a general local SVG
renderer. Conversion enables ANSI colors and uses smaller character sets for
smaller images. Cache and fallback artwork support degraded network conditions.
The enhanced dashboard currently hardcodes parts of the avatar request.

## Interfaces and dependencies

- CLI/config: Cobra, Viper; logging: Logrus.
- Persistence: `mattn/go-sqlite3` (CGO).
- Terminal: Bubble Tea, Bubbles, Lip Gloss, `x/term`.
- Images: `TheZoraiz/ascii-image-converter` v1.13.1, DiceBear HTTP service.
- API: Gorilla Mux. GitHub sync invokes Git subprocesses.
- Upstream also adds network retry/circuit breaker code, object pools, memory
  monitoring, shared utilities, and lifecycle CLI commands (`status`, `doctor`,
  `cleanup`, `uninstall`). Their registration is present; current build blockers
  prevent fresh runtime verification. The PNG converter still uses plain http.Get.
- Homepage: static HTML/CSS/JavaScript under `docs/`, deployed separately.

The repository also contains multiple generations of TUI and enhanced
calculation helpers. Trace callers before extending or deleting a similarly
named implementation. See [STATUS](STATUS.md) for gaps in the active paths.
