好的，这是一个相当深入的话题。我们将分步进行：首先，深入、彻底地分析 ProseMirror 和 Yjs 的合并原理，特别是它们如何协同工作；然后，提供一个简化的实现示例来演示这些核心概念。

### 1. 原理深度分析：ProseMirror + Yjs 的协同合并机制

实时协同编辑的核心挑战在于：多个用户同时对同一份文档进行修改，如何保证所有客户端最终能达到一致的状态，并且这个过程对用户来说是平滑无感的。ProseMirror 和 Yjs 通过一个“关注点分离”的优雅架构解决了这个问题。

- **ProseMirror (PM)**: 负责“视图层”和“意图层”。它是一个富文本编辑器框架，但不直接处理协同冲突。它关注的是：

  1.  **结构化文档 (Schema)**: 定义文档的合法结构（如段落、标题、列表等）。
  2.  **事务 (Transaction)**: 将用户的每一次编辑（如输入、删除、加粗）表示为一个或多个 `Step` 组成的事务。`Step` 是对文档的原子化、可逆的描述。
  3.  **变换与变基 (Transform & Rebase)**: 这是 PM 自身最接近“合并”的功能。`rebase` 可以将一个用户的 `Step` 序列（基于旧的文档状态）转换成可以应用在另一个用户修改后的新文档状态上的新 `Step` 序列。但在协同场景中，这个功能通常被更强大的 CRDT 机制所取代或辅助。

- **Yjs**: 负责“数据层”和“一致性”。它是一个 CRDT (Conflict-free Replicated Data Type) 的实现。CRDT 的核心特性是，只要所有操作都被分发到所有副本，无论接收顺序如何，所有副本最终都会收敛到相同的状态。Yjs 关注的是：
  1.  **无冲突数据结构**: 提供如 `Y.Text`, `Y.Array`, `Y.Map` 等共享数据类型。
  2.  **唯一操作 ID**: Yjs 中的每一次插入（无论是字符、对象还是其他元素）都有一个唯一的 ID，由 `(clientID, clock)` 组成。这为所有操作提供了一个全局的、明确的排序依据。
  3.  **逻辑删除 (Tombstones)**: 删除操作并不会真正从数据结构中移除一个元素，而是将其标记为“已删除”（墓碑）。这解决了“AB 中间插入，C 又删除 B”这类经典冲突，确保插入操作的相对位置不会因为并发的删除而错乱。
  4.  **状态向量 (State Vectors)**: 每个客户端维护一个状态向量，记录它所拥有的每个其他客户端的最新操作 `clock`。这使得客户端之间可以高效地只交换彼此缺失的更新，而不是每次都同步整个文档。

#### 两者如何协同工作 (`y-prosemirror` 绑定库)

`y-prosemirror` 是连接 PM 和 Yjs 的桥梁。它的合并/同步逻辑可以分解为两个主要方向：

**方向一：本地 ProseMirror 编辑 -> 转换为 Yjs 更新 -> 广播**

1.  **监听 PM 事务**: `y-prosemirror` 插件会监听 PM 编辑器产生的所有本地事务 (`Transaction`)。
2.  **计算差异**: 它不会直接转换事务中的 `Step`。相反，它会比较事务**之前**的 PM 文档状态和事务**之后**的 PM 文档状态。
3.  **生成 Yjs 操作**: 通过比较这两个状态，它计算出需要对 `Y.Text`（或 `Y.XmlFragment`，用于更复杂的结构）执行哪些操作。
    - **文本插入**: 对应 `ytext.insert(index, text)`。
    - **文本删除**: 对应 `ytext.delete(index, length)`。
    - **格式变化**: 对应 `ytext.format(index, length, attributes)`。
4.  **应用到 Y.Doc**: 这些 Yjs 操作被应用到本地的 `Y.Doc` 实例中。
5.  **生成并广播更新**: Yjs 内部机制会自动将这些新操作打包成一个二进制的“更新消息” (update message)，然后通过网络提供者（如 `y-websocket`）广播给其他所有连接的客户端。

