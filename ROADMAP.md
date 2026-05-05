# skild Roadmap

## Current MVP Status

- Global setup command (`skild setup`) with config file at `~/.config/skild/config.json`.
- Core commands: `list`, `install`, `update`, `repo-sync`, `config`, `version`.
- Local install target support (`openCodeDir`).
- Global install target support (`--global`) using OpenCode path conventions:
  - `XDG_CONFIG_HOME/opencode/skills`
  - fallback `~/.config/opencode/skills`
- Quiet git output by default, with verbose mode via `SKILD_VERBOSE=1`.

## Next High-Priority Work

### 1) Update Behavior Control

Current behavior updates all discovered skills. Improve control to avoid unintended updates.

- Change `update` default behavior to update only installed skills.
- Add `update <skill-name>` for single-skill updates.
- Keep broad update behavior behind explicit `update --all`.

### 2) Skill Versioning and Reproducibility

Introduce lockfile-based state tracking.

- Add `skild-lock.json`.
- Track, at minimum, per installed skill:
  - `name`
  - `sourceRepo`
  - `repoRef`
  - `resolvedCommit`
  - `installedAt`
- Update lockfile on `install` and `update`.

### 3) Restore and Drift Visibility

- Add `restore` command to reinstall exactly what lockfile defines.
- Add `check` command to report available updates by skill.

## Medium-Priority Work

- Add automated tests for:
  - setup flow
  - global/local install resolution
  - update scope behavior (`default`, `--all`, single skill)
  - lockfile read/write and restore flow
- Improve error messaging with consistent actionable hints.

## Nice-to-Have Improvements

- Optional per-skill content hashing for finer change detection.
- Non-interactive `setup` flags for CI/bootstrap workflows.
- Optional JSON output mode for machine-readable command results.
