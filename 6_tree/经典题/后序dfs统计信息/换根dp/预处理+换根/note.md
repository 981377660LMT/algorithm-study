- !换根 dp 不一定要 Rerooting 框架做，可以手动换根，灵活利用 dfs 序+更强大的数据结构做.
  对每条边(u,v)，预处理出以 u 为根时 v 子树内的答案,以 v 为根时 u 子树内的答案 `edgeInfo`.
  计算出 0 为根时的答案，减去 0->node 之间的边的贡献，加上边 node->0 之间的边的贡献，就是以 node 为根的答案.
