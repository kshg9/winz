# install.ps1 — Bulletproof one-liner installer for Winz
# Usage: irm https://winz.vercel.app/install.ps1 | iex

$ErrorActionPreference = "Stop"

$GITHUB_USER = "axiahq"
$GITHUB_REPO = "winz"
$BINARY_NAME = "winz"

# 1. Ask Windows what architecture it's running
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "386" }

# 2. Build the exact filename based on your goreleaser template
$GITHUB_ASSET_NAME = "${BINARY_NAME}_windows_${arch}.exe"
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"

Write-Step "Detected Architecture: $arch"

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "  ✓ $msg" -ForegroundColor Green }
function Write-Fail($msg) { Write-Host "  ✗ $msg" -ForegroundColor Red; exit 1 }

# ── 1. Download binary (Directly, Bypassing API Rate Limits) ──────────────────
Write-Step "Downloading latest release..."

# The magic GitHub URL that automatically redirects to the latest release asset
$DownloadUrl = "https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/latest/download/$GITHUB_ASSET_NAME"

if (-not (Test-Path -Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
}

$dest = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"

try {
    # -PassThru lets us inspect the headers of the file we just downloaded
    $response = Invoke-WebRequest -Uri $DownloadUrl -OutFile $dest -UseBasicParsing -PassThru

    # ── CAPTIVE PORTAL DEFENSE ──
    # If the network intercepted the request, the content type will be HTML, not an executable
    $contentType = $response.Headers["Content-Type"]
    if ($contentType -match "text/html") {
        Remove-Item -Path $dest -Force # Delete the corrupted HTML file
        Write-Fail "Captive portal detected! Please log into the college network first."
    }

} catch {
    Write-Fail "Download failed. Check your internet connection or ensure '$GITHUB_ASSET_NAME' exists in the latest release."
}

Write-Ok "Downloaded to $dest"

# ── 2. Add to PATH (current user, persistent) ─────────────────────────────────
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
