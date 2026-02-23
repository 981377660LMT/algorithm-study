# ProseMirror 进阶 API 与细节补充

本文档是对基础 API 的补充，涵盖了 Schema 高级属性、Plugin 系统深层机制、NodeView 完整生命周期以及 Transform 核心原理。

## 1. Schema 高级定义 (Advanced Schema)

在 `nodes` 和 `marks` 的定义中，除了 `content` 和 `toDOM`，还有一些控制编辑器行为的关键属性。

### Node Spec 属性

- **`defining: boolean`**
  - _含义_: 如果为 `true`，当用户全选该节点内容并粘贴新内容时，会保留该节点的结构（而不是替换整个节点）。
  - _场景_: 标题 (`heading`)、引用 (`blockquote`)。你希望粘贴文字到标题里时，它还是标题。
- **`isolating: boolean`**
  - _含义_: 如果为 `true`，该节点被视为“孤岛”，普通的编辑操作（如退格键删除）不会跨越该节点边界。
  - _场景_: 表格单元格 (`table_cell`)。你在单元格里按退格键，删光文字后光标应该停在格子里，而不是把前一个单元格的内容吸进来。
- **`draggable: boolean`**
  - _含义_: 该节点是否可拖拽。
- **`selectable: boolean`**
  - _含义_: 该节点是否可被选中（作为 `NodeSelection`）。默认为 `true`。
- **`code: boolean`**
  - _含义_: 标识该节点包含代码。这会影响某些输入规则和拼写检查行为。

### Mark Spec 属性

- **`excludes: string`**
  - _含义_: 定义该 Mark 与哪些 Mark 互斥。默认为空（允许共存）。
  - _示例_: `excludes: "_"` 表示该 Mark 排斥所有其他 Mark。`excludes: "strong"` 表示不能同时加粗。
- **`spanning: boolean`**
  - _含义_: 是否允许该 Mark 跨越多个节点。默认为 `true`。

## 2. 插件系统全解 (Plugin System)

Plugin 不仅仅是响应事件，它还能拥有自己的 State，并拦截/修改 Transaction。

### Plugin Spec 核心字段

- **`state`**: 定义插件的独立状态。
  - `init(config, instance)`: 初始化状态。
  - `apply(tr, value, oldState, newState)`: **核心**。每次文档变化都会触发。必须返回新的状态值。
    - _注意_: 这里是计算密集型操作的热点，务必优化性能。
- **`key`**: `PluginKey`。用于在其他地方通过 `key.getState(state)` 获取该插件的状态。
- **`view(editorView)`**: 当插件被注册到 View 时调用。
  - 返回一个对象 `{ update(view, prevState), destroy() }`。
  - _场景_: 用于管理与插件相关的外部 UI（如气泡菜单的显隐）。
- **`appendTransaction(transactions, oldState, newState)`**:
  - _含义_: 允许插件在每次 Transaction 应用后，追加一个新的 Transaction。
  - _场景_: 强制文档规范化（如：确保文档末尾总有一个空段落）。
  - _注意_: 避免死循环（追加的 Transaction 又触发 appendTransaction）。
- **`filterTransaction(tr, state)`**:
  - _含义_: 返回 `false` 可阻止 Transaction 应用。
  - _场景_: 冻结文档（只读模式的一种实现）、权限控制。

## 3. NodeView 生命周期详解

`NodeView` 是连接 Model 和 DOM 的桥梁，其生命周期方法决定了交互体验。

