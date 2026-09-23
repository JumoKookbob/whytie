$ErrorActionPreference = "Stop"

$InstallDir = Join-Path $env:LOCALAPPDATA "WhyTie\bin"
$InstallRoot = Join-Path $env:LOCALAPPDATA "WhyTie"

Write-Host ""
Write-Host "WhyTie Uninstaller"
Write-Host "=================="
Write-Host ""

$userPath = [Environment]::GetEnvironmentVariable(
    "Path",
    [EnvironmentVariableTarget]::User
)

if (-not [string]::IsNullOrWhiteSpace($userPath)) {
    $pathEntries = $userPath.Split(
        ";",
        [System.StringSplitOptions]::RemoveEmptyEntries
    )

    $newEntries = @(
        $pathEntries | Where-Object {
            $_.TrimEnd("\") -ine $InstallDir.TrimEnd("\")
        }
    )

    if ($newEntries.Count -ne $pathEntries.Count) {
        $newUserPath = $newEntries -join ";"

        [Environment]::SetEnvironmentVariable(
            "Path",
            $newUserPath,
            [EnvironmentVariableTarget]::User
        )

        Write-Host "Removed WhyTie from your user PATH."
    }
    else {
        Write-Host "WhyTie was not present in your user PATH."
    }
}

if (Test-Path -LiteralPath $InstallRoot) {
    Remove-Item -LiteralPath $InstallRoot -Recurse -Force
    Write-Host "Removed:"
    Write-Host "  $InstallRoot"
}
else {
    Write-Host "WhyTie installation directory was not found."
}

Write-Host ""
Write-Host "WhyTie has been uninstalled."
Write-Host ""
Write-Host "Existing project .whytie directories were not deleted."
Write-Host ""