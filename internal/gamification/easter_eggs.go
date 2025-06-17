package gamification

import (
	"fmt"
	"math/rand"
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
			},
			Probability: 0.8,
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
			},
			Probability: 0.6,
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
			},
			Probability: 0.9,
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
			},
			Probability: 0.7,
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
			},
			Probability: 0.5,
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
			},
			Probability: 0.3,
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
			},
			Probability: 0.4,
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
			},
			Probability: 0.6,
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
			},
			Probability: 0.8,
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
			},
			Probability: 0.7,
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
			},
			Probability: 0.2,
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
			},
			Probability: 0.3,
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
			},
			Probability: 0.4,
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
			},
			Probability: 0.6,
		},
	}
}

// FormatEasterEggMessage formats the easter egg message with proper styling
func FormatEasterEggMessage(message string) string {
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