```javascript
{
  dom: HTMLElement, // 必选。节点对应的根 DOM。
  contentDOM: HTMLElement, // 可选。如果节点有子内容，指定子内容渲染的容器。

  // 1. 更新
  update(node, decorations, innerDecorations) {
    // 返回 true: 表示我成功更新了 DOM，ProseMirror 不需要重新渲染。
    // 返回 false: 表示我处理不了（比如节点类型变了），ProseMirror 会销毁当前 NodeView 并重建。
    if (node.type !== myType) return false;
    // 手动更新 DOM 属性...
    return true;
  },

  // 2. 选中状态变化
  selectNode() {
    this.dom.classList.add("ProseMirror-selectednode");
  },
  deselectNode() {
    this.dom.classList.remove("ProseMirror-selectednode");
  },

  // 3. 事件控制
  stopEvent(event) {
    // 返回 true 阻止事件冒泡给编辑器。
    // 场景: 节点内部有 input 输入框，或者点击按钮。
    return /mousedown|keypress/.test(event.type);
  },

  ignoreMutation(mutation) {
    // 返回 true 告诉 ProseMirror: "这个 DOM 变化是我自己搞的（比如 React 渲染），你别管"。
    // 如果返回 false，ProseMirror 会认为 DOM 被外部破坏了，会尝试重绘修复。
    // 场景: 只要是 contentDOM 以外的变化，通常都应该返回 true。
    return true;
  },

  // 4. 销毁
  destroy() {
    // 清理工作，如卸载 React 组件、移除事件监听。
  }
}
```

## 4. Transform 与 Mapping (变换与映射)

这是 ProseMirror 处理协同编辑和复杂编辑逻辑的核心。

### Step (步骤)

Transaction 本质上是一系列 Step 的集合。

- **`ReplaceStep`**: 最基础的步骤，替换一段范围的内容。
- **`AddMarkStep` / `RemoveMarkStep`**: 添加/移除标记。

### Mapping (映射)

当文档发生变化后，旧的位置（Pos）在旧文档中有效，但在新文档中可能已经移动或消失。`Mapping` 用于计算这种位置变换。

- **`map(pos, bias)`**: 将旧位置映射到新位置。
  - `bias`: 当位置正好在插入内容的边界时，倾向于向左(-1)还是向右(1)移动。
- **`Mappable` 接口**: `Transaction` 实现了 `Mappable`。
  - `tr.mapping.map(oldPos)`: 获取经过该 Transaction 所有步骤后的新位置。

## 5. Decorations (装饰器) 三种类型

Decorations 属于 View 层，不改变文档数据。

1.  **Widget Decoration**: 在文档中插入一个 HTML 元素（不属于文档内容）。
    - _场景_: 占位符、光标位置的协作头像、行号。
    - _API_: `Decoration.widget(pos, dom, spec)`
2.  **Inline Decoration**: 给现有文本添加样式或属性（类似 Mark，但临时的）。
    - _场景_: 搜索高亮、拼写错误波浪线、评论高亮。
    - _API_: `Decoration.inline(from, to, attrs, spec)`
3.  **Node Decoration**: 给现有节点（Node）的 DOM 添加属性或类名。
    - _场景_: 给当前选中的段落加背景色、给表格行加样式。
    - _API_: `Decoration.node(from, to, attrs, spec)`

## 6. ParseRule (解析规则) 细节

在 `Schema.parseDOM` 中定义。

- **`priority: number`**: 匹配优先级。默认为 50。
  - _场景_: 如果你有两个规则都能匹配 `<div>`，优先级高的生效。
- **`context: string`**: 上下文限制。
  - _示例_: `context: "list/"` 表示只有在列表里的 `p` 标签才匹配这个规则。
- **`skip: boolean`**: 如果为 `true`，不创建节点，但继续解析其子内容。
  - _场景_: 解析 `<div><p>Text</p></div>`，如果 `div` 设为 `skip`，则直接把 `p` 提升上来。
- **`getContent(dom, schema)`**: 自定义如何从 DOM 提取内容。

## 7. 常用辅助方法补充

- **`Node.rangeHasMark(from, to, type)`**: 检查范围内是否包含某 Mark。
- **`Node.nodesBetween(from, to, callback)`**: 遍历范围内的所有节点。
  - _技巧_: `callback` 返回 `false` 可以停止遍历该节点的子节点。
