#!/bin/bash

# 🔔 智能通知系统测试脚本

echo "🔔 智能通知系统测试"
echo "=================="
echo ""

# 测试1: macOS系统通知
echo "📱 测试1: macOS系统通知"
if command -v osascript >/dev/null 2>&1; then
    echo "✅ osascript 可用，发送系统通知..."
    osascript -e 'display notification "🎉 这是一个系统通知测试!" with title "Termonaut" subtitle "Easter Egg" sound name "Glass"'
    echo "   检查右上角是否出现了系统通知"
else
    echo "❌ osascript 不可用"
fi
echo ""

# 测试2: 终端标题通知
echo "📱 测试2: 终端标题通知"
echo "✅ 设置终端标题..."
printf "\033]0;🎉 Termonaut Easter Egg - 标题通知测试!\007"
echo "   检查终端标题栏是否显示了彩蛋消息"
echo "   3秒后恢复原标题..."
sleep 3
printf "\033]0;Terminal\007"
echo "✅ 标题已恢复"
echo ""

# 测试3: 终端铃声 + 消息
echo "📱 测试3: 终端铃声 + 消息"
echo "✅ 播放铃声并显示消息..."
echo -e "\a🎉 听到铃声了吗？这是一个彩蛋通知！"
echo ""

# 测试4: tmux状态栏 (如果可用)
echo "📱 测试4: tmux状态栏通知"
if [ -n "$TMUX" ]; then
    echo "✅ tmux 环境检测到，发送状态栏消息..."
    tmux display-message -d 3000 "🎉 Termonaut Easter Egg - tmux通知!"
    echo "   检查tmux状态栏是否显示了消息"
else
    echo "❌ 不在tmux环境中"
fi
echo ""

# 测试5: 安全的内联消息
echo "📱 测试5: 安全内联消息"
echo "✅ 显示安全的内联彩蛋消息..."
echo ""
echo "🎉 这是一个安全的内联彩蛋通知！"
echo "   这种方式不会干扰你的终端操作"
echo ""

# 总结
echo "📊 测试总结"
echo "=========="
echo ""
echo "✅ 系统通知: $(command -v osascript >/dev/null 2>&1 && echo '可用' || echo '不可用')"
echo "✅ 终端标题: 总是可用"
echo "✅ 终端铃声: 总是可用"
echo "✅ tmux状态栏: $([ -n "$TMUX" ] && echo '可用' || echo '不可用')"
echo "✅ 内联消息: 总是可用"
echo ""

echo "💡 推荐的通知优先级:"
echo "   1. 系统通知 (最佳用户体验)"
echo "   2. 终端标题 (不干扰操作)"
echo "   3. 终端铃声 + 消息 (兼容性好)"
echo "   4. 内联消息 (安全备选)"
echo ""

echo "🎯 这些方法都不会干扰用户的终端操作！"
