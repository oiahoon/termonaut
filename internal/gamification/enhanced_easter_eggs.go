package gamification

import (
	"strings"
	"time"
)

// getEnhancedEasterEggTriggers returns additional easter egg triggers
func getEnhancedEasterEggTriggers() []EasterEggTrigger {
	return []EasterEggTrigger{
		// 编程语言特定彩蛋
		{
			ID: "rust_safety",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "cargo") || 
					   strings.Contains(ctx.LastCommand, "rustc")
			},
			Messages: []string{
				"🦀 Rust: Where memory safety meets blazing speed! Zero-cost abstractions FTW!",
				"🔒 Rust detected! Your code is safer than a bank vault wrapped in bubble wrap!",
				"⚡ Cargo building... Time to grab some coffee and contemplate ownership!",
				"🦀 Rust: Making C++ developers question their life choices since 2010!",
			},
			Probability: 0.15,
		},
		{
			ID: "go_gopher",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "go run") || 
					   strings.Contains(ctx.LastCommand, "go build") ||
					   strings.Contains(ctx.LastCommand, "go mod")
			},
			Messages: []string{
				"🐹 Go Gopher says: Simple, fast, and reliable! That's the Go way!",
				"⚡ Go: Concurrency made easy! Goroutines are dancing in your CPU!",
				"🚀 Building with Go... Rob Pike would be proud!",
				"🐹 Go: Less is more, except when it comes to performance!",
			},
			Probability: 0.15,
		},
		{
			ID: "python_zen",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "python") || 
					   strings.Contains(ctx.LastCommand, "pip") ||
					   strings.Contains(ctx.LastCommand, "poetry")
			},
			Messages: []string{
				"🐍 Python: Beautiful is better than ugly! The Zen of Python guides you!",
				"✨ Python detected! Remember: There should be one obvious way to do it!",
				"🐍 Pythonic code ahead! Simple is better than complex!",
				"🎯 Python: Readability counts! Your future self will thank you!",
			},
			Probability: 0.12,
		},
		
		// 时间相关彩蛋
		{
			ID: "late_night_coding",
			Condition: func(ctx *EasterEggContext) bool {
				hour := time.Now().Hour()
				return hour >= 23 || hour <= 5
			},
			Messages: []string{
				"🌙 Late night coding session detected! Don't forget to blink!",
				"☕ 3 AM and still coding? You're either very dedicated or very caffeinated!",
				"🦉 Night owl mode activated! Remember: bugs love the darkness!",
				"🌟 Coding under the stars! Your dedication is astronomical!",
				"😴 Pro tip: Sleep is not a bug, it's a feature! Consider implementing it!",
			},
			Probability: 0.25,
		},
		{
			ID: "early_bird",
			Condition: func(ctx *EasterEggContext) bool {
				hour := time.Now().Hour()
				return hour >= 5 && hour <= 7
			},
			Messages: []string{
				"🌅 Early bird catches the bug! Morning coding session initiated!",
				"☀️ Rise and code! The early developer gets the clean commits!",
				"🐦 5 AM coding? You're either very disciplined or couldn't sleep!",
				"🌄 Dawn patrol! Your code is as fresh as the morning dew!",
			},
			Probability: 0.2,
		},
		
		// 工作流程彩蛋
		{
			ID: "test_driven",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "test") ||
					   strings.Contains(ctx.LastCommand, "jest") ||
					   strings.Contains(ctx.LastCommand, "pytest") ||
					   strings.Contains(ctx.LastCommand, "rspec")
			},
			Messages: []string{
				"🧪 Testing detected! Red, Green, Refactor - the TDD mantra!",
				"✅ Tests are love, tests are life! Your future self sends thanks!",
				"🔬 Science mode activated! Hypothesis: Your code works. Let's test it!",
				"🛡️ Testing: The shield against the dark arts of production bugs!",
				"🎯 Good tests are like good friends - they tell you when you're wrong!",
			},
			Probability: 0.18,
		},
		{
			ID: "deployment_anxiety",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "deploy") ||
					   strings.Contains(ctx.LastCommand, "kubectl apply") ||
					   strings.Contains(ctx.LastCommand, "terraform apply")
			},
			Messages: []string{
				"🚀 Deployment detected! May the force be with your servers!",
				"😰 Deploying on Friday? Living dangerously, I see!",
				"🎲 Deployment: The ultimate test of your backup strategy!",
				"🙏 Deploying... Time to sacrifice a rubber duck to the DevOps gods!",
				"⚡ Production deployment! Remember: It works on my machine™!",
			},
			Probability: 0.3,
		},
		
		// 创意和幽默彩蛋
		{
			ID: "stack_overflow",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "curl") && 
					   strings.Contains(ctx.LastCommand, "stackoverflow")
			},
			Messages: []string{
				"📚 Stack Overflow detected! The sacred texts are being consulted!",
				"🤔 Ah, the ancient ritual of copy-paste from Stack Overflow begins!",
				"👨‍💻 Stack Overflow: Where developers go to feel both smart and stupid!",
				"📖 Consulting the hive mind... May the accepted answer be with you!",
			},
			Probability: 0.4,
		},
		{
			ID: "rubber_duck",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.CommandsInTimeframe > 50 && ctx.TimeframeDuration < time.Hour
			},
			Messages: []string{
				"🦆 Rubber duck says: Maybe it's time to explain your code to me?",
				"🤔 Lots of commands detected! Have you tried rubber duck debugging?",
				"🦆 Quack! Sometimes the best debugger has yellow feathers!",
				"💡 Pro tip: Explaining your code to a duck often reveals the bug!",
			},
			Probability: 0.1,
		},
		{
			ID: "coffee_break",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.IdleDuration > 30*time.Minute && ctx.IdleDuration < 2*time.Hour
			},
			Messages: []string{
				"☕ Coffee break detected! Caffeine levels restored to optimal!",
				"🍵 Tea time! The best code is written with a warm beverage nearby!",
				"☕ Coffee: The fuel that powers the internet! Welcome back!",
				"🥤 Hydration complete! Ready to tackle those bugs with renewed vigor!",
				"☕ Coffee break: The unofficial debugging technique that actually works!",
			},
			Probability: 0.2,
		},
		
		// 技术栈特定彩蛋
		{
			ID: "kubernetes_complexity",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "kubectl") ||
					   strings.Contains(ctx.LastCommand, "k8s") ||
					   strings.Contains(ctx.LastCommand, "helm")
			},
			Messages: []string{
				"⚓ Kubernetes detected! Welcome to the container orchestration maze!",
				"🐳 K8s: Making simple things complex since 2014! But it scales!",
				"⚙️ Kubernetes: Where YAML goes to become sentient!",
				"🎭 K8s: The art of making 3 containers look like rocket science!",
				"🚢 Sailing the Kubernetes seas! May your pods be ever ready!",
			},
			Probability: 0.2,
		},
		{
			ID: "ai_assistance",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "chatgpt") ||
					   strings.Contains(ctx.LastCommand, "copilot") ||
					   strings.Contains(ctx.LastCommand, "claude") ||
					   strings.Contains(ctx.LastCommand, "openai")
			},
			Messages: []string{
				"🤖 AI assistance detected! The robots are helping us code now!",
				"🧠 AI pair programming! Your silicon buddy is ready to help!",
				"🚀 AI-powered coding! The future is now, and it autocompletes!",
				"🤝 Human + AI collaboration! Together we're unstoppable!",
				"🎯 AI assistance: Making developers 10x faster at asking questions!",
			},
			Probability: 0.15,
		},
		
		// 情感支持彩蛋
		{
			ID: "frustration_support",
			Condition: func(ctx *EasterEggContext) bool {
				return strings.Contains(ctx.LastCommand, "kill") ||
					   strings.Contains(ctx.LastCommand, "pkill") ||
					   (ctx.CommandsInTimeframe > 20 && ctx.TimeframeDuration < 10*time.Minute)
			},
			Messages: []string{
				"😤 Feeling frustrated? Take a deep breath! Every bug is a learning opportunity!",
				"🧘 Debugging zen: The bug is not your enemy, it's your teacher!",
				"💪 Tough debugging session? You've got this! Every expert was once a beginner!",
				"🌟 Remember: The best developers aren't those who never encounter bugs, but those who fix them gracefully!",
				"🎯 Debugging is like being a detective in a crime movie where you're also the murderer!",
			},
			Probability: 0.25,
		},
		{
			ID: "productivity_celebration",
			Condition: func(ctx *EasterEggContext) bool {
				return ctx.CommandsInTimeframe > 100 && ctx.TimeframeDuration < 2*time.Hour
			},
			Messages: []string{
				"🎉 Productivity mode: ACTIVATED! You're on fire today!",
				"⚡ Command velocity: MAXIMUM! Your keyboard is smoking!",
				"🚀 Terminal ninja detected! Your CLI-fu is strong!",
				"🏆 Productivity champion! You're making the terminal proud!",
				"💫 You're in the zone! The flow state is strong with this one!",
			},
			Probability: 0.3,
		},
		
		// 季节性和时间特定彩蛋
		{
			ID: "monday_motivation",
			Condition: func(ctx *EasterEggContext) bool {
				return time.Now().Weekday() == time.Monday && ctx.IsFirstCommandToday
			},
			Messages: []string{
				"💪 Monday motivation! New week, new bugs to squash!",
				"☀️ Monday morning! Time to turn coffee into code!",
				"🚀 Monday launch sequence initiated! Let's make this week awesome!",
				"🎯 Monday mindset: Fresh start, clean slate, infinite possibilities!",
			},
			Probability: 0.4,
		},
		{
			ID: "friday_feeling",
			Condition: func(ctx *EasterEggContext) bool {
				return time.Now().Weekday() == time.Friday && time.Now().Hour() > 15
			},
			Messages: []string{
				"🎉 Friday afternoon! The weekend is calling, but first... one more commit!",
				"🍻 TGIF! Time to deploy and pray... or maybe just pray!",
				"🏖️ Friday vibes! Weekend mode loading... 99% complete!",
				"😎 Friday afternoon coding! Living dangerously close to the weekend!",
			},
			Probability: 0.3,
		},
	}
}

