/* eslint-disable @typescript-eslint/no-unused-vars */
// 寻找前驱后继/区间删除

import { SortedList } from '../../离线查询/根号分治/SortedList/SortedList'
import { onlineBfs } from '../OnlineBfs-在线bfs'

/**
 * 使用有序列表寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 */
const Finder = SortedList

// https://leetcode.cn/problems/minimum-reverse-operations/
function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
  const finder = [new Finder(), new Finder()]
  for (let i = 0; i < n; i++) {
    finder[i & 1].add(i)
  }
  banned.forEach(i => {
    finder[i & 1].discard(i)
  })

  const getRange = (i: number): [number, number] => [
    Math.max(i - k + 1, k - i - 1),
    Math.min(i + k - 1, 2 * n - k - i - 1)
  ]

  const setUsed = (u: number) => {
    finder[u & 1].discard(u)
  }

  const findUnused = (u: number): number | null => {
    const [left, right] = getRange(u)
    const next = finder[(u + k + 1) & 1].ceiling(left)
    if (next != null && left <= next && next <= right) {
      return next
    }
    return null
  }

  const dist = onlineBfs(n, p, setUsed, findUnused)
  return dist.map(d => (d === INF ? -1 : d))
}

export { Finder }
