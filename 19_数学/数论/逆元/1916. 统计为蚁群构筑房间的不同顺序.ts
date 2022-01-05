import { comb } from './逆元求comb'

function waysToBuildRooms(prevRoom: number[]): number {
  const n = prevRoom.length
  const MOD = 1e9 + 7
  const adjMap = new Map<number, number[]>()
  for (let i = 0; i < n; i++) {
    !adjMap.has(prevRoom[i]) && adjMap.set(prevRoom[i], [])
    adjMap.get(prevRoom[i])!.push(i)
  }

  function dfs(cur: number): [nodeCount: number, orderCount: number] {
    if (!adjMap.has(cur)) return [1, 1]

    let curNodeCount = 0
    let curSortCount = 1

    for (const next of adjMap.get(cur)!) {
      const [subNodeCount, subSortCount] = dfs(next)
      curNodeCount += subNodeCount
      // 新的拓扑排序方案数为: 当前的拓扑排序方案数 * 从nodes个位置里选nodes_个位置分配给该子树 * 子树的拓扑排序方案数
      curSortCount =
        (((subSortCount * curSortCount) % MOD) * comb(curNodeCount, subNodeCount)) % MOD
    }

    return [curNodeCount + 1, curSortCount]
  }

  return Number(dfs(0)[1])
}

console.log(waysToBuildRooms([-1, 0, 0, 1, 2]))
