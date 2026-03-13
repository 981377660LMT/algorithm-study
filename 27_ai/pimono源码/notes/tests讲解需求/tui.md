你好！作为新手，阅读和理解一个库的 README 文档是掌握它的最好方式。这是一份非常详细的终端用户界面（TUI）框架 `@mariozechner/pi-tui` 的说明文档。

我会带你**逐段对照、透彻讲解**这份文档，确保你弄懂所有的核心概念。

---

### 1. 标题与简介 (Title & Intro)

> Minimal terminal UI framework with differential rendering and synchronized output for flicker-free interactive CLI applications.

**讲解**：
这是一个**轻量级的终端 UI (TUI) 框架**。它的核心卖点是：

- **Differential rendering (差量渲染)**：只更新屏幕上发生变化的部分，而不是每次都重绘整个屏幕。
- **Synchronized output (同步输出)**：保证画面一次性输出完毕，这让基于命令行的交互应用（CLI）**不会出现画面闪烁 (flicker-free)**。

---

### 2. 核心特性 (Features)

- **Differential Rendering (差量渲染)**: 这是一种三策略渲染系统，能智能决定当前到底该怎么重绘界面以节省性能。
- **Synchronized Output (同步输出)**: 使用了 CSI 2026 这种特殊的控制码机制，让画面更新是“原子操作”（即瞬间完成，不闪屏）。
- **Bracketed Paste Mode (括号粘贴模式)**: 当用户在终端里一次性粘贴一大段文本（超过10行）时，它能够正确处理并添加标记。
- **Component-based (基于组件)**: 和 React/Vue 类似，它的设计是基于组件的。任何实现了 `render()` 方法的对象都可以是一个组件。
- **Theme Support (主题支持)**: 组件支持传入主题配置，方便你自定义颜色和样式。
- **Built-in Components (内置组件)**: 直接提供了现成好用的组件，比如文本(`Text`)、输入框(`Input`)、多行编辑器(`Editor`)、`Markdown` 渲染、加载动画(`Loader`)等。
- **Inline Images (内联图片)**: 支持在比较高级的终端（如 Kitty 或 iTerm2）里直接显示真彩图片！
- **Autocomplete Support (自动补全支持)**: 输入框支持文件路径和 `/` 斜杠命令的自动补全。

---

### 3. 快速开始 (Quick Start)

它给出了一个最小化运行的例子：

```typescript
import { TUI, Text, Editor, ProcessTerminal } from '@mariozechner/pi-tui'

// 1. 创建终端实例 (一般直接对接Node的输入输出流 stdin/stdout)
const terminal = new ProcessTerminal()

// 2. 创建 TUI 主机
const tui = new TUI(terminal)

// 3. 往里面塞入组件
tui.addChild(new Text('Welcome to my app!')) // 添加一行文本

const editor = new Editor(tui, editorTheme) // 创建一个多行编辑器组件
editor.onSubmit = text => {
  // 监听用户敲下回车的提交事件
  console.log('Submitted:', text)
  tui.addChild(new Text(`You said: ${text}`)) // 把用户输入的话再展示到屏幕上
}
tui.addChild(editor) // 把编辑器添加到界面

// 4. 正式启动 TUI
tui.start()
```

**讲解**：这个模式和前端开发很像：创建实例 -> 挂载组件 -> 启动渲染。

---

### 4. 核心 API (Core API)

#### TUI (主容器)

管理所有的组件和渲染进程。常用方法有 `addChild` (添加)，`removeChild` (移除)，`start`/`stop` 以及手动请求重绘 `tui.requestRender()`。快捷键 `Shift+Ctrl+D` 会触发它的 Debug 模式。

#### Overlays (悬浮层)

用于在现有内容“之上”渲染内容，且不会替换掉底下的内容。这对于实现**弹窗 (Dialogs)**、**菜单 (Menus)** 这种模态 UI 非常有用。

- 你可以非常详细地设置它的宽高（支持绝对像素或百分比）、位置（`row/col`）或锚点（如 `bottom-right` 右下角）。
- 甚至支持响应式隐藏功能：比如当终端宽度太小（`termWidth < 100`）时，自动隐藏该弹窗。

#### Component Interface (组件接口)

