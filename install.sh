#!/bin/bash
set -e

# Termonaut Installation Script
# This script addresses common installation issues and provides better error handling

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
USER_INSTALL_DIR="$HOME/.local/bin"

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

# Check if running as root
is_root() {
    [ "$(id -u)" = "0" ]
}

# Check if directory is writable
is_writable() {
    [ -w "$1" ] 2>/dev/null
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
            print_error "Supported platforms: Linux, macOS (Darwin), Windows"
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
            print_error "Supported architectures: x86_64, arm64, i386, armv7l"
            exit 1
            ;;
    esac

    echo "${platform}-${arch}"
}

# Get the latest release version
get_latest_version() {
    print_status "Fetching latest version information..."

    local version=""
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "${GITHUB_RELEASES}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "${GITHUB_RELEASES}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        print_error "  Ubuntu/Debian: sudo apt install curl"
        print_error "  CentOS/RHEL: sudo yum install curl"
        print_error "  macOS: curl should be pre-installed"
        exit 1
    fi

    if [ -z "$version" ]; then
        print_error "Failed to fetch version information from GitHub API"
        print_error "Please check your internet connection and try again"
        exit 1
    fi

    echo "$version"
}

# Determine best installation directory
choose_install_dir() {
    local target_dir=""

    # If user specified a directory, use it
    if [ -n "$CUSTOM_INSTALL_DIR" ]; then
        target_dir="$CUSTOM_INSTALL_DIR"
        print_status "Using custom installation directory: $target_dir"
    # If running as root or /usr/local/bin is writable, use it
    elif is_root || is_writable "/usr/local/bin"; then
        target_dir="/usr/local/bin"
        print_status "Using system installation directory: $target_dir"
    # Otherwise, use user directory
    else
        target_dir="$USER_INSTALL_DIR"
        print_status "Using user installation directory: $target_dir"
        print_warning "You may need to add $target_dir to your PATH"
    fi

    # Create directory if it doesn't exist
    if [ ! -d "$target_dir" ]; then
        print_status "Creating directory: $target_dir"
        if ! mkdir -p "$target_dir" 2>/dev/null; then
            if [ "$target_dir" = "/usr/local/bin" ] && command -v sudo >/dev/null 2>&1; then
                sudo mkdir -p "$target_dir"
            else
                print_error "Failed to create directory: $target_dir"
                exit 1
            fi
        fi
    fi

    echo "$target_dir"
}

