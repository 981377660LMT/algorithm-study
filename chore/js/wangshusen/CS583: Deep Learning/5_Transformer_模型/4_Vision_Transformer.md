# Vision_Transformer

https://www.bilibili.com/video/BV18Q4y1o7NY

这份 PDF 课件是 Shusen Wang 教授关于 **Vision Transformer (ViT)** 的讲解。这篇论文（ICLR 2021: _An Image is Worth 16x16 Words_）是将 Transformer 架构从 NLP 领域无缝迁移到计算机视觉领域的开山之作。

以下是对该课件的**深入分析与解构**：

### 1. 核心理念：把图像当作文本 (The Paradigm Shift)

在 ViT 之前，CV 领域的绝对统治者是 CNN（如 ResNet）。CNN 利用卷积核提取局部特征，具有平移不变性等归纳偏置（Inductive Bias）。
ViT 的核心思想非常激进且简单：**完全抛弃卷积，直接把图像切碎，像处理一句话中的单词一样处理图像块。**

### 2. 系统架构解构

ViT 的处理流程可以被严格拆解为以下 5 个步骤：

#### **第一步：切块与展平 (Patchify & Flatten)**

- **操作**：给定一张 $H \times W \times C$ 的图像，将其切成 $N$ 个固定大小的块（Patches）。
- **参数**：假设 Patch 大小为 $16 \times 16$，那么原图就被切成了一个 $16 \times 16$ 的网格。
- **向量化**：每个 Patch 包含 $16 \times 16 \times 3$ (RGB) 个像素值。ViT 直接将其拉平（Flatten）为一个一维向量。
- _隐喻_：这相当于把图片变成了一篇由 N 个“单词”组成的文章。

#### **第二步：线性投影 (Linear Projection)**

- **操作**：每个 Patch 向量通过一个全连接层（Dense Layer），映射到一个固定维度 $D$ 的潜在空间。
- **权重共享**：所有 Patch 共用同一个映射矩阵（类似于 CNN 的卷积核共享，或者是 Word Embedding 矩阵）。
- _结果_：得到一系列 Patch Embeddings：$\mathbf{z}_1, \mathbf{z}_2, \dots, \mathbf{z}_N$。

#### **第三步：注入位置信息 (Positional Encoding)**

- **痛点**：Transformer 的 Self-Attention 机制是**排列不变**的。如果你打乱 Patch 的顺序，Transformer 算出来的结果是一模一样的。但对于图像来说，猫的头在哪、腿在哪至关重要。
- **解决**：为每个 Patch Embedding **加上**（Add）一个可学习的位置向量。
- _效果_：模型因此能区分“左上角”和“右下角”。课件 Slide 20-21 直观展示了打乱的 Patch 如何通过位置编码被逻辑还原。

#### **第四步：[CLS] Token (借鉴 BERT)**

- **设计**：ViT 并没有简单地对所有 Patch 的输出取平均，而是照搬了 BERT 的设计。
- **操作**：在 Patch 序列的最前面强行插入一个特殊的、可学习的向量，称为 `[CLS]`（Classification Token）。
  - 输入序列变成了：`[CLS], Patch_1, Patch_2, ..., Patch_N`。
- **目的**：在经过多层 Self-Attention 的“混合双打”后，这个 `[CLS]` 向量会聚合全图的信息，专门用于最终的分类决策。

#### **第五步：Transformer Encoder**

- **主体**：这就是标准的 Transformer Encoder（堆叠的 Multi-Head Self-Attention + MLP）。这里的架构与 NLP 中的完全一致，没有任何针对视觉的特殊改动。
- **输出**：只取 `[CLS]` 对应的最终输出向量 $\mathbf{c}$，送入 Softmax 分类器。

### 3. 关键结论：数据量的饥渴 (Scale is All You Need)

课件的最后部分（Slide 32-34）揭示了 ViT 最重要的特性——**对数据规模的依赖**。

- **小数据 (ImageNet-1k, 1.3M)**：ViT 表现 **不如** ResNet。
  - _原因_：CNN 自带“先验知识”（Inductive Bias），比如它知道像素之间有局部联系、物体平移后还是同一个物体。ViT 不知道这些，它需要从头学习。在数据不够时，ViT 学不到这些几何特性。
- **中数据 (ImageNet-21k, 14M)**：ViT 与 ResNet 打平。
- **超大数据 (JFT-300M, 3 亿张图)**：ViT **超越** ResNet。
  - _原因_：当数据量大到一定程度，ViT 强大的建模能力（Global Context）开始显现优势，而不再受限于缺乏先验知识。

### 总结

ViT 的胜利证明了**通用架构 (General Purpose Architectures)** 的潜力。它不再依赖人工设计的特征提取器（卷积核），而是依赖**大规模数据预训练**让模型自己学会“如何看图”。

- **输入**：Image Patches Sequence
- **核心**：Self-Attention (计算 Patch 之间的全局关系)
- **代价**：需要极大的预训练数据集才能“涌现”出高性能。
