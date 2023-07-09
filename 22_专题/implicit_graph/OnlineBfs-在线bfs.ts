/* eslint-disable no-constant-condition */

// 注意还有一种求拓扑图最长路的问题
// !这种问题不能用在线bfs(无法保证边数为O(n))
// !需要利用线段树建图/线段树维护值域最大值来解决
// https://leetcode.cn/problems/maximum-number-of-jumps-to-reach-the-last-index/

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
  setUsed(start)
  let queue1 = new Uint32Array(n)
  let queue2 = new Uint32Array(n)

  queue1[0] = start
  let curQueue = queue1
  let nextQueue = queue2
  let curPtr = 1
  let nextPtr = 0

  while (curPtr) {
    for (let i = 0; i < curPtr; i++) {
      const cur = curQueue[i]
      while (true) {
        const next = findUnused(cur)
        if (next == null) {
          break
        }
        dist[next] = dist[cur] + 1 // weight
        nextQueue[nextPtr++] = next
        setUsed(next)
      }
    }
    const tmp = curQueue
    curQueue = nextQueue
    nextQueue = tmp
    curPtr = nextPtr
    nextPtr = 0
  }

  return dist
}

export { onlineBfs }
