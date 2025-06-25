#!/bin/bash

# Termonaut Safe Shell Installation Script
# Demonstrates the new safe configuration management system

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔒 Termonaut Safe Shell Installation Demo${NC}"
echo "=============================================="
echo

# Check if termonaut is available
if ! command -v ./termonaut &> /dev/null; then
    echo -e "${RED}❌ Termonaut binary not found. Please build it first:${NC}"
    echo "   go build -o termonaut cmd/termonaut/*.go"
    exit 1
fi

echo -e "${BLUE}📋 Safe Installation Features${NC}"
echo "──────────────────────────────"
echo "✅ Automatic backup creation before any changes"
echo "✅ Precise configuration block detection and removal"
echo "✅ Atomic file operations (write to temp, then rename)"
echo "✅ Shell syntax validation after modifications"
echo "✅ Automatic rollback on any failure"
echo "✅ Clean empty line management"
echo "✅ Support for corrupted/incomplete installations"
echo

# Detect current shell
SHELL_TYPE=$(basename "$SHELL")
echo -e "${BLUE}🐚 Detected Shell: ${SHELL_TYPE}${NC}"

case "$SHELL_TYPE" in
    "zsh")
        CONFIG_FILE="$HOME/.zshrc"
        ;;
    "bash")
        if [[ -f "$HOME/.bashrc" ]]; then
            CONFIG_FILE="$HOME/.bashrc"
        else
            CONFIG_FILE="$HOME/.bash_profile"
        fi
        ;;
    *)
        echo -e "${YELLOW}⚠️ Unsupported shell: $SHELL_TYPE${NC}"
        echo "This demo supports zsh and bash only"
        exit 1
        ;;
esac

echo "Configuration file: $CONFIG_FILE"
echo

# Show current config file status
echo -e "${BLUE}📄 Current Configuration Status${NC}"
echo "─────────────────────────────────"
if [[ -f "$CONFIG_FILE" ]]; then
    echo "✅ Configuration file exists"
    echo "   Size: $(wc -l < "$CONFIG_FILE") lines"

    if grep -q "termonaut" "$CONFIG_FILE" 2>/dev/null; then
        echo "⚠️  Existing Termonaut installation detected"

        # Show the existing block
        echo
        echo "Existing Termonaut block:"
        echo "========================"
        grep -A 20 -B 2 "# Termonaut shell integration" "$CONFIG_FILE" | head -25
        echo "========================"
    else
        echo "✅ No existing Termonaut installation"
    fi
else
    echo "⚠️  Configuration file does not exist (will be created)"
fi
echo

# Test 1: Clean installation
echo -e "${BLUE}🧪 Test 1: Clean Installation${NC}"
echo "─────────────────────────────"
echo "Installing Termonaut hook using safe method..."

if ./termonaut advanced shell install; then
    echo -e "${GREEN}✅ Installation successful${NC}"

    # Verify installation
    if grep -q "termonaut" "$CONFIG_FILE"; then
        echo "✅ Hook properly added to configuration"

        # Check for backup
        BACKUP_FILES=$(ls "${CONFIG_FILE}.termonaut_backup_"* 2>/dev/null | wc -l)
        if [[ $BACKUP_FILES -gt 0 ]]; then
            echo "✅ Backup created successfully"
            echo "   Backup files: $BACKUP_FILES"
        fi
    else
        echo -e "${RED}❌ Hook not found in configuration${NC}"
    fi
else
    echo -e "${RED}❌ Installation failed${NC}"
fi
echo

# Test 2: Force reinstallation
echo -e "${BLUE}🧪 Test 2: Force Reinstallation${NC}"
echo "─────────────────────────────────"
echo "Testing force reinstallation (should replace existing hook)..."

if ./termonaut advanced shell install --force; then
    echo -e "${GREEN}✅ Force reinstallation successful${NC}"

    # Count Termonaut blocks (should be only one)
    BLOCK_COUNT=$(grep -c "# Termonaut shell integration" "$CONFIG_FILE" 2>/dev/null || echo "0")
    if [[ $BLOCK_COUNT -eq 1 ]]; then
        echo "✅ Only one Termonaut block present (no duplicates)"
    else
        echo -e "${YELLOW}⚠️ Found $BLOCK_COUNT Termonaut blocks${NC}"
    fi
else
    echo -e "${RED}❌ Force reinstallation failed${NC}"
fi
echo

# Test 3: Syntax validation
echo -e "${BLUE}🧪 Test 3: Syntax Validation${NC}"
echo "────────────────────────────"
echo "Testing shell syntax validation..."

