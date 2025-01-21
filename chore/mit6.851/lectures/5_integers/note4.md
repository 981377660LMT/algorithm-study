下面是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 14** (_Integer Sorting_) 的**详细讲解与总结**。本讲的核心主题是**如何在字长较大的 Word RAM 模型中实现线性或近线性时间的整数排序**，尤其是当 \(w\)（机器字长）足够大（如 \(w \ge (\log^{2 + \epsilon} n)\)）时，可以用较为复杂的**Signature Sort**算法实现在 \(O(n)\) 时间排序。

---

## 1. 背景与动机

### 1.1 整数排序与优先队列

Thorup [7] 表明，若能在时间 \(O(nS(n,w))\) 内对 \(n\) 个 \(w\)-bit 整数进行排序，就可构造一个支持插入、删除和取最小值都在 \(O(S(n,w))\) **最坏情况**时间的优先队列（priority queue）。因此，若要实现**常数时间优先队列**，就需要在 \(O(n)\) 内对整数进行排序——也就是**线性时间排序**成为开放问题。

**开放问题**：能否在任意 \(w\) 值下做真正的 \(O(n)\) 时间整数排序？

### 1.2 已知结果

- **比较模型**：\(O(n\log n)\) 排序下界。
- **Counting Sort**：\(O(n + 2^w)\) 时间，但若 \(w = O(\log n)\) 才可视为线性 \(O(n)\)。
- **Radix Sort**：\(O\bigl(n \frac{w}{\log n}\bigr)\)，当 \(w=O(\log n)\) 时即 \(O(n)\) 线性时间。
- **van Emde Boas (vEB) 排序** [6]：\(O(n \log\log w)\) 或改进到 \(O\bigl(n \frac{\log w}{\log n}\bigr)\)。
- **Signature Sort** [2]：在 \(w \ge (\log^2 n)\) (更精确：\(w \ge \log^{2+\epsilon} n\)) 时可 \(O(n)\) 时间排序。
- **Han [4]**：\(O(n \log\log n)\) in AC\(^0\) 模型。
- **Han & Thorup [5,6]**：在所有 \(w\) 中，能以随机化 \(O\bigl(n \sqrt{\frac{\log w}{\log n}}\bigr)\) 完成排序。

**开放问题**：对于中间区间的 \(w\)，如何得到最优或更好的结果？目前最优算法依赖随机化。

---

## 2. 在 \(w = \Omega(\log^{2+\epsilon} n)\) 时的线性时间排序

**Signature Sort** 算法 ([2]) 能在 \(O(n)\) 时间内排序 \(n\) 个大小 \(\le w\)（具体 \(w_0\le \log^{2+\epsilon}n\)）的整数，当 Word RAM 中字长 \(w\) 满足：\(w \ge (\log^{2+\epsilon}n)\cdot \log\log n\)。其构造比较复杂，包含以下关键步骤与子算法：

1. **Bitonic Sequences**
2. **Logarithmic Merge**
3. **Packed Sorting**
4. **Signature Sort**（主体）

下文按自底向上的方式，先讲 bitonic sequence sorting，再讲快速 merge（并行合并），然后用它来做 packed sort，最后组合到 Signature Sort。

---

## 3. Bitonic Sequences 与并行比较

### 3.1 Bitonic 序列

序列若形如“先单调递增，再单调递减”，或其**循环位移**，即称**bitonic**。可在并行模型中用“bitonic mergesort”在 \(O(\log n)\) 时间内排序该序列；而在 Word RAM 中，如果我们把每步 swap 并行化，就可以在 \(O(\log n)\) 轮完成。这里做为“并行比较”的启发性技术。

### 3.2 Logarithmic Merge on words

假设有 2 个**字**(word)，每个含 \(k\) 个大小为 \(b\) 位的有序元素，总共 \(2k\) 个元素；想用**并行思路**在 \(O(\log k)\) 时间合并成一个有序序列：

1. **Reverse** 第二个 word（元素顺序翻转）；
2. 把它与第一个 word 拼接形成**bitonic**序列；
3. 用**bitonic sort**的方式在并行下只要 \(\log k\) 步；
4. 整个过程同理，能在 \(\log k\) 时间完成对这 \(2k\) 元素的合并。

关键在于：每层比较/交换都可以在**常数**位操作中并行完成（用 bitmask + shift + 标志位，等 trick）。

---

## 4. Packed Sorting

### 4.1 模型假设

- 现有 \(n\) 个大小 \(b\) 位的整数；
- 字长 \(w \ge 2(b+1)\log n \log\log n\)；
- 因此可在一个 \(w\)-bit word 中**打包** \(k = \log n \log\log n\) 个元素，并在元素之间留一些空位（如 1 bit gap）。

### 4.2 算法

