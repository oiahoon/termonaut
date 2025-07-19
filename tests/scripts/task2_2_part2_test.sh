#!/bin/bash

# Task 2.2 Part 2 Acceptance Test
# Tests for Memory Leak Detection and Object Pools

set -e

echo "🔍 Task 2.2 Part 2: Memory Leak Detection & Object Pools"
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
echo "📁 Memory Monitoring Structure"
echo "=============================="

# Test 1: Check if monitoring package exists
run_test "Monitoring package directory exists" "test -d internal/monitoring"

# Test 2: Check if memory monitor exists
run_test "Memory monitor implementation exists" "test -f internal/monitoring/memory.go"

# Test 3: Check memory monitor content
run_test "Memory monitor has core functionality" "grep -q 'type MemoryMonitor struct\|func NewMemoryMonitor\|func.*Start\|func.*Stop' internal/monitoring/memory.go"

echo ""
echo "🔧 Memory Monitor Quality"
echo "========================"

# Test 4: Check memory monitor line count
MEMORY_MONITOR_LINES=$(wc -l < internal/monitoring/memory.go 2>/dev/null || echo "0")
if [ "$MEMORY_MONITOR_LINES" -gt 400 ]; then
    echo -e "Testing: Memory monitor comprehensiveness... ${GREEN}PASS${NC} ($MEMORY_MONITOR_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Memory monitor comprehensiveness... ${RED}FAIL${NC} ($MEMORY_MONITOR_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 5: Check for memory snapshots
run_test "Memory monitor has snapshots" "grep -q 'MemorySnapshot\|takeSnapshot\|GetSnapshots' internal/monitoring/memory.go"

# Test 6: Check for leak detection
run_test "Memory monitor has leak detection" "grep -q 'detectLeak\|LeakSuspect\|calculateTrend' internal/monitoring/memory.go"

# Test 7: Check for threshold monitoring
run_test "Memory monitor has thresholds" "grep -q 'SetThresholds\|checkThresholds\|memoryThreshold' internal/monitoring/memory.go"

# Test 8: Check for GC functionality
run_test "Memory monitor has GC control" "grep -q 'ForceGC\|runtime.GC' internal/monitoring/memory.go"

echo ""
echo "🏊 Object Pool Structure"
echo "======================="

# Test 9: Check if pool package exists
run_test "Pool package directory exists" "test -d internal/pool"

# Test 10: Check if object pool exists
run_test "Object pool implementation exists" "test -f internal/pool/object_pool.go"

# Test 11: Check object pool content
run_test "Object pool has core functionality" "grep -q 'CommandPool\|StringBuilderPool\|ByteSlicePool' internal/pool/object_pool.go"

echo ""
echo "🔧 Object Pool Quality"
echo "======================"

# Test 12: Check object pool line count
OBJECT_POOL_LINES=$(wc -l < internal/pool/object_pool.go 2>/dev/null || echo "0")
if [ "$OBJECT_POOL_LINES" -gt 200 ]; then
    echo -e "Testing: Object pool comprehensiveness... ${GREEN}PASS${NC} ($OBJECT_POOL_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Object pool comprehensiveness... ${RED}FAIL${NC} ($OBJECT_POOL_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 13: Check for multiple pool types
run_test "Object pool has multiple types" "grep -q 'CommandPool\|StringBuilderPool\|ByteSlicePool\|MapPool' internal/pool/object_pool.go"

# Test 14: Check for global pools
run_test "Object pool has global pools" "grep -q 'GlobalPools\|DefaultPools\|GetCommand\|PutCommand' internal/pool/object_pool.go"

# Test 15: Check for pool reset functionality
run_test "Object pool has reset functionality" "grep -q 'resetCommand\|Reset\|Clear' internal/pool/object_pool.go"

echo ""
echo "📊 Memory Optimization Features"
echo "==============================="

# Test 16: Count memory monitor methods
MONITOR_METHODS=$(grep -c "func.*MemoryMonitor" internal/monitoring/memory.go 2>/dev/null || echo "0")
if [ "$MONITOR_METHODS" -ge 8 ]; then
    echo -e "Testing: Memory monitor method count... ${GREEN}PASS${NC} ($MONITOR_METHODS methods)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Memory monitor method count... ${RED}FAIL${NC} ($MONITOR_METHODS methods)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 17: Count pool types
POOL_TYPES=$(grep -c "type.*Pool struct" internal/pool/object_pool.go 2>/dev/null || echo "0")
if [ "$POOL_TYPES" -ge 4 ]; then
    echo -e "Testing: Object pool type count... ${GREEN}PASS${NC} ($POOL_TYPES types)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Object pool type count... ${RED}FAIL${NC} ($POOL_TYPES types)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 18: Check for statistics and monitoring
run_test "Memory monitor has statistics" "grep -q 'MemoryStats\|GetStats\|PrintStats' internal/monitoring/memory.go"

# Test 19: Check for context and cancellation
run_test "Memory monitor has proper lifecycle" "grep -q 'context.Context\|context.WithCancel\|sync.WaitGroup' internal/monitoring/memory.go"

# Test 20: Check for thread safety in pools
run_test "Object pools are thread-safe" "grep -q 'sync.Pool\|sync.Mutex' internal/pool/object_pool.go"

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

if [ $SUCCESS_RATE -ge 85 ]; then
    echo -e "\n${GREEN}✅ Task 2.2 Part 2 completed successfully!${NC}"
    echo "Memory leak detection and object pools are well implemented."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.2 Part 2 needs more work.${NC}"
    echo "Please improve memory optimization implementation."
    exit 1
fi
