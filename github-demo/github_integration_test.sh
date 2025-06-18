#!/bin/bash

# Termonaut GitHub Integration Complete Test
# 测试所有GitHub相关功能

echo "🚀 Termonaut GitHub Integration Complete Test"
echo "============================================="
echo

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0

# 测试函数
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_pattern="$3"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${BLUE}Testing: ${test_name}${NC}"
    echo "Command: $command"

    # 运行命令并捕获输出
    output=$(eval "$command" 2>&1)
    exit_code=$?

    # 检查结果
    if [ $exit_code -eq 0 ] && [[ "$output" =~ $expected_pattern ]]; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC}"
        echo "Expected pattern: $expected_pattern"
        echo "Actual output: $output"
        echo "Exit code: $exit_code"
    fi
    echo
}

# 确保在正确的目录
cd "$(dirname "$0")/.."

echo "📊 1. Badge Generation Tests"
echo "─────────────────────────────"

# 测试基本badge生成
run_test "Basic Badge Generation" \
    "tn github badges generate" \
    "Commands.*Streak.*Productivity"

# 测试JSON格式输出
run_test "JSON Badge Export" \
    "tn github badges generate --format json" \
    '"Commands".*"Streak".*"Productivity"'

# 测试文件输出
mkdir -p test-output
run_test "Badge File Output" \
    "tn github badges generate --format json --output test-output/badges.json" \
    ""

# 验证文件是否创建
if [ -f "test-output/badges.json" ]; then
    echo -e "${GREEN}✅ Badge file created successfully${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ Badge file creation failed${NC}"
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo
echo "📝 2. Profile Generation Tests"
echo "──────────────────────────────"

# 测试基本profile生成
run_test "Basic Profile Generation" \
    "tn github profile generate" \
    "My Termonaut Profile.*Stats.*Overview"

# 测试markdown格式输出
run_test "Markdown Profile Export" \
    "tn github profile generate --format markdown" \
    "# 🚀 My Termonaut Profile"

# 测试profile文件输出
run_test "Profile File Output" \
    "tn github profile generate --format markdown --output test-output/profile.md" \
    ""

# 验证profile文件
if [ -f "test-output/profile.md" ]; then
    echo -e "${GREEN}✅ Profile file created successfully${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}❌ Profile file creation failed${NC}"
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo
echo "🔥 3. Heatmap Tests"
echo "──────────────────"

# 测试heatmap生成
run_test "Heatmap Generation" \
    "tn heatmap" \
    "Weekly Productivity Heatmap.*Peak Hours"

# 测试JSON heatmap
run_test "JSON Heatmap Export" \
    "tn heatmap --json" \
    '"optimal_hours".*"peak_day"'

echo
echo "🤖 4. GitHub Actions Tests"
echo "──────────────────────────"

# 测试actions help
run_test "GitHub Actions Help" \
    "tn github actions --help" \
    "GitHub Actions integration.*workflow"

# 测试actions list (应该提示配置)
run_test "GitHub Actions List" \
    "tn github actions list" \
    "Configure GitHub repository first"

echo
echo "⚙️ 5. Configuration Tests"
echo "─────────────────────────"

# 测试配置查看
run_test "Configuration Display" \
    "tn config get" \
    "Termonaut Configuration.*Display Mode.*Theme"

# 测试主题设置
run_test "Theme Configuration" \
    "tn config set theme emoji" \
    "Configuration updated.*theme.*emoji"

# 恢复默认主题
tn config set theme minimal > /dev/null 2>&1

echo
echo "📈 6. Stats Integration Tests"
echo "────────────────────────────"

# 测试基本统计
run_test "Basic Stats" \
    "tn stats" \
    "Termonaut Stats.*Total Commands.*Level"

# 测试版本信息
run_test "Version Information" \
    "tn version" \
    "Termonaut.*Commit.*Built"

echo
echo "🔧 7. File Validation Tests"
echo "──────────────────────────"

# 验证生成的文件内容
echo "Validating generated files..."

if [ -f "test-output/badges.json" ]; then
    # 检查JSON格式
    if python3 -m json.tool test-output/badges.json > /dev/null 2>&1; then
        echo -e "${GREEN}✅ badges.json is valid JSON${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ badges.json is invalid JSON${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    # 检查必需的badge类型
    required_badges=("Commands" "Streak" "Productivity" "XP")
    for badge in "${required_badges[@]}"; do
        if grep -q "\"$badge\"" test-output/badges.json; then
            echo -e "${GREEN}✅ $badge badge found${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}❌ $badge badge missing${NC}"
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    done
fi

if [ -f "test-output/profile.md" ]; then
    # 检查profile内容
    required_sections=("# 🚀 My Termonaut Profile" "## 📊 Stats" "## 📈 Overview" "## 🏆 Achievements")
    for section in "${required_sections[@]}"; do
        if grep -q "$section" test-output/profile.md; then
            echo -e "${GREEN}✅ Profile section found: $section${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}❌ Profile section missing: $section${NC}"
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    done
fi

echo
echo "🎨 8. Badge URL Validation"
echo "─────────────────────────"

# 验证badge URL格式
if [ -f "test-output/badges.json" ]; then
    echo "Validating badge URLs..."

    # 检查URL格式
    if grep -q "https://img.shields.io/badge/" test-output/badges.json; then
        echo -e "${GREEN}✅ Badge URLs use shields.io format${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Invalid badge URL format${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    # 检查样式参数
    if grep -q "style=flat-square" test-output/badges.json; then
        echo -e "${GREEN}✅ Badge style parameter found${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Badge style parameter missing${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    # 检查logo参数
    if grep -q "logo=terminal" test-output/badges.json; then
        echo -e "${GREEN}✅ Terminal logo parameter found${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Terminal logo parameter missing${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
fi

echo
echo "📱 9. Social Integration Tests"
echo "─────────────────────────────"

# 测试profile的社交媒体友好性
if [ -f "test-output/profile.md" ]; then
    # 检查emoji使用
    if grep -q "🚀\|📊\|🏆\|🔥" test-output/profile.md; then
        echo -e "${GREEN}✅ Profile contains emojis${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Profile lacks visual elements${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    # 检查链接格式
    if grep -q "\[.*\](https://.*)" test-output/profile.md; then
        echo -e "${GREEN}✅ Profile contains proper markdown links${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ Profile lacks proper links${NC}"
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
fi

echo
echo "🧹 10. Cleanup and Summary"
echo "─────────────────────────"

# 显示生成的文件
echo "Generated files:"
ls -la test-output/ 2>/dev/null || echo "No files generated"

# 清理测试文件
echo "Cleaning up test files..."
rm -rf test-output/

echo
echo "📋 Test Summary"
echo "═══════════════"
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$((TOTAL_TESTS - PASSED_TESTS))${NC}"

if [ $PASSED_TESTS -eq $TOTAL_TESTS ]; then
    echo -e "${GREEN}🎉 All tests passed! GitHub integration is working perfectly.${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️ Some tests failed. Please check the output above.${NC}"
    success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "Success rate: ${BLUE}$success_rate%${NC}"
    exit 1
fi