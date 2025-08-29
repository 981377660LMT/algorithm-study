好的，我们来对 `prosemirror-commands` 进行一次彻底的、深入的剖析。这个包是 ProseMirror 生态中“用户意图”到“状态变更”的转换层。它提供了一套标准的、可组合的函数，用于执行常见的编辑操作，如切换样式、插入节点、删除内容等。

理解 `prosemirror-commands` 的设计，是构建编辑器 UI（如工具栏按钮、快捷键）和自定义编辑行为的关键。

我们将从以下四个关键角度来解构它：

1.  **核心概念：`Command` 函数的本质**
2.  **命令的分类与实践：一个功能导览**
3.  **组合的艺术：`chainCommands` 的威力**
4.  **创建你自己的命令：扩展编辑能力**

---

### 1. 核心概念：`Command` 函数的本质

在 ProseMirror 中，一个“命令”（Command）就是一个遵循特定签名的函数。它既是一个**查询**，也是一个**动作**。

#### a. 命令的函数签名

```typescript
type Command = (
  state: EditorState,
  dispatch?: (tr: Transaction) => void,
  view?: EditorView
) => boolean
```

让我们逐一分解这个签名：

- **`state: EditorState`**: 命令的**输入**。任何命令都必须基于当前的编辑器状态来做判断和操作。
- **`dispatch?: (tr: Transaction) => void`**: 命令的**执行器**。这是一个可选的回调函数。
  - 如果 `dispatch` **被提供了**，命令在执行其逻辑时，会创建一个 `Transaction`，并通过调用 `dispatch(tr)` 将其发送出去，从而改变编辑器的状态。
  - 如果 `dispatch` **未被提供**（为 `undefined`），命令则进入“**只读查询模式**”。它只检查在当前 `state` 下，该命令**是否可以被执行**，而不会产生任何实际效果。
- **`view?: EditorView`**: 可选的视图实例。大多数命令是纯粹基于状态的，但少数需要与视图交互的命令（如处理滚动）可能会用到它。
- **`=> boolean`**: 命令的**返回值**。
  - 返回 `true` 表示：命令**成功执行**（如果 `dispatch` 存在），或者**可以被执行**（如果 `dispatch` 不存在）。
  - 返回 `false` 表示：在当前状态下，该命令**无法执行**（例如，在一个空行上尝试“解除列表”）。

#### b. 命令的双重角色

这种设计非常巧妙，使得同一个命令函数可以用于两个不同的场景：

1.  **执行操作**:

    ```typescript
    // 假设有一个工具栏按钮
    boldButton.addEventListener('click', () => {
      // 调用命令，并传入 view.dispatch 来实际执行它
      toggleMark(schema.marks.strong)(view.state, view.dispatch)
    })
    ```

2.  **检查状态（更新 UI）**:
    ```typescript
    // 在每次编辑器状态更新时，检查按钮是否应该被激活
    function updateToolbar(state: EditorState) {
      // 调用命令，但不传入 dispatch，只检查其可用性
      const isBoldActive = toggleMark(schema.marks.strong)(state)
      boldButton.classList.toggle('active', isBoldActive)
    }
    ```
    _注意：`toggleMark` 比较特殊，它的检查模式也会检查标记是否已激活。_

---

### 2. 命令的分类与实践：一个功能导览

`prosemirror-commands` 提供了大量预设的命令，可以按功能分为几类：

#### a. 标记命令 (Mark Commands)

- **`toggleMark(markType, attrs?)`**: 这是最常用的命令之一。如果选区内没有此标记，则添加它；如果已有此标记，则移除它。非常适合用于实现加粗、斜体、下划线等按钮。

#### b. 节点包裹与提升命令 (Wrapping and Lifting)

- **`wrapIn(nodeType, attrs?)`**: 用指定的节点类型包裹当前选中的文本块。例如，`wrapIn(schema.nodes.blockquote)` 可以将一个或多个段落包裹成一个引用块。
- **`lift(target)`**: 将选中的文本块从其父容器中“提升”出来。例如，在一个列表项中使用 `lift` 可以将其变成一个普通的段落，脱离列表。
- **`setBlockType(nodeType, attrs?)`**: 将当前选中的文本块的类型改变为指定的节点类型。例如，`setBlockType(schema.nodes.heading, { level: 1 })` 可以将一个段落变成一级标题。

