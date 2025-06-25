#!/bin/bash

# Script to fix TUI layout issues
# This addresses tab visibility, content overflow, and layout problems

set -e

echo "🔧 Fixing TUI layout issues..."

# Backup the original file
cp /Users/huangyuyao/OwnWork/termonaut/internal/tui/enhanced/dashboard.go /Users/huangyuyao/OwnWork/termonaut/internal/tui/enhanced/dashboard.go.backup

echo "✅ Backup created: dashboard.go.backup"

# The main issues to fix:
# 1. Content height calculation
# 2. Tab navigation visibility
# 3. Proper content area sizing
# 4. Better responsive layout

echo "📝 Applying layout fixes..."

# We'll create a patch for the main View function
cat > /tmp/tui_layout_fix.patch << 'EOF'
--- a/internal/tui/enhanced/dashboard.go
+++ b/internal/tui/enhanced/dashboard.go
@@ -246,12 +246,25 @@ func (d *EnhancedDashboard) View() string {
 	// Footer
 	footer := d.renderFooter()
 	
+	// Calculate available content height
+	headerHeight := lipgloss.Height(header)
+	tabNavHeight := lipgloss.Height(tabNav)
+	footerHeight := lipgloss.Height(footer)
+	availableHeight := d.windowHeight - headerHeight - tabNavHeight - footerHeight - 2 // 2 for margins
+	
+	// Ensure minimum height
+	if availableHeight < 10 {
+		availableHeight = 10
+	}
+	
+	// Apply height constraint to content
+	content = lipgloss.NewStyle().Height(availableHeight).Render(content)
+	
 	// Combine all parts
 	return lipgloss.JoinVertical(
 		lipgloss.Left,
 		header,
 		tabNav,
 		content,
 		footer,
 	)
EOF

echo "🎯 Layout fixes identified:"
echo "  1. Content height calculation and constraint"
echo "  2. Proper space allocation for header, tabs, and footer"
echo "  3. Minimum height guarantee"
echo "  4. Responsive layout improvements"

echo ""
echo "📋 Manual fixes needed in dashboard.go:"
echo ""
echo "1. In the View() function, add content height calculation:"
echo "   - Calculate header, tab, and footer heights"
echo "   - Set available content height"
echo "   - Apply height constraint to content area"
echo ""
echo "2. In renderTabNavigation(), ensure proper width handling:"
echo "   - Handle tab overflow for narrow terminals"
echo "   - Add responsive tab display"
echo ""
echo "3. In content rendering functions, add height awareness:"
echo "   - Limit content to available space"
echo "   - Add scrolling or truncation for overflow"

echo ""
echo "🚀 Quick test commands:"
echo "  termonaut tui --mode compact   # Test compact mode"
echo "  termonaut tui --mode full      # Test full mode"
echo "  termonaut tui --mode smart     # Test smart mode (default)"

echo ""
echo "🔍 Debug terminal size:"
echo "  Current terminal: $(tput cols)x$(tput lines)"
echo "  Recommended minimum: 80x24"
echo "  Optimal size: 120x30+"
