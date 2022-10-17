这类问题的做法是

1. **预处理 queryGroup**,即先记录下每个点处的所有查询
2. dfs 遍历到每个点处,取出当前结点对应的 query 进行处理

```Python
queryGroup = [[] for _ in range(n)]  # !记录每个结点处的查询(qi,qv)
for qi, (node, qv) in enumerate(queries):
   queryGroup[node].append((qi, qv))


# dfs过程中
for qi, qv in queryGroup[cur]:
    res[qi] = trie.search(qv)
```
