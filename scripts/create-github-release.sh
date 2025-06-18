#!/bin/bash

# Create GitHub Release Script for Termonaut
# This script creates a GitHub release and uploads the artifacts

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
VERSION="v0.9.0"
BUILD_VERSION="0.9.0"  # Version used in file names
RELEASE_TITLE="Termonaut v0.9.0 - Official Stable Release with Complete GitHub Integration & Enhanced UX"
REPO_OWNER="oiahoon"
REPO_NAME="termonaut"

echo -e "${BLUE}🚀 Creating GitHub Release: ${VERSION}${NC}"
echo -e "════════════════════════════════════════════════"
echo

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI (gh) is not installed${NC}"
    echo -e "Install it with: brew install gh"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${RED}❌ Not authenticated with GitHub CLI${NC}"
    echo -e "Run: gh auth login"
    exit 1
fi

# Check if release files exist
if [[ ! -d "dist" ]] || [[ ! -f "dist/RELEASE_NOTES_${BUILD_VERSION}.md" ]]; then
    echo -e "${RED}❌ Release artifacts not found${NC}"
    echo -e "Run the release build script first: ./scripts/release-${BUILD_VERSION}.sh"
    exit 1
fi

# Create the release
echo -e "${BLUE}📝 Creating GitHub release...${NC}"

# Check if release already exists
if gh release view "${VERSION}" --repo "${REPO_OWNER}/${REPO_NAME}" &> /dev/null; then
    echo -e "${YELLOW}⚠️  Release ${VERSION} already exists${NC}"
    echo -e "Deleting existing release to recreate..."
    gh release delete "${VERSION}" --repo "${REPO_OWNER}/${REPO_NAME}" --yes
fi

# Create the release with notes from file
gh release create "${VERSION}" \
    --repo "${REPO_OWNER}/${REPO_NAME}" \
    --title "${RELEASE_TITLE}" \
    --notes-file "dist/RELEASE_NOTES_${BUILD_VERSION}.md" \
    --generate-notes

echo -e "${GREEN}✅ Release created successfully${NC}"

# Upload all release assets
echo -e "${BLUE}📦 Uploading release assets...${NC}"

# Upload binaries
for file in dist/termonaut-${BUILD_VERSION}-*; do
    if [[ -f "$file" && ! "$file" =~ \.md$ && ! "$file" =~ checksums ]]; then
        echo -e "Uploading $(basename "$file")..."
        gh release upload "${VERSION}" "$file" --repo "${REPO_OWNER}/${REPO_NAME}"
    fi
done

# Upload checksums
echo -e "Uploading checksums..."
gh release upload "${VERSION}" "dist/termonaut-${BUILD_VERSION}-checksums.txt" --repo "${REPO_OWNER}/${REPO_NAME}"

echo -e "${GREEN}✅ All assets uploaded successfully${NC}"

# Display release information
echo
echo -e "${GREEN}🎉 Release ${VERSION} created successfully!${NC}"
echo
echo -e "${BLUE}📋 Release Details:${NC}"
echo -e "• URL: https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/${VERSION}"
echo -e "• Title: ${RELEASE_TITLE}"
echo -e "• Type: Stable Release"
echo
echo -e "${YELLOW}📋 Next Steps:${NC}"
echo -e "1. Test the release thoroughly"
echo -e "2. Announce stable release to community for feedback"
echo -e "3. Collect feedback and iterate"
echo -e "4. Prepare for v1.0.0 stable release"
echo
echo -e "${BLUE}🔗 Release URL: https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/${VERSION}${NC}"