# VS Code 博客归档

https://code.visualstudio.com/blogs/archive

## 2026

### 一月（January）

- [为 Agent 赋予视觉表达：VS Code 中的 MCP 应用支持 / Giving Agents a Visual Voice: MCP Apps Support in VS Code](https://code.visualstudio.com/blogs/2026/01/26/mcp-apps-support)
- [构建 docfind：使用 Rust 和 WebAssembly 的高速客户端搜索 / Building docfind: Fast Client-Side Search with Rust and WebAssembly](https://code.visualstudio.com/blogs/2026/01/15/docfind)

## 给 AI 代理赋予可视化表达：VS Code 中的 MCP Apps 支持

https://code.visualstudio.com/blogs/2026/01/26/mcp-apps-support

**精简讲解：MCP Apps 是"代理的可视化对话界面"**

- **问题**：VS Code 的 AI 代理虽然能执行代码编辑/命令/数据库查询，但**只能用文本回应**——需要重复文字来完成"列表重排""图表探索""确认破坏性操作"等任务；来回沟通效率低。
- **方案**：MCP Apps（Model Context Protocol 的官方扩展），让 MCP 服务器的工具调用可以返回**交互式 UI 组件**（仪表板/表单/可视化），直接在 agent 聊天面板里渲染。
- **三个核心用例**：
  - **列表拖放排序**：代理提出排序建议 → 用户直接在 UI 里拖拽调整，而不是来回讲"把 A 移到 B 上面"。
  - **性能分析火焰图**：代理分析 CPU profile → 渲染交互式火焰图，用户可以深钻调用栈、hover 看时间戳，自己验证瓶颈。
  - **功能开关选择器**：代理列出特性开关配置 → 显示可搜索的 picker，带 prod/staging/dev 环境标签，用户选完即可生成 SDK 代码，一步完成。
- **实现方式**：MCP 工具调用的返回值里除了文本，还可以包含 UI 组件描述；这些组件在 agent 聊天里直接渲染，用户交互反馈回代理做多轮对话。
- **生态推进**：Storybook（开源 MCP 服务）已支持 MCP Apps，现在可以在 agent 聊天里直接预览设计系统组件，无需切换到 Storybook 网页。
- **可用性**：2026/1/26 在 VS Code Insiders（日构建）已上线；预计下周稳定版也会有。
- **设计思想**：VS Code 一直是"文本编辑器+"（扩展能加 UI）；Jupyter notebook 展示了"代码+富输出"如何改变工作流；现在 GitHub Copilot agents 有了 MCP Apps，就是"代理可以说出来、看出来、交互"，不仅是告诉，还能示意。

## 构建 docfind：使用 Rust 和 WebAssembly 实现快速客户端搜索

https://code.visualstudio.com/blogs/2026/01/15/docfind

**精简逻辑版讲解：docfind 是怎么做“秒出结果”的站内搜索**

- **问题**：VS Code 文档站原来搜索要跳转传统搜索引擎；目标是像 VS Code `Ctrl+P` 一样“边输边出结果”，且不想引入服务端/运维成本。
- **调研结论**：Algolia/TypeSense 要服务端或运维；Lunr.js 在 ~3MB 文档量下索引约 10MB 太大；Stork 索引仍偏大且维护状态不理想 → 决定自研纯前端方案。
- **核心思路**：把“文档→关键词→相关文档”的检索完全搬到浏览器里，用 **Rust + WebAssembly** 做高性能执行。
- **三件套技术**：
  - **FST（有限状态转导器）**：用紧凑自动机存大量排序关键词，实现快速匹配；
  - **RAKE**：从每篇文档自动提取关键词并打相关性分；
  - **FSST**：把标题/分类/片段等短字符串压缩存储，显示时再按需解压。
- **索引怎么构建**：构建期 CLI 读 `documents.json` → RAKE 提取关键词 → 建 FST（关键词→索引）+ `keyword -> [(doc, score)]` 映射 → 文档字符串 FSST 压缩 → 打包成二进制索引。
- **关键工程技巧**：索引不是单独文件，而是**直接嵌进 WASM 模块**；为了避免“文档一改就重编 WASM”，用“预编译 WASM 模板 + CLI patch 二进制”的方式，把索引作为 data segment 写入，并把 `INDEX_BASE/INDEX_LEN` 两个占位全局变量从 `0xdead_beef` 改成真实地址/长度。
- **查询怎么跑**：用户输入时加载 WASM（代码+索引），用 **前缀匹配 + Levenshtein automaton（容错拼写）** 命中关键词，合并多关键词得分排序，只对命中结果解压必要字符串。
- **效果（文中数据）**：VS Code 网站（~3MB markdown、~3,700 文档分段）索引约 **5.9MB**，Brotli 后约 **2.7MB**；查询约 **0.4ms/次**（M2 MacBook Air）；**只在用户有搜索意图时下载一个 WASM 资源**，无需服务端、无需 API key。

## 2025

### 十二月（December）

- [推介 VS Code Insiders 播客 / Introducing the VS Code Insiders Podcast](https://code.visualstudio.com/blogs/2025/12/03/introducing-vs-code-insiders-podcast)

## 推介 VS Code Insiders 播客

**核心概念：内容载体升级，从文本博客到音频对话**

- **解决的问题**：开发者很难了解 VS Code 背后的设计决策、优先级选择、以及与社区的互动。博客是单向文本传播；播客可以呈现更丰富的观点碰撞和人物故事。
- **播客内容**：由 VS Code 官方发起，邀请编辑、产品经理、社区贡献者等做深度对话，涵盖从可访问性设计、Planning Agent 实现原理、开源社区建设，到 GitHub Universe 关键发布等议题。
- **价值定位**：
  - **内幕视角**：听到一线工程师讨论技术权衡、为什么某些特性被优先级降低、如何在约束下创新。
  - **多元观点**：不同角色（设计师、工程师、产品）的思考碰撞，比单一作者博客更有深度。
  - **AI 趋势讨论**：直接讨论 AI 对编程的影响、代理模式的未来等前沿话题。
- **传播策略**：免费、多平台（Spotify、Apple Podcasts 等），降低收听门槛；将博客文章升级成音频对话，形成"视频+音频+文本"的多维传播矩阵。
- **目标受众**：不只是开发者，还包括产品思考者、开源贡献者、AI 爱好者等。
- **启示**：媒体多元化是长期社区建设的必要条件；一个成熟的开发者平台需要同步维护博客、播客、视频、文档、GitHub 讨论等多个传播渠道。

### 十一月（November）

- [推介 Visual Studio Code 私有市场：团队的安全策展扩展中心 / Introducing the Visual Studio Code Private Marketplace: Your Team's Secure, Curated Extension Hub](https://code.visualstudio.com/blogs/2025/11/18/PrivateMarketplace)
- [开源 AI 编辑器：第二个里程碑 / Open Source AI Editor: Second Milestone](https://code.visualstudio.com/blogs/2025/11/04/openSourceAIEditorSecondMilestone)
- [统一的编码 Agent 体验 / A Unified Experience for all Coding Agents](https://code.visualstudio.com/blogs/2025/11/03/unified-agent-experience)

## Visual Studio Code 私有市场

**核心概念：企业级扩展管理平台，让团队掌握自己的工具生态**

- **痛点**：企业开发团队面临的困境是"公开市场太开放"——无法控制开发者安装什么扩展、无法内部审计、无法离线部署、对安全敏感的代码无法使用第三方扩展。
- **解决方案**：私有市场（Private Marketplace）是专为 GitHub Enterprise 客户的一个"受管制的扩展应用商店"。企业 IT/安全团队可以：
  - **策展**：精选哪些扩展对团队可见（内部工具 + 审核过的公开扩展）。
  - **本地化**：把公开扩展重新下载、扫描、再托管到自己的私有市场，适配离线/气隙环境。
  - **集中分发**：新员工入职后自动看到团队批准的扩展列表，消除"手动安装 zip 问题"。
- **使用体验**：从用户视角无感——扩展搜索、安装、更新的流程保持一致，但数据来源和版本管理都由企业控制。
- **安全考量**：
  - 专有扩展（含 IP）可以托管在防火墙内，不泄露代码。
  - 对公开扩展进行安全审计后再分发，降低供应链风险。
  - 符合 SOC2、HIPAA 等合规要求。
- **目标用户**：GitHub Enterprise、Copilot Business/Enterprise 用户，对安全治理要求高的大型企业。
- **市场视角**：这是"开发工具企业化"的关键一步，VS Code 从"个人开发工具"升级到"企业级工作平台"。

## 开源 AI 编辑器：第二个里程碑

**核心概念：从"开源 Chat 扩展"到"开源内联建议"，AI 功能逐步透明化**

- **背景**：2025 年 5 月，VS Code 宣布开源 Copilot Chat 扩展（agent mode）；2025 年 6 月达到第一个里程碑（Chat 开源）；2025 年 11 月达到第二个里程碑——**内联建议（inline suggestions）也开源了**。
- **内联建议的重要性**：
  - 这是"边输代码边显示 AI 补全"的核心能力，用户最高频接触的 AI 功能。
  - 之前被视为"黑盒"（内部实现细节），现在完全开源透明。
- **扩展整合目标**：
  - 把两个扩展合并成一个——GitHub Copilot Chat 扩展现在包括 inline suggestions（之前由 GitHub Copilot 扩展提供）。
  - **好处**：用户体验一致、工程维护集中、社区贡献单一入口。
  - **风险管理**：旧的 Copilot 扩展将在 2026 年早期从市场下架；提供临时回退选项以防兼容性问题。
- **性能优化故事**：开源化过程中，团队发现了优化机会：
  - **"按建议输入"检测**：如果用户的输入匹配上一条建议，直接继续显示，不需新请求。
  - **缓存复用**：多个相似输入可以共享缓存建议。
  - **请求复用**：若上一次 LLM 请求还在流式返回中，直接复用而非取消再新建。
  - **多行智能**：自动判断是显示单行还是多行补全，基于置信度和上下文。
- **开源社区贡献**：代码在 [vscode-copilot-chat](https://github.com/microsoft/vscode-copilot-chat) 仓库公开，贡献指南、PR 审查流程都透明，邀请开源社区改进。
- **后续计划**：进一步重构 AI 功能组件从扩展移入 VS Code 核心；持续监控迭代计划并邀请反馈。
- **深层含义**：AI 功能开源不仅是代码透明，还是"信任建立"——开发者可以审计 Copilot 怎样运作、数据如何处理、隐私如何保障。

## 统一的编码 Agent 体验

**核心概念：多代理编排，把 VS Code 升级成"AI 团队管理中心"**

- **问题描述**：2025 年代理爆炸——GitHub Copilot、Copilot Coding Agent、GitHub Copilot CLI、OpenAI Codex，加上用户可能的自定义代理。结果是"选择瘫痪"和"上下文混乱"。
- **四大解决方案**：

  **1. Agent Sessions（代理会话面板）**
  - 在 VS Code 侧栏新增一个"Agent Sessions"视图，统一管理所有正在运行的代理。
  - 用户可以看到哪些代理在运行、运行状态、快速切换会话。
  - **聊天编辑器（Chat Editors）**：对每个代理的聊天可在标签页中打开，用户可"在线修正"代理的计划——不需等待或取消，直接发送更新给代理调整策略。
  - **任务委托（Delegate）**：在任何聊天中可以一键把任务转移给其他代理，保留上下文。

  **2. OpenAI Codex 集成**
  - VS Code 现在支持用 GitHub Copilot Pro+ 订阅直接使用 OpenAI Codex。
  - 无需单独的 OpenAI 账号；Copilot 统一管理模型调用和速率限制。
  - 用户可以在多个代理之间快速切换，尝试不同的代理风格。

  **3. Planning Agent（规划代理）**
  - 一个新的内置代理，专门把"懒惰提示"转化成详细计划。
  - 例如用户只说"加拖拽排序"，Planning Agent 会问：
    - 加到什么组件？
    - 用什么库？（推荐对比 React Beautiful DnD vs React DnD）
    - 还有什么考虑？
  - **Handoff（交接）特性**：Planning Agent 整理好计划后，用户可以选择"在编辑器打开完整计划"或"继续执行"，灵活决策。
  - **模型选择建议**：文中透露 Claude 模型在识别缺失上下文、提出正确问题上表现更好。

  **4. Subagents（子代理）**
  - 解决"上下文膨胀"问题的新工具 `#runSubagent`。
  - 子代理运行在隔离的上下文中，只知道你明确传入的信息，不会被主聊天的历史污染。
  - **使用场景**：主聊天在开发 API，需要研究认证方案 → 启动子代理专门研究 → 研究完毕只有最终结论返回主聊天。
  - **性能优势**：子代理不会因为主聊天的冗长历史而变慢；可以并行运行多个子代理，加快探索。

- **工程思想**：
  - 从"一个代理统治一切"转向"多代理生态 + 统一编排"。
  - 从"单线程对话"升级到"多线程+隔离上下文"的并发执行。
  - 定位 VS Code 为"代理编排的控制中心"而非仅仅是聊天界面。

- **一年的进度**：2024 年 11 月发布 Copilot Edits；2025 年 11 月已有 4 大代理 + Planning/Subagents。作者戏言"12 个月后会是什么样"，预示 Agent 生态将继续加速。

### 十月（October）

- [使用自带密钥扩展 VS Code 中的模型选择 / Expanding Model Choice in VS Code with Bring Your Own Key](https://code.visualstudio.com/blogs/2025/10/22/bring-your-own-key)

### 九月（September）

- [推介自动模型选择（预览）/ Introducing auto model selection (preview)](https://code.visualstudio.com/blogs/2025/09/15/autoModelSelection)

## 自动模型选择（Auto Model Selection）

**核心概念：智能降级 + 成本优化，让用户无需关心模型细节**

- **问题背景**：
  - 2025 年 LLM 模型选择爆炸——Claude Sonnet 4、GPT-5、GPT-5 mini、新的开源模型。
  - 用户面临"选模型焦虑"：选大模型贵但聪明，选小模型快但可能失败。
  - 系统面临"容量管理"：VS Code 日均百万级 agent 请求，需要智能分流。

- **解决方案**：自动模型选择（Auto）
  - **自动调度**：基于当前系统容量、模型负载、用户历史行为，VS Code 自动选择最佳模型（Claude Sonnet 4 → GPT-5 → GPT-5 mini）。
  - **用户感知**：从用户角度，只需选 "Auto"，无需关心底层用了哪个模型。可以悬停查看本次调用用的具体模型。
  - **成本激励**：
    - 免费用户：享受速率限制内的最佳模型。
    - 付费用户：Auto 选择时请求成本打 9 折（Claude Sonnet 4 计为 0.9x）；若选小模型（GPT-5 mini）则不扣费（0x）。
  - **容量韧性**：若用户用完了付费请求配额，Auto 会自动降到 0x 模型（免费可用），确保体验连贯。

- **工程角度的意义**：
  - **负载均衡**：避免单一大模型过载；通过多模型库存平衡流量。
  - **成本可控**：企业客户可以预测成本（10% 折扣更好估算）。
  - **用户友好**：无需用户学习 Claude vs GPT-5 的优缺点，系统智能应对。

- **未来方向**（文中透露）：
  - **任务自适应**：根据任务复杂度动态选模型（简单 completion 用 mini，复杂 reasoning 用 Sonnet）。
  - **扩展模型库**：集成更多模型提供商。
  - **免费用户支持**：让免费层也能通过 Auto 享受新模型。
  - **UI 改进**：让用户更直观地看到成本和模型映射。

- **深层启示**：这是从"用户选择模型"到"系统智能调度"的范式转变，类似于 SQL 查询优化器或操作系统调度器的逻辑——用户表达意图，系统优化执行。

### 八月（August）

- [VS Code Dev Days – 参加附近活动了解 AI 辅助开发 / VS Code Dev Days – Join an event near you to learn about AI-assisted development](https://code.visualstudio.com/blogs/2025/08/27/vscode-dev-days)

### 七月（July）

- [在 VS Code 中启用 GitHub 编码 Agent / Command GitHub's Coding Agent from VS Code](https://code.visualstudio.com/blogs/2025/07/17/copilot-coding-agent)

### 六月（June）

- [开源 AI 编辑器：第一个里程碑 / Open Source AI Editor: First Milestone](https://code.visualstudio.com/blogs/2025/06/30/openSourceAIEditorFirstMilestone)
- [完整的 MCP 体验：VS Code 中的全规范支持 / The Complete MCP Experience: Full Specification Support in VS Code](https://code.visualstudio.com/blogs/2025/06/12/full-mcp-spec-support)

## 在 VS Code 中启用 GitHub 编码 Agent

**核心概念：把云端 Agent 的能力拉进编辑器，实现"任务委托"工作流**

- **GitHub Copilot Coding Agent 是什么**：
  - 一个在 GitHub 云端运行的自主 AI 开发者，可以被分配 GitHub issue。
  - Agent 会：探索代码库、修改文件、编译、运行测试、打开 PR、审查反馈、迭代。
  - 运行在临时隔离的开发环境中，可以访问数据库、云服务（通过 MCP）等外部工具。
  - 工作流：你在 GitHub 上把 issue 分配给 Copilot → Agent 工作并打开 PR → 你审查反馈 → Agent 迭代 → 合并。

- **为什么这篇文章关键**：Coding Agent 原本只能在 GitHub Web 上操作；现在可以**直接从 VS Code 中命令 Agent**，大幅降低上下文切换。

- **VS Code 中的四大集成**：

  **1. Pull Requests 视图新增"Copilot on My Behalf"**
  - 专门的查询列表，显示所有 Copilot 代理正在处理的 issue 和 PR。
  - 点击"View Session"可以看到 Agent 的实时进度——每一步命令执行、每一条决策理由，完全透明。
  - 可以随时终止 Agent（若方向错误）。

  **2. PR 审查和迭代**
  - Agent 完成时自动分配 PR 给你做审查。
  - PR 包含截图（如 UI 变更），无需本地检出即可快速验证。
  - 如果预览服务支持（Netlify/Vercel/Azure Static Web Apps），还能看到实时部署预览。
  - 你在 VS Code 里留评论 → Agent 接收反馈 → 自动更新 PR。

  **3. MCP 服务配置**
  - 若某些任务（如数据库迁移）需要 Agent 完全自主，可在 GitHub 仓库设置中配置 MCP 服务。
  - 例如：为 Agent 注册 Supabase MCP 服务器，Agent 就可以直接修改数据库结构，而非只能生成迁移脚本。
  - 这让 Agent 的自主性从"代码生成"升到"端到端任务完成"。

  **4. 从 Chat 直接委托任务**
  - 在 Copilot Chat 对话中，任何时刻都可以点"Delegate to Coding Agent"。
  - 聊天的全部上下文会传给 Agent，Agent 据此打开 PR 并开始工作。
  - 跳过了 GitHub issue 这一步，直接从 VS Code 的思路流进入实现流。

  **5. 从 Sidebar 快速分配**
  - GitHub Pull Requests 扩展的侧栏里，可以直接选择 issue 并分配给 Copilot Coding Agent。

- **成本模型**：每个 Agent session 计为**一次高级请求**，包含代码探索、修改、测试的全部工作，非常划算。

- **VS Code 团队自己用**：他们在 VS Code 仓库中已有大量 Copilot 打开的 PR（可 GitHub 搜索 `head:copilot/`），验证了生产就绪度。

- **后续投资方向**：
  - PR 性能和渲染优化。
  - 集成的 Chat 视图来跟踪 Agent 会话。
  - Agents 命令中心（类似 Agent Sessions）。
  - 自定义指令在 Coding Agent 和 VS Code 间共享。
  - 更多文档。

- **工程意义**：这不是"用机器人替代开发者"，而是"开发者和机器人的协作编程"——你负责高层决策和审查，Agent 负责机械性实现，效率显著提升。

## 开源 AI 编辑器：第一个里程碑

**核心概念：VS Code AI 功能全面开源，build in public，社区参与**

- **战略背景**：2025 年 5 月，VS Code 团队宣布把 VS Code 打造成"开源 AI 编辑器"。目标是：
  - 社区驱动创新：vs Code 在编辑器领域成功，AI 功能也应该开源。
  - 数据透明：让开发者可以审计 AI 怎么运作、数据去向。
  - 整体开放度：与 VS Code 核心一致。

- **第一个里程碑（2025/6/30）**：GitHub Copilot Chat 扩展开源到 MIT 许可证，在 GitHub 公开。
  - 开发者可以：阅读代码理解 agent mode 实现、看到系统 prompt、了解数据处理流程。
  - 可以贡献代码、提 Issue、参与设计讨论。
  - 长期目标：把这些代码重构后集成到 VS Code 核心。

- **社区参与方式**：
  - 代码开源在 [vscode-copilot-chat](https://github.com/microsoft/vscode-copilot-chat)。
  - Issue 跟踪在主 [vscode](https://github.com/microsoft/vscode) 仓库。
  - 有完整的贡献指南（CONTRIBUTING.md）。
  - 可以提 PR，参与迭代。

- **下一步**：
  - 继续开源更多 AI 组件（内联建议 ← 这在第二个里程碑完成了）。
  - 逐步把扩展功能重构为 VS Code 核心的一部分。
  - 保持原有的 GitHub Copilot 扩展（提供 inline completions），但计划在 2026 早期迁移到单一 Chat 扩展。
  - 邀请开源社区参与，确保计划涵盖真实开源场景。

- **哲学思考**：VS Code 成功的秘诀是"开放 + 社区"；AI 时代应该延续这个传统，而不是把 AI 当黑盒商业产品。透明化 AI 运作方式，才能建立长期信任。

## 完整的 MCP 体验：VS Code 中的全规范支持

**核心概念：Model Context Protocol 从"工具集"升级成"完整代理基础设施"**

- **MCP 背景**：Model Context Protocol 是 Anthropic 2024 年发起的开放标准，定义了"LLM 和外部系统的通信协议"。VS Code 之前只支持部分功能（工具、基础工作区感知）；现在支持完整规范。

- **五大核心 MCP 原语（Primitives）**：

  **1. Authorization（授权）— 最重要的新增**
  - 新的授权规范由 Microsoft、Anthropic、Okta/Auth0、Stytch、Descope 等共同设计。
  - **核心创新**：把"MCP 服务器"和"授权服务器"分离。服务器无需自己实现 OAuth，而是委托给专门的身份提供商。
  - **VS Code 应用**：GitHub MCP 服务器现在是远程服务，利用 VS Code 既有的 GitHub 认证；用户登一次，所有 MCP 服务都能访问。
  - **企业价值**：远程 MCP 服务可以独立扩展，同时保持企业级安全。

  **2. Prompts（提示）**
  - 不是静态模板，而是"动态工作流"。MCP 服务器可以根据当前工作区状态、项目上下文等，动态生成提示。
  - 在 VS Code 的斜杠命令（slash commands）中显示为 `/mcp.servername.promptname`。
  - **例子**：Gistpad MCP 提供的 prompt 可能会根据你打开的文件推荐相关的 snippet。

  **3. Resources（资源）**
  - 代表"有语义的信息对象"，可以直接在 VS Code 中交互。
  - **例子**：
    - Playwright MCP 截屏 → 图片成为资源 → 你可以拖进工作区、注解、分享。
    - 调试工具返回日志 → 日志流式更新到 VS Code → 不用切换浏览器。
  - 这让外部工具的输出无缝集成到工作区。

  **4. Sampling（采样）— 最受欢迎的新功能**
  - 允许 MCP 服务器自己调用 LLM（之前只有客户端能调 LLM）。
  - **好处**：服务器无需管理 AI SDK 和 API 密钥；直接用 VS Code 订阅的模型。
  - **用途**：复杂推理、多代理协调。
  - **控制**：VS Code 提供模型选择器，用户可以指定 MCP 服务器能用哪些模型（比如限制只用 Claude，不能用 GPT）。

  **5. Tools（工具）**
  - 原有支持，现在只是规范的一部分。

- **生态应用**：
  - Storybook（开源组件库工具）已支持 MCP Apps（MCP 的 UI 扩展）。现在在 Agent 聊天里直接预览设计系统组件，而非切换到 Storybook 网页。
  - GitHub MCP 服务器：远程服务，支持授权、prompt、resource 等原语。

- **学习曲线**：
  - 新手：还是只用 tools 就够。
  - 进阶：定义 prompts 来封装复杂工作流。
  - 高级：用 sampling 让服务器自主推理，building multi-agent systems。

- **战略意义**：MCP 从"API 规范"升到"代理基础设施"。结合 VS Code 的 Agent Sessions、Subagents，这套体系可以支撑复杂的多代理编程工作流，是 AI-native 开发环境的基石。

### 五月（May）

- [通过 AI + 远程开发提升生产力 / Enhance productivity with AI + Remote Dev](https://code.visualstudio.com/blogs/2025/05/27/ai-and-remote)
- [VS Code：开源 AI 编辑器 / VS Code: Open Source AI Editor](https://code.visualstudio.com/blogs/2025/05/19/openSourceAIEditor)
- [超越工具，在 VS Code 中集成 MCP / Beyond the tools, adding MCP in VS Code](https://code.visualstudio.com/blogs/2025/05/12/agent-mode-meets-mcp)

### 四月（April）

- [Agent 模式：对所有用户可用且支持 MCP / Agent mode: available to all users and supports MCP](https://code.visualstudio.com/blogs/2025/04/07/agentMode)

### 三月（March）

- [上下文至关重要：使用自定义指令获得更优 AI 结果 / Context is all you need: Better AI results with custom instructions](https://code.visualstudio.com/blogs/2025/03/26/custom-instructions)

### 二月（February）

- [推介 GitHub Copilot Agent 模式（预览）/ Introducing GitHub Copilot agent mode (preview)](https://code.visualstudio.com/blogs/2025/02/24/introducing-copilot-agent-mode)
- [Copilot 下一步编辑建议（预览）/ Copilot Next Edit Suggestions (preview)](https://code.visualstudio.com/blogs/2025/02/12/next-edit-suggestions)

## 2024

### 十二月（December）

- [宣布 VS Code 免费 GitHub Copilot / Announcing a free GitHub Copilot for VS Code](https://code.visualstudio.com/blogs/2024/12/18/free-github-copilot)

### 十一月（November）

- [推介 GitHub Copilot for Azure（预览）/ Introducing GitHub Copilot for Azure (preview)](https://code.visualstudio.com/blogs/2024/11/15/introducing-github-copilot-for-azure)
- [推介 Copilot 编辑（预览）/ Introducing Copilot Edits (preview)](https://code.visualstudio.com/blogs/2024/11/12/introducing-copilot-edits)

### 六月（June）

- [GitHub Copilot 扩展就是一切 / GitHub Copilot Extensions are all you need](https://code.visualstudio.com/blogs/2024/06/24/extensions-are-all-you-need)

## 宣布 VS Code 免费 GitHub Copilot

**核心概念：AI 开发工具的"消费者化"，降低准入门槛，拉动用户规模**

- **商业策略转变**：
  - 2024 年 12 月，GitHub Copilot 推出**免费层** — 所需仅 GitHub 账号，无需订阅、无需试用期、无需信用卡。
  - 免费配额：每月 2000 次代码补全（约 80/工作日）+ 50 次聊天请求 + GPT-4o + Claude 3.5 Sonnet 访问。
  - 付费升级：Pro 计划无限制，新增 o1 和 Gemini 模型。

- **为什么开放免费**：
  - **扩大基数**：AI 已是编码基础设施，vs Code 需要锁定开发者生态。
  - **体验优先**：让开发者先尝到甜头，再考虑付费；这是 SaaS 标准玩法。
  - **竞争压力**：GitHub Copilot 面对 Claude、GPT-4o 等竞品，免费是防守策略。

- **新功能综合展示**（文中强调 2024 年新增）：

  **1. Copilot Edits（多文件编辑体验）**
  - 不同于"补全"和"对话"，Copilot Edits 是"多文件协作编辑"的新范式。
  - 用户指定"Working Set"（一组要编辑的文件）+ 自然语言指令 → Copilot 提议跨文件变更。
  - **核心特色**：在编辑器里直接显示修改（inline），用户可逐条 Accept/Discard。
  - **迭代流**：支持 Undo/Redo，用户可反复调整，直到满意。可边运行单元测试边审查改动。
  - **技术栈**：双模型架构 — 基础模型（GPT-4o/o1/Claude Sonnet）生成改动提议 + 推测解码端点快速应用。
  - **应用例**：一位无 Swift 经验的产品经理用 Copilot Edits 从零写了 macOS 应用；VS Code 工程师用它做大范围重构。

  **2. 模型选择自由**
  - 用户可在 Chat/Inline Chat/Copilot Edits 间自由选择模型。
  - 例如：用 GPT-4o 规划，再用 Claude Sonnet 编码（充分利用两者优势）。

  **3. 自定义指令（Custom Instructions）**
  - 在编辑器或项目级指定行为偏好（代码风格、框架选择等）。
  - 自动从 `.github/copilot-instructions.md` 加载，方便团队协作。
  - **例子**：指定"React 18 + hooks + TypeScript，对象属性使用简写"等规范。

  **4. 全项目上下文感知（@workspace 参与者）**
  - `@workspace` 是一个"代码库领域专家"，可用 `@` 语法调用。
  - Copilot 自动检测意图，若问题需要全项目上下文会自动引入 @workspace。
  - 用户可输入 `/help` 看所有可用参与者。

  **5. 智能重命名**
  - 按 `F2` 进行重命名时，Copilot 根据符号实现和使用场景给出建议，解决"命名难"问题。

  **6. 语音聊天**
  - 点击麦克风图标启动语音输入（基于本地运行的开源语音模型，无需第三方 app）。
  - 配合 Copilot Edits 做快速原型，"说话"就能写出 demo。

  **7. 终端智能**
  - 在终端按 Cmd/Ctrl + I，Copilot 帮你写 shell 命令、解释错误、自动修复。
  - 例如：无需记 ffmpeg 语法，直接问"怎么从视频提取帧"。

  **8. 提交消息生成**
  - 无需再写"changes"，Copilot 根据改动 diff 和历史提交自动建议提交消息。
  - 支持自定义指令定制格式。

  **9. 扩展生态优化**
  - 每个扩展都可利用 Copilot API 提供 AI 体验（已有 100+ 扩展）。
  - MongoDB、Stripe 等工具扩展提供 `@mongodb` 等参与者，无需切换工具。

  **10. Vision（视觉预览）**
  - 可根据截图/Figma 设计生成 UI 代码。
  - 可生成图片的 alt text。
  - 目前需自带 OpenAI/Anthropic/Gemini API key；后续免费集成。

- **定位**：VS Code 从"编辑器"升级到"AI-native 编辑器"，Copilot 成为核心功能而非插件。
- **2025 展望**：作者暗示"Copilot 无处不在"，还有更多跨产品集成（VS、GitHub.com 等）。

## GitHub Copilot for Azure（预览）

**核心概念：Copilot 的"云平台领域专家"，把 Azure 文档和操作拉进编辑器**

- **问题场景**：
  - 开发者在 VS Code 写代码，突然需要部署到 Azure 或了解 Azure 服务 → 要切换到 Azure Portal。
  - 需要查文档、找命令、配置基础设施，多个标签页来回切。

- **解决方案**：GitHub Copilot for Azure 扩展
  - 一个 `@azure` 参与者，直接在 VS Code Chat 里与 Azure 交互。
  - 集成 Azure 账户、文档、资源管理 API。

- **四大应用场景**：

  **1. 学习 Azure**
  - 提问"Azure AI Search 是什么？"、"什么 Azure 服务能跑我的容器？"
  - Copilot 从最新文档提取、解释，无需搜索多个源。
  - 对新手降低学习曲线；对老手快速查阅新服务。

  **2. 部署应用**
  - 推荐 azd（Azure Developer CLI）模板、自动化部署配置。
  - 帮你写 YAML、选择资源、配置 CI/CD。
  - 例如：问"怎么用 Python 建 RAG chat app"，Copilot 推荐模板、步骤、依赖。
  - 无需手工查 CLI 参数或 sample repo。

  **3. 故障诊断**
  - 分析日志、识别性能瓶颈。
  - 例如："为什么我的 Kubernetes 集群变慢？"→ Copilot 查诊断数据、推荐调整（副本数、资源限制等）。
  - "网站返回 500 错误"→ 查日志、定位代码问题、建议修复。

  **4. 运营管理**
  - 查询资源列表、成本分析。
  - 例如："我有多少个免费层 App Service Plan"、"存储账户按大小按地区分类统计"。
  - 无需登 Portal，直接数据驱动决策。

- **斜杠命令（Slash Commands）**：
  - `/help` — 查看功能列表
  - `/learn` — 学习 Azure
  - `/resources` — 资源查询
  - `/diagnose` — 故障诊断
  - `/changeTenant` — 切换租户

- **状态**：2024 年 11 月预览上线，可从 VS Marketplace 安装。
- **深意**：这是"AI + 云平台"融合的范式——开发者不需要分裂注意力在多个 UI 上，一个聊天窗口就是"云控制台"。

## Copilot Edits（预览）

**核心概念：从"对话"和"补全"进化到"多文件协作编辑"的新人机交互模式**

- **问题描述**：
  - 之前 Copilot 要么用补全（单行/几行），要么用 Chat（问答）。
  - 需要一次性修改多个文件时（如大重构），要么手工逐文件改，要么 Chat 里长篇大论描述 。
  - 缺少"我告诉 Copilot 方向，它并行修改多个文件，我逐条审查"的工作流。

- **Copilot Edits 的设计**：
  - 打开 Copilot Edits 面板（Secondary Side Bar，默认右边）。
  - 定义"Working Set"：拖拽文件或按 `#` 添加要编辑的文件集合。
  - 输入自然语言指令。
  - Copilot 在编辑器里直接显示每个文件的改动（inline diff 风格）。
  - **逐条决策**：每个改动旁有 Accept/Discard 按钮，用户可混合接受/拒绝。
  - **迭代**：不满意可 Undo、调整指令、继续迭代，直到完美。
  - **验证**：支持边跑单测边审查（左边测试面板，右边编辑面板）。

- **技术架构**：
  - **双模型系统**：
    - 基础模型（GPT-4o/o1-preview/o1-mini/Claude 3.5 Sonnet 任选）生成编辑建议。
    - 推测解码端点（speculative decoding）快速在编辑器里渲染改动。
  - **推测解码优化**：比常规模型快，但还在改进中。

- **核心特色**：
  - **精确控制**：Working Set 确保改动只在指定文件，除非创建新文件。
  - **语音友好**：支持语音输入，可以边说边迭代，像与同事 pair program 的体验。
  - **无缝过渡**：可从 Chat 切到 Edits，保留上下文（后续特性）。
  - **建议工作集**：后续会自动推荐哪些文件应加入 Working Set。

- **应用例**：
  - 产品经理（无编码经验）用 Copilot Edits 快速迭代产品想法的早期代码。
  - VS Code 团队用它做大规模重构（跨百多个文件）。
  - 有人零基础 Swift 经验，用 Copilot Edits 从零搭建 macOS 应用，每次迭代都运行验证。

- **当前状态**：预览版，所有 GitHub Copilot 用户可用；后续计划改进推测解码性能、支持块级 Undo。
- **哲学**：从"AI 作为问答工具"升到"AI 作为协作编辑器"，更贴近真实开发流程（提意见、看方案、逐步调整）。

## GitHub Copilot Extensions are all you need

**核心概念：开放 Copilot 能力给第三方扩展，让开发工具生态具有 AI 原生特性**

- **战略背景**：
  - 2024 年 Build 大会，Microsoft 公布两个新 API：Chat API 和 Language Model API。
  - 这两个 API 把 GitHub Copilot 的能力开放给任何 VS Code 扩展。
  - 标题隐喻"扩展就是一切"（向 Google "Attention is All You Need" 论文致敬）。

- **两个核心 API**：

  **1. Chat API**
  - 让扩展在 Chat 视图里注册"聊天参与者"（Chat Participant）。
  - 例如 `@stripe` 参与者，就是 Stripe 官方扩展提供的。
  - 用户 `@stripe 怎么生成 Stripe 支付链接` → 参与者处理请求、返回 Stripe 文档+示例代码。
  - 使用 Language Model API 来处理自然语言并生成响应。

  **2. Language Model API**
  - 扩展可直接调用 Copilot 提供的 LLM（无需自己管理 API 密钥）。
  - 功能：
    - 直接从 LLM 获取响应。
    - 融入 VS Code 上下文（当前文件内容、工作区信息等）。
    - 支持多种应用场景：编辑器上下文菜单、源代码管理（生成 commit message）、Copilot 驱动的重命名等。

- **两种扩展方式**：

  **方式 1：VS Code 扩展（本地）**
  - 直接在 VS Code 扩展里用 Chat + Language Model API。
  - 例子：MongoDB for VS Code 扩展提供 `@mongodb` 参与者，能生成复杂 MongoDB 查询。
  - Stripe 扩展提供 `@stripe` 参与者，教你 Stripe API 用法、生成集成代码。
  - PostgreSQL Chat Participant 扩展让你 `@pg` 和数据库对话。

  **方式 2：GitHub App（远程）**
  - 基于后端服务的 GitHub App，跨 GitHub.com、VS Code、Visual Studio 等多平台。
  - 需加入"Copilot Partner Program"；不能完全访问 VS Code API（权限受限）。

- **生态案例**：
  - **Stripe**：`@stripe` 参与者，从文档生成支付集成代码，减少查阅时间。
  - **MongoDB**：`@mongodb` 参与者，根据自然语言生成 MongoDB 查询、分析查询性能、给出文档架构建议。
  - **Parallels**：`@parallels` 参与者，用自然语言操作虚拟机（"启动 Windows 11 VM"）。
  - **PostgreSQL**：`@pg` 参与者，学习 PostgreSQL、生成 SQL、生成数据访问代码。

- **已有 100+ 扩展基于这些 API 建立**，涵盖数据库、支付、基础设施、框架等。

- **后续计划**（2024 年中的预告）：
  - 意图检测（Intent Detection）：自动唤起合适的参与者，无需用户手工 `@xxx`。
  - GPT-4o 支持。
  - 提高 token 限制（当时是 4K）。
  - Chat 参与者支持编辑器 inline chat、终端、笔记本。
  - Variables Resolving API：扩展贡献变量，提供领域特定上下文。
  - Tools API：自然语言转工具调用，多参与者协作。

- **核心哲学**：
  - 不是"Copilot 垄断 AI"，而是"Copilot 是基础设施，扩展生态赋能于此"。
  - 类比：vs Code 通过 Extension API 统治编辑器 → Copilot 通过 Chat/LM API 统治 AI 编程。
  - 最强大的 AI 体验来自于"深层领域知识"（Stripe、MongoDB）+ "通用推理能力"（LLM）的结合。
- [使用 WebAssembly 进行扩展开发 - 第二部分 / Using WebAssembly for Extension Development - Part Two](https://code.visualstudio.com/blogs/2024/06/07/wasm-part2)

### 五月（May）

- [使用 WebAssembly 进行扩展开发 / Using WebAssembly for Extension Development](https://code.visualstudio.com/blogs/2024/05/08/wasm)

### 四月（April）

- [VS Code Day 2024 终极指南 / Your Ultimate Guide to VS Code Day 2024](https://code.visualstudio.com/blogs/2024/04/15/vscode-day)

## 2023

### 十一月（November）

- [追求 VS Code 中的卓越智能 / Pursuit of wicked smartness in VS Code](https://code.visualstudio.com/blogs/2023/11/13/vscode-copilot-smarter)

### 七月（July）

- [通过名称改编缩小 VS Code / Shrinking VS Code with name mangling](https://code.visualstudio.com/blogs/2023/07/20/mangling-vscode)

## 通过名称混淆缩小 VS Code

https://code.visualstudio.com/blogs/2023/07/20/mangling-vscode

**一句话版**

- VS Code 团队发现“代码已压缩但变量/属性名仍很长”，于是用“构建时自动改短名字（name mangling）”在几乎不改源码的前提下，把核心 JS 体积砍了约 20%，并带来启动更快（少扫描源码文本）。

**核心概念：什么是 mangling**

- 在不改变程序语义的前提下，把长标识符改成短标识符（如 `someLongVariableName` → `x`），从而减少 JavaScript 源代码文本长度。
- 这对 Web 场景尤其有价值：下载更小、磁盘更小、启动时 JS 引擎要“扫描/解析”的文本更少。

**为什么“默认压缩”还不够**

- 工具（文中是 esbuild）默认只会在“非常确定安全”的情况下 mangling。
- 它通常会改“局部变量/参数名”，但会避开很多“属性名/导出名”，因为 JS 里这些名字经常会被“字符串/动态访问”依赖，一改就可能炸，比如：
  - 动态属性访问：`obj[prop]`
  - 序列化/反序列化：JSON 期待固定字段名
  - 对外 API：调用方不知道你改了名
  - DOM/外部库 API：名字不能乱动

**他们踩过的坑：只靠规则匹配太危险**

- 试过“把以下划线开头的属性当私有属性去改短”（esbuild 支持按正则 mangling 属性名）。
- 问题是：`_` 只是约定，不等于真私有；而且代码里可能有“外部访问 private（TypeScript 的 private 运行时并不真的私有）”、测试里可能需要访问等，靠命名规则很难 100% 保证正确。

**关键做法：用 TypeScript 做“安全网”，在 TS 源码层改名**
他们的突破点是：不要在编译后的 JS 上瞎改，而是：

- 先用 TypeScript AST 找出“private/protected 属性”（语义层面更可靠）。
- 对这些符号做 TypeScript 的“重命名 refactor”（能自动更新所有引用）。
- 把生成的重命名编辑应用回 TS 源码，再正常编译打包。
- 这样 TypeScript 编译本身就会帮你抓住“改名导致的引用遗漏/冲突”，比只跑运行时测试更有把握。

**需要额外处理的边界情况**

- 新名字不仅要在当前类里唯一，还要在继承链（父类/子类）里避免冲突（因为 TS 的 `private` 只是编译期约束，运行时还是同名属性）。
- 有些地方子类把继承来的 `protected` “变成 public”（不管是历史遗留还是设计），这些场景需要禁用 mangling 或特殊处理。

**结果（文章里的数字）**

- 仅 mangling 私有属性：`workbench.js` 约从 12.3MB → 10.6MB（接近 14%）。
- 进一步 mangling “仅内部使用的导出符号名”（避免碰扩展 API、避免碰未类型化 JS 调用的入口等）：`workbench.js` 10.6MB → 9.8MB。
- 总体：`workbench.js` 相比不做 mangling 约小 20%；全 VS Code 产物少了约 3.9MB JS。
- 还带来约 5% 的加载提速（原因很朴素：少扫描/解析源码文本）。

**这篇文章想传达的“工程方法论”**

- 不是“为了极致压缩而牺牲可维护性”，而是：
  - 先承认体积增长是长期趋势，并持续监控；
  - 找到“收益大、风险可控”的切入口（私有属性、内部导出）；
  - 用强工具链（TypeScript 的重命名 + 再编译）把风险降到可接受；
  - 让优化“几乎对开发者透明”，才值得长期维护。

### 六月（June）

- [在 VS Code for Web 中运行 WebAssembly / Run WebAssemblies in VS Code for the Web](https://code.visualstudio.com/blogs/2023/06/05/vscode-wasm-wasi)

### 四月（April）

- [VS Code Day：为编辑器举办的活动？/ VS Code Day: An event for an editor?](https://code.visualstudio.com/blogs/2023/04/13/vscode-day)

### 三月（March）

- [Visual Studio Code 与 GitHub Copilot / Visual Studio Code and GitHub Copilot](https://code.visualstudio.com/blogs/2023/03/30/vscode-copilot)

## 2022

### 十二月（December）

- [远程开发更好了 / Remote Development Even Better](https://code.visualstudio.com/blogs/2022/12/07/remote-even-better)

### 十一月（November）

- [VS Code 迁移到进程沙箱 / Migrating VS Code to Process Sandboxing](https://code.visualstudio.com/blogs/2022/11/28/vscode-sandbox)

### 十月（October）

- [VS Code 社区讨论：为扩展开发者 / VS Code Community Discussions for Extension Authors](https://code.visualstudio.com/blogs/2022/10/04/vscode-community-discussions)

### 九月（September）

- [自定义开发容器特性 / Custom Dev Container Features](https://code.visualstudio.com/blogs/2022/09/15/dev-container-features)

### 八月（August）

- [推介 Markdown 语言服务器 / Introducing the Markdown Language Server](https://code.visualstudio.com/blogs/2022/08/16/markdown-language-server)

### 七月（July）

- [Visual Studio Code 服务器 / The Visual Studio Code Server](https://code.visualstudio.com/blogs/2022/07/07/vscode-server)

### 五月（May）

- [开发容器 CLI / The dev container CLI](https://code.visualstudio.com/blogs/2022/05/18/dev-container-cli)

### 四月（April）

- [使用容器从本地开发迁移到远程开发 / Using Containers to move from Local to Remote Development](https://code.visualstudio.com/blogs/2022/04/04/increase-productivity-with-containers)

### 三月（March）

- [教程的问题 / The problem with tutorials](https://code.visualstudio.com/blogs/2022/03/08/the-tutorial-problem)

## 2021

### 十一月（November）

- [Notebook，VS Code 风格 / Notebooks, Visual Studio Code style](https://code.visualstudio.com/blogs/2021/11/08/custom-notebooks)

### 十月（October）

- [vscode.dev(!) / vscode.dev(!)](https://code.visualstudio.com/blogs/2021/10/20/vscode-dev)
- [Visual Studio Code Webview UI 工具包 / Webview UI Toolkit for Visual Studio Code](https://code.visualstudio.com/blogs/2021/10/11/webview-ui-toolkit)

### 九月（September）

- [括号对着色快 10,000 倍 / Bracket pair colorization 10,000x faster](https://code.visualstudio.com/blogs/2021/09/29/bracket-pair-colorization)

### 八月（August）

- [Notebook 的时代来临 / The Coming of Age of Notebooks](https://code.visualstudio.com/blogs/2021/08/05/notebooks)

### 七月（July）

- [工作区信任 / Workspace Trust](https://code.visualstudio.com/blogs/2021/07/06/workspace-trust)

### 六月（June）

- [远程仓库 / Remote Repositories](https://code.visualstudio.com/blogs/2021/06/10/remote-repositories)
- [Visual Studio Code 在 Build 2021 / Visual Studio Code at Build 2021](https://code.visualstudio.com/blogs/2021/06/02/build-2021)

### 二月（February）

- [使用平分法解决扩展问题 / Resolving extension issues with bisect](https://code.visualstudio.com/blogs/2021/02/16/extension-bisect)

## 2020

### 十二月（December）

- [在 Chromebook 上使用 VS Code 学习 / Learning with VS Code on Chromebooks](https://code.visualstudio.com/blogs/2020/12/03/chromebook-get-started)

### 七月（July）

- [教育中的开发容器：讲师指南 / Development Containers in Education: A Guide for Instructors](https://code.visualstudio.com/blogs/2020/07/27/containers-edu)
- [在 WSL 2 中使用开发容器 / Using Dev Containers in WSL 2](https://code.visualstudio.com/blogs/2020/07/01/containers-wsl)

### 六月（June）

- [Go 体验的下一阶段 / The next phase of the Go experience](https://code.visualstudio.com/blogs/2020/06/09/go-extension)

### 五月（May）

- [Visual Studio Code 在 Build 2020 / Visual Studio Code at Build 2020](https://code.visualstudio.com/blogs/2020/05/14/vscode-build-2020)
- [推介 GitHub Issues 集成 / Introducing GitHub Issues Integration](https://code.visualstudio.com/blogs/2020/05/06/github-issues-integration)

### 三月（March）

- [在 WSL 2 中使用 Docker / Using Docker in WSL 2](https://code.visualstudio.com/blogs/2020/03/02/docker-in-wsl2)

### 二月（February）

- [自定义数据格式：发展 HTML 和 CSS 语言特性 / Custom Data Format: Evolving HTML and CSS language features](https://code.visualstudio.com/blogs/2020/02/24/custom-data-format)
- [优化 CI 构建时间 / Improving CI Build Times](https://code.visualstudio.com/blogs/2020/02/18/optimizing-ci)

## 2019

### 十月（October）

- [使用 VS Code 检查容器 / Inspecting Containers with VS Code](https://code.visualstudio.com/blogs/2019/10/31/inspecting-containers)
- [Remote SSH：技巧和窍门 / Remote SSH: Tips and Tricks](https://code.visualstudio.com/blogs/2019/10/03/remote-ssh-tips-and-tricks)

### 九月（September）

- [VS Code 与 WSL 2 / WSL 2 with Visual Studio Code](https://code.visualstudio.com/blogs/2019/09/03/wsl2)

### 七月（July）

- [Remote SSH 与 Visual Studio Code / Remote SSH with Visual Studio Code](https://code.visualstudio.com/blogs/2019/07/25/remote-ssh)

### 五月（May）

- [严格空值检查 Visual Studio Code / Strict null checking Visual Studio Code](https://code.visualstudio.com/blogs/2019/05/23/strict-null)
- [VS Code 远程开发 / Remote Development with VS Code](https://code.visualstudio.com/blogs/2019/05/02/remote-development)

### 二月（February）

- [语言服务器索引格式（LSIF）/ The Language Server Index Format (LSIF)](https://code.visualstudio.com/blogs/2019/02/19/lsif)

## 2018

### 十二月（December）

- [首次查看丰富的代码导航体验 / First look at a rich code navigation experience](https://code.visualstudio.com/blogs/2018/12/04/rich-navigation)

### 十一月（November）

- [Event-Stream 包安全更新 / Event-Stream Package Security Update](https://code.visualstudio.com/blogs/2018/11/26/event-stream)

### 九月（September）

- [Visual Studio Code 使用 Azure Pipelines / Visual Studio Code using Azure Pipelines](https://code.visualstudio.com/blogs/2018/09/12/engineering-with-azure-pipelines)
- [在 Visual Studio Code 中的 GitHub Pull Requests / GitHub Pull Requests in Visual Studio Code](https://code.visualstudio.com/blogs/2018/09/10/introducing-github-pullrequests)

### 八月（August）

- [Debug Adapter Protocol 新主页 / New home for the Debug Adapter Protocol](https://code.visualstudio.com/blogs/2018/08/07/debug-adapter-protocol-website)

### 七月（July）

- [推介日志点和自动附加 / Introducing Logpoints and auto-attach](https://code.visualstudio.com/blogs/2018/07/12/introducing-logpoints-and-auto-attach)

### 五月（May）

- [Visual Studio Live Share 公开预览 / Visual Studio Live Share Public Preview](https://code.visualstudio.com/blogs/2018/05/07/live-share-public-preview)

### 四月（April）

- [VS Code 中的必应驱动设置搜索 / Bing-powered settings search in VS Code](https://code.visualstudio.com/blogs/2018/04/25/bing-settings-search)

### 三月（March）

- [文本缓冲区重新实现 / Text Buffer Reimplementation](https://code.visualstudio.com/blogs/2018/03/23/text-buffer-reimplementation)

## 2017

### 十二月（December）

- [Chrome 调试的新增功能 / What's new for Chrome debugging](https://code.visualstudio.com/blogs/2017/12/20/chrome-debugging)

### 十一月（November）

- [Visual Studio Code 在 Connect(); 2017 / Visual Studio Code at Connect(); 2017](https://code.visualstudio.com/blogs/2017/11/16/connect)
- [推介 Visual Studio Live Share / Introducing Visual Studio Live Share](https://code.visualstudio.com/blogs/2017/11/15/live-share)

### 十月（October）

- [图标之旅 / The Icon Journey](https://code.visualstudio.com/blogs/2017/10/24/theicon)
- [集成终端性能改进 / Integrated Terminal Performance Improvements](https://code.visualstudio.com/blogs/2017/10/03/terminal-renderer)

### 九月（September）

- [使用 VS Code 调试 Java 应用 / Using VS Code to Debug Java Applications](https://code.visualstudio.com/blogs/2017/09/28/java-debug)

### 八月（August）

- [Visual Studio Code 中的 Emmet 2.0 / Emmet 2.0 in Visual Studio Code](https://code.visualstudio.com/blogs/2017/08/07/emmet)

### 六月（June）

- [新鲜油漆 - 给 VS Code 一个新外观 / Fresh Paint - Give VS Code a New Look](https://code.visualstudio.com/blogs/2017/06/20/great-looking-editor-roundup)

### 五月（May）

- [Build 2017 演示 / Build 2017 Demo](https://code.visualstudio.com/blogs/2017/05/10/build-2017-demo)

### 四月（April）

- [Sublime Text 扩展综述 / Sublime Text Extension Roundup](https://code.visualstudio.com/blogs/2017/04/10/sublime-text-roundup)

### 三月（March）

- [扩展包 / Extension Packs](https://code.visualstudio.com/blogs/2017/03/07/extension-pack-roundup)

### 二月（February）

- [使用 CodeLens 的扩展 / Extensions using CodeLens](https://code.visualstudio.com/blogs/2017/02/12/code-lens-roundup)
- [语法突出显示优化 / Optimizations in Syntax Highlighting](https://code.visualstudio.com/blogs/2017/02/08/syntax-highlighting-optimizations)

### 一月（January）

- [Node.js 开发与 Visual Studio Code 和 Azure / Node.js Development with Visual Studio Code and Azure](https://code.visualstudio.com/blogs/2017/01/15/connect-nina-e2e)

## 2016

### 十二月（December）

- [自定义 VS Code 扩展综述 / Customize VS Code Extension Roundup](https://code.visualstudio.com/blogs/2016/12/12/roundup-customize)

### 十一月（November）

- [Hot Exit 来到 Insiders / Hot Exit Comes to Insiders](https://code.visualstudio.com/blogs/2016/11/30/hot-exit-in-insiders)
- [创建格式化程序扩展 / Creating a Formatter Extension](https://code.visualstudio.com/blogs/2016/11/15/formatters-best-practices)
- [1.7 Rollback 事故报告 / 1.7 Rollback Incident Report](https://code.visualstudio.com/blogs/2016/11/3/rollback)

### 十月（October）

- [JavaScript 扩展第 2 部分 / JavaScript Extensions Part 2](https://code.visualstudio.com/blogs/2016/10/31/js_roundup_2)

### 九月（September）

- [JavaScript 扩展第 1 部分 / JavaScript Extensions Part 1](https://code.visualstudio.com/blogs/2016/09/14/js_roundup_1)
- [VS Code 中的文件和文件夹图标！/ File and Folder Icons in VS Code!](https://code.visualstudio.com/blogs/2016/09/08/icon-themes)

### 八月（August）

- [Windows 和 Mac 上的 iOS Web 调试 / iOS Web Debugging on Windows and Mac](https://code.visualstudio.com/blogs/2016/08/22/introducing-ios-web-debugging-for-vs-code-on-windows-and-mac)
- [再见 User Voice，你好 GitHub Reactions！/ Goodbye User Voice, Hello GitHub Reactions!](https://code.visualstudio.com/blogs/2016/08/19/goodbyeuservoice)
- [宣布推介视频 / Announcing Intro Videos](https://code.visualstudio.com/blogs/2016/08/15/introvideos)
- [扩展综述 - Git 的乐趣 / Extensions Roundup - Fun with Git](https://code.visualstudio.com/blogs/2016/07/29/extensions-roundup-git)

### 六月（June）

- [语言的通用协议 / A Common Protocol for Languages](https://code.visualstudio.com/blogs/2016/06/27/common-language-protocol)

### 五月（May）

- [Insiders Build 的演变 / Evolution of the Insiders Build](https://code.visualstudio.com/blogs/2016/05/23/evolution-of-insiders)
- [2016 年 4 月版本 / April 2016 Release](https://code.visualstudio.com/blogs/2016/05/09/April2016Release)
- [扩展综述 / Extensions Roundup](https://code.visualstudio.com/blogs/2016/05/04/extension-roundup-may)

### 四月（April）

- [Visual Studio Code 1.0！/ Visual Studio Code 1.0!](https://code.visualstudio.com/blogs/2016/04/14/vscode-1.0)

### 三月（March）

- [VS Code 扩展 / VS Code Extensions](https://code.visualstudio.com/blogs/2016/03/11/ExtensionsRoundup)
- [2016 年 2 月恢复版本 / February 2016 Recovery Release](https://code.visualstudio.com/blogs/2016/03/14/Feb2016Recovery)
- [2016 年 2 月版本 / February 2016 Release](https://code.visualstudio.com/blogs/2016/03/07/Feb2016Release)

### 二月（February）

- [为 VS Code 推介 Chrome 调试器 / Introducing Chrome Debugging for VS Code](https://code.visualstudio.com/blogs/2016/02/23/introducing-chrome-debugger-for-vs-code)
- [推介 Insiders Build / Introducing the Insiders Build](https://code.visualstudio.com/blogs/2016/02/01/introducing_insiders_build)
