# AI coding entry point

## Start here

1. Read [README.md](README.md), [docs/STATUS.md](docs/STATUS.md), then
   [docs/RESUME.md](docs/RESUME.md).
2. Run `git status --short` and inspect changes relevant to the request. Preserve
   existing user changes and untracked files.
3. Current synchronized upstream `7d20cf1` does not compile. STATUS records the
   observed blockers; earlier successful tests apply only to `823ed6a`. Do not
   present old binaries as current verification or silently repair unrelated code.
4. Read [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and the smallest relevant
   source files. Source code and current command registration outrank old docs.
5. For build/validation use [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).

This repository is being preserved for optional future work. Do not interpret
historical roadmaps or unfinished feature lists as authorization to implement them.
Root entry rules apply to CLI, internal packages, scripts, docs, and homepage.

## Scope and evidence

- Keep changes focused on the user's request. Prefer a small working core over
  expansion into API, sharing, or additional UI variants.
- Distinguish stored data, calculated estimates, static demo output, and runtime
  verification. Never describe a placeholder panel or HTTP success response as
  proof of working persistence.
- Audit actual Cobra registration: Go files containing a command do not prove
  that the command is reachable. Check fresh `--help` output.
- Do not use archived completion reports as evidence that features work.
- Never promise secret removal based only on regex sanitization. Use synthetic
  commands in tests and never print real command-history secrets.

## Safe development

- Build `./cmd/termonaut`; `make build` currently points at `./cmd` and is stale.
- `go-sqlite3` needs CGO and a C compiler; do not switch to `CGO_ENABLED=0`
  as a shortcut for a functioning runtime.
- Use a temporary `HOME` for CLI smoke tests (see development guide).
  Even config loading can create files. Avoid invoking setup/init/quickstart,
  shell install/uninstall, or GitHub sync as routine verification.
- Changes to shell hooks must be checked in a disposable interactive shell;
  never source generated hooks into the developer's normal shell automatically.
- Keep tests proportional. For behavior changes add a focused regression test;
  for documentation-only work validate links, facts, and `git diff --check`.
- Report exactly what was tested and what remains unverified. Local Go tests do
  not establish shell compatibility, deployed release health, or live sync.

## Keep context useful

- Update STATUS when an implementation gap is fixed or newly confirmed.
- Update architecture/development docs only when stable structure or commands change.
- Keep restart priorities and prompts in RESUME; do not add more completion reports.
- Keep historical docs under `docs/archive/`; consult them only for relevant history.
- Do not commit, push, publish, tag a release, install hooks, or modify the user's
  real data merely as a consequence of a documentation or code-reading request.
