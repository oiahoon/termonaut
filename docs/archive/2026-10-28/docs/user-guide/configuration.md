# ⚙️ Termonaut 配置指南

本指南详细介绍了 Termonaut 的所有配置选项，帮助你个性化你的终端追踪体验。

## 📁 配置文件位置

Termonaut 的配置文件位于：
```
~/.termonaut/config.toml
```

如果文件不存在，Termonaut 会在首次运行时创建默认配置。

## 🔧 配置管理命令

```bash
# 查看所有配置
termonaut config get

# 查看特定配置
termonaut config get theme

# 设置配置值
termonaut config set theme emoji

# 重置为默认配置
termonaut config reset
```

## 📋 完整配置选项

### 显示和主题设置

```toml
[display]
# 显示模式：off, enter, ps1, floating
display_mode = "enter"

# 主题：minimal, emoji, ascii
theme = "emoji"

# 是否显示游戏化元素
show_gamification = true

# 默认 TUI 模式：smart, compact, full, minimal
default_tui_mode = "smart"
```

#### 显示模式说明
- `off`: 关闭所有显示
- `enter`: 空命令时显示统计（推荐）
- `ps1`: 在命令提示符中显示
- `floating`: 浮动显示（实验性）

#### 主题说明
- `minimal`: 纯文本，无装饰
- `emoji`: 使用表情符号（推荐）
- `ascii`: ASCII 艺术风格

### 追踪行为设置

```toml
[tracking]
# 会话超时时间（分钟）
idle_timeout_minutes = 10

# 是否追踪 Git 仓库信息
track_git_repos = true

# 是否自动分类命令
command_categories = true

# 是否追踪工作目录
track_working_directory = true

# 是否记录命令执行时间
track_execution_time = true
```

### 游戏化系统设置

```toml
[gamification]
# 是否启用 XP 系统
enable_xp = true

# 是否启用成就系统
enable_achievements = true

# 是否启用连击系统
enable_streaks = true

# 是否显示等级进度
show_level_progress = true

# 新命令 XP 倍数
new_command_multiplier = 2.0

# 连击 XP 倍数
streak_multiplier = 1.5
```

### 头像系统设置

```toml
[avatar]
# 头像样式：pixel-art, bottts, adventurer, avataaars
style = "pixel-art"

# 头像大小：small, medium, large, auto
size = "auto"

# 是否启用头像进化
enable_evolution = true

# 头像缓存时间（小时）
cache_duration = 24
```

### 隐私设置

```toml
[privacy]
# 匿名模式（隐藏个人路径）
anonymous_mode = false

# 要忽略的命令模式
opt_out_commands = ["password", "secret", "token"]

# 是否记录命令参数
record_arguments = true

# 是否记录环境变量
record_environment = false

# 数据保留天数（0 = 永久保留）
data_retention_days = 0
```

### GitHub 集成设置

```toml
[github]
# 是否启用同步
sync_enabled = false

# 同步仓库
sync_repo = "username/termonaut-profile"

# 同步频率：hourly, daily, weekly
sync_frequency = "daily"

# 徽章更新频率
badge_update_frequency = "daily"

# GitHub Token（可选，用于私有仓库）
# token = "ghp_xxxxxxxxxxxx"
```

### 高级设置

```toml
[advanced]
# 数据库路径
database_path = "~/.termonaut/termonaut.db"

# 日志级别：debug, info, warn, error
log_level = "info"

# 日志文件路径
log_file = "~/.termonaut/termonaut.log"

# API 服务器端口
api_port = 8080

# 是否启用 API 服务器
enable_api = false

# 备份频率：daily, weekly, monthly
backup_frequency = "weekly"
```

## 🎨 主题自定义

### 表情符号主题配置

```toml
[theme.emoji]
command_icon = "⚡"
session_icon = "📱"
time_icon = "⏱️"
streak_icon = "🔥"
level_icon = "🚀"
achievement_icon = "🏆"
```

### ASCII 主题配置

```toml
[theme.ascii]
use_colors = true
progress_char = "█"
empty_char = "░"
border_style = "rounded"
```

## 🔄 环境变量配置

你也可以通过环境变量来配置 Termonaut：

```bash
# 设置主题
export TERMONAUT_THEME=emoji

# 设置显示模式
export TERMONAUT_DISPLAY_MODE=enter

# 启用调试模式
export TERMONAUT_DEBUG=true

# 设置配置文件路径
export TERMONAUT_CONFIG_PATH=/custom/path/config.toml
```

## 📊 命令分类配置

Termonaut 自动将命令分为以下类别，每个类别有不同的 XP 倍数：

```toml
[categories]
[categories.git]
commands = ["git", "gh", "hub"]
xp_multiplier = 1.5
icon = "🔀"

[categories.docker]
commands = ["docker", "docker-compose", "podman"]
xp_multiplier = 1.3
icon = "🐳"

[categories.kubernetes]
commands = ["kubectl", "k9s", "helm"]
xp_multiplier = 1.4
icon = "☸️"

[categories.development]
commands = ["npm", "yarn", "pip", "cargo", "go"]
xp_multiplier = 1.2
icon = "💻"

[categories.system]
commands = ["ls", "cd", "pwd", "find", "grep"]
xp_multiplier = 1.0
icon = "🖥️"
```

## 🎯 配置示例

### 最小化配置（性能优先）

```toml
display_mode = "off"
theme = "minimal"
show_gamification = false
track_git_repos = false
command_categories = false
anonymous_mode = true
```

### 完整体验配置（功能丰富）

```toml
display_mode = "enter"
theme = "emoji"
show_gamification = true
track_git_repos = true
command_categories = true
enable_achievements = true
avatar_style = "pixel-art"
sync_enabled = true
```

### 隐私优先配置

```toml
anonymous_mode = true
record_arguments = false
record_environment = false
data_retention_days = 30
opt_out_commands = ["ssh", "scp", "curl", "wget", "password", "secret"]
```

## 🔧 配置验证

```bash
# 验证配置文件语法
termonaut config validate

# 查看配置文件路径
termonaut config path

# 备份当前配置
termonaut config backup

# 恢复配置备份
termonaut config restore
```

## 🚨 故障排除

### 配置文件损坏

```bash
# 重置为默认配置
termonaut config reset

# 或者手动删除配置文件
rm ~/.termonaut/config.toml
```

### 权限问题

```bash
# 修复配置目录权限
chmod 755 ~/.termonaut
chmod 644 ~/.termonaut/config.toml
```

### 配置不生效

```bash
# 重新加载配置
termonaut config reload

# 或者重启终端会话
```

## 📚 配置最佳实践

1. **逐步配置**: 先使用默认配置，然后逐步调整
2. **备份配置**: 在大幅修改前备份配置文件
3. **测试更改**: 修改后测试功能是否正常
4. **版本控制**: 可以将配置文件加入个人的 dotfiles 仓库

## 🔄 配置迁移

### 从旧版本升级

```bash
# Termonaut 会自动迁移旧配置
# 如果遇到问题，可以手动重置
termonaut config migrate
```

### 多设备同步

```bash
# 导出配置
termonaut config export > termonaut-config.toml

# 在另一台设备导入
termonaut config import < termonaut-config.toml
```

---

**💡 提示**: 配置更改会立即生效，无需重启 Termonaut。如果遇到问题，可以随时重置为默认配置。
