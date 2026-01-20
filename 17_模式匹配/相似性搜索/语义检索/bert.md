BERT（Bidirectional Encoder Representations from Transformers）是 NLP 领域的里程碑式模型，由 Google 在 2018 年提出。这一模型彻底改变了自然语言处理的范式，开启了“预训练 (Pre-training) + 微调 (Fine-tuning)”的时代。

以下是对 BERT 的深入讲解，涵盖背景、架构、原理、训练任务及应用。

---

## 1. BERT 的核心定位

在 BERT 出现之前，NLP 主流模型（如 Word2Vec, GloVe）生成的词向量是**静态**的（Static Embedding），即“bank”这个词无论是在“河岸”还是“银行”的语境下，向量都是一样的。

虽然后来的 ELMo 和 OpenAI GPT 引入了**上下文相关**（Contextual）的表示，但它们存在局限性：

- **ELMo**: 使用双向 LSTM，特征提取能力不如 Transformer。
- **GPT (v1)**: 使用 Transformer Decoder，受限于自回归（Autoregressive）特性，只能单向（从左到右）看语境。

**BERT 的核心突破在于：利用 Transformer Encoder 实现了真正的双向（Bidirectional）上下文理解。**

---

## 2. 模型架构

BERT 的主体结构非常简单，它是 **Transformer 的 Encoder 部分** 的堆叠。

### 2.1 结构参数

BERT 主要发布了两个版本：

- **BERT-Base**: L=12 (层数), H=768 (隐藏层维度), A=12 (注意力头数), 参数量 ≈ 1.1 亿。
- **BERT-Large**: L=24, H=1024, A=16, 参数量 ≈ 3.4 亿。

### 2.2 输入表示 (Input Embeddings)

BERT 的输入不仅仅是 Token 的 Embedding，它是三个 Embedding 的加和：

$$Input = TokenEmb + SegmentEmb + PositionEmb$$

1.  **Token Embeddings**:
    - 使用 WordPiece 分词（例如 `playing` -> `play` + `##ing`）。
    - 开头添加特殊标记 `[CLS]`（用于分类任务的聚合表示）。
    - 句子两端或中间添加 `[SEP]`（分隔符）。
2.  **Segment Embeddings**:
    - 用来区分句子对（Sentence A 和 Sentence B）。如果是单句任务，则全为 A。
3.  **Position Embeddings**:
    - 因为 Transformer 是并行计算，无法感知位置信息，需要把位置编码自动学习进去（注意：BERT 学习的是绝对位置 embedding，最长支持 512）。

---

## 3. 预训练任务 (Pre-training Tasks)

BERT 之所以强大，在于它在大规模无标注语料上进行的两个自监督学习任务。

### 3.1 Masked Language Model (MLM) —— 完形填空

这是 BERT 实现“双向”理解的关键。

- **机制**：随机 Mask 掉输入句子中 15% 的 Token，然后让模型预测这些被 Mask 掉的词。
- **训练目标**：最大化被 Mask 词的概率。
- **Mask 策略细节**：为了减小预训练（有 [MASK]）和微调（无 [MASK]）之间的不匹配（Mismatch），对于那 15% 被选中的词：
  - 80% 的概率替换为 `[MASK]` token。
  - 10% 的概率替换为随机的一个词（让模型学会纠错）。
  - 10% 的概率保持不变（让模型学会偏向真实分布）。

> **对比 GPT**：GPT 预测下一个词（必须盖住后面的），BERT 预测缺失的词（可以看到前后的词），因此 BERT 能利用全向信息。

### 3.2 Next Sentence Prediction (NSP) —— 是否由下一句承接

为了让模型理解句子间的逻辑关系（对问答、推理任务很重要）。

- **机制**：输入两个句子 A 和 B。
- **标签**：
  - `IsNext`: B 确实是 A 的下一句（50% 概率）。
  - `NotNext`: B 是语料库中随机抽取的句子（50% 概率）。
- **利用**：取 `[CLS]` 位置输出的向量进行二分类预测。

---

## 4. 微调 (Fine-tuning)

预训练完成后，BERT 就像一个通读了大量书籍的“博学者”。微调就是给它布置具体的考试任务。
BERT 的输出主要利用两个部分：`[CLS]` 向量和 **Sequence Vectors**（每个 token 的向量）。

