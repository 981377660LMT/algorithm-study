const hasCycle = (n: number, prerequisites: [number, number][]) => {
  if (prerequisites.length === 0) return false

  const inDegrees = Array<number>(n).fill(0)
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  for (const [cur, pre] of prerequisites) {
    inDegrees[cur]++
    adjList[pre].push(cur)
  }

  const queue: number[] = []
  inDegrees.forEach((v, i) => !v && queue.push(i))

  let count = 0
  while (queue.length) {
    const cur = queue.shift()!
    count++
    for (const next of adjList[cur]) {
      inDegrees[next]--
      if (inDegrees[next] === 0) queue.push(next)
    }
  }

  return count !== n
}

export {}
