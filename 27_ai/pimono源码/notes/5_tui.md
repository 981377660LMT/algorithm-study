# TUI 包精读笔记

**这个 TUI 的本质是什么？**

一个围绕"行数组"的差分渲染引擎。组件 `render(width) → string[]`，TUI 逐行比较新旧输出，只重绘变化行。用 Synchronized Output 包裹写入消除闪烁。简单到不需要虚拟 DOM，但足够构建完整的终端 IDE。

**最值得学的 13 个洞见：**

1. **渲染模型**：行数组而非字符矩阵，差分粒度是行不是字符
2. **布局**：宽度向下传，高度向上汇报——一维流式布局
3. **Synchronized Output**（`CSI ? 2026 h/l`）是消除终端闪烁的根本解法
4. **宽度溢出是 TUI 的"段错误"**——溢出 = 自动换行 = 所有行号错位 = 全盘崩溃
5. **Kitty 协议双层解析**：先尝试现代协议，自动 fallback 传统序列，上层透明
6. **StdinBuffer 碎片重组** + **drainInput() SSH 泄漏防护**
7. **CURSOR_MARKER 用 APC 序列**：终端忽略但应用可检测，实现软硬光标分离
8. **AnsiCodeTracker 只在行尾重置 underline**（其他样式不泄漏到 padding）
9. **Overlay 合成是 "单行列替换"** + 样式 reset 隔离
10. **requestRender 用 nextTick 合并**——React 批量更新的同构思想
11. **VirtualTerminal 用 xterm.js headless** 做确定性测试
12. **Editor 的 kill ring / undo 用 structuredClone**——简单粗暴但有效
13. **fd 做文件补全** > Node fs 遍历——快且尊重 .gitignore

---

> `@mariozechner/pi-tui` — 一个用于 AI CLI 工具 (pi) 的终端 UI 框架
> 核心理念：**差分渲染 (differential rendering)** + **组件化** + **协议级键盘适配**

---

## 一、架构全景

```
┌────────────────────────────────────────────────────┐
│  TUI (extends Container)                           │
│  ┌──────────────────────────────────────────────┐  │
│  │  Component 树 (children)                     │  │
│  │  ┌─────┐ ┌────┐ ┌──────────┐ ┌───────┐      │  │
│  │  │Text │ │Box │ │Markdown  │ │Editor │ ...   │  │
│  │  └─────┘ └────┘ └──────────┘ └───────┘      │  │
│  └──────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────┐  │
│  │  Overlay 栈 (modal 层叠窗口)                 │  │
│  │  SelectList / SettingsList / ...              │  │
│  └──────────────────────────────────────────────┘  │
│  差分渲染引擎 (doRender)                           │
│  Focus 管理 / Input Pipeline                       │
├────────────────────────────────────────────────────┤
│  Terminal 抽象层                                    │
│  ├── ProcessTerminal (真实 stdin/stdout)            │
│  └── VirtualTerminal (xterm.js headless, 测试用)   │
├────────────────────────────────────────────────────┤
│  基础设施                                           │
│  ├── keys.ts        Kitty 协议 + 传统序列解析       │
│  ├── stdin-buffer.ts  输入分割/缓冲                │
│  ├── utils.ts       visibleWidth / ANSI 处理       │
│  ├── terminal-image.ts  Kitty/iTerm2 图片协议      │
│  ├── keybindings.ts     动作 → 快捷键映射          │
│  ├── kill-ring.ts       Emacs 风格剪切环           │
│  └── undo-stack.ts      structuredClone 撤销栈     │
└────────────────────────────────────────────────────┘
```

---

## 二、核心接口 `Component`

```ts
interface Component {
  render(width: number): string[] // 给定宽度 → 输出行数组
  handleInput?(data: string): void // 可选：处理键盘输入
  invalidate(): void // 清除缓存，下次重新渲染
  wantsKeyRelease?: boolean // 是否接收 Kitty key release
}
```

