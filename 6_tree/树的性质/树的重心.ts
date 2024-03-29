// 计算以无根树每个点为根节点时的最大子树大小，这个值最小的点称为无根树的重心。
// 1.某个点是树的重心等价于它最大子树大小不大于整棵树大小的一半。
// 2.树至多有两个重心。如果树有两个重心，那么它们相邻。此时树一定有偶数个节点，且可以被划分为两个大小相等的分支，每个分支各自包含一个重心。
// 3.树中所有点到某个点的距离和中，到重心的距离和是最小的；如果有两个重心，那么到它们的距离和一样。反过来，距离和最小的点一定是重心。

// 在 DFS 中计算每个子树的大小，记录“向下”的子树的最大大小，利用总点数 - 当前子树（这里的子树指有根树的子树）的大小得到“向上”的子树的大小
// 利用性质1，dfs即可

function findCentroids(n: number, edges: [cur: number, next: number][], root = 0): number[] {
  const res: number[] = []
  const weight = new Uint32Array(n + 1) // 最大连通块大小,即此节点为割点分割之后两半的最大大小
  const subSize = new Uint32Array(n + 1) // 子树的大小,即向`下面`走可以到多少个结点

  const adjList: number[][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v]) => {
    adjList[u].push(v)
    adjList[v].push(u)
  })

  dfs(root, -1)
  return res

  function dfs(cur: number, parent: number): void {
    subSize[cur] = 1

    for (const next of adjList[cur]) {
      if (next === parent) continue
      dfs(next, cur)
      // 后序,更新cur:此时cur可以拿到各个next的信息
      subSize[cur] += subSize[next]
      weight[cur] = Math.max(weight[cur], subSize[next])
    }

    // cur准备回退了，检查cur是否合法
    weight[cur] = Math.max(weight[cur], n - subSize[cur])
    if (weight[cur] <= n / 2) res.push(cur)
  }
}

export {}

if (require.main === module) {
  console.log(
    findCentroids(4, [
      [1, 0],
      [1, 2],
      [1, 3]
    ])
  )
}