- **`EditorView.someProp(propName, f)`**: 查找并执行某个 Prop（从插件列表里找）。

## 8. InputRules (输入规则) 深度定制

`InputRules` 用于实现类似 Markdown 的自动格式化（如输入 `* ` 变列表）。

- **原理**: 监听 `textInput` 事件，当输入的文本匹配正则表达式时，触发 Transaction。
- **`new InputRule(match, handler)`**:
  - `match`: 正则表达式。**必须以 `$` 结尾**，确保匹配的是刚输入的内容。
  - `handler`: `(state, match, start, end) => Transaction | null`。
    - `start`, `end`: 匹配到的文本在文档中的范围。
    - 返回 `null` 表示不处理。
- **常用技巧**:
  - **捕获组**: 正则中的捕获组可以在 `match` 数组中获取，用于提取参数（如 `#{1,6}` 提取标题级别）。
  - **`wrappingInputRule`**: 官方提供的辅助函数，用于将当前块包裹在某节点中（如引用）。
  - **`textblockTypeInputRule`**: 用于改变当前块的类型（如标题）。

## 9. Keymap (按键映射) 与命令链

- **优先级**: 插件列表 (`state.plugins`) 中，**排在前面**的插件定义的 Keymap 优先级更高。
- **`chainCommands(cmd1, cmd2, ...)`**:
  - 组合多个命令。依次执行 `cmd1`, `cmd2`...
  - 如果某个命令返回 `true`，则停止后续执行。
  - _场景_: `Enter` 键。先尝试“跳出列表”，如果不行则“插入新段落”。
- **修饰符**:
  - `Mod`: Mac 上是 Cmd，Windows/Linux 上是 Ctrl。
  - `Shift-Mod-z`: 常见撤销重做快捷键。

## 10. 剪贴板与 Slice 处理 (Clipboard)

ProseMirror 的复制粘贴不仅仅是 HTML，还涉及 `Slice` 的处理。

- **`transformPasted(slice)`**:
  - _场景_: 粘贴时清洗数据。例如，从 Word 粘贴过来包含大量无用样式，或者将纯文本 URL 转换为卡片节点。
  - _操作_: 遍历 `slice.content`，返回一个新的 `Slice`。
- **`transformPastedHTML(html)`**:
  - _场景_: 在解析 HTML 之前对其进行字符串级别的清洗。
- **`clipboardTextSerializer(slice)`**:
  - _场景_: 当用户复制内容到纯文本编辑器（如记事本）时，控制生成的文本格式。

## 11. 性能优化核心点 (Performance)

- **`view.updateState` 节流**: 虽然 PM 内部有 Diff，但如果外部频繁触发（如 React 的 `useEffect` 依赖项没写好），会导致性能问题。
- **`DecorationSet` 优化**:
  - 如果 Decoration 很多（如代码高亮），不要每次都重新创建整个 Set。
  - 使用 `DecorationSet.map` 将旧的 Set 映射到新位置，然后只计算变动部分的 Decoration，最后用 `add` / `remove` 更新。
- **大文档的 `Plugin.apply`**:
  - 避免在 `apply` 中遍历整个 `doc`。
  - 利用 `tr.steps` 判断变更范围，只更新受影响的部分状态。

## 12. DOMSerializer 与高级序列化

除了 `toDOM`，有时我们需要更精细地控制 HTML 输出（例如导出时）。

- **`DOMSerializer.fromSchema(schema)`**: 创建序列化器。
- **`serializeFragment(fragment, options)`**:
  - `document`: 指定创建节点的 document 上下文。
  - **`node`**: 自定义特定节点的序列化逻辑（覆盖 `toDOM`）。
    ```javascript
    serializer.serializeFragment(doc.content, {
      node(node) {
        if (node.type.name === 'image') {
          // 导出时给图片加个 wrapper
          return ['div', { class: 'img-wrap' }, ['img', { src: node.attrs.src }]]
        }
        // 返回 null 使用默认 toDOM
        return null
      }
    })
    ```

