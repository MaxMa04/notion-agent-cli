#!/usr/bin/env node
"use strict";

const { execFileSync, execSync } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");

const binaryName =
  os.platform() === "win32" ? "notion-agent.exe" : "notion-agent";
const binaryPath = path.join(__dirname, "bin", binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error("notion-agent binary not found. Running install...");
  try {
    execSync("node install.js", { cwd: __dirname, stdio: "inherit" });
  } catch {
    // install.js already prints error messages
  }
  if (!fs.existsSync(binaryPath)) {
    console.error(
      "Installation failed. Please run: npm install -g @vibelabsio/notion-agent-cli"
    );
    process.exit(1);
  }
}

try {
  execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  process.exit(err.status || 1);
}
