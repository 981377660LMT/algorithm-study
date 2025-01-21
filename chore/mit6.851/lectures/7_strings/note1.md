下面是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 15** (_Data Structures for Static Trees_) 的**详细讲解与总结**。本讲聚焦于对**静态树**进行预处理，以支持以下三种查询在（接近）**常数时间**内完成，并且保持**线性空间**：

1. **RMQ** (Range Minimum Query)
2. **LCA** (Lowest Common Ancestor)
3. **LA** (Level Ancestor)

下文先描述这三个问题与其相互关系，然后分别介绍如何在线性空间中构造能在常数时间查询的 RMQ 与 LCA 结构，最后讨论相对更难的 LA 问题，并展示一系列数据结构在不同时间-空间复杂度下的实现。

---

## 1. 问题与目标

### 1.1 Range Minimum Query (RMQ)

给定一个数组 \(A\)（大小 \(n\)），需要在预处理后，在查询 \(\mathrm{RMQ}(i,j)\) 时返回区间 \([i,j]\) 中最小元素的**索引**（或者该最小值本身）。

- 例如，\(\mathrm{RMQ}(3,7)\) 找到从 \(A[3]\) 到 \(A[7]\) 之间最小值的位置。

### 1.2 Lowest Common Ancestor (LCA)

给定一颗**有根树** \(T\)（大小 \(n\)），在查询 LCA(x,y) 时，返回**节点 x 与 y 在树中的最低公共祖先**。即它是 x 与 y 的公共祖先中，距离二者最近（或深度最大的）者。

### 1.3 Level Ancestor (LA)

同样给定一颗有根树 \(T\)，`LA(x, k)` 查询返回节点 x 的**第 k 级祖先**（亦即从 x 往上走 k 步得到的祖先节点）。形式化可写：若深度(x)=\(d\)，则 `LA(x, k)` 是深度 \(d-k\) 的那个祖先。

---

## 2. RMQ 与 LCA 之间的等价性

**重要结论**：LCA 与 RMQ 可以**双向**相互归约，因而在研究其中一个问题的线性时间/常数查询解法时，也能为另一个问题提供方案。

1. **RMQ → LCA**：构造**Cartesian Tree**：

   - 对数组 \(A\)，找到全局最小值 \(A[i]\)，令其作为根；左子树处理 \(\{A[\dots i-1]\}\)，右子树处理 \(\{A[i+1\dots]\}\)。
   - 这样得到一个最小堆性质的二叉树。可以证明：\(\mathrm{RMQ}(i,j)\) 在数组上，对应此树中**LCA(节点 i, 节点 j)**。
   - Cartesian Tree 可在 \(O(n)\) 时间用“维护右脊 (right spine)”的方法构造。

2. **LCA → RMQ**：对树做**欧拉游程 (Euler Tour)**或**DFS遍历**，记录每次访问节点的“深度”序列，然后 LCA(x,y) 就变成在此深度数组上做 \(\mathrm{RMQ}(\mathrm{inOrder}(x), \mathrm{inOrder}(y))\)（两节点在遍历中的出现区间的最小值位置）。

另外，对于要在 RMQ 中只比较相对大小，也可把数值映射到 \(\{0,\dots,n-1\}\) 使之成为离散区间的问题。

---

## 3. 常数查询时间的 RMQ 与 LCA（线性空间）

1. **LCA/±1-RMQ** 关系：

   - Bender & Farach-Colton (2000) 显示当序列中相邻元素之差为±1 时（称“±1 RMQ”），可以在**线性空间**和**常数查询**下实现 RMQ。
   - 结合 LCA 到 ±1 RMQ 的构造：对树做“欧拉序列”，深度变化每次 ±1，因而对应 ±1 RMQ。

2. **普通 RMQ**：

   - 先通过 Cartesian Tree + 上述欧拉遍历 -> ±1 RMQ，得到**常数查询, 线性空间**解法。