所有你想在这个 TUI 中显示的元素，都必须满足以下接口定义：

- `render(width: number): string[]`：这是最重要的方法。框架会传入当前终端的宽度 `width`，然后**组件必须要返回一个字符串数组，数组的每一项代表终端里的一行**。需要注意：**返回的任何一行字符长度绝对不能超过 `width`**，否则会报错。
- `handleInput(data: string)`（可选）：当组件处于焦点状态并收到用户键盘输入时被调用。
- `invalidate()`（可选）：用于清除组件内部的缓存，强制下一次渲染时完全重绘。

#### Focusable Interface (可聚焦与输入法 IME 支持)

如果你写了一个输入框组件，你需要让它实现 `Focusable` 接口。这很重要，因为它配合一个特殊的标记 `CURSOR_MARKER`，能够**让中文/日文/韩文输入法的候选词选框出现在正确的光标位置**，而不会飘在屏幕左上角。

---

### 5. 内置好用的组件 (Built-in Components)

文档罗列了作者已经写好的常用组件：

- **Container**: 容器，用来把别的组件打包塞在一起。
- **Box**: 带 padding (内边距) 和背景色的容器。
- **Text / TruncatedText**: 多行文本（自动换行）和 单行文本（超出宽度会自动省略为 `...`）。
- **Input**: 单行输入框，支持各种常用快捷键（Ctrl+A回到行首，Ctrl+W删词等）。
- **Editor**: 多行编辑器，支持文字折行展示、自动补全、大段粘贴（还能告诉你粘了多少行），当文本超过终端高度还会自动滚动！**它是整个框架最强大的组件。**
- **Markdown**: 支持把 Markdown 语法直接渲染得很好看（有加粗、颜色、列表甚至代码高亮）。
- **Loader / CancellableLoader**: 在做异步请求时展示的 Loading 动画（后者支持按 Escape 键取消加载）。
- **SelectList / SettingsList**: 交互式列表。你可以用方向键上下选择目标，按回车确认（非常适合做设置面板或下拉菜单）。
- **Spacer**: 用来占据一行或多行空白空间的占位组件。
- **Image**: 能在终端里渲染真正的图片（需要你的终端软件支持）。

---

### 6. 其他高级特性

#### 自动补全 (Autocomplete)

`CombinedAutocompleteProvider` 可以同时支持两种补全：按下 `/` 弹出命令补全列表，按下 `Tab` 补全电脑里的文件路径。

#### 键盘按键检测 (Key Detection)

为了处理用户按下了什么键，框架提供了一个 `matchesKey` 辅助函数。比如 `matchesKey(data, Key.ctrl("c"))` 就能判断用户是不是按了 Ctrl+C。

#### 差量渲染策略 (Differential Rendering)

这里作者解释了不闪屏的魔法原理：

1. **首次渲染**：老老实实输出所有行。
2. **终端窗口大小改变**：或者画面滚动发生了变化，由于`很难算差异，所以清空全屏重新渲染。`
3. **日常普通更新**：光标移动到发生改变的第一行，只把有变化的部分清掉重写，其他一概不碰。

#### 实用工具函数 (Utilities)

因为在终端开发中，含有颜色的文字会有一些隐藏的控制字符（ANSI codes，比如 `\x1b[31m` 代表红色）。这会导致 Javascript 去获取字符串长度时算错。
所以作者提供了非常有用的工具：

- `visibleWidth(str)`：算出这串具有颜色的文字**肉眼看起来**有多长。
- `truncateToWidth(str, width)`：截断这串具有颜色的文字，且非常安全地保证颜色不会乱掉。
- `wrapTextWithAnsi(str, width)`：让带颜色的文字安全地换行。

### 总结

这就是 README 的完整说明。作为新手，你只需要记住：

1. 先创建 `Terminal` 和 `TUI`。
2. 往 `TUI` 里 `addChild()` 添加类似 `Text` 或者 `Editor` 这样的现成组件。
3. `tui.start()` 跑起来即可。
   遇到复杂的样式或逻辑，多看看内置组件的 API 或自己尝试写个 `render()`。

希望这个透彻的讲解能帮你理清脉络！如果你想测试，可以自己在代码里跑跑文档末尾提到的 `npx tsx test/chat-simple.ts` 这个 Demo！
