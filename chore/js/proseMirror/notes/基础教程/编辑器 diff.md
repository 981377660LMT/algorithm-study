## 编辑器 diff

- json diff：JSON Patch：描述 JSON 文档变更的格式
  ![alt text](image-35.png)
  jsondiffpatch
- 纯文本 diff：diff-match-patch、Unidiff

---

好的，我们来系统地、深入地讲解“编辑器 diff”这个主题。这不仅仅是比较两个文本文件，更是现代编辑器，尤其是协同编辑和版本控制功能的核心基石。

我们将从以下几个层面来解构这个概念：

1.  **为什么需要 Diff：核心应用场景**
2.  **两大流派：基于状态的 Diff vs. 基于操作的 Diff**
3.  **技术详解：`diff-match-patch`, `Unidiff`, `JSON Patch`**
4.  **结构化文档的挑战与 ProseMirror 的选择**

---

### 1. 为什么需要 Diff：核心应用场景

“Diff”（Difference）的本质是**用一种结构化的方式，描述从状态 A 到状态 B 所需的变更**。在编辑器领域，它主要用于：

- **版本控制 (Version Control)**：像 Git 这样的系统，其核心就是存储和应用 diff（在 Git 中称为 patch），从而高效地记录每次提交的变更，而不是存储整个文件的无数个副本。
- **协同编辑 (Collaborative Editing)**：当多个用户同时编辑时，每个用户的修改都需要以 diff 的形式发送给其他人。其他人接收到 diff 后，将其应用（合并）到自己的文档版本上。
- **追踪修订 (Track Changes)**：在 Word 或 Google Docs 中，你看到的红色删除线、绿色下划线等，本质上就是将 diff 可视化地呈现在文档上，而不是直接应用它。
- **撤销/重做 (Undo/Redo)**：用户的每次操作都可以被看作一个 diff。撤销就是应用这个 diff 的“逆操作”，重做就是再次应用它。

---

### 2. 两大流派：基于状态的 Diff vs. 基于操作的 Diff

这是理解编辑器 diff 的最核心分野，它决定了 diff 算法的设计哲学。

#### A. 基于状态的 Diff (State-based Diff)

这种方法不关心“如何”从 A 变成 B，只关心“结果”。它通过直接比较两个完整的文档状态（快照），计算出它们之间的差异。

- **工作方式**：`diff(stateA, stateB) -> patch`
- **典型代表**：`diff-match-patch`, `Unidiff`
- **优点**：
  - **无状态**：算法本身不需要知道任何历史信息，给两个版本就能算，非常纯粹。
  - **健壮**：即使中间过程丢失，只要有最终状态，总能算出差异。
- **缺点**：
  - **意图丢失**：它无法区分用户的真实意图。例如，用户“删除一个段落再输入一个新段落”和“修改了原段落的每个字”，在状态 diff 看来可能结果完全一样。
  - **性能开销**：对于大文档，每次都进行全文比较，成本可能很高。
  - **冲突解决困难**：在协同编辑中，合并两个用户基于同一旧版本的 diff（即 `patch(patch(stateA, stateB1), stateB2)`）非常困难且容易出错，这就是著名的“三方合并”难题。

#### B. 基于操作的 Diff (Operation-based Diff)

这种方法记录的是用户完成变更所执行的**一系列具体操作（Operations）**。它关心的是“过程”而非“结果”。

- **工作方式**：监听用户事件，直接生成操作序列 `[op1, op2, ...]`。
- **典型代表**：**操作变换 (Operational Transformation, OT)**，如 ProseMirror 和 Google Docs 使用的；**无冲突复制数据类型 (CRDTs)**。
- **优点**：
  - **意图保留**：每个操作都精确地代表了用户的意图（例如 `updateStyle`, `splitNode`），信息量更丰富。
  - **轻量级**：操作通常非常小，网络传输开销低。
  - **为协同而生**：OT 和 CRDTs 的核心就是设计了一套数学模型来保证在并发操作下，不同客户端应用操作序列后能收敛到一致的状态。
- **缺点**：
  - **实现复杂**：需要实现复杂的转换（transform）或合并（merge）函数来处理并发操作的冲突。
  - **有状态**：通常需要维护版本号或操作历史，对实现要求更高。

---

### 3. 技术详解

现在我们来看你在笔记中提到的具体技术，它们都属于**基于状态的 Diff**。

#### `diff-match-patch` (by Google, Neil Fraser)

这是一个非常著名且强大的**纯文本**比较库。它包含三个核心部分：

