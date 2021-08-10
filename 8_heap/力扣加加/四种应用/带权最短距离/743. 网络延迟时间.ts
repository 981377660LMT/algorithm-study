/**
 * @param {number[][]} times
 * @param {number} n
 * @param {number} k
 * @return {number}
 * @description
 * 有 n 个网络节点，标记为 1 到 n。
 * 从某个节点 K 发出一个信号。需要多久才能使所有节点都收到信号？如果不能使所有节点收到信号，返回 -1 。
 */
const networkDelayTime = function (times: number[][], n: number, k: number): number {
  let max = -1
  const adjList = Array.from<number, Map<number, number>>({ length: n }, () => new Map())
  const visited = Array<number>(n).fill(-1)
  times.forEach(item => adjList[item[0] - 1].set(item[1] - 1, item[2]))
  visited[k - 1] = 1

  const dfs = (cur: number, sum: number) => {
    for (const [next, weight] of adjList[cur]) {
      if (visited[next] === -1) {
        visited[next] = 1
        max = Math.max(max, weight + sum)
        dfs(next, weight + sum)
      }
    }
  }
  dfs(k - 1, 0)
  console.log(visited, adjList)
  return visited.includes(-1) ? -1 : max
}

console.log(
  networkDelayTime(
    [
      [2, 1, 1],
      [2, 3, 1],
      [3, 4, 1],
    ],
    4,
    2
  )
)
// 输出：2
