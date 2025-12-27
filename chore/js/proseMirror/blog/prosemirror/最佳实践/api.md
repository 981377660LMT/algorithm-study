ProseMirror 是一个功能极其强大但也相对复杂的富文本编辑器框架。它不像 Quill 或 CKEditor 那样开箱即用，而更像是一套构建编辑器的“乐高积木”。

结合实际业务场景（如文档协作、CMS 后台、评论系统等），以下是 ProseMirror 的核心概念深入解析及最佳实践：

### 1. 核心架构与业务映射

ProseMirror 的四大核心模块对应着不同的业务需求层面：

- **`prosemirror-model` (数据层)**: 定义文档结构（Schema）。

  - _业务场景_: 你的业务是否允许嵌套列表？是否允许在引用中插入图片？是否需要自定义“提及(@user)”节点？
  - _最佳实践_: **严格定义 Schema**。不要试图在渲染层解决结构问题。如果业务规定“标题内不能有图片”，就在 Schema 中通过 `content` 属性限制它（例如 `content: "text*"` 而不是 `content: "inline*"`）。

- **`prosemirror-state` (状态层)**: 管理文档内容、选区(Selection)和插件状态。

  - _业务场景_: 实现撤销/重做、高亮搜索结果、字数统计。
  - _最佳实践_: **状态不可变性 (Immutability)**。永远不要直接修改 `state.doc`。所有的变更必须通过 `Transaction`（事务）来完成。这对于实现协同编辑至关重要。

- **`prosemirror-view` (视图层)**: 将状态渲染为 DOM，并处理用户交互。

  - _业务场景_: 自定义图片上传进度的展示、点击 `@user` 弹出用户信息卡片。
  - _最佳实践_: 使用 **Node Views** 处理复杂交互。对于非文本内容（如视频、复杂的数学公式、React/Vue 组件嵌入），不要依赖默认渲染，而是编写自定义的 `NodeView`。

- **`prosemirror-transform` (变更层)**: 处理文档的修改步骤 (Steps)。
  - _业务场景_: 协同编辑冲突解决、Markdown 粘贴转换。
  - _最佳实践_: 理解 **Mapping**。当你保留一个位置（比如评论的位置）时，如果文档前面插入了文字，你的位置索引需要通过 `mapping` 进行更新，否则会错位。

---

### 2. 常见业务场景的最佳实践

#### 场景 A：自定义数据结构 (Schema Design)

**需求**: 业务需要一个“警告框(Callout)”组件，包含一个图标和文本内容。

- **错误做法**: 使用 CSS class 给 `<p>` 标签加样式。
- **最佳实践**: 定义一个新的 Node 类型。
  ```javascript
  // Schema 定义
  callout: {
    content: "block+", // 允许包含段落等块级元素
    group: "block",
    attrs: { type: { default: "info" } }, // 存储类型：info, warning, error
    toDOM(node) { return ["div", { class: `callout ${node.attrs.type}` }, 0] }, // 0 代表内容坑位
    parseDOM: [{ tag: "div.callout", getAttrs: dom => ({ type: dom.classList.contains("error") ? "error" : "info" }) }]
  }
  ```
  _理由_: 这样在数据层面上它就是独立的结构，方便后续做转换（如转 Markdown）或特定逻辑处理。

#### 场景 B：协同编辑 (Collaborative Editing)

**需求**: 多人同时编辑同一文档。

- **核心机制**: ProseMirror 内置了基于 **Steps** 的协同机制（类似 OT 算法）。
- **最佳实践**:
  1.  **Central Authority (服务端)**: 服务端必须是“真理的来源”。它接收 Steps，按顺序应用，并广播给其他客户端。
  2.  **Rebasing (变基)**: 客户端 A 提交了 Step 1，但在服务端确认前，收到了客户端 B 的 Step 2。客户端 A 必须撤销本地的 Step 1，应用 Step 2，然后重新计算 Step 1（变基），再尝试提交。
  3.  **使用 `prosemirror-collab`**: 不要自己手写同步逻辑，直接使用官方插件。它处理了复杂的 Step 追踪和变基逻辑。

#### 场景 C：UI 框架集成 (React/Vue Integration)

**需求**: 在编辑器中插入一个复杂的 React 组件（例如一个动态图表或日历控件）。

