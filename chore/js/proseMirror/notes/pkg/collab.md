好的，我们来对 `prosemirror-collab` 的源码进行一次逐行逐句、深入骨髓的剖析，确保您能彻底理解其内部的每一个细节和设计决策。

我们将按照以下逻辑顺序来解构这份源码，这比单纯从上到下阅读更符合理解算法的流程：

1.  **数据结构 (Data Structures)**: `CollabState` 和 `Rebaseable`，它们是整个协同算法的“名词”。
2.  **`collab()` 插件**: 它是如何作为“记录员”，在本地捕捉和管理用户变更的。
3.  **`sendableSteps()`**: 它是如何作为“信使”，打包本地变更准备发往服务器的。
4.  **`receiveTransaction()`**: 它是如何作为“外交官”，处理从服务器收到的远程变更的。
5.  **`rebaseSteps()`**: 它是整个算法的“心脏”，执行最关键的“变基”魔法。

---

### 1. 数据结构 (Data Structures): 协同算法的基石

#### `class CollabState`

```typescript
class CollabState {
  constructor(readonly version: number, readonly unconfirmed: readonly Rebaseable[]) {}
}
```

这是 `collab` 插件自身的状态。它非常纯粹，只包含两个核心信息：

- `version: number`: **我同步到的服务器版本号**。这代表了当前客户端的文档状态与中央服务器权威历史的同步点。每次成功接收并应用一批来自服务器的 `steps` 后，这个版本号就会增加。
- `unconfirmed: readonly Rebaseable[]`: **我本地`已应用`、但尚未被服务器确认的变更列表**。这是“乐观更新”策略的体现。用户操作后，变更会立即加入此列表并在本地生效，然后等待服务器的“确认回执”。

#### `class Rebaseable`

```typescript
class Rebaseable {
  constructor(readonly step: Step, readonly inverted: Step, readonly origin: Transform) {}
}
```

为什么 `unconfirmed` 列表里存的不是简单的 `Step`，而是这个 `Rebaseable` 对象？因为在执行“变基”操作时，我们需要更多信息：

- `step: Step`: 变更本身。这是最基本的数据。
- `inverted: Step`: **该变更的逆操作**。这是 `rebase` 魔法的关键之一。当需要临时“撤销”本地未确认的变更时，就是通过应用这个 `inverted` 步骤来实现的。
- `origin: Transform`: **产生这个 `step` 的原始事务 (`Transaction`)**。注意，`Transaction` 继承自 `Transform`。为什么需要它？`sendableSteps` 的注释给出了答案：为了查找时间戳等元数据。一个事务上可能附加了 `time`、`userID` 等 `meta` 信息。即使 `step` 本身在多次 `rebase` 后被变换得“面目全非”，但它的 `origin` 始终指向最初创建它的那个事务，从而保留了其原始上下文。

#### `function unconfirmedFrom(transform: Transform)`

```typescript
function unconfirmedFrom(transform: Transform) {
  let result = []
  for (let i = 0; i < transform.steps.length; i++)
    result.push(
      new Rebaseable(transform.steps[i], transform.steps[i].invert(transform.docs[i]), transform)
    )
  return result
}
```

这个辅助函数的作用就是将一个事务（`transform`）中的所有 `steps` 转换为 `Rebaseable` 对象数组。它遍历事务中的每个 `step`，并为每个 `step` 计算出其逆操作 `inverted`，然后将 `step`、`inverted` 和事务本身 `transform` 一起打包成一个新的 `Rebaseable` 实例。

---

### 2. `collab()` 插件: 本地变更的“记录员”

```typescript
export function collab(config: CollabConfig = {}): Plugin {
  // ... config initialization ...

  return new Plugin({
    key: collabKey,

    state: {
      init: () => new CollabState(conf.version, []),
      apply(tr, collab) {
        let newState = tr.getMeta(collabKey)
        if (newState) return newState
        if (tr.docChanged)
          return new CollabState(collab.version, collab.unconfirmed.concat(unconfirmedFrom(tr)))
        return collab
      }
    },

    config: conf,
    historyPreserveItems: true
  })
}
```

这是模块的主入口，一个标准的插件工厂。我们重点关注 `state.apply` 方法，它是处理本地编辑的核心。

`apply(tr, collab)` 在每次 `EditorState` 更新时被调用，`tr` 是导致这次更新的事务，`collab` 是旧的 `CollabState`。

