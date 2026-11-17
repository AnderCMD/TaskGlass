<#
.SYNOPSIS
  Builds native/tap/tap.c into dist/taskglass_tap.dll.

.DESCRIPTION
  Needs a MinGW-w64 toolchain (gcc or clang) on PATH.
  Install one with: winget install MartinStorsjo.LLVM-MinGW.UCRT
#>
param(
    [string]$Out = (Join-Path $PSScriptRoot "..\dist\taskglass_tap.dll")
)
$ErrorActionPreference = "Stop"
$src = Join-Path $PSScriptRoot "..\native\tap"

# winget installs LLVM-MinGW as a portable package that is not always on PATH yet.
if (-not (Get-Command x86_64-w64-mingw32-clang, x86_64-w64-mingw32-gcc, gcc -ErrorAction SilentlyContinue)) {
    $pkg = Get-ChildItem "$env:LOCALAPPDATA\Microsoft\WinGet\Packages" -Directory -ErrorAction SilentlyContinue |
        Where-Object Name -like "*LLVM-MinGW*" | Select-Object -First 1
    if ($pkg) {
        $bin = Join-Path (Get-ChildItem $pkg.FullName -Directory | Select-Object -First 1).FullName "bin"
        $env:Path = "$bin;$env:Path"
    }
}
$cc = @("x86_64-w64-mingw32-clang", "x86_64-w64-mingw32-gcc", "gcc") |
    Where-Object { Get-Command $_ -ErrorAction SilentlyContinue } | Select-Object -First 1
if (-not $cc) { throw "No MinGW-w64 compiler found. Install one: winget install MartinStorsjo.LLVM-MinGW.UCRT" }

New-Item -ItemType Directory -Force (Split-Path $Out) | Out-Null
& $cc -shared -O2 -Wall -Wextra -o $Out "$src\tap.c" "$src\tap.def" `
    -lole32 -luuid -lruntimeobject -loleaut32 -static-libgcc
if ($LASTEXITCODE -ne 0) { throw "compilation failed" }
Write-Host "built $Out ($((Get-Item $Out).Length) bytes) with $cc"
