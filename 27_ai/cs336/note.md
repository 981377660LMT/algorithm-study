# 斯坦福 CS336: Language Modeling from Scratch | 大模型从 0 到 1

> Stanford / Spring 2025
> 讲师：Tatsunori Hashimoto, Percy Liang
> 224n 入门，cs336 精通

## 参考链接 / References

- 课程主页: <http://cs336.stanford.edu/spring2025/>
- YouTube 播放列表: <https://www.youtube.com/playlist?list=PLoROMvodv4rOY23Y0BoGoBGgQ1zmU_MT_>
- 知乎笔记: <https://zhuanlan.zhihu.com/p/1992175622886269951>
- CSDN 笔记: <https://blog.csdn.net/weixin_43807749/article/details/154741242>

1. **双语课程目录** — 19 讲中英文对照表
2. **作业一览** — 5 个作业的内容说明与 GitHub 链接
3. **五阶段学习路线**：
   - **阶段一 (L1-4, A1)**: 基础打桩 — 从零构建 Transformer + 分词器 + 优化器
   - **阶段二 (L5-8, A2)**: 系统优化 — GPU 原理、Triton kernel、分布式并行训练
   - **阶段三 (L9-12, A3)**: 缩放与推理 — 缩放定律、推理优化、评估方法
   - **阶段四 (L13-14, A4)**: 数据工程 — Common Crawl 处理、过滤、去重
   - **阶段五 (L15-17, A5)**: 对齐 — SFT、RLHF、RL + DPO
4. **先修知识检查清单**、**GPU 资源价格表**、**每日学习建议**

---

## 课程目录 / Course Schedule

| #   | 主题 (中文)              | Topic (English)                 | 讲师  | 文件夹                                  |
| --- | ------------------------ | ------------------------------- | ----- | --------------------------------------- |
| 1   | 概述与分词               | Overview & Tokenization         | Percy | `01_overview_and_tokenization/`         |
| 2   | PyTorch 与资源核算       | PyTorch & Resource Accounting   | Percy | `02_pytorch_and_resource_accounting/`   |
| 3   | 模型架构与超参数         | Architectures & Hyperparameters | Tatsu | `03_architectures_and_hyperparameters/` |
| 4   | 混合专家模型             | Mixture of Experts (MoE)        | Tatsu | `04_mixture_of_experts/`                |
| 5   | GPU 原理                 | GPUs                            | Tatsu | `05_gpus/`                              |
| 6   | GPU 内核与 Triton        | Kernels & Triton                | Tatsu | `06_kernels_and_triton/`                |
| 7   | 并行计算 (上)            | Parallelism I                   | Tatsu | `07_parallelism_1/`                     |
| 8   | 并行计算 (下)            | Parallelism II                  | Percy | `08_parallelism_2/`                     |
| 9   | 缩放定律 (上)            | Scaling Laws I                  | Tatsu | `09_scaling_laws_1/`                    |
| 10  | 推理                     | Inference                       | Percy | `10_inference/`                         |
| 11  | 缩放定律 (下)            | Scaling Laws II                 | Tatsu | `11_scaling_laws_2/`                    |
| 12  | 评估                     | Evaluation                      | Percy | `12_evaluation/`                        |
| 13  | 数据 (上)                | Data I                          | Percy | `13_data_1/`                            |
| 14  | 数据 (下)                | Data II                         | Percy | `14_data_2/`                            |
| 15  | 对齐 - SFT / RLHF        | Alignment - SFT / RLHF          | Tatsu | `15_alignment_sft_rlhf/`                |
| 16  | 对齐 - 强化学习 (上)     | Alignment - RL I                | Tatsu | `16_alignment_rl_1/`                    |
| 17  | 对齐 - 强化学习 (下)     | Alignment - RL II               | Percy | `17_alignment_rl_2/`                    |
| 18  | 客座讲座 - 林俊洋 (Qwen) | Guest Lecture - Junyang Lin     | Guest | `18_guest_lecture_junyang_lin/`         |
| 19  | 客座讲座 - Mike Lewis    | Guest Lecture - Mike Lewis      | Guest | `19_guest_lecture_mike_lewis/`          |

