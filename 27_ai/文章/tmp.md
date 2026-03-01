Alexander Embiricos
OpenAI Codex 产品负责人
访谈对象
Alexander Embiricos，OpenAI Codex 产品负责人
背景履历
Dropbox PM → 创办 Multi（Remotion，获 Greylock $13M 融资）→ 2024 被 OpenAI 收购
访谈平台
20VC 播客（Harry Stebbings 主持），2026年2月21日
核心主题
AI 编程范式跃迁、智能体生态重构、硬件异构化、市场竞争、SaaS 冲击

写在前面：为什么值得听
2026 年 2 月，AI 编程工具市场的竞争进入白热化。Cursor 月收入据传已逼近 10 亿美元量级，Claude Code 凭借终端极客口碑迅速蚕食市场，而 OpenAI 在编程领域的存在感一度被质疑。就在这个节点，OpenAI Codex 的产品负责人 Alexander Embiricos 坐到了 Harry Stebbings 的播客面前，做了一次极为坦诚的长对话。
之所以值得精读这场访谈，除了因为 Alex 手握 OpenAI CodeX 的方向盘，更因为他的个人经历赋予了他一个极其罕见的复合视角——他既是产品经理，又是连续创业者，还是亲手被 AI 诏安的人。

一、前置导读：Alexander Embiricos 是谁
1.1 Dropbox 时代
Alexander Embiricos 是 OpenAI Codex 的产品负责人，但是他的职业起点是 Dropbox 的产品经理，他的关键洞察源自 Dropbox 时期对 Slack 崛起的近距离观察。他亲眼看到一个产品层面「更优」的方案如何输给了一个「更有引力」的方案——这个教训后来深刻影响了他对 AI 工具产品形态的判断。
"For a while, we thought people should comment on documents in Dropbox... it was more optimal. However, what we saw is that Slack is just such a center of gravity of people talking to each other. Nobody wants to comment on the document. I just want to Slack you." —— Alex

他后来在讨论 AI 工具市场终局时，几乎完整复用了这个框架：最优的功能不一定赢，成为「日常重力中心」的产品才赢。当主持人问他「被 Dropbox 教会的最大教训是什么」时，他的回答是：
"If we build our agent purely as workflow automation, then it's always going to be like pulling teeth to get that thing started. But if you can build a system that people just love using even if they only use it for partial tasks, over time they'll get better and better at using it." —— Alex

1.2 创业五年：Multi（Remotion）的成与败
离开 Dropbox 后，Alex 联合创立了视频协作平台 Multi（前身 Remotion）。这家公司基于 Zoom 基础设施构建，主打远程结对编程与企业协作，获得了 Greylock 等顶级风投 1300 万美元融资。五年的创业经历，给了他关于「如何在不确定性中保持方向感」的最深刻教育。
在访谈最开头的个人动机问题中，他分享了创业至暗时刻：
"One of my darkest moments... was recognizing that I had spent the past few months trying to avoid losing and all of a sudden, I was like, oh my god, that is why I'm so unhappy and that's probably why the startup isn't going well." —— Alex

这种从「避免失败」到「追求胜利」的心理转换，他说自己需要「every now and then」反复校正。但他补充说，比「赢」更底层的驱动力是：「I just love building things and building things for people.」
2024 年，Multi 被 OpenAI 以人才收购（Acqui-hire）形式整合进企业级解决方案部门。Alex 的团队在开发者协作、多人在线交互方面的经验，恰好填补了 OpenAI 从「模型」走向「产品」过程中缺失的能力。

1.3 在 OpenAI：从「误判」到找到方向
Alex 在访谈快火环节坦诚了他最初的产品直觉是完全错误的：
"When I joined OpenAI, I thought that we would all just be hanging out with our computers screen sharing... within a year we'd have this agent that we're just talking to. That was completely wrong. The rate of progress in multimodal models was slower than I expected." —— Alex
多模态（视频和音频）的路走不通，真正跑通的路径是「Agent 通过代码操作电脑」。这个认知转变直接决定了 Codex 的产品方向——放弃追求科幻式的多模态交互，聚焦打磨代码执行能力。
人物小结
Dropbox 教会他「重力中心」，创业五年教会他在至暗时刻保持方向感，加入 OpenAI 后的认知纠偏教会他「正确比自负更重要」。这三段经历的叠加，让他成为了一个极其务实、不迷信直觉的产品领导者。理解这个人，才能理解他后续所有产品决策的底层逻辑。

