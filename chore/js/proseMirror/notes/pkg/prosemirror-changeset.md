好的，我们来对 `prosemirror-changeset` 进行一次彻底的、深入的剖析。这个库虽然不是 ProseMirror 的核心模块之一，但它是由 ProseMirror 的作者 Marijn Haverbeke 亲自编写的，用于解决一个非常具体且重要的问题：**如何高效、结构化地比较两个 ProseMirror 文档，并找出它们之间的差异？**

它常被用于实现“追踪修订”（Track Changes）或“差异对比”（Diffing）等高级功能。

我们将从以下四个关键角度来解构它：

1.  **核心问题：超越文本 Diff**
2.  **`ChangeSet` 对象：差异的结构化描述**
3.  **算法揭秘：它如何工作？**
4.  **实践应用：构建一个“追踪修订”视图**

---

### 1. 核心问题：超越文本 Diff

为什么我们不能简单地将两个 ProseMirror 文档序列化成 HTML 或文本，然后用标准的文本 diff 工具（如 `diff-match-patch`）来比较呢？

答案在于**结构（Structure）**。

- **文本 Diff 的局限**: 文本 diff 工具将所有内容视为一维的字符流。如果你将一个段落的样式从 `paragraph` 改为 `heading`，文本 diff 可能会认为你删除了整行旧文本，又添加了整行新文本，因为它不理解“块级节点类型变更”这个概念。
- **`prosemirror-changeset` 的优势**: 它是一个**结构化 diff 工具**。它直接在 ProseMirror 的 `Node` 树上进行操作，因此它能理解：
  - 节点属性的变化（如 `heading` 的 `level` 改变）。
  - 节点类型的变化（如 `paragraph` 变成 `code_block`）。
  - 标记（Mark）的添加或移除。
  - 文本内容的增删。

它产生的差异描述，远比“第 X 行有变动”要精确和丰富得多。

---

### 2. `ChangeSet` 对象：差异的结构化描述

`prosemirror-changeset` 的核心是 `ChangeSet` 类。你通过比较两个文档来创建一个实例：

```typescript
import { ChangeSet } from "prosemirror-changeset";
import { Node } from "prosemirror-model";

const oldDoc: Node = /* ... */;
const newDoc: Node = /* ... */;

// 计算差异
const changeset: ChangeSet = ChangeSet.create(oldDoc, newDoc);
```

这个 `changeset` 对象本身不包含文档内容，它只包含一个描述差异的核心数据结构：**一个 `Span` 对象的数组**。

#### `Span` 对象：差异的原子单元

`Span` 描述了**新文档（`newDoc`）**中的一段连续范围，以及这段范围与**旧文档（`oldDoc`）**的关系。

一个 `Span` 对象有三个关键属性：

- **`length: number`**: 这个 `Span` 在**新文档**中覆盖的长度。
- **`data: any`**: 这是一个标记，用于指示这个 `Span` 的类型。它有三种可能的值：
  - **`null`**: 表示这个范围内的内容是**新增的（Inserted）**。它在旧文档中不存在。
  - **`true`**: 表示这个范围内的内容是**未改变的（Unchanged）**。它在旧文档和新文档中完全相同。
  - **一个 `Change` 对象**: 表示这个范围内的内容是**已改变的（Changed）**。这个 `Change` 对象本身包含了更详细的关于如何从旧范围变成新范围的信息。

`ChangeSet` 的 `spans` 数组将整个新文档完美地分割成一系列连续的、不重叠的 `Span`。

#### `Change` 对象：深入“已改变”的细节

当一个 `Span` 的 `data` 是一个 `Change` 对象时，表示这部分内容发生了变化。这个 `Change` 对象提供了更深层的信息：

- `fromA: number`, `toA: number`: 这段内容在**旧文档**中的范围。
- `fromB: number`, `toB: number`: 这段内容在**新文档**中的范围。
- `deleted: Slice`: 从旧文档中被删除的内容片段。
- `inserted: Slice`: 在新文档中被插入的内容片段。

通过检查 `deleted` 和 `inserted`，你可以精确地知道是文本变了，还是标记变了，或是节点属性变了。

---

### 3. 算法揭秘：它如何工作？

`prosemirror-changeset` 使用了一种高效的、针对树状结构的 diff 算法。其基本思想是：

