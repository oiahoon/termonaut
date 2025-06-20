package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ShellType represents different shell types
type ShellType string

const (
	// Zsh shell
	Zsh ShellType = "zsh"
	// Bash shell
	Bash ShellType = "bash"
	// Fish shell
	Fish ShellType = "fish"
	// PowerShell
	PowerShell ShellType = "powershell"
)

// HookInstaller handles shell hook installation
type HookInstaller struct {
	shellType  ShellType
	configFile string
	binaryPath string
}

// SafeConfigManager handles safe modification of shell configuration files
type SafeConfigManager struct {
	configFile string
	shellType  ShellType
}

// ConfigBackup represents a backup of the configuration file
type ConfigBackup struct {
	OriginalPath string
	BackupPath   string
	Timestamp    time.Time
}

// NewHookInstaller creates a new hook installer
func NewHookInstaller(binaryPath string) (*HookInstaller, error) {
	shellType, configFile, err := detectShell()
	if err != nil {
		return nil, fmt.Errorf("failed to detect shell: %w", err)
	}

	return &HookInstaller{
		shellType:  shellType,
		configFile: configFile,
		binaryPath: binaryPath,
	}, nil
}

// NewSafeConfigManager creates a new safe configuration manager
func NewSafeConfigManager(configFile string, shellType ShellType) *SafeConfigManager {
	return &SafeConfigManager{
		configFile: configFile,
		shellType:  shellType,
	}
}

// Install installs the shell hook
func (h *HookInstaller) Install() error {
	return h.InstallWithForce(false)
}

// InstallWithForce installs the shell hook with force option
func (h *HookInstaller) InstallWithForce(force bool) error {
	// Create tn symlink if it doesn't exist
	if err := h.createShortcutSymlink(); err != nil {
		// Don't fail installation if symlink creation fails
		fmt.Printf("Warning: Could not create 'tn' shortcut: %v\n", err)
	}

	// Create safe config manager
	configManager := NewSafeConfigManager(h.configFile, h.shellType)

	// Check if already installed (unless force is used)
	if !force {
		_, startIdx, _, err := configManager.GetTermonautBlock()
		if err != nil && !strings.Contains(err.Error(), "incomplete") {
			return fmt.Errorf("failed to check installation status: %w", err)
		}
		if startIdx != -1 {
			return fmt.Errorf("termonaut hook is already installed")
		}
	} else {
		// Force mode: safely remove existing installation first
		if err := configManager.RemoveTermonautBlock(); err != nil {
			fmt.Printf("Warning: Failed to remove existing hook: %v\n", err)
		}
	}

	// Generate hook content based on shell type
	var hookContent string
	switch h.shellType {
	case Zsh:
		hookContent = h.generateZshHook()
	case Bash:
		hookContent = h.generateBashHook()
	case Fish:
		hookContent = h.generateFishHook()
	case PowerShell:
		hookContent = h.generatePowerShellHook()
	default:
		return fmt.Errorf("unsupported shell: %s", h.shellType)
	}

	// Safely add the hook
	if err := configManager.AddTermonautBlock(hookContent); err != nil {
		return fmt.Errorf("failed to install hook: %w", err)
	}

	fmt.Printf("✅ Termonaut hook installed successfully in %s\n", h.configFile)
	return nil
}

// generateZshHook generates the Zsh hook content
func (h *HookInstaller) generateZshHook() string {
	return fmt.Sprintf(`# Termonaut shell integration (v0.9.3 Safe)
termonaut_preexec() {
    # Complete job control suppression - eliminate ALL job messages
    {
        # Disable job control notifications globally
        setopt NO_NOTIFY 2>/dev/null || true
        setopt NO_HUP 2>/dev/null || true
        setopt NO_BG_NICE 2>/dev/null || true
        setopt NO_CHECK_JOBS 2>/dev/null || true

        # Method 1: Use subshell with complete isolation
        (
            # Run in completely isolated subshell
            %s log-command "$1" >/dev/null 2>&1 &
            disown %%%% 2>/dev/null || true
        ) >/dev/null 2>&1 &

        # Method 2: Disown the subshell itself
        disown %%%% 2>/dev/null || true

        # Method 3: Clear job table
        jobs >/dev/null 2>&1 | while read job; do
            disown "$job" 2>/dev/null || true
        done

        # Method 4: Reset job control state
        setopt NO_MONITOR 2>/dev/null || true

    } >/dev/null 2>&1
}

# Check if preexec_functions exists, if not create it
if [[ -z "${preexec_functions+x}" ]]; then
    preexec_functions=()
fi

# Add our function to preexec_functions if not already present
if [[ ! " ${preexec_functions[@]} " =~ " termonaut_preexec " ]]; then
    preexec_functions+=(termonaut_preexec)
fi`, h.binaryPath)
}

