好的，我们来对 `prosemirror-example-setup` 进行一次彻底的、深入的剖析。这个包在 ProseMirror 生态中扮演着一个非常特殊的角色：它本身不提供任何新的核心功能，但它却是**几乎所有初学者接触 ProseMirror 的第一站**。

`prosemirror-example-setup` 的本质是一个**“最佳实践集合”**和**“快速启动器”**。它将构建一个功能齐全的基础编辑器所需的各种插件和配置项打包成一个单一、易于使用的函数。

我们将从以下四个关键角度来解构它，并理解为什么在生产项目中你最终可能会“抛弃”它：

1.  **`exampleSetup` 函数：它到底做了什么？**
2.  **配置的“大礼包”：包含哪些插件？**
3.  **Schema 的角色：基础与扩展**
4.  **从“示例”到“生产”：为什么以及如何“毕业”**

---

### 1. `exampleSetup` 函数：它到底做了什么？

`exampleSetup` 是这个包暴露的唯一核心函数。它的签名如下：

```typescript
function exampleSetup(options: {
  schema: Schema
  mapKeys?: object
  menuBar?: boolean
  history?: boolean
  floatingMenu?: boolean
  menuContent?: MenuItem[][]
}): Plugin[]
```

它的作用非常直接：接收一个配置对象，返回一个**`Plugin` 数组**。这个数组包含了构建一个“开箱即用”的编辑器所需的所有基础插件。

**使用方法**:
你只需将这个函数返回的插件数组，在创建 `EditorState` 时传入 `plugins` 字段即可。

```typescript
import { EditorState } from 'prosemirror-state'
import { schema } from 'prosemirror-schema-basic'
import { exampleSetup } from 'prosemirror-example-setup'

const state = EditorState.create({
  schema,
  plugins: exampleSetup({ schema }) // 核心用法
})
```

---

### 2. 配置的“大礼包”：包含哪些插件？

调用 `exampleSetup({ schema })` 就像是打开了一个包含了 ProseMirror “全家桶”的礼包。它为你自动配置并整合了以下核心模块的插件：

- **`prosemirror-inputrules`**:

  - **功能**: 实现输入时自动转换。
  - **具体配置**:
    - 输入 `## ` 转换成二级标题。
    - 输入 `* ` 或 `- ` 转换成无序列表。
    - 输入 `1. ` 转换成有序列表。
    - 输入 `>` 转换成引用块。
    - 输入 `` `...` `` 转换成代码块。

- **`prosemirror-keymap`**:

  - **功能**: 绑定核心的键盘快捷键。
  - **具体配置**:
    - `Enter`: 调用 `splitBlock` 来创建新段落。
    - `Backspace` 和 `Delete`: 调用 `chainCommands` 组合的删除和合并命令。
    - `Mod-b`, `Mod-i`: 调用 `toggleMark` 来切换粗体/斜体。
    - `Mod-z` (撤销) 和 `Mod-y` (重做): 绑定到 `prosemirror-history` 的命令。
    - 以及大量其他符合直觉的文本编辑快捷键。
  - 它还包含了 `prosemirror-commands` 中的 `baseKeymap`，这是所有平台通用的基础按键绑定。

- **`prosemirror-history`**:

  - **功能**: 提供撤销和重做功能。
  - **具体配置**: 默认启用，可以通过 `options.history = false` 来禁用。

- **`prosemirror-commands`**:

  - **功能**: 提供基础的编辑命令。
  - **具体配置**: `exampleSetup` 内部大量使用了 `prosemirror-commands` 提供的命令来配置 `keymap`。

- **`prosemirror-dropcursor`**:

  - **功能**: 在拖放内容时，显示一个光标来指示将被放置的位置。

- **`prosemirror-gapcursor`**:

  - **功能**: 在那些无法放置普通文本光标的位置（如两个表格之间）提供一个可见的“间隙光标”。

- **`prosemirror-menu`** (可选):
  - **功能**: 提供一个简单的、可配置的工具栏。
  - **具体配置**:
    - 如果 `options.menuBar` 不为 `false`，它会创建一个 `menuBar` 插件，在编辑器顶部显示一个工具栏。
    - 工具栏的内容由 `options.menuContent` 决定，`exampleSetup` 内部会根据传入的 `schema` 自动生成一套合理的默认菜单项（如加粗、斜体、标题、列表等）。

---

### 3. Schema 的角色：基础与扩展

`exampleSetup` 的配置高度依赖于你传入的 `schema`。它会检查 `schema` 中是否存在特定的节点和标记类型，来决定是否启用相关的快捷键和菜单项。

例如：

- 如果你的 `schema` 中包含了 `schema.marks.strong`，那么 `Mod-b` 快捷键和工具栏中的“加粗”按钮才会被启用。
- 如果你的 `schema` 中包含了 `schema.nodes.bullet_list`，那么输入 `* ` 自动转换为列表的 `inputRule` 和相关的菜单项才会被添加。

