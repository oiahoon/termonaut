#!/bin/bash

# Task 2.1 Part 2 Acceptance Test
# Tests for integration tests

set -e

echo "🔗 Task 2.1 Part 2: Integration Tests"
echo "====================================="

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
echo "📁 Integration Test Structure"
echo "============================"

# Test 1: Check if integration test directory exists
run_test "Integration test directory exists" "test -d tests/integration"

# Test 2: Check if workflow test exists
run_test "Workflow integration test exists" "test -f tests/integration/workflow_test.go"

# Test 3: Check integration test content
run_test "Integration test has comprehensive tests" "grep -q 'TestFullWorkflow\|TestDatabaseIntegration\|TestConfigIntegration' tests/integration/workflow_test.go"

echo ""
echo "🔧 Integration Test Quality"
echo "=========================="

# Test 4: Check integration test line count
INTEGRATION_TEST_LINES=$(wc -l < tests/integration/workflow_test.go 2>/dev/null || echo "0")
if [ "$INTEGRATION_TEST_LINES" -gt 200 ]; then
    echo -e "Testing: Integration test comprehensiveness... ${GREEN}PASS${NC} ($INTEGRATION_TEST_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Integration test comprehensiveness... ${RED}FAIL${NC} ($INTEGRATION_TEST_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 5: Check for workflow testing
run_test "Full workflow test implemented" "grep -q 'TestFullWorkflow' tests/integration/workflow_test.go"

# Test 6: Check for database integration testing
run_test "Database integration test implemented" "grep -q 'TestDatabaseIntegration' tests/integration/workflow_test.go"

# Test 7: Check for config integration testing
run_test "Config integration test implemented" "grep -q 'TestConfigIntegration' tests/integration/workflow_test.go"

echo ""
echo "📊 Test Coverage Analysis"
echo "========================"

# Test 8: Count integration test functions
INTEGRATION_TEST_FUNCTIONS=$(grep -c "func Test" tests/integration/workflow_test.go 2>/dev/null || echo "0")
if [ "$INTEGRATION_TEST_FUNCTIONS" -ge 3 ]; then
    echo -e "Testing: Integration test function count... ${GREEN}PASS${NC} ($INTEGRATION_TEST_FUNCTIONS functions)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Integration test function count... ${RED}FAIL${NC} ($INTEGRATION_TEST_FUNCTIONS functions)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 9: Check for proper test structure
run_test "Integration tests use proper structure" "grep -q 'tempDir.*TempDir\|defer.*Close' tests/integration/workflow_test.go"

# Test 10: Check for comprehensive assertions
run_test "Integration tests have assertions" "grep -q 'if.*Error\|t.Errorf\|t.Fatalf' tests/integration/workflow_test.go"

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

if [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "\n${GREEN}✅ Task 2.1 Part 2 completed successfully!${NC}"
    echo "Integration tests have been implemented."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.1 Part 2 needs more work.${NC}"
    echo "Please improve integration test coverage."
    exit 1
fi