- **挑战**: ProseMirror 控制 DOM，React 也想控制 DOM，两者会打架。
- **最佳实践**: 使用 **Node Views** 作为桥梁。
  1.  创建一个 `NodeView` 类。
  2.  在 `dom` 属性中创建一个容器元素。
  3.  使用 `ReactDOM.render` (或 `createRoot`) 将 React 组件挂载到这个容器中。
  4.  **关键点**: 在 `ignoreMutation` 和 `stopEvent` 方法中告诉 ProseMirror 忽略 React 组件内部的 DOM 变化和事件，防止编辑器重绘破坏 React 状态。

#### 场景 D：元数据与装饰器 (Decorations)

**需求**: 实现“拼写检查”或“搜索高亮”，但不改变文档实际内容。

- **错误做法**: 修改文档模型，插入 `<span class="highlight">` 节点。
- **最佳实践**: 使用 **Decorations**。
  - Decorations 是视图层的概念，它在渲染时“假装”修改了 DOM，但不会污染 `state.doc`。
  - 通过 Plugin 的 `props: { decorations(state) { ... } }` 返回 DecorationSet。
  - 这样，当你保存文档数据时，这些高亮不会被保存进去。

#### 场景 E：大文档性能优化

**需求**: 编辑几万字的长文档不卡顿。

- **最佳实践**:
  1.  **不要频繁全量序列化**: 只有在保存时才调用 `state.doc.toJSON()`。
  2.  **Plugin 优化**: 确保你的 Plugin 中的 `apply` 方法非常快。如果需要重计算（如语法高亮），考虑防抖 (Debounce) 或增量更新。
  3.  **NodeView 更新**: 在 NodeView 的 `update` 方法中，仔细比对新旧节点。如果属性没变，直接返回 `true` 阻止 DOM 重绘。

### 3. 代码组织建议

不要把所有代码写在一个文件里。建议按功能模块划分：

- `/schema`: 定义 nodes 和 marks。
- `/plugins`: 独立的逻辑模块（如 `KeymapPlugin`, `HistoryPlugin`, `UploadPlugin`）。
- `/commands`: 可复用的编辑命令（如 `insertImage`, `toggleBold`）。
- `/nodeviews`: 自定义渲染逻辑。

### 4. 调试技巧

- **ProseMirror DevTools**: 这是一个 Chrome 插件（也有库版本），必装。它可以让你实时查看 State、Document 结构、Selection 位置和 History 栈。
- **理解 Index**: 能够区分 `Node Index` (在父节点中的位置) 和 `Document Position` (全局位置) 是调试这类问题的关键。

### 总结

ProseMirror 的学习曲线陡峭，是因为它把**控制权**完全交给了开发者。在业务开发中，**先设计 Schema，再通过 Plugin 扩展行为，最后用 NodeView 完善表现**，是通用的最佳路径。

### 4. 调试技巧

- **ProseMirror DevTools**: 这是一个 Chrome 插件（也有库版本），必装。它可以让你实时查看 State、Document 结构、Selection 位置和 History 栈。
- **理解 Index**: 能够区分 `Node Index` (在父节点中的位置) 和 `Document Position` (全局位置) 是调试这类问题的关键。

### 5. 决策指南：什么时候用什么？ (When to use What?)

在 ProseMirror 开发中，经常面临技术选型的困惑。以下是针对常见场景的决策树：

#### Node vs Mark

- **用 Node (节点)**: 当内容是**结构化**的，通常独占空间，或者不能被分割时。
  - _例子_: 图片、代码块、水平分割线、视频卡片。
  - _判断标准_: 如果你删除了它的一半，剩下的一半还能存在吗？如果不能（比如半张图片没意义），那就是 Node。
- **用 Mark (标记)**: 当内容是**附着在文本上**的元数据或样式，且可以跨越多个节点时。
  - _例子_: 加粗、斜体、链接、文本颜色、行内评论。
  - _判断标准_: 它可以像衣服一样穿在文字身上吗？如果是，那就是 Mark。

#### toDOM vs NodeView

- **用 toDOM**: 当节点只需要**静态渲染**为 HTML，且没有复杂的交互逻辑时。
  - _例子_: 简单的段落、标题、列表、引用。
  - _优势_: 性能最好，代码最少。
- **用 NodeView**: 当节点需要**自定义交互**、复杂的 DOM 结构、或者需要集成前端框架（React/Vue）时。
  - _例子_: 可点击展开的代码块、带进度条的上传占位符、嵌入的看板卡片。
  - _优势_: 完全控制 DOM 生命周期，可以处理事件（`stopEvent`）、控制光标行为（`ignoreMutation`）。

