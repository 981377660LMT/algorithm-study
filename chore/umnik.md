# Things I don't know (by um_nik)

https://codeforces.com/blog/entry/92248

Stop learning useless algorithms,
go and solve some problems,
learn how to use binary search

---

Things that I have heard of, but don't know (imagine how many things I haven't even heard of):

- Li-Chao Segment Tree
- Segment Tree Beats
- RMQ in O(n)/O(1)
- Any self-balancing tree except treap
- Link-cut tree
- Wavelet tree
- Mergesort tree
- Binomial heap
- Fibonacci heap
- Leftist heap
- Dominator tree
- 3-connected components in O(n)
- k-th shortest path
- Matching in general graph
- Weighted matching in general graph
- Preflow-push
- MCMF in O(poly(V,E))
- Minimum arborescence (directed MST) in O(ElogV)
- Suffix tree
- Online convex hull in 2D
- Convex hull in 3D
- Halfplane intersection
- Voronoi diagram / Delaunay triangulation
- Operation on formal power series (exp, log, sqrt, ...) (I know the general idea of Newton method)
- How to actually use generating functions to solve problems
- Lagrange Inversion formula
- That derivative magic by Elegia
- That new subset convolution derivative magic by Elegia
- How Elegia's mind works
- Sweepline Mo
- Matroid intersection

---

下面是一份“百科式”的简要扫盲清单，涵盖了列表中每个主题的核心概念、用途，以及在算法竞赛或实际应用中的位置。由于篇幅所限，这些解释只能算是概要介绍，帮助你建立初步印象。如需深入，建议根据关键词查阅具体教程、论文或权威书籍。

---

## 1. **Li-Chao Segment Tree**

- **用途**：在一个一维区间上维护一组线性函数 \((f_i(x) = a_ix + b_i)\)，并支持查询给定 \(x\) 点时的函数最小值/最大值。
- **原理**：类似线段树，但每个节点只记录“一条最优线/直线”。在插入新线时，通过与当前节点所记录的线比较它们在该区间的分界，递归更新左右子区间。
- **特点**：实现起来比“Convex Hull Trick”方式更直接，可以在 \(\mathbb{R}\) 上维护上下凸壳的最值，是一种简化的线段树结构。
- **复杂度**：若坐标或查询点范围很大（甚至是 \([-10^{9},\,10^{9}]\)），也可用“动态开点”的方式实现，插入/查询都在 \(O(\log(\text{range}))\) 左右。

---

## 2. **Segment Tree Beats**

- **用途**：解决一些“区间更新 + 区间查询”非常复杂的场景，如同时支持区间 \(\min\)/\(\max\) 截断、区间加法、区间求和、查询区间最大/最小等。
- **核心思路**：在传统线段树维护信息的基础上，额外保存区间最值、次最值以及计数等统计数据。在进行「\(\min(a_i, x)\) / \(\max(a_i, x)\)」类型的更新时，通过这些统计识别哪些子区间完全满足被截断的条件，哪些需要继续递归下探。
- **实现难点**：需要非常仔细地处理各种“分情况更新”，但一旦写成模板，对复杂区间操作非常有用。

---

## 3. **RMQ in \(O(n)\)/\(O(1)\)**

- **RMQ**：Range Minimum (or Maximum) Query，即在一个静态数组中多次查询区间最小值/最大值。
- **Sparse Table**：预处理 \(O(n \log n)\) 后可在 \(O(1)\) 查询区间最值；
- **Cartesian Tree + LCA** / **RMQ in \(O(n)\)**：还有一种线性预处理算法，可以做到整体预处理 \(O(n)\) 并使得每个查询 \(O(1)\)。
- **常用场景**：需要频繁地在静态数组上做“区间最值”查询，比如处理某些 DP、或分割问题等。

---

## 4. **任意自平衡搜索树 (Self-balancing BST) 除了 Treap**

- **常见种类**：AVL 树、Red-Black Tree、Splay Tree、Treap、Scapegoat Tree、B-Tree（及其变体）等等。
- **思路**：通过旋转 (rotation) 或者重新平衡 (rebalance) 的方式，在插入/删除节点后保持树高在 \(O(\log n)\)。
- **应用**：在标准库里常见的 `std::set` / `std::map` (C++ 使用 Red-Black Tree)；AVL 则更严格地控制高度；Splay 在访问局部性强时效率好。
- **Treap**：以随机优先级维护堆性质 + 以节点键值维护 BST 性质；这里题主已经提到过，不再赘述。

---

## 5. **Link-Cut Tree**

