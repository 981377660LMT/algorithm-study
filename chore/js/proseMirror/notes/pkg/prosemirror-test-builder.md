好的，我们来对 `prosemirror-test-builder` 进行一次彻底的、深入的剖析。这个包在 ProseMirror 生态中不直接用于生产环境，但它对于任何严肃的 ProseMirror 开发来说，都是一个**不可或缺的开发工具**。它极大地简化了编写测试用例的复杂度，让测试变得直观、可读且易于维护。

我们将从以下四个关键角度来解构它：

1.  **核心问题：测试 ProseMirror 的痛苦之处**
2.  **解决方案：使用“带标签的模板字符串”来构建文档**
3.  **标签的魔力：`<cursor>`、`<anchor>` 和其他**
4.  **实践应用：编写一个真实的测试用例**

---

### 1. 核心问题：测试 ProseMirror 的痛苦之处

为 ProseMirror 编写测试的主要难点在于**构建测试所需的 `EditorState`**。想象一下，你想测试一个命令，这个命令只在一个包含粗体字的段落、且光标位于粗体字中间时才生效。

如果**不使用** `prosemirror-test-builder`，你需要手动构建这个状态：

```typescript
// 手动构建，非常繁琐和不直观
import { schema } from './schema'
import { EditorState, TextSelection } from 'prosemirror-state'

const p = schema.nodes.paragraph.create(null, [
  schema.text('some '),
  schema.text('bold', [schema.marks.strong.create()]),
  schema.text(' text')
])
const doc = schema.nodes.doc.create(null, [p])

// 计算光标位置： "some " (5) + "b" (1) = 6
const selection = TextSelection.create(doc, 6)

const state = EditorState.create({
  doc,
  selection
})
```

这种方式有几个巨大的缺点：

- **极其繁琐**: 创建每一个节点和标记都需要调用相应的 `create` 方法。
- **可读性差**: 很难从代码中一眼看出最终的文档结构和光标位置。
- **位置计算易错**: 手动计算光标或选区的 `pos` 是一场噩梦，尤其是在文档结构复杂时。任何微小的文本改动都可能需要重新计算所有位置。

`prosemirror-test-builder` 的诞生就是为了彻底解决这个问题。

---

### 2. 解决方案：使用“带标签的模板字符串”来构建文档

`prosemirror-test-builder` 的核心是一个名为 `builders` 的工厂函数。它利用了 JavaScript ES6 的一个强大特性：**带标签的模板字符串 (Tagged Template Literals)**。

#### a. `builders` 函数

首先，你需要用你的 `Schema` 来创建一个“构建器”集合：

```typescript
import { builders } from 'prosemirror-test-builder'
import { mySchema } from './schema'

// builders 函数返回一个对象，其键是 Schema 中的节点名
// 例如：{ doc: Function, p: Function, h1: Function, ... }
const { doc, p, strong } = builders(mySchema)
```

#### b. 带标签的模板字符串

现在，你可以像下面这样使用返回的 `doc` 和 `p` 函数：

```typescript
// 使用构建器，非常直观
const myDoc = doc(p('some ', strong('bold'), ' text'))
```

这已经比手动构建好很多了。但 `prosemirror-test-builder` 的真正魔力在于，你可以直接使用模板字符串，并通过特殊的**“标签” (`<tag>`)** 来定义光标和选区。

```typescript
// 最终的、最直观的形式
const myDocWithSelection = doc(p('some <b><cursor></b>old text'))
```

这行代码做了什么？

- 它解析了这个类似 HTML 的字符串。
- 它创建了一个包含 `<p>some bold text</p>` 的 ProseMirror `Node`。
- 它识别了 `<b>` 标签，并将其转换成 `strong` 标记。
- 最重要的是，它识别了 `<cursor>` 标签，并自动计算出该位置的 `pos`，然后创建了一个 `TextSelection` 在那里。
- 最终，`myDocWithSelection` 不仅仅是一个 `Node`，它是一个包含了 `doc` 和 `selection` 的完整测试状态对象！

