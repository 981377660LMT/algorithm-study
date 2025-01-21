下面是一份对 **MIT 6.851/6.897: Advanced Data Structures (Fall 2017) Lecture 11** (_Integer Data Structures & Predecessor Queries_) 的**详细讲解与总结**。本讲聚焦于**整数数据结构**以及如何在不同模型下实现高效的**前驱/后继查询**（Predecessor/Successor），包括 **van Emde Boas 树**与 **Y-fast 树**。

---

## 1. 背景与模型

### 1.1 固定宇宙假设 (Fixed Universe)

在该主题中，我们常假设**整数宇宙**大小为 \(u\)，即我们有 \(w\)-bit 整数 \([0,1,\dots,u-1]\)（通常 \(u=2^w\)）。很多数据结构可在 \(\text{O}(\log \log u)\) 时间内完成前驱、后继、插入、删除等操作，但要基于某些模型或技巧。

### 1.2 常见模型

1. **Transdichotomous RAM**

   - 内存大小 \(S\)，每个单元宽度 \(w\)；要求 \(w \ge \log S\)。
   - 允许有限集合的操作（+，-，\*，/，% 等），每次操作只访问 \(O(1)\) 个字。

2. **Word RAM**

   - 与 Transdichotomous RAM 类似，但常设定特定操作集合（如 +, -, &, |, <<, >>, <, > 等）
   - 是较常见分析模型，聚焦每个操作在一个机器字（\(w\) bit）上完成。

3. **Cell-Probe 模型**
   - 强模型：对内存的每次读或写算作一次“探针（probe）”，其余计算视为免费。
   - 常用于**下界**分析。

不同模型下，对前驱/后继查询的性能会不同，也会影响空间占用和时间复杂度的**上下界**。

---

## 2. 前驱/后继查询在各种模型下的结果

### 2.1 BST 模型

在简单的比较型二叉搜索树（BST）模型里，每次查询或更新需 \(\Theta(\log n)\) 时间。该结果是经典且**最优**的（符合比较复杂度下界）。

### 2.2 Word RAM 模型

- **van Emde Boas**（vEB）结构：能在 \(O(\log \log u)\) 时间完成查询，需 \(\Theta(u)\) 空间。在 \(u=2^w\) 时，则 \(\log \log u = \log w\)。
- vEB 可做改进：使用哈希减少空间到 \(\Theta(n)\)，但代价是**高概率**的查询时间 \(O(\log \log u)\)。
- **Y-fast 树**（Willard）也可在 \(\Theta(n)\) 空间内达 \(O(\log \log u)\) 时间**高概率**；构造思路是 x-fast 树加分层/间接技巧 (indirection)。
- **Fusion 树**：可以在 \(\Theta(u)\) 空间达 \(O(\log_w n)\) 时间（下一讲详细讨论）。

若我们组合 vEB 与 Fusion Tree 并选择最优者，可达 \(O(\sqrt{\log n \cdot \log \log n})\) 等时间——具体取决于 \(w\) 与 \(n\) 的关系。通常可认为 vEB 在 \(w = O(\log n)\) 时更好，而 Fusion Tree 在 \(w\) 很大时更优。

### 2.3 Cell-Probe 模型下界

已知下界为 \(\Omega(\min\{\log_w n,\ \log \log u\})\) 级别，同时空间若想是 \(n^{\poly \log n}\) 或更少，就会受该下界限制。这与 vEB / Fusion Tree 的上界在大范围内是相匹配的。

---

## 3. van Emde Boas 树

**van Emde Boas** 结构（简称 vEB）是经典的**固定宇宙**前驱/后继数据结构，能在 \(\text{O}(\log \log u)\) 时间完成插入、删除、查询。

### 3.1 结构递归定义

- 把宇宙规模 \(u\) 视为 \(\sqrt{u}\) 个“簇（cluster）”，每个簇包含 \(\sqrt{u}\) 大小；
- 用一个**summary** vEB 结构（规模 \(\sqrt{u}\)）记录哪些簇非空；
- 每个**cluster**本身是一个 vEB 结构，大小也 \(\sqrt{u}\)；
- 存放全局的最小值 (min) 以便快速查询最小值，不用递归查询。

