const useUnionFindMap = () => {
  let count = 0 // 一开始的联通分量个数
  const parent = new Map<number, number>()

  const add = (key: number) => {
    if (parent.has(key)) return
    parent.set(key, key)
    count++
  }

  const find = (val: number) => {
    while (parent.get(val) !== val) {
      val = parent.get(val)!
    }

    return val
  }

  const union = (key1: number, key2: number) => {
    if (isConnected(key1, key2)) return
    const root1 = find(key1)
    const root2 = find(key2)
    if (root1 === root2) return
    // rank优化:总是让大的根指向小的根
    parent.set(Math.max(root1, root2), Math.min(root1, root2))
    count--
  }

  const isConnected = (key1: number, key2: number) => find(key1) === find(key2)

  const getCount = () => count

  return { add, union, find, isConnected, getCount }
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
  const island = new Set<number>()
  const uf = useUnionFindMap()

  for (const [x, y] of positions) {
    uf.add(x * n + y)
    island.add(x * n + y)

    for (const [dx, dy] of [
      [0, 1],
      [0, -1],
      [1, 0],
      [-1, 0],
    ]) {
      const [nextI, nextY] = [x + dx, y + dy]
      if (nextI >= 0 && nextI < m && nextY >= 0 && nextY < n && island.has(nextI * n + nextY))
        uf.union(x * n + y, nextI * n + nextY)
    }

    res.push(uf.getCount())
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