**洞见 1：渲染模型是 "行数组"，不是字符矩阵。**
每个组件返回 `string[]`，每条 string 是一行（可含 ANSI 转义）。TUI 只需逐行对比就能做差分。
这比 blessed/ink 的"虚拟 DOM → 字符矩阵"模型简单得多，但足以构建完整的 CLI IDE。

**洞见 2：宽度向下传递，高度向上汇报。**
组件接收 `width` 参数但不接收 `height`——高度由渲染结果的行数隐式决定。
这是一维流式布局，与浏览器 CSS 的 block flow 本质相同。

---

## 三、差分渲染引擎 (`doRender`)

这是 TUI 最核心的算法。渲染流程：

```
render(width)          → 组件树输出 newLines: string[]
compositeOverlays()    → 在 newLines 上合成 overlay 层
extractCursorPosition()→ 找 CURSOR_MARKER 并剥离
applyLineResets()      → 每行尾追加样式重置
diff(prev, new)        → 找到 firstChanged / lastChanged
输出变化行             → 用 CSI 序列移动光标，逐行覆写
```

### 关键优化策略

| 场景     | 策略                                           |
| -------- | ---------------------------------------------- |
| 首次渲染 | 直接输出所有行（不 clear）                     |
| 宽度变化 | 全量 clear + 重绘（换行点全变）                |
| 内容缩短 | 可选 clearOnShrink（清除残留行）               |
| 局部变化 | **只重绘 firstChanged → lastChanged 之间的行** |
| 只追加   | 从 previousLines.length 处开始写               |
| 无变化   | 仅更新硬件光标位置                             |

**洞见 3：Synchronized Output（`CSI ? 2026 h/l`）消除闪烁。**
所有写入被包裹在 synchronized output 对中，终端会缓冲直到收到结束标记后一次性刷新。
这是现代终端闪烁问题的根本解法。

**洞见 4：差分粒度是 "行" 而非 "字符"。**
行级差分 + synchronized output 在实践中已经足够，不需要字符级 diff。
因为终端写入的瓶颈在 I/O 而非计算，减少写入行数比精细 diff 更有效。

### 严格的宽度安全

渲染后会检查每行的 `visibleWidth` 是否超出终端宽度。如果超出：

1. 写入 crash log 到 `~/.pi/agent/pi-crash.log`
2. 调用 `stop()` 清理状态
3. 抛出 Error

**洞见 5：宽度溢出是 TUI 的"段错误"。**
当一行超出终端宽度，终端会自动换行，导致所有行号偏移，之后的差分对比全部错位。
这就是为什么代码中到处都有 `visibleWidth` 检查和 `sliceByColumn` 截断。

---

## 四、终端抽象与 Kitty 键盘协议

### `ProcessTerminal` 的启动流程

```
start()
  → setRawMode(true)           // 原始模式，逐字符接收
  → 启用 bracketed paste       // \x1b[?2004h
  → 发送 SIGWINCH              // 刷新终端尺寸
  → Windows: 启用 VT Input     // koffi 调 Win32 API
  → 查询 Kitty 协议            // CSI ? u
  → 设置 StdinBuffer           // 输入分割
```

### Kitty 键盘协议

传统终端键盘序列有歧义（如 `Ctrl+[` ≡ `ESC`），Kitty 协议解决了这个问题：

- 查询：发送 `CSI ? u`，终端回复 `CSI ? <flags> u`
- 启用：发送 `CSI > 7 u`（flag 1+2+4 = 按键消歧 + 事件类型 + 基础布局键）
- 好处：可区分 press/repeat/release、支持非拉丁键盘布局

**洞见 6：`keys.ts` 实现了完整的双层键盘解析。**

```ts
matchesKey(data, 'ctrl+shift+a') // 同时支持两种协议
```