3. **具体实现思路**：
   - 有一个更简单的解决 RMQ 用 `Sparse Table`：预处理 \(O(n\log n)\)，查询 \(O(1)\) 但空间 \(O(n\log n)\)。
   - 要实现**线性空间, O(1) 查询**，可以采用**分块/分层**(indirection)技巧：把数组分成大小 \(\frac12\log n\) 块，对每种块形态建lookup table,... 并在“父数组”用 Sparse Table(小规模)处理；最总合成 => \(O(1)\) 查询, \(O(n)\) 空间。

---

## 4. Level Ancestor (LA) 问题

静态树中，给定节点 \(x\) 与整数 \(k\)，要找 `LA(x, k)`——x的第k级祖先。下面概述几种方法，各有不同的时间/空间特性 ([Bender-Farach-Colton 2004])。

### 4.1 算法 A: \(\langle O(n^2), O(1)\rangle\)

- 对每个节点 v 与每个可能 k (到根距离) 建一个表：table[v][k] = LA(v,k)。空间 \(O(n^2)\)。
- 查询 O(1)。
- 预处理 O(n^2)。

### 4.2 算法 B: \(\langle O(n \log n), O(\log n)\rangle\)

- 维护**跳指针 (jump pointers)**：对每节点 v，记录 v 的 1-祖先, 2-祖先, 4-祖先, ...
- 查询时自顶向下依次跳，最多 \(\log n\) 次。
- 预处理构造跳表 O(n\log n)，查询 O(\log n)。

### 4.3 算法 C: \(\langle O(n), O(\sqrt{n})\rangle\)

- **Longest path decomposition**：找树中最长路径为一组，然后除去这条路径后在剩余子树递归。
- 在每条路径上存成数组结构(把节点按深度顺序)；一个节点要找 k 级祖先可以先跳到对应路径的第 k 位置，若超过就跳向父路径...
- 可能会跳 O(\sqrt{n}) 次(最坏案例：一条路径 size \(\approx\sqrt{n}\)，加多级嵌套...)。

### 4.4 算法 D: \(\langle O(n), O(\log n)\rangle\)

- **Ladder decomposition**：类似Longest path，但是每条路径延伸一倍(合并) => 使“爬”次数变 \(\log n\)。
- 每次 query 最多 \(\log n\) 跳。

### 4.5 算法 E: \(\langle O(n \log n), O(1)\rangle\)

- 在 Ladder decomposition 基础上，加上**跳指针**思路。
- Query 时可在一次跳指针再 ladder 跳 -> O(1)。
- 预处理 O(n \log n)。

### 4.6 算法 F: \(\langle O(n), O(1)\rangle\)

- 由 Dietz (1991) 实现**最优**预处理与查询的 LA。
- 利用**leaf-based** jump pointers + ladder decomposition + micro-tree trick... 大体思路：只对“叶子数”打 jump pointers，如果树里 leaf 的数L足够小(\(\approx n/\log n\)) => O(n) 预处理 + O(1) query；
- 通过**宏观-微观**分解(路径/子树)等方法控制 leaf 数...

最终可得到**线性时间构造 + O(1) 查询** LA 数据结构，但实现颇为复杂。

---

## 5. 结论

- **静态LCA** (Lowest Common Ancestor) 与 **RMQ** (Range Minimum Query) 是**等价**问题（在预处理 + 常数查询 + 线性空间层面）。常见方法包括**Cartesian Tree** + **欧拉遍历** + **分块**或**Sparse Table**等，可在 \(O(n)\) 预处理后 \(O(1)\) 查询。
- **LA** (Level Ancestor) 问题表面类似 LCA，但无已知简单线性构造/常数查询做法；仍有一系列方法：
  - Jump pointers (log query), Ladder decomposition (log query), or combining them for (1) query but \(\ge O(n\log n)\) 构造；
  - 有**Dietz** 等人的复杂技巧，可在**\(O(n)\) 预处理 + O(1)\) 查询**实现。

总结：静态树上的 RMQ / LCA / LA 都能做到**常数查询、线性空间**，只是 LA 的实现最难，也最复杂细致。
