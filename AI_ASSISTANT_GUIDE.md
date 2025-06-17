# 🤖 AI Assistant Development Guide for Termonaut

This document provides comprehensive guidance for AI assistants working on the Termonaut project. It includes project context, development principles, coding standards, and specific instructions for maintaining consistency across development sessions.

## 📖 Project Overview

### What is Termonaut?
Termonaut is a **gamified terminal productivity tracker** that transforms command-line usage into an engaging RPG-like experience. It combines:
- **Command logging** with minimal performance impact
- **XP and leveling system** with space-themed progression
- **Achievement badges** for usage milestones and discoveries
- **Rich statistics** about terminal habits and productivity
- **GitHub integration** for shareable profile badges
- **Privacy-first approach** with local-only data storage by default

### Core Philosophy
- **Terminal-Native**: Everything happens in the CLI, no web dashboards
- **Performance-First**: < 1ms logging overhead, async operations
- **Privacy-Focused**: Local SQLite storage, optional sharing only
- **Gamification-Driven**: Make terminal usage fun and rewarding
- **Developer-Friendly**: Built by developers, for developers

## 🎯 Development Principles

### 1. User Experience First
- **Minimal Setup**: `termonaut init` should be all users need
- **Instant Gratification**: Show immediate value after first commands
- **Non-Intrusive**: Never interrupt or slow down normal terminal usage
- **Discoverable**: Features should be easy to find and understand

### 2. Technical Excellence
- **Performance**: Profile and benchmark everything
- **Reliability**: Handle edge cases gracefully
- **Maintainability**: Clean, documented, testable code
- **Compatibility**: Support major shells (Bash, Zsh, Fish)

### 3. Progressive Enhancement
- **Core First**: Basic logging and stats before gamification
- **Graceful Degradation**: Features should fail safely
- **Configurable**: Users can disable features they don't want
- **Extensible**: Architecture supports future enhancements

## 🏗️ Architecture Guidelines

### System Components
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Shell Hook    │───▶│  Data Capture    │───▶│   SQLite DB     │
│   (preexec)     │    │     Layer        │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
┌─────────────────┐    ┌──────────────────┐             │
│   CLI Frontend  │◀───│  Stats Engine    │◀────────────┘
│   (commands)    │    │                  │
└─────────────────┘    └──────────────────┘
        │                       │
        ▼                       ▼
┌─────────────────┐    ┌──────────────────┐
│ Gamification    │    │  Export/Sync     │
│    Engine       │    │     Module       │
└─────────────────┘    └──────────────────┘
```

### Key Design Patterns

**1. Command Pattern for CLI**
```go
type Command interface {
    Execute(args []string) error
    Help() string
}
```

**2. Repository Pattern for Data**
```go
type CommandRepository interface {
    Store(cmd Command) error
    GetStats(filter StatsFilter) (*Stats, error)
    GetSessions(limit int) ([]Session, error)
}
```

**3. Observer Pattern for Gamification**
```go
type EventHandler interface {
    HandleCommandLogged(cmd Command)
    HandleSessionStarted(session Session)
}
```

## 🗄️ Database Design

### Schema Principles
- **Normalized**: Avoid data duplication
- **Indexed**: Fast queries on common operations
- **Extensible**: Easy to add new fields
- **Efficient**: Minimal storage overhead

### Key Tables
```sql
-- Core data
commands (id, timestamp, session_id, command, exit_code, cwd, duration_ms)
sessions (id, start_time, end_time, terminal_pid, shell_type)

-- Gamification
user_progress (total_xp, current_level, commands_count, current_streak)
achievements (id, name, description, earned_at, xp_bonus)

-- Performance optimization
daily_stats (date, commands_count, session_count, active_time_minutes)
```

### Query Patterns
- Use prepared statements for security
- Implement connection pooling for performance
- Cache frequently accessed data (daily stats)
- Use transactions for consistency

## 🎮 Gamification System

### XP Calculation Rules
```go
const (
    BaseXPPerCommand = 1
    NewCommandBonus = 5
    XPPerActiveMinute = 0.1
)

