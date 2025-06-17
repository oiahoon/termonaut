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
	case Fish:
		return h.installFishHook()
	case PowerShell:
		return h.installPowerShellHook()
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
	case Fish:
		return h.uninstallFishHook()
	case PowerShell:
		return h.uninstallPowerShellHook()
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
	hook := fmt.Sprintf(`
# Termonaut shell integration
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
            %s log-command "$1" &
            
            # Force exit to prevent any shell interaction
            exit 0
        ) &
        disown %%- 2>/dev/null || true
    } 2>/dev/null
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
        # Silent background execution with comprehensive job control suppression
        {
            # Create a completely detached subshell
            (
                # Disable all job control and output
                set +m 2>/dev/null
                unset HISTFILE 2>/dev/null
                exec </dev/null >/dev/null 2>&1
                
                # Run termonaut in completely isolated environment
                %s log-command "$BASH_COMMAND" &
                
                # Force exit to prevent any shell interaction
                exit 0
            ) &
            disown $! 2>/dev/null || true
        } 2>/dev/null
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

// installFishHook installs the Fish shell hook
func (h *HookInstaller) installFishHook() error {
	hook := fmt.Sprintf(`
# Termonaut shell integration
function termonaut_preexec --on-event fish_preexec
    %s log-command "$argv" >/dev/null 2>&1 &
    disown
end
`, h.binaryPath)

	return h.appendToConfigFile(hook)
}

// installPowerShellHook installs the PowerShell hook
func (h *HookInstaller) installPowerShellHook() error {
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

	return h.appendToConfigFile(hook)
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