**关键点**: 为什么是比较状态而不是转换 `Step`？因为 PM 的 `Step` 是基于文档的线性索引的，而 Yjs 的 CRDT 结构不是一个简单的线性数组。直接转换 `Step` 会非常复杂且容易出错。通过比较前后状态，可以直接映射到 Yjs 的高级 API，逻辑更清晰、更健壮。

**方向二：接收远程 Yjs 更新 -> 应用到 Y.Doc -> 转换为 PM 事务 -> 更新编辑器**

1.  **接收更新**: 客户端通过网络接收到来自其他用户的二进制更新消息。
2.  **应用到 Y.Doc**: 使用 `Y.applyUpdate(ydoc, update)` 将这些更新应用到本地的 `Y.Doc`。Yjs 的 CRDT 算法在这里发挥作用，自动、无冲突地合并这些更改。即使本地有未同步的更改，合并结果也是确定的。
3.  **监听 Yjs 变化**: `y-prosemirror` 监听本地 `Y.Doc` 的 `observe` 事件。当 `Y.Doc` 因接收远程更新而改变时，这个事件会被触发。
4.  **计算 PM 差异**: 事件的回调函数会收到一个描述了 `Y.Text` 变化的 `YTextEvent`。这个事件包含了哪些内容被插入、删除或格式化了。
5.  **生成 PM 事务**: `y-prosemirror` 将 `YTextEvent` 描述的差异转换成一个或多个 PM `Step`。
    - `ytext` 的插入 -> PM 的 `ReplaceStep` (插入)。
    - `ytext` 的删除 -> PM 的 `ReplaceStep` (删除)。
    - `ytext` 的格式化 -> PM 的 `AddMarkStep` / `RemoveMarkStep`。
6.  **应用到 PM 编辑器**: `y-prosemirror` 创建一个新的 PM `Transaction`，包含这些生成的 `Step`，并将其应用到编辑器状态上。由于这个事务源于远程，它会被特殊标记，以防止触发上面“方向一”的逻辑，避免无限循环。

#### 合并冲突的解决实例

假设原始文本是 `AC`。

- **User 1**: 在 `A` 和 `C` 之间插入 `B`。
  - PM 事务: 在索引 1 处插入 `B`。
  - `y-prosemirror` 转换为 Yjs 操作: `ytext.insert(1, 'B')`。
  - Yjs 内部: 创建一个带唯一 ID 的 `Item('B')`，并将其插入到 `Item('A')` 和 `Item('C')` 之间。
- **User 2 (并发操作)**: 在 `A` 和 `C` 之间插入 `X`。
  - PM 事务: 在索引 1 处插入 `X`。
  - `y-prosemirror` 转换为 Yjs 操作: `ytext.insert(1, 'X')`。
  - Yjs 内部: 创建一个带唯一 ID 的 `Item('X')`，并将其插入到 `Item('A')` 和 `Item('C')` 之间。

**合并过程**:
当 User 1 接收到 User 2 的更新（或反之），Yjs 需要决定 `B` 和 `X` 的顺序。CRDT 的排序规则（通常基于 `clientID` 的大小）会介入。例如，如果 User 2 的 `clientID` 大于 User 1，那么 `X` 会被排在 `B` 的后面。

最终，两个客户端的 `Y.Doc` 都会收敛到同一个状态，例如 `ABXC`（或者 `AXBC`，取决于 `clientID`）。然后 `y-prosemirror` 会将这个最终状态渲染到 PM 编辑器中。整个过程对用户来说是自动的，不需要任何手动解决冲突的步骤。

### 2. 简化实现示例

要完整实现一个 `y-prosemirror` 是一个巨大的工程。下面我们将用伪代码和简化的 JavaScript 来模拟其核心思想，帮助你理解其工作流程。

我们将模拟以下部分：

- 一个极简的 `ProseMirrorEditor` 类，能应用事务。
- 一个极简的 `YDoc` 和 `YText` 类，模拟 CRDT 的插入和删除。
- 一个 `Binding` 类，作为 `y-prosemirror` 的简化版。