- **用途**：动态维护树形结构，支持将一棵树的某条边“切断”，或将若干棵树“连”到一起，并可快速（\(O(\log n)\)）查询节点到根的路径信息（例如路径和、路径最大值）。
- **结构**：一种高级的动态树数据结构，与 **Euler Tour Tree**、**ET-Tree** 等同属“动态树”范畴。
- **核心思路**：将节点间的父子关系拆解为若干辅助数据结构（如 Splay、Treap 等）维护，并能在不同链接关系下随意分割/合并。
- **应用**：在涉及“动态改变树结构”或“反复合并/分裂”的问题（如在线网络连通、维护树上路径计算等）中十分有用。

---

## 6. **Wavelet Tree**

- **用途**：在一个静态序列上实现多种查询：如“在子区间 \([L,R]\) 内第 k 小的元素、出现次数、排名”等。
- **结构**：Wavelet Tree 分层地对序列的数值进行二进制划分，用树状的方式存储“哪些位置落入左子区间/右子区间”。
- **优点**：在 \(\log(\text{Range})\) 或 \(\log(n)\) 的层数内，就可以完成各种“排名/选择/数值频次”查询。
- **应用**：常见于压缩后的离散值域；在数据分析、竞赛题中用于“离线处理 + 快速查询统计”。

---

## 7. **Merge Sort Tree**

- **用途**：在静态数组上建线段树，每个节点存“该区间的元素的有序列表”；可用来做区间内第 k 小值或小于 x 的个数等查询。
- **原理**：构造时，每个节点对应区间 \([L,R]\) 的元素排序好的列表，内部通过合并左右子节点的列表构建（**类似归并排序 Merge 的过程**）。
- **复杂度**：构造 \(O(n \log n)\)，查询单次 \(O(\log n \cdot \log n)\)（二分 + 下探）。
- **应用**：竞赛中有时比 Fenwick/Segment Tree + offline trick 更易写且直观。

---

## 8. **Binomial Heap**

- **用途**：一种可合并 (meld) 的优先队列结构，支持插入、取最小值、合并两个堆等操作。
- **结构**：由一系列“二项树 (binomial tree)”构成，每棵树的节点数是 \(2^k\)；通过链接/合并这些树来实现操作。
- **特点**：插入/合并在摊还意义下可以很高效，最坏复杂度一般在 \(O(\log n)\)。
- **应用**：在需要频繁合并优先队列的场景很有用，不过实际实现中常用 **Fibonacci Heap** 或更简洁的 **Pairing Heap** 等。

---

## 9. **Fibonacci Heap**

- **用途**：同样是一种可合并堆，Dijkstra 算法的经典理论实现提到过它可以使得最短路达成 \(O(E + V \log V)\)。
- **结构**：采用“最小树根 + 双向链表 + rank”管理，“删节点”或“合并节点”都以延迟合并的策略进行。
- **特点**：摊还复杂度插入 \(O(1)\)，合并 \(O(1)\)，提取最小 \(O(\log n)\)。
- **实践**：虽然理论上好，但实际常用 Pairing Heap、Binomial Heap；实现相对简单且常数因素更小。

---

## 10. **Leftist Heap**

- **用途**：一种合并堆(merge heap)，通过维护“左偏性” (leftist property) 来让合并操作在最坏情况下保持 \(O(\log n)\)。
- **结构**：类似二叉堆，但每个节点存一个“距离（或零距离）”记录左侧最短路径，保证合并时优先倾向左。
- **优势**：实现简单，常见于教学或竞赛用，可快速合并两个堆。
- **对比**：效率不及 Fibonacci Heap 的摊还性能，但比二叉堆更适合“合并”操作。

---

## 11. **Dominator Tree**

- **用途**：在有向图（一般是控制流图 CFG）中，对某个起始点 \(s\)，若一个节点 \(d\) 在所有从 \(s\) 出发到 \(v\) 的路径上都出现，那么 \(d\) 称为 \(v\) 的 dominator。Dominator Tree 则把节点按谁支配谁的关系组织成一棵树。
- **应用**：编译器优化、程序分析、图理论；在代码流分析中非常常见。
- **算法**：有 Tarjan 等提出的基于“半支配 (semi-dominator) 定理”的线性时间算法。

---

## 12. **3-connected Components in \(O(n)\)**

- **概念**：3-连通分量是指图在移除不超过 2 个顶点后仍保持连通的最大部分；这是图论中更高级的连通性概念。
- **经典算法**：能在 \(O(V + E)\) 线性时间内分解一个平面图或一般图的 3-连通分量（使用 DFS + low-link 技巧等）。
- **意义**：在网络可靠性、分割问题等方面有应用，跟双连通分量（桥与割点的延伸）类似但更高级。

