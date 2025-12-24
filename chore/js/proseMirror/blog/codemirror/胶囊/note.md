基于对工作区代码的深入分析，富文本编辑器中的“胶囊”（Capsule/Pill）功能主要集中在 **AI Chat（对话机器人）** 场景中，基于 **CodeMirror 6** 实现。而在 **Lander Code Editor（低代码编辑器）** 场景中，则采用了不同的技术路线（语法高亮与智能提示）。

以下是针对这两种不同产品形态的技术细节分析：

### 1. AI Chat 场景：原子化胶囊（Atomic Capsule）

在 AI 对话输入框中，用户引用的上下文、变量或指令需要被视为一个整体，不可部分编辑，只能整体删除或移动。

#### 核心实现方案

- **底层引擎**：CodeMirror 6
- **数据存储**：**特殊标记文本**。胶囊在底层数据模型（`Doc`）中实际上是一段包含特殊字符的纯文本。
  - **格式**：`@名称\u200B`。
  - **关键点**：使用 **零宽空格 (`\u200B`)** 作为胶囊的结束边界标记（Boundary），配合 `@` 前缀，使得正则可以精准匹配，同时对用户视觉不可见。
- **视图渲染**：**Decoration & Widget**。
  - 使用 `ViewPlugin` 扫描文档内容。
  - 使用 `Decoration.replace({ widget })` 将匹配到的文本范围（如 `@变量A\u200B`）在视图层替换为一个自定义的 DOM 节点（胶囊组件）。

#### 关键技术细节

1.  **原子化光标控制 (`Atomic Ranges`)**

    - 为了让胶囊表现得像一个“字符”，必须禁止光标停留在胶囊文本的中间。
    - **实现**：在 `createCapsulePlugin` 中通过 `provide` 属性注册 `EditorView.atomicRanges`。
    - **效果**：当用户移动光标时，CodeMirror 会自动跳过整个胶囊的文本范围，直接从胶囊头跳到尾，反之亦然。

2.  **正则匹配与渲染循环**

    - **代码位置**：`packages/chat-bot-sdk/src/components/Input/TextArea/extensions/capsule-plugin.ts`
    - **逻辑**：
      ```typescript
      // 伪代码逻辑
      const regex = /@([^@\s\u200B]+)\u200B/g; // 匹配 @name + 零宽空格
      while ((match = regex.exec(doc)) !== null) {
          // 创建替换型装饰器，将文本替换为 Widget
          const widget = new MentionCapsuleWidget(name, ...);
          decorations.push(Decoration.replace({ widget }).range(from, to));
      }
      ```

3.  **自定义 Widget (`MentionCapsuleWidget`)**

    - **代码位置**：`packages/chat-bot-sdk/src/components/Input/TextArea/widgets/MentionCapsuleWidget.ts`
    - **功能**：
      - **渲染**：生成包含图标、文本、删除按钮的 DOM 结构。
      - **交互**：接管点击事件（`onclick`），触发外部回调（如查看引用详情）。
      - **删除**：内置删除按钮逻辑，点击时计算出该 Widget 对应的文档位置（`from`, `to`），调用 `view.dispatch` 删除底层文本。
    - **状态同步**：根据 `globalState` 判断胶囊是否有效（如引用是否已失效），并渲染不同样式（灰色/高亮）。

4.  **输入法（IME）兼容性**
    - **问题**：在中文输入过程中，文档内容频繁变化，如果实时触发正则匹配和 DOM 替换，会导致输入法选词框抖动或中断。
    - **解决**：在 `ViewPlugin` 的 `update` 方法中检查 `composingRef.current`（是否正在合成输入）。如果是，则跳过重新计算装饰器的逻辑，仅做位置映射，直到 `compositionend` 事件触发。

---

### 2. Lander Code Editor 场景：语法高亮与绑定（Syntax Highlighting）

在低代码平台的 JS/SQL 编辑器中，用户需要编写复杂的表达式（如 `{{ table1.data }}`）。这里的“变量”更强调**可编辑性**和**代码语义**，而不是原子化的 UI 块。

#### 核心实现方案

- **底层引擎**：CodeMirror 6
- **数据存储**：**纯代码文本**。如 `select * from table1`。
- **视图渲染**：**语法高亮 (`HighlightStyle`)**。

#### 关键技术细节

