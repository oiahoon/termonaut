# 🚀 Termonaut 快速开始指南

欢迎使用 Termonaut！这个指南将帮助你快速上手，开始你的终端生产力追踪之旅。

## 📋 前提条件

确保你已经完成了 [安装](installation.md)。如果还没有，请先安装 Termonaut：

```bash
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

## 🎯 第一步：初始化设置

### 新用户推荐方式

```bash
# 交互式设置向导（推荐）
termonaut setup
```

这个向导会引导你完成：
- Shell 集成设置
- 基本配置选项
- 隐私设置
- 主题选择

### 快速设置

```bash
# 一键快速设置
termonaut quickstart
```

使用默认配置快速开始，稍后可以调整。

## 📊 第二步：查看你的第一个统计

```bash
# 查看今天的统计
termonaut stats

# 查看本周统计
termonaut stats --weekly

# 查看本月统计
termonaut stats --monthly
```

示例输出：
```
🚀 Today's Terminal Stats (2024-01-15)
─────────────────────────────────────
Commands Executed: 127 🎯
Active Time: 3h 42m ⏱️
Session Count: 4 📱
New Commands: 3 ⭐
Current Streak: 12 days 🔥

🎮 Level 8 Astronaut (2,150 XP)
Progress to Level 9: ████████████░░░░ 75%
```

## 🖥️ 第三步：探索交互式界面

```bash
# 启动交互式 TUI 界面
termonaut tui

# 或者使用短命令
tn tui
```

TUI 界面包含多个标签页：
- **📊 Stats**: 详细统计信息
- **🏆 Achievements**: 成就和徽章
- **📈 Trends**: 使用趋势分析
- **⚙️ Settings**: 配置选项

### TUI 模式选择

```bash
# 智能模式（自动适应终端大小）
termonaut tui --mode smart

# 紧凑模式（小终端）
termonaut tui --mode compact

# 完整模式（大终端）
termonaut tui --mode full

# 最小模式（纯文本）
termonaut tui --mode minimal
```

## 🎮 第四步：了解游戏化系统

### XP 和等级
- 每个命令都会获得 XP
- 新命令获得额外奖励
- 连续使用获得连击奖励
- 不同类别的命令有不同的 XP 倍数

### 成就系统
```bash
# 查看所有成就
termonaut tui  # 然后切换到 Achievements 标签
```

一些容易获得的成就：
- 🚀 **First Launch**: 执行第一个命令
- 🌟 **Explorer**: 使用 50 个不同命令
- 🔥 **Streak Keeper**: 保持 7 天使用连击

### 头像系统
你的头像会随着等级提升而进化！

## ⚙️ 第五步：个性化配置

```bash
# 查看当前配置
termonaut config get

# 启用表情符号主题
termonaut config set theme emoji

# 调整头像样式
termonaut config set avatar_style pixel-art

# 设置隐私模式
termonaut config set anonymous_mode true
```

### 配置文件位置
配置文件位于 `~/.termonaut/config.toml`，你可以直接编辑：

```toml
# 显示和主题
display_mode = "enter"
theme = "emoji"
show_gamification = true

# 追踪行为
idle_timeout_minutes = 10
track_git_repos = true
command_categories = true

# 隐私设置
anonymous_mode = false
opt_out_commands = ["password", "secret"]
```

## 🔄 第六步：GitHub 集成（可选）

如果你想在 GitHub 上展示你的终端统计：

```bash
# 设置 GitHub 同步
termonaut github sync setup

# 生成个人资料
termonaut github profile generate

# 生成动态徽章
termonaut github badges generate
```

然后你可以在 README 中添加徽章：
```markdown
![Commands](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/commands.json)
![Level](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/level.json)
```

## 💡 第七步：日常使用技巧

### 快速查看统计
```bash
# 在空命令行按回车查看快速统计（如果启用）
# 或者使用短命令
tn stats
```

### 使用别名
```bash
# 设置 'tn' 作为 'termonaut' 的别名
termonaut alias create

# 现在可以使用短命令
tn tui
tn stats
tn config get
```

### 空命令统计
如果启用了空命令统计功能，在终端中按回车（空命令）会显示快速统计。

## 🎯 常用命令速查

```bash
# 基本命令
termonaut stats              # 今日统计
termonaut tui                # 交互界面
termonaut config get         # 查看配置

# 设置命令
termonaut setup              # 设置向导
termonaut quickstart         # 快速设置
termonaut alias create       # 创建别名

# GitHub 集成
termonaut github sync now    # 立即同步
termonaut github profile     # 生成个人资料

# 高级功能
termonaut advanced analytics # 高级分析
termonaut advanced api       # API 服务器
```

## 🔍 探索更多功能

### 高级分析
```bash
# 启动高级分析
termonaut advanced analytics

# 查看生产力趋势
termonaut advanced productivity
```

### API 服务器
```bash
# 启动 API 服务器（用于集成）
termonaut advanced api --port 8080
```

### 数据导出
```bash
# 导出数据（计划中的功能）
termonaut export --format json
termonaut export --format csv
```

## 🎉 你已经准备好了！

现在你已经掌握了 Termonaut 的基本使用方法。继续使用终端，观察你的统计数据增长，解锁新的成就，提升你的等级！

## 📚 下一步阅读

- [配置指南](configuration.md) - 详细的配置选项
- [功能文档](../features/) - 了解所有功能
- [故障排除](../TROUBLESHOOTING.md) - 解决常见问题
- [贡献指南](../CONTRIBUTING.md) - 参与项目开发

## 💬 获取帮助

- 查看内置帮助：`termonaut --help`
- 访问 [GitHub Issues](https://github.com/oiahoon/termonaut/issues)
- 阅读 [故障排除指南](../TROUBLESHOOTING.md)

---

**🚀 开始你的终端生产力追踪之旅吧！每个命令都是向精通迈进的一步。**