// generateBashHook generates the Bash hook content
func (h *HookInstaller) generateBashHook() string {
	return fmt.Sprintf(`# Termonaut shell integration (v0.9.3 Safe)
termonaut_log_command() {
    if [ -n "$BASH_COMMAND" ]; then
        # Complete job control suppression - eliminate ALL job messages
        {
            # Disable job control globally
            set +m 2>/dev/null || true
            set +b 2>/dev/null || true

            # Method 1: Use exec with complete redirection
            (
                exec %s log-command "$BASH_COMMAND" >/dev/null 2>&1 &
            ) 2>/dev/null &

            # Method 2: Disown all background jobs
            disown $! 2>/dev/null || true
            jobs | awk '{print $1}' | while read job; do
                disown "$job" 2>/dev/null || true
            done

        } 2>/dev/null
    fi
}

# Set up DEBUG trap
trap 'termonaut_log_command' DEBUG`, h.binaryPath)
}

// generateFishHook generates the Fish hook content
func (h *HookInstaller) generateFishHook() string {
	return fmt.Sprintf(`# Termonaut shell integration (v0.9.3 Safe)
function termonaut_preexec --on-event fish_preexec
    %s log-command "$argv" >/dev/null 2>&1 &
    disown
end`, h.binaryPath)
}

// generatePowerShellHook generates the PowerShell hook content
func (h *HookInstaller) generatePowerShellHook() string {
	return fmt.Sprintf(`# Termonaut shell integration (v0.9.3 Safe)
function Invoke-TermonautLogging {
    param($Command)
    try {
        Start-Job -ScriptBlock {
            param($BinaryPath, $Cmd)
            & $BinaryPath log-command $Cmd 2>$null
        } -ArgumentList "%s", $Command | Out-Null
    } catch {
        # Silently ignore errors
    }
}

# PowerShell command history hook
$PSDefaultParameterValues['*:Verbose'] = $false
$PSDefaultParameterValues['*:Debug'] = $false

# Override the prompt to capture commands
function global:prompt {
    $history = Get-History -Count 1 -ErrorAction SilentlyContinue
    if ($history -and $history.CommandLine) {
        Invoke-TermonautLogging -Command $history.CommandLine
    }

    # Return original prompt
    "PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) "
}`, h.binaryPath)
}

// Uninstall removes the shell hook
func (h *HookInstaller) Uninstall() error {
	// Use safe config manager for removal
	configManager := NewSafeConfigManager(h.configFile, h.shellType)
	return configManager.RemoveTermonautBlock()
}

// IsInstalled checks if the hook is already installed
func (h *HookInstaller) IsInstalled() (bool, error) {
	content, err := os.ReadFile(h.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read config file: %w", err)
	}

	contentStr := string(content)

	// Check for specific hook markers based on shell type
	switch h.shellType {
	case Zsh:
		return strings.Contains(contentStr, "termonaut_preexec") &&
			strings.Contains(contentStr, "preexec_functions"), nil
	case Bash:
		return strings.Contains(contentStr, "termonaut_log_command") &&
			strings.Contains(contentStr, "trap 'termonaut_log_command' DEBUG"), nil
	case Fish:
		return strings.Contains(contentStr, "function termonaut_preexec"), nil
	case PowerShell:
		return strings.Contains(contentStr, "Invoke-TermonautLogging"), nil
	default:
		// Fallback to generic check
		return strings.Contains(contentStr, "termonaut"), nil
	}
}