#### Plugin vs Command

- **用 Command (命令)**: 当你需要执行一个**一次性**的操作，通常由用户主动触发。
  - _例子_: 点击工具栏按钮变粗体、按下回车键换行、粘贴内容。
  - _本质_: 一个函数 `(state, dispatch) => boolean`。
- **用 Plugin (插件)**: 当你需要**持续监听**编辑器状态，或者拦截/修改默认行为时。
  - _例子_: 实时字数统计、Markdown 快捷键输入转换（输入 `# ` 变标题）、协同编辑同步、撤销重做栈。
  - _本质_: 一个包含 `state`、`view` 钩子和 `props` 的对象。

#### Decoration vs Schema

- **用 Schema (数据层)**: 当这个特性是**文档内容的一部分**，需要被保存到数据库并在其他地方（如阅读页）展示时。
  - _例子_: 待办事项的勾选状态、图片的 src 属性。
- **用 Decoration (装饰器)**: 当这个特性只是**临时的视觉效果**，或者是基于本地状态的辅助信息时。
  - _例子_: 搜索关键词高亮、当前行高亮、协同编辑时别人的光标、拼写错误波浪线。

### 总结

ProseMirror 的学习曲线陡峭，是因为它把**控制权**完全交给了开发者。在业务开发中，**先设计 Schema，再通过 Plugin 扩展行为，最后用 NodeView 完善表现**，是通用的最佳路径。

---

ProseMirror 的强大在于其极高的可定制性，但这也带来了选择困难。在实际业务开发中，最常见的问题就是：“实现这个功能，我到底该用 Node 还是 Mark？用 Plugin 还是 Command？”

以下是结合业务场景的详细决策指南与最佳实践。

### 1. 决策指南：什么时候用什么？ (When to use What?)

这是开发中最常遇到的技术选型路口，以下是详细的判断标准：

#### A. Node (节点) vs Mark (标记)

- **用 Node**: 当内容是**结构化**的，通常独占空间，或者**不可分割**时。

  - _场景_: 图片、视频卡片、代码块、水平分割线、数学公式（块级）。
  - _判断标准_: 如果你删除了它的一半，剩下的一半还能存在吗？如果不能（比如半张图片没意义），那就是 Node。
  - _细节_: Node 可以有属性（attrs），比如图片的 `src`，代码块的 `language`。

- **用 Mark**: 当内容是**附着在文本上**的元数据或样式，且可以**跨越**多个节点时。
  - _场景_: 加粗、斜体、链接、文本颜色、行内评论、高亮。
  - _判断标准_: 它可以像衣服一样穿在文字身上吗？如果是，那就是 Mark。
  - _细节_: Mark 也可以有属性，比如链接的 `href`。

#### B. toDOM (简单渲染) vs NodeView (自定义视图)

- **用 toDOM**: 当节点只需要**静态渲染**为 HTML，且没有复杂的交互逻辑时。

  - _场景_: 简单的段落 (`<p>`)、标题 (`<h1>`)、列表 (`<ul>`)、引用 (`<blockquote>`)。
  - _优势_: 性能最好，代码最少，ProseMirror 自动处理更新。

- **用 NodeView**: 当节点需要**自定义交互**、复杂的 DOM 结构、或者需要集成前端框架（React/Vue）时。
  - _场景_:
    - **交互**: 点击展开/折叠的代码块、带进度条的上传占位符。
    - **框架集成**: 在编辑器里嵌入一个 React 的看板卡片、日历控件。
    - **生命周期**: 需要在节点销毁时清理定时器或事件监听。
  - _关键 API_: `stopEvent` (阻止事件冒泡给编辑器), `ignoreMutation` (告诉编辑器忽略内部 DOM 变化)。

#### C. Plugin (插件) vs Command (命令)

- **用 Command**: 当你需要执行一个**一次性**的操作，通常由用户主动触发。

  - _场景_: 点击工具栏按钮变粗体、按下回车键换行、粘贴内容、点击菜单插入图片。
  - _本质_: 一个函数 `(state, dispatch) => boolean`。如果返回 `true`，表示命令已执行（或适用）。

