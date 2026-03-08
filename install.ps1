# install.ps1 — One-liner installer for <yourapp>
# Usage: irm https://raw.githubusercontent.com/YOU/REPO/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$GITHUB_USER = "kshg9"
$GITHUB_REPO = "winz"
$BINARY_NAME = "winz"          # what you want it called in PATH
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"   # no admin needed

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "  ✓ $msg" -ForegroundColor Green }
function Write-Fail($msg) { Write-Host "  ✗ $msg" -ForegroundColor Red; exit 1 }

# ── 1. Fetch latest release info from GitHub API ──────────────────────────────
Write-Step "Fetching latest release..."
try {
    $release = Invoke-RestMethod "https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/releases/latest"
} catch {
    Write-Fail "Could not reach GitHub API. Check your connection."
}

$asset = $release.assets | Where-Object { $_.name -like "*windows*amd64*" -or $_.name -like "*windows*" } | Select-Object -First 1
if (-not $asset) { Write-Fail "No Windows binary found in latest release assets." }

Write-Ok "Found: $($asset.name) (release $($release.tag_name))"

# ── 2. Set execution policy (current user only, no admin) ─────────────────────
Write-Step "Setting execution policy..."
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
Write-Ok "Execution policy set (CurrentUser only)"

# ── 3. Download binary ────────────────────────────────────────────────────────
Write-Step "Downloading binary..."
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
$dest = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"

try {
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $dest -UseBasicParsing
} catch {
    Write-Fail "Download failed: $_"
}
Write-Ok "Downloaded to $dest"

# ── 4. Add to PATH (current user, persistent) ─────────────────────────────────
Write-Step "Adding to PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$INSTALL_DIR", "User")
    Write-Ok "Added $INSTALL_DIR to your PATH"
} else {
    Write-Ok "Already in PATH"
}

# also patch current session so it works immediately
$env:PATH += ";$INSTALL_DIR"

# ── Done ──────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  $BINARY_NAME installed! Run: $BINARY_NAME" -ForegroundColor Yellow
Write-Host "  (If PATH doesn't work yet, restart your terminal)" -ForegroundColor Gray
Write-Host ""
