#!/usr/bin/env bash
# install.sh — one-line installer for paddle-ball
# Usage:  curl -fsSL https://raw.githubusercontent.com/subhadeeproy3902/paddle-ball/main/install.sh | bash
set -euo pipefail

REPO="subhadeeproy3902/paddle-ball"
BINARY="paddle-ball"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# ── colour helpers ──────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; CYAN='\033[0;36m'; BOLD='\033[1m'; RESET='\033[0m'
info()    { echo -e "${CYAN}[paddle-ball]${RESET} $*"; }
success() { echo -e "${GREEN}[paddle-ball]${RESET} $*"; }
error()   { echo -e "${RED}[paddle-ball] ERROR:${RESET} $*" >&2; exit 1; }

# ── detect OS ───────────────────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)   OS="linux"   ;;
  darwin)  OS="darwin"  ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *)       error "Unsupported OS: $OS" ;;
esac

# ── detect arch ─────────────────────────────────────────────────────────────
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)             error "Unsupported architecture: $ARCH" ;;
esac

# ── fetch latest tag ────────────────────────────────────────────────────────
info "Fetching latest release…"
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
[ -z "$LATEST" ] && error "Could not determine latest version. Check your internet connection."
info "Latest version: ${BOLD}${LATEST}${RESET}"

# ── build download URL ───────────────────────────────────────────────────────
EXT="tar.gz"
[ "$OS" = "windows" ] && EXT="zip"
FILENAME="${BINARY}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

# ── download ────────────────────────────────────────────────────────────────
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

info "Downloading ${FILENAME}…"
curl -fsSL "$URL" -o "$TMP/$FILENAME" \
  || error "Download failed. URL: $URL"

# ── extract ─────────────────────────────────────────────────────────────────
info "Extracting…"
if [ "$EXT" = "zip" ]; then
  command -v unzip >/dev/null 2>&1 || error "unzip not found — install it first"
  unzip -q "$TMP/$FILENAME" -d "$TMP"
else
  tar -xzf "$TMP/$FILENAME" -C "$TMP"
fi

# ── install ─────────────────────────────────────────────────────────────────
BIN_PATH="$TMP/$BINARY"
[ ! -f "$BIN_PATH" ] && BIN_PATH=$(find "$TMP" -name "$BINARY" -type f | head -1)
[ -z "$BIN_PATH" ] && error "Binary not found in archive"

chmod +x "$BIN_PATH"

if [ -w "$INSTALL_DIR" ]; then
  mv "$BIN_PATH" "$INSTALL_DIR/$BINARY"
else
  info "Requesting sudo to install to $INSTALL_DIR…"
  sudo mv "$BIN_PATH" "$INSTALL_DIR/$BINARY"
fi

# ── verify ──────────────────────────────────────────────────────────────────
INSTALLED_VER="$("$INSTALL_DIR/$BINARY" version 2>/dev/null | head -1 || echo '?')"
success "Installed ${BOLD}${BINARY}${RESET} → ${INSTALL_DIR}/${BINARY}"
success "Version: ${INSTALLED_VER}"
echo
echo -e "  Run ${CYAN}paddle-ball${RESET} to play!"
echo -e "  Run ${CYAN}paddle-ball --help${RESET} for all commands."