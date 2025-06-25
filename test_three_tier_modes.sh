#!/bin/bash

# Three-Tier Viewing Modes Test Script
echo "🎯 Testing Termonaut Three-Tier Viewing Modes"
echo "=============================================="

# Build the project
echo "📦 Building Termonaut with three-tier modes..."
cd /Users/johuang/Work/termonaut
go build -o termonaut cmd/termonaut/*.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    
    echo ""
    echo "🎨 Three-Tier Architecture:"
    echo "  1️⃣  极简模式: termonaut stats (shell输出，最快速)"
    echo "  2️⃣  普通模式: termonaut tui-compact (紧凑TUI，平衡体验)"
    echo "  3️⃣  完整模式: termonaut tui-enhanced (完整TUI，沉浸体验)"
    echo ""
    echo "🚀 Avatar Technology Stack:"
    echo "  • DiceBear API - 优秀的头像生成服务"
    echo "  • ascii-image-converter - Go生态最佳ASCII转换库"
    echo "  • 动态尺寸支持: 8x4 到 65x32 字符"
    echo ""
    
    echo "🧪 Testing Mode 1: 极简模式 (Shell Output)"
    echo "─────────────────────────────────────────────"
    echo "Command: ./termonaut stats --today"
    echo ""
    ./termonaut stats --today
    echo ""
    
    echo "Press Enter to test Mode 2: 普通模式 (Compact TUI)..."
    read -r
    
    echo ""
    echo "🧪 Testing Mode 2: 普通模式 (Compact TUI)"
    echo "─────────────────────────────────────────────"
    echo "Command: ./termonaut tui-compact"
    echo ""
    echo "Features:"
    echo "  • Small avatars (8-25 characters wide)"
    echo "  • Optimized for quick viewing"
    echo "  • Works on smaller terminals (40+ chars)"
    echo "  • Fast loading and response"
    echo ""
    echo "Navigation: Tab=next, r=refresh, q=quit"
    echo ""
    
    ./termonaut tui-compact
    
    echo ""
    echo "Press Enter to test Mode 3: 完整模式 (Enhanced TUI)..."
    read -r
    
    echo ""
    echo "🧪 Testing Mode 3: 完整模式 (Enhanced TUI)"
    echo "─────────────────────────────────────────────"
    echo "Command: ./termonaut tui-enhanced"
    echo ""
    echo "Features:"
    echo "  • Large avatars (35-70 characters wide)"
    echo "  • Immersive experience"
    echo "  • Full feature set"
    echo "  • Best for wide terminals (100+ chars)"
    echo ""
    echo "Navigation: Tab=next, r=refresh, q=quit"
    echo ""
    
    ./termonaut tui-enhanced
    
    echo ""
    echo "🎉 Three-Tier Mode Testing Complete!"
    echo ""
    echo "📊 Mode Comparison Summary:"
    echo "┌─────────────┬─────────────┬─────────────┬─────────────┐"
    echo "│    Mode     │   Speed     │ Avatar Size │ Best For    │"
    echo "├─────────────┼─────────────┼─────────────┼─────────────┤"
    echo "│ 极简模式    │   最快      │    3行      │ 快速查看    │"
    echo "│ 普通模式    │   快速      │  8-25字符   │ 日常监控    │"
    echo "│ 完整模式    │   丰富      │ 35-70字符   │ 深度分析    │"
    echo "└─────────────┴─────────────┴─────────────┴─────────────┘"
    echo ""
    echo "💡 Usage Recommendations:"
    echo "  🏃‍♂️ Quick check: termonaut stats --today"
    echo "  📊 Daily monitoring: termonaut tui-compact"
    echo "  🎮 Deep analysis: termonaut tui-enhanced"
    echo ""
    echo "🎨 Avatar Technology:"
    echo "  ✅ DiceBear API integration"
    echo "  ✅ ascii-image-converter library"
    echo "  ✅ Dynamic size adaptation (8x4 to 65x32)"
    echo "  ✅ Fallback default avatars"
    echo "  ✅ Real-time terminal size detection"
    
else
    echo "❌ Build failed!"
    echo "Please check the error messages above."
    exit 1
fi
