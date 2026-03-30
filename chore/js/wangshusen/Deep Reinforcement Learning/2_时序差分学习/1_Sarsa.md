# Sarsa

根据王树森老师的课件内容，我为你整理了关于 **SARSA 算法** 的核心讲解。

这节课的重点是推导并讲解 **Temporal Difference (TD)** 算法家族中的 **SARSA**。

---

### 1. 数学基石：从“回报”到“TD Target”

理解 SARSA 的核心在于理解 **TD Target（时序差分目标）** 是怎么来的。

- **折扣回报公式 (Discounted Return):**
  $t$ 时刻的回报 $U_t$ 可以递归地写成：
  $$U_t = R_t + \gamma \cdot U_{t+1}$$
  即：**今天的总价值 = 今天的奖励 + 衰减后的明天的总价值**。

- **动作价值函数 (Q 函数):**
  $Q_\pi(s_t, a_t)$ 本质上是对回报 $U_t$ 的期望。根据上面的递归公式，我们可以得出：
  $$Q_\pi(s_t, a_t) \approx \mathbb{E}[R_t + \gamma \cdot Q_\pi(s_{t+1}, a_{t+1})]$$

- **蒙特卡洛近似 (Monte Carlo Approximation):**
  数学期望 $\mathbb{E}$ 很难直接求，但在实际训练中，我们可以用**一次真实的观测（采样）**来代替期望。
  我们将观测到的奖励 $r_t$ 和下一个状态动作的估值 $Q(s_{t+1}, a_{t+1})$ 代入，就得到了 **TD Target ($y_t$)**：
  $$y_t = r_t + \gamma \cdot Q(s_{t+1}, a_{t+1})$$

### 2. 核心思想：TD Learning (时序差分学习)

这一部分是算法的灵魂：**为什么要更新 Q 值？**

- **现状**：$Q(s_t, a_t)$ 是模型对当前价值的**估计（Guess）**。
- **目标**：$y_t$ (TD Target) 包含了一部分**真实发生的奖励 $r_t$** 加上一部分未来的估计。
- **结论**：因为 $y_t$ 含有真实观测数据，虽然不完全准，但比单纯的估计 $Q(s_t, a_t)$ 更可靠。
- **做法**：**让 $Q(s_t, a_t)$ 去通过更新，逼近 $y_t$。**

---

### 3. SARSA 算法详解

SARSA 是一种具体的 TD 算法。

#### (1) 名字由来

它的名字来源于一次更新所需要的**五元组**数据：
**S**tate, **A**ction, **R**eward, **S**tate (next), **A**ction (next)
$\rightarrow$ **S A R S A**

#### (2) 两种实现形式

**A. 表格形式 (Tabular SARSA)**

- **适用场景**：状态和动作数量有限（比如走迷宫）。
- **流程**：
  1.  观测到五元组 $(s_t, a_t, r_t, s_{t+1}, a_{t+1})$。
  2.  计算 **TD Target**: $y_t = r_t + \gamma \cdot Q(s_{t+1}, a_{t+1})$ （查表得到）。
  3.  计算 **TD Error**: $\delta_t = Q(s_t, a_t) - y_t$。
  4.  **更新表格**: $Q(s_t, a_t) \leftarrow Q(s_t, a_t) - \alpha \cdot \delta_t$ （$\alpha$ 为学习率）。

**B. 神经网络形式 (Value Network)**

- **适用场景**：状态空间巨大（如围棋、图像输入），无法建表，需要用神经网络近似 Q 函数。
- **流程**：
  1.  把神经网络记为 $Q(s, a; w)$，参数为 $w$。
  2.  计算 **TD Target**: $y_t = r_t + \gamma \cdot Q(s_{t+1}, a_{t+1}; w)$。
      _注意：这里把 $y_t$ 视为常数（Label），不求梯度。_
  3.  **损失函数**:
      $$Loss = \frac{1}{2} (Q(s_t, a_t; w) - y_t)^2$$
      本质上就是最小化 TD Error 的平方。
  4.  **梯度下降**: 对 $w$ 求导更新，使网络预测值靠近 $y_t$。

---

### 4. 总结

- **SARSA 也是 Critic**：在 Actor-Critic 算法中，用来训练 Critic（价值网络）的方法本质上就是 SARSA。
- **区别预告**：SARSA 在计算 $y_t$ 时，使用的是**实际采取的下一个动作** $a_{t+1}$（On-policy）。下一节课将讲 Q-Learning，它在计算 $y_t$ 时使用的是**理论上最优的动作**（Off-policy）。
