# @mariozechner/pi-tui 全面讲解

**架构层面**：整体分层图（TUI → 组件 → Terminal 抽象 → 底层工具），组件接口设计

**核心模块**（每个都结合了源码原理 + 测试用例说明）：

- **TUI 渲染引擎**：三策略差分渲染、Overlay 叠加系统、IME 光标支持
- **Terminal 抽象**：ProcessTerminal（Kitty 协议/bracketed paste）、VirtualTerminal（测试）
- **keys.ts**：键盘解析（支持 Kitty + 传统序列 + 非拉丁布局）
- **stdin-buffer.ts**：输入缓冲与序列拆分
- **utils.ts**：可见宽度计算（CJK/emoji）、ANSI 感知截断/换行
- **fuzzy.ts**：模糊匹配算法

**内置组件**（7 个）：Editor、Input、Markdown、SelectList、TruncatedText、Image、Autocomplete，每个都列出了对应测试覆盖的场景

## 一、项目定位

`pi-tui` 是一个**轻量级终端 UI 框架**，核心卖点是**差分渲染 + 同步输出**，实现终端应用的**无闪烁**交互体验。它被用来构建 `pi`（一个 AI 终端聊天工具）的整个终端界面。

> 类比：React 之于浏览器 DOM，pi-tui 之于终端屏幕。它用组件树描述 UI，用差分算法只更新变化的行。

---

## 二、整体架构

```
┌─────────────────────────────────────────────┐
│                    TUI                       │  ← 主容器，管理组件树 + 渲染循环
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │ Text     │  │ Editor   │  │ Markdown  │  │  ← 内置组件
│  │ Input    │  │ Loader   │  │ SelectList│  │
│  │ Box      │  │ Image    │  │ Spacer    │  │
│  └──────────┘  └──────────┘  └───────────┘  │
│  ┌─────────────────────────────────────────┐ │
│  │           Overlay Stack                 │ │  ← 弹窗/菜单叠加层
│  └─────────────────────────────────────────┘ │
├─────────────────────────────────────────────┤
│              Terminal 接口                    │  ← 抽象层
│  ProcessTerminal (真实终端)                   │
│  VirtualTerminal (测试用虚拟终端)              │
├─────────────────────────────────────────────┤
│              底层工具                         │
│  keys.ts  utils.ts  fuzzy.ts  stdin-buffer  │
│  kill-ring  undo-stack  terminal-image      │
└─────────────────────────────────────────────┘
```

---

## 三、核心模块详解

### 3.1 Component 接口 —— 一切组件的基础

```ts
interface Component {
  render(width: number): string[] // 返回行数组，每行不能超过 width
  handleInput?(data: string): void // 有焦点时接收键盘输入
  invalidate(): void // 清除缓存，强制下次重新渲染
}
```

**关键约束**：`render()` 返回的每一行的**可见宽度**（不含 ANSI 转义码）不能超过 `width`，否则 TUI 会直接 crash 并写入错误日志。

**测试验证**（`tui-render.test.ts`）：

- 测试差分渲染只更新变化的行
- 测试 resize 后触发全量重绘
- 测试内容缩短时正确清除多余行
- 测试每行末尾自动追加 SGR reset（`\x1b[0m\x1b]8;;\x07`），防止样式泄漏到下一行

---

### 3.2 TUI 类 —— 渲染引擎

TUI 继承自 `Container`（组件容器），是整个框架的核心。

#### 差分渲染三策略

| 场景                              | 策略                           | 触发条件                           |
| --------------------------------- | ------------------------------ | ---------------------------------- |
| **首次渲染**                      | 直接输出所有行                 | `previousLines` 为空               |
| **终端宽度变化 / 变化在视口上方** | 清屏 + 全量重绘                | `width` 改变或首个变化行在视口上方 |
| **正常更新**                      | 移到首个变化行，只渲染变化范围 | 默认情况                           |

所有输出都用**同步输出协议**（CSI 2026）包裹，确保终端原子性刷新，不闪烁。

**测试验证**（`tui-render.test.ts`）：

```
"TUI resize handling" → 宽度变化触发全量重绘
"TUI differential rendering" - 验证只有变化行被重写:
  - "only changed lines are rendered" → spinner 类场景，只中间行变化
  - "first line changed" → 仅首行变化
  - "last line changed" → 仅末行变化
  - "non-adjacent changes" → 不相邻行变化
```

#### Overlay（叠加层）系统

Overlay 是渲染在基础内容之上的弹窗层，支持：

- **锚点定位**：`center`、`top-left`、`bottom-right` 等 9 种
- **百分比定位**：`row: "25%"`, `col: "50%"`
- **绝对定位**：`row: 5, col: 10`
- **尺寸控制**：`width`、`minWidth`、`maxHeight`（支持百分比）
- **边距**：`margin`
- **响应式可见性**：`visible: (w, h) => w >= 100`

