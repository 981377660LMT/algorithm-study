在**位向量（bitvector）的 Rank/Select 数据结构**发展史中，**Jacobson’s Rank**（或称 **Jacobson’s Structure**）是一个里程碑式的方案。它最早由 **Guy Jacobson** 在其 1989 年博士论文 _“Space-efficient static trees and graphs”_ 中提出，演示了如何在只增加很少的冗余空间的情况下，实现对位向量的 `rank` 查询在 **\(O(1)\)** 时间完成。该想法后来被许多后续研究借鉴和推广，比如 **Munro**、**Raman** 等人的工作。下面我们就系统地解读一下 **Jacobson’s Rank** 的核心思路、实现细节和影响。

---

## 1. 问题背景

给定一个长度为 \(n\) 的位向量（bitvector）\(B\)，其元素分别是 0 或 1。我们希望支持如下操作：

- **Rank(1, i)**：统计 \(B[0..i]\)（或 \(B[1..i]\)，具体看实现的下标约定）区间内出现 1 的次数。
- （可选）**Rank(0, i)**：统计 0 的次数，可以用 \( \text{rank}(1, i)\) 的结果做简单转换：\(\text{rank}(0, i) = (i+1) - \text{rank}(1, i)\)（如果下标从 0 开始数起）。

朴素做法中，每次查询若都从头扫描到 \(i\)，需要 \(O(i)\) 时间，若 \(i\) 可能接近 \(n\)，就等于 \(O(n)\) 级别，难以满足快速查询需求。Jacobson 给出的结构能在接近线性空间内，将 `rank(1, i)` 时间降到 \(O(1)\)。

---

## 2. 主要思想：分级索引 + 查表

Jacobson’s Rank 采用**两层或多层的分块**（blocking）与**查表**（lookup）技术来在常数时间内计算 rank。它与现在常见的“分块 + 前缀和 + 小块查表”方案是一脉相承的。其核心包括：

1. **划分大块 (Superblock)**

   - 每个大块大小大约设为 \(\frac{\log n}{2}\) 位或其他与 \(\log n\) 相关的量。
   - 对于第 \(j\) 个大块，保存它之前所有位中 1 的数量（前缀和），存在一个数组 `R_s` 里。这样查询时可迅速知道“到大块开头为止有多少个 1”。

2. **划分小块 (Block)**

   - 每个大块又继续被划分为若干个小块，每个小块大小约为 \(\frac{1}{2}\log n\) 位。
   - 在大块内部，用另一个数组记录该小块相对于大块开头累计的 1 数，如 `R_b`。这样我们就能精确到“到小块开头为止有多少 1”。

3. **在小块内查表**
   - 对于小块内部剩下的部分（可能十几或几十个比特），可以直接用预先构建的查表（lookup table）来回答“这段比特里从头到第 \(r\) 位有多少个 1”之类的问题。
   - 若小块大小是 \(t\) 比特，则可以用一个大小为 \(2^t\) 的表，为每个可能的 bit-pattern 存储其前缀 1 数量。这样查表只需 \(O(1)\)。
   - Jacobson 文中提出的做法，还探讨了如何在更紧凑的空间内保留同样的查表功能。

通过以上三步，就可以在 \(\text{rank}(1, i)\) 查询时：

1. 找到 \(i\) 所属的“superblock”（大块）\(s\)，获取 `R_s[s]`；
2. 找到大块内的小块 \(b\)，加上 `R_b[b]`；
3. 剩余部分用查表获取在该小块里的 1 数。
4. 三者相加就是 \(\text{rank}(1, i)\) 的结果。

总的时间 \(\approx O(1)\)，而空间主要是：存储大块前缀和 + 存储小块前缀和 + 小块查表所需的常量空间。

---

## 3. 空间复杂度与常数优化

### 3.1 Jacobson 结构的空间分析

- **大块前缀数组**：有 \(\frac{n}{\log n/2} = O(\frac{n}{\log n})\) 个大块，每个记录一个数值，数值大小 \(\log n\) bits 足矣（因为最多记录到 n 个 1），所以总空间约为 \(O(\frac{n}{\log n} \times \log n) = O(n)\) bits。
- **小块前缀数组**：每个大块里再分成 \(\frac{\log n/2}{\frac{1}{2}\log n} = O(1)\) 个小块，因此全局也有 \(O(\frac{n}{\log n})\) 小块，每个存 \(\log(\log n)\) bits 的累积值，合计依然 \(O(\frac{n}{\log n} \times \log \log n) = O(\frac{n \log \log n}{\log n})\)。相较于 \(n\) 仍然低阶，可视为 \(o(n)\) 对于大规模数据。
- **小块查表**：如果小块大小是 \(t = \frac{1}{2}\log n\) bits，则查表大小是 \(2^t \times t\) bits，可能看起来是 \(2^{\frac{1}{2}\log n} = n^{1/2}\) 规模，乘以 \(\log n\) 级别仍是 \(O(n^{1/2} \log n)\)，这比 \(n\) 小得多，对大规模 \(n\) 而言是次线性的。
- 综上，Jacobson’s Rank 在原位向量之外增加的**冗余空间**控制在**线性**甚至**低于 n**的级别。

### 3.2 常数时间保证

