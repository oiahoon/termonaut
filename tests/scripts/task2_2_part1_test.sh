#!/bin/bash

# Task 2.2 Part 1 Acceptance Test
# Tests for LRU Cache Implementation

set -e

echo "🧠 Task 2.2 Part 1: LRU Cache Implementation"
echo "============================================"

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
echo "📁 LRU Cache Structure"
echo "====================="

# Test 1: Check if cache package exists
run_test "Cache package directory exists" "test -d internal/cache"

# Test 2: Check if LRU cache file exists
run_test "LRU cache implementation exists" "test -f internal/cache/lru.go"

# Test 3: Check LRU cache content
run_test "LRU cache has core functionality" "grep -q 'type LRUCache struct\|func NewLRUCache\|func.*Get\|func.*Set' internal/cache/lru.go"

echo ""
echo "🔧 LRU Cache Quality"
echo "==================="

# Test 4: Check LRU cache line count
LRU_CACHE_LINES=$(wc -l < internal/cache/lru.go 2>/dev/null || echo "0")
if [ "$LRU_CACHE_LINES" -gt 400 ]; then
    echo -e "Testing: LRU cache comprehensiveness... ${GREEN}PASS${NC} ($LRU_CACHE_LINES lines)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: LRU cache comprehensiveness... ${RED}FAIL${NC} ($LRU_CACHE_LINES lines)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 5: Check for thread safety
run_test "LRU cache has thread safety" "grep -q 'sync.RWMutex\|mutex.Lock\|mutex.RLock' internal/cache/lru.go"

# Test 6: Check for TTL support
run_test "LRU cache supports TTL" "grep -q 'TTL\|ExpiresAt\|time.Duration' internal/cache/lru.go"

# Test 7: Check for statistics
run_test "LRU cache has statistics" "grep -q 'Stats\|hits\|misses\|HitRate' internal/cache/lru.go"

# Test 8: Check for cleanup functionality
run_test "LRU cache has cleanup" "grep -q 'CleanupExpired\|removeElement' internal/cache/lru.go"

echo ""
echo "🔗 Database Integration"
echo "======================"

# Test 9: Check if database uses cache
run_test "Database imports cache package" "grep -q 'github.com/oiahoon/termonaut/internal/cache' internal/database/database.go"

# Test 10: Check if database has LRU cache field
run_test "Database has LRU cache field" "grep -q 'lruCache.*LRUCache\|lruCache.*cache.LRUCache' internal/database/database.go"

# Test 11: Check for cache methods in database
run_test "Database has cache methods" "grep -q 'getCachedResult\|setCachedResult\|GetCacheStats' internal/database/database.go"

# Test 12: Check for batch operations
run_test "Database has batch operations" "grep -q 'StoreCommandsBatch\|WithTransaction' internal/database/database.go"

echo ""
echo "📊 Cache Features Analysis"
echo "=========================="

# Test 13: Count cache methods
CACHE_METHODS=$(grep -c "func.*LRUCache" internal/cache/lru.go 2>/dev/null || echo "0")
if [ "$CACHE_METHODS" -ge 10 ]; then
    echo -e "Testing: Cache method count... ${GREEN}PASS${NC} ($CACHE_METHODS methods)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Testing: Cache method count... ${RED}FAIL${NC} ($CACHE_METHODS methods)"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Test 14: Check for memory estimation
run_test "Cache has memory estimation" "grep -q 'EstimateMemoryUsage\|estimateSize' internal/cache/lru.go"

# Test 15: Check for TTL cache wrapper
run_test "TTL cache wrapper exists" "grep -q 'type TTLCache struct\|NewTTLCache' internal/cache/lru.go"

echo ""
echo "📈 Results Summary"
echo "=================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

SUCCESS_RATE=$(( TESTS_PASSED * 100 / TOTAL_TESTS ))
echo "Success Rate: ${SUCCESS_RATE}%"

if [ $SUCCESS_RATE -ge 85 ]; then
    echo -e "\n${GREEN}✅ Task 2.2 Part 1 completed successfully!${NC}"
    echo "LRU Cache implementation is comprehensive and well-integrated."
    exit 0
else
    echo -e "\n${RED}❌ Task 2.2 Part 1 needs more work.${NC}"
    echo "Please improve LRU cache implementation."
    exit 1
fi