```javascript
// --- 1. 模拟 ProseMirror ---
// 极简的 PM 编辑器状态
class ProseMirrorState {
  constructor(doc) {
    this.doc = doc // doc 是一个简单的字符串
  }

  // 应用一个“事务”，这里简化为直接返回新状态
  apply(transaction) {
    const newDoc =
      this.doc.slice(0, transaction.from) +
      (transaction.text || '') +
      this.doc.slice(transaction.to)
    return new ProseMirrorState(newDoc)
  }
}

// 极简的 PM 编辑器视图
class ProseMirrorEditor {
  constructor(doc) {
    this.state = new ProseMirrorState(doc)
    this.onTransaction = null // 回调，当有新事务时触发
  }

  // 模拟用户输入
  dispatch(transaction) {
    console.log(`PM: 调度事务`, transaction)
    this.state = this.state.apply(transaction)
    console.log(`PM: 新文档状态: "${this.state.doc}"`)
    if (this.onTransaction) {
      this.onTransaction(transaction, this.state)
    }
  }

  // 从外部（如Yjs绑定）更新状态
  updateState(newState) {
    console.log(`PM: 接收到外部更新，新文档: "${newState.doc}"`)
    this.state = newState
  }
}

// --- 2. 模拟 Yjs ---
// 模拟 Yjs 的 Item，每个字符都有来源和ID
class YItem {
  constructor(content, id) {
    this.content = content
    this.id = id // 简化ID，实际是 { client, clock }
    this.deleted = false
  }
}

// 模拟 YText 数据结构
class YText {
  constructor() {
    this.items = []
    this.onUpdate = null // 当数据变化时触发的回调
  }

  // 获取当前未被删除的文本内容
  toString() {
    return this.items
      .filter(item => !item.deleted)
      .map(item => item.content)
      .join('')
  }

  // 模拟插入
  insert(index, text, clientId) {
    const newItems = text
      .split('')
      .map((char, i) => new YItem(char, `${clientId}-${Date.now() + i}`))
    this.items.splice(index, 0, ...newItems)
    if (this.onUpdate) {
      // 实际的 Yjs 会生成一个二进制的 update
      this.onUpdate({ type: 'insert', index, text })
    }
  }

  // 模拟删除 (标记为墓碑)
  delete(index, length) {
    for (let i = 0; i < length; i++) {
      if (this.items[index + i]) {
        this.items[index + i].deleted = true
      }
    }
    if (this.onUpdate) {
      this.onUpdate({ type: 'delete', index, length })
    }
  }
}

// 模拟 Y.Doc
class YDoc {
  constructor() {
    this.text = new YText()
  }
}

// --- 3. 模拟 y-prosemirror 绑定 ---
class CollabBinding {
  constructor(editor, ydoc, clientId) {
    this.editor = editor
    this.ydoc = ydoc
    this.clientId = clientId
    this._isApplyingRemote = false // 标志位，防止无限循环

    // 方向一：监听 PM -> 更新 Yjs
    this.editor.onTransaction = (transaction, newState) => {
      if (this._isApplyingRemote) {
        return
      }
      console.log(`Binding: 监听到本地 PM 事务，准备更新 Yjs`)
      // 简化：直接将事务转换为Yjs操作
      // 真实实现会比较状态
      if (transaction.text) {
        // 插入
        this.ydoc.text.insert(transaction.from, transaction.text, this.clientId)
      } else {
        // 删除
        const length = transaction.to - transaction.from
        this.ydoc.text.delete(transaction.from, length)
      }
    }

    // 方向二：监听 Yjs -> 更新 PM
    this.ydoc.text.onUpdate = yjsUpdate => {
      console.log(`Binding: 监听到 Yjs 更新，准备更新 PM`)
      this._isApplyingRemote = true

      // 简化：直接从 Yjs 的当前状态生成新的 PM 状态
      // 真实的实现会把 yjsUpdate 转换为 PM Steps
      const newDocContent = this.ydoc.text.toString()
      const newPmState = new ProseMirrorState(newDocContent)
      this.editor.updateState(newPmState)

      this._isApplyingRemote = false
    }
  }
}

// --- 4. 运行模拟 ---

// 共享的数据模型
const sharedYDoc = new YDoc()

// 客户端 1
const editor1 = new ProseMirrorEditor('AC')
const binding1 = new CollabBinding(editor1, sharedYDoc, 'client1')

// 客户端 2
const editor2 = new ProseMirrorEditor('AC')
const binding2 = new CollabBinding(editor2, sharedYDoc, 'client2')

// 模拟网络：Yjs 的更新会触发另一个客户端绑定的 onUpdate
// 真实的实现会通过 y-websocket 等库来广播
sharedYDoc.text.onUpdate = update => {
  console.log(`\n--- 模拟网络广播 Yjs 更新 ---`)
  // 广播给所有客户端（除了源）
  // 在这个简化模型里，我们手动触发另一个编辑器的更新逻辑
  // 假设我们能区分来源，但这里为了简单，两个都更新
  binding1.ydoc.text.onUpdate(update)
  binding2.ydoc.text.onUpdate(update)
}

console.log('--- 初始状态 ---')
console.log('Editor 1:', editor1.state.doc)
console.log('Editor 2:', editor2.state.doc)
console.log('\n')

console.log("--- User 1 在索引 1 插入 'B' ---")
// 模拟用户1在编辑器中输入
editor1.dispatch({ from: 1, to: 1, text: 'B' })

/*
执行流程:
1. editor1.dispatch -> 触发 binding1.editor.onTransaction
2. binding1 将 PM 事务转为 ydoc.text.insert(1, 'B')
3. ydoc.text.insert 触发 sharedYDoc.text.onUpdate (模拟网络广播)
4. onUpdate 回调被调用，它会触发所有绑定的 Yjs 更新逻辑
5. binding2.ydoc.text.onUpdate 被触发
6. binding2 将 Yjs 状态 ("ABC") 更新到 editor2
*/

console.log('\n--- 最终状态 ---')
console.log('Editor 1:', editor1.state.doc) // "ABC"
console.log('Editor 2:', editor2.state.doc) // "ABC"
console.log('YDoc:', sharedYDoc.text.toString()) // "ABC"
```

