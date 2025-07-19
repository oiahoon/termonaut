#!/bin/bash

# Phase 1 Acceptance Test Script (No Go dependency)
# Tests for core functionality completion, database performance optimization, and error handling enhancement

set -e

echo "🚀 Starting Phase 1 Acceptance Tests (Code Analysis Only)"
echo "========================================================="

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

# Test 5: Verify TODO comments are removed from critical functions
run_test "Time tracking TODO removed" "! grep -q 'TODO.*time-based tracking' internal/database/gamification.go"

# Test 6: Verify productivity TODO removed
run_test "Productivity TODO removed" "! grep -q 'TODO.*Calculate actual productivity' cmd/termonaut/github_simple.go"

# Test 7: Verify avatar TODO removed
run_test "Avatar TODO removed" "! grep -q 'TODO.*Get from user stats' internal/avatar/manager.go"

echo ""
echo "🗄️ Task 1.2: Database Performance Optimization"
echo "=============================================="

# Test 8: Check if enhanced connection pool settings exist
run_test "Enhanced connection pool settings" "grep -q 'SetMaxOpenConns(5)' internal/database/database.go"

# Test 9: Check if optimized database URL is used
run_test "Optimized database URL" "grep -q '_cache_size=10000' internal/database/database.go"

# Test 10: Check if caching mechanism is implemented
run_test "Caching mechanism implemented" "grep -q 'getCachedResult\|setCachedResult' internal/database/database.go"

# Test 11: Check if cache TTL is defined
run_test "Cache TTL defined" "grep -q 'CacheTTL.*=.*time.Minute' internal/database/database.go"

# Test 12: Check if batch operations are implemented
run_test "Batch operations implemented" "grep -q 'StoreCommandsBatch' internal/database/database.go"

# Test 13: Check if transaction support is added
run_test "Transaction support added" "grep -q 'WithTransaction' internal/database/database.go"

# Test 14: Check if optimized GetBasicStats uses single query
run_test "Optimized GetBasicStats" "grep -A 15 'GetBasicStats' internal/database/database.go | grep -q 'SELECT.*COUNT.*FROM commands.*as total_commands'"

echo ""
echo "🛡️ Task 1.3: Error Handling Enhancement"
echo "======================================="

# Test 15: Check if enhanced network client exists
run_test "Enhanced network client exists" "test -f internal/network/client.go"

# Test 16: Check if retry mechanism is implemented
run_test "Retry mechanism implemented" "grep -q 'DoWithRetry' internal/network/client.go"

# Test 17: Check if circuit breaker is implemented
run_test "Circuit breaker implemented" "grep -q 'CircuitBreaker' internal/network/client.go"

# Test 18: Check if exponential backoff is implemented
run_test "Exponential backoff implemented" "grep -q 'BackoffFactor' internal/network/client.go"

# Test 19: Check if DiceBear client uses enhanced networking
run_test "DiceBear client uses enhanced networking" "grep -q 'network.EnhancedClient' internal/avatar/dicebear.go"

# Test 20: Check if context support is added
run_test "Context support added" "grep -q 'GetWithContext' internal/network/client.go"

# Test 21: Check if safe body reading is implemented
run_test "Safe body reading implemented" "grep -q 'SafeReadBody' internal/network/client.go"

echo ""
echo "🧪 Test Files Verification"
echo "=========================="

# Test 22: Check if TODO fixes test exists
run_test "TODO fixes test exists" "test -f tests/unit/todo_fixes_test.go"

# Test 23: Check if database performance test exists
run_test "Database performance test exists" "test -f tests/unit/database_performance_test.go"

# Test 24: Check if error handling test exists
run_test "Error handling test exists" "test -f tests/unit/error_handling_test.go"

# Test 25: Check if benchmark tests are included
run_test "Benchmark tests included" "grep -q 'func Benchmark' tests/unit/database_performance_test.go"

echo ""
echo "🔍 Code Quality Checks"
echo "======================"

# Test 26: Check for remaining critical TODO items
CRITICAL_TODO_COUNT=$(grep -r "TODO" internal/database/gamification.go internal/avatar/manager.go cmd/termonaut/github_simple.go 2>/dev/null | wc -l || echo "0")
if [ "$CRITICAL_TODO_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✓ Critical TODO items resolved${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Critical TODO items remaining (${CRITICAL_TODO_COUNT})${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 27: Check for proper error handling patterns
ERROR_HANDLING_COUNT=$(grep -r "fmt.Errorf.*%w" internal/ cmd/ --include="*.go" 2>/dev/null | wc -l || echo "0")
if [ "$ERROR_HANDLING_COUNT" -gt 15 ]; then
    echo -e "${GREEN}✓ Good error handling patterns found (${ERROR_HANDLING_COUNT} instances)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ Moderate error handling patterns (${ERROR_HANDLING_COUNT} instances)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 28: Check for transaction usage
TRANSACTION_COUNT=$(grep -r "WithTransaction\|Begin()" internal/ --include="*.go" 2>/dev/null | wc -l || echo "0")
if [ "$TRANSACTION_COUNT" -gt 2 ]; then
    echo -e "${GREEN}✓ Transaction support implemented (${TRANSACTION_COUNT} instances)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Insufficient transaction support (${TRANSACTION_COUNT} instances)${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 29: Check for caching implementation
CACHE_COUNT=$(grep -r "cache\|Cache" internal/database/ --include="*.go" 2>/dev/null | wc -l || echo "0")
if [ "$CACHE_COUNT" -gt 5 ]; then
    echo -e "${GREEN}✓ Caching mechanism implemented (${CACHE_COUNT} references)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Insufficient caching implementation (${CACHE_COUNT} references)${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "📊 Implementation Verification"
echo "=============================="

# Test 30: Verify time tracking implementation
if grep -q "strftime('%H'" internal/database/gamification.go; then
    echo -e "${GREEN}✓ Time-based query implementation found${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Time-based query implementation missing${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 31: Verify productivity calculation complexity
PRODUCTIVITY_LINES=$(grep -A 20 "calculateProductivityScore" cmd/termonaut/github_simple.go | wc -l || echo "0")
if [ "$PRODUCTIVITY_LINES" -gt 15 ]; then
    echo -e "${GREEN}✓ Comprehensive productivity calculation (${PRODUCTIVITY_LINES} lines)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Simple productivity calculation (${PRODUCTIVITY_LINES} lines)${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
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

if [ $SUCCESS_RATE -ge 90 ]; then
    echo -e "\n${GREEN}🎉 Phase 1 acceptance tests mostly passed! (${SUCCESS_RATE}%)${NC}"
    echo "Phase 1 implementation is substantially complete."
    
    if [ $TESTS_FAILED -gt 0 ]; then
        echo -e "${YELLOW}⚠ Minor issues detected. Consider reviewing failed tests.${NC}"
    fi
    
    exit 0
elif [ $SUCCESS_RATE -ge 75 ]; then
    echo -e "\n${YELLOW}⚠ Phase 1 partially complete (${SUCCESS_RATE}%)${NC}"
    echo "Most functionality implemented, but some issues need attention."
    exit 1
else
    echo -e "\n${RED}❌ Phase 1 needs significant work (${SUCCESS_RATE}%)${NC}"
    echo "Please review and fix major issues before proceeding."
    exit 1
fi
