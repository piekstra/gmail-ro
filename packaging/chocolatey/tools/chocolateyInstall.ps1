$ErrorActionPreference = 'Stop'

$version = $env:ChocolateyPackageVersion
$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Checksum injected by release workflow - DO NOT EDIT MANUALLY
$checksumAmd64 = 'CHECKSUM_AMD64_PLACEHOLDER'

# Only 64-bit Windows is supported (arm64 builds not available)
if (-not [Environment]::Is64BitOperatingSystem) {
    throw "32-bit Windows is not supported. gmail-ro requires 64-bit Windows."
}

$baseUrl = "https://github.com/open-cli-collective/gmail-ro/releases/download/v${version}"
$zipFile = "gmail-ro_v${version}_windows_amd64.zip"
$url = "${baseUrl}/${zipFile}"

Write-Host "Installing gmail-ro ${version} for Windows x64..."
Write-Host "URL: ${url}"
Write-Host "Checksum (SHA256): ${checksumAmd64}"

Install-ChocolateyZipPackage -PackageName $env:ChocolateyPackageName `
    -Url $url `
    -UnzipLocation $toolsDir `
    -Checksum $checksumAmd64 `
    -ChecksumType 'sha256'

# Exclude non-executables from shimming
New-Item "$toolsDir\LICENSE.ignore" -Type File -Force | Out-Null
New-Item "$toolsDir\README.md.ignore" -Type File -Force | Out-Null

Write-Host "gmail-ro installed successfully. Run 'gmail-ro --help' to get started."