叠加原理：在基础行的指定列位置"切开"，插入 overlay 内容，两侧保留原始内容。

**测试验证**（`overlay-options.test.ts`）：

```
"width overflow protection" → overlay 内容超宽时自动截断，防止 crash
  - 含复杂 ANSI 序列、超链接、CJK 宽字符的边界处理
"anchor positioning" → top-left/bottom-right/top-center 各锚点位置正确
"percentage positioning" → row/col 百分比计算（0%=顶/左，100%=底/右）
"maxHeight" → 超高内容被截断
"stacked overlays" → 多层叠加，后者覆盖前者，hideOverlay 移除最上层
```

**测试验证**（`overlay-short-content.test.ts`）：

```
当基础内容不足以填满终端时，overlay 仍能正确渲染
```

**测试验证**（`tui-overlay-style-leak.test.ts`）：

```
基础内容的 ANSI 样式代码不会因为 overlay 切割而泄漏到下一行
（SGR reset 被 overlay 截断的回归测试）
```

#### Focusable 接口 + IME 支持

实现了 `Focusable` 接口的组件在获得焦点时：

1. TUI 设置 `focused = true`
2. 组件在渲染时在光标位置插入 `CURSOR_MARKER`（一个零宽 APC 转义序列 `\x1b_pi:c\x07`）
3. TUI 扫描渲染结果，找到 marker 并定位硬件光标

这使得 CJK 输入法（IME）候选窗能出现在正确位置。

---

### 3.3 Terminal 接口 —— 终端抽象

```ts
interface Terminal {
  start(onInput, onResize): void
  stop(): void
  write(data: string): void
  columns: number
  rows: number
  // ... 光标操作、清屏操作
}
```

两个实现：

- **`ProcessTerminal`**：真实终端（`process.stdin/stdout`），支持 raw mode、bracketed paste、Kitty 键盘协议
- **`VirtualTerminal`**：测试用，基于 `@xterm/headless`，模拟终端行为

#### Kitty 键盘协议

`ProcessTerminal` 启动时会查询终端是否支持 Kitty 协议：

1. 发送 `\x1b[?u` 查询
2. 如果收到 `\x1b[?<flags>u` 响应 → 支持
3. 发送 `\x1b[>7u` 启用（flag 1+2+4 = 消歧义 + 事件类型 + 备选键）

Kitty 协议的好处：能区分按下/重复/释放、支持非拉丁键盘布局。

---

### 3.4 keys.ts —— 键盘输入解析

提供 `matchesKey(data, keyId)` 和 `parseKey(data)` 两个核心函数，同时支持传统终端序列和 Kitty 协议。

```ts
// 类型安全的键标识符
Key.ctrl('c') // → "ctrl+c"
Key.ctrlShift('p') // → "ctrl+shift+p"
Key.escape // → "escape"
Key.alt('left') // → "alt+left"
```

**测试验证**（`keys.test.ts`）：

```
"Kitty protocol with alternate keys" → 非拉丁键盘支持:
  - 西里尔布局的 Ctrl+C（物理 C 键）被正确识别
  - Dvorak 布局的字母映射
  - 修饰键组合（Ctrl+Shift+字母）

"Legacy key matching" → 传统终端序列:
  - Ctrl+C = ASCII 3
  - 方向键 = \x1b[A/B/C/D
  - F1-F12 功能键
  - Alt+方向键、Shift+Tab
```

---

### 3.5 stdin-buffer.ts —— 输入缓冲与序列拆分

**问题**：终端的 stdin 数据可能被拆成多个 chunk 到达，例如鼠标事件 `\x1b[<35;20;5m` 可能分 3 次到达。

`StdinBuffer` 的职责：

1. 缓冲不完整的转义序列，等待后续数据
2. 将批量到达的数据拆分成独立事件
3. 识别 bracketed paste 标记，将粘贴内容作为整体事件发出

**测试验证**（`stdin-buffer.test.ts`）：

```
"Regular Characters" → 普通字符直接透传
"Complete Escape Sequences" → 完整的鼠标/方向键/功能键序列直接发出
"Partial Escape Sequences" → 不完整序列被缓冲，超时后发出
"Split sequences" → 跨 chunk 的序列被正确拼接
"Kitty Keyboard Protocol" → 11 个测试覆盖 press/release/修饰键/批量序列
"Mouse Events" → 5 个测试覆盖 SGR 鼠标按下/释放/移动/分段事件
"Bracketed Paste" → 5 个测试覆盖粘贴开始/结束标记、分段粘贴
```

