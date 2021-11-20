const useUnionFindArray = (size: number) => {
  const parent = Array.from<number, number>({ length: size }, (_, i) => i)

  const find = (key: number) => {
    while ((parent[key] !== undefined && parent[key]) !== key) {
      // 进行路径压缩
      parent[key] = parent[parent[key]]
      key = parent[key]
    }
    return key
  }

  const union = (key1: number, key2: number) => {
    const root1 = find(key1)
    const root2 = find(key2)
    if (root1 === root2) return
    // 简单rank优化:总是让大的根指向小的根
    parent[Math.max(root1, root2)] = Math.min(root1, root2)
  }

  const isConnected = (key1: number, key2: number) => find(key1) === find(key2)

  return { union, find, isConnected }
}

/**
 *
 * @param m
 * @param n
 * @param positions
 * 给你一个大小为 m x n 的二进制网格 grid
 * 0 表示水，1 表示陆地
 * 其中 positions[i] = [ri, ci] 是要执行第 i 次操作的位置 (ri, ci) 。
 * 操作将某个位置的水转换成陆地
 * 返回第i次转换为陆地后，地图中岛屿的数量
 * @summary
 * 1.初始化并查集和visited陆地集合
   2.每新增一个陆地就把它加入并查集和陆地集合中
   3.对新加的陆地，连接它和四周的陆地
   4.获取并查集里的联通分量
 */
function numIslands2(m: number, n: number, positions: number[][]): number[] {
  const res: number[] = []
  const visited = new Set<number>()
  const uf = useUnionFindArray(m * n)
  let count = 0

  for (const [x, y] of positions) {
    // 重复的岛屿
    if (visited.has(x * n + y)) {
      res.push(count)
      continue
    }

    visited.add(x * n + y)
    count++

    for (const [dx, dy] of [
      [0, 1],
      [0, -1],
      [1, 0],
      [-1, 0],
    ]) {
      const [nextX, nextY] = [x + dx, y + dy]
      if (nextX >= 0 && nextX < m && nextY >= 0 && nextY < n && visited.has(nextX * n + nextY)) {
        if (uf.isConnected(x * n + y, nextX * n + nextY)) continue
        uf.union(x * n + y, nextX * n + nextY)
        count--
      }
    }

    res.push(count)
  }

  return res
}

console.log(
  numIslands2(3, 3, [
    [0, 0],
    [0, 1],
    [1, 2],
    [2, 1],
  ])
)
// 输出：[1,1,2,3]
export {}
