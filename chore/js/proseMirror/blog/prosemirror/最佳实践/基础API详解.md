# ProseMirror 基础 API 深度解析

ProseMirror 的 API 庞大且精细，掌握核心 API 的细节是开发高质量编辑器的基础。本文将从数据模型、状态管理、视图交互三个维度进行盘点。

## 1. 数据模型 (Model)

### Schema (模式)

Schema 定义了文档允许包含哪些节点和标记。

- **`nodes`**: 定义节点类型（如 `doc`, `paragraph`, `text`）。
  - `content`: 内容表达式（如 `"inline*"`, `"block+"`, `"(paragraph | heading)+"`）。
  - `group`: 节点分组（如 `"block"`, `"inline"`），用于简化 `content` 定义。
  - `inline`: 是否为内联节点（默认为 `false`）。
  - `atom`: 是否为原子节点（内容不可编辑，光标不能进入，如图片）。
  - `attrs`: 定义节点属性（如 `{ src: { default: null } }`）。
  - `toDOM`: 定义如何渲染为 HTML。
  - `parseDOM`: 定义如何从 HTML 解析。
- **`marks`**: 定义标记类型（如 `strong`, `em`, `link`）。
  - `inclusive`: 光标在标记末尾时，输入内容是否继承该标记（默认 `true`，链接通常设为 `false`）。

### Node (节点)

文档树的基本构建块。

- **`type`**: 节点的类型 (`NodeType`)。
- **`attrs`**: 节点的属性对象。
- **`content`**: 节点的子内容 (`Fragment`)。
- **`marks`**: 应用于该节点的标记数组。
- **`text`**: 文本节点的文本内容。
- **核心方法**:
  - `node.childCount`: 子节点数量。
  - `node.child(index)`: 获取指定索引的子节点。
  - `node.textContent`: 获取纯文本内容。
  - `node.copy(content)`: 复制节点结构，但替换内容。
  - `node.slice(from, to)`: 切割节点，返回 `Slice`。
  - `node.rangeHasMark(from, to, type)`: 检查范围内是否有指定标记。

### Fragment (片段)

表示一组节点（通常是某个父节点的 `content`）。

- 类似于数组，但不可变且针对树结构优化。
- **核心方法**:
  - `fragment.size`: 内容的总长度（以 token 为单位）。
  - `fragment.forEach((node, offset, index) => ...)`: 遍历子节点。
  - `fragment.append(otherFragment)`: 连接片段。

### Slice (切片)

表示文档的一部分，通常用于复制/粘贴或拖拽。

- **`content`**: 片段内容 (`Fragment`)。
- **`openStart`**: 开始处的开放深度（即左侧有多少层父节点被切开了）。
- **`openEnd`**: 结束处的开放深度。
- _理解_: 复制一个列表项中的文字，`openStart` 和 `openEnd` 可能不为 0，因为你切断了 `<ul>` 和 `<li>`。

## 2. 状态管理 (State)

### EditorState (编辑器状态)

包含文档的完整状态，不可变对象。

- **`doc`**: 当前文档 (`Node`)。
- **`selection`**: 当前选区 (`Selection`)。
- **`storedMarks`**: 当前存储的标记（如按下 Ctrl+B 后，尚未输入文字时的加粗状态）。
- **`schema`**: 文档使用的 Schema。
- **`plugins`**: 激活的插件列表。
- **`tr`**: **获取一个新的 Transaction**（这是修改状态的唯一入口）。

### Transaction (事务)

继承自 `Transform`，用于描述状态变更。

- **修改文档**:
  - `tr.insertText(text, from, to)`: 插入文本。
  - `tr.replaceWith(from, to, node)`: 用节点替换范围。
  - `tr.delete(from, to)`: 删除范围。
  - `tr.insert(pos, node)`: 插入节点。
- **修改标记**:
  - `tr.addMark(from, to, mark)`: 添加标记。
  - `tr.removeMark(from, to, markType)`: 移除标记。
- **修改选区**:
  - `tr.setSelection(selection)`: 更新选区。