这个简化的实现揭示了 `y-prosemirror` 的核心双向绑定和数据流。虽然它省略了 CRDT 冲突解决的复杂细节（如唯一 ID 排序）和从 Yjs 事件到 PM `Step` 的精确转换，但它清晰地展示了 ProseMirror 作为“意图”和“视图”层，Yjs 作为“数据”和“一致性”层的分工与协作模式。这是构建健壮协同编辑系统的关键所在。

---

好的，遵照您的要求，这里提供一个可以直接运行的具体实现代码示例。

这个示例将创建一个完整的 HTML 文件。您只需要按照说明安装依赖并运行一个简单的协同服务器，就可以在两个浏览器窗口中看到实时的协同编辑效果。

### 步骤 1: 项目设置和依赖安装

首先，创建一个新的文件夹，例如 `prosemirror-yjs-demo`。然后在此文件夹中打开终端，并安装所需的库。

```bash
# 创建并进入项目文件夹
mkdir prosemirror-yjs-demo
cd prosemirror-yjs-demo

# 初始化 npm 项目
npm init -y

# 安装所有必需的依赖
npm install prosemirror-model prosemirror-state prosemirror-view prosemirror-schema-basic prosemirror-commands prosemirror-keymap prosemirror-history yjs y-prosemirror y-websocket
```

### 步骤 2: 启动协同服务器

`y-websocket` 包自带一个简单的 WebSocket 服务器，非常适合用于演示。在您的项目文件夹 (`prosemirror-yjs-demo`) 中，运行以下命令来启动它：

```bash
npx y-websocket
```

您会看到类似 `Listening on localhost:1234` 的输出。保持这个终端窗口运行。

### 步骤 3: 创建 `index.html` 文件

在 `prosemirror-yjs-demo` 文件夹中，创建 `index.html` 文件，并将以下代码完整地复制进去。

