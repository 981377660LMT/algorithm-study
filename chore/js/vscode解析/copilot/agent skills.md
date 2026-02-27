针对 VS Code 1.109 引入的 **Agent Skills（代理技能）** 功能，以下是深度分析与实战指南：

### 1. 核心定义：从“指令”到“技能包”的进化

**Agent Skills** 不仅仅是提示词（Prompts），它是一个**功能完备的文件夹**。
与传统的“自定义指令（Custom Instructions）”相比，它的最大区别在于其**模块化**和**资源集成能力**。

- **自定义指令**：主要定义全局编码风格（如“始终使用 TypeScript”）。
- **Agent Skills**：定义特定任务的能力（如“如何根据公司规范编写 Webview 测试”），并且可以携带脚本、代码模板和示例文件。

### 2. 核心架构：三级渐进式加载 (Progressive Disclosure)

为了防止一次性给 AI 投喂太多信息导致上下文窗口爆炸，Agent Skills 采用了三层加载机制：

1.  **Level 1：技能发现 (Discovery)**：AI 仅仅读取 `SKILL.md` 的 `name` 和 `description`。这是轻量级的，AI 借此判断当前任务是否需要该技能。
2.  **Level 2：指令加载 (Instructions)**：一旦匹配，AI 会读取 `SKILL.md` 正文中的详细指令。
3.  **Level 3：资源访问 (Resources)**：只有当 AI 需要执行特定任务时，才会去读取同目录下的脚本（如 `test-template.js`）或文档。

### 3. 如何创建一个技能

一个标准的技能文件夹结构如下：

```text
.github/skills/webapp-testing/
├── SKILL.md            # 核心定义文件
├── test-template.js    # 辅助脚本
└── examples/           # 示例参考
```

#### `SKILL.md` 编写规范：

```markdown
---
name: webapp-testing
description: 当用户需要为 React 组件编写端到端测试时使用此技能。
argument-hint: [组件名称]
---

# Web 应用测试指南

1. 始终使用 Playwright 框架。
2. 引用 [测试模板](./test-template.js) 生成代码。
```

### 4. 关键配置项深度解析

- **`user-invokable` (默认 true)**：是否允许用户在 Chat 中输入 `/webapp-testing` 手动触发。
- **`disable-model-invocation` (默认 false)**：如果设为 `true`，AI 就不会根据需求自动加载它，必须由用户手动输入斜杠命令。这适用于那些比较“重”或者仅在特定时刻才需要的工具。
- **`chat.agentSkillsLocations`**：你可以在 VS Code 设置中指定自定义的技能存放位置，实现跨项目共享。

### 5. 开发者与插件作者的福音

Agent Skills 是一个**开放标准 (agentskills.io)**，这意味着：

- **可移植性**：你在 VS Code 写的技能，同样可以运行在 Copilot CLI 或其他支持该标准的 AI Agent 中。
- **插件分发**：扩展开发者可以通过 package.json 的 `chatSkills` 贡献点向用户提供“专业技能”。

```json
// package.json
"contributes": {
  "chatSkills": [
    { "path": "./skills/my-skill/SKILL.md" }
  ]
}
```

### 6. 实战建议

- **整理公司内部规范**：将“架构规范”、“部署流程”、“内部 API 文档”分别做成不同的 Skills。
- **配合子智能体 (Subagents)**：你可以让一个专门负责测试的子智能体去加载 `webapp-testing` 技能，而主智能体保持“干净”的上下文。
- **安全性**：在使用社区分享的技能时，务必检查其中的脚本。你可以通过 `chat.tools.edits.autoApprove` 设置来限制 AI 对敏感脚本的修改。

**总结**：Agent Skills 让 AI 从一个“通才”变成了可以按需加载专业能力的“专家组”。它把知识、工具和指令打包在一起，是构建复杂 AI 工作流的关键一环。
