# 两种类型

- 静态，配合 dfs；动态，配合 bit：
  预处理 distToRoot 数组，dfs 出每个点到根结点的值
  树链上维护阿贝尔群
  `RangeAdd，pointGet`，链上加单点求职参考树状数组的差分 api

  点的修改:u++,v++,lca--,fa[lca]--
  边的修改:u++,v++,lca-=2
  TODO:设计成一个 Wrapper 或 func，因为大多场合这都是其中一步

---

https://www.luogu.com.cn/blog/LawrenceSivan/shu-shang-ci-fen-zong-jie

点差分：https://www.luogu.com.cn/problem/P3128
边差分：https://www.luogu.com.cn/problem/P2680

---

https://zhuanlan.zhihu.com/p/30677133

---

# 树链求并。

按照 dfn 排序后，对每个 i，将 ui 到根节点的链上点+1，将 lca(u[i],u[i+1])到根节点的链上点-1。
现在问题转化为了：" 链上加 " & " 单点求值 "。
可以使用树上差分将问题转化为：" 单点加 " & " 子树求和 "。

树上差分解决问题类型
