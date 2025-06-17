package environment

import (
	"os"
	"strings"
)

// CIEnvironment represents different CI/CD platforms
type CIEnvironment string

const (
	// CI Platforms
	CIUnknown        CIEnvironment = "unknown"
	CIGitHubActions  CIEnvironment = "github_actions"
	CIGitLabCI       CIEnvironment = "gitlab_ci"
	CIJenkins        CIEnvironment = "jenkins"
	CITravisCI       CIEnvironment = "travis_ci"
	CICircleCI       CIEnvironment = "circle_ci"
	CIBuildkite      CIEnvironment = "buildkite"
	CITeamCity       CIEnvironment = "teamcity"
	CIBamboo         CIEnvironment = "bamboo"
	CIAppVeyor       CIEnvironment = "appveyor"
	CIAzurePipelines CIEnvironment = "azure_pipelines"
	CICodeBuild      CIEnvironment = "codebuild"
	CIDrone          CIEnvironment = "drone"
	CISemaphore      CIEnvironment = "semaphore"
	CIGeneric        CIEnvironment = "generic_ci"
)

// EnvironmentInfo contains information about the current environment
type EnvironmentInfo struct {
	IsCI                 bool
	CIEnvironment        CIEnvironment
	SupportColors        bool
	SupportInteractive   bool
	SupportEmojis        bool
	TerminalType         string
	ShouldSuppressOutput bool
}

// Detector handles environment detection
type Detector struct{}

// NewDetector creates a new environment detector
func NewDetector() *Detector {
	return &Detector{}
}

// DetectEnvironment detects the current environment and its capabilities
func (d *Detector) DetectEnvironment() *EnvironmentInfo {
	info := &EnvironmentInfo{
		IsCI:                 d.detectCI(),
		CIEnvironment:        d.detectCIEnvironment(),
		SupportColors:        d.detectColorSupport(),
		SupportInteractive:   d.detectInteractiveSupport(),
		SupportEmojis:        d.detectEmojiSupport(),
		TerminalType:         d.detectTerminalType(),
		ShouldSuppressOutput: d.shouldSuppressOutput(),
	}

	return info
}

// detectCI checks if running in a CI environment
func (d *Detector) detectCI() bool {
	// Check common CI environment variables
	ciIndicators := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"BUILD_NUMBER",
		"CI_SERVER",
		"TF_BUILD", // Azure Pipelines
	}

	for _, indicator := range ciIndicators {
		if value := os.Getenv(indicator); value != "" &&
			(strings.ToLower(value) == "true" || strings.ToLower(value) == "1") {
			return true
		}
	}

	// Check for specific CI platform variables
	ciPlatformVars := []string{
		"GITHUB_ACTIONS",
		"GITLAB_CI",
		"JENKINS_URL",
		"TRAVIS",
		"CIRCLECI",
		"BUILDKITE",
		"TEAMCITY_VERSION",
		"bamboo_buildNumber",
		"APPVEYOR",
		"TF_BUILD",
		"CODEBUILD_BUILD_ID",
		"DRONE",
		"SEMAPHORE",
	}

	for _, platformVar := range ciPlatformVars {
		if os.Getenv(platformVar) != "" {
			return true
		}
	}

	return false
}

// detectCIEnvironment identifies the specific CI platform
func (d *Detector) detectCIEnvironment() CIEnvironment {
	// GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return CIGitHubActions
	}

	// GitLab CI
	if os.Getenv("GITLAB_CI") == "true" {
		return CIGitLabCI
	}

	// Jenkins
	if os.Getenv("JENKINS_URL") != "" {
		return CIJenkins
	}

	// Travis CI
	if os.Getenv("TRAVIS") == "true" {
		return CITravisCI
	}

	// CircleCI
	if os.Getenv("CIRCLECI") == "true" {
		return CICircleCI
	}

	// Buildkite
	if os.Getenv("BUILDKITE") == "true" {
		return CIBuildkite
	}

	// TeamCity
	if os.Getenv("TEAMCITY_VERSION") != "" {
		return CITeamCity
	}

	// Bamboo
	if os.Getenv("bamboo_buildNumber") != "" {
		return CIBamboo
	}

	// AppVeyor
	if os.Getenv("APPVEYOR") == "True" {
		return CIAppVeyor
	}

	// Azure Pipelines
	if os.Getenv("TF_BUILD") == "True" {
		return CIAzurePipelines
	}

	// AWS CodeBuild
	if os.Getenv("CODEBUILD_BUILD_ID") != "" {
		return CICodeBuild
	}

	// Drone CI
	if os.Getenv("DRONE") == "true" {
		return CIDrone
	}

	// Semaphore
	if os.Getenv("SEMAPHORE") == "true" {
		return CISemaphore
	}

	// Generic CI detection
	if d.detectCI() {
		return CIGeneric
	}

	return CIUnknown
}

