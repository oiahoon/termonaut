# Troubleshooting

| Symptom | First check |
| --- | --- |
| Missing internal/cache or Go type/API compile errors | Known upstream `7d20cf1` blockers; see STATUS before runtime testing |
| `make build` fails with no Go files in `cmd` | Use `go build -o /tmp/termonaut-dev ./cmd/termonaut`; Makefile target is stale |
| SQLite driver fails at runtime | Build with CGO and a working C compiler; avoid `CGO_ENABLED=0` |
| Build cache permission failure | Set `GOCACHE` to a writable scratch directory and retry |
| Avatar is slow or fallback art appears | Check DiceBear/network/cache; PNG download has no explicit timeout |
| XP or session totals seem inaccurate | Inspect capture/session path; actual exit status and duration are not collected |
| Achievements or sessions look prepopulated | Some enhanced dashboard panels are static examples |
| `stats --json` is not valid JSON | Current path uses Go formatting; fix/use a verified export path |
| Minimal TUI mode crashes | Known nil-command path; use `stats --today` for simple output |
| Config value saves but behavior is unchanged | Trace the consuming path; several settings are not wired through |
| Command absent despite an implementation file | Check Cobra registration and fresh `--help` output |

For hook problems inspect generated configuration and existing hooks in a
disposable shell first. Do not repeatedly force-install or overwrite normal shell
configuration. Collect synthetic reproductions instead of posting raw history.

See [STATUS](STATUS.md) for source findings and [DEVELOPMENT](DEVELOPMENT.md)
for an isolated baseline. Historical troubleshooting advice is preserved in
[archive](archive/README.md) and must be rechecked before use.
