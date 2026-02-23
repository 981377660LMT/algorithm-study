# Claude Code 扩展能力

## Skill

Anthropic 对 ModelContextProtocol 的完全体的落地。
在提出 MCP 的时候, specifications 中三个主要组成部分：prompt/tool/resources​ 但是随着发展和社区影响，mcp 更像是 openapi/restfulapi 的代理层，大家都在关注工具（tools），prompts、resources 几乎从 mcp server 生态看不到，甚至大部分 mcp client 是没有支持的 prompt 和 resources 的，并不是 real model context。
​Skill 本身还是三大组成部分：Prompt/Tool/Resources，相比于 MCP 的高成本实现和接入，skill 以更轻量的方式（web 上下载个压缩包、coding 的时候拷贝一下）实现了 model context engineering.
mcp 的协议有些过于理想化，导致很多协议落地困难。而在 claude 本身的 sp 迭代中可以看出有很多类似 skill 的设计模式（遇到 x 触发 y），基于实践的经验，他们慢慢迭代出这一套更适合实际场景的简化协议

---

这种理解非常深刻，精准地捕捉到了 **“理想化协议 (MCP)”** 与 **“工程化实践 (Skills)”** 之间的演进张力。可以从以下四个维度来深入理解这一转变：

### 1. 从“协议优先”转向“上下文工程 (Context Engineering) 优先”

- **MCP 的困境**：MCP 试图建立一套通用的、跨语言的 C/S 架构协议。但在实际落地中，为了兼容性，它被迫变得“重”。开发者需要维护独立的 Server 进程，而模型与 Server 之间的通信（尤其是 Resources 和 Prompts）往往因为缺乏统一的渲染标准，最终退化成了简单的 API 调用（Tools）。
- **Skills 的突破**：Skill 不再纠结于通信协议，而是回归到 **“模型如何最高效地获取知识”**。它将代码、文档、指令统一为文件系统中的“静态资产”。通过 `SKILL.md`，它直接在模型的上下文窗口内进行“渐进式披露”，这本质上是更高阶的 Prompt Engineering。

### 2. 对 MCP 三大支柱的“降维打击”式实现

Skill 并没有抛弃 MCP 的 Prompt/Tool/Resource 理念，而是用更符合 Agent 直觉的方式重新实现了它们：

- **Prompt (入口)**：`SKILL.md` 本身就是一段高质量的 System Prompt 片段，它定义了能力的边界和触发时机。
- **Resources (知识)**：Skill 文件夹下的 `.md` 或 `.txt` 文件就是 Resource。模型通过 `Read` 工具按需读取，避免了 MCP 中 Resource 预加载导致的上下文污染。
- **Tools (执行)**：Skill 文件夹下的脚本（Python/JS）通过本地 Bash 直接执行。相比 MCP Server 的远程调用，这种方式让模型拥有了“即兴发挥”和“自我修正”的能力。

### 3. 渐进式加载：解决“上下文效率”的银弹

MCP 的一个痛点是：如果连接了 10 个 Server，模型往往难以在启动时就理清所有 Server 的所有能力。

- **Skill 的逻辑**： 1. **轻量索引**：启动时只给模型看 Skill 的名称和简短描述（Low Context Cost）。 2. **按需触发**：当模型判断需要该能力时，才调用 `Skill` 工具加载完整的 `SKILL.md`。 3. **深度探索**：根据 `SKILL.md` 的指引，进一步读取代码或文档。
  这种**三层渐进式披露**（Index -> Guide -> Implementation）是目前处理大规模 Agent 能力集的最优解。

### 4. 落地门槛：从“运维”转向“创作”

- **MCP**：需要写代码、打包 Docker、配置环境变量、维护运行状态。这更像是**后端开发**。
- **Skills**：只需要写 Markdown 和脚本，然后拷贝到 `.claude/skills` 目录。这更像是**内容创作**。
  这种低摩擦的特性，使得 Skill 能够像“插件”或“乐高积木”一样在社区快速流通。

### 总结

**MCP 是 Anthropic 给行业画的“同心圆”（试图标准化），而 Skill 是 Anthropic 给自己留的“后门”（追求极致体验）。**

Skill 的出现标志着 AI 开发者的重心正在从 **“如何连接工具”** 转向 **“如何为模型构建知识体系”**。它证明了在 Agent 领域，**“文件即协议”** 往往比 **“接口即协议”** 更加高效和灵活。

---

Skill 所有内容会组织在一个文件夹内，常见的目录结构，其中 SKILL.md 为 Skill 的入口，可以视为 main 函数。
Skill.md 的描述文件案例，由上方的一小段 meta 数据和下方的一整段 md 文档构成.
Skill 运行的本质是通过一个 Skill 工具进行渐进式加载：

1. Skill 工具的描述中附带所有具体 Skill 的基础信息（名称+描述），通常该信息非常简短。
2. 模型自动判断需要触发何种 skill，触发 skill 工具，入参为 skill 的名称（例如 document-skills:pdf)
3. tool_result 中返回剩余 SKILL.md 内容
4. 按照 SKILL.md 文件中的介绍，按需继续阅读 skill 文件夹下的文件（例如 scripts/xxx.py)
5. 模型决策执行代码文件 or 继续读取更多文件
6. 继续 Agent 主循环

## Claude.md

CLAUDE.md 包含项目特定的领域知识

## Command

本质上是结构化的可复用的提示词
目前 claude code 中 Command 可以被其他几乎所有能力主动调用（Subagents、其他 Command 等）

## Subagents

独立上下文的并行工作单元，唯一有并行能力的扩展功能

## Hooks

Hooks 是在代理生命周期特定点执行的 Python 函数，提供确定性处理和自动化能力。

## MCP

连接代理到外部工具和数据源，外部集成专用，上下文效率最低（启动时加载所有内容）
适用场景：

- 外部数据的引入：例如连接 JIRA、查询数据库、获取实时天气数据
- 本地工具：例如 Chrome 插件等（未来可能会被 skills 替代）

## Plugin

本质上是一个融合了以上所有功能的扩展性市场，用于打包和分发功能集合
由于以上所有能力都是基于“文件”+“配置”的方式进行的，因此 Plugin 本质就是将远程的能力下载到本地并且修改 setting.json 文件进行应用
并且大部分的远端都管理在 github 中，能很好复用 git 的版本控制和权限管理能力

# Claude Agent SDK 接口能力

Agent SDK 的执行方式为，将 node 的 npm 作为子进程，python 以 text stream 的方式用 stdin 和 sdtout 进行交互
