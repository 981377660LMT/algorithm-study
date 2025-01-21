下面是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 12** (_Fusion Trees_) 的**详细讲解与总结**。本讲继承了之前对**整数数据结构**以及前驱/后继（Predecessor/Successor）问题的讨论，重点介绍**Fusion Tree**的数据结构：它在 **Word RAM** 模型下能以 \(O(\log_w n)\) 时间做前驱/后继查询，显著快于 van Emde Boas 树的 \(O(\log \log u)\)（当 \(w\) 大于 \(\log u\) 级别时更有优势）。

---

## 1. Fusion Tree 的概览

### 1.1 背景与目标

- **模型**：在本讲中，我们使用 **Word RAM**（C风格操作集，包括加减、位运算、移位等）进行分析；每个字长为 \(w\) bit。
- **问题**：维护一个静态集合 \(S\)（大小 \(n\)，每个元素是 \(w\) 位整数），支持前驱、后继查询。
- **Fusion Tree**：由 Fredman 和 Willard 发明 [1][2]，在 **Word RAM** 下，可在  
  \[
  O(\log_w n)
  \]  
  时间完成查询，而空间为 \(O(n)\)。

在之前，我们见过 van Emde Boas (vEB) 结构在 \(O(\log \log u)\) 时间内查询，但当 \(u = 2^w\) 时，\(\log \log u = \log w\)。Fusion Tree 可在较大 \(w\)（当 \(w\) 超过 \(\log n\)）时，更胜一筹，实现更快查询。

### 1.2 Fusion Tree 的核心思路

Fusion Tree 的主体是一个**B树结构**——用 **\(k=w^{1/5}\)** 级的分支因子，从而树高度 \(\approx \log\_{w^{1/5}} n = O(\frac{\log n}{\log w}) = O(\log_w n)\)。

- 关键挑战：**如何在一个节点内用 \(O(1)\) 时间决定要走哪条分支**？因为节点可能含 \(k\) 个关键字，各关键字本身有 \(w\) bits，直接比较似乎要花 \(k \cdot w\) bits。
- 解决方法：利用**Sketching**与**Parallel Comparison**技术，实现对 \(k\) 个关键字与查询值的相对位置判断在 \(O(1)\) 时间内完成。
- 另外，还用到**常数时间定位字的最高设位 (MSB)** 的技巧。

下文详述这些技巧。

---

## 2. Fusion Tree 数据结构的整体思路

Fusion 树是一个**多叉搜索树**（类似 B-树），每个节点可有 **\(k=w^{1/5}\)** 个关键字、\(k+1\) 子指针，整树高度 \(\approx \log*{k} n = \log*{w^{1/5}} n = O(\log_w n)\)。

- 在节点内做一次查找，就可以决定向哪条分支走；理想上只用 \(O(1)\) 时间决定。
- 构造手段：预先对节点内的关键字做**预处理**（sketching等）以加速比较。

操作时间主要由查找的树高决定：\(O(\log_w n)\)。因此，若能在**节点查找**中保证 \(O(1)\)，就可达以上整体复杂度。

---

## 3. Sketching（提取关键位）

### 3.1 Motivating Problem

在节点内有 \(k = w^{1/5}\) 个关键字 \(x*0 < x_1 < \cdots < x*{k-1}\)。要比较某查询 \(q\) 与这些关键字，单纯地看 \(k\cdot w\) bits 太大。但注意到“要区分 \(k\) 个关键字”只需要 \(\log k \approx \frac{1}{5}\log w\) bits信息。

- 因此，可把每个关键字的真正“分歧位”提取到少量 bit ——**Sketch**：这样不必保留关键字全部 \(w\) bit，就能比较它们次序。

### 3.2 Perfect Sketch

如果把树形路径（深度 \(w\)）想象成各关键字的二进制表示，它们之间最多有 \(k-1\) 处“分叉”。把这些分叉对应的 bit 位置提取出来，即可在一个 \(O(k)\approx w^{1/5}\) bit 的 sketch 中保留足够区分的位。

- 在理想情况下，这个**完美sketch**可在常数时间获取（在 AC\(^0\) 模型中视为一个多路位选择电路操作），保序性得以保持：若 \(x_i < x_j\)，则 `sketch(x_i) < sketch(x_j)`。

然而，在 Word RAM 中无法直接做 AC\(^0\) 的完美一次性操作，需要**用乘法与移位**等操作，做成一个近似的sketch（多出来一些空位，但顺序仍保留不变）。

### 3.3 Desketchify

当查询 \(q\) 不属于节点关键字集合时，仅通过 `sketch(q)` 可能会产生歧义：`sketch(q)` 在 sketch 空间里的顺序不一定与 \(q\) 在原空间一致。  
解决方法：找出与 \(q\) 相邻的 sketch neighbor—— `xi, xi+1`（满足 `sketch(xi) <= sketch(q) <= sketch(xi+1)`），然后用**Desketchify**来判断实际在原 \(w\)-bit 空间中，\(q\) 与 \(xi, xi+1\) 的相对位置，从而得知前驱/后继。

