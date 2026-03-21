# CS336 Lecture 3: LLM 架构与超参数

> "Everything You Didn't Want to Know About LM Architecture and Training"
>
> 核心方法论：我们无法自己训练所有 Transformer，所以要**从他人的经验中学习**——通过分析过去几年 19+ 个开放模型的架构选择，发现**趋同演化**的规律。

---

## 一、架构变体 (Architecture Variations)

### 1.1 Pre-Norm vs Post-Norm

这是**几乎所有人都达成共识**的第一个选择，也是最早被确立的"铁律"。

**原始 Transformer (Post-Norm)**：

```
x → Multi-Head Attention → Add(residual) → LayerNorm → FFN → Add(residual) → LayerNorm
```

LayerNorm 放在残差连接**之后**（即在残差流内部）。

**现代做法 (Pre-Norm)**：

```
x → LayerNorm → Multi-Head Attention → Add(residual) → LayerNorm → FFN → Add(residual)
```

LayerNorm 放在每个子模块**之前**（在残差流外部）。

**为什么 Pre-Norm 更好？**

1. **梯度传播**：残差连接提供了从网络顶部到底部的恒等映射(identity connection)，梯度可以无损传播。如果在残差流中间插入 LayerNorm，会破坏这条直通路径，导致**梯度衰减或爆炸**（Post-Norm 的梯度 norm 远高于 Pre-Norm，表现为图中的橙色飙升曲线）。
2. **训练稳定性**：Post-Norm 需要精心设计的 learning rate warmup 才能稳定训练；Pre-Norm 可以直接训练，loss spike 更少、gradient norm 更低更平稳。
3. **实验证据**：Xiong et al. 2020, Salazar & Yuan 等多篇论文在 MT 和 BERT 场景下反复验证了这一点。

> 唯一的例外：OPT-350M 用了 Post-Norm（大概是失误）。

**新发展——Double Norm（2024-2025）**：

在 Pre-Norm 基础上，**在子模块之后也加一个 LayerNorm**：

- **Grok, Gemma 2**：attention 和 FFN 前后都加 LayerNorm
- **OLMo-2**：仅在 FFN 后加额外 LayerNorm

动机是进一步增强训练稳定性，在更大规模模型上被验证有效。

---

### 1.2 RMSNorm vs LayerNorm

**LayerNorm**：

$$\text{LayerNorm}(x) = \gamma \cdot \frac{x - \mu}{\sqrt{\sigma^2 + \epsilon}} + \beta$$

其中 $\mu = \frac{1}{d}\sum x_i$，$\sigma^2 = \frac{1}{d}\sum(x_i - \mu)^2$，$\gamma$ 和 $\beta$ 是可学习参数。

**RMSNorm**：

$$\text{RMSNorm}(x) = \gamma \cdot \frac{x}{\sqrt{\frac{1}{d}\sum x_i^2 + \epsilon}}$$

去掉了**减均值**和**加偏置** $\beta$ 两个操作。

**为什么几乎所有现代模型都用 RMSNorm？**

关键洞察不在于 FLOPs，而在于**内存移动(memory movement)**：

| 操作类型 | 占总 FLOPs 比例 | 占总运行时间比例 |
|---------|---------------|---------------|
| Tensor contractions (矩阵乘) | 99.8% | ~75% |
| Normalization (LayerNorm 等) | 0.17% | **~25%** |

> 来源：Ivanov et al. 2023 ("Memory Movement is All You Need")

Normalization 操作虽然 FLOPs 极少，但因为**内存移动开销大**，占了约 25% 的运行时间。RMSNorm 省掉了减均值和加偏置的内存访问，实测更快：

- Vanilla Transformer: 3.5 steps/sec
- RMSNorm 版本: 3.68 steps/sec（Narang et al. 2020）
- 且最终 loss 更低或持平

使用 RMSNorm 的代表模型：LLaMA 全系列, PaLM, Chinchilla, T5, Qwen, DeepSeek 等。
例外：Cohere (Command-A, R+) 仍用 LayerNorm。

---

### 1.3 去掉偏置项 (Dropping Bias Terms)

原始 FFN：$\text{FFN}(x) = \text{ReLU}(xW_1 + b_1)W_2 + b_2$

现代做法：$\text{FFN}(x) = \text{ReLU}(xW_1)W_2$

**原因**：
1. **性能持平**：bias 对模型质量影响极小，矩阵乘本身就够了
2. **训练稳定性**：经验表明去掉 bias 可以**减少训练不稳定性**，这是更重要的原因
3. **系统效率**：少了参数需要从内存加载

