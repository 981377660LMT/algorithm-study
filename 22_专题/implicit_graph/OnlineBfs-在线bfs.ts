/* eslint-disable no-constant-condition */

// 在线bfs求最短路.
// !这种边数很多的图论问题，通用的思路是，如果求最最短路，可以在线bfs，如果求最长路，可以线段树优化dp
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
): [dist: number[], pre: Int32Array] {
  const dist = Array(n).fill(INF)
  dist[start] = 0
  setUsed(start)
  const pre = new Int32Array(n).fill(-1)
  const queue = new Uint32Array(n)
  let left = 0
  let right = 0
  queue[right++] = start
  while (left < right) {
    const cur = queue[left++]
    while (true) {
      const next = findUnused(cur)
      if (next == null) break
      dist[next] = dist[cur] + 1
      pre[next] = cur
      setUsed(next)
      queue[right++] = next
    }
  }

  return [dist, pre]
}

export { onlineBfs }