二、编程会被「自动化」吗？
主持人 Harry 一上来就抛出了一个尖锐的行业问题：
Harry：Elon said that coding is one of the first professions to be largely automated. Do you agree given your position and what you see day to day?
Alex：For sure, I would agree that coding is one of the first domains where LLMs are really good. But what does it mean for coding to be automated? It's kind of a heavy statement, right?

Alex 做了一系列历史类比。
他先追溯了「computer」这个词的起源——在 Bletchley Park 时代，computer 指的是人类计算员，而非机器：
"There were humans who would punch out punch cards and put them into the machine and do a bunch of tabulated math... the first spreadsheet software was kind of loosely based off this idea that you would have an office full of desks arranged in a grid and people doing tabulations and then passing their sheets to the next person." —— Alex

然后他引出 —— 每一次「自动化」旧任务，都引发了对输出的需求爆炸：
"All these things, those specific tasks have become automated, but every time that's happened, there's been an explosion in demand for the output and so you need many more people actually to do that kind of work even if the specific task has changed." —— Alex

当 Harry 追问「所以你认为 5 年后工程师会更多而非更少？」时：
"Yeah. Sometimes we change what terms mean, right? The term computer now refers to something else. Now we have the term software engineer. I definitely think we'll have many more builders." —— Alex

2.1 人才栈压缩：前后端分工正在消融
但「更多 builder」不意味着工种不变。Alex 观察到一个结构性变化——人才技能栈正在被压缩（compression of the talent stack）：
"You still need software engineers today. You still need designers. I'm a PM. Do you need PMs? You could have some fun jokes about that. I don't think you need them." —— Alex

他以 Codex 团队为例说明，前端/后端的严格分工已经大幅减少，全栈化成为常态。行业数据印证了这一趋势：OpenAI 内部的 Sora Android 应用完全依赖 Codex 辅助，在 18 天内构建完成并在 10 天后登顶 App Store——这种效率在传统分工模式下不可想象。

2.2 「你不需要太多产品经理」
Harry ：为什么你觉得这个世界不需要 PM？
"I think it's incredibly hard to define what a PM is. I kind of think of the role as actually explicitly undefined and your goal is just to adapt to whatever the team or business needs." —— Alex

他进一步解释，PM 通常做的事情——「退后一步看全局、和 GTM 团队协作、做团队啦啦队和质量把关人」——这些都可以被一个足够强的工程 lead 或注重产品的设计师替代：
"All of those things I just described, which are maybe my current role, could be done by a really strong eng lead or a designer who thinks a lot about product. And so I think it's like often useful to have product managers, but you probably don't want many of them until the team is really large." —— Alex

[图片]

三、AGI 的真实瓶颈：不是模型，而是人类自己
Harry 在预研 Alex 过往言论时发现了一个有意思的观点，当场追问：
Harry：You said that human typing speed and validation work is the key bottleneck to AGI, not model compute or architecture. Help me understand what you really meant by that.
Alex：That's a fun one. Let's do this slightly Socratically. How many times would you say you use AI today?
Harry：30 plus times a day.
Alex：How many times do you think, assuming zero energy expenditure from you, AI could help you per day?
Harry：I think we'll have inference running 24 hours a day across every single thing.

接着 Alex 描述了他在 OpenAI 工程师身上看到的现象：
"I hear things from engineers who are telling me: I constantly have Codex running. I never close my laptop and if it's not running while I'm in a meeting I'm like wasting my time. I need to make sure Codex always has work that it's doing. And that's super cool but that's a lot of work to manage these agents." —— Alex

然后他回到核心论点——如果 AI 理论上可以帮你数万次，但你实际只用了 30 次，瓶颈显然不在机器端，而在人类端：
"I think AI should be helping us tens of thousands of times per day... but I'm too lazy to type out that many prompts and I am too uncreative to figure out all the ways that AI can help me." —— Alex

这段自嘲式的坦诚非常有力——连做这个产品的人自己都承认「太懒了」、「太缺乏创意了」去充分利用 AI，那普通用户呢？他因此给出了产品层面的解答：
"The world we want to get to is one where to use AI you don't really need to figure out the right way to prompt. It's just super easy for you and you don't even need to recognize that AI could help you. It's just connected to your context and chimes in helpfully." —— Alex