---

## 13. **\(k\)-th Shortest Path**

- **问题**：找到从起点到终点的第 \(k\) 条最短路（按路径长度从小到大排序）。
- **常见方法**：
  - A\* + best-first search + 堆
  - Yen’s Algorithm, Eppstein’s Algorithm 等
- **复杂度**：往往不如一次最短路那么好做，因为要维护很多可能分支；在竞赛题中经常出现需要输出前 \(k\) 小路径。

---

## 14. **Matching in General Graph / Weighted Matching in General Graph**

- **含义**：在一般图（非二分图）中寻找最大匹配 / 最小边覆盖；若加上权值则是最小花费 / 最大权匹配。
- **经典算法**：
  - **Blossom Algorithm**（Edmonds’ Algorithm）：解决无权一般图最大匹配
  - **Hungarian-like Algorithm/Edmonds-Karp** 变形：在加权一般图匹配也可以用更复杂的 Blossom + 费用流。
- **应用**：调度、排课、网络路由、运营研究中大量运用；实现起来通常比二分图匹配复杂许多。

---

## 15. **Preflow-Push**

- **用途**：最大流算法家族之一，相比经典的 Ford-Fulkerson / Edmond-Karp / Dinic，Preflow-Push 在某些稠密图中表现更好。
- **原理**：通过“预流 (preflow)”的概念让中间状态中可能出现局部流量“过满”，再通过“推送 (push)”和“提升 (relabel)”操作使流最终达成有效最大流。
- **复杂度**：典型的最高标号 (highest label) 推送实现能达 \(O(V^2 \sqrt{E})\) 或更好；在实践中往往也不差。

---

## 16. **MCMF in \(O(\text{poly}(V, E))\)**

