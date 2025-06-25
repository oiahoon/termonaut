#!/bin/bash

# Safe Enhanced TUI Test Script
echo "🚀 Testing Termonaut Enhanced TUI (Safe Version)"
echo "================================================="

# Build the project
echo "📦 Building Termonaut..."
cd /Users/johuang/Work/termonaut
go build -o termonaut cmd/termonaut/*.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    
    # First, let's check if the basic command works
    echo ""
    echo "🔍 Testing basic command functionality..."
    ./termonaut --help > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ Basic command works"
    else
        echo "❌ Basic command failed"
        exit 1
    fi
    
    # Check if enhanced TUI command is available
    echo ""
    echo "🔍 Checking enhanced TUI command..."
    ./termonaut tui-enhanced --help > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ Enhanced TUI command is available"
    else
        echo "❌ Enhanced TUI command not found"
        exit 1
    fi
    
    echo ""
    echo "🎨 Enhanced TUI Features:"
    echo "  • 🏠 Home - Personal dashboard with avatar and quick stats"
    echo "  • 📊 Analytics - Deep productivity insights and charts"
    echo "  • 🎮 Gamification - XP, levels, achievements, and progress"
    echo "  • 🔥 Activity - Time-based activity visualization and heatmaps"
    echo "  • 🛠️ Tools - Utility functions and integrations"
    echo "  • ⚙️ Settings - Theme customization and configuration"
    echo ""
    echo "🎯 Navigation:"
    echo "  • Tab/L/→ - Next tab"
    echo "  • Shift+Tab/H/← - Previous tab"
    echo "  • R/F5 - Refresh data"
    echo "  • S - Jump to settings"
    echo "  • Q/Ctrl+C - Quit"
    echo ""
    echo "🎨 Available Themes:"
    echo "  • Space (default) - Purple space theme"
    echo "  • Cyberpunk - Neon colors"
    echo "  • Minimal - Clean black & white"
    echo "  • Retro - Vintage colors"
    echo "  • Nature - Green nature theme"
    echo ""
    echo "🚀 Starting Enhanced TUI..."
    echo "   (The interface should load safely now)"
    echo "   (Press 'q' to quit when you're done testing)"
    echo ""
    
    # Launch enhanced TUI with error handling
    ./termonaut tui-enhanced
    exit_code=$?
    
    echo ""
    if [ $exit_code -eq 0 ]; then
        echo "🎉 Enhanced TUI test completed successfully!"
    else
        echo "⚠️  Enhanced TUI exited with code: $exit_code"
    fi
    
    echo ""
    echo "💡 Next steps:"
    echo "  1. Test different tabs with Tab/Shift+Tab"
    echo "  2. Try the refresh function with 'r'"
    echo "  3. Check the responsive layout by resizing your terminal"
    echo "  4. Compare with original TUI: ./termonaut tui"
    echo ""
    echo "🐛 If you encounter issues:"
    echo "  1. Check terminal size (recommended: 100+ columns)"
    echo "  2. Ensure database is initialized: ./termonaut init"
    echo "  3. Try the original TUI: ./termonaut tui"
    
else
    echo "❌ Build failed!"
    echo "Please check the error messages above."
    exit 1
fi
