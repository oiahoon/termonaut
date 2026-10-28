# Historical documentation

Consolidated on **2026-10-28**. Original Markdown files are preserved byte-for-byte
under `2026-10-28/`, retaining their original repository-relative paths except the ignored release
summary filename, which is renamed and mapped in the manifest.
These are historical claims and plans, not current feature guarantees.
Internal links inside historical files retain their original text and may be stale.

Current sources of truth: [status](../STATUS.md), [architecture](../ARCHITECTURE.md),
[development](../DEVELOPMENT.md), [resume checklist](../RESUME.md), and live source.
The original changelog remains at [CHANGELOG](../CHANGELOG.md).

No application source, shell integration, site assets, or release workflow was
changed. The pre-existing untracked `GITHUB_RELEASE_DESCRIPTION.md` was untouched.

## Consolidation map

- Old README, overview, planning, completeness/phase reports → current README and STATUS.
- Old development/AI guides and implementation reports → AGENTS, ARCHITECTURE, DEVELOPMENT, RESUME.
- Duplicate quick starts and user guides → QUICK_START, DEVELOPMENT, TROUBLESHOOTING.
- Avatar/TUI specs and analysis → ARCHITECTURE and STATUS; detail retained here.
- Homepage/release/GitHub guides → OPERATIONS; original instructions retained here.

## Original files

[Checksum manifest](2026-10-28/manifest.json) records original path, archive path,
and SHA-256 for recovery. To restore a document, copy its archived bytes to its
original path only after reviewing any current replacement. Do not bulk overwrite
current entry points. Git history also retains tracked originals.