// Category multipliers
var CategoryBonus = map[string]int{
    "git":    2,  // Version control
    "vim":    2,  // Editors
    "make":   3,  // Build tools
    "docker": 2,  // Containers
}
```

### Level Progression
- **Formula**: `level = sqrt(total_xp / 100)`
- **Progression**: Exponential growth, early levels come quickly
- **Titles**: Space-themed (Cadet, Navigator, Commander, etc.)

### Achievement Categories
1. **Milestones**: Command counts, levels, time-based
2. **Discovery**: New commands, categories explored
3. **Consistency**: Daily/weekly streaks
4. **Efficiency**: Speed, productivity patterns
5. **Special**: Easter eggs, seasonal events

## 💻 CLI Interface Standards

### Command Structure
```bash
termonaut <subcommand> [flags] [arguments]
```

### Subcommand Categories
- **Core**: `stats`, `sessions`, `xp`, `badges`
- **Config**: `config get/set/reset`, `init`
- **Data**: `export`, `import`, `backup`
- **Advanced**: `debug`, `migrate`, `profile`

### Output Guidelines
- **Human-readable by default**: Rich formatting, colors, emojis
- **Machine-readable option**: `--json` flag for automation
- **Consistent formatting**: Use templates for similar outputs
- **Progress indicators**: Show progress for long operations

### Error Handling
```go
// Good: Specific, actionable errors
return fmt.Errorf("failed to initialize shell hook: %w", err)

// Bad: Generic, unhelpful errors
return fmt.Errorf("something went wrong")
```

## 🧪 Testing Strategy

### Test Categories
1. **Unit Tests**: Individual functions, pure logic
2. **Integration Tests**: Database operations, CLI commands
3. **Performance Tests**: Latency, memory usage, concurrency
4. **End-to-End Tests**: Full workflows, shell integration

### Testing Guidelines
- **Test naming**: `Test<Function>_<Scenario>_<Expected>`
- **Table-driven tests**: For multiple input/output pairs
- **Test isolation**: Each test should be independent
- **Mock external dependencies**: Database, filesystem, shell

### Coverage Targets
- **Unit tests**: 80%+ coverage
- **Integration tests**: All major workflows
- **Performance tests**: Critical path operations
- **E2E tests**: Happy path scenarios

## 📝 Coding Standards

### Go Specific (If Using Go)

**Project Structure:**
```
cmd/           # CLI entry points
internal/      # Private application code
pkg/           # Public library code
scripts/       # Development/deployment scripts
tests/         # Test files and fixtures
```

**Code Style:**
- Follow `gofmt` formatting
- Use `golangci-lint` for linting
- Document public functions and types
- Handle errors explicitly
- Use structured logging (zap/logrus)

**Patterns:**
```go
// Good: Clear error handling
result, err := operation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Good: Structured logging
logger.Info("command processed",
    zap.String("command", cmd),
    zap.Duration("duration", elapsed))

// Good: Interface design
type StatsCalculator interface {
    CalculateDaily(date time.Time) (*DailyStats, error)
    CalculateWeekly(week int) (*WeeklyStats, error)
}
```

### Python Specific (If Using Python)

**Project Structure:**
```
termonaut/     # Main package
├── cli/       # CLI commands
├── core/      # Core business logic
├── db/        # Database operations
└── utils/     # Shared utilities
tests/         # Test files
scripts/       # Development scripts
```

**Code Style:**
- Follow PEP 8 formatting (black, isort)
- Use type hints for public APIs
- Document with docstrings
- Use dataclasses for data structures
- Prefer pathlib over os.path

**Patterns:**
```python
# Good: Type hints and error handling
def calculate_xp(commands: List[Command]) -> int:
    """Calculate total XP from command list."""
    try:
        return sum(cmd.base_xp + cmd.bonus_xp for cmd in commands)
    except Exception as e:
        logger.error(f"XP calculation failed: {e}")
        raise

# Good: Dataclass usage
@dataclass
class CommandStats:
    total_commands: int
    unique_commands: int
    most_used: str
    least_used: str
