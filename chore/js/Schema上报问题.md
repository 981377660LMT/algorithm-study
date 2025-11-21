- 客户端存储：大的 Schema 可能有几 MB，可以考虑存到 indexDB，然后有一个面板可以看到历史几十条 Schema 记录；
- 但是排障的话就必须要用户自己的电脑了。

---

这是一个非常棒的架构思路，完全可行，而且是很多高级前端应用（如 Figma、VS Code Web 版）采用的策略。

这种方案的核心在于：**“本地全量记录（低成本） + 远程按需同步（高价值）”**。

你可以通过**“导出/分享快照”**的机制来解决“必须要在用户电脑上排障”的问题。

### 核心方案：IndexedDB + 按需上报

#### 1. 架构设计

- **日常运行（本地闭环）**：

  - 每次 Schema 变更，自动存入浏览器 **IndexedDB**。
  - IndexedDB 容量大（通常是磁盘剩余空间的 50%以上），存几十条甚至几百条 5MB 的 JSON 毫无压力。
  - 提供一个**“调试面板”**（Time Travel Panel），读取 IndexedDB，让用户可以回滚或查看历史版本的 Schema。

- **排障时刻（按需桥接）**：
  - **场景 A（人工排障）**：用户在调试面板中选中某条记录，点击“**生成排障链接**”或“**上报当前状态**”。前端将这条特定的 JSON 上传到服务器（OSS/S3），生成一个临时 ID，发给开发人员。
  - **场景 B（AI Chat）**：当 AI 需要分析时，前端从 IndexedDB 取出当前最新的那条，发送给 AI（或者先上传换取 URL 给 AI）。

#### 2. 方案优势

1.  **零平时带宽成本**：用户正常操作时，不需要把巨大的 Schema 传给服务器。
2.  **“时光机”功能**：用户操作错了，可以通过本地历史记录回退（Undo/Redo 的超级版）。
3.  **隐私安全**：数据默认在本地，只有用户主动点击“上报”时才流出。

#### 3. 代码实现思路

我们可以使用 `idb` (一个轻量级的 IndexedDB 包装库) 来简化操作。

```typescript
import { openDB } from 'idb'

const DB_NAME = 'LanderSchemaDB'
const STORE_NAME = 'snapshots'
const MAX_HISTORY = 50 // 保留最近50条

// 1. 初始化数据库
const dbPromise = openDB(DB_NAME, 1, {
  upgrade(db) {
    // 创建存储库，使用时间戳作为主键
    if (!db.objectStoreNames.contains(STORE_NAME)) {
      db.createObjectStore(STORE_NAME, { keyPath: 'timestamp' })
    }
  }
})

// 2. 保存快照 (自动清理旧数据)
export async function saveSnapshot(schema: any) {
  const db = await dbPromise
  const timestamp = Date.now()

  await db.put(STORE_NAME, {
    timestamp,
    schema,
    summary: `Update at ${new Date(timestamp).toLocaleTimeString()}` // 用于列表展示的摘要
  })

  // 清理多余的历史记录
  const count = await db.count(STORE_NAME)
  if (count > MAX_HISTORY) {
    const keys = await db.getAllKeys(STORE_NAME)
    // 删除最旧的记录 (keys 默认是排序的)
    const keysToDelete = keys.slice(0, keys.length - MAX_HISTORY)
    const tx = db.transaction(STORE_NAME, 'readwrite')
    await Promise.all(keysToDelete.map(key => tx.store.delete(key)))
    await tx.done
  }
}

// 3. 获取历史列表 (用于调试面板)
export async function getHistoryList() {
  const db = await dbPromise
  // 只获取摘要信息，不加载巨大的 schema 内容，保证面板渲染快
  const all = await db.getAll(STORE_NAME)
  return all.map(({ timestamp, summary }) => ({ timestamp, summary })).reverse()
}

// 4. 【关键】上报特定记录用于排障
export async function uploadSnapshotForSupport(timestamp: number) {
  const db = await dbPromise
  const record = await db.get(STORE_NAME, timestamp)

  if (!record) throw new Error('Snapshot not found')

  // 模拟上传到 OSS/S3
  const formData = new FormData()
  formData.append('file', new Blob([JSON.stringify(record.schema)]), 'schema.json')

  // 真实场景：调用你的后端上传接口
  const response = await fetch('/api/debug/upload-snapshot', {
    method: 'POST',
    body: formData
  })

  const { snapshotId } = await response.json()
  return snapshotId // 返回 ID，用户把这个 ID 发给客服/研发，或者发给 Chat
}
```

#### 4. 交互流程设计

为了解决“必须要在用户电脑上”的问题，你可以在低代码平台的界面右下角做一个 **"Debug Mode"**：

1.  **用户遇到问题**。
2.  用户打开 **"历史记录面板"**。
3.  面板列出了最近 50 次操作的时间点。
4.  用户点击最近的一条，选择 **"分享给技术支持"**。
5.  前端执行 `uploadSnapshotForSupport`，上传成功后弹出一个 ID：`Debug-ID: 89757`。
6.  用户把这个 ID 发给 Chat 机器人，或者发给人工客服。
7.  **排障解决**：Chat/客服 根据 ID 从云端拉取那份巨大的 Schema 进行分析。

