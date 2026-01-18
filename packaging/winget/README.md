# Winget Package for gmail-ro

This directory contains the Winget manifest templates for distributing gmail-ro on Windows via the [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs) repository.

## Manifest Structure

Winget requires three manifest files (schema version 1.10.0):

```
packaging/winget/
├── OpenCliCollective.gmail-ro.yaml              # Version manifest
├── OpenCliCollective.gmail-ro.installer.yaml    # Installer manifest
├── OpenCliCollective.gmail-ro.locale.en-US.yaml # Locale manifest
└── README.md                                    # This file
```

| File | Purpose |
|------|---------|
| `*.yaml` (version) | Package identifier and version |
| `*.installer.yaml` | Download URLs, checksums, installer type |
| `*.locale.en-US.yaml` | Description, tags, license info |

## How It Works

### Placeholder Strategy

The manifests contain placeholders that are replaced at release time:

| Placeholder | Location | Replaced With |
|-------------|----------|---------------|
| `0.0.0` | All manifests | Actual version (e.g., `1.0.5`) |
| 64 zeros | Installer manifest | SHA256 checksum from release |
| `v0.0.0` in URL | Installer manifest | Actual tag (e.g., `v1.0.5`) |

### Release Workflow

The GitHub Actions workflow:

1. Downloads `checksums.txt` from the GitHub release
2. Extracts the SHA256 hash for `gmail-ro_v{VERSION}_windows_amd64.zip`
3. Checks if package exists in winget-pkgs (new vs update)
4. For **new packages**: Generates manifests from templates and uses `wingetcreate submit`
5. For **updates**: Uses `wingetcreate update` command
6. Creates a PR to microsoft/winget-pkgs

### New vs Update Detection

```powershell
$response = Invoke-WebRequest -Uri "https://api.github.com/repos/microsoft/winget-pkgs/contents/manifests/o/OpenCliCollective/gmail-ro" -Method Head -SkipHttpErrorCheck
if ($response.StatusCode -eq 200) {
    # Package exists - use wingetcreate update
} else {
    # New package - use wingetcreate submit with generated manifests
}
```

## Architecture Support

Only x64 (amd64) is supported. The GoReleaser configuration explicitly excludes Windows arm64 builds:

```yaml
ignore:
  - goos: windows
    goarch: arm64
```

## Local Validation

To validate the manifests locally (requires Windows with winget):

1. **Create a test directory** with processed manifests:
   ```powershell
   $testDir = "winget-test"
   $testVersion = "0.0.1"
   $testHash = "0000000000000000000000000000000000000000000000000000000000000001"

   New-Item -ItemType Directory -Path $testDir -Force

   # Process each manifest
   foreach ($file in Get-ChildItem "packaging/winget/*.yaml") {
       $content = Get-Content $file -Raw
       $content = $content -replace "0\.0\.0", $testVersion
       $content = $content -replace "0{64}", $testHash
       Set-Content "$testDir/$($file.Name)" $content
   }
   ```

2. **Validate**:
   ```powershell
   winget validate --manifest $testDir
   ```

## Installation (End Users)

Once the package is published to winget-pkgs:

```powershell
winget install OpenCliCollective.gmail-ro
```

## Related

- [Winget Manifest Documentation](https://github.com/microsoft/winget-pkgs/tree/master/doc/manifest)
- [wingetcreate Tool](https://github.com/microsoft/winget-create)
- [Winget Package Submission](https://learn.microsoft.com/en-us/windows/package-manager/package/repository)
