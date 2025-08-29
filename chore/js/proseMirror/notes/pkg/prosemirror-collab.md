prosemirror-collab 是一个 ProseMirror 插件，它为编辑器实现协同编辑功能提供了一个健壮的框架。它本身不处理网络通信，而是提供了一套机制来管理本地和远程的变更，使得将 ProseMirror 集成到任何协同编辑后端成为可能。

### 1. 核心思想：中央集权与乐观更新

ProseMirror 的协同模型基于一个**中央集权**的架构。这意味着有一个权威的中央服务（通常是你的服务器）负责接收所有客户端的变更，并确定它们最终的、唯一的顺序。

客户端则采用**乐观更新**的策略：用户的操作会立即在本地编辑器上生效，而不需要等待服务器的确认。当来自其他协作者的变更到达时，prosemirror-collab 模块负责将本地尚未被服务器确认的变更，在远程变更的基础上进行“变基”（Rebase），从而使所有客户端最终收敛到一致的状态。

### 2. 核心概念与 API

- **`collab(config)`**: 这是模块的主要入口，一个插件工厂函数。你需要将它添加到编辑器的插件列表中。它接收 `version` 和 `clientID` 作为配置。

  - `version`: 文档的初始版本号。
  - `clientID`: 当前客户端的唯一标识符。

- **`CollabState`**: 这是协同插件所管理的状态，包含两个关键部分：

  - `version`: 当前客户端已从中央权威处同步到的最新版本号。
  - `unconfirmed`: 一个 `Rebaseable` 对象的数组，代表了本地已应用但尚未被中央权威确认的变更步骤（`Step`）。

- **`sendableSteps(state)`**: 当你需要将本地变更发送到服务器时，调用此函数。它会返回一个包含 `version`、`steps` 和 `clientID` 的对象。如果返回 `null`，则表示没有需要发送的变更。你的应用代码负责将这个对象通过网络发送到中央服务器。

- **[`receiveTransaction(state, steps, clientIDs)`](prosemirror-collab/src/collab.ts)**: 当你的应用从服务器接收到一批新的变更时，调用此函数。它会返回一个新的 `Transaction`，你需要将这个事务应用到编辑器状态上。
  - `steps`: 从服务器接收到的 `Step` 数组。
  - `clientIDs`: 与 `steps` 数组一一对应的、产生这些步骤的客户端 ID 数组。

### 3. 工作流程与“变基”魔法

协同编辑的完整流程如下：

1.  **本地编辑**: 用户在编辑器中进行操作。ProseMirror 生成一个事务（`Transaction`），其中包含相应的 `Step`。`collab` 插件会将这些 `Step` 包装成 `Rebaseable` 对象，存入 `unconfirmed` 列表中。

2.  **发送变更**: 你的应用代码定期（或在每次编辑后）调用 `sendableSteps`，获取未确认的步骤，并将它们连同当前文档版本号发送到服务器。

3.  **服务器处理**: 服务器接收到客户端的变更。它会检查客户端提交的 `version` 是否与服务器上当前的最新版本匹配。

    - 如果匹配，服务器接受这些 `Step`，将它们应用到中央文档上，增加版本号，然后将这些 `Step` 广播给所有其他连接的客户端。
    - 如果不匹配，说明客户端的文档版本落后了，服务器会拒绝此次提交。

4.  **接收变更**: 客户端从服务器接收到一批广播来的 `Step`。此时，它调用 `receiveTransaction`。

5.  **变基 (Rebasing)**: 这是最关键的一步，在 `receiveTransaction` 内部发生。
    - 插件首先检查收到的 `steps` 中，哪些是自己之前发送并被服务器确认的。这些步骤会被从本地的 `unconfirmed` 列表中移除。
    - 剩下的 `steps` 是来自其他协作者的。
    - 此时，本地可能还有一些 `unconfirmed` 的步骤（例如，在发送和接收的间隙中，用户又进行了新的编辑）。
    - 插件会使用 `rebaseSteps` 函数执行变基操作：它会暂时**撤销**本地未确认的步骤，然后**应用**所有来自远程的步骤，最后再**重新应用**本地的步骤。在重新应用时，这些本地步骤会通过 prosemirror-transform 的 `map` 方法进行坐标变换，以适应远程变更所带来的文档结构变化。

