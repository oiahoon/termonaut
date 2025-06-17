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
			ID: "hidden_command",
			Condition: func(ctx *EasterEggContext) bool {
				secretCommands := []string{"hero", "secret", "konami", "xyzzy", "plugh"}
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
