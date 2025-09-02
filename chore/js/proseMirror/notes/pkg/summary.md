- prosemirror-model 为 ProseMirror 编辑器提供了一个健壮、结构化的文档模型。它通过 Schema 保证了文档的合法性，并通过 Parser 和 Serializer 实现了模型与浏览器 DOM 之间的解耦和双向转换，这是 ProseMirror 强大功能和灵活性的基础。
- prosemirror-transform 提供了一套强大而精确的文档变更描述和操作机制。它将复杂的文档修改分解为简单、可控的步骤，是 ProseMirror 实现其声明式、可预测状态管理和高级功能的基石。
- prosemirror-state 通过其不可变的 EditorState 对象和基于 Transaction 的更新机制，为编辑器提供了一个集中、可预测且易于扩展的状态管理核心。这种设计是实现撤销/重做、协同编辑等高级功能的基石。
- prosemirror-view 通过 EditorView 作为控制器，利用 ViewDesc 这一抽象层，高效地将 EditorState 同步到真实的 DOM 上。它通过 NodeView 和 Decoration 提供了强大的自定义和扩展能力，使得开发者可以构建出功能丰富、交互复杂的编辑器。同时，它也封装了大量处理 contenteditable 跨浏览器兼容性问题的复杂逻辑，让开发者可以更专注于业务功能的实现。

---

- prosemirror-history 通过将变更分解为带映射信息的 Item，并利用 prosemirror-transform 的映射（Mapping）和变换（Transform）能力，实现了一个强大且支持协同编辑的选择性历史记录系统。
- prosemirror-collab 是一个 ProseMirror 插件，它为编辑器实现协同编辑功能提供了一个健壮的框架。它本身不处理网络通信，而是提供了一套机制来管理本地和远程的变更，使得将 ProseMirror 集成到任何协同编辑后端成为可能。

---

- prosemirror-commands 提供了一系列标准的、可组合的“命令”函数，这些函数是连接用户意图（如点击按钮、按快捷键）与编辑器状态变更的桥梁。其“查询/执行”双重角色设计，使得 UI 更新与逻辑执行可以共享同一套代码。
- prosemirror-keymap 提供了一个插件，用于将键盘快捷键（如 `Mod-B`）映射到特定的命令。它通过标准化的按键名称实现了跨平台兼容，并通过插件顺序实现了快捷键的优先级和覆盖。
- prosemirror-inputrules 提供了一个插件，用于实现“输入时自动格式化”（如输入 `## ` 变成标题）。它通过匹配用户输入的文本模式，自动触发相应的文档转换。
- prosemirror-dropcursor 和 prosemirror-gapcursor 是两个提升用户体验的 UI 插件。`dropcursor` 在拖放操作时显示一个精确的插入位置指示器；`gapcursor` 则解决了在节点之间无法放置光标的“死胡同”问题，确保文档中处处可编辑。
- prosemirror-menu 提供了一套基础的、由状态驱动的菜单栏 UI 组件。它是一个极佳的范例，展示了如何将编辑器的命令和状态与 UI 元素的可用性/激活状态进行绑定。
- prosemirror-markdown 充当了 ProseMirror 结构化文档与 Markdown 纯文本之间的“翻译官”。它利用 `markdown-it` 库，提供了高度可定制的解析器（Markdown -> ProseMirror）和序列化器（ProseMirror -> Markdown）。
- prosemirror-schema-basic, prosemirror-schema-list, prosemirror-tables 提供了构建编辑器 Schema 所需的节点和标记规格。它们体现了 ProseMirror 模块化、可组合的 Schema 设计思想，允许开发者像搭积木一样构建出符合自己需求的文档结构。
- prosemirror-changeset 提供了一个结构化的文档比较工具。它能精确地找出两个文档版本之间的增、删、改，是实现“追踪修订”或“差异对比”等高级功能的基础。
- prosemirror-test-builder 是一个开发工具，它通过“带标签的模板字符串”极大地简化了测试用例的编写，让开发者能以所见即所得的方式快速构建复杂的测试状态。

---

### ProseMirror 架构关系深度解析

ProseMirror 的架构设计极其优雅，它将一个复杂的富文本编辑器分解为多个职责单一、低耦合、高内聚的模块。我们可以将这些模块想象成一个层次分明的金字塔结构，从底层的数据模型到顶层的用户界面。

#### **第一层：数据核心 (The Core Data Layer)**

这是整个架构的基石，定义了“文档是什么”以及“如何安全地改变它”。

1.  **`prosemirror-model`**: **【宪法】**
    这是最底层。它定义了文档的结构（`Schema`）、内容（`Node`, `Fragment`）和剪贴板数据（`Slice`）。它像一部宪法，规定了文档的一切合法形态，保证了数据的结构化和有效性。**所有其他模块都建立在 `prosemirror-model` 提供的不可变数据结构之上。**

2.  **`prosemirror-transform`**: **【立法机构】**
    它建立在 `model` 之上，定义了**如何修改文档**。它将变更抽象为原子化的、可逆的 `Step`，并提供了强大的位置映射（`Mapping`）能力。它就像立法机构，制定了所有变更的“法律程序”。没有它，撤销/重做和协同编辑都无从谈起。

3.  **`prosemirror-state`**: **【中央政府】**
    这是数据核心的“大脑”。它将 `doc` (来自 `model`) 和变更逻辑 (来自 `transform`) 包装成一个完整的、不可变的 `EditorState` 快照。它引入了 `Transaction`（继承自 `Transform`）作为状态演变的唯一途径，并通过 `Plugin` 系统提供了强大的扩展能力。它像中央政府，统管着编辑器的所有状态，并为其他功能模块（插件）提供了运行的土壤。

#### **第二层：功能与交互逻辑 (The Feature & Interaction Layer)**

