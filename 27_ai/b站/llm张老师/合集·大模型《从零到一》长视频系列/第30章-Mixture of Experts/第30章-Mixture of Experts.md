Mixture of Experts (MoE) 通过"稀疏激活"实现了一个看似矛盾的目标——用 12B 激活参数达到 70B 效果，让模型拥有海量参数但每次推理只用其中一小部分，这正是 Mixtral 8x7B 和 DeepSeek-V3 背后的核心架构思想。

- 在标准 Transformer 中，每个 Block 包含两部分：

  Self-Attention（注意力层）
  FFN（前馈网络）

  `MoE 的改动很简单：把 FFN 替换成 MoE 层。`
