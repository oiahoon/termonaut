package gamification

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// EasterEggTrigger represents different trigger conditions for easter eggs
type EasterEggTrigger struct {
	ID          string
	Condition   func(context *EasterEggContext) bool
	Messages    []string
	Probability float64 // 0.0 to 1.0
}

// EasterEggContext contains context information for easter egg evaluation
type EasterEggContext struct {
	CommandsInTimeframe int
	TimeframeDuration   time.Duration
	IdleDuration        time.Duration
	IsFirstCommandToday bool
	LastCommand         string
	CommandHistory      []string
	QuotesMismatched    bool
}

// EasterEggManager handles easter egg triggering and display
type EasterEggManager struct {
	triggers []EasterEggTrigger
	rand     *rand.Rand
}

// NewEasterEggManager creates a new easter egg manager
func NewEasterEggManager() *EasterEggManager {
	return &EasterEggManager{
		triggers: getEasterEggTriggers(),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CheckForEasterEgg checks if any easter egg should be triggered
func (eem *EasterEggManager) CheckForEasterEgg(context *EasterEggContext) string {
	for _, trigger := range eem.triggers {
		if trigger.Condition(context) && eem.rand.Float64() < trigger.Probability {
			return eem.selectRandomMessage(trigger.Messages)
		}
	}
	return ""
}

// isCommonCommand checks if a command is a common shell command
func isCommonCommand(cmd string) bool {
	commonCommands := []string{
		"ls", "cd", "pwd", "cat", "echo", "grep", "find", "cp", "mv", "rm", "mkdir", "rmdir",
		"chmod", "chown", "ps", "top", "kill", "which", "man", "git", "vim", "nano", "less",
		"more", "head", "tail", "sort", "uniq", "wc", "tar", "zip", "curl", "wget", "ssh",
		"scp", "sudo", "su", "systemctl", "service", "apt", "yum", "brew", "pip", "npm",
		"docker", "kubectl", "make", "gcc", "python", "node", "go", "java", "mysql", "psql",
	}

	for _, common := range commonCommands {
		if strings.EqualFold(cmd, common) {
			return true
		}
	}
	return false
}

// selectRandomMessage selects a random message from the list
func (eem *EasterEggManager) selectRandomMessage(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[eem.rand.Intn(len(messages))]
}

// getEasterEggTriggers returns all available easter egg triggers
func getEasterEggTriggers() []EasterEggTrigger {
	return []EasterEggTrigger{
		{
			ID: "speed_run",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.CommandsInTimeframe >= 10 && ctx.TimeframeDuration <= 5*time.Second
			},
			Messages: []string{
				"Speed run mode: 🏃‍♂️ calm down!",
				"Whoa there, speed demon! 🏃‍♀️💨",
				"Terminal ninja detected! 🥷⚡",
				"Slow down, Flash! ⚡🛑",
				"Are you typing or playing piano? 🎹💨",
				"Speedrun.com would be proud! 🏆⚡",
				"Terminal Olympics gold medal! 🥇💨",
			},
			Probability: 0.15,
		},
		{
			ID: "coffee_break",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.IdleDuration > 5*time.Minute
			},
			Messages: []string{
				"Time flies! Want a coffee break? ☕",
				"Welcome back! Miss me? 😊",
				"Long time no see! 👋",
				"Did you get lost in the real world? 🌍",
				"Back from the coffee machine? ☕😴",
				"Productivity pause detected! 🛑☕",
				"The terminal was getting lonely... 🥺",
				"Did you solve world hunger while away? 🌍✨",
				"Welcome back, keyboard warrior! ⌨️🛡️",
			},
			Probability: 0.25,
		},
		{
			ID: "new_day",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.IsFirstCommandToday
			},
			Messages: []string{
				"New day, new XP — let's go! 🌅✨",
				"Good morning, terminal warrior! 🌞⚔️",
				"Fresh start, fresh commands! 🆕💪",
				"Time to level up today! 📈🎯",
				"Another day, another shell! 🐚☀️",
				"Ready to conquer the command line? 👑💻",
				"Today's forecast: 100% chance of productivity! 🌤️💪",
				"The early bird catches the... commits? 🐦📝",
				"Rise and grind, code ninja! 🥷☀️",
			},
			Probability: 0.4,
		},
		{
			ID: "quote_mismatch",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.QuotesMismatched
			},
			Messages: []string{
				"Quotes left open... Are you a wizard? 🧙‍♂️✨",
				"Unmatched quotes detected! 📝❓",
				"Quote escape artist! 🎪📜",
				"Did you forget to close something? 🔓💭",
				"Quote limbo! How low can you go? 🤸‍♀️",
				"Syntax error in 3... 2... 1... 💥",
				"The shell is confused by your poetry! 📝🤔",
			},
			Probability: 0.3,
		},
		{
			ID: "git_commit_typo",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "git commit") && strings.Contains(cmd, "fix")
			},
			Messages: []string{
				"Another 'fix' commit? 🔧😅",
				"Fix the fix that fixed the fix? 🔄🛠️",
				"Git gud at commit messages! 😎📝",
				"'Fix' #666: Achievement unlocked! 👹🎯",
				"The eternal cycle of fixes begins... 🔄♾️",
				"Plot twist: This fix breaks something else! 🎭💥",
				"Commit message creativity level: 0/10 😂",
				"'Fix' - the most overused word in git history! 📚",
			},
			Probability: 0.2,
		},
		{
			ID: "docker_whale",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(strings.ToLower(ctx.LastCommand), "docker")
			},
			Messages: []string{
				"🐳 Docker whale says hello!",
				"Container magic in progress! 📦✨",
				"Dockerizing all the things! 🐳🚀",
				"Whale, whale, whale... what do we have here? 🐋",
				"Setting sail with containers! ⛵📦",
				"It works on my machine... in a container! 📦😏",
				"Containerization nation! 🏗️🐳",
				"Docker: Making deployment less scary since 2013! 🎭",
			},
			Probability: 0.15,
		},
		{
			ID: "kubernetes_captain",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "kubectl") || strings.Contains(cmd, "k8s")
			},
			Messages: []string{
				"Ahoy, Kubernetes Captain! ⚓🚢",
				"Orchestrating like a maestro! 🎼🎭",
				"Pod life chose you! 🫛💫",
				"Sailing the cluster seas! 🌊⛵",
				"May the pods be with you! 🫛⭐",
				"Kubectl-ing like a boss! 👑⚓",
				"Cluster management level: Expert! 🎯🏆",
				"YAML files fear your power! 📝⚡",
			},
			Probability: 0.2,
		},
		{
			ID: "vim_escape",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(strings.ToLower(ctx.LastCommand), "vim")
			},
			Messages: []string{
				"Entering the Vim dimension... 🌌📝",
				"Remember: ESC is your friend! ⌨️🆘",
				"Vim: Where legends are made! ⚔️📜",
				"Good luck escaping! 😅🚪",
				"Welcome to the text editor maze! 🧩✏️",
				"hjkl warriors unite! ⚔️🎮",
				"Modal editing: confusing newcomers since 1976! 🤯📅",
				"Vim or Emacs? The eternal question... 🤔⚡",
			},
			Probability: 0.25,
		},
		{
			ID: "exit_command",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(strings.TrimSpace(ctx.LastCommand))
				return cmd == "exit" || cmd == "logout" || cmd == "quit"
			},
			Messages: []string{
				"🎵 You can check out anytime you like... but you can never leave! 🏨",
				"👋 See you later, space cowboy! 🤠",
				"🚪 Goodbye! Don't forget to save your work! 💾",
				"🌅 Until next time, terminal warrior! ⚔️",
				"🎭 Exit, stage left! 🎪",
				"🚀 Safe travels through the digital cosmos! ✨",
			},
			Probability: 0.7,
		},
		{
			ID: "consecutive_errors",
			Condition: func(ctx *EasterEggContext) bool {
				if len(ctx.CommandHistory) < 3 {
					return false
				}
				// Check last 3 commands for consecutive errors (simulated)
				errorCommands := 0
				for i := len(ctx.CommandHistory) - 3; i < len(ctx.CommandHistory); i++ {
					cmd := ctx.CommandHistory[i]
					// Simple heuristic: commands with typos or non-existent commands
					if strings.Contains(cmd, "command not found") ||
						len(strings.Fields(cmd)) == 1 && !isCommonCommand(cmd) {
						errorCommands++
					}
				}
				return errorCommands >= 3
			},
			Messages: []string{
				"🤖 Don't give up, shell-ronin! 🥷",
				"💪 Every master was once a disaster! 🌟",
				"🎯 Practice makes perfect! Keep trying! 🔥",
				"🧠 Errors are just learning opportunities in disguise! 📚",
				"⚡ The terminal believes in you! 💖",
				"🦾 You're building character with every mistake! 💎",
			},
			Probability: 0.8,
		},
		{
			ID: "special_time_420",
			Condition: func(ctx *EasterEggContext) bool {
				now := time.Now()
				return (now.Hour() == 4 && now.Minute() == 20) ||
					(now.Hour() == 16 && now.Minute() == 20) // 4:20 PM too
			},
			Messages: []string{
				"🌿 Wakey wakey, hacker! 👁️",
				"⏰ 4:20 - Time to elevate your coding! 🚀",
				"🍃 High-level programming detected! 🔥",
				"🌱 Growing your skills one command at a time! 📈",
				"💚 Blaze it... with productivity! ⚡",
			},
			Probability: 0.9,
		},
		{
			ID: "challenge_completed",
			Condition: func(ctx *EasterEggContext) bool {
				// Trigger when user completes a significant milestone
				return ctx.CommandsInTimeframe >= 100 && ctx.TimeframeDuration <= 24*time.Hour
			},
			Messages: []string{
				"🧠 You are becoming... unstoppable! 💪⚡",
				"🏆 Command line mastery level: LEGENDARY! 👑",
				"🚀 Houston, we have a terminal genius! 🛸",
				"⭐ Your skills are reaching cosmic levels! 🌌",
				"💎 Forged in the fires of the terminal! 🔥",
			},
			Probability: 1.0,
		},
		{
			ID: "midnight_hacker",
			Condition: func(ctx *EasterEggContext) bool {
				hour := time.Now().Hour()
				return hour >= 0 && hour <= 2 // Between midnight and 2 AM
			},
			Messages: []string{
				"🌙 Midnight oil burning bright! 🔥",
				"🦉 Night owl mode activated! 🌃",
				"⚡ The terminal never sleeps! 💻",
				"🌟 Coding under the stars! ✨",
				"🦇 Creature of the night! 🌒",
				"☕ Coffee level: MAXIMUM! ☕",
			},
			Probability: 0.6,
		},
		{
			ID: "git_push_force",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "git push") &&
					(strings.Contains(cmd, "--force") || strings.Contains(cmd, "-f"))
			},
			Messages: []string{
				"⚠️ Force push detected! May the Git be with you! 🌟",
				"💥 Going nuclear on that repository! 💣",
				"🚨 Red alert! Force push in progress! 🚨",
				"⚡ With great power comes great responsibility! 🕷️",
				"🙏 Hoping your teammates forgive you! 😅",
				"Force push: The nuclear option of Git! ☢️💥",
				"Your commit history just got rewritten! 📝✨",
				"Breaking: Local developer breaks production! 📰💥",
			},
			Probability: 0.3,
		},
		{
			ID: "massive_list",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.HasPrefix(cmd, "ls") &&
					(strings.Contains(cmd, "/") || strings.Contains(cmd, "*"))
			},
			Messages: []string{
				"📂 Exploring the file system like a true explorer! 🗺️",
				"🔍 Detective mode: ON! 🕵️‍♂️",
				"📋 Cataloging the digital universe! 🌌",
				"🗃️ File archaeology in progress! ⛏️",
				"📊 Data discovery mission activated! 🚀",
			},
			Probability: 0.3,
		},
		{
			ID: "rm_danger",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.HasPrefix(cmd, "rm") &&
					(strings.Contains(cmd, "-r") || strings.Contains(cmd, "*"))
			},
			Messages: []string{
				"⚠️ Danger zone! Hope you know what you're doing! 😰",
				"💣 Destructive power activated! 💥",
				"🗑️ Spring cleaning or digital chaos? 🤔",
				"⚡ With great rm comes great responsibility! 🕷️",
				"🙏 RIP files... may they rest in /dev/null 👻",
				"Files deleted faster than you can say 'oops'! 💨🗑️",
				"Marie Kondo would be proud... or terrified! 🧹😱",
			},
			Probability: 0.25,
		},
		{
			ID: "productivity_beast",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.CommandsInTimeframe >= 50 && ctx.TimeframeDuration <= 30*time.Minute
			},
			Messages: []string{
				"🔥 Productivity BEAST mode! 🦾",
				"⚡ Terminal velocity achieved! 🚀",
				"💨 Speed of light coding! ⚡",
				"🌪️ Command line tornado! 🌪️",
				"🏃‍♂️ Gotta go fast! 💨",
			},
			Probability: 0.9,
		},
		{
			ID: "ascii_art_celebration",
			Condition: func(ctx *EasterEggContext) bool {
				// Random chance for ASCII art
				return time.Now().Second()%42 == 0 // Every 42nd second
			},
			Messages: []string{
				"🎨 ASCII Art Time!\n" +
					"     🚀\n" +
					"    /|\\\n" +
					"   / | \\\n" +
					"  |  T  |\n" +
					"  |     |\n" +
					"  ||   ||\n" +
					"  /\\   /\\\n" +
					"Termonaut Power!",
				"🎭 Command Line Theater!\n" +
					"  ╭─────────╮\n" +
					"  │ > Hello │\n" +
					"  │   World │\n" +
					"  ╰─────────╯\n" +
					"The terminal speaks!",
				"🎪 ASCII Circus!\n" +
					"    ∩───∩\n" +
					"   (  ◕   ◕ )\n" +
					"    \\   ▽  /\n" +
					"     \\     /\n" +
					"      \\___/\n" +
					"Code Bear says hi!",
			},
			Probability: 0.05,
		},
		{
			ID: "hidden_command",
			Condition: func(ctx *EasterEggContext) bool {
				secretCommands := []string{"hero", "secret", "konami", "xyzzy", "plugh", "42", "sudo make me a sandwich"}
				cmd := strings.ToLower(ctx.LastCommand)
				for _, secret := range secretCommands {
					if strings.Contains(cmd, secret) {
						return true
					}
				}
				return false
			},
			Messages: []string{
				"🎮 Secret command detected! You found an easter egg!",
				"🥚 Konami code activated! Up, up, down, down...",
				"🔍 Hidden feature unlocked! You're a true explorer!",
				"🏆 Secret achievement: Command line archaeologist!",
				"✨ Magic word detected! Abracadabra!",
				"🤖 sudo: make me a sandwich? 🥪 (Nice try!)",
				"📖 42: The answer to life, universe, and everything! 🌌",
			},
			Probability: 1.0, // Always trigger for secret commands
		},
		{
			ID: "late_night_coding",
			Condition: func(ctx *EasterEggContext) bool {
				hour := time.Now().Hour()
				return hour >= 23 || hour <= 5 // 11 PM to 5 AM
			},
			Messages: []string{
				"🌙 Midnight oil burning bright!",
				"🦉 Night owl mode activated!",
				"☕ Coffee level: Critical!",
				"🌃 The code never sleeps!",
				"💻 3 AM thoughts hit different...",
				"🌟 Night shift productivity!",
			},
			Probability: 0.4,
		},
		{
			ID: "friday_mood",
			Condition: func(ctx *EasterEggContext) bool {
				return time.Now().Weekday() == time.Friday
			},
			Messages: []string{
				"🎉 TGIF! Weekend loading... 🔄",
				"🍺 Friday vibes detected!",
				"🏖️ Weekend.exe starting...",
				"😎 Friday feels! Almost there!",
				"🎊 Last sprint before freedom!",
				"Friday deploy? Living dangerously! 🎲💥",
				"Weekend warrior mode: Activating! 🏹⚔️",
			},
			Probability: 0.1,
		},
		{
			ID: "monday_blues",
			Condition: func(ctx *EasterEggContext) bool {
				return time.Now().Weekday() == time.Monday
			},
			Messages: []string{
				"☕ Monday motivation loading... ⏳",
				"💪 Monday warrior mode: ON!",
				"🌟 New week, new opportunities!",
				"⚡ Monday energy charging... 🔋",
				"🎯 Week one, day one. Let's go!",
				"Monday: The boss battle of weekdays! ⚔️👔",
				"Coffee levels: Critical! Productivity: Pending! ☕⏳",
			},
			Probability: 0.15,
		},
		{
			ID: "long_command",
			Condition: func(ctx *EasterEggContext) bool {
				return len(ctx.LastCommand) > 100
			},
			Messages: []string{
				"📏 That's one long command! Novel in progress?",
				"✍️ War and Peace of commands!",
				"📚 Command line literature!",
				"🎭 Shakespearean command detected!",
				"📖 TL;DR version available? 😅",
				"Command length: Enterprise edition! 🏢📏",
				"Breaking: Developer discovers the spacebar! 📰⌨️",
			},
			Probability: 0.2,
		},
		{
			ID: "python_snake",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.HasPrefix(cmd, "python") || strings.Contains(cmd, ".py")
			},
			Messages: []string{
				"🐍 Python detected! Ssssslithering into code!",
				"Import this: Beautiful is better than ugly! 🌟📜",
				"Pythonic vibes activated! 🐍✨",
				"Zen of Python loading... 🧘‍♂️🐍",
				"Snake charmer at work! 🎵🐍",
				"Life's too short for semicolons! 😏🐍",
			},
			Probability: 0.1,
		},
		{
			ID: "javascript_chaos",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "node") || strings.Contains(cmd, "npm") ||
					strings.Contains(cmd, ".js") || strings.Contains(cmd, "yarn")
			},
			Messages: []string{
				"JavaScript: Making the impossible... possible! 🎭💫",
				"undefined is not a function... yet! 🤷‍♂️💥",
				"Node.js: JavaScript everywhere! 🌍⚡",
				"NPM install: Downloading the internet... 📦🌐",
				"== vs === : The eternal struggle! ⚖️😅",
				"Callback hell survivors club! 🔥😈",
			},
			Probability: 0.1,
		},
		{
			ID: "database_queries",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "mysql") || strings.Contains(cmd, "psql") ||
					strings.Contains(cmd, "mongo") || strings.Contains(cmd, "redis") ||
					strings.Contains(cmd, "sqlite")
			},
			Messages: []string{
				"Database whisperer detected! 🗄️🔮",
				"SELECT * FROM awesome WHERE you = 'amazing'! 🏆📊",
				"Joining tables like a relationship counselor! 💒📋",
				"SQL: Structured Query Language or Squirrel? 🐿️🤔",
				"NoSQL? More like NoProblems! 📊😎",
				"ACID compliance: Not just for chemistry! ⚗️📊",
			},
			Probability: 0.15,
		},
		{
			ID: "testing_dedication",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "test") || strings.Contains(cmd, "jest") ||
					strings.Contains(cmd, "pytest") || strings.Contains(cmd, "rspec") ||
					strings.Contains(cmd, "mocha")
			},
			Messages: []string{
				"Testing in production? How adventurous! 🎢🧪",
				"Red, Green, Refactor - the holy trinity! 🔴🟢🔄",
				"99 bugs in the code, take one down... 🐛🎵",
				"Test coverage: Aiming for the stars! ⭐📊",
				"Quality assurance: Because YOLO isn't a strategy! 🎯✅",
				"Debugging: Being a detective for your own crimes! 🕵️‍♂️🔍",
			},
			Probability: 0.12,
		},
		{
			ID: "ai_commands",
			Condition: func(ctx *EasterEggContext) bool {
				cmd := strings.ToLower(ctx.LastCommand)
				return strings.Contains(cmd, "chatgpt") || strings.Contains(cmd, "claude") ||
					strings.Contains(cmd, "copilot") || strings.Contains(cmd, "ai") ||
					strings.Contains(cmd, "gpt") || strings.Contains(cmd, "llm")
			},
			Messages: []string{
				"AI assistant detected! Hello, fellow digital being! 🤖👋",
				"Humans and AI, coding together! 🤝💻",
				"The future is collaborative! 🚀🤖",
				"AI: Artificial Intelligence or Actually Intelligent? 🧠✨",
				"Prompt engineering: The new coding skill! 💬⚡",
				"Beep boop: AI translation successful! 🤖🔄",
			},
			Probability: 0.08,
		},
	}
}

