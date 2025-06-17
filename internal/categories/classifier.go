package categories

import (
	"regexp"
	"strings"
)

// Category represents a command category
type Category string

const (
	Git         Category = "git"
	Development Category = "development"
	System      Category = "system"
	Navigation  Category = "navigation"
	Network     Category = "network"
	Docker      Category = "docker"
	Kubernetes  Category = "kubernetes"
	Cloud       Category = "cloud"
	Database    Category = "database"
	Package     Category = "package"
	Build       Category = "build"
	Test        Category = "test"
	Text        Category = "text"
	Archive     Category = "archive"
	Security    Category = "security"
	Monitoring  Category = "monitoring"
	Unknown     Category = "unknown"
)

// CategoryInfo holds category metadata
type CategoryInfo struct {
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	XPBonus     float64 `json:"xp_bonus"`
}

// CommandClassifier classifies commands into categories
type CommandClassifier struct {
	patterns map[Category][]*regexp.Regexp
	metadata map[Category]*CategoryInfo
}

// NewCommandClassifier creates a new command classifier
func NewCommandClassifier() *CommandClassifier {
	cc := &CommandClassifier{
		patterns: make(map[Category][]*regexp.Regexp),
		metadata: make(map[Category]*CategoryInfo),
	}
	cc.initializePatterns()
	cc.initializeMetadata()
	return cc
}

// initializePatterns sets up regex patterns for command classification
func (cc *CommandClassifier) initializePatterns() {
	// Git commands
	cc.addPatterns(Git, []string{
		`^git\s+`,
		`^gh\s+`,
		`^hub\s+`,
		`^gitk$`,
		`^tig\s+`,
	})

	// Development commands
	cc.addPatterns(Development, []string{
		`^(node|npm|yarn|pnpm)\s+`,
		`^(python|python3|py)\s+`,
		`^(go|gofmt|goimports)\s+`,
		`^(cargo|rustc)\s+`,
		`^(javac?|mvn|gradle)\s+`,
		`^(gcc|g\+\+|clang)\s+`,
		`^(php|composer)\s+`,
		`^(ruby|gem|bundle)\s+`,
		`^(dotnet|nuget)\s+`,
		`^(swift|swiftc)\s+`,
		`^(scala|sbt)\s+`,
		`^(kotlin|kotlinc)\s+`,
		`^vim?\s+`,
		`^(emacs|nano|code|subl)\s+`,
		`^(ipython|jupyter)\s+`,
	})

	// System commands
	cc.addPatterns(System, []string{
		`^(ls|ll|la|dir)\s*`,
		`^(ps|top|htop|btop)\s*`,
		`^(kill|killall|pkill)\s+`,
		`^(sudo|su)\s+`,
		`^(systemctl|service)\s+`,
		`^(mount|umount)\s+`,
		`^(df|du|free)\s*`,
		`^(uname|whoami|id)\s*`,
		`^(date|uptime|w)\s*`,
		`^(chmod|chown|chgrp)\s+`,
		`^(ln|cp|mv|rm|mkdir|rmdir)\s+`,
		`^(find|locate|which|whereis)\s+`,
		`^(crontab|at|jobs)\s*`,
		`^(env|export|set)\s*`,
	})

	// Navigation commands
	cc.addPatterns(Navigation, []string{
		`^cd\s*`,
		`^pwd\s*`,
		`^pushd\s+`,
		`^popd\s*`,
		`^dirs\s*`,
	})

	// Network commands
	cc.addPatterns(Network, []string{
		`^(curl|wget|http)\s+`,
		`^(ping|traceroute|nslookup|dig)\s+`,
		`^(ssh|scp|rsync|sftp)\s+`,
		`^(netstat|ss|lsof)\s*`,
		`^(tcpdump|wireshark)\s+`,
		`^(nc|netcat|telnet)\s+`,
		`^(iptables|ufw|firewall-cmd)\s+`,
	})

	// Docker commands
	cc.addPatterns(Docker, []string{
		`^docker\s+`,
		`^docker-compose\s+`,
		`^podman\s+`,
	})

	// Kubernetes commands
	cc.addPatterns(Kubernetes, []string{
		`^kubectl\s+`,
		`^k9s\s*`,
		`^helm\s+`,
		`^minikube\s+`,
		`^kind\s+`,
		`^k3s\s+`,
	})

	// Cloud commands
	cc.addPatterns(Cloud, []string{
		`^(aws|gcloud|az)\s+`,
		`^(terraform|terragrunt)\s+`,
		`^(ansible|ansible-playbook)\s+`,
		`^(pulumi|cdk)\s+`,
		`^(sam|serverless)\s+`,
	})

	// Database commands
	cc.addPatterns(Database, []string{
		`^(mysql|psql|sqlite3)\s+`,
		`^(redis-cli|mongo)\s+`,
		`^(pg_dump|mysqldump)\s+`,
	})

	// Package managers
	cc.addPatterns(Package, []string{
		`^(apt|yum|dnf|pacman|brew|port)\s+`,
		`^(pip|pip3|conda)\s+`,
		`^(snap|flatpak|appimage)\s+`,
	})

	// Build tools
	cc.addPatterns(Build, []string{
		`^(make|cmake|ninja)\s+`,
		`^(webpack|rollup|vite)\s+`,
		`^(gulp|grunt)\s+`,
		`^(bazel|buck)\s+`,
	})

	// Testing
	cc.addPatterns(Test, []string{
		`^(jest|mocha|pytest|rspec)\s+`,
		`^(phpunit|junit|nunit)\s+`,
		`^(karma|cypress|selenium)\s+`,
	})

	// Text processing
	cc.addPatterns(Text, []string{
		`^(grep|egrep|fgrep|rg|ag)\s+`,
		`^(sed|awk|cut|sort|uniq)\s+`,
		`^(head|tail|less|more|cat)\s+`,
		`^(tr|wc|diff|patch)\s+`,
		`^(jq|yq|xmllint)\s+`,
	})

	// Archive tools
	cc.addPatterns(Archive, []string{
		`^(tar|zip|unzip|gzip|gunzip)\s+`,
		`^(7z|rar|unrar)\s+`,
		`^(xz|bzip2|compress)\s+`,
	})

	// Security tools
	cc.addPatterns(Security, []string{
		`^(gpg|openssl|keytool)\s+`,
		`^(nmap|masscan|nikto)\s+`,
		`^(hashcat|john)\s+`,
		`^(vault|age|sops)\s+`,
	})

	// Monitoring tools
	cc.addPatterns(Monitoring, []string{
		`^(prometheus|grafana|influx)\s+`,
		`^(datadog|newrelic)\s+`,
		`^(nagios|zabbix)\s+`,
		`^(elastic|logstash|kibana)\s+`,
	})
}

