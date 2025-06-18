#!/bin/bash

# Termonaut Release Build Script
# Builds cross-platform binaries for Homebrew distribution

set -e

VERSION=${1:-"v0.9.0"}
BUILD_DIR="dist"
BINARY_NAME="termonaut"

echo "🚀 Building Termonaut ${VERSION} for release..."

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build metadata
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Build flags
LDFLAGS="-s -w -X github.com/oiahoon/termonaut/cmd/termonaut.Version=${VERSION} -X github.com/oiahoon/termonaut/cmd/termonaut.BuildTime=${BUILD_TIME} -X github.com/oiahoon/termonaut/cmd/termonaut.GitCommit=${GIT_COMMIT} -X github.com/oiahoon/termonaut/cmd/termonaut.GitBranch=${GIT_BRANCH}"

echo "📦 Building binaries..."

# macOS (Intel)
echo "  🍎 Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ./cmd/termonaut

# macOS (Apple Silicon)
echo "  🍎 Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ./cmd/termonaut

# Linux (x86_64)
echo "  🐧 Building for Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ./cmd/termonaut

# Linux (ARM64)
echo "  🐧 Building for Linux (ARM64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ./cmd/termonaut

# Create archives
echo "📦 Creating release archives..."

cd ${BUILD_DIR}

# macOS Intel
tar -czf ${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
# macOS Apple Silicon
tar -czf ${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
# Linux x86_64
tar -czf ${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
# Linux ARM64
tar -czf ${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64

# Generate checksums
echo "🔐 Generating checksums..."
if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 *.tar.gz > ${BINARY_NAME}-${VERSION}-checksums.txt
elif command -v sha256sum >/dev/null 2>&1; then
    sha256sum *.tar.gz > ${BINARY_NAME}-${VERSION}-checksums.txt
else
    echo "Warning: Neither shasum nor sha256sum found. Checksums not generated."
fi

cd ..

echo "✅ Release build complete!"
echo ""
echo "📂 Built files:"
ls -la ${BUILD_DIR}/*.tar.gz
echo ""
echo "🔐 Checksums:"
cat ${BUILD_DIR}/${BINARY_NAME}-${VERSION}-checksums.txt
echo ""
echo "📋 Next steps:"
echo "1. Upload the .tar.gz files to GitHub Release"
echo "2. Use the checksums in the Homebrew formula"
echo "3. Update the Homebrew formula with the new version"