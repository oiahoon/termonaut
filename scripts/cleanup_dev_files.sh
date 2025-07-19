#!/bin/bash

# Development Cleanup Script
# Removes temporary files and redundant documentation from development process

echo "🧹 Cleaning up development files..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counter for cleaned files
CLEANED_COUNT=0

# Function to remove files if they exist
cleanup_files() {
    local pattern="$1"
    local description="$2"
    
    files=$(find . -name "$pattern" 2>/dev/null)
    if [ -n "$files" ]; then
        echo -e "${YELLOW}Removing $description...${NC}"
        echo "$files" | while read -r file; do
            rm -f "$file"
            echo "  Removed: $file"
            CLEANED_COUNT=$((CLEANED_COUNT + 1))
        done
    fi
}

# Clean temporary files
cleanup_files "*.tmp" "temporary files"
cleanup_files "*.bak" "backup files"
cleanup_files "*~" "editor backup files"
cleanup_files ".DS_Store" "macOS system files"
cleanup_files "Thumbs.db" "Windows thumbnail files"

# Clean Go build artifacts (if any)
cleanup_files "*.test" "Go test binaries"
cleanup_files "coverage.out" "Go coverage files"

# Clean log files that might have been created during development
cleanup_files "debug.log" "debug log files"
cleanup_files "test.log" "test log files"

# Remove empty directories in tests (if any)
find tests/ -type d -empty -delete 2>/dev/null || true

# Organize documentation
echo -e "${YELLOW}Organizing documentation...${NC}"

# Ensure docs directory structure is clean
mkdir -p docs/reports
mkdir -p docs/development

# Move completion reports to reports directory
if [ -f "docs/TASK2_1_COMPLETION_REPORT.md" ]; then
    mv docs/TASK2_1_COMPLETION_REPORT.md docs/reports/
    echo "  Moved: Task 2.1 completion report"
fi

if [ -f "docs/TASK2_2_COMPLETION_REPORT.md" ]; then
    mv docs/TASK2_2_COMPLETION_REPORT.md docs/reports/
    echo "  Moved: Task 2.2 completion report"
fi

if [ -f "docs/PHASE2_START_REPORT.md" ]; then
    mv docs/PHASE2_START_REPORT.md docs/reports/
    echo "  Moved: Phase 2 start report"
fi

if [ -f "docs/PHASE2_COMPLETION_REPORT.md" ]; then
    mv docs/PHASE2_COMPLETION_REPORT.md docs/reports/
    echo "  Moved: Phase 2 completion report"
fi

# Clean up any duplicate or redundant documentation
echo -e "${YELLOW}Checking for redundant documentation...${NC}"

# Remove any duplicate README files
find . -name "README*.md" -not -path "./README.md" | while read -r file; do
    if [ "$file" != "./README.md" ]; then
        echo "  Found duplicate README: $file"
        # Don't auto-remove, just report
    fi
done

# Verify test scripts are executable
echo -e "${YELLOW}Ensuring test scripts are executable...${NC}"
find tests/scripts/ -name "*.sh" -exec chmod +x {} \;

# Create a summary of the cleanup
echo ""
echo -e "${GREEN}✅ Cleanup completed!${NC}"
echo ""
echo "📊 Cleanup Summary:"
echo "=================="

# Count files by type
echo "Test files: $(find tests/ -name "*.go" | wc -l)"
echo "Test scripts: $(find tests/scripts/ -name "*.sh" | wc -l)"
echo "Source files: $(find internal/ -name "*.go" | wc -l)"
echo "Documentation: $(find docs/ -name "*.md" | wc -l)"

# Show directory structure
echo ""
echo "📁 Final Directory Structure:"
echo "============================"
echo "Project Structure:"
tree -d -L 3 2>/dev/null || find . -type d -not -path "./.git*" | head -20

echo ""
echo -e "${GREEN}🎉 Development cleanup completed successfully!${NC}"
echo "The project is now clean and ready for production use."
