/**
 * @param {number[][]} times
 * @param {number} n  有 n 个网络节点，标记为 1 到 n。
 * @param {number} k
 * @return {number}
 */
function networkDelayTime(times: number[][], n: number, k: number): number {
  const dist = Array<number>(n + 1).fill(Infinity)
  dist[k] = 0

  for (let i = 0; i < n; i++) {
    for (const [u, v, w] of times) {
      if (dist[u] + w < dist[v]) {
        dist[v] = dist[u] + w
      }
    }
  }

  let res = -1
  for (let i = 1; i <= n; i++) {
    res = Math.max(res, dist[i])
  }

  return res === Infinity ? -1 : res
}

export default 1

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
