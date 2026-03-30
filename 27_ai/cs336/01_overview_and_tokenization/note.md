# Lecture 1: Overview and Tokenization | 概述与分词

> 讲师：Percy Liang
> 视频：<https://www.youtube.com/watch?v=SQ3fZ1sAqXI>
> 主题：课程哲学、语言模型发展简史、分词器原理与 BPE 算法

**核心洞见**：

1. **为什么从零构建** — 抽象层是 leaky 的，研究者正与底层技术脱节
2. **三种知识类型** — Mechanics (可教) / Mindset (可灌输) / Intuitions (部分可教)
3. **苦涩教训的正确解读** — 不是"堆钱就行"，而是 "Algorithms at Scale"

**技术解构**：

- Tokenization 4 种方案对比 (Character/Byte/Word/BPE )
- BPE 算法完整流程 + 可视化示例
- 压缩比的效率意义 ($O(n^2)$ Attention)

**关键反思**：

- 效率是系统设计的核心约束，不是可选优化
- "神圣仁慈"：很多架构决策没有理论解释，只能靠实验
- 从 Tokenizer 开始 = 从最底层理解

---

## 核心洞见 / Core Insights

### 1. 为什么从零构建？| Why Build from Scratch?

**危机感 (Crisis)**：研究者正与底层技术脱节

| 时间线 | 研究者的日常        |
| ------ | ------------------- |
| 8 年前 | 自己实现并训练模型  |
| 6 年前 | 下载 BERT 并微调    |
| 现在   | Prompt 一个闭源 API |

> **Percy 的观点**：抽象层是好的，但它们是 **leaky** 的。与操作系统不同，你甚至不知道这个抽象到底是什么——只是 string in, string out。真正的基础研究需要 **撕开整个栈**，co-design data、systems、model。

**洞见**：Prompt 工程没有错，但它无法推动 **基础性** 研究。如果你想改变游戏规则，你必须理解规则是如何制定的。

---

### 2. 小规模实验的陷阱 | The Trap of Small-Scale Experiments

**问题 1：FLOPs 分布随规模变化**

```
模型规模        Attention FLOPs : MLP FLOPs
-----------------------------------------
小模型          ≈ 1:1
175B 参数       MLPs 完全主导
```

→ 如果你在小规模优化 Attention，你可能在优化 **错误的东西**。

**问题 2：涌现现象 (Emergent Behavior)**

Jason Wei (2022) 的经典图：

- 在某个 FLOPs 阈值之前，准确率 ≈ 0
- 突然"涌现"出 in-context learning 等能力

→ 如果你只在小规模徘徊，你会得出"语言模型不行"的错误结论。

---

### 3. 三种知识类型 | Three Types of Knowledge

| 类型                  | 可教授程度  | 示例                                                       |
| --------------------- | ----------- | ---------------------------------------------------------- |
| **Mechanics (机制)**  | ✅ 完全可教 | Transformer 结构、模型并行如何利用 GPU                     |
| **Mindset (心态)**    | ✅ 可以灌输 | "榨干硬件的每一滴性能"、"认真对待 scaling"                 |
| **Intuitions (直觉)** | ⚠️ 部分可教 | 什么数据/架构决策能产出好模型——小规模的答案 ≠ 大规模的答案 |

> **洞见**：OpenAI 不是发明了 Transformer，而是将 **scaling mindset** 推到极致。这种心态比任何单一技术都重要。

---

### 4. 苦涩的教训 (The Bitter Lesson) 的正确解读

**错误理解**：只有 Scale 重要，算法无所谓，堆钱就行。

**正确理解**：**Algorithms at Scale** 才是答案。

$$
\text{Accuracy} = \text{Efficiency} \times \text{Resources}
$$

- 规模越大，效率越重要（花 1 亿美元，你不能浪费）
- OpenAI 的 GPU 利用率远高于学术界
- 2012-2019 年间，ImageNet 训练的 **算法效率提升 44×**（比摩尔定律还快）

