package privacy

import (
	"regexp"
	"strings"
)

// CommandSanitizer handles command sanitization for privacy protection
type CommandSanitizer struct {
	config *SanitizationConfig
}

// SanitizationConfig holds configuration for command sanitization
type SanitizationConfig struct {
	Enabled            bool     `json:"enabled"`
	SanitizePasswords  bool     `json:"sanitize_passwords"`
	SanitizeURLs       bool     `json:"sanitize_urls"`
	SanitizeFilePaths  bool     `json:"sanitize_file_paths"`
	SanitizeTokens     bool     `json:"sanitize_tokens"`
	SanitizeEmails     bool     `json:"sanitize_emails"`
	IgnoreCommands     []string `json:"ignore_commands"`     // Commands to completely ignore
	SensitivePatterns  []string `json:"sensitive_patterns"`  // Custom patterns to sanitize
	PreservePrefixes   []string `json:"preserve_prefixes"`   // Prefixes to keep (like git, npm)
	MaxArgLength       int      `json:"max_arg_length"`      // Max length for arguments
}

// NewCommandSanitizer creates a new command sanitizer
func NewCommandSanitizer(config *SanitizationConfig) *CommandSanitizer {
	if config == nil {
		config = DefaultSanitizationConfig()
	}
	return &CommandSanitizer{config: config}
}

// DefaultSanitizationConfig returns default sanitization configuration
func DefaultSanitizationConfig() *SanitizationConfig {
	return &SanitizationConfig{
		Enabled:           true,
		SanitizePasswords: true,
		SanitizeURLs:      true,
		SanitizeFilePaths: true,
		SanitizeTokens:    true,
		SanitizeEmails:    true,
		IgnoreCommands: []string{
			"sudo", "su", "passwd", "ssh-keygen", "gpg", 
			"openssl", "keychain", "security",
		},
		SensitivePatterns: []string{
			`password\s*[=:]\s*\S+`,
			`token\s*[=:]\s*\S+`,
			`key\s*[=:]\s*\S+`,
			`secret\s*[=:]\s*\S+`,
			`auth\s*[=:]\s*\S+`,
		},
		PreservePrefixes: []string{
			"git", "npm", "pip", "docker", "kubectl", "yarn", "cargo",
			"go", "python", "node", "mvn", "gradle",
		},
		MaxArgLength: 50,
	}
}

// SanitizeCommand sanitizes a command for privacy protection
func (cs *CommandSanitizer) SanitizeCommand(command string) (string, bool) {
	if !cs.config.Enabled {
		return command, false
	}

	originalCommand := command
	command = strings.TrimSpace(command)
	
	if command == "" {
		return "", true
	}

	// Check if command should be completely ignored
	if cs.shouldIgnoreCommand(command) {
		return "", true // Return empty string to ignore this command
	}

	// Parse command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", true
	}

	baseCommand := parts[0]
	args := parts[1:]

	// Preserve known safe prefixes
	if cs.isPreservedPrefix(baseCommand) {
		sanitizedArgs := cs.sanitizeArguments(args, baseCommand)
		return baseCommand + " " + strings.Join(sanitizedArgs, " "), false
	}

	// Apply general sanitization
	sanitizedCommand := cs.applySanitization(command)
	
	return sanitizedCommand, sanitizedCommand != originalCommand
}

// shouldIgnoreCommand checks if a command should be completely ignored
func (cs *CommandSanitizer) shouldIgnoreCommand(command string) bool {
	for _, ignored := range cs.config.IgnoreCommands {
		if strings.HasPrefix(strings.ToLower(command), strings.ToLower(ignored)) {
			return true
		}
	}
	return false
}

// isPreservedPrefix checks if a command starts with a preserved prefix
func (cs *CommandSanitizer) isPreservedPrefix(baseCommand string) bool {
	for _, prefix := range cs.config.PreservePrefixes {
		if strings.EqualFold(baseCommand, prefix) {
			return true
		}
	}
	return false
}

