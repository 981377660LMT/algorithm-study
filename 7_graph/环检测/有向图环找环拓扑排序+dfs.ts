const hasCycle = (adjList: number[][]) => {
  const n = adjList.length
  const indegrees = Array<number>(n).fill(0)
  for (const [_, next] of adjList) {
    indegrees[next]++
  }

  const queue: number[] = []
  indegrees.forEach((degree, i) => degree === 0 && queue.push(i))

  while (queue.length > 0) {
    const cur = queue.shift()!
    for (const next of adjList[cur]) {
      indegrees[next]--
      if (indegrees[next] === 0) queue.push(next)
    }
  }

  return indegrees.some(degree => degree !== 0)
}

export {}
