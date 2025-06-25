# 🎭 浮动彩蛋通知系统 - 可行性研究

## 🎯 目标效果
在用户触发彩蛋时，在shell窗口顶部显示浮动通知，3-5秒后自动消失，类似桌面通知系统。

## 🔍 技术方案分析

### 方案1: ANSI转义序列实现 ⭐⭐⭐
```go
type FloatingNotification struct {
    message     string
    duration    time.Duration
    position    Position
    style       NotificationStyle
}

func (fn *FloatingNotification) Show() {
    // 1. 保存当前终端状态
    fmt.Print("\033[s")           // 保存光标位置
    fmt.Print("\033[?25l")        // 隐藏光标
    
    // 2. 移动到顶部显示通知
    fmt.Printf("\033[1;1H")       // 移动到第1行第1列
    fmt.Print("\033[K")           // 清除当前行
    
    // 3. 显示彩蛋通知
    styledMessage := fn.applyStyle(fn.message)
    fmt.Print(styledMessage)
    
    // 4. 设置自动消失
    go func() {
        time.Sleep(fn.duration)
        fn.Hide()
    }()
}

func (fn *FloatingNotification) Hide() {
    fmt.Printf("\033[1;1H")       // 移动到通知位置
    fmt.Print("\033[K")           // 清除通知行
    fmt.Print("\033[u")           // 恢复光标位置
    fmt.Print("\033[?25h")        // 显示光标
}
```

**优点:**
- ✅ 跨平台兼容
- ✅ 无外部依赖
- ✅ 实现简单
- ✅ 可定制样式

**缺点:**
- ❌ 可能干扰用户输入
- ❌ 与其他程序输出冲突
- ❌ 不是真正的浮动层

### 方案2: 终端特定功能 ⭐⭐⭐⭐
```go
type TerminalNotifier interface {
    ShowNotification(message string, duration time.Duration) error
    IsSupported() bool
}

// iTerm2 实现
type ITerm2Notifier struct{}
func (i *ITerm2Notifier) ShowNotification(msg string, duration time.Duration) error {
    // 使用 iTerm2 的 OSC 序列
    osc := fmt.Sprintf("\033]9;%s\007", msg)
    fmt.Print(osc)
    return nil
}

// Warp 实现
type WarpNotifier struct{}
func (w *WarpNotifier) ShowNotification(msg string, duration time.Duration) error {
    // Warp 支持的通知方式
    return exec.Command("warp-cli", "notify", msg).Run()
}

// VS Code 终端实现
type VSCodeNotifier struct{}
func (v *VSCodeNotifier) ShowNotification(msg string, duration time.Duration) error {
    // VS Code 终端的通知API
    vscode := fmt.Sprintf("\033]633;A\033\\%s\033]633;B\033\\", msg)
    fmt.Print(vscode)
    return nil
}
```

**优点:**
- ✅ 真正的通知效果
- ✅ 不干扰用户操作
- ✅ 原生终端集成
- ✅ 更好的用户体验

**缺点:**
- ❌ 终端兼容性限制
- ❌ 需要检测终端类型
- ❌ 功能可能有限

### 方案3: 混合实现 ⭐⭐⭐⭐⭐
```go
type SmartEasterEggNotifier struct {
    notifiers []TerminalNotifier
    fallback  *ANSINotifier
}

func (sen *SmartEasterEggNotifier) ShowEasterEgg(egg EasterEgg) {
    // 1. 尝试现代终端通知
    for _, notifier := range sen.notifiers {
        if notifier.IsSupported() {
            if err := notifier.ShowNotification(egg.Message, egg.Duration); err == nil {
                return // 成功显示，退出
            }
        }
    }
    
    // 2. 降级到ANSI实现
    sen.fallback.ShowNotification(egg.Message, egg.Duration)
}
```

## 🎨 视觉效果设计

### 通知样式
```go
type NotificationStyle struct {
    Background  Color
    Foreground  Color
    Border      BorderStyle
    Animation   AnimationType
    Position    Position
}

// 彩蛋主题样式
var EasterEggStyles = map[string]NotificationStyle{
    "celebration": {
        Background: ColorYellow,
        Foreground: ColorBlack,
        Border:     BorderDouble,
        Animation:   AnimationBounce,
        Position:    PositionTop,
    },
    "achievement": {
        Background: ColorGreen,
        Foreground: ColorWhite,
        Border:     BorderRounded,
        Animation:   AnimationFade,
        Position:    PositionTop,
    },
    "humor": {
        Background: ColorCyan,
        Foreground: ColorBlack,
        Border:     BorderDashed,
        Animation:   AnimationSlide,
        Position:    PositionTop,
    },
}
```

