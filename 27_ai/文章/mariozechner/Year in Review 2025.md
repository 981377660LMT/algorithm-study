# [Year in Review 2025](https://mariozechner.at/posts/2025-12-22-year-in-review-2025/)

# 深入分析：Mario Zechner 的 2025 年度回顾

这是一篇来自奥地利开发者、教练和演讲者 **Mario Zechner** 的年度技术回顾。这篇文章远不只是一份项目清单——它是一份关于 AI 编码代理（coding agents）实践现状、开源精神、技术行动主义以及个人成长的深刻反思。以下是我从多个维度提炼出的关键洞见。

---

## 1. 关于 AI 编码代理：「没有人知道该怎么正确使用它们」

这是全文最核心的判断：

> "Nobody knows yet how to do this properly. We are all just throwing shit at the wall, declaring victory, while secretly crying over all the tech debt we introduced into our codebases by letting agents run amok."

### 🔍 洞见

Mario 将自己归类为 **"紧绳派"（tight leash）**——始终保持人在回路中（human-in-the-loop），而非"代理军团派"（orchestrate armies of agents）。他的理由很实际：

- **"代理军团"模式对他从没真正奏效过**，除了研究类任务。让代理大量写代码最终是"灾难的配方"。
- 他指出了一个尖锐的观察：**声称掌握"代理军团"方法的人很少公开发布他们的代码**，而他自己尽量开源一切。这暗示着那些宣称的"成功"可能经不起公开审视。
- 他最终放弃了 Claude Code，自己写了一个极简编码代理 **pi**，核心工具只有四个：`read`、`write`、`edit`、`bash`。其他一切都是扩展层。

**深层启示**：当前 AI 编码工具的价值更多体现在 **降低启动新项目的心理门槛**，而非真正的"效率倍增"。Mario 自己也承认：

> "Did my productivity increase via LLMs? I don't actually know... It's all just vibes."

---

## 2. MCP vs CLI：简单胜于复杂

Mario 对 MCP（Model Context Protocol）的探索和放弃路径极具参考价值：

1. **初探**：为 Claude Code 构建了 MCP 服务器（`vs-claude`、`mailcp`）
2. **发现问题**：MCP 服务器的输出只能通过 LLM 在上下文中处理，这会**大量消耗 context window**
3. **转向 CLI**：bash 管道天然支持组合性（composability），CLI 工具更容易创建、维护和调试
4. **实证对比**：做了非科学性评估，结论是——**对大多数场景 MCP 和 CLI 差别不大，但 CLI 更实用**

### 🔍 洞见

> "What matters more is that whatever tool you use, it doesn't shit a ton of tokens into your context."

这与当前 AI 社区对 MCP 的狂热形成了鲜明对比。Mario 的经验指向一个更根本的原则：**上下文窗口是稀缺资源，工具设计的首要原则是节约 token**。

---

## 3. 领域专家 + AI = 真正的杠杆点

教妻子（语言学家）使用 Claude Code 的案例是全文最温暖也最有洞察力的部分：

- 妻子 Steffi 需要分析 32 位说话者、18,000+ 行语言学标注数据
- 不到两晚的教学时间，她就能独立用 Claude Code 构建可复现的数据分析管道
- 关键不是教她写 Python，而是教她 **"把工作结构化为小脚本组成的管道"** 的思维方式

### 🔍 洞见

> "Agents are really only effective in the hands of domain experts."

但紧跟着一个未解难题：

> "I still struggle to come up with teaching materials that could be deployed to a bigger audience of non-technical domain experts... Sitting down with everyone for two nights is not a scalable approach."

**这是整个 AI 教育领域的核心瓶颈**：每个领域专家的技术基线不同，需要高度个性化的指导。这也解释了为什么 AI 工具的真正普及比想象中更难——不是技术不够好，而是**桥接知识鸿沟的pedagogical方法还没有找到**。

---

## 4. Sitegeist 与浏览器代理安全的根本矛盾

Mario 构建的浏览器代理 Sitegeist 比 Claude for Chrome、OpenAI Atlas 等大厂产品表现更好，但他同时发现了一个**无法解决的安全问题**：

