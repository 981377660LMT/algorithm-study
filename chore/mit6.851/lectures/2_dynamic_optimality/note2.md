以下是一份对 **MIT 6.851: Advanced Data Structures (Spring 2012) Lecture 6** (_Dynamic Optimality of BSTs, part 2_) 的**详细讲解与总结**。本讲延续了上一讲（Lecture 5）的主题：**在几何视角下研究二叉搜索树（BST）的动态最优性**，并聚焦于**利用独立矩形（Independent Rectangles）技术推导BST访问的下界**。另外，还介绍了 **Tango Trees** (一种 \(O(\log\log n)\)-competitive 的自适应 BST)，以及 **Signed Greedy** 算法如何在几何模型中提供下界。

---

## 1. 回顾与概览

在 Lecture 5 中，我们看到：

1. **BST 与 Arborally Satisfied 点集的对应**

   - 对给定访问序列 \((x_1,\ldots,x_m)\)，在 \((\text{key},\text{time})\) 平面上放置输入点 \((x_i, i)\)。
   - 如果某 BST 算法在时刻 \(i\) 访问了节点\(k\)，就对 \((k,i)\) 放置一个点。
   - 该点集满足“任意不在同一行或列的两点所构成的矩形都包含第三点”时，称为 Arborally Satisfied。正好对应一个合法的 BST 执行。
   - **寻找最小 Arborally Satisfied 点集**与“寻找最佳 BST 执行”等价。

2. **离线 Greedy** 构造
   - 逐行（自下而上）“填充”必要点，使点集保持 Arborally Satisfied；
   - 猜想这种贪心构造可在常数因子上逼近最优点集，从而若能被在线模拟，也会给出一个常数竞争的动态最优 BST。
   - 但尚无完整证明。

**本讲**：主要讨论**下界**。采用 “**独立矩形**(Independent Rectangles)” 的思想给出对任何 BST 执行（或等价的 Arborally Satisfied 集）的**信息论级**下界。

- 然后介绍了三种著名下界：Wilber 1、Wilber 2，以及“Signed Greedy”算法产生的下界。
- 并带出**Tango Trees**，一种达到了 \(O(\log\log n)\) 竞争比的在线自适应 BST。

---

## 2. 独立矩形法（Independent Rectangles）推导下界

### 2.1 定义与核心结论

**独立矩形**大致含义：在点集 \(\mathcal{P}\) 上，取两对点 \((a,b)\) 与 \((c,d)\) 所构成的两个矩形，若它们彼此不满足“会被第三点填充”（即都不 arborally satisfied），且其中一个矩形的任何点都不处于另一个矩形的内部（使之“相对独立”），则称这两个矩形是独立的。

- 用“\(+\) 矩形”指对角线正斜率（即 “\(/\)” 形）构成的矩形；
- 用“\(-\) 矩形”指对角线负斜率（即 “\(\backslash\)” 形）构成的矩形。

**结论** (Independent Rectangle Bounds)：

> 若某点集 \(\mathcal{P}\) 要被补充成 Arborally Satisfied 的最小超集 \(\mathrm{OPT}\)，则  
> \[
> |\mathrm{OPT}| \; \ge\; |\mathcal{P}| \;+\; \tfrac12 \,\bigl|\mathrm{MAXIND}\bigr|,
> \]  
> 其中 \(\mathrm{MAXIND}\) 是 \(\mathcal{P}\) 所能生成的最多独立矩形的子集大小。

这表示：如果能在 \(\mathcal{P}\) 中找到许多彼此独立的矩形，就逼迫我们在 Arborally Satisfied 的超集里必须额外添加足够多的点来填充这些矩形。

**因此**，对于 BST 执行而言，“\(\mathrm{OPT}\)” 对应最优执行的访问点数，\(\mathrm{MAXIND}\) 给出了一个**下界**。每个**独立矩形**“暗示”了至少要多加 1 个新点(从大局计要 \(\frac{1}{2}\) 的系数)，从而形成一个下界。

---

### 2.2 证明思路

Lecture 中给出了详细的证明大纲（以 + 矩形为例）。核心是一个“充电” (charging) 的方法：

