以下是对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 19** (_Dynamic Trees: Link-Cut Trees_) 的**详细讲解与总结**。

---

## 1. 引言与概述

这一讲主要关注**动态树（Dynamic Trees）**。在以往，我们处理过在静态树上做查询（如 LCA、RMQ 等），但对于需要**改变**树结构的场景，还缺乏高效的数据结构。**Link-Cut Tree**（Sleator 与 Tarjan 1983）就是处理这类动态树操作的经典数据结构，可在**对数时间**完成各种操作。

**应用背景**：

- 在**网络流（Network Flow）算法**中，需要频繁修改树/森林结构（如在 Edmond–Karp 或其他增广路径流中处理森林结构等）。
- 在**动态连通性（Dynamic Connectivity）**或“动态生成树”问题中，也可作为内部组件。
- 在**理论**上，它也是自适应、平衡结构的一种典型。

**目标**：我们要维持**一片森林**（多棵有根树），让每个节点可有多个无序的子节点。支持下列操作（全部在摊还 \(O(\log n)\) 时间）：

1. `make_tree()`：返回一个新建单节点树(顶点)。
2. `link(v,w)`：将森林中**另一个**树的根节点 `v` 与 `w`(在其他树中)相连，使 `v` 成为 `w` 的子节点（v 与 w 原本在不同树中，v 为根）。
3. `cut(v)`：切断 `v` 与其父节点之间的边，使 `v` 成为一棵单独子树（v 不是根）；
4. `find_root(v)`：返回 `v` 所在树的根节点；
5. `path_aggregate(v)`：在**根到 v 的路径**上做某种聚合（如对边权取 max/min/sum 等）。我们也可以扩展到节点数据的聚合。

**Link-Cut Tree**核心：

- 通过将原树**分解成一系列优先路径（preferred path）**并用**Splay 树**作为辅助数据结构进行维护。这与之前讲“Tango Trees”类似，也有“preferred child”的概念。
- 每次操作都要调用一个关键子过程 `access(v)`，将树结构调整到等效形态，以便操作执行。

---

## 2. Link-Cut Trees之路径分解 (Preferred-path Decomposition)

### 2.1 基本定义

- **represented tree**：原始的根森林中的一棵树，含根节点、子节点等传统意义下的结构。
- **preferred child**：对任意节点 `v`，其子节点中**最后一次被访问到的**(或访问过子树内节点)的那一支子树(=preferred child)；若没有访问过子树则无优先子。
- **preferred edge**：节点与其 preferred child 之间的边。
- **preferred path**：由preferred edge连续连接成的路径；或者单节点(若该节点无 preferred child 或它本身是无父节点?
- 整颗 represented tree 的**所有节点**被**不交叠**的 preferred path 们分割覆盖。

**Link-Cut Tree**的表示：

- 对每条 preferred path，用一棵**Splay 树**（即**auxiliary tree**）存储其中节点，这些节点按其在 represented tree 中的深度来当做关键字( key )。
- 这些 auxiliary trees 之间通过**path-parent**指针连起来：每棵 auxiliary tree 的 root 节点记下其在 represented tree 中的 parent（若存在）的优先路径结构，从而全局组合成**“一棵辅助森林”**。

### 2.2 关键操作列表 (再次回顾)

1. `make_tree()`: 新建单节点 => trivially done。
2. `link(v,w)`: 要将**两棵树**合并，让 `v` 成为 `w` 的子节点 (其中 `v` 原是根)。
3. `cut(v)`: 将 `v` 与其父节点断开。
4. `find_root(v)`: 找到 `v` 在 represented tree 中的根。
5. `path_aggregate(v)`: 计算代表树中从根到 v 的路径上某种聚合值( sum / min / max )。

### 2.3 Access(v)

链接/切断/查找根/路径聚合 都要用到**`access(v)`**子过程。`access(v)`会对 represented tree 做**调整**：

- 让 `v` 成为其树中的 preferred path 最下端节点(从根一路下来都变成了preferred edges)；
- 具体做法：
  1. 先把 `v` 在其 auxiliary tree 中`Splay(v)`，使 `v`成Splay根；
     - 并把 `v` 的右子树（意味着 deeper node ) 断开 => 这样 v 的子在 represented tree 中不再是其 preferred child(相当于preferred child = NULL)；
     - 这些断开的子节点形成新的 auxiliary tree 并由 `v` 做 `path-parent` 指针指向 v。
  2. 现在要向上（represented tree 中的 parent）再把 v 设为该 parent 的 preferred child，这要求把 parent 之前的下层优先路径断开 => splay 父节点 => 断开其 right child =>把 v 变为其 right child => splay v 使 v 在新的 auxiliary tree 中做根。
  3. 重复往上直到抵达 represented tree 的 root 或无 parent => 形成 v 连接到 root 的单条优先路径。

