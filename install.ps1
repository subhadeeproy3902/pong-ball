# install.ps1 — one-line installer for paddle-ball on Windows.
# Usage:  irm https://raw.githubusercontent.com/subhadeeproy3902/paddle-ball/main/install.ps1 | iex
$ErrorActionPreference = 'Stop'

$repo = 'subhadeeproy3902/paddle-ball'
$dir  = "$env:LOCALAPPDATA\Programs\paddle-ball"
$url  = "https://github.com/$repo/releases/latest/download/paddle-ball_windows_amd64.zip"
$zip  = "$env:TEMP\paddle-ball.zip"

Write-Host "[paddle-ball] downloading latest release..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing

New-Item -ItemType Directory -Force -Path $dir | Out-Null
Expand-Archive -Path $zip -DestinationPath $dir -Force
Remove-Item $zip -Force

# Add to the user PATH (once), and to the current session.
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$dir*") {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$dir", 'User')
}
$env:Path += ";$dir"

Write-Host "[paddle-ball] installed to $dir" -ForegroundColor Green
& "$dir\paddle-ball.exe" version
Write-Host "[paddle-ball] open a NEW terminal, then run: paddle-ball" -ForegroundColor Green