## 13. History 模块的高级控制

`prosemirror-history` 提供了比简单的 undo/redo 更细粒度的控制。

- **`closeHistory(tr)`**:
  - _含义_: 强制关闭当前的历史记录组。后续的变更将被放入一个新的撤销项中。
  - _场景_: 用户输入了一段话，停顿了很久，你希望把停顿前后的输入分开撤销。
- **`addToHistory: boolean`**:
  - _含义_: Transaction 的元属性。如果设为 `false`，该次变更不会进入撤销栈。
  - _场景_: 协同编辑中接收别人的变更，或者自动格式化（如自动闭合括号）。
- **合并撤销项**:
  - 默认情况下，连续的输入或删除会被合并为一个撤销步。这是通过 `historyEventSwap` 等内部机制控制的，通常不需要手动干预，但了解这一点有助于调试“为什么撤销了一大段”。

## 14. 选区与光标的高级处理

- **GapCursor (间隙光标)**:
  - _问题_: 当两个块级节点（如两个图片）紧挨着时，用户很难把光标放到它们中间。
  - _解决_: 使用 `prosemirror-gapcursor` 插件。它会渲染一个假的水平光标，允许用户在这些“不可能”的位置输入。
- **`TextSelection.between($anchor, $head)`**:
  - _技巧_: 智能创建选区。它会自动处理方向（从左到右或从右到左）。
- **`NodeSelection.create(doc, from)`**:
  - _技巧_: 强制选中某个节点（如点击图片时）。

## 15. Node 与 Fragment 的底层操作

在编写复杂 Command 时，直接操作 Node 对象比操作 Transaction 更高效（用于计算）。

- **`node.cut(from, to)`**:
  - 获取节点的一部分（返回一个新的 Node，如果切到了中间，会自动补全层级结构）。
- **`node.replace(from, to, slice)`**:
  - 在节点内部执行替换（不生成 Transaction，仅返回新 Node）。用于预计算结果。
- **`ContentMatch`**:
  - _含义_: Schema 验证的核心状态机。
  - _API_: `node.contentMatchAt(index)`。
  - _场景_: 检查在当前位置插入某节点是否合法 (`match.matchType(nodeType)`).

## 16. 测试辅助工具

ProseMirror 官方提供了 `prosemirror-test-builder`，极大简化了测试用例的编写。

- **构建器模式**:

  ```javascript
  import { builders } from 'prosemirror-test-builder'
  const { doc, p, strong } = builders(mySchema)

  // 创建一个文档对象，不用手写 JSON
  const myDoc = doc(p('Hello ', strong('World')))
  ```

- **断言**:
  - `expect(newState.doc).toEqual(myDoc)`: 直接比较 Node 结构。

## 17. 协同编辑核心 API (Collab)

`prosemirror-collab` 模块虽然封装了大部分逻辑，但理解其底层 API 对于调试至关重要。

- **`sendableSteps(state)`**:
  - _作用_: 获取当前客户端已产生但尚未发送给服务端的 Steps。
  - _返回_: `{ version, steps, clientID, origin }`。
  - _场景_: 定时轮询或 WebSocket 发送数据时调用。
- **`receiveTransaction(state, steps, clientIDs)`**:
  - _作用_: 将服务端广播过来的 Steps 应用到本地状态。
  - _返回_: 一个新的 Transaction。
  - _注意_: 这个 Transaction 会自动处理 Rebase（变基），即如果本地有未提交的 Steps，它会先撤销本地 Steps，应用远程 Steps，再重新应用本地 Steps。
- **`collab({ version })`**:
  - _配置_: `version` 必须与服务端当前的文档版本一致。如果版本错乱，会导致协同崩溃。

## 18. Schema Content 表达式详解