// GetShellType returns the detected shell type
func (h *HookInstaller) GetShellType() ShellType {
	return h.shellType
}

// detectShell detects the current shell and returns its config file
func detectShell() (ShellType, string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", "", fmt.Errorf("SHELL environment variable not set")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if strings.Contains(shell, "zsh") {
		configFile := filepath.Join(homeDir, ".zshrc")
		return Zsh, configFile, nil
	} else if strings.Contains(shell, "bash") {
		// Try .bashrc first, then .bash_profile
		bashrc := filepath.Join(homeDir, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			return Bash, bashrc, nil
		}
		bashProfile := filepath.Join(homeDir, ".bash_profile")
		return Bash, bashProfile, nil
	} else if strings.Contains(shell, "fish") {
		configDir := filepath.Join(homeDir, ".config", "fish")
		configFile := filepath.Join(configDir, "config.fish")
		return Fish, configFile, nil
	} else if strings.Contains(shell, "pwsh") || strings.Contains(shell, "powershell") {
		// PowerShell profile path varies by OS
		if err := os.MkdirAll(filepath.Join(homeDir, "Documents", "PowerShell"), 0755); err == nil {
			configFile := filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
			return PowerShell, configFile, nil
		}
	}

	return "", "", fmt.Errorf("unsupported shell: %s", shell)
}

// installZshHook installs the Zsh preexec hook
func (h *HookInstaller) installZshHook() error {
	return h.installZshHookWithForce(false)
}

// installBashHook installs the Bash DEBUG trap hook
func (h *HookInstaller) installBashHook() error {
	return h.installBashHookWithForce(false)
}

// installFishHook installs the Fish shell hook
func (h *HookInstaller) installFishHook() error {
	return h.installFishHookWithForce(false)
}

// installPowerShellHook installs the PowerShell hook
func (h *HookInstaller) installPowerShellHook() error {
	return h.installPowerShellHookWithForce(false)
}

// uninstallZshHook removes the Zsh hook
func (h *HookInstaller) uninstallZshHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "fi")
}

// uninstallBashHook removes the Bash hook
func (h *HookInstaller) uninstallBashHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "trap 'termonaut_log_command' DEBUG")
}

// uninstallFishHook removes the Fish hook
func (h *HookInstaller) uninstallFishHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "end")
}

// uninstallPowerShellHook removes the PowerShell hook
func (h *HookInstaller) uninstallPowerShellHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "}")
}

// appendToConfigFile appends content to the shell config file
func (h *HookInstaller) appendToConfigFile(content string) error {
	return h.appendToConfigFileWithForce(content, false)
}

// appendToConfigFileWithForce appends content to the shell config file with force option
func (h *HookInstaller) appendToConfigFileWithForce(content string, force bool) error {
	// Check if already installed (unless force is used)
	if !force {
		installed, err := h.IsInstalled()
		if err != nil {
			return fmt.Errorf("failed to check if hook is installed: %w", err)
		}
		if installed {
			return fmt.Errorf("termonaut hook is already installed")
		}
	}

	// Create config file if it doesn't exist
	if _, err := os.Stat(h.configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(h.configFile), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		if _, err := os.Create(h.configFile); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// Append hook content
	file, err := os.OpenFile(h.configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write hook to config file: %w", err)
	}

	return nil
}

// removeFromConfigFile removes content between start and end markers
func (h *HookInstaller) removeFromConfigFile(startMarker, endMarker string) error {
	content, err := os.ReadFile(h.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	var inBlock bool

	for _, line := range lines {
		if strings.Contains(line, startMarker) {
			inBlock = true
			continue
		}
		if inBlock && strings.Contains(line, endMarker) {
			inBlock = false
			continue
		}
		if !inBlock {
			newLines = append(newLines, line)
		}
	}

	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(h.configFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	return nil
}

// GetTerminalPID returns the current terminal process ID
func GetTerminalPID() int {
	return os.Getppid()
}

// GetCurrentWorkingDir returns the current working directory
func GetCurrentWorkingDir() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return ""
}

// GetBinaryPath returns the path to the termonaut binary
func GetBinaryPath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		// Fallback: try to find in PATH
		if path, err := exec.LookPath("termonaut"); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("failed to locate termonaut binary: %w", err)
	}
	return executable, nil
}

