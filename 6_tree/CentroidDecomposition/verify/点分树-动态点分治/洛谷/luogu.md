https://www.luogu.com.cn/training/78277

- 对任意两点 u,v，
  dist(u,v) = dist(u,lca) + dist(v,lca)
  u,v 在点分树上的 lca 一定在 u,v 的路径上
- 它的高度与点分治的深度一样，只有 logn 级别

---

note: 没有必要学习，用处小