3.1「打字问题」与「验证速度」
Alex 的观察可以被概括为两道瓶颈。
第一道是「打字问题」（The Typing Problem）——当机器能在几秒内生成数千行复杂代码时，要求人类通过键盘输入详尽 Prompt 是严重的效率倒退。
第二道更致命的瓶颈是「验证速度」——阅读并审查 AI 生成的数百行重构代码，其认知负担往往超过自己从头写的成本。这种「生成快于验证」的不对称性意味着，不改变审查工作流，AI 创造的生产力将被卡在人类 QA 环节。

3.2 三阶段解法路线

Alex 将解决方案拆解为三个递进阶段：
"First, let's have agents work really well for software engineering and coding because LLMs happen to be good at that. Next, let's realize that for an agent to be useful more generally, using a computer is super valuable. And also we'll realize that all agents are actually coding agents because coding is just the best way for an agent to use a computer. Then finally, once we see what's working, let's build that productization where you have highly specific features that just work immediately out of the box for people." —— Alex
他特别强调这三个阶段会快速走完：「I think we're going to speedrun this entire one, two, three journey in the next months.」

3.3 对「产品化提示词」的看法

当 Harry 问是否 OpenAI 的工作就是「产品化人类的提示词和行动」时：
"I think it is our job to make sure that we have models with amazing capabilities and then eventually to get to a world where this is highly productized, and so you just have this magic text box or audio input or whatever, or you can just add AI to your group chat and it just starts to help." —— Alex

但他强调当前不应该急于产品化特定工作流，而是先提供灵活的开放工具。他引用了 Claude Code 首发时的成功作为正面案例：
"I think this was the genius of when Claude Code first shipped. What they really got right was they had this tool that was super easy to use in whatever context you wanted, just in your terminal, and people started experimenting with where to use it." —— Alex

四、产品战术：从「结对编程」到「完全委托」

4.1 GPT 5.2 Codex：拐点时刻

当 Harry 问 OpenAI 内部有多少代码是由 Codex 写的时，Alex 描述了一个剧变：
"Most people that I know are basically not opening editors anymore." —— Alex

他详细解释了这个转变的时间线：
"Before GPT 5.2 Codex, the kinds of AI features we were using to write code were tab completion or maybe pair programming with the model. You still needed to be at your laptop with your hands on the keyboard. Then at the time of GPT 5.2 Codex in December, we switched to: actually I'm just going to fully delegate this task. I'm going to have a plan with it, make sure we like the spec, and then I'm just going to let it cook." —— Alex

「Let it cook」—— 从「人类驾驶、AI 辅助」到「人类规划、AI 执行」，代码本身不再由人类编写：
"The code itself is not being written by humans anymore." —— Alex

4.2 为什么要做独立应用而非 IDE 插件

Harry 问出了很多人的疑问：24 个月后我们还会有 IDE 吗？
"If that's the answer, then yes, you could even argue the Codex app is an IDE. I don't think it is. For me, I think of an IDE as a really powerful editor. We explicitly didn't build editing into the Codex app because we wanted it to be really clear how you're meant to use it." —— Alex

Codex App 被精确定位为「多智能体任务编排中心」而非代码编辑器。它的三大核心能力是：

- 并行任务管理：同时运行多个智能体——一个跑集成测试，一个重构数据库，一个扫描安全漏洞。
- 计划模式（Plan Mode）：Agent 先输出架构计划并向人类提问，批准后才执行。Alex 类比为「新员工在编码前提交 RFC」。
- Skills 开放标准：Agent 可执行的非编码任务脚本，如分诊任务、监控部署等。

关于 Plan Mode 的价值，Alex 有一段解释：
"We recently shipped a very prominent plan mode where you have the agent go off and propose how it's going to do something. It's quite a long plan and then it asks you questions about if you agree. This is very similar to if you had a new hire who was new to your codebase and they had to present a request for comments to the rest of the team before they started doing the work." —— Alex

4.3 AI 代码自查的实践
面对「AI slop」（AI 垃圾代码）的问题，Alex 透露了 OpenAI 的应对方案：
"A common practice with Codex is to have Codex review its own PR. Codex is actually incredibly good at this. We've explicitly trained the model to be good at code review, including making sure it's really good at creating high signal feedback, so it'll basically have few false positives of criticism, which means you can really trust when it has feedback. Nearly all code at OpenAI is reviewed by Codex automatically whenever you push it to a repo." —— Alex
「高信号、低误报」——这意味着 AI 审查的可信度已经高到可以作为默认工具，这种模式可以概括为 Agent-to-Agent 的自动化审查闭环，即，人类从微观的代码编写者上升为系统级的指挥者。