# Download and install Termonaut
install_termonaut() {
    local platform="$1"
    local version="$2"
    local install_dir="$3"
    local temp_dir="/tmp/termonaut-install-$$"
    local binary_name="termonaut"

    if [[ "$platform" == *"windows"* ]]; then
        binary_name="termonaut.exe"
    fi

    local download_url="https://github.com/${REPO}/releases/download/${version}/termonaut-${platform}"
    local temp_binary="${temp_dir}/${binary_name}"
    local final_binary="${install_dir}/termonaut"

    print_status "Creating temporary directory..."
    mkdir -p "$temp_dir"

    print_status "Downloading Termonaut ${version} for ${platform}..."
    print_status "Download URL: $download_url"

    # Download with better error handling
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L --fail --silent --show-error -o "$temp_binary" "$download_url"; then
            print_error "Failed to download from: $download_url"
            print_error "Please check if the release exists for your platform"
            rm -rf "$temp_dir"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget --quiet -O "$temp_binary" "$download_url"; then
            print_error "Failed to download from: $download_url"
            print_error "Please check if the release exists for your platform"
            rm -rf "$temp_dir"
            exit 1
        fi
    fi

    # Verify download
    if [[ ! -f "$temp_binary" ]] || [[ ! -s "$temp_binary" ]]; then
        print_error "Downloaded file is missing or empty"
        rm -rf "$temp_dir"
        exit 1
    fi

    # Make binary executable
    chmod +x "$temp_binary"

    # Test the binary
    print_status "Testing downloaded binary..."
    if ! "$temp_binary" --version >/dev/null 2>&1; then
        print_warning "Binary test failed, but continuing with installation"
    fi

    # Install binary
    print_status "Installing Termonaut to ${install_dir}..."

    # Try different installation methods based on permissions
    if is_writable "$install_dir"; then
        # Direct copy if directory is writable
        if ! cp "$temp_binary" "$final_binary"; then
            print_error "Failed to copy binary to $final_binary"
            rm -rf "$temp_dir"
            exit 1
        fi
    elif command -v sudo >/dev/null 2>&1; then
        # Use sudo if available
        print_status "Root privileges required for installation to ${install_dir}"
        if ! sudo cp "$temp_binary" "$final_binary"; then
            print_error "Failed to install with sudo"
            rm -rf "$temp_dir"
            exit 1
        fi
    else
        # Fallback: suggest manual installation
        print_error "Cannot install to $install_dir (no write permission and sudo not available)"
        print_error "Please run one of the following commands manually:"
        print_error "  sudo cp $temp_binary $final_binary"
        print_error "  OR"
        print_error "  mkdir -p ~/.local/bin && cp $temp_binary ~/.local/bin/termonaut"
        print_error "  export PATH=\"\$HOME/.local/bin:\$PATH\""
        exit 1
    fi

    # Set proper permissions
    if [ "$install_dir" = "/usr/local/bin" ] && ! is_writable "$final_binary"; then
        sudo chmod 755 "$final_binary"
    else
        chmod 755 "$final_binary"
    fi

    # Clean up
    rm -rf "$temp_dir"

    print_success "Termonaut ${version} installed successfully to ${final_binary}!"
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."

    if command -v termonaut >/dev/null 2>&1; then
        local installed_version
        installed_version=$(termonaut --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation verified: ${installed_version}"

        # Test basic functionality
        if termonaut --help >/dev/null 2>&1; then
            print_success "Basic functionality test passed"
        else
            print_warning "Basic functionality test failed"
        fi

        return 0
    else
        print_error "Installation verification failed. Termonaut not found in PATH."
        print_error "You may need to:"
        print_error "  1. Restart your terminal"
        print_error "  2. Add the installation directory to your PATH"
        print_error "  3. Run: export PATH=\"\$PATH:/usr/local/bin\""
        return 1
    fi
}

# Update PATH if needed
update_path() {
    local install_dir="$1"

    # Check if install_dir is in PATH
    if echo "$PATH" | grep -q "$install_dir"; then
        print_success "Installation directory is already in PATH"
        return 0
    fi

    # Add to PATH for user installations
    if [ "$install_dir" = "$USER_INSTALL_DIR" ]; then
        print_status "Adding $install_dir to PATH..."

        # Determine shell config file
        local shell_config=""
        if [ -n "$ZSH_VERSION" ] || [[ "$SHELL" == *"zsh"* ]]; then
            shell_config="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ] || [[ "$SHELL" == *"bash"* ]]; then
            if [ -f "$HOME/.bashrc" ]; then
                shell_config="$HOME/.bashrc"
            else
                shell_config="$HOME/.bash_profile"
            fi
        fi

        if [ -n "$shell_config" ]; then
            if ! grep -q "$install_dir" "$shell_config" 2>/dev/null; then
                echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$shell_config"
                print_success "Added $install_dir to PATH in $shell_config"
                print_warning "Please restart your terminal or run: source $shell_config"
            fi
        else
            print_warning "Please add $install_dir to your PATH manually"
        fi
    fi
}

# Setup shell integration
setup_shell_integration() {
    print_status "Setting up shell integration..."

    if termonaut advanced shell install >/dev/null 2>&1; then
        print_success "Shell integration setup completed!"
        print_warning "Please restart your terminal or source your shell config to activate Termonaut"
    else
        print_warning "Shell integration setup failed. You can set it up later with:"
        print_warning "  termonaut advanced shell install"
    fi
}

# Show usage instructions
show_usage() {
    echo
    print_success "🚀 Termonaut has been installed successfully!"
    echo
    echo "Quick start:"
    echo "  1. Restart your terminal or source your shell config:"
    echo "     source ~/.bashrc   # for bash users"
    echo "     source ~/.zshrc    # for zsh users"
    echo
    echo "  2. Start tracking your terminal activity:"
    echo "     termonaut stats"
    echo "     tn stats           # short alias"
    echo
    echo "  3. Check your progress:"
    echo "     termonaut xp"
    echo "     termonaut badges"
    echo
    echo "  4. Configure the tool:"
    echo "     termonaut config set theme emoji"
    echo "     termonaut config set empty_command_stats true"
    echo
    echo "  5. Get help:"
    echo "     termonaut --help"
    echo "     termonaut advanced --help"
    echo
    echo "For more information:"
    echo "  📖 Documentation: https://github.com/oiahoon/termonaut"
    echo "  🐛 Issues: https://github.com/oiahoon/termonaut/issues"
    echo
}

# Show troubleshooting info
show_troubleshooting() {
    echo
    print_error "Installation failed. Here are some troubleshooting steps:"
    echo
    echo "1. Check your internet connection and try again"
    echo "2. Try manual installation:"
    echo "   wget https://github.com/oiahoon/termonaut/releases/latest/download/termonaut-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
    echo "   chmod +x termonaut-*"
    echo "   sudo mv termonaut-* /usr/local/bin/termonaut"
    echo
    echo "3. Use Homebrew (macOS/Linux):"
    echo "   brew tap oiahoon/termonaut"
    echo "   brew install termonaut"
    echo
    echo "4. Build from source:"
    echo "   git clone https://github.com/oiahoon/termonaut.git"
    echo "   cd termonaut"
    echo "   go build -o termonaut cmd/termonaut/*.go"
    echo "   sudo mv termonaut /usr/local/bin/"
    echo
    echo "5. Report the issue at: https://github.com/oiahoon/termonaut/issues"
    echo
}

# Main installation function
main() {
    echo
    echo "🚀 Termonaut Installation Script"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo

    # Show system information
    print_status "System information:"
    echo "  OS: $(uname -s)"
    echo "  Architecture: $(uname -m)"
    echo "  Shell: $SHELL"
    echo

    # Check if Termonaut is already installed
    if command -v termonaut >/dev/null 2>&1; then
        local current_version
        current_version=$(termonaut --version 2>/dev/null | head -n1 || echo "unknown")
        print_warning "Termonaut is already installed: ${current_version}"
        echo
        read -p "Do you want to update to the latest version? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Installation cancelled."
            exit 0
        fi
        echo
    fi

    # Detect platform
    print_status "Detecting platform..."
    local platform
    platform=$(detect_platform)
    print_success "Detected platform: ${platform}"

    # Get latest version
    local version
    version=$(get_latest_version)
    print_success "Latest version: ${version}"

    # Choose installation directory
    local install_dir
    install_dir=$(choose_install_dir)
    print_success "Installation directory: ${install_dir}"

    # Install
    install_termonaut "$platform" "$version" "$install_dir"

    # Update PATH if needed
    update_path "$install_dir"

    # Verify installation
    if verify_installation; then
        # Setup shell integration
        setup_shell_integration

        # Show usage instructions
        show_usage
    else
        show_troubleshooting
        exit 1
    fi
}

# Handle command line arguments
CUSTOM_INSTALL_DIR=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --install-dir)
            CUSTOM_INSTALL_DIR="$2"
            if [ -z "$CUSTOM_INSTALL_DIR" ]; then
                print_error "Missing argument for --install-dir"
                exit 1
            fi
            shift 2
            ;;
        --help|-h)
            echo "Termonaut Installation Script"
            echo
            echo "Usage: $0 [options]"
            echo
            echo "Options:"
            echo "  --install-dir DIR    Install to custom directory"
            echo "                       (default: /usr/local/bin or ~/.local/bin)"
            echo "  --help, -h          Show this help message"
            echo
            echo "Examples:"
            echo "  $0                           # Auto-detect best installation directory"
            echo "  $0 --install-dir ~/.local/bin # Install to user directory"
            echo "  $0 --install-dir /opt/bin     # Install to custom directory"
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