1.  **基于语法树的着色**

    - **代码位置**：`packages/code-editor/src/extensions/style/jsHighlight.tsx`
    - **逻辑**：使用 `Lezer` 解析器生成语法树（Syntax Tree），通过 `HighlightStyle.define` 为不同的语法节点（`tags.keyword`, `tags.string` 等）定义颜色。
    - **差异**：不像 Chat 场景那样替换 DOM，而是给文本套上 `<span>` 并应用 CSS 类名。

2.  **智能提示与补全 (`Completion`)**

    - **代码位置**：`packages/code-editor/src/extensions/completion/UI/CompletionInfoCard.tsx`
    - **逻辑**：当用户输入变量名时，弹出的不是胶囊选择器，而是代码补全列表。选中后插入的是纯文本。
    - **增强**：通过 `HoverTooltip`，当鼠标悬停在变量文本上时，显示该变量的当前值（Evaluation）或类型定义，达到类似胶囊的信息展示效果，但不牺牲编辑灵活性。

3.  **插值语法处理 (`JSBindingInput`)**
    - 在 `JSBindingInput` 组件中，处理 `{{ ... }}` 这种插值语法。虽然视觉上可能通过背景色区分代码块和普通字符串，但本质上仍然是文本编辑，允许用户修改变量名的任意字符（如把 `table1` 改为 `table2`）。

### 总结对比

| 特性         | AI Chat (胶囊方案)                            | Lander Code Editor (代码方案)    |
| :----------- | :-------------------------------------------- | :------------------------------- |
| **视觉形态** | 独立的 UI 块（圆角、图标、背景色）            | 彩色文本                         |
| **编辑能力** | **原子化**（整体删除/移动，不可修改内部文字） | **流式**（可自由修改任意字符）   |
| **底层数据** | 带特殊标记的文本 (`@name\u200B`)              | 纯代码文本                       |
| **实现技术** | `Decoration.replace` + `WidgetType`           | `HighlightStyle` + `Syntax Tree` |
| **光标行为** | 跳过整个胶囊 (`Atomic Ranges`)                | 字符级移动                       |
| **适用场景** | 引用实体、上下文、指令、提及用户              | 编写逻辑、SQL 查询、变量绑定     |

这种差异化的实现方案精准地匹配了不同场景下的用户心智：聊天时关注“引用对象”，编程时关注“逻辑表达”。

---

这是一个**非常成熟且巧妙**的工程化方案，特别适合在 CodeMirror 6 中实现“类 Notion”或“低代码公式编辑器”的胶囊效果。

使用 `\u200B` (Zero Width Space, ZWSP) 作为锚点（Anchor）是解决纯文本编辑器中“胶囊边界判定”问题的经典手段。

以下是对该方案的深度评审，包括它的**优势**、**潜在风险**以及**关键实现细节**。

### 一、 方案核心优势分析

1.  **边界判定极其精准 (Precision)**

    - **痛点：** 纯正则匹配 `@User` 最大的问题是贪婪匹配或边界模糊。比如用户输入 `email@User.com`，正则很难区分这是邮箱还是胶囊。或者用户想在胶囊后紧接着输入文字，光标往往“粘”在胶囊文本上。
    - **ZWSP 的作用：** `\u200B` 充当了一个**隐形的物理屏障**。
      - 正则可以写成：`/@(.+?)\u200B/g`。
      - 只有当 `\u200B` 存在时，才渲染为胶囊。这保证了用户正在输入 `@User` 时（还未确认），它只是普通文本；一旦确认（插入 `\u200B`），它瞬间变为胶囊。

2.  **光标导航体验极佳 (Navigation)**

    - CodeMirror 的 `Atomic Ranges` 配合 ZWSP 简直是绝配。
    - 当光标移动到 `\u200B` 之后时，实际上已经跳出了胶囊的逻辑范围。这解决了“光标在胶囊内还是外”的薛定谔状态。

3.  **数据序列化简单 (Serialization)**
    - 存储到数据库的就是纯字符串：`"Result = @Price\u200B * @Count\u200B"`。
    - 不需要复杂的 JSON 结构，不需要维护额外的 Offset Map。任何支持 UTF-8 的后端都能存。

---

### 二、 潜在风险与挑战 (Risk Assessment)

虽然方案很棒，但在落地时必须处理以下“坑”：

#### 1. 复制粘贴的“脏数据”问题

- **场景：** 用户复制了包含胶囊的文本，粘贴到 Excel 或 VS Code 中。
- **后果：** 粘贴出来的文本里会夹杂着看不见的 `\u200B`。
  - 在 Excel 中可能导致公式报错。
  - 在代码中可能导致语法错误（Invisible Character）。