1. **选取最大宽度的矩形**；观察所有与它相交的其它矩形，它们可分为共享顶点、共享底点或都不共享；彼此之间因为独立性产生某些严格的排序关系。这样能保证一个特定垂直分割线把矩形分割成左右两部分，各自“干净”互不冲突。
2. **在最优 arborally satisfied 超集中**，必须为此大矩形添加一个点来将之“填充”。实际会找出一对点 \((p, q)\) 在同一行，通过一次计费，把此矩形“费用”记在这对点上。因为这种匹配是双不交叉的，所以当我们计费给若干对点，就需要额外的新点满足这些对点之间的行间隙，最终推导出 \(\mathrm{OPT}\) 至少要添加那么多个点。

**综合**：得到 \(\mathrm{OPT} \ge |\mathcal{P}| + \tfrac12 |\mathrm{MAXIND}|\)。

---

## 3. Wilber 下界与其它下界技术

### 3.1 Wilber 的第二下界 (Wilber 2)

**Wilber [1]** 提出了两个著名的 BST 下界。**Wilber 2** 的定义方式如下：

1. 给定输入点集 \(\{(x_i, i)\}\)；
2. 对于每个点 \((x_i, i)\)，找出那些“能正交可见（orthogonally visible）且位于其下方”的点(也就是在矩形里不被别的点阻挡的那些点)；按其纵坐标排序；
3. 在这些“下方可见点”的排序中，统计左右方向的切换次数（left-right alternations）。
4. 对全部点 \((x_i, i)\) 的计数之和即“Wilber2 值”。

**大意**：如果访问序列中，某点的下方可见点们左右切换频繁，就暗示结构必须做大量“旋转”才能紧凑访问，这构成了对 BST 访问的一个下界。

**猜想**：\(\mathrm{OPT} = \Theta(\text{Wilber2})\)，意味着 Wilber2 完全刻画了最佳 BST 的访问代价。

- 到目前尚未证明或否定；这是 BST 动态最优猜想的核心一部分。

### 3.2 Key-independent Optimality

Iacono [2] 提出了一种“键无关(key independent)最优性”设定：如果键值完全是“无意义或随机的”，则期望最优访问代价就等价于“工作集上界(Working Set Bound)”。

- 这样可推出**Splay 树**在这种键无关情境下是最优的（因为它满足 Working Set Property）。
- 这为 Splay 树提供了另一种有限场景下的“最优性”支持。

### 3.3 Wilber 的第一下界 (Wilber 1)

- 给定固定形状的“下界树” \(P\)（往往是一个完整二叉树），把访问序列 \((x_1,\dots,x_m)\) 中的键都映射到 \(P\) 中。
- 对于 \(P\) 中每个节点 \(y\)，统计访问序列中“在其左子树与右子树之间切换的次数”总和；所有节点的切换次数之和就是 Wilber1。
- Wilber1 也为任意 BST 算法给出了一个下界：任意执行若要访问那些序列，也需要至少对应这么多旋转/指针行走代价。

**示例**：bit-reversal sequence (位反转序列) 在一棵完美二叉树中会导致 \(\Theta(n\log n)\) 的 Wilber1 值。

**开放**：是否对任意访问序列，都能找到一个形状 \(P\)，使得 \(\mathrm{OPT}=\Theta(\text{Wilber1})\)?

- 若答案是“是”，则 Wilber1 下界也能精确刻画 \(\mathrm{OPT}\)。

---

## 4. Tango Trees

### 4.1 动机

**Tango Trees** (Demaine, Harmon, Iacono, Patrascu [DHIP04]) 是已知最好的在线 BST 算法之一，有 **\(O(\log\log n)\)-competitive** 的性能。它说明了：

- 一般平衡树只能保证 \(O(\log n)\) 稳定查询，但没有竞争比意义（因为最优离线算法可能利用序列结构更快）。
- Tango Trees 则保证 **相对于 Wilber1**，只差一个 \(O(\log\log n)\) 的因子，是第一个打破 \(O(\log n)\) 竞争比的自适应 BST 结构。

### 4.2 核心思想