1. **初始打包**：将 \(n\) 个 \(b\)-bit 整数放入 \(\frac{n}{k}\) 或 \(O(\frac{n}{k})\) 个 word 里，每 word 含 \(k\) 个元素（可能最后不满），花 \(O(n)\) 时间。
2. 对每个 word 内部**排序**，可递归“split in half + merge” (logarithmic merge) 在 \(O(k)\) 时间把 \(k\) 个元素排好；对 \(\frac{n}{k}\) 个 word 需 \(O(\frac{n}{k}\cdot k)=O(n)\) 时间。
3. 然后对 \(\frac{n}{k}\) 个 word 做“归并排序”，每次 merge 两个有序 word（共 2k 个元素），用 bitonic 并行 merge，耗时 \(O(\log k)\)。
   - \(\frac{n}{k}\) 个 word 进行合并排序会有 \(\log(\frac{n}{k})\) 层，每层花 \(O(\frac{n}{k}\log k)\) 时间。 由于 \(\log k = \log(\log n \log\log n) = O(\log\log n)\)， 所以整体为 \(O(n)\)。

### 4.3 结论

若 \(w \ge 2(b+1)\log n \log\log n\)，则**Packed Sorting**可在 \(O(n)\) 时间将 \(n\) 个 \(b\)-bit 整数排序。

---

## 5. Signature Sort 主算法

处理**一般**情形： \(n\) 个整数，每个 \(w_0 \le \log^{2+\epsilon} n\) bits，机器字长 \(w \ge \log^{2+\epsilon}n \cdot \log\log n\)（或类似条件）即可。

### 5.1 总体思路

1. **将每个整数拆成** \(\log^\epsilon n\) 块（chunk），每块约 \(\frac{w_0}{\log^\epsilon n}\) bits；
2. 对每块做**哈希**(静态完美哈希)成 ~\(O(\log n)\) bits“签名”（signature），这不保序，但可保“相等性区分”；
3. 如此得到**签名串**(Signature)大小 \(\approx \log^{1+\epsilon} n\)，再在 Word RAM 中保证 \(w \ge 2(\ldots)\log n\log\log n\) 便可用**Packed Sorting**在 \(O(n)\) 时间内对签名串排序；
4. 签名本身不保序，需要额外结构**compressed trie**(压缩字典树)来恢复真实顺序。
   - 建立一个压缩 trie，对所有签名串做插入，in-order 即得到签名的排序顺序；
   - 保证此过程在 \(O(n)\) 完成。
5. 由于**签名无法直接对真实数值排序**（哈希不保顺序），必须在 trie 中额外存储**原块**信息。对同一路径节点的 edge 做**递归排序**（再次 chunk + hash + sorted 过程），反复 \(\frac{1}{\epsilon}+1\) 层即把剩余 bits 降到 \(\approx 2(b+1)\log n\log\log n\) 范围，然后**packed sort**收尾。
6. 最终得到按原数值顺序的排序结果。

### 5.2 复杂度

- 拆分 + 哈希：\(O(n)\)。
- 用 packed sort 对签名排序：\(O(n)\)。
- 建**压缩 trie**：插入 \(n\) 串，每次 longest common prefix (通过 MSB(xor)) + walk + create edge => total \(O(n)\)。
- 递归：深度 \(O(\frac{1}{\epsilon}+1)\)，每层花 \(O(n)\)，总 \(O(n)\)。
- 最终**in-order** 输出：\(O(n)\)。

故**Signature Sort**在 \(O(n)\) 时间内完成。当 \(w_0 \le \log^{2+\epsilon} n\) 且 \(w \ge \log^{2+\epsilon} n \cdot \log\log n\)，达到**线性时间**整数排序。

---

## 6. 总结

**Signature Sort** ([2]) 在字长足够大的 **Word RAM** 上实现了**线性时间**排序，填补了当 \(w\approx \log^{2+\epsilon}n\) 时的空白。主要用到：

1. **Bitonic Sequence** 并行排序技巧；
2. **Logarithmic Merge** 并行合并 \(\to O(\log k)\)；
3. **Packed Sort**：在 \(w \ge 2(b+1)\log n\log\log n\) 时可 \(O(n)\) 时间排序；
4. **Signature Sort**：进一步对一般 \(\le \log^{2+\epsilon}n\) 位的整数，通过**分块 + hash** + **压缩 trie**来处理，从而在 \(O(n)\) 时间完成。

**开放问题**：是否能在**所有** \(w\) / \(n\) 场景下（特别在介于 \(\log n\) 到 \(\log^{2+\epsilon}n\) 之间）也达到线性或近线性时间整数排序？目前最好成果是 Han-Thorup [5,6] 的 \(O\bigl(n\sqrt{\frac{\log w}{\log n}}\bigr)\) 随机算法等。

---

### 参考文献

- [2] Andersson, Hagerup, Nilsson, Raman. _Sorting in Linear Time?_ JCSS 57(1):74–93 (1998).
- [4] Y. Han. _Deterministic Sorting in O(n log log n) Time and Linear Space._, J. Algorithms 50(1):96–105 (2004).
- [5] Y. Han, M. Thorup. _Integer Sorting in O(n sqrt(log log n)) Expected Time and Linear Space._, FOCS 2002.
- [6] Kirkpatrick, Reisch. _Upper Bounds for Sorting Integers on Random Access Machines._ TCS 28:263–276 (1984).
- [7] M. Thorup. _Equivalence between Priority Queues and Sorting._ FOCS 2002:125–134.
