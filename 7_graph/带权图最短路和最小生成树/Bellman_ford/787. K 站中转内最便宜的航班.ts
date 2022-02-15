/**
 * @param {number} n
 * @param {number[][]} flights
 * @param {number} src
 * @param {number} dst
 * @param {number} k  最多经过 k 站中转的路线
 * @return {number}
 * 找到出一条最多经过 k 站中转的路线，使得从 src 到 dst 的 价格最便宜 ，
 * 并返回该价格。 如果不存在这样的路线，则输出 -1。
 * @summary
 * 带限制的最短路径
 */
const findCheapestPrice = function (
  n: number,
  flights: number[][],
  src: number,
  dst: number,
  k: number
): number {
  const dist = Array<number>(n).fill(Infinity)
  dist[src] = 0

  // 最多k个中转点，即更新k+1次
  for (let i = 0; i < k + 1; i++) {
    const clone = dist.slice()
    for (const [u, v, w] of flights) {
      if (clone[u] + w < dist[v]) {
        dist[v] = clone[u] + w
      }
    }
  }

  return dist[dst] === Infinity ? -1 : dist[dst]
}

console.log(
  findCheapestPrice(
    3,
    [
      [0, 1, 100],
      [1, 2, 100],
      [0, 2, 500],
    ],
    0,
    2,
    1
  )
)
