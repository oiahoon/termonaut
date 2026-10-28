# Operations

## Homepage

The existing homepage remains in `docs/index.html`, `docs/assets/`, and associated
site configuration. Serve it without changing application data:

```sh
python3 -m http.server 8000 --directory docs
```

See `.github/workflows/deploy-homepage.yml` for deployment behavior and
`scripts/dev-homepage.sh` for optional development helpers (it may open a browser).
The homepage's marketing claims and external links were not revalidated by the
2026-10-28 documentation cleanup. Do not infer measured performance or working
features from its demo output. No deployment was performed.

## Releases

Read the actual workflow before release work:

- `.github/workflows/release.yml`: tag/manual-triggered builds with CGO, unit tests,
  packaging, and GitHub release publication.
- `.github/workflows/manual-release.yml`: separate manual workflow.
- `.github/workflows/update-homebrew.yml` and `Formula/`: Homebrew integration.
- `scripts/install/` and `scripts/maintenance/`: installation/release helpers.

Source version defaults to `v0.10.2`; build linker flags can override version,
commit, and date. Release scripts, install asset names, Formula version, docs,
and workflow toolchains may have drifted; compare them before publishing.
This cleanup did not audit remote tags, installed binaries, credentials, or CI runs.

For a future authorized release: verify the exact commit/version, inspect build
matrix and C cross-compilers, run tests, validate actual asset names against the
installer/Formula, then publish only within the requested release scope. Verify
published artifacts and remote state afterwards. Old release checklists are
historical references, not automatically safe instructions.

## User data and external effects

Shell setup writes shell configuration and may add a shortcut. GitHub sync can
clone repositories, generate files, commit, and push. API persistence is incomplete.
Do not use any of these operations as incidental checks while restoring context.
Preserving the project locally does not disable remote workflows or archive GitHub.

Original operations guides are retained in [archive](archive/README.md).
