// 1 <= pairs.length <= 105
// 请你返回 任意一个 pairs 的合法重新排列。
// 数据保证至少存在一个 pairs 的合法重新排列。
/**
 *
 * @param pairs
 * @returns
 * @summary 有向图的欧拉路径
 */
function validArrangement(pairs: number[][]): number[][] {
  const adjMap = new Map<number, number[]>()
  buildGraph()

  const start = getStart(pairs)
  const eulerLoop = getEulerPath(adjMap, start)

  const res: number[][] = []
  let i = 0
  while (res.length < pairs.length) {
    res.push([eulerLoop[i], eulerLoop[i + 1]])
    i++
  }

  return res

  function buildGraph() {
    for (const [u, v] of pairs) {
      !adjMap.has(u) && adjMap.set(u, [])
      adjMap.get(u)!.push(v)
    }
  }

  // 题目给定的图一定满足以下二者之一：
  // 所有点入度等于出度；
  // 恰有一个点出度 = 入度 + 1（欧拉路径的起点），且恰有一个点入度 = 出度 + 1（欧拉路径的终点），其他点入度等于出度。
  function getStart(pairs: number[][]): number {
    const outdegree = new Map<number, number>()
    const indegree = new Map<number, number>()

    for (const [u, v] of pairs) {
      outdegree.set(u, (outdegree.get(u) || 0) + 1)
      indegree.set(v, (indegree.get(v) || 0) + 1)
    }

    const oddStartPoint = [...new Set(pairs.flat())].filter(
      key => (outdegree.get(key) || 0) - (indegree.get(key) || 0) === 1
    )

    if (oddStartPoint.length > 0) return oddStartPoint[0]
    else return pairs[0][0]
  }

  function getEulerPath(adjMap: Map<number, number[]>, start: number): number[] {
    let cur = start
    const stack: number[] = [cur]
    const res: number[] = []

    while (stack.length > 0) {
      if (adjMap.has(cur) && adjMap.get(cur)!.length > 0) {
        stack.push(cur)
        const next = adjMap.get(cur)!.pop()!
        cur = next
      } else {
        res.push(cur)
        cur = stack.pop()!
      }
    }

    return res.reverse()
  }
}

// console.log(
//   validArrangement([
//     [5, 1],
//     [4, 5],
//     [11, 9],
//     [9, 4],
//   ])
// )
// console.log(
//   validArrangement([
//     [1, 3],
//     [3, 2],
//     [2, 1],
//   ])
// )
console.log(
  validArrangement([
    [13, 6],
    [17, 13],
    [8, 11],
    [1, 19],
    [16, 6],
    [19, 0],
    [3, 4],
    [11, 9],
    [5, 3],
    [9, 15],
    [6, 15],
    [14, 10],
    [2, 1],
    [6, 2],
    [4, 8],
    [0, 5],
    [15, 16],
    [10, 17],
  ])
)
// 解释：
// 输出的是一个合法重新排列，因为每一个 endi-1 都等于 starti 。
// end0 = 9 == 9 = start1
// end1 = 4 == 4 = start2
// end2 = 5 == 5 = start3