// ValidateShellInstallation checks if the shell integration is working
func ValidateShellInstallation() error {
	installer, err := NewHookInstaller("")
	if err != nil {
		return fmt.Errorf("failed to create hook installer: %w", err)
	}

	installed, err := installer.IsInstalled()
	if err != nil {
		return fmt.Errorf("failed to check installation status: %w", err)
	}

	if !installed {
		return fmt.Errorf("termonaut hook is not installed")
	}

	return nil
}

// createShortcutSymlink creates a symbolic link to the termonaut binary
func (h *HookInstaller) createShortcutSymlink() error {
	// Try to find a suitable directory in PATH for the symlink
	pathDirs := strings.Split(os.Getenv("PATH"), ":")
	var targetDir string

	// Prefer /usr/local/bin if it exists and is writable
	for _, dir := range pathDirs {
		if dir == "/usr/local/bin" {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				targetDir = dir
				break
			}
		}
	}

	// If /usr/local/bin is not available, try other directories
	if targetDir == "" {
		for _, dir := range pathDirs {
			if strings.Contains(dir, "bin") && !strings.Contains(dir, "sbin") {
				if info, err := os.Stat(dir); err == nil && info.IsDir() {
					targetDir = dir
					break
				}
			}
		}
	}

	if targetDir == "" {
		return fmt.Errorf("no suitable directory found in PATH")
	}

	linkPath := filepath.Join(targetDir, "tn")

	// Check if the symlink already exists
	if _, err := os.Lstat(linkPath); err == nil {
		return nil // Already exists
	}

	// Create the symlink
	if err := os.Symlink(h.binaryPath, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink at %s: %w", linkPath, err)
	}

	fmt.Printf("✅ Created 'tn' shortcut at %s\n", linkPath)
	return nil
}

// installZshHookWithForce installs the Zsh preexec hook with force option
func (h *HookInstaller) installZshHookWithForce(force bool) error {
	hook := fmt.Sprintf(`
# Termonaut shell integration (v0.9.0 Stable)
termonaut_preexec() {
    # Complete job control suppression - eliminate ALL job messages
    {
        # Disable job control notifications globally
        setopt NO_NOTIFY 2>/dev/null || true
        setopt NO_HUP 2>/dev/null || true
        setopt NO_BG_NICE 2>/dev/null || true
        setopt NO_CHECK_JOBS 2>/dev/null || true

        # Method 1: Use subshell with complete isolation
        (
            # Run in completely isolated subshell
            %s log-command "$1" >/dev/null 2>&1 &
            disown %%%% 2>/dev/null || true
        ) >/dev/null 2>&1 &

        # Method 2: Disown the subshell itself
        disown %%%% 2>/dev/null || true

        # Method 3: Clear job table
        jobs >/dev/null 2>&1 | while read job; do
            disown "$job" 2>/dev/null || true
        done

        # Method 4: Reset job control state
        setopt NO_MONITOR 2>/dev/null || true

    } >/dev/null 2>&1
}

# Check if preexec_functions exists, if not create it
if [[ -z "${preexec_functions+x}" ]]; then
    preexec_functions=()
fi

# Add our function to preexec_functions if not already present
if [[ ! " ${preexec_functions[@]} " =~ " termonaut_preexec " ]]; then
    preexec_functions+=(termonaut_preexec)
fi
`, h.binaryPath)

	return h.appendToConfigFileWithForce(hook, force)
}

// installBashHookWithForce installs the Bash DEBUG trap hook with force option
func (h *HookInstaller) installBashHookWithForce(force bool) error {
	hook := fmt.Sprintf(`
# Termonaut shell integration (v0.9.0 Stable)
termonaut_log_command() {
    if [ -n "$BASH_COMMAND" ]; then
        # Complete job control suppression - eliminate ALL job messages
        {
            # Disable job control globally
            set +m 2>/dev/null || true
            set +b 2>/dev/null || true

            # Method 1: Use exec with complete redirection
            (
                exec %s log-command "$BASH_COMMAND" >/dev/null 2>&1 &
            ) 2>/dev/null &

            # Method 2: Disown all background jobs
            disown $! 2>/dev/null || true
            jobs | awk '{print $1}' | while read job; do
                disown "$job" 2>/dev/null || true
            done

        } 2>/dev/null
    fi
}

# Set up DEBUG trap
trap 'termonaut_log_command' DEBUG
`, h.binaryPath)

	return h.appendToConfigFileWithForce(hook, force)
}

