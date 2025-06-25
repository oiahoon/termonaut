# 🤝 Contributing to Termonaut

Thank you for your interest in contributing to Termonaut! This document provides guidelines and information for contributors.

## 🎯 Ways to Contribute

### 🐛 Bug Reports
- Search existing issues before creating new ones
- Use the bug report template
- Include steps to reproduce, expected behavior, and system info
- Add relevant labels (bug, priority, etc.)

### 💡 Feature Requests
- Check roadmap and existing issues first
- Use the feature request template
- Explain the use case and benefit
- Consider implementation complexity

### 💻 Code Contributions
- Start with "good first issue" labels
- Fork the repository and create feature branches
- Follow coding standards and write tests
- Update documentation as needed

### 📚 Documentation
- Fix typos, improve clarity, add examples
- Update guides when features change
- Translate documentation (future)
- Write tutorials and blog posts

## 🚀 Getting Started

### Development Setup

1. **Fork and Clone**
```bash
git clone https://github.com/yourusername/termonaut.git
cd termonaut
```

2. **Setup Environment**
```bash
# For Go development
make dev-setup-go

# For Python development
make dev-setup-python
```

3. **Install Pre-commit Hooks**
```bash
pip install pre-commit
pre-commit install
```

4. **Verify Setup**
```bash
make test
make build
```

### Development Workflow

1. **Create Issue** (for new features/bugs)
2. **Create Branch** from main: `git checkout -b feature/your-feature`
3. **Make Changes** following our standards
4. **Add Tests** for new functionality
5. **Update Docs** if needed
6. **Commit Changes** with conventional commits
7. **Push Branch** and create Pull Request
8. **Address Reviews** and iterate

## 📝 Coding Standards

### Commit Messages
Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

feat: add streak calculation for gamification
fix: resolve shell hook timing issue
docs: update installation instructions
test: add integration tests for export
refactor: simplify XP calculation logic
perf: optimize database queries
```

**Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `test`: Test additions/changes
- `refactor`: Code restructuring without feature changes
- `perf`: Performance improvements
- `style`: Code style changes (formatting, etc.)
- `chore`: Build process, dependencies, etc.

### Code Style

**General Principles:**
- Write clear, readable code
- Add comments for complex logic
- Handle errors gracefully
- Follow language-specific style guides
- Write comprehensive tests

**Go Style:**
- Follow `gofmt` formatting
- Use `golangci-lint` for linting
- Document public functions and types
- Handle errors explicitly

**Python Style:**
- Follow PEP 8 with `black` formatting
- Use type hints for public APIs
- Write docstrings for functions/classes
- Prefer pathlib over os.path

### Testing Requirements
- Unit tests for all new functions
- Integration tests for features
- Performance tests for critical paths
- Maintain 80%+ code coverage
- All tests must pass before merging

## 🔄 Pull Request Process

### Before Submitting
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] PR description explains changes

### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Refactoring

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests pass
```

### Review Process
1. **Automated Checks**: CI/CD pipeline runs tests and linting
2. **Maintainer Review**: Core team reviews code and design
3. **Address Feedback**: Make requested changes
4. **Final Approval**: Merge when approved and tests pass

### Merge Requirements
- At least one maintainer approval
- All CI checks passing
- No merge conflicts
- Up-to-date with main branch

## 🎖️ Recognition

### Contributor Types
- **Code Contributors**: Implement features and fix bugs
- **Documentation Contributors**: Improve docs and guides
- **Community Contributors**: Help users, triage issues
- **Beta Testers**: Test pre-releases and provide feedback

### Recognition Methods
- GitHub contributor graph
- CONTRIBUTORS.md acknowledgment
- Release notes mentions
- Special badges for significant contributions

## 📋 Issue Guidelines

### Bug Reports
**Use this template:**
```markdown
**Describe the bug**
Clear description of the problem

**To Reproduce**
Steps to reproduce the behavior

**Expected behavior**
What you expected to happen

**Environment**
- OS: [e.g. macOS 12.0]
- Shell: [e.g. zsh 5.8]
- Termonaut version: [e.g. 0.1.0]

**Additional context**
Any other relevant information
```

### Feature Requests
**Use this template:**
```markdown
**Is your feature request related to a problem?**
Description of the problem

**Describe the solution you'd like**
Clear description of desired feature

**Describe alternatives considered**
Other solutions you've considered

**Additional context**
Any other relevant information
```

## 🚨 Code of Conduct

### Our Pledge
We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Expected Behavior
- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior
- Trolling, insulting/derogatory comments, and personal attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate

### Enforcement
Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team. All complaints will be reviewed and investigated and will result in a response that is deemed necessary and appropriate to the circumstances.

## 📞 Getting Help

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Discord** (future): Real-time community chat
- **Email**: Direct contact with maintainers

### Asking Questions
- Search existing issues and discussions first
- Provide context and details
- Use appropriate channels for different topics
- Be patient and respectful

### Mentorship
- New contributors welcome!
- Ask for help with setup or understanding codebase
- Pair programming sessions available for complex features
- Regular office hours for questions

## 🎯 Project Priorities

### Current Focus (v0.1-0.2)
- Core command logging functionality
- Basic CLI interface
- SQLite database setup
- Shell integration (Zsh/Bash)

### Near-term Priorities (v0.3-0.5)
- Gamification system (XP, levels, badges)
- Rich terminal UI
- Advanced statistics
- Configuration management

### Contribution Areas
**High Priority:**
- Shell hook optimization
- Database performance
- Error handling improvements
- Test coverage expansion

**Medium Priority:**
- CLI UX enhancements
- Documentation improvements
- Additional shell support
- Configuration options

**Nice to Have:**
- Advanced visualizations
- Social features
- Plugin system
- Mobile companion app

## 📈 Success Metrics

### Quality Indicators
- Test coverage > 80%
- No critical bugs in releases
- Performance benchmarks met
- Documentation completeness

### Community Health
- Response time to issues < 48 hours
- Active contributor growth
- Positive community interactions
- Regular release cadence

---

Thank you for contributing to Termonaut! Your efforts help make terminal productivity more engaging for developers worldwide. 🚀

For more detailed technical information, see our [Development Guide](DEVELOPMENT.md).