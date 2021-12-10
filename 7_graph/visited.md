`visited 最好在进入时就添加`
next 时只需要看在不在 visited 就行
这是效率最高的

```JS
function dfs(cur: number, visited: Set<number>): void {
    // 1. 添加
    visited.add(cur)

    for (const next of adjList[cur]) {
      // 2. 检查
      if (visited.has(next)) continue
      dfs(next, visited)
    }
  }
```
