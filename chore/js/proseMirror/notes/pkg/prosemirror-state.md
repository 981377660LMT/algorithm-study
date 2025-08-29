好的，我们来对 `prosemirror-state` 进行最终的深度剖析。如果说 `prosemirror-model` 是名词（文档是什么），`prosemirror-transform` 是动词（如何改变文档），那么 `prosemirror-state` 就是**形容词和副词的集合**，它描述了**某一时刻编辑器的完整情状**，并规定了状态之间如何**平滑地过渡**。

它是连接数据模型与用户界面的中央枢纽。我们将从以下四个关键概念来彻底解构它：

1.  **`EditorState`：不可变的“世界快照”**
2.  **`Transaction`：状态演变的唯一途径**
3.  **`Selection`：用户的意图焦点**
4.  **`Plugin`：注入行为与状态的强大机制**

---

### 1. `EditorState`：不可变的“世界快照”

`EditorState` 是 ProseMirror 的核心概念。它是一个**不可变的（Immutable）** 对象，包含了在任何给定时间点，描述编辑器所需的一切信息。

#### a. `EditorState` 的构成

一个 `EditorState` 实例主要由以下几部分组成：

- **`doc: Node`**: 当前的文档。这是一个来自 `prosemirror-model` 的 `Node` 对象。
- **`selection: Selection`**: 当前用户的选区。它描述了光标位置或选中的范围。
- **`plugins: Plugin[]`**: 应用于此状态的所有插件的有序列表。
- **`storedMarks: readonly Mark[] | null`**: 当选区为空时，存储的“激活”标记。例如，你点击了加粗按钮但还未输入文本，`strong` 标记就会被存储在这里，你接下来输入的文本将自动应用这个标记。

#### b. 不可变性的核心价值

`EditorState` 的不可变性是 ProseMirror 设计哲学的基石，它带来了巨大的好处：

- **可预测的更新**: 状态更新的唯一方式是创建一个新的 `EditorState` 对象。这使得状态变化清晰、可追溯，杜绝了“幽灵更新”和副作用。
- **高效的集成**: 在与 React、Vue 等现代 UI 框架集成时，你只需比较 `oldState === newState` 即可知道状态是否改变，从而触发高效的视图重渲染。
- **可靠的历史记录**: 撤销操作仅仅是将当前状态指针指回一个旧的 `EditorState` 实例，简单、快速且绝对可靠。
- **协同编辑的基础**: 协同编辑需要在不同状态之间转换和合并变更，不可变的状态快照为此提供了稳定的基础。

```typescript
import { EditorState } from 'prosemirror-state'
import { schema } from 'prosemirror-schema-basic'

// 创建一个初始状态
const initialState = EditorState.create({
  schema // 必须提供 schema
  // doc: myInitialDoc, // 可以提供一个初始文档
  // plugins: [history(), keymap(...)], // 可以提供插件
})
```

---

### 2. `Transaction`：状态演变的唯一途径

既然 `EditorState` 是不可变的，我们如何改变它？答案是**事务（Transaction）**。`Transaction` 是从一个状态到下一个状态的描述。

#### a. `Transaction` vs `Transform`

`Transaction` **继承**自 `prosemirror-transform` 的 `Transform` 类。这意味着：

- 它拥有 `Transform` 的所有能力：可以累积 `Step` 来修改文档 (`tr.insertText(...)`, `tr.delete(...)`)。
- 它在此基础上增加了**管理编辑器状态**的能力。

#### b. `Transaction` 的额外职责

除了修改文档，一个 `Transaction` 还可以：

- **更新选区**: `tr.setSelection(newSelection)`。
- **管理激活标记**: `tr.setStoredMarks(...)`, `tr.ensureMarks(...)`。
- **附加元数据**: `tr.setMeta(pluginKey, data)`。这是插件之间以及插件与外部世界通信的关键机制。例如，历史记录插件通过元数据来区分一个事务是普通的编辑还是“撤销”操作。
- **控制滚动**: `tr.scrollIntoView()`。标记这个事务完成后，视图应该滚动以确保选区可见。

#### c. 状态更新的生命周期

ProseMirror 的核心数据流是一个简单而强大的循环：

1.  基于当前状态创建一个事务：`const tr = state.tr;`
2.  对事务进行一系列修改（添加步骤、设置选区等）。
3.  应用事务，生成一个新的状态：`const newState = state.apply(tr);`

