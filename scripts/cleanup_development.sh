#!/bin/bash

# Development Cleanup Script
# Cleans up temporary files and development artifacts

set -e

echo "🧹 Starting Development Cleanup"
echo "==============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Cleanup counters
FILES_REMOVED=0
DIRS_REMOVED=0
BYTES_SAVED=0

# Helper function to remove files/directories
cleanup_item() {
    local item="$1"
    local description="$2"
    
    if [ -e "$item" ]; then
        local size=$(du -sb "$item" 2>/dev/null | cut -f1 || echo "0")
        
        if [ -d "$item" ]; then
            echo -n "Removing directory: $description... "
            rm -rf "$item"
            DIRS_REMOVED=$((DIRS_REMOVED + 1))
        else
            echo -n "Removing file: $description... "
            rm -f "$item"
            FILES_REMOVED=$((FILES_REMOVED + 1))
        fi
        
        BYTES_SAVED=$((BYTES_SAVED + size))
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "Skipping: $description (not found)"
    fi
}

echo ""
echo "🗑️ Cleaning Temporary Files"
echo "==========================="

# Remove common temporary files
cleanup_item "*.tmp" "Temporary files"
cleanup_item "*.temp" "Temp files"
cleanup_item ".DS_Store" "macOS metadata files"
cleanup_item "Thumbs.db" "Windows thumbnail cache"

# Remove build artifacts (if any)
cleanup_item "bin/" "Build directory"
cleanup_item "dist/" "Distribution directory"
cleanup_item "build/" "Build artifacts"

# Remove test artifacts
cleanup_item "coverage.out" "Go coverage file"
cleanup_item "coverage.html" "Coverage HTML report"
cleanup_item "*.test" "Test binaries"

# Remove editor temporary files
cleanup_item "*.swp" "Vim swap files"
cleanup_item "*.swo" "Vim swap files"
cleanup_item "*~" "Editor backup files"
cleanup_item ".vscode/settings.json" "VS Code local settings"

echo ""
echo "📝 Cleaning Development Logs"
echo "============================"

# Remove development logs
cleanup_item "debug.log" "Debug log file"
cleanup_item "error.log" "Error log file"
cleanup_item "*.log" "Log files"

# Remove test databases
cleanup_item "test.db" "Test database"
cleanup_item "*.db-journal" "SQLite journal files"
cleanup_item "*.db-wal" "SQLite WAL files"
cleanup_item "*.db-shm" "SQLite shared memory files"

echo ""
echo "🔧 Cleaning Development Tools"
echo "============================="

# Remove development tool artifacts
cleanup_item ".golangci-lint" "Linter cache"
cleanup_item ".mypy_cache/" "Python type checker cache"
cleanup_item "__pycache__/" "Python bytecode cache"
cleanup_item "*.pyc" "Python compiled files"
cleanup_item "node_modules/" "Node.js dependencies (if any)"

echo ""
echo "📊 Cleaning Analysis Files"
echo "=========================="

# Remove analysis and profiling files
cleanup_item "cpu.prof" "CPU profile"
cleanup_item "mem.prof" "Memory profile"
cleanup_item "profile.out" "Profile output"
cleanup_item "trace.out" "Trace output"

echo ""
echo "🧪 Cleaning Test Artifacts"
echo "=========================="

# Remove test-specific temporary files
cleanup_item "testdata/tmp/" "Test temporary data"
cleanup_item "tests/tmp/" "Test temporary directory"
cleanup_item "*.test.log" "Test log files"

# Keep important test files but remove temporary ones
find tests/ -name "*.tmp" -type f -delete 2>/dev/null || true
find tests/ -name "temp_*" -type f -delete 2>/dev/null || true

echo ""
echo "📁 Organizing Documentation"
echo "==========================="

# Check for duplicate or outdated documentation
if [ -f "docs/OLD_README.md" ]; then
    cleanup_item "docs/OLD_README.md" "Old README file"