---

## 作业 / Assignments

| 作业 | 主题 (中文)                                         | Topic (English)                                              | 代码仓库                                                        | 文件夹                      |
| ---- | --------------------------------------------------- | ------------------------------------------------------------ | --------------------------------------------------------------- | --------------------------- |
| A1   | 基础：分词器、模型架构、优化器                      | Basics: Tokenizer, Model, Optimizer                          | [code](https://github.com/stanford-cs336/assignment1-basics)    | `assignments/a1_basics/`    |
| A2   | 系统：性能分析、FlashAttention (Triton)、分布式训练 | Systems: Profiling, FlashAttention2, Distributed Training    | [code](https://github.com/stanford-cs336/assignment2-systems)   | `assignments/a2_systems/`   |
| A3   | 缩放：理解 Transformer 各组件、拟合缩放定律         | Scaling: Understand Transformer Components, Fit Scaling Laws | [code](https://github.com/stanford-cs336/assignment3-scaling)   | `assignments/a3_scaling/`   |
| A4   | 数据：Common Crawl 处理、过滤、去重                 | Data: Common Crawl Processing, Filtering, Deduplication      | [code](https://github.com/stanford-cs336/assignment4-data)      | `assignments/a4_data/`      |
| A5   | 对齐与推理 RL：SFT、强化学习、DPO (可选)            | Alignment & Reasoning RL: SFT, RL, DPO (optional)            | [code](https://github.com/stanford-cs336/assignment5-alignment) | `assignments/a5_alignment/` |

---

## 学习路线 / Study Roadmap

### 阶段一：基础打桩 (Lecture 1-4, Assignment 1) — 预计 2-3 周

**目标**: 从零构建一个可训练的 Transformer 语言模型

1. **Lecture 1 - 概述与分词**
   - 理解语言模型的整体流程 (数据 → 分词 → 模型 → 训练 → 评估)
   - 实现 BPE (Byte Pair Encoding) 分词器
   - 学习资源：原始 BPE 论文、SentencePiece/tiktoken 源码
2. **Lecture 2 - PyTorch 与资源核算**
   - 掌握 PyTorch 张量操作、自动求导机制
   - 学会计算模型的 FLOPs、内存占用、通信量
   - 练习：手算一个小 Transformer 的参数量和 FLOPs
3. **Lecture 3 - 模型架构与超参数**
   - 深入理解 Transformer 结构：Multi-Head Attention, FFN, LayerNorm, Positional Encoding
   - 掌握 RoPE、GQA、SwiGLU 等现代变体
   - 阅读：Attention Is All You Need, LLaMA 论文
4. **Lecture 4 - 混合专家模型 (MoE)**
   - 理解 MoE 的路由机制、负载均衡
   - 了解 Mixtral、Switch Transformer 的设计
5. **完成 Assignment 1**
   - 实现完整的：分词器 + Transformer 模型 + AdamW 优化器
   - 在小规模数据上训练模型

### 阶段二：系统优化 (Lecture 5-8, Assignment 2) — 预计 2-3 周

**目标**: 让模型在 GPU 上高效运行，掌握分布式训练

6. **Lecture 5 - GPU 原理**
   - GPU 架构：SM、Warp、内存层次 (SRAM/HBM)
   - 理解 compute-bound vs. memory-bound
7. **Lecture 6 - 内核与 Triton**
   - 学习如何用 Triton 编写自定义 GPU kernel
   - 实现 FlashAttention2 的 Triton 版本
8. **Lecture 7-8 - 并行计算**
   - 数据并行 (DDP, FSDP/ZeRO)
   - 张量并行 (Tensor Parallelism)
   - 流水线并行 (Pipeline Parallelism)
   - 阅读：ZeRO 论文、Megatron-LM 论文
9. **完成 Assignment 2**
   - 性能分析 (profiling) 模型各层
   - 用 Triton 实现 FlashAttention2
   - 构建内存高效的分布式训练代码

### 阶段三：缩放与推理 (Lecture 9-12, Assignment 3) — 预计 2 周

**目标**: 理解缩放定律，掌握模型评估与推理优化

10. **Lecture 9 & 11 - 缩放定律**
    - Chinchilla 缩放定律：最优模型大小 vs. 数据量 vs. 计算量
    - 如何用小模型预测大模型性能
    - 阅读：Kaplan et al. (2020), Hoffmann et al. (2022)
11. **Lecture 10 - 推理**
    - KV Cache、投机采样 (Speculative Decoding)
    - 量化 (Quantization)、蒸馏 (Distillation)
    - 推理引擎：vLLM、TensorRT-LLM
12. **Lecture 12 - 评估**
    - 困惑度 (Perplexity) 与下游任务评估
    - Benchmark 套件：MMLU、HellaSwag、HumanEval 等
    - 评估的陷阱与局限性
13. **完成 Assignment 3**
    - 消融实验理解 Transformer 各组件的贡献
    - 通过训练 API 拟合缩放定律曲线

### 阶段四：数据工程 (Lecture 13-14, Assignment 4) — 预计 2 周

**目标**: 掌握预训练数据的全流程处理

14. **Lecture 13-14 - 数据**
    - Common Crawl 数据处理流水线
    - 质量过滤：启发式方法、分类器过滤
    - 去重：Exact dedup (hash) → Near dedup (MinHash/LSH)
    - 数据混合策略、域特定数据的重要性
    - 阅读：RefinedWeb、Dolma、FineWeb 论文
15. **完成 Assignment 4**
    - 实现从原始 Common Crawl dump 到可用训练数据的全流程
    - 实现过滤和去重方法

### 阶段五：对齐 (Lecture 15-17, Assignment 5) — 预计 2 周

**目标**: 掌握从预训练模型到可用助手的对齐技术

16. **Lecture 15 - SFT / RLHF**
    - 监督微调 (Supervised Fine-Tuning) 的数据与方法
    - RLHF 流程：奖励模型 → PPO 训练
17. **Lecture 16-17 - 强化学习对齐**
    - GRPO、DPO 等直接偏好优化方法
    - 推理增强 RL (让模型学会推理)
    - 安全对齐
18. **完成 Assignment 5**
    - 实现 SFT + RL 训练流程
    - 在数学推理任务上训练模型
    - (可选) 实现 DPO 安全对齐

---

## 先修知识检查清单 / Prerequisites Checklist

- [ ] Python 熟练 (重实现、少脚手架)
- [ ] PyTorch 深度学习框架
- [ ] 线性代数基础 (矩阵运算)
- [ ] 概率统计基础
- [ ] 机器学习/深度学习基础 (CS224N 或同等水平)
- [ ] GPU 与内存层次的基本理解

## GPU 资源 (自学用) / GPU Compute for Self-Study

| 平台        | 价格 (H100 80GB)      | 链接                                                |
| ----------- | --------------------- | --------------------------------------------------- |
| RunPod      | $1.99-$2.99/hr        | <https://runpod.io/pricing>                         |
| Lambda Labs | $2.49-$3.29/hr        | <https://lambda.ai/pricing>                         |
| Paperspace  | $2.24/hr              | <https://www.paperspace.com/pricing>                |
| Together AI | $2.85/hr (min 8 GPUs) | <https://www.together.ai/blog/instant-gpu-clusters> |

> 建议：先在 CPU 上调试代码正确性，再用 GPU 跑训练和 benchmark。

---

## 每日学习建议 / Daily Study Tips

1. **看讲座视频** → 做笔记到对应文件夹的 `note.md`
2. **阅读推荐论文** → 重点关注方法论和实验设计
3. **动手做作业** → 这是本课精髓，代码量很大，务必亲自实现
4. **复盘总结** → 每个阶段结束写一篇回顾，连接各知识点

---

# 斯坦福 CS336:大模型从 0 到 1

https://zhuanlan.zhihu.com/p/1992175622886269951
https://blog.csdn.net/weixin_43807749/article/details/154741242

224n入门，cs336精通

https://www.youtube.com/playlist?list=PLoROMvodv4rOY23Y0BoGoBGgQ1zmU_MT_
