#!/bin/bash

# Enhanced Termonaut GitHub Integration Demo
# Demonstrates all GitHub integration features including sync, badges, and workflows

echo "🚀 Enhanced Termonaut GitHub Integration Demo"
echo "=============================================="
echo

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if termonaut is available
if ! command -v ./termonaut &> /dev/null; then
    echo -e "${RED}❌ Termonaut binary not found. Please build it first:${NC}"
    echo "   go build -o termonaut cmd/termonaut/*.go"
    exit 1
fi

echo -e "${BLUE}📊 1. Current Stats Overview${NC}"
echo "──────────────────────────────"
./termonaut stats
echo

echo -e "${BLUE}📋 2. GitHub Sync Configuration${NC}"
echo "─────────────────────────────────"
echo "Current sync configuration:"
echo "• Sync enabled: $(./termonaut config get sync_enabled)"
echo "• Sync repository: $(./termonaut config get sync_repo)"
echo "• Update frequency: $(./termonaut config get badge_update_frequency)"
echo

echo -e "${BLUE}🎯 3. Badge Generation${NC}"
echo "─────────────────────────"
echo "Generating dynamic badges for GitHub profile..."
echo

# Generate badges in different formats
echo -e "${YELLOW}📝 URL Format:${NC}"
./termonaut github badges generate --format url
echo

echo -e "${YELLOW}📝 JSON Format:${NC}"
./termonaut github badges generate --format json
echo

echo -e "${YELLOW}📝 Markdown Format:${NC}"
./termonaut github badges generate --format markdown
echo

echo -e "${BLUE}📄 4. Profile Generation${NC}"
echo "────────────────────────────"
echo "Generating comprehensive GitHub profile..."
echo

# Generate profile
./termonaut github profile generate > /tmp/termonaut_profile.md
echo -e "${GREEN}✅ Profile generated! Preview:${NC}"
echo
head -20 /tmp/termonaut_profile.md
echo "..."
echo

echo -e "${BLUE}🤖 5. GitHub Actions Workflows${NC}"
echo "─────────────────────────────────"
echo "Available workflow templates:"
echo

./termonaut github actions list
echo

echo -e "${BLUE}🔧 6. Workflow Generation${NC}"
echo "────────────────────────────"
echo "Generating GitHub Actions workflow files..."
echo

# Create demo directory for workflows
mkdir -p github-demo/.github/workflows

echo -e "${YELLOW}Generating stats update workflow...${NC}"
./termonaut github actions generate termonaut-stats-update
echo

echo -e "${YELLOW}Generating profile sync workflow...${NC}"
./termonaut github actions generate termonaut-profile-sync
echo

echo -e "${YELLOW}Generating weekly report workflow...${NC}"
./termonaut github actions generate termonaut-weekly-report
echo

echo -e "${BLUE}💾 7. Export Examples${NC}"
echo "───────────────────────"
echo "Creating export examples..."

# Create exports directory
mkdir -p examples/exports

# Export badges
echo "Exporting badges..."
./termonaut github badges generate --format json --output examples/exports/badges.json
./termonaut github badges generate --format markdown --output examples/exports/badges.md

# Export profile
echo "Exporting profile..."
./termonaut github profile generate --format markdown --output examples/exports/profile.md
./termonaut github profile generate --format json --output examples/exports/profile.json

echo -e "${GREEN}✅ Exports saved to examples/exports/${NC}"
ls -la examples/exports/
echo

echo -e "${BLUE}🔄 8. Sync Feature Demo${NC}"
echo "─────────────────────────"
echo "Testing GitHub synchronization..."
echo

# Test sync status
echo -e "${YELLOW}Sync Status:${NC}"
./termonaut advanced github sync status
echo

# Test sync now (placeholder)
echo -e "${YELLOW}Manual Sync:${NC}"
./termonaut advanced github sync now
echo

echo -e "${BLUE}📋 9. Setup Instructions${NC}"
echo "────────────────────────────"
echo -e "${GREEN}🎉 GitHub Integration Setup Complete!${NC}"
echo
echo "📋 Next Steps for Full Integration:"
echo
echo "1️⃣  Repository Setup:"
echo "   • Create a GitHub repository for your profile (e.g., username/username)"
echo "   • Or use a dedicated stats repository (e.g., username/termonaut-stats)"
echo
echo "2️⃣  Badge Integration:"
echo "   • Copy badge URLs from the generated output"
echo "   • Add them to your GitHub profile README.md"
echo "   • Example:"
echo "     ![Commands](https://img.shields.io/badge/Commands-116-green?style=flat-square&logo=terminal&logoColor=white)"
echo
echo "3️⃣  GitHub Actions Automation:"
echo "   • Copy the generated workflow files to your repository"
echo "   • Commit them to .github/workflows/ directory"
echo "   • Configure repository secrets if needed"
echo
echo "4️⃣  Profile Enhancement:"
echo "   • Use the generated profile markdown in your README"
echo "   • Customize the content to match your style"
echo "   • Set up automatic updates via GitHub Actions"
echo
echo "5️⃣  Social Sharing:"
echo "   • Share your terminal productivity stats on social media"
echo "   • Use the generated snippets for different platforms"
echo "   • Show off your command-line skills!"
echo

echo -e "${BLUE}🔗 10. Useful Links${NC}"
echo "──────────────────────"
echo "• Shields.io Documentation: https://shields.io/"
echo "• GitHub Actions Documentation: https://docs.github.com/en/actions"
echo "• GitHub Profile README Guide: https://docs.github.com/en/github/setting-up-and-managing-your-github-profile/managing-your-profile-readme"
echo "• Termonaut Repository: https://github.com/oiahoon/termonaut"
echo

echo -e "${GREEN}🎯 Demo Complete!${NC}"
echo "Check the generated files in:"
echo "• examples/exports/ - Badge and profile exports"
echo "• .github/workflows/ - GitHub Actions workflows"
echo
echo "Happy coding! 🚀"