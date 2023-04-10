import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'
import { Finder } from '../../../22_专题/implicit_graph/RangeFinder/Finder-fastset'

const INF = 2e15

function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
  const finder = [new Finder(n), new Finder(n)]
  for (let i = 0; i < n; i++) {
    finder[(i & 1) ^ 1].erase(i)
  }
  banned.forEach(i => {
    finder[i & 1].erase(i)
  })

  // 反转长度为k的子数组,从i可以到达的左右边界闭区间
  const getRange = (i: number): [number, number] => [
    Math.max(i - k + 1, k - i - 1),
    Math.min(i + k - 1, 2 * n - k - i - 1)
  ]

  // 将u位置标记为已经访问过.
  const setUsed = (u: number) => {
    finder[u & 1].erase(u)
  }

  // 找到一个未使用的位置.如果不存在,返回null.
  const findUnused = (u: number): number | null => {
    const [left, right] = getRange(u)
    const next = finder[(u + k + 1) & 1].next(left)
    if (next != null && left <= next && next <= right) {
      return next
    }
    return null
  }

  const dist = onlineBfs(n, p, setUsed, findUnused)
  return dist.map(d => (d === INF ? -1 : d))
}