> **心态框架**：给定计算和数据预算，能训练出的 **最佳模型** 是什么？这个问题在任何规模都有意义。

---

## 课程架构 | Course Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     效率优先 (Efficiency First)               │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
   ┌─────────┐          ┌─────────┐          ┌─────────┐
   │ Basics  │          │ Systems │          │ Scaling │
   │ (基础)   │          │ (系统)   │          │ (缩放)   │
   │         │          │         │          │         │
   │ Tokenizer│         │ Kernels │          │ Scaling │
   │ Model   │          │ Parallel│          │  Laws   │
   │ Training│          │ Inference│         │         │
   └─────────┘          └─────────┘          └─────────┘
        │                     │                     │
        ▼                     ▼                     ▼
   ┌─────────┐          ┌─────────┐
   │  Data   │          │Alignment│
   │ (数据)   │          │ (对齐)   │
   │         │          │         │
   │ Crawl   │          │  SFT    │
   │ Filter  │          │  RLHF   │
   │ Dedupe  │          │  GRPO   │
   └─────────┘          └─────────┘
```

**设计决策背后的效率逻辑**：

| 决策                      | 效率理由                                          |
| ------------------------- | ------------------------------------------------- |
| 激进的数据过滤            | 不浪费 precious compute 在垃圾数据上              |
| Tokenization (而非 Bytes) | Bytes 压缩比 = 1，序列太长，Attention 是 $O(n^2)$ |
| 单 epoch 训练             | 我们 compute-constrained，要尽快看更多数据        |
| Scaling Laws              | 用少量计算推断大模型最优超参                      |
| Alignment                 | 投资对齐可以让 **更小的** base model 达到同等效果 |

---

## Tokenization 深度解构 | Tokenization Deep Dive

### 为什么需要 Tokenization？

**目标**：将 Unicode 字符串 → 固定维度整数序列（模型输入）

**核心权衡**：

| 方案            | 词表大小 | 压缩比    | 问题                                |
| --------------- | -------- | --------- | ----------------------------------- |
| Character-based | ~128K+   | ~1.5      | 稀有字符浪费词表空间                |
| Byte-based      | 256      | 1         | 序列太长，效率灾难                  |
| Word-based      | 无界     | 高        | OOV (未登录词) 问题，必须用 `<UNK>` |
| **BPE**         | 可控     | 适中 (~4) | ✅ 自适应分配词表空间               |

**压缩比 (Compression Ratio)** = $\frac{\text{Bytes}}{\text{Tokens}}$

- GPT-2：约 1.6 bytes/token
- 目标：高压缩比 = 短序列 = 更高效的 Attention

---

### BPE 算法精解 | BPE Algorithm Explained

**起源**：1994 年 Phillip Gage 发明的 **数据压缩** 算法
**NLP 引入**：2015 年神经机器翻译 (Sennrich et al.)
**LLM 标准化**：GPT-2 (2019)

#### 核心思想

> 训练 Tokenizer 本身，让它从数据中学习如何分词。常见序列 → 1 个 token，稀有序列 → 多个 token。

#### 算法流程

```python
def train_bpe(text: str, num_merges: int):
    """
    输入：原始文本，合并次数
    输出：merges dict, vocab dict
    """
    # 1. 转成 bytes
    indices = list(text.encode('utf-8'))  # e.g., [116, 104, 101, ...]

    merges = {}  # {(a, b): new_token}
    vocab = {i: bytes([i]) for i in range(256)}  # 初始词表 = 256 bytes

    for i in range(num_merges):
        # 2. 统计相邻 pair 出现次数
        pair_counts = count_pairs(indices)

        # 3. 找出现最多的 pair
        best_pair = max(pair_counts, key=pair_counts.get)

        # 4. 创建新 token
        new_token = 256 + i
        merges[best_pair] = new_token
        vocab[new_token] = vocab[best_pair[0]] + vocab[best_pair[1]]

        # 5. 在 indices 中执行替换
        indices = merge(indices, best_pair, new_token)

    return merges, vocab
