这类问题的做法是**预处理 queries**
dfs 遍历到每个点处,取出当前结点对应的 query 进行处理

```Python
queries = [[] for _ in range(n)]  # !记录每个结点处的查询(qi,qv)
```
