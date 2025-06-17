#!/bin/bash
set -e

# Termonaut Installation Script
# This script automatically detects your platform and installs the latest version of Termonaut

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# GitHub repository
REPO="oiahoon/termonaut"
GITHUB_API="https://api.github.com/repos/${REPO}"
GITHUB_RELEASES="${GITHUB_API}/releases/latest"

# Default installation directory
INSTALL_DIR="/usr/local/bin"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect the platform
detect_platform() {
    local platform=""
    local arch=""

    case "$(uname -s)" in
        Linux*)
            platform="linux"
            ;;
        Darwin*)
            platform="darwin"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            platform="windows"
            ;;
        *)
            print_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        i386|i686)
            arch="386"
            ;;
        armv7l)
            arch="arm"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${platform}-${arch}"
}

# Get the latest release version
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "${GITHUB_RELEASES}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "${GITHUB_RELEASES}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download and install Termonaut
install_termonaut() {
    local platform="$1"
    local version="$2"
    local temp_dir="/tmp/termonaut-install"
    local binary_name="termonaut"

    if [[ "$platform" == *"windows"* ]]; then
        binary_name="termonaut.exe"
    fi

    local download_url="https://github.com/${REPO}/releases/download/${version}/termonaut-${platform}"
    local temp_binary="${temp_dir}/${binary_name}"

    print_status "Creating temporary directory..."
    mkdir -p "$temp_dir"

    print_status "Downloading Termonaut ${version} for ${platform}..."
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_binary" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_binary" "$download_url"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi

    # Check if download was successful
    if [[ ! -f "$temp_binary" ]]; then
        print_error "Failed to download Termonaut binary"
        exit 1
    fi

    # Make binary executable
    chmod +x "$temp_binary"

    # Install binary
    print_status "Installing Termonaut to ${INSTALL_DIR}..."

    # Try to install with sudo if needed
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$temp_binary" "${INSTALL_DIR}/termonaut"
    else
        print_status "Root privileges required for installation to ${INSTALL_DIR}"
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$temp_binary" "${INSTALL_DIR}/termonaut"
        else
            print_error "sudo is not available. Please run as root or install manually."
            print_status "Manual installation: cp ${temp_binary} ${INSTALL_DIR}/termonaut"
            exit 1
        fi
    fi

    # Clean up
    rm -rf "$temp_dir"

    print_success "Termonaut ${version} installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v termonaut >/dev/null 2>&1; then
        local installed_version=$(termonaut --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation verified: ${installed_version}"
        return 0
    else
        print_error "Installation verification failed. Termonaut not found in PATH."
        return 1
    fi
}

# Setup shell integration
setup_shell_integration() {
    print_status "Setting up shell integration..."

    if termonaut init >/dev/null 2>&1; then
        print_success "Shell integration setup completed!"
        print_warning "Please restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc) to activate Termonaut."
    else
        print_warning "Shell integration setup failed. You can set it up later with 'termonaut init'."
    fi
}

# Show usage instructions
show_usage() {
    echo
    print_success "🚀 Termonaut has been installed successfully!"
    echo
    echo "Quick start:"
    echo "  1. Restart your terminal or source your shell config:"
    echo "     source ~/.bashrc   # or ~/.zshrc for zsh users"
    echo
    echo "  2. Start tracking your terminal activity:"
    echo "     termonaut stats"
    echo
    echo "  3. Check your progress:"
    echo "     termonaut xp"
    echo "     termonaut badges"
    echo
    echo "  4. Explore features:"
    echo "     termonaut quests       # Daily challenges"
    echo "     termonaut multipliers  # XP bonuses"
    echo "     termonaut heatmap generate  # Activity heatmap"
    echo
    echo "  5. Configure display mode:"
    echo "     termonaut display mode rich    # Full features"
    echo "     termonaut display mode minimal # Simple output"
    echo
    echo "For more information:"
    echo "  termonaut --help"
    echo "  https://github.com/oiahoon/termonaut"
    echo
}

# Main installation function
main() {
    echo
    echo "🚀 Termonaut Installation Script"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo

    # Check if Termonaut is already installed
    if command -v termonaut >/dev/null 2>&1; then
        local current_version=$(termonaut --version 2>/dev/null | head -n1 || echo "unknown")
        print_warning "Termonaut is already installed: ${current_version}"
        read -p "Do you want to update to the latest version? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Installation cancelled."
            exit 0
        fi
    fi

    # Detect platform
    print_status "Detecting platform..."
    local platform=$(detect_platform)
    print_success "Detected platform: ${platform}"

    # Get latest version
    print_status "Fetching latest version..."
    local version=$(get_latest_version)
    if [[ -z "$version" ]]; then
        print_error "Failed to fetch latest version information"
        exit 1
    fi
    print_success "Latest version: ${version}"

    # Install
    install_termonaut "$platform" "$version"

    # Verify
    if verify_installation; then
        # Setup shell integration
        setup_shell_integration

        # Show usage
        show_usage
    else
        exit 1
    fi
}

# Handle command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help|-h)
            echo "Termonaut Installation Script"
            echo
            echo "Usage: $0 [options]"
            echo
            echo "Options:"
            echo "  --install-dir DIR    Install to custom directory (default: /usr/local/bin)"
            echo "  --help, -h          Show this help message"
            echo
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information."
            exit 1
            ;;
    esac
done

# Run main installation
main