- **元数据**:
  - `tr.setMeta(key, value)`: 附加元数据（用于插件通信）。
  - `tr.time`: 事务发生的时间戳。
  - `tr.docChanged`: 是否修改了文档内容。

### Selection (选区)

- **`from` / `to`**: 选区的起始和结束位置。
- **`$from` / `$to`**: 解析后的位置对象 (`ResolvedPos`)，包含上下文信息。
- **`anchor` / `head`**: 锚点（不动点）和头部（动点）。
- **子类**:
  - **`TextSelection`**: 文本选区。
  - **`NodeSelection`**: 节点选区（选中图片、代码块等）。
  - **`AllSelection`**: 全选。

### ResolvedPos (解析位置)

通过 `doc.resolve(pos)` 获取，提供位置的上下文信息。

- **`parent`**: 直接父节点。
- **`depth`**: 深度（根节点为 0）。
- **`node(depth)`**: 获取指定深度的祖先节点。
- **`index(depth)`**: 在指定深度父节点中的索引。
- **`pos`**: 绝对位置。
- **`textOffset`**: 在文本节点内的偏移量。

## 3. 视图与交互 (View)

### EditorView (编辑器视图)

连接 DOM 和 State。

- **`state`**: 当前编辑器状态。
- **`dom`**: 编辑器的根 DOM 元素。
- **`updateState(state)`**: 更新视图以匹配新状态。
- **`dispatch(tr)`**: 派发事务（通常会调用 `updateState`）。
- **`posAtCoords({left, top})`**: 根据屏幕坐标获取文档位置。
- **`coordsAtPos(pos)`**: 根据文档位置获取屏幕坐标。
- **`domAtPos(pos)`**: 获取文档位置对应的 DOM 节点和偏移量。

### Props (属性)

传递给 `EditorView` 的配置对象，也可以通过插件定义。

- **`handleDOMEvents`**: 拦截原生 DOM 事件。
- **`handleKeyDown`**: 拦截键盘事件。
- **`handlePaste`**: 拦截粘贴事件。
- **`handleDrop`**: 拦截拖拽事件。
- **`nodeViews`**: 自定义节点渲染逻辑。
- **`decorations`**: 返回装饰器集合 (`DecorationSet`)。

## 4. 常用工具

### DOMParser

- `DOMParser.fromSchema(schema).parse(domNode)`: 将 DOM 转换为 ProseMirror 文档。

### DOMSerializer

- `DOMSerializer.fromSchema(schema).serializeFragment(fragment)`: 将 ProseMirror 片段转换为 DOM。

### Keymap

- `keymap({ "Mod-b": toggleMark(schema.marks.strong) })`: 绑定快捷键。

### InputRules

- `inputRules({ rules: [...] })`: 定义输入规则（如 Markdown 自动转换）。

## 5. 核心概念辨析：Index vs Position

- **Position (Pos)**: 全局绝对坐标。
  - 文档开始为 0。
  - 每个节点（非文本）占用 1 个单位（开始标签）+ 内容长度 + 1 个单位（结束标签）。
  - 文本节点占用其字符长度。
- **Index**: 相对索引。
  - 在父节点的 `content` 数组中的下标（第几个子节点）。
  - `node.child(index)` 使用的是 Index。

## 6. 最佳实践细节

1.  **不要手动修改 DOM**: 除非在 `NodeView` 的受控区域内，否则永远不要直接操作 `view.dom`。一切变更必须通过 `Transaction`。
2.  **理解 `tr.mapping`**: 当一个事务包含多个步骤时，后续步骤的位置需要通过 `mapping` 映射，以指向文档变化后的正确位置。
3.  **使用 `canPlaceholder`**: 在执行 `replace` 或 `insert` 前，可以使用 `tr.doc.content.matchAt(pos)` 检查 Schema 是否允许在此处插入该节点，避免产生无效文档。

---