---

### 3.6 utils.ts —— 字符串宽度与 ANSI 处理

核心工具函数：

| 函数                           | 功能                                             |
| ------------------------------ | ------------------------------------------------ |
| `visibleWidth(str)`            | 计算终端可见宽度（忽略 ANSI 码，处理 CJK/emoji） |
| `truncateToWidth(str, width)`  | 截断到指定宽度（保留 ANSI 码，加省略号）         |
| `wrapTextWithAnsi(str, width)` | 按宽度换行（保留 ANSI 跨行样式）                 |

**宽度计算的难点**：

- emoji 宽 2 列（如 👍）
- CJK 字符宽 2 列
- ANSI 转义码宽 0 列
- 组合字符（ZWJ emoji 序列）需要 `Intl.Segmenter` 分割

**测试验证**（`wrap-ansi.test.ts`）：

```
ANSI styled text 换行时样式在新行继续
嵌套样式（bold + color）正确传递
连字符位置换行
```

**测试验证**（`regression-regional-indicator-width.test.ts`）：

```
旗帜 emoji（Regional Indicator）的流式渲染问题：
  - 部分旗帜字符 "🇨" 在流式输出中是中间状态
  - 如果测量为宽度 1 而终端渲染为宽度 2，差分渲染会漂移
  - 解决方案：Regional Indicator 单例始终按宽度 2 计算
```

---

### 3.7 fuzzy.ts —— 模糊匹配

用于自动补全的模糊搜索算法：

- 查询字符必须按序出现在目标中（不要求连续）
- **评分规则**：连续匹配加分、词边界匹配加分、间隔越大扣分越多
- 支持字母数字 token 交换（如 `codex52` 能匹配 `gpt-5.2-codex`）

**测试验证**（`fuzzy.test.ts`）：

```
空查询匹配所有、查询过长不匹配
连续匹配得分优于零散匹配
词边界匹配得分更优
数字字母交换匹配（codex52 → 5.2-codex）
fuzzyFilter 支持自定义 getText 函数
```

---

### 3.8 kill-ring.ts + undo-stack.ts —— 编辑辅助

**KillRing**：Emacs 风格的剪切环

- `push(text, { prepend, accumulate })` → 添加或合并到环中
- `peek()` → 查看最近一条
- `rotate()` → 循环切换（yank-pop）
- 连续的 Ctrl+W 会 accumulate 成一条

**UndoStack**：通用撤销栈

- `push(state)` → 深拷贝（`structuredClone`）状态快照
- `pop()` → 弹出最近快照
- 支持操作合并（连续输入合并为一次 undo 单元）

---

## 四、内置组件详解

### 4.1 Editor —— 多行编辑器

最复杂的组件，功能包括：

- 多行编辑 + 自动换行
- 命令历史（Up/Down 导航，最多 100 条）
- Kill ring（Ctrl+W/U/K 删除 → Ctrl+Y 粘贴 → Alt+Y 循环）
- Undo/Redo（连续输入合并、删除/粘贴各自为独立单元）
- Slash 命令自动补全（输入 `/`）
- 文件路径补全（按 Tab）
- 大段粘贴处理（>10 行变为折叠标记）
- 字符跳转（Ctrl+] 向前跳到指定字符）
- Sticky column（垂直移动时保持目标列）

**测试验证**（`editor.test.ts`）—— 最大的测试文件：

```
"Prompt history navigation" (12 个测试):
  - Up/Down 循环历史
  - 空/重复条目不入历史
  - 历史限制 100 条

"Kill ring" (15 个测试):
  - Ctrl+W 连续删除 accumulate 到同一条
  - Ctrl+U (删到行首) + Ctrl+K (删到行尾)
  - Ctrl+Y yank + Alt+Y yank-pop 循环
  - 不同删除操作之间不 accumulate

"Undo" (25+ 个测试):
  - 连续输入合并为一个 undo 单元
  - 删除操作（Ctrl+W, Ctrl+U, 粘贴）各自独立
  - 光标移动开启新 undo 单元
  - Ctrl+Z 撤销 + Ctrl+Shift+Z / Ctrl+Y 重做

"Grapheme-aware text wrapping" (8 个测试):
  - emoji (👍宽2) + 普通字符混合换行
  - 双宽字符在行尾时整体移到下一行
  - ZWJ 序列（如 👩‍👩‍👧‍👦）正确处理

"Word wrapping" (11 个测试):
  - 在空格/标点处断行
  - 超长单词强制断行
  - URL 处理

"Sticky column" (18 个测试):
  - 从长行移到短行再移回，光标恢复原列
  - 通过换行的宽字符行

"Character jump (Ctrl+])" (数个测试):
  - 跳到指定字符、反向跳转
```