这意味着 `exampleSetup` 具有一定的**自适应性**。它不仅仅是为 `prosemirror-schema-basic` 服务的，如果你使用 `prosemirror-schema-list` 扩展了你的 `schema`，`exampleSetup` 会自动识别并添加列表相关的快捷键和菜单。

---

### 4. 从“示例”到“生产”：为什么以及如何“毕业”

`prosemirror-example-setup` 是一个极好的学习工具和项目启动器，但当你的项目变得复杂时，你很可能会遇到它的局限性，并需要从它“毕业”。

#### a. 为什么要“毕业”？

1.  **定制化需求**: `exampleSetup` 的配置是固定的。如果你想修改一个快捷键、禁用某个 `inputRule`，或者彻底改变菜单栏的样式和行为，你会发现直接修改 `exampleSetup` 的源码或者在其外部覆盖它的配置非常困难和混乱。
2.  **包体积优化 (Bundle Size)**: `exampleSetup` 引入了 `prosemirror-menu` 等可能在你的最终产品中并不需要的依赖。为了精细化控制打包体积，你需要手动挑选你真正需要的模块。
3.  **UI 框架集成**: `prosemirror-menu` 是一个纯粹的、基于 DOM 的菜单实现。在现代前端项目中，你几乎肯定会使用 React, Vue, Svelte 等框架来构建你的 UI。在这种情况下，`prosemirror-menu` 就显得格格不入，你需要自己构建与框架集成的工具栏。
4.  **代码清晰度与可维护性**: 当你的编辑器逻辑变得复杂时，一个单一的、黑盒的 `exampleSetup` 函数会掩盖很多实现细节。将插件配置显式地、模块化地写在你的代码中，会大大提高项目的可读性和可维护性。

#### b. 如何“毕业”？

“毕业”的过程其实就是**“手动实现 `exampleSetup`”**。

