// 2 <= n <= 15
// 请你返回一个大小为 n-1 的数组，其中第 d 个元素（下标从 1 开始）是城市间 最大距离 恰好等于 d 的`子树`数目。
// 1.求每个点到所有点的最短距离--多源最短路径算法 floyd
// 2.二进制枚举所有状态（子集）
function countSubgraphsForEachDiameter(n: number, edges: number[][]): number[] {
  // 构建dist矩阵
  const dist = Array.from<number, number[]>({ length: n }, () => Array(n).fill(Infinity))

  for (const [u, v] of edges) {
    dist[u - 1][v - 1] = 1 // 相邻的边权重为1
    dist[v - 1][u - 1] = 1 // 相邻的边权重为1
  }

  for (let k = 0; k < n; k++) {
    for (let i = 0; i < n; i++) {
      for (let j = 0; j < n; j++) {
        dist[i][j] = Math.min(dist[i][j], dist[i][k] + dist[k][j])
      }
    }
  }

  const res = Array<number>(n).fill(0)

  // 枚举子集
  for (let state = 0; state < 1 << n; state++) {
    let maxDist = 0
    let edges = 0

    for (let u = 0; u < n; u++) {
      for (let v = u + 1; v < n; v++) {
        if (((state >> u) & 1) === 1 && ((state >> v) & 1) === 1) {
          // 判断是不是边
          if (dist[u][v] === 1) edges++
          maxDist = Math.max(maxDist, dist[u][v])
        }
      }
    }

    const vertex = countOne(state)
    if (isTree(edges, vertex)) {
      res[maxDist]++
    }
  }

  return res.slice(1)

  function countOne(n: number) {
    let res = 0

    while (n > 0) {
      res++
      n &= n - 1
    }

    return res
  }

  function isTree(edges: number, vertex: number) {
    return vertex === edges + 1
  }
}

console.log(
  countSubgraphsForEachDiameter(4, [
    [1, 2],
    [2, 3],
    [2, 4],
  ])
)
// 输出：[3,4,0]