五、竞争格局

5.1 「我们的工作是分发智能」

这是整场访谈中最「反直觉」的段落。当 Harry 说「你的工作是 Codex 的成功」时，Alex 的回应：
Harry：Your job is the success of Codex.
Alex：Actually our job is the distribution of intelligence. This is really unintuitive but we put all this effort into training these models and then we serve these models to our competitors.
Harry：This is so difficult for me as a venture capitalist to understand.
Alex：Because we're playing such a long game, if the competition gets better, we learn. It's actually helpful for us.

在具体实践中，Codex 推行了 agents.md 和 agents/skills 两个开放标准。Alex 含蓄地吐槽了竞争对手的态度：
"We helped push for putting skills to be stored in a neutral named folder called agents instead of in like Codex or something. And again, everyone has jumped on it except the usual suspect." —— Alex

行业数据显示，截至 2026 年 Q1，全球已有超 60,000 个 GitHub 仓库采用 AGENTS.md，引入后 AI 智能体运行时间平均缩短 28.64%。当上下文配置被标准化后，竞争无情地回归到模型核心能力——这恰恰是 OpenAI 的想要主导的领域。

5.2 对 Cursor 的评价

Harry 抛出了一个挑衅性的预测：「我认为 Cursor 今年会丢掉一半收入。」Alex 的反应很有趣：
Harry：I think Cursor will lose half of their revenue this year. Agree or disagree?
Alex：Can I just like no comment?... I think they've built a really successful business. We see them a lot when we're in enterprise.

Harry 追问 Cursor 还是 Claude Code 谁是更大的竞争对手：
"I see Cursor a lot more than Claude Code. My sort of narrative is that you have to meet people where they're at. What's coolest about Cursor from my perspective is that it meets developers exactly where they are. It's like: you used to be using VS Code, switch to Cursor, almost nothing is broken about your workflow, everything works, certain aspects just got better." —— Alex

关于 Cursor 自研模型的决策：
"I feel like there is a gap in the market for that kind of model... My thesis for how they win is that they meet everyone where they are and make it really easy to step up into using more advanced agentic workflows." —— Alex

5.3 对 Claude Code 的评价

Alex 对 Claude Code 的评价是整场访谈中最值得细品的段落之一。当 Harry 问「Claude Code 做对了什么」时：
"Way back last year they made something that was really easy to use and just worked with all your tools with zero setup by running it locally in your terminal. When we started investing much more in the Codex CLI and shipped great models for it like GPT-5, our growth exploded. So I think that idea of just meet people where they're at, give them something easy to use, let them ramp from there has been awesome. That's probably the biggest learning we've had from them." —— Alex

接着 Harry 追问了一个极有杀伤力的问题：
Harry：What mistake do you think they made that you've learned from having had the benefit of seeing them make it?
Alex：I think that they overindexed on their initial success with their command line interface tool. I think at the end of the day it's not the friendliest UI and it makes it hard to extend beyond pure builders. And it makes it difficult to truly delegate to agents because to delegate through that kind of interface you have to be a power user of your terminal or tmux or something.

这段分析对 Anthropic 的产品团队极具参考价值：Claude Code 的 CLI 优势同时也是它的天花板。Codex App 的诞生，本质上就是对 Claude Code 这个「战略失误（?）」的有意修正。

5.4 数据护城河：Anthropic 并没有你想象的优势

一位投资人通过 Harry 提了一个尖锐问题：Anthropic 是否在编程数据上有了不可逾越的优势？
"I definitely don't think they have a significant advantage in terms of data on coding... I feel like we have plenty enough data to build really good coding models. I actually think the place that's more interesting for getting data now is knowledge work tasks. That's data that's not really available most places on the internet." —— Alex

他甚至暗示了获取知识工作数据的激进方式：
"Maybe you have to pay people to simulate doing tasks so that you can learn these trajectories for the model. Maybe you should acquire startups that are no longer in business but have a lot of data, like say they're a Slack or something." —— Alex

5.5 市场终局：不是三分天下
"I think this might end up with fewer providers capturing a lot in the long run... I think we're in this temporary phase where we have agents that are really good at coding. But I think that's probably temporary and over time we're going to end up with agents that can do anything for you." —— Alex

