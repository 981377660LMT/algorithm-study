下面是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 22** (_Memory Models_) 的**详细讲解与总结**。本讲讨论了各种**内存模型**以及在外存计算、分层内存、Cache 忽略模型（cache-oblivious）等环境下的复杂度结果，以及**红蓝棋子游戏 (red-blue pebble game)** 对空间复杂度和 I/O 复杂度的揭示。本讲也衔接了之前介绍的外存模型、cache-oblivious 模型，拓展到多层次层级内存 (hierarchical memory) 及其算法效率的相关内容。

---

## 1. 引言：为什么要研究内存模型？

随着计算机体系结构的发展，现实中存在多级内存层次结构（如 L1/L2 cache、主存、外存等），而传统的 RAM 模型通常只考虑单一级内存和无限随机访问，而**外存模型**（External Memory Model / I/O Model）或**cache-oblivious 模型**能更贴近实际情况：

- 在外存模型中，每次 I/O 访问读取/写入一个大小为 \(B\) 的数据块 (block)；
- 在 cache-oblivious 模型中，算法不知道 \(B\) 和内存大小 \(M\)，却依然能在多层次缓存中取得良好表现。

另一个角度是**Red-Blue Pebble Game**将**CPU高速缓存**与**外部存储**分别用不同颜色标记，以捕捉实际中“访问高速缓存 vs. 访问外存代价迥异”的特征，并从中推导各种问题的 I/O 下界或 I/O 上界。

---

## 2. 二级存储模型 (Two-Level Storage)

### 2.1 理想化两级存储

- 模型：CPU + RAM(外存)
- 块大小 B，每次操作可从 RAM 中最多读取 B 项到缓冲或写 B 项到某块
- 我们关心**块操作**(block operation) 的次数 (I/O 复杂度)

#### 主要结论：

1. **Permutation** \(N\) 个项到 \(\frac{N}{B}\)个指定块中，需要 **\(\Omega(\frac{N}{B}\log B)\)** 块操作(平均情况下, tall-disk: \(\frac{N}{B} > B\))
2. **Sorting** 大小N的序列 => \(\Theta(\frac{N}{B}\log\_{M/B}\frac{N}{B})\) 块操作 => 先前外存算法中有同结论
   - Lower bound: 一次读取\(O(\log(B))\)信息 => 需\(O(\log \frac{N}{B} / \log(B))\), 结合… => \(\Omega(\frac{N}{B}\log\_{M/B}\frac{N}{B})\).
   - Upper bound: 用 **k 路归并排序**(k ~ M/B) => 递归 => \(\Theta(\frac{N}{B}\log\_{M/B}\frac{N}{B})\).

### 2.2 Red-Blue Pebble Game

更进一步，为建模“缓存” vs. “外存”区别：

- **红棋子** (red) = 缓存中的数据；**蓝棋子**(blue)=外存中的数据；
- 操作：从 CPU 角度读/写 => red↔blue 交换； “若所有前驱有红，才能在节点上放红”；
- 研究在 DAG 计算中 red-blue pebble game 来看 I/O complexity。(H. Kung等1981).

对于FFT, 矩阵乘, sorting等：给出**IO复杂度**上下界, 并与 block-size, cache-size 相关。

---

## 3. 外存模型更一般

### 3.1 Scanning, Searching, Sorting

- **扫描**：线性扫描N个元素 => \(\Theta(\frac{N}{B})\) 转移
- **搜索**：在大小N的有序数据中找key => \(\Theta(\log\_{B+1}N)\)
- **排序**： \(\Theta\bigl(\frac{N}{B}\log\_{\frac{M}{B}}\frac{N}{B}\bigr)\)

### 3.2 排列 (Permutation)

- 需 \(\Theta(\min\{N, \frac{N}{B}\log\_{M/B}\frac{N}{B}\})\) 次块操作.
- 因为若 N < sorting bound, 就可无须做全排序，只是把项目放到目标块中即可。

---

## 4. 更多层次下的 I/O => UMH / HMM

我们扩展2-level到**k-level**或者多级层次 ( L1->L2->...->主存->外存 )，为捕捉真实情况。不过分析也变得复杂。

- **UMH**(Uniform Memory Hierarchy) model: 指定若干参数(\(\beta, \alpha\))控制块大小递增关系。E.g. block size随层指数增长… => merge sort, distribution sort, matrix multiply…
- **cache-oblivious**：算法本身不知 \(B\) 与 \(M\)，却希望自动取得 near-optimal I/O 性能。
- “扫描( scanning )” “分块( blocking )” “cache-oblivious mergesort” “van Emde Boas layout for searching”…

**结论**：

- cache-oblivious 下： scanning 仍是 \(\Theta(\frac{N}{B})\)
- searching 需要 \(\Theta(\log_B N)\) blocks => \(\Theta(\log N / \log B)\).
- sorting： 需 \(\Theta\bigl(\frac{N}{B}\log\_{M/B}\frac{N}{B}\bigr)\) blocks(若 tall cache: \(M \ge B^{1+\epsilon}\))
- 这些与外存结论相同，但在CO模型**构造**更精妙(“proactive blocking”)

---

## 总结

Lecture 22 重点在：

- **Idea**: 不同内存模型(二级 / 多级)对I/O复杂度影响
- **主要结论**:
  - Ideal 2-level => Permutation下界 \(\Omega(\frac{N}{B}\log B)\), Sorting \(\Theta(\frac{N}{B}\log\_{M/B}\frac{N}{B})\)
  - Red-Blue pebble game => 常见问题如FFT/矩阵乘/排序 => I/O bounds.
  - UMH/HMM => 大范围多级 cache => mergesort / distribution sort => bounds up/down.
  - Cache-Oblivious => 不显式B,M却能达近似外存效率( under tall-cache assumption )。

这些**模型**与**分析**提供了在IO/外存/多层次存储环境中进行**算法设计与复杂度测评**的理论框架，也是后续研究**IO效率**的基础。

---

**参考**：(同讲义)

- [1] Jia-Wei, Kung, 1981 (Red-Blue Pebble)
- [2] Aggarwal & Vitter, 1988 (外存sort)
- [3] Alpern, Carter, Selker’s UMH model, 1994...
- ... etc.

(完)