1.  **Diff**:

    - **算法**: 它使用 Myers 差分算法及其变体，在效率和精度上做了很多优化。
    - **输出**: 生成一个由 `[操作, 文本]` 元组组成的数组。操作有三种：`DIFF_INSERT` (1), `DIFF_DELETE` (-1), `DIFF_EQUAL` (0)。
    - **示例**: `diff("Apple", "Crabapple")`
      - 输出: `[ [1, "Crab"], [0, "apple"] ]` (含义：插入 "Crab"，然后是公共部分 "apple")

2.  **Match**: 在文本中进行模糊匹配，即使目标文本有少量错误也能找到。
3.  **Patch**:
    - **功能**: 基于 Diff 的结果，生成一种可应用的补丁格式，并能将补丁应用到源文本上。
    - **特点**: 它的 Patch 格式包含了上下文信息，使得应用补丁时更加健壮。即使目标文本与生成补丁时的源文本不完全一样，只要上下文能对上，补丁依然有很大概率成功应用。

**局限性**: `diff-match-patch` 是**语义盲目**的。它只关心字符。对于富文本编辑器，`<b>Apple</b>` 变成 `<i>Apple</i>`，它可能会识别为删除了 `<b>` 和 `</b>`，并插入了 `<i>` 和 `</i>`，而无法理解这是“样式变更”。

#### `Unidiff` (Unified Diff Format)

这是一种**标准化的文本补丁格式**，而不是一个库。你每天在 `git diff` 或 GitHub Pull Requests 中看到的就是它。

- **特点**:
  - **面向行 (Line-oriented)**：它的最小比较单位是“行”。
  - **人类可读**: 格式清晰，易于阅读。
  - **包含上下文**: `@@ -1,5 +1,6 @@` 这样的头部信息定义了变更的起始行和影响行数。`+` 表示新增行，`-` 表示删除行，空格开头的行是未变的上下文。
- **示例**:
  ```diff
  --- a/file.txt
  +++ b/file.txt
  @@ -1,3 +1,4 @@
   Apple
  -Banana
  +Blueberry
   Cherry
  +Durian
  ```
  这个 diff 清晰地表示了：删除了 "Banana"，新增了 "Blueberry"，并在末尾新增了 "Durian"。

**局限性**: 因为是面向行的，它完全不适用于富文本。它无法表示一行中某个单词从普通变为粗体的操作。它主要用于代码和纯文本文档。

#### `JSON Patch` (RFC 6902)

这是一种用于描述 JSON 文档变更的**标准化格式**。

- **工作方式**: 它将 diff 表示为一个由多个“操作对象”组成的数组。每个对象描述一个原子操作。
- **核心操作**: `add`, `remove`, `replace`, `move`, `copy`, `test`。
- **路径表示**: 使用 JSON Pointer (RFC 6901) 来精确定位文档中的某个位置，例如 `/users/1/name`。
- **示例**:
  - 源 JSON: `{ "name": "Alice", "contact": { "email": "a@a.com" } }`
  - 目标 JSON: `{ "name": "Bob", "contact": { "email": "a@a.com", "phone": "123" } }`
  - JSON Patch:
    ```json
    [
      { "op": "replace", "path": "/name", "value": "Bob" },
      { "op": "add", "path": "/contact/phone", "value": "123" }
    ]
    ```

**与编辑器的关系**: 如果你的编辑器文档模型本身就是一个 JSON 结构（或者可以被看作 JSON），那么 JSON Patch 是一个非常好的**基于状态**的 diff 格式。它比纯文本 diff 更具结构性。

---

### 4. 结构化文档的挑战与 ProseMirror 的选择

富文本编辑器的文档不是纯文本，而是一个**结构化的树形文档**（Document Object Model）。

- **挑战**:

  1.  **语义**: “将两个段落合并”和“删除第一个段落结尾的换行符”在文本 diff 上可能一样，但在结构上是完全不同的操作。
  2.  **位置**: 当文档结构变化时，基于字符偏移量的位置会失效。例如，在文档开头加一个字，后面所有位置都需要更新。
  3.  **并发**: 状态 diff 难以解决并发冲突。

- **ProseMirror 的解决方案**:
  ProseMirror 彻底抛弃了基于状态的 diff 思想，全面拥抱**基于操作的 Diff (OT)**。
  1.  **`Step`**: 用户的每个最小操作（如 `addMark`, `replace`）都被抽象成一个可序列化、可逆的 `Step` 对象。这就是最原子的 diff。
  2.  **`Transaction`**: 一系列 `Step` 组成一个事务。
  3.  **`Transform`**: `Transaction` 的父类，核心是 `map` 方法。它定义了如何让一个位置或一个 `Step` 在另一个 `Step` 发生后进行“坐标变换”。
  4.  **`rebaseSteps`**: 正如我们之前深入分析的，`collab` 模块的核心就是利用 `map` 能力，智能地合并远程和本地的操作序列（`Step` 数组）。

