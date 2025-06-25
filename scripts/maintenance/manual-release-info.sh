#!/bin/bash

# Manual GitHub Release Information Generator
# This script generates the information needed to manually create a GitHub release

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

VERSION="v0.9.0"
BUILD_VERSION="0.9.0"
RELEASE_TITLE="Termonaut v0.9.0 - Official Stable Release with Complete GitHub Integration & Enhanced UX"

echo -e "${BLUE}🚀 Manual GitHub Release Information${NC}"
echo -e "════════════════════════════════════════════════"
echo

echo -e "${GREEN}📋 Release Details:${NC}"
echo -e "• Tag: ${VERSION}"
echo -e "• Title: ${RELEASE_TITLE}"
echo -e "• Target: main branch"
echo -e "• Type: Latest release (not pre-release)"
echo

echo -e "${BLUE}📝 Release Notes File:${NC}"
echo -e "dist/RELEASE_NOTES_${BUILD_VERSION}.md"
echo

echo -e "${YELLOW}📦 Files to Upload:${NC}"
echo -e "The following files should be uploaded as release assets:"
echo

for file in dist/termonaut-${BUILD_VERSION}-*; do
    if [[ -f "$file" ]]; then
        size=$(ls -lh "$file" | awk '{print $5}')
        echo -e "• $(basename "$file") (${size})"
    fi
done

if [[ -f "dist/termonaut-${BUILD_VERSION}-checksums.txt" ]]; then
    echo -e "• termonaut-${BUILD_VERSION}-checksums.txt"
fi

echo

echo -e "${GREEN}🔗 Manual Steps:${NC}"
echo -e "1. Go to: https://github.com/oiahoon/termonaut/releases"
echo -e "2. Click 'Create a new release' or edit existing v0.9.0"
echo -e "3. Set tag: ${VERSION}"
echo -e "4. Set title: ${RELEASE_TITLE}"
echo -e "5. Copy content from: dist/RELEASE_NOTES_${BUILD_VERSION}.md"
echo -e "6. Upload all the files listed above"
echo -e "7. Uncheck 'Set as a pre-release' (this is a stable release)"
echo -e "8. Click 'Publish release'"

echo

echo -e "${BLUE}📊 Build Information:${NC}"
echo -e "• Build Date: $(date -u +'%Y-%m-%dT%H:%M:%SZ')"
echo -e "• Git Commit: $(git rev-parse --short HEAD)"
echo -e "• Branch: $(git rev-parse --abbrev-ref HEAD)"

echo

echo -e "${YELLOW}✅ Verification Steps:${NC}"
echo -e "After creating the release, verify:"
echo -e "1. Download links work for all platforms"
echo -e "2. Checksums match the uploaded files"
echo -e "3. Release notes display correctly"
echo -e "4. Installation instructions work"

echo

echo -e "${GREEN}🎉 Release ${VERSION} ready for manual creation!${NC}"