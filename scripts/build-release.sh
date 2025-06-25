#!/bin/bash

# 🚀 Termonaut Release Build Script
# Builds binaries for multiple platforms

set -e

VERSION="v0.10.1"
BUILD_DIR="build/release"
BINARY_NAME="termonaut"

echo "🚀 Building Termonaut $VERSION for multiple platforms"
echo "=================================================="

# Create build directory
mkdir -p "$BUILD_DIR"

# Build information
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD)

# Build flags
LDFLAGS="-X main.version=$VERSION -X main.commit=$GIT_COMMIT -X main.buildTime=$BUILD_TIME"

# Platform configurations
declare -a PLATFORMS=(
    "darwin/amd64"      # macOS Intel
    "darwin/arm64"      # macOS Apple Silicon
    "linux/amd64"       # Linux x64
    "linux/arm64"       # Linux ARM64
    "windows/amd64"     # Windows x64
)

echo "📦 Building for ${#PLATFORMS[@]} platforms..."
echo ""

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    output_name="$BINARY_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        output_name="$output_name.exe"
    fi
    
    echo "🔨 Building $GOOS/$GOARCH -> $output_name"
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags "$LDFLAGS" \
        -o "$BUILD_DIR/$output_name" \
        cmd/termonaut/*.go
    
    # Verify build
    if [ -f "$BUILD_DIR/$output_name" ]; then
        size=$(du -h "$BUILD_DIR/$output_name" | cut -f1)
        echo "   ✅ Success ($size)"
    else
        echo "   ❌ Failed"
        exit 1
    fi
done

echo ""
echo "📊 Build Summary"
echo "==============="
echo "Version: $VERSION"
echo "Commit: $GIT_COMMIT"
echo "Build Time: $BUILD_TIME"
echo "Output Directory: $BUILD_DIR"
echo ""

echo "📁 Generated Files:"
ls -la "$BUILD_DIR/"

echo ""
echo "🎯 Release Artifacts Ready!"
echo ""
echo "Next steps:"
echo "1. Test binaries on target platforms"
echo "2. Create GitHub release with these artifacts"
echo "3. Update installation scripts"
echo "4. Announce release"
