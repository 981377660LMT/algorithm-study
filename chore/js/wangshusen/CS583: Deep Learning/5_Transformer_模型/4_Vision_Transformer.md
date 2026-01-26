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

---

这节课由王树森老师讲解，介绍了 **Vision Transformer (ViT)**，这是 Transformer 架构在计算机视觉（CV）领域的开创性应用。ViT 的出现打破了 CNN 在图像分类领域的长期统治。

以下是逻辑清晰、深入且不遗漏的分析：

---

### 第一部分：背景与定位

#### 1. 什么是 Vision Transformer (ViT)?

- **发布时间**：2020 年 10 月（arXiv），2021 年正式发表（ICLR）。此论文名为《An Image is Worth 16x16 Words》。
- **地位**：它是目前图片分类（Image Classification）领域的最强模型之一，在大规模数据预训练下，性能超越了最先进的 CNN（如 ResNet）。
- **核心思想**：Transformer 原本是为 NLP（处理单词序列）设计的，ViT 并没有改变 Transformer 的架构，而是**设法把一张二维图片变成了一个一维的序列**，从而直接套用 Transformer。

---

### 第二部分：图片序列化 (Preprocessing)

如何把一张图片喂给只吃序列的 Transformer？

#### 1. 分块 (Patch Partitioning)

- **做法**：将整张图切分成大小相同的若干个小块（Patches）。
  - 例如：图片 $224 \times 224$，Patch 大小设为 $16 \times 16$。
  - 总共切得 $N = (224/16) \times (224/16) = 14 \times 14 = 196$ 个块。
- **Stride**：通常设置 Stride = Patch Size（无重叠切分）。如果 Stride < Patch Size，会有重叠，块数会变多，计算量变大。

#### 2. 线性投影 (Linear Projection of Flattened Patches)

- 把每个 $16 \times 16 \times 3$ (RGB) 的小块拉平成一个长向量。
- 通过一个**共享权重**的全连接层（Linear Layer），将这个长向量映射到 Transformer 需要的维度 $D$（例如 768）。
- 现在，这 196 个 $D$ 维向量，就完全等同于 NLP 里的一句包含 196 个单词的句子。

---

### 第三部分：Embeddings 处理

Transformer 的输入除了 Patch Embeddings，还需要两个关键的加持：

#### 1. Positional Encoding (位置编码)

- **必要性**：Transformer 的自注意力机制是**位置无关**的（Permutation Invariant）。如果你把拼图打乱，Transformer 看到的特征是一模一样的。但对图片来说，哪块在左上角，哪块在右下角非常重要。
- **做法**：给每个 Patch 向量加上一个可学习的**位置向量**。
- **影响**：如果不加位置编码，准确率会掉 3% 左右。

#### 2. Class Token (分类符 [CLS])

- **借鉴 BERT**：ViT 借鉴了 BERT 的设计，在输入的 Patch 序列最前面，硬塞了一个特殊的向量 **`[CLS]` Token ($z_0$)**。
- **作用**：这个 Token 不代表任何图片块，它的任务是在经过多层 Self-Attention 后，**汇总全图的信息**。
- **结果**：也就是说，不管你切了多少个块，最终主要拿第 0 个位置的输出向量 ($c_0$) 去做分类预测。

---

### 第四部分：模型架构 (Transformer Encoder)

ViT 的主体就是一个标准的 **Transformer Encoder**（跟上节课讲的一模一样）。

1.  **输入**：Sequence of Vectors ($N+1$ 个向量，包含 Patches 和 CLS Token)。
2.  **堆叠层**：
    - Multi-Head Self-Attention (MSA)：所有 Patch 互相“看”一遍，捕捉全局关联（比如“狗头”Patch 和“狗尾巴”Patch 的关系）。
    - MLP (Feed Forward Network)。
    - Layer Norm & Skip Connections (残差连接)。
3.  **输出**：同样是 $N+1$ 个向量。
4.  **分类头 (Classification Head)**：
    - 抛弃 $c_1$ 到 $c_N$（因为它们只代表局部信息）。
    - **只取 $c_0$**（即 `[CLS]` 位置的输出），输入一个 Softmax 分类器，输出类别概率（如 Dog: 0.4, Cat: 0.1 ...）。

---

### 第五部分：训练策略与实验结论 (Scale is All You Need)

ViT 并不是在任何情况下都比 CNN 强，它对**数据规模**有极高的渴望。

#### 1. 预训练与微调 (Pre-training & Fine-tuning)

- **流程**：
  1.  **Pre-train**：在一个超大数据集 A（如 JFT-300M，3 亿张图）上训练。
  2.  **Fine-tune**：迁移到任务数据集 B（如 ImageNet，130 万张图）上微调权重。
  3.  **Test**：在 B 的测试集上评估。

#### 2. 与 ResNet (CNN) 的对比

- **小数据 (ImageNet-1K, 130 万图)**：ViT **弱于** ResNet。
  - 原因：CNN 自带**归纳偏置 (Inductive Bias)**，天生懂得“平移不变性”和“局部相关性”，所以在小数据上学得快。ViT 这种“白纸一张”的模型还没学明白。
- **中数据 (ImageNet-21K, 1400 万图)**：ViT **持平** ResNet。
- **大数据 (JFT-300M, 3 亿图)**：ViT **显著超越** ResNet。
  - 结论：随着数据量无限增长，CNN 会因为其结构的局限性遭遇性能瓶颈（Saturate），而 ViT 的性能曲线仍在上升。

### 总结

1.  **ViT = Patches + Transformer Encoder**。
2.  **没有花哨的创新**：ViT 的核心贡献在于证明了 Transformer 不需要针对视觉做复杂的魔改（不像以前的 CNN+Attention 混合体），只要把图切了，纯 Transformer 就能称霸 CV 界。
3.  **大数据的胜利**：ViT 是典型的“大力出奇迹”。它牺牲了 CNN 的高效归纳偏置，换来了在大规模数据下更高的上限。这标志着 CV 和 NLP 的架构正在走向大一统。