```
# ProseMirror 基础 API 深度解析

ProseMirror 的 API 庞大且精细，掌握核心 API 的细节是开发高质量编辑器的基础。本文将从数据模型、状态管理、视图交互三个维度进行盘点。

## 1. 数据模型 (Model)

### Schema (模式)
Schema 定义了文档允许包含哪些节点和标记。
- **`nodes`**: 定义节点类型（如 `doc`, `paragraph`, `text`）。
  - `content`: 内容表达式（如 `"inline*"`, `"block+"`, `"(paragraph | heading)+"`）。
  - `group`: 节点分组（如 `"block"`, `"inline"`），用于简化 `content` 定义。
  - `inline`: 是否为内联节点（默认为 `false`）。
  - `atom`: 是否为原子节点（内容不可编辑，光标不能进入，如图片）。
  - `attrs`: 定义节点属性（如 `{ src: { default: null } }`）。
  - `toDOM`: 定义如何渲染为 HTML。
  - `parseDOM`: 定义如何从 HTML 解析。
- **`marks`**: 定义标记类型（如 `strong`, `em`, `link`）。
  - `inclusive`: 光标在标记末尾时，输入内容是否继承该标记（默认 `true`，链接通常设为 `false`）。

### Node (节点)
文档树的基本构建块。
- **`type`**: 节点的类型 (`NodeType`)。
- **`attrs`**: 节点的属性对象。
- **`content`**: 节点的子内容 (`Fragment`)。
- **`marks`**: 应用于该节点的标记数组。
- **`text`**: 文本节点的文本内容。
- **核心方法**:
  - `node.childCount`: 子节点数量。
  - `node.child(index)`: 获取指定索引的子节点。
  - `node.textContent`: 获取纯文本内容。
  - `node.copy(content)`: 复制节点结构，但替换内容。
  - `node.slice(from, to)`: 切割节点，返回 `Slice`。
  - `node.rangeHasMark(from, to, type)`: 检查范围内是否有指定标记。

### Fragment (片段)
表示一组节点（通常是某个父节点的 `content`）。
- 类似于数组，但不可变且针对树结构优化。
- **核心方法**:
  - `fragment.size`: 内容的总长度（以 token 为单位）。
  - `fragment.forEach((node, offset, index) => ...)`: 遍历子节点。
  - `fragment.append(otherFragment)`: 连接片段。

### Slice (切片)
表示文档的一部分，通常用于复制/粘贴或拖拽。
- **`content`**: 片段内容 (`Fragment`)。
- **`openStart`**: 开始处的开放深度（即左侧有多少层父节点被切开了）。
- **`openEnd`**: 结束处的开放深度。
- *理解*: 复制一个列表项中的文字，`openStart` 和 `openEnd` 可能不为 0，因为你切断了 `<ul>` 和 `<li>`。

## 2. 状态管理 (State)

### EditorState (编辑器状态)
包含文档的完整状态，不可变对象。
- **`doc`**: 当前文档 (`Node`)。
- **`selection`**: 当前选区 (`Selection`)。
- **`storedMarks`**: 当前存储的标记（如按下 Ctrl+B 后，尚未输入文字时的加粗状态）。
- **`schema`**: 文档使用的 Schema。
- **`plugins`**: 激活的插件列表。
- **`tr`**: **获取一个新的 Transaction**（这是修改状态的唯一入口）。

### Transaction (事务)
继承自 `Transform`，用于描述状态变更。
- **修改文档**:
  - `tr.insertText(text, from, to)`: 插入文本。
  - `tr.replaceWith(from, to, node)`: 用节点替换范围。
  - `tr.delete(from, to)`: 删除范围。
  - `tr.insert(pos, node)`: 插入节点。
- **修改标记**:
  - `tr.addMark(from, to, mark)`: 添加标记。
  - `tr.removeMark(from, to, markType)`: 移除标记。
- **修改选区**:
  - `tr.setSelection(selection)`: 更新选区。
- **元数据**:
  - `tr.setMeta(key, value)`: 附加元数据（用于插件通信）。
  - `tr.time`: 事务发生的时间戳。
  - `tr.docChanged`: 是否修改了文档内容。

### Selection (选区)
- **`from` / `to`**: 选区的起始和结束位置。
- **`$from` / `$to`**: 解析后的位置对象 (`ResolvedPos`)，包含上下文信息。
- **`anchor` / `head`**: 锚点（不动点）和头部（动点）。
- **子类**:
  - **`TextSelection`**: 文本选区。
  - **`NodeSelection`**: 节点选区（选中图片、代码块等）。
  - **`AllSelection`**: 全选。

### ResolvedPos (解析位置)
通过 `doc.resolve(pos)` 获取，提供位置的上下文信息。
- **`parent`**: 直接父节点。
- **`depth`**: 深度（根节点为 0）。
- **`node(depth)`**: 获取指定深度的祖先节点。
- **`index(depth)`**: 在指定深度父节点中的索引。
- **`pos`**: 绝对位置。
- **`textOffset`**: 在文本节点内的偏移量。

## 3. 视图与交互 (View)

### EditorView (编辑器视图)
连接 DOM 和 State。
- **`state`**: 当前编辑器状态。
- **`dom`**: 编辑器的根 DOM 元素。
- **`updateState(state)`**: 更新视图以匹配新状态。
- **`dispatch(tr)`**: 派发事务（通常会调用 `updateState`）。
- **`posAtCoords({left, top})`**: 根据屏幕坐标获取文档位置。
- **`coordsAtPos(pos)`**: 根据文档位置获取屏幕坐标。
- **`domAtPos(pos)`**: 获取文档位置对应的 DOM 节点和偏移量。

### Props (属性)
传递给 `EditorView` 的配置对象，也可以通过插件定义。
- **`handleDOMEvents`**: 拦截原生 DOM 事件。
- **`handleKeyDown`**: 拦截键盘事件。
- **`handlePaste`**: 拦截粘贴事件。
- **`handleDrop`**: 拦截拖拽事件。
- **`nodeViews`**: 自定义节点渲染逻辑。
- **`decorations`**: 返回装饰器集合 (`DecorationSet`)。

## 4. 常用工具

### DOMParser
- `DOMParser.fromSchema(schema).parse(domNode)`: 将 DOM 转换为 ProseMirror 文档。

### DOMSerializer
- `DOMSerializer.fromSchema(schema).serializeFragment(fragment)`: 将 ProseMirror 片段转换为 DOM。

### Keymap
- `keymap({ "Mod-b": toggleMark(schema.marks.strong) })`: 绑定快捷键。

### InputRules
- `inputRules({ rules: [...] })`: 定义输入规则（如 Markdown 自动转换）。

## 5. 核心概念辨析：Index vs Position

- **Position (Pos)**: 全局绝对坐标。
  - 文档开始为 0。
  - 每个节点（非文本）占用 1 个单位（开始标签）+ 内容长度 + 1 个单位（结束标签）。
  - 文本节点占用其字符长度。
- **Index**: 相对索引。
  - 在父节点的 `content` 数组中的下标（第几个子节点）。
  - `node.child(index)` 使用的是 Index。

## 6. 最佳实践细节

1.  **不要手动修改 DOM**: 除非在 `NodeView` 的受控区域内，否则永远不要直接操作 `view.dom`。一切变更必须通过 `Transaction`。
2.  **理解 `tr.mapping`**: 当一个事务包含多个步骤时，后续步骤的位置需要通过 `mapping` 映射，以指向文档变化后的正确位置。
3.  **使用 `canPlaceholder`**: 在执行 `replace` 或 `insert` 前，可以使用 `tr.doc.content.matchAt(pos)` 检查 Schema 是否允许在此处插入该节点，避免产生无效文档。

```

