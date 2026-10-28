# 📦 Termonaut 安装指南

本指南将帮助你在各种平台上安装 Termonaut。

## 🚀 快速安装

### 自动安装脚本（推荐）

```bash
# 一键安装（支持 macOS 和 Linux）
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

这个脚本会：
- 自动检测你的操作系统和架构
- 下载适合的二进制文件
- 安装到系统路径
- 设置必要的权限

### Homebrew 安装（macOS）

```bash
# 添加我们的 tap
brew tap oiahoon/termonaut

# 安装 Termonaut
brew install termonaut

# 或者从 homebrew-core 安装（即将支持）
brew install termonaut
```

## 🔧 手动安装

### 1. 下载二进制文件

访问 [GitHub Releases](https://github.com/oiahoon/termonaut/releases/latest) 页面，下载适合你系统的版本：

- **macOS Intel**: `termonaut-darwin-amd64`
- **macOS Apple Silicon**: `termonaut-darwin-arm64`
- **Linux x64**: `termonaut-linux-amd64`
- **Linux ARM64**: `termonaut-linux-arm64`
- **Windows x64**: `termonaut-windows-amd64.exe`

### 2. 安装到系统路径

#### macOS/Linux:
```bash
# 下载后重命名并移动到系统路径
chmod +x termonaut-*
sudo mv termonaut-* /usr/local/bin/termonaut

# 验证安装
termonaut version
```

#### Windows:
```powershell
# 将 termonaut-windows-amd64.exe 重命名为 termonaut.exe
# 移动到 PATH 中的目录，如 C:\Windows\System32\
# 或者添加到用户 PATH 环境变量
```

## ⚙️ 初始化设置

安装完成后，需要设置 shell 集成：

### 新用户（推荐）

```bash
# 交互式设置向导
termonaut setup

# 或者快速设置
termonaut quickstart
```

### 手动设置

```bash
# 安装 shell 钩子
termonaut init

# 重新加载 shell 配置
source ~/.bashrc  # 或 ~/.zshrc
```

## 🔍 验证安装

```bash
# 检查版本
termonaut version

# 查看帮助
termonaut --help

# 测试基本功能
termonaut stats

# 启动交互界面
termonaut tui
```

## 🛠️ 高级安装选项

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/oiahoon/termonaut.git
cd termonaut

# 构建
go build -o termonaut cmd/termonaut/*.go

# 安装
sudo mv termonaut /usr/local/bin/
```

### Docker 使用

```bash
# 拉取镜像（如果有的话）
docker pull oiahoon/termonaut:latest

# 运行
docker run -it --rm oiahoon/termonaut:latest
```

## 🔧 安装脚本选项

我们的安装脚本支持多种选项：

```bash
# 安全模式安装（推荐）
./scripts/install/safe-shell-install.sh

# 验证安装
./scripts/install/verify-install.sh

# Homebrew 集成设置
./scripts/install/setup-homebrew-integration.sh
```

## ❗ 故障排除

### 常见问题

1. **权限错误**
   ```bash
   # 确保有执行权限
   chmod +x termonaut
   ```

2. **找不到命令**
   ```bash
   # 检查 PATH
   echo $PATH
   
   # 添加到 PATH（临时）
   export PATH=$PATH:/usr/local/bin
   ```

3. **Shell 集成不工作**
   ```bash
   # 重新安装 shell 钩子
   termonaut init --force
   
   # 重新加载配置
   source ~/.bashrc
   ```

### 获取帮助

如果遇到安装问题：

1. 查看 [故障排除指南](../TROUBLESHOOTING.md)
2. 搜索 [GitHub Issues](https://github.com/oiahoon/termonaut/issues)
3. 提交新的 Issue

## 🔄 更新 Termonaut

### Homebrew 用户
```bash
brew update
brew upgrade termonaut
```

### 手动安装用户
```bash
# 重新运行安装脚本
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

## 🗑️ 卸载

```bash
# 删除二进制文件
sudo rm /usr/local/bin/termonaut

# 删除配置文件（可选）
rm -rf ~/.termonaut

# 移除 shell 钩子（手动编辑 ~/.bashrc 或 ~/.zshrc）
```

---

**🎉 安装完成！** 现在你可以开始使用 Termonaut 来追踪和游戏化你的终端使用了！

下一步：查看 [快速开始指南](quick-start.md) 了解基本使用方法。