1.  `let newState = tr.getMeta(collabKey)`: 这是一个“后门”或“快捷通道”。`receiveTransaction` 函数在处理完远程变更后，会直接计算出一个全新的 `CollabState`，并通过 `tr.setMeta(collabKey, ...)` 将其附加到事务上。这里就是用来接收这个新状态的。如果存在，就直接使用它，跳过后续所有逻辑。

2.  `if (newState) return newState`: 如果是通过 `meta` 传入的新状态，直接返回，应用它。

3.  `if (tr.docChanged)`: 如果不是通过 `meta` 传入，并且这个事务改变了文档（通常意味着是用户的本地编辑），则执行以下操作：

    - `unconfirmedFrom(tr)`: 将当前事务中的所有 `steps` 转换为 `Rebaseable` 对象。
    - `collab.unconfirmed.concat(...)`: 将这些新的 `Rebaseable` 对象追加到旧的 `unconfirmed` 列表末尾。
    - `new CollabState(collab.version, ...)`: 创建一个新的 `CollabState`。注意，`version` 保持不变，因为这只是本地变更，尚未与服务器同步。

4.  `return collab`: 如果事务没有改变文档（如仅移动光标），则状态不变。

`historyPreserveItems: true` 这个配置是告诉 `prosemirror-history` 插件：“不要合并来自协同编辑事务的步骤”。这确保了每个 `Step` 都是独立的，可以被精确地 `rebase` 和追踪，这对于协同编辑至关重要。

---

### 3. `sendableSteps()`: 打包变更的“信使”

```typescript
export function sendableSteps(state: EditorState): { ... } | null {
  let collabState = collabKey.getState(state) as CollabState
  if (collabState.unconfirmed.length == 0) return null
  return {
    version: collabState.version,
    steps: collabState.unconfirmed.map(s => s.step),
    clientID: (collabKey.get(state)!.spec as any).config.clientID,
    get origins() { ... }
  }
}
```

这个函数非常直白，它扮演了信使的角色，负责将需要发送给服务器的数据打包好。

1.  获取当前的 `collabState`。
2.  如果 `unconfirmed` 列表为空，说明没有需要发送的变更，返回 `null`。
3.  如果 `unconfirmed` 列表不为空，就构建一个包含三/四个关键部分的对象：
    - `version`: 客户端当前的文档版本。服务器会用它来检查提交是否基于最新版本。
    - `steps`: 从 `unconfirmed` 列表中的 `Rebaseable` 对象里提取出的纯 `Step` 数组。
    - `clientID`: 当前客户端的 ID。
    - `origins`: 一个 getter 属性，懒加载地从 `Rebaseable` 对象中提取出原始事务数组。

你的应用代码需要做的就是：定期调用 `sendableSteps()`，如果返回值不为 `null`，就将其通过 WebSocket 或其他方式发送给服务器。

---

### 4. `receiveTransaction()`: 处理远程变更的“外交官”

这是处理服务器广播的核心函数，也是逻辑最复杂的部分之一。

```typescript
export function receiveTransaction(
  state: EditorState,
  steps: readonly Step[],
  clientIDs: readonly (string | number)[]
  // ... options
) {
  // 1. 获取当前状态
  let collabState = collabKey.getState(state)
  let version = collabState.version + steps.length
  let ourID = ...

  // 2. 识别并分离自己的变更
  let ours = 0
  while (ours < clientIDs.length && clientIDs[ours] == ourID) ++ours
  let unconfirmed = collabState.unconfirmed.slice(ours)
  steps = ours ? steps.slice(ours) : steps

  // 3. 如果全是自己的变更，快速处理
  if (!steps.length) return state.tr.setMeta(collabKey, new CollabState(version, unconfirmed))

  // 4. 执行变基 (Rebase)
  let nUnconfirmed = unconfirmed.length
  let tr = state.tr
  if (nUnconfirmed) {
    unconfirmed = rebaseSteps(unconfirmed, steps, tr)
  } else {
    for (let i = 0; i < steps.length; i++) tr.step(steps[i])
    unconfirmed = []
  }

  // 5. 构建最终事务并返回
  let newCollabState = new CollabState(version, unconfirmed)
  // ... mapSelectionBackward logic ...
  return tr
    .setMeta('rebased', nUnconfirmed)
    .setMeta('addToHistory', false)
    .setMeta(collabKey, newCollabState)
}
```

让我们一步步分解这个“外交”过程：

1.  **获取当前状态**: 获取当前的 `collabState`、`clientID`，并预先计算出应用完这批 `steps` 后的新 `version`。