已进行更改。

我已经为你整理了 ProseMirror 的基础 API 详解，并保存到了 `基础API详解.md` 文件中。这份文档涵盖了数据模型、状态管理、视图交互等核心模块。

为了让你更直观地理解这些 API 如何协同工作，我补充一个**实战代码示例**，展示如何编写一个自定义命令（Command），以及一些文档中难以展开的**位置导航细节**。

### 1. 实战示例：编写一个“插入/切换引用块”的命令

这个例子展示了 `Schema`, `Selection`, `Transaction`, `Node` 是如何串联的。

```javascript
import { setBlockType } from 'prosemirror-commands'

// 命令函数的签名总是：(state, dispatch, view) => boolean
const toggleBlockquote = (state, dispatch) => {
  // 1. 从 Schema 中获取节点类型
  const blockquote = state.schema.nodes.blockquote
  if (!blockquote) return false // 如果当前 schema 不支持引用块，直接返回

  // 2. 检查当前选区是否已经是引用块
  // $from 是 ResolvedPos，提供了上下文信息
  const { $from } = state.selection
  // 检查当前深度的父节点是否是 blockquote
  const isBlockquote = $from.parent.type === blockquote

  // 3. 如果 dispatch 为空，说明只是询问“当前能否执行此命令”（用于更新 UI 按钮状态）
  if (!dispatch) return true

  // 4. 执行逻辑
  if (isBlockquote) {
    // 如果已经是引用块，则转换为普通段落 (Paragraph)
    // 这里直接复用官方提供的 setBlockType 命令
    return setBlockType(state.schema.nodes.paragraph)(state, dispatch)
  } else {
    // 如果不是，则转换为引用块
    // wrapIn (prosemirror-commands) 是更通用的包裹命令，但这里演示手动构建 Transaction

    // 创建一个 Transaction
    const tr = state.tr

    // 获取选区范围
    const { from, to } = state.selection

    // 使用 setBlockType 逻辑（简化版）：将选区内的块级节点类型修改为 blockquote
    tr.setBlockType(from, to, blockquote)

    // 5. 派发事务，触发视图更新
    dispatch(tr)
    return true
  }
}
```

