# Current implementation status

Source review: **2026-10-28**, refreshed after synchronizing upstream `7d20cf1`. Default source version: `v0.10.2` in
`cmd/termonaut/main.go`; release builds can override it with linker flags.
This is a code/verification snapshot, not a production or installed-version audit.

## Direction

Preserve this project so it can be resumed without reconstructing context.
No feature expansion is committed. If work resumes, first decide whether the
personal tracking/XP/avatar loop is still useful. Historical plans are reference
material, not an active backlog. No remote repository archive setting or CI
schedule was changed by the documentation cleanup.

## Source paths (current runtime verification blocked)

| Area | What the source implements | Verification boundary |
| --- | --- | --- |
| Capture | Shell hooks launch background `log-command`; filtering and SQLite insertion | Current build blocked; shell behavior unverified |
| Persistence | Command/session queries, user progress, achievement storage | Earlier baseline passed; current suite blocked |
| Progression | Classification, streaks, XP, achievement checks | Upstream adds tests; current suite blocked |
| Stats | Command totals, today count, common commands | Not all time-range helpers are complete |
| UI | Text stats and classic/enhanced Bubble Tea dashboards | Some panels are static examples |
| Avatars | DiceBear fetching, PNG-to-ASCII, caching, fallback | No live network or resize verification in this review |
| GitHub | Badge/profile/heatmap generation, clone/update/commit/push code | No authenticated remote round trip verified |
| API | HTTP routing, several stats/search handlers | Multiple endpoints remain skeletons |

## Confirmed gaps to recheck before development

1. **Capture fidelity:** `runLogCommandCommand` assigns `ExitCode: 0` and no actual
   duration. Its session call passes Zsh as shell type regardless of actual shell.
   `GetTerminalPID` returns the immediate parent PID; background subshells may
   affect grouping. Success rate, failure penalties, and timing require real
   execution evidence before being trusted.
2. **Shell side effects:** generated Bash/Zsh hooks manipulate job control;
   the Bash hook installs a DEBUG trap. Compatibility with existing shell
   customization needs a disposable interactive-shell test.
3. **Static dashboard panels:** enhanced `renderAchievements`,
   `renderSessionHistory`, and `renderActivityHeatmap` display hardcoded examples.
   Analytics daily average divides total commands by 30, not observed elapsed days.
4. **Statistics gaps:** basic `GetWeeklyStats`/`GetMonthlyStats` return zero
   placeholders; first/last command-time lookup returns nil. Separate analytics
   code exists, so verify which path a caller actually uses.
5. **Configuration wiring:** the capture path constructs default sanitizer settings
   instead of mapping the loaded privacy settings. Enhanced TUI initial avatar
   request hardcodes username `user` and style `pixel-art`. Config field existence
   does not establish that a particular execution path respects it.
6. **Command reachability/output:** several standalone commands in `main.go` are
   explicitly not registered. `stats --json` uses `%+v`, not a JSON encoder.
   `tui --mode minimal` calls `runStatsCommand(nil, ...)`, which then dereferences
   `cmd.Flags()`; this path needs correction and a regression test.
7. **API skeletons:** command-by-ID and realtime endpoints return not-implemented;
   CSV export uses a static example; config update acknowledges without saving.
8. **Network robustness:** avatar PNG download uses `http.Get` without an explicit
   timeout. Validate timeout/error behavior before treating fetching as bounded.
9. **Build/docs drift:** Makefile `build-go` targets `./cmd`, but the main package
   is `./cmd/termonaut`. Python setup targets are historical scaffolding.
   Some help text, website claims, and release notes overstate implemented behavior.

Sources: [capture/CLI](../cmd/termonaut/main.go),
[shell](../internal/shell/shell.go), [TUI routing](../cmd/termonaut/tui.go),
[enhanced dashboard](../internal/tui/enhanced/dashboard.go),
[stats](../internal/stats/stats.go), [API](../internal/api/server.go),
[converter](../internal/avatar/converter.go), [Makefile](../Makefile).

## Verification recorded on 2026-10-28

### Current synchronized baseline: upstream `7d20cf1`

Before pushing documentation, fetching origin revealed three previously absent
upstream commits (`2f7ef21`, `ea4f254`, `7d20cf1`). Their application code is
preserved unchanged. **The current checkout does not pass `go test ./...`.**
With cached dependencies and a writable GOCACHE, observed blockers include:

- `internal/database/database.go` imports a missing local `internal/cache` package.
- `internal/pool/object_pool.go` references nonexistent `Command.DurationMs` and
  `Command.Category` fields (the existing duration field is `DurationMS`).
- `internal/avatar/manager.go` assigns the two-result cache Get call to one variable.
- `internal/config/config.go` declares an unused `cfg` variable in Get.
- Source review additionally shows `logger` declared twice in the DB struct;
  the missing-package failure may mask further compiler errors.

These are upstream source problems, not a failure introduced by this docs commit.
Restore a compiling baseline before relying on fresh CLI smoke tests. Do not
mistake an old binary in /tmp or the checkout for a build of synchronized HEAD.
The error list is not guaranteed exhaustive until compilation progresses.

The synchronized tree has nine Go test files across unit, integration, and
benchmark directories (configuration, privacy, network, database, and workflow
coverage). Their presence is not proof they pass: the suite is blocked at compile
and package setup. Lifecycle commands and config helpers have also been added.

### Earlier local baseline: `823ed6a`

Before remote synchronization, `GOCACHE=/private/tmp/termonaut-go-cache
GOPROXY=off go test ./...` passed with Go 1.27.3 on macOS arm64. That checkout had
only `tests/unit/basic_test.go` (database basics and config save/load).
A fresh source build, version/help inspection, and disposable-HOME smoke test
(config get → synthetic log-command → stats --today showing Commands: 1) passed
on that earlier baseline only. This is historical evidence, not current HEAD proof.

All archived originals are checksum-verified; current documentation local link
targets are checked. No real shell config, user activity database, or remote sync
was exercised. Update this document when compilation is repaired and re-run the
current tests before upgrading the verification claim.
