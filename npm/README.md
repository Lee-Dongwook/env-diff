# env-diff

Compare the development environment on your machine with the one used in CI.

`env-diff` captures OS, architecture, toolchain versions, and environment variable names. It does not store environment variable values.

## Install

```bash
npm install --save-dev env-diff
```

Or run it without installing:

```bash
npx env-diff --help
```

## Usage

Capture your local environment:

```bash
npx envdiff capture --out local.envdiff.json
```

Capture the CI environment by adding this step to your workflow:

```yaml
- name: Capture CI environment
  run: npx envdiff capture --out ci.envdiff.json
```

Compare the snapshots:

```bash
npx envdiff diff local.envdiff.json ci.envdiff.json
```

To make the command exit with an error when differences are found:

```bash
npx envdiff diff --fail-on-diff local.envdiff.json ci.envdiff.json
```

## Supported platforms

The npm package includes prebuilt binaries for:

- macOS: x64, arm64
- Linux: x64, arm64
- Windows: x64, arm64

## Privacy

Snapshot files contain environment variable names and whether their values are empty. They do not contain the values themselves. Review the snapshot before sharing it, since variable names can still reveal information about your setup.

## License

MIT
