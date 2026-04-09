#!/bin/sh
set -e

REPO="lucasenlucas/NetForce"
VERSION="V09.04.2026"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="netforce"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    echo "For Windows download manually from: https://github.com/${REPO}/releases"
    exit 1
    ;;
esac

FILENAME="${BIN_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="${BASE_URL}/${FILENAME}"

echo ""
echo "  NetForce Installer"
echo "  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Version  : ${VERSION}"
echo "  Platform : ${OS}/${ARCH}"
echo "  Source   : ${DOWNLOAD_URL}"
echo ""

# Download
TMP="$(mktemp)"
echo "  Downloading..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP"
chmod +x "$TMP"

# Install
echo "  Installing to ${INSTALL_DIR}/${BIN_NAME} ..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "${INSTALL_DIR}/${BIN_NAME}"
else
  sudo mv "$TMP" "${INSTALL_DIR}/${BIN_NAME}"
fi

echo ""
echo "  ✓ NetForce installed successfully!"
echo "  Run: netforce -f explain"
echo ""
echo "  ⚠  For authorized testing only."
echo ""
