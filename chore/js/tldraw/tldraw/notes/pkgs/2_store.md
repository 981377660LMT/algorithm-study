好的，我们来深入详细地讲解您工作区中的 **`@tldraw/store`** 包。

这个包是 `tldraw` 数据管理的核心，可以被理解为一个**专为 `tldraw` 设计的、高性能的、支持撤销/重做的响应式内存数据库**。它构建于我们之前讨论的 `@tldraw/state` 之上。

我们将从其核心职责、架构设计和关键代码文件三个方面进行剖析。

---

### **1. 核心职责**

如 [`CONTEXT.md`](packages/store/CONTEXT.md) 所述，`@tldraw/store` 的主要职责是：

1.  **存储所有数据**: 它是一个容器，保存着 `tldraw` 画布上的所有“记录 (Records)”，包括图形、页面、相机、用户状态等。
2.  **提供原子性操作**: 提供 `put` (添加/更新) 和 `remove` (删除) 等方法来修改数据，并保证这些操作是原子性的（要么全部成功，要么全部失败）。
3.  **管理历史记录**: 自动记录每一次数据变更，并提供 `undo()` 和 `redo()` 方法，实现撤销和重做功能。
4.  **响应式更新**: 当数据发生变化时，能够自动通知所有依赖该数据的地方进行更新。
5.  **数据迁移**: 在加载旧版本数据时，能够自动将其更新到最新的数据结构。

---

### **2. 架构设计与核心概念**

`@tldraw/store` 的设计精妙地结合了响应式编程和不可变数据的思想。

#### **a. Store 与 Record**

- **`Store` 类**: 这是整个包的中心。一个 `Store` 实例就代表一个完整的画布数据库。它内部维护着所有的数据记录和一个历史记录栈。
- **`BaseRecord`**:
  - **文件**: [`src/lib/BaseRecord.ts`](packages/store/src/lib/BaseRecord.ts)
  - **概念**: 这是所有存储在 `Store` 中数据的基类。任何东西，无论是图形、资源还是相机，都必须继承自 `BaseRecord`。它强制要求每个记录都必须有一个唯一的 `id` 和一个 `typeName`。

#### **b. 不可变性 (Immutability) 与 `ImmutableMap`**

- **文件**: [`src/lib/ImmutableMap.ts`](packages/store/src/lib/ImmutableMap.ts)
- **概念**: 这是理解 `Store` 性能和历史记录机制的关键。`Store` 内部不使用原生的 `Map` 或对象来存储数据，而是使用一个自定义的 `ImmutableMap`。
- **工作原理**: 当你调用 `store.put(newRecord)` 时，`Store` **不会**直接修改现有的数据集合。相反，它会创建一个**新的 `ImmutableMap` 实例**，这个新实例包含了旧实例的所有数据以及你的新记录。这种写时复制 (Copy-on-Write) 的策略有几个巨大优势：
  1.  **历史追溯**: 旧版本的数据永远不会丢失，可以轻松地在不同版本之间切换，这是实现撤销/重做的基础。
  2.  **性能优化**: `ImmutableMap` 内部使用持久化数据结构，创建新版本时可以共享大部分未改变的内部节点，使得复制操作非常快速和节省内存。
  3.  **响应式优化**: 可以通过简单的引用比较 (`oldMap === newMap`) 来快速判断数据是否发生了变化。

#### **c. 变更集 (Diff) 与历史记录**

- **文件**: [`src/lib/RecordsDiff.ts`](packages/store/src/lib/RecordsDiff.ts)
- **概念**: `Store` 并不直接记录每一次操作的指令，而是记录操作前后的**数据差异 (Diff)**。
- **`RecordsDiff` (或 `HistoryEntry`)**: 每当一个事务结束时，`Store` 会计算出这个事务中：
  - `added`: 添加了哪些记录。
  - `updated`: 更新了哪些记录（同时包含更新前 `[from]` 和更新后 `[to]` 的值）。
  - `removed`: 删除了哪些记录。
    这个 `RecordsDiff` 对象就是历史记录栈中的一个条目。
