function maxTargetNodes(edges1: number[][], edges2: number[][], k: number): number[] {
  const n = edges1.length + 1
  const m = edges2.length + 1
  const tree1: number[][] = Array.from({ length: n }, () => [])
  for (const [u, v] of edges1) {
    tree1[u].push(v)
    tree1[v].push(u)
  }
  const tree2: number[][] = Array.from({ length: m }, () => [])
  for (const [u, v] of edges2) {
    tree2[u].push(v)
    tree2[v].push(u)
  }

  let maxNode = 0
  if (k > 0) {
    for (let i = 0; i < m; i++) {
      const count = bfs(tree2, i, k - 1)
      if (count > maxNode) {
        maxNode = count
      }
    }
  }

  const res: number[] = []
  for (let i = 0; i < n; i++) {
    const count = bfs(tree1, i, k)
    res.push(count + maxNode)
  }
  return res
}

function bfs(tree: number[][], start: number, maxDepth: number): number {
  if (maxDepth < 0) return 0
  const visited = new Set<number>()
  let queue: [number, number][] = []
  queue.push([start, 0])
  visited.add(start)
  let count = 1
  while (queue.length > 0) {
    const nextQueue: [number, number][] = []
    for (const [node, depth] of queue) {
      if (depth >= maxDepth) continue
      for (const neighbor of tree[node]) {
        if (!visited.has(neighbor)) {
          visited.add(neighbor)
          nextQueue.push([neighbor, depth + 1])
          count++
        }
      }
    }
    queue = nextQueue
  }
  return count
}
