#!/bin/bash

# 🚀 快速浮动通知测试脚本

echo "🎭 快速测试浮动通知功能"
echo "========================"
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
echo "🎯 开始浮动通知测试..."
echo "请观察终端顶部的通知效果"
echo ""

# 等待2秒让用户准备
sleep 2

# 运行浮动通知测试
./termonaut easter-egg --floating

echo ""
echo "✅ 测试完成!"
echo ""
echo "💡 如果你看到了浮动在顶部的通知，说明功能正常工作!"
echo "🔧 如需详细测试，请运行: ./test_floating_notifications.sh"
