# Install channel status — pong-ball

**The definitive, honest status of every install channel in [PUBLISH.md](PUBLISH.md).**
As of **2026-06-19**, release **v1.0.0**. Numbers that matter: the repo is at **0 stars / 0 forks / 0 watchers** — that single fact is the *only* thing blocking the two "gated" rows below.

Three buckets:

- ✅ **Works today** — a real person can install it right now.
- ⏳ **Submitted** — done on our side, now waiting on a volunteer/automated reviewer. No code work left.
- 🔒 **Gated** — blocked on popularity or a sponsor, *not* on code quality. There's a working fallback for each.

---

## ✅ Working today

| Command | Platform | Notes |
|---|---|---|
| `go install github.com/subhadeeproy3902/pong-ball@latest` | any OS w/ Go | #15 — module is public |
| `curl -sSL https://raw.githubusercontent.com/subhadeeproy3902/pong-ball/main/install.sh \| sh` | macOS · Linux | #16 — auto-detects OS/arch, pulls the latest release |
| `brew install https://raw.githubusercontent.com/subhadeeproy3902/pong-ball/main/Formula/pong-ball.rb` | macOS · Linux | one-command Homebrew, straight from the formula |
| `scoop bucket add pong-ball https://github.com/subhadeeproy3902/pong-ball` then `scoop install pong-ball` | Windows | #10 — installs from this repo's own Scoop bucket |
| Download `.deb` / `.rpm` / `.apk` / `.zip` from [Releases](https://github.com/subhadeeproy3902/pong-ball/releases) | Linux · Windows | built by GoReleaser every tag |

These are what the README, website, and `llms.txt` now advertise — nothing on those pages is aspirational anymore.

---

## ⏳ Submitted — will work once a reviewer approves (no code work left)

| # | Command | Where it's pending | Status | What's left |
|---|---|---|---|---|
| 9 | `choco install pong-ball` | [community.chocolatey.org/packages/pong-ball/1.0.0](https://community.chocolatey.org/packages/pong-ball/1.0.0) | **Published**, pending automated review + moderation | Watch the moderation email. If an automated check flags something, respond in the **Review Comments** box on the package page (not Disqus). Repushing the same 1.0.0 is allowed until approved. |
| 8 | `winget install pong-ball` | [microsoft/winget-pkgs #390476](https://github.com/microsoft/winget-pkgs/pull/390476) | **CLA signed ✓**, validation pipeline running, `New-Package` label, no error labels | Nothing — wait for the bot + a moderator to merge. Only act if a `Needs-Author-Feedback` label appears. |

Both manifests are validated and their checksums match the v1.0.0 release exactly. Future versions auto-publish: Chocolatey via the new `chocolatey` CI job (Windows runner), WinGet would need a fresh PR per version.

---

## 🔒 Gated — blocked on popularity or a sponsor, not on code

Each of these has a **working fallback today** (listed above). The "gate" is a deliberate notability/maintenance bar set by the ecosystem.

| # | Command | The gate | What you can do |
|---|---|---|---|
| 6 | `brew install pong-ball` (no tap) | Homebrew core wants **30 forks / 30 watchers / 75 stars**; self-submissions held higher | Use the one-command formula install (above). Grow stars, then `brew create` + PR to `Homebrew/homebrew-core`. |
| 10 | `scoop install pong-ball` (bare, no bucket) | Scoop **Extras requires 100 stars / 50 forks** (a *required* checkbox on their package-request form) | Use `scoop bucket add` (above). [PR #18079](https://github.com/ScoopInstaller/Extras/pull/18079) can't merge yet — **recommend closing it** and resubmitting once you clear 100★ (their bot/maintainers close under-threshold PRs anyway). |
| 2 | `sudo dnf install pong-ball` | Fedora official review + a packager **sponsor** (months) | **Self-service stepping stone:** [Fedora COPR](https://copr.fedorainfracloud.org) — same-day RPM repo. Users run `dnf copr enable subhadeeproy3902/pong-ball` once, then `dnf install pong-ball`. |
| 5 | `sudo apk add pong-ball` | A merge request to Alpine's `aports`, reviewed by a maintainer (weeks–months) | Write an `APKBUILD`, open an MR to `aports` targeting the `community` branch. The `.apk` already builds on every release. |
| 1 | `sudo apt install pong-ball` | Official Debian needs a **Debian Developer sponsor** — an identity/trust process, months | Self-host an APT repo ([Cloudsmith](https://cloudsmith.com) / [Gemfury](https://fury.co), both GoReleaser-integrated) so it works after one `add-apt-repository`. The `.deb` already builds every release. |
| 4 | `sudo pacman -S pong-ball` | Arch's `core`/`extra` are curated by Arch staff — **no public submission queue** | Publish to the **AUR** (zero review queue): add an `aurs:` block + an `AUR_KEY` secret, then `yay -S pong-ball-bin`. |
| 12 | `sudo snap install pong-ball` | Self-service, but needs the `snapcraft` toolchain in CI | Register the name on snapcraft.io, add a `snap/snapcraft.yaml` + a snapcraft step to CI. Goes live same-day once pushed. |
| 13 | `flatpak install pong-ball` | A submission to Flathub, reviewed by volunteers | Submit a Flatpak manifest to `flathub/flathub`. |

---

## What changed this round

- **Chocolatey** — v1.0.0 published (and repushed with a CDN icon to pre-clear the most common moderator nitpick); added a dedicated `chocolatey` job (Windows runner) for future versions.
- **WinGet** — CLA signed; PR is in normal validation flow.
- **Scoop** — discovered Extras now enforces a 100★/50-fork gate, so the docs switched from the bare command to the working `scoop bucket add` form.
- **Homebrew (macOS)** — collapsed to a single `brew install <formula-url>` command everywhere.

## Bottom line

- **5 channels work today**, **2 are submitted and just need a reviewer**, and the **5 gated ones all have a working fallback**.
- The single highest-leverage move for the gated rows is **stars** — clearing 75★ unlocks Homebrew core and 100★ unlocks Scoop Extras, with no further packaging work.
- Quickest remaining self-service wins: **Fedora COPR**, **AUR** (`-bin`), and **Snap**.
