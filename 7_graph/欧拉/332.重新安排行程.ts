/**
 * @param {string[][]} tickets
 * @return {string[]}
 * 行程必须从 JFK 开始。如果存在多种有效的行程，请你按字典排序返回最小的行程组合。
 * 例如，行程 ["JFK", "LGA"] 与 ["JFK", "LGB"] 相比就更小，排序更靠前。
 * 假定所有机票至少存在一种合理的行程。且所有的机票 必须都用一次 且 只能用一次。(欧拉路径)
 * @summary Hierholzer算法(插入回路法)
 */
const findItinerary = function (tickets: string[][]): string[] {
  const res: string[] = []
  const adjMap = new Map<string, string[]>()

  // 字典序小的排在后边，遍历时先pop出来
  for (const [v, w] of tickets.sort().reverse()) {
    if (!adjMap.has(v)) adjMap.set(v, [])
    adjMap.get(v)!.push(w)
  }

  // 求欧拉路径的算法
  let cur = 'JFK'
  const stack: string[] = ['JFK']
  while (stack.length) {
    if (adjMap.get(cur)?.length) {
      stack.push(cur)
      cur = adjMap.get(cur)?.pop()!
    } else {
      res.push(cur)
      cur = stack.pop()!
    }
    console.log(stack, res)
  }

  console.log(adjMap)
  return res.reverse()
}

console.log(
  findItinerary([
    ['JFK', 'SFO'],
    ['JFK', 'ATL'],
    ['SFO', 'ATL'],
    ['ATL', 'JFK'],
    ['ATL', 'SFO'],
  ])
)

export {}