**结论**: 对于像 ProseMirror 这样复杂的结构化编辑器，简单的文本 diff 或 JSON diff 无法满足其对语义、性能和协同能力的要求。因此，它选择了更复杂但功能更强大的**操作变换（OT）**作为其 diff 机制的理论基础，将 diff 的核心从“比较结果”转移到了“记录和变换过程”。

---

在前端开发中，“diff”通常指两类：一是用于比较文本内容差异的库，二是虚拟 DOM（Virtual DOM）中用于计算新旧 VNode 树差异的算法库。

以下是这两类中一些常用的库：

### 1. 文本差异比较 (Text Diff)

这类库常用于代码比对、版本控制、富文本编辑器等场景。

- **[jsdiff](https://github.com/kpdecker/jsdiff)**

  - 最流行和通用的文本比较库之一。
  - 基于 Myers 差分算法实现。
  - 可以比较行、单词、CSS、JSON 对象等，功能非常强大。

- **[diff-match-patch](https://github.com/google/diff-match-patch)**

  - 由 Google 开发和维护，非常健壮和高效。
  - 提供了 diff（差异）、match（匹配）和 patch（补丁）三种功能。
  - 被广泛应用于 Google Docs 等产品中。

- **[fast-diff](https://github.com/jhchen/fast-diff)**
  - 一个专注于性能的简单文本 diff 库，在处理大量文本时速度很快。

### 2. 虚拟 DOM 差异比较 (Virtual DOM Diff)

这类库是现代前端框架（如 React, Vue）的核心部分，但也有一些优秀的独立实现，可用于构建自己的框架或高性能 UI 组件。

- **[Snabbdom](https://github.com/snabbdom/snabbdom)**

  - 一个非常著名、模块化且性能极高的 Virtual DOM 库。
  - 它的设计思想深刻影响了 Vue 2.x 的 VDOM diff 算法。
  - 通过其丰富的模块和钩子（hooks）系统，可以轻松扩展功能。

- **[virtual-dom](https://github.com/Matt-Esch/virtual-dom)**

  - 可以说是 Virtual DOM 概念的早期和经典的实现之一，启发了很多后来的框架。
  - 虽然现在更新较少，但其历史地位和参考价值仍然很高。

- **[preact](https://preactjs.com/)**
  - 虽然 Preact 是一个轻量级的 React 替代品，但其内部的 diff 算法非常高效和简洁，可以作为学习和研究的对象。

通常情况下，如果你在使用 React 或 Vue 等框架，你不需要直接关心 VDOM diff 库，因为框架已经内置了。如果你需要做文本内容的比较，**`jsdiff`** 是最常见的选择。

当然，对于比较 JSON 对象的差异，也有许多优秀的库。这类库在 API 测试、数据同步、配置管理等场景中非常有用。

以下是一些常用的 JSON diff 库：

- **[deep-diff](https://github.com/flitbit/diff)**

  - 非常流行和强大的库，可以深入比较嵌套的对象和数组。
  - 它能详细地报告出变化的类型（新增、删除、编辑）、路径以及变化前后的值。
  - 输出是一个包含所有差异的数组，非常适合进行程序化处理。

- **[jsondiffpatch](https://github.com/benjamine/jsondiffpatch)**

  - 功能非常全面，不仅能找出差异（diff），还能生成一个“补丁”（patch）对象。
  - 这个补丁对象可以被用来更新原始 JSON，或者从新 JSON 还原回原始 JSON。
  - 它还提供了非常漂亮的 HTML 格式化工具，可以直观地展示差异，非常适合在前端页面上显示。

- **[json-diff](https://github.com/andreyvit/json-diff)**

  - 一个相对简单直接的库，专注于生成一个易于阅读的差异报告。
  - 它会输出一个带有特殊标记（如 `__old` 和 `__new`）的 JSON 对象来表示差异，或者彩色的文本输出。

- **[fast-json-patch](https://github.com/Starcounter-Jack/JSON-Patch)**
  - 这个库严格遵循 **JSON Patch (RFC 6902)** 标准。
  - 它不生成自定义的 diff 格式，而是生成一个标准的操作序列（如 `add`, `remove`, `replace`）。
  - 这使得它非常适合用于 API 通信，因为客户端和服务器可以都遵循同一个标准来交换数据变更。

### 总结

- 如果你需要**最详细的变更记录**来进行逻辑处理，推荐使用 **`deep-diff`**。
- 如果你需要**生成补丁**并应用，或者需要**漂亮的 HTML 可视化**，**`jsondiffpatch`** 是最佳选择。
- 如果你需要遵循 **RFC 6902 标准**来与后端或其他系统通信，应该使用 **`fast-json-patch`**。
