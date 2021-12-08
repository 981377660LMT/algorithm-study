**直接取出 id** 省去了对 query 做 map 的处理

```Python
按顺序处理query
for id in sorted(range(n), key=lambda x: -queries[x][0]):
    ...
    query[id] ...
    res[id] = ...
    ...
return res
```

离线查询排序的作用:

1. 题目要求，逐步更新
2. 后面的查询可以命中缓存