---

### 3. 标签的魔力：`<cursor>`、`<anchor>` 和其他

`prosemirror-test-builder` 的核心就是这些在模板字符串中使用的“标签”。它们让你能够以一种所见即所得的方式，精确地定义测试场景中的选区。

#### 最常用的标签：

- **`<cursor>`**:

  - **作用**: 定义一个光标（一个零宽度的选区）。
  - **等价于**: `TextSelection.create(doc, pos)`。
  - **示例**: `p("Hello <cursor>World")`

- **`<anchor>`** 和 **`<head>`**:

  - **作用**: 定义一个有范围的 `TextSelection`。`anchor` 是选区的固定点，`head` 是活动点。
  - **示例**: `p("Select <anchor>this<head> text")`

- **`<node>`**:

  - **作用**: 定义一个 `NodeSelection`，选中紧随其后的那个节点。
  - **示例**: `p("Select the image: <node><img>")`

- **`<all>`**:

  - **作用**: 定义一个 `AllSelection`，选中整个文档。
  - **示例**: `doc(<all>p("Select all"))`

- **自定义标签**:
  - 你还可以使用任意名称的标签，如 `<mytag>`。这些标签的位置会被记录下来，你可以在测试中断言它们的位置。
  - **示例**: `p("A<start>B<end>C")`
  - 在测试中，你可以通过 `doc.tag.start` 和 `doc.tag.end` 来获取 `A` 和 `B` 之间的位置，非常适合用来验证某个命令是否正确地修改了特定范围。

---

### 4. 实践应用：编写一个真实的测试用例

让我们用 `prosemirror-test-builder` 来测试 `prosemirror-commands` 中的 `toggleMark` 命令。

**测试场景**:

1.  给定一个段落，其中部分文本是粗体。
2.  将光标置于粗体文本的**外部**。
3.  执行 `toggleMark(strong)` 命令。
4.  断言：光标所在单词应该被加上粗体，而原有的粗体文本不受影响。

```typescript
import { builders } from 'prosemirror-test-builder'
import { toggleMark } from 'prosemirror-commands'
import { mySchema } from '../src/schema' // 你的编辑器 Schema

// 1. 创建构建器
const { doc, p, strong } = builders(mySchema, {
  // 你可以为标记提供简写，比如用 <b> 代替 <strong '...'>
  b: schema.marks.strong
})

describe('toggleMark command', () => {
  it('should make a word bold', () => {
    // 2. 使用模板字符串定义初始状态
    // 'state' 是一个包含 doc 和 selection 的对象
    const state = doc(p("Don't touch <b>this</b>, but make <cursor>that bold"))

    // 3. 执行命令
    let resultState = state
    const command = toggleMark(mySchema.marks.strong)
    // 在测试中，我们手动 dispatch
    command(state, tr => {
      resultState = state.apply(tr)
    })

    // 4. 定义期望的最终状态
    const expected = doc(p("Don't touch <b>this</b>, but make <b>that</b> bold"))

    // 5. 断言
    // prosemirror-model 的 toJSON 方法可以用来比较文档结构
    expect(resultState.doc.toJSON()).toEqual(expected.doc.toJSON())
  })
})
```

这个测试用例清晰、简洁、可读性极高。任何人都能一眼看出测试的初始状态和期望结果，这在没有 `prosemirror-test-builder` 的情况下是几乎不可能做到的。

### 总结

`prosemirror-test-builder` 是一个通过牺牲极小的运行时性能（仅在测试环境中）来换取巨大开发效率和代码可维护性的典范。它通过**带标签的模板字符串**这一巧妙的抽象，将繁琐、易错的测试状态构建过程，变成了一种直观、声明式的“绘图”过程。对于任何需要编写单元测试或集成测试的 ProseMirror 项目，它都应该是你的首选工具。