你可以直接查看 `prosemirror-example-setup` 的[源代码](https://github.com/ProseMirror/prosemirror-example-setup/blob/master/src/index.js)。它的代码非常简短和清晰。

“毕业”的步骤如下：

1.  **移除 `prosemirror-example-setup` 依赖**。
2.  **在你的 `plugins` 数组中，手动添加你需要的插件**：

    ```typescript
    import { history } from 'prosemirror-history'
    import { keymap } from 'prosemirror-keymap'
    import { baseKeymap } from 'prosemirror-commands'
    import { dropCursor } from 'prosemirror-dropcursor'
    import { gapCursor } from 'prosemirror-gapcursor'
    import { buildInputRules } from './my-input-rules' // 你自己的输入规则
    import { buildKeymap } from './my-keymap' // 你自己的快捷键映射

    const plugins = [
      buildInputRules(schema),
      keymap(buildKeymap(schema)),
      keymap(baseKeymap),
      dropCursor(),
      gapCursor(),
      history()
      // ... 你自己的其他插件，比如用于自定义菜单的插件
    ]
    ```

3.  **创建你自己的 `buildKeymap` 和 `buildInputRules` 函数**。你可以直接从 `prosemirror-example-setup` 的源码中复制粘贴开始，然后根据你的需求进行修改、添加或删除。
4.  **构建你自己的 UI**。使用你选择的前端框架，构建工具栏、菜单等。在 UI 组件中，调用 `prosemirror-commands` 提供的命令来与编辑器交互，并根据 `EditorState` 来更新 UI 的激活状态。

### 总结

`prosemirror-example-setup` 是一个优秀的“脚手架”和“导师”。它向初学者展示了如何将 ProseMirror 生态中的各个部分组合起来，构建一个功能完备的编辑器。但它被设计为**一个起点，而非终点**。理解它的内部构成，并学会在合适的时机“毕业”，手动接管插件的配置，是每一个 ProseMirror 开发者从入门到精通的必经之路。`...` ``` 转换成代码块。

- **`prosemirror-keymap`**:

  - **功能**: 绑定核心的键盘快捷键。
  - **具体配置**:
    - `Enter`: 调用 `splitBlock` 来创建新段落。
    - `Backspace` 和 `Delete`: 调用 `chainCommands` 组合的删除和合并命令。
    - `Mod-b`, `Mod-i`: 调用 `toggleMark` 来切换粗体/斜体。
    - `Mod-z` (撤销) 和 `Mod-y` (重做): 绑定到 `prosemirror-history` 的命令。
    - 以及大量其他符合直觉的文本编辑快捷键。
  - 它还包含了 `prosemirror-commands` 中的 `baseKeymap`，这是所有平台通用的基础按键绑定。

- **`prosemirror-history`**:

  - **功能**: 提供撤销和重做功能。
  - **具体配置**: 默认启用，可以通过 `options.history = false` 来禁用。

- **`prosemirror-commands`**:

  - **功能**: 提供基础的编辑命令。
  - **具体配置**: `exampleSetup` 内部大量使用了 `prosemirror-commands` 提供的命令来配置 `keymap`。

- **`prosemirror-dropcursor`**:

  - **功能**: 在拖放内容时，显示一个光标来指示将被放置的位置。

- **`prosemirror-gapcursor`**:

  - **功能**: 在那些无法放置普通文本光标的位置（如两个表格之间）提供一个可见的“间隙光标”。

- **`prosemirror-menu`** (可选):
  - **功能**: 提供一个简单的、可配置的工具栏。
  - **具体配置**:
    - 如果 `options.menuBar` 不为 `false`，它会创建一个 `menuBar` 插件，在编辑器顶部显示一个工具栏。
    - 工具栏的内容由 `options.menuContent` 决定，`exampleSetup` 内部会根据传入的 `schema` 自动生成一套合理的默认菜单项（如加粗、斜体、标题、列表等）。

---

### 3. Schema 的角色：基础与扩展

`exampleSetup` 的配置高度依赖于你传入的 `schema`。它会检查 `schema` 中是否存在特定的节点和标记类型，来决定是否启用相关的快捷键和菜单项。

例如：

- 如果你的 `schema` 中包含了 `schema.marks.strong`，那么 `Mod-b` 快捷键和工具栏中的“加粗”按钮才会被启用。
- 如果你的 `schema` 中包含了 `schema.nodes.bullet_list`，那么输入 `* ` 自动转换为列表的 `inputRule` 和相关的菜单项才会被添加。

这意味着 `exampleSetup` 具有一定的**自适应性**。它不仅仅是为 `prosemirror-schema-basic` 服务的，如果你使用 `prosemirror-schema-list` 扩展了你的 `schema`，`exampleSetup` 会自动识别并添加列表相关的快捷键和菜单。

---

### 4. 从“示例”到“生产”：为什么以及如何“毕业”

`prosemirror-example-setup` 是一个极好的学习工具和项目启动器，但当你的项目变得复杂时，你很可能会遇到它的局限性，并需要从它“毕业”。

#### a. 为什么要“毕业”？

1.  **定制化需求**: `exampleSetup` 的配置是固定的。如果你想修改一个快捷键、禁用某个 `inputRule`，或者彻底改变菜单栏的样式和行为，你会发现直接修改 `exampleSetup` 的源码或者在其外部覆盖它的配置非常困难和混乱。
2.  **包体积优化 (Bundle Size)**: `exampleSetup` 引入了 `prosemirror-menu` 等可能在你的最终产品中并不需要的依赖。为了精细化控制打包体积，你需要手动挑选你真正需要的模块。
3.  **UI 框架集成**: `prosemirror-menu` 是一个纯粹的、基于 DOM 的菜单实现。在现代前端项目中，你几乎肯定会使用 React, Vue, Svelte 等框架来构建你的 UI。在这种情况下，`prosemirror-menu` 就显得格格不入，你需要自己构建与框架集成的工具栏。
4.  **代码清晰度与可维护性**: 当你的编辑器逻辑变得复杂时，一个单一的、黑盒的 `exampleSetup` 函数会掩盖很多实现细节。将插件配置显式地、模块化地写在你的代码中，会大大提高项目的可读性和可维护性。

#### b. 如何“毕业”？

“毕业”的过程其实就是**“手动实现 `exampleSetup`”**。

你可以直接查看 `prosemirror-example-setup` 的[源代码](https://github.com/ProseMirror/prosemirror-example-setup/blob/master/src/index.js)。它的代码非常简短和清晰。

“毕业”的步骤如下：

1.  **移除 `prosemirror-example-setup` 依赖**。
2.  **在你的 `plugins` 数组中，手动添加你需要的插件**：

    ```typescript
    import { history } from 'prosemirror-history'
    import { keymap } from 'prosemirror-keymap'
    import { baseKeymap } from 'prosemirror-commands'
    import { dropCursor } from 'prosemirror-dropcursor'
    import { gapCursor } from 'prosemirror-gapcursor'
    import { buildInputRules } from './my-input-rules' // 你自己的输入规则
    import { buildKeymap } from './my-keymap' // 你自己的快捷键映射

    const plugins = [
      buildInputRules(schema),
      keymap(buildKeymap(schema)),
      keymap(baseKeymap),
      dropCursor(),
      gapCursor(),
      history()
      // ... 你自己的其他插件，比如用于自定义菜单的插件
    ]
    ```

3.  **创建你自己的 `buildKeymap` 和 `buildInputRules` 函数**。你可以直接从 `prosemirror-example-setup` 的源码中复制粘贴开始，然后根据你的需求进行修改、添加或删除。
4.  **构建你自己的 UI**。使用你选择的前端框架，构建工具栏、菜单等。在 UI 组件中，调用 `prosemirror-commands` 提供的命令来与编辑器交互，并根据 `EditorState` 来更新 UI 的激活状态。

### 总结

`prosemirror-example-setup` 是一个优秀的“脚手架”和“导师”。它向初学者展示了如何将 ProseMirror 生态中的各个部分组合起来，构建一个功能完备的编辑器。但它被设计为**一个起点，而非终点**。理解它的内部构成，并学会在合适的时机“毕业”，手动接管插件的配置，是每一个 ProseMirror 开发者从入门到精通的必经之路。