`content` 属性定义了节点的子节点约束，语法类似正则表达式。

- **基本符号**:
  - `paragraph`: 必须包含一个 paragraph 节点。
  - `paragraph+`: 一个或多个。
  - `paragraph*`: 零个或多个。
  - `paragraph?`: 零个或一个。
- **组合与分组**:
  - `heading paragraph+`: 一个标题后跟至少一个段落。
  - `(paragraph | blockquote)+`: 段落或引用块的混合序列。
- **Group (组)**:
  - 在 Node 定义中设置 `group: "block"`。
  - 表达式中使用 `block+` 即可匹配所有组名为 block 的节点。
  - _特殊组_: `inline` (内联节点), `text` (文本节点)。
- **严格模式**:
  - 默认情况下，ProseMirror 会自动修正不符合 Schema 的内容（如在不允许的地方插入图片，它可能会被丢弃）。

## 19. EditorProps 视图配置细节

`EditorView` 的 `props` 选项中还有一些控制视图行为的细节。

- **`attributes`**:
  - _作用_: 给编辑器根 DOM 元素 (`contenteditable`) 添加属性。
  - _用法_: 可以是对象 `{ class: "my-editor" }` 或函数 `(state) => ({ class: ... })`。
- **`scrollThreshold` / `scrollMargin`**:
  - _作用_: 控制光标移动到边缘时自动滚动的触发阈值和边距。
  - _场景_: 想要类似 VS Code 的滚动体验时调整这些值。
- **`editable`**:
  - _作用_: 控制编辑器是否只读。
  - _用法_: `(state) => boolean`。
  - _注意_: 即使返回 `false`，插件依然可以派发 Transaction，只是用户无法通过 DOM 输入。

## 20. 自定义 Step (Custom Steps)

这是 ProseMirror 最硬核的功能之一，允许你定义全新的编辑原语。

- **场景**: 比如你需要实现一个“折叠/展开”功能，但不想修改文档内容，只想修改一个持久化的元数据，且这个元数据需要随协同编辑同步。
- **实现**:
  1. 继承 `Step` 类。
  2. 实现 `apply(doc)`: 返回 `StepResult`。
  3. 实现 `invert(doc)`: 返回逆向 Step 用于撤销。
  4. 实现 `map(mapping)`: 处理位置映射。
  5. **注册**: 使用 `Step.jsonID` 注册，以便能够被序列化传输。

## 21. 节点与标记的集合操作

在处理复杂的 Schema 转换时，这些底层方法很有用。

- **`MarkType.isInSet(marks)`**:
  - _作用_: 检查某个 MarkType 是否存在于 Mark 数组中。
- **`Mark.addToSet(marks)`**:
  - _作用_: 将当前 Mark 添加到数组中，如果已存在同类型 Mark，则根据配置决定是替换还是共存。
- **`Mark.removeFromSet(marks)`**:
  - _作用_: 从数组中移除当前 Mark。
- **`NodeType.create(attrs, content, marks)`**:
  - _作用_: 手动创建节点实例。比 `schema.node(...)` 更底层。

## 22. 高级变换方法 (Advanced Transforms)

`Transform` 类（以及 `Transaction`）提供了一系列高级方法，用于执行复杂的文档结构调整。

- **`lift(range, targetDepth)`**:
  - _作用_: 将选区内的内容提升到父级。
  - _场景_: 从列表中移出（Shift-Tab），或者取消引用块。
- **`wrap(range, wrappers)`**:
  - _作用_: 将选区内的内容包裹在新的节点中。
  - _场景_: 选中段落变为引用块，或变为列表项。
  - _参数_: `wrappers` 是一组 `{type, attrs}` 对象，通常通过 `findWrapping` 辅助函数计算得出。
- **`split(pos, depth, typesAfter)`**:
  - _作用_: 在指定位置分割节点。
  - _场景_: 按下回车键时，将一个段落切成两个。
