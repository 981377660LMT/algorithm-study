在静态位图（bitvector）上实现 **\(O(1)\) 时间**的 **rank** 和 **select** 操作，并且总的额外空间只需 \(o(n)\)（与原位向量大小 \(n\) 相比是低阶的增长）— 这是经典的 **Succinct Data Structure** 研究目标之一。  
**Guy Jacobson**（1989）首先提出可在 **\(n + o(n)\) bits** 的空间里实现 **rank(1, i) = O(1)**；之后 **David Clark**（1996）等人提出了可在同级空间里也让 **select(1, k) = O(1)**，与 Jacobson 方案结合就能得到**同时**支持 \(O(1)\) rank & select 的静态位图结构。

---

## 一、背景与目标

给定一个长度为 \(n\) 的**静态位向量（bitvector）** \(B\)，其中每个位置存放 0 或 1，不允许修改（或极少更新）。我们希望在尽量紧凑的空间下，支持以下查询：

1. \(\mathrm{rank}(1, i)\)：统计 \(B\) 在区间 \([0..i]\) 内出现的 1 的数量（假设下标从 0 开始）；
2. \(\mathrm{select}(1, k)\)：找到第 \(k\) 个出现的 1 的位置索引。

朴素实现每次都要扫描 \(O(i)\) 或 \(O(n)\)，太慢。  
若有一套**额外索引**，可把查询时间降到 \(O(1)\)。但关键是：

- **不能**过度膨胀空间，占用几倍于 \(n\) 的存储；
- **理想目标**：在**接近线性的空间**（\(n + o(n)\) bits）里完成。

两位主要贡献者：

- **Jacobson**：在 1989 年博士论文提出**rank** in \(O(1)\)（又称“Jacobson’s Rank”）；
- **Clark**：在 1996 年博士论文里为 **select** 也达到 \(O(1)\) 给出典型构造方案（“Clark’s Select”）。

最终，这两部分结合，就得到**bitvectors 上 rank/select 都是 \(O(1)\)，额外空间 \(o(n)\)** 的经典结果。

---

## 二、Jacobson’s Rank — \(O(1)\) 的 rank 实现

Jacobson 首先解决 \(\mathrm{rank}(1, i)\) 在 \(O(1)\) 时间的问题。其基本思路是**分块 + 前缀累加 + 小范围查表**：

1. **分大块 (superblock)**

   - 令大块大小 \(\approx \log^2 n\) 或其它与 \(\log n\) 相关的量；把 bitvector 分为 \(\frac{n}{\text{superblockSize}}\) 个大块。
   - 对每个大块，记录“到该大块开头为止”的 1 总数（即前缀和）。存在一个数组 `superCount[]` 中。
   - 这样当要计算 rank(1, i)，可以先立刻得知“到 i 所属大块开头，一共多少个 1”。

2. **分小块 (block)**

   - 在大块内再切成若干小块，每块大小可设 \(\log n\) 左右；
   - 记录每个小块“相对于大块开头”的 1 数。数组 `blockCount[]` 来存。
   - 这样进一步定位到“大块 + 小块”后，就能知道到小块开头的 1 数。

3. **小块余位**
   - 小块剩余的部分（往往不超过 \(\log n\) bits），可以**查表**（precomputed lookup table）或用内置 `popcount` 在常数步内数完其中 1 的个数。
   - 表大小约 \(2^{\log n} = n\) 的规模，或更精巧的区块方法避免过大。

这样，\(\mathrm{rank}(1, i)\) 查询时只需：

1. 计算 i 所在大块 ID => 读 `superCount[]`;
2. 计算在大块内的小块 ID => 读 `blockCount[]`;
3. 对小块里余下 \(\le \log n\) bits 用查表 / popcount => 合计得到准确的 rank。

> **时间**：只需 **\(O(1)\)** 的数组访问 + 小块查表 / popcount。  
> **空间**：大块 + 小块的索引规模约 \(O(\frac{n}{\log n} \times \log n) = O(n)\)，再加上一些 \(o(n)\) 级别的查表或打包技巧，合计仍是 **\(n + o(n)\)** bits 范围内。

Jacobson’s Rank 在当时是突破性的，证明了在不牺牲太多空间的前提下也能将 rank 提速到 \(O(1)\)。后续很多方案（如 RRR 结构、sds-libs）都以此为基础做更多压缩或工程化实现。

---

## 三、Clark’s Select — \(O(1)\) 的 select 实现

在有了 \(\mathrm{rank}\) in \(O(1)\) 后，可以用**二分**在 \(O(\log n)\) 时间做 \(\mathrm{select}\)。  
但 **Clark**（1996）给出一套额外索引，让 \(\mathrm{select}\) 也达成 \(O(1)\)。思路大致是：