### 4.1 句子对分类 / 文本分类 (Text Classification)

- **任务**：情感分析、自然语言推理 (MNLI)。
- **做法**：取 `[CLS]` 的输出向量 $C \in \mathbb{R}^H$，接一个全连接层 $W \in \mathbb{R}^{K \times H}$，在这个层上做 Softmax 分类。

### 4.2 序列标注 (Sequence Labeling)

- **任务**：命名实体识别 (NER)、词性标注 (POS)。
- **做法**：将每个 Token 的输出向量 $T_i$ 接分类层，预测该 Token 的标签（如 B-PER, I-LOC）。

### 4.3 问答任务 (SQuAD)

- **任务**：给定文章和问题，找出答案在文章中的起始和结束位置。
- **做法**：学习两个向量 $S$ (Start) 和 $E$ (End)。计算文章中每个 token 向量 $T_i$ 与 $S, E$ 的点积，Softmax 最大的位置即为起止点。

---

## 5. BERT 的优缺点总结

### 优点

1.  **真·双向**：相比 LSTM 的双向（两个单向拼接），Attention 机制能同时看到所有词。
2.  **通用性强**：统一了 NLP 任务架构，几乎所有 NLP 任务都可以刷到 SOTA。
3.  **特征提取能力强**：Transformer 结构深，能够提取复杂的语义和句法特征。

### 缺点 & 局限性

1.  **预训练-微调差异 (Prompt Mismatch)**：预训练时有大量的 `[MASK]`，微调时没有，导致分布不一致（这也是后来 XLNet 改进的点）。
2.  **不适合生成任务 (NLG)**：BERT 是 Encoder 结构，本质是理解模型，不擅长像 GPT 那样逐词生成长文本。
3.  **输入长度限制**：Self-Attention 的复杂度是 $O(n^2)$，导致 BERT 一般限制序列长度为 512，长文本处理困难。
4.  **训练成本高**：显存占用大，训练慢。

---

## 6. BERT 的家族变体 (Evolution)

BERT 发布后，出现了很多改进版本：

1.  **RoBERTa (Robustly optimized BERT approach)**:
    - Facebook 提出。
    - **改进**：去掉了 NSP 任务（发现没啥用）；动态 Masking；更大的 Batch Size；更多的训练数据。效果普遍优于 BERT。
2.  **ALBERT (A Lite BERT)**:
    - Google 提出。
    - **改进**：参数共享（跨层共享参数，大幅减少参数量）；Embedding 层矩阵分解；把 NSP 改为 SOP (Sentence Order Prediction，预测句子顺序)。
3.  **DistilBERT**:
    - Hugging Face 提出。
    - **改进**：使用知识蒸馏（Knowledge Distillation），由大 BERT 此时教小 BERT，保留了 97% 的性能但速度快 60%。
4.  **ELECTRA**:
    - **改进**：类似 GAN 的生成-判别器结构。生成器生成词替换 [MASK]，判别器（Discriminator）判断每个词是原词还是替换词。计算效率极高。

## 7. 代码视角 (Python/PyTorch)

使用 `transformers` 库极其简单：

```python
from transformers import BertTokenizer, BertModel
import torch

# 1. 加载预训练模型和分词器
tokenizer = BertTokenizer.from_pretrained('bert-base-uncased')
model = BertModel.from_pretrained('bert-base-uncased')

# 2. 准备输入
text = "The quick brown fox jumps over the lazy dog."
# 自动加上 [CLS], [SEP] 并转为 Tensor
inputs = tokenizer(text, return_tensors="pt")

# inputs 包含:
# - input_ids: token 的 id
# - token_type_ids: 区分句子 A/B (这里只有 A，全 0)
# - attention_mask: 区分 padding

# 3. 前向传播
with torch.no_grad():
    outputs = model(**inputs)

# 4. 获取输出
# last_hidden_state: [batch_size, seq_len, hidden_size(768)] -> 每个 token 的向量
# pooler_output: [batch_size, hidden_size(768)] -> [CLS] 经过处理后的向量，常用于句级分类
last_hidden_states = outputs.last_hidden_state
cls_vector = outputs.pooler_output

print(f"Token vectors shape: {last_hidden_states.shape}")
print(f"Sentence vector shape: {cls_vector.shape}")
```
