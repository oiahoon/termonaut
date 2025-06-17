#!/bin/bash

# Empty Command Stats Demo Script for Termonaut
# This script demonstrates the new empty command stats feature

set -e

echo "🎯 Termonaut Empty Command Stats Demo"
echo "════════════════════════════════════════"
echo

# Build the binary first
echo "🔨 Building Termonaut..."
if ! go build -o termonaut cmd/termonaut/*.go; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"
echo

echo "📖 About Empty Command Stats Feature:"
echo "────────────────────────────────────"
echo "• When you press Enter on an empty command line, Termonaut shows quick stats"
echo "• No more need to type 'termonaut stats' - just hit Enter!"
echo "• Configurable display modes: minimal vs rich"
echo "• Respects your privacy and display settings"
echo

# 1. Check current configuration
echo "1️⃣  Current Configuration"
echo "─────────────────────"
./termonaut config get empty_command_stats || echo "Feature status: unknown"
echo

# 2. Enable the feature
echo "2️⃣  Enabling Empty Command Stats"
echo "─────────────────────────────"
./termonaut config set empty_command_stats true
echo

# 3. Show configuration options
echo "3️⃣  Display Mode Options"
echo "─────────────────────"
echo "Current theme: $(./termonaut config get theme)"
echo "Current display mode: $(./termonaut config get display_mode)"
echo
echo "🔹 Available themes:"
echo "   • minimal: Clean text-only display"
echo "   • emoji: Rich display with emojis and progress bars"
echo "   • ascii: ASCII art styling"
echo

# 4. Test different display modes
echo "4️⃣  Testing Different Display Modes"
echo "───────────────────────────────────"

echo "📱 Setting minimal theme..."
./termonaut config set theme minimal
echo "Test: Empty command will show minimal stats like:"
echo "📊 Lv.1 | 0 cmds | 0 streak | 100 XP"
echo

echo "🎨 Setting emoji theme..."
./termonaut config set theme emoji
echo "Test: Empty command will show rich stats like:"
echo "🚀 Level 1 [████████] 150 XP"
echo "🎯 5 commands today | 🔥 3 day streak"
echo "👑 git status (10x)"
echo

# 5. Show how to disable the feature
echo "5️⃣  Disabling the Feature"
echo "─────────────────────"
echo "To disable empty command stats:"
echo "  termonaut config set empty_command_stats false"
echo

# 6. Integration with shell hooks
echo "6️⃣  Shell Hook Integration"
echo "──────────────────────"
echo "📋 How it works:"
echo "  1. Your shell hook captures the empty command"
echo "  2. Termonaut detects trimmed command is empty"
echo "  3. If empty_command_stats=true, shows quick stats"
echo "  4. If disabled, silently ignores empty commands"
echo

echo "🔧 Shell Integration Status:"
if ./termonaut advanced shell status 2>/dev/null; then
    echo "✅ Shell hooks are installed and working"
else
    echo "❌ Shell hooks need to be installed"
    echo "   Run: termonaut advanced shell install"
fi
echo

# 7. Easter Egg integration
echo "7️⃣  Integration with Easter Eggs"
echo "────────────────────────────"
echo "🥚 Empty commands don't trigger Easter Eggs, but regular commands do!"
echo "🎮 Easter Eggs status: $(./termonaut config get easter_eggs_enabled)"
echo

# 8. Configuration examples
echo "8️⃣  Configuration Examples"
echo "─────────────────────────"
echo "🎯 Minimal productivity setup:"
echo "  termonaut config set theme minimal"
echo "  termonaut config set empty_command_stats true"
echo "  termonaut config set display_mode enter"
echo

echo "🎮 Full gamification experience:"
echo "  termonaut config set theme emoji"
echo "  termonaut config set empty_command_stats true"
echo "  termonaut config set easter_eggs_enabled true"
echo "  termonaut config set show_gamification true"
echo

echo "😶 Quiet/stealth mode:"
echo "  termonaut config set empty_command_stats false"
echo "  termonaut config set display_mode off"
echo "  termonaut config set easter_eggs_enabled false"
echo

# Summary
echo "📋 Feature Summary"
echo "─────────────────"
echo "✅ 1. Empty Command Detection - Automatically detects when you press Enter"
echo "✅ 2. Quick Stats Display - Shows level, commands, streak, and XP instantly"
echo "✅ 3. Configurable Themes - Minimal, Emoji, or ASCII styling"
echo "✅ 4. Privacy Aware - Respects your display and privacy settings"
echo "✅ 5. Shell Integration - Works seamlessly with existing hooks"
echo "✅ 6. Performance Optimized - Fast, silent background operation"
echo

echo "🎉 Demo Complete!"
echo
echo "💡 Next Steps:"
echo "   1. Install shell hooks: 'termonaut advanced shell install'"
echo "   2. Configure your preferred theme: 'termonaut config set theme emoji'"
echo "   3. Enable the feature: 'termonaut config set empty_command_stats true'"
echo "   4. Press Enter on an empty command line to see your stats!"
echo
echo "🔗 More Info:"
echo "   • Run 'termonaut config --help' for all configuration options"
echo "   • Check 'termonaut --help' for full command reference"
echo "   • Visit the docs for advanced configuration tips" 