// installFishHookWithForce installs the Fish shell hook with force option
func (h *HookInstaller) installFishHookWithForce(force bool) error {
	hook := fmt.Sprintf(`
# Termonaut shell integration
function termonaut_preexec --on-event fish_preexec
    %s log-command "$argv" >/dev/null 2>&1 &
    disown
end
`, h.binaryPath)

	return h.appendToConfigFileWithForce(hook, force)
}

// installPowerShellHookWithForce installs the PowerShell hook with force option
func (h *HookInstaller) installPowerShellHookWithForce(force bool) error {
	hook := fmt.Sprintf(`
# Termonaut shell integration
function Invoke-TermonautLogging {
    param($Command)
    try {
        Start-Job -ScriptBlock {
            param($BinaryPath, $Cmd)
            & $BinaryPath log-command $Cmd 2>$null
        } -ArgumentList "%s", $Command | Out-Null
    } catch {
        # Silently ignore errors
    }
}

# PowerShell command history hook
$PSDefaultParameterValues['*:Verbose'] = $false
$PSDefaultParameterValues['*:Debug'] = $false

# Override the prompt to capture commands
function global:prompt {
    $history = Get-History -Count 1 -ErrorAction SilentlyContinue
    if ($history -and $history.CommandLine) {
        Invoke-TermonautLogging -Command $history.CommandLine
    }

    # Return original prompt
    "PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) "
}
`, h.binaryPath)

	return h.appendToConfigFileWithForce(hook, force)
}

