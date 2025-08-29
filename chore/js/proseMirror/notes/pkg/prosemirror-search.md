好的，我们来对 `prosemirror-search` 进行一次彻底的、深入的剖析。这个包提供了一个健壮的、结构化的查找与替换功能，它远比在纯文本上执行 `String.prototype.indexOf` 要强大得多。它被设计用来在复杂的、结构化的 ProseMirror 文档中进行精确操作。

我们将从以下四个关键角度来解构它：

1.  **核心问题：超越简单的文本搜索**
2.  **`SearchQuery`：搜索逻辑的封装**
3.  **从查询到高亮：`Decoration` 的应用**
4.  **API 与实践：如何驱动查找与替换**

---

### 1. 核心问题：超越简单的文本搜索

在一个简单的 `<textarea>` 中，查找和替换相对容易。但在 ProseMirror 中，文档是一棵结构化的节点树，这带来了新的挑战：

- **结构感知**: 一个简单的文本搜索会忽略文档结构。例如，搜索 `doc.textContent` 无法区分文本是在段落中、标题中，还是在代码块中。更重要的是，它返回的索引位置无法直接映射回 ProseMirror 文档中有效的 `pos`。
- **原子节点**: 文档中可能包含图片、视频等“原子”节点，这些节点会中断文本流，必须被正确地跳过。
- **性能**: 对大型文档进行暴力字符串匹配可能会很慢。需要一种高效的方式来迭代文本内容。

`prosemirror-search` 通过直接在 ProseMirror 的数据模型上操作，解决了所有这些问题。它逐个遍历文本节点，在节点内部进行匹配，并正确计算出每个匹配项在整个文档中的绝对 `pos`。

---

### 2. `SearchQuery`：搜索逻辑的封装

`SearchQuery` 类是 `prosemirror-search` 的大脑。它将一个搜索请求（包括要查找的内容、大小写敏感性等）封装成一个可执行的对象。

#### `SearchQuery` 的构造函数

```typescript
new SearchQuery(
  regexp: RegExp,
  caseSensitive?: boolean,
  wholeWord?: boolean // 这个参数在较新版本中可能被废弃，逻辑通常由 regexp 本身处理
)
```

- **`regexp: RegExp`**: 这是最核心的参数。`prosemirror-search` 的强大之处在于它基于**正则表达式**，而不仅仅是纯字符串。
  - **重要约定**: 你传递的正则表达式**必须**包含 `g` (global) 标志，这样才能找到所有匹配项，而不仅仅是第一个。
  - **大小写**: 如果你想进行不区分大小写的搜索，你需要给正则表达式添加 `i` 标志。`caseSensitive` 参数只是一个辅助，最终它也是通过修改正则表达式来实现的。

**示例**:

```typescript
// 查找所有 "ProseMirror"，区分大小写
const query1 = new SearchQuery(/ProseMirror/g)

// 查找所有 "prosemirror"，不区分大小写
const query2 = new SearchQuery(/prosemirror/gi)

// 查找所有看起来像十六进制颜色的代码
const query3 = new SearchQuery(/#[0-9a-fA-F]{6}/g)
```

这个 `SearchQuery` 对象一旦创建，就可以被传递给插件或命令，用于执行实际的搜索操作。

---

### 3. 从查询到高亮：`Decoration` 的应用

`prosemirror-search` 的查找功能在视觉上是通过 ProseMirror 的 `Decoration` 系统实现的。当执行一个查找时，它并**不改变**文档内容或选区，而是为所有匹配项“贴上”一个装饰。

**工作流程**:

1.  **触发搜索**: 用户在搜索框中输入内容，你的应用代码创建一个 `SearchQuery` 对象。
2.  **通知插件**: 你通过 `setMeta` 方法将这个新的 `SearchQuery` 对象发送给 `search` 插件。
    ```typescript
    view.dispatch(view.state.tr.setMeta(searchPluginKey, { query: newQuery }))
    ```
