# GitHub Actions自动化发布系统

## 🚀 概述

Termonaut现在使用完全自动化的GitHub Actions工作流来构建和发布跨平台二进制文件。这消除了手动构建的复杂性，确保所有平台的一致性和可靠性。

## 🌍 支持的平台

| 平台 | 架构 | 状态 |
|------|------|------|
| macOS | Intel (x64) | ✅ 完全支持 |
| macOS | Apple Silicon (ARM64) | ✅ 完全支持 |
| Linux | x86_64 | ✅ 完全支持 |
| Linux | ARM64 | ✅ 完全支持 |
| Windows | x86_64 | ✅ 完全支持 |

## 📋 工作流文件

### 1. 主发布工作流 (`.github/workflows/release.yml`)

**触发条件:**
- 推送git标签 (`v*`)
- 手动触发 (workflow_dispatch)

**功能:**
- 🏗️ 跨平台二进制构建 (5个平台)
- 🧪 自动化测试执行
- 📦 GitHub Release创建
- 🔐 SHA256校验和生成
- 📝 自动生成发布说明

### 2. Homebrew更新工作流 (`.github/workflows/update-homebrew.yml`)

**触发条件:**
- GitHub Release发布后自动触发
- 手动触发更新特定版本

**功能:**
- 🍺 自动更新外部Homebrew tap (`oiahoon/homebrew-termonaut`)
- 📁 同步更新本地Formula作为备份
- 🔄 提交和推送更改到两个位置
- ✅ Formula语法验证
- 🛡️ 失败时使用本地Formula作为后备

### 3. 手动发布工作流 (`.github/workflows/manual-release.yml`)

**触发条件:**
- 手动触发，用于快速发布

**功能:**
- 🔖 更新代码中的版本号
- 🏷️ 创建和推送git标签
- 🔄 自动触发主发布工作流

## 🍺 Homebrew集成方案

### 方案A: 自动更新外部Tap (推荐)

如果你有 `oiahoon/homebrew-termonaut` 仓库:

**优势:**
- ✅ 标准的Homebrew tap结构
- ✅ 用户安装简单: `brew install oiahoon/termonaut/termonaut`
- ✅ 完全自动化更新
- ✅ 符合Homebrew最佳实践

**工作流程:**
1. 发布新版本时自动触发
2. 更新 `oiahoon/homebrew-termonaut/termonaut.rb`
3. 同时更新本地 `Formula/termonaut.rb` 作为备份
4. 用户通过标准tap安装

### 方案B: 本地Formula (备选)

如果没有外部tap仓库:

**特点:**
- 📁 使用项目内的 `Formula/termonaut.rb`
- 👤 用户需要: `brew install Formula/termonaut.rb`
- 🔄 仍然自动更新
- ⚠️ 稍微不如标准tap方便

### 🛠️ 设置Homebrew集成

运行设置脚本:
```bash
./scripts/setup-homebrew-integration.sh
```

这个脚本会:
- 🔍 检查你的homebrew-termonaut仓库状态
- 📋 显示可用的集成选项
- 🧪 提供测试功能
- 📖 给出详细的设置指导

## 🛠️ 如何发布新版本

### 方法1: 一键发布 (推荐)

