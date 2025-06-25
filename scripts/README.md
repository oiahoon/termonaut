# 📜 Termonaut 脚本中心

这里包含了 Termonaut 项目的所有脚本文件，按功能进行了分类整理。

## 📁 目录结构

```
scripts/
├── README.md              # 本文档
├── build/                 # 🔨 构建脚本
├── install/               # 📦 安装脚本
├── test/                  # 🧪 测试脚本
├── maintenance/           # 🔧 维护脚本
└── archive/               # 📚 归档脚本
    ├── releases/          # 历史发布脚本
    └── deprecated/        # 废弃脚本
```

## 🔨 构建脚本 (build/)

用于项目构建和发布的脚本：

- `build-release.sh` - 构建发布版本
- `create-github-release.sh` - 创建 GitHub 发布

**使用方法**:
```bash
cd scripts/build
./build-release.sh
```

## 📦 安装脚本 (install/)

用于项目安装和部署的脚本：

- `safe-shell-install.sh` - 安全的 Shell 集成安装
- `verify-install.sh` - 验证安装是否成功
- `setup-homebrew-integration.sh` - 设置 Homebrew 集成
- `update-homebrew-formula.sh` - 更新 Homebrew 配方
- `release-homebrew.sh` - 发布到 Homebrew

**使用方法**:
```bash
cd scripts/install
./safe-shell-install.sh
```

## 🧪 测试脚本 (test/)

用于测试项目功能的脚本：

- `test-tui-layout.sh` - 测试 TUI 布局
- `test-homebrew-formula-generation.sh` - 测试 Homebrew 配方生成
- 其他测试脚本...

**使用方法**:
```bash
cd scripts/test
./test-tui-layout.sh
```

## 🔧 维护脚本 (maintenance/)

用于项目维护的脚本：

- `verify-release.sh` - 验证发布版本
- `manual-release-info.sh` - 手动发布信息

**使用方法**:
```bash
cd scripts/maintenance
./verify-release.sh
```

## 📚 归档脚本 (archive/)

### 历史发布脚本 (archive/releases/)
包含各个版本的发布脚本，仅供参考：
- `release-0.9.0-rc.sh`
- `release-0.9.0.sh`
- `release-0.9.1.sh`
- `release-0.9.2.sh`
- `release-0.9.4.sh`
- `release-0.10.0.sh`
- `re-release-v0.10.1.sh`

### 废弃脚本 (archive/deprecated/)
已废弃的脚本，不建议使用：
- `fix-documentation.sh`
- `fix-tui-layout.sh`
- `fix-homebrew-tap.sh`

## 🚀 快速使用

### 开发者常用脚本
```bash
# 构建项目
./scripts/build/build-release.sh

# 安装到本地
./scripts/install/safe-shell-install.sh

# 验证安装
./scripts/install/verify-install.sh

# 运行测试
./scripts/test/test-tui-layout.sh
```

### 维护者常用脚本
```bash
# 验证发布
./scripts/maintenance/verify-release.sh

# 更新 Homebrew
./scripts/install/update-homebrew-formula.sh
```

## ⚠️ 注意事项

1. **执行权限**: 确保脚本有执行权限
   ```bash
   chmod +x scripts/category/script-name.sh
   ```

2. **工作目录**: 大部分脚本需要在项目根目录执行
   ```bash
   cd /path/to/termonaut
   ./scripts/category/script-name.sh
   ```

3. **依赖检查**: 某些脚本可能需要特定的依赖工具，请查看脚本内容了解要求

4. **归档脚本**: `archive/` 目录下的脚本仅供参考，不建议在生产环境使用

## 📝 脚本开发规范

如果你要添加新脚本，请遵循以下规范：

1. **命名规范**: 使用小写字母和连字符，如 `build-project.sh`
2. **分类放置**: 根据功能放入相应的子目录
3. **添加说明**: 在脚本开头添加功能说明注释
4. **错误处理**: 添加适当的错误处理和退出码
5. **更新文档**: 在本 README 中添加脚本说明

## 🤝 贡献

如果你想改进现有脚本或添加新脚本，请：

1. 查看 [贡献指南](../docs/CONTRIBUTING.md)
2. 遵循脚本开发规范
3. 提交 Pull Request

---

**💡 提示**: 如果脚本执行遇到问题，请检查：
- 是否有执行权限
- 是否在正确的工作目录
- 是否满足脚本的依赖要求