- **对策：**
  - **拦截 Copy 事件：** 监听编辑器的 `copy` 或 `cut` 事件，获取选区文本，执行 `text.replace(/\u200B/g, "")` 清洗后再写入剪贴板。
  - **或者：** 接受这个特性，如果你的产品生态是封闭的（只在你的产品内粘贴），保留 ZWSP 反而是好事，能保留胶囊状态。

#### 2. 删除行为的“半截”问题

- **场景：** 用户在胶囊后按 Backspace。
- **期望：** 整个胶囊被删除。
- **实际（若未处理）：** 删除了 `\u200B`。
  - 正则 `/ @(.+?)\u200B /` 匹配失效。
  - 胶囊瞬间“炸开”，变回了普通文本 `@名称`。
- **对策：**
  - 必须配置 `EditorView.atomicRanges`。将 `@名称\u200B` 整个范围标记为原子性。这样 Backspace 会一次性删除整个范围，而不会只删掉末尾的字符。

#### 3. 搜索匹配干扰

- **场景：** 用户按 Ctrl+F 搜索“名称”。
- **问题：** 胶囊内的文本是 `@名称\u200B`。如果搜索逻辑是纯文本匹配，能搜到；但如果用户想搜“@名称”，可能会因为漏打 ZWSP 而搜不到（取决于搜索实现）。通常这不是大问题。

---

### 三、 关键代码实现细节 (Implementation Guide)

这是基于你描述方案的 CodeMirror 6 核心实现代码：

#### 1. 正则匹配器 (Matcher)

```typescript
import { MatchDecorator, Decoration, ViewPlugin, EditorView, WidgetType } from '@codemirror/view'

// 1. 定义胶囊外观
class CapsuleWidget extends WidgetType {
  constructor(readonly label: string) {
    super()
  }

  eq(other: CapsuleWidget) {
    return other.label === this.label
  }

  toDOM() {
    const span = document.createElement('span')
    span.className = 'cm-capsule'
    span.textContent = this.label // 显示时不包含 @ 和 \u200B
    return span
  }

  ignoreEvent() {
    return true
  }
}

// 2. 定义匹配逻辑
const capsuleMatcher = new MatchDecorator({
  // 核心正则：匹配 @ 开头，\u200B 结尾的非贪婪片段
  // 注意：这里假设名称中不包含 \u200B
  regexp: /@([^@\u200B]+)\u200B/g,

  decoration: match => {
    const label = match[1] // 提取名称
    return Decoration.replace({
      widget: new CapsuleWidget(label),
      // 关键：inclusive: false 确保光标不会落在替换区域的边界内侧
      inclusive: false
    })
  }
})

// 3. 视图插件
export const capsulePlugin = ViewPlugin.fromClass(
  class {
    decorations: any
    constructor(view: EditorView) {
      this.decorations = capsuleMatcher.createDeco(view)
    }
    update(update: any) {
      if (update.docChanged || update.viewportChanged) {
        this.decorations = capsuleMatcher.updateDeco(update, this.decorations)
      }
    }
  },
  {
    decorations: v => v.decorations,
    // 4. 关键：原子范围，保证删除时整体删除
    provide: plugin =>
      EditorView.atomicRanges.of(view => {
        return view.plugin(plugin)?.decorations || Decoration.none
      })
  }
)
```

#### 2. 插入胶囊的辅助函数 (Helper)

你需要提供一个命令或函数，供外部（如点击菜单）调用来插入胶囊。

```typescript
import { EditorView } from '@codemirror/view'

export function insertCapsule(view: EditorView, name: string) {
  const transaction = view.state.update({
    changes: {
      from: view.state.selection.main.head,
      // 插入格式化文本：@ + 名称 + ZWSP
      insert: `@${name}\u200B`
    },
    // 插入后，光标移动到 ZWSP 后面
    selection: { anchor: view.state.selection.main.head + name.length + 2 }
  })
  view.dispatch(transaction)
  view.focus()
}
```

### 四、 结论

**评分：9/10**

这个方案是 CodeMirror 生态中处理“结构化内联数据”的**最优解之一**。

- 它避开了复杂的 State Field 维护。
- 它利用了 ZWSP 完美解决了正则匹配的边界问题。
- 它利用了 Atomic Ranges 完美解决了光标交互问题。

**唯一需要注意的**就是数据的“出口清洗”（Export Sanitization），即在数据离开编辑器（复制、保存到后端）时，是否需要剥离 `\u200B`，或者后端是否能兼容这个字符。只要处理好这一点，这就是一个生产级可用的坚固方案。

