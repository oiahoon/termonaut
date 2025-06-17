# 🔧 Termonaut Troubleshooting Guide

This guide helps you resolve common issues with Termonaut installation, configuration, and usage.

## 📋 Table of Contents

- [Installation Issues](#installation-issues)
- [Shell Integration Problems](#shell-integration-problems)
- [Performance Issues](#performance-issues)
- [Database Problems](#database-problems)
- [Configuration Issues](#configuration-issues)
- [TUI Interface Problems](#tui-interface-problems)
- [Logging and Debugging](#logging-and-debugging)
- [Platform-Specific Issues](#platform-specific-issues)
- [Getting Help](#getting-help)

## 🔧 Installation Issues

### Binary Not Found After Installation

**Problem**: `command not found: termonaut` after installation

**Solutions**:
1. **Check PATH**: Ensure `/usr/local/bin` is in your PATH
   ```bash
   echo $PATH | grep "/usr/local/bin"
   ```

2. **Verify Installation Location**:
   ```bash
   which termonaut
   ls -la /usr/local/bin/termonaut
   ```

3. **Re-install with Explicit Path**:
   ```bash
   curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
   ```

4. **Manual Installation**:
   ```bash
   # Download for your platform
   wget https://github.com/oiahoon/termonaut/releases/latest/download/termonaut-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
   chmod +x termonaut-*
   sudo mv termonaut-* /usr/local/bin/termonaut
   ```

### Permission Denied Errors

**Problem**: Permission errors during installation

**Solutions**:
1. **Use sudo for System Installation**:
   ```bash
   sudo mv termonaut /usr/local/bin/
   ```

2. **Install to User Directory**:
   ```bash
   mkdir -p ~/.local/bin
   mv termonaut ~/.local/bin/
   export PATH="$HOME/.local/bin:$PATH"
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   ```

### Platform Compatibility Issues

**Problem**: Binary doesn't work on your platform

**Solutions**:
1. **Check Platform Support**:
   ```bash
   uname -s -m  # Should show supported platform
   ```

2. **Build from Source**:
   ```bash
   git clone https://github.com/oiahoon/termonaut.git
   cd termonaut
   go build -o termonaut cmd/termonaut/*.go
   ./termonaut version
   ```

## 🐚 Shell Integration Problems

### Hook Not Working

**Problem**: Commands not being tracked after `termonaut advanced shell install`

**Solutions**:
1. **Verify Hook Installation**:
   ```bash
   grep -n "termonaut" ~/.zshrc ~/.bashrc 2>/dev/null
   ```

2. **Restart Shell or Source Config**:
   ```bash
   # For Zsh
   source ~/.zshrc
   
   # For Bash
   source ~/.bashrc
   ```

3. **Manual Hook Installation**:
   ```bash
   # For Zsh - add to ~/.zshrc
   preexec() {
     /usr/local/bin/termonaut log-command "$1" &
   }
   
   # For Bash - add to ~/.bashrc
   trap 'termonaut log-command "$BASH_COMMAND" &' DEBUG
   ```

4. **Check Hook Function**:
   ```bash
   # Test the hook
   echo "test command"
   termonaut stats
   ```

### Job Control Messages

**Problem**: Seeing job control messages like `[1] 12345 done`

**Solutions**:
1. **Use Updated Hook** (if available):
   ```bash
   termonaut advanced shell install --force  # Re-install with latest hook
   ```

2. **Manual Fix for Zsh**:
   ```bash
   # Add to ~/.zshrc
   preexec() {
     { nohup /usr/local/bin/termonaut log-command "$1" > /dev/null 2>&1 & } 2>/dev/null
   }
   ```

3. **Manual Fix for Bash**:
   ```bash
   # Add to ~/.bashrc
   trap '{ nohup termonaut log-command "$BASH_COMMAND" > /dev/null 2>&1 & } 2>/dev/null' DEBUG
   ```

### Commands Not Being Categorized

**Problem**: All commands show as "Unknown" category

**Solutions**:
1. **Enable Command Categories**:
   ```bash
   termonaut config set command_categories true
   ```

2. **Check Configuration**:
   ```bash
   termonaut config get command_categories
   ```

## ⚡ Performance Issues

### Slow Command Execution

**Problem**: Terminal feels slower after installing Termonaut

**Solutions**:
1. **Check Logging Performance**:
   ```bash
   time echo "test"  # Should be < 10ms additional
   ```

2. **Disable Gamification Temporarily**:
   ```bash
   termonaut config set show_gamification false
   ```

3. **Reduce Database Size**:
   ```bash
   # Export data first
   termonaut export backup.json
   
   # Clean old data (keep last 30 days)
   sqlite3 ~/.termonaut/termonaut.db "DELETE FROM commands WHERE timestamp < datetime('now', '-30 days');"
   ```

### High Memory Usage

**Problem**: Termonaut using excessive memory

**Solutions**:
1. **Check Database Size**:
   ```bash
   du -h ~/.termonaut/termonaut.db
   ```

2. **Optimize Database**:
   ```bash
   sqlite3 ~/.termonaut/termonaut.db "VACUUM;"
   ```

3. **Limit Data Retention**:
   ```bash
   termonaut config set data_retention_days 30
   ```

## 💾 Database Problems

### Database Corruption

**Problem**: Database errors or corruption messages

**Solutions**:
1. **Check Database Integrity**:
   ```bash
   sqlite3 ~/.termonaut/termonaut.db "PRAGMA integrity_check;"
   ```

2. **Backup and Recreate**:
   ```bash
   # Backup existing data
   cp ~/.termonaut/termonaut.db ~/.termonaut/termonaut.db.backup
   
   # Try to export data
   termonaut export recovery.json
   
   # Remove corrupted database
   rm ~/.termonaut/termonaut.db
   
   # Restart termonaut (will create new database)
   termonaut stats
   
   # Import backup if successful
   termonaut import recovery.json
   ```

### Permission Issues with Database

**Problem**: Cannot write to database

**Solutions**:
1. **Check Directory Permissions**:
   ```bash
   ls -la ~/.termonaut/
   chmod 755 ~/.termonaut/
   chmod 644 ~/.termonaut/termonaut.db
   ```

2. **Check Disk Space**:
   ```bash
   df -h ~/.termonaut/
   ```

## ⚙️ Configuration Issues

### Configuration Not Loading

**Problem**: Settings not being applied

**Solutions**:
1. **Verify Config File**:
   ```bash
   cat ~/.termonaut/config.toml
   ```

2. **Validate TOML Syntax**:
   ```bash
   # Test with simple config
   termonaut config set theme emoji
   termonaut config get theme
   ```

3. **Reset to Defaults**:
   ```bash
   mv ~/.termonaut/config.toml ~/.termonaut/config.toml.backup
   termonaut config get  # Will create new default config
   ```

### GitHub Integration Not Working

**Problem**: GitHub badges not updating

**Solutions**:
1. **Check GitHub Configuration**:
   ```bash
   termonaut config get sync_enabled
   termonaut config get sync_repo
   ```

2. **Verify Repository Access**:
   ```bash
   termonaut github badges generate --test
   ```

3. **Check Badge URLs**:
   ```bash
   curl -I "https://shields.io/badge/Commands-123-brightgreen"
   ```

## 🖥️ TUI Interface Problems

### TUI Not Starting

**Problem**: `termonaut tui` fails to start

**Solutions**:
1. **Check Terminal Compatibility**:
   ```bash
   echo $TERM
   tput colors  # Should show 256 or more
   ```

2. **Use Simple Mode**:
   ```bash
   termonaut dashboard  # Alternative to TUI
   ```

3. **Check Terminal Size**:
   ```bash
   stty size  # Should be at least 24x80
   ```

### TUI Display Issues

**Problem**: Garbled text or broken layout

**Solutions**:
1. **Force UTF-8 Encoding**:
   ```bash
   export LANG=en_US.UTF-8
   export LC_ALL=en_US.UTF-8
   termonaut tui
   ```

2. **Disable Colors**:
   ```bash
   termonaut config set theme minimal
   termonaut tui
   ```

3. **Resize Terminal**:
   ```bash
   # Make terminal larger (at least 80x24)
   clear && termonaut tui
   ```

## 🐛 Logging and Debugging

### Enable Debug Logging

```bash
# Set debug log level
termonaut config set log_level debug

# Check logs
tail -f ~/.termonaut/termonaut.log
```

### Verbose Command Execution

```bash
# Run with verbose output
TERMONAUT_DEBUG=1 termonaut stats

# Test specific command
TERMONAUT_DEBUG=1 termonaut log-command "test command"
```

### Database Inspection

```bash
# Check command history
sqlite3 ~/.termonaut/termonaut.db "SELECT * FROM commands ORDER BY timestamp DESC LIMIT 10;"

# Check achievements
sqlite3 ~/.termonaut/termonaut.db "SELECT * FROM achievements;"

# Check sessions
sqlite3 ~/.termonaut/termonaut.db "SELECT * FROM sessions ORDER BY start_time DESC LIMIT 5;"
```

## 🖥️ Platform-Specific Issues

### macOS Issues

**macOS Catalina+ Security**:
```bash
# If binary is quarantined
sudo xattr -d com.apple.quarantine /usr/local/bin/termonaut
```

**macOS Shell Hook**:
```bash
# Zsh is default on macOS 10.15+
echo $SHELL
# Should be /bin/zsh or /usr/local/bin/zsh
```

### Linux Issues

**Missing Dependencies**:
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install sqlite3 curl

# CentOS/RHEL
sudo yum install sqlite curl

# Arch Linux
sudo pacman -S sqlite curl
```

### Windows (WSL) Issues

**Path Issues in WSL**:
```bash
# Ensure Linux path format
export PATH="/usr/local/bin:$PATH"
```

**WSL Shell Integration**:
```bash
# Make sure using Linux shell, not Windows
echo $0
# Should show bash or zsh, not cmd or powershell
```

## 🆘 Getting Help

### Collecting Debug Information

When reporting issues, include:

```bash
# System information
uname -a
echo $SHELL
echo $TERM

# Termonaut version and config
termonaut version
termonaut config get

# Database status
ls -la ~/.termonaut/
sqlite3 ~/.termonaut/termonaut.db "SELECT COUNT(*) FROM commands;"

# Recent logs
tail -20 ~/.termonaut/termonaut.log
```

### Common Diagnostic Commands

```bash
# Full system check
termonaut debug --system-info
termonaut debug --shell-integration
termonaut debug --database-check
```

### Report Issues

1. **GitHub Issues**: https://github.com/oiahoon/termonaut/issues
2. **Include debug information** from above
3. **Describe steps to reproduce**
4. **Specify your environment** (OS, shell, terminal)

### Community Support

- **Discussions**: GitHub Discussions for questions
- **Documentation**: Check README.md and other docs
- **Examples**: Look in `examples/` directory

---

## 🔍 Quick Diagnostics Checklist

- [ ] Binary installed and in PATH
- [ ] Shell hooks properly configured
- [ ] Database file exists and is writable
- [ ] Configuration file is valid TOML
- [ ] Terminal supports colors and UTF-8
- [ ] Sufficient disk space available
- [ ] No permission issues with ~/.termonaut/

If all checks pass and you're still having issues, try:
1. Restart your terminal completely
2. Run `termonaut advanced shell install --force` to reinstall hooks
3. Enable debug logging and check for errors
4. Report the issue with debug information 