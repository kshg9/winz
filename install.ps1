$ErrorActionPreference = "Stop"

$GITHUB_USER = "kshg9"
$GITHUB_REPO = "winz"
$BINARY_NAME = "winz"
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "  ✓ $msg" -ForegroundColor Green }
function Write-Fail($msg) { Write-Host "  ✗ $msg" -ForegroundColor Red; exit 1 }

# ── 1. Fetch latest release info from GitHub API ──────────────────────────────
Write-Step "Fetching latest release..."
try {
    # WARNING: Subject to 60 requests/hr rate limit from shared IPs
    $release = Invoke-RestMethod "https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/releases/latest"
} catch {
    Write-Fail "Could not reach GitHub API. You might be rate-limited. Try again later."
}

$asset = $release.assets | Where-Object { $_.name -like "*windows*amd64*" -or $_.name -like "*windows*" } | Select-Object -First 1
if (-not $asset) { Write-Fail "No Windows binary found in latest release assets." }

Write-Ok "Found: $($asset.name) (release $($release.tag_name))"

# ── 2. Download binary ────────────────────────────────────────────────────────
Write-Step "Downloading binary..."
if (-not (Test-Path -Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
}

$dest = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"

try {
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $dest -UseBasicParsing
} catch {
    Write-Fail "Download failed: $_"
}
Write-Ok "Downloaded to $dest"

# ── 3. Add to PATH (current user, persistent) ─────────────────────────────────
Write-Step "Adding to PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$INSTALL_DIR", "User")
    Write-Ok "Added $INSTALL_DIR to your PATH"
} else {
    Write-Ok "Already in PATH"
}

# Patch current session so it works immediately without restarting the terminal
$env:PATH += ";$INSTALL_DIR"

# ── Done ──────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  $BINARY_NAME installed! Run: $BINARY_NAME" -ForegroundColor Yellow
Write-Host ""