他回到 Dropbox/Slack 的类比来论证：
"If there is a single agent you can use for nearly anything, there will just be this giant pull and everyone will talk about how they use that one agent for things. Teams will share best practices with each other, there'll be hackathons around how to use that best thing. And you'll end up with just a handful of these." —— Alex

六、硬件异构化：速度如何重塑一切
Harry 问 Alex 与 Cerebras 的合作意味着什么：
Harry：You partnered with Cerebras. How important is speed for developers when using Codex?
Alex：The simple answer is: it's super important.
Harry：Is it like an inference monopoly? You have it now and competitors don't?
Alex：I don't think we're going to end up in a monopolistic world. There's so much competitive pressure. But I will say that we have news coming about that partnership soon and I'm very excited for these things to ship.

行业背景
2026 年 1 月，OpenAI 与 Cerebras 达成超 100 亿美元、750 兆瓦部署规模的战略协议——全球最大规模的高速 AI 推理专用硬件部署。OpenAI 专门研发了极速模型 Codex-Spark（超 1000 Tokens/秒，是全量模型的 15 倍），与全量 GPT-5.3 形成「双模并发」架构。
Alex 还透露了模型效率的持续提升：
"With GPT 5.3 Codex, that model is significantly more efficient than prior models. We recently rolled out a change where in the API those models are served 40% faster and in Codex they're served 25% faster. Speed matters a lot and we're approaching it from all angles: both the hardware, how you do inference, and the model level." —— Alex

七、企业市场与 SaaS 生死论

7.1 自下而上 vs 自上而下的部署之争
Harry 挑战了 Alex 关于「不需要 FTE 介入」的观点：
Harry：Data security, permissioning, access provisions is really freaking hard and people are much less intelligent and confident than we give them credit for. You actually need an FTE to go in and custom fit solutions. Am I wrong?
Alex：You're right if you're trying to go all the way from zero to one... But what I've seen is that when we do these things top down, we end up massively underleveraging the potential of AI.

他用一个假设来论证为什么需要同步给一线员工提供 AI 工具：
"Imagine if you work in a customer support role and AI is being brought into your role and starting to automate meaningful chunks of your work, but you've never heard of ChatGPT nor are you allowed to use it. In this world you have no intuition for what this thing is. Whereas in a world where you've been using ChatGPT for work, you have much more intuition for how this works and you feel much more empowered." —— Alex

7.2 Atlas 浏览器：赋能与安全的灰犀牛
Alex 透露了 OpenAI 自己的浏览器：
"In OpenAI we're building a browser, Atlas. I think one of the key reasons is that by building a browser and controlling it tightly end to end, we can build safe agentic browsing for enterprise that is a way to access things agentically that are otherwise not yet built out by FTEs." —— Alex
安全警示：Atlas 的企业级风险
行业安全评估显示，Atlas 允许 Agent 在认证会话中自主执行操作，但传统 WAF/DLP 无法区分授权 Agent 行为与 Prompt 注入导致的失控。目前尚未通过 SOC 2/ISO 认证，缺乏 SIEM 集成，被官方禁止处理受监管数据（PHI/PCI），企业管理后台中 Agent 高级功能默认关闭。

7.3 SaaS 是否已死？
Harry 直接问：Salesforce 和 ServiceNow 跌了 20-40%，它们不应该跌吗？
"Does this SaaS company own a relationship with a human on the other end? If it does, then I suspect it's not going away. Or does it own some really important system of record? Probably not going away. On the other hand, is the SaaS company a kind of glue layer but it doesn't own either of those two things? I'm more nervous about that kind of company." —— Alex
他对几个具体类别给出了直接判断：

- Salesforce / ServiceNow 的下跌：「I don't think they should be [down that much].」
- Dropbox：「Respectfully, I think Dropbox is in a very difficult position.」
- 客户支持类 SaaS：「I do think you're going to come for customer support. I wouldn't want to be in that category.」
- Monday.com 式的项目管理工具：认为大多数用户不会真的去 vibe code 一个 to-do list，因此安全。

  7.4 给投资人的建议

Harry：If you were on my team as an investor, how would you think about areas for us to invest in?
Alex：Number one, I look for things with physical infrastructure. I don't think you're going into energy supply. And two is the fintech and banking integrations, gnarly financial products. I don't think OpenAI is going to go into building 500 relationships with banks in Southeast Asia.