此过程会修改许多“preferred child”并相应修改 link-cut tree 中 auxiliary tree 的结构(通过Splay,split,join等操作)。最终 v 成为 root->v 的优先路径末端。

---

## 3. 具体实现及操作

### 3.1 find_root(v)

通过 `access(v)`, 然后在 v 所在 auxiliary tree 里找到最深(key最小)节点(往左走到空)，那即represent tree的根，然后 splay之 => \(O(\log n)\) 摊还。

### 3.2 path_aggregate(v)

`access(v)`先把v到根做成一条优先路径 => v所在辅助树包含此路径 => 用Splay树中“节点扩展”可存子树(尤其是 `root->v` 路径)信息 => 结果 O(\log n) 完成。

### 3.3 cut(v)

要切断 `v->parent` 这个represented edge。

- `access(v)` => v becomes top of path => 其 left subtree in splay tree包含v之上的节点；
- 然后 let `left(v)=null` => represented tree断开 => v与其原父分成两块。(Splay 结构会更新 pointer)

### 3.4 link(v,w)

让 `v`(是一棵树的根) 变为 `w` 的子。

- `access(v)`, `access(w)` => v,w各是其所在辅助树根 => connect => `left(v) = w` => done.

---

## 4. 性能分析

### 4.1 费时在 `access(v)`

所有操作 reduce to `access(v)` plus O(\(\log n\)). Focus on `access(v)` complexity:

- `access(v)` 内部做的splay操作(多次) + loop iteration (每次可能做1~2 splay).
- splay 单次摊还 O(\log n)。
- 需 bounding loop iteration => loop迭代数 = # of “preferred child changes” => 可能大 or small ???

**先给出** \(O(\log^2 n)\) 摊还界：

- 认为 loop 最多 \(\log n\) 次(heavy-light argument, see below) each splay \(\log n\) => total \(\log^2 n\).

### 4.2 强化: O(\log n) 摊还

实际上 loop iteration 的“preferred child changes”**总数**在 m 次操作后 <= O(m\(\log n\)) => each operation average \(\log n\). Meanwhile each splay is \(\log n\). => combined O(\log n) amortized.

#### 4.2.1 用Heavy-Light分解

**Heavy edge**: edge to child c if child’s subtree > half of node’s entire subtree => at most one heavy child  
**Light edge**: otherwise.

一次 `access(v)` => 修改**preferred edges**在 root->v 路径 => path中 \(\le \log n\) 处 light edges => “light edge becomes preferred” => sum over all operations => <= m\(\log n\).  
heavy edge becomes preferred can happen rarely => bounding with potential analysis => leads to amortized O(\log n).

---

## 5. 小结

Link-Cut Trees实现了**动态树**的一系列操作 (`link, cut, findRoot, pathAggregate`)均能在**摊还 O(\log n)**时间进行。其核心是**preferred path**分解+**Splay**辅助结构+**重轻分解**(heavy-light decomposition)或潜能分析来控制时间。

**小结**：

- **Link-Cut Tree** = forest of **auxiliary splay trees** + `path-parent` pointers = 优先路径分解
- `access(v)` 核心操作 => 加速后续 `findRoot, pathAggr, link, cut`等
- 分析: **heavy-light** decomposition + splay => each operation O(\log n) amortized.

这是**动态树数据结构**的一大里程碑，也广泛应用于网络流、动态连通性、Euler Tour等领域算法的实现中。

---

**参考**：

- [Sleator, Tarjan 1983,1984], _A Data Structure for Dynamic Trees_ + _Data Structures and Network Algorithms_
- Tarjan, _Amortized complexity..._
- 优化：Link-Cut Trees 变体 e.g. Top Trees, Euler Tour Trees, etc.
