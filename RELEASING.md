# Releasing skild

`skild` releases are automated with GitHub Actions and triggered by semantic version tags.

## Supported Release Targets

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`

Each release publishes:

- `skild_<os>_<arch>.tar.gz`
- `SHA256SUMS`

## Create a Release

1. Ensure `main` is green in CI.
2. Create and push a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

3. Wait for the `Release` workflow to finish.
4. Verify assets in GitHub Releases.

## Verify Downloaded Artifacts

```bash
shasum -a 256 -c SHA256SUMS
```
