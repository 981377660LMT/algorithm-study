以下是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 20** (_Dynamic Connectivity_) 的**详细讲解与总结**。本讲是上一讲动态树（Link-Cut Trees）的继续扩展，将目光放到更一般的**动态图(Dynamic Graph)** 的连通性问题上，尤其是如何维护一个图在插边/删边/查询连通等操作时仍能快速答复 `connected(u,v)` 并且在不同约束场景（如只删边的**Decremental**、只加边的**Incremental**、和完全动态**Fully Dynamic**）中提供有效的数据结构。

以下按照本次讲座的脉络，进行详细梳理。

---

## 1. 问题及已有结果概述

### 1.1 问题：动态连通性 (Dynamic Connectivity)

给定一个**无向图 G**，允许我们对其进行**插入边、删除边、以及插入/删除独立顶点**等操作(据具体定义而定)。要求支持查询 `connected(u,v)`：判断顶点u和v在当前图中是否同一连通分量，以及可能的 `connectedness` 相关扩展（如全图连通与否）。

按更新操作范围分类：

1. **Fully Dynamic**：可以插入和删除边/顶点。
2. **Incremental**：只插入，不删除。
3. **Decremental**：只删除，不插入。

### 1.2 相关已有成果

**各类图**：

- **树 (Tree)**：使用上一讲**Link-Cut**或者**Euler Tour Tree**可达 \(O(\log n)\) 时间，对运算(增删边/查连通)都OK。
  - 如果只删(Decremental)时，我们会介绍一种**常数摊还**时间算法(本讲后部分)。
- **平面图 (Planar Graph)**：Eppstein 等达 \(O(\log n)\) 时间。[2]

**一般图**：

- Thorup (2000) [3]：\(O(\log n(\log\log n)^3)\) 更新，\(O(\frac{\log n}{\log\log\log n})\) 查询
- Holm–Lichtenberg–Thorup (2001) [4]：\(O((\log n)^2)\) 更新, \(O(\frac{\log n}{\log\log n})\) 查询
- **Incremental**：最优 \(\approx \alpha(n)\) (不严格: Union-Find [14])
- **Decremental**：可达 \(O(m\log n + \dots + q)\) => 在稠密图时 ~ \(O(\log n)\) [7]
- **Open**: 是否能 \(O(\log n)\) 对所有更新+查询？

**Worst-case** vs. **Amortized**：

- Eppstein等(1997) [6]：Worst-case \(O(\sqrt{n})\) 更新, \(O(1)\) 查询。
- 其它结构大多是**amortized**结果。
- Worst-case \(O(\polylog n)\)更新 + \(O(\polylog n)\)查询仍是**开放**。

**Patrascu-Demaine (2005)** [5] 给出了一些对称的下界结果(时间复杂度之间的Trade-off等)，具体留到后续**下界**讨论。

---

## 2. Euler-Tour Trees

作为**Link-Cut Trees**的一个较简单替代结构，**Euler-Tour Tree** (Henzinger + King [1]) 也能在**树**(或者森林)的动态场景下实现 \(O(\log n)\) 时间的**Cut/Link/Connected**等操作。它的结构更容易实现或分析，但它对**路径**上的聚合不如 Link-Cut 强大，不过**子树**相关的操作非常适合(后面会用到做子树聚合).

### 2.1 基本思想

针对一棵**根树**，做**欧拉游(Euler Tour)**(也等价于DFS序列)走过每条边两次，从根出发到根结束，把访问节点的顺序记下(节点可能出现多次)。

- 用**平衡BST**(AVL、Red-Black或 splay, etc.)存储这个“欧拉序列”，每个元素对应一次对某节点的访问。
- 每个真实节点 `v` 在 Euler Tour 中可能出现多次(第一次出现指针、最后一次出现指针...都可以存)。
- 当**cut**(v)时，对 Euler Tour 的那段“从 v 第一次出现到 v 最后一次出现”的连续区间做**split**(切割成独立 BST) => 形成新的欧拉序列表 => 表示切断子树。
- 当**link(u,v)**时，需要把 U 子树插入到 V 访问序列相邻位置(在 BST 做相应合并)。
- `connected(u,v)`：判断是否在**同一棵 Euler Tour**(可检查**findRoot**(u) == findRoot(v))
- 其余类似**path**聚合/子树聚合都可在 BST 里做range query。

