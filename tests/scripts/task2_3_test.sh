#!/bin/bash

# Task 2.3 Acceptance Test
# Tests for Code Refactoring

set -e

echo "🔧 Task 2.3: Code Refactoring"
echo "============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=0

# Helper function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -n "Testing: $test_name... "
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if eval "$test_command" >/dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}FAIL${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

echo ""
echo "📁 Refactored Component Structure"
echo "================================="

# Test 1: Check if components package exists
run_test "Components package directory exists" "test -d internal/tui/enhanced/components"

# Test 2: Check if home tab component exists
run_test "Home tab component exists" "test -f internal/tui/enhanced/components/home_tab.go"

# Test 3: Check if theme component exists
run_test "Theme component exists" "test -f internal/tui/enhanced/components/theme.go"

# Test 4: Check if utils package exists
run_test "Utils package directory exists" "test -d internal/utils"

# Test 5: Check if common utils exist
run_test "Common utils implementation exists" "test -f internal/utils/common.go"

echo ""
echo "🔧 Component Quality"
echo "==================="

# Test 6: Check home tab component content
run_test "Home tab has proper structure" "grep -q 'type HomeTabComponent struct\|func NewHomeTabComponent\|func.*Render' internal/tui/enhanced/components/home_tab.go"

# Test 7: Check theme component content
run_test "Theme has multiple themes" "grep -q 'NewSpaceTheme\|NewCyberpunkTheme\|NewMinimalTheme' internal/tui/enhanced/components/theme.go"

# Test 8: Check utils content
run_test "Utils has utility functions" "grep -q 'StringUtils\|TimeUtils\|NumberUtils\|ProgressUtils' internal/utils/common.go"

echo ""
echo "📊 Code Quality Metrics"
echo "======================="

# Test 9: Check home tab component line count
HOME_TAB_LINES=$(wc -l < internal/tui/enhanced/components/home_tab.go 2>/dev/null || echo "0")
if [ "$HOME_TAB_LINES" -gt 200 ] && [ "$HOME_TAB_LINES" -lt 400 ]; then
    echo -e "Testing: Home tab component size... ${GREEN}PASS${NC} ($HOME_TAB_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Home tab component size... ${RED}FAIL${NC} ($HOME_TAB_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 10: Check theme component line count
THEME_LINES=$(wc -l < internal/tui/enhanced/components/theme.go 2>/dev/null || echo "0")
if [ "$THEME_LINES" -gt 200 ]; then
    echo -e "Testing: Theme component comprehensiveness... ${GREEN}PASS${NC} ($THEME_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Theme component comprehensiveness... ${RED}FAIL${NC} ($THEME_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 11: Check utils line count
UTILS_LINES=$(wc -l < internal/utils/common.go 2>/dev/null || echo "0")
if [ "$UTILS_LINES" -gt 300 ]; then
    echo -e "Testing: Utils comprehensiveness... ${GREEN}PASS${NC} ($UTILS_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Utils comprehensiveness... ${RED}FAIL${NC} ($UTILS_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "🎨 Theme System"
echo "==============="

# Test 12: Check for multiple theme implementations
THEME_COUNT=$(grep -c "func New.*Theme" internal/tui/enhanced/components/theme.go 2>/dev/null || echo "0")
if [ "$THEME_COUNT" -ge 3 ]; then
    echo -e "Testing: Multiple themes available... ${GREEN}PASS${NC} ($THEME_COUNT themes)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Multiple themes available... ${RED}FAIL${NC} ($THEME_COUNT themes)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 13: Check for theme color definitions
run_test "Themes have color definitions" "grep -q 'Primary.*lipgloss.Color\|Secondary.*lipgloss.Color' internal/tui/enhanced/components/theme.go"

# Test 14: Check for component styles
run_test "Themes have component styles" "grep -q 'ContentBox\|SectionBox\|StatBox' internal/tui/enhanced/components/theme.go"

echo ""
echo "🛠️ Utility Functions"
echo "==================="

# Test 15: Count utility types
UTIL_TYPES=$(grep -c "type.*Utils struct" internal/utils/common.go 2>/dev/null || echo "0")
if [ "$UTIL_TYPES" -ge 5 ]; then
    echo -e "Testing: Utility type count... ${GREEN}PASS${NC} ($UTIL_TYPES types)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Utility type count... ${RED}FAIL${NC} ($UTIL_TYPES types)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 16: Check for string utilities
run_test "String utilities implemented" "grep -q 'TruncateString\|PadString\|CenterString\|WrapText' internal/utils/common.go"

# Test 17: Check for time utilities
run_test "Time utilities implemented" "grep -q 'FormatDuration\|FormatRelativeTime' internal/utils/common.go"

# Test 18: Check for number utilities
run_test "Number utilities implemented" "grep -q 'FormatNumber\|FormatBytes\|FormatPercentage' internal/utils/common.go"

# Test 19: Check for progress utilities
run_test "Progress utilities implemented" "grep -q 'CreateProgressBar\|CreateProgressBarWithPercentage' internal/utils/common.go"

echo ""
echo "🧹 Code Organization"
echo "==================="

# Test 20: Check for proper package structure
run_test "Components have proper package declaration" "grep -q '^package components$' internal/tui/enhanced/components/home_tab.go"

# Test 21: Check for proper imports
run_test "Components have proper imports" "grep -q 'github.com/charmbracelet/lipgloss' internal/tui/enhanced/components/theme.go"

# Test 22: Check for global utility instances
run_test "Utils have global instances" "grep -q 'String.*=.*StringUtils\|Time.*=.*TimeUtils' internal/utils/common.go"

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

# Additional metrics
echo ""
echo "📊 Code Metrics"
echo "==============="

# Count total refactored lines
TOTAL_REFACTORED_LINES=$((HOME_TAB_LINES + THEME_LINES + UTILS_LINES))
echo "Total refactored code: $TOTAL_REFACTORED_LINES lines"

# Check for function count
TOTAL_FUNCTIONS=$(grep -r "func " internal/tui/enhanced/components/ internal/utils/ 2>/dev/null | wc -l || echo "0")
echo "Total functions: $TOTAL_FUNCTIONS"

# Check for struct count
TOTAL_STRUCTS=$(grep -r "type.*struct" internal/tui/enhanced/components/ internal/utils/ 2>/dev/null | wc -l || echo "0")
echo "Total structs: $TOTAL_STRUCTS"

if [ $SUCCESS_RATE -ge 85 ]; then
    echo -e "\n${GREEN}✅ Task 2.3 completed successfully!${NC}"
    echo "Code refactoring has improved structure and maintainability."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.3 needs more work.${NC}"
    echo "Please improve code refactoring implementation."
    exit 1
fi