#### 表示方法

常用**高位**/ **低位**划分：对键 \(x\) 的 \(w\)-bit 表示中，前 \(\frac{w}{2}\) bits 称为高位 `high(x)`，后 \(\frac{w}{2}\) bits 为 `low(x)`。

- vEB 将 `high(x)` 用来索引簇 ID，在 `cluster[ high(x) ]` 中插入 `low(x)`；
- 同时在 `summary` 里记录这个簇的 ID。

#### 查询 successor

给定 x：

1. 若 x < min，则 successor(x) = min。
2. 否则在簇 `cluster[ high(x) ]` 中找大于 `low(x)` 的最小值；若存在则组合回返回；
3. 否则找 `c' = successor( summary, high(x) )`；若存在则其 successor = (c', min(cluster[c']))；
4. 否则无后继。

复杂度遵循 \(T(u)= T(\sqrt{u}) + O(1)\)；解得 \(T(u)= O(\log \log u)\)。

### 3.2 插入/删除

插入也遵循类似思路：更新 min/max；在对应 cluster 中递归插入 `low(x)`；若该 cluster 原先为空，则在 summary 里插入 `high(x)`。  
删除时，也要维护 min/max，若 cluster 变空则从 summary 删除等。

### 3.3 另一种视角：二叉树

可以把 \([0,u-1]\) 视为深度 \(w=\log u\) 的完全二叉树，其每个叶子代表一个数字。**vEB** 就是只存储必要的节点(有1标记的路径)，并在每个节点保留对左右子树整合信息（OR值等）。可在 \(\log \log u\) 上下完成搜索，但要做一些额外技巧或指针。  
某些实现能得到 \(O(u)\) 空间与 \(O(\log\log u)\) 查询、\(O(\log\log u)\) 更新。

---

## 4. 缩减空间：x-fast / y-fast 树

### 4.1 x-fast 树

- 只存储所有（key 的前缀）在一个哈希结构中，可 O(1) 查找节点信息。
- 实现 \(\log w\) 时间查询（与 vEB 相同阶），但需 \(O(n w)\) 空间。
- 更新需在每个前缀节点更新信息，成本 \(O(w)\)。

### 4.2 y-fast 树

**y-fast** 树由 **x-fast 树 + 分层/间接**（indirection）构成：

- 将 \(n\) 个元素分成 \(\approx n/w\) 个块，每块大小约 \(w\)，用一个平衡 BST 存之。
- 对每块的一个代表元素建 x-fast 树。
- 这样在查询时先用 x-fast 树（大小 \(\approx n/w\)）找代表块，再在块内以 \(\log w\) 时间找具体前驱/后继。
- 因为块大小 \(\Theta(w)\)，\(w = \log u\) => 块大小 \(\Theta(\log u)\)，查询总为 \(O(\log \log u)\)。
- 更新时：若块过大/过小，则分裂/合并并更新 x-fast 树代表，摊还代价 \(O(\log \log u)\)。
- **空间**：因为 x-fast 树只有 \(\approx n/w\) 块 \* \(w\) 前缀 => \(O(n)\) 总空间。

因此 y-fast 树能在 \(O(\log \log u)\) 时间完成前驱/后继、插入、删除，空间 \(O(n)\)。

---

## 5. 小结

1. **前驱/后继**在比较模型下是 \(\Omega(\log n)\)，但在**Word RAM**或**Cell-Probe**模型下，通过处理整数的位结构可以实现**更快**的 \(\log\log\) 级访问。
2. **van Emde Boas**：\(\Theta(u)\) 空间实现 \(O(\log \log u)\) 操作；可通过哈希或进一步技巧将空间降到 \(O(n)\)，代价是随机化 (high probability)。
3. **y-fast 树**：通过将宇宙分成小块，使用 x-fast 树管理代表元素，达到 \(O(n)\) 空间、 \(O(\log \log u)\) 时间（高概率）性能。

下次将介绍**Fusion 树**（Fusion Tree），在大字长 \(w\) 时可获 \(O(\log_w n)\) 查询复杂度，也是一种整数数据结构的重要实例。
