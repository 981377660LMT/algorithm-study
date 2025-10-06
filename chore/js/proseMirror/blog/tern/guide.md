https://ternjs.net/doc/manual.html
好的，这是一份对您提供的 Tern 官方参考手册的详细讲解。我会将这份技术性很强的文档分解成更易于理解的部分，并对关键概念进行阐述。

这份手册详尽地描述了 Tern 的所有组成部分，从高层的编辑器插件到底层的类型推断引擎。

---

### 1. Tern 服务器 (The Tern Server)

这部分描述了 Tern 的核心——后台服务器进程。

- **启动方式**: 通常，你不需要手动启动服务器。编辑器插件（如 VS Code、Vim 的插件）会在你打开一个 JS 项目时自动为你启动 `bin/tern` 这个 Node.js 脚本。
- **配置**:
  - 服务器启动后，会向上查找 `.tern-project` 文件来获取项目配置。
  - 如果没有找到，它会使用你主目录下的 `.tern-config` 文件作为全局默认配置。
- **通信**: 服务器会监听一个随机端口，并通过 HTTP 提供一个简单的 JSON API。编辑器插件就是通过向这个端口发送 JSON 数据来与服务器通信的。
- **命令行参数**:
  - `--port`: 指定一个固定端口，而不是随机端口。
  - `--host`: 指定监听的主机地址（默认为本地 `127.0.0.1`）。
  - `--persistent`: 默认情况下，服务器在闲置 5 分钟后会自动关闭。此参数可以禁用该行为。
  - `--ignore-stdin`: 默认情况下，当启动它的进程（通常是编辑器）关闭时，服务器也会关闭。此参数可以禁用该行为。
  - `--verbose`: 输出详细的日志，用于调试。

---

### 2. JSON 协议 (JSON Protocol)

这是编辑器插件与 Tern 服务器之间通信的“语言规范”。所有请求都是通过向服务器端口发送一个 HTTP POST 请求，请求体为 JSON 文档。

一个请求主要包含三个部分：`query`（你要问什么）、`files`（你提供的代码文件）和 `timeout`（超时限制）。

#### 核心查询类型 (`query.type`)

- **`completions` (自动补全)**:

  - **作用**: 在代码的特定位置请求补全建议。
  - **关键参数**: `file` 和 `end` 指定位置；`docs`, `urls`, `origins` 可以请求返回更丰富的信息（文档、链接、定义来源）；`filter` 控制是在服务器端还是客户端过滤结果。
  - **返回**: 一个包含 `start`, `end` 和 `completions` 数组的对象。

- **`type` (查询类型)**:

  - **作用**: 查询光标处表达式的类型，通常用于鼠标悬停提示。
  - **返回**: 包含 `type` (类型描述字符串), `name` (变量/属性名), `doc` (文档) 等信息的对象。`guess` 字段表示这个类型是推断出来的还是猜测的。

- **`definition` (跳转到定义)**:

  - **作用**: 查找变量或属性的定义位置。
  - **返回**: 包含 `file`, `start`, `end` (定义位置) 和 `context` (定义处的上下文代码) 的对象。

- **`refs` (查找引用)**:

  - **作用**: 查找一个变量或属性在整个项目中的所有引用。
  - **返回**: 包含 `name` (变量名) 和一个 `refs` 数组，数组中每个对象都指明了引用的文件和位置。

- **`rename` (重命名)**:

  - **作用**: 对一个变量进行作用域安全的重命名。
  - **返回**: 一个 `changes` 数组，描述了需要对哪些文件的哪些位置进行文本替换。**Tern 只提供修改方案，实际的文件修改由编辑器插件完成。**

- **`files` (获取文件列表)**:
  - **作用**: 获取服务器当前分析的所有文件的列表。

---

### 3. 编程接口 (Programming Interface)

这部分是为那些想在自己的程序中（而不是通过 HTTP）直接使用 Tern 服务器功能的开发者准备的。

- **`new Server(options)`**: 创建一个服务器实例。你可以通过 `options` 对象传入配置，如 `defs` (类型定义), `plugins` (插件), `getFile` (一个函数，告诉 Tern 如何读取文件内容)。
- **核心方法**:
  - `addFile(name, text)`: 向服务器添加一个文件。
  - `delFile(name)`: 从服务器移除一个文件。
  - `request(doc, callback)`: 发送一个编程版的 JSON 请求。
  - `on(event, handler)` / `off(event, handler)`: 监听服务器事件（如 `beforeLoad`, `afterLoad`），这是插件工作的关键。

---

### 4. JSON 类型定义 (JSON Type Definitions)

这是 Tern 的一个核心特性。因为 Tern 无法凭空知道浏览器环境的 `document` 对象或 Node.js 的 `fs` 模块是什么，所以需要一种方式来“告诉”它。JSON 类型定义就是这个方式。