---

## 4. Parallel Comparison（并行比较）

### 4.1 问题

在节点中有 \(k\approx w^{1/5}\) 个（sketch 后的）值，每个 \(\approx w^{1/5}\) bit，总共 \(\approx w^{2/5}\) bits。要在常数时间内找 `rank(q)`（确定 \(q\) 在这些keys的顺序位置）。

### 4.2 方法

- 将全部 \(k\) 个 sketch 合并成一个“node sketch”——在每个关键字前加个标志位，然后拼接。
- 把 `sketch(q)` 复制 \(k\) 份合并为 `sketch(q)^k`。
- 做一次**减法**： `node_sketch - sketch(q)^k`；
  - 若对应关键字的 sketch >= `sketch(q)`，则结果块会产生进位或特定标记；
  - 然后用适当的 mask (AND 操作) 提取各块高位 bit 得到并行比较结果。
- 得到一个 bit 向量，标明对每个关键字是大于还是小于 `sketch(q)`。再通过一次**popcount**或**most significant bit**操作找到分界点即可。

这种并行比较在**乘法**+**位运算**的帮助下，只需 \(O(1)\) 时间完成 \(k\) 路比较。

---

## 5. 最高位 (MSB) 运算的常数时间实现

最后一步往往需找“最高设位” (most significant set bit) 或“最左边 1 bit”。这在传统上要 \(\Theta(\log w)\) 操作，但 Fusion Tree 提供一个**复杂的位操作**技巧，在 \(O(1)\) 完成：

- 将 \(w\) bits分成 \(\sqrt{w}\) 段；先找哪一段非空，再在该段内找具体 bit；用类似Sketch/Parallel Comparison手段并行实现。
- 整个过程只需有限次乘法、AND、XOR、移位，即可在 **Word RAM** 下做到常数时间。

---

## 6. 总体实现与复杂度

### 6.1 单节点查找

Fusion Tree 节点包含 \(k=w^{1/5}\) 个关键字：

1. 对关键字做**sketch**并存成 `node_sketch`；
2. 对查询 `q` 做 approximate sketch `sketch(q)`；
3. 并行比较获取 \(q\) 在关键字排序中的位置 `rank(q)`；
4. 通过 `desketchify` 处理判断实际相邻关键字；
5. 得到要走的子指针或在节点内得到前驱/后继结果。

以上核心步骤都能在 **\(O(1)\)** 时间内完成。

### 6.2 整棵树操作

- 整棵 Fusion Tree 是个多叉树，高度 \(\approx \frac{\log n}{\log k}=\frac{\log n}{\log w^{1/5}}=O(\log_w n)\)。
- 每层节点查找 \(O(1)\)，故前驱/后继查询总 \(O(\log_w n)\) 时间。
- 构造或插入等在**静态或动态改进版本**可分别达 \(O(n)\) 建树、或者 \(O(\log_w n)\) 更新的摊还复杂度（在一些实现中增加 \(O(\log\log n)\) 因素）。

---

## 7. 其它模型和变体

- **AC\(^0\) RAM**：若不允许乘法（深度 \(\log w\) 电路），Fusion Tree 有变体 [3] 在更复杂环境下依然可行。
- **动态 Fusion 树**：
  - Exponential Trees [4] 可在 \(O(\log_w n + \log\log n)\) 更新。
  - 基于哈希 [5] 提供期望 \(O(\log_w n)\) 更新等。
- **开放问题**：能否得到**高概率** \(\log_w n\) 的更新？

---

## 8. 总结

**Fusion Tree** 用几项关键技巧（**Sketching**、**Parallel Comparison**、**MSB 常数时间**）在 **Word RAM** 模型下实现**子节点查找只要 O(1)**，最终在多叉搜索树结构中达成**\(O(\log_w n)\) 前驱/后继查询**，空间 \(O(n)\)。对大字长 \(w\)，Fusion Tree 比 vEB (\(\log \log u\)) 更优，提供了一种利用“位并行”加速搜索的数据结构思路。

### 参考文献

- [1] M. L. Fredman, D. E. Willard, _Blasting Through The Information Theoretic Barrier With Fusion Trees_, STOC 1990.
- [2] M. L. Fredman, D. E. Willard, _Surpassing the information theoretic barrier with fusion trees_, JCSS 1993.
- [3] A. Andersson, P. B. Miltersen, M. Thorup, _Fusion trees can be implemented with AC0 instructions only_, TCS 1999.
- [4] A. Andersson, M. Thorup, _Dynamic ordered sets with exponential search trees_, J. ACM 2007.
- [5] R. Raman, _Priority queues: small, monotone, and trans-dichotomous_, ESA 1996.

(完)