3.  **插件执行搜索**:
    - `search` 插件的 `state.apply` 方法接收到这个元数据。
    - 它会启动一个“游标”（cursor），这个游标会智能地遍历文档中的所有**文本节点**。
    - 在每个文本节点内部，它会使用 `SearchQuery` 中的正则表达式来查找匹配项。
    - 对于每一个找到的匹配项，它会记录下其在整个文档中的绝对 `from` 和 `to` 位置。
4.  **创建 `Decoration`**:
    - 插件收集所有匹配项的 `{ from, to }` 范围。
    - 对于每一个范围，它会创建一个 `Decoration.inline`。
    - `Decoration.inline(from, to, { class: 'search-match' })`
5.  **更新状态与视图**:
    - 插件将所有这些 `Decoration` 组合成一个 `DecorationSet`，并存入自己的状态中。
    - ProseMirror 视图检测到 `Decoration` 发生了变化，于是重新渲染视图，将 `search-match` 这个 CSS 类应用到所有匹配的文本片段上。
6.  **CSS 生效**: 你需要在你的 CSS 文件中定义 `.search-match` 的样式，例如给它一个黄色的背景，从而让用户看到高亮效果。

```css
.search-match {
  background-color: #fffa80;
}
```

---

### 4. API 与实践：如何驱动查找与替换

`prosemirror-search` 提供了一套命令式的 API 来驱动整个流程。

#### a. `search(state, dispatch?, query?)`

这是 `search` 插件的工厂函数。你通常在初始化编辑器时只调用一次，不带任何参数，来获取插件实例和它的 `key`。

```typescript
import { search, searchKey as searchPluginKey } from 'prosemirror-search'

const mySearchPlugin = search()

const plugins = [
  // ...
  mySearchPlugin
]
```

#### b. `findNext(state, dispatch?)` 和 `findPrev(state, dispatch?)`

- 这两个是 `Command` 函数，用于在已高亮的搜索结果之间跳转。
- 当被调用时，它们会从插件状态中读取当前的搜索结果，找到当前光标之后（或之前）的第一个匹配项，然后创建一个事务将编辑器的选区移动到那里。
- 你可以将它们绑定到 "Next" 和 "Previous" 按钮上。

#### c. `replace(state, dispatch?, content?)`

- 这是执行“替换当前项”的命令。
- 它会检查当前选区是否恰好覆盖了一个搜索结果（插件内部会做这个判断）。
- 如果是，它会用你提供的 `content` 字符串替换掉这个匹配项。

#### d. `replaceAll(state, dispatch?, content?)`

- 执行“全部替换”。
- 它会遍历插件状态中所有的搜索结果，并创建一个包含多个 `ReplaceStep` 的单一 `Transaction`，一次性地替换掉所有匹配项。

**一个简单的实现流程**:

1.  **UI**: 创建一个搜索框、一个替换框、"Find", "Next", "Replace", "Replace All" 按钮。
2.  **Find**:
    - 当用户点击 "Find" 或在搜索框输入时，根据输入内容创建一个 `SearchQuery`。
    - 通过 `tr.setMeta(searchPluginKey, { query })` 将查询发送给插件，触发高亮。
3.  **Next**:
    - 按钮的 `onClick` 调用 `findNext(view.state, view.dispatch)`。
4.  **Replace**:
    - 按钮的 `onClick` 调用 `replace(view.state, view.dispatch, replacementText)`。
5.  **Replace All**:
    - 按钮的 `onClick` 调用 `replaceAll(view.state, view.dispatch, replacementText)`。
6.  **关闭搜索**:
    - 发送一个空的 `SearchQuery` (`new SearchQuery(//g)`) 或一个特定的关闭信号给插件，插件会清空 `DecorationSet`，高亮随之消失。

### 总结

`prosemirror-search` 是一个功能完备的、基于 ProseMirror 核心机制构建的搜索与替换解决方案。它通过：

- **一个强大的 `SearchQuery` 对象**，利用正则表达式提供灵活的匹配能力。
- **与 `Decoration` 系统的深度集成**，实现了非侵入式、高性能的结果高亮。
- **一套清晰的命令式 API**，让开发者可以轻松地构建出功能齐全的查找/替换用户界面。

它完美地展示了如何利用 ProseMirror 的插件、元数据事务和装饰器系统来构建复杂而强大的编辑器功能。
