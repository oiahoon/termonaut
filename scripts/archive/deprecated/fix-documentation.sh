#!/bin/bash

# Documentation Consistency Fix Script
# Updates outdated command references across documentation files

echo "🔧 Fixing documentation consistency..."

# Replace 'termonaut init' with 'termonaut advanced shell install'
find . -name "*.md" -type f -exec sed -i.bak 's/termonaut init/termonaut advanced shell install/g' {} \;

# Clean up backup files
find . -name "*.bak" -type f -delete

echo "✅ Updated all documentation files"
echo "📋 Fixed references:"
echo "   termonaut init → termonaut advanced shell install"

# Verify changes
echo
echo "🔍 Verification:"
grep -r "termonaut advanced shell install" docs/ --include="*.md" | wc -l | xargs echo "Updated references found:" 