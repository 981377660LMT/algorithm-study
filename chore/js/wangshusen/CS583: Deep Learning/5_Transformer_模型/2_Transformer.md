# Transformer

这份 PDF 是 Shusen Wang 教授关于 **Transformer 模型** 的第二部分，主题为 **“From Shallow to Deep” (从浅层到深层)**。

如果说第一部分（1/2）是造出了“积木”（Attention Layer），那么这一部分（2/2）就是教我们如何把积木搭成一座宏伟的摩天大楼。它完整地展示了 Transformer 的整体架构。

以下是对该课件的**深入分析与解构**：

### 1. 核心进化：Multi-Head Attention (多头注意力)

课件首先将简单的 Self-Attention 升级为 Multi-Head Self-Attention。

- **动机**：人看东西时，可能同时关注不同方面（比如既看颜色，也看形状）。Single-Head（单头）注意力只能捕捉一种类型的相关性。
- **实现**：
  - 不做一次 Attention，而是并行做 $l$ 次。
  - 如果你有 $l$ 个头，你就有 $l$ 组不同的参矩阵 $(\mathbf{W}_Q, \mathbf{W}_K, \mathbf{W}_V)$。
  - **Concatenation**：将 $l$ 个头的输出拼接起来。如果每个头输出 $d$ 维向量，拼接后就得到 $l \times d$ 维。
  - **Linear Projection**：虽然课件中简化了这一点，但在拼接后通常还会接一个线性层（Dense）把维度融合。

### 2. 积木单元：Encoder Block 和 Decoder Block

这是 Transformer 架构中最具标志性的“模块化”设计。

#### **Encoder Block (编码器模块)**

一个标准的 Encoder Block 由两部分组成（课件中为了教学简化，略去了 Residual Connection 和 Layer Normalization，专注于核心逻辑）：

1.  **Multi-Head Self-Attention**：让输入的序列“自己理解自己”，捕捉词与词之间的依赖。
2.  **Dense Layer (Feed Forward Network)**：对每个位置的向量进行独立的非线性变换。公式为 $\mathbf{u} = \text{ReLU}(\mathbf{W}\mathbf{c})$。
    - _输入_：$512 \times m$ 矩阵（$m$ 是序列长度）。
    - _输出_：$512 \times m$ 矩阵。
    - **特点**：输入和输出形状完全一致，这意味着我们可以无限**堆叠 (Stack)** 这些模块。

#### **Decoder Block (解码器模块)**

Decoder 稍微复杂一点，因为它不仅要生成序列，还要参考 Encoder 的信息。一个 Block 包含三层：

1.  **Multi-Head Self-Attention**：Decoder 自己看自己（Masked，防止偷看后面，虽然课件此处未详细展开 Mask 细节，但这是生成式任务必须的）。
2.  **Multi-Head Attention (Cross-Attention)**：这是关键桥梁。
    - **Query** 来自 Decoder 上一层的输出。
    - **Key & Value** 来自 **Encoder 的最终输出**（$\mathbf{U}$）。
    - _作用_：生成翻译时，回头看原文（Encoder）里哪里重要。
3.  **Dense Layer**：最后的处理。

### 3. 宏伟蓝图：整体架构 (Put Everything Together)

课件通过一系列图示展示了完整的 Transformer 架构：

- **Encoder 堆叠**：将 6 个 Encoder Block 串联。输入 $\mathbf{X}$ 经过 6 次“Self-Attention + Dense”的洗礼，通过不断抽象，最终输出矩阵 $\mathbf{U}$。$\mathbf{U}$ 包含了原文的深度语义表示。
- **Decoder 堆叠**：将 6 个 Decoder Block 串联。
  - 每个 Decoder Block 都有一根线连到 Encoder 的输出 $\mathbf{U}$ 上（用于 Cross-Attention）。
  - 输入是已经生成的单词序列（$\mathbf{X}'$）。
- **生成过程**：
  1.  Encoder 一次性处理完整个英文句子 $\mathbf{X}$，得到 $\mathbf{U}$。
  2.  Decoder 先接收 `Start Sign`。
  3.  结合 $\mathbf{U}$ 和 `Start`，Decoder 算出第一个词的概率分布 $\mathbf{y}_1$。
  4.  采样出第一个词（比如 "Ich"），把它加到输入里。
  5.  Decoder 接收 `Start, Ich`，算出第二个词的概率...
  6.  直到生成 `Stop Sign`。

### 4. Transformer vs. RNN

课件在 Slide 34 做了一个直观的对比：

- **RNN/LSTM**：水平向右的时间轴依赖。必须按顺序算，$\mathbf{h}_t$ 依赖 $\mathbf{h}_{t-1}$。
- **Transformer**：垂直向上的深度依赖。
  - 没有水平箭头！
  - 层与层之间是并行的，序列中的每个词在同一层是独立计算的（通过 Attention 交互）。
  - 这使得 Transformer 极度适合 GPU 并行加速。

### 总结

这份课件揭示了 Deep Learning 时代模型设计的两个核心哲学：

1.  **模块化与深度**：与其设计一个复杂的巨型网络，不如设计一个精巧的 Block，然后把这个 Block 堆叠 6 次、12 次甚至 100 次（如 GPT-4）。
2.  **注意力就是一切**：通过 Multi-Head 机制，模型可以在不同的语义子空间（Subspaces）里并行捕捉信息，彻底抛弃了循环连接，实现了并行化与长距离依赖的完美统一。
