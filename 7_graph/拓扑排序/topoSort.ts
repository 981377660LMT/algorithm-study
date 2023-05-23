/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * 返回图的拓扑排序结果, 若不存在则返回空列表
 */
function topoSort(n: number, adjList: number[][], deg: number[], directed = true): number[] {
  const startDeg = directed ? 0 : 1
  const res: number[] = []
  const stack: number[] = []
  for (let i = 0; i < n; i++) {
    if (deg[i] === 0) stack.push(i)
  }
  const visited = new Uint8Array(n)

  while (stack.length) {
    const cur = stack.pop()!
    res.push(cur)
    visited[cur] = 1
    for (const next of adjList[cur]) {
      if (visited[next]) continue
      deg[next]--
      if (deg[next] === startDeg) stack.push(next)
    }
  }

  return res.length === n ? res : []
}

/**
 * 返回图的拓扑排序结果, 若不存在则返回空列表
 */
function topoSort2<T extends PropertyKey>(
  allVertex: Set<T>,
  adjMap: Map<T, T[]>,
  deg: Map<T, number>,
  directed = true
): T[] {
  const startDeg = directed ? 0 : 1
  const res: T[] = []
  const stack: T[] = []
  for (const v of allVertex) {
    if (deg.get(v) === 0) stack.push(v)
  }
  const visited = new Set<T>()

  while (stack.length) {
    const cur = stack.pop()!
    res.push(cur)
    visited.add(cur)
    for (const next of adjMap.get(cur)!) {
      if (visited.has(next)) continue
      deg.set(next, deg.get(next)! - 1)
      if (deg.get(next) === startDeg) stack.push(next)
    }
  }

  return res.length === allVertex.size ? res : []
}

export { topoSort, topoSort2 }
