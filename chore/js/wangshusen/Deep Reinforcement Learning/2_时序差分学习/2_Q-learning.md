# Q-learning

## 1. Q-learning 与 Sarsa 的核心区别

Q-learning 和 Sarsa 都是时序差分（TD）算法，但它们的目标和更新方式不同：

| 特性               | Sarsa                                                                          | Q-learning                                                              |
| :----------------- | :----------------------------------------------------------------------------- | :---------------------------------------------------------------------- |
| **学习目标**       | 动作价值函数 $Q_\pi$ (对应策略 $\pi$)                                          | **最优**动作价值函数 $Q^*$ (对应最优策略 $\pi^*$)                       |
| **TD Target 计算** | 使用实际执行的动作 $a_{t+1}$ <br> $y_t = r_t + \gamma Q_\pi(s_{t+1}, a_{t+1})$ | 使用最大化价值的动作 <br> $y_t = r_t + \gamma \max_{a} Q^*(s_{t+1}, a)$ |
| **性质**           | On-policy (同策略)                                                             | Off-policy (异策略)                                                     |

- **Sarsa**: 很多时候用于 Actor-Critic 中更新价值网络。
- **Q-learning**: 直接学习最优价值，是 **DQN** (Deep Q-Network) 的基础算法。

## 2. Q-learning 的推导

Q-learning 的 TD Target 源自 **贝尔曼最优方程 (Bellman Optimality Equation)**。

1.  **期望形式**：
    最优动作价值函数 $Q^*$ 可以写成期望的形式：
    $$Q^*(s_t, a_t) = \mathbb{E} [R_t + \gamma \max_{a} Q^*(s_{t+1}, a)]$$

    - 这里 $\max_{a} Q^*(s_{t+1}, a)$ 代表在下一状态 $s_{t+1}$ 下，做出最优动作所能获得的价值。
    - 因为 $a_{t+1}$ 是最优策略选出的动作，所以它天然满足最大化 $Q^*$ 的条件。

2.  **蒙特卡洛近似 (TD Target)**：
    直接求期望很困难，我们使用观测到的样本 $(s_t, a_t, r_t, s_{t+1})$ 进行近似：

    - 把期望中的 $R_t$ 替换为观测奖励 $r_t$。
    - 把期望中的 $S_{t+1}$ 替换为观测状态 $s_{t+1}$。

    得到 **TD Target** ($y_t$)：
    $$y_t = r_t + \gamma \max_{a} Q^*(s_{t+1}, a)$$

    - $y_t$ 部分基于真实观测（$r_t$），比单纯的预测更靠谱，因此作为训练的目标。

## 3. 算法形式

### 3.1 表格形式 (Tabular Q-learning)

适用于状态 $s$ 和动作 $a$ 数量有限的情况。$Q^*$ 是一个表格。

1.  观测到一个 Transition: $(s_t, a_t, r_t, s_{t+1})$。
2.  **计算 Target**: 查表找到 $s_{t+1}$ 对应的一行，找出最大值 $\max_a Q(s_{t+1}, a)$，计算 $y_t = r_t + \gamma \max_a Q(s_{t+1}, a)$。
3.  **计算 TD Error**: $\delta_t = y_t - Q(s_t, a_t)$。
4.  **更新表格**:
    $$Q(s_t, a_t) \leftarrow Q(s_t, a_t) + \alpha \cdot \delta_t$$

### 3.2 神经网络形式 (DQN)

当状态空间很大时，用神经网络 $Q(s, a; w)$ 来近似 $Q^*(s, a)$。

1.  **输入**: 状态 $s$。
2.  **输出**: 所有可能动作的打分（分数越高代表动作越好）。
3.  **动作选择**: Agent 根据 $Q$ 网络输出的分数选择动作（通常选分数最高的 $a_t$）。
4.  **参数更新**:
    - 观测 Transition $(s_t, a_t, r_t, s_{t+1})$。
    - 计算 TD Target: $y_t = r_t + \gamma \max_{a} Q(s_{t+1}, a; w)$。
    - 计算 TD Error (Loss): $\delta_t = Q(s_t, a_t; w) - y_t$。
    - **梯度下降**: 最小化 $\delta_t^2$，更新参数 $w$。

## 4. 总结

- Q-learning 旨在利用 TD Target 学习 **最优动作价值函数** $Q^*$。
- TD Target 的关键在于其中的 $\max$ 操作：**$y_t = r_t + \gamma \max_a Q(s_{t+1}, a)$**。
- DQN 就是 Q-learning 的神经网络版本，通过最小化预测值与 TD Target 的差距来训练网络。
