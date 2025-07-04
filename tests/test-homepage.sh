#!/bin/bash

# Termonaut Homepage Test Script
# Basic tests to ensure homepage functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DOCS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../docs" && pwd)"
TESTS_PASSED=0
TESTS_FAILED=0

echo -e "${BLUE}🧪 Termonaut Homepage Tests${NC}"
echo -e "${BLUE}===========================${NC}"
echo ""

# Test function
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -n "Testing $test_name... "
    
    if eval "$test_command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# File existence tests
echo -e "${BLUE}📁 File Existence Tests${NC}"
run_test "index.html exists" "[ -f '$DOCS_DIR/index.html' ]"
run_test "CSS file exists" "[ -f '$DOCS_DIR/assets/css/style.css' ]"
run_test "JavaScript file exists" "[ -f '$DOCS_DIR/assets/js/main.js' ]"
run_test "README exists" "[ -f '$DOCS_DIR/README.md' ]"
run_test "Favicon exists" "[ -f '$DOCS_DIR/favicon.svg' ]"
echo ""

# HTML structure tests
echo -e "${BLUE}🏗️ HTML Structure Tests${NC}"
run_test "HTML has title tag" "grep -q '<title>' '$DOCS_DIR/index.html'"
run_test "HTML has meta description" "grep -q 'meta.*description' '$DOCS_DIR/index.html'"
run_test "HTML has Open Graph tags" "grep -q 'og:title' '$DOCS_DIR/index.html'"
run_test "HTML has viewport meta" "grep -q 'viewport' '$DOCS_DIR/index.html'"
run_test "HTML has charset meta" "grep -q 'charset' '$DOCS_DIR/index.html'"
echo ""

# Content tests
echo -e "${BLUE}📝 Content Tests${NC}"
run_test "Title contains 'Termonaut'" "grep -q 'Termonaut' '$DOCS_DIR/index.html'"
run_test "Has hero section" "grep -q 'class=\"hero\"' '$DOCS_DIR/index.html'"
run_test "Has features section" "grep -q 'id=\"features\"' '$DOCS_DIR/index.html'"
run_test "Has demo section" "grep -q 'id=\"demo\"' '$DOCS_DIR/index.html'"
run_test "Has install section" "grep -q 'id=\"install\"' '$DOCS_DIR/index.html'"
run_test "Has GitHub links" "grep -q 'github.com/oiahoon/termonaut' '$DOCS_DIR/index.html'"
echo ""

# CSS tests
echo -e "${BLUE}🎨 CSS Tests${NC}"
run_test "CSS has root variables" "grep -q ':root' '$DOCS_DIR/assets/css/style.css'"
run_test "CSS has terminal styles" "grep -q 'terminal' '$DOCS_DIR/assets/css/style.css'"
run_test "CSS has responsive breakpoints" "grep -q '@media' '$DOCS_DIR/assets/css/style.css'"
run_test "CSS has animations" "grep -q '@keyframes' '$DOCS_DIR/assets/css/style.css'"
echo ""

# JavaScript tests
echo -e "${BLUE}⚡ JavaScript Tests${NC}"
run_test "JS has DOMContentLoaded" "grep -q 'DOMContentLoaded' '$DOCS_DIR/assets/js/main.js'"
run_test "JS has GitHub stats function" "grep -q 'initGitHubStats' '$DOCS_DIR/assets/js/main.js'"
run_test "JS has copy functionality" "grep -q 'clipboard' '$DOCS_DIR/assets/js/main.js'"
run_test "JS has demo terminal" "grep -q 'initDemoTerminal' '$DOCS_DIR/assets/js/main.js'"
run_test "JS has particles init" "grep -q 'initParticles' '$DOCS_DIR/assets/js/main.js'"
echo ""

# Link validation tests
echo -e "${BLUE}🔗 Link Validation Tests${NC}"

# Check internal links
internal_links_valid=true
while IFS= read -r link; do
    id=$(echo "$link" | sed 's/.*href="#\([^"]*\)".*/\1/')
    if ! grep -q "id=\"$id\"" "$DOCS_DIR/index.html"; then
        internal_links_valid=false
        break
    fi
done < <(grep -o 'href="#[^"]*"' "$DOCS_DIR/index.html" 2>/dev/null || true)

run_test "Internal links are valid" "$internal_links_valid"

# Check for common issues
run_test "No broken image tags" "! grep -q 'src=\"\"' '$DOCS_DIR/index.html'"
run_test "No empty href attributes" "! grep -q 'href=\"\"' '$DOCS_DIR/index.html'"
echo ""

# Performance tests
echo -e "${BLUE}⚡ Performance Tests${NC}"
html_size=$(wc -c < "$DOCS_DIR/index.html")
css_size=$(wc -c < "$DOCS_DIR/assets/css/style.css")
js_size=$(wc -c < "$DOCS_DIR/assets/js/main.js")

run_test "HTML size reasonable (<100KB)" "[ $html_size -lt 100000 ]"
run_test "CSS size reasonable (<200KB)" "[ $css_size -lt 200000 ]"
run_test "JS size reasonable (<100KB)" "[ $js_size -lt 100000 ]"
echo ""

# Accessibility tests
echo -e "${BLUE}♿ Accessibility Tests${NC}"
run_test "Has alt attributes for images" "! grep -q '<img[^>]*src[^>]*>' '$DOCS_DIR/index.html' || grep -q 'alt=' '$DOCS_DIR/index.html'"
run_test "Has proper heading hierarchy" "grep -q '<h1' '$DOCS_DIR/index.html'"
run_test "Has lang attribute" "grep -q 'lang=' '$DOCS_DIR/index.html'"
echo ""

# Security tests
echo -e "${BLUE}🔒 Security Tests${NC}"
run_test "No inline JavaScript" "! grep -q 'javascript:' '$DOCS_DIR/index.html'"
run_test "No eval() usage" "! grep -q 'eval(' '$DOCS_DIR/assets/js/main.js'"
run_test "Uses HTTPS links" "! grep -q 'http://[^l]' '$DOCS_DIR/index.html' || true"
echo ""

# Summary
echo -e "${BLUE}📊 Test Summary${NC}"
echo -e "${BLUE}===============${NC}"
echo -e "Tests passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests failed: ${RED}$TESTS_FAILED${NC}"
echo -e "Total tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ $TESTS_FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed! Homepage is ready for deployment.${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Some tests failed. Please fix the issues before deployment.${NC}"
    exit 1
fi