这个过程确保了即使用户正在输入，远程的变更也能被平滑地插入，而用户的输入（未确认的变更）则会被智能地调整到正确的新位置上，最终所有客户端的文档内容都会收敛一致。

prosemirror-collab 与 `prosemirror-history` 也能很好地协同工作，确保在协同编辑环境下的撤销/重做行为符合预期。

---

好的，我们来深入探讨 `prosemirror-collab`。这个模块是 ProseMirror 官方提供的、用于实现协同编辑的客户端解决方案。理解它，不仅能让你学会如何构建一个协同编辑器，更能让你深刻体会到 `prosemirror-transform` 中 `Step` 和 `Mappable` 设计的精妙之处。

与 `prosemirror-history` 类似，`prosemirror-collab` 也是一个 `Plugin`，但它解决的问题要复杂得多：**如何在多个用户同时编辑时，保证所有人的文档最终能够收敛到一致的状态？**

我们将从以下四个关键角度来解构它：

1.  **核心模型：中央集权式（Centralized Authority）**
2.  **协同协议：客户端与服务器的“舞蹈”**
3.  **`rebase`（变基）：解决并发冲突的核心魔法**
4.  **API 与实现：`collab` 插件的内部运作**

---

### 1. 核心模型：中央集权式（Centralized Authority）

`prosemirror-collab` 采用的是一种经典的、基于**中央服务器**的协同编辑模型。这与 Yjs 等基于 CRDT 的去中心化模型不同。

在这个模型中：

- **服务器是唯一真理的来源（Single Source of Truth）**: 服务器维护着文档的“官方”版本号和一份权威的、线性的 `Step` 历史记录。
- **客户端的修改需要被确认**: 客户端提交的任何修改（`Step`）都必须被服务器接收、确认并整合到官方历史中，才算作最终生效。
- **所有变更由服务器广播**: 一旦服务器接受了一组 `Step`，它会将这些 `Step` 广播给所有连接的客户端（包括提交者自己）。

这个模型就像一个有严格管理员的 Git 仓库：你不能直接 `push`，你必须提交一个 Pull Request（发送你的 `Step`），管理员（服务器）在确认无误后将其合并（`merge`）到主分支，然后通知所有人去 `pull` 最新的主分支。

---

### 2. 协同协议：客户端与服务器的“舞蹈”

`prosemirror-collab` 本身**只实现了客户端的逻辑**。它定义了一套协议，你需要自己实现一个遵循此协议的服务器。这个协议的交互流程（“舞蹈”）如下：

#### a. 客户端状态

每个客户端需要维护两个关键状态：

- **`version`**: 我所知道的、来自服务器的最新文档版本号。
- **`unconfirmedSteps: Step[]`**: 我已经在我本地应用、但尚未被服务器确认的 `Step` 序列。

#### b. 交互流程

1.  **初始连接**: 客户端连接到服务器，请求文档的当前版本号和所有历史 `Step`，然后在本地构建出最新的 `EditorState`。

2.  **客户端做出修改**:

    - 用户在编辑器中进行操作，产生一个 `Transaction`。
    - `prosemirror-collab` 插件捕获这个事务，将其中的 `Step` 存储在本地的 `unconfirmedSteps` 数组中。
    - 客户端**立即**在本地应用这些 `Step`，为用户提供即时反馈。
    - 客户端将 `{ version, steps, clientID }` 发送给服务器，尝试提交这些变更。

