/**
 * @description 寻找桥(强连通)
 */
const criticalConnections = (n: number, connections: number[][]) => {
  const adjList = new Map<number, number[]>()
  const visited = new Set<number>()
  const bridge: [number, number][] = []
  const order = Array<number>(n).fill(Infinity)
  const lower = Array<number>(n).fill(Infinity)
  let step = 0
  connections.forEach(([pre, next]) => {
    adjList.set(pre, adjList.get(pre)?.concat([next]) || [next])
  })

  const dfs = (cur: number, pre: number) => {
    visited.add(cur)
    order[cur] = step
    lower[cur] = step
    step++
    for (const next of adjList.get(cur) || []) {
      if (next === pre) continue
      if (!visited.has(next)) {
        dfs(next, cur)
        // dfs回溯
        // 不问pre节点
      }
      lower[cur] = Math.min(lower[cur], lower[next])
      order[cur] < lower[next] && bridge.push([cur, next])
    }
  }
  for (const v of adjList.keys()) {
    !visited.has(v) && dfs(v, v)
  }
  console.log(lower, order, adjList)

  return bridge
}

console.dir(
  criticalConnections(5, [
    [1, 0],
    [2, 0],
    [3, 2],
    [4, 2],
    [4, 3],
    [3, 0],
    [4, 0],
  ]),
  { depth: null }
)

export {}
