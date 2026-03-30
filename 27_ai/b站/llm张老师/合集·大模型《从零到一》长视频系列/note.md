https://www.bilibili.com/video/BV1EJ4m1t7Zs
https://www.waylandz.com/llm-transformer-book/

## Transformer 架构：从直觉到实现 (目录)

### 核心部分概览

- **Part 1：建立直觉** (第 1-3 章)
- **Part 2：核心组件** (第 4-7 章)
- **Part 3：Attention 机制** (第 8-12 章)
- **Part 4：完整架构** (第 13-17 章)
- **Part 5：代码实现** (第 18-20 章)
- **Part 6：生产优化** (第 21-22 章)
- **Part 7：架构变体** (第 23-25 章)
- **Part 8：部署与微调** (第 26-27 章)
- **Part 9：前沿进展** (第 28-32 章)

---

### 详细章节目录

- **[前言](https://www.waylandz.com/llm-transformer-book/%E5%89%8D%E8%A8%80/)**
- **第 1 章：[GPT 是什么 — LLM 发展简史与核心思想](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC01%E7%AB%A0-GPT%E6%98%AF%E4%BB%80%E4%B9%88-LLM%E5%8F%91%E5%B1%95%E7%AE%80%E5%8F%B2%E4%B8%8E%E6%A0%B8%E5%BF%83%E6%80%9D%E6%83%B3/)**
- **第 2 章：[大模型的本质 — 就是两个文件](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC02%E7%AB%A0-%E5%A4%A7%E6%A8%A1%E5%9E%8B%E7%9A%84%E6%9C%AC%E8%B4%A8-%E5%B0%B1%E6%98%AF%E4%B8%A4%E4%B8%AA%E6%96%87%E4%BB%B6/)**
- **第 3 章：[Transformer 全景图](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC03%E7%AB%A0-Transformer%E5%85%A8%E6%99%AF%E5%9B%BE/)**
- **第 4 章：[Tokenization — 文字如何变成数字](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC04%E7%AB%A0-Tokenization-%E6%96%87%E5%AD%97%E5%A6%82%E4%BD%95%E5%8F%98%E6%88%90%E6%95%B0%E5%AD%97/)**
- **第 5 章：[Positional Encoding — 给文字加位置](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC05%E7%AB%A0-Positional-Encoding-%E7%BB%99%E6%96%87%E5%AD%97%E5%8A%A0%E4%BD%8D%E7%BD%AE/)**
- **第 6 章：[LayerNorm 与 Softmax — 数字的缩放与概率化](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC06%E7%AB%A0-LayerNorm%E4%B8%8ESoftmax-%E6%95%B0%E5%AD%97%E7%9A%84%E7%BC%A9%E6%94%BE%E4%B8%8E%E6%A6%82%E7%8E%87%E5%8C%96/)**
- **第 7 章：[神经网络层 — 不需要懂也能理解 Transformer](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC07%E7%AB%A0-%E7%A5%9E%E7%BB%8F%E7%BD%91%E7%BB%9C%E5%B1%82-%E4%B8%8D%E9%9C%80%E8%A6%81%E6%87%82%E4%B9%9F%E8%83%BD%E7%90%86%E8%A7%A3Transformer/)**
- **第 8 章：[线性变换的几何意义 — 矩阵乘法的本质](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC08%E7%AB%A0-%E7%BA%BF%E6%80%A7%E5%8F%98%E6%8D%A2%E7%9A%84%E5%87%A0%E4%BD%95%E6%84%8F%E4%B9%89-%E7%9F%A9%E9%98%B5%E4%B9%98%E6%B3%95%E7%9A%84%E6%9C%AC%E8%B4%A8/)**
- **第 9 章：[Attention 的几何逻辑 — 为什么是点积](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC09%E7%AB%A0-Attention%E7%9A%84%E5%87%A0%E4%BD%95%E9%80%BB%E8%BE%91-%E4%B8%BA%E4%BB%80%E4%B9%88%E6%98%AF%E7%82%B9%E7%A7%AF/)**
- **第 10 章：[QKV 到底是什么 — Attention 的三个主角](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC10%E7%AB%A0-QKV%E5%88%B0%E5%BA%95%E6%98%AF%E4%BB%80%E4%B9%88-Attention%E7%9A%84%E4%B8%89%E4%B8%AA%E4%B8%BB%E8%A7%92/)**
- **第 11 章：[Multi-Head Attention — 多视角理解](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC11%E7%AB%A0-Multi-Head-Attention-%E5%A4%9A%E8%A7%86%E8%A7%92%E7%90%86%E8%A7%A3/)**
- **第 12 章：[QKV 输出的本质](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC12%E7%AB%A0-QKV%E8%BE%93%E5%87%BA%E7%9A%84%E6%9C%AC%E8%B4%A8/)**
- **第 13 章：[残差连接与 Dropout — 训练稳定的秘密](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC13%E7%AB%A0-%E6%AE%8B%E5%B7%AE%E8%BF%9E%E6%8E%A5%E4%B8%8EDropout-%E8%AE%AD%E7%BB%83%E7%A8%B3%E5%AE%9A%E7%9A%84%E7%A7%98%E5%AF%86/)**
- **第 14 章：[词嵌入与位置信息的深层逻辑 — 为什么相加而不是拼接](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC14%E7%AB%A0-%E8%AF%8D%E5%B5%8C%E5%85%A5%E4%B8%8E%E4%BD%8D%E7%BD%AE%E4%BF%A1%E6%81%AF%E7%9A%84%E6%B7%B1%E5%B1%82%E9%80%BB%E8%BE%91-%E4%B8%BA%E4%BB%80%E4%B9%88%E7%9B%B8%E5%8A%A0%E8%80%8C%E4%B8%8D%E6%98%AF%E6%8B%BC%E6%8E%A5/)**
- **第 15 章：[Transformer 完整前向传播 — 从输入到输出](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC15%E7%AB%A0-Transformer%E5%AE%8C%E6%95%B4%E5%89%8D%E5%90%91%E4%BC%A0%E6%92%AD-%E4%BB%8E%E8%BE%93%E5%85%A5%E5%88%B0%E8%BE%93%E5%87%BA/)**
- **第 16 章：[训练与推理的异同 — 为什么推理要一个字一个字生成](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC16%E7%AB%A0-%E8%AE%AD%E7%BB%83%E4%B8%8E%E6%8E%A8%E7%90%86%E7%9A%84%E5%BC%82%E5%90%8C-%E4%B8%BA%E4%BB%80%E4%B9%88%E6%8E%A8%E7%90%86%E8%A6%81%E4%B8%80%E4%B8%AA%E5%AD%97%E4%B8%80%E4%B8%AA%E5%AD%97%E7%94%9F%E6%88%90/)**
- **第 17 章：[学习率的理解 — 训练稳定的关键](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC17%E7%AB%A0-%E5%AD%A6%E4%B9%A0%E7%8E%87%E7%9A%84%E7%90%86%E8%A7%A3-%E8%AE%AD%E7%BB%83%E7%A8%B3%E5%AE%9A%E7%9A%84%E5%85%B3%E9%94%AE/)**
- **第 18 章：[手写 Model.py — 模型定义](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC18%E7%AB%A0-%E6%89%8B%E5%86%99Model.py-%E6%A8%A1%E5%9E%8B%E5%AE%9A%E4%B9%89)**
- **第 19 章：[手写 Train.py — 训练循环](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC19%E7%AB%A0-%E6%89%8B%E5%86%99Train.py-%E8%AE%AD%E7%BB%83%E5%BE%AA%E7%8E%AF)**
- **第 20 章：[手写 Inference.py — 推理逻辑](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC20%E7%AB%A0-%E6%89%8B%E5%86%99Inference.py-%E6%8E%A8%E7%90%86%E9%80%BB%E8%BE%91)**
- **第 21 章：[Flash Attention — 内存优化原理](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC21%E7%AB%A0-Flash-Attention-%E5%86%85%E5%AD%98%E4%BC%98%E5%8C%96%E5%8E%9F%E7%90%86/)**
- **第 22 章：[KV Cache — 推理加速](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC22%E7%AB%A0-KV-Cache-%E6%8E%A8%E7%90%86%E5%8A%A0%E9%80%9F/)**
- **第 23 章：[MHA 到 MQA 到 GQA 演进](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC23%E7%AB%A0-MHA%E5%88%B0MQA%E5%88%B0GQA%E6%BC%94%E8%BF%9B/)**
- **第 24 章：[Sparse 与 Infinite Attention](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC24%E7%AB%A0-Sparse%E4%B8%8EInfinite-Attention/)**
- **第 25 章：[位置编码演进 — Sinusoidal 到 RoPE 到 ALiBi](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC25%E7%AB%A0-%E4%BD%8D%E7%BD%AE%E7%BC%96%E7%A0%81%E6%BC%94%E8%BF%9B-Sinusoidal%E5%88%B0RoPE%E5%88%B0ALiBi/)**
- **第 26 章：[LoRA 与 QLoRA — 高效微调](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC26%E7%AB%A0-LoRA%E4%B8%8EQLoRA-%E9%AB%98%E6%95%88%E5%BE%AE%E8%B0%83/)**
- **第 27 章：[模型量化 — GPTQ / AWQ / GGUF](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC27%E7%AB%A0-%E6%A8%A1%E5%9E%8B%E9%87%8F%E5%8C%96-GPTQ-AWQ-GGUF/)**
- **第 28 章：[Prompt Engineering — 提示工程实战](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC28%E7%AB%A0-Prompt-Engineering-%E6%8F%90%E7%A4%BA%E5%B7%A5%E7%A8%8B%E5%AE%9E%E6%88%98/)**
- **第 29 章：[RLHF 与偏好学习 — 让模型对齐人类](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC29%E7%AB%A0-RLHF%E4%B8%8E%E5%81%8F%E5%A5%BD%E5%AD%A6%E4%B9%A0-%E8%AE%A9%E6%A8%A1%E5%9E%8B%E5%AF%B9%E9%BD%90%E4%BA%BA%E7%B1%BB/)**
- **第 30 章：[Mixture of Experts — 稀疏激活的秘密](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC30%E7%AB%A0-Mixture-of-Experts-%E7%A8%80%E7%96%8F%E6%BF%80%E6%B4%BB%E7%9A%84%E7%A7%98%E5%AF%86/)**
- **第 31 章：[推理模型革命 — 从 o1 到 R1](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC31%E7%AB%A0-%E6%8E%A8%E7%90%86%E6%A8%A1%E5%9E%8B%E9%9D%A9%E5%91%BD-%E4%BB%8Eo1%E5%88%B0R1/)**
- **第 32 章：[后 Transformer 架构 — Mamba 与混合模型](https://www.waylandz.com/llm-transformer-book/%E7%AC%AC32%E7%AB%A0-%E5%90%8ETransformer%E6%9E%B6%E6%9E%84-Mamba%E4%B8%8E%E6%B7%B7%E5%90%88%E6%A8%A1%E5%9E%8B/)**
- **附录 A：[Scaling Law 与计算量估算](https://www.waylandz.com/llm-transformer-book/%E9%99%84%E5%BD%95A-Scaling-Law%E4%B8%8E%E8%AE%A1%E7%AE%97%E9%87%8F%E4%BC%B0%E7%AE%97/)**
- **附录 B：[解码策略详解](https://www.waylandz.com/llm-transformer-book/%E9%99%84%E5%BD%95B-%E8%A7%A3%E7%A0%81%E7%AD%96%E7%95%A5%E8%AF%A6%E8%A7%A3/)**
- **附录 C：[常见问题 FAQ](https://www.waylandz.com/llm-transformer-book/%E9%99%84%E5%BD%95C-%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98FAQ/)**
