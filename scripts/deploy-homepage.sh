#!/bin/bash

# Termonaut Homepage Deployment Script
# This script helps deploy the homepage to GitHub Pages

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCS_DIR="$REPO_ROOT/docs"
BRANCH="main"

echo -e "${BLUE}🚀 Termonaut Homepage Deployment${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Function to check prerequisites
check_prerequisites() {
    echo -e "${BLUE}🔍 Checking prerequisites...${NC}"
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}❌ Error: Not in a git repository${NC}"
        exit 1
    fi
    
    # Check if docs directory exists
    if [ ! -d "$DOCS_DIR" ]; then
        echo -e "${RED}❌ Error: docs directory not found${NC}"
        exit 1
    fi
    
    # Check if index.html exists
    if [ ! -f "$DOCS_DIR/index.html" ]; then
        echo -e "${RED}❌ Error: index.html not found in docs directory${NC}"
        exit 1
    fi
    
    # Check git status
    if [ -n "$(git status --porcelain)" ]; then
        echo -e "${YELLOW}⚠️  Warning: You have uncommitted changes${NC}"
        echo -e "${YELLOW}   Consider committing them before deployment${NC}"
        echo ""
        git status --short
        echo ""
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}🛑 Deployment cancelled${NC}"
            exit 0
        fi
    fi
    
    echo -e "${GREEN}✅ Prerequisites check passed${NC}"
    echo ""
}

# Function to validate homepage
validate_homepage() {
    echo -e "${BLUE}🔍 Validating homepage...${NC}"
    
    # Check HTML structure
    local errors=0
    
    # Check for required elements
    if ! grep -q "<title>" "$DOCS_DIR/index.html"; then
        echo -e "${RED}❌ Missing <title> tag${NC}"
        errors=$((errors + 1))
    fi
    
    if ! grep -q "meta.*description" "$DOCS_DIR/index.html"; then
        echo -e "${RED}❌ Missing meta description${NC}"
        errors=$((errors + 1))
    fi
    
    if ! grep -q "og:title" "$DOCS_DIR/index.html"; then
        echo -e "${RED}❌ Missing Open Graph tags${NC}"
        errors=$((errors + 1))
    fi
    
    # Check for required assets
    if [ ! -f "$DOCS_DIR/assets/css/style.css" ]; then
        echo -e "${RED}❌ Missing CSS file${NC}"
        errors=$((errors + 1))
    fi
    
    if [ ! -f "$DOCS_DIR/assets/js/main.js" ]; then
        echo -e "${RED}❌ Missing JavaScript file${NC}"
        errors=$((errors + 1))
    fi
    
    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}✅ Homepage validation passed${NC}"
    else
        echo -e "${RED}❌ Homepage validation failed with $errors errors${NC}"
        exit 1
    fi
    echo ""
}

# Function to show deployment info
show_deployment_info() {
    echo -e "${BLUE}📊 Deployment Information${NC}"
    echo -e "${BLUE}========================${NC}"
    echo "Repository: $(git remote get-url origin 2>/dev/null || echo 'No remote origin')"
    echo "Branch: $BRANCH"
    echo "Current commit: $(git rev-parse --short HEAD)"
    echo "Docs directory: $DOCS_DIR"
    echo "GitHub Pages URL: https://$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\)\/\([^.]*\).*/\1.github.io\/\2/' 2>/dev/null || echo 'username.github.io/repository')"
    echo ""
}

# Function to commit and push changes
deploy_changes() {
    echo -e "${BLUE}🚀 Deploying homepage...${NC}"
    
    # Add docs directory to git
    git add "$DOCS_DIR"
    
    # Check if there are changes to commit
    if git diff --cached --quiet; then
        echo -e "${YELLOW}⚠️  No changes to deploy${NC}"
        echo -e "${GREEN}✅ Homepage is already up to date${NC}"
        return 0
    fi
    
    # Create commit message
    local commit_msg="🚀 Update Termonaut homepage"
    if [ -n "${1:-}" ]; then
        commit_msg="$1"
    fi
    
    # Commit changes
    echo -e "${BLUE}📝 Committing changes...${NC}"
    git commit -m "$commit_msg"
    
    # Push to remote
    echo -e "${BLUE}📤 Pushing to GitHub...${NC}"
    git push origin "$BRANCH"
    
    echo -e "${GREEN}✅ Homepage deployed successfully!${NC}"
    echo ""
}

# Function to wait for deployment
wait_for_deployment() {
    echo -e "${BLUE}⏳ Waiting for GitHub Pages deployment...${NC}"
    echo -e "${YELLOW}   This may take a few minutes${NC}"
    
    local github_pages_url
    github_pages_url="https://$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\)\/\([^.]*\).*/\1.github.io\/\2/' 2>/dev/null)"
    
    echo -e "${BLUE}🌐 Your homepage will be available at:${NC}"
    echo -e "${GREEN}   $github_pages_url${NC}"
    echo ""
    
    # Optional: Open browser
    read -p "Open homepage in browser? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        case "$(uname -s)" in
            Darwin)  # macOS
                open "$github_pages_url"
                ;;
            Linux)
                if command -v xdg-open > /dev/null; then
                    xdg-open "$github_pages_url"
                fi
                ;;
            CYGWIN*|MINGW32*|MSYS*|MINGW*)  # Windows
                start "$github_pages_url"
                ;;
        esac
    fi
}

# Function to show post-deployment tips
show_tips() {
    echo -e "${BLUE}💡 Post-Deployment Tips${NC}"
    echo -e "${BLUE}=======================${NC}"
    echo "• GitHub Pages deployment may take 5-10 minutes"
    echo "• Check GitHub Actions tab for deployment status"
    echo "• Test the live site on different devices"
    echo "• Monitor Lighthouse scores for performance"
    echo "• Update social media links if needed"
    echo ""
    echo -e "${BLUE}🔗 Useful Links:${NC}"
    echo "• GitHub Actions: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\)\/\([^.]*\).*/\1\/\2/' 2>/dev/null)/actions"
    echo "• Repository Settings: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\)\/\([^.]*\).*/\1\/\2/' 2>/dev/null)/settings/pages"
    echo ""
}

# Main function
main() {
    local commit_message="${1:-}"
    
    check_prerequisites
    validate_homepage
    show_deployment_info
    
    # Confirm deployment
    echo -e "${YELLOW}🚀 Ready to deploy homepage to GitHub Pages${NC}"
    read -p "Continue with deployment? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}🛑 Deployment cancelled${NC}"
        exit 0
    fi
    
    deploy_changes "$commit_message"
    wait_for_deployment
    show_tips
    
    echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Termonaut Homepage Deployment Script"
        echo ""
        echo "Usage: $0 [OPTIONS] [COMMIT_MESSAGE]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --validate     Only validate homepage and exit"
        echo "  --info         Show deployment info and exit"
        echo ""
        echo "Examples:"
        echo "  $0                                    # Deploy with default commit message"
        echo "  $0 'Update hero section'             # Deploy with custom commit message"
        echo "  $0 --validate                        # Only validate homepage"
        echo "  $0 --info                            # Show deployment information"
        exit 0
        ;;
    --validate)
        check_prerequisites
        validate_homepage
        echo -e "${GREEN}✅ Homepage is ready for deployment${NC}"
        exit 0
        ;;
    --info)
        show_deployment_info
        exit 0
        ;;
    --*)
        echo -e "${RED}❌ Error: Unknown option $1${NC}"
        echo "Use --help for usage information"
        exit 1
        ;;
esac

# Run main function
main "$@"
