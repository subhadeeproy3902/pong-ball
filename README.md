# 🏓 paddle-ball

> A physics-based neon terminal paddleball game — smooth delta-time movement,
> particle effects, power-ups, persistent score history, and four visual themes.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Release](https://img.shields.io/github/v/release/subhadeeproy3902/paddle-ball?style=flat&color=00FFFF)](https://github.com/subhadeeproy3902/paddle-ball/releases)
[![License](https://img.shields.io/badge/license-MIT-C3E88D?style=flat)](LICENSE)
[![Twitter](https://img.shields.io/badge/Twitter-mvp__Subha-1DA1F2?style=flat&logo=twitter)](https://twitter.com/mvp_Subha)

---

## ✨ Features

- **Real physics** — delta-time ball movement, paddle angle + spin transfer
- **5 difficulty phases** — Warm Up → INSANE, auto-escalating by score
- **4 game modes** — Classic, Arcade (lives + power-ups), Zen, Time Trial
- **6 power-ups** — Wide Paddle, Slow Mo, Fire Paddle, Iron Shield, Ghost Ball, Bomb
- **Particle system** — wall impacts, paddle flash, score pops, game-over explosion
- **Score history** — persistent JSON store at `~/.paddle-ball/scores.json`
- **4 color themes** — Neon Arcade, Monochrome, Sunset, Ocean Night
- **Single binary** — zero runtime dependencies, one command install

---

## 🚀 Install

### One-line (macOS / Linux)
```bash
curl -fsSL https://raw.githubusercontent.com/subhadeeproy3902/paddle-ball/main/install.sh | bash
```

### Homebrew (macOS / Linux)
```bash
brew tap subhadeeproy3902/paddle-ball
brew install paddle-ball
```

### apt (Debian / Ubuntu)
```bash
curl -fsSL https://packages.paddle-ball.dev/install.sh | sudo bash
# or manually:
sudo apt install paddle-ball
```

### Snap (Linux)
```bash
sudo snap install paddle-ball
```

### Scoop (Windows)
```powershell
scoop bucket add paddle-ball https://github.com/subhadeeproy3902/scoop-paddle-ball
scoop install paddle-ball
```

### WinGet (Windows)
```powershell
winget install subhadeeproy3902.paddle-ball
```

### Go install
```bash
go install github.com/subhadeeproy3902/paddle-ball@latest
```

### Docker (no install)
```bash
docker run --rm -it ghcr.io/subhadeeproy3902/paddle-ball:latest
```

---

## 🎮 Usage

```
paddle-ball                     # Title screen
paddle-ball play                # Jump straight in (Classic)
paddle-ball play --mode arcade  # Arcade mode with power-ups
paddle-ball play --mode zen     # Infinite lives
paddle-ball play --mode timed   # 60-second blitz
paddle-ball scores              # View leaderboard
paddle-ball scores --all        # Full history
paddle-ball scores --json       # Raw JSON dump
paddle-ball reset               # Wipe saved scores
paddle-ball version             # Version info
```

## 🕹  Controls

| Key | Action |
|---|---|
| `↑` `W` `K` | Move paddle up |
| `↓` `S` `J` | Move paddle down |
| `P` / `Space` | Pause / Resume |
| `T` | Cycle color theme |
| `?` | Help overlay |
| `R` | Restart (pause / game over) |
| `Q` / `Ctrl+C` | Quit |

---

## 🏆 Difficulty Phases

| Phase | Score | Speed | Paddle |
|---|---|---|---|
| 🌱 Warm Up | 0+ | 100 % | 7 cells |
| 🔥 Heating Up | 10+ | 125 % | 6 cells |
| 💥 On Fire | 25+ | 155 % | 5 cells |
| ⚡ Blazing | 50+ | 190 % | 4 cells |
| 🏆 INSANE | 100+ | 235 % | 3 cells |

---

## 🛠  Development

```bash
git clone https://github.com/subhadeeproy3902/paddle-ball
cd paddle-ball
go mod tidy
make run          # run game
make build        # binary → bin/paddle-ball
make test         # unit tests
make snapshot     # local release build (goreleaser)
```

### Release
```bash
git tag v1.0.0
git push origin v1.0.0   # triggers GitHub Actions → goreleaser
```

---

## 📁 Project Structure

```
paddle-ball/
├── main.go                 Entry point (version injection)
├── cmd/root.go             All Cobra CLI commands
├── game/
│   ├── types.go            All game types and enums
│   ├── model.go            Bubbletea Model struct + NewModel
│   ├── update.go           Bubbletea Update (game loop)
│   ├── view.go             Bubbletea View (full renderer)
│   ├── physics.go          Ball/paddle collision math
│   ├── particles.go        Particle system
│   └── scoring.go          Points, streaks, phase checks
├── ui/theme.go             4 color themes + lipgloss helpers
├── store/store.go          Score persistence (atomic JSON)
├── .goreleaser.yaml        Cross-compile + publish pipeline
└── Makefile
```

---

## 🔗 Links

- **GitHub** — [github.com/subhadeeproy3902/paddle-ball](https://github.com/subhadeeproy3902/paddle-ball)
- **Twitter** — [@mvp_Subha](https://twitter.com/mvp_Subha)
- **LinkedIn** — [subhadeep3902](https://linkedin.com/in/subhadeep3902)

---

## 📄 License

MIT © [Subhadeep Roy](https://github.com/subhadeeproy3902)

*Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) +
[Lip Gloss](https://github.com/charmbracelet/lipgloss) from the
[Charm](https://charm.sh) ecosystem.*