case "$SHELL_TYPE" in
    "zsh")
        if zsh -n "$CONFIG_FILE" 2>/dev/null; then
            echo -e "${GREEN}✅ Zsh syntax validation passed${NC}"
        else
            echo -e "${RED}❌ Zsh syntax validation failed${NC}"
            echo "Error details:"
            zsh -n "$CONFIG_FILE" 2>&1 | head -5
        fi
        ;;
    "bash")
        if bash -n "$CONFIG_FILE" 2>/dev/null; then
            echo -e "${GREEN}✅ Bash syntax validation passed${NC}"
        else
            echo -e "${RED}❌ Bash syntax validation failed${NC}"
            echo "Error details:"
            bash -n "$CONFIG_FILE" 2>&1 | head -5
        fi
        ;;
esac
echo

# Test 4: Clean uninstallation
echo -e "${BLUE}🧪 Test 4: Clean Uninstallation${NC}"
echo "──────────────────────────────"
echo "Testing safe uninstallation..."

# Create a backup before uninstalling for demo purposes
cp "$CONFIG_FILE" "${CONFIG_FILE}.demo_backup"

if ./termonaut advanced shell uninstall; then
    echo -e "${GREEN}✅ Uninstallation successful${NC}"

    # Verify removal
    if ! grep -q "termonaut" "$CONFIG_FILE" 2>/dev/null; then
        echo "✅ Hook completely removed from configuration"

        # Check that other content is preserved
        if [[ -f "${CONFIG_FILE}.demo_backup" ]]; then
            ORIGINAL_LINES=$(grep -v "termonaut" "${CONFIG_FILE}.demo_backup" | grep -c "." || echo "0")
            CURRENT_LINES=$(grep -c "." "$CONFIG_FILE" || echo "0")

            if [[ $CURRENT_LINES -ge $((ORIGINAL_LINES - 5)) ]]; then
                echo "✅ Other configuration content preserved"
            else
                echo -e "${YELLOW}⚠️ Some content may have been lost${NC}"
                echo "   Original: $ORIGINAL_LINES lines, Current: $CURRENT_LINES lines"
            fi
        fi
    else
        echo -e "${RED}❌ Hook still present in configuration${NC}"
    fi
else
    echo -e "${RED}❌ Uninstallation failed${NC}"
fi
echo

# Test 5: Backup and recovery demonstration
echo -e "${BLUE}🧪 Test 5: Backup and Recovery${NC}"
echo "─────────────────────────────────"
echo "Demonstrating backup and recovery features..."

# Reinstall for demo
./termonaut advanced shell install >/dev/null 2>&1

echo "Current backup files:"
ls -la "${CONFIG_FILE}.termonaut_backup_"* 2>/dev/null || echo "No backup files found"
echo

# Show backup file content preview
LATEST_BACKUP=$(ls -t "${CONFIG_FILE}.termonaut_backup_"* 2>/dev/null | head -1)
if [[ -n "$LATEST_BACKUP" ]]; then
    echo "Latest backup content preview:"
    echo "=============================="
    head -10 "$LATEST_BACKUP"
    echo "... (truncated)"
    echo "=============================="
fi
echo

# Clean up demo backup
rm -f "${CONFIG_FILE}.demo_backup"

# Summary
echo -e "${BLUE}📊 Safety Features Summary${NC}"
echo "─────────────────────────────"
echo "✅ 1. Automatic Backup Creation"
echo "     • Timestamped backups before any modification"
echo "     • Automatic cleanup on successful operations"
echo "     • Preserved on failures for manual recovery"
echo
echo "✅ 2. Precise Block Detection"
echo "     • Shell-specific start/end markers"
echo "     • Handles incomplete or corrupted blocks"
echo "     • Prevents accidental deletion of other content"
echo
echo "✅ 3. Atomic Operations"
echo "     • Write to temporary file first"
echo "     • Atomic rename for consistency"
echo "     • No partial writes or corruption"
echo
echo "✅ 4. Syntax Validation"
echo "     • Shell-specific syntax checking"
echo "     • Automatic rollback on syntax errors"
echo "     • Prevents broken shell configurations"
echo
echo "✅ 5. Error Recovery"
echo "     • Automatic restoration from backup on failure"
echo "     • Detailed error messages with recovery instructions"
echo "     • Graceful handling of edge cases"
echo

echo -e "${GREEN}🎉 Safe Shell Installation Demo Complete!${NC}"
echo
echo -e "${BLUE}💡 Key Benefits:${NC}"
echo "   • No more broken shell configurations"
echo "   • Safe reinstallation and updates"
echo "   • Automatic backup and recovery"
echo "   • Works with existing shell customizations"
echo "   • Handles edge cases and corrupted installations"
echo
echo -e "${BLUE}🔧 Usage:${NC}"
echo "   • Install: termonaut advanced shell install"
echo "   • Force reinstall: termonaut advanced shell install --force"
echo "   • Uninstall: termonaut advanced shell uninstall"
echo "   • Check status: termonaut advanced shell status"
echo