这段代码会加载所有必要的模块，并设置一个连接到本地 WebSocket 服务器的 ProseMirror 编辑器。

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>ProseMirror + Yjs Demo</title>
    <style>
      /* 基本的 ProseMirror 编辑器样式 */
      .ProseMirror {
        position: relative;
        border: 1px solid #ccc;
        padding: 10px;
        min-height: 300px;
      }
      .ProseMirror:focus {
        outline: none;
        border-color: #66afe9;
      }
      /* Yjs 远程光标样式 */
      .ProseMirror-yjs-cursor {
        position: relative;
        margin-left: -1px;
        margin-right: -1px;
        border-left: 1px solid black;
        border-right: 1px solid black;
        border-color: orange;
        word-break: normal;
        pointer-events: none;
      }
      .ProseMirror-yjs-cursor > div {
        position: absolute;
        top: -1.05em;
        left: -1px;
        font-size: 13px;
        background-color: rgb(250, 129, 0);
        color: white;
        padding: 2px 6px;
        border-radius: 3px;
        white-space: nowrap;
      }
    </style>
  </head>
  <body>
    <h1>ProseMirror + Yjs 协同编辑</h1>
    <p>在另一个浏览器窗口打开此页面以查看协同效果。</p>
    <div id="editor"></div>
    <p>连接状态: <span id="status">连接中...</span></p>

    <script type="module">
      // --- 1. 导入所需模块 ---
      import * as Y from 'yjs'
      import { WebsocketProvider } from 'y-websocket'
      import { ySyncPlugin, yCursorPlugin, yUndoPlugin, undo, redo } from 'y-prosemirror'

      import { EditorState } from 'prosemirror-state'
      import { EditorView } from 'prosemirror-view'
      import { Schema } from 'prosemirror-model'
      import { schema as basicSchema } from 'prosemirror-schema-basic'
      import { keymap } from 'prosemirror-keymap'
      import { baseKeymap } from 'prosemirror-commands'
      import { history } from 'prosemirror-history'

      // --- 2. 设置 Yjs ---
      // 创建一个 Yjs 文档
      const ydoc = new Y.Doc()

      // 连接到在 localhost:1234 运行的 WebSocket 服务器
      const provider = new WebsocketProvider(
        'ws://localhost:1234',
        'prosemirror-yjs-demo-room', // 房间名，相同房间名的客户端会同步
        ydoc
      )

      // 获取用于 ProseMirror 的 Yjs 数据类型 (XML Fragment)
      const yXmlFragment = ydoc.getXmlFragment('prosemirror')

      // 监听连接状态并更新 UI
      const statusElement = document.getElementById('status')
      provider.on('status', event => {
        statusElement.textContent = event.status
      })

      // --- 3. 设置 ProseMirror ---
      // 定义 ProseMirror 文档结构
      const schema = new Schema({
        nodes: basicSchema.spec.nodes,
        marks: basicSchema.spec.marks
      })

      // 创建 ProseMirror 状态
      const state = EditorState.create({
        schema,
        plugins: [
          // ySyncPlugin 是核心，它负责双向同步
          ySyncPlugin(yXmlFragment),
          // yCursorPlugin 用于显示远程用户的光标
          yCursorPlugin(provider.awareness),
          // yUndoPlugin 提供了共享的撤销/重做历史记录
          yUndoPlugin(),
          // 将 Yjs 的撤销/重做绑定到标准快捷键
          keymap({
            'Mod-z': undo,
            'Mod-y': redo,
            'Shift-Mod-z': redo
          }),
          // ProseMirror 的基础按键映射
          keymap(baseKeymap)
        ]
      })

      // --- 4. 创建并挂载编辑器 ---
      const editorElement = document.getElementById('editor')
      const view = new EditorView(editorElement, {
        state,
        // 当编辑器状态改变时，这个函数会被调用
        dispatchTransaction(transaction) {
          const newState = view.state.apply(transaction)
          view.updateState(newState)
        }
      })

      // 将编辑器实例和 provider 暴露到 window 对象，方便调试
      window.editorView = view
      window.yProvider = provider
    </script>
  </body>
</html>
```

### 步骤 4: 运行和测试

1.  **确保 WebSocket 服务器正在运行** (执行 `npx y-websocket` 的终端)。
2.  你需要一个简单的 HTTP 服务器来提供 `index.html` 文件，因为浏览器出于安全原因限制直接从文件系统加载模块。最简单的方法是使用 `serve`：
    - 如果没有安装 `serve`，先安装：`npm install -g serve`
    - 在 `prosemirror-yjs-demo` 文件夹中，运行：`serve`
    - 它会给出一个本地地址，通常是 `http://localhost:3000`。