1.  **逐层比较**: 算法从两个文档的根节点开始，递归地向下比较。
2.  **寻找共同前缀和后缀**: 在比较任意两个节点的内容（`Fragment`）时，它首先快速扫描并跳过开头和结尾完全相同的部分。这些部分会被标记为**未改变的（`data: true`）**。
3.  **处理中间差异**: 对于中间不同的部分，算法会更深入地分析。它会尝试将一个节点的变更（如属性变化）与另一个节点的增删匹配起来，以找到最小的、最符合逻辑的差异描述。
4.  **利用 `Node.eq`**: 它大量使用 ProseMirror `Node` 对象自带的 `.eq()` 方法来快速判断两个节点是否完全相同，这比深度比较对象要快得多。
5.  **优化**: 算法经过了优化，可以快速处理长文档，避免在大型未更改区域浪费计算资源。

最终，它生成一个线性的 `Span` 数组，这个数组是对新文档的一种“着色”，精确地标记了每一部分是新增、未变还是已改。

---

### 4. 实践应用：构建一个“追踪修订”视图

`prosemirror-changeset` 最常见的用途就是创建一个类似 Google Docs 或 Word 的“追踪修订”视图。这通常通过 ProseMirror 的 `Decoration` 来实现。

**实现步骤**:

1.  **获取文档**: 你需要有 `oldDoc` 和 `newDoc` 两个版本的文档。

2.  **计算 `ChangeSet`**:

    ```typescript
    const changeset = ChangeSet.create(oldDoc, newDoc)
    ```

3.  **创建 `Decoration`**: 遍历 `changeset` 的 `spans` 和被删除的部分，为它们创建 `Decoration`。

    ```typescript
    import { Decoration, DecorationSet } from 'prosemirror-view'

    const decorations: Decoration[] = []
    let newDocPos = 0

    // 1. 为新增和修改的部分创建装饰器
    for (const span of changeset.spans) {
      if (span.data === null) {
        // 新增
        decorations.push(
          Decoration.inline(newDocPos, newDocPos + span.length, { class: 'insertion' })
        )
      } else if (span.data instanceof Change) {
        // 修改
        // 可以进一步分析 Change 对象，但简单起见，也标记为修改
        decorations.push(Decoration.inline(newDocPos, newDocPos + span.length, { class: 'change' }))
      }
      newDocPos += span.length
    }

    // 2. 为删除的部分创建“小部件”装饰器
    // changeset.deleted 包含了所有被删除的片段及其在旧文档中的位置
    // 我们需要将这些旧位置映射到新文档中，以决定在哪里显示“删除线”
    changeset.deleted.forEach(del => {
      // mapPos(pos) 将旧文档位置映射到新文档
      const pos = changeset.mapPos(del.pos)
      if (pos != null) {
        // 创建一个 Widget Decoration，它会在指定位置插入一个 DOM 元素
        // 这个 DOM 元素可以是一个带有删除线的、不可编辑的 span
        decorations.push(Decoration.widget(pos, createDeletedNodeView(del.slice)))
      }
    })

    const decorationSet = DecorationSet.create(newDoc, decorations)
    ```

4.  **应用 `Decoration`**: 将这个 `decorationSet` 通过插件提供给 `EditorView`。

    ```typescript
    new Plugin({
      state: {
        init() {
          return decorationSet
        },
        apply(tr) {
          // 当文档变化时，需要重新计算并更新 DecorationSet
          return decorationSet.map(tr.mapping, tr.doc)
        }
      },
      props: {
        decorations(state) {
          return this.getState(state)
        }
      }
    })
    ```

5.  **添加 CSS**:
    ```css
    .insertion {
      background-color: #e6ffed;
    }
    .change {
      background-color: #e6f7ff;
    }
    .deleted-widget {
      text-decoration: line-through;
      color: #f5222d;
    }
    ```

### 总结

`prosemirror-changeset` 是一个强大而专一的工具。它不是用于实时协同编辑的（那是 `prosemirror-transform` 和 OT/CRDT 的领域），而是用于**异步的、事后的文档比较**。它通过提供一个精确、结构化的差异描述（`ChangeSet`），使得构建复杂的、对用户友好的“追踪修订”和“差异对比”功能成为可能，完美地弥补了传统文本 diff 工具在处理富文本时的不足。
