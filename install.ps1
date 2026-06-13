# install.ps1 - one-line installer AND updater for pong-ball on Windows.
# Usage:  irm https://raw.githubusercontent.com/subhadeeproy3902/pong-ball/main/install.ps1 | iex
#
# Running it again upgrades in place: it stops a running game, removes older
# copies that would shadow the new one on PATH (the usual "it won't update"
# trap), installs the latest release, and puts itself FIRST on PATH so the new
# binary always wins.
$ErrorActionPreference = 'Stop'

$repo = 'subhadeeproy3902/pong-ball'
$dir  = "$env:LOCALAPPDATA\Programs\pong-ball"
$exe  = Join-Path $dir 'pong-ball.exe'
$url  = "https://github.com/$repo/releases/latest/download/pong-ball_windows_amd64.zip"
$zip  = Join-Path $env:TEMP 'pong-ball.zip'

Write-Host "[pong-ball] installing the latest release..." -ForegroundColor Cyan

# Stop any running instance so the binary isn't locked while we replace it.
Get-Process pong-ball -ErrorAction SilentlyContinue | Stop-Process -Force

# Remove older copies anywhere else on PATH that would shadow the new install.
# A stale "go install" build in %GOPATH%\bin is the most common culprit: it sits
# earlier on PATH, so `pong-ball` keeps launching the old version forever.
$target = [System.IO.Path]::GetFullPath($exe)
foreach ($cmd in @(Get-Command pong-ball -All -ErrorAction SilentlyContinue)) {
    $src = $cmd.Source
    if ($src -and [System.IO.Path]::GetFullPath($src) -ne $target) {
        try {
            Remove-Item -LiteralPath $src -Force -ErrorAction Stop
            Write-Host "[pong-ball] removed stale copy: $src" -ForegroundColor DarkYellow
        } catch {
            Write-Host "[pong-ball] heads up: an old copy at $src is locked - close it and remove it manually" -ForegroundColor DarkYellow
        }
    }
}

# Download + extract the latest release.
Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
New-Item -ItemType Directory -Force -Path $dir | Out-Null
Expand-Archive -Path $zip -DestinationPath $dir -Force
Remove-Item $zip -Force

# Put our dir FIRST on the user PATH (and de-dupe) so it always wins, then
# mirror that into the current session.
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$rest = @($userPath -split ';' | Where-Object { $_ -and $_ -ne $dir })
[Environment]::SetEnvironmentVariable('Path', (@($dir) + $rest) -join ';', 'User')
$env:Path = "$dir;" + ($env:Path -split ';' | Where-Object { $_ -and $_ -ne $dir } | Select-Object -Unique) -join ';'

Write-Host "[pong-ball] installed to $dir" -ForegroundColor Green
& $exe version

# Sanity check: make sure `pong-ball` now resolves to what we just installed.
$resolved = (Get-Command pong-ball -ErrorAction SilentlyContinue).Source
if ($resolved -and ([System.IO.Path]::GetFullPath($resolved) -eq $target)) {
    Write-Host "[pong-ball] ready - just run: pong-ball" -ForegroundColor Green
} else {
    Write-Host "[pong-ball] installed, but PATH still resolves to: $resolved" -ForegroundColor DarkYellow
    Write-Host "[pong-ball] open a NEW terminal, then run: pong-ball" -ForegroundColor DarkYellow
}
