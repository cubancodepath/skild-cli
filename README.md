# skild

`skild` is a lightweight Go CLI to install and manage OpenCode skills from a single Git repository.

It is designed to be simple, reproducible, and team-friendly:

- Global configuration in `~/.config/skild/config.json`
- Local and global skill installation targets
- Minimal command set for setup, list, install, and update
- Automated CI and release pipeline via GitHub Actions

## Features

- Install skills from one configurable repository
- Discover skills by scanning directories that contain `SKILL.md`
- Install one skill or all skills
- Update installed skill set from remote repository
- Install locally (project) or globally (`--global`)
- Quiet git output by default (verbose mode available)

## Requirements

- Go 1.22+
- Git

## Build

```bash
go build -o skild ./cmd/skild
```

Optional Ubuntu build (x86_64):

```bash
GOOS=linux GOARCH=amd64 go build -o skild-linux-amd64 ./cmd/skild
```

## Quick Start

1. Build the binary:

```bash
go build -o skild ./cmd/skild
```

2. Run setup:

```bash
./skild setup
```

3. List available skills:

```bash
./skild list
```

4. Install a skill:

```bash
./skild install <skill-name>
```

5. Install globally (OpenCode global path):

```bash
./skild install <skill-name> --global
```

## Commands

- `skild setup` Configure global skild settings
- `skild config` Show resolved configuration
- `skild list` List available skills from configured repository
- `skild install <skill-name>` Install a single skill
- `skild install --all` Install all discovered skills
- `skild install <skill-name> --global` Install globally
- `skild update` Sync repo and reinstall all skills (local)
- `skild update --global` Sync repo and reinstall all skills (global)
- `skild repo-sync` Clone/update cached source repository
- `skild version` Show CLI version

## Configuration

Configuration is stored globally at:

- `~/.config/skild/config.json`

Example:

```json
{
  "repoUrl": "https://github.com/your-org/company-skills.git",
  "repoRef": "main",
  "cacheDir": "~/.cache/skild",
  "rootPath": "skills",
  "openCodeDir": ".opencode/skills",
  "globalOpenCodeDir": "~/.config/opencode/skills",
  "installMode": "copy"
}
```

Notes:

- `globalOpenCodeDir` follows OpenCode conventions:
  - `$XDG_CONFIG_HOME/opencode/skills`
  - fallback: `~/.config/opencode/skills`

## Verbose Mode

By default, git command output is suppressed.

To enable verbose output:

```bash
SKILD_VERBOSE=1 skild list
```

Accepted truthy values: `1`, `true`, `yes`.

## CI/CD and Releases

This project includes:

- CI workflow: vet, test, build on push/PR to `main`
- Release workflow: build tarballs for Linux/macOS on tags `v*`

Release assets:

- `skild_linux_amd64.tar.gz`
- `skild_linux_arm64.tar.gz`
- `skild_darwin_amd64.tar.gz`
- `skild_darwin_arm64.tar.gz`
- `SHA256SUMS`

See `RELEASING.md` for release steps.

## Project Status

MVP in active development.

See `MVP.md` and `ROADMAP.md` for scope and next milestones.