| Original repository path | Preserved original |
| --- | --- |
| `README.md` | [original](2026-10-28/README.md) |
| `RELEASE_SUMMARY_0.9.2.md` | [original](2026-10-28/historical-release-summary-0.9.2.md) |
| `docs/API.md` | [original](2026-10-28/docs/API.md) |
| `docs/AVATAR_SYSTEM_SPEC.md` | [original](2026-10-28/docs/AVATAR_SYSTEM_SPEC.md) |
| `docs/CONTRIBUTING.md` | [original](2026-10-28/docs/CONTRIBUTING.md) |
| `docs/DEVELOPMENT.md` | [original](2026-10-28/docs/DEVELOPMENT.md) |
| `docs/GITHUB_ACTIONS_RELEASE.md` | [original](2026-10-28/docs/GITHUB_ACTIONS_RELEASE.md) |
| `docs/GITHUB_INTEGRATION.md` | [original](2026-10-28/docs/GITHUB_INTEGRATION.md) |
| `docs/HOMEBREW_RELEASE.md` | [original](2026-10-28/docs/HOMEBREW_RELEASE.md) |
| `docs/HOMEPAGE_IMPLEMENTATION.md` | [original](2026-10-28/docs/HOMEPAGE_IMPLEMENTATION.md) |
| `docs/PHASE1_SUMMARY.md` | [original](2026-10-28/docs/PHASE1_SUMMARY.md) |
| `docs/PROJECT_COMPLETENESS_REPORT.md` | [original](2026-10-28/docs/PROJECT_COMPLETENESS_REPORT.md) |
| `docs/PROJECT_OVERVIEW.md` | [original](2026-10-28/docs/PROJECT_OVERVIEW.md) |
| `docs/PROJECT_PLANNING.md` | [original](2026-10-28/docs/PROJECT_PLANNING.md) |
| `docs/QUICK_START.md` | [original](2026-10-28/docs/QUICK_START.md) |
| `docs/README.md` | [original](2026-10-28/docs/README.md) |
| `docs/SAFE_SHELL_CONFIGURATION.md` | [original](2026-10-28/docs/SAFE_SHELL_CONFIGURATION.md) |
| `docs/TROUBLESHOOTING.md` | [original](2026-10-28/docs/TROUBLESHOOTING.md) |
| `docs/analysis/PROJECT_STRUCTURE_ANALYSIS.md` | [original](2026-10-28/docs/analysis/PROJECT_STRUCTURE_ANALYSIS.md) |
| `docs/analysis/README_ANALYSIS_REPORT.md` | [original](2026-10-28/docs/analysis/README_ANALYSIS_REPORT.md) |
| `docs/analysis/TUI_IMPLEMENTATION_STATUS.md` | [original](2026-10-28/docs/analysis/TUI_IMPLEMENTATION_STATUS.md) |
| `docs/analysis/USER_ISSUES_ANALYSIS.md` | [original](2026-10-28/docs/analysis/USER_ISSUES_ANALYSIS.md) |
| `docs/development/AI_ASSISTANT_GUIDE.md` | [original](2026-10-28/docs/development/AI_ASSISTANT_GUIDE.md) |
| `docs/development/AVATAR_LAYOUT_IMPROVEMENTS.md` | [original](2026-10-28/docs/development/AVATAR_LAYOUT_IMPROVEMENTS.md) |
| `docs/development/COMMIT_READY_CHECKLIST.md` | [original](2026-10-28/docs/development/COMMIT_READY_CHECKLIST.md) |
| `docs/development/EASTER_EGGS_AND_AVATAR_IMPROVEMENTS.md` | [original](2026-10-28/docs/development/EASTER_EGGS_AND_AVATAR_IMPROVEMENTS.md) |
| `docs/development/ENHANCED_TUI_DEMO.md` | [original](2026-10-28/docs/development/ENHANCED_TUI_DEMO.md) |
| `docs/development/FINAL_COMMIT_AND_RELEASE_GUIDE.md` | [original](2026-10-28/docs/development/FINAL_COMMIT_AND_RELEASE_GUIDE.md) |
| `docs/development/LATEST_CHANGES_SUMMARY.md` | [original](2026-10-28/docs/development/LATEST_CHANGES_SUMMARY.md) |
| `docs/development/MUCH_WIDER_AVATAR_UPDATE.md` | [original](2026-10-28/docs/development/MUCH_WIDER_AVATAR_UPDATE.md) |
| `docs/development/ONBOARDING_SYSTEM_COMPLETE.md` | [original](2026-10-28/docs/development/ONBOARDING_SYSTEM_COMPLETE.md) |
| `docs/development/PANIC_FIX_SUMMARY.md` | [original](2026-10-28/docs/development/PANIC_FIX_SUMMARY.md) |
| `docs/development/PERMISSION_FIX_COMPLETE.md` | [original](2026-10-28/docs/development/PERMISSION_FIX_COMPLETE.md) |
| `docs/development/PROJECT_INTEGRATION_COMPLETE.md` | [original](2026-10-28/docs/development/PROJECT_INTEGRATION_COMPLETE.md) |
| `docs/development/TUI_ENHANCEMENT_IMPLEMENTATION.md` | [original](2026-10-28/docs/development/TUI_ENHANCEMENT_IMPLEMENTATION.md) |
| `docs/development/TUI_GUIDE.md` | [original](2026-10-28/docs/development/TUI_GUIDE.md) |
| `docs/development/UI_ENHANCEMENT_PLAN.md` | [original](2026-10-28/docs/development/UI_ENHANCEMENT_PLAN.md) |
| `docs/features/command-structure.md` | [original](2026-10-28/docs/features/command-structure.md) |
| `docs/features/enhanced-features-v0.10.1.md` | [original](2026-10-28/docs/features/enhanced-features-v0.10.1.md) |
| `docs/features/tui-modes.md` | [original](2026-10-28/docs/features/tui-modes.md) |
| `docs/releases/RELEASE_PREPARATION_v0.10.0.md` | [original](2026-10-28/docs/releases/RELEASE_PREPARATION_v0.10.0.md) |
| `docs/releases/current/v0.10.1/RELEASE_CHECKLIST_v0.10.1.md` | [original](2026-10-28/docs/releases/current/v0.10.1/RELEASE_CHECKLIST_v0.10.1.md) |
| `docs/releases/current/v0.10.1/RELEASE_NOTES_v0.10.1.md` | [original](2026-10-28/docs/releases/current/v0.10.1/RELEASE_NOTES_v0.10.1.md) |
| `docs/releases/current/v0.10.1/RELEASE_NOTES_v0.10.1_UPDATED.md` | [original](2026-10-28/docs/releases/current/v0.10.1/RELEASE_NOTES_v0.10.1_UPDATED.md) |
| `docs/releases/current/v0.10.2/RELEASE_NOTES_v0.10.2.md` | [original](2026-10-28/docs/releases/current/v0.10.2/RELEASE_NOTES_v0.10.2.md) |
| `docs/user-guide/configuration.md` | [original](2026-10-28/docs/user-guide/configuration.md) |
| `docs/user-guide/installation.md` | [original](2026-10-28/docs/user-guide/installation.md) |
| `docs/user-guide/quick-start.md` | [original](2026-10-28/docs/user-guide/quick-start.md) |