### 4.2 Input —— 单行输入

Editor 的单行版本，支持水平滚动、kill ring、undo。

**测试验证**（`input.test.ts`）：

```
"Kill ring" (9 个测试) → 与 Editor 类似
"Undo" (10 个测试) → 与 Editor 类似
基础输入、退格、提交
```

### 4.3 Markdown —— Markdown 渲染器

将 Markdown 文本渲染为带 ANSI 样式的终端行。

**测试验证**（`markdown.test.ts`）：

```
"Nested lists" (5 个测试):
  - 简单/深层嵌套、有序/无序/混合列表
  - LLM 输出中代码块不缩进时的编号保持

"Tables" (11 个测试):
  - 列宽计算、对齐、行分隔符
  - 窄终端下的表格换行
  - 复杂单元格内容

"Blockquotes" (5 个测试):
  - 多行内容、lazy continuation
  - 内联格式（加粗、代码）

"Spacing" (4 个测试):
  - 代码块/分隔线/标题/引用后的空行间距

"Pre-styled text":
  - 灰色斜体（thinking trace）+ 内联代码后样式恢复

"HTML tags" → 渲染为纯文本，不隐藏内容
```

### 4.4 SelectList —— 选择列表

带键盘导航的选项列表。

**测试验证**（`select-list.test.ts`）：

```
多行 description 规范化为单行
```

### 4.5 TruncatedText —— 单行截断文本

用于状态栏、标题等场景。

**测试验证**（`truncated-text.test.ts`）：

```
精确宽度填充（不足补空格）
垂直 padding
长文本截断 + 省略号
ANSI 样式保留 + reset 在省略号前
空文本、换行符（只取第一行）
```

### 4.6 Image —— 终端内联图片

支持 Kitty graphics protocol 和 iTerm2 inline images，不支持的终端显示占位文本。

**测试验证**（`terminal-image.test.ts`）：

```
"iTerm2 image protocol" → 检测 ESC ]1337;File= 序列
"Kitty image protocol" → 检测 ESC _G 序列
"Bug regression" → 300KB+ 超长行、ANSI 码前后的图片序列
"Negative cases" → 普通文本、文件路径不误判
"Mixed content" → 多种图片协议混合
```

**测试验证**（`bug-regression-isimageline-startswith-bug.test.ts`）：

```
回归测试：旧实现用 startsWith() 检测图片序列，
当图片不在行首时崩溃（300KB+ 行触发 "exceeds terminal width" 错误）
```

### 4.7 Autocomplete —— 自动补全

`CombinedAutocompleteProvider` 支持 slash 命令补全和文件路径补全。

**测试验证**（`autocomplete.test.ts`）：

```
路径前缀提取（从编辑器输入中识别路径）
fd 命令的文件/目录建议
大小写不敏感过滤、目录优先排序
空格路径自动加引号
隐藏文件处理（包含 .pi/.github，排除 .git）
引号路径补全
多级路径补全
```

---

## 五、测试架构

### 虚拟终端（`virtual-terminal.ts`）

所有测试都使用 `VirtualTerminal`——基于 `@xterm/headless` 的虚拟终端实现。它不需要真实终端，可以：

- 程序化注入输入（模拟按键）
- 读取渲染结果（获取屏幕上的文字）
- 控制终端尺寸

### 测试主题（`test-themes.ts`）

提供无样式的测试用主题，所有颜色函数返回原字符串（`(s) => s`），方便断言。

### 交互式 demo（`chat-simple.ts`）

一个完整的终端聊天界面 demo：

- Markdown 消息 + 自定义背景色
- 加载动画
- Editor + 自动补全 + slash 命令
- 消息间的间距

运行方式：`npx tsx test/chat-simple.ts`

---

## 六、设计亮点总结

| 设计点         | 实现方式                                             |
| -------------- | ---------------------------------------------------- |
| 无闪烁渲染     | CSI 2026 同步输出 + 差分更新                         |
| ANSI 安全      | 每行追加 SGR/OSC reset，overlay 切割时防止样式泄漏   |
| CJK/emoji 支持 | `Intl.Segmenter` + `eastAsianWidth` + RGI emoji 正则 |
| IME 支持       | `CURSOR_MARKER` 零宽标记 + 硬件光标定位              |
| 跨平台键盘     | Kitty 协议 + 传统序列双路解析                        |
| 非拉丁键盘     | Kitty alternate key 支持（base layout key）          |
| 流式安全       | Regional Indicator 宽度保守策略，防止差分漂移        |
| 大段粘贴       | Bracketed paste + >10 行折叠标记                     |
| 测试友好       | Terminal 接口抽象，VirtualTerminal 实现完整模拟      |
