# 🍺 Homebrew Release Guide

This guide walks you through the process of releasing Termonaut to Homebrew.

## 📋 Prerequisites

1. **GitHub CLI** installed and authenticated:
   ```bash
   brew install gh
   gh auth login
   ```

2. **Homebrew** installed for testing:
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

3. **Project** ready for release with all tests passing

## 🚀 Release Process

### Step 1: Build Release Binaries

```bash
# Build cross-platform binaries and generate checksums
./scripts/build-release.sh v0.9.0
```

This creates:
- `dist/termonaut-v0.9.0-darwin-amd64.tar.gz` (macOS Intel)
- `dist/termonaut-v0.9.0-darwin-arm64.tar.gz` (macOS Apple Silicon)
- `dist/termonaut-v0.9.0-linux-amd64.tar.gz` (Linux x86_64)
- `dist/termonaut-v0.9.0-linux-arm64.tar.gz` (Linux ARM64)
- `dist/termonaut-v0.9.0-checksums.txt` (SHA256 checksums)

### Step 2: Create GitHub Release

```bash
# Create a draft GitHub release with all artifacts
./scripts/create-github-release.sh v0.9.0
```

This will:
- Create a draft release on GitHub
- Upload all binary archives
- Include detailed release notes
- Provide next steps

### Step 3: Test Homebrew Formula Locally

```bash
# Test the formula before submission
brew install --formula Formula/termonaut.rb

# Test the installation
termonaut --version
termonaut --help
termonaut stats

# Uninstall for clean testing
brew uninstall termonaut
```

### Step 4: Submit to Homebrew

There are two main approaches:

#### Option A: Submit to homebrew-core (Recommended)

1. **Fork the homebrew-core repository**:
   ```bash
   gh repo fork Homebrew/homebrew-core
   ```

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/homebrew-core.git
   cd homebrew-core
   ```

3. **Create the formula**:
   ```bash
   # Copy our formula to the correct location
   cp ../termonaut/Formula/termonaut.rb Formula/termonaut.rb
   ```

4. **Test the formula**:
   ```bash
   brew install --formula Formula/termonaut.rb
   brew test termonaut
   brew audit --strict termonaut
   ```

5. **Submit PR**:
   ```bash
   git checkout -b add-termonaut
   git add Formula/termonaut.rb
   git commit -m "termonaut 0.9.0 (new formula)"
   git push origin add-termonaut
   gh pr create --title "termonaut 0.9.0 (new formula)" --body "Adds termonaut, a gamified terminal productivity tracker"
   ```

#### Option B: Create Your Own Tap (Faster)

1. **Create a tap repository**:
   ```bash
   gh repo create homebrew-termonaut --public
   git clone https://github.com/YOUR_USERNAME/homebrew-termonaut.git
   cd homebrew-termonaut
   ```

2. **Add the formula**:
   ```bash
   cp ../termonaut/Formula/termonaut.rb termonaut.rb
   git add termonaut.rb
   git commit -m "Add termonaut formula"
   git push origin main
   ```

3. **Install from your tap**:
   ```bash
   brew tap YOUR_USERNAME/termonaut
   brew install termonaut
   ```

## 🧪 Testing Checklist

Before submitting, ensure:

- [ ] Formula installs successfully on macOS (Intel)
- [ ] Formula installs successfully on macOS (Apple Silicon)
- [ ] Formula installs successfully on Linux (if supported)
- [ ] All binaries work correctly after installation
- [ ] `brew test termonaut` passes
- [ ] `brew audit --strict termonaut` passes
- [ ] Version information is correct
- [ ] Dependencies are minimal and justified
- [ ] Post-install message is helpful

## 📝 Formula Requirements

The Homebrew formula must meet these requirements:

1. **Clear description** - Concise but informative
2. **Correct license** - Must match the project license
3. **Platform support** - macOS required, Linux optional
4. **Checksums** - SHA256 for all downloads
5. **Tests** - Basic functionality verification
6. **Minimal dependencies** - Only essential dependencies
7. **Proper installation** - Files in correct locations

## 🔄 Update Process

For future releases:

1. **Update checksums** in `Formula/termonaut.rb`
2. **Update version** number
3. **Update URLs** to point to new release
4. **Test the updated formula**
5. **Submit PR** with version bump

## 🎯 Release Checklist

- [ ] All tests pass locally
- [ ] Version tag created and pushed
- [ ] GitHub release published
- [ ] Homebrew formula tested locally
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Homebrew submission completed

## 🆘 Troubleshooting

### Formula Audit Errors

**Issue**: `brew audit` reports style violations
**Solution**: Follow Homebrew's style guide and run `brew style termonaut`

**Issue**: Checksums don't match
**Solution**: Rebuild binaries and update checksums in formula

**Issue**: Dependencies missing
**Solution**: Add required dependencies to the formula

### Installation Errors

**Issue**: Binary not found after installation
**Solution**: Check that binary is properly installed to `bin.install`

**Issue**: Permission denied
**Solution**: Ensure binary has execute permissions in the archive

### Submission Issues

**Issue**: PR rejected due to naming conflicts
**Solution**: Check if similar formula exists, consider alternative naming

**Issue**: Formula doesn't meet quality standards
**Solution**: Address all feedback and re-submit with improvements

## 📚 Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Python Cookbook](https://docs.brew.sh/Python-for-Formula-Authors)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Homebrew Maintainer Guide](https://docs.brew.sh/Maintainer-Guidelines)

---

## 🎉 Post-Release

Once your formula is accepted:

1. **Update README.md** with Homebrew installation instructions
2. **Announce the release** on social media and developer communities
3. **Monitor issues** and be prepared to fix any problems
4. **Plan next release** based on user feedback

Congratulations on releasing to Homebrew! 🎉 