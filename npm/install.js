#!/usr/bin/env node
"use strict";

const https = require("https");
const http = require("http");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");

const REPO = "MaxMa04/notion-agent-cli";
const BINARY_NAME = "notion-agent";

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  const platformMap = {
    linux: "linux",
    darwin: "darwin",
    win32: "windows",
  };

  const archMap = {
    x64: "amd64",
    arm64: "arm64",
  };

  const mappedPlatform = platformMap[platform];
  const mappedArch = archMap[arch];

  if (!mappedPlatform) {
    throw new Error(`Unsupported platform: ${platform}`);
  }
  if (!mappedArch) {
    throw new Error(`Unsupported architecture: ${arch}`);
  }

  return { os: mappedPlatform, arch: mappedArch };
}

function getVersion() {
  const pkg = JSON.parse(
    fs.readFileSync(path.join(__dirname, "package.json"), "utf8")
  );
  return pkg.version;
}

function getDownloadUrl(version, platform) {
  const ext = platform.os === "windows" ? "zip" : "tar.gz";
  const filename = `notion-agent-cli_${version}_${platform.os}_${platform.arch}.${ext}`;
  return `https://github.com/${REPO}/releases/download/v${version}/${filename}`;
}

function download(url) {
  return new Promise((resolve, reject) => {
    const get = url.startsWith("https:") ? https.get : http.get;
    get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        return download(res.headers.location).then(resolve, reject);
      }
      if (res.statusCode !== 200) {
        return reject(new Error(`Download failed: HTTP ${res.statusCode} for ${url}`));
      }
      const chunks = [];
      res.on("data", (chunk) => chunks.push(chunk));
      res.on("end", () => resolve(Buffer.concat(chunks)));
      res.on("error", reject);
    }).on("error", reject);
  });
}

async function extractTarGz(buffer, destDir) {
  const tmpFile = path.join(os.tmpdir(), `notion-agent-${Date.now()}.tar.gz`);
  fs.writeFileSync(tmpFile, buffer);
  try {
    execSync(`tar xzf "${tmpFile}" -C "${destDir}"`, { stdio: "pipe" });
  } finally {
    fs.unlinkSync(tmpFile);
  }
}

async function extractZip(buffer, destDir) {
  const tmpFile = path.join(os.tmpdir(), `notion-agent-${Date.now()}.zip`);
  fs.writeFileSync(tmpFile, buffer);
  try {
    if (os.platform() === "win32") {
      execSync(
        `powershell -Command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${destDir}' -Force"`,
        { stdio: "pipe" }
      );
    } else {
      execSync(`unzip -o "${tmpFile}" -d "${destDir}"`, { stdio: "pipe" });
    }
  } finally {
    fs.unlinkSync(tmpFile);
  }
}

async function main() {
  const platform = getPlatform();
  const version = getVersion();
  const url = getDownloadUrl(version, platform);

  const binDir = path.join(__dirname, "bin");
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const binaryName =
    platform.os === "windows" ? `${BINARY_NAME}.exe` : BINARY_NAME;
  const binaryPath = path.join(binDir, binaryName);

  // Skip if binary already exists and is the right version
  if (fs.existsSync(binaryPath)) {
    try {
      const output = execSync(`"${binaryPath}" --version`, {
        encoding: "utf8",
        stdio: "pipe",
      }).trim();
      if (output.includes(version)) {
        console.log(`notion-agent v${version} already installed.`);
        return;
      }
    } catch {
      // Binary exists but is broken, re-download
    }
  }

  console.log(`Downloading notion-agent v${version} for ${platform.os}/${platform.arch}...`);

  const buffer = await download(url);

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "notion-agent-"));
  try {
    if (platform.os === "windows") {
      await extractZip(buffer, tmpDir);
    } else {
      await extractTarGz(buffer, tmpDir);
    }

    // Find the binary in the extracted files
    const extractedBinary = path.join(tmpDir, binaryName);
    if (!fs.existsSync(extractedBinary)) {
      throw new Error(
        `Binary not found in archive. Expected: ${binaryName}`
      );
    }

    fs.copyFileSync(extractedBinary, binaryPath);
    if (platform.os !== "windows") {
      fs.chmodSync(binaryPath, 0o755);
    }
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }

  console.log(`notion-agent v${version} installed successfully.`);
}

main().catch((err) => {
  console.error(`Failed to install notion-agent: ${err.message}`);
  console.error(
    "You can download the binary manually from: https://github.com/MaxMa04/notion-agent-cli/releases"
  );
  process.exit(1);
});