给定一个固定对手树 \(P\)（可能是完美二叉树），在执行过程中，我们为 \(P\) 中的若干**首选子指针**(preferred child) 构造“首选路径”(preferred path)。这些路径分段作为**辅助树**(auxiliary tree) 维护，保证当访问序列导致路径变化时，我们能在 \(O(\log\log n)\) 时间内更新。

- **首选子指针**：对每个节点 \(y\)，其首选子是最近被访问过的子树的一侧。访问序列会导致这些指针动态改变，形成“首选路径”的断开与重连。
- 每个辅助树包含至多 \(O(\log n)\) 个节点，对辅助树的搜索或分裂/合并操作都在 \(O(\log(\log n))\) 完成。
- 每访问一个节点，就可能有 \(k\) 次辅助树切换。若 Wilber1 下界给出每一步只有有限次优先子指针切换，那么整次访问就消耗 \(O(k\log\log n)\) 时间。

通过详细分析，Tango Tree 最终得到 \(O(\log \log n)\)-competitive 性能。

**结果**：这也说明 Wilber1 下界距离真最优 \(\mathrm{OPT}\) 可能存在 \(O(\log \log n)\) 的差异（至少对在线算法而言）。因为 Tango 树只能保证与 Wilber1 相差 \(O(\log \log n)\)，如果 Wilber1 与 \(\mathrm{OPT}\) 总是同阶，则这意味着 Tango 就实现了动态最优，否则 Wilber1 并未完全刻画 \(\mathrm{OPT}\)。

---

## 5. Signed Greedy 算法与下界

**Signed Greedy** 是对 Lecture 5 中“离线贪心（Greedy）”的一个变体，用来生成**下界**（而不是上界）。它只满足对 \(+\) 矩形或 \(-\) 矩形的一边进行贪心填充，而不管另一种斜线矩形从而引发新的未填矩形。

- **Signed Greedy** 不对应实际的 BST（因为只对半部分矩形类型进行满足），故它不是可行的执行方案，但它在几何上与独立矩形下界有着紧密联系。
- 该算法产生的点集规模是 \(|\mathcal{P}| + \bigl|\text{some half-rectangles}\bigr|\)，并与独立矩形下界 (MAXIND) 有个常数因子关系。

**最后结论**：

1. 普通 Greedy（上界，是真 BST）
2. Signed Greedy（下界，不是真 BST）
3. 独立矩形下界 (MAXIND)

这三者**处于一个常数因子的“夹逼”关系**。如果有证据表明 Signed Greedy 与普通 Greedy 相差常数因子，那么就说明**独立矩形下界与 Greedy 上界**也是常数因子，进而表明 Greedy 在常数竞争比上逼近最佳 BST 执行（解决了动态最优问题）。这是当代研究中的一个重要猜想。

---

## 6. 总结与展望

- **独立矩形法**提供了一个强有力的**下界**分析工具：  
  \[
  \text{OPT} \;\ge\; |\text{input}| \;+\; \tfrac12 \,\bigl|\mathrm{MAXIND}\bigr|.
  \]
- **Wilber 1** 与 **Wilber 2**：两种著名下界度量，对 BST 动态最优带来深刻洞见，仍有不少开放问题：
  - 是否 \(\mathrm{OPT} = \Theta(\text{Wilber2})\)？
  - 是否有对任意序列都能构造一棵静态下界树 \(P\) 使 Wilber1 = \(\mathrm{OPT}\)？
- **Key-independent Optimality**：在键值无意义时，Splay 等满足 Working Set 性质的数据结构可达期望最优。
- **Tango Tree**：通过首选路径+辅助树将在线竞争比降到 \(O(\log\log n)\)。
- **Signed Greedy**：仅满足 \(+\) 或 \(-\) 矩形的一半贪心，提供一个**几何下界**，与**独立矩形**等价的下界互相印证出常数因子差。

这些结果综合了对 BST 动态最优性的多视角理解；**动态最优**仍然是一个“世纪难题”之一。后续章节（Lecture 7 及以后）会继续阐述更多关于自适应 BST、Splay 猜想及相关进展。
