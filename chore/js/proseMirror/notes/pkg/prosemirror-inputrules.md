好的，我们来对 `prosemirror-inputrules` 进行一次彻底的、深入的剖析。这个包是提升 ProseMirror 编辑器“智能化”和“效率”的关键模块之一。它实现了在用户输入特定文本模式时，自动触发相应格式转换的功能，也就是我们常说的“输入时自动格式化”（As-You-Type Formatting）。

我们将从以下四个关键角度来解构它：

1.  **核心概念：从“文本模式”到“状态变更”的魔法**
2.  **`InputRule` 类：规则的原子定义**
3.  **`inputRules` 插件：驱动规则的引擎**
4.  **实践与定制：创建你自己的输入规则**

---

### 1. 核心概念：从“文本模式”到“状态变更”的魔法

`prosemirror-inputrules` 的核心思想非常直观：**监听用户的文本输入，如果输入的内容匹配了预设的正则表达式，就自动执行一个相应的命令（`Transaction`）**。

这与 `prosemirror-keymap` 有着本质的区别：

- `prosemirror-keymap` 响应的是**特定的按键组合**（如 `Ctrl+B`）。
- `prosemirror-inputrules` 响应的是**特定的文本内容模式**（如输入了 `## `）。

常见的例子包括：

- 输入 `## ` 自动转换成一个二级标题。
- 输入 `1. ` 自动开启一个有序列表。
- 输入 `---` 自动转换成一条水平分割线。
- 输入 `(c)` 自动转换成版权符号 `©`。

它通过在用户输入后立即提供格式反馈，极大地提升了写作效率，是现代富文本编辑器（如 Notion, Slack）的标配功能。

---

### 2. `InputRule` 类：规则的原子定义

所有输入规则的基础是 `InputRule` 类。一个 `InputRule` 实例定义了一个完整的“匹配-执行”逻辑。

#### `InputRule` 的构造函数

```typescript
new InputRule(
  match: RegExp,
  handler: (
    state: EditorState,
    match: RegExpExecArray,
    start: number,
    end: number
  ) => Transaction | null
)
```

让我们来剖析这两个至关重要的参数：

#### a. `match: RegExp` (正则表达式)

这是规则的“触发器”。这个正则表达式有一个**非常重要的约定**：它**必须以 `$` 结尾**。

- **为什么必须以 `$` 结尾？** 因为 `prosemirror-inputrules` 的工作方式是检查**从当前文本块的开头到当前光标位置**的这段文本。`$` 确保了你的正则表达式只在模式恰好出现在光标之前时才匹配成功。
- **示例**:
  - 匹配 `## ` 的规则：`/^(##\s)$/`
  - 匹配三个或更多破折号的规则：`/^(\-\-\-)$/`
  - 匹配 `(c)` 的规则：`/\(c\)$/`

#### b. `handler` (处理器函数)

这是规则的“执行器”。当 `match` 成功匹配时，这个函数会被调用，并由它来生成一个描述如何变更文档的 `Transaction`。

- **参数**:
  - `state: EditorState`: 当前的编辑器状态。
  - `match: RegExpExecArray`: 正则表达式 `exec()` 方法的匹配结果数组。`match[0]` 是整个匹配的字符串，`match[1]`, `match[2]` 等是捕获组。
  - `start: number`: 匹配文本在文档中的起始位置。
  - `end: number`: 匹配文本在文档中的结束位置（也就是当前光标位置）。
- **返回值**:
  - 如果规则应该被执行，函数必须返回一个**`Transaction` 对象**。这个事务通常会删除匹配到的文本（如 `## `），并应用新的节点或标记。
  - 如果出于某种原因，即使匹配成功也不应执行（例如，在代码块中不应转换标题），函数可以返回 `null` 来中止操作。

---

### 3. `inputRules` 插件：驱动规则的引擎

`InputRule` 实例本身只是一个数据结构。要让它们真正工作，你需要使用 `inputRules` 工厂函数来创建一个插件。

