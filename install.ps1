$ErrorActionPreference = "Stop"

$GITHUB_USER = "axiahq"
$GITHUB_REPO = "winz"
$BINARY_NAME = "_"
$INSTALL_DIR = "$env:USERPROFILE\.local\bin"

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "  ✓ $msg" -ForegroundColor Green }
function Write-Fail($msg) { Write-Host "  ✗ $msg" -ForegroundColor Red; exit 1 }

$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "386" }
$GITHUB_ASSET_NAME = "winz_windows_${arch}.exe"

Write-Step "Detected Architecture: $arch"
Write-Step "Downloading latest release..."

$DownloadUrl = "https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/latest/download/$GITHUB_ASSET_NAME"

if (-not (Test-Path -Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
}

$dest = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"

try {
    $response = Invoke-WebRequest -Uri $DownloadUrl -OutFile $dest -UseBasicParsing -PassThru
    if ($response.Headers["Content-Type"] -match "text/html") {
        Remove-Item -Path $dest -Force
        Write-Fail "Captive portal detected! Please log into the network first."
    }
} catch {
    Write-Fail "Download failed. GitHub returned an error (likely 404 Not Found)."
}

Write-Ok "Downloaded to $dest"

Write-Step "Adding to PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$INSTALL_DIR", "User")
    Write-Ok "Added $INSTALL_DIR to your PATH"
} else {
    Write-Ok "Already in PATH"
}

$env:PATH += ";$INSTALL_DIR"

Write-Host "`n  Installed! Run: $BINARY_NAME`n" -ForegroundColor Yellow
