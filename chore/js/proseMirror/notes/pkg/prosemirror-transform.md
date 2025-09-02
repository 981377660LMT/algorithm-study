

好的，我们来深入 `prosemirror-transform` 的内部世界。如果说 `prosemirror-model` 是静态的骨架，那么 `prosemirror-transform` 就是驱动骨架运动的动态肌肉和肌腱。它是 ProseMirror 中最具数学美感和工程挑战性的部分，也是实现撤销/重做和协同编辑等高级功能的核心。

我们将从以下四个方面来解构它：

1.  **`Step`：变革的原子**
2.  **`Transform`：构建变更序列的工厂**
3.  **`Mappable` 与位置映射：在变化中保持同步的“GPS”**
4.  **可逆性与操作转换（OT）：高级功能的基础**

---

### 1. `Step`：变革的原子

`Step` 是对文档的一次**最小、原子性、可序列化**的修改。任何复杂的编辑操作，最终都会被分解成一个或多个 `Step`。

#### a. `Step` 的设计哲学

- **原子性**: 一个 `Step` 要么完全应用，要么完全不应用，不存在中间状态。
- **可序列化**: `Step` 可以被转换成 JSON (`step.toJSON()`)，这对于通过网络发送或存储在数据库中至关重要。
- **可逆性**: 每个 `Step` 都有一个可计算的逆操作 (`step.invert(doc)`)，这是实现撤销功能的基础。

#### b. 核心 `Step` 类型

最核心的 `Step` 是 `ReplaceStep`。事实上，大部分编辑操作（插入、删除）都可以被看作是 `ReplaceStep` 的特例。

- **`ReplaceStep`**:
  - **定义**: `new ReplaceStep(from, to, slice)`
  - **作用**: 将文档中从 `from` 到 `to` 的内容替换为给定的 `Slice`。
  - **插入**: `new ReplaceStep(pos, pos, slice)` (用一个非空切片替换一个空范围)。
  - **删除**: `new ReplaceStep(from, to, Slice.empty)` (用一个空切片替换一个非空范围)。

其他 `Step` 类型主要用于处理标记（Marks），因为直接用 `ReplaceStep` 替换带有不同标记的相同文本会很低效。

- `AddMarkStep(from, to, mark)`: 在指定范围内添加一个标记。
- `RemoveMarkStep(from, to, mark)`: 在指定范围内移除一个标记。

---

### 2. `Transform`：构建变更序列的工厂

`Transform` 对象是开发者与之交互的主要接口。它是一个“变更构建器”，始于一个初始文档，然后通过一系列方法调用来累积 `Step`。

#### a. `Transform` 的工作流程

1.  **创建**: `const tr = new Transform(doc);` (在 `prosemirror-state` 中，我们通常使用 `state.tr`，它会创建一个继承自 `Transform` 的 `Transaction` 对象)。
2.  **构建**: 调用 `tr` 的方法来描述变更，例如 `tr.insertText('!', 10)` 或 `tr.delete(5, 7)`。这些方法在内部会创建相应的 `Step` 并添加到 `tr` 中。
3.  **获取结果**:
    - `tr.doc`: 获取应用了所有 `Step` 之后**新的**文档对象。
    - `tr.steps`: 获取累积的 `Step` 数组。
    - `tr.docChanged`: 一个布尔值，表示文档是否真的发生了改变。

```typescript
import { Transform } from 'prosemirror-transform'
import { schema, doc, p } from 'prosemirror-test-builder' // 假设的测试工具

const myDoc = doc(p('hello')) // 初始文档

// 1. 创建一个 Transform
const tr = new Transform(myDoc)

// 2. 构建变更
tr.insertText(' world', 6) // 在 'hello' 后面插入
tr.addMark(1, 6, schema.marks.strong.create()) // 将 'hello' 加粗

// 3. 查看结果
console.log(tr.steps.length) // 2 (一个 ReplaceStep, 一个 AddMarkStep)
console.log(tr.doc.toString()) // doc(p(strong("hello"), " world")) -> 新的文档
console.log(myDoc.toString()) // doc(p("hello")) -> 原始文档保持不变
```

#### b. `Transaction` vs `Transform`

在实际应用中，我们几乎总是使用 `prosemirror-state` 提供的 `Transaction` 对象。`Transaction` **继承**自 `Transform`，并在其基础上增加了管理编辑器状态（如选区、插件元数据、滚动位置）的能力。但所有核心的文档修改逻辑都源于 `Transform`。

