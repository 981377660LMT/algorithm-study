这份文档介绍了 **Deep Agents CLI**，这是一个用于构建和运行 Deep Agents 的交互式命令行界面。它允许你在终端中直接与具备持久记忆、文件操作和代码执行能力的 Agent 进行交互。

以下是 Deep Agents CLI 的核心功能和使用指南总结：

### 1. 核心能力
Deep Agents CLI 集成了多种工具，使其能够像开发者一样工作：
*   **文件操作**: 读、写、编辑项目文件。
*   **Shell 执行**: 运行测试、构建项目、管理依赖或 Git 操作。
*   **网络搜索**: 使用 Tavily 搜索最新文档和信息。
*   **任务规划**: 内置 Todo 系统，将复杂任务拆解为步骤并追踪进度。
*   **持久记忆**: 跨会话存储和检索信息（如项目规范）。
*   **人在回路**: 敏感操作（如修改代码）需要人工批准。

### 2. 快速开始

**步骤 1: 设置 API Key**
默认使用 Anthropic Claude Sonnet 4。
```bash
export ANTHROPIC_API_KEY="your-api-key"
# 可选：启用网络搜索
export TAVILY_API_KEY="your-key"
```

**步骤 2: 运行 CLI**
推荐使用 `uvx` 直接运行：
```bash
uvx deepagents-cli
```
或者本地安装后运行：
```bash
pip install deepagents-cli
deepagents-cli
```

**步骤 3: 下达任务**
```bash
> Create a Python script that prints "Hello, World!"
```

### 3. 记忆系统 (Memory)
这是该 CLI 的一大亮点。Agent 会将信息以 Markdown 文件的形式存储在 `~/.deepagents/AGENT_NAME/memories/` 目录下。

*   **学习能力**: 你只需教 Agent 一次项目规范（例如：“API 使用 snake_case”），它会自动保存到记忆中。
*   **跨会话应用**: 在未来的会话中，Agent 会自动检索并遵守这些规范，无需重复提示。

### 4. 远程沙箱 (Remote Sandboxes)
为了安全地执行代码，CLI 支持在隔离的远程环境中运行操作。

*   **支持提供商**: Runloop, Daytona, Modal。
*   **优势**: 保护本地机器安全、环境隔离、并行执行。
*   **使用方法**:
    ```bash
    # 配置提供商 Key 后
    uvx deepagents-cli --sandbox runloop --sandbox-setup ./setup.sh
    ```

### 5. 交互技巧
*   **Slash 命令**:
    *   `/tokens`: 查看 Token 使用情况。
    *   `/clear`: 清除对话历史。
    *   `/exit`: 退出。
*   **直接执行 Shell**: 在命令前加 `!` (例如 `!git status`)。
*   **快捷键**:
    *   `Ctrl+T`: 切换“自动批准”模式（跳过确认提示）。
    *   `Ctrl+E`: 打开外部编辑器。

### 6. 常用命令参数
*   `--agent NAME`: 使用特定名称的 Agent（拥有独立的记忆空间）。
*   `--auto-approve`: 自动批准所有工具调用。
*   `deepagents reset --agent NAME`: 重置指定 Agent 的记忆。