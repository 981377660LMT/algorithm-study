好的，我们来深入分析 Yjs 的 UndoManager.js 的实现原理。

Yjs 是一个高性能的 CRDT (无冲突复制数据类型) 库，用于构建实时协作应用。它的 `UndoManager` 是一个非常精巧的设计，它必须在 CRDT 这种非线性的、保留所有历史的数据模型之上，实现用户所期望的线性“撤销/重做”体验。

### 核心思想：撤销不是“逆操作”，而是“反向操作”

在传统的编辑器中，撤销一个操作（比如插入 "hello"）通常是执行一个逆操作（删除 "hello"）。但在 Yjs (CRDT) 中，所有的数据结构都是只增的（append-only），删除一个项只是将其标记为“已删除”，而不是从数据结构中物理移除。

因此，`UndoManager` 的核心思想是：

1.  **捕获变更**: 监听一次或多次事务（Transaction），将期间发生的所有**插入**和**删除**记录下来。
2.  **创建快照**: 将这些变更集合（`insertions` 和 `deletions`）打包成一个 `StackItem`，存入 `undoStack`。
3.  **执行撤销 (Undo)**: 当用户撤销时，从 `undoStack` 中取出一个 `StackItem`。然后执行一个“反向操作”：
    - **删除** `StackItem` 中记录的**所有插入项**。
    - **恢复** `StackItem` 中记录的**所有删除项**（Yjs 中称为 `redoItem`）。
4.  **执行重做 (Redo)**: 撤销操作本身也会产生一个新的变更集合（删除了之前插入的，恢复了之前删除的）。这个新的变更集合被打包成一个新的 `StackItem` 放入 `redoStack`。当用户重做时，执行与撤销类似的反向操作。

---

### 关键数据结构与概念

1.  **`StackItem`**: 这是撤销/重做栈中的基本单元。

    ```javascript
    class StackItem {
      constructor(deletions, insertions) {
        this.insertions = insertions // 本次操作中所有“新插入”项的集合
        this.deletions = deletions // 本次操作中所有“被删除”项的集合
        this.meta = new Map() // 用于存储元数据，如光标位置
      }
    }
    ```

    - `insertions` 和 `deletions` 都是 `DeleteSet` 类型。`DeleteSet` 是 Yjs 中一个高效的数据结构，用于表示一组被删除的 `Item`（Yjs 中的原子操作单元）。尽管名字叫 `DeleteSet`，但 `insertions` 也是用它来存储的，因为它能高效地表示一系列连续的 `Item` ID。

2.  **`afterTransactionHandler`**: 这是 `UndoManager` 的心脏。它在每次 Yjs 文档事务结束后被调用，负责捕获变更并创建 `StackItem`。

    - **事务来源过滤 (`trackedOrigins`)**: 这是实现“只撤销自己操作”的关键。每个 Yjs 事务都可以有一个 `origin`。`UndoManager` 只会捕获那些 `origin` 在 `trackedOrigins` 集合中的事务。默认情况下，这个集合只包含 `null`（代表本地用户操作）和 `UndoManager` 实例本身。来自远程协作者的事务会有不同的 `origin`，因此会被忽略。
    - **变更合并 (`captureTimeout`)**: 为了获得更好的用户体验，连续快速的操作（如打字）应该被合并成一个撤销步骤。`UndoManager` 通过 `captureTimeout`（默认为 500ms）实现这一点。如果在超时时间内发生了新的变更，它会与上一个 `StackItem` 合并，而不是创建一个新的。
    - **栈管理**:
      - 当一个**普通操作**发生时，它会**清空 `redoStack`**。这是标准的撤销逻辑。
      - 当一个**撤销操作**发生时，它产生的“反向变更”会被打包成 `StackItem` 并推入 `redoStack`。
      - 当一个**重做操作**发生时，它产生的“反向变更”会被打包成 `StackItem` 并推入 `undoStack`。