```typescript
import { inputRules, textblockTypeInputRule } from 'prosemirror-inputrules'
import { schema } from './schema'

const myInputRulesPlugin = inputRules({
  rules: [
    // 将 ## 转换为 h2
    textblockTypeInputRule(/^(##\s)$/, schema.nodes.heading, { level: 2 })
    // ... 其他规则
  ]
})
```

这个插件的内部工作流程如下：

1.  **监听事务**: 插件通过 `state.apply` 方法监听每一个进入编辑器的 `Transaction`。
2.  **识别文本输入**: 它专门寻找那些由用户输入文本而产生的事务（通常是包含一个 `ReplaceStep` 并且 `slice` 中有文本的事务）。
3.  **触发检查**: 当检测到文本输入时，它会触发检查逻辑。
4.  **执行匹配**: 它会获取当前光标所在的文本块，并从块的开头到光标位置截取字符串。然后，它会遍历你提供的 `rules` 数组，用每个规则的 `RegExp` 去匹配这段字符串。
5.  **调用 `handler`**: 一旦找到第一个匹配的规则，它就会调用该规则的 `handler` 函数，并传入所有必要的参数。
6.  **分发新事务**: 如果 `handler` 返回了一个 `Transaction` (`tr`)，插件会获取这个 `tr` 中的所有 `Step`，并将它们附加到当前正在处理的主事务上。
7.  **完成转换**: 当主事务最终被应用到视图时，它已经包含了输入规则产生的变更，用户看到的就是一个无缝的格式转换。

---

### 4. 实践与定制：创建你自己的输入规则

`prosemirror-inputrules` 自带了几个非常有用的规则构造器：

- **`textblockTypeInputRule(regexp, nodeType, attrs?)`**: 用于将匹配的文本转换成一个指定类型的**块级节点**。这是创建标题、代码块规则最简单的方式。
- **`wrappingInputRule(regexp, nodeType, attrs?, join?)`**: 用于将当前文本块**包裹**在另一个节点中。这是创建列表、引用块规则最简单的方式。
- **`smartQuotes`**: 一个预设的规则数组，用于将直引号 (`'`, `"`) 转换成弯引号（`‘ ’`, `“ ”`）。

#### 示例：创建一个将 `(c)` 转换为 `©` 的自定义规则

这是一个典型的“文本替换”规则，我们需要手动创建。

```typescript
import { InputRule } from 'prosemirror-inputrules'
import { EditorState, Transaction } from 'prosemirror-state'

// 1. 定义正则表达式，必须以 $ 结尾
const copyrightRule = new InputRule(
  /\(c\)$/,
  (state: EditorState, match: RegExpExecArray, start: number, end: number): Transaction | null => {
    // 2. 检查是否在不应转换的地方（可选，但推荐）
    // 例如，我们不希望在代码块中进行转换
    const $start = state.doc.resolve(start)
    if (!$start.parent.type.spec.marks) {
      // 检查父节点是否允许标记（一个简单的判断是否为代码块的方法）
      return null
    }

    // 3. 创建一个事务
    const tr = state.tr

    // 4. 在事务中，用 © 符号替换掉匹配到的 (c) 文本
    tr.replaceWith(start, end, state.schema.text('©'))

    // 5. 返回这个事务
    return tr
  }
)

// 然后将它加入到你的插件配置中
const myPlugin = inputRules({
  rules: [
    ...smartQuotes,
    copyrightRule
    // ...其他规则
  ]
})
```

### 总结

`prosemirror-inputrules` 是一个强大而灵活的自动化工具。它通过一个清晰的**“正则匹配 + 事务处理器”**模型，让开发者可以轻松地定义各种“输入时自动格式化”的行为。理解 `InputRule` 的构造、特别是 `$` 在正则表达式中的关键作用，以及如何编写 `handler` 函数来生成 `Transaction`，是掌握这个模块并为你的编辑器创造高效、流畅输入体验的核心。