---

在 CodeMirror 6 (CM6) 中，除了“胶囊（Capsule）”这种利用 `ReplaceDecoration` + `Widget` 的技巧外，还有许多类似的“黑魔法”可以用来构建富交互的编辑器。

这些技巧的核心逻辑通常都是：**分离“数据（Doc）”与“表现（View）”**。

以下是几个经典的技巧模式及其伴随的挑战：

### 技巧一：幽灵文本 (Ghost Text / In-line Suggestion)

**场景：** AI 代码补全（如 GitHub Copilot）、输入提示占位符。
**原理：** 在光标后插入一段**不存在于文档中**的灰色文本。

- **实现方式：**
  - 使用 `Decoration.widget({ side: 1 })`。
  - `side: 1` 确保它永远紧贴在光标右侧。
  - 这个 Widget 不替换任何文本，只是“挤”在两个字符中间。
- **挑战：**
  - **换行处理：** 如果幽灵文本很长，是应该换行显示还是截断？CM 的 Widget 默认是 inline-block，处理多行对齐非常痛苦。
  - **光标穿透：** 用户按右键时，是应该把幽灵文本“吃进去”（采纳建议），还是仅仅穿过它？这需要拦截 `keydown` 事件并自定义 Transaction。

### 技巧二：代码折叠作为“数据隐藏” (Folding as Data Hiding)

**场景：** 隐藏复杂的元数据（如 YAML Frontmatter），或者将长 URL 缩短显示。
**原理：** 利用 CM 的折叠（Folding）机制，但不是为了折叠代码块，而是为了隐藏“脏数据”。

- **实现方式：**
  - 底层文本：`[链接描述](https://very-long-url.com/token=xyz...)`
  - 视图：使用 `Decoration.replace` 将 `(...)` 部分替换为一个小的 `🔗` 图标 Widget。
- **挑战：**
  - **搜索（Find）：** 浏览器原生的 Ctrl+F 或 CM 的 Search 插件会搜到底层文本。如果用户搜到 `token=xyz`，视图会滚动到那里，但用户看到的是一个图标，会非常困惑。
  - **编辑边界：** 用户想修改 URL 的一部分怎么办？通常需要实现“点击图标展开编辑，失焦自动折叠”的逻辑。

### 技巧三：块级挂件 (Block Widgets / Line Injections)

**场景：** 在代码行之间插入图片预览、SQL 查询结果表格、错误提示条。
**原理：** 在两行文本之间“硬塞”一个 DOM 块。

- **实现方式：**
  - `Decoration.widget({ block: true, side: 1 })`。
  - 这会让 Widget 独占一行高度，不影响行号计算（视觉上在行与行之间）。
- **挑战：**
  - **高度抖动 (Layout Thrashing)：** 如果 Widget 内部加载图片（异步），加载完成后高度突变，会导致编辑器滚动条跳动。必须预设高度或使用 `requestMeasure` 通知 CM 重绘。
  - **选区跨越：** 当用户全选（Cmd+A）或上下拖拽选区经过这个块时，视觉反馈是什么？CM 默认会忽略块级 Widget，导致选区看起来“断开”了。

### 技巧四：只读区间 (Read-only Ranges)

**场景：** 模板编辑器，其中某些变量 `{{company_name}}` 不允许用户修改，只能整体删除。
**原理：** 结合 `EditorState.readOnly` 和 `Filter`。

- **实现方式：**
  - 不仅仅是 `contenteditable="false"`。
  - 需要使用 `EditorState.transactionFilter`。
  - 监听每一个 Transaction，如果它试图修改只读范围内的文本，就拦截并修改该 Transaction（例如阻止修改，或者将修改范围扩大到删除整个变量）。
- **挑战：**
  - **协同编辑冲突：** 如果 A 用户锁定了某行，B 用户在协同中删除了该行，Filter 逻辑可能会导致协同算法（OT/CRDT）崩溃或状态不一致。

---

### 核心挑战总结：当“所见”不等于“所得”

所有上述技巧（包括胶囊）都会面临以下三个共性挑战：

#### 1. 坐标映射地狱 (Coordinate Mapping)

- **问题：** 屏幕上的 `(x, y)` 坐标不再直接对应文档的 `index`。
- **场景：** 你想在鼠标悬停在胶囊上时显示 Tooltip。
- **API：** 必须熟练使用 `view.posAtCoords` 和 `view.coordsAtPos`。如果使用了 `ReplaceDecoration`，视觉高度和逻辑行高不一致，计算非常容易出错。