2.  **识别并分离自己的变更**: 这是非常关键的一步。服务器广播的 `steps` 列表，其开头部分可能就是我们自己刚刚提交并被服务器确认的。

    - `while (ours < clientIDs.length && clientIDs[ours] == ourID) ++ours`: 这个循环从 `clientIDs` 数组的开头开始，计算有多少个连续的 `step` 是由我们自己（`ourID`）产生的。
    - `let unconfirmed = collabState.unconfirmed.slice(ours)`: 我们从本地的 `unconfirmed` 列表中，**移除**掉这些已经被服务器确认的 `steps`。剩下的 `unconfirmed` 就是在“发送-接收”间隙中用户新产生的、真正未确认的变更。
    - `steps = ours ? steps.slice(ours) : steps`: 同样地，从服务器广播来的 `steps` 列表中，也移除掉我们自己的部分。剩下的 `steps` 就是纯粹来自**其他协作者**的变更。

3.  **快速路径**: `if (!steps.length)`，如果经过上一步分离后，发现所有 `steps` 都是我们自己的，那就意味着没有远程变更需要处理。我们只需更新 `CollabState`（包含新的 `version` 和缩短后的 `unconfirmed` 列表），然后通过 `setMeta` 快速返回一个事务即可。

4.  **执行变基 (Rebase)**: 这是核心冲突处理。

    - `if (nUnconfirmed)`: 如果我们本地还有未确认的变更（`unconfirmed`），就必须执行变基。
      - `unconfirmed = rebaseSteps(unconfirmed, steps, tr)`: 调用 `rebaseSteps` 函数（下一节详述），它会：
        1.  在 `tr` 上撤销本地的 `unconfirmed` 步骤。
        2.  在 `tr` 上应用远程的 `steps`。
        3.  在 `tr` 上重新应用经过坐标变换后的本地步骤。
        4.  返回变换后的、新的 `unconfirmed` 列表。
    - `else`: 如果本地没有未确认的变更，情况就简单了，直接将所有远程 `steps` 应用到 `tr` 上即可。

5.  **构建最终事务**:
    - `let newCollabState = new CollabState(version, unconfirmed)`: 创建最终的、全新的 `CollabState`。
    - `tr.setMeta(...)`: 给这个事务打上各种标记：
      - `'rebased'`: 告诉其他插件（如果需要）有多少个步骤被变基了。
      - `'addToHistory', false`: **极其重要**。告诉 `prosemirror-history` 插件，不要将这个事务记录到撤销栈中。因为这个事务包含的变更（无论是远程的还是变基后的本地变更）已经在各自的来源处被记录过了，在这里重复记录会导致历史记录混乱。
      - `collabKey, newCollabState`: 将最终的 `CollabState` 通过“后门”传递给 `apply` 方法。

---

### 5. `rebaseSteps()`: 算法的“心脏”

```typescript
export function rebaseSteps(
  steps: readonly Rebaseable[], // 本地未确认的变更 (Local)
  over: readonly Step[], // 远程发来的变更 (Remote)
  transform: Transform // 一个空的事务，用于执行操作
) {
  // 1. 撤销本地变更
  for (let i = steps.length - 1; i >= 0; i--) transform.step(steps[i].inverted)

  // 2. 应用远程变更
  for (let i = 0; i < over.length; i++) transform.step(over[i])

  let result = []
  // 3. 重新应用并转换本地变更
  for (let i = 0, mapFrom = steps.length; i < steps.length; i++) {
    // 3a. 映射 Step
    let mapped = steps[i].step.map(transform.mapping.slice(mapFrom))
    mapFrom--

    // 3b. 尝试应用变换后的 Step
    if (mapped && !transform.maybeStep(mapped).failed) {
      // 3c. 设置镜像，用于优化后续映射
      transform.mapping.setMirror(mapFrom, transform.steps.length - 1)
      // 3d. 将成功应用的、新的 Rebaseable 对象存入结果
      result.push(
        new Rebaseable(
          mapped,
          mapped.invert(transform.docs[transform.docs.length - 1]),
          steps[i].origin
        )
      )
    }
  }
  return result
}
```

这是 OT (Operational Transformation) 算法最直观的体现。让我们用一个例子来走一遍：

- 初始文档: `"A"`
- `steps` (Local): `[ insert(1, "B") ]` (用户输入了 B)
- `over` (Remote): `[ insert(0, "X") ]` (别人在你之前输入了 X)
- `transform`: 一个基于 `"A"` 的空事务。

1.  **撤销本地变更**:

    - `transform.step(steps[0].inverted)`: 应用 `insert(1, "B")` 的逆操作，即 `delete(1, "B")`。
    - `transform` 后的文档状态回到 `"A"`。
    - `transform.steps` 现在是 `[ delete(1, "B") ]`。

