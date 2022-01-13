function topoSort(adjList: number[], indegrees: number[]): number[] {
  const n = adjList.length
  const res: number[] = []
  let queue: number[] = []
  indegrees.forEach((degree, id) => degree === 0 && queue.push(id))

  while (queue.length > 0) {
    const len = queue.length
    const nextQueue: number[] = []

    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      const next = adjList[cur]
      indegrees[next]--
      if (indegrees[next] === 0) nextQueue.push(next)
    }

    queue = nextQueue
  }

  return res.length === n ? res : []
}

/**
 * @description 计算每个点在拓扑排序中的最大深度
 *
 **/
function topoSortDepth(adjList: number[], indegrees: number[]): number[] {
  const n = adjList.length
  const topoLevels = Array<number>(n).fill(0)
  let level = 0
  let queue: number[] = []
  indegrees.forEach((degree, id) => degree === 0 && queue.push(id))

  while (queue.length > 0) {
    const len = queue.length
    const nextQueue: number[] = []
    level++

    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      const next = adjList[cur]
      indegrees[next]--
      if (indegrees[next] === 0) nextQueue.push(next)
      topoLevels[next] = level
    }

    queue = nextQueue
  }

  return topoLevels
}

export { topoSort, topoSortDepth }
// 数组用s结尾比较好