1. 访问 [GitHub Actions](https://github.com/oiahoon/termonaut/actions)
2. 选择 "Manual Release" 工作流
3. 点击 "Run workflow"
4. 输入版本号 (例如: `0.9.5`)
5. 确认选项并点击 "Run workflow"

**这将自动:**
- ✅ 更新 `cmd/termonaut/main.go` 中的版本
- ✅ 创建git标签 `v0.9.5`
- ✅ 触发完整的发布流程
- ✅ 构建所有平台的二进制文件
- ✅ 创建GitHub Release
- ✅ 更新Homebrew formula (外部tap + 本地备份)

### 方法2: 传统git标签方式

```bash
# 1. 更新版本号
sed -i 's/version = "[^"]*"/version = "0.9.5"/' cmd/termonaut/main.go

# 2. 提交更改
git add cmd/termonaut/main.go
git commit -m "🔖 Bump version to 0.9.5"
git push

# 3. 创建和推送标签
git tag -a v0.9.5 -m "Release v0.9.5"
git push origin v0.9.5
```

### 方法3: 直接触发发布工作流

1. 访问 [Release Workflow](https://github.com/oiahoon/termonaut/actions/workflows/release.yml)
2. 点击 "Run workflow"
3. 输入版本号 (例如: `v0.9.5`)
4. 点击 "Run workflow"

## 📊 构建矩阵详情

```yaml
strategy:
  matrix:
    include:
      # macOS builds
      - os: macos-latest
        goos: darwin
        goarch: amd64
        name: darwin-amd64
      - os: macos-latest
        goos: darwin
        goarch: arm64
        name: darwin-arm64

      # Linux builds
      - os: ubuntu-latest
        goos: linux
        goarch: amd64
        name: linux-amd64
      - os: ubuntu-latest
        goos: linux
        goarch: arm64
        name: linux-arm64

      # Windows builds
      - os: windows-latest
        goos: windows
        goarch: amd64
        name: windows-amd64
        ext: .exe
```

## 🔧 技术细节

### CGO处理
- **macOS**: 原生CGO支持
- **Linux**: 使用交叉编译工具链
  - x64: `gcc-multilib`
  - ARM64: `gcc-aarch64-linux-gnu`
- **Windows**: 原生支持

### 二进制命名规范
```
termonaut-{version}-{platform}-{arch}[.exe]
```

例如:
- `termonaut-0.9.5-darwin-amd64`
- `termonaut-0.9.5-linux-arm64`
- `termonaut-0.9.5-windows-amd64.exe`

### 构建标志
```bash
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE}"
```

### Homebrew集成技术细节

**外部Tap更新流程:**
1. Checkout主仓库获取release信息
2. Checkout `oiahoon/homebrew-termonaut` 仓库
3. 下载并计算SHA256校验和
4. 更新 `termonaut.rb` formula
5. 提交并推送到外部tap
6. 同步更新本地Formula作为备份

**权限配置:**
- 使用 `HOMEBREW_TAP_TOKEN` (推荐) 或 `GITHUB_TOKEN`
- 需要对homebrew-termonaut仓库的写权限

## 🔍 监控和调试

### 查看构建状态
- [Actions页面](https://github.com/oiahoon/termonaut/actions)
- [Releases页面](https://github.com/oiahoon/termonaut/releases)

### 常见问题

**Q: 构建失败怎么办？**
A: 检查Actions日志，通常是依赖问题或代码语法错误

**Q: 某个平台的二进制文件缺失？**
A: 检查构建矩阵中对应平台的构建日志

**Q: Homebrew formula没有自动更新？**
A: 检查update-homebrew工作流的执行状态，可能需要配置HOMEBREW_TAP_TOKEN

**Q: 如何添加新平台支持？**
A: 在`.github/workflows/release.yml`的构建矩阵中添加新的平台配置

**Q: homebrew-termonaut仓库访问失败？**
A: 确保仓库存在且公开，或配置HOMEBREW_TAP_TOKEN密钥

## 🎯 优势

### 相比手动发布
- ✅ **一致性**: 所有平台使用相同的构建环境
- ✅ **可靠性**: 自动化减少人为错误
- ✅ **效率**: 并行构建，节省时间
- ✅ **可追溯**: 完整的构建日志和工件
- ✅ **自动化**: 从版本更新到发布的完整流程

### 相比本地Docker构建
- ✅ **原生性能**: 每个平台在对应的原生环境构建
- ✅ **无环境依赖**: 不需要本地Docker或交叉编译工具
- ✅ **并行构建**: 5个平台同时构建
- ✅ **集成化**: 与GitHub生态系统深度集成

### Homebrew集成优势
- ✅ **双重保障**: 外部tap + 本地formula备份
- ✅ **标准化**: 符合Homebrew最佳实践
- ✅ **用户友好**: 简单的安装命令
- ✅ **自动维护**: 无需手动更新多个仓库

## 🚀 未来改进

- [ ] 添加更多Linux发行版支持
- [ ] 集成自动化测试覆盖率报告
- [ ] 添加性能基准测试
- [ ] 支持预发布版本 (beta/rc)
- [ ] 集成安全扫描
- [ ] 添加更多包管理器支持 (apt, yum, etc.)
- [ ] 支持多个Homebrew tap同时更新

## 📞 支持

如果在使用自动化发布系统时遇到问题:

1. 查看 [GitHub Actions日志](https://github.com/oiahoon/termonaut/actions)
2. 运行 `./scripts/setup-homebrew-integration.sh` 检查配置
3. 检查 [Issues页面](https://github.com/oiahoon/termonaut/issues)
4. 创建新的Issue描述问题

---

**Happy automated releasing! 🎉**