### 总结

**可以这样做，且非常推荐。**

这实际上是把 IndexedDB 当作了**“本地的 Git”**，而上报接口就是 `git push`。只有在需要协作（排障）时才 Push，平时都在本地 Commit。

---

这个“**本地全量记录 (IndexedDB) + 按需上报**”的方案已经非常成熟，但针对**低代码平台**和**排障**这两个特定场景，还有几个关键点可以进一步优化，以提升性能、用户体验和排障效率。

以下是 5 个维度的改进建议：

### 1. 存储优化：从“全量快照”改为“增量差异 (Diff)”

**问题**：如果 Schema 有 5MB，用户操作 20 次就是 100MB。虽然 IndexedDB 存得下，但写入和读取会变慢，且占用用户磁盘空间。
**改进**：

- **基准 + 补丁**：只存储第一个版本的全量 Schema，后续操作只存储 **JSON Patch** (差异)。
- **定期快照**：每隔 10 次操作存一个全量快照（避免恢复时计算量过大）。
- **库推荐**：使用 `immer` (自带 patch 生成) 或 `fast-json-patch`。

```typescript
// 伪代码示例
import { produceWithPatches } from 'immer'

// 产生变更时，只存 patches
const [nextState, patches, inversePatches] = produceWithPatches(currentState, draft => {
  draft.components[0].props.color = 'red'
})

// 存入 IndexedDB 的只有极小的 patches 数组，而不是整个 state
await db.add('history', { type: 'patch', data: patches, timestamp: Date.now() })
```

### 2. 性能优化：将存储操作移至 Web Worker

**问题**：JSON 的序列化（`JSON.stringify`）和 IndexedDB 的写入操作虽然是异步的，但在处理大对象（MB 级别）时，仍然可能占用主线程 CPU，导致页面瞬间掉帧或卡顿。
**改进**：

- **Off-main-thread**：创建一个专门的 `SchemaWorker`。
- **流程**：主线程将 Schema 对象 `postMessage` 给 Worker -> Worker 进行压缩/Diff 计算 -> Worker 写入 IndexedDB。
- **结果**：UI 线程完全不受影响，用户操作丝般顺滑。

### 3. 排障上下文：增加“视觉快照”和“操作日志”

**问题**：有时候 Schema 看起来是没问题的，但渲染出来是错的（比如 CSS 冲突、浏览器兼容性）。光看代码（Schema）找不到原因。
**改进**：

- **视觉快照**：在上报时，利用 `html2canvas` 截取当前画布的一张截图，随 Schema 一起上传。客服一眼就能看到“哪里坏了”。
- **操作面包屑 (Breadcrumbs)**：记录用户最近的点击、报错日志。
  - _Schema 只是结果，操作日志是过程。_
- **库推荐**：`rrweb`。这是一个神器，它可以录制 DOM 变化。你可以选择只录制最近 30 秒的操作回放，随 Schema 一起上报。

### 4. 隐私与生命周期：设置数据过期 (TTL)

**问题**：用户上传的 Schema 可能包含敏感业务数据。如果永久存储在 OSS/S3 上，会有合规风险。
**改进**：

- **自动过期**：上传到云端的文件设置 **TTL (Time To Live)**，例如 7 天后自动删除。
- **阅后即焚**：或者在工单结单后自动触发删除。
- **加密**：在前端使用简单的对称加密（如 AES）打包数据，密钥只显示在用户屏幕上（`Debug-ID: 89757#Key123`），服务端只存密文。只有拿到完整 ID 的人才能解密查看。

### 5. 闭环体验：Deep Link (深度链接) 复现

**问题**：开发人员拿到 Schema JSON 后，还需要手动导入到本地环境复现，很麻烦。
**改进**：

- **一键复现 URL**：
  上报成功后，生成的不仅仅是一个 ID，而是一个链接：
  `https://lowcode.internal/debug?snapshot_url=https://oss.../schema.json`
- **Debug 模式**：
  当开发人员打开这个链接时，平台自动进入“只读调试模式”，从 URL 加载 Schema 覆盖当前状态。这样开发人员点击链接，直接就能看到用户当时的界面。

### 总结架构图

建议的最终形态如下：

```mermaid
graph TD
    User[用户操作] -->|变更| MainThread
    MainThread -->|postMessage| Worker[Web Worker]

    subgraph "Web Worker (不卡顿 UI)"
        Worker -->|计算 Diff| Patcher
        Worker -->|压缩| Compressor
        Compressor -->|存储| IDB[(IndexedDB)]
    end

    User -->|点击"一键上报"| Uploader

    subgraph "上报包 (Payload)"
        Schema[当前 Schema]
        Screenshot[界面截图]
        Logs[控制台报错/操作日志]
    end

    IDB -->|读取| Uploader
    Uploader -->|上传| OSS[云存储 (7天过期)]
    OSS -->|返回 URL| ChatBot
    ChatBot -->|分析| AI_Dev[AI / 研发人员]
```

---

**补充：也可以存日志数据库，定期删除过期的**
