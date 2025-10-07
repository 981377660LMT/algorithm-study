好的，我们来详细讲解 `tldraw` 中关于 **Persistence (持久化)** 的概念。这部分文档是实现画布内容保存、加载以及多人协作功能的关键。

---

### **什么是持久化 (Persistence)？**

在 `tldraw` 中，**持久化**指的是将编辑器当前的状态（包括所有图形、页面、相机位置等）存储到某个地方（如浏览器本地数据库或你的后端服务器），以便之后可以恢复。

这是实现以下功能的基础：

1.  **本地保存**：关闭浏览器标签页后，下次打开时内容依然存在。
2.  **云端同步/多人协作**：将画布的变更发送到服务器，并接收来自其他协作者的变更，实现实时同步。

---

### **1. 最简单的方式：`persistenceKey`**

如果你只需要实现**本地持久化**和**跨标签页同步**，最简单的方法就是使用 `persistenceKey` prop。

```tsx
import { Tldraw } from 'tldraw'
import 'tldraw/tldraw.css'

export default function () {
  return (
    <div style={{ position: 'fixed', inset: 0 }}>
      {/* 提供一个唯一的 key */}
      <Tldraw persistenceKey="my-project-document-1" />
    </div>
  )
}
```

**工作原理**：

- 当你提供了 `persistenceKey`，`tldraw` 会自动将画布的所有内容保存到浏览器的 **IndexedDB** 中。
- 更强大的是，如果你在另一个浏览器标签页中打开一个使用**相同 `persistenceKey`** 的 `<Tldraw>` 组件，它们的内容会自动实时同步。这是通过 `BroadcastChannel` API 实现的。

**注意**：同步的是**文档内容**（如图形、页面），而每个编辑器实例的状态（如当前选中的图形、光标位置）是独立的。

---

### **2. 手动控制：快照 (Snapshots)**

如果你需要将数据保存到你自己的后端服务器，或者进行更精细的控制，可以使用**快照 (Snapshot)**。快照是编辑器状态的一个 JSON 表示。

#### **获取快照 (`getSnapshot`)**

你可以使用 `getSnapshot` 函数来获取当前编辑器状态的 JSON 对象。

```tsx
import { useEditor, getSnapshot } from 'tldraw'

function SaveButton({ documentId, userId }) {
  const editor = useEditor()
  return (
    <button
      onClick={async () => {
        // 从 editor.store 获取快照
        const { document, session } = getSnapshot(editor.store)

        // 在多用户应用中，你应该分开存储 document 和 session
        // document: 画布内容，所有用户共享
        // session: 用户特定的状态（如当前页面、选择等），每个用户独立
        await saveDocumentStateToApi(documentId, document)
        await saveSessionStateToApi(documentId, userId, session)
      }}
    >
      保存
    </button>
  )
}
```

#### **加载快照 (`loadSnapshot`)**

使用 `loadSnapshot` 函数可以将之前保存的 JSON 对象加载回编辑器。

```tsx
import { useEditor, loadSnapshot } from 'tldraw'

function LoadButton({ documentId, userId }) {
  const editor = useEditor()
  return (
    <button
      onClick={async () => {
        const document = await loadDocumentStateFromApi(documentId)
        const session = await loadSessionStateFromApi(documentId, userId)

        // 注意：工具状态需要单独重置
        editor.setCurrentTool('select')
        // 将快照加载到 editor.store
        loadSnapshot(editor.store, { document, session })
      }}
    >
      加载
    </button>
  )
}
```

你也可以在组件初始化时通过 `snapshot` prop 来设置初始状态。

---

### **3. 高级用法：`store` Prop 和实时同步**

对于复杂的场景，比如多人协作，你需要更底层的控制。

#### **监听变更 (`store.listen`)**

你可以监听 `store` 的变化，以获取增量的更新数据。这对于将用户的每一步操作发送到后端服务器至关重要。

```javascript
// 监听用户对文档的修改
const unlisten = editor.store.listen(
  update => {
    // `update` 对象包含了 added, updated, removed 的记录
    // 将这个 update 对象发送到你的 WebSocket 服务器
    console.log('增量更新:', update)
    sendUpdateToBackend(update)
  },
  // scope: 'document' -> 只关心文档内容的变化
  // source: 'user' -> 只关心当前用户的操作（忽略来自远程的变更）
  { scope: 'document', source: 'user' }
)

// 当组件卸载时，记得取消监听
// unlisten()
```

#### **合并远程变更 (`store.mergeRemoteChanges`)**

当从服务器收到其他协作者的变更时，你需要将这些变更应用到当前用户的编辑器中。使用 `mergeRemoteChanges` 可以确保这些变更被标记为 `source: 'remote'`，从而不会再次触发 `store.listen` (如果你设置了 `source: 'user'`)，避免无限循环。

```javascript
// 假设这是从 WebSocket 收到的远程变更
myWebSocket.on('message', remoteUpdate => {
  editor.store.mergeRemoteChanges(() => {
    // 在这个回调函数中应用所有远程变更
    // tldraw 内部提供了 applyHistoryEntry 等方法来应用 update
    editor.store.applyHistoryEntry(remoteUpdate)
  })
})
```

---

### **4. 数据迁移 (Migrations)**

随着 `tldraw` 或你的应用的迭代，图形的数据结构可能会发生变化。例如，你可能想给一个自定义图形增加一个新的 `color` 属性。**迁移 (Migrations)** 机制就是用来处理这个问题的。

当 `loadSnapshot` 加载一个旧版本的快照时，它会自动运行所有必要的迁移脚本，将旧数据安全地更新到最新结构。

#### **为自定义图形添加迁移**

如果你创建了自定义图形，并且修改了它的 `props` 结构，你应该在对应的 `ShapeUtil` 类中定义迁移规则。

```typescript
import { createShapePropsMigrationSequence, createShapePropsMigrationIds, ShapeUtil } from 'tldraw'

// 1. 定义版本号
const Versions = createShapePropsMigrationIds('my-card-shape', {
  AddColor: 1
  // AddBorderRadius: 2, ...
})

class MyCardShapeUtil extends ShapeUtil<MyCardShape> {
  static override type = 'my-card-shape'

  // 2. 定义迁移序列
  static override migrations = createShapePropsMigrationSequence({
    sequence: [
      {
        id: Versions.AddColor,
        up(props) {
          // 这个函数会修改旧数据的 props
          // 这里为旧的卡片图形添加一个默认的 color 属性
          props.color = 'black'
        }
        // down 迁移（降级）不再需要
      }
    ]
  })

  // ... 其他 ShapeUtil 的方法
}
```

通过这种方式，即使你以后修改了 `MyCardShape` 的数据结构，之前保存的旧数据也能够被正确加载和显示。
