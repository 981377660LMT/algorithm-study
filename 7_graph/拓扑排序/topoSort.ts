function topoSort(adjMap: Map<string, Set<string>>, charSet: Set<string>) {
  const indegrees = new Map<string, number>()
  for (const charSet of adjMap.values()) {
    for (const char of charSet) {
      indegrees.set(char, (indegrees.get(char) || 0) + 1)
    }
  }

  // 入度为0的点
  const queue = [...charSet].filter(char => !indegrees.has(char))
  const res: string[] = []
  while (queue.length) {
    console.log(queue, indegrees, adjMap)
    const cur = queue.shift()!
    res.push(cur)
    for (const next of adjMap.get(cur) || []) {
      indegrees.set(next, indegrees.get(next)! - 1)
      if (indegrees.get(next) === 0) {
        queue.push(next)
      }
    }
  }

  return res.length === charSet.size ? res : []
}

export { topoSort }
