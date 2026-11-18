<#
.SYNOPSIS
  Builds dist/TaskGlass-<version>-Setup.exe with Inno Setup.

.DESCRIPTION
  Expects dist/taskglass.exe and dist/taskglass_tap.dll to exist (build them
  first). Needs Inno Setup 6: winget install JRSoftware.InnoSetup
#>
param([Parameter(Mandatory)][string]$Version)
$ErrorActionPreference = "Stop"
$root = Resolve-Path (Join-Path $PSScriptRoot "..")

$iscc = @(
    (Get-Command ISCC.exe -ErrorAction SilentlyContinue).Source,
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe",
    "$env:LOCALAPPDATA\Programs\Inno Setup 6\ISCC.exe"
) | Where-Object { $_ -and (Test-Path $_) } | Select-Object -First 1
if (-not $iscc) { throw "Inno Setup 6 not found. Install it: winget install JRSoftware.InnoSetup" }
foreach ($f in "taskglass.exe", "taskglass_tap.dll") {
    if (-not (Test-Path (Join-Path $root "dist\$f"))) { throw "dist\$f is missing; build it first" }
}
& $iscc "/DAppVersion=$Version" (Join-Path $root "installer\taskglass.iss")
if ($LASTEXITCODE -ne 0) { throw "ISCC failed" }
Write-Host "built dist\TaskGlass-$Version-Setup.exe"