Jacobson 通过巧妙的编码把大块和小块信息紧凑地打包到一块，而且将查表也做了可常数访问的存储布局，整体使 `rank(1, i)` 查询在常数步内完成。

- 这背后需要一些位运算和分层指针技巧，文章中有详细讨论。
- 在后续文献中，通常引用 Jacobson 的结果，表示**可以在 \(n + o(n)\) bits 空间里实现 rank 在 \(O(1)\) 时间回答**。

---

## 4. 与后续工作的关系

Jacobson’s Rank 结构为后来的**Succinct Data Structures**奠定了基础，常被称为“**Succinct Rank/Select**”的开端之一。后续很多工作在此之上做了改进或变体：

1. **Munro 等人的改进**

   - 进一步讨论了 select 查询的实现，把 rank/select 都做到 \(O(1)\)。
   - 也讨论了对 0 的 rank/select，以及在不同场合如何做平衡分块。

2. **RRR (Raman-Raman-Rao) 结构**

   - 在保证 \(O(1)\) rank/select 的基础上，还尝试在**位向量稀疏度**等情况下做**组合编码**（combinational encoding），使空间更接近信息论极限（\(\log \binom{n}{r}\) bits）。
   - 依然是受 Jacobson’s Rank 启发，分块+表格+组合数的高级压缩方案。

3. **工程实现**
   - 在实际库（如 [SDSL-lite](https://github.com/simongog/sdsl-lite)）中，会看到 `sdsl::bit_vector` + `rank_support` + `select_support` 等结构，大多遵从 Jacobson 风格的分块思想，并辅以各种小块查表优化（如 POPCOUNT CPU 指令等）。

---

## 5. 核心要点小结

- **目标**：在**近线性空间**（n + o(n) bits）里，实现对一个静态位向量的 \(\text{rank}(1,i)\) 或 \(\text{rank}(0,i)\) 查询在 **\(O(1)\)** 时间完成。
- **方法**：分大块 (superblock) + 小块 (block)，大块、小块存前缀 1 数，最末端用查表处理小块剩余比特的 1 计数。
- **关键**：通过巧妙编码和层次化布局，让每一步只需做常数次的数组/表访问和位运算。
- **效果**：这被公认为是**Succinct / 压缩数据结构**最典型的 Rank 技巧之一，对后续的**Select、Wavelet Tree、FM-Index**等影响深远。

---

## 6. 例子：简化示意

假设我们有一个长度 \(n=32\) 的小型 bitvector（仅举例示范）。Jacobson’s 方法可能像这样（仅概念性说明，实际论文中可能更复杂些）：

1. 设大块大小 \(B=8\) 位，小块大小 \(b=4\) 位。
2. 位向量被切为 4 个 superblock，每个 8 bits；每个 superblock 用一个数组记录“到本 superblock 前”为止的 1 数。
3. 每个 superblock 内，再把 8 bits 划为两个 4 bits 小块；对小块记录相对于 superblock 开头的 1 数。
4. 对 4 bits 的所有 \(2^4=16\) 种情况预制一张查表 `popcount[16]`，告诉你“这 4 bits 有多少个 1”，并可细分到前缀 1 数。
5. 查询 `rank(1, i=25)` 时：
   - 找 i=25 属于 superblock index=3（因为 0~7, 8~15, 16~23, 24~31 四个 superblock），查 `prefixSuperblock[3]` 得到 superblock=3 之前 24 bits 的 1 总数；
   - 再看在 superblock=3 里的小块 index，25-24=1 => i 在这个 superblock 的第 1 位；假设它落在小块=0 or 1，需要具体计算；
   - 查相应的 `prefixBlock[]` 得到在该小块开始时累计的 1 数；
   - 最后查 “小块前缀 1” 查表，得到精确在小块内 offset=1 时段落中有多少个 1；
   - 累加即可得到 rank(1,25)。每一步仅做少量的数组访问 + 查表。

对大规模 n 时，可拓展思路，重点是大块与小块大小依赖 \(\log n\) 级别，从而保留 `O(1)` 访问特性。

---

## 7. 结语

**Jacobson’s Rank** 是位向量快速检索结构中最早且最具代表性的方案之一，开启了**Succinct Data Structures**领域的系统研究。其“**分块 + 前缀 + 查表**”的设计理念在后来的许多压缩索引结构、文本检索算法（如 FM-Index）的基础组件“rank/select”中广泛应用。该方案在**理论与实践**上都极具价值，至今仍是实现**静态位向量**快查的经典范式。

---

### 延伸阅读

1. **Jacobson, Guy.** _Space-efficient static trees and graphs._ In _Foundations of Computer Science, 1989. Proceedings., 30th Annual Symposium on_, pp. 549-554. IEEE, 1989.
2. **Clark, D. R.** _Compact Pat Trees._ (PhD thesis), University of Waterloo, 1996. (也对 rank/select 有讨论)
3. **Munro, J. Ian, and V. Raman.** _Succinct Representation of Data Structures._ In _Handbook of Data Structures and Applications_, 2004.
4. **Navarro, G.** _Compact Data Structures: A practical approach._ Cambridge University Press, 2016. (对 Jacobson、RRR、Wavelet Tree 等有详细介绍)

若对**进一步压缩**或**同时支持 select** 也要 \(O(1)\) 时间感兴趣，可查看 **RRR 结构** 和各类后续改进方案。
