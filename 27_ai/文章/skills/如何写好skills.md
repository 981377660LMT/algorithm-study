## 基础

1. 设计、实现、验证、迭代
2. 为什么 Skills 不是通过语义相似度或 BM25 召回的，而是将 Skills 的 Frontmatter 信息常驻在上下文里呢？这是因为我们希望 Skill 是原子的、可二次组合的
3. Skill 必须由 AI 来编写，而不是你来写
   /skill-creator
4. 当前 SKILL.md 太长了，缺乏 progressive-disclosure 机制，但是也请不要滥用。
   Skill 里可以“机械式”执行的部分（如果有）请帮我用 Python 实现。
   Bash 命令和 CLI 是最好的工具
   保留 Human-in-the-Loop 的交互(AskUserQuestion)
5. /skill-creator是可进化(Agent RL)的：
   生成 → 你在真实任务里用它 → 失败/别扭的地方反馈 → 让 skill-creator 改 SKILL.md。
   每一轮都是一次人工 RLHF。skill-creator 自带了 eval（评估）能力
   Claude Code 会帮你自动生成测试用例，并用多个 Sub-agent 并行运行评估
6. 先写eval，再写skill；把一个能跑的例子改造成skill；allowed-tools 收窄权限

---

## 并行多任务

1. 为什么复杂skill要多任务并行
   并行子任务的价值拆解为三个维度——性能、隔离、多样性
   - 性能
     多 Agent 并行不是在做同一件事的“多线程”，而是在同一份数据上投入了“多份智力”
   - **上下文隔离**
   - 多样性
2. 何时不该使用 -> 依赖
3. 如何在skill中提示
   为提高效率，分析阶段应显式使用 spawn 这个单词来启发模型启动并行子任务（
4. single-skill 到 multi-skill
   把复杂子任务独立成一个子 Skill，然后在主 Skill 中通过 spawn 子 Agent 时指定使用哪个 Skill 来执行