- **MCMF**：Minimum Cost Maximum Flow（最小费用最大流）。
- **算法**：Successive Shortest Path、Cycle Canceling、Network Simplex 等都在多项式时间可解；常见实现复杂度如 \(O(\maxFlow \times E \log V)\) 等，都归入 \(O(\text{poly}(V,E))\)。

  [Network Simplex](https://codeforces.com/blog/entry/94190)

- **应用**：线路规划、运输调度、匹配问题（加权）等。

---

## 17. **Minimum Arborescence (directed MST) in \(O(E \log V)\)**

- **问题**：在有根有向图里，找出一棵以给定根 r 出发的“向外辐射最小生成树”（也称 Arborescence）。
- **算法**：最知名的是 **Edmond’s Algorithm**（又称 Chu–Liu/Edmonds），可以实现 \(O(E \log V)\) 或 \(O(V^2)\) 等不同版本。
- **用途**：当需要选出一组最小代价有向边以保证可达性时。

---

## 18. **Suffix Tree**

- **用途**：在字符串中快速完成各种查询，如任意子串出现次数、最长重复子串、子串匹配等。
- **结构**：一棵紧凑前缀树 (trie) 的变体，将所有后缀插入到一棵树中，并使用边压缩。
- **构造算法**：Ukkonen 算法能在线性时间内构造 Suffix Tree；另一种常用结构是 Suffix Array + LCP，但它只需 \(O(n \log n)\) 容易实现。
- **应用**：字符串模式匹配、基因序列分析、数据压缩等。

---

## 19. **Online Convex Hull in 2D**

- **问题**：在 2D 平面上，点不断到来（在线），需要维持它们的凸包或支持查询当前凸包。
- **算法**：常见做法是“动态维护”凸包数据结构，如平衡树 + 分割线；也有一些更复杂的结构（如 Li-Chao-like 但在二维很 tricky）。
- **注意**：在线维护凸包比一次性构造（如 Graham Scan, Andrew’s Monotone Chain）要难得多。

---

## 20. **Convex Hull in 3D**

- **问题**：给定三维空间中的点集，求其凸包（多面体）。
- **算法**：常见构造法如 Quickhull、Gift Wrapping (也叫 Jarvis March 类比)、Incremental 等；若点数是 n，一般可以做到 \(O(n \log n)\) 或 \(O(n^2)\) 视实现而定。
- **难点**：3D 情况下要维护面、边、顶点的拓扑结构，相对 2D 多了不少细节。

---

## 21. **Halfplane Intersection**

- **概念**：给定平面上若干条线，每条线确定一个半平面，找所有半平面的交集区域（若区域非空就是一个凸多边形）。
- **算法**：主流是“半平面交”做法：将这些直线按极角排序，再用一个双端队列 (deque) 依次做切割，复杂度 \(O(n)\) 或 \(O(n \log n)\) 取决于是否需排序。
- **应用**：几何中判断公共可行区域，多边形裁剪，线性规划 (LP) 的几何视角等。

---

## 22. **Voronoi Diagram / Delaunay Triangulation**

- **Voronoi Diagram**：给定一组点（站点），对平面进行划分，使得每个点所在的“胞腔”中的所有位置与该点距离最小。
- **Delaunay Triangulation**：与 Voronoi 对偶的一种三角剖分，使得生成的三角形尽量不含钝角；应用于插值、网格生成、最近邻查询等。
- **算法**：可以在 \(O(n \log n)\) 时间构造（Divide & Conquer、Fortune’s Sweep Line 等）。
- **用途**：计算几何、地图、游戏 AI、数值分析等。

---

## 23. **操作 Formal Power Series（FPS）如 exp、log、sqrt 等**

- **背景**：形式幂级数 \(f(x) = \sum\_{n=0}^{\infty} a_n x^n\)，在竞赛中常用多项式或幂级数来表示组合恒等式/生成函数。
- **Newton 迭代法**：用类似牛顿法的手段来计算多项式的 `exp(f)`, `log(f)`, `sqrt(f)` 等，逐次逼近并截断到所需项数（一般到 \(x^n\) 的 \(n\)-次项）。
- **复杂度**：若做 FFT 卷积，常见实现到 \(O(n \log n)\) ~ \(O(n (\log n)^2)\) 之间。

---

## 24. **如何用 Generating Functions（GF）解题**

- **思路**：很多计数/组合问题可用生成函数“表达式”来代表解，然后通过展开/系数提取来求答案。
- **操作**：典型操作包括卷积（对应序列加法/组合）、幂次、指数、对数、差分等。
- **竞赛技巧**：常见于一些高级组合计数题，比如多重约束下的数目计算，可以把子问题转成多项式的乘法、复合或拉格朗日反演。

---

## 25. **Lagrange Inversion Formula**

- **概念**：给定形式幂级数的函数反演公式，用来求 \([x^n]F^{-1}(x)\) 或类似的复合函数系数。
- **在组合中的用法**：当你有 \(y = x \cdot f(y)\) 这种函数关系时，想求 \(y\) 的幂级数展开，可以用拉格朗日反演做快速系数提取。
- **难点**：公式本身不长，但要熟练应用需要对 GF 思想和分式展开有较好掌握。

---

## 26. **Elegia 的微分魔法 / 子集卷积微分魔法**

- **背景**：Elegia 是一位算法竞赛爱好者/博主，在讨论高阶技巧（特别是多项式、子集卷积、FFT、微分等）时提出一些巧妙公式。
- **大意**：通过对多项式(或子集卷积函数)做微分、再通过一些代数操作（如将结果乘上某个因子、再积分）简化卷积或反演的过程，从而加速某些 DP/计数问题。
- **使用门槛**：需要对多项式操作、卷积变换、生成函数非常熟悉，这些“魔法”技巧通常非常特定场景下奏效。

---

## 27. **Sweepline Mo**

- **Mo’s Algorithm**：主要处理数组上区间查询，通过块状分解 + 排序实现离线查询。
- **Sweepline Mo**：有时是指“在二维(或更高维)上，通过类似 Mo’s sqrt decomposition + sweep line”来处理更复杂的离线问题。例如“把一个维度看作时间或排序维度，另一个维度做离线分块”。
- **原理**：仍是离线 + 重排查询顺序来减少搬移成本，但在几何或其他多维情境下难度更高。

---

## 28. **Matroid Intersection**

- **概念**：母题是“给定两个独立性系统（两个拟阵），寻找它们共同的最大独立集”。拟阵 (matroid) 是一类具有可交换性、增量性质的抽象结构，广泛应用于贪心算法理论。
- **算法**：有多项式时间算法 (Edmonds)，比如在图中找最大生成森林或最大匹配都能视为某些拟阵问题。
- **实际意义**：在更复杂的约束组合优化问题中可用；竞赛中如果遇到“两个条件都必须满足的独立集”，可以尝试拟阵交的思路。

---

## 29. **“How Elegia's mind works”**

- **背景**：这是一句玩笑话，表达对某些算法大神（如 Elegia）的“天马行空思路”感到佩服，想知道他是如何灵光一现、推导出种种高阶花式公式。
- **现实**：多看、多想、多做题，站在数学/组合理论的基础之上，有时能产出一些巧妙的构造或变形；这是大量积累和天赋的结合。

---
