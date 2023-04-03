/* eslint-disable no-constant-condition */

const INF = 2e15

/**
 * 在线bfs.
 * 不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
 * @param setUsed 将 u 标记为已访问。
 * @param findUnused 找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `null`。
 * https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/
 */
function onlineBfs(
  n: number,
  start: number,
  setUsed: (cur: number) => void,
  findUnused: (cur: number) => number | null
): number[] {
  const dist = new Array(n).fill(INF)
  dist[start] = 0
  let queue = [start]
  setUsed(start)

  while (queue.length) {
    const nextQueue: number[] = []
    for (let i = 0; i < queue.length; i++) {
      const cur = queue[i]
      while (true) {
        const next = findUnused(cur)
        if (next == null) {
          break
        }
        dist[next] = dist[cur] + 1 // weight
        nextQueue.push(next)
        setUsed(next)
      }
    }

    queue = nextQueue
  }

  return dist
}

export { onlineBfs }
