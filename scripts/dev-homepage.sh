#!/bin/bash

# Termonaut Homepage Development Server
# This script provides a local development environment for the homepage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PORT=${PORT:-8000}
DOCS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../docs" && pwd)"

echo -e "${BLUE}🚀 Termonaut Homepage Development Server${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Check if docs directory exists
if [ ! -d "$DOCS_DIR" ]; then
    echo -e "${RED}❌ Error: docs directory not found at $DOCS_DIR${NC}"
    exit 1
fi

# Check if index.html exists
if [ ! -f "$DOCS_DIR/index.html" ]; then
    echo -e "${RED}❌ Error: index.html not found in docs directory${NC}"
    exit 1
fi

echo -e "${GREEN}📁 Serving from: $DOCS_DIR${NC}"
echo -e "${GREEN}🌐 Port: $PORT${NC}"
echo ""

# Function to check if port is available
check_port() {
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Port $PORT is already in use${NC}"
        echo -e "${YELLOW}   Trying to find an available port...${NC}"
        
        # Find next available port
        while lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; do
            PORT=$((PORT + 1))
        done
        
        echo -e "${GREEN}✅ Using port $PORT instead${NC}"
    fi
}

# Function to open browser
open_browser() {
    local url="http://localhost:$PORT"
    echo -e "${BLUE}🌐 Opening browser at $url${NC}"
    
    # Detect OS and open browser
    case "$(uname -s)" in
        Darwin)  # macOS
            open "$url"
            ;;
        Linux)
            if command -v xdg-open > /dev/null; then
                xdg-open "$url"
            elif command -v gnome-open > /dev/null; then
                gnome-open "$url"
            fi
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)  # Windows
            start "$url"
            ;;
    esac
}

# Function to start Python server
start_python_server() {
    echo -e "${BLUE}🐍 Starting Python HTTP server...${NC}"
    cd "$DOCS_DIR"
    
    # Try Python 3 first, then Python 2
    if command -v python3 > /dev/null; then
        python3 -m http.server $PORT
    elif command -v python > /dev/null; then
        python -m SimpleHTTPServer $PORT
    else
        echo -e "${RED}❌ Python not found${NC}"
        return 1
    fi
}

# Function to start Node.js server
start_node_server() {
    echo -e "${BLUE}📦 Starting Node.js server...${NC}"
    cd "$DOCS_DIR"
    
    if command -v npx > /dev/null; then
        npx serve . -p $PORT
    else
        echo -e "${RED}❌ Node.js/npx not found${NC}"
        return 1
    fi
}

# Function to start PHP server
start_php_server() {
    echo -e "${BLUE}🐘 Starting PHP built-in server...${NC}"
    cd "$DOCS_DIR"
    php -S localhost:$PORT
}

# Function to validate HTML
validate_html() {
    echo -e "${BLUE}🔍 Validating HTML...${NC}"
    
    if command -v tidy > /dev/null; then
        if tidy -q -e "$DOCS_DIR/index.html" 2>/dev/null; then
            echo -e "${GREEN}✅ HTML validation passed${NC}"
        else
            echo -e "${YELLOW}⚠️  HTML validation warnings (non-critical)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  HTML Tidy not installed, skipping validation${NC}"
        echo -e "${YELLOW}   Install with: brew install tidy-html5 (macOS) or apt-get install tidy (Linux)${NC}"
    fi
}

# Function to check links
check_links() {
    echo -e "${BLUE}🔗 Checking internal links...${NC}"
    
    # Basic link checking
    local broken_links=0
    
    # Extract internal links
    grep -o 'href="#[^"]*"' "$DOCS_DIR/index.html" | while read -r link; do
        id=$(echo "$link" | sed 's/href="#\([^"]*\)"/\1/')
        if grep -q "id=\"$id\"" "$DOCS_DIR/index.html"; then
            echo -e "${GREEN}✅ $link${NC}"
        else
            echo -e "${RED}❌ $link (target not found)${NC}"
            broken_links=$((broken_links + 1))
        fi
    done
    
    if [ $broken_links -eq 0 ]; then
        echo -e "${GREEN}✅ All internal links are valid${NC}"
    fi
}

# Function to show development tips
show_tips() {
    echo ""
    echo -e "${BLUE}💡 Development Tips:${NC}"
    echo -e "${BLUE}==================${NC}"
    echo "• Press Ctrl+C to stop the server"
    echo "• Edit files in $DOCS_DIR"
    echo "• Refresh browser to see changes"
    echo "• Check browser console for JavaScript errors"
    echo "• Use browser dev tools to test responsive design"
    echo ""
    echo -e "${BLUE}🔧 Useful URLs:${NC}"
    echo "• Homepage: http://localhost:$PORT"
    echo "• Direct CSS: http://localhost:$PORT/assets/css/style.css"
    echo "• Direct JS: http://localhost:$PORT/assets/js/main.js"
    echo ""
}

# Function to cleanup on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}🛑 Shutting down development server...${NC}"
    echo -e "${GREEN}✅ Thanks for developing Termonaut! 🚀${NC}"
}

# Set trap for cleanup
trap cleanup EXIT

# Main execution
main() {
    # Validate HTML first
    validate_html
    echo ""
    
    # Check links
    check_links
    echo ""
    
    # Check port availability
    check_port
    
    # Show tips
    show_tips
    
    # Start server based on available tools
    echo -e "${BLUE}🚀 Starting development server...${NC}"
    echo ""
    
    # Open browser after a short delay
    (sleep 2 && open_browser) &
    
    # Try different server options
    if command -v python3 > /dev/null || command -v python > /dev/null; then
        start_python_server
    elif command -v npx > /dev/null; then
        start_node_server
    elif command -v php > /dev/null; then
        start_php_server
    else
        echo -e "${RED}❌ No suitable server found${NC}"
        echo -e "${YELLOW}Please install one of the following:${NC}"
        echo "• Python 3: python3 -m http.server"
        echo "• Node.js: npx serve"
        echo "• PHP: php -S localhost:8000"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Termonaut Homepage Development Server"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --port PORT    Set custom port (default: 8000)"
        echo "  --validate     Only validate HTML and exit"
        echo "  --check-links  Only check links and exit"
        echo ""
        echo "Environment Variables:"
        echo "  PORT          Custom port number"
        echo ""
        echo "Examples:"
        echo "  $0                    # Start server on port 8000"
        echo "  $0 --port 3000       # Start server on port 3000"
        echo "  PORT=9000 $0         # Start server on port 9000"
        echo "  $0 --validate        # Only validate HTML"
        exit 0
        ;;
    --port)
        if [ -n "${2:-}" ]; then
            PORT="$2"
            shift 2
        else
            echo -e "${RED}❌ Error: --port requires a port number${NC}"
            exit 1
        fi
        ;;
    --validate)
        validate_html
        exit 0
        ;;
    --check-links)
        check_links
        exit 0
        ;;
    --*)
        echo -e "${RED}❌ Error: Unknown option $1${NC}"
        echo "Use --help for usage information"
        exit 1
        ;;
esac

# Run main function
main "$@"
