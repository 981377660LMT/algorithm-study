/* eslint-disable no-constant-condition */
/* eslint-disable no-loop-func */

const INF = 2e15

/**
 * BellmanFord算法`O(VE)`求解带负权边的单源最短路,并求出每个点的前驱.
 * @returns [dist, pre] (起点到各点的最短距离,每个点的前驱).
 * 距离为`-INF`表示经过负环到达.
 * 距离为`INF`表示不可达.
 */
function bellmanFord(
  n: number,
  adjList: [to: number, cost: number][][] | number[][][],
  start: number
): [dist: number[], pre: Int32Array] {
  const dist = Array(n)
  const pre = new Int32Array(n)
  for (let i = 0; i < n; ++i) {
    dist[i] = INF
    pre[i] = -1
  }
  dist[start] = 0
  let loop = 0

  while (true) {
    loop++
    let updated = false
    for (let from = 0; from < n; from++) {
      if (dist[from] === INF) continue
      adjList[from].forEach(([to, cost]) => {
        let cand = dist[from] + cost
        if (cand < -INF) cand = -INF
        if (cand < dist[to]) {
          updated = true
          pre[to] = from
          if (loop >= n) cand = -INF
          dist[to] = cand
        }
      })
    }

    if (!updated) break
  }

  return [dist, pre]
}

export { bellmanFord }

if (require.main === module) {
  const n = 4
  const edges = [
    [0, 1, 2],
    [0, 2, 3],
    [1, 2, -5],
    [1, 3, 1],
    [2, 3, 2]
  ]
  const adjList = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v, w]) => {
    adjList[u].push([v, w])
  })
  const start = 1
  const [dist, pre] = bellmanFord(n, adjList, start)
  console.log(dist)
}