- **`undo()` / `redo()` 的实现**:
  - **`undo()`**: `Store` 取出历史记录栈顶的 `RecordsDiff`，并应用它的**“反向操作”**。例如，对于一个 `updated` 操作，它会将记录从 `[to]` 的值改回 `[from]` 的值。
  - **`redo()`**: `Store` 应用之前被撤销的 `RecordsDiff` 的**“正向操作”**。

#### **d. 响应式查询 (`executeQuery`)**

- **文件**: [`src/lib/executeQuery.ts`](packages/store/src/lib/executeQuery.ts)
- **概念**: `Store` 允许你创建对数据的**响应式查询**。这些查询本质上是基于 `@tldraw/state` 的 `computed` 信号。
- **工作原理**:
  1.  你定义一个查询函数，例如 `() => store.getAllShapesOnCurrentPage()`。
  2.  `executeQuery` 会创建一个 `computed` 信号来执行这个函数。
  3.  这个 `computed` 信号会自动追踪它所依赖的所有 `Store` 内部数据。
  4.  当任何被依赖的数据（例如，一个图形被移动，或者页面被切换）发生变化时，这个查询会自动重新执行，并返回最新的结果。
      这使得 UI 能够以极高的性能自动响应数据的变化。

---

### **3. 关键代码文件解析**

- **[`src/index.ts`](packages/store/src/index.ts)**: 这是包的公共 API 入口。它导出了最重要的工厂函数 `createTLStore`，以及 `loadSnapshot` 和 `getSnapshot` 等用于持久化的工具函数。

- **`Store.ts` (核心类)**:

  - `constructor`: 接收一个初始快照和 `TLSchema`。`TLSchema` 用于验证和迁移数据。
  - `put(records)` / `remove(ids)`: 修改数据的核心方法。它们内部会启动一个事务，计算 `diff`，并更新内部的 `ImmutableMap`。
  - `listen(callback)`: 允许外部代码订阅 `Store` 的变更。每当一个事务提交后，它就会带着 `RecordsDiff` 调用所有监听器。这是实现多人协作同步的关键。
  - `undo()` / `redo()`: 操作历史记录栈。
  - `get(id)`: 从 `Store` 中获取单个记录。

- **[`src/lib/migrate.ts`](packages/store/src/lib/migrate.ts)**:

  - **作用**: 负责数据迁移。当 `loadSnapshot` 加载一个带有旧版本 `schema` 的快照时，`Store` 会调用这里的迁移逻辑，逐个版本地运行迁移脚本，直到数据结构与当前 `schema` 匹配为止。

- **[`src/lib/devFreeze.ts`](packages/store/src/lib/devFreeze.ts)**:
  - **作用**: 一个仅在开发模式下启用的工具。它会深度冻结 (deep-freeze) 从 `Store` 中取出的所有记录。如果你尝试直接修改一个记录的属性（例如 `shape.x = 100`），而不是通过 `store.put(shape)`，程序会立即抛出错误。这强制开发者遵循不可变数据的编程范式，避免了许多难以追踪的 bug。

### **总结：一次用户操作的完整流程**

1.  用户在画布上拖动一个矩形。
2.  UI 层触发 `editor.updateShapes(...)`。
3.  `Editor` 最终调用 `store.put([{...updatedShape}])`。
4.  `Store` 在一个事务中：
    a. 创建一个新的 `ImmutableMap`，其中包含更新后的矩形记录。
    b. 计算出 `RecordsDiff`，内容为 `{ updated: { [shapeId]: [oldShape, updatedShape] } }`。
    c. 将这个 `diff` 推入历史记录栈。
    d. 更新内部的、代表所有记录的 `atom` 信号。
5.  `@tldraw/state` 的调度器启动：
    a. 所有依赖这个矩形位置的 `computed` 信号（例如，一个显示选中框位置的查询）被标记为“脏”。
    b. 当 UI 重新渲染并访问这些 `computed` 信号时，它们会重新计算并返回新值。
6.  `Store` 的 `listen` 回调被触发，将 `diff` 发送出去（如果是多人协作场景，此时会通过 WebSocket 发送给其他客户端）。

通过这种设计，`@tldraw/store` 实现了一个既健壮又高效的数据层，完美支撑了 `tldraw` 复杂的交互和协作需求。
