/**
 * 
 * @param n   1 <= n <= 100
 * @param red_edges 
 * @param blue_edges 
 * 有向图
   返回长度为 n 的数组 answer，
   其中 answer[X] 是从节点 0 到节点 X 的红色边和蓝色边交替出现的最短路径的长度。
   如果不存在这样的路径，那么 answer[x] = -1
   @summary 
   最短路径 bfs
   我们并不知道节点0到达节点i的最短路径，是'红蓝红...'还是'蓝红蓝...'，
   所以我们需要都找出来，用nx2的数组dist保存，最后再选短的那个。
 */
function shortestAlternatingPaths(
  n: number,
  red_edges: number[][],
  blue_edges: number[][]
): number[] {
  // 建图
  const res = Array<number>(n).fill(Infinity)
  const redAdjList = Array.from<number, number[]>({ length: n }, () => [])
  const blueAdjList = Array.from<number, number[]>({ length: n }, () => [])
  for (const [v, w] of red_edges) redAdjList[v].push(w)
  for (const [v, w] of blue_edges) blueAdjList[v].push(w)
  const adjList: Record<number, number[][]> = {
    0: redAdjList,
    1: blueAdjList,
  }

  // bfs
  const visited = Array.from<unknown, [boolean, boolean]>({ length: n }, () => [false, false])
  const queue: [cur: number, color: number, steps: number][] = [
    [0, 0, 0],
    [0, 1, 0],
  ]

  while (queue.length) {
    const [cur, color, steps] = queue.shift()!
    res[cur] = Math.min(res[cur], steps)

    for (const next of adjList[color][cur]) {
      if (visited[next][color]) continue
      visited[next][color] = true
      queue.push([next, color ^ 1, steps + 1])
    }
  }

  return res.map(v => (v === Infinity ? -1 : v))
}

// console.log(shortestAlternatingPaths(3, [[0, 1]], [[2, 1]]))
console.log(
  shortestAlternatingPaths(
    5,
    [
      [0, 1],
      [1, 2],
      [2, 3],
      [3, 4],
    ],
    [
      [1, 2],
      [2, 3],
      [3, 1],
    ]
  )
)
// [0,1,-1]