## Additional history found during remote synchronization

Remote main `7d20cf1` contained three commits absent from the initial checkout.
Their additional reports and README are preserved byte-for-byte below; the
[additional manifest](remote-main-7d20cf1/manifest.json) records paths and checksums.
Remote code is retained unchanged. Current build failures are recorded in STATUS.

| Original repository path | Preserved original |
| --- | --- |
| `docs/DEVELOPMENT_OPTIMIZATION_PLAN.md` | [original](remote-main-7d20cf1/docs/DEVELOPMENT_OPTIMIZATION_PLAN.md) |
| `docs/UX_ANALYSIS_REPORT.md` | [original](remote-main-7d20cf1/docs/UX_ANALYSIS_REPORT.md) |
| `docs/PHASE1_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/PHASE1_COMPLETION_REPORT.md) |
| `docs/PROJECT_STATUS_FINAL.md` | [original](remote-main-7d20cf1/docs/PROJECT_STATUS_FINAL.md) |
| `docs/PHASE4_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/PHASE4_COMPLETION_REPORT.md) |
| `docs/README_ANALYSIS_REPORT.md` | [original](remote-main-7d20cf1/docs/README_ANALYSIS_REPORT.md) |
| `docs/PHASE3_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/PHASE3_COMPLETION_REPORT.md) |
| `docs/OPTIMIZATION_COMPLETION_SUMMARY.md` | [original](remote-main-7d20cf1/docs/OPTIMIZATION_COMPLETION_SUMMARY.md) |
| `docs/OPTIMIZATION_ANALYSIS_REPORT.md` | [original](remote-main-7d20cf1/docs/OPTIMIZATION_ANALYSIS_REPORT.md) |
| `docs/PROJECT_FINAL_COMPLETION.md` | [original](remote-main-7d20cf1/docs/PROJECT_FINAL_COMPLETION.md) |
| `docs/PHASE4_DEVELOPMENT_PLAN.md` | [original](remote-main-7d20cf1/docs/PHASE4_DEVELOPMENT_PLAN.md) |
| `docs/PHASE3_DEVELOPMENT_PLAN.md` | [original](remote-main-7d20cf1/docs/PHASE3_DEVELOPMENT_PLAN.md) |
| `docs/reports/TASK2_1_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/reports/TASK2_1_COMPLETION_REPORT.md) |
| `docs/reports/PHASE2_START_REPORT.md` | [original](remote-main-7d20cf1/docs/reports/PHASE2_START_REPORT.md) |
| `docs/reports/PHASE2_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/reports/PHASE2_COMPLETION_REPORT.md) |
| `docs/reports/TASK2_2_COMPLETION_REPORT.md` | [original](remote-main-7d20cf1/docs/reports/TASK2_2_COMPLETION_REPORT.md) |
| `README.md` | [original](remote-main-7d20cf1/README.md) |
