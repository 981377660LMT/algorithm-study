Deep & Cross Networks (DCN) 译作“深度交叉网络”，可以用于召回双塔模型、粗排三塔模型、精排模型。DCN 由一个深度网络和一个交叉网络组成，交叉网络的基本组成单元是交叉层 (Cross Layer)。这节课最重点的部分就是交叉层。

---

这节课的内容从简单的 FM 模型升级到了更复杂的 **DCN V2 (Deep Cross Network Version 2)**。

如果说 FM 解决了“二阶特征交叉”的问题，那么 DCN 就是为了解决“**高阶 (High-Order) 特征交叉**”的问题。它的设计初衷是在保留深度神经网络（DNN）强大学习能力的同时，显式地、高效地构造多阶特征交叉。

以下是深度逻辑拆解与总结：

### 1. 核心定位：Embed & Cross

- **DNN (Deep Neural Network)**: 擅长隐式学习特征之间的非线性关系，但对“特征交叉（如 A x B x C）”这种显式乘法关系的捕捉效率并不高。
- **DCN (Deep Cross Network)**: 是一种**混合架构**。它由两部分组成：
  1.  **Cross Network (交叉网络)**: 专门负责显式地构造特征交叉。
  2.  **Deep Network (深度网络)**: 就是普通的 MLP (多层感知机)，负责隐式学习复杂非线性关系。

### 2. 交叉层 (Cross Layer) 的魔法

这是 DCN 的灵魂。一个 Cross Layer 的公式如下（基于 V2 版本）：

$$ \mathbf{x}\_{l+1} = \mathbf{x}\_0 \odot (\mathbf{W}\_l \mathbf{x}\_l + \mathbf{b}\_l) + \mathbf{x}\_l $$

让我们拆解这个公式：

- $\mathbf{x}_0$: **最底层的输入向量** (原始特征 embedding 拼接后的向量)。**注意：每一层都要把这个 $x_0$ 拿进来乘。**
- $\mathbf{x}_l$: 第 $l$ 层的输入向量。
- $\mathbf{W}_l, \mathbf{b}_l$: 第 $l$ 层的权重和偏置。
- $\mathbf{W}_l \mathbf{x}_l + \mathbf{b}_l$: 这是一个简单的线性变换。
- $\odot$: **Hadamard Product (逐元素相乘)**。
- $+ \mathbf{x}_l$: **Skip Connection (残差连接)**，保留上一层的信息，防止梯度消失。

**本质原理**：

- 第 1 层：$\mathbf{x}_1$ 包含了 $\mathbf{x}_0 \times \mathbf{x}_0$ (二阶交叉)。
- 第 2 层：$\mathbf{x}_2$ 包含了 $\mathbf{x}_0 \times \mathbf{x}_1$ $\to$ $\mathbf{x}_0 \times (\mathbf{x}_0 \times \mathbf{x}_0)$ (三阶交叉)。
- 第 $L$ 层：包含了 $\mathbf{x}_0$ 的 $L+1$ 阶交叉。
- 通过这种递归结构，DCN 可以用很少的参数显式地构造出高阶特征交叉。

### 3. DCN 的整体架构 (Parallel Structure)

DCN 通常采用**并行结构**将 Cross Network 和 Deep Network 结合：

1.  **输入层**：所有特征 Embeddings 拼接成一个长向量 $\mathbf{x}_0$。
2.  **分流**：
    - 一路进入 **Cross Network** (堆叠 $N$ 个 Cross Layer)。
    - 一路进入 **Deep Network** (堆叠 $M$ 个全连接层)。
3.  **融合**：将 Cross Network 的输出和 Deep Network 的输出 **Concat** (拼接) 起来。
4.  **输出层**：经过一个 Logits 层得到预测结果。

### 4. DCN 在推荐系统中的位置

DCN 不是一个独立的业务模型（像召回模型或排序模型），而是一个**通用的神经网络组件**（像 CNN 或 Transformer）。它可以替换掉任何现有模型中的“普通全连接层 (MLP)”。

- **在双塔模型中**：可以用 DCN 替换单纯的 MLP 来增强 User Tower 或 Item Tower 的表达能力。
- **在多目标排序 (Shared-Bottom) 中**：可以用 DCN 替换底层的 Shared-Bottom 网络。
- **在 MMoE 中**：可以用 DCN 作为每一个 Expert 网络。

### 5. 总结与评价

- **DCN V1 vs V2**：
  - V1 (2017): Cross Layer 公式稍微不同 ($\mathbf{x}_0 \mathbf{x}_l^T \mathbf{w}_l$)，计算的是一个标量乘法，表达能力受限（秩为 1）。
  - V2 (Mixture of LoRA / Matrix): 引入了矩阵权重 $\mathbf{W}_l$（低秩分解版），大大增强了表达能力。王老师课中讲的 Hadamard 积形式也是 V2 的一种变体或是为了便于理解的简化描述（实际 V2 主要是通过矩阵乘法实现更复杂的混合）。
- **优势**：
  - **显式高阶交叉**：比 FM 的二阶更强。
  - **参数高效**：Cross Network 的参数量远小于同深度的 MLP。
  - **即插即用**：可以无缝集成到现有的推荐架构中。

这节课的核心 take-away 是：**全连接层 (DNN) 不是万能的**，它在特征交叉这件事情上效率不高。通过引入专门设计的 Cross Network，我们能在参数量几乎不增加的情况下，显著提升模型的特征捕获能力。