```

#### 可视化示例

```
原始: "the cat and the hat"

Step 0 (bytes):
[116, 104, 101, 32, 99, 97, 116, 32, 97, 110, 100, 32, 116, 104, 101, 32, 104, 97, 116]
 t    h    e   _   c   a   t   _   a   n   d   _   t    h    e   _   h   a   t

Step 1: merge (116, 104) → 256  # "th" 出现 2 次
[256, 101, 32, 99, 97, 116, 32, 97, 110, 100, 32, 256, 101, 32, 104, 97, 116]
 th   e   _   c   a   t   _   a   n   d   _   th   e   _   h   a   t

Step 2: merge (256, 101) → 257  # "the" 出现 2 次
[257, 32, 99, 97, 116, 32, 97, 110, 100, 32, 257, 32, 104, 97, 116]
 the  _   c   a   t   _   a   n   d   _  the  _   h   a   t

Step 3: merge (257, 32) → 258  # "the " 出现 2 次
[258, 99, 97, 116, 32, 97, 110, 100, 32, 258, 104, 97, 116]
the_  c   a   t   _   a   n   d   _  the_ h   a   t

压缩比: 19 bytes → 13 tokens = 1.46
```

#### Encode 过程

```python
def encode(text: str, merges: dict) -> list[int]:
    indices = list(text.encode('utf-8'))

    # 按训练顺序回放 merges
    for (a, b), new_token in merges.items():
        indices = merge(indices, (a, b), new_token)

    return indices
