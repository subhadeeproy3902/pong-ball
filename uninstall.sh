#!/usr/bin/env bash
# uninstall.sh - removes paddle-ball and ALL its data.
# Usage:  curl -fsSL https://raw.githubusercontent.com/subhadeeproy3902/paddle-ball/main/uninstall.sh | bash
set -uo pipefail

CYAN='\033[0;36m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RESET='\033[0m'
info()    { echo -e "${CYAN}[paddle-ball]${RESET} $*"; }
success() { echo -e "${GREEN}[paddle-ball]${RESET} $*"; }

echo
echo -e "${YELLOW}This will permanently remove paddle-ball and ALL its data:${RESET}"
echo "  - the paddle-ball binary (from your PATH)"
echo "  - saved scores & settings  (~/.paddle-ball)"
echo "  - cached sound files       (\${TMPDIR:-/tmp}/paddle-ball-sfx)"
echo
printf "Remove everything? (Y/N) "
# Read from the terminal even when run via 'curl ... | bash'.
read -r ans </dev/tty

case "$ans" in
  [Yy]*)
    BIN="$(command -v paddle-ball 2>/dev/null || true)"
    if [ -n "$BIN" ]; then
      if [ -w "$BIN" ] || [ -w "$(dirname "$BIN")" ]; then
        rm -f "$BIN"
      else
        info "Requesting sudo to remove $BIN…"
        sudo rm -f "$BIN"
      fi
      info "removed $BIN"
    fi

    rm -rf "$HOME/.paddle-ball"
    rm -rf "${TMPDIR:-/tmp}/paddle-ball-sfx"

    echo
    success "paddle-ball uninstalled. Thanks for playing!"
    echo "  (Installed via Homebrew? 'brew uninstall paddle-ball' clears the formula receipt too.)"
    ;;
  *)
    info "Aborted - nothing was removed."
    ;;
esac