> 这与 RMSNorm 去掉 $\beta$ 是一脉相承的思路：能省则省，纯矩阵乘就够了。

---

### 1.4 激活函数与门控线性单元 (Activations & GLU)

#### 激活函数 Zoo

| 激活函数 | 公式 | 主要用户 |
|---------|------|---------|
| ReLU | $\max(0, x)$ | 原始 Transformer |
| GeLU | $x \cdot \Phi(x)$ ($\Phi$ 是高斯 CDF) | GPT-1/2/3, GPT-J |
| Swish | $x \cdot \sigma(x)$ ($\sigma$ 是 sigmoid) | PaLM, LLaMA 系列 |

GeLU 和 Swish 在零点附近不是单调的（有一小段负值凸起），但实际训练中这不是问题——因为优化器用了高学习率+动量，激活值不会被困在零点。

#### 门控线性单元 (GLU) — 核心创新

普通 MLP：

$$\text{FFN}(x) = \sigma(xW_1)W_2$$

GLU 版本（以 SwiGLU 为例）：

$$\text{FFN}(x) = [\text{Swish}(xW_1) \odot (xV)] W_2$$

关键变化：新增了一个**门控矩阵** $V$，$xV$ 产生的向量与 $\text{Swish}(xW_1)$ **逐元素相乘(element-wise)**，形成门控机制。

各种 GLU 变体只是非线性函数不同：
- **ReGLU**：用 ReLU
- **GEGLU**：用 GeLU（T5 V1.1, Gemma 2/3 使用）
- **SwiGLU**：用 Swish（**最流行**：LLaMA 全系, PaLM, OLMo, Qwen, DeepSeek）

**实验证据**：

- Shazeer (原始 GLU 论文)：GLU 变体在 KoLA, SST2 等任务上一致优于非门控版本，且差异显著
- Narang et al. 2020：在 T5 框架下，GLU 变体一致取得更低 loss（加粗行全是 GLU）

**参数匹配**：GLU 多了矩阵 $V$，为保持总参数量不变，通常将隐藏层维度缩小为非门控版本的 $\frac{2}{3}$。

> GLU 不是**必需的**（GPT-3 没用也很强，Memotron 340B 用 squared ReLU，Falcon 211B 用 ReLU），但它是**一致有增益的 (consistent gains)**，所以成了默认选择。

---

### 1.5 串行 vs 并行层 (Serial vs Parallel Layers)

**串行（标准做法）**：

```
x → Attention(x) → Add(residual) → FFN(result) → Add(residual)
```

Attention 的输出作为 MLP 的输入，是**顺序**执行的。

**并行**：

```
x → [Attention(x) + FFN(x)] → Add(residual)
```

Attention 和 FFN **同时计算**，结果加在一起再加回残差流。

- 由 GPT-J 首创，PaLM 在大规模上验证
- 优势：可以**融合 LayerNorm 和矩阵乘**，获得系统并行效率增益
- 劣势：串行更有表达力（组合两个计算 vs 仅相加）
- 近一年大多数模型回归串行；例外：Cohere Command-A/R+, Falcon Q111B

---

### 1.6 位置编码 (Position Embeddings) — 已收敛至 RoPE

#### 历史演变

| 方法 | 代表模型 | 时期 |
|-----|---------|------|
| Sinusoidal | 原始 Transformer | 2017 |
| Absolute (学习的) | GPT-1/2/3, OPT | 2018-2022 |
| Relative (加法) | T5, Gopher | 2019-2021 |
| ALiBi | BLOOM | 2022 (短暂流行) |
| **RoPE** | **几乎所有 2023+ 模型** | **2021至今** |

#### RoPE 的核心思想

**目标**：设计一个嵌入函数 $f(x, i)$，使得：

$$\langle f(x_m, m), f(x_n, n) \rangle = g(x_m, x_n, m-n)$$

即两个 token 嵌入的内积只依赖**相对位置** $m-n$，不依赖绝对位置。

**关键洞察**：内积对旋转不变 → 用**旋转**来编码位置。

**2D 直觉**：
- "we" 在位置 0 → 不旋转
- "know" 在位置 1 → 旋转 1 个单位角度
- 在 "of course we know" 中，"we" 在位置 2 旋转 2 次，"know" 在位置 3 旋转 3 次
- 两者的**相对角度不变** → 内积相同 ✓