他还指出了 AI 时代创业者素质的变化：
"There was this phase where you would invest in the person who can just build good product and you could ignore if they had a good thesis around a customer or go to market. I think that was an anomaly. Now maybe that kind of founder is not the founder you should invest in because it's relatively easier to build good product and you need to go back to investing in the founder who's thought through distribution." —— Alex

[图片]

八、推理成本、聊天 UI 与 Agent 间交互
8.1 推理是新的销售与营销？
Harry 引用了 Jason Lemkin 的观点「推理是新的销售与营销」，Alex 表示不太同意：
"I struggle with that. In this new world where anyone can build and it's increasingly easy to build things, having a good relationship with customer, knowing what they need is as hard as ever, maybe even harder as there's just more stuff in the market to choose from." —— Alex

8.2 聊天是否是终局 UI
Harry 引用 a16z 合伙人的观点——大多数人想要的是浏览器式的发现交互，而非聊天。Alex 的回答是一个二元架构：
"The simple answer is yes [chat will endure]. But there's two components. It's going to be some entity that I can talk to however I want about whatever I want. I shouldn't have to navigate to a place where I work with my coding AI and then have a different place for my sales AI. It's just like I'm going to talk to a thing and it's just going to help." —— Alex

但他补充说，专业用户需要专属的图形界面：
"If you had an executive assistant but you can only work by talking to them, that's super annoying. At some point you want to get to the show notes and look at them yourself and edit them yourself. So I think we'll pair chat with functional graphical interfaces that are bespoke to what someone needs." —— Alex

8.3 Agent-to-Agent 交互
当 Harry 问到 Agent 间交互的设计时，Alex 给出了一个务实的观察：
"The best interfaces for Codex to do work also tend to be the best interfaces for humans. So when people ask how to make their codebase more efficient for the agent, the answer is: have you looked at it yourself? Is it easy for a human to work with?" —— Alex

他举了一个具体例子：
"If you just set up most test runners, they emit all the outputs of all tests. As a human it's really annoying because you have to find the one that failed among hundreds of thousands of lines. Turns out that's terrible for AI as well. But if you filter it down to only emit the failed test — better for humans, also better for agents." —— Alex

九、增长数据、定价教训与团队势头

9.1 增长曲线

Alex 分享了少见的公开增长数据：
"Since August, we grew by like 20x. And then even late in the year, we doubled from December to now." —— Alex
结合行业数据（每周处理数万亿 Token），Codex 已是 OpenAI 调用量最大的编码模型。更重要的是，它已开始向免费 ChatGPT 用户开放部分功能，并投放了超级碗广告。

9.2 北极星指标
Harry：What is the metric you use as the defining north star?
Alex：It's actually not revenue as the primary. The primary is active users.
Harry：Weekly active? Is that frequent enough? If this is replacing the IDE, is daily active not better?
Alex：I think daily active will be better soon. I actually agree with the criticism. We should probably just be a daily.

9.3 定价教训
"For a while Codex Cloud was effectively unlimited. Every day that we left it that way we knew it would be harder to wind back. When we wound back that unlimited use to some more reasonable limit, there was a lot of blowback from users... The lesson I learned the hard way is you can't make things unlimited for too long." —— Alex

9.4 团队感受到的势头变化
Harry 问团队是否能感受到「风向变化」：
"Absolutely. If you look at the history of Codex, the first thing we launched last year was this amazing idea that people were super excited about — give the agent its own computer in the cloud, have as many as you want work in parallel. Super great idea. To be honest, it didn't work as well as what we shipped later. It was not the best." —— Alex
他坦诚地描述了改进路径：
"We had feedback around our model being slower and less fun to work with and being less good at communicating with you while it was working. We addressed that feedback... with the app, the feedback has been resounding from the market that this is a really high quality experience. It's unintuitively simple and people are just loving it. Even our biggest critics are converted." —— Alex

十、人才战略与对年轻工程师的建议
10.1 残酷的人才竞争现实
"The war for talent is incredibly fierce right now. At OpenAI we have an incredibly strong brand so we're able to attract a lot of talent. But even so, we put a ton of effort into closing candidates that we're really excited about. Even we feel that — it's not like you don't just get whoever you want for free." —— Alex

10.2 PM 招聘的极端标准
"I think [PMs] have to be the perfect fit. And if you have someone who's not the perfect fit, they might just do more harm than good. So it's kind of means that we're way more selective than I might have been in other roles." —— Alex

