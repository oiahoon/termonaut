# Development

**Current blocker:** after synchronizing upstream `7d20cf1`, package setup and
compilation fail (missing local cache package and type/API mismatches). See
[STATUS](STATUS.md). Commands below describe the intended verification path;
repair compilation first. Do not run a stale existing binary as current HEAD proof.

## Prerequisites and build

Use Go compatible with `go.mod` (Go 1.23.0; toolchain hint Go 1.24.4), CGO,
and a C compiler. On macOS a working Command Line Tools installation supplies
Clang. The SQLite driver needs CGO at runtime.

```sh
go version
go env CGO_ENABLED CC
go build -o /tmp/termonaut-dev ./cmd/termonaut
/tmp/termonaut-dev version
/tmp/termonaut-dev --help
go test ./...
```

Use the explicit package path above. `make build` currently targets `./cmd` and
needs repair. Do not use the old Python setup targets or run `go mod init` in this
existing Go module. Cross-compilation with SQLite requires a suitable C toolchain;
setting `CGO_ENABLED=0` is not sufficient for a working application.

## Isolated local smoke test

Build first, then override HOME only for individual application invocations.
This keeps the Go module cache/toolchain and your real application data separate.
`log-command` below records a synthetic string; it does not execute `echo`.

```sh
go build -o /tmp/termonaut-dev ./cmd/termonaut
termonaut_test_home=$(mktemp -d)
HOME="$termonaut_test_home" /tmp/termonaut-dev config get
HOME="$termonaut_test_home" /tmp/termonaut-dev log-command 'echo termonaut-smoke'
HOME="$termonaut_test_home" /tmp/termonaut-dev stats --today
printf 'Disposable data: %s
' "$termonaut_test_home"
```

Expect `Commands: 1` in today's statistics. Inspect the temporary directory and
remove it when finished. Do not call `setup`, `quickstart`, `init`, or shell
installation for this smoke test. Default `stats` and TUI may fetch avatars;
`stats --today` avoids that path.

If the environment blocks the normal Go build cache, point `GOCACHE` at a writable
scratch directory. `GOPROXY=off` works only when required dependencies already
exist locally; offline dependency misses are environment failures, not test results.

## Verification by change type

| Change | Appropriate verification |
| --- | --- |
| Docs only | Local link targets, factual source checks, `git diff --check` |
| DB / XP / privacy | Focused regression test, then `go test ./...` |
| CLI registration / flags | Fresh binary `--help`, focused invocation with temp HOME |
| TUI | Disposable HOME, real terminal, resize/navigation checks; identify static panels |
| Shell hook | Generated-hook inspection plus disposable interactive shell with existing hooks |
| GitHub / API | Synthetic or disposable fixtures first; real remote operations require task scope |

Many shell scripts under `tests/` and `examples/` are demonstrations, perform
installation, or assume an existing binary. Read them before execution. They are
not automatically a safe unattended regression suite.

## Contribution workflow

Inspect `git status` before changing files. Use focused `codex/` branches for new
agent work if a branch is needed. Preserve unrelated changes. Use `gofmt` on
changed Go files, test affected behavior, and check whitespace. Update
[STATUS](STATUS.md) when closing a confirmed gap and [architecture](ARCHITECTURE.md)
when stable structure changes. Do not create another generic completion report.

Release and homepage details live in [OPERATIONS](OPERATIONS.md).
