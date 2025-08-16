# Release Guide

This document describes the release process for phpier.

## Release Process Overview

1. **Prepare Release**
   - Update version numbers
   - Update CHANGELOG.md
   - Test functionality
   - Update documentation

2. **Build and Test**
   - Build for all platforms
   - Test binaries
   - Generate checksums

3. **Create Release**
   - Tag repository
   - Push to GitHub
   - Create GitHub release
   - Publish binaries

## Prerequisites

### Tools Required
- Go 1.20+
- Git
- GitHub CLI (`gh`) - for GitHub releases
- Make (optional, for convenience)

### Permissions
- Write access to the repository
- Permission to create releases
- GitHub token with repo scope

## Release Types

### Patch Release (v1.0.1)
- Bug fixes
- Security patches
- Documentation updates
- No breaking changes

### Minor Release (v1.1.0)
- New features
- Enhancements
- Non-breaking changes
- Deprecations

### Major Release (v2.0.0)
- Breaking changes
- Major architecture changes
- API changes

## Step-by-Step Release Process

### 1. Prepare the Release

#### Update Version Information
```bash
# Check current version
./phpier version

# Choose new version (follow semantic versioning)
NEW_VERSION="v1.0.0"
```

#### Update CHANGELOG.md
```markdown
## [1.0.0] - 2024-01-15

### Added
- New feature descriptions

### Changed
- Changes to existing features

### Fixed
- Bug fix descriptions
```

#### Update Documentation
- Update README.md if needed
- Update any version references
- Update installation instructions

### 2. Test the Release

#### Run Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint
```

#### Test Build Process
```bash
# Test local build
make quick

# Test cross-platform build (dry run)
make release-dry VERSION=$NEW_VERSION
```

#### Manual Testing
```bash
# Test init command
mkdir test-release
cd test-release
../phpier init 8.3

# Test up command (if Docker is available)
../phpier up -d
../phpier down
```

### 3. Create the Release

#### Option A: Using Scripts (Recommended)

##### Local Build and Test
```bash
# Build for all platforms (dry run)
./scripts/release.sh $NEW_VERSION --dry-run --checksums --zip

# If everything looks good, build for real
./scripts/release.sh $NEW_VERSION --checksums --zip
```

##### GitHub Release
```bash
# Create GitHub release with binaries
./scripts/release.sh $NEW_VERSION --checksums --zip --github
```

#### Option B: Using Makefile
```bash
# Build release
make release VERSION=$NEW_VERSION

# Create GitHub release
make github-release VERSION=$NEW_VERSION
```

#### Option C: Using Git Tags (Automated via GitHub Actions)
```bash
# Commit all changes
git add .
git commit -m "Release $NEW_VERSION"

# Create and push tag
git tag $NEW_VERSION
git push origin $NEW_VERSION

# GitHub Actions will automatically:
# - Run tests
# - Build for all platforms
# - Create GitHub release
```

### 4. Verify the Release

#### Check GitHub Release
1. Go to GitHub releases page
2. Verify all binaries are attached
3. Verify checksums.txt is present
4. Test download links

#### Test Downloaded Binaries
```bash
# Download and test a binary
curl -L -o phpier-test https://github.com/your-org/phpier/releases/download/$NEW_VERSION/phpier-linux-amd64
chmod +x phpier-test
./phpier-test version
./phpier-test --help
```

#### Verify Checksums
```bash
# Download checksums
curl -L -o checksums.txt https://github.com/your-org/phpier/releases/download/$NEW_VERSION/checksums.txt

# Verify checksum
sha256sum -c checksums.txt
```

### 5. Post-Release Tasks

#### Update Documentation
- Update installation instructions with new version
- Update any getting started guides
- Update Docker Hub descriptions if applicable

#### Announce Release
- Create announcement (social media, blog, etc.)
- Update project documentation sites
- Notify users through appropriate channels

#### Prepare for Next Release
- Create new development branch if needed
- Update version to next development version
- Plan next release features

## Automated Release (GitHub Actions)

The project includes GitHub Actions workflows for automated releases:

### CI Workflow (`.github/workflows/ci.yml`)
- Runs on every push and PR
- Tests multiple Go versions
- Runs linting and tests
- Builds binary for verification

### Release Workflow (`.github/workflows/release.yml`)
- Triggers on version tags (v*)
- Builds for all supported platforms
- Generates checksums and archives
- Creates GitHub release automatically

### Using Automated Release
```bash
# Simply create and push a tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will handle the rest
```

## Supported Platforms

phpier is built for the following platforms:

| OS      | Architecture | Binary Name               |
|---------|-------------|---------------------------|
| Linux   | AMD64       | phpier-linux-amd64     |
| Linux   | ARM64       | phpier-linux-arm64     |
| macOS   | AMD64       | phpier-darwin-amd64    |
| macOS   | ARM64       | phpier-darwin-arm64    |
| Windows | AMD64       | phpier-windows-amd64.exe |

## Build Configuration

### Build Flags
The release builds include the following information:
- **Version**: Git tag or specified version
- **Commit**: Git commit hash
- **Date**: Build timestamp
- **Go Version**: Go version used for build

### LDFLAGS
```bash
LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
```

## Troubleshooting

### Build Failures
```bash
# Clean and retry
make clean
go mod tidy
make build
```

### Missing Dependencies
```bash
# Reinstall dependencies
make deps
```

### Cross-Platform Build Issues
```bash
# Update Go version
go version  # Should be 1.20+

# Verify platform support
go tool dist list | grep -E "(linux|darwin|windows)"
```

### GitHub Release Issues
```bash
# Check GitHub CLI authentication
gh auth status

# Re-authenticate if needed
gh auth login
```

## Release Checklist

Use this checklist for each release:

### Pre-Release
- [ ] All tests passing
- [ ] CHANGELOG.md updated
- [ ] Version numbers updated
- [ ] Documentation updated
- [ ] Manual testing completed

### Release
- [ ] Tag created and pushed
- [ ] GitHub release created
- [ ] All binaries built successfully
- [ ] Checksums generated
- [ ] Release notes written

### Post-Release
- [ ] Binaries tested on different platforms
- [ ] Installation instructions verified
- [ ] Announcement published
- [ ] Next release planned

## Contact

For questions about the release process, please:
- Open an issue on GitHub
- Contact the maintainers
- Check the project documentation

---

**Remember**: Always test releases thoroughly before publishing to ensure quality and stability for users.