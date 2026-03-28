## 深入解读：How VS Code Builds with AI

**发布日期**: 2026 年 3 月 13 日
**作者**: Pierce Boggan（GitHub Copilot PM）& Peng Lyu（VS Code 工程经理）

### 文章定位

这篇文章的重要性远超一般的技术博客。它是 VS Code 团队首次系统性地公开**自身如何使用 AI Agent 来构建 VS Code** 的深度内幕。更重要的是，它揭示了一个里程碑事件的幕后原因：**VS Code 在坚持了十年月度发布后，转向了周度发布（Weekly Stable Releases）**——而这一转型的驱动力正是 AI Agent。

这篇文章本质上是一份**"Agent 驱动的软件工程实践白皮书"**，值得逐节深入分析。

---

### 1. 核心事件：十年月度发布 → 周度发布

#### 背景

VS Code 自 2016 年以来一直保持严格的月度发布节奏。每月经历的完整周期：

```
规划 → 构建 → 测试 → Endgame Week（交叉测试） → 撰写发布说明 → 发布
```

这套运作模式已经成为团队文化的一部分，每个人轮流负责不同角色。月度节奏给了充足的缓冲时间。

#### 为什么能转向周度？

文章给出了直接答案：**Agent 自动化了过去需要一周才能完成的开销性工作**。具体来说：

- 问题分类（Triage）→ Agent 驱动
- 提交总结 → Agent 自动生成
- 发布说明 → Agent 辅助撰写
- 代码审查 → Copilot Code Review 前置
- 回归测试 → Agent + Playwright 自动验证
- 文档更新 → 探索 Agent 自动检测过时文档

**💡 深入分析**：这不仅仅是"用 AI 写代码更快了"这么简单。VS Code 团队认识到，**发布节奏的瓶颈不是编码速度，而是围绕编码的一切运维开销（overhead）**。100+ 个日提交、每天数百个 issue、跨团队交叉测试——这些才是月度发布的真正约束。Agent 自动化的是这些"编码之外的编码工作"。

这也解释了为什么 v1.111 的 release notes 中突然宣布转向 Weekly Stable：团队已经在内部验证了这套流程足够可靠。

#### 对用户的实际影响

> A bug fix that used to wait three weeks now ships in days.  
> A feature merged on Monday can be in developers' editors that same week.

"发布 → 反馈 → 迭代"的循环从**月级**压缩到**天级**。这也是为什么最近几个版本（v1.110 → v1.111 → v1.113）的功能演进速度明显加快。

---

### 2. 六条核心实践经验

文章提炼了 VS Code 团队在 Agent 驱动开发中总结的六条经验，每一条都值得深入探讨：

#### A. 并行化自己（Parallelize Yourself）

> 在切换上下文之前，先启动 3-4 个 Agent 会话。

**Peng Lyu 的日常工作流**：

1. 更新 VS Code Insiders（每天两次 Insiders 构建）
2. 运行自定义 Agent 通过 Work IQ 拉取当日会议日程
3. Agent 生成任务快照 → 分为"自己做（需要人的决策）"和"开放代码任务（可委托给 Agent）"
4. 在进入第一个会议前，已经有多个 Agent 在并行执行任务

**实现手段**：

- 多个 VS Code 窗口 / Git Worktree
- Copilot CLI 后台 Agent（云端运行）
- Claude Agent 并发会话

**💡 深入分析**：这彻底颠覆了 Paul Graham 经典的"Maker's Schedule vs Manager's Schedule"二元对立。传统观点认为管理者的碎片化日程与开发者需要的大块专注时间不兼容。但 Agent 让管理者能在会议间隙并行推进编码任务——**不需要大块时间，因为 Agent 在你开会时替你"做"，你只需要在会后"审"**。

这也意味着**工程管理者的角色正在向"Agent 调度员"演进**——核心技能从"自己写代码"变成"分解任务 + 启动 Agent + 审查产出"。

#### B. 跳过中间产物（Skip the Intermediate Artifacts）

传统流程：

```
会议 → 会议纪要 → Issue → 规格文档 → 代码 → PR
```

现在的流程：

```
会议 → Agent 会话 → 代码 → PR
```

> "I don't write down meeting notes anymore. I'm kicking off the agents directly." — Peng Lyu

**💡 深入分析**：这是整篇文章中最激进的观点。它主张**文档作为中间产物正在被 Agent 会话替代**。传统的 spec/PRD 本质上是一种"假设"——你在编码之前猜测体验应该是什么样的。而现在 PM 可以直接让 Agent 把假设变成原型 PR，用实际体验代替文字描述。

