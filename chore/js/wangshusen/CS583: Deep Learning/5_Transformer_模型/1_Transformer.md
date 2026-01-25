# Transformer

这份 PDF 课件是 Shusen Wang 教授关于 **Transformer 模型** 的第一部分，主题为 **“Attention without RNN” (没有 RNN 的注意力机制)**。

这一节课的主要逻辑是“去依附”：它展示了如何将 Attention 机制从 RNN 架构中剥离出来，使其成为一个独立的组件（Layer），并最终演化出 Self-Attention（自注意力）。

以下是对该课件的**深入分析与解构**：

### 1. 核心思想：从“插件”到“主角”

在 Seq2Seq RNN 模型中，Attention 只是一个辅助“插件”，用于解决长序列遗忘问题。Transformer 的核心突破在于——**如果我们完全扔掉 RNN，只用 Attention 会怎样？**

为了做到这一点，必须重新定义 Attention 的输入和输出，使其不依赖于 RNN 的隐藏状态（Hidden States）。

### 2. 概念解构一：Q、K、V 的提炼

课件首先回顾了 RNN 中的 Attention，意在提炼出三个核心概念（类似于数据库查询的隐喻）：

- **Query ($\mathbf{q}$)**：你现在的关注点是什么？（来自 Decoder 的状态）。
- **Key ($\mathbf{k}$)**：信息库里的索引标签。（来自 Encoder 的隐藏状态）。
- **Value ($\mathbf{v}$)**：索引对应的内容本体。（同样来自 Encoder 的隐藏状态）。

在 RNN 中：
$$ \mathbf{q} = \mathbf{W}\_Q \cdot \mathbf{s} \quad (\text{Decoder state}) $$
$$ \mathbf{k} = \mathbf{W}\_K \cdot \mathbf{h}, \quad \mathbf{v} = \mathbf{W}\_V \cdot \mathbf{h} \quad (\text{Encoder state}) $$
Attention 的本质就是：**拿 Q 去匹配所有的 K，算出权重，然后对 V 进行加权求和。**

### 3. 概念解构二：通用的 Attention Layer

为了摆脱 RNN，课件定义了一个通用的 **Attention Layer**（在 Transformer 中通常指 **Cross-Attention** 或 Encoder-Decoder Attention）：

- **输入来源**：不再是 RNN 的状态 $\mathbf{h}$ 或 $\mathbf{s}$，而是直接基于输入向量序列（比如 Word Embeddings）。
  - **Source Inputs ($\mathbf{X}$)**：比如英语句子的向量序列 $\mathbf{x}_1, \dots, \mathbf{x}_m$。
  - **Target Inputs ($\mathbf{X}'$)**：比如德语句子的向量序列 $\mathbf{x}'_1, \dots, \mathbf{x}'_t$。
- **生成 QKV**：
  - **Key & Value** 来自 source $\mathbf{X}$：$\mathbf{k}_j = \mathbf{W}_K \mathbf{x}_j$, $\mathbf{v}_j = \mathbf{W}_V \mathbf{x}_j$。
  - **Query** 来自 target $\mathbf{X}'$：$\mathbf{q}_i = \mathbf{W}_Q \mathbf{x}'_i$。
- **计算逻辑**：
  1.  每个 Decoder 向量算出一个 Query。
  2.  这个 Query 去和 Encoder 所有向量生成的 Keys 算相似度（点积）。
  3.  Softmax 归一化得到权重 $\boldsymbol{\alpha}$。
  4.  加权求和 Encoder 的 Values 得到 Context Vector $\mathbf{c}_i$。

**结论**：这个层实现了**“用目标（Target）去源（Source）中查找相关信息”**的功能，这完全不再需要时序上的 RNN 连接。

### 4. 概念解构三：Self-Attention Layer (自注意力)

这是 Transformer 的灵魂。如果 Source 和 Target 是**同一个序列**，会发生什么？
这就是 **Self-Attention**：$\text{Attn}(\mathbf{X}, \mathbf{X})$。

- **输入**：只有一个序列 $\mathbf{X} = [\mathbf{x}_1, \dots, \mathbf{x}_m]$。
- **生成 QKV**：三者都源自同一个 $\mathbf{x}$。
  - $\mathbf{q}_i = \mathbf{W}_Q \mathbf{x}_i$
  - $\mathbf{k}_i = \mathbf{W}_K \mathbf{x}_i$
  - $\mathbf{v}_i = \mathbf{W}_V \mathbf{x}_i$
- **物理意义**：
  - 对于序列中的每一个词（比如第 $i$ 个词），它都充当一次“查询者”（Query）。
  - 它去询问序列中的所有其他词（Keys）：_“你们谁跟我有关系？”_
  - 它根据关系的强弱，聚合所有词的信息（Values）来更新自己的表示。
- **结果**：输出的 $\mathbf{c}_i$ 不再只是第 $i$ 个词原本的含义，而是**融合了整个句子上下文信息的、关于第 $i$ 个位置的深层表达**。

### 5. 对比总结 (Slide 45-46)

课件最后给出了极简的公式对比，一针见血：

- **Standard Attention (Cross-Attention)**:
  $$ \mathbf{C} = \text{Attn}(\mathbf{X}_{\text{source}}, \mathbf{X}_{\text{target}}) $$
  - _用途_：机器翻译中 Decoder 看 Encoder。
- **Self-Attention**:
  $$ \mathbf{C} = \text{Attn}(\mathbf{X}, \mathbf{X}) $$
  - _用途_：Encoder 内部或 Decoder 内部，自己理解自己，捕捉内部元素间的距离依赖。

### 深度解读：为什么要这么做？

虽然这份课件没有深入讲 Transformer 的所有细节（这是 1/2 部分），但它清晰地展示了 Transformer 底层的**原子操作**演变：

1.  **并行计算**：RNN 必须等 $t-1$ 算完才能算 $t$。而 Self-Attention 中，$\mathbf{Q}, \mathbf{K}, \mathbf{V}$ 的计算以及最后的加权求和，对于序列中所有位置 $i$ 都是独立的，可以**矩阵运算并行化**。
2.  **全局视野**：RNN 需要一步步传递信息，甚至还要依靠双向 RNN 才能看全。Self-Attention 每个词直接与所有词“握手”，距离永远是 1，极大地增强了对上下文的捕捉能力。

这节课是理解 Transformer 架构之前的“数学铺垫”，解释了 Attention 作为一个独立 Layer 的合法性和在不同场景下的变体。