- 要让浏览器代理真正有用，需要通过 **debugger API** 与网站交互
- 这意味着模型可以访问 **HTTP-only cookies、认证凭据和浏览器端机密**
- 通过 prompt injection，攻击者可以轻易诱导模型将这些数据发送到外部端点
- 他尝试了沙箱隔离（JailJS），但**无法覆盖所有攻击面**

### 🔍 洞见

> "I personally think that Anthropic putting Claude for Chrome into the hands of normies is a really bad idea."

但他也自我调侃式地指出：

> "Maybe the great prompt injection and exfiltration crisis is averted simply by the fact that these agents aren't fit for the general populace yet."

**这是一个"安全通过难用来保障"的悖论**——目前浏览器代理的安全性本质上依赖于它们还不够易用，普通人不会用。一旦产品体验改善，安全问题就会爆发。

---

## 5. 用「破烂号指标」调查新闻业 AI 使用

Mario 用 **em-dash/en-dash 频率分析** 来检测奥地利新闻媒体是否使用 LLM 写文章，方法论优雅而巧妙：

1. **基本假设**：LLM 偏好使用 em-dash 和 en-dash
2. **控制变量**：如果 dash 增加的同时 hyphen 减少，可能是 CMS 自动替换；如果 dash 爆增而 hyphen 不变，则是新的 dash 密集内容被引入（即 LLM 生成）
3. **差异分析**：
   - Falter 和 Krone 呈现 CMS 模式（替换）
   - OE24、Heute、Exxpress 呈现 LLM 模式（新增内容）
4. **与公共资金关联**：这些媒体接受了政府的"数字化转型"资助

### 🔍 洞见

这种调查方法的价值在于**将模糊的怀疑转化为可量化的证据**。它也揭示了一个政策黑洞：**奥地利纳税人资助了这些 AI 部署，但没有任何媒体向读者披露 LLM 的使用**。

---

## 6. 硬件项目 Boxie：技术教育的种子

Boxie（为三岁儿子制作的离线音频播放器）是 Mario 最自豪的项目。从闪烁 LED 开始，到自行设计 PCB、焊接 SMT 元件、移植 Doom 到嵌入式系统——这是一段完整的硬件学习旅程。

### 🔍 洞见

> "Anytime something breaks... the boy and I get to unscrew the device together, look inside its internals, and dig around in them. It is my hope that this early exposure to tech will spark in him the interest to learn how these things work, just like my Game Boy sparked that interest in me."

这份对 **"可打开的技术"** 的执着——在一个万物封闭的时代——体现了一种深刻的教育哲学：**理解始于拆解**。

---

## 7. 关于开源和商业化的内心挣扎

贯穿全文的一个潜在主题是 Mario 在 **开源 vs 商业化** 之间的反复犹豫：

- **Yakety**（语音转录）：原计划商业化 → 被其他事分心 → 开源
- **Sitegeist**（浏览器代理）：原计划商业化 → 觉得给技术用户收费"不对" → 即将开源
- **pi**（编码代理）：直接开源 → 获得社区贡献 → "I'm kinda glad I'm back in the open source game"

### 🔍 洞见

Mario 最终选择通过开源建立影响力，再通过影响力为慈善项目（Cards for Ukraine）筹款。这是一种**非传统的"开源→社交资本→慈善"变现路径**。

---

## 8. 宏观判断

全文读完，Mario Zechner 的 2025 年可以浓缩为几个核心观点：

| 主题                   | 判断                                          |
| ---------------------- | --------------------------------------------- |
| AI 编码代理            | 有用但被高估，"紧绳"比"放手"更靠谱            |
| MCP                    | 过度炒作，CLI 在大多数场景下更实用            |
| AI 安全                | 浏览器代理的 prompt injection 问题本质上无解  |
| AI 教育                | 领域专家受益最大，但教学方法不可规模化        |
| AI 在媒体/教育中的滥用 | 严重且缺乏监管                                |
| 个人生产力             | "It's all just vibes" —— 感觉有提升但无法证明 |

---

## 最后

这篇文章最令人钦佩的品质是**诚实**。在一个充斥着"10x productivity"宣言的年代，Mario 说的是："我不知道我是否真的更高效了。" 在一个人人兜售 AI 解决方案的时代，他说的是："没有人知道该怎么正确做这件事。如果有人说他找到了答案，别信。"

这种诚实本身就是最有价值的洞见。
