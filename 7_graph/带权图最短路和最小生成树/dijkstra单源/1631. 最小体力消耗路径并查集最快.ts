// 将每个格子抽象成图中的一个点；
// 将每两个相邻的格子之间连接一条边，
// 长度为这两个格子本身权值的差的绝对值；
// 需要找到一条从左上角到右下角的「最短路径」，

import { useUnionFindArray } from '../../../14_并查集/useUnionFind'

// 其中路径的长度定义为路径上所有边的权值的最大值。

// 我们可以将所有边按照长度进行排序并依次添加进并查集，直到左上角和右下角连通为止
function minimumEffortPath(heights: number[][]): number {
  const m = heights.length
  const n = heights[0].length

  const edges: [number, number, number][] = []
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      const key = i * n + j
      if (i + 1 < m) edges.push([Math.abs(heights[i + 1][j] - heights[i][j]), key, key + n])
      if (j + 1 < n) edges.push([Math.abs(heights[i][j + 1] - heights[i][j]), key, key + 1])
    }
  }
  edges.sort((a, b) => a[0] - b[0])

  const uf = useUnionFindArray(m * n)
  for (const edge of edges) {
    uf.union(edge[1], edge[2])
    if (uf.isConnected(0, m * n - 1)) {
      return edge[0]
    }
  }

  return 0
}

console.log(
  minimumEffortPath([
    [1, 2, 2],
    [3, 8, 2],
    [5, 3, 5],
  ])
)
