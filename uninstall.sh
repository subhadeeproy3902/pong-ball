#!/usr/bin/env bash
# uninstall.sh - removes pong-ball and ALL its data.
# Usage:  curl -fsSL https://raw.githubusercontent.com/subhadeeproy3902/pong-ball/main/uninstall.sh | bash
set -uo pipefail

CYAN='\033[0;36m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RESET='\033[0m'
info()    { echo -e "${CYAN}[pong-ball]${RESET} $*"; }
success() { echo -e "${GREEN}[pong-ball]${RESET} $*"; }

echo
echo -e "${YELLOW}This will permanently remove pong-ball and ALL its data:${RESET}"
echo "  - the pong-ball binary (from your PATH)"
echo "  - saved scores & settings  (~/.pong-ball)"
echo "  - cached sound files       (\${TMPDIR:-/tmp}/pong-ball-sfx)"
echo
printf "Remove everything? (Y/N) "
# Read from the terminal even when run via 'curl ... | bash'.
read -r ans </dev/tty

case "$ans" in
  [Yy]*)
    BIN="$(command -v pong-ball 2>/dev/null || true)"
    if [ -n "$BIN" ]; then
      if [ -w "$BIN" ] || [ -w "$(dirname "$BIN")" ]; then
        rm -f "$BIN"
      else
        info "Requesting sudo to remove $BIN…"
        sudo rm -f "$BIN"
      fi
      info "removed $BIN"
    fi

    rm -rf "$HOME/.pong-ball"
    rm -rf "${TMPDIR:-/tmp}/pong-ball-sfx"

    echo
    success "pong-ball uninstalled. Thanks for playing!"
    echo "  (Installed via Homebrew? 'brew uninstall pong-ball' clears the formula receipt too.)"
    ;;
  *)
    info "Aborted - nothing was removed."
    ;;
esac