// addPatterns compiles and adds regex patterns for a category
func (cc *CommandClassifier) addPatterns(category Category, patterns []string) {
	for _, pattern := range patterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			cc.patterns[category] = append(cc.patterns[category], regex)
		}
	}
}

// initializeMetadata sets up category metadata
func (cc *CommandClassifier) initializeMetadata() {
	cc.metadata[Git] = &CategoryInfo{
		Name:        "Git & Version Control",
		Icon:        "🌿",
		Color:       "green",
		Description: "Version control operations",
		XPBonus:     1.5,
	}

	cc.metadata[Development] = &CategoryInfo{
		Name:        "Development",
		Icon:        "💻",
		Color:       "blue",
		Description: "Programming and development tools",
		XPBonus:     1.3,
	}

	cc.metadata[System] = &CategoryInfo{
		Name:        "System Administration",
		Icon:        "⚙️",
		Color:       "yellow",
		Description: "System management and administration",
		XPBonus:     1.0,
	}

	cc.metadata[Navigation] = &CategoryInfo{
		Name:        "Navigation",
		Icon:        "🧭",
		Color:       "cyan",
		Description: "Directory navigation and exploration",
		XPBonus:     0.8,
	}

	cc.metadata[Network] = &CategoryInfo{
		Name:        "Network & Internet",
		Icon:        "🌐",
		Color:       "purple",
		Description: "Network operations and internet tools",
		XPBonus:     1.2,
	}

	cc.metadata[Docker] = &CategoryInfo{
		Name:        "Docker & Containers",
		Icon:        "🐳",
		Color:       "blue",
		Description: "Container management and orchestration",
		XPBonus:     1.4,
	}

	cc.metadata[Kubernetes] = &CategoryInfo{
		Name:        "Kubernetes",
		Icon:        "☸️",
		Color:       "blue",
		Description: "Kubernetes cluster management",
		XPBonus:     1.6,
	}

	cc.metadata[Cloud] = &CategoryInfo{
		Name:        "Cloud & Infrastructure",
		Icon:        "☁️",
		Color:       "white",
		Description: "Cloud services and infrastructure tools",
		XPBonus:     1.5,
	}

	cc.metadata[Database] = &CategoryInfo{
		Name:        "Database",
		Icon:        "🗄️",
		Color:       "orange",
		Description: "Database operations and management",
		XPBonus:     1.3,
	}

	cc.metadata[Package] = &CategoryInfo{
		Name:        "Package Management",
		Icon:        "📦",
		Color:       "brown",
		Description: "Package installation and management",
		XPBonus:     1.1,
	}

	cc.metadata[Build] = &CategoryInfo{
		Name:        "Build & Deployment",
		Icon:        "🔨",
		Color:       "red",
		Description: "Build tools and deployment systems",
		XPBonus:     1.2,
	}

	cc.metadata[Test] = &CategoryInfo{
		Name:        "Testing",
		Icon:        "🧪",
		Color:       "magenta",
		Description: "Testing frameworks and tools",
		XPBonus:     1.3,
	}

	cc.metadata[Text] = &CategoryInfo{
		Name:        "Text Processing",
		Icon:        "📝",
		Color:       "black",
		Description: "Text manipulation and processing",
		XPBonus:     1.0,
	}

	cc.metadata[Archive] = &CategoryInfo{
		Name:        "Archive & Compression",
		Icon:        "🗜️",
		Color:       "gray",
		Description: "File compression and archiving",
		XPBonus:     0.9,
	}

	cc.metadata[Security] = &CategoryInfo{
		Name:        "Security",
		Icon:        "🔒",
		Color:       "red",
		Description: "Security and cryptography tools",
		XPBonus:     1.4,
	}

	cc.metadata[Monitoring] = &CategoryInfo{
		Name:        "Monitoring & Analytics",
		Icon:        "📊",
		Color:       "green",
		Description: "System monitoring and analytics",
		XPBonus:     1.2,
	}

	cc.metadata[Unknown] = &CategoryInfo{
		Name:        "Unknown",
		Icon:        "❓",
		Color:       "gray",
		Description: "Unclassified commands",
		XPBonus:     1.0,
	}
}