// CreateBackup creates a timestamped backup of the configuration file
func (scm *SafeConfigManager) CreateBackup() (*ConfigBackup, error) {
	if _, err := os.Stat(scm.configFile); os.IsNotExist(err) {
		// File doesn't exist, no backup needed
		return nil, nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.termonaut_backup_%s", scm.configFile, timestamp)

	// Read original file
	content, err := os.ReadFile(scm.configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file for backup: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	return &ConfigBackup{
		OriginalPath: scm.configFile,
		BackupPath:   backupPath,
		Timestamp:    time.Now(),
	}, nil
}

// RestoreFromBackup restores the configuration file from backup
func (backup *ConfigBackup) RestoreFromBackup() error {
	if backup == nil {
		return fmt.Errorf("no backup to restore from")
	}

	content, err := os.ReadFile(backup.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	if err := os.WriteFile(backup.OriginalPath, content, 0644); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	return nil
}

// CleanupBackup removes the backup file
func (backup *ConfigBackup) CleanupBackup() error {
	if backup == nil {
		return nil
	}
	return os.Remove(backup.BackupPath)
}

// GetTermonautBlock extracts the Termonaut configuration block from the file
func (scm *SafeConfigManager) GetTermonautBlock() ([]string, int, int, error) {
	content, err := os.ReadFile(scm.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, -1, -1, nil // File doesn't exist
		}
		return nil, -1, -1, fmt.Errorf("failed to read config file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var blockLines []string
	startIdx := -1
	endIdx := -1

	// Define precise markers for different shells
	var startMarker, endMarker string
	switch scm.shellType {
	case Zsh:
		startMarker = "# Termonaut shell integration"
		endMarker = "fi"
	case Bash:
		startMarker = "# Termonaut shell integration"
		endMarker = "trap 'termonaut_log_command' DEBUG"
	case Fish:
		startMarker = "# Termonaut shell integration"
		endMarker = "end"
	case PowerShell:
		startMarker = "# Termonaut shell integration"
		endMarker = "}"
	default:
		return nil, -1, -1, fmt.Errorf("unsupported shell type: %s", scm.shellType)
	}

	inBlock := false
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.Contains(trimmedLine, startMarker) && !inBlock {
			inBlock = true
			startIdx = i
			blockLines = append(blockLines, line)
			continue
		}

		if inBlock {
			blockLines = append(blockLines, line)

			// More precise end detection
			if scm.isBlockEndLine(trimmedLine, endMarker) {
				endIdx = i
				break
			}
		}
	}

	if startIdx != -1 && endIdx == -1 {
		// Block started but never ended - this is a problem
		return blockLines, startIdx, -1, fmt.Errorf("incomplete Termonaut block found")
	}

	return blockLines, startIdx, endIdx, nil
}

// isBlockEndLine determines if a line marks the end of the Termonaut block
func (scm *SafeConfigManager) isBlockEndLine(line, endMarker string) bool {
	switch scm.shellType {
	case Zsh:
		// For Zsh, look for the closing 'fi' of our conditional block
		return strings.HasPrefix(line, "fi") &&
			!strings.Contains(line, "if") // Make sure it's not "if...fi" on same line
	case Bash:
		// For Bash, look for the trap command
		return strings.Contains(line, "trap") &&
			strings.Contains(line, "termonaut_log_command") &&
			strings.Contains(line, "DEBUG")
	case Fish:
		// For Fish, look for 'end'
		return line == "end"
	case PowerShell:
		// For PowerShell, look for closing brace
		return line == "}"
	default:
		return strings.Contains(line, endMarker)
	}
}

// RemoveTermonautBlock safely removes the Termonaut configuration block
func (scm *SafeConfigManager) RemoveTermonautBlock() error {
	// Create backup first
	backup, err := scm.CreateBackup()
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// If we created a backup, set up cleanup
	if backup != nil {
		defer func() {
			// Only cleanup backup on success
			if err == nil {
				backup.CleanupBackup()
			}
		}()
	}

	content, err := os.ReadFile(scm.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to remove
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	_, startIdx, endIdx, err := scm.GetTermonautBlock()
	if err != nil {
		// If there's an incomplete block, try to recover
		if strings.Contains(err.Error(), "incomplete") && backup != nil {
			return fmt.Errorf("incomplete Termonaut block detected, backup created at %s: %w", backup.BackupPath, err)
		}
		return err
	}

	// If no block found, nothing to remove
	if startIdx == -1 {
		return nil
	}

	// Remove the block
	var newLines []string
	newLines = append(newLines, lines[:startIdx]...)
	if endIdx != -1 && endIdx+1 < len(lines) {
		newLines = append(newLines, lines[endIdx+1:]...)
	}

	// Clean up empty lines around the removed block
	newLines = scm.cleanupEmptyLines(newLines, startIdx)

	// Write the new content atomically
	newContent := strings.Join(newLines, "\n")
	if err := scm.writeConfigFileAtomically(newContent); err != nil {
		// Restore from backup on failure
		if backup != nil {
			if restoreErr := backup.RestoreFromBackup(); restoreErr != nil {
				return fmt.Errorf("failed to write config and restore backup: write error: %w, restore error: %v", err, restoreErr)
			}
		}
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	// Validate the syntax of the modified file
	if err := scm.validateShellSyntax(); err != nil {
		// Restore from backup on syntax error
		if backup != nil {
			if restoreErr := backup.RestoreFromBackup(); restoreErr != nil {
				return fmt.Errorf("syntax validation failed and restore failed: syntax error: %w, restore error: %v", err, restoreErr)
			}
		}
		return fmt.Errorf("syntax validation failed, restored from backup: %w", err)
	}

	return nil
}

// AddTermonautBlock safely adds the Termonaut configuration block
func (scm *SafeConfigManager) AddTermonautBlock(content string) error {
	// Create backup first
	backup, err := scm.CreateBackup()
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// If we created a backup, set up cleanup
	if backup != nil {
		defer func() {
			// Only cleanup backup on success
			if err == nil {
				backup.CleanupBackup()
			}
		}()
	}

	// Check if block already exists
	_, startIdx, _, err := scm.GetTermonautBlock()
	if err != nil && !strings.Contains(err.Error(), "incomplete") {
		return err
	}
	if startIdx != -1 {
		return fmt.Errorf("Termonaut block already exists")
	}

	// Create config file if it doesn't exist
	if _, err := os.Stat(scm.configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(scm.configFile), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create with basic shell setup
		initialContent := scm.getInitialConfigContent()
		if err := os.WriteFile(scm.configFile, []byte(initialContent), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// Read current content
	currentContent, err := os.ReadFile(scm.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Add the new block with proper spacing
	newContent := string(currentContent)
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += "\n" + content + "\n"

	// Write atomically
	if err := scm.writeConfigFileAtomically(newContent); err != nil {
		// Restore from backup on failure
		if backup != nil {
			if restoreErr := backup.RestoreFromBackup(); restoreErr != nil {
				return fmt.Errorf("failed to write config and restore backup: write error: %w, restore error: %v", err, restoreErr)
			}
		}
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	// Validate the syntax of the modified file
	if err := scm.validateShellSyntax(); err != nil {
		// Restore from backup on syntax error
		if backup != nil {
			if restoreErr := backup.RestoreFromBackup(); restoreErr != nil {
				return fmt.Errorf("syntax validation failed and restore failed: syntax error: %w, restore error: %v", err, restoreErr)
			}
		}
		return fmt.Errorf("syntax validation failed, restored from backup: %w", err)
	}

	return nil
}

// writeConfigFileAtomically writes content to the config file atomically
func (scm *SafeConfigManager) writeConfigFileAtomically(content string) error {
	// Write to temporary file first
	tempFile := scm.configFile + ".tmp"
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, scm.configFile); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// validateShellSyntax validates the syntax of the shell configuration file
func (scm *SafeConfigManager) validateShellSyntax() error {
	var cmd *exec.Cmd

	switch scm.shellType {
	case Zsh:
		// Use zsh -n to check syntax without execution
		cmd = exec.Command("zsh", "-n", scm.configFile)
	case Bash:
		// Use bash -n to check syntax without execution
		cmd = exec.Command("bash", "-n", scm.configFile)
	case Fish:
		// Fish has --parse-only option
		cmd = exec.Command("fish", "--parse-only", scm.configFile)
	case PowerShell:
		// PowerShell syntax check
		cmd = exec.Command("pwsh", "-NoProfile", "-Command", fmt.Sprintf("Get-Content '%s' | Out-Null", scm.configFile))
	default:
		// Skip validation for unknown shells
		return nil
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("syntax validation failed: %s", string(output))
	}

	return nil
}

// cleanupEmptyLines removes excessive empty lines around the modification point
func (scm *SafeConfigManager) cleanupEmptyLines(lines []string, modifyIdx int) []string {
	// Remove excessive empty lines before and after the modification point
	result := make([]string, 0, len(lines))

	for i, line := range lines {
		// Keep the line if it's not empty
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
			continue
		}

		// For empty lines, be more selective
		prevNonEmpty := -1
		nextNonEmpty := -1

		// Find previous non-empty line
		for j := i - 1; j >= 0; j-- {
			if strings.TrimSpace(lines[j]) != "" {
				prevNonEmpty = j
				break
			}
		}

		// Find next non-empty line
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) != "" {
				nextNonEmpty = j
				break
			}
		}

		// Keep at most one empty line between non-empty content
		if prevNonEmpty != -1 && nextNonEmpty != -1 {
			// Check if this is the first empty line in a sequence
			isFirstEmpty := i == 0 || strings.TrimSpace(lines[i-1]) != ""
			if isFirstEmpty {
				result = append(result, line)
			}
		} else if prevNonEmpty != -1 || nextNonEmpty != -1 {
			// Keep empty lines at the beginning or end
			result = append(result, line)
		}
	}

	return result
}

// getInitialConfigContent returns initial content for a new config file
func (scm *SafeConfigManager) getInitialConfigContent() string {
	switch scm.shellType {
	case Zsh:
		return "# Zsh configuration\n"
	case Bash:
		return "# Bash configuration\n"
	case Fish:
		return "# Fish configuration\n"
	case PowerShell:
		return "# PowerShell configuration\n"
	default:
		return "# Shell configuration\n"
	}
}