**高维推广**：将 $D$ 维向量切成 $D/2$ 个**二维子空间**，每个子空间有不同的旋转频率 $\theta_i$：
- 高频 $\theta$：捕捉近距离相对位置信息
- 低频 $\theta$：捕捉远距离相对位置信息
- （类比 sinusoidal 编码的多频率设计）

**数学形式**：对每对维度 $(2i, 2i+1)$ 应用旋转矩阵：

$$R(\theta_i \cdot m) = \begin{pmatrix} \cos(m\theta_i) & -\sin(m\theta_i) \\ \sin(m\theta_i) & \cos(m\theta_i) \end{pmatrix}$$

**与传统位置编码的关键区别**：RoPE **不在 embedding 层添加**，而是在**每一层的 attention 计算时**介入——对 Q 和 K 做旋转，V 不做旋转。

```python
# LLaMA 风格的 RoPE 实现伪代码
Q, K, V = linear(x)              # 正常的 QKV 投影
cos, sin = compute_rope_angles()   # 基于位置计算旋转角度
Q_rot = apply_rotation(Q, cos, sin)  # 旋转 Q
K_rot = apply_rotation(K, cos, sin)  # 旋转 K
attn = softmax(Q_rot @ K_rot.T / sqrt(d)) @ V  # 用旋转后的 Q,K 做 attention
```

> $\theta_i$ 是**固定的**（不可学习），按照预定义的频率表确定。旋转本质上就是一个固定矩阵乘，不引入额外训练问题。

---

## 二、超参数选择 (Hyperparameters)

### 2.1 FFN 隐藏层维度 $d_{ff}$ vs 模型维度 $d_{model}$

**共识规则**：

| 类型 | 比例 | 说明 |
|-----|------|------|
| 非门控 MLP (ReLU/GeLU) | $d_{ff} = 4 \times d_{model}$ | 自原始 Transformer 以来的约定 |
| 门控 MLP (SwiGLU 等) | $d_{ff} = \frac{8}{3} d_{model} \approx 2.67 \times d_{model}$ | 因为多了门控矩阵 $V$，为参数匹配而缩小 |

**实验支撑**（Kaplan et al. Scaling Laws 论文）：

$d_{ff}/d_{model}$ 比值在 1~10 之间都接近最优，4 正好在这个平坦区域内。不是自然法则，但被广泛验证为合理默认值。

**遵循此规则的模型**：LLaMA 1, Qwen, DeepSeek, T5 V1.1 等。

**例外**：
- **T5 (11B)**：$d_{model}=1024$, $d_{ff}=65536$，比值高达 **64x**！理由是更宽的矩阵乘可以获得更好的并行效率。但后继版本 T5 V1.1 回归到标准的 2.5x 比值并取得了更好结果。
- **Gemma 2**：用了 8x 比值。
- **PaLM, Mistral, LLaMA**：比标准 2.67 略大，但量级一致。

---

### 2.2 注意力头维度

**共识**：$d_{model} = n_{heads} \times d_{head}$（比值为 1）

增加头数时**分割**维度，而不是扩大总维度。GPT-3, T5, LaMDA, PaLM, LLaMA-2 都遵循这一规则。

理论上头数太多 → 每头维度太低 → 注意力低秩瓶颈（Bhojanapalli et al. 2020），但实践中比值为 1 的模型表现良好。

---

### 2.3 宽深比 (Aspect Ratio): $d_{model}$ / $n_{layers}$

**共识甜点**：约 **128** 维/层

- Kaplan et al. 实验：在 50M / 274M / 1.5B 三个规模上，$d_{model}/n_{layers} \approx 100$ 时 loss 最低
- **跨规模稳定**：最优宽深比在多个数量级上基本不变（利好消息）
- GPT-3 和 LLaMA 系列模型都大致遵循此比例

**系统考量**：
- 更深 → 更适合 Pipeline Parallelism（层间切分到不同设备）
- 更宽 → 更适合 Tensor Parallelism（矩阵切分到不同设备）
- 网络带宽约束会反过来影响宽深选择

**下游任务注意**（Tay et al., Google）：
- 在 **pre-training loss** 上，宽深比影响不大，总参数量才是关键
- 在 **downstream accuracy**（如 fine-tuned SuperGLUE）上，更深模型可能在同 FLOPs 下更好

---

### 2.4 词表大小 (Vocabulary Size)

| 时期/类型 | 词表大小范围 | 代表 |
|----------|-----------|------|
| 早期/单语 | 30K-50K | GPT-2, 早期 LLaMA |
| 现代/多语/生产 | 100K-256K | GPT-4 (~100K), LLaMA 3 (128K), Cohere Command-A (256K) |

