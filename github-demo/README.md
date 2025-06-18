# 🚀 Termonaut GitHub Integration Demo

这是一个展示Termonaut GitHub集成功能的演示页面。

## 📊 实时统计 Badges

![Commands](https://img.shields.io/badge/Commands-93-lightgrey?style=flat-square&logo=terminal&logoColor=white) ![Streak](https://img.shields.io/badge/Streak-2+days-red?style=flat-square&logo=terminal&logoColor=white) ![Productivity](https://img.shields.io/badge/Productivity-80.0%25-green?style=flat-square&logo=terminal&logoColor=white) ![Achievements](https://img.shields.io/badge/Achievements-5%2F10-blue?style=flat-square&logo=terminal&logoColor=white) ![Last Active](https://img.shields.io/badge/Last+Active-4h+ago-green?style=flat-square&logo=terminal&logoColor=white) ![XP](https://img.shields.io/badge/XP-Level+3+%28779%2F1600%29-lightgrey?style=flat-square&logo=terminal&logoColor=white)

## 🎮 我的终端游戏化之旅

### 📈 当前状态
- **等级**: 3 级 (XP: 779/1600)
- **总命令数**: 93 条
- **独特命令**: 63 种
- **当前连击**: 2 天
- **最长连击**: 2 天
- **今日命令**: 36 条

### 🏆 解锁成就
- ✅ **First Steps**: 执行第一个命令
- ⏳ **Command Master**: 执行1000个命令 (进行中...)

### 🔥 最常用命令
1. `curl -s example.com` (5 次) ████████████████████
2. `source ~/.zshrc` (4 次) ████████████████
3. `./termonaut stats` (3 次) ████████████
4. `git status` (3 次) ████████████
5. `gst` (3 次) ████████████

## 🛠️ 功能特性

### ✨ GitHub集成功能
- **动态Badges**: 实时显示终端使用统计
- **Profile生成**: 自动生成个人档案
- **Heatmap可视化**: GitHub风格的活动热力图
- **Actions自动化**: 自动更新统计数据

### 🎯 生产力追踪
- **命令分类**: 自动识别命令类型
- **时间分析**: 识别最佳工作时间
- **连击系统**: 鼓励持续使用
- **隐私保护**: 敏感信息自动脱敏

## 📊 活动热力图

```
🔥 Weekly Productivity Heatmap
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

       00  03  06  09  12  15  18  21
Mon
Tue                       🔥
Wed                   🔥
Thu
Fri
Sat
Sun

Legend: ░ Light  ▓ Medium  █ High  🔥 Peak
Focus Score: 60.8/100 (Well Focused)
```

## 🚀 快速开始

### 1. 安装Termonaut
```bash
# 使用install脚本
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash

# 或下载二进制文件
sudo cp termonaut-0.9.0-rc2-darwin-amd64 /usr/local/bin/termonaut
```

### 2. 初始化
```bash
# 初始化shell集成
tn init

# 重新加载shell配置
source ~/.zshrc  # 或 ~/.bashrc
```

### 3. 开始使用
```bash
# 查看统计信息
tn stats

# 生成GitHub badges
tn github badges generate

# 生成个人档案
tn github profile generate

# 查看活动热力图
tn heatmap
```

## 🔧 配置选项

```bash
# 启用emoji主题
tn config set theme emoji

# 启用空命令统计
tn config set empty_command_stats true

# 查看所有配置
tn config get
```

## 📁 文件导出

### Badges JSON
```json
{
  "Commands": "https://img.shields.io/badge/Commands-93-lightgrey?style=flat-square&logo=terminal&logoColor=white",
  "Streak": "https://img.shields.io/badge/Streak-2+days-red?style=flat-square&logo=terminal&logoColor=white",
  "Productivity": "https://img.shields.io/badge/Productivity-80.0%25-green?style=flat-square&logo=terminal&logoColor=white"
}
```

### Profile Markdown
完整的个人档案可以导出为Markdown格式，包含：
- 实时统计badges
- 详细使用概览
- 成就解锁状态
- 最常用命令排行

## 🤖 GitHub Actions集成

### 自动更新工作流
```yaml
name: Update Termonaut Stats
on:
  schedule:
    - cron: '0 */6 * * *'  # 每6小时更新一次
  workflow_dispatch:

jobs:
  update-stats:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Generate Badge Data
      run: |
        mkdir -p badges
        termonaut github badges generate --format json --output badges/
    - name: Commit changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add badges/
        git commit -m "🚀 Update Termonaut stats" || exit 0
        git push
```

## 🎨 自定义样式

Badges支持多种样式：
- `flat-square` (默认)
- `flat`
- `plastic`
- `for-the-badge`
- `social`

颜色主题：
- 🟢 绿色：高生产力、活跃状态
- 🔴 红色：低连击、需要改进
- 🔵 蓝色：成就、特殊状态
- ⚪ 灰色：中性数据

## 📈 数据洞察

### 生产力指标
- **命令多样性**: 使用不同命令的比例
- **时间分布**: 一天中的活跃时段
- **一致性评分**: 使用习惯的稳定性
- **专业化程度**: 特定领域的深度

### 趋势分析
- 每日/每周/每月命令趋势
- 最佳工作时间识别
- 命令类别偏好分析
- 生产力模式识别

## 🔗 相关链接

- [Termonaut GitHub仓库](https://github.com/oiahoon/termonaut)
- [完整文档](https://github.com/oiahoon/termonaut/blob/main/README.md)
- [安装指南](https://github.com/oiahoon/termonaut/blob/main/docs/QUICK_START.md)
- [故障排除](https://github.com/oiahoon/termonaut/blob/main/docs/TROUBLESHOOTING.md)

---

*由 [Termonaut](https://github.com/oiahoon/termonaut) 自动生成 - 终端生产力追踪工具*
*最后更新: 2025年6月18日*