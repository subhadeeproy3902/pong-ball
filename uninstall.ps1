# uninstall.ps1 - removes pong-ball and ALL its data from Windows.
# Usage:  irm https://raw.githubusercontent.com/subhadeeproy3902/pong-ball/main/uninstall.ps1 | iex
$ErrorActionPreference = 'SilentlyContinue'

$dir = "$env:LOCALAPPDATA\Programs\pong-ball"

Write-Host ""
Write-Host "This will permanently remove pong-ball and ALL its data:" -ForegroundColor Yellow
Write-Host "  - the pong-ball binary (every copy on your PATH)"
Write-Host "  - the install folder + its PATH entry"
Write-Host "  - saved scores & settings  ($env:APPDATA\pong-ball)"
Write-Host "  - cached sound files       ($env:TEMP\pong-ball-sfx)"
Write-Host ""
$ans = Read-Host "Remove everything? (Y/N)"

if ($ans -match '^[Yy]') {
    # Stop a running instance so files aren't locked.
    Get-Process pong-ball -ErrorAction SilentlyContinue | Stop-Process -Force

    # Remove every pong-ball.exe on PATH (installer dir, %GOPATH%\bin, etc.).
    foreach ($cmd in @(Get-Command pong-ball -All -ErrorAction SilentlyContinue)) {
        if ($cmd.Source -and (Test-Path -LiteralPath $cmd.Source)) {
            Remove-Item -LiteralPath $cmd.Source -Force
            Write-Host "[pong-ball] removed $($cmd.Source)" -ForegroundColor DarkGray
        }
    }

    # Remove the install directory.
    if (Test-Path -LiteralPath $dir) {
        Remove-Item -LiteralPath $dir -Recurse -Force
        Write-Host "[pong-ball] removed $dir" -ForegroundColor DarkGray
    }

    # Drop the install dir from the user PATH.
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ($userPath) {
        $kept = @($userPath -split ';' | Where-Object { $_ -and $_ -ne $dir })
        [Environment]::SetEnvironmentVariable('Path', ($kept -join ';'), 'User')
    }

    # Remove saved scores/config and the cached sound files.
    Remove-Item -LiteralPath "$env:APPDATA\pong-ball" -Recurse -Force
    Remove-Item -LiteralPath "$env:TEMP\pong-ball-sfx" -Recurse -Force

    Write-Host ""
    Write-Host "[pong-ball] uninstalled. Thanks for playing!" -ForegroundColor Green
    Write-Host "[pong-ball] open a new terminal for the PATH change to take effect." -ForegroundColor DarkGray
    Write-Host "[pong-ball] (installed via Scoop instead? run: scoop uninstall pong-ball)" -ForegroundColor DarkGray
} else {
    Write-Host "[pong-ball] aborted - nothing was removed." -ForegroundColor Cyan
}