- **格式**: 一种特殊的 JSON 格式，用来描述对象、函数、属性及其类型。
- **特殊属性 (以 `!` 开头)**:
  - `!name`: 定义这个库的名称，作为类型的来源标识。
  - `!type`: 指定一个对象的具体类型，如 `fn(...)` 表示函数。
  - `!proto`: 指定对象的原型。
  - `!doc` / `!url`: 为类型附加文档字符串和链接。
  - `!define`: 定义可以在文件内部复用的局部类型。
- **类型字符串**:
  - `"number"`, `"string"`, `"bool"`: 原始类型。
  - `"fn(arg1: type) -> rettype"`: 函数类型。
  - `"[type]"`: 数组类型。
  - `"+MyConstructor"`: 表示 `MyConstructor` 的一个实例。

---

### 5. 服务器插件 (Server Plugins)

插件用于扩展 Tern 的核心功能，使其能够理解特定的环境、库或模块系统。

- **工作原理**: 插件通过在服务器初始化时注册事件处理器（hooks）来工作。例如，监听 `preParse` 事件可以在解析前修改代码（如从 HTML 中提取 JS），监听 `postInfer` 事件可以在类型推断后进行补充。
- **重要插件**:
  - **`doc_comment`**: 解析 JSDoc 风格的注释 (`/** ... */`) 来获取类型和文档信息。
  - **`commonjs`**: 让 Tern 理解 `require`, `module`, `exports`。
  - **`node`**: 加载 Node.js 的内置模块类型（如 `fs`, `http`）和全局变量（如 `process`），并启用 `commonjs` 插件。
  - **`es_modules`**: 支持 ES6 的 `import` 和 `export`。
  - **`requirejs`**: 支持 RequireJS (AMD) 模块系统。
  - **`webpack`**: 利用 webpack 的 `enhance-resolve` 模块来理解 `webpack.config.js` 中的别名和解析规则。
  - **`angular`**: 尝试理解 Angular.js 的依赖注入。
- **第三方插件**: 你可以自己编写插件并通过 npm 分发。包名通常遵循 `tern-pluginname` 的格式。

---

### 6. 项目配置 (`.tern-project`)

这是用户配置 Tern 如何分析项目的入口文件，是一个 JSON 文件。

- `libs`: 加载哪些内置的 JSON 类型定义。例如 `["browser", "jquery"]`。
- `loadEagerly`: 一个文件或 glob 模式数组。Tern 会在启动时就主动加载和分析这些文件，以加快后续响应。
- `dontLoad`: 告诉 Tern 忽略某些文件或目录，通常用于排除 node_modules。
- `plugins`: 配置要加载的插件及其选项。
- `ecmaVersion`: 指定解析时使用的 JS 版本（如 5, 6, 8...）。
- `dependencyBudget`: 一个数字，限制 Tern 为加载一个依赖项（如通过 `require`）所愿意分析的代码量，防止因分析一个巨大的依赖而卡死。

---

### 7. 推断引擎 (Inference Engine)

这部分深入到了 Tern 的最底层，即 `lib/infer.js`。这是所有类型推断魔法发生的地方。

- **核心概念**:
  - **Context (上下文)**: 保存了全局作用域和所有类型信息的状态容器。
  - **Analysis (分析)**: `infer.analyze(ast)` 是核心函数，它遍历由 **Acorn** 解析器生成的 AST（抽象语法树），并构建类型关系。
  - **Types (类型)**: `Obj`, `Fn`, `Arr` 等是内部表示 JavaScript 类型的对象。
  - **Abstract Values (抽象值, AVal)**: 这是关键。一个 AVal 代表一个“类型集合”。由于 JS 的动态性，一个变量可能在不同地方有不同类型，AVal 就用来存储所有这些可能性。
  - **Constraints (约束)**: 类型传播的机制。当分析 `let x = y` 时，Tern 会建立一个约束，使得 `y` 的 AVal 中所有类型都会“传播”到 `x` 的 AVal 中。整个推断过程就是建立和求解一个庞大的约束网络。
- **工具函数**:
  - `infer.findExpressionAt(ast, ...)`: 在 AST 中根据位置查找对应的表达式节点，这是所有查询的第一步。
  - `infer.expressionType(...)`: 推断一个表达式节点的类型。

---

### 8. 编辑器插件 (Editor Plugins)

这部分列出了将 Tern 集成到各种编辑器的插件。它们是 Tern 的“前端”，负责将用户在编辑器中的操作转换成对 Tern 服务器的 JSON 请求，并展示返回的结果。

总而言之，这份手册是一份非常全面的开发者文档，它不仅告诉了终端用户如何配置 Tern (`.tern-project`)，也为希望扩展 Tern（编写插件）或理解其工作原理（推断引擎）的开发者提供了深入的指南。
