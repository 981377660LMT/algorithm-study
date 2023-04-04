// 寻找前驱后继/区间删除
// 非常快

import { FastSet } from '../../../24_高级数据结构/珂朵莉树/FastSet'
import { onlineBfs } from '../OnlineBfs-在线bfs'

/**
 * 使用FastSet寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 */
const Finder = FastSet

// https://leetcode.cn/problems/minimum-reverse-operations/
function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
  const INF = 2e15
  const finder = [new Finder(n), new Finder(n)]
  for (let i = 0; i < n; i++) {
    finder[i & 1].insert(i)
  }
  banned.forEach(i => {
    finder[i & 1].erase(i)
  })

  const getRange = (i: number): [number, number] => [
    Math.max(i - k + 1, k - i - 1),
    Math.min(i + k - 1, 2 * n - k - i - 1)
  ]

  const setUsed = (u: number) => {
    finder[u & 1].erase(u)
  }

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

if (require.main === module) {
  //   4
  // 0
  // [1,2]
  // 4
  console.log(minReverseOperations(4, 0, [1, 2], 4))
}

export { Finder }
