`visited 最好在进入时就添加`
next 时只需要看在不在 visited 就行
这是效率最高的
保险起见**进来和 next 都添加**

```JS
function dfs(cur: number, visited: Set<number>): void {
    // 1. 添加
    visited.add(cur)

    for (const next of adjList[cur]) {
      // 2. 检查
      if (visited.has(next)) continue
      visited.add(next)
      dfs(next, visited)
    }
  }
```

# NOTE: 对于 n\*m 的网格图，BFS 最多只占用 O(min(n,m)) 的空间，而 DFS 最多会占用 O(nm) 的空间