它先尝试 Kitty CSI-u 解析（`\x1b[<cp>;<mod>u`），失败则回退到传统序列查找表。
这使得上层代码完全不需要关心终端类型。

### `StdinBuffer`：输入碎片重组

stdin 数据可能碎片化到达（特别是 SSH 场景）：

```
\x1b[<35;20;5m 可能拆成: "\x1b" → "[<35" → ";20;5m"
```

`StdinBuffer` 用状态机判断 `isCompleteSequence()`，缓冲不完整序列，超时后释放。

**洞见 7：`drainInput()` 解决 SSH 键盘泄漏。**
退出时，Kitty release 事件可能在 SSH 慢连接中延迟到达。
`drainInput()` 先禁用协议，然后等待 stdin 静默，防止转义序列泄漏到父 shell。

---

## 五、ANSI/Unicode 宽度计算（`utils.ts`）

这一层解决的核心问题：**终端中一个"字符"实际占多少列？**

```ts
visibleWidth('hello') // 5  (ASCII 快速路径)
visibleWidth('你好') // 4  (CJK 宽字符，每个 2 列)
visibleWidth('👨‍👩‍👧') // 2  (ZWJ emoji, 2 列)
visibleWidth('\x1b[31mhi') // 2  (ANSI 转义不占宽)
```

### 实现要点

1. **ASCII 快速路径**：纯 ASCII 直接 `str.length`
2. **Intl.Segmenter**：按 grapheme cluster 分割（处理组合字符）
3. **Emoji 预筛**：`couldBeEmoji()` 快速判断，避免昂贵的 `\p{RGI_Emoji}` regex
4. **东亚宽度**：`get-east-asian-width` 库判断全角/半角
5. **LRU 缓存**：512 条目的 Map 缓存非 ASCII 计算结果

### `AnsiCodeTracker`：样式追踪器

跨行渲染时需要知道"当前激活了哪些样式"。`AnsiCodeTracker` 逐条解析 SGR 码：

```ts
tracker.process('\x1b[1;31m') // bold=true, fgColor="31"
tracker.getActiveCodes() // → "\x1b[1;31m"
tracker.getLineEndReset() // → "" (只有 underline 会"泄漏"到 padding)
```

**洞见 8：`getLineEndReset()` 只重置 underline。**
因为在终端中，只有 underline 会视觉上延伸到行尾的空白区域。
颜色、粗体等不会泄漏，所以不需要在每行末尾完全重置——这减少了输出量。

---

## 六、Overlay 系统

Overlay 是在主内容之上合成的浮动组件（弹窗、菜单等）。

```ts
const handle = tui.showOverlay(component, {
  anchor: 'center', // 锚点：9 种位置
  width: '50%', // 支持百分比
  maxHeight: '50%',
  margin: 2, // 边距
  visible: (w, h) => w > 60 // 条件可见性
})
handle.setHidden(true) // 临时隐藏
handle.hide() // 永久移除
```

### 合成算法 (`compositeLineAt`)

对于每一行 overlay 覆盖的区域：

```
base:    [  before  |  overlayed  |  after  ]
overlay:            [  content    ]
result:  [  before  + content     + after   ]
```

用 `extractSegments()` **单趟**扫描 baseLine 提取 before 和 after 段，
然后拼接 overlay 内容，中间插入样式重置 `\x1b[0m\x1b]8;;\x07`。

**洞见 9：Overlay 合成是"单行分时复用"而非像素合成。**
不需要 z-buffer 或透明度——在文本终端中，overlay 就是在特定列范围内替换文本。
样式隔离通过 reset code 实现。

### Focus 栈

Overlay 有独立的 focus 栈：

- push overlay → 保存 `preFocus`，focus 到 overlay
- pop overlay → 恢复 `preFocus`
- overlay 可以条件性不可见 → focus 自动回退到最顶部可见 overlay

---

## 七、组件一览

### 原子组件