文章举了一个真实案例：Pierce（PM）自己用 Agent 生成了 Copilot Chat 的 Fork 功能 PR，然后和工程师 Justin 一起在办公室审查、调整 CSS、合并。这个功能现在已经发布在 VS Code 中。

这对传统软件工程流程的冲击是深远的。**如果 PM 都能产出可合并的 PR，那么 spec 文档存在的意义是什么？** 答案是：spec 的目的是对齐认知，但 PR 是更好的认知对齐工具——它是具体的、可运行的、可测试的。

#### C. 自动化随速度扩展的开销（Automate the Overhead That Scales with Velocity）

这是最有工程实践价值的一节。VS Code 团队自动化了以下流水线：

| 流水线                   | 实现方式                             | 触发条件     |
| ------------------------ | ------------------------------------ | ------------ |
| **提交总结**             | 自定义 Slash Command + 快速模型      | 每日 / 按需  |
| **Insiders 更新日志**    | 同上，自动生成并推送                 | 提交到 main  |
| **X（Twitter）自动发帖** | Copilot SDK + GitHub Actions         | 提交到 main  |
| **Issue 分类**           | Agent Loop in GitHub Actions         | Issue 被创建 |
| **重复检测**             | Agent + 置信度评分                   | Issue 被创建 |
| **自动分配 Owner**       | Agent 读取 ownership docs + 历史模式 | Issue 被创建 |
| **标签建议**             | 同上                                 | Issue 被创建 |

技术栈核心：**Copilot CLI + Copilot SDK + GitHub Actions**

**量化结果**（2025 vs 2026 年 1-3 月同比）：

- 提交量：2,339 → 5,104（**增长 2.2 倍**）
- 关闭 Issue 数：2,916 → 8,402（**增长 2.9 倍**）

**💡 深入分析**：这组数据非常有说服力。提交量翻倍意味着代码产出量翻倍，而 Issue 关闭量接近 3 倍说明不仅产出增加了，**响应社区反馈的速度也显著提升**。这与"Agent 自动化 triage"直接相关——更快的分类 → 更快找到正确的工程师 → 更快修复。

更值得注意的是，团队还专门建了一个 **Chrome 扩展**，在 GitHub Issue 页面直接展示 triage 建议（重复检测、推荐 owner、标签建议）和跨团队的 Issue 状态仪表盘。这说明他们不是把 Agent 当玩具，而是真正构建了**生产级的 Agent 辅助基础设施**。

#### D. 在追求速度之前投资测试（Invest in Harnesses Before Speed）

> "Without the right harness, for the first week or two your productivity is really high. Then you quickly reach a ceiling where you keep regressing." — Peng Lyu

这是最清醒的一条经验。三层质量保障体系：

**第一层：自动化验证（Playwright Agent）**

- 自定义 Agent 使用 Playwright MCP Server
- 启动 VS Code → 导航到被测功能 → 截图 → 评估是否符合预期
- 如果截图显示问题，**Agent 自动修复并重新验证**
- 截图存档供人工复核

**第二层：测试套件 + 黄金场景（Golden Scenarios）**

- 传统单元/集成测试是底线
- "黄金场景"：核心用户流的预期行为规格
- 之前在月度 Endgame Week 中手动测试
- 现在交给 Agent 作为自动化的合并后验证
- 探索中：自动生成 demo 视频（PR 合并 → 视频生成 → 直接用于 changelog / 推文）

**第三层：代码审查**

- 每个 PR 自动触发 Copilot Code Review
- **工程师必须先解决 Copilot 的评论，才能请求人工审查**
- 六个月前这条规则不可行（噪音太大），现在模型质量大幅提升后成为标准流程
- Slack 协作：Bot 发 PR 链接 + CI/Review 状态实时更新
- 文化准则："give one, take one"——提交一个 PR，就去审查一个

**💡 深入分析**：Peng 的那句话是整篇文章最重要的警告。他坦率地承认：**Agent 驱动的高速开发如果没有足够的质量护栏，很快就会自我崩溃**。这不是理论——这是 VS Code 团队在实际中踩过的坑。

Playwright Agent 自动验证回路特别值得关注。这正是 v1.110 引入的 Browser Tools 的内部最佳实践——Agent 不仅写代码，还自己验证、自己修复。VS Code 团队自己就是这个功能的第一批重度用户。

#### E. 所有权在演变（Ownership is Evolving）

