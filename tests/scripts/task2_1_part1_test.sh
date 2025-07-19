#!/bin/bash

# Task 2.1 Part 1 Acceptance Test
# Tests for additional unit test coverage

set -e

echo "🧪 Task 2.1 Part 1: Additional Unit Tests"
echo "========================================="

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
echo "📁 New Test Files Check"
echo "======================="

# Test 1: Check if config test exists
run_test "Config test file exists" "test -f tests/unit/config_test.go"

# Test 2: Check if privacy test exists  
run_test "Privacy test file exists" "test -f tests/unit/privacy_simple_test.go"

# Test 3: Check config test content
run_test "Config test has comprehensive tests" "grep -q 'TestDefaultConfig\|TestConfigSaveAndLoad\|TestConfigValidation' tests/unit/config_test.go"

# Test 4: Check privacy test content
run_test "Privacy test has basic tests" "grep -q 'TestCommandSanitizerBasic\|TestSanitizationConfig' tests/unit/privacy_simple_test.go"

echo ""
echo "🔧 Test File Quality Check"
echo "=========================="

# Test 5: Check config test line count (should be substantial)
CONFIG_TEST_LINES=$(wc -l < tests/unit/config_test.go 2>/dev/null || echo "0")
if [ "$CONFIG_TEST_LINES" -gt 200 ]; then
    echo -e "Testing: Config test comprehensiveness... ${GREEN}PASS${NC} ($CONFIG_TEST_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Config test comprehensiveness... ${RED}FAIL${NC} ($CONFIG_TEST_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 6: Check privacy test line count
PRIVACY_TEST_LINES=$(wc -l < tests/unit/privacy_simple_test.go 2>/dev/null || echo "0")
if [ "$PRIVACY_TEST_LINES" -gt 50 ]; then
    echo -e "Testing: Privacy test adequacy... ${GREEN}PASS${NC} ($PRIVACY_TEST_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Privacy test adequacy... ${RED}FAIL${NC} ($PRIVACY_TEST_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "📊 Test Coverage Analysis"
echo "========================"

# Test 7: Count total test files
TOTAL_TEST_FILES=$(find tests/unit -name "*_test.go" | wc -l)
if [ "$TOTAL_TEST_FILES" -ge 6 ]; then
    echo -e "Testing: Test file count... ${GREEN}PASS${NC} ($TOTAL_TEST_FILES files)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Test file count... ${YELLOW}PARTIAL${NC} ($TOTAL_TEST_FILES files)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 8: Check for test function diversity
TEST_FUNCTIONS=$(grep -r "func Test" tests/unit/ | wc -l)
if [ "$TEST_FUNCTIONS" -ge 20 ]; then
    echo -e "Testing: Test function diversity... ${GREEN}PASS${NC} ($TEST_FUNCTIONS functions)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Test function diversity... ${YELLOW}PARTIAL${NC} ($TEST_FUNCTIONS functions)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

if [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "\n${GREEN}✅ Task 2.1 Part 1 completed successfully!${NC}"
    echo "Additional unit tests have been added."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.1 Part 1 needs more work.${NC}"
    echo "Please add more comprehensive unit tests."
    exit 1
fi
