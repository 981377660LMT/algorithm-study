![树的重心](../../images/447bf9065511f5b136989797fc52f50f60c64eaf5b265a87c3fbf12f25c66860.png)

模板：

```JS
function findCentre(n: number, edges: [next: number, weight: number][][]): number[] {
  const res: number[] = []
  // 最大子树大小
  const maxSizeOfSubtree = Array<number>(n).fill(Infinity)
  // 树的大小
  const treeSize = Array<number>(n).fill(Infinity)

  function dfs(cur: number, parent: number) {
    treeSize[cur] = 1
    maxSizeOfSubtree[cur] = 0

    for (const [next, _] of edges[cur]) {
      if (next === parent) continue
      dfs(next, cur)
      // dfs更新阶段
      // 后序,更新cur:此时cur可以拿到各个next的信息，更新cur
      maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], treeSize[next])
      treeSize[cur] += treeSize[next]
    }

    // dfs结算阶段:准备离开cur了，最后检查cur是否合法
    maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], n - treeSize[cur])
    if (maxSizeOfSubtree[cur] <= n / 2) res.push(cur)
  }
  dfs(0, -1)

  console.log(treeSize, maxSizeOfSubtree)
  return res
}
```

后序 dfs 的两个地方：`要在脑海里想象递归树`

```JS
function dfs(cur: number, parent: number) {
    treeSize[cur] = 1
    maxSizeOfSubtree[cur] = 0

    for (const [next, _] of edges[cur]) {
      if (next === parent) continue
      dfs(next, cur)

      // 1....
      // 根据next的信息更新cur
    }

    // 2...
    // 准备撤离cur，结算cur
  }
```

1 那个地方是 cur 用他下面相邻的 next 结点信息更新自己
2 那个地方是 dfs 准备永远离开 cur 结点 做最后的收尾工作
