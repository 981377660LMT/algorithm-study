Transformer 的 O(N^2) 复杂度限制了它处理超长序列的能力。State Space Models (SSM) 用 O(N) 的线性复杂度实现序列建模，而 Mamba 通过"选择性"机制让 SSM 首次在语言任务上匹敌 Transformer。混合架构如 Jamba 结合了两者优势，代表了后 Transformer 时代的重要方向。

---

O(N^2) 复杂度限制长序列处理
KV Cache 内存随序列线性增长
这两个问题在超长上下文时变得严重

---

Transformer 统治了过去 7 年，但没有永恒的王者。Mamba、混合架构、以及未来更多的创新，都在告诉我们：保持学习，保持开放。