- **`join(pos, depth)`**:
  - _作用_: 合并两个相邻的同类型节点。
  - _场景_: 按下 Delete 键删除段落间的换行符时。

## 23. 节点类型判断详解 (Node Flags)

`Node` 和 `NodeType` 上有一系列布尔属性，容易混淆，需明确区分。

- **`isBlock`**: 是否为块级节点（如 `paragraph`, `heading`）。
- **`isInline`**: 是否为内联节点（如 `text`, `image`）。
- **`isText`**: 是否为纯文本节点。
- **`isLeaf`**: 是否为叶子节点（不允许有子内容）。`text` 和 `atom: true` 的节点都是叶子。
- **`isAtom`**: 是否为原子节点。原子节点在编辑器中被视为一个整体，光标无法进入其内部（除非定义了 NodeView 且处理了光标）。
- **`isTextblock`**: 是否为文本块（即直接包含 inline 内容的 block，如 `paragraph`）。只有 Textblock 才能包含光标。

## 24. Plugin Props 的强大之处

`Plugin` 不仅可以管理 State，还可以通过 `props` 属性直接向 `EditorView` 注入 Props。这是插件修改编辑器行为的主要方式。

- **拦截交互**:
  ```javascript
  new Plugin({
    props: {
      handleKeyDown(view, event) {
        if (event.key === 'Enter') {
          // 拦截回车键
          return true
        }
      },
      handleDOMEvents: {
        mouseover(view, event) {
          // 监听鼠标悬停
        }
      }
    }
  })
  ```
- **渲染装饰器**: `decorations(state)` 是最常用的 prop，用于动态高亮。
- **自定义节点视图**: 插件甚至可以注入 `nodeViews`，允许插件自带特定节点的渲染逻辑。

## 25. DOMParser 与 Slice 解析

在处理粘贴或拖拽时，我们经常需要手动解析 HTML 片段。

- **`DOMParser.parseSlice(dom, options)`**:
  - _作用_: 将 DOM 节点解析为 `Slice`，而不是完整的 `Node`。
  - _场景_: 粘贴处理。`Slice` 保留了“开放深度”信息（openStart/openEnd），这对于保持粘贴内容的上下文结构至关重要。
  - _参数_: `preserveWhitespace: true` 可以保留 HTML 中的换行和空格（默认会折叠）。

## 26. 常用生态库推荐 (Utils)

虽然不是核心 API，但这些库是实际开发中的“标准库”。

- **`prosemirror-utils`**:
  - 提供了大量查找和操作节点的辅助函数，如 `findParentNode`, `setTextSelection`, `removeSelectedNode`。强烈建议使用，避免重复造轮子。
- **`prosemirror-keymap`**:
  - 几乎必装，用于绑定快捷键。
- **`prosemirror-dropcursor`**:
  - 拖拽时显示插入位置的指示线。
- **`prosemirror-gapcursor`**:
  - 允许光标聚焦到两个块级节点之间。

## 27. Transaction Meta 约定

`tr.setMeta(key, value)` 用于在插件间通信。有一些约定俗成的 Key：

- **`"addToHistory"`**: (Boolean) 控制是否加入撤销栈。
- **`"pointer"`**: (Boolean) 标识该事务是由鼠标点击触发的（用于处理选区更新）。
- **`"paste"`**: (Boolean) 标识该事务是由粘贴触发的。
- **自定义 Key**: 建议使用 Plugin 实例本身作为 Key (`pluginKey`)，以避免命名冲突。
  ```javascript
  const myPluginKey = new PluginKey('myPlugin')
  tr.setMeta(myPluginKey, { someData: 123 })
  // 在插件的 apply 中读取
  const meta = tr.getMeta(myPluginKey)
  ```

## 28. Fragment 高级操作与 Diff 算法

`Fragment` 是 ProseMirror 文档树的基石，除了基本的遍历，它还提供了一些用于高效更新的 Diff 方法。