// sanitizeArguments sanitizes command arguments based on command type
func (cs *CommandSanitizer) sanitizeArguments(args []string, baseCommand string) []string {
	var sanitizedArgs []string
	
	for _, arg := range args {
		sanitizedArg := arg
		
		// Skip flags and options (starting with -)
		if strings.HasPrefix(arg, "-") {
			sanitizedArgs = append(sanitizedArgs, arg)
			continue
		}
		
		// Apply sanitization based on context
		if cs.config.SanitizeURLs && cs.isURL(arg) {
			sanitizedArg = cs.sanitizeURL(arg)
		} else if cs.config.SanitizeFilePaths && cs.isFilePath(arg) {
			sanitizedArg = cs.sanitizeFilePath(arg)
		} else if cs.config.SanitizeEmails && cs.isEmail(arg) {
			sanitizedArg = "[EMAIL]"
		} else if cs.config.SanitizeTokens && cs.isToken(arg) {
			sanitizedArg = "[TOKEN]"
		} else {
			// Apply general pattern sanitization
			sanitizedArg = cs.sanitizePatterns(arg)
		}
		
		// Truncate long arguments
		if len(sanitizedArg) > cs.config.MaxArgLength {
			sanitizedArg = sanitizedArg[:cs.config.MaxArgLength-3] + "..."
		}
		
		sanitizedArgs = append(sanitizedArgs, sanitizedArg)
	}
	
	return sanitizedArgs
}

// applySanitization applies general sanitization patterns
func (cs *CommandSanitizer) applySanitization(command string) string {
	sanitized := command
	
	if cs.config.SanitizePasswords {
		sanitized = cs.sanitizePasswords(sanitized)
	}
	
	if cs.config.SanitizeTokens {
		sanitized = cs.sanitizeTokens(sanitized)
	}
	
	// Apply custom sensitive patterns
	for _, pattern := range cs.config.SensitivePatterns {
		if re, err := regexp.Compile(`(?i)`+pattern); err == nil {
			sanitized = re.ReplaceAllString(sanitized, "[REDACTED]")
		}
	}
	
	return sanitized
}

// Helper methods for detection and sanitization
func (cs *CommandSanitizer) isURL(arg string) bool {
	return strings.HasPrefix(arg, "http://") || 
		   strings.HasPrefix(arg, "https://") || 
		   strings.HasPrefix(arg, "ftp://") ||
		   strings.Contains(arg, "://")
}

func (cs *CommandSanitizer) sanitizeURL(url string) string {
	// Keep protocol and domain, hide path and parameters
	re := regexp.MustCompile(`^(https?://[^/]+).*`)
	if matches := re.FindStringSubmatch(url); len(matches) > 1 {
		return matches[1] + "/[PATH]"
	}
	return "[URL]"
}

func (cs *CommandSanitizer) isFilePath(arg string) bool {
	return strings.Contains(arg, "/") || 
		   strings.Contains(arg, "\\") ||
		   strings.HasPrefix(arg, "~") ||
		   strings.HasPrefix(arg, ".")
}

func (cs *CommandSanitizer) sanitizeFilePath(path string) string {
	// Keep just the filename or directory name
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return "[PATH]/" + parts[len(parts)-1]
	}
	return path
}

func (cs *CommandSanitizer) isEmail(arg string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(arg)
}

func (cs *CommandSanitizer) isToken(arg string) bool {
	// Detect common token patterns
	if len(arg) > 20 && (strings.Contains(arg, "_") || 
		regexp.MustCompile(`^[a-zA-Z0-9+/=]{20,}$`).MatchString(arg)) {
		return true
	}
	// Common prefixes for tokens
	prefixes := []string{"ghp_", "gho_", "ghu_", "ghs_", "sk-", "pk_", "rk_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}

func (cs *CommandSanitizer) sanitizePasswords(command string) string {
	patterns := []string{
		`password\s*[=:]\s*\S+`,
		`passwd\s+\S+`,
		`-p\s+\S+`,
		`--password\s+\S+`,
	}
	
	result := command
	for _, pattern := range patterns {
		if re, err := regexp.Compile(`(?i)`+pattern); err == nil {
			result = re.ReplaceAllString(result, "[PASSWORD]")
		}
	}
	return result
}

func (cs *CommandSanitizer) sanitizeTokens(command string) string {
	patterns := []string{
		`token\s*[=:]\s*\S+`,
		`--token\s+\S+`,
		`-t\s+\S+`,
		`bearer\s+\S+`,
	}
	
	result := command
	for _, pattern := range patterns {
		if re, err := regexp.Compile(`(?i)`+pattern); err == nil {
			result = re.ReplaceAllString(result, "[TOKEN]")
		}
	}
	return result
}

func (cs *CommandSanitizer) sanitizePatterns(arg string) string {
	// Apply custom sensitive patterns
	result := arg
	for _, pattern := range cs.config.SensitivePatterns {
		if re, err := regexp.Compile(`(?i)`+pattern); err == nil {
			if re.MatchString(arg) {
				return "[REDACTED]"
			}
		}
	}
	return result
} 