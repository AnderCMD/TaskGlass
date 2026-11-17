<#
.SYNOPSIS
  Re-renders the logo, banner, social preview and icons from assets/src.

.DESCRIPTION
  Uses Microsoft Edge (headless) to rasterise the SVG/HTML sources, then
  tools/genicon to pack the .ico files. Outputs are committed, so this only
  needs to run when something in assets/src changes.
#>
$ErrorActionPreference = "Stop"
$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$src = Join-Path $root "assets\src"
$img = Join-Path $root "docs\images"
$tmp = Join-Path ([IO.Path]::GetTempPath()) "taskglass-render"
New-Item -ItemType Directory -Force $img, $tmp | Out-Null

$edge = @(
    "${env:ProgramFiles(x86)}\Microsoft\Edge\Application\msedge.exe",
    "$env:ProgramFiles\Microsoft\Edge\Application\msedge.exe",
    "$env:ProgramFiles\Google\Chrome\Application\chrome.exe"
) | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $edge) { throw "Microsoft Edge or Chrome is required to render the assets." }

function Render($file, $out, $w, $h, [switch]$Transparent) {
    $args = @("--headless=new", "--disable-gpu", "--hide-scrollbars", "--force-device-scale-factor=1",
        "--user-data-dir=$tmp\profile-$([guid]::NewGuid().ToString('N'))", "--window-size=$w,$h", "--screenshot=$out")
    Remove-Item $out -ErrorAction SilentlyContinue
    if ($Transparent) { $args += "--default-background-color=00000000" }
    $args += ([Uri](Resolve-Path $file).Path).AbsoluteUri
    # Edge prints status lines to stderr, which -ErrorAction Stop would treat as failure.
    $prev = $ErrorActionPreference; $ErrorActionPreference = "Continue"
    & $edge @args *> $null
    $ErrorActionPreference = $prev
    for ($i = 0; $i -lt 40 -and -not (Test-Path $out); $i++) { Start-Sleep -Milliseconds 250 }
    if (-not (Test-Path $out)) { throw "render failed: $out" }
}

# Wrap a bare SVG in a page that fills the viewport, so it can be screenshotted at any size.
function RenderSvg($svg, $out, $size) {
    $page = Join-Path $tmp ("svg-" + [IO.Path]::GetFileNameWithoutExtension($svg) + "-$size.html")
    $uri = ([Uri](Resolve-Path $svg).Path).AbsoluteUri
    "<!doctype html><meta charset=utf-8><style>html,body{margin:0;background:transparent;overflow:hidden}img{display:block;width:${size}px;height:${size}px}</style><img src='$uri'>" |
        Set-Content $page -Encoding utf8
    Render $page $out $size $size -Transparent
}

Render "$src\banner.html" "$img\banner.png" 1280 320
Render "$src\social-preview.html" "$img\social-preview.png" 1280 640
RenderSvg "$src\logo.svg" "$img\logo.png" 512

foreach ($s in 16, 20, 24, 32, 40, 48, 64, 128, 256) {
    RenderSvg "$src\logo.svg" "$tmp\color-$s.png" $s
    RenderSvg "$src\glyph-light.svg" "$tmp\light-$s.png" $s
    RenderSvg "$src\glyph-dark.svg" "$tmp\dark-$s.png" $s
}
Push-Location $root
try { go run ./tools/genicon -in $tmp -out (Join-Path $root "assets") } finally { Pop-Location }
Write-Host "assets rendered"
