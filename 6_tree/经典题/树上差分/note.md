# 两种类型

1. 静态：预处理 distToRoot 数组，dfs 出每个点到根结点的值
2. 动态：树链上维护阿贝尔群
   点的修改:u++,v++,lca--,fa[lca]--
   边的修改:u++,v++,lca-=2

---

https://www.luogu.com.cn/blog/LawrenceSivan/shu-shang-ci-fen-zong-jie

点差分：https://www.luogu.com.cn/problem/P3128
边差分：https://www.luogu.com.cn/problem/P2680