- **用 Plugin**: 当你需要**持续监听**编辑器状态，或者拦截/修改默认行为时。
  - _场景_:
    - **监听**: 实时字数统计、光标位置监控。
    - **拦截**: 阻止用户删除特定的只读节点、Markdown 快捷键输入转换（输入 `# ` 变标题）。
    - **增强**: 协同编辑同步（监听所有 Transaction）、撤销重做栈。
  - _本质_: 一个包含 `state`、`view` 钩子和 `props` 的对象，生命周期伴随编辑器始终。

#### D. Decoration (装饰器) vs Schema (数据模型)

- **用 Schema (数据层)**: 当这个特性是**文档内容的一部分**，需要被保存到数据库并在其他地方（如阅读页）展示时。

  - _场景_: 待办事项的勾选状态、图片的 src 属性、表格的列宽。

- **用 Decoration (视图层)**: 当这个特性只是**临时的视觉效果**，或者是基于本地状态的辅助信息时。
  - _场景_:
    - **搜索高亮**: 搜索词变黄，保存文档时不需要保存黄色背景。
    - **协同光标**: 看到别人的光标位置。
    - **代码查错**: 拼写错误下方的波浪线。
    - **UI 辅助**: 选中表格单元格时的蓝色边框。

---

### 2. 核心业务场景最佳实践

#### 场景一：实现“提及”功能 (@User)

- **选型**: `Node` (Inline Node) + `NodeView`。
- **为什么是 Node?** 虽然它看起来像文本，但它是一个整体（原子性）。你通常不希望用户把光标放在 `@User` 的名字中间修改它。
- **Schema 定义**:

  ```javascript
  mention: {
    group: "inline",
    inline: true,
    atom: true, // 关键：原子节点，光标无法进入内部
    attrs: { id: {}, name: {} },
    toDOM: node => ["span", { class: "mention", "data-id": node.attrs.id }, "@" + node.attrs.name]
  }
  ```

- **交互**: 使用 `InputRule` 插件监听 `@` 输入，触发一个弹出层（建议用 `Decoration` 渲染弹出层锚点，或者使用 `prosemirror-suggest` 库）。

#### 场景二：大文档性能优化

- **问题**: 文档超过 5 万字，输入卡顿。
- **优化策略**:
  1.  **减少 `toDOM` 重绘**: 在 `NodeView` 的 `update` 方法中，严格比对新旧节点属性。如果属性没变，直接返回 `true`，阻止 DOM 重绘。
  2.  **按需序列化**: 不要每次 `doc` 变化都调用 `doc.toJSON()`。使用 `debounce` 防抖，或者只在保存按钮点击时序列化。
  3.  **Plugin 瘦身**: 检查所有 Plugin 的 `apply` 方法。这里是热路径，每次按键都会触发。不要在这里做复杂的计算（如全量正则扫描）。

#### 场景三：协同编辑 (Collaborative Editing)

- **核心**: 必须使用 `prosemirror-collab` 插件。
- **服务端**: 服务端必须是“真理的来源 (Source of Truth)”。
- **坑点**:
  - **Step 丢失**: 客户端断网重连时，必须重新拉取服务端最新版本，而不是盲目提交本地积压的 Steps。
  - **ID 冲突**: 如果你的 Node 有 `id` 属性，协同编辑时两个人同时插入节点可能导致 ID 重复。需要在服务端解决 ID 冲突，或者使用 UUID。

#### 场景四：混合排版（React/Vue 组件嵌入）

- **场景**: 在文档里插入一个可交互的“投票组件”。
- **最佳实践**:
  - 创建一个 `NodeView`。
  - 在 `dom` 属性中创建一个 `div` 容器。
  - 使用 `ReactDOM.createRoot(div).render(...)` 挂载组件。
  - **重要**: 在 `destroy()` 方法中调用 `root.unmount()` 防止内存泄漏。
  - **通信**: React 组件修改数据时，不要直接改 props，而是调用 `view.dispatch` 发送 Transaction 修改 ProseMirror 的 Node 属性，ProseMirror 更新后再通过 `update` 方法流回 React 组件（单向数据流）。

### 3. 调试技巧

- **ProseMirror DevTools**: 必装 Chrome 插件。可以查看当前的 Document 树结构、Selection 位置、Plugin 状态。
- **理解 Index**:
  - `pos`: 绝对位置，从文档开头算起（0, 1, 2...）。
  - `$pos.index(depth)`: 在当前层级父节点中的索引（第几个子节点）。
  - 调试时打印 `view.state.selection.$from` 可以看到详细的位置信息。