2.  **应用远程变更**:

    - `transform.step(over[0])`: 应用 `insert(0, "X")`。
    - `transform` 后的文档状态变为 `"XA"`。
    - `transform.steps` 现在是 `[ delete(1, "B"), insert(0, "X") ]`。

3.  **重新应用并转换本地变更**:

    - `i = 0`, `mapFrom = 1`。
    - **3a. 映射 Step**:
      - `transform.mapping`: 此刻它包含了 `delete(1, "B")` 和 `insert(0, "X")` 两个步骤的完整映射信息。
      - `transform.mapping.slice(mapFrom)`: `slice(1)` 表示我们只关心**第二个及之后**的 `step` (`insert(0, "X")`) 所产生的坐标变换。为什么？因为第一个 `step` 是我们自己撤销自己的操作，我们不希望它影响我们自己后续的重做。
      - `steps[0].step.map(...)`: 对 `insert(1, "B")` 进行映射。位置 `1` 经过 `insert(0, "X")` 的变换后，向后移动了 1 位，变成了 `2`。
      - `mapped` 变量现在是 `insert(2, "B")`。
    - **3b. 尝试应用**:
      - `transform.maybeStep(mapped)`: 将 `insert(2, "B")` 应用到当前文档 `"XA"` 上。
      - 应用成功，文档变为 `"XAB"`。`transform.steps` 变为 `[ delete(1, "B"), insert(0, "X"), insert(2, "B") ]`。
    - **3c. 设置镜像**: `setMirror` 是一个性能优化，它告诉 `mapping` 对象：“第 0 个原始 `step` 的变换结果对应于最终 `steps` 数组中的第 2 个 `step`”。这在有多个 `step` 需要 `rebase` 时可以避免重复计算。
    - **3d. 存入结果**: 创建一个新的 `Rebaseable` 对象，其中包含了变换后的 `mapped` 步骤及其新的逆操作，但 `origin` 保持不变。然后将其推入 `result` 数组。

4.  **返回**: 函数返回 `result`，即 `[ new Rebaseable(insert(2, "B"), ...) ]`。这个返回值会成为 `receiveTransaction` 中新的 `unconfirmed` 列表。

至此，`rebaseSteps` 完成了它的使命：在保留远程变更的同时，智能地、无冲突地重放了本地的乐观更新，并计算出了新的、等待提交的 `unconfirmed` 列表。整个协同编辑的核心魔法就在于此。

---

好的，我们换一种全新的方式，像一个侦探一样，拿着放大镜，一步一步地追踪这个循环里的数据流动。这次我们不用比喻，只看代码和数据，保证彻底弄懂。

### 案发现场设定

- **初始文档**: `"A"`
- **你的本地操作 (`steps`)**: 你连续输入了两个字符 "B" 和 "C"。
  - `L0` (本地步骤 0): `insert(1, "B")` (文档变为 "AB")
  - `L1` (本地步骤 1): `insert(2, "C")` (文档变为 "ABC")
- **同事的远程操作 (`over`)**: 在你操作时，他在开头插入了 "X"。
  - `R0` (远程步骤 0): `insert(0, "X")`

### `rebaseSteps` 函数的前两步

在进入我们关心的这个 `for` 循环之前，函数已经做了两件事：

1.  **撤销你的操作**: `transform` 对象依次执行了 `L1` 和 `L0` 的逆操作。
    - `transform.step( Invert(L1) )` (即 `delete(2, "C")` )
    - `transform.step( Invert(L0) )` (即 `delete(1, "B")` )
2.  **应用同事的操作**: `transform` 对象执行了 `R0`。
    - `transform.step( R0 )` (即 `insert(0, "X")` )

**关键时刻**: 此刻，进入 `for` 循环之前，`transform` 对象内部的 `steps` 列表是这样的：

```
// transform.steps 列表内容:
[
  // 索引 0: 你的 L1 的逆操作
  delete(2, "C"),

  // 索引 1: 你的 L0 的逆操作
  delete(1, "B"),

  // 索引 2: 同事的 R0 操作
  insert(0, "X")
]
```

而 `transform` 对象所管理的**文档内容**已经变成了 `"XA"`。

---

### 侦探开始追踪 `for` 循环

#### **第一次循环 (i = 0): 重做你的第一个操作 `L0`**

- **`i = 0`**, `mapFrom` 被初始化为 `steps.length`，即 **`2`**。

