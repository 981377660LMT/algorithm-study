## 1. ProseMirror 中，model、state、view、transform 能力是什么，依赖是怎么样的，如何深刻理解

ProseMirror 是一个富文本编辑器框架，它通过一系列精心设计的模块来提供其功能。理解这些模块及其之间的关系对于深刻理解 ProseMirror 至关重要。这些核心模块包括 Model、State、View 和 Transform。

### Model（模型）

- **能力**：Model 定义了编辑器的文档模型，包括文档的结构、节点（如段落、标题）、标记（如加粗、斜体）等。它是 ProseMirror 编辑器的基础，用于表示文档的内容和结构。
- **依赖**：Model 是基础模块，不依赖其他模块。

### State（状态）

- **能力**：State 表示编辑器的当前状态，包括当前的文档内容（由 Model 定义）、选择区域、插件状态等。State 是不可变的，任何状态的变更都会产生一个新的 State 实例。
- **依赖**：State 依赖于 Model，因为它需要 Model 来表示文档的内容。

### View（视图）

- **能力**：View 负责将 State 渲染到浏览器中，并处理用户的输入（如键盘事件、鼠标事件），将这些输入转换为对文档的操作。View 是 State 和用户之间的桥梁。
- **依赖**：View 依赖于 State，因为它需要根据 State 来渲染编辑器的界面。同时，View 也会影响 State，因为用户的操作会导致 State 的变更。

### Transform（变换）

- **能力**：Transform 提供了一种方式来描述和应用对文档的变更，如插入文本、删除节点等。这些变更通过事务（Transaction）来表示，每个事务可以包含多个步骤（Step），每个步骤描述了一个具体的变更操作。
- **依赖**：Transform 依赖于 Model，因为它需要知道文档的结构来正确地应用变更。同时，Transform 的结果（事务）通常会被应用到 State 上，导致 State 的更新。

### 理解它们的关系

- **Model**是一切的基础，定义了文档的结构和内容。
- **State**在 Model 的基础上增加了当前的编辑状态，如选择区域、插件状态等。
- **View**根据 State 来渲染编辑器界面，并将用户的操作转换为对 State 的变更。
- **Transform**提供了一种机制来描述和应用对 Model 的变更，这些变更通过 State 来管理和追踪。

---

## 2. proseMirrot 中有哪些核心模块，哪些扩展模块

### 核心模块

核心模块是 ProseMirror 的基础，提供了编辑器的主要功能。

1. **prosemirror-model**：定义了编辑器的文档模型，包括节点（Node）和标记（Mark）的概念。这个模块是处理文档结构的基础。

2. **prosemirror-state**：管理编辑器的状态，包括文档的当前内容、选区（Selection）、以及通过插件（Plugin）扩展的状态。

3. **prosemirror-view**：负责将编辑器状态渲染到 DOM 中，并处理用户的输入事件。这个模块是编辑器与用户交互的桥梁。

4. **prosemirror-transform**：提供了一系列操作来变更文档内容，如添加、删除、移动节点等。这些操作通过事务（Transaction）来组织和应用。

5. **prosemirror-schema**：定义了文档的模式（Schema），包括可用的节点类型、它们的属性、以及节点之间的关系。模式是约束和指导文档结构的规则。

### 扩展模块

扩展模块提供了额外的功能，使得开发者可以根据需要为编辑器添加更多特性。

1. **prosemirror-commands**：提供了一组常用的编辑命令，如加粗文本、创建列表等。

2. **prosemirror-history**：添加了撤销（undo）和重做（redo）的功能。

3. **prosemirror-inputrules**：允许定义输入规则，自动转换输入的文本，如在输入`*`后自动创建一个列表。

4. **prosemirror-keymap**：用于定义键盘快捷键。

5. **prosemirror-dropcursor**：在拖放操作中显示光标位置的视觉反馈。

6. **prosemirror-gapcursor**：提供了一个特殊的光标，可以在块级节点之间移动。

7. **prosemirror-menu**：提供了构建编辑器菜单的工具。

8. **prosemirror-collab**：添加了协作编辑功能，允许多个用户同时编辑同一个文档。

这些扩展模块使得 ProseMirror 非常灵活和强大，开发者可以根据项目的需求选择性地引入这些模块来构建富文本编辑器。通过组合不同的模块，可以实现从简单的文本编辑器到复杂的协作编辑平台等各种应用。