这一层是各种插件，它们作为“部门”，寄生在 `prosemirror-state` 的 `Plugin` 系统中，为编辑器添加具体功能。

4.  **`prosemirror-history`**: **【历史档案馆】**
    它是一个插件，利用 `transform` 的 `step.invert()` 能力，将事务的逆操作存入 `done` 栈，实现了撤销/重做。它完全依赖 `state` 的事务机制和 `transform` 的可逆性。

5.  **`prosemirror-collab`**: **【外交部】**
    它也是一个插件，利用 `transform` 的 `rebase` 和 `map` 能力，处理本地与远程的 `Step` 冲突，是实现协同编辑的客户端框架。它与 `history` 插件能良好协作，共同处理复杂的变更历史。

6.  **`prosemirror-commands`**: **【行政命令集】**
    它提供了一系列标准化的函数（`Command`），用于执行具体的编辑操作。这些命令是连接上层 UI 与底层 `state` 变更的“标准公文”。它们创建 `Transaction`，但自己不分发，而是交由调用者（如 `keymap` 或 `menu`）处理。

7.  **`prosemirror-keymap` & `prosemirror-inputrules`**: **【自动化办公系统】**
    这两个插件是“事件监听器”，它们将用户的物理输入（键盘按键、文本模式）翻译成对 `commands` 的调用。`keymap` 负责快捷键，`inputrules` 负责输入时自动格式化。它们是实现高效编辑体验的关键。

#### **第三层：视图与用户界面 (The View & UI Layer)**

这一层负责将抽象的 `EditorState` 渲染为用户能看到和交互的 DOM，并处理来自浏览器的原始事件。

8.  **`prosemirror-view`**: **【渲染引擎与事件中心】**
    这是连接抽象状态与具体 DOM 的桥梁。它接收 `EditorState`，通过高效的 diff/patch 算法更新 DOM。同时，它捕获所有浏览器事件（点击、输入、拖拽），并将它们翻译成 `Transaction`，从而完成数据流的闭环。**它是唯一直接与 DOM 交互的核心模块**。

9.  **`prosemirror-dropcursor` & `prosemirror-gapcursor`**: **【UI 体验优化器】**
    这两个插件直接与 `view` 配合，通过 `Decoration` 系统在视图层“画”出视觉提示，而完全不影响 `model` 和 `state`。它们是典型的视图层插件，专门用于提升交互体验。

10. **`prosemirror-menu`**: **【官方 UI 套件（示例）】**
    它是一个建立在 `view` 和 `commands` 之上的 UI 示例。它展示了如何读取 `EditorState` 来更新按钮状态（通过调用命令的“查询模式”），以及如何通过 `dispatch` 来执行命令，是构建自定义 UI 的学习典范。

#### **第四层：生态与工具 (The Ecosystem & Tooling Layer)**

这些模块不直接参与编辑器的核心运行时，但为开发、测试和数据交换提供了便利。

11. **Schema 定义模块 (`-schema-basic`, `-schema-list`, `-schema-tables`)**: **【词汇表】**
    它们为 `prosemirror-model` 的 `Schema` 提供了预设的“词汇”（节点和标记规格），让开发者可以快速搭建文档结构。

12. **数据交换模块 (`prosemirror-markdown`)**: **【翻译官】**
    它在 ProseMirror 的 `model` 和外部数据格式（Markdown）之间进行双向转换，是编辑器与外部世界沟通的桥梁。

13. **辅助工具 (`prosemirror-changeset`, `prosemirror-test-builder`)**: **【开发工具箱】**
    `changeset` 用于比较两个 `model` 实例的差异；`test-builder` 则简化了测试中构建 `state` 的过程。它们是开发和测试高级功能时的得力助手。

**总结关系图:**

```
+-------------------------------------------------------------------+
|                     第四层: 生态与工具 (Ecosystem & Tooling)          |
|  [menu] [markdown] [schema-list] [test-builder] [changeset] ...   |
+-------------------------------------------------------------------+
      ^                                                 ^
      | (提供UI/数据交换/Schema定义)                        | (用于开发测试)
      v                                                 v
+-------------------------------------------------------------------+
|                     第三层: 视图与UI (View & UI)                    |
|                                                                   |
|   +-----------------------+      +------------------------------+   |
|   | prosemirror-view      | <--> | dropcursor, gapcursor, etc.  |   |
|   | (渲染引擎, 事件中心)    |      | (UI体验插件)                 |   |
|   +-----------------------+      +------------------------------+   |
+-------------------------------------------------------------------+
      ^ (Dispatch(tr))                                  ^ (Decorations)
      |                                                 |
      v (State -> DOM)                                  v (State -> UI)
+-------------------------------------------------------------------+
|                  第二层: 功能与交互逻辑 (Features & Interaction)    |
|                                                                   |
|   [history] [collab] [keymap] [inputrules] [commands] ...         |
|   (所有这些都是在 Plugin 系统上运行的 "部门")                       |
|                                                                   |
+-------------------------------------------------------------------+
      ^ (作为插件运行在 State 之上)
      |
      v (通过 Transaction 与 State 交互)
+-------------------------------------------------------------------+
|                     第一层: 数据核心 (Core Data)                    |
|                                                                   |
|   +-----------------------+ (包含并控制)                            |
|   | prosemirror-state     |                                         |
|   | (中央政府: EditorState) |                                         |
|   +-----------------------+                                         |
|             ^                                                       |
|             | (使用)                                                |
|   +-----------------------+ (基于)  +-----------------------------+ |
|   | prosemirror-transform | ----> | prosemirror-model           | |
|   | (立法机构: Step, Map) |       | (宪法: Schema, Node)        | |
|   +-----------------------+       +-----------------------------+ |
+-------------------------------------------------------------------+
```
