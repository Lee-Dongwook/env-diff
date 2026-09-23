# env-diff

## Runtime Env Report

Compare your local development environment with the one used in CI.

`runtime-env-report` captures OS, architecture, toolchain versions, and environment variable names. It does not store environment variable values.

## Install

```bash
npm install --save-dev runtime-env-report
```

After installation, run the CLI as `envdiff`:

```bash
npx envdiff --help
```

To try it without adding it to your project:

```bash
npx --yes --package=runtime-env-report envdiff --help
```

## Usage

Capture your local environment:

```bash
npx envdiff capture --out local.envdiff.json
```

Capture the CI environment:

```yaml
- name: Capture CI environment
  run: npx --yes --package=runtime-env-report envdiff capture --out ci.envdiff.json
```

Compare the snapshots:

```bash
npx envdiff diff local.envdiff.json ci.envdiff.json
```

Make the command exit with an error when differences are found:

```bash
npx envdiff diff --fail-on-diff local.envdiff.json ci.envdiff.json
```

## Supported platforms

The npm package includes prebuilt binaries for:

- macOS: x64 and arm64
- Linux: x64 and arm64
- Windows: x64 and arm64

## Privacy

Snapshot files contain environment variable names and whether their values are empty. They do not contain environment variable values. Variable names can still reveal information about your setup, so review snapshots before sharing them.

## License

MIT