3.  **服务器处理提交**:

    - **情况 A (无冲突)**: 服务器收到客户端的提交，发现其 `version` 与服务器当前的 `version` 一致。
      - 服务器接受这些 `Step`。
      - 服务器将这些 `Step` 应用到自己的权威文档上。
      - 服务器将自己的 `version` 加一。
      - 服务器将 `{ steps, newVersion }` **广播**给所有客户端。
    - **情况 B (冲突)**: 服务器收到客户端的提交，但发现其 `version` **落后于**服务器当前的 `version`。这意味着在客户端提交期间，已经有其他人的修改被服务器接受了。
      - 服务器**拒绝**这次提交，并向该客户端返回一个“冲突”或“过时”的响应（通常是 HTTP 409 Conflict）。

4.  **客户端处理广播**:

    - 当客户端收到来自服务器的广播（包含新 `Step` 和新 `version`）时：
      - **如果提交者是自己**: 客户端发现广播来的 `Step` 就是自己之前提交的，它就知道自己的修改已被确认。于是，它从 `unconfirmedSteps` 中移除这些已确认的 `Step`，并更新自己的 `version`。
      - **如果提交者是别人**: 这是最关键的一步。客户端不能直接应用这些新来的 `Step`，因为它自己还有 `unconfirmedSteps`。直接应用会导致数据错乱。此时，客户端必须执行 **`rebase`（变基）** 操作（见下一节）。

5.  **客户端处理被拒**:
    - 如果客户端的提交被服务器拒绝（情况 B），它什么也不做，只是静静等待服务器关于新版本的广播。一旦收到广播，它就会执行 `rebase`，然后用变基后的新 `Step` 和新 `version` 重新尝试提交。

---

### 3. `rebase`（变基）：解决并发冲突的核心魔法

`rebase` 是整个协同算法的核心，它完美地利用了 `prosemirror-transform` 的 `Mappable` 接口。

**场景**:

- 你的本地文档是 `"A"`。
- 你输入了 `"B"`，现在本地是 `"AB"`。你的 `unconfirmedSteps` 是 `[insert(1, "B")]`。
- 这时，你收到了来自服务器的广播：别人在你之前插入了 `"X"`，操作是 `[insert(0, "X")]`。

**`rebase` 的过程**:

1.  **回滚本地未确认的修改**: 客户端先临时撤销掉自己的 `unconfirmedSteps` (`[insert(1, "B")]`)，文档状态回到 `"A"`。
2.  **应用来自服务器的权威修改**: 将服务器广播来的 `Step` (`[insert(0, "X")]`) 应用到文档上。文档现在变成了 `"XA"`。
3.  **转换（Transform）自己的未确认修改**: 这是最关键的一步。客户端需要计算出：我原来的操作 `insert(1, "B")`，在一个 `insert(0, "X")` 发生之后，应该变成什么？
    - 它使用 `step.map(mapping)` 方法。这里的 `mapping` 就是 `insert(0, "X")` 这个 `Step` 自身。
    - `insert(1, "B")` 中的位置 `1`，经过 `insert(0, "X")` 的映射后，向后移动了一位，变成了 `2`。
    - 所以，原来的 `Step` 被**转换**成了一个新的 `Step`：`insert(2, "B")`。
4.  **应用转换后的修改**: 将这个新的 `Step` (`insert(2, "B")`) 应用到当前文档 `"XA"` 上，得到最终结果 `"XAB"`。
5.  **更新状态**: 客户端的 `unconfirmedSteps` 现在被更新为转换后的 `Step` (`[insert(2, "B")]`)。文档内容是 `"XAB"`。

现在，客户端的本地状态与一个先执行 `insert(0, "X")` 再执行 `insert(1, "B")` 的结果完全一致，数据收敛了！之后，它会用这个新的 `unconfirmedSteps` 再次尝试向服务器提交。

---

### 4. API 与实现：`collab` 插件的内部运作

`prosemirror-collab` 提供了实现上述客户端逻辑所需的所有工具。

- **`collab({ version, clientID })`**: 这是创建协同插件的工厂函数。

  - `version`: 初始的文档版本号。
  - `clientID`: 当前客户端的唯一标识符，用于区分是自己的提交还是他人的提交。