#### 2. 差异化更新 (Diffing & Re-rendering)

- **问题：** Widget 是 DOM 节点。如果用户打字很快，CM 会频繁触发 `update`。
- **性能陷阱：** 如果你的 `Widget.toDOM` 每次都返回新的 DOM 元素，编辑器会疯狂闪烁且性能低下。
- **解决：** 必须在 `Widget` 类中实现 `eq(otherWidget)` 方法。只有当数据真正变化时，才允许 CM 销毁旧 DOM 创建新 DOM。

#### 3. 协同编辑 (Collab) 中的位置偏移

- **问题：** 你在位置 10 插入了一个胶囊。队友在位置 0 插入了 5 个字。你的胶囊应该自动变到位置 15。
- **机制：**
  - 如果是 `MatchDecorator`（基于正则实时扫描），这通常不是问题，因为文本变了，正则会重新匹配。
  - 如果是 **StateField** 维护的装饰器（手动添加的），必须在 `update` 方法中调用 `decoration.map(transaction.changes)`。忘记这一步，胶囊就会“钉”在原地，与文本错位。

### 进阶建议：StateField vs ViewPlugin

- **ViewPlugin (推荐用于胶囊)：** 适合“由文本内容决定外观”的场景（如正则匹配）。逻辑简单，性能好，不污染 State。
- **StateField (推荐用于复杂交互)：** 适合“由用户操作决定状态”的场景（如点击行号添加断点、折叠代码）。因为这些状态不是写在文档文本里的，必须存在 StateField 里才能持久化或同步。

---

## review

根据建议修复重构代码

这段代码实现了一个基本的胶囊插件逻辑，但在**性能**和**CodeMirror 最佳实践**方面存在几个严重的问题，特别是在处理大文档或高频更新时可能会导致卡顿。

以下是详细的代码 Review 和重构建议。

### 🛑 关键问题 (Critical Issues)

#### 1. 严重的性能隐患：`doc.toString()`

```typescript
// ❌ 每一帧 update 都会执行
const doc = view.state.doc.toString()
```

- **问题**：`view.state.doc.toString()` 会将整个文档（可能是几万行代码）拼接成一个巨大的字符串。
- **后果**：每次用户打字（触发 `update`），都会触发全量字符串转换和全量正则匹配（`matchCapsules`）。在长文档中，这将导致严重的输入延迟（Lag）。
- **最佳实践**：CodeMirror 6 设计为**增量更新**。应该只扫描视口（Viewport）内的可见区域，或者使用 `MatchDecorator`（它会自动处理视口优化）。

#### 2. 删除逻辑中的“位置过期”风险

```typescript
// ❌ originalFrom 是创建 Widget 时的位置，不是点击时的位置
const handleDeleteCapsule = (name: string, originalFrom: number, view: EditorView) => { ... }
```

- **问题**：`originalFrom` 是在创建装饰器那一刻捕获的。如果用户在胶囊**之前**输入了文本，胶囊的实际位置会后移，但 `originalFrom` 还是旧值。
- **后果**：点击删除时，可能会删除错误的位置，或者因为找不到匹配而失效。
- **修复**：不要把位置传给 Widget。在 Widget 的事件处理函数中，使用 `view.posAtDOM(event.target)` 来获取当前最新的位置。

#### 3. Widget 设计反模式

```typescript
// ❌ Widget 不应该持有 absolute position (from, to)
const widget = new MentionCapsuleWidget(name, from, to, ...);
```

- **原理**：在 CM6 中，Widget 应该是纯粹的渲染单元。位置信息由 `Decoration` 管理。当文档变化时，Decoration 会自动映射位置，但 Widget 实例内部的 `from` 属性不会自动更新。

#### 4. 排序效率

```typescript
// ⚠️ 每次都重新排序
decorations.sort((a, b) => a.from - b.from || a.to - b.to)
```

- 如果使用 `RangeSetBuilder` 或 `MatchDecorator`，通常不需要手动全量排序，或者构建过程本身就是有序的。

---

### ✅ 重构建议：使用 `MatchDecorator`

CodeMirror 官方提供了 `MatchDecorator` 类，专门用于解决“正则匹配 -> 生成装饰器”的场景，它内置了视口优化和增量更新，能完美替代你手写的 `createCapsuleDecorations` 逻辑。

#### 重构后的代码

