#!/bin/bash

# Phase 1 Acceptance Test Script
# Tests for core functionality completion, database performance optimization, and error handling enhancement

set -e

echo "🚀 Starting Phase 1 Acceptance Tests"
echo "===================================="

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

# Helper function to run a test with output
run_test_with_output() {
    local test_name="$1"
    local test_command="$2"
    
    echo "Testing: $test_name"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

echo ""
echo "📋 Task 1.1: Critical TODO Items Completion"
echo "============================================"

# Test 1: Check if time tracking functions exist
run_test "Time tracking functions exist" "grep -q 'getEarlyBirdCommands\|getNightOwlCommands' internal/database/gamification.go"

# Test 2: Check if productivity calculation is implemented
run_test "Productivity calculation implemented" "grep -q 'calculateProductivityScore' cmd/termonaut/github_simple.go"

# Test 3: Check if achievements counting is implemented
run_test "Achievements counting implemented" "grep -q 'getAchievementsCount' cmd/termonaut/github_simple.go"

# Test 4: Check if avatar level fetching is improved
run_test "Avatar level fetching improved" "grep -q 'getUserLevel' internal/avatar/manager.go"

echo ""
echo "🗄️ Task 1.2: Database Performance Optimization"
echo "=============================================="

# Test 5: Check if enhanced connection pool settings exist
run_test "Enhanced connection pool settings" "grep -q 'SetMaxOpenConns(5)' internal/database/database.go"

# Test 6: Check if caching mechanism is implemented
run_test "Caching mechanism implemented" "grep -q 'getCachedResult\|setCachedResult' internal/database/database.go"

# Test 7: Check if batch operations are implemented
run_test "Batch operations implemented" "grep -q 'StoreCommandsBatch' internal/database/database.go"

# Test 8: Check if transaction support is added
run_test "Transaction support added" "grep -q 'WithTransaction' internal/database/database.go"

echo ""
echo "🛡️ Task 1.3: Error Handling Enhancement"
echo "======================================="

# Test 9: Check if enhanced network client exists
run_test "Enhanced network client exists" "test -f internal/network/client.go"

# Test 10: Check if retry mechanism is implemented
run_test "Retry mechanism implemented" "grep -q 'DoWithRetry' internal/network/client.go"

# Test 11: Check if circuit breaker is implemented
run_test "Circuit breaker implemented" "grep -q 'CircuitBreaker' internal/network/client.go"

# Test 12: Check if DiceBear client uses enhanced networking
run_test "DiceBear client uses enhanced networking" "grep -q 'network.EnhancedClient' internal/avatar/dicebear.go"

echo ""
echo "🧪 Unit Tests Execution"
echo "======================="

# Test 13: Run TODO fixes tests
if [ -f "tests/unit/todo_fixes_test.go" ]; then
    run_test_with_output "TODO fixes tests" "cd tests/unit && go test -v -run TestTimeTracking"
else
    echo -e "${YELLOW}⚠ TODO fixes tests not found${NC}"
fi

# Test 14: Run database performance tests
if [ -f "tests/unit/database_performance_test.go" ]; then
    run_test_with_output "Database performance tests" "cd tests/unit && go test -v -run TestCache"
else
    echo -e "${YELLOW}⚠ Database performance tests not found${NC}"
fi

# Test 15: Run error handling tests
if [ -f "tests/unit/error_handling_test.go" ]; then
    run_test_with_output "Error handling tests" "cd tests/unit && go test -v -run TestEnhancedClient"
else
    echo -e "${YELLOW}⚠ Error handling tests not found${NC}"
fi

echo ""
echo "📊 Performance Benchmarks"
echo "========================="

# Test 16: Run performance benchmarks if available
if [ -f "tests/unit/database_performance_test.go" ]; then
    echo "Running database performance benchmarks..."
    cd tests/unit && go test -bench=BenchmarkCommandLogging -benchmem || true
    cd ../..
else
    echo -e "${YELLOW}⚠ Performance benchmarks not available${NC}"
fi

echo ""
echo "🔍 Code Quality Checks"
echo "======================"

# Test 17: Check for remaining TODO items in critical files
TODO_COUNT=$(grep -r "TODO" internal/ cmd/ --include="*.go" | grep -v "test" | wc -l || echo "0")
if [ "$TODO_COUNT" -lt 5 ]; then
    echo -e "${GREEN}✓ TODO items reduced (${TODO_COUNT} remaining)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Too many TODO items remaining (${TODO_COUNT})${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 18: Check for proper error handling patterns
ERROR_HANDLING_COUNT=$(grep -r "fmt.Errorf.*%w" internal/ cmd/ --include="*.go" | wc -l || echo "0")
if [ "$ERROR_HANDLING_COUNT" -gt 10 ]; then
    echo -e "${GREEN}✓ Good error handling patterns found (${ERROR_HANDLING_COUNT} instances)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Insufficient error handling patterns (${ERROR_HANDLING_COUNT} instances)${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All Phase 1 acceptance tests passed!${NC}"
    echo "Phase 1 is ready for completion."
    exit 0
else
    echo -e "\n${RED}❌ Some tests failed. Please review and fix issues before proceeding.${NC}"
    echo "Success rate: $(( TESTS_PASSED * 100 / TOTAL_TESTS ))%"
    exit 1
fi
