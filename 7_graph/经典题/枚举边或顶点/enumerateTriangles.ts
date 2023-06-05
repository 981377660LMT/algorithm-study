/* eslint-disable no-inner-declarations */
// EnumeratingC3
// 遍历无向图中的所有三元环(a,b,c) a<b<c
//  O(E^1.5)
// https://kopricky.github.io/code/Graph/EnumeratingC3.html

/**
 * 遍历无向图中的所有三元环(a,b,c),`a<b<c`.
 * @param n 顶点数.
 * @param edges 边集.每条边为无向边'u<->v',范围为[0,n).边数<=2e5.
 * @param callback 回调函数.参数为三元环的三个顶点.
 * @complexity O(E^1.5).
 */
function enumerateTriangles(
  n: number,
  edges: [u: number, v: number][] | number[][],
  callback: (a: number, b: number, c: number) => void
): void {
  let edgeCount = 0
  const adjList: number[][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v]) => {
    // 无向边定向
    if (u < v) adjList[u].push(v)
    else adjList[v].push(u)
    edgeCount++
  })

  const threshold = Math.sqrt(edgeCount / 4) | 0
  const memo: number[][] = Array(n)
  for (let i = 0; i < n; i++) memo[i] = []
  const visited = new Uint8Array(n)

  // processHighDegree
  for (let i = 0; i < n; i++) {
    const adjI = adjList[i]
    if (adjI.length <= threshold) continue
    for (let j = 0; j < adjI.length; j++) visited[adjI[j]] = 1
    for (let j = 0; j < adjI.length; j++) {
      const u = adjI[j]
      const adjU = adjList[u]
      for (let k = 0; k < adjU.length; k++) {
        const v = adjU[k]
        if (visited[v]) callback(i, u, v)
      }
    }
    for (let j = 0; j < adjI.length; j++) visited[adjI[j]] = 0
  }

  // processLowDegree
  for (let i = 0; i < n; i++) {
    const adjI = adjList[i]
    if (adjI.length > threshold) continue
    for (let j = 0; j < adjI.length; j++) {
      const u = adjI[j]
      for (let k = 0; k < adjI.length; k++) {
        const v = adjI[k]
        if (v > u) memo[u].push(i * n + v)
      }
    }
  }
  for (let i = 0; i < n; i++) {
    const adjI = adjList[i]
    for (let j = 0; j < adjI.length; j++) visited[adjI[j]] = 1
    for (let j = 0; j < memo[i].length; j++) {
      const hash = memo[i][j]
      const a = (hash / n) | 0
      const b = hash % n
      if (visited[b]) callback(a, i, b)
    }
    for (let j = 0; j < adjI.length; j++) visited[adjI[j]] = 0
  }
}

export { enumerateTriangles }

if (require.main === module) {
  // 1761. 一个图中连通三元组的最小度数
  // https://leetcode.cn/problems/minimum-degree-of-a-connected-trio-in-a-graph/submissions/
  function minTrioDegree(n: number, edges: number[][]): number {
    const deg = new Uint16Array(n)
    for (let i = 0; i < edges.length; i++) {
      edges[i][0]--
      edges[i][1]--
      deg[edges[i][0]]++
      deg[edges[i][1]]++
    }

    let res = 2e15
    enumerateTriangles(n, edges, (a, b, c) => {
      res = Math.min(res, deg[a] + deg[b] + deg[c] - 6)
    })
    return res === 2e15 ? -1 : res
  }

  // n = 6, edges = [[1,2],[1,3],[3,2],[4,1],[5,2],[3,6]]
  console.log(
    minTrioDegree(6, [
      [1, 2],
      [1, 3],
      [3, 2],
      [4, 1],
      [5, 2],
      [3, 6]
    ])
  )
}
