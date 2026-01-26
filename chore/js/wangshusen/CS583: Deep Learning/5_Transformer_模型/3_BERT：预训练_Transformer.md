Transformer 模型是目前机器翻译等 NLP 问题最好的解决办法，比 RNN 有大幅提高。Bidirectional Encoder Representations from Transformers (BERT) 是预训练 Transformer 最常用的方法，可以大幅提升 Transformer 的表现。

# BERT：预训练\_Transformer

这份 PDF 课件是 Shusen Wang 教授关于 **BERT (Bidirectional Encoder Representations from Transformers)** 的讲解。BERT 是 NLP 领域继 Transformer 之后的又一里程碑，标志着 NLP 进入了“预训练+微调”的新时代。

以下是对该课件的**深入分析与解构**：

### 1. 核心定位：Transformer 的“编码器”独立版

课件开篇明确指出了 BERT 的本质：**BERT 就是 Transformer 的 Encoder 部分**。

- **对比**：
  - 标准的 Transformer（机器翻译）：Encoder + Decoder。
  - GPT：Decoder Stack（自回归生成，只能看上文）。
  - **BERT**：Encoder Stack（双向理解，同时看上下文）。
- **目的**：BERT 的目标不是为了“生成”文本（如翻译），而是为了**预训练 (Pre-training)** 一个能够深度“理解”语言的模型。

### 2. 预训练任务解构 (The Pre-training Tasks)

BERT 之所以强大，在于它设计了两个巧妙的**自监督 (Self-supervised)** 任务，让模型可以在海量无标注文本上自己学习。

#### **Task 1: Masked Language Model (MLM, 完形填空)**

这是 BERT 的灵魂。

- **操作**：随机遮盖输入句子中的一个或多个词（用 `[MASK]` 符号代替）。
  - 例如：`the cat sat on the mat` $\rightarrow$ `the [MASK] sat on the mat`。
- **预测**：
  1.  将遮盖后的句子输入 Transformer Encoder。
  2.  取出 `[MASK]` 位置对应的输出向量（Context Vector $\mathbf{u}_{\text{mask}}$）。
  3.  通过一个 Softmax 分类器，预测这个位置原本的词是 "cat"。
- **深度解读**：
  - 传统的语言模型（如 RNN 或 GPT）是单向的（预测下一个词）。
  - BERT 的 MLM 强迫模型利用**上下文（Context）**——即 `[MASK]` 左边的 `the` 和右边的 `sat on...`——来推断中间的词。这使得 BERT 能够学习到真正的**双向语义表示**。

#### **Task 2: Next Sentence Prediction (NSP, 下句预测)**

为了让模型理解句子之间的逻辑关系（这对问答、推理等任务至关重要）。

- **输入**：一对句子 `(Sentence A, Sentence B)`。
- **特殊标记**：
  - `[SEP]`：插入在两个句子中间，作为分隔符。
  - `[CLS]`：插入在整个序列的最开头。这个 token 的输出向量 $\mathbf{c}$ 被专门设计用来代表**整个输入对的语义**。
- **预测**：
  - 取 `[CLS]` 位置的输出向量。
  - 通过一个二分类器（Binary Classifier），判断 Sentence B 是否真的是 Sentence A 的下一句。
  - _正例_：`Calculus is a branch of math` + `It was developed by Newton...` (True)
  - _负例_：`Calculus is a branch of math` + `Panda is native to China` (False)

### 3. 联合训练 (Joint Training)

BERT 并不是分两步训练，而是**同时**做这两个任务。

- **Input Representation**：
  `[CLS] Sentence A [SEP] Sentence B`（其中随机混入 `[MASK]`）
- **Loss Function**：
  $$ Loss = Loss*{\text{NSP}} + Loss*{\text{MLM1}} + Loss\_{\text{MLM2}} + \dots $$
- **优化**：通过一次梯度下降同时更新所有参数。

### 4. 数据与成本 (Data & Cost)

- **数据优势**：BERT 不需要人工标注数据（"Unsupervised" in label sense）。它只需要爬取海量的文本（如 Wikipedia），然后自动生成 `[MASK]` 和句子对标签。
- **计算昂贵**：虽然数据廉价，但算力昂贵。
  - **BERT Base (110M 参数)**：16 个 TPU 跑 4 天。
  - **BERT Large (235M 参数)**：64 个 TPU 跑 4 天。
  - 这也揭示了当前 AI 领域的趋势：大力出奇迹（Scale matters）。

### 总结

这份课件清晰地展示了 BERT 的极简美学：仅仅通过**堆叠 Transformer Encoder** 和设计**两个简单的填空/判断任务**，就构建出了当时最强大的语言理解模型。

- **Masked Word** 训练了它对**词级**语义的理解。
- **Next Sentence** 训练了它对**句级**逻辑的理解。
- **Encoder 架构** 赋予了它**双向**并行的处理能力。
