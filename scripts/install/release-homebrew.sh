#!/bin/bash

# Complete Homebrew Release Script for Termonaut
# This script handles the complete release process for Homebrew

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
VERSION=${1:-"v0.9.0"}
BUILD_VERSION=$(echo "$VERSION" | sed 's/^v//')

echo -e "${BLUE}🚀 Termonaut Homebrew Release Pipeline${NC}"
echo "=========================================="
echo "Version: $VERSION"
echo

# Step 1: Build Release Binaries
print_status "Step 1: Building release binaries..."
if [[ ! -f "scripts/build-release.sh" ]]; then
    print_error "Build script not found: scripts/build-release.sh"
    exit 1
fi

./scripts/build-release.sh "$VERSION"
print_success "Release binaries built successfully"

# Step 2: Verify checksums file exists
CHECKSUMS_FILE="dist/termonaut-${VERSION}-checksums.txt"
if [[ ! -f "$CHECKSUMS_FILE" ]]; then
    print_error "Checksums file not found: $CHECKSUMS_FILE"
    exit 1
fi

print_success "Checksums file verified"

# Step 3: Update Homebrew Formula
print_status "Step 2: Updating Homebrew formula..."
./scripts/update-homebrew-formula.sh "$VERSION"
print_success "Homebrew formula updated"

# Step 4: Test the Formula (if Homebrew is available)
if command -v brew >/dev/null 2>&1; then
    print_status "Step 3: Testing Homebrew formula..."

    # Uninstall any existing version
    if brew list termonaut >/dev/null 2>&1; then
        print_warning "Uninstalling existing termonaut..."
        brew uninstall termonaut
    fi

    # Test installation
    print_status "Testing formula installation..."
    if brew install --formula Formula/termonaut.rb; then
        print_success "Formula installation test passed"

        # Test basic functionality
        print_status "Testing basic functionality..."
        if termonaut --version >/dev/null 2>&1; then
            print_success "Basic functionality test passed"
        else
            print_error "Basic functionality test failed"
            exit 1
        fi

        # Test short alias
        if tn --version >/dev/null 2>&1; then
            print_success "Short alias test passed"
        else
            print_warning "Short alias test failed (non-critical)"
        fi

        # Run brew test
        print_status "Running brew test..."
        if brew test termonaut; then
            print_success "Brew test passed"
        else
            print_error "Brew test failed"
            exit 1
        fi

        # Run brew audit
        print_status "Running brew audit..."
        if brew audit --strict termonaut; then
            print_success "Brew audit passed"
        else
            print_warning "Brew audit failed (check warnings above)"
        fi

        # Clean up
        print_status "Cleaning up test installation..."
        brew uninstall termonaut

    else
        print_error "Formula installation test failed"
        exit 1
    fi
else
    print_warning "Homebrew not found, skipping formula tests"
fi

# Step 5: Show next steps
echo
print_success "🎉 Homebrew release preparation complete!"
echo
echo -e "${BLUE}📋 Next Steps:${NC}"
echo "1. Commit the updated formula:"
echo "   git add Formula/termonaut.rb"
echo "   git commit -m \"Update Homebrew formula to ${VERSION}\""
echo
echo "2. Create GitHub Release (if not done already):"
echo "   ./scripts/create-github-release.sh ${VERSION}"
echo
echo "3. Option A - Create your own tap:"
echo "   gh repo create homebrew-termonaut --public"
echo "   git clone https://github.com/YOUR_USERNAME/homebrew-termonaut.git"
echo "   cp Formula/termonaut.rb homebrew-termonaut/termonaut.rb"
echo "   cd homebrew-termonaut && git add . && git commit -m \"Add termonaut ${VERSION}\" && git push"
echo
echo "4. Option B - Submit to homebrew-core:"
echo "   # Fork homebrew-core on GitHub"
echo "   git clone https://github.com/YOUR_USERNAME/homebrew-core.git"
echo "   cp Formula/termonaut.rb homebrew-core/Formula/termonaut.rb"
echo "   cd homebrew-core"
echo "   git checkout -b add-termonaut-${BUILD_VERSION}"
echo "   git add Formula/termonaut.rb"
echo "   git commit -m \"termonaut ${BUILD_VERSION} (new formula)\""
echo "   git push origin add-termonaut-${BUILD_VERSION}"
echo "   # Create PR on GitHub"
echo
echo "5. Test your tap (if using Option A):"
echo "   brew tap YOUR_USERNAME/termonaut"
echo "   brew install termonaut"
echo
echo -e "${GREEN}🍺 Ready for Homebrew release!${NC}"

# Step 6: Show release artifacts
echo
print_status "Release artifacts created:"
ls -la dist/termonaut-${VERSION}-*
echo
print_status "Checksums:"
cat "$CHECKSUMS_FILE"