- **`sendableSteps(state)`**:

  - 这是一个辅助函数，用于检查当前状态，看是否有可以发送给服务器的未确认 `Step`。
  - 如果 `state.plugins` 中 `collab` 插件的内部状态表明有未确认的 `Step`，它会返回一个对象：`{ version, steps, clientID }`。否则返回 `null`。
  - 你的网络层代码需要定期调用这个函数，并将返回的数据发送给服务器。

- **`receiveTransaction(state, steps, clientIDs)`**:
  - 这是你的网络层在收到服务器广播后需要调用的核心函数。
  - 它接收当前的 `state`、服务器广播的 `steps` 数组以及这些 `steps` 对应的 `clientIDs` 数组。
  - 它会为你生成一个新的 `Transaction`。这个 `Transaction` 内部已经处理了所有的 `rebase` 逻辑。你只需要将这个 `Transaction` 应用到你的 `EditorView` 即可。

### 总结

`prosemirror-collab` 是一个精炼的、客户端侧的 OT（操作转换）实现。它本身不是一个完整的解决方案，而是一个**框架**和**协议**。它要求开发者自己实现服务器端逻辑和网络通信层。

它的优点是逻辑清晰、基于 ProseMirror 核心概念构建、可靠性高。它的缺点是需要一个中央服务器，并且对网络延迟比较敏感。

在现代前端开发中，由于 Yjs 等 CRDT 方案的成熟（它们能更好地处理离线编辑、网络分区，并且不强制要求中央服务器），直接使用 `prosemirror-collab` 的场景相对减少了。但学习它的原理，是理解 OT 算法和 ProseMirror 强大转换能力的最佳途径之一。

---