3.  **`popStackItem`**: 这是执行撤销/重做的核心函数。
    - **输入**: `undoManager` 实例、要操作的栈 (`undoStack` 或 `redoStack`)、事件类型 (`'undo'` 或 `'redo'`)。
    - **核心逻辑**:
      1.  从指定的栈中弹出一个 `stackItem`。
      2.  **遍历 `stackItem.insertions`**: 对于这个 `StackItem` 记录的每一个“插入项”，执行删除操作 (`item.delete(transaction)`)。
      3.  **遍历 `stackItem.deletions`**: 对于这个 `StackItem` 记录的每一个“删除项”，执行恢复操作 (`redoItem(transaction, struct, ...)`）。`redoItem` 是一个内部函数，它会找到被标记为删除的项并“复活”它。
      4.  **处理冲突**: `redoItem` 内部逻辑非常复杂，它需要处理在原始项被删除后，其上下文可能已经发生变化的情况（例如，一个被删除的文本字符，其左右的字符可能已经被修改）。
      5.  **过滤 (`deleteFilter`)**: 在执行删除时，会调用用户提供的 `deleteFilter` 函数，允许用户阻止某些特定项被删除。
    - **原子性**: 所有这些反向操作都在一个**单独的 Yjs 事务**中完成，保证了撤销/重做操作的原子性。

### 工作流程串讲

假设用户输入了 "AB"。

1.  **捕获**:

    - 用户输入 'A'，触发一个事务。`afterTransactionHandler` 捕获到这个插入操作。
    - 它创建一个 `StackItem1`，其中 `insertions` 包含了代表 'A' 的 `Item`。`StackItem1` 被推入 `undoStack`。
    - 用户在 100ms 内输入 'B'，触发第二个事务。
    - `afterTransactionHandler` 再次被调用。它发现距离上次操作（`this.lastChange`）的时间小于 `captureTimeout`。
    - 因此，它**不会创建新 `StackItem`**，而是将代表 'B' 的 `Item` 合并到 `StackItem1` 的 `insertions` 集合中。`undoStack` 中仍然只有一个 `StackItem`，但它现在代表了插入 "AB" 这个完整的操作。

2.  **撤销 (Undo)**:

    - 用户点击撤销，调用 `undo()` 方法。
    - `undo()` 内部调用 `popStackItem(this, this.undoStack, 'undo')`。
    - `popStackItem` 从 `undoStack` 中弹出 `StackItem1`。
    - 它开启一个新事务，并遍历 `StackItem1.insertions`，对代表 'A' 和 'B' 的 `Item` 执行 `delete()` 操作。
    - 这个撤销事务本身也产生了变更：一个包含了 'A' 和 'B' 的 `deleteSet`。
    - 由于当前 `this.undoing` 为 `true`，`afterTransactionHandler` 会将这个新的变更（删除了 "AB"）打包成一个新的 `StackItem2`，并推入 `redoStack`。

3.  **重做 (Redo)**:
    - 用户点击重做，调用 `redo()` 方法。
    - `redo()` 内部调用 `popStackItem(this, this.redoStack, 'redo')`。
    - `popStackItem` 从 `redoStack` 中弹出 `StackItem2`。
    - 它开启一个新事务，并遍历 `StackItem2.deletions`（这里面是 'A' 和 'B'），对它们执行 `redoItem()` 操作，将 'A' 和 'B' 恢复。
    - 这个重做事务本身也产生了变更：一个包含了 'A' 和 'B' 的 `insertions` 集合。
    - 由于当前 `this.redoing` 为 `true`，`afterTransactionHandler` 会将这个新变更打包成 `StackItem3`（内容与 `StackItem1` 相同），并推入 `undoStack`。

至此，系统回到了撤销之前的状态，并且可以再次进行撤销。这个精巧的闭环设计保证了在复杂的协作环境下，单用户的撤销/重做功能依然能够正确、可靠地工作。

---

好的，我们继续基于您提供的 UndoManager.js 源码，深入讲解更多实现细节。

### 构造函数与精细化配置 (`constructor`)

`UndoManager` 的构造函数提供了丰富的选项，让开发者可以精细化地控制其行为：

- **`typeScope`**: 这是最重要的参数，它定义了 `UndoManager` 的“管辖范围”。

  - 如果传入整个 `Doc` 实例，它会捕获文档上的所有变更。
  - 如果传入一个或多个具体的 `Y.Type`（如 `Y.Text`、`Y.Array`），它就只会捕获这些特定类型及其子节点上的变更。这对于实现局部撤销（例如，只撤销某个输入框的修改，而不影响页面其他部分）至关重要。

- **`captureTimeout` (默认 500ms)**: 这是实现**操作合并**的关键。如之前所述，在 500ms 内连续发生的操作会被合并成一个撤销步骤。

- **`captureTransaction`**: 一个函数，允许你基于事务（`Transaction`）的某些属性来决定是否要捕获这次变更。例如，你可以忽略由特定插件或系统行为产生的事务。

- **`deleteFilter` (默认 `() => true`)**: 一个过滤器函数。在执行撤销/重做时，如果要删除某个 `Item`，会先调用此函数。如果返回 `false`，则该 `Item` 不会被删除。这提供了一个“保护”某些内容不被撤销操作删除的机制。

- **`trackedOrigins` (默认 `new Set([null])`)**: 一个集合，定义了哪些“来源”的事务应该被追踪。
  - `null` 是本地用户操作的默认 `origin`。
  - `UndoManager` 实例自身也会被加入这个集合，因为撤销/重做操作本身产生的事务需要被捕获以放入对方的栈中。
  - 来自远程协作者的事务通常有不同的 `origin`，因此默认不会被 `UndoManager` 捕获，从而实现了“只撤销自己的操作”。你可以通过 `addTrackedOrigin` 方法添加自定义的来源。

### 核心事件处理器 (`afterTransactionHandler`)

这是 `UndoManager` 的“大脑”，它在每次事务结束后被触发，决定如何处理这次变更。

1.  **过滤 (L251-257)**: 首先，它会进行一系列检查，如果满足以下任一条件，则直接忽略本次事务：

    - `captureTransaction` 函数返回 `false`。
    - 本次事务影响的类型不在 `UndoManager` 的 `scope` 内。
    - 事务的 `origin` 不在 `trackedOrigins` 集合中。

2.  **栈管理 (L262-266)**:

    - 如果当前是**普通操作**（非撤销/重做），它会调用 `this.clear(false, true)` **清空 `redoStack`**。这是符合直觉的标准行为：一旦有了新的操作，之前的“重做”路径就无效了。
    - 如果当前是**撤销操作** (`undoing`)，它会将产生的“反向变更”推入 `redoStack`。
    - 如果当前是**重做操作** (`redoing`)，它会将产生的“反向变更”推入 `undoStack`。

3.  **捕获变更 (L268-276)**:

    - **`insertions`**: 通过比较事务前后的状态 `transaction.afterState` 和 `transaction.beforeState`，计算出本次事务中所有新插入的 `Item`。
    - **`deletions`**: 直接使用 `transaction.deleteSet`，它记录了本次事务中所有被删除的 `Item`。

4.  **合并或新建 `StackItem` (L279-287)**:

    - **合并**: 如果距离上次捕获操作的时间小于 `captureTimeout`，并且当前不是撤销/重做操作，它会找到 `undoStack` 的最后一个 `StackItem`，并将本次的 `insertions` 和 `deletions` 合并进去。
    - **新建**: 否则，它会用本次的 `deletions` 和 `insertions` 创建一个新的 `StackItem` 并推入栈中。

5.  **内存管理 - 防止垃圾回收 (L294-298)**:
    - 这是非常关键的一步。对于所有被删除的 `Item`，它会调用 `keepItem(item, true)`。
    - **作用**: Yjs 有一个垃圾回收 (GC) 机制，会清理掉被删除且不再被任何地方引用的 `Item`。`keepItem(true)` 相当于给这个 `Item` 打上一个“请勿回收”的标记。因为即使用户删除了某些内容，`UndoManager` 仍然需要保留这些 `Item` 的信息，以便未来可以通过“撤销”操作将它们“复活”。

### 撤销/重做执行器 (`popStackItem`)

这个函数是实际执行撤销/重做的“手臂”。

1.  **循环与空操作处理 (L64)**: `while (stack.length > 0 && undoManager.currStackItem === null)` 这个循环确保了即使栈顶的 `StackItem` 是一个“空操作”（即它引用的 `Item` 已经被远程用户永久删除了，导致撤销它不会产生任何实际效果），`UndoManager` 也会继续弹出下一个，直到找到一个能产生实际变更的 `StackItem` 或者栈为空。

2.  **分离待办事项 (L74-106)**:

    - 它遍历 `stackItem.insertions`，找到所有需要被**删除**的 `Item`，放入 `itemsToDelete` 数组。
    - 它遍历 `stackItem.deletions`，找到所有需要被**恢复**的 `Item`，放入 `itemsToRedo` 集合。

3.  **执行恢复 (L107)**:

    - 调用 `redoItem(...)` 来“复活”`itemsToRedo` 中的 `Item`。`redoItem` 是 Yjs 内部一个非常复杂的函数，它负责将被标记为删除的 `Item` 重新整合回文档结构中，并智能地处理其上下文变化（比如它原来的邻居节点已经不存在了）。

4.  **执行删除 (L110-115)**:
    - **倒序删除**: `for (let i = itemsToDelete.length - 1; i >= 0; i--)`。这里采用倒序遍历删除，这是一个经典的树操作技巧。因为在 Yjs 的结构中，子节点总是在父节点之后创建。倒序删除可以确保先删除子节点再删除父节点，避免了因父节点先被删除而导致子节点信息不完整的问题。
    - 在删除前会通过 `undoManager.deleteFilter(item)` 进行最后一次校验。

### 内存管理 - 释放内存 (`clear` 和 `clearUndoManagerStackItem`)

- 当调用 `undoManager.clear()` 清空 `undoStack` 或 `redoStack` 时，会遍历栈中所有的 `StackItem`。
- `clearUndoManagerStackItem` (L43) 函数会被调用，它内部执行 `keepItem(item, false)`。
- 这会移除之前设置的“请勿回收”标记，告诉 Yjs 的 GC：“`UndoManager` 不再需要这些 `Item` 了，如果它们没有被其他地方引用，你现在可以安全地回收它们了”。这有效地防止了内存泄漏。
