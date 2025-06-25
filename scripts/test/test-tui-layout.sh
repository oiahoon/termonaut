#!/bin/bash

# Script to test TUI layout in different terminal sizes
# This helps verify the layout fixes work correctly

set -e

echo "🧪 Testing TUI layout fixes..."

# Build the latest version
echo "🔨 Building termonaut..."
cd /Users/huangyuyao/OwnWork/termonaut
go build -o termonaut cmd/termonaut/*.go

echo "✅ Build completed"

# Test different terminal sizes
echo ""
echo "📏 Testing different terminal sizes:"
echo ""

echo "1. Current terminal size: $(tput cols)x$(tput lines)"
echo "   - Should show appropriate layout for this size"
echo ""

echo "2. Layout modes to test:"
echo "   termonaut tui --mode compact   # Compact layout"
echo "   termonaut tui --mode full      # Full layout"  
echo "   termonaut tui --mode smart     # Smart adaptive layout"
echo ""

echo "3. Expected behavior by terminal size:"
echo "   < 60 cols:  Icon-only tabs, minimal content"
echo "   60-80 cols: Short tab names, compact layout"
echo "   80-100 cols: Full tab names, standard layout"
echo "   > 100 cols: Wide layout with side-by-side content"
echo ""

echo "4. Layout components that should be visible:"
echo "   ✓ Header with title and level"
echo "   ✓ Tab navigation (responsive)"
echo "   ✓ Main content area (height-constrained)"
echo "   ✓ Footer with help text"
echo ""

echo "5. Common issues that should be fixed:"
echo "   ✓ Tabs not visible in narrow terminals"
echo "   ✓ Content overflowing terminal height"
echo "   ✓ Vertical sections cut off"
echo "   ✓ Poor space utilization"
echo ""

echo "🚀 Quick test (3 seconds each mode):"
echo ""

# Test compact mode
echo "Testing compact mode..."
echo "termonaut tui --mode compact"
echo "(Press 'q' to quit, or wait 3 seconds)"
echo ""

# Test full mode  
echo "Testing full mode..."
echo "termonaut tui --mode full"
echo "(Press 'q' to quit, or wait 3 seconds)"
echo ""

# Test smart mode
echo "Testing smart mode (default)..."
echo "termonaut tui"
echo "(Press 'q' to quit, or wait 3 seconds)"
echo ""

echo "🔍 Manual testing checklist:"
echo "□ Can you see all tab names or icons?"
echo "□ Can you navigate between tabs with Tab/Shift+Tab?"
echo "□ Is the content area properly sized?"
echo "□ Are all sections visible without scrolling?"
echo "□ Does the layout adapt when you resize the terminal?"
echo ""

echo "🐛 If you still see issues:"
echo "1. Try resizing your terminal to at least 80x24"
echo "2. Use 'termonaut tui --mode compact' for small terminals"
echo "3. Check that your terminal supports the required features"
echo ""

echo "✅ TUI layout test script ready!"
echo "Run the commands above to test the fixes."