// FormatEasterEggMessage formats the easter egg message with proper styling
func FormatEasterEggMessage(message string) string {
	// Check terminal capabilities for enhanced formatting
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	colorterm := os.Getenv("COLORTERM")

	// Enhanced formatting for modern terminals
	if termProgram == "warp" || termProgram == "iterm.app" || termProgram == "vscode" ||
		termProgram == "alacritty" || termProgram == "kitty" || termProgram == "hyper" ||
		colorterm == "truecolor" {
		// Use enhanced formatting with better spacing and colors
		return fmt.Sprintf("\n\033[38;5;214m🥚\033[0m \033[1m%s\033[0m\n", message)
	}

	// Fallback for basic terminals
	return fmt.Sprintf("\n🥚 %s\n", message)
}

// GetRandomMotivationalQuote returns a random motivational quote
func (eem *EasterEggManager) GetRandomMotivationalQuote() string {
	quotes := []string{
		"Code like nobody's watching! 💻✨",
		"Debugging is like being a detective! 🔍🕵️‍♂️",
		"Every command brings you closer to mastery! 🎯",
		"Terminal wisdom: One command at a time! 🧘‍♂️",
		"You're building the future, one line at a time! 🏗️🚀",
		"Embrace the command line, become one with the shell! 🐚🕉️",
		"Error messages are just the computer's way of teaching! 📚💡",
		"Persistence beats resistance! Keep coding! 💪⚡",
	}
	return quotes[eem.rand.Intn(len(quotes))]
}

