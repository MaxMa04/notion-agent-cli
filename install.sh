#!/bin/sh
set -e

REPO="MaxMa04/notion-agent-cli"
BINARY="notion-agent"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *)
    echo "Error: Unsupported OS: $OS"
    echo "Please download manually from https://github.com/${REPO}/releases"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture: $ARCH"
    echo "Please download manually from https://github.com/${REPO}/releases"
    exit 1
    ;;
esac

# Get latest version
echo "Fetching latest version..."
if command -v curl > /dev/null 2>&1; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')
elif command -v wget > /dev/null 2>&1; then
  VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')
else
  echo "Error: Neither curl nor wget found. Please install one of them."
  exit 1
fi

if [ -z "$VERSION" ]; then
  echo "Error: Could not determine latest version."
  echo "Please download manually from https://github.com/${REPO}/releases"
  exit 1
fi

# Construct download URL
FILENAME="${BINARY}-cli_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

# Determine install directory
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "Installing ${BINARY} v${VERSION} for ${OS}/${ARCH}..."
echo "  From: ${URL}"
echo "  To:   ${INSTALL_DIR}/${BINARY}"

# Create temp directory
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Download
if command -v curl > /dev/null 2>&1; then
  curl -fsSL "$URL" -o "$TMPDIR/archive.tar.gz"
else
  wget -qO "$TMPDIR/archive.tar.gz" "$URL"
fi

# Extract
tar xzf "$TMPDIR/archive.tar.gz" -C "$TMPDIR"

# Install
if [ -w "$INSTALL_DIR" ]; then
  install -m 755 "$TMPDIR/${BINARY}" "$INSTALL_DIR/${BINARY}"
else
  echo "Note: ${INSTALL_DIR} requires elevated permissions."
  sudo install -m 755 "$TMPDIR/${BINARY}" "$INSTALL_DIR/${BINARY}"
fi

echo ""
echo "${BINARY} v${VERSION} installed successfully!"
"${INSTALL_DIR}/${BINARY}" --version 2>/dev/null || true
