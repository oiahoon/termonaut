#!/bin/bash

# Test script for Enhanced TUI
echo "🚀 Testing Termonaut Enhanced TUI"
echo "=================================="

# Build the project
echo "📦 Building Termonaut..."
cd /Users/johuang/Work/termonaut
go build -o termonaut cmd/termonaut/*.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    
    echo ""
    echo "🎨 Available commands:"
    echo "  ./termonaut tui           - Original TUI"
    echo "  ./termonaut tui-enhanced  - New Enhanced TUI (Beta)"
    echo ""
    echo "🎯 Enhanced TUI Features:"
    echo "  • Modern responsive design"
    echo "  • Multiple themes (Space, Cyberpunk, Minimal, Retro, Nature)"
    echo "  • Avatar system integration"
    echo "  • Improved navigation with Tab/Shift+Tab"
    echo "  • Real-time data updates"
    echo "  • Keyboard shortcuts (R=refresh, S=settings, Q=quit)"
    echo ""
    echo "🚀 Starting Enhanced TUI in 3 seconds..."
    sleep 3
    
    # Launch enhanced TUI
    ./termonaut tui-enhanced
else
    echo "❌ Build failed!"
    exit 1
fi
