# TUI 包全面解析

> `@mariozechner/pi-tui` —— 一个**差分渲染**的终端 UI 框架。
> 核心哲学：组件只管输出字符串数组，框架自动找出哪些行变了，只重绘那几行。

---

## 目录

1. [整体架构：一句话说清楚](#1-整体架构)
2. [核心接口 Component](#2-核心接口-component)
3. [TUI 类：调度中枢](#3-tui-类调度中枢)
4. [差分渲染算法 doRender](#4-差分渲染算法-dorender)
5. [终端抽象层 Terminal](#5-终端抽象层-terminal)
6. [键盘输入处理管线](#6-键盘输入处理管线)
7. [ANSI / Unicode 宽度工具集](#7-ansi--unicode-宽度工具集)
8. [Overlay 弹层系统](#8-overlay-弹层系统)
9. [内置组件一览](#9-内置组件一览)
10. [终端图片协议](#10-终端图片协议)
11. [关键设计洞见](#11-关键设计洞见)

---

## 1. 整体架构

```
用户操作 → stdin → StdinBuffer(缓冲/拼接转义序列)
         → TUI.handleInput(分发给 focused 组件)
         → 组件修改自身状态, 调用 invalidate/requestRender
         → TUI.doRender():
              1. 所有组件 render(width) → string[]
              2. 合成 overlay 覆盖层
              3. 与上一帧逐行比对
              4. 只把变化的行写入 stdout (包裹在 synchronized output 中)
```

**数据流是单向的**：输入 → 状态变更 → 重新渲染 → 差分输出。
组件不需要知道"自己在屏幕哪一行"，只需要返回 `string[]`。

---

## 2. 核心接口 Component

```ts
interface Component {
  render(width: number): string[] // 必须实现：给定宽度，返回各行
  handleInput?(data: string): void // 可选：处理键盘输入
  invalidate?(): void // 可选：清除缓存，强制重绘
  children?: Component[] // 可选：子组件列表
}
```

### 为什么这个设计好？

- **极简合约**：组件只需关心「给我宽度，我还你行」。不需要知道自己被放在屏幕哪里、被谁包含。
- **无状态渲染**：`render()` 是纯函数式的 —— 每次调用产生完整输出，框架负责差分。
- **单向数据流**：组件只管输出，不管重绘调度。

### Focusable 接口

```ts
interface Focusable {
  focused: boolean // 由 TUI 设置
}
```

只有同时实现 `Component` 和 `Focusable` 的组件才能接收键盘输入。
TUI 通过 `setFocus(component)` 将 `focused = true` 设给目标，将上一个焦点的 `focused = false`。

---

## 3. TUI 类：调度中枢

TUI 继承自 `Container`（一个简单的纵向堆叠容器），是整个 UI 的根节点。

### 核心状态

```ts
class TUI extends Container {
  terminal: Terminal // 终端抽象
  private focusedComponent // 当前焦点组件
  private overlayStack // 弹层栈
  private previousLines // 上一帧渲染结果（用于差分）
  private previousWidth // 上一帧宽度
  private cursorRow // 当前光标行（逻辑）
  private hardwareCursorRow // 硬件光标行（实际终端位置）
  private maxLinesRendered // 历史最大行数
}
```

### 启动流程

```
tui.start()
  → terminal.start()              // 进入 raw mode, 隐藏光标
  → terminal.enableKittyProtocol  // 如果终端支持，启用 Kitty 键盘协议
  → terminal.setBracketedPaste    // 启用粘贴模式
  → 绑定 stdin 的 data 事件      // 开始监听输入
```

### 请求渲染的巧妙设计

```ts
requestRender(): void {
  if (this.renderScheduled) return;
  this.renderScheduled = true;
  process.nextTick(() => {
    this.renderScheduled = false;
    this.doRender();
  });
}
```

**关键洞见**：用 `process.nextTick` 把同一事件循环中的多次渲染请求合并为一次。
比如一个事件同时修改了 3 个组件的状态，它们各自 `requestRender()`，但实际只渲染一次。

### 输入分发

```
handleInput(data)
  → 过滤掉 Kitty key release 事件 (":3" 后缀)
  → 过滤掉 Kitty 协议响应
  → 过滤掉 cell size 响应
  → 如果有 overlay 且 overlay 有焦点组件 → 分发给 overlay
  → 否则分发给 focusedComponent.handleInput(data)
  → 每次输入后 requestRender()
```

---

## 4. 差分渲染算法 doRender

这是整个 TUI 框架的核心，也是最精妙的部分。

### 渲染流水线

```
doRender()
  ┌── 1. 收集所有行 ──────────────────────────┐
  │  rootLines = this.render(width)            │  所有组件纵向输出
  │  overlayLines = compositeOverlays(...)     │  叠加弹层
  │  newLines = overlayLines                   │
  └────────────────────────────────────────────┘
  ┌── 2. 提取光标位置 ────────────────────────┐
  │  扫描 newLines 中的 CURSOR_MARKER         │  一个特殊 APC 序列
  │  记录 {row, col}, 然后从行中删除 marker  │
  └────────────────────────────────────────────┘
  ┌── 3. 行尾 reset 处理 ────────────────────┐
  │  对每行末尾添加 ANSI reset              │  只 reset 下划线
  │  (只有下划线会"渗透"到右侧 padding)     │  其他样式不需要
  └────────────────────────────────────────────┘
  ┌── 4. 差分比对 ────────────────────────────┐
  │  逐行对比 newLines vs previousLines      │
  │  找到 firstChanged 和 lastChanged        │
  └────────────────────────────────────────────┘
  ┌── 5. 输出 ────────────────────────────────┐
  │  只输出 [firstChanged..lastChanged] 的行 │
  │  包裹在 synchronized output 中           │
  │  最后定位硬件光标(用于 IME 候选窗)       │
  └────────────────────────────────────────────┘
```

### 完全重绘 vs 差分更新

两种情况会触发完全重绘 (`fullRender`):

1. **宽度变化** — 终端 resize 时所有行都需要重排
2. **首个变化行在可视区域之上** — 无法局部更新，需要完整重绘
3. **新增行数超过终端高度** — 直接全部重写更高效

差分更新时：

- 用 CSI 序列移动光标到 `firstChanged` 所在行
- 只重写 `firstChanged` 到 `lastChanged` 之间的行
- 如果旧行比新行多，用 `\x1b[2K` 清除多余行

### Synchronized Output

```
\x1b[?2026h   ← 开始同步输出 (终端缓冲所有写入)
  ...写入所有变化行...
\x1b[?2026l   ← 结束同步输出 (终端一次性刷新)
```

这个 DEC 私有模式序列让支持的终端把所有输出攒到一起，最后一次性刷新到屏幕，**彻底消除闪烁**。

### CURSOR_MARKER 的精巧设计

```ts
const CURSOR_MARKER = '\x1b_pi:c\x07' // APC (Application Program Command) 序列
```

组件在 `render()` 中把这个零宽标记插到光标应该在的位置。
框架在渲染后扫描所有输出行，提取标记位置，然后把硬件光标移过去。

**为什么用 APC？** 因为 APC 是零宽度的，不影响排版计算，终端也不会显示它。
**用途**：硬件光标定位给 IME（输入法）候选窗提供正确的位置。

### 视口管理

当内容超过终端高度时，TUI 只渲染底部可见的行（自然滚动到底部）：

```ts
const viewportTop = Math.max(0, newLines.length - height)
```

光标追踪使用两个变量：

- `cursorRow`：逻辑光标行（内容末尾）
- `hardwareCursorRow`：真实终端光标行（用于计算移动距离）

---

## 5. 终端抽象层 Terminal

### Terminal 接口

```ts
interface Terminal {
  start(): void // 进入 raw mode
  stop(): void // 恢复正常模式
  write(data: string) // 写入 stdout
  cols: number // 列数
  rows: number // 行数
  showCursor(): void
  hideCursor(): void
  clearScreen(): void
  // ...
}
```

### ProcessTerminal 启动序列

```
start()
  → process.stdin.setRawMode(true)       // 字符级输入
  → 隐藏光标
  → 启用 bracketed paste mode            // 区分用户输入和粘贴
  → 监听 SIGWINCH                        // 窗口大小变化
  → Windows: setupWindowsVtInput()       // Windows 需要额外的 VT 输入处理
  → 查询 Kitty 键盘协议支持              // CSI ? u
  → 如果支持，启用 flags=7               // CSI > 7 u (disambiguate + report events + report alternates)
  → 初始化 StdinBuffer                   // 缓冲/拼接转义序列
```

### drainInput() — 优雅退出

退出时调用，等待 50ms 读取并丢弃所有剩余输入。
**目的**：Kitty 协议的 key release 事件可能延迟到达。如果不 drain，这些字节会泄漏到父 shell，
导致用户看到乱码。在 SSH 连接中尤其重要，因为网络延迟会加剧这个问题。

---

## 6. 键盘输入处理管线

### 分层架构

```
stdin 原始字节
  → StdinBuffer (拼接/缓冲转义序列)
     → 检测序列完整性 (CSI/OSC/DCS/APC)
     → 超时释放不完整序列 (10ms)
  → TUI.handleInput
     → matchesKey() / parseKey()
        → 优先匹配 Kitty CSI-u 格式
        → 回退到传统 VT 序列表
```

### StdinBuffer 的必要性

终端的一次 `data` 事件可能包含多个完整序列，也可能只包含一个序列的一部分。
StdinBuffer 解决这个问题：

```ts
function isCompleteSequence(data: Buffer, startPos: number): number | null
```

这是一个**状态机**，识别 CSI（`ESC [`）、OSC（`ESC ]`）、DCS（`ESC P`）、APC（`ESC _`）序列的完整边界。
当缓冲区中凑齐完整序列时立即释放，否则等待 10ms 后释放（处理单独的 ESC 键）。

### Kitty 键盘协议

传统终端协议的问题：

- 无法区分 `Ctrl+I` 和 `Tab`
- 无法区分 `Ctrl+M` 和 `Enter`
- 大写字母 + Ctrl 无法表达

Kitty 协议格式：`CSI keycode ; modifiers u`  
其中 modifiers 编码了 Shift, Alt, Ctrl, Super 的组合。

框架的策略：**先尝试 Kitty 解析，失败则回退到传统序列匹配**。
这保证了在不支持 Kitty 的终端上也能正常工作。

### 按键释放过滤

Kitty 协议报告按键按下 **和** 释放事件。释放事件的 modifier 包含 `:3`。
TUI 在 `handleInput` 的最前面就过滤掉所有 release 事件，只处理按下。

---

## 7. ANSI / Unicode 宽度工具集

### visibleWidth() — 终端可见宽度计算

这是整个 TUI 中**调用最频繁**的函数，也是精心优化的：

```ts
function visibleWidth(str: string): number
```

计算策略：

1. **ASCII 快速路径**：如果字符串全是 ASCII 且无 ANSI，直接返回 `.length`
2. **跳过 ANSI 转义序列**：CSI、OSC、APC 等序列宽度为 0
3. **grapheme 分割**：用 `Intl.Segmenter` 提取 grapheme cluster
4. **宽字符检测**：CJK 字符宽度为 2（通过 `get-east-asian-width`）
5. **Emoji 检测**：组合 emoji（如 👨‍👩‍👧‍👦）按 grapheme 正确计算
6. **LRU 缓存**：最近 512 个结果缓存

### AnsiCodeTracker — 跨行 ANSI 状态追踪

当文本需要换行时，ANSI 样式码必须在新行重新激活。
`AnsiCodeTracker` 维护当前激活的 SGR 代码集合：

```ts
class AnsiCodeTracker {
  process(ansiCode: string): void // 记录遇到的 SGR 代码
  getActiveCodes(): string // 返回当前所有激活样式的 ESC 序列
  getLineEndReset(): string // 行尾需要的 reset（只 reset 下划线）
  reset(): void // 清空状态
}
```

**关键洞见**：`getLineEndReset()` **只重置下划线**。
为什么？因为在终端中，只有下划线会「渗透」到行尾 padding 空格。
粗体、斜体、颜色等不会影响到字符外的空间。这个细节减少了大量不必要的 ANSI 输出。

### wrapTextWithAnsi() — ANSI 安全的文本换行

```ts
function wrapTextWithAnsi(text: string, maxWidth: number): string[]
```

该函数在做文本换行的同时：

- 正确处理多字节字符（不会在 grapheme 中间切断）
- 在换行处关闭所有 ANSI 样式
- 在新行开头重新开启上一行的 ANSI 样式
- 正确计算 CJK/emoji 的列宽

### extractSegments() — Overlay 合成辅助

```ts
function extractSegments(
  line: string,
  insertCol: number,
  insertWidth: number
): { before: string; after: string }
```

从一行中提取「overlay 之前」和「overlay 之后」的部分，处理了：

- ANSI 序列在切割点的正确闭合/重开
- 宽字符被切成半边时用空格填充
- 保持切割后的列宽精确

---

## 8. Overlay 弹层系统

Overlay 是 TUI 的弹出层机制，用于实现对话框、下拉菜单等浮在主内容之上的 UI。

### 数据结构

```ts
interface OverlayOptions {
  component: Component
  width: number | string // 绝对列数或百分比 "80%"
  anchor: 'top' | 'bottom' // 锚定位置
  margin?: number // 距锚点的边距
  visible?: () => boolean // 动态显示/隐藏
}

interface OverlayHandle {
  hide(): void // 关闭弹层
}
```

### 工作原理

1. 调用 `tui.showOverlay(options)` → 返回 `OverlayHandle`
2. Overlay 被推入 `overlayStack`（可以多层叠加）
3. 当前焦点被保存，焦点转移到 overlay 的组件
4. 渲染时，先渲染底层内容，然后按栈顺序合成 overlay

### 合成算法 compositeOverlays

```
对每个 overlay（从栈底到栈顶）:
  1. 计算 overlay 的行范围 (startRow, endRow)
  2. 计算列范围 (居中放置)
  3. 对 overlay 覆盖区域内的每一行:
     - 取底层行的「左侧」部分
     - 取 overlay 行
     - 取底层行的「右侧」部分
     - 拼接起来
```

**居中公式**：`startCol = Math.floor((width - overlayWidth) / 2)`

### 焦点保存/恢复

```
showOverlay → 保存当前 focusedComponent → 设置 overlay 组件为焦点
hideOverlay → 恢复之前保存的 focusedComponent
```

---

## 9. 内置组件一览

### Container — 纵向堆叠容器

```ts
class Container implements Component {
  children: Component[] = []
  render(width: number): string[] {
    // 将所有子组件的输出垂直拼接
    return children.flatMap(c => c.render(width))
  }
}
```

TUI 本身就继承自 Container，它就是根容器。

### Text — 文本显示

带 padding、background、word wrap 的文本组件。
使用 `wrapTextWithAnsi()` 做 ANSI 安全的文本换行。
有渲染缓存：text + width 不变就直接返回上次结果。

### TruncatedText — 单行截断文本

只显示一行，超出宽度时截断（用 `truncateToWidth()`）。
适用于标题栏、状态栏等固定单行场景。

### Box — 带 padding/background 的容器

```ts
class Box implements Component {
  children: Component[];
  paddingX, paddingY;
  bgFn?: (text: string) => string;  // 背景色函数
}
```

渲染所有子组件后，添加 padding，并用 `applyBackgroundToLine()` 施加背景色。
背景色应用在 padding 阶段而非内联，这样可以让背景色延伸到行的完整宽度。

### Editor — 多行文本编辑器

最复杂的组件（~1000+ 行），功能包括：

| 功能       | 实现方式                                                      |
| ---------- | ------------------------------------------------------------- |
| 自动换行   | `wordWrapLine()` 切成 `TextChunk[]`，保持光标位置映射         |
| 垂直滚动   | 可见行数 = 终端高度 \* 30%，最少 5 行                         |
| 自动补全   | AutocompleteProvider 接口 + SelectList 浮层                   |
| 命令历史   | 上下箭头浏览，navigateHistory()                               |
| Kill Ring  | Emacs 风格 Ctrl+K/U 删除，Ctrl+Y 粘贴，Meta+Y 轮转            |
| 撤销       | UndoStack<EditorState>，structuredClone 快照                  |
| 粘贴处理   | 支持 bracketed paste，超大粘贴自动替换为 `[[paste:N]]` 占位符 |
| 字符跳转   | 类 vim `f`/`F` 的跳转到指定字符                               |
| Kitty 解码 | `decodeKittyPrintable()` 从 CSI-u 序列提取可打印字符          |

光标渲染：用 `\x1b[7m` (反色) 高亮当前字符，末尾用反色空格。

### Input — 单行输入

和 Editor 类似但只有一行，支持水平滚动。
同样有 Kill Ring、撤销、grapheme 级光标移动。

### SelectList — 选择列表

可滚动的选项列表，箭头上下选择，回车确认。
支持 `setFilter()` 过滤。

### SettingsList — 设置列表

比 SelectList 更复杂：

- 左侧标签 + 右侧当前值
- 支持切换值 / 打开子菜单
- 可选搜索过滤（使用 fuzzy matching）

### Loader — 加载动画

继承 Text，用 `setInterval(80ms)` 旋转 braille 字符 `⠋⠙⠹...`。

### Image — 终端图片

利用 Kitty Graphics Protocol 或 iTerm2 内联图片协议在终端显示图片。
不支持的终端显示 fallback 文字。

### Markdown — Markdown 渲染

用 `marked` 的 lexer 解析 Markdown token，然后用 theme 函数渲染：

- heading / paragraph / code / blockquote / list / hr 等
- 背景色在 padding 阶段施加
- 有渲染缓存 (text + width)

### Spacer — 空白行

最简单的组件。`render()` 返回 N 个空字符串。

---

## 10. 终端图片协议

### 支持的协议

| 协议           | 检测方式                   | 编码方式                      |
| -------------- | -------------------------- | ----------------------------- |
| Kitty Graphics | `KITTY_WINDOW_ID` 环境变量 | Base64 分块传输，4096 字节/块 |
| iTerm2 Inline  | `TERM_PROGRAM=iTerm.app`   | OSC 1337 编码                 |

### Kitty 图片传输

```
ESC_Gf=100,a=T,m=1,i=ID,q=2;chunk1 ESC\    ← 第一块 (m=1 还有后续)
ESC_Gm=1;chunk2 ESC\                         ← 中间块
ESC_Gm=0;chunkN ESC\                         ← 最后一块 (m=0 传输完成)
```

### 行高计算

图片转换为终端行数：

```ts
imageRows = Math.ceil(heightPx / cellHeightPx)
```

`cellHeightPx` 通过终端的 cell size 响应获取。

### 渲染技巧

图片占据 N 行。前 N-1 行是空行（让 TUI 知道要为图片留足空间），
最后一行包含 `\x1b[${N-1}A`（移回顶部）+ 图片 escape 序列。

---

## 11. 关键设计洞见

### 1. 「字符串数组」作为渲染原语

不用 virtual DOM，不用 canvas，就用 `string[]`。  
每个字符串就是终端的一行，天然对齐终端的工作方式。
差分就是逐行字符串比较，简单高效。

### 2. 渲染合并

`requestRender()` 用 `process.nextTick` 自动批处理。
无论多少组件在同一事件循环中调用 `requestRender()`，实际只渲染一次。

### 3. 只有下划线会渗透

`getLineEndReset()` 只重置下划线。这是一个通过实际测试发现的终端行为：
只有下划线样式会视觉上延伸到行尾空白区域。省略其他 reset 减少了输出量。

### 4. APC 序列做光标标记

`CURSOR_MARKER` 使用 APC（Application Program Command）格式 `\x1b_pi:c\x07`。
APC 在终端中零宽度、不可见、不影响排版。组件只需在渲染输出中插入这个标记，
框架就能知道硬件光标该放在哪里——这样 IME 候选窗能出现在正确位置。

### 5. Kitty 协议的优雅降级

启动时发送 `CSI ? u` 查询。如果终端响应了，就启用。
后续按键解析优先尝试 Kitty 格式，失败则回退传统序列表。
这种「先问再用，失败回退」的策略保证了最广泛的终端兼容性。

### 6. drain 输入防泄漏

退出时 drain 50ms 的 stdin —— 因为 Kitty 协议的 key release 事件可能延迟到达。
如果不 drain，这些字节会被父 shell 解释为命令输入，造成"幽灵按键"。
在 SSH 场景下这个问题尤其严重（网络延迟放大了 release 的延迟）。

### 7. 图片行的特殊处理

`isImageLine()` 快速检测某行是否包含图片 escape 序列。
差分渲染时，图片行需要特殊处理 —— 不能像普通行一样做字符串比较和清除/重写，
因为图片的 escape 序列可能非常长，且包含 base64 数据。

### 8. structuredClone 做撤销

Editor 的 UndoStack 用 `structuredClone()` 做状态快照。
虽然不是最节省内存的方案，但代码极简且正确——对于编辑器场景（状态较小）完全够用。

### 9. Kill Ring 的连续积累

连续的 kill 操作会合并到同一个 ring entry：

- 向后删除 → prepend
- 向前删除 → append

停止 kill 后再 kill 则创建新 entry。这与 Emacs 行为完全一致。

### 10. 缓存与失效

各组件普遍采用「缓存上次渲染 + 输入变化时失效」：

- Box: 缓存 childLines + width + bgSample
- Text: 缓存 text + width
- Markdown: 缓存 text + width
- Image: 缓存 base64Data + width

TUI 调 `invalidate()` 会级联到所有子组件。

---

## 附录：关键文件索引

| 文件                           | 职责             | 核心导出                                              |
| ------------------------------ | ---------------- | ----------------------------------------------------- |
| `tui.ts`                       | 框架核心         | `TUI`, `Component`, `Container`, `CURSOR_MARKER`      |
| `terminal.ts`                  | 终端抽象         | `Terminal`, `ProcessTerminal`                         |
| `utils.ts`                     | 字符串/ANSI 工具 | `visibleWidth`, `wrapTextWithAnsi`, `AnsiCodeTracker` |
| `keys.ts`                      | 按键解析         | `Key`, `matchesKey`, `parseKey`, `KeyId`              |
| `stdin-buffer.ts`              | 输入缓冲         | `StdinBuffer`, `isCompleteSequence`                   |
| `keybindings.ts`               | 快捷键绑定       | `EditorKeybindingsManager`, `EditorAction`            |
| `terminal-image.ts`            | 终端图片         | `renderImage`, `detectCapabilities`                   |
| `kill-ring.ts`                 | Kill Ring        | `KillRing`                                            |
| `undo-stack.ts`                | 撤销栈           | `UndoStack`                                           |
| `fuzzy.ts`                     | 模糊匹配         | `fuzzyMatch`, `fuzzyFilter`                           |
| `autocomplete.ts`              | 补全逻辑         | 路径补全、fd 集成                                     |
| `components/editor.ts`         | 多行编辑器       | `Editor`                                              |
| `components/input.ts`          | 单行输入         | `Input`                                               |
| `components/select-list.ts`    | 选择列表         | `SelectList`                                          |
| `components/settings-list.ts`  | 设置列表         | `SettingsList`                                        |
| `components/box.ts`            | 容器盒子         | `Box`                                                 |
| `components/text.ts`           | 文本显示         | `Text`                                                |
| `components/markdown.ts`       | Markdown         | `Markdown`                                            |
| `components/loader.ts`         | 加载动画         | `Loader`                                              |
| `components/image.ts`          | 图片显示         | `Image`                                               |
| `components/truncated-text.ts` | 截断文本         | `TruncatedText`                                       |
| `components/spacer.ts`         | 空白             | `Spacer`                                              |
