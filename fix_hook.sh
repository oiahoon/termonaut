#!/bin/bash

# Termonaut Hook Fix Script
# This script updates your shell hook to use silent background execution
# Fixes the [1] job status output issue

echo "🔧 Termonaut Hook Fix Script"
echo "=========================================="
echo

# Get the current shell
current_shell=$(basename "$SHELL")
echo "Detected shell: $current_shell"

# Determine config file
if [[ "$current_shell" == "zsh" ]]; then
    config_file="$HOME/.zshrc"
elif [[ "$current_shell" == "bash" ]]; then
    if [[ -f "$HOME/.bashrc" ]]; then
        config_file="$HOME/.bashrc"
    else
        config_file="$HOME/.bash_profile"
    fi
else
    echo "❌ Unsupported shell: $current_shell"
    exit 1
fi

echo "Config file: $config_file"
echo

# Check if Termonaut hook exists
if ! grep -q "termonaut" "$config_file"; then
    echo "❌ No Termonaut hook found in $config_file"
    echo "Run 'termonaut init' to install the hook first."
    exit 1
fi

echo "✅ Found existing Termonaut hook"
echo

# Backup the config file
backup_file="${config_file}.backup.$(date +%Y%m%d_%H%M%S)"
cp "$config_file" "$backup_file"
echo "📋 Created backup: $backup_file"

# Remove old hook
if [[ "$current_shell" == "zsh" ]]; then
    # Remove Zsh hook
    sed -i '' '/# Termonaut shell integration/,/^fi$/d' "$config_file"
else
    # Remove Bash hook
    sed -i '' '/# Termonaut shell integration/,/trap.*termonaut_log_command.*DEBUG/d' "$config_file"
fi

echo "🗑️  Removed old hook"

# Get Termonaut binary path
if command -v termonaut >/dev/null 2>&1; then
    termonaut_path=$(which termonaut)
else
    # Try local binary
    if [[ -f "./termonaut" ]]; then
        termonaut_path="$(pwd)/termonaut"
    else
        echo "❌ Termonaut binary not found"
        echo "Please ensure 'termonaut' is in your PATH or run this script from the Termonaut directory"
        exit 1
    fi
fi

echo "🎯 Using Termonaut binary: $termonaut_path"

# Add new silent hook
if [[ "$current_shell" == "zsh" ]]; then
    # Add Zsh hook with enhanced job control suppression (v0.9.0 RC)
    cat >> "$config_file" << EOF

# Termonaut shell integration (v0.9.0 RC - Enhanced)
termonaut_preexec() {
    # Silent background execution with comprehensive job control suppression
    {
        # Create a completely detached subshell
        (
            # Disable all job control and output
            set +m 2>/dev/null
            unset HISTFILE 2>/dev/null
            exec </dev/null >/dev/null 2>&1
            
            # Run termonaut in completely isolated environment
            $termonaut_path log-command "\$1" &
            
            # Force exit to prevent any shell interaction
            exit 0
        ) &
        disown %% 2>/dev/null || true
    } 2>/dev/null
}

# Check if preexec_functions exists, if not create it
if [[ -z "\${preexec_functions+x}" ]]; then
    preexec_functions=()
fi

# Add our function to preexec_functions if not already present
if [[ ! " \${preexec_functions[@]} " =~ " termonaut_preexec " ]]; then
    preexec_functions+=(termonaut_preexec)
fi
EOF
else
    # Add Bash hook with enhanced job control suppression (v0.9.0 RC)
    cat >> "$config_file" << EOF

# Termonaut shell integration (v0.9.0 RC - Enhanced)
termonaut_log_command() {
    if [ -n "\$BASH_COMMAND" ]; then
        # Silent background execution with comprehensive job control suppression
        {
            # Create a completely detached subshell
            (
                # Disable all job control and output
                set +m 2>/dev/null
                unset HISTFILE 2>/dev/null
                exec </dev/null >/dev/null 2>&1
                
                # Run termonaut in completely isolated environment
                $termonaut_path log-command "\$BASH_COMMAND" &
                
                # Force exit to prevent any shell interaction
                exit 0
            ) &
            disown \$! 2>/dev/null || true
        } 2>/dev/null
    fi
}

# Set up DEBUG trap
trap 'termonaut_log_command' DEBUG
EOF
fi

echo "✅ Added new silent hook"
echo

echo "🎉 Hook fix completed!"
echo
echo "The fix includes:"
echo "  • Silent background execution ({ command & } 2>/dev/null)"
echo "  • Redirects stdout and stderr to /dev/null"
echo "  • Eliminates job status notifications [1] xxx done"
echo
echo "To apply the changes:"
if [[ "$current_shell" == "zsh" ]]; then
    echo "  source ~/.zshrc"
else
    echo "  source $config_file"
fi
echo
echo "Or restart your terminal."
echo
echo "Your original config has been backed up to: $backup_file"