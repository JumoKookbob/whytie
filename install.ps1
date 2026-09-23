$ErrorActionPreference = "Stop"

$AppName = "WhyTie"
$InstallDir = Join-Path $env:LOCALAPPDATA "WhyTie\bin"
$TargetExe = Join-Path $InstallDir "whytie.exe"
$SourceExe = Join-Path $PSScriptRoot "whytie.exe"

Write-Host ""
Write-Host "WhyTie Installer"
Write-Host "==============="
Write-Host ""

if (-not (Test-Path -LiteralPath $SourceExe)) {
    Write-Error "whytie.exe was not found next to install.ps1."
    exit 1
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Copy-Item -LiteralPath $SourceExe -Destination $TargetExe -Force

$userPath = [Environment]::GetEnvironmentVariable(
    "Path",
    [EnvironmentVariableTarget]::User
)

$pathEntries = @()

if (-not [string]::IsNullOrWhiteSpace($userPath)) {
    $pathEntries = $userPath.Split(
        ";",
        [System.StringSplitOptions]::RemoveEmptyEntries
    )
}

$alreadyInPath = $false

foreach ($entry in $pathEntries) {
    if ($entry.TrimEnd("\") -ieq $InstallDir.TrimEnd("\")) {
        $alreadyInPath = $true
        break
    }
}

if (-not $alreadyInPath) {
    $newUserPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
        $InstallDir
    }
    else {
        $userPath.TrimEnd(";") + ";" + $InstallDir
    }

    [Environment]::SetEnvironmentVariable(
        "Path",
        $newUserPath,
        [EnvironmentVariableTarget]::User
    )

    Write-Host "Added WhyTie to your user PATH."
}
else {
    Write-Host "WhyTie is already in your user PATH."
}

# Make WhyTie available in this PowerShell session too.
$currentEntries = $env:Path.Split(
    ";",
    [System.StringSplitOptions]::RemoveEmptyEntries
)

$currentHasPath = $false

foreach ($entry in $currentEntries) {
    if ($entry.TrimEnd("\") -ieq $InstallDir.TrimEnd("\")) {
        $currentHasPath = $true
        break
    }
}

if (-not $currentHasPath) {
    $env:Path = $env:Path.TrimEnd(";") + ";" + $InstallDir
}

Write-Host ""
Write-Host "Installed:"
Write-Host "  $TargetExe"
Write-Host ""

& $TargetExe version

Write-Host ""
Write-Host "WhyTie is ready."
Write-Host ""
Write-Host "Try:"
Write-Host "  whytie help"
Write-Host "  whytie init"
Write-Host "  whytie scan ."
Write-Host "  whytie list"
Write-Host ""