| 组件                | 功能                      | 关键设计                              |
| ------------------- | ------------------------- | ------------------------------------- |
| `Text`              | 包裹文字、padding、背景色 | `wrapTextWithAnsi` 跨行保持 ANSI 样式 |
| `TruncatedText`     | 单行截断（带省略号）      | 不换行，只取第一行                    |
| `Spacer`            | 空行                      | 最简单的组件                          |
| `Loader`            | 转转转动画                | 80ms 间隔 setInterval + requestRender |
| `CancellableLoader` | Loader + AbortSignal      | `handleInput` 监听 Escape             |
| `Image`             | Kitty/iTerm2 显示图片     | 多行占位 + 光标回退 + 图片序列        |

### 容器组件

| 组件        | 功能                      | 关键设计                                 |
| ----------- | ------------------------- | ---------------------------------------- |
| `Container` | 子组件纵向排列            | TUI 自身的基类                           |
| `Box`       | 子组件 + padding + 背景色 | 缓存渲染结果，通过 bgFn 采样检测主题变化 |

### 交互组件

| 组件           | 功能          | 关键设计                                  |
| -------------- | ------------- | ----------------------------------------- |
| `Input`        | 单行输入框    | 水平滚动、kill ring、undo                 |
| `Editor`       | 多行编辑器    | word wrap、自动补全、粘贴占位、历史       |
| `SelectList`   | 单选列表      | 箭头循环、滚动窗口                        |
| `SettingsList` | 设置面板      | 子菜单委托、搜索过滤                      |
| `Markdown`     | Markdown 渲染 | 用 `marked` lexer 解析，自定义 token 渲染 |

---

## 八、Editor 组件深度分析

Editor 是最复杂的组件（≈1000 行），实现了一个完整的终端文本编辑器：

### 渲染流程

```
state.lines[] (逻辑行)
  → wordWrapLine() → LayoutLine[] (视觉行，含光标映射)
  → scrollOffset 裁剪 → 可见行
  → 光标渲染（反色块 + CURSOR_MARKER）
  → 上下边框（带滚动提示 "↑ N more"）
  → 自动补全列表（如果激活）
```

**洞见 10：光标用 `\x1b[7m`（反色）渲染为"方块"。**
CURSOR_MARKER (`\x1b_pi:c\x07`) 是 APC 序列——终端会忽略它，但 TUI 能找到它的位置。
这实现了"软件光标"（视觉）和"硬件光标"（IME 输入法候选框定位）的分离。

### Emacs 风格操作

- **Kill Ring**：`Ctrl+K` / `Ctrl+U` 删除的文本进入 ring，`Ctrl+Y` 粘贴，`Alt+Y` 循环
- **连续 kill 累积**：连续的 kill 操作合并为一条（`accumulate: true`）
- **字符跳转**：`Ctrl+]` → 输入目标字符 → 跳到该字符

### 自动补全

```ts
editor.setAutocompleteProvider(provider)
// Provider 接口：
interface AutocompleteProvider {
  getCompletions(prefix: string): AutocompleteItem[]
}
```

内置了 `CombinedAutocompleteProvider`，支持合并多个补全源。
`autocomplete.ts` 还实现了：

- 文件路径补全（用 `fd` 工具搜索）
- `@"..."` 引号路径语法
- slash command 补全

### 大粘贴处理

大文本粘贴时，编辑器会：

1. 检测 bracketed paste（`\x1b[200~` ... `\x1b[201~`）
2. 大于 1KB 的粘贴替换为占位符 `[[paste:N]]`
3. 提交时 `getExpandedText()` 展开占位符

**洞见 11：编辑器是"提示输入框"而非通用编辑器。**
它的设计服务于 AI CLI 的 prompt 输入场景：
单行为主 → Shift+Enter 多行 → Enter 提交 → 历史浏览 → 自动补全 → 粘贴文件引用

---

## 九、测试策略