```ts
import { Plugin, PluginKey, TextSelection, EditorState, Transaction } from 'prosemirror-state'
import { Step, Transform } from 'prosemirror-transform'

class Rebaseable {
  constructor(readonly step: Step, readonly inverted: Step, readonly origin: Transform) {}
}

/// Undo a given set of steps, apply a set of other steps, and then
/// redo them @internal
export function rebaseSteps(
  steps: readonly Rebaseable[],
  over: readonly Step[],
  transform: Transform
) {
  for (let i = steps.length - 1; i >= 0; i--) transform.step(steps[i].inverted)
  for (let i = 0; i < over.length; i++) transform.step(over[i])
  let result = []
  for (let i = 0, mapFrom = steps.length; i < steps.length; i++) {
    let mapped = steps[i].step.map(transform.mapping.slice(mapFrom))
    mapFrom--
    if (mapped && !transform.maybeStep(mapped).failed) {
      transform.mapping.setMirror(mapFrom, transform.steps.length - 1)
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

// This state field accumulates changes that have to be sent to the
// central authority in the collaborating group and makes it possible
// to integrate changes made by peers into our local document. It is
// defined by the plugin, and will be available as the `collab` field
// in the resulting editor state.
class CollabState {
  constructor(
    // The version number of the last update received from the central
    // authority. Starts at 0 or the value of the `version` property
    // in the option object, for the editor's value when the option
    // was enabled.
    readonly version: number,
    // The local steps that havent been successfully sent to the
    // server yet.
    readonly unconfirmed: readonly Rebaseable[]
  ) {}
}

function unconfirmedFrom(transform: Transform) {
  let result = []
  for (let i = 0; i < transform.steps.length; i++)
    result.push(
      new Rebaseable(transform.steps[i], transform.steps[i].invert(transform.docs[i]), transform)
    )
  return result
}

const collabKey = new PluginKey('collab')

type CollabConfig = {
  /// The starting version number of the collaborative editing.
  /// Defaults to 0.
  version?: number

  /// This client's ID, used to distinguish its changes from those of
  /// other clients. Defaults to a random 32-bit number.
  clientID?: number | string
}

/// Creates a plugin that enables the collaborative editing framework
/// for the editor.
export function collab(config: CollabConfig = {}): Plugin {
  let conf: Required<CollabConfig> = {
    version: config.version || 0,
    clientID: config.clientID == null ? Math.floor(Math.random() * 0xffffffff) : config.clientID
  }

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

    // This is used to notify the history plugin to not merge steps,
    // so that the history can be rebased.
    historyPreserveItems: true
  })
}

/// Create a transaction that represents a set of new steps received from
/// the authority. Applying this transaction moves the state forward to
/// adjust to the authority's view of the document.
export function receiveTransaction(
  state: EditorState,
  steps: readonly Step[],
  clientIDs: readonly (string | number)[],
  options: {
    /// When enabled (the default is `false`), if the current
    /// selection is a [text selection](#state.TextSelection), its
    /// sides are mapped with a negative bias for this transaction, so
    /// that content inserted at the cursor ends up after the cursor.
    /// Users usually prefer this, but it isn't done by default for
    /// reasons of backwards compatibility.
    mapSelectionBackward?: boolean
  } = {}
) {
  // Pushes a set of steps (received from the central authority) into
  // the editor state (which should have the collab plugin enabled).
  // Will recognize its own changes, and confirm unconfirmed steps as
  // appropriate. Remaining unconfirmed steps will be rebased over
  // remote steps.
  let collabState = collabKey.getState(state)
  let version = collabState.version + steps.length
  let ourID: string | number = (collabKey.get(state)!.spec as any).config.clientID

  // Find out which prefix of the steps originated with us
  let ours = 0
  while (ours < clientIDs.length && clientIDs[ours] == ourID) ++ours
  let unconfirmed = collabState.unconfirmed.slice(ours)
  steps = ours ? steps.slice(ours) : steps

  // If all steps originated with us, we're done.
  if (!steps.length) return state.tr.setMeta(collabKey, new CollabState(version, unconfirmed))

  let nUnconfirmed = unconfirmed.length
  let tr = state.tr
  if (nUnconfirmed) {
    unconfirmed = rebaseSteps(unconfirmed, steps, tr)
  } else {
    for (let i = 0; i < steps.length; i++) tr.step(steps[i])
    unconfirmed = []
  }

  let newCollabState = new CollabState(version, unconfirmed)
  if (options && options.mapSelectionBackward && state.selection instanceof TextSelection) {
    tr.setSelection(
      TextSelection.between(
        tr.doc.resolve(tr.mapping.map(state.selection.anchor, -1)),
        tr.doc.resolve(tr.mapping.map(state.selection.head, -1)),
        -1
      )
    )
    ;(tr as any).updated &= ~1
  }
  return tr
    .setMeta('rebased', nUnconfirmed)
    .setMeta('addToHistory', false)
    .setMeta(collabKey, newCollabState)
}

/// Provides data describing the editor's unconfirmed steps, which need
/// to be sent to the central authority. Returns null when there is
/// nothing to send.
///
/// `origins` holds the _original_ transactions that produced each
/// steps. This can be useful for looking up time stamps and other
/// metadata for the steps, but note that the steps may have been
/// rebased, whereas the origin transactions are still the old,
/// unchanged objects.
export function sendableSteps(state: EditorState): {
  version: number
  steps: readonly Step[]
  clientID: number | string
  origins: readonly Transaction[]
} | null {
  let collabState = collabKey.getState(state) as CollabState
  if (collabState.unconfirmed.length == 0) return null
  return {
    version: collabState.version,
    steps: collabState.unconfirmed.map(s => s.step),
    clientID: (collabKey.get(state)!.spec as any).config.clientID,
    get origins() {
      return (
        (this as any)._origins ||
        ((this as any)._origins = collabState.unconfirmed.map(s => s.origin))
      )
    }
  }
}

/// Get the version up to which the collab plugin has synced with the
/// central authority.
export function getVersion(state: EditorState): number {
  return collabKey.getState(state).version
}
```