> When PMs, engineers from other areas, community contributors, and agents can contribute to any component, traditional ownership models need to adapt. Accountability for outcomes still rests with engineers.

**💡 深入分析**：这触及了 Agent 时代最敏感的组织问题。当 PM 能提交 PR、Agent 能跨组件修改代码、社区贡献者能深入核心模块时，"谁拥有这段代码？"的答案变得模糊。

VS Code 团队的立场很清晰：**产出可以由任何人或 Agent 完成，但工程师对结果负责**。这意味着工程师的角色从"代码编写者"转变为"代码质量守门人"——你不一定写了这段代码，但你负责确保它是正确的、可维护的、符合架构的。

#### F. 人类把关品味（Keep Humans in the Loop for Taste）

> Agents check correctness. Humans evaluate delight.

Pierce 提出了一个有意思的概念——**基于品味的评分（Taste-based Grading）**：

1. 写下你期望的**定性体验**（不是功能规格，而是"用起来应该是什么感觉"）
2. 让 Agent 评估实现是否匹配
3. 约 80% 的 Agent 观察有用，20% 需要忽略

**💡 深入分析**：这是对"AI 替代人类"叙事的最好回应。VS Code 团队明确区分了两类评估：

- **正确性（Correctness）**→ Agent 可以做，而且比人快
- **愉悦感（Delight）**→ 只有人类能判断

一个功能可以通过所有测试、没有 bug，但"用起来怪怪的"。这种微妙的产品判断力目前完全无法自动化。文章暗示 Endgame Week 不会被取消，只是被压缩——"交叉体验测试"仍然是人类的工作。

---

### 3. Agent-Ready 代码库评估

文章最后提出了一个开放性概念：**Agent-Ready Codebase Assessment（代码库的 Agent 就绪度评估）**。

核心标准：

- **结构**：Agent 能否找到正确的组件？
- **文档**：Agent 能否理解组件的职责和接口？
- **测试覆盖**：Agent 能否检测回归？

检验方法：如果一个 PM 把问题扔给 Agent 就能得到 reasonable 的 PR，说明代码库的 Agent 就绪度高。如果 Agent 挣扎，那就是一个信号——需要改善结构、文档或测试。

**💡 深入分析**：这个概念可能会成为未来衡量代码库质量的新维度。传统的代码质量指标包括测试覆盖率、圈复杂度、技术债务等。但 "Agent 就绪度"是一个全新维度——它衡量的是**代码库对自动化理解和修改的友好程度**。

几个实际推论：

- **好的代码结构**不仅对人类有益，对 Agent 也有益
- **好的注释和文档**不再只是"锦上添花"，而是 Agent 能否正确工作的关键
- **好的测试**不再只是防回归，而是 Agent 自我验证的基础设施

---

### 4. 总结：VS Code 团队揭示的 Agent 时代软件工程范式

| 维度           | 传统模式                   | Agent 驱动模式                      |
| -------------- | -------------------------- | ----------------------------------- |
| **发布节奏**   | 月度（受限于人工开销）     | 周度（开销被 Agent 自动化）         |
| **工作模式**   | 串行（一次做一件事）       | 并行（多个 Agent 会话 + 人工审查）  |
| **产出链条**   | 会议 → 文档 → Issue → 代码 | 会议 → Agent → PR                   |
| **PM 角色**    | 写 spec，交给工程师        | 直接产出原型 PR                     |
| **工程师角色** | 编码为主                   | 审查 + 架构 + 品味把关              |
| **代码所有权** | 按组件分配给个人           | 任何人/Agent 可贡献，工程师负责结果 |
| **质量保障**   | 人工 Endgame Week          | Agent 自动验证 + 人工品味评估       |
| **Issue 管理** | 人工轮值 triage            | Agent 自动分类 + 分配 + 查重        |
| **代码审查**   | 纯人工                     | Copilot 前置审查 → 人工终审         |

**核心洞察**：VS Code 团队展示的不是"用 AI 写代码"，而是**用 Agent 重构了整个软件交付流水线**。编码只是其中一个环节，triage、review、testing、release notes、documentation 都被 Agent 渗透。当这些环节都加速后，月度发布的瓶颈消失了，周度发布自然成为可能。

**最深层的启示**：Agent 不是让个人效率提升 10 倍的魔法，而是**让组织的运转摩擦趋近于零**的系统性工具。会议到代码的中间产物被消除，信息在人和 Agent 之间无缝流动，质量保障从人工瓶颈变成自动化流水线。这才是"AI 原生软件工程"的真正含义。