// MergeEasterEggTriggers merges original and enhanced triggers
func MergeEasterEggTriggers() []EasterEggTrigger {
	original := getEasterEggTriggers()
	enhanced := getEnhancedEasterEggTriggers()
	
	// Combine both sets
	allTriggers := make([]EasterEggTrigger, 0, len(original)+len(enhanced))
	allTriggers = append(allTriggers, original...)
	allTriggers = append(allTriggers, enhanced...)
	
	return allTriggers
}

// GetEasterEggStats returns statistics about easter egg triggers
func GetEasterEggStats() map[string]interface{} {
	triggers := MergeEasterEggTriggers()
	
	stats := map[string]interface{}{
		"total_triggers":     len(triggers),
		"categories": map[string]int{
			"programming_languages": 3,
			"time_based":           4,
			"workflow":             2,
			"creative_humor":       3,
			"tech_stack":           2,
			"emotional_support":    2,
			"seasonal":             2,
		},
		"average_probability": calculateAverageProbability(triggers),
		"high_probability":    countHighProbabilityTriggers(triggers),
	}
	
	return stats
}

func calculateAverageProbability(triggers []EasterEggTrigger) float64 {
	total := 0.0
	for _, trigger := range triggers {
		total += trigger.Probability
	}
	return total / float64(len(triggers))
}

func countHighProbabilityTriggers(triggers []EasterEggTrigger) int {
	count := 0
	for _, trigger := range triggers {
		if trigger.Probability > 0.2 {
			count++
		}
	}
	return count
}
