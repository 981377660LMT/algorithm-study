好的，我们来详细讲解 `tldraw` 中关于 **Collaboration (协作)** 的功能。这部分文档是构建实时多人白板应用的核心。

`tldraw` SDK 为实时协作提供了深度支持。你可以选择使用官方推荐的 `@tldraw/sync` 库快速实现，也可以使用其底层的 API 来集成任何你自己的或其他第三方的后端服务。

---

### **1. 推荐方式：使用 `@tldraw/sync`**

`@tldraw/sync` 是 `tldraw` 官方为实现快速、高效的多人协作而专门构建的库。它也是官方旗舰应用 [tldraw.com](https://tldraw.com) 背后的协作引擎。

#### **快速上手：`useSyncDemo` (用于原型开发)**

为了让你能轻松体验协作功能，`@tldraw/sync` 库提供了一个 `useSyncDemo` Hook，它会连接到一个由 `tldraw` 官方托管的**临时演示服务器**。

假设你有一个简单的本地应用：

```tsx
import { Tldraw } from 'tldraw'

function MyApp() {
  return <Tldraw />
}
```

要将其变为协作应用，只需三步：

**第一步：安装 `@tldraw/sync` 库**

```bash
npm install @tldraw/sync
```

**第二步：在代码中使用 `useSyncDemo`**
导入 `useSyncDemo` Hook，用一个**唯一的房间 ID (`roomId`)** 调用它，然后将返回的 `store` 传递给 `<Tldraw>` 组件。

```tsx
import { Tldraw } from 'tldraw'
import { useSyncDemo } from '@tldraw/sync'

function MyApp() {
  // 使用一个唯一的房间 ID 来创建一个支持同步的 store
  const store = useSyncDemo({ roomId: 'my-unique-room-id' })

  // 将这个 store 传递给 Tldraw 组件
  return <Tldraw store={store} />
}
```

**第三步：测试**
在你的浏览器中打开一个隐身窗口，或者使用另一台设备，访问同一个应用的 URL。只要 `roomId` 相同，这些应用就会进入一个共享的协作会话，你可以看到彼此的光标和实时编辑。

**演示版的局限性 (重要！)**
`useSyncDemo` 非常适合用于原型设计和功能体验，但**绝对不能用于生产环境**。

- 演示服务器上的数据最多只保留 24 小时。
- 房间是公开的，任何知道 `roomId` 的人都可以进入并编辑。

#### **在生产环境中使用 `@tldraw/sync`**

要在生产环境中使用 `@tldraw/sync`，你需要**自己托管 `tldraw sync server`**。官方文档中有专门的文章详细介绍如何部署。`tldraw` 官方目前不提供生产级别的托管服务。

---

### **2. 手动集成：使用其他后端**

`tldraw` 的 SDK 被设计为**后端无关**的，这意味着你可以将其与任何数据同步方案集成。例如，[Liveblocks](https://liveblocks.io/) 就提供了官方的 `tldraw` 集成方案，许多团队也已成功将 `tldraw` 接入了自己内部的实时数据系统。

手动集成主要涉及三个方面：

#### **a. 数据同步 (Synchronizing Data)**

这是协作的核心。你需要实现一套机制来分发和合并所有用户的编辑操作。

- **发送变更**: 使用 `editor.store.listen` 监听当前用户的操作，并将增量更新（`update` 对象）发送到你的后端（例如通过 WebSocket）。
- **接收变更**: 从后端接收到其他用户的变更后，使用 `editor.store.mergeRemoteChanges` 将这些变更应用到本地的 `store` 中。

这部分详细内容在 **Persistence (持久化)** 文档中有深入讲解。

#### **b. 用户状态 (User Presence)**

“用户状态”指的是在画布上显示其他协作者的信息，例如他们的**光标位置、姓名、颜色、当前选择的图形**等。

- 这些信息在 `store` 中以 `instance_presence` 记录的形式存在。
- `tldraw` 提供了一个 `createPresenceStateDerivation` 辅助函数，用于生成当前用户的状态信息。你需要将这个信息通过你的后端发送给其他所有协作者。
- 当编辑器收到一个 `userId` 与当前用户不同的 `instance_presence` 记录时，就会将其视为一个“协作者”，并在 UI 上显示其状态。
- 你可以通过 `editor.getCollaborators()` 和 `editor.getCollaboratorsOnCurrentPage()` 这两个响应式方法来获取协作者的信息。

#### **c. 协作相关的 UI 组件**

`tldraw` 库内置了一些专门为协作设计的 UI 组件。你可以通过 `<Tldraw>` 的 `components` prop 来移除、替换或自定义它们。

- **协作者光标**: `TLComponents.CollaboratorCursor`
- **光标聊天气泡**: `TLComponents.CursorChatBubble`
- **离线指示器**: `TLComponents.TopPanel`
- **协作者列表**: `TLComponents.SharePanel`
