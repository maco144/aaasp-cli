#!/usr/bin/env sh
set -e

REPO="maco144/aaasp-cli"
INSTALL_DIR="/usr/local/bin"
BINARY="aaasp"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect arch
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  arm64 | aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Resolve version (latest release if not set)
VERSION="${AAASP_CLI_VERSION:-}"
ASSET="aaasp_${OS}_${ARCH}"

if [ -z "$VERSION" ]; then
  URL="https://github.com/$REPO/releases/latest/download/$ASSET"
else
  URL="https://github.com/$REPO/releases/download/v$VERSION/$ASSET"
fi

echo "Installing aaasp CLI ${VERSION:-latest} ($OS/$ARCH)..."

TMP="$(mktemp)"
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

# Install (use sudo if needed)
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$INSTALL_DIR/$BINARY"
else
  echo "Requesting sudo to install to $INSTALL_DIR"
  sudo mv "$TMP" "$INSTALL_DIR/$BINARY"
fi

echo ""
echo "aaasp installed to $INSTALL_DIR/$BINARY"
echo ""
echo "Get started:"
echo "  export AAASP_API_KEY=your_key_here"
echo "  aaasp whoami"
