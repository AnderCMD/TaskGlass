<#
.SYNOPSIS
  Compiles cmd/taskglass/taskglass.rc into resource_windows_amd64.syso.

.DESCRIPTION
  The Go linker picks the .syso up automatically, embedding the icon, version
  information and manifest into taskglass.exe. The .syso is committed, so this
  only needs to run when the .rc, the manifest or the icon change.
  Needs windres (MinGW-w64): winget install MartinStorsjo.LLVM-MinGW.UCRT
#>
$ErrorActionPreference = "Stop"
$dir = Join-Path $PSScriptRoot "..\cmd\taskglass"

if (-not (Get-Command x86_64-w64-mingw32-windres, windres -ErrorAction SilentlyContinue)) {
    $pkg = Get-ChildItem "$env:LOCALAPPDATA\Microsoft\WinGet\Packages" -Directory -ErrorAction SilentlyContinue |
        Where-Object Name -like "*LLVM-MinGW*" | Select-Object -First 1
    if ($pkg) {
        $bin = Join-Path (Get-ChildItem $pkg.FullName -Directory | Select-Object -First 1).FullName "bin"
        $env:Path = "$bin;$env:Path"
    }
}
$windres = @("x86_64-w64-mingw32-windres", "windres") |
    Where-Object { Get-Command $_ -ErrorAction SilentlyContinue } | Select-Object -First 1
if (-not $windres) { throw "windres not found. Install MinGW-w64: winget install MartinStorsjo.LLVM-MinGW.UCRT" }

Push-Location $dir
try {
    & $windres --target=pe-x86-64 -O coff -i taskglass.rc -o resource_windows_amd64.syso
    if ($LASTEXITCODE -ne 0) { throw "windres failed" }
} finally { Pop-Location }
Write-Host "wrote cmd/taskglass/resource_windows_amd64.syso"