- **`findDiffStart(other: Fragment, pos: number)`**:
  - _作用_: 比较当前 Fragment 和另一个 Fragment，返回第一个不同节点的起始位置。
  - _场景_: 优化渲染。当你需要知道从哪里开始更新 DOM 时。
- **`findDiffEnd(other: Fragment, pos: number, otherPos: number)`**:
  - _作用_: 返回最后一个不同节点的结束位置。
  - _结合使用_: 通过 `findDiffStart` 和 `findDiffEnd`，你可以精确锁定变化的范围，从而只更新这一小部分 DOM，而不是重绘整个文档。
- **`forEach(callback)`**:
  - _细节_: 遍历的是直接子节点。如果需要深度遍历，请使用 `nodesBetween`。

## 29. ReplaceAroundStep (高级替换)

这是 `ReplaceStep` 的变体，也是 ProseMirror 中最难理解的 Step 之一。

- **作用**: 替换一个范围的内容，但**保留**中间的某个切片结构。
- **场景**: 比如你有一个 `<blockquote><p>Text</p></blockquote>`，你想把它变成 `<div class="wrapper"><p>Text</p></div>`。
  - 如果用普通的 `ReplaceStep`，你需要删除整个 blockquote 再插入 div，这会导致 `<p>` 节点被销毁重建（光标丢失，React 组件卸载）。
  - 使用 `ReplaceAroundStep`，你可以只替换“外壳”（blockquote -> div），而保留内部的 `<p>` 节点对象不变。
- **生成**: 通常不需要手动创建，`lift` 和 `wrap` 等高级变换会自动生成它。

## 30. Mark 的相等性判断细节

在处理样式切换时，判断两个 Mark 是否“相同”有细微差别。

- **`mark.eq(otherMark)`**:
  - _严格相等_: 类型相同，且**所有属性** (`attrs`) 都深度相等。
  - _场景_: 判断当前选区是否已经应用了完全一致的样式（包括颜色值等）。
- **`mark.type.isInSet(marks)`**:
  - _类型存在_: 只检查数组中是否有该类型的 Mark，不关心属性。
  - _场景_: 判断“这里是否有链接”，不管链接地址是什么。
- **`mark.addToSet(marks)`**:
  - _智能合并_: 如果数组中已有同类型 Mark：
    - 如果属性完全相同，不重复添加。
    - 如果属性不同，通常会替换旧的（取决于 Schema 定义的 `excludes`）。

## 31. NodeRange (节点范围)

`Selection` 对象中经常会用到 `NodeRange`，它描述了同一层级下的一系列兄弟节点。

- **`$from.blockRange($to)`**:
  - _作用_: 尝试获取包含 `$from` 到 `$to` 的最小块级范围。
  - _返回_: `NodeRange` 对象，包含 `start`, `end`, `depth`, `parent`。
  - _场景_: 用户选中了两个列表项，你想对这两个列表项整体执行操作（如缩进）。如果选区跨越了不同层级（如从列表内选到了列表外），这个方法可能返回 `null`。

## 32. 坐标系统与命中测试细节

`view.posAtCoords` 和 `view.coordsAtPos` 是处理鼠标交互的核心。

- **`posAtCoords({ left, top })`**:
  - _返回_: `{ pos, inside }`。
  - **`pos`**: 最接近该坐标的文档位置。
  - **`inside`**: 如果坐标落在某个节点内部（而不是两个节点之间），返回该节点的起始位置。这对于判断点击的是哪个具体节点（如点击了哪个单元格）非常有用。
- **`coordsAtPos(pos, side)`**:
  - _参数 `side`_:
    - `-1`: 获取位置左侧字符的右边缘坐标。
    - `1`: 获取位置右侧字符的左边缘坐标。
    - _场景_: 当光标在行尾换行处时，位置相同，但视觉坐标可能在上一行末尾，也可能在下一行开头。`side` 参数决定了你获取的是哪一个。

