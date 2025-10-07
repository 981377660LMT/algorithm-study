好的，我们来详细讲解 `tldraw` 的官方实时协作解决方案：**`tldraw sync`**。

`tldraw sync` 是 `tldraw` 官方为实现快速、容错的共享文档同步而构建的库。它也是官方旗舰应用 [tldraw.com](https://tldraw.com) 在生产环境中使用的核心技术。

你可以使用它为你的 `tldraw` 应用添加实时的多人协作功能。

---

### **1. 如何部署 `tldraw sync`**

`tldraw sync` 提供了一个用于原型设计的**演示服务器**，但要在生产环境中使用，你**必须自己托管后端服务**。

主要有两种托管方式：

#### **方式一：使用 Cloudflare 模板 (推荐)**

这是最快、最推荐的入门方式。`tldraw` 官方提供了一个 Cloudflare 模板，可以一键部署一个生产级别的最小化后端系统。

- **技术栈**:

  - **Durable Objects**: 为每个协作房间提供一个唯一的、有状态的 WebSocket 服务器。这完美地解决了“一个房间只应有一个全局权威实例”的问题。
  - **R2**: 用于持久化存储文档快照，以及存储图片、视频等大型二进制资源。

- **你需要自己实现的功能**:
  - **身份验证和授权**: 决定谁可以访问哪个房间。
  - **速率和大小限制**: 防止滥用资源上传。
  - **历史快照存储**: 实现文档的长期版本历史。
  - **房间列表和搜索**: 让用户能找到和管理他们的协作空间。

**[点击这里开始使用 Cloudflare 模板](https://github.com/tldraw/tldraw-sync-template)**

#### **方式二：集成到你自己的后端**

虽然官方推荐 Cloudflare，但 `@tldraw/sync-core` 库可以集成到任何支持 WebSocket 的 JavaScript 服务器环境中（如 Node.js, Bun, Deno 等）。

官方提供了一个**[简单的服务器示例](https://github.com/tldraw/tldraw/tree/main/apps/sync-server-example)**，你可以参考它来了解如何将各个部分组合在一起。

---

### **2. `tldraw sync` 后端架构**

一个完整的 `tldraw sync` 后端由 2 到 3 个部分组成：

1.  **WebSocket 服务器**:

    - 为每个共享文档提供一个“房间 (Room)”。
    - 负责在所有客户端之间同步文档状态，并持久化存储。
    - 这是协作的核心。

2.  **资源存储提供者 (Asset Storage Provider)**:

    - 用于存储和检索大型二进制文件，如图片和视频。

3.  **(可选) 链接解析服务 (Unfurling Service)**:
    - 如果你想使用内置的书签 (bookmark) 图形，你需要一个服务来解析 URL 并提取其元数据（如标题、描述、缩略图）。

---

### **3. 前端与后端的整合**

在前端，你只需要使用 `@tldraw/sync` 包中的 `useSync` Hook。

下面是一个简化的客户端实现示例：

```tsx
import { Tldraw, TLAssetStore, Editor } from 'tldraw'
import { useSync } from '@tldraw/sync'
import { uploadFileAndReturnUrl } from './assets' // 你需要自己实现的资源上传逻辑
import { convertUrlToBookmarkAsset } from './unfurl' // 你需要自己实现的链接解析逻辑

// 1. 实现你自己的资源存储逻辑
const myAssetStore: TLAssetStore = {
  // 当用户上传文件时，这个函数会被调用
  async upload(file, asset) {
    // 将文件上传到你的服务器（如 S3），并返回永久 URL
    const url = await uploadFileAndReturnUrl(file)
    return url
  },
  // 当 tldraw 需要显示资源时，这个函数会被调用
  async resolve(asset) {
    // 直接返回资源记录中存储的 src
    return asset.props.src
  }
}

// 2. 实现链接解析逻辑 (可选)
function registerUrlHandler(editor: Editor) {
  editor.registerExternalAssetHandler('url', async ({ url }) => {
    return await convertUrlToBookmarkAsset(url)
  })
}

// 3. 在你的组件中使用 useSync
function MyEditorComponent({ myRoomId }: { myRoomId: string }) {
  // useSync Hook 会创建并管理与服务器的 WebSocket 连接
  const store = useSync({
    // 告诉 sync client 连接到哪个服务器和房间
    uri: `wss://my-custom-backend.com/connect/${myRoomId}`,
    // 告诉 sync client 如何处理资源的上传和下载
    assets: myAssetStore
  })

  // 将 store 传递给 Tldraw，并注册链接处理器
  return <Tldraw store={store} onMount={registerUrlHandler} />
}
```

在服务器端，核心是 `@tldraw/sync-core` 包导出的 `TLSocketRoom` 类。你需要为每个文档/房间创建一个 `TLSocketRoom` 实例。

---

### **4. 支持自定义图形和绑定**

这是一个**非常重要**的环节。为了让不同版本的客户端能够无缝协作，并确保数据在服务器端的有效性，你**必须**让客户端和服务器都知道你添加的任何自定义图形或绑定。

#### **在客户端**

将你的 `customShapeUtils` 和 `customBindingUtils` 传递给 `useSync` Hook。**注意**：与 `<Tldraw>` 组件不同，`useSync` 不会自动包含默认图形，所以你需要显式地将它们也加进去。

```tsx
import { Tldraw, defaultShapeUtils, defaultBindingUtils } from 'tldraw'
import { useSync } from '@tldraw/sync'
import { useMemo } from 'react'

function MyApp() {
  const store = useSync({
    uri: '...',
    assets: myAssetStore,
    // 合并自定义和默认的 shape utils
    shapeUtils: useMemo(() => [...customShapeUtils, ...defaultShapeUtils], []),
    // 合并自定义和默认的 binding utils
    bindingUtils: useMemo(() => [...customBindingUtils, ...defaultBindingUtils], [])
  })

  return <Tldraw store={store} shapeUtils={customShapeUtils} bindingUtils={customBindingUtils} />
}
```

#### **在服务器端**

在服务器上，你需要使用 `createTLSchema` 创建一个包含你自定义项的 `schema`，并将其传递给 `TLSocketRoom`。服务器端不关心完整的 `ShapeUtil`，只关心 `props` (用于验证) 和 `migrations` (用于数据迁移)。

```typescript
import { createTLSchema, defaultShapeSchemas, defaultBindingSchemas } from '@tldraw/tlschema'
import { TLSocketRoom } from '@tldraw/sync-core'

const schema = createTLSchema({
  shapes: {
    ...defaultShapeSchemas, // 包含所有默认图形的 schema

    // 添加你的自定义图形 schema
    myCustomShape: {
      props: myCustomShapeProps, // 用于服务器端验证
      migrations: myCustomShapeMigrations // 用于数据迁移
    }
  },
  bindings: {
    ...defaultBindingSchemas
    // ... 如果有自定义绑定，也在这里添加
  }
})

// 在你的服务器逻辑中，创建房间时传入 schema
const room = new TLSocketRoom({
  schema: schema
  // ... 其他配置
})
```

---

### **5. 部署注意事项**

- **版本匹配**: 你**必须**确保客户端使用的 `@tldraw` 版本与服务器端使用的 `@tldraw/sync-core` 版本完全匹配。
- **更新策略**: `tldraw` 不保证永久的向后兼容性。因此，更新时应遵循以下策略：**先部署并运行新的后端服务，然后再发布新的客户端**。这样可以确保新客户端连接时，后端已经准备就绪。
