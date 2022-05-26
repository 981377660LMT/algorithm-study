// 带weight(每个根节点的权重)的并查集
const useUnionFindArray = (size: number) => {
  let count = size
  const parent = Array.from<number, number>({ length: size }, (_, i) => i)
  const weight = Array<number>(size).fill(1)

  const find = (key: number) => {
    while (parent[key] !== key) {
      // 进行路径压缩
      parent[key] = parent[parent[key]]
      key = parent[key]
    }
    return key
  }

  const union = (key1: number, key2: number) => {
    let root1 = find(key1)
    let root2 = find(key2)
    if (root1 === root2) return
    if (weight[root1] > weight[root2]) {
      ;[root1, root2] = [root2, root1]
    }
    parent[root1] = root2
    weight[root2] += weight[root1]
    count--
  }

  const isConnected = (key1: number, key2: number) => find(key1) === find(key2)

  const getCount = () => count

  return { union, find, isConnected, getCount, weight }
}

// 给你一个大小为 n x n 二进制矩阵 grid 。最多 只能将一格 0 变成 1 。
// 返回执行此操作后，grid 中最大的岛屿面积是多少？
// 岛屿 由一组上、下、左、右四个方向相连的 1 形成。
function largestIsland(grid: number[][]): number {
  const [m, n] = [grid.length, grid[0].length]
  const uf = useUnionFindArray(m * n)

  // 初始化并查集
  for (let x = 0; x < m; x++) {
    for (let y = 0; y < n; y++) {
      if (grid[x][y] === 1) {
        for (const [dx, dy] of [
          [-1, 0],
          [0, -1],
        ]) {
          const [nextX, nextY] = [x + dx, y + dy]
          if (nextX >= 0 && nextX < m && nextY >= 0 && nextY < n && grid[nextX][nextY] === 1) {
            uf.union(x * n + y, nextX * n + nextY)
          }
        }
      }
    }
  }

  // 当前最大的区域
  let res = Math.max(...uf.weight, 1)
  for (let x = 0; x < m; x++) {
    for (let y = 0; y < n; y++) {
      // 把一块海洋，变成陆地
      if (grid[x][y] === 0) {
        let add = 1
        const visitedRoot = new Set<number>()

        for (const [dx, dy] of [
          [0, 1],
          [0, -1],
          [1, 0],
          [-1, 0],
        ]) {
          const [nextX, nextY] = [x + dx, y + dy]
          if (nextX >= 0 && nextX < m && nextY >= 0 && nextY < n && grid[nextX][nextY] === 1) {
            const parentRoot = uf.find(nextX * n + nextY)
            if (visitedRoot.has(parentRoot)) continue
            visitedRoot.add(parentRoot)
            add += uf.weight[parentRoot]
          }
        }

        res = Math.max(res, add)
      }
    }
  }

  return res
}

console.log(
  largestIsland([
    [1, 0],
    [0, 1],
  ])
)
// 输出: 3
// 解释: 将一格0变成1，最终连通两个小岛得到面积为 3 的岛屿。

export {}