1. **对所有 `1` 做分组**

   - 每 \(\alpha\) 个 `1` 形成一组（这里 \(\alpha\) 可能设置为 \((\log n)^2\) 或别的合适值）；
   - 记录每组开头的那一位 `1` 在 bitvector 中的下标 (pos)。
   - 如果想找第 \(k\) 个 `1`，就先算出它在哪一组 `g = \lfloor (k-1)/\alpha \rfloor`，拿到组基准位置 `pos_g`。

2. **组内细分**

   - 在组内可能还有 \(\alpha\) 个 `1`，可以再每 \(\beta\) 个 `1` 分一层索引（\(\beta<\alpha\)），或结合**rank**做快速定位。
   - 这样逐层 refine，直到最后在一个小范围（几十 bits）内确定下标，此时可用 popcount + 跳步 / 查表在 \(O(1)\) 内搞定。

3. **综合结构**
   - 因为能**常数时间**得到 rank(1, x)（来自 Jacobson’s 结构），再加上分层中 “跳到某位置时，知道已经跳过多少个 1”，就能轻松判断要再往后移动多少位。
   - 在 2~3 层分组/索引的帮助下，select(1,k) 也能固定步数完成 => \(O(1)\)。

**空间**：

- 额外索引记录了每 \(\alpha\) 个 `1` 的位置，需要大约 \(\frac{\#1}{\alpha} \times \log n\) bits；再加子分组 \(\beta\)；调参数后总和仍是 \(o(n)\)。

---

## 四、整体：Bitvectors in \(O(1)\) with \(n + o(n)\) bits

综合 Jacobson 与 Clark，可以在**静态** bitvector 上得到：

1. **rank(1,i) in \(O(1)\)**（Jacobson）
2. **select(1,k) in \(O(1)\)**（Clark）
3. **总空间** = \(n + o(n)\) bits（额外索引与查表都只需低阶开销）。

由此我们获得了**最经典**的 Succinct Dictionary 结果之一：

> > 在不改变原 bitvector 的内容（只读），再加上 \(o(n)\) bits 的索引，就能让 rank & select 都变成 \(O(1)\) 时间。

很多后续研究都基于此继续做**压缩**（例如 RRR 组合编码等），或将**动态**（例如可插/删）的场景扩展成更复杂的数据结构，但 Jacobson/Clark 的成果仍是最早、最具代表性的奠基方案。

---

## 五、简要示例流程

设 bitvector 长 \(n=32\)，其中约 10 个 `1`。

- **Jacobson**：
  - 设 superblockSize=8, blockSize=4；
  - 建 superBlockCount[] 和 blockCount[]；
  - rank(1, i=23) => “superBlock 2 => + block 1 => +小块 popcount => \(\mathrm{rank}\)”。
- **Clark**：
  - 把 10 个 `1` 分成每 \(\alpha=4\) 个一组 => 组 1(第1~4个1), 组2(第5~8个1), 组3(第9~10个1)；
  - 记录每组开头 `1` 的位置；再在组内分 \(\beta=2\) 做子索引；最后配合 rank(1, x) 细分。
  - 求 select(1,k=7) => 第7个 `1` => 属于组2 => 在组内是(7-4)=3 => 查子索引+ rank => \(O(1)\)。

---

## 六、总结

1. **Jacobson’s Rank**：

   - 通过“分块 + 前缀计数 + 小块查表”在 \(n + o(n)\) bits 空间内，让 \(\mathrm{rank}(1,i)\) 变成 \(O(1)\)。
   - 是 Succinct Data Structures 的起点之一。

2. **Clark’s Select**：

   - 在已有的 rank(1, ·) \(O(1)\) 基础上，再对所有 `1` 做层次化索引，让 \(\mathrm{select}(1,k)\) 也变成 \(O(1)\)。
   - 整体空间仍为 \(n + o(n)\)。

3. **结果**：
   - bitvector + Jacobson + Clark => \(\mathrm{rank}\) & \(\mathrm{select}\) 同时 \(O(1)\)，空间仅多 \(o(n)\)；
   - 为**文本索引**（Wavelet Tree, FM-Index）、**图结构**（压缩邻接）等提供了非常高效的基础字典结构。

在后续文献，如 **Munro, Raman, Navarro** 等都围绕此做了更多优化、工程化或压缩变体，但 Jacobson 与 Clark 是最早系统实现**rank/select in \(O(1)\)** with \(n + o(n)\) extra bits 的核心方案，至今在各种 Succinct / Compressed Data Structures 中扮演关键角色。
