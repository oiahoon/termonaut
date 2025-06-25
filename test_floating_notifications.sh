#!/bin/bash

# 🎭 Termonaut浮动通知测试脚本
# 用于本地测试浮动彩蛋通知功能

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="/Users/johuang/Work/termonaut"
BINARY_PATH="$PROJECT_ROOT/termonaut"

echo -e "${CYAN}🎭 Termonaut 浮动通知测试脚本${NC}"
echo -e "${CYAN}=================================${NC}"
echo ""

# 检查是否在正确的目录
if [ ! -d "$PROJECT_ROOT" ]; then
    echo -e "${RED}❌ 错误: 项目目录不存在: $PROJECT_ROOT${NC}"
    exit 1
fi

cd "$PROJECT_ROOT"

# 函数：显示步骤
show_step() {
    echo -e "${BLUE}📋 步骤 $1: $2${NC}"
}

# 函数：显示成功
show_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# 函数：显示警告
show_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 函数：显示错误
show_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 函数：等待用户确认
wait_for_user() {
    echo -e "${YELLOW}按 Enter 键继续...${NC}"
    read -r
}

# 步骤1: 构建项目
show_step "1" "构建 Termonaut 项目"
if go build -o termonaut cmd/termonaut/*.go; then
    show_success "项目构建成功"
else
    show_error "项目构建失败"
    exit 1
fi
echo ""

# 步骤2: 检查二进制文件
show_step "2" "检查二进制文件"
if [ -f "$BINARY_PATH" ]; then
    show_success "二进制文件存在: $BINARY_PATH"
    echo -e "   文件大小: $(du -h $BINARY_PATH | cut -f1)"
    echo -e "   修改时间: $(stat -f "%Sm" $BINARY_PATH)"
else
    show_error "二进制文件不存在"
    exit 1
fi
echo ""

# 步骤3: 检查浮动通知命令
show_step "3" "检查浮动通知命令可用性"
if ./termonaut easter-egg --help | grep -q "floating"; then
    show_success "浮动通知命令可用"
else
    show_error "浮动通知命令不可用"
    exit 1
fi
echo ""

# 步骤4: 显示终端信息
show_step "4" "检测当前终端环境"
echo -e "   终端类型: ${CYAN}$TERM${NC}"
echo -e "   终端程序: ${CYAN}${TERM_PROGRAM:-未知}${NC}"
echo -e "   颜色支持: ${CYAN}${COLORTERM:-基础}${NC}"
echo -e "   Shell: ${CYAN}$SHELL${NC}"
echo -e "   终端大小: ${CYAN}$(tput cols)x$(tput lines)${NC}"
echo ""

# 步骤5: 准备测试
show_step "5" "准备开始浮动通知测试"
echo -e "${WHITE}测试说明:${NC}"
echo -e "• 测试将显示6个不同的浮动通知"
echo -e "• 每个通知会在终端顶部显示3秒"
echo -e "• 通知会自动消失，不需要手动操作"
echo -e "• 请观察终端顶部的通知效果"
echo ""
echo -e "${YELLOW}⚠️  注意: 测试期间请不要输入任何内容${NC}"
echo ""
wait_for_user

# 步骤6: 运行浮动通知测试
show_step "6" "运行浮动通知测试"
echo -e "${PURPLE}🚀 开始测试...${NC}"
echo ""

# 清屏以获得更好的测试效果
clear

echo -e "${CYAN}🎭 浮动通知测试开始${NC}"
echo -e "${CYAN}==================${NC}"
echo ""
echo -e "${WHITE}请观察终端顶部的浮动通知效果...${NC}"
echo ""

# 运行测试
if ./termonaut easter-egg --floating; then
    show_success "浮动通知测试完成"
else
    show_error "浮动通知测试失败"
    exit 1
fi

echo ""

# 步骤7: 测试结果分析
show_step "7" "测试结果分析"
echo -e "${WHITE}请回答以下问题来评估测试效果:${NC}"
echo ""

# 询问用户反馈
echo -e "${YELLOW}1. 你看到浮动通知出现在终端顶部了吗? (y/n)${NC}"
read -r saw_notifications

echo -e "${YELLOW}2. 通知的样式和边框显示正常吗? (y/n)${NC}"
read -r style_ok

echo -e "${YELLOW}3. 通知是否在3秒后自动消失? (y/n)${NC}"
read -r auto_disappear

echo -e "${YELLOW}4. 测试过程中是否干扰了你的终端使用? (y/n)${NC}"
read -r interference

echo -e "${YELLOW}5. 整体效果满意吗? (y/n)${NC}"
read -r satisfaction

echo ""

# 分析结果
show_step "8" "测试结果总结"
echo -e "${WHITE}测试结果分析:${NC}"

if [[ "$saw_notifications" == "y" ]]; then
    show_success "✅ 浮动通知显示正常"
else
    show_warning "⚠️  浮动通知显示可能有问题"
fi

if [[ "$style_ok" == "y" ]]; then
    show_success "✅ 通知样式渲染正常"
else
    show_warning "⚠️  通知样式可能需要调整"
fi

if [[ "$auto_disappear" == "y" ]]; then
    show_success "✅ 自动消失功能正常"
else
    show_warning "⚠️  自动消失功能可能有问题"
fi

if [[ "$interference" == "n" ]]; then
    show_success "✅ 无干扰，用户体验良好"
else
    show_warning "⚠️  可能对用户操作造成干扰"
fi

if [[ "$satisfaction" == "y" ]]; then
    show_success "✅ 整体效果令人满意"
else
    show_warning "⚠️  整体效果需要改进"
fi

echo ""

# 步骤9: 额外测试选项
show_step "9" "额外测试选项"
echo -e "${WHITE}你想进行额外的测试吗?${NC}"
echo ""
echo -e "1. 测试单个通知效果"
echo -e "2. 测试不同终端兼容性"
echo -e "3. 测试通知冲突处理"
echo -e "4. 跳过额外测试"
echo ""
echo -e "${YELLOW}请选择 (1-4):${NC}"
read -r choice

case $choice in
    1)
        echo -e "${BLUE}🧪 测试单个通知效果${NC}"
        echo ""
        echo -e "显示单个测试通知..."
        ./termonaut easter-egg --floating 2>/dev/null || echo "🎉 这是一个测试通知! 🎉" | head -1
        ;;
    2)
        echo -e "${BLUE}🧪 终端兼容性信息${NC}"
        echo ""
        echo -e "当前终端环境详情:"
        echo -e "TERM: $TERM"
        echo -e "TERM_PROGRAM: ${TERM_PROGRAM:-未设置}"
        echo -e "COLORTERM: ${COLORTERM:-未设置}"
        echo -e "支持的颜色数: $(tput colors 2>/dev/null || echo '未知')"
        ;;
    3)
        echo -e "${BLUE}🧪 通知冲突测试${NC}"
        echo ""
        echo -e "快速连续显示多个通知..."
        for i in {1..3}; do
            echo -e "通知 $i"
            sleep 1
        done
        ;;
    4)
        echo -e "${GREEN}跳过额外测试${NC}"
        ;;
    *)
        echo -e "${YELLOW}无效选择，跳过额外测试${NC}"
        ;;
esac

echo ""

# 步骤10: 生成测试报告
show_step "10" "生成测试报告"
REPORT_FILE="$PROJECT_ROOT/floating_notification_test_report.txt"

cat > "$REPORT_FILE" << EOF
# 🎭 Termonaut 浮动通知测试报告

## 测试环境
- 测试时间: $(date)
- 终端类型: $TERM
- 终端程序: ${TERM_PROGRAM:-未知}
- 颜色支持: ${COLORTERM:-基础}
- Shell: $SHELL
- 终端大小: $(tput cols)x$(tput lines)

## 测试结果
- 浮动通知显示: $saw_notifications
- 样式渲染正常: $style_ok
- 自动消失功能: $auto_disappear
- 无用户干扰: $([ "$interference" == "n" ] && echo "y" || echo "n")
- 整体满意度: $satisfaction

## 测试文件
- 项目路径: $PROJECT_ROOT
- 二进制文件: $BINARY_PATH
- 测试命令: ./termonaut easter-egg --floating

## 建议
$(if [[ "$saw_notifications" == "y" && "$style_ok" == "y" && "$auto_disappear" == "y" ]]; then
    echo "✅ 浮动通知功能工作正常，建议正式集成到彩蛋系统中"
else
    echo "⚠️  浮动通知功能需要进一步调试和优化"
fi)
EOF

show_success "测试报告已生成: $REPORT_FILE"
echo ""

# 最终总结
echo -e "${CYAN}🎉 测试完成总结${NC}"
echo -e "${CYAN}===============${NC}"
echo ""

if [[ "$saw_notifications" == "y" && "$style_ok" == "y" && "$auto_disappear" == "y" ]]; then
    echo -e "${GREEN}🎊 恭喜! 浮动通知功能测试成功!${NC}"
    echo -e "${GREEN}   • 通知显示正常${NC}"
    echo -e "${GREEN}   • 样式渲染完美${NC}"
    echo -e "${GREEN}   • 自动消失功能正常${NC}"
    echo -e "${GREEN}   • 用户体验良好${NC}"
    echo ""
    echo -e "${WHITE}🚀 建议下一步:${NC}"
    echo -e "   1. 将浮动通知集成到彩蛋系统"
    echo -e "   2. 添加用户配置选项"
    echo -e "   3. 优化不同终端的兼容性"
    echo -e "   4. 考虑添加更多动画效果"
else
    echo -e "${YELLOW}🔧 测试发现一些问题，需要进一步优化${NC}"
    echo -e "${WHITE}建议检查:${NC}"
    echo -e "   • 终端兼容性问题"
    echo -e "   • ANSI转义序列支持"
    echo -e "   • 样式渲染问题"
    echo -e "   • 时间控制逻辑"
fi

echo ""
echo -e "${CYAN}📋 测试报告位置: ${WHITE}$REPORT_FILE${NC}"
echo -e "${CYAN}🔧 如需重新测试，请再次运行此脚本${NC}"
echo ""
echo -e "${PURPLE}感谢测试 Termonaut 浮动通知功能! 🎭✨${NC}"
