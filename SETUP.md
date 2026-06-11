# 🏓 paddle-ball — Setup & Initialisation Guide

Everything you need to do **before** adding the source files to your repo.
Follow these steps in order — the whole thing takes about 5 minutes.

---

## Step 1 — Install Go 1.21+

```bash
# macOS (Homebrew)
brew install go

# Ubuntu / Debian
sudo apt install golang-go
# Or use the official installer: https://go.dev/dl/

# Windows — download from https://go.dev/dl/
# Then confirm:
go version   # must print go1.21 or higher
```

---

## Step 2 — Create and enter the project directory

```bash
mkdir paddle-ball
cd paddle-ball
```

> **Important:** the directory name must be `paddle-ball` to match the module path.

---

## Step 3 — Initialise the Go module

```bash
go mod init github.com/subhadeeproy3902/paddle-ball
```

This creates `go.mod`. **Do not run `go mod tidy` yet** — wait until all source
files are in place (Step 7).

---

## Step 4 — Create all subdirectories

```bash
mkdir -p cmd game ui store .github/workflows
```

---

## Step 5 — Paste / copy the source files

Place every file in its correct path. The complete tree:

```
paddle-ball/
├── main.go
├── go.mod                        ← replace the one go mod init created
├── LICENSE
├── Makefile
├── Dockerfile
├── install.sh
├── README.md
├── .goreleaser.yaml
├── .github/
│   └── workflows/
│       └── release.yml
├── cmd/
│   └── root.go
├── game/
│   ├── types.go
│   ├── model.go
│   ├── update.go
│   ├── view.go
│   ├── physics.go
│   ├── particles.go
│   └── scoring.go
├── ui/
│   └── theme.go
└── store/
    └── store.go
```

> **Tip on Windows:** use VS Code, Cursor, or any editor that handles
> Unix line endings (LF). Avoid Notepad.

---

## Step 6 — Make install.sh executable (macOS / Linux)

```bash
chmod +x install.sh
```

---

## Step 7 — Fetch dependencies

```bash
go mod tidy
```

This reads every `import` in your source files, resolves them against the module
proxy, and writes `go.sum`. You need internet access for this step.

Expected output — something like:
```
go: finding module for package github.com/charmbracelet/bubbletea
go: finding module for package github.com/charmbracelet/lipgloss
...
go: added github.com/charmbracelet/bubbletea v0.25.0
go: added github.com/charmbracelet/lipgloss v0.10.0
...
```

---

## Step 8 — Build and test locally

```bash
# Run directly (no binary needed)
go run . play

# Or build a binary first
make build        # → bin/paddle-ball
./bin/paddle-ball play
```

If the game launches and you can see the title screen, everything is working.

---

## Step 9 — Push to GitHub

```bash
git init
git add .
git commit -m "feat: initial paddleball game"
git remote add origin git@github.com:subhadeeproy3902/paddle-ball.git
git push -u origin main
```

---

## Step 10 — Cut your first release

```bash
# Tag v1.0.0
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions picks up the tag, runs GoReleaser, and publishes:
- Pre-built binaries for Linux / macOS / Windows (amd64 + arm64)
- `.deb` and `.rpm` packages
- Snap package
- Docker image on `ghcr.io`
- Homebrew formula (if `HOMEBREW_TOKEN` secret is set)

---

## Optional — Set up Homebrew tap

1. Create a new GitHub repo named `homebrew-paddle-ball` under your account.
2. Add a secret named `HOMEBREW_TOKEN` in the main repo settings
   (Settings → Secrets → Actions) with a GitHub PAT that has `repo` scope.
3. Next time you push a tag, GoReleaser auto-updates the tap.

Users can then install with:
```bash
brew tap subhadeeproy3902/paddle-ball
brew install paddle-ball
```

---

## Optional — Run a local snapshot build (no tag needed)

```bash
# Install goreleaser first
go install github.com/goreleaser/goreleaser@latest

# Build all targets locally (skips publishing)
make snapshot
# → binaries appear in dist/
```

---

## Troubleshooting

| Problem | Fix |
|---|---|
| `package not found` errors | Run `go mod tidy` again |
| `undefined: ui.ThemeCount` | Make sure `ui/theme.go` is saved |
| Game renders garbage characters | Your terminal must support UTF-8 and 256 colours |
| Terminal too small message | Resize to at least 80 × 24 |
| `goreleaser: command not found` | `go install github.com/goreleaser/goreleaser@latest` |
| Windows: colours broken | Use Windows Terminal (not the old cmd.exe) |

---

## Quick Reference

```bash
go run . play               # run the game
go run . scores             # view leaderboard
go run . --help             # all commands
make test                   # run unit tests
make lint                   # golangci-lint
make snapshot               # local release build
git tag vX.Y.Z && git push origin vX.Y.Z   # publish release
```

---

*For issues or questions open a GitHub issue at*
*https://github.com/subhadeeproy3902/paddle-ball/issues*