```

## 🔄 Development Workflow

### Feature Development Process
1. **Issue Creation**: Clear description, acceptance criteria
2. **Branch Creation**: `feature/descriptive-name`
3. **Implementation**: Follow TDD when possible
4. **Testing**: Unit + integration tests
5. **Documentation**: Update relevant docs
6. **PR Review**: Code review + automated checks
7. **Merge**: Squash commits for clean history

### Git Conventions
**Commit Messages:**
```
feat: add streak calculation for achievements
fix: resolve shell hook timing issue
docs: update installation instructions
test: add integration tests for export functionality
refactor: simplify XP calculation logic
perf: optimize database query performance
```

**Branch Naming:**
- `feature/feature-name`
- `bugfix/issue-description`
- `hotfix/critical-fix`
- `docs/documentation-update`

### Code Review Checklist
- [ ] Functionality works as specified
- [ ] Tests cover new/modified code
- [ ] Documentation is updated
- [ ] Performance impact is acceptable
- [ ] Security implications considered
- [ ] Error handling is comprehensive
- [ ] Code follows project style guidelines

## 📚 Documentation Standards

### Code Documentation
- **Public APIs**: Complete docstrings/comments
- **Complex Logic**: Inline comments explaining why
- **Examples**: Usage examples in documentation
- **Edge Cases**: Document known limitations

### User Documentation
- **README**: Quick start, installation, basic usage
- **CLI Help**: Built-in help for all commands
- **Configuration**: All options documented with examples
- **Troubleshooting**: Common issues and solutions

### Developer Documentation
- **Architecture**: High-level system design
- **Contributing**: Setup, workflow, standards
- **API Reference**: Generated from code comments
- **Changelog**: All notable changes between versions

## 🚀 Release Management

### Version Strategy
- **Semantic Versioning**: MAJOR.MINOR.PATCH
- **Pre-release**: alpha, beta, rc suffixes
- **Development**: -dev suffix for unreleased versions

### Release Checklist
**Pre-Release:**
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Performance benchmarks met
- [ ] Security scan clean
- [ ] Changelog updated

**Release:**
- [ ] Version bumped
- [ ] Git tag created
- [ ] Binaries built for all platforms
- [ ] Package managers updated
- [ ] Release notes published

**Post-Release:**
- [ ] Monitor for issues
- [ ] Community feedback collection
- [ ] Plan next iteration

## 🔍 Debugging and Troubleshooting

### Debug Mode
```bash
export TERMONAUT_DEBUG=1
export TERMONAUT_LOG_LEVEL=debug
termonaut stats --debug
```

### Common Issues
1. **Shell Integration**: Hook not installed or wrong shell
2. **Performance**: Slow database queries or excessive logging
3. **Data Corruption**: Database locks or interrupted writes
4. **Compatibility**: Shell version or OS differences

### Diagnostic Tools
```bash
# Database inspection
sqlite3 ~/.termonaut/termonaut.db ".schema"
sqlite3 ~/.termonaut/termonaut.db "SELECT * FROM commands LIMIT 5;"

# Log analysis
tail -f ~/.termonaut/termonaut.log | grep ERROR

# Performance profiling
termonaut stats --profile
```

## 🤖 AI Assistant Instructions

### When Starting a New Session
1. **Review Current State**: Check latest commits, open issues, current milestone
2. **Understand Context**: Read relevant documentation sections
3. **Identify Goals**: What needs to be accomplished this session?
4. **Plan Approach**: Break down work into manageable steps

### Development Guidelines for AI
1. **Follow Architecture**: Respect the established patterns and structure
2. **Maintain Quality**: Write tests, handle errors, document changes
3. **Consider Performance**: Profile critical paths, avoid blocking operations
4. **Update Documentation**: Keep docs in sync with code changes
5. **Think Incrementally**: Small, focused changes over large refactors

### Code Generation Best Practices
- **Start with tests**: Write tests first when adding new functionality
- **Handle edge cases**: Consider error conditions and boundary cases
- **Follow patterns**: Use established patterns in the codebase
- **Add logging**: Include appropriate logging for debugging
- **Document decisions**: Explain complex logic in comments

### When Making Changes
1. **Verify current state**: Run tests, check existing functionality
2. **Plan the change**: Understand impact on other components
3. **Implement incrementally**: Make small, testable changes
4. **Update documentation**: Keep all docs current
5. **Validate the change**: Test thoroughly, check performance

### Session Handoff
When ending a session, document:
- **What was accomplished**: Summary of changes made
- **Current state**: What's working, what needs attention
- **Next steps**: Immediate priorities for next session
- **Known issues**: Any problems discovered but not fixed
- **Documentation updates**: What docs were modified

## 📖 Key Reference Documents

### Essential Reading
1. **README.md**: Project overview and user instructions
2. **DEVELOPMENT.md**: Technical architecture and setup
3. **PROJECT_PLANNING.md**: Roadmap and milestones
4. **Database Schema**: Current table structures and relationships

### Quick Reference
- **CLI Commands**: All subcommands and their purposes
- **Configuration Options**: Available settings and defaults
- **Achievement Definitions**: All badges and unlock criteria
- **XP Calculation**: Formulas and multipliers

### External Resources
- [SQLite Documentation](https://sqlite.org/docs.html)
- [Cobra CLI Framework](https://cobra.dev/) (if using Go)
- [Click Framework](https://click.palletsprojects.com/) (if using Python)
- [Semantic Versioning](https://semver.org/)

---

This guide should be updated as the project evolves. When making significant architectural changes or adding new patterns, update this document to maintain consistency for future development sessions.