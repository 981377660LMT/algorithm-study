## 应用结构

这份文档介绍了构建和部署 **LangGraph** 应用程序的标准结构和配置方式，特别是针对 **LangSmith Deployment** 的要求。

以下是核心要素总结：

### 1. 推荐目录结构
一个标准的 LangGraph 项目通常包含源代码目录、依赖文件和配置文件。

```plaintext
my-app/
├── src/
│   ├── agent.ts        # 定义和编译图的主要入口
│   └── utils/          # 工具函数、节点逻辑、状态定义等
├── package.json        # 项目依赖
├── .env                # 环境变量
└── langgraph.json      # LangGraph 专用配置文件
```

### 2. 核心配置文件 (`langgraph.json`)
这是部署 LangGraph 应用的关键文件，用于告诉运行时环境如何加载你的图。主要包含以下字段：

*   **`dependencies`**: 指定依赖项来源，通常设为 `["."]` 以使用根目录下的 package.json。
*   **`graphs`**: 定义应用中包含的图。
    *   格式为 `"图名称": "文件路径:导出变量名"`。
    *   例如：`"agent": "./src/agent.ts:graph"` 表示名为 `agent` 的图位于 `src/agent.ts` 文件中，且导出的变量名为 `graph`。
*   **`env`**: (可选) 定义默认的环境变量，通常用于本地开发。

### 3. 配置示例

**langgraph.json**
```json
{
  "dependencies": ["."],
  "graphs": {
    "my_agent": "./src/agent.ts:app",
    "research_assistant": "./src/research.ts:graph"
  },
  "env": {
    "OPENAI_API_KEY": "sk-..."
  }
}
```

### 4. 关键点
*   **Graphs 定义**: 必须明确指向编译后的图对象（CompiledGraph）或生成图的函数。
*   **依赖管理**: 确保 package.json 中列出了所有运行时需要的库。
*   **环境变量**: 生产环境部署时，环境变量通常在部署平台（如 LangSmith）的控制台中配置，而不是硬编码在文件中。

## 测试

这份文档介绍了如何使用 `vitest` 对 **LangGraph** Agent 进行单元测试。

以下是三种主要的测试模式总结：

### 1. 基础全流程测试 (End-to-End)
由于 Agent 依赖状态，建议在每个测试中重新构建图。
*   **模式**：创建一个 `createGraph` 工厂函数。
*   **执行**：在测试中实例化 `MemorySaver`，编译图，然后调用 `invoke` 检查最终状态。

### 2. 单节点测试 (Unit Testing Nodes)
你可以跳过整个图的流程，直接测试某个节点的逻辑。
*   **方法**：通过 `compiledGraph.nodes['nodeName']` 获取节点引用。
*   **执行**：直接调用 `.invoke(input)` 并断言输出。注意这种方式会绕过 Checkpointer。

### 3. 部分路径执行 (Partial Execution)
用于测试图中特定的一段流程（例如从 Node B 到 Node C），而不需要运行整个图。
*   **模拟起始状态**：使用 `updateState` 方法，并指定 `asNode` 参数（设置为你想开始测试的节点**之前**的那个节点名），以此来模拟“刚刚执行完前一个节点”的状态。
*   **控制结束点**：在 `invoke` 时配置 `interruptAfter`（或 `interruptBefore`），让图在执行完目标节点后暂停。

### 代码示例：部分路径执行
以下代码展示了如何只运行 `node2` 到 `node3` 的逻辑：

```typescript
test('partial execution from node2 to node3', async () => {
  const uncompiledGraph = createGraph();
  const checkpointer = new MemorySaver();
  const compiledGraph = uncompiledGraph.compile({ checkpointer });

  // 1. 模拟状态：假装刚刚执行完 node1，准备进入 node2
  await compiledGraph.updateState(
    { configurable: { thread_id: '1' } },
    { my_key: 'initial_value' }, // 模拟 node1 的输出
    'node1' // asNode: 标记这些状态是由 node1 产生的
  );

  // 2. 执行并截断：从 node2 开始，在 node3 之后停止
  const result = await compiledGraph.invoke(
    null, // 传入 null 表示从当前状态（即 node2）继续
    {
      configurable: { thread_id: '1' },
      interruptAfter: ['node3'] // 在 node3 执行完后强制暂停
    },
  );

  expect(result.my_key).toBe('hello from node3');
});
```

## LangSmith Studio 
这份文档介绍了 **LangSmith Studio**，这是一个用于在本地开发、调试和可视化 LangChain/LangGraph Agent 的免费可视化界面。

它允许你连接本地运行的 Agent，实时查看其执行步骤（Prompt、工具调用、结果），并进行交互测试。

以下是设置和使用 LangSmith Studio 的核心步骤：

### 1. 准备工作
*   **账号**: 注册 [LangSmith](https://smith.langchain.com) 账号。
*   **API Key**: 获取 LangSmith API Key。
*   **环境**: 需要 Python >= 3.11（用于安装 CLI）。

### 2. 安装与配置

#### 第一步：安装 LangGraph CLI
CLI 提供了一个本地开发服务器（Agent Server）来连接 Studio。
```bash
pip install --upgrade "langgraph-cli[inmem]"
```

#### 第二步：配置环境变量
在项目根目录创建 `.env` 文件，填入 API Key：
```bash
LANGSMITH_API_KEY=lsv2...
```

#### 第三步：创建配置文件 (`langgraph.json`)
此文件告诉 CLI 如何找到你的 Agent 代码。
```json
{
  "dependencies": ["."],
  "graphs": {
    "agent": "./src/agent.py:agent" 
  },
  "env": ".env"
}
```
*   `graphs`: 映射图的名称到代码路径（格式为 `文件路径:导出变量名`）。

### 3. 启动与使用

#### 启动开发服务器
在终端运行以下命令：
```bash
langgraph dev
```

#### 访问 Studio
*   服务器启动后，可以通过浏览器访问显示的 URL（通常是 `https://smith.langchain.com/studio/...`）。
*   **Safari 用户注意**: Safari 可能会阻止连接 `localhost`。请使用 `--tunnel` 参数启动：
    ```bash
    langgraph dev --tunnel
    ```

### 4. 主要功能
*   **可视化执行**: 查看 Agent 的完整 Trace，包括 Prompt、工具参数和返回值。
*   **调试**: 当发生错误时，可以查看异常及当时的完整状态。
*   **热重载 (Hot-reloading)**: 修改本地代码（如 Prompt 或工具逻辑）后，Studio 会立即反映更改，无需重启。
*   **交互测试**: 在 UI 中直接输入测试用例，无需编写额外的测试代码。
