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

O(n^2)加边,O(1)查询
[2642. 设计可以求最短路径的图类-更新边的 floyd](<2642. 设计可以求最短路径的图类-更新边的floyd.go>)