// detectColorSupport checks if the terminal supports colors
func (d *Detector) detectColorSupport() bool {
	// If in CI, check CI-specific color support
	if d.detectCI() {
		return d.detectCIColorSupport()
	}

	term := strings.ToLower(os.Getenv("TERM"))

	// Check for color-capable terminals
	colorTerms := []string{
		"xterm",
		"xterm-color",
		"xterm-256color",
		"screen",
		"screen-256color",
		"tmux",
		"tmux-256color",
		"rxvt",
		"konsole",
		"ansi",
	}

	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	// Check COLORTERM environment variable
	if os.Getenv("COLORTERM") != "" {
		return true
	}

	// Check if TERM=dumb (explicitly no color support)
	if term == "dumb" {
		return false
	}

	return true // Default to supporting colors
}

// detectCIColorSupport checks color support in CI environments
func (d *Detector) detectCIColorSupport() bool {
	ciEnv := d.detectCIEnvironment()

	switch ciEnv {
	case CIGitHubActions:
		// GitHub Actions supports colors
		return true
	case CIGitLabCI:
		// GitLab CI supports colors
		return true
	case CITravisCI:
		// Travis CI supports colors
		return true
	case CICircleCI:
		// CircleCI supports colors
		return true
	case CIAzurePipelines:
		// Azure Pipelines supports colors
		return true
	default:
		// Conservative approach for unknown CI
		return false
	}
}

// detectInteractiveSupport checks if the environment supports interactive features
func (d *Detector) detectInteractiveSupport() bool {
	// CI environments are typically non-interactive
	if d.detectCI() {
		return false
	}

	// Check if stdin is a terminal
	if fi, err := os.Stdin.Stat(); err == nil {
		return (fi.Mode() & os.ModeCharDevice) != 0
	}

	return true // Default to interactive
}

// detectEmojiSupport checks if the environment supports emoji display
func (d *Detector) detectEmojiSupport() bool {
	// CI environments typically don't need emojis
	if d.detectCI() {
		return false
	}

	// Check terminal capabilities
	term := strings.ToLower(os.Getenv("TERM"))

	// Modern terminals typically support emojis
	modernTerms := []string{
		"xterm-256color",
		"screen-256color",
		"tmux-256color",
	}

	for _, modernTerm := range modernTerms {
		if strings.Contains(term, modernTerm) {
			return true
		}
	}

	// Check for specific terminal emulators known to support emojis
	if os.Getenv("ITERM_PROFILE") != "" || // iTerm2
		os.Getenv("TERM_PROGRAM") == "iTerm.app" ||
		os.Getenv("TERM_PROGRAM") == "vscode" || // VS Code terminal
		os.Getenv("WT_SESSION") != "" { // Windows Terminal
		return true
	}

	// Default to supporting emojis for interactive environments
	return d.detectInteractiveSupport()
}

// detectTerminalType identifies the terminal type
func (d *Detector) detectTerminalType() string {
	if term := os.Getenv("TERM"); term != "" {
		return term
	}

	if termProgram := os.Getenv("TERM_PROGRAM"); termProgram != "" {
		return termProgram
	}

	return "unknown"
}

// shouldSuppressOutput determines if output should be suppressed
func (d *Detector) shouldSuppressOutput() bool {
	// Suppress output in CI unless explicitly enabled
	if d.detectCI() {
		// Check if explicitly enabled in CI
		if os.Getenv("TERMONAUT_CI_VERBOSE") == "true" {
			return false
		}
		return true
	}

	// Check for explicit suppression
	if os.Getenv("TERMONAUT_QUIET") == "true" {
		return true
	}

	return false
}

// GetRecommendedDisplayMode returns the recommended display mode for the environment
func (d *Detector) GetRecommendedDisplayMode(info *EnvironmentInfo) string {
	if info.IsCI || info.ShouldSuppressOutput {
		return "quiet"
	}

	if info.SupportEmojis && info.SupportColors && info.SupportInteractive {
		return "rich"
	}

	return "minimal"
}

// IsCI returns true if running in a CI environment
func (d *Detector) IsCI() bool {
	return d.detectCI()
}

// ShouldBeQuiet returns true if output should be minimal/suppressed
func (d *Detector) ShouldBeQuiet() bool {
	return d.shouldSuppressOutput()
}