使用 `VirtualTerminal`（基于 `@xterm/headless`）做确定性测试：

```ts
const terminal = new VirtualTerminal(80, 24)
const tui = new TUI(terminal)
tui.start()

terminal.sendInput('ctrl+a') // 模拟输入
await terminal.flush() // 等待渲染
const viewport = terminal.getViewport() // 获取"屏幕"内容
```

**洞见 12：用真实终端模拟器做测试保证了正确性。**
xterm.js 的 headless 模式精确模拟终端行为（光标移动、滚动、ANSI 渲染），
而不是自己写一个简陋的模拟器。这使得测试结果与真实终端完全一致。

---

## 十、关键设计决策与取舍

### 1. 为什么不用 React/Ink 式虚拟 DOM？

Ink 用 Yoga (Flexbox) 做布局 + 虚拟 DOM diff + 字符矩阵输出。
pi-tui 选择了更简单的路径：

- 组件直接返回字符串数组（行）
- 差分以行为单位
- 布局是纯线性堆叠（除了 overlay）

**代价**：无横向并排布局（只有纵向流）
**收益**：极低复杂度，调试容易，性能更好

### 2. 为什么把 ANSI 处理内化而不用 strip-ansi 等库？

因为需要做的事情比 strip 更多：

- 需要在任意列位置切割含 ANSI 的字符串（`sliceByColumn`）
- 需要跟踪样式跨行延续（`AnsiCodeTracker`）
- 需要处理 OSC 8 超链接、APC 序列等
- 第三方库没有覆盖 overlay 合成这种场景

### 3. 为什么 Kitty 协议是可选的？

- 支持 Kitty 协议的终端还不是全部（虽然 Ghostty、WezTerm、Kitty 都支持）
- 传统终端下仍需可用 → 双层解析（Kitty 优先，fallback 到传统序列）
- 协议检测是自动的：查询 → 响应 → 启用，对用户透明

### 4. 图片协议的巧妙处理

图片在终端中是特殊的——它们不参与正常的文本流。
pi-tui 的处理方式：

- 用 N-1 个空行 + 1 个光标回退行来"占位"
- `isImageLine()` 检测使差分渲染跳过图片行
- 支持 Kitty（分块传输）和 iTerm2（单次内联）两种协议

---

## 十一、数据流总结

```
用户按键
  ↓
stdin → StdinBuffer (碎片重组) → TUI.handleInput
  ↓                                      ↓
  ├─ inputListeners (全局拦截)     cellSize 响应解析
  ├─ overlay focus 校验
  └─ focusedComponent.handleInput()
       ↓
     组件修改自身状态 → requestRender()
       ↓
     process.nextTick → doRender()
       ↓
     render() → compositeOverlays() → extractCursor() → diff → write()
       ↓
     终端显示更新
```

**洞见 13：`requestRender()` 用 `process.nextTick` 合并。**
一个事件循环 tick 内多次 `requestRender()` 只触发一次实际渲染。
这是 React `setState` 批量更新的相同思想。

---

## 十二、可复用的模式与技巧

1. **APC 序列做应用级标记**：`CURSOR_MARKER = "\x1b_pi:c\x07"` 是零宽度的，终端忽略它，但应用可以检测。
2. **Synchronized Output**：把所有输出包在 `CSI ? 2026 h/l` 中消除闪烁。
3. **样式继承用 tracker 而非嵌套**：`AnsiCodeTracker` 在行边界处重放样式，避免嵌套 ANSI 序列。
4. **crash log 而非静默失败**：宽度溢出写 log 后 throw，不隐藏 bug。
5. **Emacs kill ring** 是一种优雅的多剪贴板设计。
6. **structuredClone undo**：简单粗暴但有效，比 command pattern 省代码。
7. **fd 做文件补全**：外部工具做脏活（快、尊重 .gitignore），比 Node.js fs 遍历好 10 倍。