趋势：**词表越来越大**，因为：
1. 多语言支持：更大词表让低资源语言用更少 token 表示 → 推理成本更低
2. 模型变大后有能力有效利用更多 vocab 元素
3. Emoji、代码等多模态需求

---

### 2.5 正则化：Dropout 与 Weight Decay

#### Dropout — 已过时

Pre-training 通常**只做一个 epoch**（数据太多用不完），不存在过拟合问题 → dropout 没有理论必要性 → **现代模型基本不用 dropout**。

#### Weight Decay — 仍然广泛使用，但原因出人意料

**直觉上的矛盾**：既然不会过拟合，为什么还需要 weight decay（一种正则化技术）？

**答案：Weight decay 不是为了防止过拟合，而是为了获得更好的训练 loss。**

实验观察：
1. 不同 weight decay 下，train-val loss gap 几乎不变 → 不是在控制过拟合
2. 高 weight decay 的模型在**高学习率阶段训练较慢**
3. 但当学习率衰减（如 cosine decay 末期）时，高 weight decay 的模型会**急速下降**，最终达到更好的训练 loss

> 这是 weight decay 与学习率调度之间的**复杂交互作用**：weight decay 在训练末端产生某种隐式加速效果，使模型在 learning rate 降低时能更高效地优化。

---

## 三、训练稳定性技巧 (Stability Tricks) — 近一年最重要的新发展

### 3.1 问题根源：Softmax

Transformer 中有两个 softmax：
1. **输出层 softmax**：将 logits 转换为 token 概率
2. **注意力层 softmax**：计算注意力权重

Softmax 的数值问题：
- 涉及指数运算 → 容易上溢
- 涉及除法 → 可能除以零
- 是训练不稳定（gradient norm 飙升/loss spike）的主要来源

### 3.2 Z-Loss — 稳定输出 Softmax

**出发点**：输出 softmax $p(y) = \frac{e^{u(y)}}{Z}$，其中 $Z = \sum_y e^{u(y)}$。

如果 $Z \to 1$（即 $\log Z \to 0$），那么 softmax 退化为简单的指数运算，数值稳定。

**方法**：在损失函数中加辅助项：

$$\mathcal{L}_{total} = \mathcal{L}_{CE} + \alpha \cdot (\log Z)^2$$

其中 $\alpha$ 通常很小（PaLM 用 $10^{-4}$）。

**效果**：鼓励 $\log Z \approx 0$，使 softmax 归一化器保持良好行为。

> PaLM (2022) 首创 → BiTran2, DCLM, OLMo-2 等跟进。

### 3.3 QK-Norm — 稳定注意力 Softmax

**方法**：在计算 $QK^T$ 之前，对 Q 和 K 分别做 LayerNorm：

```
Q = LayerNorm(xW_Q)
K = LayerNorm(xW_K)
attn_scores = Q @ K.T / sqrt(d)
```

**效果**：控制 softmax 的**输入**大小有界 → 避免 attention logits 爆炸。

**来源**：最初是视觉/多模态社区的创新（Dehghani 2023, ViT-22B），后被文本 LLM 采用：
- Gemma 2, DCLM, OLMo-2, Chameleon 等

**对比实验**（NVIDIA）：
- Baseline 困惑度: 11.19
- + Soft capping: **更差**
- + QK-Norm: **更好**（因为可以使用更激进的学习率）

### 3.4 Logit Soft Capping — 软裁剪注意力 Logits

$$\text{logits}' = \text{softcap} \cdot \tanh\left(\frac{\text{logits}}{\text{softcap}}\right)$$

当 logits 远超 softcap 时，$\tanh$ 将其裁剪到 $\pm 1$ → 输出上界为 $\pm \text{softcap}$。

Gemma 2, OLMo-2 使用。但不如 QK-Norm 流行，且 NVIDIA 实验显示效果不如 QK-Norm。

### 3.5 LayerNorm 的惊人有效性

**一个贯穿始终的主题**：LayerNorm 在稳定训练方面**极其有效**——
- Pre-Norm: 残差块前加 LN ✓
- Double Norm: 残差块前后都加 LN ✓
- QK-Norm: Q/K 上也加 LN ✓

每次加 LayerNorm 都提升了稳定性，且几乎不影响模型性能。

---

## 四、注意力头变体 (Attention Variants)

### 4.1 推理瓶颈：KV Cache 的内存问题