**在 BST 里**:

- `findRoot(v)`: 就看 BST 里的最小或最大( Euler tour ) => \(O(\log n)\)
- `cut(v)`: 识别 v 对应区间 => split => done.
- `link(u,v)`: 先判断联通性；若不同 => merge BST => done.
- `connected(u,v)`: 就查 findRoot(u)==findRoot(v).
- `subtree_aggregate(v)`: 在 BST 里 `range(firstOccur(v), lastOccur(v))` => aggregated result.

若 BST 有**范式**(size ~ n, depth ~ \(\log n\)) => 每操作 \(O(\log n)\) (合并split查询等) => Euler-Tour Tree get \(O(\log n)\) amortized per op.

---

## 3. Decremental Connectivity (只删边) - O(1)摊还

只允许**删边**（或删点）且一旦删除就不恢复，能否加速查询**connected(u,v)**？

**结论**：在Tree的情况，我们可**常数摊还**完成相同(Connectivity)操作(这节介绍)...

3.1 关键技术：**Leaf Trimming** + Euler Tour

- 逐渐删除边 => leaf trimming => 只要有大子树(≥ \(\log n\)大小)就保留在一个**EulerTour**(或LC-tree) => handle O(n/\(\log n\)) nodes => O(\log n) op => total O(n).
- smaller subtree(≤ \(\log n\))用**bitset** + lookup table => O(1) query.
- 结果: amortized O(1).

**原理**：当一条边删后，可能某子树变 smaller，或者某节点变 leaf => 一步步把大的树(≥\(\log n\)size)用Euler Tour Trees管理，小树(\(\le \log n\)size)直接bitmask / adjacency array => constant time handle.

---

## 4. 用 Euler Tour 树做 Fully Dynamic Connectivity => \(O(\log^2 n)\)

Holm-Lichtenberg-Thorup(2001) [4]中大致方法：

1. 维护**生成森林**(SF)对当前图(保留仅 enough edges in each connected comp.);
2. 需要处理**delete(u,v)** => 若( u,v )在 SF => remove => 可能断开 => 需从graph找替代边(若有) => 需要快速找“能连接u所在组件和v所在组件的最小level边” => EulerTour + multiple-level decomposition =>Time \(\log^2 n\).

具体见[4]，或者本讲仅介绍思路: “仅\(\log\)levels, each a subgraph containing edges at or below that level, each with a spanning forest. Edges can demote to lower level... merging subtrees with replacement edges ... achieve O(\log^2 n) amort.”

---

## 5. 其他动态图问题概述

### 5.1 k-Connectivity

k连通、k-edge连通等：

- planar图 decremental 有 O(\log^2 n).
- smaller k=2,3 ... etc 有特殊 bounds.

### 5.2 Minimum Spanning Forest

Fully dynamic MSF => update O(\log^4 n) or O(\sqrt{n}) worst-case.  
Planar => O(\log n).  
Open for better?

### 5.3 Directed Graph connectivity

有向图的可达性(directed connectivity)更复杂：最优 known ~ O(n^2) update, O(1) query, or other tradeoffs in partially dynamic( BFS-based)...

### 5.4 All-Pairs Shortest Paths, planarity test, etc.

更多更高阶问题(最短路径, planarity, etc.) 都有一些成果(见references)...

---

## 6. 总结

**Euler-Tour Tree**是一种相对简单的**动态树**结构，对**路径**聚合不如 Link-Cut 强，但对**子树**聚合非常自然(因 Euler Tour 中 subtree对应连续区段)。

- Fully dynamic **Connectivity** 在一般图中最优达 \((\log^2 n)\) 摊还更新( Holm+ ) + \(\log n/\log\log n\) 查询Time。
- Decremental树的连通性可达**O(1)** 摊还(leaf trimming + bitmask trick).
- 对一般图, Decremental connectivity ~ O(\log n).
- 整体**动态连通性**仍有很多 open 问题, especially **Worst-case** bounds or o(\(\log n\)) are open.

**参考**：

- Henzinger-King(1995) [1], Eppstein(1997) [6], Holm-Lichtenberg-Thorup(2001) [4], Thorup(2000) [3], Patrascu-Demaine(2005) [5], etc.

(完)
