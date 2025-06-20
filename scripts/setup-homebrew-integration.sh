#!/bin/bash

# Termonaut Homebrew Integration Setup Script
# This script helps setup automatic Homebrew formula updates

set -e

echo "🍺 Termonaut Homebrew Integration Setup"
echo "======================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -f "cmd/termonaut/main.go" ]; then
    echo -e "${RED}❌ Error: This script must be run from the termonaut project root${NC}"
    exit 1
fi

echo -e "${BLUE}📋 Current Setup Status:${NC}"
echo ""

# Check if homebrew-termonaut repo exists
HOMEBREW_REPO="oiahoon/homebrew-termonaut"
echo -e "🔍 Checking Homebrew tap repository: ${HOMEBREW_REPO}"

if curl -s "https://api.github.com/repos/${HOMEBREW_REPO}" | grep -q '"name"'; then
    echo -e "${GREEN}✅ Homebrew tap repository exists${NC}"
    REPO_EXISTS=true
else
    echo -e "${YELLOW}⚠️  Homebrew tap repository not found${NC}"
    REPO_EXISTS=false
fi

echo ""
echo -e "${BLUE}🛠️  Integration Options:${NC}"
echo ""

if [ "$REPO_EXISTS" = true ]; then
    echo -e "${GREEN}Option 1: Automatic Integration (Recommended)${NC}"
    echo "  ✅ Your homebrew-termonaut repo will be automatically updated"
    echo "  ✅ Formula updates happen with each release"
    echo "  ✅ Zero manual work required"
    echo ""
    echo -e "${BLUE}Setup Steps:${NC}"
    echo "1. The GitHub Actions workflow will automatically:"
    echo "   - Update your homebrew-termonaut/termonaut.rb"
    echo "   - Update local Formula/termonaut.rb as backup"
    echo "   - Commit and push changes"
    echo ""
    echo "2. Users install with:"
    echo -e "${GREEN}   brew install oiahoon/termonaut/termonaut${NC}"
    echo ""

    read -p "🤔 Do you want to test the integration now? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        echo -e "${BLUE}🧪 Testing Homebrew Integration...${NC}"
        echo ""

        # Get current version
        CURRENT_VERSION=$(grep 'version = ' cmd/termonaut/main.go | sed 's/.*version = "\([^"]*\)".*/\1/')
        echo "Current version: ${CURRENT_VERSION}"

        echo ""
        echo "To test the integration, you can:"
        echo "1. Go to GitHub Actions: https://github.com/oiahoon/termonaut/actions"
        echo "2. Run 'Update Homebrew Formula' workflow manually"
        echo "3. Use current version: v${CURRENT_VERSION}"
        echo ""
        echo -e "${GREEN}✅ Integration is ready to use!${NC}"
    fi

else
    echo -e "${YELLOW}Option 1: Create Homebrew Tap Repository${NC}"
    echo "  📝 Create: https://github.com/oiahoon/homebrew-termonaut"
    echo "  📁 Add: termonaut.rb formula file"
    echo "  🔄 Enable automatic updates"
    echo ""
    echo -e "${BLUE}Option 2: Local Formula Only${NC}"
    echo "  📁 Use: Formula/termonaut.rb in this repo"
    echo "  👤 Users install with: brew install Formula/termonaut.rb"
    echo "  ⚠️  Less convenient for users"
    echo ""

    read -p "🤔 Do you want to create the homebrew-termonaut repository? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        echo -e "${BLUE}📝 Creating Homebrew Tap Repository...${NC}"
        echo ""
        echo "Please follow these steps:"
        echo ""
        echo "1. Go to: https://github.com/new"
        echo "2. Repository name: homebrew-termonaut"
        echo "3. Description: Homebrew tap for Termonaut - gamified terminal productivity tracker"
        echo "4. Make it public"
        echo "5. Initialize with README"
        echo "6. Create repository"
        echo ""
        echo "7. After creation, the GitHub Actions will automatically:"
        echo "   - Create and update termonaut.rb formula"
        echo "   - Handle all releases automatically"
        echo ""
        echo -e "${GREEN}✅ Then re-run this script to complete setup!${NC}"
    else
        echo ""
        echo -e "${BLUE}📁 Using Local Formula Only${NC}"
        echo ""
        echo "Your Formula/termonaut.rb will be updated automatically."
        echo "Users can install with:"
        echo -e "${GREEN}  brew install Formula/termonaut.rb${NC}"
        echo ""
        echo -e "${GREEN}✅ Local formula setup is ready!${NC}"
    fi
fi

echo ""
echo -e "${BLUE}🔧 Advanced Configuration (Optional):${NC}"
echo ""
echo "For enhanced security, you can create a Personal Access Token:"
echo "1. Go to: https://github.com/settings/tokens"
echo "2. Generate new token (classic)"
echo "3. Select scopes: repo, workflow"
echo "4. Add as repository secret: HOMEBREW_TAP_TOKEN"
echo ""
echo "This allows more reliable access to your homebrew-termonaut repo."
echo ""

echo -e "${GREEN}🎉 Homebrew Integration Setup Complete!${NC}"
echo ""
echo -e "${BLUE}📖 What happens next:${NC}"
echo "• Every time you release a new version, the Homebrew formula updates automatically"
echo "• Users get the latest version with: brew upgrade termonaut"
echo "• Zero manual maintenance required!"
echo ""
echo -e "${BLUE}🚀 Ready to release? Try:${NC}"
echo "  GitHub Actions → Manual Release → Enter version → Done!"
echo ""
echo -e "${GREEN}Happy brewing! 🍺${NC}"