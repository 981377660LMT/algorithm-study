**换根 dp(全方位木 DP)**
换根 DP，又叫二次扫描，是树形 DP 的一种。

其相比于一般的树形 DP 具有以下特点：

- `以树上的不同点作为根，其解不同`。
- 故为求解答案，不能单求某点的信息，`需要求解每个节点的信息`。
- 故无法通过一次搜索完成答案的求解，因为一次搜索只能得到一个节点的答案。

**换根 dp 的技巧**

1. 指定某个节点为根节点。

2. 第一次搜索完成预处理（如子树大小等），同时得到该节点的解。
   eg:`子树更新父结点`向下的最远距离,求出根 0 的答案

   - `后序 dfs 先处理出每个节点的信息`

3. 第二次搜索进行换根的动态规划，由已知解的节点推出相连节点的解。
   eg:`非子树(父结点+兄弟结点)更新子树`向上的最远距离

   - `前序 dfs 换根求解其他位置的答案`

参考:
https://qiita.com/keymoon/items/2a52f1b0fb7ef67fb89e
https://zhuanlan.zhihu.com/p/348349531
https://algo-logic.info/tree-dp/

ps:
`什么是 monoid`

> monoid 是一个二元运算和一个单位元组成的集合，满足结合律和单位元的性质。

---

1. 注意换根 dp 求出来的答案**不包含根自己**，如果求顶点的话最后需要加上自己的贡献.
2. 换根 dp 的视角始终是 0 号根节点的树，只要想着解决一个问题就好了，其实也就是树形 dp。