```typescript
// state 是当前的 EditorState
function addHelloWorld(state: EditorState, dispatch: (tr: Transaction) => void) {
  // 1. 创建事务
  const tr = state.tr

  // 2. 修改事务
  tr.insertText('Hello World', state.selection.from)

  // 3. 通过 dispatch 函数将事务应用到视图
  dispatch(tr)
}
```

`dispatch` 函数通常由 `prosemirror-view` 提供，它负责接收事务、应用它来创建新状态，并更新视图。

---

### 3. `Selection`：用户的意图焦点

`Selection` 类描述了用户的选区，它也是状态的一部分。

- **`TextSelection`**: 最常见的类型，用于文本内容。它有两个关键属性：

  - `anchor`: 选区的“锚点”，即选区开始的地方（不随光标移动而改变）。
  - `head`: 选区的“头”，即光标当前所在的位置。
  - 当 `anchor === head` 时，它是一个折叠的选区（即光标）。

- **`NodeSelection`**: 用于选中一个单一的“块”节点，例如一张图片或一个视频。

- **`AllSelection`**: 用于选中整个文档。

选区的更新也必须通过事务来完成，确保了状态的一致性。

---

### 4. `Plugin`：注入行为与状态的强大机制

如果说 `EditorState` 是大脑，那么 `Plugin` 就是注入到大脑中的各种“思想”、“技能”和“反射弧”。几乎所有非核心功能都是通过插件实现的。

一个 `Plugin` 可以向编辑器系统注入两样东西：**自有状态**和**行为属性**。

#### a. 插件状态 (`PluginSpec.state`)

一个插件可以拥有自己独立的状态。这个状态与 `EditorState` 绑定，并随着每一个事务而更新。

- **`init()`**: 初始化插件的状态。
- **`apply(tr, value, oldState, newState)`**: 这是插件状态的核心。每当一个事务被应用时，这个函数就会被调用，它接收事务和旧的插件状态 (`value`)，并**必须返回一个新的插件状态**。这使得插件状态也能保持不可变和可预测。

**例子**：`prosemirror-history` 插件的状态就是两个数组：`undoStack` 和 `redoStack`。当一个事务发生时，它的 `apply` 方法会决定是将该事务的逆操作推入 `undoStack`，还是从栈中取出一个操作来执行。

#### b. 插件属性 (`PluginSpec.props`)

`props` 是插件向 `prosemirror-view` 注入行为的方式。一些最常用的 `props` 包括：

- `handleKeyDown(view, event)`: 拦截键盘事件，并可以将其转换为一个事务。这是实现快捷键的核心。
- `handleClick(view, pos, event)`: 处理鼠标点击。
- `decorations(state)`: 动态地向视图添加装饰器（如高亮、小部件）。这是实现搜索高亮、协同光标、评论标记等功能的关键。
- `nodeViews`: 为特定类型的节点提供自定义的渲染逻辑（例如，用一个 React 组件来渲染某个节点）。

```typescript
import { Plugin, PluginKey } from 'prosemirror-state'

const myPluginKey = new PluginKey('myPlugin')

const myPlugin = new Plugin({
  key: myPluginKey,
  state: {
    init() {
      return { clickCount: 0 }
    },
    apply(tr, oldState) {
      // 如果事务的元数据中有 'incrementClick'，则更新状态
      if (tr.getMeta(myPluginKey) === 'incrementClick') {
        return { clickCount: oldState.clickCount + 1 }
      }
      return oldState
    }
  },
  props: {
    handleClick(view, pos, event) {
      // 当点击时，分发一个带有元数据的事务来更新我们的插件状态
      const tr = view.state.tr.setMeta(myPluginKey, 'incrementClick')
      view.dispatch(tr)
      return true // 表示我们已经处理了这次点击
    }
  }
})
```

### 总结

`prosemirror-state` 是 ProseMirror 架构的粘合剂和控制器。它通过**不可变的 `EditorState`** 提供了稳定、可预测的状态管理；通过**事务驱动的更新机制**保证了状态演变的原子性和可追溯性；并通过**强大而灵活的 `Plugin` 系统**，将核心与功能解耦，赋予了编辑器无限的扩展能力。掌握了 `prosemirror-state`，就等于掌握了控制 ProseMirror 编辑器行为的“总开关”。
