/** 图的深度优先遍历，递归转迭代. */
function dfsNoRecursion(
  graph: number[][],
  start: number,
  consumer: (node: number) => void | boolean
): void {
  const n = graph.length
  const curEdge = new Int32Array(n)
  const stack = new Int32Array(n)
  stack[0] = start
  let top = 0
  while (top >= 0) {
    const u = stack[top]
    if (curEdge[u] === 0) {
      if (consumer(u)) return
    }
    if (curEdge[u] < graph[u].length) {
      const v = graph[u][curEdge[u]++]
      if (curEdge[v] === 0) {
        stack[++top] = v
      }
    } else {
      top--
    }
  }
}

export {}

if (require.main === module) {
  const n = 4
  const edges = [
    [0, 1],
    [1, 2],
    [2, 3],
    [3, 0]
  ]
  const adjList: number[][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v]) => {
    adjList[u].push(v)
    adjList[v].push(u)
  })

  dfsNoRecursion(adjList, 1, node => {
    console.log(node)
  })
}
