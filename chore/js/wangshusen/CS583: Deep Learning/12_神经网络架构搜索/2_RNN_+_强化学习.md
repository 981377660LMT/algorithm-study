# RNN*+*强化学习

这节课是王树森老师讲解 **神经网络结构搜索（Neural Architecture Search, NAS）** 的第二部分，重点介绍了一种**基于循环神经网络（RNN）和强化学习（RL）** 的 NAS 方法。

这篇 ICLR 2017 的论文（通常指 Zoph & Le 的 _Neural Architecture Search with Reinforcement Learning_）是 NAS 领域的开山之作，虽然现在已不再是主流（因为太慢了），但其思想极其重要。

以下是深度、有逻辑的课程复盘：

---

### 第一部分：核心思想与定位 (Core Idea)

- **历史地位**：这是 NAS 领域的开创性工作，提出了用机器学习模型来自动设计神经网络结构的想法。
- **基本逻辑**：
  - 用一个 **Controller RNN**（控制器）来生成 **Child Network**（子网络，比如一个 CNN）的结构描述（超参数）。
  - 训练这个 Child Network，得到验证集准确率（Validation Accuracy）。
  - 把这个准确率当由 **Reward（奖励）**，通过 **Reinforcement Learning（强化学习）** 反向更新 Controller RNN，让它下次生成的结构更好。
- **现状**：已过时。准确率不如新方法，且计算量大得惊人（几百块 GPU 跑几周）。只需理解思想，不必深究代码细节。

---

### 第二部分：Controller RNN 的工作流程 (How to generate Arch?)

Controller RNN 的任务是生成一串离散的 Token，这些 Token 描述了一个 CNN 的完整结构。

#### 1. 生成过程是“序列预测”

- 假设要生成 20 层 CNN，每层需要预测 3 个超参数（如卷积核数量、核大小、步长），那么总共需要预测 $20 \times 3 = 60$ 个超参数。
- RNN 运行 60 步，每一步输出一个超参数。

#### 2. RNN 内部机制

- **输入**：
  - $t=0$ 时：随机初始化的 $X_0$ 和状态 $H_0$。
  - $t>0$ 时：上一时刻的预测结果（Embedding 后的向量 $X_{t-1}$）。
- **输出**：每一时刻输出隐状态 $H_t$，接一个 **Softmax 分类器**。
- **决策**：
  - Softmax 输出概率分布 $P_t$（例如：预测核大小时，输出 3x3, 5x5, 7x7 的概率）。
  - 根据概率 **采样 (Sample)** 一个动作 $a_t$（例如选中 3x3），或者选概率最大的 Argmax。
- **注意点**：
  - **Softmax 不共享**：预测“核数量”的 Softmax 和预测“核大小”的 Softmax 是独立的参数矩阵。
  - **Embedding 共享**：预测相同类型参数（如第 1 层核大小、第 2 层核大小）时，输入端的 Embedding 层可以共享。

---

### 第三部分：训练策略 —— 为什么用强化学习？ (Why RL?)

这是一个非常关键的逻辑转折。

- **目标函数**：我们想最大化生成的 CNN 在验证集上的准确率 $R$。
- **困难**：
  - 准确率 $R$ 依赖于 CNN 的结构（即 RNN 的输出动作序列）。
  - 但是，准确率 $R$ 关于 RNN 的参数 $\theta$ 是**不可微 (Non-differentiable)** 的。你不能通过 $R$ 直接求导传回 RNN。
  - （因为从“离散的结构参数”到“准确率”中间隔了一个复杂的训练过程，且结构选择本身是离散的）。
- **解决方案**：**强化学习 (RL)**。
  - `RL 专门处理这种“目标不可微”、“奖励延迟”的序列决策问题。`
  - 这里把 Validation Accuracy 当作黑盒环境给出的 **Reward**。

---

### 第四部分：具体的强化学习算法 (REINFORCE Algorithm)

- **映射关系**：
  - **Agent**：Controller RNN。
  - **State**：当前的隐状态 $H_t$ 和上一轮选择 $X_t$。
  - **Action**：选择具体的超参数（如 3x3, Strides=2）。
  - **Reward**：生成的 CNN 训练收敛后的 Validation Accuracy。
- **奖励设置**：
  - RNN 跑了 60 步，生成了完整结构。
  - 然后去训练这个 CNN（可能要几小时）。
  - 得到准确率 $Acc$。
  - 这 60 步的每一步，其回报（Return） $U_t$ 都设为 $Acc$（除了最后一步才拿到奖励，前面奖励视为 0，累积回报即为最终奖励）。
- **参数更新**：
  - 使用 **REINFORCE (Policy Gradient)** 算法。
  - 公式：$\nabla_\theta J(\theta) \approx \sum_{t=1}^{60} \nabla_\theta \log \pi(a_t | s_t) \cdot U_t$。
  - **直觉**：如果这次生成的结构准确率高（$U_t$ 大），就通过梯度上升，增加生成这一串超参数动作的概率；反之则抑制。

---

### 第五部分：致命缺陷 —— 计算量 (Computation Cost)

这是劝退几乎所有普通研究者的原因。

- **Sample Complexity 高**：RL 需要大量的样本（Reward）才能收敛。
- **获取 Reward 极贵**：每一个 Reward 意味着要**从头训练一个 CNN**。
- **算账**：
  - 训练一个 CNN：假设 1 小时。
  - 需要尝试：10,000 次架构组合（这在 RL 里算少的）。
  - 总时间：10,000 GPU 小时。单卡要跑 1 年。
- **Google 的豪横**：这篇论文当时用了 800 块 GPU 并发跑了几个星期。

---

### 总结

1.  **方法论**：**Transformer / RNN 生成结构 -> 训练子网络 -> 获取准确率 -> RL (REINFORCE) 更新生成器**。
2.  **贡献**：证明了神经网络结构是可以被自动搜索出来的，而且能比人设计的还好。
3.  **缺点**：极其昂贵，完全不适合平民玩家。
4.  **下节课预告**：将介绍 **Differentiable Architecture Search (DARTS)**，这是一种更现代、更高效（单卡一天就能跑完）的方法，不再需要 RL，而是将离散搜索变为连续优化。
