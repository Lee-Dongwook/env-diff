#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");

const binaries = {
  "darwin-arm64": "envdiff-darwin-arm64",
  "darwin-x64": "envdiff-darwin-x64",
  "linux-arm64": "envdiff-linux-arm64",
  "linux-x64": "envdiff-linux-x64",
  "win32-arm64": "envdiff-win32-arm64.exe",
  "win32-x64": "envdiff-win32-x64.exe",
};

const platformKey = `${process.platform}-${process.arch}`;
const binaryName = binaries[platformKey];

if (!binaryName) {
  console.error(`envdiff does not support this platform: ${platformKey}`);
  process.exit(1);
}

const binaryPath = path.join(__dirname, "..", "vendor", "bin", binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error(`envdiff binary is missing for ${platformKey}`);
  console.error(`Expected: ${binaryPath}`);
  process.exit(1);
}

const result = spawnSync(binaryPath, process.argv.slice(2), {
  stdio: "inherit",
});

if (result.error) {
  console.error(`Failed to start envdiff: ${result.error.message}`);
  process.exit(1);
}

if (result.signal) {
  console.error(`envdiff was terminated by signal: ${result.signal}`);
  process.exit(1);
}

process.exit(result.status ?? 1);