### 2. 补充细节：ResolvedPos 的导航魔法

在开发中，你经常需要回答“光标所在的段落是从哪里开始的？”或者“当前列表项的父级列表在哪里？”。`ResolvedPos` (`$pos`) 提供了极其强大的导航能力：

假设文档结构如下：

```html
0:
<doc>
  1:
  <blockquote>
    2:
    <p>3: He|llo <-- 光标在 "He" 和 "llo" 之间 (pos = 5) 8:</p>
    9:
  </blockquote>
  10:
</doc>
```

此时 `$pos = state.doc.resolve(5)`：

- **`$pos.depth`**: `2` (根 -> blockquote -> p，深度为 2)
- **`$pos.parent`**: `<p>...</p>` 节点对象。
- **`$pos.node(1)`**: `<blockquote>...</blockquote>` 节点对象（第 1 层祖先）。
- **`$pos.start(1)`**: `2` (blockquote 的内容起始位置)。
- **`$pos.end(1)`**: `9` (blockquote 的内容结束位置)。
- **`$pos.before(1)`**: `1` (blockquote 标签的起始位置)。
- **`$pos.after(1)`**: `10` (blockquote 标签的结束位置)。
- **`$pos.index(1)`**: `0` (blockquote 是 doc 的第 0 个子节点)。

**应用场景**:
如果你想删除光标所在的整个引用块：

```javascript
const depth = $pos.depth // 假设我们找到了 blockquote 的深度是 1
const from = $pos.before(1)
const to = $pos.after(1)
dispatch(state.tr.delete(from, to))
```

### 3. 补充细节：Mapping (映射)

当你连续做两个操作时，第二个操作的位置必须基于第一个操作**之后**的文档。

```javascript
const tr = state.tr
// 假设文档是 "Hello World"
// 操作 1: 在开头插入 "Hi " -> "Hi Hello World"
tr.insertText('Hi ', 0)

// 操作 2: 想删除原文档中的 "World" (原位置 6-11)
// 错误做法: tr.delete(6, 11); // 这会删除 "ello "，因为位置变了！

// 正确做法 1: 手动计算偏移 (不推荐，容易算错)
// tr.delete(6 + 3, 11 + 3);

// 正确做法 2: 使用 Mapping
const map = tr.mapping
const newFrom = map.map(6) // 自动映射为 9
const newTo = map.map(11) // 自动映射为 14
tr.delete(newFrom, newTo)
```

ProseMirror 的 `Transaction` 会自动处理这种连续调用的映射，所以通常你直接链式调用 `tr.insertText(...).delete(...)` 是安全的。但在编写异步逻辑或处理协同编辑 Steps 时，理解 Mapping 至关重要。
