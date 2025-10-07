正在收集工作区信息好的，我们来详细讲解一下 `tldraw` 的代码库。

`tldraw` 是一个用于在 React 中创建无限画布体验的软件库。它也是 [tldraw.com](https://tldraw.com) 背后的核心引擎。

### **仓库概览**

这是一个 monorepo（单一代码仓库），管理着多个相互关联的包和应用程序。

- **packages**: 存放构成 `tldraw` SDK 的核心库。
- **apps**: 存放使用 `tldraw` SDK 构建的应用程序，例如官方网站、文档和 VS Code 扩展。
- **templates**: 为 `create-tldraw` CLI 工具提供项目模板。
- **internal**: 包含用于开发和构建的内部脚本和配置。

### **核心架构**

`tldraw` 的核心架构分为三层，如 CONTEXT.md 中所述：

1.  **`@tldraw/editor`**: 纯粹的画布引擎，不包含任何具体的图形、工具或 UI。它提供了 `ShapeUtil`（用于定义图形行为）和 `StateNode`（用于实现工具状态机）等基础系统。
2.  **`@tldraw/store`**: 响应式数据层，负责状态管理、撤销/重做、持久化（使用 IndexedDB）和数据迁移。
3.  **`@tldraw/tldraw`**: 完整的 SDK，它在 `@tldraw/editor` 的基础上封装了全套 UI、默认图形（文本、箭头、几何图形等）和工具集。这是大多数开发者直接使用的包。

### **关键概念**

#### 1. Shape（图形）系统

每个图形的行为都由一个 `ShapeUtil` 类定义。例如，`FrameShapeUtil` (`packages/tldraw/api-report.api.md`) 负责处理框架（Frame）的渲染、交互和几何计算。这使得系统具有高度可扩展性，你可以通过创建自己的 `ShapeUtil` 来添加自定义图形。

#### 2. Tool（工具）系统

工具被实现为状态机（`StateNode`）。例如，选择工具 (`SelectTool`) 包含多个子状态，如 `brushing`（框选）和 `translating`（平移）。这种事件驱动的架构清晰地管理着用户的交互和工具的生命周期。

#### 3. 文件格式 (`.tldr`)

`tldraw` 使用自定义的 JSON 文件格式，后缀为 `.tldr`。

- **解析**: `parseTldrawJsonFile` 函数用于解析 `.tldr` 文件的 JSON 字符串。
- **序列化**: `serializeTldrawJson` 函数用于将编辑器当前的内容序列化为 JSON 字符串。

### **如何使用 tldraw**

在你的 React 项目中集成 `tldraw` 非常简单。

1.  **安装依赖**:

    ```sh
    npm i tldraw
    ```

2.  **在组件中使用**:
    在你的应用中，导入 `Tldraw` 组件和它的 CSS 文件。

    ```tsx
    // filepath: apps/examples/src/examples/basic/BasicExample.tsx
    import { Tldraw } from 'tldraw'
    import 'tldraw/tldraw.css'

    export default function BasicExample() {
      return (
        <div className="tldraw__editor">
          <Tldraw />
        </div>
      )
    }
    ```

### **主要应用 (apps)**

- `apps/examples/`: 包含了超过 130 个示例，是开发和测试新功能的主要环境。
- `apps/dotcom/client/`: [tldraw.com](https://tldraw.com) 网站的前端应用，实现了文件管理、用户认证和多人协作等功能。
- `apps/vscode/`: VS Code 扩展，允许用户直接在编辑器中查看和编辑 `.tldr` 文件。
- `apps/docs/`: 官方文档网站 ([tldraw.dev](https://tldraw.dev))。

### **本地开发**

如果你想为 `tldraw` 贡献代码或在本地进行开发，可以按照 README.md 中的步骤操作：

1.  **克隆仓库并安装依赖**:

    ```sh
    git clone https://github.com/tldraw/tldraw.git
    cd tldraw
    npm i -g corepack
    yarn
    ```

2.  **启动开发服务器**:
    ```sh
    yarn dev
    ```
    这会启动 examples 应用，你可以在 `http://localhost:5420` 查看，对代码库的任何修改都会实时反映出来。

### **许可**

`tldraw` SDK 在 tldraw license 下提供。你可以在商业或非商业项目中使用它，但默认情况下画布上会带有 "Made with tldraw" 的水印。可以通过购买商业许可证来移除水印。