- **`let mapped = steps[0].step.map(transform.mapping.slice(mapFrom))`**

  - `steps[0]` 是 `L0`，也就是 `insert(1, "B")`。
  - `transform.mapping.slice(2)`: 这句话是关键，它的意思是：“**只看 `transform.steps` 列表里从索引 2 开始的那些步骤所产生的坐标变化**”。
  - 查看我们的列表，从索引 2 开始只有 `insert(0, "X")`。
  - 所以，这行代码的意思就变成了：“`insert(1, "B")` 这个操作，如果受到了 `insert(0, "X")` 的影响，会变成什么样？”
  - **计算**: 在开头插入 "X"，会让原来的位置 1 向后移动 1 位，变成位置 2。
  - **结果**: `mapped` 变量现在是 `insert(2, "B")`。

- **`mapFrom--`**: `mapFrom` 变成了 **`1`**。

- **`if (mapped && ...)`**: `insert(2, "B")` 是有效的，`transform.maybeStep` 把它应用到当前文档 `"XA"` 上，文档变为 `"XAB"`。

- **`transform.mapping.setMirror(mapFrom, transform.steps.length - 1)`**

  - `mapFrom` 现在是 `1`。
  - `transform.steps.length - 1` 是刚刚添加进去的 `mapped` 步骤的索引，即索引 3。
  - 这行代码的意思是：“**告诉系统，列表里索引为 1 的步骤 (`delete(1, "B")`) 和索引为 3 的步骤 (`insert(2, "B")`) 是一对镜像！**”
  - **为什么？** 因为 `insert(2, "B")` 正是 `delete(1, "B")` 在经历了 `insert(0, "X")` 的变换后，再反转过来的结果。系统记下这个快捷方式，下次计算会更快。

- **`result.push(...)`**: 把变换后的 `insert(2, "B")` 及其新的逆操作存起来。

#### **第二次循环 (i = 1): 重做你的第二个操作 `L1`**

- **`i = 1`**, `mapFrom` 现在是 **`1`**。

- **`let mapped = steps[1].step.map(transform.mapping.slice(mapFrom))`**

  - `steps[1]` 是 `L1`，也就是 `insert(2, "C")`。
  - `transform.mapping.slice(1)`: 这次的意思是：“**看 `transform.steps` 列表里从索引 1 开始的所有步骤所产生的坐标变化**”。
  - 查看列表，从索引 1 开始有：
    - `delete(1, "B")` (索引 1)
    - `insert(0, "X")` (索引 2)
    - `insert(2, "B")` (索引 3, 上一步重做的)
  - 所以，代码的意思是：“`insert(2, "C")` 这个操作，如果受到了这三个步骤的影响，会变成什么样？”
  - **计算**:
    1.  `delete(1, "B")`：删除 "B" 会让位置 2 变为位置 1。
    2.  `insert(0, "X")`：插入 "X" 会让位置 1 变为位置 2。
    3.  `insert(2, "B")`：插入 "B" 会让位置 2 变为位置 3。
  - **结果**: `mapped` 变量现在是 `insert(3, "C")`。

- **`mapFrom--`**: `mapFrom` 变成了 **`0`**。

- **`if (mapped && ...)`**: `insert(3, "C")` 是有效的，`transform.maybeStep` 把它应用到当前文档 `"XAB"` 上，文档变为 `"XABC"`。

- **`transform.mapping.setMirror(mapFrom, transform.steps.length - 1)`**

  - `mapFrom` 现在是 `0`。
  - 新步骤的索引是 `4`。
  - 这行代码的意思是：“**告诉系统，列表里索引为 0 的步骤 (`delete(2, "C")`) 和索引为 4 的步骤 (`insert(3, "C")`) 是一对镜像！**”

- **`result.push(...)`**: 把变换后的 `insert(3, "C")` 存起来。

### 最终结论

- **`mapFrom` 的作用**: 它是一个**滑动窗口的起点**。在每次循环中，它都精确地告诉 `map` 函数：“你只需要关心从我这个位置开始的后续所有变化，我之前那些你不用管（因为那是我自己的逆操作，会干扰计算）”。`mapFrom--` 就是在滑动这个窗口。

- **`setMirror` 的作用**: 它是一个**聪明的备忘录**。它告诉系统：“我刚刚完成了一次‘撤销-变换-重做’，这个‘撤销’和这个‘重做’是一对儿，你记一下。” 当系统下次要做类似的复杂计算时，看到这个备忘录，就可以走捷径，大大提高效率。

通过这个过程，你最初的两个操作 `insert("B")` 和 `insert("C")`，被智能地转换成了 `insert("B")` 和 `insert("C")` 在新文档上的正确版本，最终得到了所有人都满意的结果 `XABC`。