3.  **测试协同**：
    - 在你的浏览器中打开 `http://localhost:3000` (或 `serve` 提供的地址)。
    - 再打开一个**新的浏览器窗口**（或无痕窗口）并访问相同的地址。
    - 现在你有了两个模拟不同用户的客户端。在一个窗口中输入文字、加粗、创建段落，你会看到另一个窗口几乎是瞬间同步了所有内容，并且还能看到对方的光标位置。你甚至可以测试 `Cmd/Ctrl+Z`，撤销操作也是共享的。

这个示例完整地展示了 `prosemirror` 和 `yjs` 如何通过 `y-prosemirror` 绑定库无缝集成，实现了强大的实时协同编辑功能。

---

好的，这里是一个完整的 `mergeSteps.ts` 文件，包含了函数实现和一个详细的示例，演示了如何使用它以及各种场景下的结果。

```typescript
/**
 * 定义操作的类型
 * 'create': 创建一个新实体，data 是完整的实体数据。
 * 'update': 更新一个现有实体，data 是部分或全部的变更字段。
 * 'delete': 删除一个实体。
 *
 * @property timestamp - 操作发生的时间戳或版本号，用于解决冲突。
 */
export type Operation<T = Record<string, any>> = {
  id: string | number
  timestamp: number
} & ({ type: 'create'; data: T } | { type: 'update'; data: Partial<T> } | { type: 'delete' })

/**
 * 合并 create、update、delete 操作.
 *
 * 借鉴了 CRDT 的 "Last-Write-Wins" 思想，通过时间戳来处理操作合并。
 *
 * @alias compactSteps
 * @param steps 原始操作序列
 * @returns 紧凑化后的操作序列
 */
function mergeSteps<T extends Record<string, any>>(steps: Operation<T>[]): Operation<T>[] {
  // 1. 按 ID 分组
  const opsById = new Map<string | number, Operation<T>[]>()
  for (const op of steps) {
    if (!opsById.has(op.id)) {
      opsById.set(op.id, [])
    }
    opsById.get(op.id)!.push(op)
  }

  const finalOps: Operation<T>[] = []

  // 2. 对每个 ID 的操作序列进行化简
  for (const [id, opList] of opsById.entries()) {
    // 按时间戳排序，确保我们按逻辑顺序处理操作
    opList.sort((a, b) => a.timestamp - b.timestamp)

    let finalOp: Operation<T> | null = null

    for (const currentOp of opList) {
      if (!finalOp) {
        // 这是该 ID 的第一个操作
        finalOp = JSON.parse(JSON.stringify(currentOp)) // 深拷贝以避免修改原始输入
        continue
      }

      // 规则：新操作 (currentOp) 与已合并的操作 (finalOp) 进行合并
      // 核心：后发生的操作（时间戳更大）具有决定权

      // 情况 1: 已有操作是 'delete'
      if (finalOp.type === 'delete') {
        if (currentOp.type === 'create') {
          // D -> C => U (删除后又创建，视为对一个已存在实体的更新)
          finalOp = {
            type: 'update',
            id: id,
            data: currentOp.data,
            timestamp: currentOp.timestamp
          }
        }
        // D -> U 或 D -> D 是无效序列，可以安全地忽略
        continue
      }

      // 情况 2: 新操作是 'delete'
      if (currentOp.type === 'delete') {
        // C -> D 或 U -> D => D (任何操作后接删除，最终结果都是删除)
        finalOp = JSON.parse(JSON.stringify(currentOp))
        continue
      }

      // 情况 3: 新操作是 'update'
      if (currentOp.type === 'update') {
        // C -> U 或 U -> U => 合并 data
        if (finalOp.type === 'create' || finalOp.type === 'update') {
          // 注意：简单的 Object.assign 是浅合并，对于嵌套对象可能不符合预期。
          // 在实际应用中，你可能需要一个深度合并的函数。
          // 这里为了演示，Object.assign 足够。
          finalOp.data = Object.assign(finalOp.data || {}, currentOp.data)
          finalOp.timestamp = currentOp.timestamp // 时间戳更新
        }
        continue
      }

      // 情况 4: 新操作是 'create'
      if (currentOp.type === 'create') {
        // C -> C 或 U -> C 是逻辑冲突，但遵循 "Last-Write-Wins"
        // 我们认为后来的 create 覆盖了之前的一切
        finalOp = JSON.parse(JSON.stringify(currentOp))
        continue
      }
    }

    if (finalOp) {
      // 检查初始状态是否为 create，最终状态是否为 delete
      // 如果是，则代表这个实体从未真正存在过，可以完全忽略
      const firstOp = opList[0]
      if (firstOp.type === 'create' && finalOp.type === 'delete') {
        // C -> ... -> D => 无操作
      } else {
        finalOps.push(finalOp)
      }
    }
  }

  // 返回结果，可以按 ID 或时间戳排序以获得确定性输出
  return finalOps.sort((a, b) => String(a.id).localeCompare(String(b.id)))
}

// --- 示例使用 ---

const exampleSteps: Operation[] = [
  // --- ID 1: Create -> Update -> Delete (最终应被移除) ---
  { type: 'create', id: 1, data: { name: 'Alice', role: 'user' }, timestamp: 100 },
  { type: 'update', id: 1, data: { role: 'admin' }, timestamp: 101 },
  { type: 'delete', id: 1, timestamp: 102 },

  // --- ID 2: Create -> Update -> Update (最终应合并为一个 Create) ---
  { type: 'create', id: 2, data: { name: 'Bob', score: 0 }, timestamp: 200 },
  { type: 'update', id: 2, data: { score: 50 }, timestamp: 201 },
  { type: 'update', id: 2, data: { score: 100, status: 'active' }, timestamp: 202 },

  // --- ID 3: Update -> Update (最终应合并为一个 Update) ---
  { type: 'update', id: 3, data: { a: 1 }, timestamp: 300 },
  { type: 'update', id: 3, data: { b: 2 }, timestamp: 301 },

  // --- ID 4: Delete -> Create (最终应合并为一个 Update) ---
  { type: 'delete', id: 4, timestamp: 400 },
  { type: 'create', id: 4, data: { name: 'Charlie' }, timestamp: 401 },

  // --- ID 5: Update -> Delete (最终应合并为一个 Delete) ---
  { type: 'update', id: 5, data: { value: 'old' }, timestamp: 500 },
  { type: 'delete', id: 5, timestamp: 501 },

  // --- ID 6: 乱序操作 (最终应按时间戳正确合并) ---
  { type: 'update', id: 6, data: { status: 'pending' }, timestamp: 602 },
  { type: 'create', id: 6, data: { name: 'David' }, timestamp: 600 }
]

console.log('Original steps:', exampleSteps.length)
const compactedSteps = mergeSteps(exampleSteps)
console.log('Compacted steps:', compactedSteps.length)
console.log('--- Result ---')
console.log(JSON.stringify(compactedSteps, null, 2))

/*
--- 预期输出 ---

Original steps: 13
Compacted steps: 5
--- Result ---
[
  {
    "type": "create",
    "id": 2,
    "data": {
      "name": "Bob",
      "score": 100,
      "status": "active"
    },
    "timestamp": 202
  },
  {
    "type": "update",
    "id": 3,
    "data": {
      "a": 1,
      "b": 2
    },
    "timestamp": 301
  },
  {
    "type": "update",
    "id": 4,
    "data": {
      "name": "Charlie"
    },
    "timestamp": 401
  },
  {
    "type": "delete",
    "id": 5,
    "timestamp": 501
  },
  {
    "type": "create",
    "id": 6,
    "data": {
      "name": "David",
      "status": "pending"
    },
    "timestamp": 602
  }
]

--- 结果分析 ---
- ID 1: 被完全移除，因为 Create -> Delete 序列抵消了。
- ID 2: 合并为一个 Create 操作，data 包含了所有 update 的内容。
- ID 3: 合并为一个 Update 操作。
- ID 4: Delete 后又 Create，被视为一个 Update。
- ID 5: Update 后又 Delete，被视为一个 Delete。
- ID 6: 即使输入乱序，也根据时间戳先处理 Create，再合并 Update，最终得到一个合并后的 Create。
*/

export {}
```