## 33. 序列化与解析的上下文 (Context)

在 `DOMParser` 解析 HTML 时，上下文至关重要。

- **`parseSlice(dom, { context })`**:
  - _问题_: 如果你解析一个 `<li>Item</li>` 字符串，默认情况下 ProseMirror 可能会把它丢弃，因为它不能直接放在 `doc` 根节点下。
  - _解决_: 提供 `context: "ul"`（或对应的节点对象），告诉解析器：“假设这个片段是放在 `<ul>` 里的”。这样 `<li>` 就会被正确保留。

## 34. 错误处理与恢复

- **`view.setProps({ handleDOMEvents })`**:
  - 可以在这里拦截 `error` 事件，防止编辑器因插件报错而白屏。
- **`State.reconfigure(config)`**:
  - 可以在运行时动态重新配置插件列表（例如用户在设置里关闭了某个插件）。

## 35. JSON 序列化与反序列化

ProseMirror 的文档是可以完全序列化为 JSON 的，这是持久化存储的基础。

- **`node.toJSON()`**:
  - _作用_: 将节点转换为 JSON 对象。
  - _默认行为_: 递归转换。输出格式通常包含 `type`, `attrs`, `content` (数组)。
  - _自定义_: 可以在 Node Spec 中定义 `toDebugString` (仅用于调试) 或覆盖 `toJSON` (极少需要)。
- **`schema.nodeFromJSON(json)`**:
  - _作用_: 将 JSON 对象还原为 Node 实例。
  - _注意_: 必须使用生成该 JSON 的同一套 Schema（或兼容的 Schema），否则会报错。
- **`mark.toJSON()` / `schema.markFromJSON(json)`**: 同理。

## 36. 视图焦点与滚动控制

- **`view.focus()`**:
  - _作用_: 让编辑器获得焦点。
  - _注意_: 在某些浏览器中，如果编辑器不在视口内，调用 `focus()` 可能会导致页面突然滚动。
- **`tr.scrollIntoView()`**:
  - _作用_: 标记该 Transaction 应用后，视图应该滚动以确保光标（或选区）可见。
  - _场景_: 插入内容后、搜索跳转后。
  - _原理_: 这是一个副作用标记，不会改变文档内容。

## 37. ResolvedPos 深度对比与导航

`ResolvedPos` (`$pos`) 提供了比较两个位置关系的强大方法。

- **`$pos.sharedDepth(otherPos)`**:
  - _作用_: 返回两个位置共有的最大深度。
  - _场景_: 判断两个光标是否在同一个段落内 (`sharedDepth >= blockDepth`)，或者是否在同一个列表项内。
- **`$pos.sameParent(otherPos)`**:
  - _作用_: 判断两个位置是否拥有同一个直接父节点。
- **`$pos.min(otherPos)` / `$pos.max(otherPos)`**:
  - _作用_: 返回较前或较后的那个 ResolvedPos 对象。

## 38. Decoration Spec 的妙用

创建 Decoration 时，第三个参数 `spec` 往往被忽略，但它非常有用。

- **过滤与查找**:
  ```javascript
  const deco = Decoration.inline(
    from,
    to,
    { class: 'error' },
    { id: 'err_123', type: 'spellcheck' }
  )
  // 后续可以通过 spec 查找
  const found = decorationSet.find(null, null, spec => spec.type === 'spellcheck')
  ```
- **`inclusiveStart` / `inclusiveEnd`**:
  - _作用_: 控制当用户在 Decoration 边缘输入时，新输入的文字是否包含在 Decoration 内（类似 Mark 的 `inclusive`）。

## 39. Node 完整性验证

- **`node.check()`**:
  - _作用_: 验证节点结构是否符合 Schema 定义。如果不符合，会抛出错误。
  - _场景_: 调试时，或者在处理不可信的外部数据构造 Node 后进行防御性检查。
