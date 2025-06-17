#!/bin/bash

# Termonaut GitHub Integration Demo
# This script demonstrates the new GitHub integration features in v0.6.0

echo "🚀 Termonaut GitHub Integration Demo - v0.6.0"
echo "============================================="
echo

# Make sure termonaut is built
if [ ! -f "./termonaut" ]; then
    echo "Building termonaut..."
    go build -o termonaut cmd/termonaut/*.go
    echo "✅ Build complete"
    echo
fi

echo "📊 1. Generating GitHub Badges"
echo "-------------------------------"
echo "This creates dynamic badges for your GitHub profile:"
echo
./termonaut github badges generate
echo
echo "You can copy these URLs and paste them into your GitHub README.md file!"
echo

echo "📈 2. Generating Profile Summary"
echo "-------------------------------"
echo "This creates a comprehensive profile with stats and achievements:"
echo
./termonaut github profile generate
echo

echo "🤖 3. Available GitHub Actions Workflows"
echo "----------------------------------------"
echo "These workflows can automatically update your badges and stats:"
echo

# Since we don't have GitHub config set up, show the available templates manually
echo "Available GitHub Actions workflow templates:"
echo
echo "• termonaut-stats-update"
echo "  Automatically update Termonaut badges and stats"
echo
echo "• termonaut-profile-sync"
echo "  Sync Termonaut profile data to repository"
echo
echo "• termonaut-weekly-report"
echo "  Generate weekly productivity reports"
echo

echo "📁 4. Export Examples"
echo "--------------------"
echo "Creating example exports..."

# Create examples directory
mkdir -p examples/exports

# Generate JSON badge data
echo "Generating JSON badge data..."
./termonaut github badges generate --format json --output examples/exports/badges.json

# Generate markdown profile
echo "Generating markdown profile..."
./termonaut github profile generate --format markdown --output examples/exports/profile.md

echo "✅ Examples saved to examples/exports/"
echo

echo "🎯 Setup Instructions"
echo "====================="
echo
echo "To set up GitHub integration:"
echo "1. Add GitHub config to your Termonaut config file:"
echo "   ~/.termonaut/config.toml"
echo
echo "   [github]"
echo "   repo_owner = \"your-username\""
echo "   repo_name = \"your-repo\""
echo
echo "2. Generate a workflow file:"
echo "   ./termonaut github actions generate termonaut-stats-update"
echo
echo "3. Commit the workflow to your repository:"
echo "   git add .github/workflows/termonaut-stats-update.yml"
echo "   git commit -m \"Add Termonaut stats automation\""
echo "   git push"
echo
echo "4. Add badges to your README.md:"
echo "   Copy the badge URLs from step 1 above"
echo

echo "📚 More Features"
echo "================"
echo "• Social media snippets: ./termonaut github profile social"
echo "• Badge endpoints for dynamic updates"
echo "• Workflow templates for automation"
echo "• Profile export in multiple formats"
echo

echo "🎉 GitHub Integration Demo Complete!"
echo "Visit the GitHub repository for more documentation and examples."