```

**关键洞见**：Encode 时必须 **按训练顺序** 回放所有 merges，这是朴素实现的性能瓶颈。

---

### Tokenization 的陷阱 | Tokenization Gotchas

1. **空格前置**：`" hello"` 和 `"hello"` 是 **完全不同** 的 token
2. **数字切分**：`"12345"` 可能被切成 `["123", "45"]`，无语义意义
3. **多语言**：非拉丁语系的压缩比通常更差
4. **可逆性**：BPE 设计为 **无损可逆**，不像传统 NLP 的分词

---

## 语言模型简史 | Brief History of LMs

```
1948  Shannon - 语言模型估计英语熵
2003  Bengio - 第一个神经语言模型
2007  Google - 5-gram, 2T tokens (比 GPT-3 还多！但是 N-gram)
2014  Seq2Seq (Ilya/Google)
2014  Adam 优化器
2015  Attention 机制
2017  "Attention Is All You Need" → Transformer
2018  ELMo, BERT
2019  GPT-2 (+ BPE tokenizer)
2020  GPT-3
2023  GPT-4 (1.8T params, $100M training cost)
2024  DeepSeek, Llama 3, Qwen 2.5...
```

**开放性光谱**：

| 类型         | 示例       | 公开内容                 |
| ------------ | ---------- | ------------------------ |
| Closed       | GPT-4      | 无                       |
| Open Weights | LLaMA      | 权重 + 架构细节 (无数据) |
| Open Source  | OLMo, DBRX | 权重 + 数据 + 训练代码   |

---

## 作业 1 预览 | Assignment 1 Preview

**实现内容**：

1. BPE Tokenizer (从零实现，要快！)
2. Transformer 模型架构
3. 交叉熵损失
4. AdamW 优化器
5. 训练循环

**数据集**：TinyStories, OpenWebText

**Leaderboard**：90 分钟 H100，最小化 OpenWebText perplexity

**警告**：

> "The entire assignment was approximately the same amount of work as **all five assignments from CS 224N plus the final project**."  
> — 课程评价

---

## 关键反思 | Key Reflections

### 1. 效率的哲学意义

Percy 反复强调：**效率不是可选的优化，而是系统设计的核心约束**。

在 compute-constrained 时代：

- 过滤数据 → 不浪费 FLOPs
- Tokenization → 缩短序列
- 单 epoch → 最大化数据覆盖

未来 data-constrained 时代：

- 多 epoch 训练？
- 新架构突破 Transformer 的 compute-efficiency 假设？

### 2. "神圣仁慈" (Divine Benevolence)

SwiGLU 论文结尾：

> "We offer no explanation except for divine benevolence."

**洞见**：很多架构决策缺乏理论解释，只是"实验表明有效"。这意味着：

- 论文中的 intuition 可能是后验的
- 唯一的真理是 **跑实验**
- 保持谦逊，但也保持好奇

### 3. 从 Tokenizer 开始的意义

为什么课程从 Tokenizer 开始？因为它是整个栈中 **最底层** 的抽象。

- 你以为输入是"文本"，其实是整数序列
- 你以为 `"hello"` 和 `" hello"` 一样，其实完全不同
- 你以为长文本只是更多 tokens，其实是 $O(n^2)$ 的 Attention

**从最底层理解，才能在最高层创新。**

---

## 学习资源 | Resources

- **Tokenizer Playground**: <https://tiktokenizer.vercel.app/>
- **Andrej Karpathy Tokenization 视频**: <https://www.youtube.com/watch?v=zduSFxRajkE>
- **BPE 原始论文**: Gage, 1994 "A New Algorithm for Data Compression"
- **NLP 中的 BPE**: Sennrich et al., 2015 "Neural Machine Translation of Rare Words with Subword Units"
- **GPT-2 论文**: Radford et al., 2019 "Language Models are Unsupervised Multitask Learners"

---

## 本讲要点总结 | TL;DR

1. **为什么从零构建**：抽象是 leaky 的，基础研究需要全栈理解
2. **三种知识**：Mechanics (可教) + Mindset (可灌输) + Intuitions (部分可教)
3. **苦涩的教训**：Algorithms at Scale 才是答案，效率比规模更重要
4. **Tokenization 本质**：在 词表大小 vs. 压缩比 之间找平衡
5. **BPE 算法**：迭代合并最频繁的相邻 pair，自适应分配词表
6. **效率驱动一切**：课程的每个设计决策都回到"如何榨干硬件"

---

_下一讲：Lecture 2 - PyTorch & Resource Accounting (FLOPs 都去哪儿了？)_

---

## CS336 第一讲：课程概览与 Tokenization

### 一、课程核心理念

**为什么要从头构建语言模型？**

Percy 指出当前 AI 研究的一个危机：研究者与底层技术越来越脱节。

| 时间  | 研究者行为             |
| ----- | ---------------------- |
| 8年前 | 自己实现和训练模型     |
| 6年前 | 下载 BERT 进行微调     |
| 现在  | 很多人只需要 prompting |

抽象层带来便利，但也有代价：

- 抽象是"泄漏的"（不像操作系统那样边界清晰）
- 基础研究需要打破抽象层，协同设计数据、系统、模型

**核心信念：** _To understand it, you have to build it._

---

### 二、关于小模型的局限性

由于前沿模型（GPT-4：1.8T参数，$100M训练成本）无法在课堂复现，我们只能训练小模型。但要警惕：

1. **FLOPs 分布差异**：小模型中 Attention vs MLP 的计算量相当，但 175B 模型中 MLP 占绝对主导
2. **涌现行为**：某些能力（如 in-context learning）只在大规模时出现

**三种可教授的知识：**

- ✅ **Mechanics**（机制）：Transformer 实现、模型并行等
- ✅ **Mindset**（思维方式）：追求效率、重视 scaling
- ⚠️ **Intuitions**（直觉）：只能部分传授（小规模结论未必适用于大规模）

---

### 三、Bitter Lesson 的正确解读

> **常见误解**：只要堆资源就行，算法不重要
>
> **正确理解**：**Algorithms at scale** 才是关键

证据：2012-2019 年间，ImageNet 训练的**算法效率提升了 44 倍**，比摩尔定律还快。

**正确的思维框架：** _给定固定的计算和数据预算，如何训练出最好的模型？_

---

### 四、课程五大模块

```
效率 (Efficiency) 是贯穿始终的主题
├── Unit 1: Basics（基础）
│   ├── Tokenizer (BPE)
│   ├── Model Architecture (Transformer 变体)
│   └── Training (AdamW, 学习率调度)
│
├── Unit 2: Systems（系统）
│   ├── Kernels (Triton, fusion, tiling)
│   ├── Parallelism (数据并行, 张量并行, FSDP)
│   └── Inference (prefill, decode, speculative decoding)
│
├── Unit 3: Scaling Laws（缩放定律）
│   └── 小规模实验 → 预测大规模超参数
│
├── Unit 4: Data（数据）
│   ├── Evaluation (perplexity, MMLU, 生成评估)
│   └── Curation (HTML→文本, 过滤, 去重)
│
└── Unit 5: Alignment（对齐）
    ├── SFT (监督微调)
    └── Learning from Feedback (DPO, GRPO)
