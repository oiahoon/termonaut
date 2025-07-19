#!/bin/bash

# Task 2.1 Part 3 Acceptance Test
# Tests for performance benchmark tests

set -e

echo "⚡ Task 2.1 Part 3: Performance Benchmark Tests"
echo "==============================================="

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
echo "📁 Benchmark Test Structure"
echo "==========================="

# Test 1: Check if benchmark test directory exists
run_test "Benchmark test directory exists" "test -d tests/benchmark"

# Test 2: Check if performance test exists
run_test "Performance benchmark test exists" "test -f tests/benchmark/performance_test.go"

# Test 3: Check benchmark test content
run_test "Benchmark test has performance tests" "grep -q 'BenchmarkCommandLogging\|BenchmarkStatsCalculation' tests/benchmark/performance_test.go"

echo ""
echo "🔧 Benchmark Test Quality"
echo "========================="

# Test 4: Check benchmark test line count
BENCHMARK_TEST_LINES=$(wc -l < tests/benchmark/performance_test.go 2>/dev/null || echo "0")
if [ "$BENCHMARK_TEST_LINES" -gt 300 ]; then
    echo -e "Testing: Benchmark test comprehensiveness... ${GREEN}PASS${NC} ($BENCHMARK_TEST_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Benchmark test comprehensiveness... ${RED}FAIL${NC} ($BENCHMARK_TEST_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 5: Check for command logging benchmark
run_test "Command logging benchmark implemented" "grep -q 'BenchmarkCommandLogging' tests/benchmark/performance_test.go"

# Test 6: Check for batch logging benchmark
run_test "Batch logging benchmark implemented" "grep -q 'BenchmarkBatchCommandLogging' tests/benchmark/performance_test.go"

# Test 7: Check for stats calculation benchmark
run_test "Stats calculation benchmark implemented" "grep -q 'BenchmarkStatsCalculation' tests/benchmark/performance_test.go"

# Test 8: Check for cache benchmark
run_test "Cache performance benchmark implemented" "grep -q 'BenchmarkStatsCalculationWithCache' tests/benchmark/performance_test.go"

# Test 9: Check for memory usage benchmark
run_test "Memory usage benchmark implemented" "grep -q 'BenchmarkMemoryUsage' tests/benchmark/performance_test.go"

echo ""
echo "📊 Benchmark Coverage Analysis"
echo "=============================="

# Test 10: Count benchmark functions
BENCHMARK_FUNCTIONS=$(grep -c "func Benchmark" tests/benchmark/performance_test.go 2>/dev/null || echo "0")
if [ "$BENCHMARK_FUNCTIONS" -ge 8 ]; then
    echo -e "Testing: Benchmark function count... ${GREEN}PASS${NC} ($BENCHMARK_FUNCTIONS functions)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Benchmark function count... ${RED}FAIL${NC} ($BENCHMARK_FUNCTIONS functions)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 11: Check for proper benchmark structure
run_test "Benchmarks use proper structure" "grep -q 'b.ResetTimer\|b.ReportAllocs' tests/benchmark/performance_test.go"

# Test 12: Check for comprehensive coverage
run_test "Benchmarks cover key operations" "grep -q 'BenchmarkCommandSanitization\|BenchmarkConfigOperations\|BenchmarkSessionManagement' tests/benchmark/performance_test.go"

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

if [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "\n${GREEN}✅ Task 2.1 Part 3 completed successfully!${NC}"
    echo "Performance benchmark tests have been implemented."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.1 Part 3 needs more work.${NC}"
    echo "Please improve benchmark test coverage."
    exit 1
fi