// ClassifyCommand determines the category of a command
func (cc *CommandClassifier) ClassifyCommand(command string) Category {
	command = strings.TrimSpace(command)
	if command == "" {
		return Unknown
	}

	// Check each category's patterns
	for category, patterns := range cc.patterns {
		for _, pattern := range patterns {
			if pattern.MatchString(command) {
				return category
			}
		}
	}

	return Unknown
}

// GetCategoryInfo returns metadata for a category
func (cc *CommandClassifier) GetCategoryInfo(category Category) *CategoryInfo {
	if info, exists := cc.metadata[category]; exists {
		return info
	}
	return cc.metadata[Unknown]
}

// GetAllCategories returns all available categories
func (cc *CommandClassifier) GetAllCategories() map[Category]*CategoryInfo {
	return cc.metadata
}

// GetXPMultiplier returns the XP multiplier for a category
func (cc *CommandClassifier) GetXPMultiplier(category Category) float64 {
	if info := cc.GetCategoryInfo(category); info != nil {
		return info.XPBonus
	}
	return 1.0
}

// GetCategoryStats calculates statistics for each category
type CategoryStats struct {
	Category   Category `json:"category"`
	Count      int      `json:"count"`
	Percentage float64  `json:"percentage"`
	TotalXP    int      `json:"total_xp"`
	LastUsed   string   `json:"last_used"`
}

// AnalyzeCategories analyzes command usage by category
func (cc *CommandClassifier) AnalyzeCategories(commands []string) map[Category]*CategoryStats {
	categoryCount := make(map[Category]int)
	totalCommands := len(commands)

	// Count commands by category
	for _, command := range commands {
		category := cc.ClassifyCommand(command)
		categoryCount[category]++
	}

	// Calculate statistics
	stats := make(map[Category]*CategoryStats)
	for category, count := range categoryCount {
		if count > 0 {
			percentage := float64(count) / float64(totalCommands) * 100
			info := cc.GetCategoryInfo(category)
			totalXP := int(float64(count) * info.XPBonus)

			stats[category] = &CategoryStats{
				Category:   category,
				Count:      count,
				Percentage: percentage,
				TotalXP:    totalXP,
				LastUsed:   "", // TODO: Implement last used tracking
			}
		}
	}

	return stats
}