fi

if [ -f "docs/DEPRECATED.md" ]; then
    cleanup_item "docs/DEPRECATED.md" "Deprecated documentation"
fi

# Remove empty directories in docs
find docs/ -type d -empty -delete 2>/dev/null || true

echo ""
echo "🔍 Checking for Large Files"
echo "==========================="

# Find and report large files (>1MB)
echo "Large files (>1MB) in project:"
find . -type f -size +1M -not -path "./.git/*" -exec ls -lh {} \; 2>/dev/null | head -10 || echo "No large files found"

echo ""
echo "📈 Cleanup Summary"
echo "=================="

# Convert bytes to human readable format
if [ $BYTES_SAVED -gt 1048576 ]; then
    SIZE_MB=$((BYTES_SAVED / 1048576))
    echo "Space saved: ${SIZE_MB}MB"
elif [ $BYTES_SAVED -gt 1024 ]; then
    SIZE_KB=$((BYTES_SAVED / 1024))
    echo "Space saved: ${SIZE_KB}KB"
else
    echo "Space saved: ${BYTES_SAVED} bytes"
fi

echo "Files removed: $FILES_REMOVED"
echo "Directories removed: $DIRS_REMOVED"

echo ""
echo "🔍 Final Project Status"
echo "======================="

# Show current project size
PROJECT_SIZE=$(du -sh . 2>/dev/null | cut -f1 || echo "Unknown")
echo "Current project size: $PROJECT_SIZE"

# Count files by type
echo ""
echo "File type summary:"
echo "Go files: $(find . -name "*.go" -not -path "./.git/*" | wc -l)"
echo "Markdown files: $(find . -name "*.md" -not -path "./.git/*" | wc -l)"
echo "Shell scripts: $(find . -name "*.sh" -not -path "./.git/*" | wc -l)"
echo "Test files: $(find . -name "*_test.go" -not -path "./.git/*" | wc -l)"

echo ""
echo "📋 Recommendations"
echo "=================="

# Check for potential issues
if [ -d ".git" ]; then
    GIT_SIZE=$(du -sh .git 2>/dev/null | cut -f1 || echo "Unknown")
    echo "Git repository size: $GIT_SIZE"
    
    # Check for large files in git history
    LARGE_FILES=$(git rev-list --objects --all | git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)' | sed -n 's/^blob //p' | sort --numeric-sort --key=2 | tail -5 2>/dev/null || echo "")
    
    if [ -n "$LARGE_FILES" ]; then
        echo -e "${YELLOW}⚠ Large files detected in git history${NC}"
        echo "Consider using git-lfs for large files"
    fi
fi

# Check for security issues
if find . -name "*.key" -o -name "*.pem" -o -name "*.p12" -not -path "./.git/*" | grep -q .; then
    echo -e "${RED}⚠ Potential security files detected${NC}"
    echo "Please review and ensure no sensitive files are committed"
fi

# Check for TODO items
TODO_COUNT=$(grep -r "TODO\|FIXME\|HACK" --include="*.go" --include="*.md" . 2>/dev/null | wc -l || echo "0")
if [ "$TODO_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}📝 TODO items remaining: $TODO_COUNT${NC}"
    echo "Consider addressing remaining TODO items"
fi

echo ""
if [ $FILES_REMOVED -gt 0 ] || [ $DIRS_REMOVED -gt 0 ]; then
    echo -e "${GREEN}✅ Cleanup completed successfully!${NC}"
    echo "Project is now clean and ready for production."
else
    echo -e "${GREEN}✅ Project was already clean!${NC}"
    echo "No cleanup was necessary."
fi

echo ""
echo "🚀 Next Steps"
echo "============"
echo "1. Run tests to ensure everything still works"
echo "2. Commit any remaining changes"
echo "3. Create a release tag if ready"
echo "4. Update documentation if needed"

exit 0