10.3 给 CS 学生的建议
Harry 设定了场景：「我是斯坦福/帝国理工/剑桥的 CS 学生，你会给我什么建议？」
"There's actually never been a better time to be an engineer because you have incredible tooling available to you. Your ability to ramp into a complex codebase has never been faster, because you can go and ask AI a ton of questions about the codebase." —— Alex
然后是关键的转折——在构建变得容易的世界里，什么变得稀缺：
"Because it's never been easier to build things, the thing that becomes scarcer is agency, taste, and quality. I would urge you to just build things and demonstrate your agency and your taste around what you build. When someone writes to me with some interesting thoughts and a link to an interesting project, that gets my attention much more than a normal resume does." —— Alex

Agency + Taste + Quality：AI 时代的人才三角
这不是空话。当 Alex 说「一个有趣项目链接比标准简历更吸引我」时，他在描述整个行业的招聘趋势变化。技术能力正在被 AI 工具拉平，而「判断该做什么」的能力无法被工具替代。构建有品味的项目并公开分享，是 AI 时代年轻工程师最有效的差异化策略。

十一、Rapid-fire 环节精选
11.1 对 Anthropic 广告的回应
Harry：Do you think the response to Anthropic's ads was the right response?
Alex：One company's being pretty negative about the future and the other company — us, OpenAI — is being really positive and just telling people they can build things and to dream. I thought that response was brilliant.

11.2 基准测试的价值
"They kind of give you a good measure of intelligence. But then you have to pair that with what it feels like to use the model. And that's a vibes thing. I'm always surprised by how vibes-based the evaluation of how it feels to work with a model is." —— Alex

11.3 利润率：先跑马圈地
"I think costs are going to come down significantly. And this is the year where agents are going to have to be connected to all these various systems and I think that's going to be very sticky. So I view this year as a race. I think you want to win that race and you should be okay taking some hit to margin in the meantime." —— Alex

11.4 五年后回望今天：什么会被淘汰
"One is just editing code by hand. Another one might even be actually managing the deployment and monitoring of systems by hand." —— Alex
他预测了一种全新的创业起步方式：
"The way you start a company is you start by getting an agent and asking it to build things and then you get more agents and then maybe eventually you add your co-founders to this service. Your main communication tool is actually your agent communication tool." —— Alex

11.5 愿景：家庭 WhatsApp 里的 AI
"What I'm most excited for is to get to a form factor for AI that means they're just helping everyone regardless of whether they're in tech and especially if they're not in tech or especially if they're older. At some point we'll add an agent to our family WhatsApp or something and it'll just start being useful to the family without anyone having to think harder about it than that." —— Alex

十二、写在最后

1. Alex 对 PM 角色的自我解构。当 AI 完成大部分执行层面工作时，PM 的价值从「管理流程」转向「判断什么值得做」。纯粹的「协调型 PM」将被淘汰，具备业务洞察力和技术感知力的「超级个体」将主导。
2. 「不要过早产品化」是 Alex 反复传达的信号。在 AI 能力快速迭代时，过早锁定工作流是风险。更好的策略是提供灵活的开放工具、观察用户创造性使用、然后基于数据做产品化。同时，纯粹的产品能力不再是护城河——分发能力、客户关系和领域专业度才是真正的壁垒。
3. Alex 对竞争格局的描述。他承认 Cursor 在企业市场存在感更强，坦诚 Codex 早期产品体验不够好，高度赞扬 Claude Code 的首发设计。这种坦率本身就是信号——他的团队对竞争态势有清醒认知，也对自身核心优势（模型能力、ChatGPT 分发、早期模型访问权、Cerebras 硬件联盟）充满信心。

时间窗口：Alex 反复强调「今年是 Agent 连接各种系统的竞赛之年」。一旦连接建立（Agent 接入 Sentry、Google Docs、内部系统等），切换成本将大幅上升。结合 AGENTS.md 标准化加速淘汰中间层、Cerebras 硬件联盟确保推理速度壁垒、Atlas 浏览器试图建立企业级入口，OpenAI 正在执行一套「基础设施→统一标准→控制入口→锁定用户」的环环相扣的战略。
2025-2026 年是 AI 编程工具市场的「跑马圈地」关键窗口期。这场竞争的终局不在于谁的代码补全更准确，而在于谁能成为知识工作者日常工作的「重力中心」。Alex 从 Dropbox 时代就理解这个道理，而他现在正在 OpenAI 把它变成现实。
