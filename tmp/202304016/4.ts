import { Tree } from '../../6_tree/重链剖分/Tree'

function isBipartite(n: number, adjList: number[][]): [colors: Int8Array, ok: boolean] {
  const colors = new Int8Array(n).fill(-1)
  for (let i = 0; i < n; i++) {
    if (colors[i] === -1 && !dfs(i, 0)) {
      return [colors, false]
    }
  }
  return [colors, true]

  function dfs(cur: number, color: number): boolean {
    colors[cur] = color
    for (let i = 0; i < adjList[cur].length; i++) {
      const next = adjList[cur][i]
      if (colors[next] === -1) {
        if (!dfs(next, color ^ 1)) {
          return false
        }
      } else if (colors[next] === color) {
        return false
      }
    }
    return true
  }
}

function minimumTotalPrice(
  n: number,
  edges: number[][],
  price: number[],
  trips: number[][]
): number {
  const adjList: number[][] = Array.from({ length: n }, () => [])
  for (const [u, v] of edges) {
    adjList[u].push(v)
    adjList[v].push(u)
  }
  const colors = isBipartite(n, adjList)[0]
  const price1 = price.slice()
  const price2 = price.slice()
  for (let i = 0; i < n; i++) {
    if (colors[i] === 0) {
      price1[i] >>= 1
    } else {
      price2[i] >>= 1
    }
  }

  const tree = new Tree(n)
  for (const [u, v] of edges) {
    tree.addEdge(u, v, 1)
  }
  tree.build(0)

  let res1 = 0
  let res2 = 0
  for (const [start, end] of trips) {
    const path = tree.getPath(start, end)
    for (const i of path) {
      res1 += price1[i]
      res2 += price2[i]
    }
  }
  return Math.min(res1, res2)
}

export {}
