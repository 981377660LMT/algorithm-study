/**
 * @param {number} n
 * @param {number[][]} connections
 * @return {number[][]}
 */
var criticalConnections = function (n, connections) {
  const adjList = new Map()
  const visited = new Set()
  const bridge = []
  const order = Array(n).fill(-1)
  const lower = []
  let step = 0
  connections.forEach(([pre, next]) => {
    adjList.set(pre, adjList.get(pre)?.concat([next]) || [next])
  })

  const dfs = (cur, pre) => {
    visited.add(cur)
    order[cur] = step
    lower[cur] = step
    step++
    adjList.get(cur)?.forEach(next => {
      // 如果没被访问到
      if (!visited.has(next)) {
        visited.add(next)
        dfs(next, cur)
        // dfs回溯,根据ne的追溯值更新cur追溯值
        lower[cur] = Math.min(lower[cur], lower[next])
        if (order[cur] < lower[next]) bridge.push([cur, next])
        // 不问pre节点
      } else if (next !== pre) {
        // 如果next已被访问到(表示存在环，也就是ne可以直接到cur, 那么与dfn[ne]比较)
        lower[cur] = Math.min(lower[cur], lower[next])
      }
    })
  }
  // for (const v of adjList.keys()) {
  //   !visited.has(v) && dfs(v, v)
  // }
  for (let v = 0; v < n; v++) {
    !visited.has(v) && dfs(v, v)
  }
  // dfs(0, 0)
  return bridge
}

console.dir(
  criticalConnections(5, [
    [1, 0],
    [2, 0],
    [3, 0],
    [4, 1],
    [4, 2],
    [4, 0],
  ])
)
