#!/bin/bash

# Termonaut Installation Verification Script
# Helps users verify their installation and configuration

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 Termonaut Installation Verification${NC}"
echo "==========================================="
echo

# Check if termonaut is installed
echo -n "✓ Checking if termonaut is installed... "
if command -v termonaut >/dev/null 2>&1; then
    echo -e "${GREEN}FOUND${NC}"
    TERMONAUT_VERSION=$(termonaut --version 2>/dev/null | head -n1 || echo "unknown")
    echo "  Version: $TERMONAUT_VERSION"
else
    echo -e "${RED}NOT FOUND${NC}"
    echo
    echo -e "${YELLOW}❌ Termonaut is not installed or not in PATH${NC}"
    echo "Please install it first:"
    echo "  brew tap oiahoon/termonaut && brew install termonaut"
    echo "  OR"
    echo "  curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash"
    exit 1
fi

echo

# Check shell integration
echo -n "✓ Checking shell integration... "
SHELL_TYPE=$(basename "$SHELL")

case "$SHELL_TYPE" in
    "zsh")
        if grep -q "termonaut" ~/.zshrc 2>/dev/null; then
            echo -e "${GREEN}CONFIGURED${NC}"
            echo "  Shell: zsh (hooks found in ~/.zshrc)"
        else
            echo -e "${YELLOW}NOT CONFIGURED${NC}"
            echo "  Run: termonaut advanced shell install"
        fi
        ;;
    "bash")
        if grep -q "termonaut" ~/.bashrc 2>/dev/null || grep -q "termonaut" ~/.bash_profile 2>/dev/null; then
            echo -e "${GREEN}CONFIGURED${NC}"
            echo "  Shell: bash (hooks found)"
        else
            echo -e "${YELLOW}NOT CONFIGURED${NC}"
            echo "  Run: termonaut advanced shell install"
        fi
        ;;
    *)
        echo -e "${YELLOW}UNKNOWN SHELL${NC}"
        echo "  Shell: $SHELL_TYPE"
        echo "  Run: termonaut advanced shell install"
        ;;
esac

echo

# Check configuration directory
echo -n "✓ Checking configuration directory... "
CONFIG_DIR="$HOME/.termonaut"
if [[ -d "$CONFIG_DIR" ]]; then
    echo -e "${GREEN}EXISTS${NC}"
    echo "  Location: $CONFIG_DIR"
    
    # Check config file
    if [[ -f "$CONFIG_DIR/config.toml" ]]; then
        echo "  Config file: ✓ Found"
    else
        echo -e "  Config file: ${YELLOW}Not found (will be created)${NC}"
    fi
    
    # Check database
    if [[ -f "$CONFIG_DIR/termonaut.db" ]]; then
        DB_SIZE=$(du -h "$CONFIG_DIR/termonaut.db" | cut -f1)
        echo "  Database: ✓ Found ($DB_SIZE)"
    else
        echo -e "  Database: ${YELLOW}Not found (will be created)${NC}"
    fi
else
    echo -e "${YELLOW}NOT FOUND${NC}"
    echo "  Will be created on first run"
fi

echo

# Test basic functionality
echo -n "✓ Testing basic functionality... "
if termonaut --help >/dev/null 2>&1; then
    echo -e "${GREEN}WORKING${NC}"
else
    echo -e "${RED}ERROR${NC}"
    echo "  Failed to run termonaut --help"
fi

echo

# Check available commands
echo "✓ Available commands:"
termonaut --help 2>/dev/null | grep -E "^  [a-z]" | head -10 | while read -r line; do
    echo "  $line"
done

echo

# Test stats (if data exists)
echo -n "✓ Testing stats command... "
if STATS_OUTPUT=$(termonaut stats 2>/dev/null); then
    echo -e "${GREEN}WORKING${NC}"
    
    # Extract command count if available
    if echo "$STATS_OUTPUT" | grep -q "Total Commands"; then
        COMMAND_COUNT=$(echo "$STATS_OUTPUT" | grep "Total Commands" | grep -o '[0-9]\+' | head -1)
        if [[ $COMMAND_COUNT -gt 0 ]]; then
            echo "  Commands tracked: $COMMAND_COUNT"
        else
            echo -e "  ${YELLOW}No commands tracked yet${NC}"
        fi
    else
        echo -e "  ${YELLOW}No data yet (start using your terminal!)${NC}"
    fi
else
    echo -e "${YELLOW}NO DATA${NC}"
    echo "  This is normal for a fresh installation"
fi

echo

# Configuration recommendations
echo -e "${BLUE}📋 Next Steps:${NC}"
echo

if ! grep -q "termonaut" ~/.zshrc ~/.bashrc ~/.bash_profile 2>/dev/null; then
    echo -e "${YELLOW}1. Set up shell integration:${NC}"
    echo "   termonaut advanced shell install"
    echo "   source ~/.zshrc  # or ~/.bashrc"
    echo
fi

echo -e "${GREEN}2. Start tracking commands:${NC}"
echo "   Just use your terminal normally!"
echo "   Commands will be automatically tracked"
echo

echo -e "${GREEN}3. Check your progress:${NC}"
echo "   termonaut stats           # View statistics"
echo "   termonaut achievements    # Check badges"
echo "   termonaut tui            # Interactive dashboard"
echo

echo -e "${GREEN}4. Customize your experience:${NC}"
echo "   termonaut config show                    # View settings"
echo "   termonaut config set theme emoji         # Enable emojis"
echo "   termonaut easter-egg --test             # Test easter eggs"
echo

echo -e "${GREEN}5. Explore advanced features:${NC}"
echo "   termonaut advanced --help               # Power user tools"
echo "   termonaut analytics                     # Deep insights"
echo "   termonaut heatmap generate              # Activity heatmap"
echo

echo -e "${BLUE}🎮 Happy terminal tracking!${NC}"
echo

# Quick tips
echo -e "${BLUE}💡 Pro Tips:${NC}"
echo "   • Run 'termonaut stats' daily to track progress"
echo "   • Achievements unlock automatically as you use your terminal"
echo "   • Try 'termonaut tui' for a beautiful interactive experience"
echo "   • Enable privacy mode if working with sensitive data"
echo

echo "For more help, visit: https://github.com/oiahoon/termonaut" 