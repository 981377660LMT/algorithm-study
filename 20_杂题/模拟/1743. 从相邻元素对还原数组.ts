// 存在一个由 n 个不同元素组成的整数数组 nums ，但你已经记不清具体内容。好在你还记得 nums 中的每一对相邻元素。
// 返回 原始数组 nums 。如果存在多种解答，返回 其中任意一个 即可。
function restoreArray(adjacentPairs: number[][]): number[] {
  const n = adjacentPairs.length + 1
  const res: number[] = []
  const adjMap = new Map<number, number[]>()
  const visited = new Set<number>()

  for (const [cur, next] of adjacentPairs) {
    !adjMap.has(cur) && adjMap.set(cur, [])
    !adjMap.has(next) && adjMap.set(next, [])
    adjMap.get(cur)!.push(next)
    adjMap.get(next)!.push(cur)
  }

  let start = Infinity
  for (const [cur, nexts] of adjMap.entries()) {
    if (nexts.length === 1) {
      start = cur
      break
    }
  }

  while (res.length < n) {
    res.push(start)
    visited.add(start)
    for (const next of adjMap.get(start) || []) {
      if (visited.has(next)) continue
      start = next
    }
  }

  return res
}

console.log(
  restoreArray([
    [4, -2],
    [1, 4],
    [-3, 1],
  ])
)
// 输出：[1,2,3,4]
// 解释：数组的所有相邻元素对都在 adjacentPairs 中。
// 特别要注意的是，adjacentPairs[i] 只表示两个元素相邻，并不保证其 左-右 顺序。
