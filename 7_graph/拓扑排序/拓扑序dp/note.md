**有向无环图 (DAG) 存在拓扑序**
（DAG, Directed Acyclic Graph）

`5300. 有向无环图中一个节点的所有祖先`
`Acwing 164. 可达性统计`

**拓扑图(DAG)的最短路可以用拓扑排序(dp) O(n)求出**
[最短路](<2297.%20Jump%20Game%20IX-%E6%8B%93%E6%89%91%E5%9B%BE%E6%9C%80%E7%9F%AD%E8%B7%AFO(n).py>)
[最长路](2050.%20%E5%B9%B6%E8%A1%8C%E8%AF%BE%E7%A8%8B%20III-%E6%8B%93%E6%89%91%E6%8E%92%E5%BA%8F%E6%9C%80%E9%95%BF%E8%B7%AF%E5%BE%84-dp.py)

```Python
queue = deque([0])
dist = defaultdict(lambda: int(1e20), {0: 0})
while queue:
    cur = queue.popleft()
    for next in adjMap[cur]:
        dist[next] = min(dist[next], dist[cur] + adjMap[cur][next])
        deg[next] -= 1
        if deg[next] == 0:
            queue.append(next)

return dist[n - 1]
```