**训练时**的算术强度 (arithmetic intensity)：

$$AI_{train} = \frac{1}{\frac{1}{K} + \frac{1}{BN}}$$

$K$ = 头数, $B$ = batch size, $N$ = 序列长度。头数和 batch 都大 → AI 高 → GPU 利用率好。

**推理时**（自回归生成，每次只处理 1 个 token）：

需要维护 **KV Cache**：逐步积累过去所有 token 的 Key 和 Value，每生成一个新 token 就新增一行。

推理算术强度：

$$AI_{infer} = \frac{1}{\frac{N}{D} + \frac{1}{B}}$$

问题：$N/D$ 项随序列长度增加而变大 → **AI 快速下降** → 推理变成内存瓶颈。

### 4.2 MQA (Multi-Query Attention)

**核心思想**：Query 保持多头，但 **K 和 V 只有一个头**（所有 Query 头共享同一组 K, V）。

效果：KV Cache 大小缩小为原来的 $1/n_{heads}$ → 内存访问大幅减少 → 算术强度改善。

### 4.3 GQA (Grouped Query Attention)

MQA 的折中版本：将 Query 头**分组**，每组共享一组 K, V。

$$n_{kv\_heads} = n_{query\_heads} / G$$

$G$ 是组大小。可以在 MHA（$G=1$）和 MQA（$G = n_{heads}$）之间灵活调节。

部分实验表明 GQA 不损失性能，而 MQA 可能损失。大型 LLaMA 模型使用 GQA。

### 4.4 长上下文：混合注意力模式

**稀疏注意力**（OpenAI, 2019）：不做全序列 attention，而是设计结构化稀疏模式（局部窗口 + 对角条纹等），GPT-3 使用。

**滑动窗口注意力 (Sliding Window Attention)**：每层只关注当前位置周围的固定窗口，有效感受野 = 窗口大小 × 层数。

**最新做法——混合块结构**（LLaMA 4, Gemma, Cohere Command-A）：

```
[Block pattern repeating every 4 layers]
Layer 1: Full self-attention, NO position embedding (no RoPE)
Layer 2: Sliding window attention + RoPE
Layer 3: Sliding window attention + RoPE
Layer 4: Sliding window attention + RoPE
```

**设计精妙之处**：
1. **全注意力层不用 RoPE** → 没有位置编码的限制 → 可以无限外推长距离依赖
2. **滑动窗口层用 RoPE** → 精确编码局部位置信息
3. **全注意力只占 1/4 层** → 控制计算/内存开销
4. LLaMA 4 声称支持 **1000 万 token** 上下文窗口

---

## 五、总结：现代 LLM 架构的默认配方

| 组件 | 共识选择 | 共识程度 |
|-----|---------|---------|
| Norm 位置 | Pre-Norm（残差流外） | ⭐⭐⭐⭐⭐ 铁律 |
| Norm 类型 | RMSNorm | ⭐⭐⭐⭐ 几乎所有 |
| 偏置项 | 去掉 | ⭐⭐⭐⭐ |
| 激活/FFN | SwiGLU | ⭐⭐⭐⭐ |
| 位置编码 | RoPE | ⭐⭐⭐⭐⭐ 完全收敛 |
| 层连接 | 串行 | ⭐⭐⭐ 主流但有例外 |
| $d_{ff}/d_{model}$ | 4x (非门控) / 2.67x (门控) | ⭐⭐⭐⭐ |
| $d_{model}/n_{layers}$ | ~128 | ⭐⭐⭐⭐ |
| 头维度分配 | $d_{model} = n_{heads} \times d_{head}$ | ⭐⭐⭐⭐ |
| 词表大小 | 100K-256K | ⭐⭐⭐ 趋势上升 |
| Dropout | 不用 | ⭐⭐⭐⭐ |
| Weight decay | 用（为了优化动力学） | ⭐⭐⭐⭐ |
| 稳定性 | Z-loss + QK-Norm | ⭐⭐⭐ 新趋势 |
| 注意力 | GQA (大模型) | ⭐⭐⭐⭐ |

### 可迁移的通用教训

1. **保持干净的残差恒等连接** — 在多种架构中反复被验证
2. **LayerNorm 是万能稳定剂** — 哪里不稳定就往哪里加
3. **门控机制一致有增益** — 不仅在 FFN，在许多架构中都是如此
4. **不要只看 FLOPs，要看内存移动** — 这是系统设计的关键视角
5. **简化能赢** — RMSNorm, 去 bias, 都是"去掉不必要的东西"
