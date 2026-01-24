# Chocolatey Package for gmail-ro

This directory contains the Chocolatey package definition for distributing gmail-ro on Windows via the [Chocolatey community repository](https://community.chocolatey.org/).

## Package Structure

```
packaging/chocolatey/
├── gmail-ro.nuspec              # Package metadata
├── tools/
│   ├── chocolateyInstall.ps1    # Installation script
│   └── chocolateyUninstall.ps1  # Uninstallation script
└── README.md                    # This file
```

## How It Works

### Checksum Injection

The install script contains a placeholder (`CHECKSUM_AMD64_PLACEHOLDER`) that is replaced at release time by the GitHub Actions workflow. The workflow:

1. Downloads `checksums.txt` from the GitHub release
2. Extracts the SHA256 hash for `gmail-ro_v{VERSION}_windows_amd64.zip`
3. Replaces the placeholder in `chocolateyInstall.ps1`
4. Packs and pushes to Chocolatey

### Version Updates

The `<version>0.0.0</version>` in the nuspec is also replaced at release time with the actual version being released.

## Local Testing

To test the package locally:

1. **Set a test version** in `gmail-ro.nuspec`:
   ```xml
   <version>1.0.0</version>
   ```

2. **Set a test checksum** in `tools/chocolateyInstall.ps1`:
   ```powershell
   $checksumAmd64 = 'actual_sha256_from_release'
   ```

3. **Pack the package**:
   ```powershell
   cd packaging/chocolatey
   choco pack
   ```

4. **Install locally** (requires admin):
   ```powershell
   choco install gmail-readonly -s . -y
   ```

5. **Test the installation**:
   ```powershell
   gmro --version
   gmro --help
   ```

6. **Uninstall**:
   ```powershell
   choco uninstall gmail-readonly -y
   ```

## Chocolatey Moderation Compliance

This package follows Chocolatey's automated moderation rules:

| Rule | Requirement | Implementation |
|------|-------------|----------------|
| CPMR0041 | projectUrl must differ from projectSourceUrl | projectUrl has `#readme` suffix |
| CPMR0055 | No custom downloaders | Uses `Install-ChocolateyZipPackage` only |
| CPMR0073 | SHA256 checksums required | Checksum injected at release time |

## Architecture Support

Only 64-bit Windows (x64/amd64) is supported. The GoReleaser configuration explicitly excludes Windows arm64 builds.

## Installation (End Users)

Once published to Chocolatey:

```powershell
choco install gmail-readonly
```

## Related

- [Chocolatey Package Documentation](https://docs.chocolatey.org/en-us/create/)
- [Chocolatey Moderation Rules](https://docs.chocolatey.org/en-us/community-repository/moderation/)