#### c. 删除与连接命令 (Deleting and Joining)

这些命令是实现 `Backspace` 和 `Delete` 键直观行为的基础。

- **`deleteSelection()`**: 删除当前选中的内容。
- **`joinBackward()`**: 在光标处向后连接。如果光标在一个文本块的开头，它会尝试将这个块与前一个块合并。例如，在段落开头按 `Backspace` 会删除换行符，使之与上一段合并。
- **`joinForward()`**: 在光标处向前连接。`Delete` 键在段落末尾的行为。
- **`selectNodeBackward` / `selectNodeForward`**: 当光标在“原子节点”（如图片）旁边时，按 `Backspace` 或 `Delete` 会选中该节点，而不是试图删除它的一部分。

#### d. 文本块创建命令

- **`splitBlock()`**: 在光标处分割当前的文本块。这是实现 `Enter` 键行为的基础。如果在一个段落中间按回车，段落会被分成两个。
- **`createParagraphNear()`**: 在当前选区附近创建一个新的空段落。

---

### 3. 组合的艺术：`chainCommands` 的威力

单个命令通常只处理一种特定情况。但用户的单个操作（如按 `Backspace` 键）可能需要应对多种情况。这时就需要 `chainCommands`。

**`chainCommands(...commands: Command[]): Command`**

`chainCommands` 接收一系列命令，并返回一个新的命令。当这个新命令被执行时，它会**按顺序**尝试执行传入的每一个命令，直到**第一个返回 `true` 的命令**为止。一旦有一个命令成功，链条就会立即停止。

**一个经典的例子：实现 `Backspace` 键的完整行为**

```typescript
import {
  chainCommands,
  deleteSelection,
  joinBackward,
  selectNodeBackward
} from 'prosemirror-commands'
import { keymap } from 'prosemirror-keymap'

const backspaceCommand = chainCommands(
  deleteSelection, // 1. 如果有选区，优先删除选区内容
  joinBackward, // 2. 如果没有选区，尝试向后合并块
  selectNodeBackward // 3. 如果前面是一个原子节点，则选中它
)

const myKeymapPlugin = keymap({
  Backspace: backspaceCommand
})
```

这个组合完美地模拟了用户对 `Backspace` 键的期望：删除选区 > 合并段落 > 选中节点。

---

### 4. 创建你自己的命令：扩展编辑能力

`prosemirror-commands` 的美妙之处在于，你自己的自定义命令与内置命令是一等公民，可以无缝地组合在一起。

**示例：创建一个命令，用于插入当前日期**

```typescript
import { EditorState, Transaction } from 'prosemirror-state'

const insertDate: Command = (state, dispatch) => {
  // 检查命令是否可以执行（这里我们总是可以插入日期）
  // 如果 dispatch 不存在，我们只做检查，直接返回 true
  if (!dispatch) {
    return true
  }

  // 创建一个事务
  const tr = state.tr

  // 获取当前日期字符串
  const date = new Date().toLocaleDateString()

  // 在当前光标位置插入文本
  // 我们不希望替换任何选区，所以只在光标处插入
  if (!tr.selection.empty) {
    tr.deleteSelection()
  }
  tr.insertText(date, tr.selection.from)

  // 分发事务
  dispatch(tr)

  // 返回 true 表示命令成功执行
  return true
}

// 现在你可以将它绑定到快捷键或按钮上
// keymap({ 'Mod-d': insertDate })
```

### 总结

`prosemirror-commands` 是一个优雅、可组合的命令层。它通过一个统一的 `Command` 函数签名，将**状态查询**和**动作执行**巧妙地结合在一起，使得构建交互式 UI 变得简单而可靠。其核心的 `chainCommands` 工具，更是将单一功能的命令串联成复杂、智能的用户行为的利器。掌握 `prosemirror-commands`，意味着你掌握了将用户意图转化为编辑器具体行为的标准方法，是从“能用”到“好用”的必经之路。
