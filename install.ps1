#Requires -Version 5.1
<#
.SYNOPSIS
    jd - Claude Code configuration manager installation script for Windows

.DESCRIPTION
    Downloads and installs jd to the specified directory.

.PARAMETER InstallDir
    Installation directory (default: $env:LOCALAPPDATA\Programs\jd)

.PARAMETER Version
    Specific version to install (default: latest)

.EXAMPLE
    # Install latest version
    irm https://raw.githubusercontent.com/itda-work/itda-jindo/main/install.ps1 | iex

.EXAMPLE
    # Install to custom directory
    $env:INSTALL_DIR = "C:\tools"; irm https://raw.githubusercontent.com/itda-work/itda-jindo/main/install.ps1 | iex

.EXAMPLE
    # Install specific version
    $env:VERSION = "v0.1.0"; irm https://raw.githubusercontent.com/itda-work/itda-jindo/main/install.ps1 | iex
#>

$ErrorActionPreference = "Stop"

$Repo = "itda-work/itda-jindo"
$Binary = "jd"
$DefaultInstallDir = Join-Path $env:LOCALAPPDATA "Programs\jd"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { $DefaultInstallDir }

function Write-Info {
    param([string]$Message)
    Write-Host "INFO: " -ForegroundColor Blue -NoNewline
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "SUCCESS: " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "WARNING: " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "ERROR: " -ForegroundColor Red -NoNewline
    Write-Host $Message
    exit 1
}

function Get-Architecture {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    switch ($arch) {
        "X64" { return "amd64" }
        "Arm64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $arch" }
    }
}

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to get latest version: $_"
    }
}

function Install-Jd {
    Write-Host ""
    Write-Host "  +--------------------------------------+"
    Write-Host "  |                                      |"
    Write-Host "  |   jd - Claude Code Config Manager    |"
    Write-Host "  |                                      |"
    Write-Host "  +--------------------------------------+"
    Write-Host ""

    $arch = Get-Architecture
    Write-Info "Detected Architecture: $arch"

    # Get version
    $version = if ($env:VERSION) {
        Write-Info "Using specified version: $env:VERSION"
        $env:VERSION
    }
    else {
        Write-Info "Fetching latest version..."
        $v = Get-LatestVersion
        Write-Info "Latest version: $v"
        $v
    }

    # Construct download URL
    $filename = "$Binary-windows-$arch.exe"
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/$filename"

    Write-Info "Downloading from: $downloadUrl"

    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        Write-Info "Creating directory: $InstallDir"
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $binaryPath = Join-Path $InstallDir "$Binary.exe"

    # Download
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath -UseBasicParsing
    }
    catch {
        Write-Error "Failed to download $downloadUrl : $_"
    }

    Write-Success "Installed $Binary to $binaryPath"

    # Check if install directory is in PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        Write-Warning "$InstallDir is not in your PATH"

        $addToPath = Read-Host "Add to PATH? (Y/n)"
        if ($addToPath -eq "" -or $addToPath -match "^[Yy]") {
            $newPath = "$userPath;$InstallDir"
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
            $env:Path = "$env:Path;$InstallDir"
            Write-Success "Added $InstallDir to PATH"
            Write-Info "Please restart your terminal for the changes to take effect"
        }
        else {
            Write-Host ""
            Write-Host "To add to PATH manually, run:"
            Write-Host ""
            Write-Host "  `$env:Path += `";$InstallDir`""
            Write-Host ""
            Write-Host "Or add it permanently via System Properties > Environment Variables"
            Write-Host ""
        }
    }

    # Verify installation
    try {
        $versionOutput = & $binaryPath --version
        Write-Info "Version: $versionOutput"
    }
    catch {
        # Ignore version check failure
    }

    Write-Host ""
    Write-Success "Installation complete!"
    Write-Host ""
    Write-Host "Get started with:"
    Write-Host "  $Binary --help"
    Write-Host ""
}

Install-Jd
