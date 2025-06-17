package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ShellType represents different shell types
type ShellType string

const (
	// Zsh shell
	Zsh ShellType = "zsh"
	// Bash shell
	Bash ShellType = "bash"
	// Fish shell (future support)
	Fish ShellType = "fish"
)

// HookInstaller handles shell hook installation
type HookInstaller struct {
	shellType  ShellType
	configFile string
	binaryPath string
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

// Install installs the shell hook
func (h *HookInstaller) Install() error {
	switch h.shellType {
	case Zsh:
		return h.installZshHook()
	case Bash:
		return h.installBashHook()
	default:
		return fmt.Errorf("unsupported shell: %s", h.shellType)
	}
}

// Uninstall removes the shell hook
func (h *HookInstaller) Uninstall() error {
	switch h.shellType {
	case Zsh:
		return h.uninstallZshHook()
	case Bash:
		return h.uninstallBashHook()
	default:
		return fmt.Errorf("unsupported shell: %s", h.shellType)
	}
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

	return strings.Contains(string(content), "termonaut"), nil
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
	}

	return "", "", fmt.Errorf("unsupported shell: %s", shell)
}

// installZshHook installs the Zsh preexec hook
func (h *HookInstaller) installZshHook() error {
	hook := fmt.Sprintf(`
# Termonaut shell integration
termonaut_preexec() {
    { %s log-command "$1" >/dev/null 2>&1 & } 2>/dev/null
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

	return h.appendToConfigFile(hook)
}

// installBashHook installs the Bash DEBUG trap hook
func (h *HookInstaller) installBashHook() error {
	hook := fmt.Sprintf(`
# Termonaut shell integration
termonaut_log_command() {
    if [ -n "$BASH_COMMAND" ]; then
        { %s log-command "$BASH_COMMAND" >/dev/null 2>&1 & } 2>/dev/null
    fi
}

# Set up DEBUG trap
trap 'termonaut_log_command' DEBUG
`, h.binaryPath)

	return h.appendToConfigFile(hook)
}

// uninstallZshHook removes the Zsh hook
func (h *HookInstaller) uninstallZshHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "fi")
}

// uninstallBashHook removes the Bash hook
func (h *HookInstaller) uninstallBashHook() error {
	return h.removeFromConfigFile("# Termonaut shell integration", "trap 'termonaut_log_command' DEBUG")
}

// appendToConfigFile appends content to the shell config file
func (h *HookInstaller) appendToConfigFile(content string) error {
	// Check if already installed
	installed, err := h.IsInstalled()
	if err != nil {
		return fmt.Errorf("failed to check if hook is installed: %w", err)
	}
	if installed {
		return fmt.Errorf("termonaut hook is already installed")
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
