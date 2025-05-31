# Homebrew Distribution Setup Guide

## 📋 Overview

This guide explains how to set up Homebrew distribution for todotui using goreleaser and homebrew-tap.

## 🚀 Prerequisites

1. **GitHub Repository**: Main repository (`yuucu/todotui`)
2. **Homebrew Tap Repository**: Create `yuucu/homebrew-tap` repository
3. **GitHub Token**: Personal access token with repo permissions

## 📝 Steps to Enable Homebrew Distribution

### 1. Create Homebrew Tap Repository

```bash
# Create a new repository named 'homebrew-tap' on GitHub
# Repository: https://github.com/yuucu/homebrew-tap
```

### 2. Setup GitHub Secrets

Add the following secret to your main repository (`yuucu/todotui`):

- **Secret Name**: `HOMEBREW_TAP_GITHUB_TOKEN`
- **Secret Value**: GitHub Personal Access Token with `repo` scope

### 3. Enable Homebrew in GoReleaser

Uncomment the homebrew configuration in `.goreleaser.yml` and update the release workflow:

```yaml
# In .github/workflows/release.yml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

### 4. Test Release Process

1. Create a test tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. Check GitHub Actions for successful release

3. Verify that Formula is created in `homebrew-tap` repository

### 5. Installation Commands

Once set up, users can install via:

```bash
# Add tap
brew tap yuucu/tap

# Install
brew install todotui

# Or in one command
brew install yuucu/tap/todotui
```

## 🧪 Testing Locally

Test goreleaser configuration:

```bash
# Install goreleaser
brew install goreleaser

# Test without releasing
goreleaser release --snapshot --clean

# Test homebrew formula (after release)
brew tap yuucu/tap
brew install --build-from-source todotui
```

## 📚 References

- [GoReleaser Homebrew Documentation](https://goreleaser.com/customization/homebrew/)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Creating Custom Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap) 