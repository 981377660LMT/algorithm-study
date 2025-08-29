好的，我们来对 `prosemirror-model` 进行一次彻底的、深入的剖析。这个包是整个 ProseMirror 生态系统的基石，理解它的设计哲学和核心类型，是精通 ProseMirror 的必经之路。

我们将从四个关键角度来解构它：

1.  **Schema：文档的类型系统与规则引擎**
2.  **Node & Fragment：不可变的文档树**
3.  **Slice：智能剪贴板的秘密**
4.  **JSON 序列化：持久化与数据交换**

---

### 1. Schema：文档的类型系统与规则引擎

`Schema` 实例是 `prosemirror-model` 的灵魂。它不仅仅是一个配置对象，更是一个编译和实例化后的**规则引擎**，定义了文档的一切可能性。

#### a. 节点 (`NodeType`) 和标记 (`MarkType`)

Schema 由 `nodes` 和 `marks` 两部分构成。

- `NodeType`: 定义了文档的结构块。每个 `NodeType` 实例都包含了它的所有元信息。
- `MarkType`: 定义了可以应用在内联内容上的样式或元数据。

让我们看一个带有详细注解的 `paragraph` 节点规格 (`NodeSpec`)：

```typescript
// 定义一个段落节点的规格
const paragraphSpec: NodeSpec = {
  // 1. 内容表达式 (Content Expression)
  // 定义了此节点可以包含哪些子节点。
  // "inline*" 意味着可以包含 0 个或多个属于 "inline" 组的节点。
  content: 'inline*',

  // 2. 分组 (Group)
  // 将节点归类，便于在内容表达式中引用。
  // 这个段落属于 "block" 组。
  group: 'block',

  // 3. 属性 (Attributes)
  // 定义了节点可以拥有的属性及其默认值。
  attrs: {
    align: { default: null } // 例如，一个对齐属性
  },

  // 4. DOM 解析规则 (Parsing from DOM)
  // 定义了如何将一个 DOM 元素解析成这个 Node。
  parseDOM: [
    {
      tag: 'p', // 匹配 <p> 标签
      getAttrs(dom: HTMLElement): object {
        // 从 DOM 元素上提取属性值
        return { align: dom.style.textAlign || null }
      }
    }
  ],

  // 5. DOM 序列化规则 (Serializing to DOM)
  // 定义了如何将这个 Node 渲染成一个 DOM 元素。
  toDOM(node: Node): DOMOutputSpec {
    const { align } = node.attrs
    // DOMOutputSpec 是一个数组，描述了 DOM 结构。
    // ['p', { style: '...' }, 0]
    // 0 表示 "内容应该渲染在这里"
    return ['p', { style: align ? `text-align: ${align}` : '' }, 0]
  }
}
```

**内容表达式的威力**：这是 Schema 最强大的特性之一。它使用一种类似正则表达式的语法来保证文档结构的合法性，从根本上杜绝了无效的文档状态。

- `"heading paragraph*"`: 一个标题，后面跟着零个或多个段落。
- `"list_item+"`: 一个或多个列表项。
- `"(ordered_list | bullet_list)"`: 一个有序列表或一个无序列表。

#### b. Schema 实例化

你定义的 `NodeSpec` 和 `MarkSpec` 对象最终会被编译成一个 `Schema` 实例。

```typescript
import { Schema } from 'prosemirror-model'

const mySchema = new Schema({
  nodes: {
    doc: { content: 'block+' },
    paragraph: paragraphSpec,
    text: { group: 'inline' }
    // ... 其他节点
  },
  marks: {
    // ... 标记定义
  }
})
```

这个 `mySchema` 对象现在就是一个功能完备的工厂和验证器。你可以用它来创建节点 (`mySchema.nodes.paragraph.create(...)`)，或者从 HTML 解析文档 (`DOMParser.fromSchema(mySchema).parse(...)`)。

---

### 2. Node & Fragment：不可变的文档树

- **`Node`**: 代表文档树中的一个节点。**它是不可变的 (Immutable)**。
- **`Fragment`**: 代表一个节点的子节点序列。它也是不可变的。

**为什么不可变性如此重要？**

1.  **可预测性**: 当你有一个 `Node` 对象时，你可以确信它不会在别处被意外修改。
2.  **高效的状态比较**: 在 `prosemirror-state` 中，判断文档是否改变只需要比较 `oldState.doc === newState.doc`，这是一个极快的引用比较。
3.  **可靠的历史记录**: 撤销操作只是将状态指回旧的 `doc` 对象，简单而可靠。
4.  **协同编辑的基础**: 不可变性使得追踪和转换（transform）变更步骤成为可能。

一个 `Node` 实例的核心属性：