### 动画效果
```go
func (fn *FloatingNotification) ShowWithAnimation() {
    switch fn.style.Animation {
    case AnimationBounce:
        fn.bounceIn()
    case AnimationFade:
        fn.fadeIn()
    case AnimationSlide:
        fn.slideDown()
    default:
        fn.Show()
    }
}

func (fn *FloatingNotification) bounceIn() {
    positions := []int{-2, -1, 0, 1, 0} // 弹跳效果
    for _, pos := range positions {
        fmt.Printf("\033[%d;1H", 1+pos)
        fmt.Print(fn.applyStyle(fn.message))
        time.Sleep(100 * time.Millisecond)
        fmt.Print("\033[K") // 清除
    }
}
```

## 🔧 实现挑战与解决方案

### 挑战1: 终端兼容性
```go
type TerminalDetector struct {
    cache map[string]bool
}

func (td *TerminalDetector) DetectTerminal() TerminalType {
    // 检测环境变量
    if term := os.Getenv("TERM_PROGRAM"); term != "" {
        switch term {
        case "iTerm.app":
            return TerminalITerm2
        case "vscode":
            return TerminalVSCode
        case "Warp":
            return TerminalWarp
        }
    }
    
    // 检测终端特性
    if td.supportsOSC() {
        return TerminalModern
    }
    
    return TerminalBasic
}
```

### 挑战2: 用户输入干扰
```go
type SafeNotifier struct {
    inputBuffer []byte
    cursorPos   int
}

func (sn *SafeNotifier) ShowSafely(message string) {
    // 1. 保存用户输入状态
    sn.saveInputState()
    
    // 2. 显示通知
    sn.showNotification(message)
    
    // 3. 恢复用户输入状态
    defer sn.restoreInputState()
}
```

### 挑战3: 多行内容处理
```go
func (fn *FloatingNotification) calculateLayout(message string) Layout {
    lines := strings.Split(message, "\n")
    maxWidth := 0
    
    for _, line := range lines {
        if width := runewidth.StringWidth(line); width > maxWidth {
            maxWidth = width
        }
    }
    
    return Layout{
        Width:  maxWidth + 4, // 边框和填充
        Height: len(lines) + 2, // 边框
        Lines:  lines,
    }
}
```

## 🎯 用户体验考虑

### 配置选项
```toml
[easter_eggs]
floating_notifications = true
notification_duration = 3000  # 毫秒
notification_position = "top" # top, bottom, center
animation_style = "bounce"    # bounce, fade, slide, none
max_notifications = 1         # 同时显示的最大数量
```

### 智能触发
```go
type SmartTrigger struct {
    lastShown    time.Time
    cooldown     time.Duration
    userActivity ActivityLevel
}

func (st *SmartTrigger) ShouldShow(egg EasterEgg) bool {
    // 1. 冷却时间检查
    if time.Since(st.lastShown) < st.cooldown {
        return false
    }
    
    // 2. 用户活动检查
    if st.userActivity == ActivityHigh {
        return false // 用户忙碌时不显示
    }
    
    // 3. 终端状态检查
    if st.isUserTyping() {
        return false // 用户正在输入时不显示
    }
    
    return true
}
```

## 📊 可行性评估

| 方案 | 兼容性 | 用户体验 | 实现复杂度 | 推荐度 |
|------|--------|----------|------------|--------|
| ANSI转义 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| 终端特定 | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 混合方案 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## 🚀 推荐实现路径

### 阶段1: 基础ANSI实现
- 实现基本的顶部通知显示
- 添加简单的样式和颜色
- 处理基本的冲突情况

### 阶段2: 终端检测和优化
- 添加终端类型检测
- 实现现代终端的原生通知
- 优化用户输入保护

### 阶段3: 高级功能
- 添加动画效果
- 实现智能触发逻辑
- 添加用户配置选项

## 💡 结论

**完全可行！** 虽然有技术挑战，但通过混合实现方案可以达到很好的效果：

1. **现代终端**: 使用原生通知API，体验最佳
2. **传统终端**: 使用ANSI转义序列，兼容性最好
3. **智能降级**: 自动选择最佳实现方式

这个功能将大大提升Termonaut的用户体验，让彩蛋真正"浮动"起来！🎉
