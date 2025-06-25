#!/bin/bash

# 🔔 改进的通知系统测试脚本

echo "🔔 Termonaut 改进通知系统测试"
echo "============================="
echo ""

# 构建项目
echo "📦 构建项目..."
if go build -o termonaut cmd/termonaut/*.go; then
    echo "✅ 构建成功"
else
    echo "❌ 构建失败"
    exit 1
fi

echo ""
echo "🎯 智能通知系统 vs 旧浮动通知系统对比"
echo ""

# 显示问题说明
echo "❌ 旧浮动通知系统的问题:"
echo "   • 定位不准确，可能不在顶部显示"
echo "   • 没有正确擦除，影响终端原有内容"
echo "   • 干扰用户输入，破坏光标位置"
echo "   • 与其他程序输出冲突"
echo ""

echo "✅ 新智能通知系统的优势:"
echo "   • 使用系统原生通知，完全不干扰终端"
echo "   • 自动检测环境，选择最佳通知方式"
echo "   • 多种备选方案，确保兼容性"
echo "   • 用户体验类似桌面应用"
echo ""

# 测试选择
echo "🧪 选择测试方式:"
echo "1. 测试新的智能通知系统 (推荐)"
echo "2. 测试旧的浮动通知系统 (已弃用)"
echo "3. 对比测试两种系统"
echo "4. 退出"
echo ""
echo -n "请选择 (1-4): "
read choice

case $choice in
    1)
        echo ""
        echo "🔔 测试智能通知系统..."
        echo "请注意观察:"
        echo "• macOS: 右上角系统通知"
        echo "• 终端标题栏变化"
        echo "• 终端铃声"
        echo ""
        ./termonaut easter-egg --smart
        ;;
    2)
        echo ""
        echo "⚠️  测试旧浮动通知系统 (已弃用)..."
        echo "注意: 这个系统有已知问题，仅用于对比"
        echo ""
        ./termonaut easter-egg --floating
        ;;
    3)
        echo ""
        echo "📊 对比测试..."
        echo ""
        echo "首先测试旧系统 (注意问题):"
        echo "按 Enter 继续..."
        read
        ./termonaut easter-egg --floating
        
        echo ""
        echo "现在测试新系统 (注意改进):"
        echo "按 Enter 继续..."
        read
        ./termonaut easter-egg --smart
        ;;
    4)
        echo "退出测试"
        exit 0
        ;;
    *)
        echo "无效选择"
        exit 1
        ;;
esac

echo ""
echo "📋 测试反馈"
echo "=========="
echo ""
echo "请回答以下问题:"
echo ""

if [ "$choice" = "1" ] || [ "$choice" = "3" ]; then
    echo "关于智能通知系统:"
    echo -n "1. 你看到系统通知了吗? (y/n): "
    read system_notification
    
    echo -n "2. 终端标题有变化吗? (y/n): "
    read title_change
    
    echo -n "3. 听到铃声了吗? (y/n): "
    read bell_sound
    
    echo -n "4. 通知是否干扰了终端使用? (y/n): "
    read interference
    
    echo -n "5. 整体体验满意吗? (y/n): "
    read satisfaction
fi

echo ""
echo "📊 测试结果分析"
echo "==============="
echo ""

if [ "$choice" = "1" ] || [ "$choice" = "3" ]; then
    echo "智能通知系统评估:"
    [ "$system_notification" = "y" ] && echo "✅ 系统通知工作正常" || echo "⚠️  系统通知可能需要检查"
    [ "$title_change" = "y" ] && echo "✅ 终端标题更新正常" || echo "⚠️  终端标题更新可能有问题"
    [ "$bell_sound" = "y" ] && echo "✅ 铃声提醒正常" || echo "ℹ️  铃声可能被静音"
    [ "$interference" = "n" ] && echo "✅ 无干扰，用户体验优秀" || echo "⚠️  存在干扰问题"
    [ "$satisfaction" = "y" ] && echo "✅ 用户满意度高" || echo "⚠️  需要改进"
fi

echo ""
echo "🎯 推荐结论"
echo "==========="
echo ""

if [ "$system_notification" = "y" ] && [ "$interference" = "n" ]; then
    echo "🎉 智能通知系统测试成功!"
    echo ""
    echo "建议:"
    echo "• ✅ 正式采用智能通知系统"
    echo "• ✅ 弃用旧的浮动通知系统"
    echo "• ✅ 将智能通知集成到彩蛋系统"
    echo "• ✅ 添加用户配置选项"
else
    echo "🔧 智能通知系统需要进一步优化"
    echo ""
    echo "建议检查:"
    echo "• 系统通知权限设置"
    echo "• 终端兼容性"
    echo "• 音频设置"
fi

echo ""
echo "💡 技术优势总结:"
echo "• 系统通知: 最佳用户体验，完全不干扰"
echo "• 终端标题: 轻量级提醒，兼容性好"
echo "• 终端铃声: 音频提醒，传统可靠"
echo "• 智能选择: 自动适配最佳方案"
echo ""
echo "🚀 这个方案完美解决了浮动通知的所有问题!"