```

---

### 五、Tokenization 详解（本讲技术核心）

**目标**：将 Unicode 字符串 ↔ 整数序列 可逆转换

#### 方案对比

| 方案            | 压缩比 | 优点               | 缺点                          |
| --------------- | ------ | ------------------ | ----------------------------- |
| Character-based | ~1.5   | 简单               | 词表巨大，稀有字符浪费槽位    |
| Byte-based      | 1.0    | 词表固定256，无OOV | 序列太长，Attention O(n²)爆炸 |
| Word-based      | 高     | 自然               | 词表无界，新词变 UNK          |
| **BPE**         | 高     | 自适应，无OOV      | 需要训练，空格处理有怪癖      |

#### BPE 算法核心思想

```
输入: 文本语料, num_merges
输出: merges 表 (pair → new_token)

1. 将文本转为字节序列 (0-255)
2. 重复 num_merges 次:
   a. 统计所有相邻 token pair 的出现次数
   b. 找出最高频的 pair
   c. 为该 pair 创建新 token (256, 257, ...)
   d. 在序列中替换所有该 pair 为新 token
```

#### 代码示例（以 "the cat and the hat" 为例）

```python
# 初始: 字节序列
indices = [116, 104, 101, 32, 99, 97, 116, ...]  # "the cat..."

# 第1次合并: (116, 104) → 256  ("th")
# 第2次合并: (256, 101) → 257  ("the")
# 第3次合并: ...

# 编码时: 按训练时的顺序重放 merges
# 解码时: 递归展开 vocab 映射
```

#### 关键观察

1. **空格是 token 的一部分**：`"hello"` 和 `" hello"` 是完全不同的 token
2. **数字被切分**：`"12345"` 可能变成 `["123", "45"]`，不按语义分组
3. **压缩比**：GPT-2 tokenizer 约 1.6 bytes/token

#### 实际体验

推荐使用 [tiktokenizer.vercel.app](https://tiktokenizer.vercel.app) 可视化不同 tokenizer 的行为。

---

### 六、Assignment 1 预告

实现完整的训练 pipeline：

- BPE tokenizer（从零实现，据称工作量最大）
- Transformer + CrossEntropy Loss
- AdamW optimizer + training loop
- 只允许使用 PyTorch 的少量基础函数

**Leaderboard**：90分钟 H100 时间内最小化 OpenWebText perplexity

---

### 七、核心 Takeaway

> **当前阶段（Compute-constrained）的设计决策都为效率服务：**
>
> - 激进过滤数据（不浪费算力在垃圾数据上）
> - Tokenization（字节模型太慢）
> - 单 epoch 训练（时间紧迫）
> - Scaling laws（用小成本找超参数）

---

下一讲将进入 **PyTorch 细节与资源计算（Resource Accounting）**，关注 FLOPs 的去向。需要我继续整理第二讲吗？
