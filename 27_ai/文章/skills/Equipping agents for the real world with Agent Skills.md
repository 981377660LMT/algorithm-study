# Equipping agents for the real world with Agent Skills

https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills

Anthropic 推出的 **Agent Skills（智能体技能）** 是一种通过文件和文件夹为 AI 智能体（如 Claude）赋予特定领域专业知识和程序性知识的新标准。

### 1. 核心理念：将“隐性知识”显性化与模块化

目前的通用大模型（General-purpose Agents）虽然强大，但缺乏特定组织或任务的上下文（Context）和流程知识（Procedural Knowledge）。
**Agent Skills 的本质是将人类的“操作手册”标准化为 AI 可读的格式。** 它不再需要为每个任务构建碎片化的专用智能体，而是像给员工发“入职指南”一样，通过挂载不同的文件夹，让通用智能体瞬间获得特定领域的专业能力。

### 2. 技术架构：渐进式披露 (Progressive Disclosure)

这是 Agent Skills 最精妙的设计原则，旨在解决**上下文窗口（Context Window）限制**与**海量知识库**之间的矛盾。它通过三个层级动态加载信息：

- **Level 1：元数据（Metadata）**

  - 每个 Skill 是一个包含 `SKILL.md` 的文件夹。
  - `SKILL.md` 顶部的 **YAML Frontmatter**（包含名称和描述）会在启动时被预加载到 System Prompt 中。
  - **作用**：让 Claude 知道“我有这个技能”，但暂不占用大量 Token。

- **Level 2：核心指令（Core Instructions）**

  - 当 Claude 判断当前任务需要该技能时，才会读取 `SKILL.md` 的正文内容。
  - **作用**：加载该技能的高层指导原则和流程。

- **Level 3：按需加载的资源（Linked Resources）**
  - 对于复杂的技能，`SKILL.md` 会引用子文件（如 `forms.md`、scripts）。
  - Claude 只有在执行到具体步骤（如“填写表单”）时，才会去读取对应的子文件。
  - **作用**：实现上下文的无限扩展能力，理论上技能包的大小可以不受限制。

### 3. 混合计算范式：LLM + 确定性代码

文章强调了一个关键的工程洞察：**并非所有任务都适合用 Token 生成来解决。**

- **LLM 的弱点**：生成式任务（如排序、复杂的数学计算、大规模文本提取）昂贵且存在概率性错误。
- **Skills 的解决方案**：允许在`技能包中封装可执行脚本（如 Python）。`
- **工作流**：Claude 不再尝试“生成”结果，而是作为**编排者（Orchestrator）**，调用技能包中预写好的 Python 脚本来处理数据（例如解析 PDF 表单字段）。这结合了 LLM 的灵活性和传统代码的确定性与效率。

### 4. 开发与安全最佳实践

- **开发流程**：建议采用“评估驱动开发”。先运行任务发现 Claude 的短板，再针对性地编写 Skill 文档。
- **迭代机制**：可以让 Claude 自我反思（Self-reflection），将它成功的操作路径固化为 Skill 代码或文档。
- **安全边界**：由于 Skills 包含可执行代码，必须像对待第三方软件库一样，仅安装受信任来源的 Skills，并审计其中的脚本和网络请求指令。

### 总结

Agent Skills 是 AI Agent 从“玩具”走向“生产力工具”的重要一步。它通过**文件系统即接口（Filesystem as Interface）** 的方式，标准化了人类知识注入 AI 的过程，并利用**渐进式加载**完美平衡了模型智力与计算成本。这预示着未来 AI 开发将从单纯的 Prompt Engineering 转向构建模块化、可复用的“技能库”。
