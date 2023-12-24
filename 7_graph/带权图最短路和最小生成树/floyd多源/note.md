floyd 剪枝(稀疏图优化)

```python
for k in range(n):
    for i in range(n):
        if dist[i][k] == INF:  # !剪枝
            continue
        for j in range(n):
            cand = dist[i][k] + dist[k][j]
            dist[i][j] = cand if dist[i][j] > cand else dist[i][j]
```