```typescript
import { MatchDecorator, Decoration, ViewPlugin, EditorView, WidgetType } from '@codemirror/view'
import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget'
import type { CapsuleGlobalState } from '../types'

// 假设你的 matchCapsules 逻辑可以转换为正则。
// 如果逻辑非常复杂无法用正则，需要用 RangeSetBuilder + 视口迭代。
// 这里假设胶囊格式为 @name\u200B (推荐方案) 或其他正则可描述的格式
const CAPSULE_REGEX = /@([^\s]+?)\u200B/g // 示例正则

export const createCapsulePlugin = (globalState: CapsuleGlobalState) => {
  // 1. 定义匹配器 (自动处理视口和增量更新)
  const matcher = new MatchDecorator({
    regexp: CAPSULE_REGEX,
    decoration: (match, view, pos) => {
      const name = match[1] // 从正则捕获组获取名称

      // 2. 删除回调：动态获取位置
      const onDelete = () => {
        // 获取当前 Widget 在 DOM 中的位置
        // 注意：这里需要 Widget 内部在触发 onDelete 时传递 DOM 元素，或者我们不传 pos
        // 更简单的做法是：Widget 内部触发事件，这里只负责生成 Widget
      }

      return Decoration.replace({
        widget: new MentionCapsuleWidget(name, globalState), // 不传 from/to
        inclusive: false
      })
    }
  })

  return ViewPlugin.fromClass(
    class CapsuleViewPlugin {
      decorations: any

      constructor(view: EditorView) {
        this.decorations = matcher.createDeco(view)
      }

      update(update: ViewUpdate) {
        // IME 处理：MatchDecorator 内部通常处理得很好，
        // 但如果你想在 IME 期间完全冻结更新，可以保留这个逻辑
        if (globalState.composingRef.current) {
          // 仅做位置映射，不重新匹配
          this.decorations = this.decorations.map(update.changes)
          return
        }

        // 标准更新：MatchDecorator 会自动判断是否需要重新扫描
        this.decorations = matcher.updateDeco(update, this.decorations)
      }
    },
    {
      decorations: v => v.decorations,
      provide: plugin =>
        EditorView.atomicRanges.of(view => {
          return view.plugin(plugin)?.decorations || Decoration.none
        })
    }
  )
}
```

### 🔧 针对 Widget 删除逻辑的修正

如果你必须在 Widget 内部有一个“X”按钮来删除胶囊，请在 `MentionCapsuleWidget` 类中这样实现：

```typescript
// MentionCapsuleWidget.ts
class MentionCapsuleWidget extends WidgetType {
  constructor(readonly name: string, readonly globalState: any) {
    super()
  }

  eq(other: MentionCapsuleWidget) {
    return other.name === this.name
  }

  toDOM(view: EditorView) {
    const span = document.createElement('span')
    span.className = 'capsule'
    span.textContent = this.name

    const deleteBtn = document.createElement('span')
    deleteBtn.textContent = 'x'
    deleteBtn.onclick = e => {
      e.preventDefault()
      // ✅ 关键：通过 DOM 反查当前在文档中的位置
      const pos = view.posAtDOM(span)
      if (pos === null) return

      // 这里的 pos 是 Widget 的起始位置
      // 我们需要知道 Widget 覆盖了多长的文本。
      // 如果是 ReplaceDecoration，通常需要计算匹配文本的长度。
      // 更好的方式是利用 Atomic Range 特性，直接选中并删除。

      // 假设胶囊文本结构是固定的 (例如 @name\u200B)
      // 或者通过 view.state.doc 获取该位置附近的文本进行校验

      // 简单粗暴删除：
      // 假设我们知道胶囊对应的文本长度，或者通过正则重新匹配一下确认长度
      // 这里为了演示，假设长度是 name.length + 2 (@ + \u200B)
      const length = this.name.length + 2

      view.dispatch({
        changes: { from: pos, to: pos + length, insert: '' }
      })
    }

    span.appendChild(deleteBtn)
    return span
  }
}
```

### 总结

1.  **必须移除 `doc.toString()`**：这是性能杀手。
2.  **使用 `MatchDecorator`**：它是官方提供的、经过优化的文本替换工具。
3.  **修正位置获取方式**：不要在闭包里存 `from`，要在交互发生时通过 `view.posAtDOM` 获取实时位置。
4.  **利用 Atomic Ranges**：如果配置了 `atomicRanges`，用户按 Backspace 就能删除胶囊，可能根本不需要显式的删除按钮，这能简化大量逻辑。