---

### 3. `Mappable` 与位置映射：在变化中保持同步的“GPS”

这是 `prosemirror-transform` 最强大也最精妙的部分。当文档内容发生变化时，文档中的位置（positions）也会随之移动。**位置映射**就是计算一个旧文档中的位置在新文档中对应位置的过程。

#### a. `Mappable` 接口

任何代表了一系列文档变化的对象，如果能提供位置映射的能力，就实现了 `Mappable` 接口。`Step` 和 `Transform` 都实现了这个接口。

核心方法是 `map(pos: number, assoc?: number): number`。

- `pos`: 旧文档中的位置。
- `assoc` (association): 一个值为 `-1` 或 `1` 的数字，用于处理边界情况。当在一个位置 `pos` 插入内容时，`pos` 这个点本身会分裂成插入内容的“前”和“后”。`assoc: -1` 表示关联到插入点的前面（位置不变），`assoc: 1` 表示关联到插入点的后面（位置向后移动）。默认是 `1`。

#### b. 为什么位置映射至关重要？

想象一下协同编辑的场景：

1.  你的文档是 `"ab"`。你的光标在位置 `2` (b 后面)。
2.  你输入了 `"c"`。你的本地 `Transform` 是 `insert(2, "c")`。你的新文档是 `"abc"`，新光标位置是 `3`。
3.  与此同时，你的同事 Alice 在她那边（也是 `"ab"`）的**位置 `1`** 插入了 `"X"`。她的 `Transform` 是 `insert(1, "X")`。
4.  你收到了 Alice 的 `Transform`。你不能直接把它应用到你当前的文档 `"abc"` 上，因为她的操作是基于 `"ab"` 的。直接应用 `insert(1, "X")` 会得到 `"Xabc"`，这是错误的！正确结果应该是 `"aXbc"`。

你需要做的是：将 Alice 的 `Transform` 在你的 `Transform` **之后**进行“变基”（rebase）。这意味着你需要计算出 Alice 的 `insert(1, "X")` 在你的文档从 `"ab"` 变成 `"abc"` 之后，应该等价于什么操作。

- Alice 的操作位置 `1`，经过你的 `insert(2, "c")` 变换后，位置 `1` 没有受到影响，所以它**映射后仍然是 `1`**。
- 所以，Alice 的操作在你当前的文档 `"abc"` 上应该被转换为 `insert(1, "X")`。应用后得到 `"aXbc"`。

反过来，当 Alice 收到你的 `insert(2, "c")` 时，她也需要做同样的事情：

- 你的操作位置 `2`，经过她的 `insert(1, "X")` 变换后，位置 `2` 向后移动了一位，**映射后变成了 `3`**。
- 所以，你的操作在她的文档 `"aXb"` 上应该被转换为 `insert(3, "c")`。应用后得到 `"aXbc"`。

最终，你们俩都得到了相同的正确结果 `"aXbc"`。这个过程就是**操作转换 (Operational Transformation, OT)** 的核心，而它完全依赖于 `Mappable` 提供的精确位置映射能力。

---

### 4. 可逆性与操作转换（OT）

`prosemirror-transform` 提供了实现这些高级功能所需的所有底层工具。

- **可逆性 (`step.invert(doc)`)**:
  每个 `Step` 都能根据它所应用的文档，生成一个“逆操作” `Step`。例如，`insert(pos, text)` 的逆操作是 `delete(pos, pos + text.length)`。`prosemirror-history` 插件就是通过存储这些逆操作来实现撤销功能的。

- **操作转换 (OT)**:
  `Transform.prototype.rebase` 和 `StepMap.map` 等方法提供了实现 OT 算法所需的核心构件。虽然 ProseMirror 本身不包含一个完整的、开箱即用的 OT 服务端，但像 `prosemirror-collab` 模块就利用了这些工具来提供一个客户端的 OT 实现框架，开发者可以基于此构建自己的协同编辑后端。

### 总结

`prosemirror-transform` 是一个纯粹的、功能强大的文档变更描述与计算引擎。它通过**原子化的 `Step`**、**链式的 `Transform` 构建器**以及**精确的 `Mappable` 位置映射系统**，为上层功能（如状态管理、视图更新、历史记录、协同编辑）提供了一个极其稳固和可靠的基础。它将编辑器中混乱、易错的变更管理，抽象成了一套优雅、可预测的数学模型。
