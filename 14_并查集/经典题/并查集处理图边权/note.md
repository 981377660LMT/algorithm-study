# 并查集处理边权：

先把所有点 union 好，再对每条边的两个端点找到对应 leader，再处理这条边，不要在 merge 过程中处理边权。

```py
uf = UnionFindArraySimple(n)
for u, v, _ in edges:
    uf.union(u, v)
groupValue = [-1] * n
for u, v, w in edges:
    groupValue[uf.find(u)] = ... # 将边权w加入到对应的分组中
```
