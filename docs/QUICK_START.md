# Quick start

**Current blocker:** after synchronizing upstream `7d20cf1`, package setup and
compilation fail (missing local cache package and type/API mismatches). See
[STATUS](STATUS.md). Commands below describe the intended verification path;
repair compilation first. Do not run a stale existing binary as current HEAD proof.

Start with a source build and the [isolated smoke test](DEVELOPMENT.md).
This verifies local recording without changing your normal shell configuration.

```sh
go build -o /tmp/termonaut-dev ./cmd/termonaut
/tmp/termonaut-dev --help
/tmp/termonaut-dev version
```

For actual daily use, put a freshly built binary in a stable location on PATH.
Run `termonaut setup` for the interactive setup or `termonaut init` for hook
installation only when you intend to modify your shell configuration.
`termonaut quickstart` also installs integration and writes defaults; it is not
a read-only preview. Restart the shell after deliberate installation.

Useful commands (check a fresh binary's help for the current command tree):

```sh
termonaut stats --today
termonaut tui --mode smart
termonaut tui --mode compact
termonaut tui --mode full
termonaut tui --mode classic
termonaut config get
```

Avoid `tui --mode minimal` until its nil-command bug is fixed. `stats --json`
currently formats Go values rather than valid JSON. Some dashboard panels are
examples; see [STATUS](STATUS.md).

Config and data default to `~/.termonaut/`. Back up this directory before changing
an existing installation; copy a live WAL database using SQLite's backup mechanism
or stop writers before copying all database files. Never share raw command history
or assume sanitizer coverage is complete. GitHub sync is opt-in and can publish
generated content. Avatar display may make network requests.

For details: [architecture](ARCHITECTURE.md), [development](DEVELOPMENT.md),
[troubleshooting](TROUBLESHOOTING.md).