- `type: NodeType`: 指向其在 Schema 中的类型定义。
- `attrs: object`: 节点的属性。
- `content: Fragment`: 包含其所有子节点的 Fragment。
- `marks: readonly Mark[]`: 应用在该节点上的标记集合（主要用于文本节点）。

**操作文档**:
由于 `Node` 是不可变的，所有修改操作都会返回一个新的 `Node` 或 `Fragment` 实例。

```typescript
// 假设 doc 是一个 Node
// 在位置 10 处插入一个新段落
const newParagraph = mySchema.nodes.paragraph.create(null, mySchema.text('new'))
const newContent = doc.content.replaceChild(10, newParagraph) // 这是一个假设的API，实际使用 transform
// 实际操作是通过 prosemirror-transform 完成的，但原理相同：
// const tr = state.tr.insert(10, newParagraph);
// const newDoc = tr.doc; // newDoc 是一个全新的 Node 实例
```

---

### 3. Slice：智能剪贴板的秘密

当用户复制内容时，我们得到的不仅仅是一段文档片段 (`Fragment`)，而是一个 `Slice` 对象。这是实现“智能粘贴”的关键。

一个 `Slice` 包含：

- `content: Fragment`: 复制的实际内容。
- `openStart: number`: 切片**开始端**的“开放深度”。
- `openEnd: number`: 切片**结束端**的“开放深度”。

**用一个例子来解释 `openStart`/`openEnd` 的魔力：**

假设文档结构是：`<ul><li><p>Hello World</p></li></ul>`

用户选中了从 "llo" 到 "Wo" 的部分文本。

1.  **复制的内容 (`content`)**: 是一个 `Fragment`，看起来像 `<p>llo Wo</p>`。
2.  **分析开放深度**:
    - 选区的**开始**位置，其祖先结构是 `doc -> ul -> li -> p`。它在 `p` 节点内部，所以有 3 层开放的父节点。但由于 `ul` 和 `li` 是在选区之外闭合的，所以实际的开放深度是 1（只开放了 `p`）。
    - 选区的**结束**位置同理，也是在 `p` 节点内部，开放深度为 1。
    - 所以，这个 `Slice` 的 `openStart` 为 1，`openEnd` 为 1。

**粘贴时会发生什么？**

当用户试图将这个 `Slice` 粘贴到一个 `<h1>` 标题中时：

- ProseMirror 查看 `Slice`，发现它的内容 (`<p>llo Wo</p>`) 不能直接放入 `<h1>`（根据 Schema）。
- 但它也看到了 `openStart=1` 和 `openEnd=1`，这告诉它：“这个内容的本来面目是段落的一部分”。
- 于是，ProseMirror 会智能地“剥离”掉外层的 `<p>` 包装，只提取出其中的内联内容 (`"llo Wo"`)，然后将这些文本粘贴到 `<h1>` 中。

**如果复制的是整个列表项 `<li><p>Hello World</p></li>` 呢？**

- `openStart` 会是 0（因为 `<li>` 本身是完整的），`openEnd` 也是 0。
- 当粘贴到一个段落中时，ProseMirror 会尝试将整个 `<li>` 插入。如果 Schema 不允许，它可能会尝试拆分当前段落，并将这个列表项插入到它们之间。

`Slice` 的开放深度信息，使得 ProseMirror 能够根据上下文，以最符合用户直觉的方式合并文档结构。

---

### 4. JSON 序列化：持久化与数据交换

`prosemirror-model` 提供了一套标准的、可预测的 JSON 格式来表示文档，这对于**将文档保存到数据库或通过网络发送至关重要。**

一个 `Node` 可以通过 `.toJSON()` 方法转换成一个 JSON 对象：

```json
{
  "type": "doc",
  "content": [
    {
      "type": "paragraph",
      "content": [
        {
          "type": "text",
          "text": "Hello "
        },
        {
          "type": "text",
          "marks": [{ "type": "strong" }],
          "text": "World"
        }
      ]
    }
  ]
}
```

你可以使用 `Node.fromJSON(schema, json)` 来从这个 JSON 对象恢复成一个完整的 `Node` 实例。

### 总结

`prosemirror-model` 是一个设计精巧、功能强大的数据建模库。它的核心特性——**严格的 Schema、不可变的数据结构、以及带有上下文信息的 Slice**——共同构成了 ProseMirror 可靠、可扩展和高性能的基础。它将富文本文档从一堆混乱的 DOM 标签，提升为了一个类型安全、结构严谨、行为可预测的程序化数据结构。
为 ProseMirror 编辑器提供了一个健壮、结构化的文档模型。它通过 Schema 保证了文档的合法性，并通过 Parser 和 Serializer 实现了模型与浏览器 DOM 之间的解耦和双向转换，这是 ProseMirror 强大功能和灵活性的基础。
