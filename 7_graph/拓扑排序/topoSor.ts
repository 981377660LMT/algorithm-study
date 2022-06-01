/**
 *
 * @param n 顶点数
 * @param adjMap 有向图
 * @param deg 每个点的入度
 * @returns 拓扑排序结果
 * @description 注意图里存在重边时不要重复添加
 */
function topoSort(n: number, adjMap: Map<number, Set<number>>, deg: number[]): number[] {
  const res: number[] = []
  let queue: number[] = []
  deg.forEach((degree, id) => degree === 0 && queue.push(id))

  while (queue.length > 0) {
    const len = queue.length
    const nQueue: number[] = []

    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      res.push(cur)
      for (const next of adjMap.get(cur) ?? []) {
        deg[next]--
        if (deg[next] === 0) nQueue.push(next)
      }
    }

    queue = nQueue
  }

  return res.length === n ? res : []
}

/**
 * @param adjMap 有向图
 * @param allVertex 有向图的所有顶点
 * @returns 拓扑排序结果
 */
function topoSort2<T extends PropertyKey>(adjMap: Map<T, Set<T>>, allVertex: Set<T>): T[] {
  const deg = new Map<T, number>()
  for (const next of adjMap.values()) {
    for (const v of next) {
      deg.set(v, (deg.get(v) ?? 0) + 1)
    }
  }

  for (const v of allVertex) if (!deg.has(v)) deg.set(v, 0)

  // 入度为0的点
  let queue = [...allVertex].filter(v => (deg.get(v) ?? 0) === 0)
  const res: T[] = []

  while (queue.length) {
    const len = queue.length
    const nQueue: T[] = []

    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      res.push(cur)
      for (const next of adjMap.get(cur) ?? []) {
        deg.set(next, deg.get(next)! - 1)
        if (deg.get(next) === 0) {
          nQueue.push(next)
        }
      }
    }

    queue = nQueue
  }

  return res.length === allVertex.size ? res : [] // 否则有环，返回[]
}

/**
 * @description 计算每个点在拓扑排序中的最大深度
 *
 **/
function topoSortDepth(n: number, adjMap: Map<number, Set<number>>, deg: number[]): number[] {
  const topoLevels = Array<number>(n).fill(0)
  let level = 0
  let queue: number[] = []
  deg.forEach((degree, id) => degree === 0 && queue.push(id))

  while (queue.length > 0) {
    const len = queue.length
    const nQueue: number[] = []
    level++

    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      for (const next of adjMap.get(cur) ?? []) {
        deg[next]--
        if (deg[next] === 0) nQueue.push(next)
        topoLevels[next] = level
      }
    }

    queue = nQueue
  }

  return topoLevels
}

export { topoSort, topoSort2 }
// 数组用s结尾比较好