// IsModernTerminal checks if the current terminal supports modern features
func IsModernTerminal() bool {
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	termProgramVersion := os.Getenv("TERM_PROGRAM_VERSION")
	colorterm := os.Getenv("COLORTERM")
	term := strings.ToLower(os.Getenv("TERM"))

	// Check for modern terminal emulators
	modernTerminals := []string{
		"warp",      // Warp Terminal
		"iterm.app", // iTerm2
		"vscode",    // VS Code Terminal
		"alacritty", // Alacritty
		"kitty",     // Kitty
		"hyper",     // Hyper
		"tabby",     // Tabby
		"terminus",  // Terminus
		"rio",       // Rio Terminal
	}

	for _, modern := range modernTerminals {
		if termProgram == modern {
			return true
		}
	}

	// Check for Windows Terminal
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	// Check for Terminal.app (macOS) with recent versions
	if termProgram == "apple_terminal" && termProgramVersion != "" {
		return true
	}

	// Check for truecolor support
	if colorterm == "truecolor" || colorterm == "24bit" {
		return true
	}

	// Check for 256-color terminals
	if strings.Contains(term, "256color") {
		return true
	}

	return false
}

// GetTerminalInfo returns information about the current terminal
func GetTerminalInfo() map[string]string {
	return map[string]string{
		"TERM":                 os.Getenv("TERM"),
		"TERM_PROGRAM":         os.Getenv("TERM_PROGRAM"),
		"TERM_PROGRAM_VERSION": os.Getenv("TERM_PROGRAM_VERSION"),
		"COLORTERM":            os.Getenv("COLORTERM"),
		"WT_SESSION":           os.Getenv("WT_SESSION"),
		"ITERM_PROFILE":        os.Getenv("ITERM_PROFILE"),
		"ITERM_SESSION_ID":     os.Getenv("ITERM_SESSION_ID"),
	}
}
